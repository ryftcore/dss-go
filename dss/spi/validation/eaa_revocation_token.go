// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/EAARevocationToken.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCIES (flagged per S2B_BRIEF.md):
//
//   - EAA (Java spi.eaa.EAA, an interface) and EAAStatusTokenIdentifier (Java
//     spi.eaa.EAAStatusTokenIdentifier) are not in this manifest (S2B_BRIEF.md only assigns
//     EAAKeyBindingPayload, EAARevocationToken and EAAValidationParameters from the spi.eaa
//     package to this chunk) but are SCC-flattened into this same Go package by a sibling chunk.
//     Their assumed shapes, inferred from every EAA/EAAStatusTokenIdentifier use in this file
//     (the only Java source calling into them within this manifest):
//
//     type EAA interface {
//     model.IdentifierBasedObject
//     ID() string // getId(), declared directly on EAA (not only inherited via IdentifierBasedObject)
//     // ... remaining getters (Filename, Signatures, EAAType, DisclosureValidations,
//     // KeyBindingSignature, KeyBindingSignaturePayload, DeviceKeyCertificateSource, Payload,
//     // SelectiveDisclosuresDigestAlgorithm) are irrelevant to this file and omitted here.
//     }
//
//     func NewEAAStatusTokenIdentifier(eaaRevocationToken *EAARevocationToken) *EAAStatusTokenIdentifier
//     // EAAStatusTokenIdentifier embeds model.TokenIdentifier, mirroring the Java
//     // "extends TokenIdentifier" relationship (see SignatureIdentifier/TimestampToken's own
//     // BuildTokenIdentifier for the same &X.TokenIdentifier pattern).
//
//   - Concrete EAA revocation tokens (e.g. dss-eaa-revocation-common's EAAStatusListToken, out of
//     scope for phase 2b) embed *EAARevocationToken and call InitToken(self) in their own
//     constructor, per the phase 2a handoff fact that every concrete Token constructor calls the
//     Init pattern. EAARevocationToken itself stays uninitialised (no InitToken call here) because
//     Java's EAARevocationToken is abstract and, like Token#getCreationDate, leaves
//     checkIsSignedBy/buildTokenIdentifier/getIssuerX500Principal available for a subclass to
//     *inherit* rather than re-override - Go achieves the same effect via method promotion once a
//     concrete subclass registers itself (or this type, if unmodified) as the model.TokenOverrides.
//     Consequently EAARevocationToken does not fully satisfy model.Token on its own (it also never
//     overrides CreationDate(), left abstract by Java too) and carries no compile-time Token
//     assertion, matching Token#getCreationDate being abstract in the Java source.
package validation

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// EAARevocationToken represents an EAA revocation representation.
type EAARevocationToken struct {
	model.TokenBase

	// encoded are the extracted binaries of the EAA Revocation Token.
	encoded *EAARevocationTokenBinary

	// signature is used to sign the EAA revocation data.
	signature AdvancedSignature

	// relatedEAA is the EAA related to this status object.
	relatedEAA EAA

	// sourceURL is the URL which was used to obtain the revocation data (online).
	sourceURL string

	// origin is the external origin (EXTERNAL or CACHED).
	origin enumerations.EAARevocationOrigin

	// status contains the revocation status of the token.
	status enumerations.EAAStatus

	// certificateSource is built on the extracted information from the EAA revocation.
	certificateSource *spi.TokenCertificateSource
}

// SetRelatedEAA sets a related EAA. Port of setRelatedEAA(EAA).
func (t *EAARevocationToken) SetRelatedEAA(relatedEAA EAA) {
	t.relatedEAA = relatedEAA
}

// SourceURL gets the source URL used to access the status token. Port of getSourceURL().
func (t *EAARevocationToken) SourceURL() string {
	return t.sourceURL
}

// SetSourceURL sets the source URL used to access the status token. Port of setSourceURL(String).
func (t *EAARevocationToken) SetSourceURL(sourceURL string) {
	t.sourceURL = sourceURL
}

// Origin gets the origin of the status token (e.g. EXTERNAL or CACHED). Port of getOrigin().
func (t *EAARevocationToken) Origin() enumerations.EAARevocationOrigin {
	return t.origin
}

// SetOrigin sets the origin of the status token (e.g. EXTERNAL or CACHED). Port of setOrigin(EAARevocationOrigin).
func (t *EAARevocationToken) SetOrigin(origin enumerations.EAARevocationOrigin) {
	t.origin = origin
}

// Signature gets the signature used to sign the EAA revocation token. Port of getSignature().
func (t *EAARevocationToken) Signature() AdvancedSignature {
	return t.signature
}

// Status gets the indication of the status of the related token (e.g. VALID, INVALID, etc.).
// Port of getStatus().
func (t *EAARevocationToken) Status() enumerations.EAAStatus {
	return t.status
}

