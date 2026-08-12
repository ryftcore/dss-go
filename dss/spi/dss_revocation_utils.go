// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/DSSRevocationUtils.java (DSS 6.5.RC1).
//
// INTEGRATION NOTE: upstream this file is a thin adapter over BouncyCastle's OCSP stack
// (org.bouncycastle.asn1.ocsp.* and org.bouncycastle.cert.ocsp.*). Go has no equivalent -
// golang.org/x/crypto/ocsp only models a single-response, by-key-hash reply and drops the
// raw bytes - so the RFC 6960 structures and the BasicOCSPResp/OCSPResp/SingleResp/
// CertificateID/RespID wrappers they are used through are defined HERE, next to the code
// that builds and converts them. Everything that consumes OCSP in phase 2a (OCSPToken,
// OCSPRef, OfflineOCSPSource, OCSPCertificateSource, ...) must consume these types rather
// than declare its own.
//
// The raw encoding of every parsed structure is retained (Encoded()/ASN1Structure()), since
// the revocation identifiers digest the response bytes.
//
// NOT PORTED: getDigestCalculator(DigestAlgorithm) returns a BouncyCastle
// org.bouncycastle.operator.DigestCalculator; it exists only to feed
// CertificateID's constructor, whose work DSSRevocationUtilsOCSPCertificateID does directly.
package spi

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"math/big"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// OCSP response status values of RFC 6960.
const (
	// OCSPResponseStatusSuccessful means the response has valid confirmations.
	OCSPResponseStatusSuccessful = 0
	// OCSPResponseStatusMalformedRequest means the request was not well formed.
	OCSPResponseStatusMalformedRequest = 1
	// OCSPResponseStatusInternalError means an internal error occurred in the responder.
	OCSPResponseStatusInternalError = 2
	// OCSPResponseStatusTryLater means the service is temporarily unavailable.
	OCSPResponseStatusTryLater = 3
	// OCSPResponseStatusSigRequired means the request must be signed.
	OCSPResponseStatusSigRequired = 5
	// OCSPResponseStatusUnauthorized means the requester is not authorised.
	OCSPResponseStatusUnauthorized = 6
)

// Certificate status alternatives of the RFC 6960 CertStatus CHOICE.
const (
	// OCSPCertStatusGood is CertStatus.good, context tag [0].
	OCSPCertStatusGood = 0
	// OCSPCertStatusRevoked is CertStatus.revoked, context tag [1].
	OCSPCertStatusRevoked = 1
	// OCSPCertStatusUnknown is CertStatus.unknown, context tag [2].
	OCSPCertStatusUnknown = 2
)

var (
	// OCSPObjectIdentifierIDPkixOcspBasic is id-pkix-ocsp-basic, the response type of a
	// BasicOCSPResponse. Corresponds to OCSPObjectIdentifiers.id_pkix_ocsp_basic.
	OCSPObjectIdentifierIDPkixOcspBasic = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 1, 1}
	// OCSPObjectIdentifierIDPkixOcspNonce is id-pkix-ocsp-nonce.
	OCSPObjectIdentifierIDPkixOcspNonce = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 1, 2}
	// OCSPObjectIdentifierIDPkixOcspArchiveCutoff is id-pkix-ocsp-archive-cutoff.
	OCSPObjectIdentifierIDPkixOcspArchiveCutoff = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 1, 6}

	// dssRevocationUtilsOIDSHA1 is the OIW id-sha1 OBJECT IDENTIFIER, the algorithm implied
	// by the sha1Hash alternative of an ESF OtherHash.
	dssRevocationUtilsOIDSHA1 = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}
)

// -----------------------------------------------------------------------------
// RFC 6960 structures (replacing org.bouncycastle.asn1.ocsp.*).
// -----------------------------------------------------------------------------

// Extension is the X.509 structure
//
//	Extension ::= SEQUENCE {
//	    extnID    OBJECT IDENTIFIER,
//	    critical  BOOLEAN DEFAULT FALSE,
//	    extnValue OCTET STRING }
//
// replacing org.bouncycastle.asn1.x509.Extension.
type Extension struct {
	// ID is the extension OID.
	ID asn1.ObjectIdentifier
	// Critical reports whether the extension is critical.
	Critical bool
	// Value holds the content octets of extnValue, i.e. the DER of the extension's value.
	Value []byte
}

// ParsedValue returns the DER encoding of the value carried by extnValue.
// Port of Extension#getParsedValue().
func (e *Extension) ParsedValue() []byte {
	return e.Value
}

// ResponderID is the RFC 6960 CHOICE
//
//	ResponderID ::= CHOICE {
//	    byName [1] Name,
//	    byKey  [2] KeyHash }
//
// replacing org.bouncycastle.asn1.ocsp.ResponderID. Exactly one field is set.
type ResponderID struct {
	// Name holds the DER encoding of the X.501 Name when the responder identified itself
	// by name, nil otherwise.
	Name []byte
	// KeyHash holds the SHA-1 hash of the responder's public key when the responder
	// identified itself by key, nil otherwise.
	KeyHash []byte
}

// RespID wraps a ResponderID, mirroring org.bouncycastle.cert.ocsp.RespID.
type RespID struct {
	responderID *ResponderID
}

// NewRespID wraps the given ResponderID.
func NewRespID(responderID *ResponderID) *RespID {
	return &RespID{responderID: responderID}
}

// ToASN1Primitive returns the wrapped ResponderID. Port of RespID#toASN1Primitive().
func (r *RespID) ToASN1Primitive() *ResponderID {
	return r.responderID
}

// CertID is the RFC 6960 structure
//
//	CertID ::= SEQUENCE {
//	    hashAlgorithm  AlgorithmIdentifier,
//	    issuerNameHash OCTET STRING,
//	    issuerKeyHash  OCTET STRING,
//	    serialNumber   CertificateSerialNumber }
type CertID struct {
	// HashAlgorithm is the algorithm the two hashes were computed with.
	HashAlgorithm *AlgorithmIdentifier
	// IssuerNameHash is the hash of the issuer's distinguished name.
	IssuerNameHash []byte
	// IssuerKeyHash is the hash of the issuer's public key.
	IssuerKeyHash []byte
	// SerialNumber is the serial number of the certificate the status is asked about.
	SerialNumber *big.Int
}

