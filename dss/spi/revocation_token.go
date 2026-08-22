// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/RevocationToken.java (DSS 6.5.RC1).
//
// RevocationToken<R> is the abstract superclass CRLToken and OCSPToken embed via
// RevocationTokenBase[R] and register with InitRevocationToken, following the same
// self-registration pattern as model.TokenBase.InitToken.
package spi

import (
	"math/big"
	"reflect"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// RevocationToken represents a revocation data token (CRL or OCSP).
//
// Java narrows Token#getDSSId() to return a RevocationTokenIdentifier via covariant override;
// Go cannot express covariant returns, so DSSID() keeps the model.Token signature and callers
// needing the narrower type assert on the result (it is always a *RevocationTokenIdentifier).
type RevocationToken[R revocation.Revocation] interface {
	model.Token

	// RevocationType returns the Revocation Token type (CRL or OCSP). Port of the abstract
	// getRevocationType().
	RevocationType() enumerations.RevocationType

	// RelatedCertificate returns a certificate token the current revocation data has been
	// issued for. Port of getRelatedCertificate().
	RelatedCertificate() *model.CertificateToken
	// SetRelatedCertificate sets the certificate token the current revocation data has been
	// issued for. Java writes the protected field directly from subclasses in a different
	// Java package (eu.europa.esig.dss.spi.x509.revocation.crl/.ocsp); the Go port flattens
	// those into this same package but keeps the setter for symmetry with the rest of the
	// Init/Set pattern used throughout this package.
	SetRelatedCertificate(certificate *model.CertificateToken)
	// RelatedCertificateID gets the DSS String Id of the related certificate, "" when there is
	// none (Java returns null). Port of getRelatedCertificateId().
	RelatedCertificateID() string

	// IssuerCertificateToken returns the issuer CertificateToken. Port of the abstract
	// getIssuerCertificateToken().
	IssuerCertificateToken() *model.CertificateToken

	// SourceURL returns the URL of the source (if available). Port of getSourceURL().
	SourceURL() string
	// SetSourceURL sets the revocation data source URL; only used for an OnlineSource. Port of
	// setSourceURL(String).
	SetSourceURL(sourceURL string)

	// Status returns the certificate status. Port of getStatus().
	Status() enumerations.CertificateStatus
	// SetStatus sets the certificate status.
	SetStatus(status enumerations.CertificateStatus)

	// ProductionDate returns the generation time of the current revocation data (when it was
	// signed). Port of getProductionDate().
	ProductionDate() time.Time
	// SetProductionDate sets the generation time of the current revocation data.
	SetProductionDate(productionDate time.Time)

	// ThisUpdate returns the date of the this update. Port of getThisUpdate().
	ThisUpdate() time.Time
	// SetThisUpdate sets the date of the this update.
	SetThisUpdate(thisUpdate time.Time)

	// NextUpdate returns the date of the next update. Port of getNextUpdate().
	NextUpdate() time.Time
	// SetNextUpdate sets the date of the next update.
	SetNextUpdate(nextUpdate time.Time)

	// RevocationDate returns the revocation date (if the token has been revoked), the zero
	// time otherwise (Java returns null). Port of getRevocationDate().
	RevocationDate() time.Time
	// SetRevocationDate sets the revocation date.
	SetRevocationDate(revocationDate time.Time)

	// CRLNumber gets the sequential number of the revocation token, when present (CRL only).
	// Port of getCRLNumber().
	CRLNumber() *big.Int
	// SetCRLNumber sets the sequential number of the revocation token.
	SetCRLNumber(crlNumber *big.Int)

	// ExpiredCertsOnCRL returns the expiredCertsOnCRL date (from CRL), the zero time when
	// absent. Port of getExpiredCertsOnCRL().
	ExpiredCertsOnCRL() time.Time
	// SetExpiredCertsOnCRL sets the expiredCertsOnCRL date.
	SetExpiredCertsOnCRL(expiredCertsOnCRL time.Time)

	// ArchiveCutOff returns the archiveCutOff date (from an OCSP Response), the zero time when
	// absent. Port of getArchiveCutOff().
	ArchiveCutOff() time.Time
	// SetArchiveCutOff sets the archiveCutOff date.
	SetArchiveCutOff(archiveCutOff time.Time)

	// CertHashPresent returns true if the certHash extension (from an OCSP Response) is
	// present. Port of isCertHashPresent().
	CertHashPresent() bool
	// SetCertHashPresent sets whether the certHash extension is present.
	SetCertHashPresent(certHashPresent bool)

	// CertHashMatch returns true if the certHash extension (from an OCSP Response) matches the
	// hash of the related certificate token. Port of isCertHashMatch().
	CertHashMatch() bool
	// SetCertHashMatch sets whether the certHash extension matches the related certificate.
	SetCertHashMatch(certHashMatch bool)

	// Reason returns the revocation reason (if the token has been revoked), "" otherwise (Java
	// returns null). Port of getReason().
	Reason() enumerations.RevocationReason
	// SetReason sets the revocation reason.
	SetReason(reason enumerations.RevocationReason)

	// CertificateSource returns a source of embedded into a revocation token certificates.
	// Port of the abstract getCertificateSource().
	CertificateSource() RevocationCertificateSource
	// Certificates returns a collection of embedded certificates; nil (Java: empty collection)
	// for a CRL. Port of the abstract getCertificates().
	Certificates() []*model.CertificateToken

	// SetExternalOrigin sets the external origin.
	//
	// Panics with the Java messages: when origin is "" (Objects.requireNonNull("The origin is
	// null")) and when origin is not EXTERNAL or CACHED (IllegalArgumentException("Only
	// external are allowed")). Port of setExternalOrigin(RevocationOrigin).
	SetExternalOrigin(origin enumerations.RevocationOrigin)
	// ExternalOrigin gets the external origin, "" when unset (Java returns null). Port of
	// getExternalOrigin().
	ExternalOrigin() enumerations.RevocationOrigin
	// IsInternal returns true if the token was not collected from an external resource
	// (online or jdbc), i.e. it comes from a signature/timestamp. Port of isInternal().
	IsInternal() bool

	// Equals reports whether both revocation tokens share the same DSS identifier and related
	// certificate. Port of equals(Object).
	//
	// Java's getClass() check has no direct Go counterpart here since the check is implicit in
	// the type system: RevocationToken[R] restricts comparisons to concrete tokens sharing the
	// same R, and today CRLToken is the sole implementation for R=CRL and OCSPToken the sole
	// one for R=OCSP, so two same-R tokens are always of the same concrete Go type. Should a
	// second implementation for the same R ever appear, RevocationTokenBase.Equals's reflect
	// check (see below) still tells them apart.
	Equals(other RevocationToken[R]) bool
}

// RevocationTokenOverrides is the private counterpart of RevocationToken kept purely for the
// self-registration call: BuildTokenIdentifier is promoted from RevocationTokenBase and reaches
// model.TokenBase.InitToken through it, exactly like model.TokenOverrides does for every other
// Token. It exists only so InitRevocationToken can type-assert the concrete token also
// implements model.TokenOverrides without leaking those dispatch-only methods into the public
// RevocationToken interface above.
type revocationTokenOverridesHolder[R revocation.Revocation] interface {
	revocationTokenBase() *RevocationTokenBase[R]
}

// RevocationTokenBase carries the state and the concrete behaviour of the Java abstract class
// RevocationToken<R>. Concrete tokens (CRLToken, OCSPToken) embed it and register themselves
// with InitRevocationToken.
type RevocationTokenBase[R revocation.Revocation] struct {
	model.TokenBase

	// overrides points back at the concrete revocation token; see InitRevocationToken. It
	// doubles as the model.Token BuildTokenIdentifier digests (Java's buildTokenIdentifier()
	// override calls this.getEncoded() through the RevocationTokenIdentifier constructor).
	overrides RevocationToken[R]

	// relatedCertificate is the CertificateToken this revocation object relates to.
	relatedCertificate *model.CertificateToken
	// sourceURL is the URL which was used to obtain the revocation data (online).
	sourceURL string
	// externalOrigin is the external origin (ONLINE or CACHED); "" stands for Java's null.
	externalOrigin enumerations.RevocationOrigin
	// status contains the revocation status of the token.
	status enumerations.CertificateStatus
	// productionDate represents the production date of the OCSP response or the thisUpdate in
	// case of CRL.
	productionDate time.Time
	// thisUpdate represents the this update date of the CRL.
	thisUpdate time.Time
	// nextUpdate represents the next update date of the CRL, or the zero time for OCSP.
	nextUpdate time.Time
	// revocationDate represents the revocation date from an X509CRLEntry or from a
	// BasicOCSPResp (if the related certificate is revoked).
	revocationDate time.Time
	// crlNumber is the sequential number of the revocation token (applicable for CRLs only).
	crlNumber *big.Int
	// expiredCertsOnCRL is the expired-certs-on-crl time extension.
	expiredCertsOnCRL time.Time
	// archiveCutOff is the archive-cut-off time extension.
	archiveCutOff time.Time
	// certHashPresent represents whether the certHash extension from an OCSP Response is
	// present (optional).
	certHashPresent bool
	// certHashMatch represents whether the certHash extension from an OCSP Response matches
	// the related certificate's hash (optional).
	certHashMatch bool
	// reason is the reason of the revocation.
	reason enumerations.RevocationReason
}

// NewRevocationTokenBase instantiates the base state of a revocation token with the Java
// default values. Port of the protected default constructor; the concrete token must still
// call InitRevocationToken.
func NewRevocationTokenBase[R revocation.Revocation]() RevocationTokenBase[R] {
	return RevocationTokenBase[R]{TokenBase: model.NewTokenBase()}
}

// InitRevocationToken registers the concrete revocation token with its base so that the base
// can dispatch to the operations Java would reach through virtual dispatch (BuildTokenIdentifier
// here, plus everything model.TokenBase itself dispatches - CheckIsSignedBy,
// IssuerX500Principal, IsSelfSigned). It must be called exactly once, by the concrete token's
// constructor, before any other method - mirroring model.TokenBase.InitToken, which it also
// drives.
//
// Panics if overrides does not implement model.TokenOverrides: CRLToken and OCSPToken satisfy
// it by combining their own CheckIsSignedBy/IssuerX500Principal with BuildTokenIdentifier and
// IsSelfSigned promoted from this struct and model.TokenBase respectively; a token missing one
// of the two direct overrides is a programmer error caught here rather than surfacing as a
// puzzling nil-overrides panic deep inside model.TokenBase.
func (t *RevocationTokenBase[R]) InitRevocationToken(overrides RevocationToken[R]) {
	t.overrides = overrides
	tokenOverrides, ok := overrides.(model.TokenOverrides)
	if !ok {
		panic("RevocationToken was not initialised: the concrete token must implement " +
			"model.TokenOverrides (CheckIsSignedBy and IssuerX500Principal)")
	}
	t.TokenBase.InitToken(tokenOverrides)
}

// revocationTokenBaseOverrides returns the registered overrides, panicking when the concrete
// token forgot to call InitRevocationToken.
func (t *RevocationTokenBase[R]) revocationTokenBaseOverrides() RevocationToken[R] {
	if t.overrides == nil {
		panic("RevocationToken was not initialised: the concrete token must call InitRevocationToken in its constructor")
	}
	return t.overrides
}

// BuildTokenIdentifier builds the token's unique identifier. Port of the buildTokenIdentifier()
// override: `return new RevocationTokenIdentifier(this);`. This is promoted to the concrete
// token (CRLToken, OCSPToken), which is what lets it satisfy model.TokenOverrides without
// providing its own BuildTokenIdentifier.
func (t *RevocationTokenBase[R]) BuildTokenIdentifier() *model.TokenIdentifier {
	identifier := NewRevocationTokenIdentifier(t.revocationTokenBaseOverrides())
	return &identifier.TokenIdentifier
}

// RelatedCertificate returns the related CertificateToken. Port of getRelatedCertificate().
func (t *RevocationTokenBase[R]) RelatedCertificate() *model.CertificateToken {
	return t.relatedCertificate
}

// SetRelatedCertificate sets the related CertificateToken.
func (t *RevocationTokenBase[R]) SetRelatedCertificate(certificate *model.CertificateToken) {
	t.relatedCertificate = certificate
}

// RelatedCertificateID gets DSS String Id of the related certificate, "" when there is none.
// Port of getRelatedCertificateId().
func (t *RevocationTokenBase[R]) RelatedCertificateID() string {
	if t.relatedCertificate != nil {
		return t.relatedCertificate.DSSIDAsString()
	}
	return ""
}

// SourceURL returns the URL of the source (if available). Port of getSourceURL().
func (t *RevocationTokenBase[R]) SourceURL() string {
	return t.sourceURL
}

// SetSourceURL sets the revocation data source URL. Port of setSourceURL(String).
func (t *RevocationTokenBase[R]) SetSourceURL(sourceURL string) {
	t.sourceURL = sourceURL
}

// Status returns the certificate status. Port of getStatus().
func (t *RevocationTokenBase[R]) Status() enumerations.CertificateStatus {
	return t.status
}

// SetStatus sets the certificate status.
func (t *RevocationTokenBase[R]) SetStatus(status enumerations.CertificateStatus) {
	t.status = status
}

// ProductionDate returns the generation time of the current revocation data. Port of
// getProductionDate().
func (t *RevocationTokenBase[R]) ProductionDate() time.Time {
	return t.productionDate
}

// SetProductionDate sets the generation time of the current revocation data.
func (t *RevocationTokenBase[R]) SetProductionDate(productionDate time.Time) {
	t.productionDate = productionDate
}

// CreationDate returns the production date. Port of the getCreationDate() override.
func (t *RevocationTokenBase[R]) CreationDate() time.Time {
	return t.productionDate
}

// ThisUpdate returns the date of the this update. Port of getThisUpdate().
func (t *RevocationTokenBase[R]) ThisUpdate() time.Time {
	return t.thisUpdate
}

// SetThisUpdate sets the date of the this update.
func (t *RevocationTokenBase[R]) SetThisUpdate(thisUpdate time.Time) {
	t.thisUpdate = thisUpdate
}

// NextUpdate returns the date of the next update. Port of getNextUpdate().
func (t *RevocationTokenBase[R]) NextUpdate() time.Time {
	return t.nextUpdate
}

// SetNextUpdate sets the date of the next update.
func (t *RevocationTokenBase[R]) SetNextUpdate(nextUpdate time.Time) {
	t.nextUpdate = nextUpdate
}

// RevocationDate returns the revocation date (if the token has been revoked). Port of
// getRevocationDate().
func (t *RevocationTokenBase[R]) RevocationDate() time.Time {
	return t.revocationDate
}

// SetRevocationDate sets the revocation date.
func (t *RevocationTokenBase[R]) SetRevocationDate(revocationDate time.Time) {
	t.revocationDate = revocationDate
}

// CRLNumber gets the sequential number of the revocation token, when present. Port of
// getCRLNumber().
func (t *RevocationTokenBase[R]) CRLNumber() *big.Int {
	return t.crlNumber
}

// SetCRLNumber sets the sequential number of the revocation token.
func (t *RevocationTokenBase[R]) SetCRLNumber(crlNumber *big.Int) {
	t.crlNumber = crlNumber
}

// ExpiredCertsOnCRL returns the expiredCertsOnCRL date (from CRL). Port of
// getExpiredCertsOnCRL().
func (t *RevocationTokenBase[R]) ExpiredCertsOnCRL() time.Time {
	return t.expiredCertsOnCRL
}

// SetExpiredCertsOnCRL sets the expiredCertsOnCRL date.
func (t *RevocationTokenBase[R]) SetExpiredCertsOnCRL(expiredCertsOnCRL time.Time) {
	t.expiredCertsOnCRL = expiredCertsOnCRL
}

// ArchiveCutOff returns the archiveCutOff date (from an OCSP Response). Port of
// getArchiveCutOff().
func (t *RevocationTokenBase[R]) ArchiveCutOff() time.Time {
	return t.archiveCutOff
}

// SetArchiveCutOff sets the archiveCutOff date.
func (t *RevocationTokenBase[R]) SetArchiveCutOff(archiveCutOff time.Time) {
	t.archiveCutOff = archiveCutOff
}

// CertHashPresent returns true if the certHash extension is present. Port of
// isCertHashPresent().
func (t *RevocationTokenBase[R]) CertHashPresent() bool {
	return t.certHashPresent
}

// SetCertHashPresent sets whether the certHash extension is present.
func (t *RevocationTokenBase[R]) SetCertHashPresent(certHashPresent bool) {
	t.certHashPresent = certHashPresent
}

// CertHashMatch returns true if the certHash extension matches the related certificate's hash.
// Port of isCertHashMatch().
func (t *RevocationTokenBase[R]) CertHashMatch() bool {
	return t.certHashMatch
}

// SetCertHashMatch sets whether the certHash extension matches the related certificate.
func (t *RevocationTokenBase[R]) SetCertHashMatch(certHashMatch bool) {
	t.certHashMatch = certHashMatch
}

// Reason returns the revocation reason (if the token has been revoked). Port of getReason().
func (t *RevocationTokenBase[R]) Reason() enumerations.RevocationReason {
	return t.reason
}

// SetReason sets the revocation reason.
func (t *RevocationTokenBase[R]) SetReason(reason enumerations.RevocationReason) {
	t.reason = reason
}

// SetExternalOrigin sets the external origin.
//
// Panics with the Java messages: origin == "" mirrors Objects.requireNonNull("The origin is
// null"); an internal origin (anything but EXTERNAL/CACHED) mirrors the
// IllegalArgumentException("Only external are allowed"). Port of
// setExternalOrigin(RevocationOrigin).
func (t *RevocationTokenBase[R]) SetExternalOrigin(origin enumerations.RevocationOrigin) {
	if origin == "" {
		panic("The origin is null")
	}
	if origin.IsInternalOrigin() {
		panic("Only external are allowed")
	}
	t.externalOrigin = origin
}

// ExternalOrigin gets the external origin, "" when unset. Port of getExternalOrigin().
func (t *RevocationTokenBase[R]) ExternalOrigin() enumerations.RevocationOrigin {
	return t.externalOrigin
}

// IsInternal returns true if the token was not collected from an external resource. Port of
// isInternal().
func (t *RevocationTokenBase[R]) IsInternal() bool {
	return t.externalOrigin == ""
}

// Equals reports whether both revocation tokens share the same DSS identifier and related
// certificate. Port of equals(Object); see the RevocationToken.Equals doc for how the Java
// getClass() check maps onto Go's type system.
func (t *RevocationTokenBase[R]) Equals(other RevocationToken[R]) bool {
	if other == nil {
		return false
	}
	otherBase := revocationTokenBaseOf[R](other)
	if otherBase == nil {
		return false
	}
	if t == otherBase {
		return true
	}
	if reflect.TypeOf(t.overrides) != reflect.TypeOf(otherBase.overrides) {
		return false
	}
	if !t.DSSID().Equals(other.DSSID()) {
		return false
	}
	if t.relatedCertificate == nil {
		return otherBase.relatedCertificate == nil
	}
	return t.relatedCertificate.Equals(otherBase.relatedCertificate)
}

// revocationTokenBase implements the unexported accessor revocationTokenBaseOf looks for; it
// is promoted to every type embedding RevocationTokenBase, mirroring model.identifierBaseOf's
// pattern for reaching an embedded base through an interface value.
func (t *RevocationTokenBase[R]) revocationTokenBase() *RevocationTokenBase[R] {
	return t
}

// revocationTokenBaseOf extracts the embedded RevocationTokenBase of any revocation token,
// which is how Equals reaches the fields Java's RevocationToken#equals compares directly.
func revocationTokenBaseOf[R revocation.Revocation](token RevocationToken[R]) *RevocationTokenBase[R] {
	if holder, ok := token.(revocationTokenOverridesHolder[R]); ok {
		return holder.revocationTokenBase()
	}
	return nil
}
