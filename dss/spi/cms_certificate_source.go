// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CMSCertificateSource.java (DSS 6.5.RC1).
//
// BouncyCastle replacements used here (see PORTING.md):
//
//   - org.bouncycastle.cms.SignerInformationStore   -> []*cmscore.SignerInfo
//   - org.bouncycastle.cms.SignerInformation        -> *cmscore.SignerInfo
//   - org.bouncycastle.util.Store<X509CertificateHolder> for SignedData.certificates ->
//     [][]byte, the encoding of every plain X.509 certificate of the set, i.e. exactly what
//     cmscore.CMS.Certificates() hands out. Store#getMatches(null) is then simply iterating
//     the slice.
//   - org.bouncycastle.asn1.cms.AttributeTable      -> cmscore.Attributes
//   - org.bouncycastle.asn1.ess.{SigningCertificate, SigningCertificateV2, ESSCertID,
//     ESSCertIDv2, OtherCertID} -> the types defined below. Go has no equivalent, so - as
//     crl_ref.go and ocsp_ref.go already do for the ESF structures - they live next to their
//     consumer rather than being invented as a new shared package.
//
// This file also carries the DSSASN1Utils methods that take a BouncyCastle CMS/ESS type this
// port only introduces here: toSignerIdentifier(SignerId), getCertificate(X509CertificateHolder),
// getAsn1Attributes(AttributeTable, ASN1ObjectIdentifier), getAsn1Encodable(Attribute) and
// getCertificateRef(OtherCertID). They keep the flattened static-utility naming
// (DSSASN1Utils<MethodName>, "get" dropped) so that folding them back into dss_asn1_utils.go
// later is a pure move.
//
// DER preservation: every certificate, attribute value and IssuerSerial is read from the bytes
// it arrived in and never re-encoded, since those bytes feed CertificateToken identifiers and
// certificate-reference digests.
package spi

import (
	"bytes"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
)

// The signed and unsigned attribute types this file reads. They are
// org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers constants upstream, hanging off
// id-aa OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840) rsadsi(113549) pkcs(1)
// pkcs-9(9) smime(16) 2 }, i.e. 1.2.840.113549.1.9.16.2; they keep their exact Java field
// name behind the "OID_" prefix, the convention oid.go already applies to DSS's own OID.java
// constants.
var (
	// OIDIdAaSigningCertificate is
	// id-aa-signingCertificate OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
	// rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 12 }
	OIDIdAaSigningCertificate = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 12}

	// OIDIdAaSigningCertificateV2 is
	// id-aa-signingCertificateV2 OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
	// rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 47 }
	OIDIdAaSigningCertificateV2 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 47}

	// OIDIdAaEtsCertificateRefs is
	// id-aa-ets-certificateRefs OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
	// rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 21 }
	OIDIdAaEtsCertificateRefs = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 21}

	// OIDIdAaEtsCertValues is
	// id-aa-ets-certValues OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
	// rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 23 }
	OIDIdAaEtsCertValues = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 23}
)

// essCertIDv2DefaultHashAlgorithm is the DEFAULT of ESSCertIDv2.hashAlgorithm, i.e.
// {algorithm id-sha256} without parameters - BouncyCastle's ESSCertIDv2.DEFAULT_ALG_ID.
var essCertIDv2DefaultHashAlgorithm = NewAlgorithmIdentifier(asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1})

// -----------------------------------------------------------------------------
// The RFC 5035 (ESS) structures the signing-certificate and certificate-reference
// attributes carry, replacing org.bouncycastle.asn1.ess.
// -----------------------------------------------------------------------------

// ESSCertID is
//
//	ESSCertID ::= SEQUENCE {
//	    certHash     Hash,                 -- SHA-1 hash of the certificate
//	    issuerSerial IssuerSerial OPTIONAL }
//
// replacing org.bouncycastle.asn1.ess.ESSCertID.
type ESSCertID struct {
	// CertHash is the SHA-1 digest of the referenced certificate.
	CertHash []byte
	// IssuerSerial identifies the referenced certificate, nil when absent.
	IssuerSerial *IssuerSerial
}

