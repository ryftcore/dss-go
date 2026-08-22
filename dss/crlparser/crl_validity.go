// Ported from dss-crl-parser/src/main/java/eu/europa/esig/dss/crl/CRLValidity.java (DSS 6.5.RC1).
package crlparser

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"math/big"
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// CRLValidity encapsulates all information related to the validity of a CRL. Use IsValid to
// check the validity.
//
// DEVIATION: upstream splits the java.security.cert.X509CRL caching into a subclass,
// X509CRLValidity (dss-crl-parser-x509crl), reached through the ICRLUtils ServiceLoader
// indirection. This port has a single native implementation, so the extra field is folded
// directly into this struct; its accessors live in x509_crl_validity.go, mirroring the
// upstream file that introduced them.
type CRLValidity struct {
	// crlBinary incorporates the CRL binaries.
	crlBinary *CRLBinary

	// url is the distributionPoint [0] DistributionPointName OPTIONAL.
	url string
	// onlyUserCerts is onlyContainsUserCerts [1] BOOLEAN DEFAULT FALSE.
	onlyUserCerts bool
	// onlyCaCerts is onlyContainsCACerts [2] BOOLEAN DEFAULT FALSE.
	onlyCaCerts bool
	// onlySomeReasonFlags is onlySomeReasons [3] ReasonFlags OPTIONAL, kept as the raw BIT
	// STRING (upstream keeps a BouncyCastle ReasonFlags, whose equality and only consumer -
	// IsUnknownCriticalExtension - only cares whether the field is present).
	onlySomeReasonFlags *asn1.BitString
	// indirectCrl is indirectCRL [4] BOOLEAN DEFAULT FALSE.
	indirectCrl bool
	// onlyAttributeCerts is onlyContainsAttributeCerts [5] BOOLEAN DEFAULT FALSE.
	onlyAttributeCerts bool

	// crlSignKeyUsage tells whether the signing certificate contains the 'cRLSign' key usage.
	crlSignKeyUsage bool
	// issuerX509PrincipalMatches tells whether the X500 Principal defined in the CRL matches
	// the value of its issuer certificate.
	issuerX509PrincipalMatches bool
	// signatureIntact tells whether the signature is valid.
	signatureIntact bool
	// signatureInvalidityReason carries the signature invalidity reason when the signature is
	// invalid, empty otherwise.
	signatureInvalidityReason string
	// signatureAlgorithm is the SignatureAlgorithm used for the signature.
	signatureAlgorithm enumerations.SignatureAlgorithm
	// issuerToken is the issuer certificate.
	issuerToken *model.CertificateToken

	// criticalExtensionsOid is the collection of critical extension OIDs.
	criticalExtensionsOid []string

	// expiredCertsOnCRL is the 'expiredCertsOnCRL' date value.
	expiredCertsOnCRL *time.Time
	// nextUpdate is the 'nextUpdate' date value.
	nextUpdate *time.Time
	// thisUpdate is the 'thisUpdate' date value.
	thisUpdate *time.Time

	// crlNumber is the CRL Number extension value.
	crlNumber *big.Int

	// x509CRL is the parsed CRL, folded in from X509CRLValidity (see x509_crl_validity.go).
	x509CRL *x509.RevocationList
}

// NewCRLValidity builds a CRLValidity over the given CRLBinary.
// Port of the CRLValidity(CRLBinary) constructor.
//
// Panics with the Java message when crlBinary is missing (Objects.requireNonNull).
func NewCRLValidity(crlBinary *CRLBinary) *CRLValidity {
	if crlBinary == nil {
		panic("CRLBinary cannot be null!")
	}
	return &CRLValidity{crlBinary: crlBinary}
}

// CrlBinary returns the binary of the CRL. Port of getCrlBinary().
func (v *CRLValidity) CrlBinary() *CRLBinary {
	return v.crlBinary
}

// DerEncoded returns the DER encoded binaries of the CRL. Port of getDerEncoded().
func (v *CRLValidity) DerEncoded() []byte {
	return v.crlBinary.Binaries()
}