// DER returns the DER encoding of the CertID.
func (c *CertID) DER() []byte {
	body := c.HashAlgorithm.DER()
	body = append(body, dssASN1UtilsWriteTLV(dssASN1UtilsTagOctetString, c.IssuerNameHash)...)
	body = append(body, dssASN1UtilsWriteTLV(dssASN1UtilsTagOctetString, c.IssuerKeyHash)...)
	body = append(body, dssASN1UtilsEncodeInteger(c.SerialNumber)...)
	return dssASN1UtilsWriteTLV(dssASN1UtilsTagSequence|dssASN1UtilsConstructed, body)
}

// CertificateID wraps a CertID, mirroring org.bouncycastle.cert.ocsp.CertificateID.
type CertificateID struct {
	id *CertID
}

// NewCertificateID wraps the given CertID.
func NewCertificateID(id *CertID) *CertificateID {
	return &CertificateID{id: id}
}

// HashAlgOID returns the OID of the hash algorithm. Port of CertificateID#getHashAlgOID().
func (c *CertificateID) HashAlgOID() asn1.ObjectIdentifier { return c.id.HashAlgorithm.Algorithm }

// IssuerNameHash returns the hash of the issuer's name.
func (c *CertificateID) IssuerNameHash() []byte { return c.id.IssuerNameHash }

// IssuerKeyHash returns the hash of the issuer's public key.
func (c *CertificateID) IssuerKeyHash() []byte { return c.id.IssuerKeyHash }

// SerialNumber returns the serial number of the certificate.
func (c *CertificateID) SerialNumber() *big.Int { return c.id.SerialNumber }

// ToASN1Primitive returns the wrapped CertID.
func (c *CertificateID) ToASN1Primitive() *CertID { return c.id }

// OCSPCertStatus is the RFC 6960 CertStatus CHOICE, replacing
// org.bouncycastle.cert.ocsp.CertificateStatus, RevokedStatus and UnknownStatus.
type OCSPCertStatus struct {
	// Tag is the CHOICE alternative: OCSPCertStatusGood, OCSPCertStatusRevoked or
	// OCSPCertStatusUnknown.
	Tag int
	// RevocationTime is set for a revoked certificate.
	RevocationTime time.Time
	// RevocationReason is the CRLReason of a revoked certificate, nil when it is absent.
	RevocationReason *int
}

// IsGood reports whether the certificate is not revoked.
func (s *OCSPCertStatus) IsGood() bool { return s.Tag == OCSPCertStatusGood }

// IsRevoked reports whether the certificate is revoked.
func (s *OCSPCertStatus) IsRevoked() bool { return s.Tag == OCSPCertStatusRevoked }

// IsUnknown reports whether the responder does not know about the certificate.
func (s *OCSPCertStatus) IsUnknown() bool { return s.Tag == OCSPCertStatusUnknown }

// HasRevocationReason reports whether a revocation reason is present.
// Port of RevokedStatus#hasRevocationReason().
func (s *OCSPCertStatus) HasRevocationReason() bool { return s.RevocationReason != nil }

// SingleResponse is the RFC 6960 structure
//
//	SingleResponse ::= SEQUENCE {
//	    certID           CertID,
//	    certStatus       CertStatus,
//	    thisUpdate       GeneralizedTime,
//	    nextUpdate  [0] EXPLICIT GeneralizedTime OPTIONAL,
//	    singleExtensions [1] EXPLICIT Extensions OPTIONAL }
type SingleResponse struct {
	// CertID identifies the certificate the status refers to.
	CertID *CertID
	// CertStatus is the reported status.
	CertStatus *OCSPCertStatus
	// ThisUpdate is the time at which the status is known to be correct.
	ThisUpdate time.Time
	// NextUpdate is the time at or before which newer information will be available,
	// nil when absent.
	NextUpdate *time.Time
	// SingleExtensions holds the response's single extensions.
	SingleExtensions []Extension
}

// SingleResp wraps a SingleResponse, mirroring org.bouncycastle.cert.ocsp.SingleResp.
type SingleResp struct {
	response *SingleResponse
}

// NewSingleResp wraps the given SingleResponse.
func NewSingleResp(response *SingleResponse) *SingleResp {
	return &SingleResp{response: response}
}

// CertID returns the identifier of the certificate. Port of SingleResp#getCertID().
func (s *SingleResp) CertID() *CertificateID { return NewCertificateID(s.response.CertID) }

// CertStatus returns the reported status. Port of SingleResp#getCertStatus().
func (s *SingleResp) CertStatus() *OCSPCertStatus { return s.response.CertStatus }

// ThisUpdate returns the thisUpdate time. Port of SingleResp#getThisUpdate().
func (s *SingleResp) ThisUpdate() time.Time { return s.response.ThisUpdate }

// NextUpdate returns the nextUpdate time, nil when absent.
// Port of SingleResp#getNextUpdate().
func (s *SingleResp) NextUpdate() *time.Time { return s.response.NextUpdate }

// Extension returns the single extension with the given OID, nil when absent.
// Port of SingleResp#getExtension(ASN1ObjectIdentifier).
func (s *SingleResp) Extension(oid asn1.ObjectIdentifier) *Extension {
	for index := range s.response.SingleExtensions {
		if s.response.SingleExtensions[index].ID.Equal(oid) {
			return &s.response.SingleExtensions[index]
		}
	}
	return nil
}

// ToASN1Primitive returns the wrapped SingleResponse.
func (s *SingleResp) ToASN1Primitive() *SingleResponse { return s.response }

// ResponseData is the RFC 6960 structure
//
//	ResponseData ::= SEQUENCE {
//	    version            [0] EXPLICIT Version DEFAULT v1,
//	    responderID            ResponderID,
//	    producedAt             GeneralizedTime,
//	    responses              SEQUENCE OF SingleResponse,
//	    responseExtensions [1] EXPLICIT Extensions OPTIONAL }
type ResponseData struct {
	// Version is the response version; 0 stands for v1.
	Version int
	// ResponderID identifies the responder.
	ResponderID *ResponderID
	// ProducedAt is the time at which the responder signed the response.
	ProducedAt time.Time
	// Responses holds one SingleResponse per requested certificate.
	Responses []*SingleResponse
	// ResponseExtensions holds the response extensions.
	ResponseExtensions []Extension
	// encoded is the DER of the ResponseData, retained because it is what the signature
	// covers.
	encoded []byte
}