// ESSCertIDv2 is
//
//	ESSCertIDv2 ::= SEQUENCE {
//	    hashAlgorithm AlgorithmIdentifier DEFAULT {algorithm id-sha256},
//	    certHash      Hash,
//	    issuerSerial  IssuerSerial OPTIONAL }
//
// replacing org.bouncycastle.asn1.ess.ESSCertIDv2.
type ESSCertIDv2 struct {
	// HashAlgorithm is the algorithm CertHash was computed with; an absent field means
	// id-sha256, the ASN.1 DEFAULT.
	HashAlgorithm *AlgorithmIdentifier
	// CertHash is the digest of the referenced certificate.
	CertHash []byte
	// IssuerSerial identifies the referenced certificate, nil when absent.
	IssuerSerial *IssuerSerial
}

// SigningCertificate is
//
//	SigningCertificate ::= SEQUENCE {
//	    certs    SEQUENCE OF ESSCertID,
//	    policies SEQUENCE OF PolicyInformation OPTIONAL }
//
// replacing org.bouncycastle.asn1.ess.SigningCertificate.
type SigningCertificate struct {
	// certs holds the encoding of every member of the certs field, decoded on demand by
	// Certs() exactly as BouncyCastle's getCerts() does.
	certs [][]byte
	// Policies is the encoding of the policies field, nil when absent.
	Policies []byte
}

// SigningCertificateV2 is
//
//	SigningCertificateV2 ::= SEQUENCE {
//	    certs    SEQUENCE OF ESSCertIDv2,
//	    policies SEQUENCE OF PolicyInformation OPTIONAL }
//
// replacing org.bouncycastle.asn1.ess.SigningCertificateV2.
type SigningCertificateV2 struct {
	// certs holds the encoding of every member of the certs field, decoded on demand by
	// Certs().
	certs [][]byte
	// Policies is the encoding of the policies field, nil when absent.
	Policies []byte
}

// OtherCertID is
//
//	OtherCertID ::= SEQUENCE {
//	    otherCertHash OtherHash,
//	    issuerSerial  IssuerSerial OPTIONAL }
//
// replacing org.bouncycastle.asn1.ess.OtherCertID. Its otherCertHash is the same CHOICE the
// ESF structures use, so it reuses the OtherHash of dss_revocation_utils.go - which already
// resolves the sha1Hash alternative to id-sha1, exactly as OtherCertID#getAlgorithmHash does.
type OtherCertID struct {
	// OtherCertHash is the digest of the referenced certificate and the algorithm it was
	// computed with.
	OtherCertHash *OtherHash
	// IssuerSerial identifies the referenced certificate, nil when absent.
	IssuerSerial *IssuerSerial
}

// ParseSigningCertificate decodes a SigningCertificate from its encoding.
// Port of SigningCertificate.getInstance(Object).
func ParseSigningCertificate(encoded []byte) (*SigningCertificate, error) {
	children, err := cmsCertificateSourceSequence(encoded, "SigningCertificate", 1, 2)
	if err != nil {
		return nil, err
	}
	certs, err := cmsCertificateSourceMembers(children[0], "SigningCertificate.certs")
	if err != nil {
		return nil, err
	}
	signingCertificate := &SigningCertificate{certs: certs}
	if len(children) > 1 {
		signingCertificate.Policies = children[1].Encoded()
	}
	return signingCertificate, nil
}

