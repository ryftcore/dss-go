// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OCSPToken.java (DSS 6.5.RC1).
//
// The token encapsulates the BasicOCSPResp of DSSRevocationUtils, which keeps the bytes the
// response was decoded from; the token identifier therefore digests exactly the binaries the
// signature covers, not a re-encoding.
//
// BouncyCastle's BasicOCSPResp#isSignatureValid(ContentVerifierProvider) has no Go
// equivalent - upstream reaches it through DSSContentVerifierProviderSecurityFactory - so
// the verification of the response signature over tbsResponseData is implemented here with
// crypto/x509.
package spi

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ocspTokenOIDIsisMttAtCertHash is
// id-isismtt-at-certHash OBJECT IDENTIFIER ::= { id-isismtt-at 13 }, i.e. 1.3.36.8.3.13,
// the Common PKI private OCSP extension carrying the digest of the covered certificate.
// Port of ISISMTTObjectIdentifiers.id_isismtt_at_certHash.
var ocspTokenOIDIsisMttAtCertHash = asn1.ObjectIdentifier{1, 3, 36, 8, 3, 13}

// ocspTokenSignatureAlgorithms maps the algorithms an OCSP response can be signed with onto
// their crypto/x509 counterparts. The algorithms crypto/x509 cannot verify (SHA-224, SHA-3,
// RIPEMD-160, PLAIN-ECDSA and Ed448 flavours) are deliberately absent: they yield a
// verification error rather than a silent success.
var ocspTokenSignatureAlgorithms = map[enumerations.SignatureAlgorithm]x509.SignatureAlgorithm{
	enumerations.SignatureAlgorithmRSAMD5:              x509.MD5WithRSA,
	enumerations.SignatureAlgorithmRSASHA1:             x509.SHA1WithRSA,
	enumerations.SignatureAlgorithmRSASHA256:           x509.SHA256WithRSA,
	enumerations.SignatureAlgorithmRSASHA384:           x509.SHA384WithRSA,
	enumerations.SignatureAlgorithmRSASHA512:           x509.SHA512WithRSA,
	enumerations.SignatureAlgorithmRSASSAPSSSHA256MGF1: x509.SHA256WithRSAPSS,
	enumerations.SignatureAlgorithmRSASSAPSSSHA384MGF1: x509.SHA384WithRSAPSS,
	enumerations.SignatureAlgorithmRSASSAPSSSHA512MGF1: x509.SHA512WithRSAPSS,
	enumerations.SignatureAlgorithmECDSASHA1:           x509.ECDSAWithSHA1,
	enumerations.SignatureAlgorithmECDSASHA256:         x509.ECDSAWithSHA256,
	enumerations.SignatureAlgorithmECDSASHA384:         x509.ECDSAWithSHA384,
	enumerations.SignatureAlgorithmECDSASHA512:         x509.ECDSAWithSHA512,
	enumerations.SignatureAlgorithmDSASHA1:             x509.DSAWithSHA1,
	enumerations.SignatureAlgorithmDSASHA256:           x509.DSAWithSHA256,
	enumerations.SignatureAlgorithmED25519:             x509.PureEd25519,
}

// OCSPToken is the OCSP signed token encapsulating a BasicOCSPResp.
type OCSPToken struct {
	RevocationTokenBase[revocation.OCSP]

	// basicOCSPResp is the encapsulated basic OCSP response.
	basicOCSPResp *BasicOCSPResp
	// latestSingleResp is the single response used with the current certificate, nil when
	// the response carries none for it.
	latestSingleResp *SingleResp
	// issuerCertificateToken is the issuer of the OCSP token.
	issuerCertificateToken *model.CertificateToken
	// certificateSource is the source of the certificates embedded in the OCSP token.
	certificateSource *OCSPCertificateSource
}