// BasicOCSPResponse is the RFC 6960 structure
//
//	BasicOCSPResponse ::= SEQUENCE {
//	    tbsResponseData      ResponseData,
//	    signatureAlgorithm   AlgorithmIdentifier,
//	    signature            BIT STRING,
//	    certs            [0] EXPLICIT SEQUENCE OF Certificate OPTIONAL }
//
// replacing org.bouncycastle.asn1.ocsp.BasicOCSPResponse.
type BasicOCSPResponse struct {
	// TBSResponseData is the signed response data.
	TBSResponseData *ResponseData
	// SignatureAlgorithm is the algorithm the response was signed with.
	SignatureAlgorithm *AlgorithmIdentifier
	// Signature is the signature value, i.e. the BIT STRING without its unused-bit count.
	Signature []byte
	// Certs holds the DER of the certificates the responder attached, in order.
	Certs [][]byte
	// encoded is the original encoding of the BasicOCSPResponse.
	encoded []byte
}

// ParseBasicOCSPResponse decodes a BasicOCSPResponse from its encoding, keeping the bytes.
func ParseBasicOCSPResponse(der []byte) (*BasicOCSPResponse, error) {
	element, rest, err := dssASN1UtilsParse(der)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the BasicOCSPResponse")
	}
	return dssRevocationUtilsBasicOCSPResponseFromElement(element)
}

// Encoded returns the original encoding of the BasicOCSPResponse.
func (b *BasicOCSPResponse) Encoded() []byte { return b.encoded }

// BasicOCSPResp wraps a BasicOCSPResponse, mirroring org.bouncycastle.cert.ocsp.BasicOCSPResp.
type BasicOCSPResp struct {
	response *BasicOCSPResponse
}

// NewBasicOCSPResp wraps the given BasicOCSPResponse.
// Port of the BasicOCSPResp(BasicOCSPResponse) constructor.
func NewBasicOCSPResp(response *BasicOCSPResponse) *BasicOCSPResp {
	return &BasicOCSPResp{response: response}
}

// ASN1Structure returns the wrapped BasicOCSPResponse.
// Port of BasicOCSPResp#toASN1Structure().
func (b *BasicOCSPResp) ASN1Structure() *BasicOCSPResponse { return b.response }

// Encoded returns the encoding of the wrapped BasicOCSPResponse.
// Port of BasicOCSPResp#getEncoded().
func (b *BasicOCSPResp) Encoded() []byte { return b.response.encoded }

// ProducedAt returns the time at which the response was produced.
// Port of BasicOCSPResp#getProducedAt().
func (b *BasicOCSPResp) ProducedAt() time.Time { return b.response.TBSResponseData.ProducedAt }

// Version returns the response version. Port of BasicOCSPResp#getVersion().
func (b *BasicOCSPResp) Version() int { return b.response.TBSResponseData.Version + 1 }

// ResponderID returns the identity of the responder. Port of BasicOCSPResp#getResponderId().
func (b *BasicOCSPResp) ResponderID() *RespID {
	return NewRespID(b.response.TBSResponseData.ResponderID)
}

// Responses returns every single response of the reply, in order.
// Port of BasicOCSPResp#getResponses().
func (b *BasicOCSPResp) Responses() []*SingleResp {
	responses := make([]*SingleResp, 0, len(b.response.TBSResponseData.Responses))
	for _, response := range b.response.TBSResponseData.Responses {
		responses = append(responses, NewSingleResp(response))
	}
	return responses
}

// SignatureAlgorithmID returns the algorithm the response was signed with.
// Port of BasicOCSPResp#getSignatureAlgorithmID().
func (b *BasicOCSPResp) SignatureAlgorithmID() *AlgorithmIdentifier {
	return b.response.SignatureAlgorithm
}

// Signature returns the signature value. Port of BasicOCSPResp#getSignature().
func (b *BasicOCSPResp) Signature() []byte { return b.response.Signature }

// TBSResponseData returns the DER of the signed response data, i.e. the bytes the signature
// covers. Port of BasicOCSPResp#getTBSResponseData().
func (b *BasicOCSPResp) TBSResponseData() []byte { return b.response.TBSResponseData.encoded }

// Extension returns the response extension with the given OID, nil when absent.
// Port of BasicOCSPResp#getExtension(ASN1ObjectIdentifier).
func (b *BasicOCSPResp) Extension(oid asn1.ObjectIdentifier) *Extension {
	extensions := b.response.TBSResponseData.ResponseExtensions
	for index := range extensions {
		if extensions[index].ID.Equal(oid) {
			return &extensions[index]
		}
	}
	return nil
}

// Certs returns the certificates attached to the response.
// Port of BasicOCSPResp#getCerts(), which yields X509CertificateHolders; the Go port
// yields parsed certificates and skips the ones crypto/x509 cannot decode.
func (b *BasicOCSPResp) Certs() []*x509.Certificate {
	certificates := make([]*x509.Certificate, 0, len(b.response.Certs))
	for _, der := range b.response.Certs {
		certificate, err := x509.ParseCertificate(der)
		if err != nil {
			continue
		}
		certificates = append(certificates, certificate)
	}
	return certificates
}

// OCSPResponse is the RFC 6960 structure
//
//	OCSPResponse ::= SEQUENCE {
//	    responseStatus  OCSPResponseStatus,
//	    responseBytes   [0] EXPLICIT ResponseBytes OPTIONAL }
//	ResponseBytes ::= SEQUENCE {
//	    responseType    OBJECT IDENTIFIER,
//	    response        OCTET STRING }
//
// replacing org.bouncycastle.asn1.ocsp.OCSPResponse.
type OCSPResponse struct {
	// ResponseStatus is one of the OCSPResponseStatus* constants.
	ResponseStatus int
	// ResponseType is the OID of the encapsulated response, nil when responseBytes is absent.
	ResponseType asn1.ObjectIdentifier
	// Response holds the content octets of the response OCTET STRING, nil when absent.
	Response []byte
	// encoded is the original encoding of the OCSPResponse.
	encoded []byte
}

// NewOCSPResponse builds a successful OCSPResponse encapsulating the given basic response
// binaries. Port of new OCSPResponse(OCSPResponseStatus, ResponseBytes).
func NewOCSPResponse(responseStatus int, responseType asn1.ObjectIdentifier, response []byte) *OCSPResponse {
	ocspResponse := &OCSPResponse{ResponseStatus: responseStatus, ResponseType: responseType, Response: response}
	ocspResponse.encoded = ocspResponse.DER()
	return ocspResponse
}

// ParseOCSPResponse decodes an OCSPResponse from its encoding, keeping the bytes.
func ParseOCSPResponse(der []byte) (*OCSPResponse, error) {
	element, rest, err := dssASN1UtilsParse(der)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the OCSPResponse")
	}
	return dssRevocationUtilsOCSPResponseFromElement(element)
}