// Certs decodes and returns the members of the certs field.
// Port of SigningCertificate#getCerts(), which likewise decodes on each call and lets a
// malformed member abort the whole conversion.
func (s *SigningCertificate) Certs() ([]*ESSCertID, error) {
	certs := make([]*ESSCertID, 0, len(s.certs))
	for _, encoded := range s.certs {
		cert, err := ParseESSCertID(encoded)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	return certs, nil
}

// ParseSigningCertificateV2 decodes a SigningCertificateV2 from its encoding.
// Port of SigningCertificateV2.getInstance(Object).
func ParseSigningCertificateV2(encoded []byte) (*SigningCertificateV2, error) {
	children, err := cmsCertificateSourceSequence(encoded, "SigningCertificateV2", 1, 2)
	if err != nil {
		return nil, err
	}
	certs, err := cmsCertificateSourceMembers(children[0], "SigningCertificateV2.certs")
	if err != nil {
		return nil, err
	}
	signingCertificate := &SigningCertificateV2{certs: certs}
	if len(children) > 1 {
		signingCertificate.Policies = children[1].Encoded()
	}
	return signingCertificate, nil
}

// Certs decodes and returns the members of the certs field.
// Port of SigningCertificateV2#getCerts().
func (s *SigningCertificateV2) Certs() ([]*ESSCertIDv2, error) {
	certs := make([]*ESSCertIDv2, 0, len(s.certs))
	for _, encoded := range s.certs {
		cert, err := ParseESSCertIDv2(encoded)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	return certs, nil
}

// ParseESSCertID decodes an ESSCertID from its encoding.
// Port of ESSCertID.getInstance(Object).
func ParseESSCertID(encoded []byte) (*ESSCertID, error) {
	children, err := cmsCertificateSourceSequence(encoded, "ESSCertID", 1, 2)
	if err != nil {
		return nil, err
	}
	if !children[0].IsUniversal(asn1ber.TagOctetString) {
		return nil, errors.New("malformed ESSCertID: certHash is not an OCTET STRING")
	}
	essCertID := &ESSCertID{CertHash: children[0].Octets()}
	if len(children) > 1 {
		issuerSerial := DSSASN1UtilsIssuerSerial(children[1].Encoded())
		if issuerSerial == nil {
			return nil, errors.New("malformed ESSCertID: the issuerSerial cannot be decoded")
		}
		essCertID.IssuerSerial = issuerSerial
	}
	return essCertID, nil
}

// ParseESSCertIDv2 decodes an ESSCertIDv2 from its encoding.
// Port of ESSCertIDv2.getInstance(Object), the DEFAULT of hashAlgorithm included: a first
// member that is already the certHash OCTET STRING means the field was omitted and the
// algorithm is id-sha256.
func ParseESSCertIDv2(encoded []byte) (*ESSCertIDv2, error) {
	children, err := cmsCertificateSourceSequence(encoded, "ESSCertIDv2", 1, 3)
	if err != nil {
		return nil, err
	}
	index := 0
	essCertID := &ESSCertIDv2{}
	if children[0].IsUniversal(asn1ber.TagOctetString) {
		essCertID.HashAlgorithm = essCertIDv2DefaultHashAlgorithm
	} else {
		hashAlgorithm, err := asn1ber.AlgorithmIdentifierFromElement(children[index])
		if err != nil {
			return nil, err
		}
		essCertID.HashAlgorithm = hashAlgorithm
		index++
	}
	if index >= len(children) {
		return nil, errors.New("malformed ESSCertIDv2: the certHash is missing")
	}
	if !children[index].IsUniversal(asn1ber.TagOctetString) {
		return nil, errors.New("malformed ESSCertIDv2: certHash is not an OCTET STRING")
	}
	essCertID.CertHash = children[index].Octets()
	index++
	if index < len(children) {
		issuerSerial := DSSASN1UtilsIssuerSerial(children[index].Encoded())
		if issuerSerial == nil {
			return nil, errors.New("malformed ESSCertIDv2: the issuerSerial cannot be decoded")
		}
		essCertID.IssuerSerial = issuerSerial
	}
	return essCertID, nil
}

// ParseOtherCertID decodes an OtherCertID from its encoding.
// Port of OtherCertID.getInstance(Object).
func ParseOtherCertID(encoded []byte) (*OtherCertID, error) {
	children, err := cmsCertificateSourceSequence(encoded, "OtherCertID", 1, 2)
	if err != nil {
		return nil, err
	}
	otherCertHash, err := ParseOtherHash(children[0].Encoded())
	if err != nil {
		return nil, err
	}
	otherCertID := &OtherCertID{OtherCertHash: otherCertHash}
	if len(children) > 1 {
		issuerSerial := DSSASN1UtilsIssuerSerial(children[1].Encoded())
		if issuerSerial == nil {
			return nil, errors.New("malformed OtherCertID: the issuerSerial cannot be decoded")
		}
		otherCertID.IssuerSerial = issuerSerial
	}
	return otherCertID, nil
}

// cmsCertificateSourceSequence parses a SEQUENCE and checks its size, reproducing the
// "Bad sequence size" guard every org.bouncycastle.asn1.ess structure opens with.
func cmsCertificateSourceSequence(encoded []byte, name string, min, max int) ([]*asn1ber.Element, error) {
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("extra data found after the %s", name)
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, fmt.Errorf("malformed %s: not a SEQUENCE", name)
	}
	if len(element.Children()) < min || len(element.Children()) > max {
		return nil, fmt.Errorf("bad sequence size in %s: %d", name, len(element.Children()))
	}
	return element.Children(), nil
}

// cmsCertificateSourceMembers returns the encoding of every member of a SEQUENCE OF field.
func cmsCertificateSourceMembers(element *asn1ber.Element, name string) ([][]byte, error) {
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, fmt.Errorf("malformed %s: not a SEQUENCE OF", name)
	}
	members := make([][]byte, 0, len(element.Children()))
	for _, child := range element.Children() {
		members = append(members, child.Encoded())
	}
	return members, nil
}