// NewOCSPToken builds the token of the certificate the given single response is about.
// Port of the OCSPToken(BasicOCSPResp, SingleResp, CertificateToken, CertificateToken)
// constructor.
//
// Panics with the Java messages when the response or the related certificate is missing
// (Objects.requireNonNull). The certificate extraction of the embedded certificate source,
// which upstream lets fail with a DSSException, is returned as an error.
func NewOCSPToken(basicOCSPResp *BasicOCSPResp, latestSingleResp *SingleResp,
	certificate *model.CertificateToken, issuer *model.CertificateToken) (*OCSPToken, error) {
	if basicOCSPResp == nil {
		panic("The OCSP Response must be defined!")
	}
	if certificate == nil {
		panic("The related certificate token cannot be null!")
	}
	token := &OCSPToken{
		RevocationTokenBase: NewRevocationTokenBase[revocation.OCSP](),
		basicOCSPResp:       basicOCSPResp,
		latestSingleResp:    latestSingleResp,
	}
	// Registers the token with RevocationTokenBase and, through it, with model.TokenBase,
	// so that the identifier and signature dispatch reach this type.
	token.InitRevocationToken(token)
	token.SetProductionDate(basicOCSPResp.ProducedAt())
	token.SetRelatedCertificate(certificate)

	if latestSingleResp != nil {
		token.SetThisUpdate(latestSingleResp.ThisUpdate())
		if nextUpdate := latestSingleResp.NextUpdate(); nextUpdate != nil {
			token.SetNextUpdate(*nextUpdate)
		}
		token.extractStatusInfo(latestSingleResp)
		token.extractArchiveCutOff(latestSingleResp)
		token.extractCertHashExtension(latestSingleResp)
	}

	if err := token.checkSignatureValidity(issuer); err != nil {
		return nil, err
	}
	// Upstream logs "OCSPToken created : {})" at debug level.
	return token, nil
}

// extractStatusInfo ports the private extractStatusInfo(SingleResp).
func (t *OCSPToken) extractStatusInfo(bestSingleResp *SingleResp) {
	certStatus := bestSingleResp.CertStatus()
	switch {
	case certStatus.IsGood():
		// Upstream logs "OCSP status is good".
		t.SetStatus(enumerations.CertificateStatusGood)
	case certStatus.IsRevoked():
		// Upstream logs "OCSP status revoked".
		t.SetStatus(enumerations.CertificateStatusRevoked)
		t.SetRevocationDate(certStatus.RevocationTime)
		reasonId := 0 // unspecified
		if certStatus.HasRevocationReason() {
			reasonId = *certStatus.RevocationReason
		}
		t.SetReason(enumerations.RevocationReasonFromInt(reasonId))
	case certStatus.IsUnknown():
		// Upstream logs "OCSP status unknown".
		t.SetStatus(enumerations.CertificateStatusUnknown)
	}
	// Upstream logs "OCSP certificate status: {}" for any other status; the RFC 6960
	// CertStatus CHOICE has no other alternative.
}

// extractArchiveCutOff ports the private extractArchiveCutOff(SingleResp), which logs and
// ignores an id_pkix_ocsp_archive_cutoff extension it cannot read.
func (t *OCSPToken) extractArchiveCutOff(bestSingleResp *SingleResp) {
	extension := bestSingleResp.Extension(OCSPObjectIdentifierIDPkixOcspArchiveCutoff)
	if extension == nil {
		return
	}
	archiveCutOff, err := DSSASN1UtilsToDate(extension.ParsedValue())
	if err != nil {
		// Upstream logs "Unable to extract id_pkix_ocsp_archive_cutoff : {}".
		return
	}
	t.SetArchiveCutOff(archiveCutOff)
}

// extractCertHashExtension extracts the CertHash extension when present.
//
// Common PKI Part 4: Operational Protocols, 3.1.2 Common PKI Private OCSP Extensions
//
//	CertHash ::= SEQUENCE {
//	    hashAlgorithm   AlgorithmIdentifier,
//	    certificateHash OCTET STRING }
//
// Port of the private extractCertHashExtension(SingleResp), which logs and ignores an
// extension it cannot read.
func (t *OCSPToken) extractCertHashExtension(bestSingleResp *SingleResp) {
	extension := bestSingleResp.Extension(ocspTokenOIDIsisMttAtCertHash)
	if extension == nil {
		return
	}
	asn1CertHash, err := ocspTokenParseCertHash(extension.ParsedValue())
	if err != nil {
		// Upstream logs "Unable to extract id_isismtt_at_certHash : {}".
		return
	}
	digestAlgo, err := enumerations.DigestAlgorithmForOID(asn1CertHash.hashAlgorithm.Algorithm.String())
	if err != nil {
		return
	}
	certHash := model.NewDigest(digestAlgo, asn1CertHash.certificateHash)

	expectedDigest, err := t.RelatedCertificate().Digest(certHash.Algorithm())
	if err != nil {
		return
	}
	t.SetCertHashPresent(true)
	t.SetCertHashMatch(bytes.Equal(expectedDigest, certHash.Value()))
}