// DER returns the DER encoding of the OCSPResponse.
func (o *OCSPResponse) DER() []byte {
	body := dssASN1UtilsWriteTLV(0x0A, []byte{byte(o.ResponseStatus)}) // ENUMERATED
	if o.ResponseType != nil {
		responseBytes := dssASN1UtilsEncodeOID(o.ResponseType)
		responseBytes = append(responseBytes, dssASN1UtilsWriteTLV(dssASN1UtilsTagOctetString, o.Response)...)
		responseBytes = dssASN1UtilsWriteTLV(dssASN1UtilsTagSequence|dssASN1UtilsConstructed, responseBytes)
		body = append(body, dssASN1UtilsWriteTLV(0xA0, responseBytes)...)
	}
	return dssASN1UtilsWriteTLV(dssASN1UtilsTagSequence|dssASN1UtilsConstructed, body)
}

// Encoded returns the original encoding of the OCSPResponse.
func (o *OCSPResponse) Encoded() []byte {
	if o.encoded == nil {
		o.encoded = o.DER()
	}
	return o.encoded
}

// OCSPResp wraps an OCSPResponse, mirroring org.bouncycastle.cert.ocsp.OCSPResp.
type OCSPResp struct {
	response *OCSPResponse
}

// NewOCSPResp wraps the given OCSPResponse. Port of the OCSPResp(OCSPResponse) constructor.
func NewOCSPResp(response *OCSPResponse) *OCSPResp {
	return &OCSPResp{response: response}
}

// NewOCSPRespFromBinaries decodes an OCSPResp from the DER of an OCSPResponse.
// Port of the OCSPResp(byte[]) constructor.
func NewOCSPRespFromBinaries(binaries []byte) (*OCSPResp, error) {
	response, err := ParseOCSPResponse(binaries)
	if err != nil {
		return nil, err
	}
	return &OCSPResp{response: response}, nil
}

// Status returns the response status. Port of OCSPResp#getStatus().
func (o *OCSPResp) Status() int { return o.response.ResponseStatus }

// Encoded returns the encoding of the wrapped OCSPResponse. Port of OCSPResp#getEncoded().
func (o *OCSPResp) Encoded() []byte { return o.response.Encoded() }

// ASN1Structure returns the wrapped OCSPResponse.
func (o *OCSPResp) ASN1Structure() *OCSPResponse { return o.response }

// ResponseObject returns the decoded response, which for id-pkix-ocsp-basic is a
// BasicOCSPResp. Port of OCSPResp#getResponseObject(); a response of another type yields
// (nil, nil), as Java yields an object of an unknown class.
func (o *OCSPResp) ResponseObject() (*BasicOCSPResp, error) {
	if o.response.ResponseType == nil {
		return nil, nil
	}
	if !o.response.ResponseType.Equal(OCSPObjectIdentifierIDPkixOcspBasic) {
		return nil, nil
	}
	basicOCSPResponse, err := ParseBasicOCSPResponse(o.response.Response)
	if err != nil {
		return nil, err
	}
	return NewBasicOCSPResp(basicOCSPResponse), nil
}

// OtherHash is the ESF CHOICE
//
//	OtherHash ::= CHOICE {
//	    sha1Hash  OtherHashValue,
//	    otherHash OtherHashAlgAndValue }
//	OtherHashAlgAndValue ::= SEQUENCE {
//	    hashAlgorithm AlgorithmIdentifier,
//	    hashValue     OtherHashValue }
//
// replacing org.bouncycastle.asn1.esf.OtherHash.
type OtherHash struct {
	// HashAlgorithm is the algorithm the value was computed with; the sha1Hash alternative
	// implies id-sha1 without parameters.
	HashAlgorithm *AlgorithmIdentifier
	// HashValue is the hash value.
	HashValue []byte
}

// ParseOtherHash decodes an OtherHash from its DER encoding.
// Port of OtherHash.getInstance(Object).
func ParseOtherHash(der []byte) (*OtherHash, error) {
	element, rest, err := dssASN1UtilsParse(der)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the OtherHash")
	}
	if element.isUniversal(dssASN1UtilsTagOctetString) {
		return &OtherHash{
			HashAlgorithm: NewAlgorithmIdentifier(dssRevocationUtilsOIDSHA1),
			HashValue:     element.octets(),
		}, nil
	}
	if !element.constructed || len(element.children) != 2 {
		return nil, errors.New("the OtherHash is neither an OCTET STRING nor an OtherHashAlgAndValue")
	}
	algorithm, err := dssASN1UtilsAlgorithmIdentifierFromElement(element.children[0])
	if err != nil {
		return nil, err
	}
	if !element.children[1].isUniversal(dssASN1UtilsTagOctetString) {
		return nil, errors.New("the OtherHashAlgAndValue hashValue is not an OCTET STRING")
	}
	return &OtherHash{HashAlgorithm: algorithm, HashValue: element.children[1].octets()}, nil
}

// -----------------------------------------------------------------------------
// The ported DSSRevocationUtils methods.
// -----------------------------------------------------------------------------

// DSSRevocationUtilsBasicOcspResp builds a BasicOCSPResp from the DER of a
// BasicOCSPResponse (RFC 2560), returning nil when the binaries cannot be decoded.
// Port of getBasicOcspResp(ASN1Sequence), whose argument is here the encoding of that
// sequence rather than a parsed object.
func DSSRevocationUtilsBasicOcspResp(asn1Sequence []byte) *BasicOCSPResp {
	basicOcspResponse, err := ParseBasicOCSPResponse(asn1Sequence)
	if err != nil {
		// Upstream logs "Impossible to create BasicOCSPResp from ASN1Sequence!".
		return nil
	}
	return NewBasicOCSPResp(basicOcspResponse)
}

// DSSRevocationUtilsOcspResp builds an OCSPResp from the DER of an OCSPResponse, returning
// nil when the binaries cannot be decoded. Port of getOcspResp(ASN1Sequence).
func DSSRevocationUtilsOcspResp(asn1Sequence []byte) *OCSPResp {
	ocspResponse, err := ParseOCSPResponse(asn1Sequence)
	if err != nil {
		// Upstream logs "Impossible to create OCSPResp from ASN1Sequence!".
		return nil
	}
	return NewOCSPResp(ocspResponse)
}

// DSSRevocationUtilsFromRespToBasic returns the BasicOCSPResp carried by an OCSPResp, or nil
// when the response carries another type. Port of fromRespToBasic(OCSPResp).
func DSSRevocationUtilsFromRespToBasic(ocspResp *OCSPResp) *BasicOCSPResp {
	basicOCSPResp, err := ocspResp.ResponseObject()
	if err != nil {
		// Upstream logs "Impossible to process OCSPResp!".
		return nil
	}
	if basicOCSPResp == nil {
		// Upstream logs "Unknown OCSP response type: {}".
		return nil
	}
	return basicOCSPResp
}