// -----------------------------------------------------------------------------
// The DSSASN1Utils methods dss_asn1_utils.go deferred to this phase.
// -----------------------------------------------------------------------------

// DSSASN1UtilsToSignerIdentifierFromSignerID transforms a CMS SignerIdentifier into a
// SignerIdentifier. Port of toSignerIdentifier(SignerId).
//
// BouncyCastle's SignerId carries the issuer name and the serial number for the
// issuerAndSerialNumber alternative and the key identifier for the subjectKeyIdentifier one,
// the unused fields being null; cmscore.SignerIdentifier models the same CHOICE, so exactly
// one of the two branches below applies.
func DSSASN1UtilsToSignerIdentifierFromSignerID(signerID *cmscore.SignerIdentifier) (*SignerIdentifier, error) {
	var issuerName []byte
	var serialNumber *big.Int
	var ski []byte
	if signerID.IssuerAndSerialNumber != nil {
		issuerName = signerID.IssuerAndSerialNumber.Issuer
		serialNumber = signerID.IssuerAndSerialNumber.SerialNumber
	} else {
		ski = signerID.SubjectKeyIdentifier
	}
	issuerX500Principal, err := DSSASN1UtilsToX500Principal(issuerName)
	if err != nil {
		return nil, err
	}
	return DSSASN1UtilsToSignerIdentifier(issuerX500Principal, serialNumber, ski), nil
}

// DSSASN1UtilsCertificate builds a CertificateToken from the encoding of an X.509
// certificate. Port of getCertificate(X509CertificateHolder), whose holder argument is here
// the certificate's DER (see the DEVIATION in dss_certificate_token_security_factory.go).
func DSSASN1UtilsCertificate(x509CertificateHolder []byte) (*model.CertificateToken, error) {
	return DSSCertificateTokenSecurityFactoryX509CertificateHolderInstance.Build(x509CertificateHolder)
}

// DSSASN1UtilsAsn1Attributes returns every attribute of the given type held by the attribute
// table, an empty slice when there is none. Port of
// getAsn1Attributes(AttributeTable, ASN1ObjectIdentifier).
//
// Java dereferences attributeTable, so a null table is a NullPointerException there; a nil
// cmscore.Attributes is simply an empty table here, and GetAll answers an empty slice. The
// two callers that may hold a null table (CMSCRLSource/CMSOCSPSource) guard it themselves.
func DSSASN1UtilsAsn1Attributes(attributeTable cmscore.Attributes, oid asn1.ObjectIdentifier) []*cmscore.Attribute {
	attributes := attributeTable.GetAll(oid)
	if attributes == nil {
		return []*cmscore.Attribute{}
	}
	return attributes
}

// DSSASN1UtilsAsn1Encodable returns the single value of an attribute, nil when the attribute
// is nil or carries anything other than exactly one value.
// Port of getAsn1Encodable(Attribute); the returned ASN1Encodable is here the parsed value,
// which keeps the bytes it was decoded from.
func DSSASN1UtilsAsn1Encodable(attribute *cmscore.Attribute) *asn1ber.Element {
	if attribute == nil {
		return nil
	}
	if len(attribute.Values) == 1 {
		return attribute.Values[0]
	}
	// Upstream logs "Only one value is allowed within attribute set! Found {}.".
	return nil
}

// DSSASN1UtilsCertificateRef converts an OtherCertID into a CertificateRef.
// Port of getCertificateRef(OtherCertID); the IllegalArgumentException
// DigestAlgorithm.forOID raises for an unknown digest algorithm is returned as an error.
func DSSASN1UtilsCertificateRef(otherCertID *OtherCertID) (*CertificateRef, error) {
	certRef := NewCertificateRef()
	digest, err := DSSRevocationUtilsDigest(otherCertID.OtherCertHash)
	if err != nil {
		return nil, err
	}
	certRef.SetCertDigest(digest)
	certRef.SetCertificateIdentifier(DSSASN1UtilsToSignerIdentifierFromIssuerSerial(otherCertID.IssuerSerial))
	return certRef, nil
}