// ToCRLInputStream opens a reader over the CRL's binaries. Port of toCRLInputStream().
func (v *CRLValidity) ToCRLInputStream() *bytes.Reader {
	return bytes.NewReader(v.DerEncoded())
}

// SignatureAlgorithm returns the used SignatureAlgorithm. Port of getSignatureAlgorithm().
func (v *CRLValidity) SignatureAlgorithm() enumerations.SignatureAlgorithm {
	return v.signatureAlgorithm
}

// SetSignatureAlgorithm sets the used SignatureAlgorithm. Port of setSignatureAlgorithm(SignatureAlgorithm).
func (v *CRLValidity) SetSignatureAlgorithm(signatureAlgorithm enumerations.SignatureAlgorithm) {
	v.signatureAlgorithm = signatureAlgorithm
}

// NextUpdate returns the 'nextUpdate' field date, nil when unset. Port of getNextUpdate().
func (v *CRLValidity) NextUpdate() *time.Time {
	return v.nextUpdate
}

// SetNextUpdate sets the 'nextUpdate' field date. Port of setNextUpdate(Date).
func (v *CRLValidity) SetNextUpdate(nextUpdate *time.Time) {
	v.nextUpdate = nextUpdate
}

// ThisUpdate returns the 'thisUpdate' field date, nil when unset. Port of getThisUpdate().
func (v *CRLValidity) ThisUpdate() *time.Time {
	return v.thisUpdate
}

// SetThisUpdate sets the 'thisUpdate' field date. Port of setThisUpdate(Date).
func (v *CRLValidity) SetThisUpdate(thisUpdate *time.Time) {
	v.thisUpdate = thisUpdate
}

// CRLNumber returns the CRL Number extension value, nil when absent. Port of getCRLNumber().
func (v *CRLValidity) CRLNumber() *big.Int {
	return v.crlNumber
}

// SetCRLNumber sets the CRL Number extension value. Port of setCRLNumber(BigInteger).
func (v *CRLValidity) SetCRLNumber(crlNumber *big.Int) {
	v.crlNumber = crlNumber
}

// ExpiredCertsOnCRL returns the 'expiredCertsOnCRL' field date, nil when unset.
// Port of getExpiredCertsOnCRL().
func (v *CRLValidity) ExpiredCertsOnCRL() *time.Time {
	return v.expiredCertsOnCRL
}

// SetExpiredCertsOnCRL sets the 'expiredCertsOnCRL' field date. Port of setExpiredCertsOnCRL(Date).
func (v *CRLValidity) SetExpiredCertsOnCRL(expiredCertsOnCRL *time.Time) {
	v.expiredCertsOnCRL = expiredCertsOnCRL
}

// IssuerX509PrincipalMatches returns whether the issuer X509 Principal matches between the one
// defined in the CRL and the corresponding value of its issuer certificate.
// Port of isIssuerX509PrincipalMatches().
func (v *CRLValidity) IssuerX509PrincipalMatches() bool {
	return v.issuerX509PrincipalMatches
}

// SetIssuerX509PrincipalMatches sets whether the issuer X509 Principal matches.
// Port of setIssuerX509PrincipalMatches(boolean).
func (v *CRLValidity) SetIssuerX509PrincipalMatches(issuerX509PrincipalMatches bool) {
	v.issuerX509PrincipalMatches = issuerX509PrincipalMatches
}

// IsSignatureIntact returns whether the signature value is valid. Port of isSignatureIntact().
func (v *CRLValidity) IsSignatureIntact() bool {
	return v.signatureIntact
}

// SetSignatureIntact sets whether the signature value is valid. Port of setSignatureIntact(boolean).
func (v *CRLValidity) SetSignatureIntact(signatureIntact bool) {
	v.signatureIntact = signatureIntact
}

// IsCrlSignKeyUsage returns whether the issuer certificate has the 'cRLSign' key usage.
// Port of isCrlSignKeyUsage().
func (v *CRLValidity) IsCrlSignKeyUsage() bool {
	return v.crlSignKeyUsage
}