// DSSRevocationUtilsFromBasicToResp wraps a BasicOCSPResp in an OCSPResp whose status is
// SUCCESSFUL. Port of fromBasicToResp(BasicOCSPResp).
func DSSRevocationUtilsFromBasicToResp(basicOCSPResp *BasicOCSPResp) *OCSPResp {
	return DSSRevocationUtilsFromBasicBinariesToResp(basicOCSPResp.Encoded())
}

// DSSRevocationUtilsEncodedFromBasicResp returns the encoding of the OCSPResp wrapping the
// given BasicOCSPResp. Port of getEncodedFromBasicResp(BasicOCSPResp).
func DSSRevocationUtilsEncodedFromBasicResp(basicOCSPResp *BasicOCSPResp) ([]byte, error) {
	if basicOCSPResp == nil {
		return nil, model.NewDSSError("Empty OCSP response")
	}
	return DSSRevocationUtilsFromBasicToResp(basicOCSPResp).Encoded(), nil
}

// DSSRevocationUtilsFromBasicBinariesToResp wraps the binaries of a BasicOCSPResponse in an
// OCSPResp whose status is SUCCESSFUL. Port of fromBasicToResp(byte[]).
func DSSRevocationUtilsFromBasicBinariesToResp(basicOCSPRespBinary []byte) *OCSPResp {
	return NewOCSPResp(NewOCSPResponse(OCSPResponseStatusSuccessful,
		OCSPObjectIdentifierIDPkixOcspBasic, basicOCSPRespBinary))
}

// DSSRevocationUtilsUsedDigestAlgorithm returns the digest algorithm used in the given
// single response. Port of getUsedDigestAlgorithm(SingleResp).
func DSSRevocationUtilsUsedDigestAlgorithm(singleResp *SingleResp) (enumerations.DigestAlgorithm, error) {
	return enumerations.DigestAlgorithmForOID(singleResp.CertID().HashAlgOID().String())
}

// DSSRevocationUtilsMatches reports whether the certificate identified by certID is the one
// the single response is about. Port of matches(CertificateID, SingleResp).
//
// Upstream comment: certId.equals fails in comparing the algoIdentifier because the
// AlgorithmIdentifier parameters are null in one case and DERNull in another; the fields are
// therefore compared one by one.
func DSSRevocationUtilsMatches(certID *CertificateID, singleResp *SingleResp) bool {
	singleRespCertID := singleResp.CertID()
	return singleRespCertID.HashAlgOID().Equal(certID.HashAlgOID()) &&
		bytes.Equal(singleRespCertID.IssuerKeyHash(), certID.IssuerKeyHash()) &&
		bytes.Equal(singleRespCertID.IssuerNameHash(), certID.IssuerNameHash()) &&
		singleRespCertID.SerialNumber().Cmp(certID.SerialNumber()) == 0
}

// DSSRevocationUtilsOCSPCertificateID returns the CertificateID of a certificate and its
// issuer's certificate.
// Port of getOCSPCertificateID(CertificateToken, CertificateToken, DigestAlgorithm).
//
// It performs the work BouncyCastle's CertificateID(DigestCalculator, X509CertificateHolder,
// BigInteger) constructor does: issuerNameHash digests the DER of the issuer's subject name
// and issuerKeyHash digests the value bits of the issuer's subjectPublicKey. The hash
// algorithm identifier carries explicit DER NULL parameters, as the DigestCalculator
// upstream builds it does.
func DSSRevocationUtilsOCSPCertificateID(cert *model.CertificateToken, issuerCert *model.CertificateToken,
	digestAlgorithm enumerations.DigestAlgorithm) (*CertificateID, error) {
	oid, err := dssASN1UtilsObjectIdentifier(digestAlgorithm.OID())
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to create CertificateID", err)
	}
	issuerName, err := DSSASN1UtilsDEREncoded(issuerCert.Certificate().RawSubject)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to create CertificateID", err)
	}
	issuerNameHash, err := dssASN1UtilsDigest(digestAlgorithm, issuerName)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to create CertificateID", err)
	}
	publicKeyData, err := dssRevocationUtilsSubjectPublicKeyBits(issuerCert)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to create CertificateID", err)
	}
	issuerKeyHash, err := dssASN1UtilsDigest(digestAlgorithm, publicKeyData)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to create CertificateID", err)
	}
	return NewCertificateID(&CertID{
		HashAlgorithm:  NewAlgorithmIdentifierWithParameters(oid, dssASN1UtilsDERNull),
		IssuerNameHash: issuerNameHash,
		IssuerKeyHash:  issuerKeyHash,
		SerialNumber:   cert.SerialNumber(),
	}), nil
}

// dssRevocationUtilsSubjectPublicKeyBits returns the value bits of the certificate's
// subjectPublicKey BIT STRING, i.e. SubjectPublicKeyInfo#getPublicKeyData().getBytes().
func dssRevocationUtilsSubjectPublicKeyBits(certificate *model.CertificateToken) ([]byte, error) {
	element, _, err := dssASN1UtilsParse(certificate.PublicKey().Encoded())
	if err != nil {
		return nil, err
	}
	if !element.constructed || len(element.children) < 2 ||
		!element.children[1].isUniversal(dssASN1UtilsTagBitString) {
		return nil, errors.New("malformed SubjectPublicKeyInfo")
	}
	return element.children[1].bitStringOctets(), nil
}

// DSSRevocationUtilsLoadOCSPBase64Encoded loads an OCSP response from a base64-encoded
// string. Port of loadOCSPBase64Encoded(String).
func DSSRevocationUtilsLoadOCSPBase64Encoded(base64Encoded string) (*BasicOCSPResp, error) {
	derEncoded := utils.FromBase64(base64Encoded)
	return DSSRevocationUtilsLoadOCSPFromBinaries(derEncoded)
}

// DSSRevocationUtilsLoadOCSPFromBinaries loads an OCSP response from its binaries.
// Port of loadOCSPFromBinaries(byte[]).
func DSSRevocationUtilsLoadOCSPFromBinaries(binaries []byte) (*BasicOCSPResp, error) {
	ocspResp, err := NewOCSPRespFromBinaries(binaries)
	if err != nil {
		return nil, err
	}
	return DSSRevocationUtilsFromRespToBasic(ocspResp), nil
}

// DSSRevocationUtilsEncoded returns the encoded binaries of the OCSP response.
// Port of getEncoded(OCSPResp).
func DSSRevocationUtilsEncoded(ocspResp *OCSPResp) []byte {
	return ocspResp.Encoded()
}

