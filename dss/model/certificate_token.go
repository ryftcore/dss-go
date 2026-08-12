// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/CertificateToken.java (DSS 6.5.RC1).
package model

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/utain/esig/dss/enumerations"
)

// certificateTokenDateFormat reproduces java.util.Date#toString(), which CertificateToken's
// toString() interpolates for the validity period.
const certificateTokenDateFormat = "Mon Jan 02 15:04:05 MST 2006"

// CertificateToken wraps an X.509 certificate encountered during signature validation,
// caching the frequently used information (identifiers, key usages, self-signed state) that
// the validation process would otherwise recompute.
type CertificateToken struct {
	TokenBase

	// x509Certificate is the encapsulated X509 certificate.
	x509Certificate *x509.Certificate
	// subject is the parsed subject name; the Java certificate caches its X500Principals too.
	subject *X500Principal
	// issuer is the parsed issuer name.
	issuer *X500Principal
	// publicKey pairs the parsed key with the certificate's SubjectPublicKeyInfo bytes.
	publicKey *PublicKey
	// entityKey is the digest of the public key and subject name (cross certificates share it).
	entityKey *EntityIdentifier
	// selfSigned stays nil until the first IsSelfSigned call.
	selfSigned *bool
	// keyUsageBits caches the certificate's key usages.
	keyUsageBits []enumerations.KeyUsageBit
	// sourceURL is the URL the certificate was downloaded from (aia.caIssuers).
	sourceURL string
}

// NewCertificateToken creates a CertificateToken wrapping the provided certificate.
//
// Panics with the Java message when the certificate is missing (Objects.requireNonNull).
// The signature algorithm is resolved from the certificate's signature algorithm OID rather
// than from its Java name, and an unknown OID - Java's IllegalArgumentException - is
// returned as an error.
func NewCertificateToken(x509Certificate *x509.Certificate) (*CertificateToken, error) {
	if x509Certificate == nil {
		panic("X509 certificate is missing")
	}

	subject, err := NewX500Principal(x509Certificate.RawSubject)
	if err != nil {
		return nil, err
	}
	issuer, err := NewX500Principal(x509Certificate.RawIssuer)
	if err != nil {
		return nil, err
	}
	// The SubjectPublicKeyInfo bytes are kept as parsed: the EntityIdentifier digests them,
	// so re-encoding the key could silently change the identifier.
	publicKey := NewPublicKeyFromEncoded(x509Certificate.RawSubjectPublicKeyInfo, x509Certificate.PublicKey)

	sigAlgOID, sigAlgParams, err := certificateTokenSigAlgOIDAndParams(x509Certificate.Raw)
	if err != nil {
		return nil, err
	}
	// The Algorithm OID is used and not the name x509Certificate.getSigAlgName().
	signatureAlgorithm, err := enumerations.SignatureAlgorithmForOIDAndParams(sigAlgOID, sigAlgParams)
	if err != nil {
		return nil, err
	}

	token := &CertificateToken{
		TokenBase:       NewTokenBase(),
		x509Certificate: x509Certificate,
		subject:         subject,
		issuer:          issuer,
		publicKey:       publicKey,
	}
	token.InitToken(token)
	token.SetSignatureAlgorithm(signatureAlgorithm)
	token.entityKey = NewEntityIdentifierBuilder(publicKey, subject).Build()
	return token, nil
}

// Abbreviation returns the DSS Id string of the token. Port of getAbbreviation().
func (c *CertificateToken) Abbreviation() string {
	return c.DSSIDAsString()
}

// EntityKey returns the identifier of the current entity key (public key plus subject name).
// Several certificates can share it (cross-certificates). Port of getEntityKey().
func (c *CertificateToken) EntityKey() *EntityIdentifier {
	return c.entityKey
}

// IssuerEntityKey returns the entity key identifier of the certificate's issuer; a
// self-signed certificate is its own issuer. Port of the getIssuerEntityKey() override.
func (c *CertificateToken) IssuerEntityKey() *EntityIdentifier {
	if c.IsSelfSigned() {
		return NewEntityIdentifierBuilder(c.PublicKey(), c.Subject().Principal()).Build()
	}
	return c.TokenBase.IssuerEntityKey()
}

// PublicKey returns the public key associated with the certificate, paired with the exact
// SubjectPublicKeyInfo bytes it was parsed from. Port of getPublicKey().
func (c *CertificateToken) PublicKey() *PublicKey {
	return c.publicKey
}