// SetCrlSignKeyUsage sets whether the issuer certificate has the 'cRLSign' key usage.
// Port of setCrlSignKeyUsage(boolean).
func (v *CRLValidity) SetCrlSignKeyUsage(crlSignKeyUsage bool) {
	v.crlSignKeyUsage = crlSignKeyUsage
}

// IssuerToken returns the issuer certificate token, nil when unset. Port of getIssuerToken().
func (v *CRLValidity) IssuerToken() *model.CertificateToken {
	return v.issuerToken
}

// SetIssuerToken sets the issuer certificate token. Port of setIssuerToken(CertificateToken).
func (v *CRLValidity) SetIssuerToken(issuerToken *model.CertificateToken) {
	v.issuerToken = issuerToken
}

// SignatureInvalidityReason returns the signature invalidity reason if the signature is
// invalid, empty otherwise. Port of getSignatureInvalidityReason().
func (v *CRLValidity) SignatureInvalidityReason() string {
	return v.signatureInvalidityReason
}

// SetSignatureInvalidityReason sets the signature invalidity reason.
// Port of setSignatureInvalidityReason(String).
func (v *CRLValidity) SetSignatureInvalidityReason(signatureInvalidityReason string) {
	v.signatureInvalidityReason = signatureInvalidityReason
}

// URL returns the distributionPoint url. Port of getUrl().
func (v *CRLValidity) URL() string {
	return v.url
}

// SetURL sets the distributionPoint url. Port of setUrl(String).
func (v *CRLValidity) SetURL(url string) {
	v.url = url
}

// SetOnlyUserCerts sets the 'onlyContainsUserCerts' value. Port of setOnlyUserCerts(boolean).
//
// Like upstream, there is no matching getter: the field only feeds IsUnknownCriticalExtension.
func (v *CRLValidity) SetOnlyUserCerts(onlyUserCerts bool) {
	v.onlyUserCerts = onlyUserCerts
}

// SetOnlyCaCerts sets the 'onlyContainsCACerts' value. Port of setOnlyCaCerts(boolean).
func (v *CRLValidity) SetOnlyCaCerts(onlyCaCerts bool) {
	v.onlyCaCerts = onlyCaCerts
}

// SetReasonFlags sets the 'onlySomeReasons' value. Port of setReasonFlags(ReasonFlags).
func (v *CRLValidity) SetReasonFlags(reasonFlags *asn1.BitString) {
	v.onlySomeReasonFlags = reasonFlags
}

// SetIndirectCrl sets the 'indirectCRL' value. Port of setIndirectCrl(boolean).
func (v *CRLValidity) SetIndirectCrl(indirectCrl bool) {
	v.indirectCrl = indirectCrl
}

// SetOnlyAttributeCerts sets the 'onlyContainsAttributeCerts' value.
// Port of setOnlyAttributeCerts(boolean).
func (v *CRLValidity) SetOnlyAttributeCerts(onlyAttributeCerts bool) {
	v.onlyAttributeCerts = onlyAttributeCerts
}

// AreCriticalExtensionsOidNotEmpty reports whether the collection of critical extension OIDs
// is not empty. Port of areCriticalExtensionsOidNotEmpty().
func (v *CRLValidity) AreCriticalExtensionsOidNotEmpty() bool {
	return len(v.criticalExtensionsOid) > 0
}

// SetCriticalExtensionsOid sets the collection of critical extension OIDs.
// Port of setCriticalExtensionsOid(Collection).
func (v *CRLValidity) SetCriticalExtensionsOid(criticalExtensionsOid []string) {
	v.criticalExtensionsOid = criticalExtensionsOid
}

// IsValid reports whether the CRL is valid. To be valid the CRL must fulfill the following
// requirements: its signature must be valid, the issuer of the certificate for which the CRL
// is used must match the CRL signing certificate, and the mandatory key usage must be present.
// Port of isValid().
func (v *CRLValidity) IsValid() bool {
	return v.issuerX509PrincipalMatches && v.signatureIntact && v.crlSignKeyUsage && !v.IsUnknownCriticalExtension()
}