// checkSignatureValidity ports the private checkSignatureValidity(CertificateToken).
func (t *OCSPToken) checkSignatureValidity(caCertificateToken *model.CertificateToken) error {
	certificateSource, err := t.certificateSourceOrError()
	if err != nil {
		return err
	}
	candidates := certificateSource.CandidatesForSigningCertificate(caCertificateToken)

	signingCertificateValidator := NewOCSPSignatureIntegrityValidator(t)
	certificateValidity := signingCertificateValidator.Validate(candidates)
	if certificateValidity != nil {
		candidates.SetTheCertificateValidity(certificateValidity)
		t.issuerCertificateToken = certificateValidity.CertificateToken()
	}
	return nil
}

// SignatureAlgorithm returns the algorithm the response was signed with, resolving it from
// the response's AlgorithmIdentifier on first use. Port of the getSignatureAlgorithm()
// override.
//
// DEVIATION: SignatureAlgorithm.forOidAndParams throws an IllegalArgumentException for an
// unknown OID; the Token contract has no error channel, so an unresolvable algorithm leaves
// the field unset, which the token renders as "?" exactly as an absent one.
func (t *OCSPToken) SignatureAlgorithm() enumerations.SignatureAlgorithm {
	if t.RevocationTokenBase.SignatureAlgorithm() == "" {
		signatureAlgorithmID := t.basicOCSPResp.SignatureAlgorithmID()
		signatureAlgorithm, err := enumerations.SignatureAlgorithmForOIDAndParams(
			signatureAlgorithmID.Algorithm.String(), signatureAlgorithmID.Parameters)
		if err == nil {
			t.SetSignatureAlgorithm(signatureAlgorithm)
		}
	}
	return t.RevocationTokenBase.SignatureAlgorithm()
}

// BasicOCSPResp returns the encapsulated basic OCSP response. Port of getBasicOCSPResp().
func (t *OCSPToken) BasicOCSPResp() *BasicOCSPResp {
	return t.basicOCSPResp
}

// LatestSingleResp returns the single response used with the related certificate, nil when
// the response carries none for it. Port of getLatestSingleResp().
func (t *OCSPToken) LatestSingleResp() *SingleResp {
	return t.latestSingleResp
}

// CertificateSource returns the source of the certificates embedded in the OCSP token.
// Port of the getCertificateSource() override, whose covariant OCSPCertificateSource return
// type Go cannot express; callers needing it assert on the result.
func (t *OCSPToken) CertificateSource() RevocationCertificateSource {
	certificateSource, err := t.certificateSourceOrError()
	if err != nil {
		return nil
	}
	return certificateSource
}

// certificateSourceOrError builds the embedded certificate source on first use. Upstream the
// OCSPCertificateSource constructor throws; the Go one returns the failure, and the token
// constructor is what surfaces it.
func (t *OCSPToken) certificateSourceOrError() (*OCSPCertificateSource, error) {
	if t.certificateSource == nil {
		certificateSource, err := NewOCSPCertificateSource(t.basicOCSPResp)
		if err != nil {
			return nil, err
		}
		t.certificateSource = certificateSource
	}
	return t.certificateSource, nil
}

// Certificates returns the certificates embedded in the OCSP token.
// Port of the getCertificates() override.
func (t *OCSPToken) Certificates() []*model.CertificateToken {
	certificateSource, err := t.certificateSourceOrError()
	if err != nil {
		return nil
	}
	return certificateSource.Certificates()
}

// Encoded returns the encoding of the OCSPResponse wrapping the basic response.
// Port of the getEncoded() override.
//
// The DSSException upstream raises for an empty response cannot occur - the constructor
// rejects a missing response - and is therefore a panic here.
func (t *OCSPToken) Encoded() []byte {
	encoded, err := DSSRevocationUtilsEncodedFromBasicResp(t.basicOCSPResp)
	if err != nil {
		panic("Empty OCSP response")
	}
	return encoded
}