// -----------------------------------------------------------------------------
// CMSCertificateSource.
// -----------------------------------------------------------------------------

// CMSCertificateSource is a CMS certificate source. Port of the abstract class
// CMSCertificateSource; the concrete sources of the later phases (CAdESCertificateSource,
// TimestampCertificateSource) embed the *CMSCertificateSource NewCMSCertificateSource returns
// and override only CertificateSourceType.
//
// The constructor already registers the source with its SignatureCertificateSource base, so
// CandidatesForSigningCertificate dispatches to ExtractCandidatesForSigningCertificate below.
// A subclass that overrides that method - none does in DSS 6.5.RC1 - has to re-register
// itself with InitSignatureCertificateSource.
type CMSCertificateSource struct {
	SignatureCertificateSource

	// signerInformations are the signers present in a CMS.
	signerInformations []*cmscore.SignerInfo

	// certificates are the certificates present within the SignedData.certificates field,
	// each as its own encoding.
	certificates [][]byte

	// currentSignerInformation is the SignerInformation of the current signature.
	currentSignerInformation *cmscore.SignerInfo
}

// NewCMSCertificateSource instantiates a CMS certificate source for a given signer and
// extracts every certificate, certificate identifier and certificate reference it holds.
// Port of the protected CMSCertificateSource(SignerInformationStore, Store<X509CertificateHolder>,
// SignerInformation) constructor.
//
// Panics with the Java message when currentSignerInformation is missing
// (Objects.requireNonNull). DEVIATION: the two other Objects.requireNonNull guards cannot be
// reproduced - a Java Store is an object that is either present or null, while its Go
// counterpart is a slice whose nil value is indistinguishable from an empty store, and
// cmscore hands out a nil slice for a CMS that carries no certificate and no signer.
//
// The DSSException toX500Principal raises for an undecodable issuer name propagates out of
// the Java constructor; here it is returned as an error. Every other failure of the
// extraction is swallowed exactly where upstream swallows it.
func NewCMSCertificateSource(signerInformations []*cmscore.SignerInfo, certificates [][]byte,
	currentSignerInformation *cmscore.SignerInfo) (*CMSCertificateSource, error) {
	if currentSignerInformation == nil {
		panic("currentSignerInformation is null, it must be provided!")
	}

	source := &CMSCertificateSource{
		signerInformations:       signerInformations,
		certificates:             certificates,
		currentSignerInformation: currentSignerInformation,
	}
	source.InitSignatureCertificateSource(source)

	if err := source.extractCertificateIdentifiers(); err != nil {
		return nil, err
	}
	source.extractSignedCertificates()
	source.extractSigningCertificateReferences()

	source.extractCertificateValues()
	source.extractCertificateRefsFromUnsignedAttribute(OIDIdAaEtsCertificateRefs,
		enumerations.CertificateRefOriginCompleteCertificateRefs)
	source.extractCertificateRefsFromUnsignedAttribute(OIDAttributeCertificateRefsOid,
		enumerations.CertificateRefOriginAttributeCertificateRefs)
	return source, nil
}

// extractCertificateIdentifiers ports the private method of the same name.
func (s *CMSCertificateSource) extractCertificateIdentifiers() error {
	currentSignerIdentifier, err := DSSASN1UtilsToSignerIdentifierFromSignerID(s.currentSignerInformation.SID)
	if err != nil {
		return err
	}
	found := false
	for _, signerInformation := range s.signerInformations {
		signerIdentifier, err := DSSASN1UtilsToSignerIdentifierFromSignerID(signerInformation.SID)
		if err != nil {
			return err
		}
		if signerIdentifier.IsEquivalent(currentSignerIdentifier) {
			signerIdentifier.SetCurrent(true)
			found = true
		}
		s.AddCertificateIdentifier(signerIdentifier, enumerations.CertificateOriginSignedData)
	}
	if !found {
		// Upstream logs "SID not found in SignerInfos".
		currentSignerIdentifier.SetCurrent(true)
		s.AddCertificateIdentifier(currentSignerIdentifier, enumerations.CertificateOriginSignedData)
	}
	return nil
}