// CertificateSource gets the certificate source built on the extracted EAA revocation
// information. Port of getCertificateSource().
func (t *EAARevocationToken) CertificateSource() *spi.TokenCertificateSource {
	return t.certificateSource
}

// SetCertificateSource sets the certificate source built on the extracted EAA revocation
// information. Port of setCertificateSource(TokenCertificateSource).
func (t *EAARevocationToken) SetCertificateSource(certificateSource *spi.TokenCertificateSource) {
	t.certificateSource = certificateSource
}

// Type gets the type of the token. Port of getType().
func (t *EAARevocationToken) Type() string {
	if t.signature != nil {
		return t.signature.SignatureType()
	}
	return ""
}

// Subject gets the subject of the token. Port of getSubject(); not implemented by default.
func (t *EAARevocationToken) Subject() string {
	// not implemented by default
	return ""
}

// SubjectMatch gets whether the subject defined in the EAA revocation token matches the
// value defined in the EAA. Port of getSubjectMatch(); returns nil (Java NULL) when not
// supported, otherwise TRUE/FALSE.
func (t *EAARevocationToken) SubjectMatch() *bool {
	// not implemented by default
	return nil
}

// ExpirationDate gets the expiration date of the token. Port of getExpirationDate(); not
// implemented by default.
func (t *EAARevocationToken) ExpirationDate() *time.Time {
	// not implemented by default
	return nil
}

// TimeToLive gets the time in seconds when a consumer should request a new token after its
// extraction. Port of getTimeToLive(); Java's Number is a polymorphic numeric wrapper, mapped
// to `any` per the same convention as claim.Claim.NumberValue; not implemented by default.
func (t *EAARevocationToken) TimeToLive() any {
	// not implemented by default
	return nil
}

// CreationDate is abstract in Java's Token and is never overridden by EAARevocationToken
// (verified against upstream: EAARevocationToken.java overrides buildTokenIdentifier,
// checkIsSignedBy and getIssuerX500Principal only). Java can still type-check a bare
// EAARevocationToken reference because the class itself is abstract and never instantiated
// directly. Go interface satisfaction is structural, so a concrete method is required here
// purely so *EAARevocationToken can be passed where model.Token is expected (e.g.
// NewEAAStatusTokenIdentifier); a concrete subclass (dss-eaa-revocation-common, out of phase
// 2b scope) is expected to define its own CreationDate(), which shadows this one via Go's
// embedding rules exactly like a Java override would. Calling this stub directly panics.
func (t *EAARevocationToken) CreationDate() time.Time {
	panic("eu.europa.esig.dss.model.x509.Token.getCreationDate: abstract, not implemented by EAARevocationToken")
}

// BuildTokenIdentifier builds the token's unique identifier.
// Port of the protected buildTokenIdentifier() override.
func (t *EAARevocationToken) BuildTokenIdentifier() *model.TokenIdentifier {
	return &NewEAAStatusTokenIdentifier(t).TokenIdentifier
}

// CheckIsSignedBy is not supported for an EAA revocation token: the signature is verified
// through the wrapping AdvancedSignature, never per public key. Port of the protected
// checkIsSignedBy(PublicKey) override, which throws UnsupportedOperationException(getClass().getName()).
func (t *EAARevocationToken) CheckIsSignedBy(publicKey *model.PublicKey) enumerations.SignatureValidity {
	panic("eu.europa.esig.dss.spi.eaa.EAARevocationToken")
}

// IssuerX500Principal returns the X500Principal of the certificate that signed this token.
// Port of getIssuerX500Principal().
func (t *EAARevocationToken) IssuerX500Principal() *model.X500Principal {
	if t.signature.SigningCertificateToken() != nil {
		return t.signature.SigningCertificateToken().Subject().Principal()
	}
	return nil
}

// RelatedEAA gets the related EAA. Port of getRelatedEAA().
func (t *EAARevocationToken) RelatedEAA() EAA {
	return t.relatedEAA
}

// RelatedEAAId gets the String identifier of the related EAA. Port of getRelatedEAAId().
func (t *EAARevocationToken) RelatedEAAId() string {
	if t.relatedEAA != nil {
		return t.relatedEAA.ID()
	}
	return ""
}

// ToString returns a string representation of the token using the given indentation.
// Port of toString(String); upstream leaves this as a TODO returning "".
func (t *EAARevocationToken) ToString(indentStr string) string {
	// TODO : to be implemented
	return ""
}

// String returns ToString(""). Port of toString(), inherited from Token in Java (never
// overridden by EAARevocationToken); ported directly per the model.Token convention every
// concrete token type follows (e.g. model.CertificateToken.String()).
func (t *EAARevocationToken) String() string {
	return t.ToString("")
}

// Encoded returns the encoded form of the wrapped token. Port of getEncoded().
func (t *EAARevocationToken) Encoded() []byte {
	return t.encoded.Binaries()
}
