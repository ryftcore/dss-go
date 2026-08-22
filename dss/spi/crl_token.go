// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/crl/CRLToken.java (DSS 6.5.RC1).
package spi

import (
	"io"
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// CRLToken represents a CRL and provides the information about its validity.
type CRLToken struct {
	RevocationTokenBase[revocation.CRL]

	// crlValidity is the reference to the related CRLValidity.
	crlValidity *crlparser.CRLValidity
}

// NewCRLToken builds the token of the certificate managed by the CRL described by the given
// validity. Port of the CRLToken(CertificateToken, CRLValidity) constructor.
//
// Panics with the Java message when the validity is missing (Objects.requireNonNull). The
// DSSExceptions the revocation status extraction raises - a CRL signed by another issuer
// than the certificate's - are returned as errors.
func NewCRLToken(certificateToken *model.CertificateToken, crlValidity *crlparser.CRLValidity) (*CRLToken, error) {
	if crlValidity == nil {
		panic("CRL Validity cannot be null")
	}
	token := &CRLToken{
		RevocationTokenBase: NewRevocationTokenBase[revocation.CRL](),
		crlValidity:         crlValidity,
	}
	// Registers the token with RevocationTokenBase and, through it, with model.TokenBase,
	// so that the identifier and signature dispatch reach this type.
	token.InitRevocationToken(token)
	token.SetRelatedCertificate(certificateToken)
	token.initInfo()
	if err := token.setRevocationStatus(certificateToken); err != nil {
		return nil, err
	}
	// Upstream logs "A CRLToken created with Id : [{}]" at debug level.
	return token, nil
}

// crlTokenTimeValue dereferences a nullable *time.Time the way crlparser.CRLValidity exposes
// its Date fields, defaulting to the zero time (RevocationTokenBase's null convention) when
// absent.
func crlTokenTimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// initInfo ports the private initInfo().
func (t *CRLToken) initInfo() {
	t.SetSignatureAlgorithm(t.crlValidity.SignatureAlgorithm())
	t.SetThisUpdate(crlTokenTimeValue(t.crlValidity.ThisUpdate()))
	// The dates are equal in case of a CRL.
	t.SetProductionDate(crlTokenTimeValue(t.crlValidity.ThisUpdate()))
	t.SetNextUpdate(crlTokenTimeValue(t.crlValidity.NextUpdate()))
	t.SetCRLNumber(t.crlValidity.CRLNumber())
	t.SetExpiredCertsOnCRL(crlTokenTimeValue(t.crlValidity.ExpiredCertsOnCRL()))

	if issuerToken := t.crlValidity.IssuerToken(); issuerToken != nil {
		t.SetPublicKeyOfTheSigner(issuerToken.PublicKey())
	}

	signatureIntact := t.crlValidity.IsSignatureIntact()
	t.SetSignatureValidity(enumerations.SignatureValidityGet(&signatureIntact))
	t.SetInvalidityReason(t.crlValidity.SignatureInvalidityReason())
	// Upstream logs "{} -> More details in trace." when the signature is not intact.
}

// setRevocationStatus ports the private setRevocationStatus(CertificateToken).
func (t *CRLToken) setRevocationStatus(certificateToken *model.CertificateToken) error {
	issuerToken := certificateToken.IssuerX500Principal()
	var crlSignerSubject *model.X500Principal
	if crlSigner := t.crlValidity.IssuerToken(); crlSigner != nil {
		crlSignerSubject = crlSigner.Subject().Principal()
	}

	if !DSSASN1UtilsX500PrincipalAreEquals(issuerToken, crlSignerSubject) {
		if !t.crlValidity.IsSignatureIntact() {
			return model.NewDSSError(t.crlValidity.SignatureInvalidityReason())
		}
		return model.NewDSSError("The CRLToken is not signed by the same issuer as the CertificateToken to be verified!")
	}

	crlEntry := crlparser.CRLUtilsRevocationInfo(t.crlValidity, certificateToken.SerialNumber())
	if crlEntry != nil {
		t.SetStatus(enumerations.CertificateStatus_REVOKED)
		t.SetRevocationDate(crlEntry.RevocationDate())
		if revocationReason := crlEntry.RevocationReason(); revocationReason != nil {
			// java.security.cert.CRLReason#ordinal() is the RFC 5280 reasonCode value.
			t.SetReason(enumerations.RevocationReasonFromInt(*revocationReason))
		}
	} else {
		t.SetStatus(enumerations.CertificateStatus_GOOD)
	}
	return nil
}

// CheckIsSignedBy is not supported for a CRL: the signature is verified once, by the CRL
// parser, and its outcome is carried by the CRLValidity.
// Port of the checkIsSignedBy(PublicKey) override, which throws an
// UnsupportedOperationException carrying the class name.
func (t *CRLToken) CheckIsSignedBy(publicKey *model.PublicKey) enumerations.SignatureValidity {
	panic("UnsupportedOperationException : eu.europa.esig.dss.spi.x509.revocation.crl.CRLToken")
}

// CertificateSource returns nil: a CRL embeds no certificate.
// Port of the getCertificateSource() override.
func (t *CRLToken) CertificateSource() RevocationCertificateSource {
	return nil
}

// Certificates returns an empty slice: a CRL embeds no certificate.
// Port of the getCertificates() override.
func (t *CRLToken) Certificates() []*model.CertificateToken {
	return nil
}

// CrlValidity returns the validity information of the CRL. Port of getCrlValidity().
func (t *CRLToken) CrlValidity() *crlparser.CRLValidity {
	return t.crlValidity
}

// IssuerX500Principal returns the subject of the CRL signer, nil when the signature is
// invalid and the issuer could therefore not be established.
// Port of the getIssuerX500Principal() override.
func (t *CRLToken) IssuerX500Principal() *model.X500Principal {
	if issuerToken := t.crlValidity.IssuerToken(); issuerToken != nil {
		return issuerToken.Subject().Principal()
	}
	return nil
}

// IssuerCertificateToken returns the certificate of the CRL signer.
// Port of the getIssuerCertificateToken() override.
func (t *CRLToken) IssuerCertificateToken() *model.CertificateToken {
	return t.crlValidity.IssuerToken()
}

// Encoded returns the DER encoding of the CRL. Port of the getEncoded() override.
func (t *CRLToken) Encoded() []byte {
	return t.crlValidity.DerEncoded()
}

// CRLStream opens a reader over the CRL binaries. Port of getCRLStream(); Java's
// InputStream becomes an io.Reader over the retained DER.
func (t *CRLToken) CRLStream() io.Reader {
	return t.crlValidity.ToCRLInputStream()
}

// IsValid reports whether the token signature is intact and the signing certificate has the
// cRLSign key usage bit set. Port of the isValid() override.
func (t *CRLToken) IsValid() bool {
	return t.crlValidity.IsValid()
}

// RevocationType returns CRL. Port of the getRevocationType() override.
func (t *CRLToken) RevocationType() enumerations.RevocationType {
	return enumerations.RevocationType_CRL
}

// Abbreviation returns the DSS abbreviation of the token, used for debugging.
// Port of the getAbbreviation() override.
func (t *CRLToken) Abbreviation() string {
	productionDate := "?"
	if !t.ProductionDate().IsZero() {
		productionDate = DSSUtilsFormatDateToRFC(t.ProductionDate())
	}
	return "CRLToken[" + productionDate + ", signedBy=" + crlTokenPrincipalString(t.IssuerX500Principal()) + "]"
}

// ToString renders the token with the given indentation.
// Port of the toString(String) override.
func (t *CRLToken) ToString(indentStr string) string {
	var out strings.Builder
	out.WriteString(indentStr)
	out.WriteString("CRLToken[\n")
	indentStr += "\t"
	out.WriteString(indentStr)
	out.WriteString("Id: ")
	out.WriteString(t.DSSIDAsString())
	out.WriteString("\n")
	out.WriteString(indentStr)
	out.WriteString("Production time: ")
	out.WriteString(crlTokenDateString(t.ProductionDate()))
	out.WriteString("\n")
	out.WriteString(indentStr)
	out.WriteString("NextUpdate time: ")
	out.WriteString(crlTokenDateString(t.NextUpdate()))
	out.WriteString("\n")
	out.WriteString(indentStr)
	out.WriteString("Signature algorithm: ")
	if t.SignatureAlgorithm() == "" {
		out.WriteString("?")
	} else {
		out.WriteString(string(t.SignatureAlgorithm()))
	}
	out.WriteString("\n")
	out.WriteString(indentStr)
	out.WriteString("Status: ")
	out.WriteString(string(t.Status()))
	out.WriteString("\n")
	out.WriteString(indentStr)
	out.WriteString("Issuer's certificate: ")
	out.WriteString(crlTokenPrincipalString(t.IssuerX500Principal()))
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
func (t *CRLToken) String() string {
	return t.ToString("")
}

// crlTokenDateString renders a date the way upstream's toString() does, i.e. "?" for an
// absent one.
func crlTokenDateString(date time.Time) string {
	if date.IsZero() {
		return "?"
	}
	return DSSUtilsFormatDateToRFC(date)
}

// crlTokenPrincipalString renders a principal the way Java's string concatenation does.
func crlTokenPrincipalString(principal *model.X500Principal) string {
	if principal == nil {
		return "null"
	}
	return principal.String()
}

// compile-time assertion: a CRLToken is a CRL revocation token.
var _ RevocationToken[revocation.CRL] = (*CRLToken)(nil)