// NotAfter returns the expiration date of the certificate. Port of getNotAfter().
func (c *CertificateToken) NotAfter() time.Time {
	return c.x509Certificate.NotAfter
}

// NotBefore returns the issuance date of the certificate. Port of getNotBefore().
func (c *CertificateToken) NotBefore() time.Time {
	return c.x509Certificate.NotBefore
}

// CreationDate returns the certificate's notBefore date. Port of getCreationDate().
func (c *CertificateToken) CreationDate() time.Time {
	return c.NotBefore()
}

// SourceURL returns the certificate's source URL. Port of getSourceURL().
func (c *CertificateToken) SourceURL() string {
	return c.sourceURL
}

// SetSourceURL sets the certificate's source URL. Port of setSourceURL(String).
func (c *CertificateToken) SetSourceURL(sourceURL string) {
	c.sourceURL = sourceURL
}

// IsValidOn reports whether the given date lies in the certificate's validity period, both
// bounds included. Port of isValidOn(Date).
//
// Java answers false for a null date; the Go port maps that to the zero time.Time.
func (c *CertificateToken) IsValidOn(date time.Time) bool {
	if c.x509Certificate == nil || date.IsZero() {
		return false
	}
	return !date.Before(c.x509Certificate.NotBefore) && !date.After(c.x509Certificate.NotAfter)
}

// IsSelfSigned reports whether the certificate is self-signed, i.e. it is self-issued and
// its signature verifies with the public key bound into it [RFC5280]. The answer is computed
// once and cached, and a positive answer marks the token's signature as VALID.
// Port of the isSelfSigned() override.
func (c *CertificateToken) IsSelfSigned() bool {
	if c.selfSigned == nil {
		selfSigned := c.IsSelfIssued()
		if selfSigned {
			if err := certificateTokenVerify(c.x509Certificate, c.PublicKey()); err == nil {
				selfSigned = true
				c.SetSignatureValidity(enumerations.SignatureValidity_VALID)
			} else {
				selfSigned = false
			}
		}
		c.selfSigned = &selfSigned
	} else if *c.selfSigned {
		c.SetSignatureValidity(enumerations.SignatureValidity_VALID)
	}
	return *c.selfSigned
}

// IsSelfIssued reports whether the certificate is self-issued, i.e. its issuer and subject
// are the same entity [RFC5280]. Port of isSelfIssued().
func (c *CertificateToken) IsSelfIssued() bool {
	return bytes.Equal(c.x509Certificate.RawSubject, c.x509Certificate.RawIssuer)
}

// IsEquivalent reports whether the given token has the same public key.
// Port of isEquivalent(CertificateToken).
func (c *CertificateToken) IsEquivalent(token *CertificateToken) bool {
	return c.PublicKey().Equals(token.PublicKey())
}

// Certificate returns the enclosed X.509 certificate. Port of getCertificate().
func (c *CertificateToken) Certificate() *x509.Certificate {
	return c.x509Certificate
}

// Encoded returns the ASN.1 DER encoded form of this certificate. Port of getEncoded().
//
// Unlike the Java method this does not clone, and the CertificateEncodingException branch
// (which raises "Unable to encode the certificate") has no Go counterpart because the DER
// is retained by crypto/x509 rather than re-encoded.
func (c *CertificateToken) Encoded() []byte {
	return c.x509Certificate.Raw
}

// SerialNumber returns the certificate serial number, the integer the CA assigns to each
// certificate it issues. Port of getSerialNumber().
func (c *CertificateToken) SerialNumber() *big.Int {
	return c.x509Certificate.SerialNumber
}

// Subject returns the subject wrapped in an X500PrincipalHelper. Port of getSubject().
func (c *CertificateToken) Subject() *X500PrincipalHelper {
	return NewX500PrincipalHelper(c.subject)
}

// Issuer returns the issuer wrapped in an X500PrincipalHelper. Port of getIssuer().
func (c *CertificateToken) Issuer() *X500PrincipalHelper {
	return NewX500PrincipalHelper(c.issuer)
}

// IssuerX500Principal returns the X500Principal of the certificate that signed this token.
// Port of the getIssuerX500Principal() override.
func (c *CertificateToken) IssuerX500Principal() *X500Principal {
	return c.issuer
}