// extractSignedCertificates ports the private method of the same name.
//
// Upstream wraps the whole loop in one try, so the first unreadable certificate ends the
// extraction rather than merely skipping that entry; the early return reproduces that.
func (s *CMSCertificateSource) extractSignedCertificates() {
	for _, x509CertificateHolder := range s.certificates {
		certificate, err := DSSASN1UtilsCertificate(x509CertificateHolder)
		if err != nil {
			// Upstream logs "Cannot extract certificates from CMS Signed Data : {}".
			return
		}
		s.AddCertificateWithOrigin(certificate, enumerations.CertificateOriginSignedData)
	}
}

// extractSigningCertificateReferences ports the private method of the same name.
func (s *CMSCertificateSource) extractSigningCertificateReferences() {
	signedAttributes := s.currentSignerInformation.SignedAttributes
	if len(signedAttributes) > 0 {
		for _, signingCertificateV1Attribute := range signedAttributes.GetAll(OIDIdAaSigningCertificate) {
			s.extractSigningCertificateV1(signingCertificateV1Attribute)
		}
		for _, signingCertificateV2Attribute := range signedAttributes.GetAll(OIDIdAaSigningCertificateV2) {
			s.extractSigningCertificateV2(signingCertificateV2Attribute)
		}
	}
}

// extractSigningCertificateV1 ports the private extractSigningCertificateV1(Attribute).
func (s *CMSCertificateSource) extractSigningCertificateV1(attribute *cmscore.Attribute) {
	for _, asn1Encodable := range attribute.Values {
		signingCertificate, err := ParseSigningCertificate(asn1Encodable.Encoded())
		if err != nil {
			// Upstream logs "SigningCertificate attribute '{}' is not well defined!".
			continue
		}
		certs, err := signingCertificate.Certs()
		if err != nil {
			continue
		}
		s.extractESSCertIDs(certs, enumerations.CertificateRefOriginSigningCertificate)
	}
}

// extractESSCertIDs ports the private extractESSCertIDs(ESSCertID[], CertificateRefOrigin).
func (s *CMSCertificateSource) extractESSCertIDs(essCertIDs []*ESSCertID, origin enumerations.CertificateRefOrigin) {
	for _, essCertID := range essCertIDs {
		certRef := NewCertificateRef()

		certHash := essCertID.CertHash
		if len(certHash) > 0 {
			certRef.SetCertDigest(model.NewDigest(enumerations.DigestAlgorithmSHA1, certHash))
			// Upstream logs "Found Certificate Hash in signingCertificateAttributeV1 {} with
			// algorithm {}" in debug.
		}
		certRef.SetCertificateIdentifier(DSSASN1UtilsToSignerIdentifierFromIssuerSerial(essCertID.IssuerSerial))
		s.AddCertificateRef(certRef, origin)
	}
}

// extractSigningCertificateV2 ports the private extractSigningCertificateV2(Attribute).
func (s *CMSCertificateSource) extractSigningCertificateV2(attribute *cmscore.Attribute) {
	for _, asn1Encodable := range attribute.Values {
		signingCertificate, err := ParseSigningCertificateV2(asn1Encodable.Encoded())
		if err != nil {
			// Upstream logs "SigningCertificateV2 attribute '{}' is not well defined!".
			continue
		}
		certs, err := signingCertificate.Certs()
		if err != nil {
			continue
		}
		if err := s.extractESSCertIDv2s(certs, enumerations.CertificateRefOriginSigningCertificate); err != nil {
			// DigestAlgorithm.forOID's IllegalArgumentException lands in the same catch.
			continue
		}
	}
}

// extractESSCertIDv2s ports the private extractESSCertIDv2s(ESSCertIDv2[], CertificateRefOrigin).
//
// Upstream builds every CertificateRef inside the caller's try, so an unknown digest
// algorithm aborts the remaining members of the attribute value; the returned error
// reproduces that, the references added so far staying in place.
func (s *CMSCertificateSource) extractESSCertIDv2s(essCertIDv2s []*ESSCertIDv2, origin enumerations.CertificateRefOrigin) error {
	for _, essCertIDv2 := range essCertIDv2s {
		certRef := NewCertificateRef()
		digestAlgorithm, err := enumerations.DigestAlgorithmForOID(essCertIDv2.HashAlgorithm.Algorithm.String())
		if err != nil {
			return err
		}
		certHash := essCertIDv2.CertHash
		certRef.SetCertDigest(model.NewDigest(digestAlgorithm, certHash))
		// Upstream logs "Found Certificate Hash in SigningCertificateV2 {} with algorithm {}"
		// in debug.
		certRef.SetCertificateIdentifier(DSSASN1UtilsToSignerIdentifierFromIssuerSerial(essCertIDv2.IssuerSerial))
		s.AddCertificateRef(certRef, origin)
	}
	return nil
}

