// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/CertificateVerifier.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY (flagged per S2B_BRIEF.md): TrustAnchorVerifier (Java
// spi.validation.TrustAnchorVerifier) is assigned to sibling chunk VAL-C (s2b_VAL-C.txt),
// which lands it in this same package. It is referenced here by name only, matching the
// shape already documented by revocation_data_verifier.go's header comment
// (NewDefaultTrustAnchorVerifier() *TrustAnchorVerifier).
//
// Java overloads setTrustedCertSources/setAdjunctCertSources for CertificateSource... and
// ListCertificateSource; Go cannot overload by parameter type, so the ListCertificateSource
// variants are named SetTrustedCertSourcesFromList/SetAdjunctCertSourcesFromList, following
// the naming convention established by model/policy/cryptographic_suite_factory.go
// (LoadCryptographicSuiteFromReader).
package validation

import (
	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/x509/aia"
)

// CertificateVerifier provides information on the sources to be used in the validation
// process in the context of a signature.
type CertificateVerifier interface {
	// CrlSource returns the CRL source associated with this verifier: the used CRL source for
	// external access (web, filesystem, cached, ...). Port of getCrlSource().
	CrlSource() spi.RevocationSource[revocation.CRL]

	// SetCrlSource defines the source of CRL used by this class. Port of setCrlSource(...).
	SetCrlSource(crlSource spi.RevocationSource[revocation.CRL])

	// OcspSource returns the OCSP source associated with this verifier: the used OCSP source
	// for external access (web, filesystem, cached, ...). Port of getOcspSource().
	OcspSource() spi.RevocationSource[revocation.OCSP]

	// SetOcspSource defines the source of OCSP used by this class. Port of setOcspSource(...).
	SetOcspSource(ocspSource spi.RevocationSource[revocation.OCSP])

	// RevocationDataLoadingStrategyFactory returns a factory used to create revocation data
	// loading strategy associated with this verifier: creates the defined strategy to fetch
	// OCSP or CRL for certificate validation. Port of
	// getRevocationDataLoadingStrategyFactory().
	RevocationDataLoadingStrategyFactory() RevocationDataLoadingStrategyFactory

	// SetRevocationDataLoadingStrategyFactory creates a strategy used to fetch OCSP or CRL for
	// certificate validation. Default: OCSPFirstRevocationDataLoadingStrategyFactory used to
	// create a strategy to extract OCSP token first and CRL after. Port of
	// setRevocationDataLoadingStrategyFactory(...).
	SetRevocationDataLoadingStrategyFactory(revocationDataLoadingStrategyFactory RevocationDataLoadingStrategyFactory)

	// RevocationDataVerifier returns a RevocationDataVerifier associated with this verifier.
	// Port of getRevocationDataVerifier().
	RevocationDataVerifier() *RevocationDataVerifier

	// SetRevocationDataVerifier sets the RevocationDataVerifier used to validate acceptance of
	// the retrieved (from offline or online sources) revocation data. This class is used to
	// verify revocation data extracted from the validating document itself, as well the
	// revocation data retrieved from remote sources during the validation process.
	//
	// NOTE: It is not recommended to use the same instance of RevocationDataVerifier within
	// different CertificateVerifiers, as it may lead to concurrency issues during the
	// execution in multi-threaded environments. Please use a new RevocationDataVerifier per
	// each CertificateVerifier.
	//
	// Port of setRevocationDataVerifier(...).
	SetRevocationDataVerifier(revocationDataVerifier *RevocationDataVerifier)

	// IsRevocationFallback returns whether revocation data still shall be returned if
	// validation of requested revocation data failed (i.e. both for OCSP and CRL). Port of
	// isRevocationFallback().
	IsRevocationFallback() bool

	// SetRevocationFallback sets whether a revocation data still have to be returned to the
	// validation process, in case validation of obtained revocation data has failed (i.e.
	// both for OCSP and CRL). Default: FALSE (invalid revocation data not returned).
	//
	// NOTE: Revocation fallback is enforced to TRUE (return even invalid revocation data, when
	// no valid found) on signature validation.
	//
	// Port of setRevocationFallback(boolean).
	SetRevocationFallback(revocationFallback bool)

	// TimestampTokenVerifier returns a TimestampTokenVerifier associated with this verifier.
	// Port of getTimestampTokenVerifier().
	TimestampTokenVerifier() *TimestampTokenVerifier

	// SetTimestampTokenVerifier sets the TimestampTokenVerifier used to validate acceptance of
	// the timestamp tokens encapsulated within the signature.
	//
	// NOTE: This object is not synchronized by default with the used XML Validation Policy.
	// Please configure the object yourself in case a customized behavior is expected for
	// acceptance of timestamp tokens.
	//
	// Port of setTimestampTokenVerifier(...).
	SetTimestampTokenVerifier(timestampTokenVerifier *TimestampTokenVerifier)

	// TrustAnchorVerifier returns a TrustAnchorVerifier associated with this verifier. Port of
	// getTrustAnchorVerifier().
	TrustAnchorVerifier() *TrustAnchorVerifier

	// SetTrustAnchorVerifier sets the TrustAnchorVerifier used to validate acceptance of the
	// trust anchors.
	//
	// NOTE: This object is not synchronized by default with the used XML Validation Policy.
	// The trusted certificate source set within CertificateVerifier is synchronized
	// automatically unless explicitly defined within TrustAnchorVerifier. Please configure the
	// object yourself in case a customized behavior is expected for acceptance of trust
	// anchors.
	//
	// Port of setTrustAnchorVerifier(...).
	SetTrustAnchorVerifier(trustAnchorVerifier *TrustAnchorVerifier)

	// TrustedCertSources returns the trusted certificate sources associated with this
	// verifier. These sources are used to identify the trusted anchors. Port of
	// getTrustedCertSources().
	TrustedCertSources() *spi.ListCertificateSource

	// SetTrustedCertSources sets multiple trusted certificate sources. Port of
	// setTrustedCertSources(CertificateSource...).
	SetTrustedCertSources(certSources ...spi.CertificateSource)

	// AddTrustedCertSources adds trusted certificate sources to an existing list of trusted
	// certificate sources. Port of addTrustedCertSources(CertificateSource...).
	AddTrustedCertSources(certSources ...spi.CertificateSource)

	// SetTrustedCertSourcesFromList sets a list of trusted certificate sources. Port of the
	// setTrustedCertSources(ListCertificateSource) overload.
	SetTrustedCertSourcesFromList(trustedListCertificateSource *spi.ListCertificateSource)

	// AdjunctCertSources returns the list of adjunct certificate sources assigned to this
	// verifier: the certificate source which contains additional certificate (missing CA,
	// ...). Port of getAdjunctCertSources().
	AdjunctCertSources() *spi.ListCertificateSource

	// SetAdjunctCertSources sets multiple adjunct certificate sources. Port of
	// setAdjunctCertSources(CertificateSource...).
	SetAdjunctCertSources(certSources ...spi.CertificateSource)

	// AddAdjunctCertSources adds adjunct certificate sources to an existing list of adjunct
	// certificate sources. Port of addAdjunctCertSources(CertificateSource...).
	AddAdjunctCertSources(certSources ...spi.CertificateSource)

	// SetAdjunctCertSourcesFromList sets a list of adjunct certificate sources. Port of the
	// setAdjunctCertSources(ListCertificateSource) overload.
	SetAdjunctCertSourcesFromList(adjunctListCertificateSource *spi.ListCertificateSource)

	// AIASource gets the AIASource used to load a model.CertificateToken's issuer by defined
	// AIA URI(s) within the token. Port of getAIASource().
	AIASource() aia.AIASource

	// SetAIASource sets the AIASource used to load a model.CertificateToken's issuer by
	// defined AIA URI(s) within the token. Port of setAIASource(...).
	SetAIASource(aiaSource aia.AIASource)

	// SetAlertOnInvalidSignature allows to change the behavior on invalid signature (T/LT/LTA
	// augmentation). Nil provides a possibility to skip check execution. Default:
	// ExceptionOnStatusAlert - throw an exception. Port of setAlertOnInvalidSignature(...).
	SetAlertOnInvalidSignature(alertOnInvalidSignature alert.StatusAlert)

	// AlertOnInvalidSignature returns the defined execution behaviour on invalid signature.
	// Port of getAlertOnInvalidSignature().
	AlertOnInvalidSignature() alert.StatusAlert

	// SetAlertOnInvalidTimestamp allows to change the behavior on invalid timestamp (LT/LTA
	// augmentation). Nil provides a possibility to skip check execution. Default:
	// ExceptionOnStatusAlert - throw an exception. Port of setAlertOnInvalidTimestamp(...).
	SetAlertOnInvalidTimestamp(alertOnInvalidTimestamp alert.StatusAlert)

	// AlertOnInvalidTimestamp returns the defined execution behaviour on invalid timestamp.
	// Port of getAlertOnInvalidTimestamp().
	AlertOnInvalidTimestamp() alert.StatusAlert

	// SetAlertOnMissingRevocationData allows to change the behavior on missing revocation data
	// (LT/LTA augmentation). Nil provides a possibility to skip check execution. Default:
	// ExceptionOnStatusAlert - throw an exception. Port of
	// setAlertOnMissingRevocationData(...).
	SetAlertOnMissingRevocationData(alertOnMissingRevocationData alert.StatusAlert)

	// AlertOnMissingRevocationData returns the defined execution behaviour on missing
	// revocation data. Port of getAlertOnMissingRevocationData().
	AlertOnMissingRevocationData() alert.StatusAlert

	// SetAlertOnRevokedCertificate allows to change the behavior on revoked certificates
	// (LT/LTA augmentation). Nil provides a possibility to skip check execution. Default:
	// ExceptionOnStatusAlert - throw an exception. Port of setAlertOnRevokedCertificate(...).
	SetAlertOnRevokedCertificate(alertOnRevokedCertificate alert.StatusAlert)

	// AlertOnRevokedCertificate returns the defined execution behaviour on revoked
	// certificate. Port of getAlertOnRevokedCertificate().
	AlertOnRevokedCertificate() alert.StatusAlert

	// SetAlertOnNoRevocationAfterBestSignatureTime allows to change the behavior on revocation
	// data issued after a control time. Nil provides a possibility to skip check execution.
	// Default: LogOnStatusAlert - log a warning. Port of
	// setAlertOnNoRevocationAfterBestSignatureTime(...).
	SetAlertOnNoRevocationAfterBestSignatureTime(alertOnNoRevocationAfterBestSignatureTime alert.StatusAlert)

	// AlertOnNoRevocationAfterBestSignatureTime returns the defined execution behaviour if no
	// revocation data obtained with an issuance time after the bestSignatureTime. Port of
	// getAlertOnNoRevocationAfterBestSignatureTime().
	AlertOnNoRevocationAfterBestSignatureTime() alert.StatusAlert

	// SetAlertOnUncoveredPOE allows to change the behavior on uncovered POE (timestamp). Nil
	// provides a possibility to skip check execution. Default: LogOnStatusAlert - log a
	// warning. Port of setAlertOnUncoveredPOE(...).
	SetAlertOnUncoveredPOE(alertOnUncoveredPOE alert.StatusAlert)

	// AlertOnUncoveredPOE returns the defined execution behaviour on uncovered POE (timestamp).
	// Port of getAlertOnUncoveredPOE().
	AlertOnUncoveredPOE() alert.StatusAlert

	// SetAlertOnExpiredCertificate allows to change a behavior on signature creation or
	// augmentation with an expired signing-certificate (notAfter is before the current time).
	// Validated the available POEs in case of existing signature augmentation. Nil provides a
	// possibility to skip check execution. Default: ExceptionOnStatusAlert - throw an
	// exception. Port of setAlertOnExpiredCertificate(...).
	SetAlertOnExpiredCertificate(alertOnExpiredCertificate alert.StatusAlert)

	// AlertOnExpiredCertificate returns the defined behavior on signature creation or
	// augmentation with an expired signing-certificate (notAfter is before the current time).
	// Validated the available POEs in case of existing signature augmentation. Port of
	// getAlertOnExpiredCertificate().
	AlertOnExpiredCertificate() alert.StatusAlert

	// SetAlertOnNotYetValidCertificate allows to change a behavior on signature creation with
	// a not yet valid signing-certificate (notBefore is after the current time). Nil provides
	// a possibility to skip check execution. Default: ExceptionOnStatusAlert - throw an
	// exception. Port of setAlertOnNotYetValidCertificate(...).
	SetAlertOnNotYetValidCertificate(alertOnNotYetValidCertificate alert.StatusAlert)

	// AlertOnNotYetValidCertificate returns the defined behavior on signature creation with a
	// not yet valid signing-certificate. Port of getAlertOnNotYetValidCertificate().
	AlertOnNotYetValidCertificate() alert.StatusAlert

	// SetAugmentationAlertOnHigherSignatureLevel allows to change the augmentation behaviour
	// for a signature of a higher level or a document containing a such signature. Nil
	// provides a possibility to skip check execution. Default: ExceptionOnStatusAlert - throw
	// an exception. Port of setAugmentationAlertOnHigherSignatureLevel(...).
	SetAugmentationAlertOnHigherSignatureLevel(augmentationAlertOnHigherSignatureLevel alert.StatusAlert)

	// AugmentationAlertOnHigherSignatureLevel returns the defined augmentation behaviour for a
	// signature of a higher level or a document containing a such signature. Port of
	// getAugmentationAlertOnHigherSignatureLevel().
	AugmentationAlertOnHigherSignatureLevel() alert.StatusAlert

	// SetAugmentationAlertOnSignatureWithoutCertificates allows to change the augmentation
	// behaviour for a signature without certificates. The alert is triggered when no
	// certificate is defined within the signature. Nil provides a possibility to skip check
	// execution. Default: ExceptionOnStatusAlert - throw an exception. Port of
	// setAugmentationAlertOnSignatureWithoutCertificates(...).
	SetAugmentationAlertOnSignatureWithoutCertificates(augmentationAlertOnSignatureWithoutCertificates alert.StatusAlert)

	// AugmentationAlertOnSignatureWithoutCertificates returns the defined augmentation
	// behaviour for a signature without certificates. Port of
	// getAugmentationAlertOnSignatureWithoutCertificates().
	AugmentationAlertOnSignatureWithoutCertificates() alert.StatusAlert

	// SetAugmentationAlertOnSelfSignedCertificateChains allows to change the augmentation
	// behaviour for a signature containing only self-signed certificate chains. The alert is
	// triggered when all used certificates are self-signed. Nil provides a possibility to skip
	// check execution. Default: ExceptionOnStatusAlert - throw an exception. Port of
	// setAugmentationAlertOnSelfSignedCertificateChains(...).
	SetAugmentationAlertOnSelfSignedCertificateChains(augmentationAlertOnSelfSignedCertificateChains alert.StatusAlert)

	// AugmentationAlertOnSelfSignedCertificateChains returns the defined augmentation
	// behaviour for a signature containing only self-signed certificate chains. Port of
	// getAugmentationAlertOnSelfSignedCertificateChains().
	AugmentationAlertOnSelfSignedCertificateChains() alert.StatusAlert

	// SetCheckRevocationForUntrustedChains allows enabling of revocation checking for
	// untrusted certificate chains. Default: FALSE (revocation data is not checked for
	// untrusted certificate chains). Port of setCheckRevocationForUntrustedChains(boolean).
	SetCheckRevocationForUntrustedChains(enable bool)

	// IsCheckRevocationForUntrustedChains returns true if revocation check is enabled for
	// untrusted certificate chains. Port of isCheckRevocationForUntrustedChains().
	IsCheckRevocationForUntrustedChains() bool
}