// IssuerX500Principal returns the subject of the OCSP token issuer, nil when the signer
// could not be established. Port of the getIssuerX500Principal() override.
func (t *OCSPToken) IssuerX500Principal() *model.X500Principal {
	if t.issuerCertificateToken != nil {
		return t.issuerCertificateToken.Subject().Principal()
	}
	return nil
}

// IssuerCertificateToken returns the certificate of the OCSP token issuer.
// Port of the getIssuerCertificateToken() override.
func (t *OCSPToken) IssuerCertificateToken() *model.CertificateToken {
	return t.issuerCertificateToken
}

// IsValid reports whether the OCSP token is valid. IsSignedBy must have been called first.
// Port of the isValid() override.
func (t *OCSPToken) IsValid() bool {
	return t.IsSignatureIntact() && t.isOCSPVersionValid()
}

// CheckIsSignedBy verifies whether the OCSP token has been signed with the given public key.
// Port of the protected checkIsSignedBy(PublicKey) override.
//
// DEVIATION: Java separates a verification that returns false (INVALID, no reason recorded)
// from one that throws (INVALID, reason recorded). crypto/x509 reports both as an error, so
// every failed verification records an invalidity reason.
func (t *OCSPToken) CheckIsSignedBy(publicKey *model.PublicKey) enumerations.SignatureValidity {
	t.SetInvalidityReason("")
	if err := ocspTokenVerifySignature(t.basicOCSPResp, publicKey); err != nil {
		// Upstream logs "An error occurred during in attempt to check signature owner : ".
		t.SetInvalidityReason(fmt.Sprintf("%T - %s", err, err.Error()))
		t.SetSignatureValidity(enumerations.SignatureValidityInvalid)
		return t.SignatureValidity()
	}
	valid := true
	t.SetSignatureValidity(enumerations.SignatureValidityGet(&valid))
	return t.SignatureValidity()
}

// OCSPTokenVersion returns the version defined within the OCSP token, i.e. the encoded value
// plus one ('1' for the default 'v1'). Port of getOCSPTokenVersion().
func (t *OCSPToken) OCSPTokenVersion() int {
	return t.basicOCSPResp.Version()
}

// isOCSPVersionValid ports the private isOCSPVersionValid(), which records an invalidity
// reason when the response syntax version is not v1 (RFC 6960).
func (t *OCSPToken) isOCSPVersionValid() bool {
	versionValid := t.OCSPTokenVersion() == 1
	if !versionValid && utils.IsStringEmpty(t.InvalidityReason()) {
		t.SetInvalidityReason("Basic OCSP Response version is invalid (shall be v1)!")
	}
	return versionValid
}

// RevocationType returns OCSP. Port of the getRevocationType() override.
func (t *OCSPToken) RevocationType() enumerations.RevocationType {
	return enumerations.RevocationTypeOCSP
}

// Abbreviation returns the DSS abbreviation of the token, used for debugging.
// Port of the getAbbreviation() override.
func (t *OCSPToken) Abbreviation() string {
	producedAt := "?"
	if t.basicOCSPResp != nil {
		producedAt = DSSUtilsFormatDateToRFC(t.basicOCSPResp.ProducedAt())
	}
	return "OCSPToken[" + producedAt + ", signedBy=" + ocspTokenPrincipalString(t.IssuerX500Principal()) + "]"
}