// extractCertificateValues ports the private extractCertificateValues().
func (s *CMSCertificateSource) extractCertificateValues() {
	// Java guards this with "unsignedAttributes != null"; a nil cmscore.Attributes stands
	// for both an absent and an empty attribute table, and DSSASN1UtilsAsn1Attributes
	// answers an empty slice for either, so the guard is redundant here.
	unsignedAttributes := s.currentSignerInformation.UnsignedAttributes
	attributes := DSSASN1UtilsAsn1Attributes(unsignedAttributes, OIDIdAaEtsCertValues)
	for _, attribute := range attributes {
		s.extractCertificateValuesFromAttribute(attribute)
	}
}

// extractCertificateValuesFromAttribute ports the private extractCertificateValues(Attribute)
// overload.
func (s *CMSCertificateSource) extractCertificateValuesFromAttribute(attribute *cmscore.Attribute) {
	attrValue := DSSASN1UtilsAsn1Encodable(attribute)
	if attrValue == nil {
		return
	}
	if attrValue.IsUniversal(asn1ber.TagSequence) && attrValue.IsConstructed() {
		for _, member := range attrValue.Children() {
			// Java re-encodes the Certificate it decoded (Certificate.getInstance(...)
			// .getEncoded()); the member's own bytes are used here instead, which is the
			// same for the DER a certificate has to be in and preserves the encoding the
			// CertificateToken identifier is computed over.
			certificate, err := DSSUtilsLoadCertificateFromBinary(member.Encoded())
			if err != nil {
				// Upstream logs "Unable to parse encapsulated certificate : {}".
				continue
			}
			s.AddCertificateWithOrigin(certificate, enumerations.CertificateOriginCertificateValues)
		}
	}
	// Upstream logs "Certificate values shall be encoded as an ASN1Sequence. Found encoding :
	// {}" otherwise.
}

// extractCertificateRefsFromUnsignedAttribute ports the private method of the same name.
func (s *CMSCertificateSource) extractCertificateRefsFromUnsignedAttribute(attributeOid asn1.ObjectIdentifier,
	origin enumerations.CertificateRefOrigin) {
	// See extractCertificateValues on the redundant "unsignedAttributes != null" guard.
	unsignedAttributes := s.currentSignerInformation.UnsignedAttributes
	attributes := DSSASN1UtilsAsn1Attributes(unsignedAttributes, attributeOid)
	for _, attribute := range attributes {
		attrValue := DSSASN1UtilsAsn1Encodable(attribute)
		if attrValue == nil {
			continue
		}
		if attrValue.IsUniversal(asn1ber.TagSequence) && attrValue.IsConstructed() {
			for _, member := range attrValue.Children() {
				otherCertID, err := ParseOtherCertID(member.Encoded())
				if err != nil {
					// Upstream logs "Unable to parse encapsulated OtherCertID : {}".
					continue
				}
				certRef, err := DSSASN1UtilsCertificateRef(otherCertID)
				if err != nil {
					continue
				}
				s.AddCertificateRef(certRef, origin)
			}
		}
		// Upstream logs "Certificate values shall be encoded as an ASN1Sequence. Found
		// encoding : {}" otherwise.
	}
}