// IsUnknownCriticalExtension reports whether the critical extensions are unknown.
// Port of isUnknownCriticalExtension().
func (v *CRLValidity) IsUnknownCriticalExtension() bool {
	return v.AreCriticalExtensionsOidNotEmpty() &&
		((v.onlyAttributeCerts && v.onlyCaCerts && v.onlyUserCerts && v.indirectCrl) ||
			v.onlySomeReasonFlags != nil || v.url == "")
}

// Equals reports whether the two CRLValidity values carry the same state.
// Port of equals(Object).
func (v *CRLValidity) Equals(other *CRLValidity) bool {
	if v == other {
		return true
	}
	if other == nil {
		return false
	}
	if v.onlyUserCerts != other.onlyUserCerts || v.onlyCaCerts != other.onlyCaCerts ||
		v.indirectCrl != other.indirectCrl || v.onlyAttributeCerts != other.onlyAttributeCerts ||
		v.crlSignKeyUsage != other.crlSignKeyUsage ||
		v.issuerX509PrincipalMatches != other.issuerX509PrincipalMatches ||
		v.signatureIntact != other.signatureIntact {
		return false
	}
	if !crlValidityCRLBinaryEquals(v.crlBinary, other.crlBinary) {
		return false
	}
	if v.url != other.url {
		return false
	}
	if !crlValidityReasonFlagsEqual(v.onlySomeReasonFlags, other.onlySomeReasonFlags) {
		return false
	}
	if v.signatureInvalidityReason != other.signatureInvalidityReason {
		return false
	}
	if v.signatureAlgorithm != other.signatureAlgorithm {
		return false
	}
	if !crlValidityIssuerTokenEquals(v.issuerToken, other.issuerToken) {
		return false
	}
	if !crlValidityStringsEqual(v.criticalExtensionsOid, other.criticalExtensionsOid) {
		return false
	}
	if !crlValidityTimeEqual(v.expiredCertsOnCRL, other.expiredCertsOnCRL) ||
		!crlValidityTimeEqual(v.nextUpdate, other.nextUpdate) ||
		!crlValidityTimeEqual(v.thisUpdate, other.thisUpdate) {
		return false
	}
	return crlValidityBigIntEqual(v.crlNumber, other.crlNumber)
}

// String returns a compact description of the CRLValidity, mirroring the fields toString()
// interpolates. Port of toString().
func (v *CRLValidity) String() string {
	var out strings.Builder
	out.WriteString("CRLValidity{DSS ID=")
	out.WriteString(v.crlBinary.AsXmlID())
	out.WriteString(", issuerX509PrincipalMatches=")
	crlValidityWriteBool(&out, v.issuerX509PrincipalMatches)
	out.WriteString(", signatureIntact=")
	crlValidityWriteBool(&out, v.signatureIntact)
	out.WriteString(", crlSignKeyUsage=")
	crlValidityWriteBool(&out, v.crlSignKeyUsage)
	out.WriteString(", unknownCriticalExtension=")
	crlValidityWriteBool(&out, v.IsUnknownCriticalExtension())
	out.WriteString(", issuerToken=")
	if v.issuerToken != nil {
		out.WriteString(v.issuerToken.String())
	} else {
		out.WriteString("null")
	}
	out.WriteString(", signatureInvalidityReason='")
	out.WriteString(v.signatureInvalidityReason)
	out.WriteString("'}")
	return out.String()
}

func crlValidityWriteBool(out *strings.Builder, b bool) {
	if b {
		out.WriteString("true")
	} else {
		out.WriteString("false")
	}
}

func crlValidityCRLBinaryEquals(a, b *CRLBinary) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.DSSID().Equals(b.DSSID())
}

func crlValidityIssuerTokenEquals(a, b *model.CertificateToken) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equals(b)
}

func crlValidityReasonFlagsEqual(a, b *asn1.BitString) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.BitLength == b.BitLength && bytes.Equal(a.Bytes, b.Bytes)
}

func crlValidityStringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func crlValidityTimeEqual(a, b *time.Time) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equal(*b)
}

func crlValidityBigIntEqual(a, b *big.Int) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Cmp(b) == 0
}