// ToString renders the token with the given indentation.
// Port of the toString(String) override.
func (t *OCSPToken) ToString(indentStr string) string {
	var out strings.Builder
	out.WriteString(indentStr)
	out.WriteString("OCSPToken[\n")
	indentStr += "\t"
	out.WriteString(indentStr)
	out.WriteString("Id: ")
	out.WriteString(t.DSSIDAsString())
	out.WriteString("\n")
	out.WriteString(indentStr)
	out.WriteString("ProductionTime: ")
	out.WriteString(DSSUtilsFormatDateToRFC(t.ProductionDate()))
	out.WriteString("; ")
	out.WriteString(indentStr)
	out.WriteString("ThisUpdate: ")
	out.WriteString(DSSUtilsFormatDateToRFC(t.ThisUpdate()))
	out.WriteString("; ")
	out.WriteString(indentStr)
	out.WriteString("NextUpdate: ")
	out.WriteString(DSSUtilsFormatDateToRFC(t.NextUpdate()))
	out.WriteString("\n")
	if t.IssuerX500Principal() != nil {
		out.WriteString(indentStr)
		out.WriteString("SignedBy: ")
		out.WriteString(t.IssuerX500Principal().String())
		out.WriteString("\n")
	}
	out.WriteString(indentStr)
	out.WriteString("Signature algorithm: ")
	// Upstream reads the field, not the lazy getter, so an unresolved algorithm prints "?".
	if signatureAlgorithm := t.RevocationTokenBase.SignatureAlgorithm(); signatureAlgorithm == "" {
		out.WriteString("?")
	} else {
		out.WriteString(signatureAlgorithm.JCEID())
	}
	out.WriteString("\n")
	if t.RelatedCertificateID() != "" {
		out.WriteString(indentStr)
		out.WriteString("Related certificate: ")
		out.WriteString(t.RelatedCertificateID())
		out.WriteString("\n")
	}
	indentStr = indentStr[1:]
	out.WriteString(indentStr)
	out.WriteString("]")
	return out.String()
}

// String returns ToString(""). Port of Token#toString().
func (t *OCSPToken) String() string {
	return t.ToString("")
}

// ocspTokenCertHash is the Common PKI CertHash structure, replacing
// org.bouncycastle.asn1.isismtt.ocsp.CertHash.
type ocspTokenCertHash struct {
	// hashAlgorithm is the algorithm the certificate digest was computed with.
	hashAlgorithm *AlgorithmIdentifier
	// certificateHash is the digest of the covered certificate.
	certificateHash []byte
}

// ocspTokenParseCertHash decodes a CertHash from its DER encoding.
// Port of CertHash.getInstance(Object).
func ocspTokenParseCertHash(der []byte) (*ocspTokenCertHash, error) {
	element, rest, err := asn1ber.Parse(der)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the CertHash")
	}
	if !element.IsConstructed() || len(element.Children()) != 2 {
		return nil, errors.New("malformed CertHash")
	}
	hashAlgorithm, err := asn1ber.AlgorithmIdentifierFromElement(element.Children()[0])
	if err != nil {
		return nil, err
	}
	if !element.Children()[1].IsUniversal(asn1ber.TagOctetString) {
		return nil, errors.New("malformed CertHash: certificateHash is not an OCTET STRING")
	}
	return &ocspTokenCertHash{hashAlgorithm: hashAlgorithm, certificateHash: element.Children()[1].Octets()}, nil
}

// ocspTokenVerifySignature verifies the response signature over the retained tbsResponseData
// bytes, standing in for BasicOCSPResp#isSignatureValid(ContentVerifierProvider).
func ocspTokenVerifySignature(basicOCSPResp *BasicOCSPResp, publicKey *model.PublicKey) error {
	if publicKey == nil {
		return errors.New("NullPointerException : public key is missing")
	}
	key := publicKey.Key()
	if key == nil {
		return errors.New("InvalidKeyException : the public key has not been parsed")
	}
	signatureAlgorithmID := basicOCSPResp.SignatureAlgorithmID()
	signatureAlgorithm, err := enumerations.SignatureAlgorithmForOIDAndParams(
		signatureAlgorithmID.Algorithm.String(), signatureAlgorithmID.Parameters)
	if err != nil {
		return err
	}
	algorithm, supported := ocspTokenSignatureAlgorithms[signatureAlgorithm]
	if !supported {
		return fmt.Errorf("NoSuchAlgorithmException : %s cannot be verified", signatureAlgorithm)
	}
	// CheckSignature verifies "signed" against the receiver's public key, so the signer's
	// key is carried by a throwaway certificate value.
	signer := &x509.Certificate{PublicKey: key}
	return signer.CheckSignature(algorithm, basicOCSPResp.TBSResponseData(), basicOCSPResp.Signature())
}

// ocspTokenPrincipalString renders a principal the way Java's string concatenation does.
func ocspTokenPrincipalString(principal *model.X500Principal) string {
	if principal == nil {
		return "null"
	}
	return principal.String()
}

// compile-time assertion: an OCSPToken is an OCSP revocation token.
var _ RevocationToken[revocation.OCSP] = (*OCSPToken)(nil)