// CheckIsSignedBy verifies the certificate's signature with the given public key, recording
// the outcome and, on failure, the reason. Port of the protected checkIsSignedBy(PublicKey).
//
// DEVIATION: upstream raises a DSSException ("No provider has been found for signature
// validation : %s") when the JCA has no provider for the algorithm; Go has no provider
// mechanism, so that branch is absent and an unsupported algorithm is reported as an
// invalidity reason like any other verification failure.
func (c *CertificateToken) CheckIsSignedBy(publicKey *PublicKey) enumerations.SignatureValidity {
	result := enumerations.SignatureValidity_INVALID
	c.SetInvalidityReason("")
	if err := certificateTokenVerify(c.x509Certificate, publicKey); err == nil {
		result = enumerations.SignatureValidity_VALID
	} else {
		c.SetInvalidityReason(certificateTokenInvalidityReason(err))
	}
	c.SetSignatureValidity(result)
	return result
}

// CheckKeyUsage reports whether the certificate carries the given key usage bit.
// Port of checkKeyUsage(KeyUsageBit).
func (c *CertificateToken) CheckKeyUsage(keyUsageBit enumerations.KeyUsageBit) bool {
	for _, bit := range c.KeyUsageBits() {
		if bit == keyUsageBit {
			return true
		}
	}
	return false
}

// KeyUsageBits returns the certificate's key usages, computed once and cached.
// Port of getKeyUsageBits().
//
// Java reads the KeyUsage extension as a boolean array and treats an absent extension as
// null; crypto/x509 collapses both an absent extension and an all-zero one to a zero
// bitmask, which yields the same empty result.
func (c *CertificateToken) KeyUsageBits() []enumerations.KeyUsageBit {
	if c.keyUsageBits == nil {
		keyUsageBits := make([]enumerations.KeyUsageBit, 0)
		keyUsage := c.x509Certificate.KeyUsage
		if keyUsage != 0 {
			for _, keyUsageBit := range enumerations.KeyUsageBitValues() {
				if keyUsage&x509.KeyUsage(1<<uint(keyUsageBit.Index())) != 0 {
					keyUsageBits = append(keyUsageBits, keyUsageBit)
				}
			}
		}
		c.keyUsageBits = keyUsageBits
	}
	return c.keyUsageBits
}

// IsCA reports whether the BasicConstraints extension marks the certificate as a CA.
// Port of isCA().
func (c *CertificateToken) IsCA() bool {
	return c.PathLenConstraint() != -1
}

// PathLenConstraint returns the pathLenConstraint value when the BasicConstraints extension
// is present and cA is true, math.MaxInt32 when it is a CA without a pathLenConstraint, and
// -1 otherwise. Port of getPathLenConstraint(), i.e.
// java.security.cert.X509Certificate#getBasicConstraints().
func (c *CertificateToken) PathLenConstraint() int {
	if !c.x509Certificate.BasicConstraintsValid || !c.x509Certificate.IsCA {
		return -1
	}
	if c.x509Certificate.MaxPathLen > 0 || c.x509Certificate.MaxPathLenZero {
		return c.x509Certificate.MaxPathLen
	}
	return math.MaxInt32
}

// Signature returns the signature value of the certificate. Port of getSignature().
func (c *CertificateToken) Signature() []byte {
	return c.x509Certificate.Signature
}

// BuildTokenIdentifier builds the token's unique identifier.
// Port of the protected buildTokenIdentifier().
func (c *CertificateToken) BuildTokenIdentifier() *TokenIdentifier {
	return &NewCertificateTokenIdentifier(c).TokenIdentifier
}