// DSSRevocationUtilsDSSResponderIDFromRespID transforms a RespID into a ResponderId.
// Port of getDSSResponderId(RespID).
func DSSRevocationUtilsDSSResponderIDFromRespID(respID *RespID) (*ResponderId, error) {
	return DSSRevocationUtilsDSSResponderID(respID.ToASN1Primitive())
}

// DSSRevocationUtilsDSSResponderID transforms a ResponderID into a ResponderId.
// Port of getDSSResponderId(ResponderID).
func DSSRevocationUtilsDSSResponderID(responderID *ResponderID) (*ResponderId, error) {
	principal, err := DSSASN1UtilsToX500Principal(responderID.Name)
	if err != nil {
		return nil, err
	}
	return NewResponderId(principal, responderID.KeyHash), nil
}

// DSSRevocationUtilsCRLRevocationTokenKeys builds the revocation token keys of the CRLs
// referenced by the certificate. Port of getCRLRevocationTokenKeys(CertificateToken).
func DSSRevocationUtilsCRLRevocationTokenKeys(certificateToken *model.CertificateToken) []string {
	revocationKeys := make([]string, 0)
	for _, crlURL := range CertificateExtensionsUtilsCRLAccessUrls(certificateToken) {
		revocationKeys = append(revocationKeys, DSSRevocationUtilsCRLRevocationTokenKey(crlURL))
	}
	return revocationKeys
}

// DSSRevocationUtilsCRLRevocationTokenKey returns the CRL key, i.e. the hex-encoded SHA-1
// digest of the URL. Port of getCRLRevocationTokenKey(String).
//
// DSSUtilsSHA1Digest's error is a Go-only surface (digest() throwing NoSuchAlgorithmException
// is not reachable for SHA-1, which is always available); Java's method declares no checked
// exception either, so an error here is reproduced as a panic.
func DSSRevocationUtilsCRLRevocationTokenKey(crlURL string) string {
	digest, err := DSSUtilsSHA1Digest(crlURL)
	if err != nil {
		panic(err.Error())
	}
	return digest
}

// DSSRevocationUtilsOcspRevocationTokenKeys builds the revocation token keys of the OCSP
// responders referenced by the certificate.
// Port of getOcspRevocationTokenKeys(CertificateToken).
func DSSRevocationUtilsOcspRevocationTokenKeys(certificateToken *model.CertificateToken) []string {
	revocationKeys := make([]string, 0)
	for _, ocspURL := range CertificateExtensionsUtilsOCSPAccessUrls(certificateToken) {
		revocationKeys = append(revocationKeys, DSSRevocationUtilsOcspRevocationKey(certificateToken, ocspURL))
	}
	return revocationKeys
}

// DSSRevocationUtilsOcspRevocationKey returns the OCSP key, i.e. the hex-encoded SHA-1
// digest of "<token DSS id>:<url>". Port of getOcspRevocationKey(CertificateToken, String).
//
// See DSSRevocationUtilsCRLRevocationTokenKey for why DSSUtilsSHA1Digest's error is panicked
// rather than propagated.
func DSSRevocationUtilsOcspRevocationKey(certificateToken *model.CertificateToken, ocspURL string) string {
	digest, err := DSSUtilsSHA1Digest(certificateToken.DSSIDAsString() + ":" + ocspURL)
	if err != nil {
		panic(err.Error())
	}
	return digest
}

// DSSRevocationUtilsLatestSingleResponse returns the most recent single response of the OCSP
// response that concerns the given certificate, or nil when there is none.
// Port of getLatestSingleResponse(BasicOCSPResp, CertificateToken, CertificateToken).
func DSSRevocationUtilsLatestSingleResponse(basicResponse *BasicOCSPResp, certificate *model.CertificateToken,
	issuer *model.CertificateToken) *SingleResp {
	singleResponses := DSSRevocationUtilsSingleResponses(basicResponse, certificate, issuer)
	if len(singleResponses) == 0 {
		return nil
	} else if len(singleResponses) == 1 {
		return singleResponses[0]
	}
	return dssRevocationUtilsLatestSingleRespInList(singleResponses)
}

// dssRevocationUtilsLatestSingleRespInList ports the private getLatestSingleRespInList.
func dssRevocationUtilsLatestSingleRespInList(singleResponses []*SingleResp) *SingleResp {
	var latest time.Time
	var latestResp *SingleResp
	for _, singleResp := range singleResponses {
		thisUpdate := singleResp.ThisUpdate()
		if latestResp == nil || thisUpdate.After(latest) {
			latestResp = singleResp
			latest = thisUpdate
		}
	}
	return latestResp
}

// DSSRevocationUtilsSingleResponses returns every single response of the OCSP response that
// concerns the given certificate.
// Port of getSingleResponses(BasicOCSPResp, CertificateToken, CertificateToken).
func DSSRevocationUtilsSingleResponses(basicResponse *BasicOCSPResp, certificate *model.CertificateToken,
	issuer *model.CertificateToken) []*SingleResp {
	result := make([]*SingleResp, 0)
	for _, singleResp := range basicResponse.Responses() {
		usedDigestAlgorithm, err := DSSRevocationUtilsUsedDigestAlgorithm(singleResp)
		if err != nil {
			continue
		}
		certID, err := DSSRevocationUtilsOCSPCertificateID(certificate, issuer, usedDigestAlgorithm)
		if err != nil {
			continue
		}
		if DSSRevocationUtilsMatches(certID, singleResp) {
			result = append(result, singleResp)
		}
	}
	return result
}

// DSSRevocationUtilsDigest converts an OtherHash into a model.Digest, returning the empty
// Digest when the argument is nil (Java returns null). Port of getDigest(OtherHash).
func DSSRevocationUtilsDigest(otherHash *OtherHash) (model.Digest, error) {
	if otherHash == nil {
		return model.Digest{}, nil
	}
	digestAlgorithm, err := enumerations.DigestAlgorithmForOID(otherHash.HashAlgorithm.Algorithm.String())
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(digestAlgorithm, otherHash.HashValue), nil
}

// dssRevocationUtilsProducedRevocation is the slice of the RevocationToken contract
// checkIssuerValidAtRevocationProductionTime relies on. It is declared structurally so the
// method accepts whichever concrete revocation token the revocation chunk defines.
type dssRevocationUtilsProducedRevocation interface {
	// ProductionDate returns the time at which the revocation data was produced.
	ProductionDate() time.Time
}