// ExtractCandidatesForSigningCertificate extracts the candidates to be the signing
// certificate. Port of the protected extractCandidatesForSigningCertificate(CertificateSource)
// override; signingCertificateSource may be nil.
//
// Two Java behaviours have no exact Go counterpart and are resolved as follows:
//   - getBySignerIdentifier returns a Set whose iteration order Java leaves unspecified;
//     the Go map is walked in a deterministic order (by DSS identifier) so that a repeated
//     run of the same input yields the same candidate.
//   - CertificateToken#getDigest cannot fail for an algorithm that was itself resolved from
//     an OID, but its Go counterpart returns an error; a failure leaves digestEqual false.
func (s *CMSCertificateSource) ExtractCandidatesForSigningCertificate(
	signingCertificateSource CertificateSource) *CandidatesForSigningCertificate {
	candidates := NewCandidatesForSigningCertificate()

	currentSignerIdentifier := s.CurrentCertificateIdentifier()
	if currentSignerIdentifier != nil && !currentSignerIdentifier.IsEmpty() {
		certificate := s.CertificateTokenBySignerIdentifier(currentSignerIdentifier)
		if certificate == nil && signingCertificateSource != nil {
			foundTokens := signingCertificateSource.BySignerIdentifier(currentSignerIdentifier)
			if len(foundTokens) > 0 {
				// Upstream logs "Resolved signing certificate by certificate identifier".
				certificate = cmsCertificateSourceFirstToken(foundTokens)
			}
		}

		var certificateValidity *CertificateValidity
		if certificate != nil {
			certificateValidity = NewCertificateValidity(certificate)
		} else {
			certificateValidity = NewCertificateValidityFromSignerIdentifier(currentSignerIdentifier)
		}

		signingCertRefs := s.SigningCertificateRefs()
		if len(signingCertRefs) > 0 {
			// first one
			signingCertRef := signingCertRefs[0]
			certDigest := signingCertRef.CertDigest()
			certificateValidity.SetDigestPresent(!certDigest.IsEmpty())

			if certificate != nil {
				certificateDigest, err := certificate.Digest(certDigest.Algorithm())
				certificateValidity.SetDigestEqual(err == nil && bytes.Equal(certificateDigest, certDigest.Value()))
			}

			sigCertIdentifier := signingCertRef.CertificateIdentifier()
			certificateValidity.SetIssuerSerialPresent(sigCertIdentifier != nil)
			if sigCertIdentifier != nil {
				if certificate != nil {
					certificateValidity.SetSerialNumberEqual(
						cmsCertificateSourceSerialsEqual(certificate.SerialNumber(), sigCertIdentifier.SerialNumber()))
					certificateValidity.SetDistinguishedNameEqual(DSSASN1UtilsX500PrincipalAreEquals(
						certificate.IssuerX500Principal(), sigCertIdentifier.IssuerName()))
				} else {
					certificateValidity.SetSerialNumberEqual(
						cmsCertificateSourceSerialsEqual(currentSignerIdentifier.SerialNumber(), sigCertIdentifier.SerialNumber()))
					certificateValidity.SetDistinguishedNameEqual(DSSASN1UtilsX500PrincipalAreEquals(
						currentSignerIdentifier.IssuerName(), sigCertIdentifier.IssuerName()))
				}
				certificateValidity.SetSignerIdMatch(currentSignerIdentifier.IsEquivalent(sigCertIdentifier))
			}
		}

		candidates.Add(certificateValidity)
		// The candidate was just added, so the "The certificate validity must be a part of the
		// candidates list!" failure setTheCertificateValidity guards against cannot happen.
		_ = candidates.SetTheCertificateValidity(certificateValidity)

	} else if signingCertificateSource != nil {
		allSignatureCertificates := signingCertificateSource.Certificates()
		// Upstream logs "No signing certificate reference found. Resolve all {} certificates
		// from the provided certificate source as signing candidates.".
		for _, certCandidate := range allSignatureCertificates {
			candidates.Add(NewCertificateValidity(certCandidate))
		}
	}

	return candidates
}

// cmsCertificateSourceFirstToken stands in for Java's foundTokens.iterator().next(): the map
// keys (DSS identifiers) are sorted so that the choice is reproducible, where Java's HashSet
// leaves it to the hash order.
func cmsCertificateSourceFirstToken(tokens map[string]*model.CertificateToken) *model.CertificateToken {
	keys := make([]string, 0, len(tokens))
	for key := range tokens {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return tokens[keys[0]]
}

// cmsCertificateSourceSerialsEqual reproduces BigInteger#equals(Object), which answers false
// for a null argument and never raises.
func cmsCertificateSourceSerialsEqual(first, second *big.Int) bool {
	if first == nil || second == nil {
		return false
	}
	return first.Cmp(second) == 0
}

// compile-time assertion: a CMSCertificateSource drives its SignatureCertificateSource base.
var _ SignatureCertificateSourceOverrides = (*CMSCertificateSource)(nil)