// ToString returns a string representation of the token using the given indentation.
// Port of toString(String).
func (c *CertificateToken) ToString(indentStr string) string {
	var out strings.Builder
	out.WriteString(indentStr)
	out.WriteString("CertificateToken[\n")
	indentStr += "\t"

	out.WriteString(indentStr)
	out.WriteString("DSS Id              : ")
	out.WriteString(c.DSSIDAsString())
	out.WriteByte('\n')

	out.WriteString(indentStr)
	out.WriteString("Identity Id         : ")
	out.WriteString(c.EntityKey().String())
	out.WriteByte('\n')

	out.WriteString(indentStr)
	out.WriteString("Validity period     : ")
	out.WriteString(c.x509Certificate.NotBefore.Format(certificateTokenDateFormat))
	out.WriteString(" - ")
	out.WriteString(c.x509Certificate.NotAfter.Format(certificateTokenDateFormat))
	out.WriteByte('\n')

	out.WriteString(indentStr)
	out.WriteString("Subject name        : ")
	out.WriteString(c.Subject().Canonical())
	out.WriteByte('\n')

	out.WriteString(indentStr)
	out.WriteString("Issuer subject name : ")
	out.WriteString(c.Issuer().Canonical())
	out.WriteByte('\n')

	out.WriteString(indentStr)
	out.WriteString("Serial Number       : ")
	out.WriteString(c.SerialNumber().String())
	out.WriteByte('\n')

	out.WriteString(indentStr)
	out.WriteString("Signature algorithm : ")
	if c.SignatureAlgorithm() == "" {
		out.WriteString("?")
	} else {
		out.WriteString(string(c.SignatureAlgorithm()))
	}
	out.WriteByte('\n')

	if c.IsSelfSigned() {
		out.WriteString(indentStr)
		out.WriteString("[SELF-SIGNED]")
		out.WriteByte('\n')
	}

	indentStr = indentStr[1:]
	out.WriteString(indentStr)
	out.WriteByte(']')
	return out.String()
}

// String returns ToString(""). Port of toString().
func (c *CertificateToken) String() string {
	return c.ToString("")
}

// Equals reports whether both tokens carry the same DSS Id. Port of Token#equals(Object),
// which additionally requires the two tokens to be of the same class.
//
// Java relies on hashCode-based collections keyed on tokens; the Go equivalent keys maps on
// DSSIDAsString().
func (c *CertificateToken) Equals(other *CertificateToken) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	return c.DSSID().Equals(other.DSSID())
}

// certificateTokenVerify verifies the certificate's own signature with the given public key,
// standing in for java.security.cert.Certificate#verify(PublicKey).
//
// DEVIATION: crypto/x509 refuses MD5-based signatures outright (InsecureAlgorithmError),
// where the JCA would verify them; MD5-signed certificates therefore never report a valid
// signature here.
func certificateTokenVerify(certificate *x509.Certificate, publicKey *PublicKey) error {
	if publicKey == nil {
		return errors.New("NullPointerException : public key is missing")
	}
	key := publicKey.Key()
	if key == nil {
		return errors.New("InvalidKeyException : the public key has not been parsed")
	}
	// CheckSignature verifies "signed" against the receiver's public key, so the signer's
	// key is carried by a throwaway certificate value.
	signer := &x509.Certificate{PublicKey: key}
	return signer.CheckSignature(certificate.SignatureAlgorithm, certificate.RawTBSCertificate, certificate.Signature)
}

// certificateTokenInvalidityReason renders a verification failure the way upstream does:
// the exception's simple class name, " : ", and its message. Go errors have no class name,
// so the error type is used instead.
func certificateTokenInvalidityReason(err error) string {
	return fmt.Sprintf("%T : %s", err, err.Error())
}

// certificateTokenSigAlgOIDAndParams extracts the signatureAlgorithm AlgorithmIdentifier of
// the outer Certificate SEQUENCE, standing in for X509Certificate#getSigAlgOID() and
// #getSigAlgParams(). Like sun.security.x509.AlgorithmId, an absent parameters field and an
// explicit ASN.1 NULL both yield nil parameters.
func certificateTokenSigAlgOIDAndParams(raw []byte) (string, []byte, error) {
	var certificate struct {
		TBSCertificate     asn1.RawValue
		SignatureAlgorithm asn1.RawValue
		SignatureValue     asn1.BitString
	}
	if _, err := asn1.Unmarshal(raw, &certificate); err != nil {
		return "", nil, err
	}
	var algorithmIdentifier struct {
		Algorithm  asn1.ObjectIdentifier
		Parameters asn1.RawValue `asn1:"optional"`
	}
	if _, err := asn1.Unmarshal(certificate.SignatureAlgorithm.FullBytes, &algorithmIdentifier); err != nil {
		return "", nil, err
	}
	params := algorithmIdentifier.Parameters.FullBytes
	if len(params) == 0 ||
		(algorithmIdentifier.Parameters.Class == asn1.ClassUniversal && algorithmIdentifier.Parameters.Tag == asn1.TagNull) {
		params = nil
	}
	return algorithmIdentifier.Algorithm.String(), params, nil
}

// compile-time interface assertions.
var (
	_ Token          = (*CertificateToken)(nil)
	_ TokenOverrides = (*CertificateToken)(nil)
)