// DSSRevocationUtilsCheckIssuerValidAtRevocationProductionTime reports whether the revocation
// data was produced within the validity range of its issuer's certificate.
// Port of checkIssuerValidAtRevocationProductionTime(RevocationToken, CertificateToken).
func DSSRevocationUtilsCheckIssuerValidAtRevocationProductionTime(
	revocationToken dssRevocationUtilsProducedRevocation, issuerCertificateToken *model.CertificateToken) bool {
	return issuerCertificateToken != nil && issuerCertificateToken.IsValidOn(revocationToken.ProductionDate())
}

// -----------------------------------------------------------------------------
// RFC 6960 decoding.
// -----------------------------------------------------------------------------

// dssRevocationUtilsOCSPResponseFromElement decodes an already parsed OCSPResponse.
func dssRevocationUtilsOCSPResponseFromElement(element *dssASN1UtilsElement) (*OCSPResponse, error) {
	if !element.constructed || len(element.children) == 0 || len(element.children) > 2 {
		return nil, errors.New("malformed OCSPResponse")
	}
	status := element.children[0]
	if status.class != 0 || (status.tagNumber != 0x0A && status.tagNumber != dssASN1UtilsTagInteger) {
		return nil, errors.New("malformed OCSPResponse: responseStatus is not an ENUMERATED")
	}
	response := &OCSPResponse{ResponseStatus: int(status.integer().Int64()), encoded: element.encoded}
	if len(element.children) == 1 {
		return response, nil
	}
	tagged := element.children[1]
	if tagged.class != 0x80 || tagged.tagNumber != 0 || !tagged.constructed || len(tagged.children) != 1 {
		return nil, errors.New("malformed OCSPResponse: responseBytes is not [0] EXPLICIT")
	}
	responseBytes := tagged.children[0]
	if !responseBytes.constructed || len(responseBytes.children) != 2 {
		return nil, errors.New("malformed ResponseBytes")
	}
	responseType, err := responseBytes.children[0].objectIdentifier()
	if err != nil {
		return nil, err
	}
	if !responseBytes.children[1].isUniversal(dssASN1UtilsTagOctetString) {
		return nil, errors.New("malformed ResponseBytes: response is not an OCTET STRING")
	}
	response.ResponseType = responseType
	response.Response = responseBytes.children[1].octets()
	return response, nil
}

// dssRevocationUtilsBasicOCSPResponseFromElement decodes an already parsed BasicOCSPResponse.
func dssRevocationUtilsBasicOCSPResponseFromElement(element *dssASN1UtilsElement) (*BasicOCSPResponse, error) {
	if !element.constructed || len(element.children) < 3 {
		return nil, errors.New("malformed BasicOCSPResponse")
	}
	responseData, err := dssRevocationUtilsResponseDataFromElement(element.children[0])
	if err != nil {
		return nil, err
	}
	signatureAlgorithm, err := dssASN1UtilsAlgorithmIdentifierFromElement(element.children[1])
	if err != nil {
		return nil, err
	}
	if !element.children[2].isUniversal(dssASN1UtilsTagBitString) {
		return nil, errors.New("malformed BasicOCSPResponse: signature is not a BIT STRING")
	}
	response := &BasicOCSPResponse{
		TBSResponseData:    responseData,
		SignatureAlgorithm: signatureAlgorithm,
		Signature:          element.children[2].bitStringOctets(),
		encoded:            element.encoded,
	}
	if len(element.children) > 3 {
		tagged := element.children[3]
		if tagged.class != 0x80 || tagged.tagNumber != 0 || !tagged.constructed || len(tagged.children) != 1 {
			return nil, errors.New("malformed BasicOCSPResponse: certs is not [0] EXPLICIT")
		}
		certificates := tagged.children[0]
		if !certificates.constructed {
			return nil, errors.New("malformed BasicOCSPResponse: certs is not a SEQUENCE")
		}
		for _, certificate := range certificates.children {
			response.Certs = append(response.Certs, certificate.encoded)
		}
	}
	return response, nil
}

// dssRevocationUtilsResponseDataFromElement decodes an already parsed ResponseData.
func dssRevocationUtilsResponseDataFromElement(element *dssASN1UtilsElement) (*ResponseData, error) {
	if !element.constructed || len(element.children) < 3 {
		return nil, errors.New("malformed ResponseData")
	}
	data := &ResponseData{encoded: element.encoded}
	index := 0
	if element.children[index].class == 0x80 && element.children[index].tagNumber == 0 {
		version := element.children[index]
		if !version.constructed || len(version.children) != 1 {
			return nil, errors.New("malformed ResponseData: version is not [0] EXPLICIT")
		}
		data.Version = int(version.children[0].integer().Int64())
		index++
	}
	if len(element.children) < index+3 {
		return nil, errors.New("malformed ResponseData")
	}

	responderID, err := dssRevocationUtilsResponderIDFromElement(element.children[index])
	if err != nil {
		return nil, err
	}
	data.ResponderID = responderID
	index++

	producedAt := element.children[index]
	if !producedAt.isUniversal(dssASN1UtilsTagGeneralizedTime) {
		return nil, errors.New("malformed ResponseData: producedAt is not a GeneralizedTime")
	}
	data.ProducedAt, err = dssASN1UtilsParseGeneralizedTime(string(producedAt.content))
	if err != nil {
		return nil, err
	}
	index++

	responses := element.children[index]
	if !responses.constructed {
		return nil, errors.New("malformed ResponseData: responses is not a SEQUENCE")
	}
	for _, response := range responses.children {
		singleResponse, err := dssRevocationUtilsSingleResponseFromElement(response)
		if err != nil {
			return nil, err
		}
		data.Responses = append(data.Responses, singleResponse)
	}
	index++

	if index < len(element.children) {
		extensions, err := dssRevocationUtilsExtensionsFromTagged(element.children[index], 1)
		if err != nil {
			return nil, err
		}
		data.ResponseExtensions = extensions
	}
	return data, nil
}

// dssRevocationUtilsResponderIDFromElement decodes the ResponderID CHOICE.
func dssRevocationUtilsResponderIDFromElement(element *dssASN1UtilsElement) (*ResponderID, error) {
	if element.class != 0x80 || !element.constructed || len(element.children) != 1 {
		return nil, errors.New("malformed ResponderID")
	}
	switch element.tagNumber {
	case 1:
		return &ResponderID{Name: element.children[0].encoded}, nil
	case 2:
		if !element.children[0].isUniversal(dssASN1UtilsTagOctetString) {
			return nil, errors.New("malformed ResponderID: byKey is not an OCTET STRING")
		}
		return &ResponderID{KeyHash: element.children[0].octets()}, nil
	}
	return nil, errors.New("malformed ResponderID: unknown CHOICE alternative")
}

// dssRevocationUtilsSingleResponseFromElement decodes an already parsed SingleResponse.
func dssRevocationUtilsSingleResponseFromElement(element *dssASN1UtilsElement) (*SingleResponse, error) {
	if !element.constructed || len(element.children) < 3 {
		return nil, errors.New("malformed SingleResponse")
	}
	certID, err := dssRevocationUtilsCertIDFromElement(element.children[0])
	if err != nil {
		return nil, err
	}
	certStatus, err := dssRevocationUtilsCertStatusFromElement(element.children[1])
	if err != nil {
		return nil, err
	}
	thisUpdate := element.children[2]
	if !thisUpdate.isUniversal(dssASN1UtilsTagGeneralizedTime) {
		return nil, errors.New("malformed SingleResponse: thisUpdate is not a GeneralizedTime")
	}
	thisUpdateDate, err := dssASN1UtilsParseGeneralizedTime(string(thisUpdate.content))
	if err != nil {
		return nil, err
	}
	response := &SingleResponse{CertID: certID, CertStatus: certStatus, ThisUpdate: thisUpdateDate}

	for _, child := range element.children[3:] {
		if child.class != 0x80 {
			return nil, errors.New("malformed SingleResponse: unexpected component")
		}
		switch child.tagNumber {
		case 0:
			if !child.constructed || len(child.children) != 1 ||
				!child.children[0].isUniversal(dssASN1UtilsTagGeneralizedTime) {
				return nil, errors.New("malformed SingleResponse: nextUpdate is not [0] EXPLICIT GeneralizedTime")
			}
			nextUpdate, err := dssASN1UtilsParseGeneralizedTime(string(child.children[0].content))
			if err != nil {
				return nil, err
			}
			response.NextUpdate = &nextUpdate
		case 1:
			extensions, err := dssRevocationUtilsExtensionsFromTagged(child, 1)
			if err != nil {
				return nil, err
			}
			response.SingleExtensions = extensions
		}
	}
	return response, nil
}

// dssRevocationUtilsCertIDFromElement decodes an already parsed CertID.
func dssRevocationUtilsCertIDFromElement(element *dssASN1UtilsElement) (*CertID, error) {
	if !element.constructed || len(element.children) != 4 {
		return nil, errors.New("malformed CertID")
	}
	hashAlgorithm, err := dssASN1UtilsAlgorithmIdentifierFromElement(element.children[0])
	if err != nil {
		return nil, err
	}
	if !element.children[1].isUniversal(dssASN1UtilsTagOctetString) ||
		!element.children[2].isUniversal(dssASN1UtilsTagOctetString) ||
		!element.children[3].isUniversal(dssASN1UtilsTagInteger) {
		return nil, errors.New("malformed CertID: unexpected component types")
	}
	return &CertID{
		HashAlgorithm:  hashAlgorithm,
		IssuerNameHash: element.children[1].octets(),
		IssuerKeyHash:  element.children[2].octets(),
		SerialNumber:   element.children[3].integer(),
	}, nil
}

// dssRevocationUtilsCertStatusFromElement decodes the CertStatus CHOICE, whose alternatives
// are implicitly tagged: [0] IMPLICIT NULL, [1] IMPLICIT RevokedInfo, [2] IMPLICIT UnknownInfo.
func dssRevocationUtilsCertStatusFromElement(element *dssASN1UtilsElement) (*OCSPCertStatus, error) {
	if element.class != 0x80 {
		return nil, errors.New("malformed CertStatus")
	}
	switch element.tagNumber {
	case OCSPCertStatusGood:
		return &OCSPCertStatus{Tag: OCSPCertStatusGood}, nil
	case OCSPCertStatusUnknown:
		return &OCSPCertStatus{Tag: OCSPCertStatusUnknown}, nil
	case OCSPCertStatusRevoked:
		if !element.constructed || len(element.children) == 0 {
			return nil, errors.New("malformed RevokedInfo")
		}
		revocationTimeElement := element.children[0]
		if !revocationTimeElement.isUniversal(dssASN1UtilsTagGeneralizedTime) {
			return nil, errors.New("malformed RevokedInfo: revocationTime is not a GeneralizedTime")
		}
		revocationTime, err := dssASN1UtilsParseGeneralizedTime(string(revocationTimeElement.content))
		if err != nil {
			return nil, err
		}
		status := &OCSPCertStatus{Tag: OCSPCertStatusRevoked, RevocationTime: revocationTime}
		if len(element.children) > 1 {
			reasonElement := element.children[1]
			if reasonElement.class != 0x80 || reasonElement.tagNumber != 0 ||
				!reasonElement.constructed || len(reasonElement.children) != 1 {
				return nil, errors.New("malformed RevokedInfo: revocationReason is not [0] EXPLICIT")
			}
			reason := int(reasonElement.children[0].integer().Int64())
			status.RevocationReason = &reason
		}
		return status, nil
	}
	return nil, errors.New("malformed CertStatus: unknown CHOICE alternative")
}

// dssRevocationUtilsExtensionsFromTagged decodes an [n] EXPLICIT Extensions component.
func dssRevocationUtilsExtensionsFromTagged(element *dssASN1UtilsElement, tagNumber uint64) ([]Extension, error) {
	if element.class != 0x80 || element.tagNumber != tagNumber || !element.constructed || len(element.children) != 1 {
		return nil, errors.New("malformed Extensions: not an EXPLICIT tagged SEQUENCE")
	}
	sequence := element.children[0]
	if !sequence.constructed {
		return nil, errors.New("malformed Extensions: not a SEQUENCE")
	}
	extensions := make([]Extension, 0, len(sequence.children))
	for _, child := range sequence.children {
		if !child.constructed || len(child.children) < 2 || len(child.children) > 3 {
			return nil, errors.New("malformed Extension")
		}
		oid, err := child.children[0].objectIdentifier()
		if err != nil {
			return nil, err
		}
		extension := Extension{ID: oid}
		valueIndex := 1
		if child.children[1].isUniversal(dssASN1UtilsTagBoolean) {
			extension.Critical = len(child.children[1].content) > 0 && child.children[1].content[0] != 0x00
			valueIndex = 2
		}
		if valueIndex >= len(child.children) || !child.children[valueIndex].isUniversal(dssASN1UtilsTagOctetString) {
			return nil, errors.New("malformed Extension: extnValue is not an OCTET STRING")
		}
		extension.Value = child.children[valueIndex].octets()
		extensions = append(extensions, extension)
	}
	return extensions, nil
}
