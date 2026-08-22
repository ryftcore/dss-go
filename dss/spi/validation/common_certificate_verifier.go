// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/CommonCertificateVerifier.java (DSS 6.5.RC1).
//
// This class provides the different sources used to verify the status of a certificate using
// the trust model. There are four different types of sources to be defined:
//   - Trusted certificates source;
//   - Adjunct certificates source (not trusted);
//   - OCSP source;
//   - CRL source;
//   - AIA source to give access to the certificates through AIA.
//
// Java's slf4j logging (LOG.info on construction, LOG.warn in assertNotTrusted) has no Go
// equivalent and is not ported.
package validation

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/x509/aia"
)

// CommonCertificateVerifier is the default CertificateVerifier implementation.
type CommonCertificateVerifier struct {
	// trustedCertSources contains the reference to multiple trusted certificate sources.
	// These sources are fixed, it means that the same sources are used for different
	// validations.
	trustedCertSources *spi.ListCertificateSource

	// adjunctCertSources contains the reference to arbitrary certificate source, can contain a
	// trust store, or any intermediate certificates.
	adjunctCertSources *spi.ListCertificateSource

	// ocspSource contains the reference to the OCSPSource.
	ocspSource spi.RevocationSource[revocation.OCSP]

	// crlSource contains the reference to the CRLSource.
	crlSource spi.RevocationSource[revocation.CRL]

	// revocationDataLoadingStrategyFactory creates a revocation data loading strategy used to
	// fetch OCSP or CRL for validating certificates.
	//
	// Default: OCSPFirstRevocationDataLoadingStrategyFactory is used to extract OCSP token
	// first and CRL after.
	revocationDataLoadingStrategyFactory RevocationDataLoadingStrategyFactory

	// revocationDataVerifier verifies the validity of retrieved revocation data (used to
	// evaluate if a new revocation should be requested).
	revocationDataVerifier *RevocationDataVerifier

	// revocationFallback defines whether a revocation data failed a verification still shall
	// be returned to the validation process, when no valid revocation has been obtained.
	revocationFallback bool

	// timestampTokenVerifier verifies validity of a timestamp tokens encapsulated within the
	// signature document.
	timestampTokenVerifier *TimestampTokenVerifier

	// trustAnchorVerifier verifies validity of trust anchors.
	trustAnchorVerifier *TrustAnchorVerifier

	// aiaSource is used to download a certificate's issuer by the AIA URI(s) defining within a
	// certificate.
	aiaSource aia.AIASource

	// alertOnInvalidSignature sets the behavior to follow in case of invalid signature
	// (augmentation process).
	//
	// Default: ExceptionOnStatusAlert - throw the exception.
	alertOnInvalidSignature alert.StatusAlert

	// alertOnInvalidTimestamp sets the behavior to follow in case of invalid timestamp
	// (augmentation process).
	//
	// Default: ExceptionOnStatusAlert - throw the exception.
	alertOnInvalidTimestamp alert.StatusAlert

	// alertOnMissingRevocationData sets the behavior to follow in case of missing revocation
	// data (augmentation process).
	//
	// Default: ExceptionOnStatusAlert - throw the exception.
	alertOnMissingRevocationData alert.StatusAlert

	// alertOnRevokedCertificate sets the behavior to follow in case of revoked certificate
	// (augmentation process).
	//
	// Default: ExceptionOnStatusAlert - throw the exception.
	alertOnRevokedCertificate alert.StatusAlert

	// alertOnNoRevocationAfterBestSignatureTime sets the behavior to follow in case of no
	// revocation data issued after the bestSignatureTime (augmentation process).
	//
	// Default: LogOnStatusAlert - log a warning message.
	alertOnNoRevocationAfterBestSignatureTime alert.StatusAlert

	// alertOnUncoveredPOE sets the behavior to follow in case of missing revocation data for a
	// POE.
	//
	// Default: LogOnStatusAlert - log a warning message.
	alertOnUncoveredPOE alert.StatusAlert

	// alertOnExpiredCertificate defines a behavior on signature creation with an expired
	// signing-certificate.
	//
	// Default: ExceptionOnStatusAlert - throw the exception.
	alertOnExpiredCertificate alert.StatusAlert

	// alertOnNotYetValidCertificate defines a behavior on signature creation with a not yet
	// valid certificate.
	//
	// Default: ExceptionOnStatusAlert - throw the exception.
	alertOnNotYetValidCertificate alert.StatusAlert

	// augmentationAlertOnHigherSignatureLevel defines behavior on signature augmentation when
	// a signature with a higher level is present, than the target one.
	//
	// Default: ExceptionOnStatusAlert - throw the exception.
	augmentationAlertOnHigherSignatureLevel alert.StatusAlert

	// augmentationAlertOnSignatureWithoutCertificates defines behavior on augmentation for a
	// signature not containing certificates.
	//
	// Default: ExceptionOnStatusAlert - throw the exception.
	augmentationAlertOnSignatureWithoutCertificates alert.StatusAlert

	// augmentationAlertOnSelfSignedCertificateChains defines behavior on augmentation for a
	// signature containing only self-signed certificate chains.
	//
	// Default: ExceptionOnStatusAlert - throw the exception.
	augmentationAlertOnSelfSignedCertificateChains alert.StatusAlert

	// checkRevocationForUntrustedChains sets the behavior to follow for revocation retrieving
	// in case of untrusted certificate chains.
	//
	// Default: false (revocation are not checked in case of certificates issued from an
	// unsure source).
	checkRevocationForUntrustedChains bool
}

// NewCommonCertificateVerifier is the default constructor. The DataLoader is created to allow
// the retrieval of certificates through AIA. Port of the no-arg constructor (this(false)).
func NewCommonCertificateVerifier() *CommonCertificateVerifier {
	return newCommonCertificateVerifier(false)
}

// NewCommonCertificateVerifierSimple allows creating of CommonCertificateVerifier without
// DataLoader. It means that only a -B profile signature can be created. simpleCreationOnly:
// if true the CommonCertificateVerifier will not contain AIASource. Port of the
// CommonCertificateVerifier(boolean) constructor.
func NewCommonCertificateVerifierSimple(simpleCreationOnly bool) *CommonCertificateVerifier {
	return newCommonCertificateVerifier(simpleCreationOnly)
}

func newCommonCertificateVerifier(simpleCreationOnly bool) *CommonCertificateVerifier {
	v := &CommonCertificateVerifier{
		trustedCertSources:                              spi.NewListCertificateSource(),
		adjunctCertSources:                              spi.NewListCertificateSource(),
		revocationDataLoadingStrategyFactory:            NewOCSPFirstRevocationDataLoadingStrategyFactory(),
		revocationFallback:                              false,
		alertOnInvalidSignature:                         alert.NewExceptionOnStatusAlert(),
		alertOnInvalidTimestamp:                         alert.NewExceptionOnStatusAlert(),
		alertOnMissingRevocationData:                    alert.NewExceptionOnStatusAlert(),
		alertOnRevokedCertificate:                       alert.NewExceptionOnStatusAlert(),
		alertOnNoRevocationAfterBestSignatureTime:       alert.NewLogOnStatusAlert(),
		alertOnUncoveredPOE:                             alert.NewLogOnStatusAlert(),
		alertOnExpiredCertificate:                       alert.NewExceptionOnStatusAlert(),
		alertOnNotYetValidCertificate:                   alert.NewExceptionOnStatusAlert(),
		augmentationAlertOnHigherSignatureLevel:         alert.NewExceptionOnStatusAlert(),
		augmentationAlertOnSignatureWithoutCertificates: alert.NewExceptionOnStatusAlert(),
		augmentationAlertOnSelfSignedCertificateChains:  alert.NewExceptionOnStatusAlert(),
		checkRevocationForUntrustedChains:               false,
	}
	if !simpleCreationOnly {
		v.aiaSource = aia.NewDefaultAIASource()
	}
	return v
}

// CrlSource returns the CRL source associated with this verifier. Port of getCrlSource().
func (v *CommonCertificateVerifier) CrlSource() spi.RevocationSource[revocation.CRL] {
	return v.crlSource
}

// SetCrlSource defines the source of CRL used by this class. Port of setCrlSource(...).
func (v *CommonCertificateVerifier) SetCrlSource(crlSource spi.RevocationSource[revocation.CRL]) {
	v.crlSource = crlSource
}

// OcspSource returns the OCSP source associated with this verifier. Port of getOcspSource().
func (v *CommonCertificateVerifier) OcspSource() spi.RevocationSource[revocation.OCSP] {
	return v.ocspSource
}

// SetOcspSource defines the source of OCSP used by this class. Port of setOcspSource(...).
func (v *CommonCertificateVerifier) SetOcspSource(ocspSource spi.RevocationSource[revocation.OCSP]) {
	v.ocspSource = ocspSource
}

// RevocationDataLoadingStrategyFactory returns a factory used to create revocation data
// loading strategy associated with this verifier. Port of
// getRevocationDataLoadingStrategyFactory().
func (v *CommonCertificateVerifier) RevocationDataLoadingStrategyFactory() RevocationDataLoadingStrategyFactory {
	return v.revocationDataLoadingStrategyFactory
}

// SetRevocationDataLoadingStrategyFactory creates a strategy used to fetch OCSP or CRL for
// certificate validation. Port of setRevocationDataLoadingStrategyFactory(...).
//
// Panics with the Java message when revocationDataLoadingStrategyFactory is nil
// (Objects.requireNonNull(..., "RevocationDataLoadingStrategyFactory shall be defined!")).
func (v *CommonCertificateVerifier) SetRevocationDataLoadingStrategyFactory(revocationDataLoadingStrategyFactory RevocationDataLoadingStrategyFactory) {
	if revocationDataLoadingStrategyFactory == nil {
		panic("RevocationDataLoadingStrategyFactory shall be defined!")
	}
	v.revocationDataLoadingStrategyFactory = revocationDataLoadingStrategyFactory
}

// RevocationDataVerifier returns a RevocationDataVerifier associated with this verifier. Port
// of getRevocationDataVerifier().
func (v *CommonCertificateVerifier) RevocationDataVerifier() *RevocationDataVerifier {
	return v.revocationDataVerifier
}

// SetRevocationDataVerifier sets the RevocationDataVerifier used to validate acceptance of
// the retrieved revocation data. Port of setRevocationDataVerifier(...).
func (v *CommonCertificateVerifier) SetRevocationDataVerifier(revocationDataVerifier *RevocationDataVerifier) {
	v.revocationDataVerifier = revocationDataVerifier
}

// IsRevocationFallback returns whether revocation data still shall be returned if validation
// of requested revocation data failed. Port of isRevocationFallback().
func (v *CommonCertificateVerifier) IsRevocationFallback() bool {
	return v.revocationFallback
}

// SetRevocationFallback sets whether a revocation data still have to be returned to the
// validation process, in case validation of obtained revocation data has failed. Port of
// setRevocationFallback(boolean).
func (v *CommonCertificateVerifier) SetRevocationFallback(revocationFallback bool) {
	v.revocationFallback = revocationFallback
}

// TimestampTokenVerifier returns a TimestampTokenVerifier associated with this verifier. Port
// of getTimestampTokenVerifier().
func (v *CommonCertificateVerifier) TimestampTokenVerifier() *TimestampTokenVerifier {
	return v.timestampTokenVerifier
}

// SetTimestampTokenVerifier sets the TimestampTokenVerifier used to validate acceptance of
// the timestamp tokens encapsulated within the signature. Port of
// setTimestampTokenVerifier(...).
func (v *CommonCertificateVerifier) SetTimestampTokenVerifier(timestampTokenVerifier *TimestampTokenVerifier) {
	v.timestampTokenVerifier = timestampTokenVerifier
}

// TrustAnchorVerifier returns a TrustAnchorVerifier associated with this verifier. Port of
// getTrustAnchorVerifier().
func (v *CommonCertificateVerifier) TrustAnchorVerifier() *TrustAnchorVerifier {
	return v.trustAnchorVerifier
}

// SetTrustAnchorVerifier sets the TrustAnchorVerifier used to validate acceptance of the
// trust anchors. Port of setTrustAnchorVerifier(...).
func (v *CommonCertificateVerifier) SetTrustAnchorVerifier(trustAnchorVerifier *TrustAnchorVerifier) {
	v.trustAnchorVerifier = trustAnchorVerifier
}

// TrustedCertSources returns the trusted certificate sources associated with this verifier.
// Port of getTrustedCertSources().
func (v *CommonCertificateVerifier) TrustedCertSources() *spi.ListCertificateSource {
	return v.trustedCertSources
}

// SetTrustedCertSources sets multiple trusted certificate sources. Port of
// setTrustedCertSources(CertificateSource...).
//
// Panics with the Java message when a certSources entry is not a trusted source
// (UnsupportedOperationException).
func (v *CommonCertificateVerifier) SetTrustedCertSources(certSources ...spi.CertificateSource) {
	v.trustedCertSources = spi.NewListCertificateSource()
	v.AddTrustedCertSources(certSources...)
}

// AddTrustedCertSources adds trusted certificate sources to an existing list of trusted
// certificate sources. Port of addTrustedCertSources(CertificateSource...).
//
// Panics with the Java message when a certSources entry is not a trusted source
// (UnsupportedOperationException).
func (v *CommonCertificateVerifier) AddTrustedCertSources(certSources ...spi.CertificateSource) {
	for _, certificateSource := range certSources {
		if certificateSource.CertificateSourceType().IsTrusted() {
			v.trustedCertSources.Add(certificateSource)
		} else {
			panic(fmt.Sprintf("The certificateSource with type [%s] is not allowed in the trustedCertSources. Please, "+
				"use CertificateSource with a type TRUSTED_STORE or TRUSTED_LIST.", certificateSource.CertificateSourceType()))
		}
	}
}

// SetTrustedCertSourcesFromList sets a list of trusted certificate sources. Port of the
// setTrustedCertSources(ListCertificateSource) overload.
//
// Panics with the Java message when trustedListCertificateSource contains a non-trusted
// source (UnsupportedOperationException).
func (v *CommonCertificateVerifier) SetTrustedCertSourcesFromList(trustedListCertificateSource *spi.ListCertificateSource) {
	if trustedListCertificateSource == nil {
		v.trustedCertSources = spi.NewListCertificateSource()
	} else if trustedListCertificateSource.AreAllCertSourcesTrusted() {
		v.trustedCertSources = trustedListCertificateSource
	} else {
		panic("The trusted ListCertificateSource must contain only trusted sources " +
			"with a type TRUSTED_STORE or TRUSTED_LIST.")
	}
}

// AdjunctCertSources returns the list of adjunct certificate sources assigned to this
// verifier. Port of getAdjunctCertSources().
func (v *CommonCertificateVerifier) AdjunctCertSources() *spi.ListCertificateSource {
	return v.adjunctCertSources
}

// SetAdjunctCertSources sets multiple adjunct certificate sources. Port of
// setAdjunctCertSources(CertificateSource...).
func (v *CommonCertificateVerifier) SetAdjunctCertSources(certSources ...spi.CertificateSource) {
	v.adjunctCertSources = spi.NewListCertificateSource()
	v.AddAdjunctCertSources(certSources...)
}

// AddAdjunctCertSources adds adjunct certificate sources to an existing list of adjunct
// certificate sources. Port of addAdjunctCertSources(CertificateSource...).
func (v *CommonCertificateVerifier) AddAdjunctCertSources(certSources ...spi.CertificateSource) {
	for _, certificateSource := range certSources {
		assertNotTrusted(certificateSource)
		v.adjunctCertSources.Add(certificateSource)
	}
}

// SetAdjunctCertSourcesFromList sets a list of adjunct certificate sources. Port of the
// setAdjunctCertSources(ListCertificateSource) overload.
func (v *CommonCertificateVerifier) SetAdjunctCertSourcesFromList(adjunctListCertificateSource *spi.ListCertificateSource) {
	if adjunctListCertificateSource == nil {
		adjunctListCertificateSource = spi.NewListCertificateSource()
	}
	for _, certificateSource := range adjunctListCertificateSource.Sources() {
		assertNotTrusted(certificateSource)
	}
	v.adjunctCertSources = adjunctListCertificateSource
}

// assertNotTrusted logs (dropped, see file header) when adjunctCertificateSource is trusted.
// Port of assertNotTrusted(...).
func assertNotTrusted(adjunctCertificateSource spi.CertificateSource) {
	_ = adjunctCertificateSource
}

// AIASource gets the AIASource used to load a certificate's issuer by defined AIA URI(s)
// within the token. Port of getAIASource().
func (v *CommonCertificateVerifier) AIASource() aia.AIASource {
	return v.aiaSource
}

// SetAIASource sets the AIASource used to load a certificate's issuer by defined AIA URI(s)
// within the token. Port of setAIASource(...).
func (v *CommonCertificateVerifier) SetAIASource(aiaSource aia.AIASource) {
	v.aiaSource = aiaSource
}

// SetAlertOnInvalidSignature allows to change the behavior on invalid signature. Port of
// setAlertOnInvalidSignature(...).
func (v *CommonCertificateVerifier) SetAlertOnInvalidSignature(alertOnInvalidSignature alert.StatusAlert) {
	v.alertOnInvalidSignature = alertOnInvalidSignature
}

// AlertOnInvalidSignature returns the defined execution behaviour on invalid signature. Port
// of getAlertOnInvalidSignature().
func (v *CommonCertificateVerifier) AlertOnInvalidSignature() alert.StatusAlert {
	return v.alertOnInvalidSignature
}

// SetAlertOnInvalidTimestamp allows to change the behavior on invalid timestamp. Port of
// setAlertOnInvalidTimestamp(...).
func (v *CommonCertificateVerifier) SetAlertOnInvalidTimestamp(alertOnInvalidTimestamp alert.StatusAlert) {
	v.alertOnInvalidTimestamp = alertOnInvalidTimestamp
}

// AlertOnInvalidTimestamp returns the defined execution behaviour on invalid timestamp. Port
// of getAlertOnInvalidTimestamp().
func (v *CommonCertificateVerifier) AlertOnInvalidTimestamp() alert.StatusAlert {
	return v.alertOnInvalidTimestamp
}

// SetAlertOnMissingRevocationData allows to change the behavior on missing revocation data.
// Port of setAlertOnMissingRevocationData(...).
func (v *CommonCertificateVerifier) SetAlertOnMissingRevocationData(alertOnMissingRevocationData alert.StatusAlert) {
	v.alertOnMissingRevocationData = alertOnMissingRevocationData
}

// AlertOnMissingRevocationData returns the defined execution behaviour on missing revocation
// data. Port of getAlertOnMissingRevocationData().
func (v *CommonCertificateVerifier) AlertOnMissingRevocationData() alert.StatusAlert {
	return v.alertOnMissingRevocationData
}

// SetAlertOnUncoveredPOE allows to change the behavior on uncovered POE (timestamp). Port of
// setAlertOnUncoveredPOE(...).
func (v *CommonCertificateVerifier) SetAlertOnUncoveredPOE(alertOnUncoveredPOE alert.StatusAlert) {
	v.alertOnUncoveredPOE = alertOnUncoveredPOE
}

// AlertOnUncoveredPOE returns the defined execution behaviour on uncovered POE (timestamp).
// Port of getAlertOnUncoveredPOE().
func (v *CommonCertificateVerifier) AlertOnUncoveredPOE() alert.StatusAlert {
	return v.alertOnUncoveredPOE
}

// SetAlertOnRevokedCertificate allows to change the behavior on revoked certificates. Port of
// setAlertOnRevokedCertificate(...).
func (v *CommonCertificateVerifier) SetAlertOnRevokedCertificate(alertOnRevokedCertificate alert.StatusAlert) {
	v.alertOnRevokedCertificate = alertOnRevokedCertificate
}

// AlertOnRevokedCertificate returns the defined execution behaviour on revoked certificate.
// Port of getAlertOnRevokedCertificate().
func (v *CommonCertificateVerifier) AlertOnRevokedCertificate() alert.StatusAlert {
	return v.alertOnRevokedCertificate
}

// SetAlertOnNoRevocationAfterBestSignatureTime allows to change the behavior on revocation
// data issued after a control time. Port of
// setAlertOnNoRevocationAfterBestSignatureTime(...).
func (v *CommonCertificateVerifier) SetAlertOnNoRevocationAfterBestSignatureTime(alertOnNoRevocationAfterBestSignatureTime alert.StatusAlert) {
	v.alertOnNoRevocationAfterBestSignatureTime = alertOnNoRevocationAfterBestSignatureTime
}

// AlertOnNoRevocationAfterBestSignatureTime returns the defined execution behaviour if no
// revocation data obtained with an issuance time after the bestSignatureTime. Port of
// getAlertOnNoRevocationAfterBestSignatureTime().
func (v *CommonCertificateVerifier) AlertOnNoRevocationAfterBestSignatureTime() alert.StatusAlert {
	return v.alertOnNoRevocationAfterBestSignatureTime
}

// SetAlertOnExpiredCertificate allows to change a behavior on signature creation or
// augmentation with an expired signing-certificate. Port of
// setAlertOnExpiredCertificate(...).
func (v *CommonCertificateVerifier) SetAlertOnExpiredCertificate(alertOnExpiredCertificate alert.StatusAlert) {
	v.alertOnExpiredCertificate = alertOnExpiredCertificate
}

// AlertOnExpiredCertificate returns the defined behavior on signature creation or
// augmentation with an expired signing-certificate. Port of getAlertOnExpiredCertificate().
func (v *CommonCertificateVerifier) AlertOnExpiredCertificate() alert.StatusAlert {
	return v.alertOnExpiredCertificate
}

// SetAlertOnNotYetValidCertificate allows to change a behavior on signature creation with a
// not yet valid signing-certificate. Port of setAlertOnNotYetValidCertificate(...).
func (v *CommonCertificateVerifier) SetAlertOnNotYetValidCertificate(alertOnNotYetValidCertificate alert.StatusAlert) {
	v.alertOnNotYetValidCertificate = alertOnNotYetValidCertificate
}

// AlertOnNotYetValidCertificate returns the defined behavior on signature creation with a not
// yet valid signing-certificate. Port of getAlertOnNotYetValidCertificate().
func (v *CommonCertificateVerifier) AlertOnNotYetValidCertificate() alert.StatusAlert {
	return v.alertOnNotYetValidCertificate
}

// SetAugmentationAlertOnHigherSignatureLevel allows to change the augmentation behaviour for
// a signature of a higher level or a document containing a such signature. Port of
// setAugmentationAlertOnHigherSignatureLevel(...).
func (v *CommonCertificateVerifier) SetAugmentationAlertOnHigherSignatureLevel(augmentationAlertOnHigherSignatureLevel alert.StatusAlert) {
	v.augmentationAlertOnHigherSignatureLevel = augmentationAlertOnHigherSignatureLevel
}

// AugmentationAlertOnHigherSignatureLevel returns the defined augmentation behaviour for a
// signature of a higher level or a document containing a such signature. Port of
// getAugmentationAlertOnHigherSignatureLevel().
func (v *CommonCertificateVerifier) AugmentationAlertOnHigherSignatureLevel() alert.StatusAlert {
	return v.augmentationAlertOnHigherSignatureLevel
}

// SetAugmentationAlertOnSignatureWithoutCertificates allows to change the augmentation
// behaviour for a signature without certificates. Port of
// setAugmentationAlertOnSignatureWithoutCertificates(...).
func (v *CommonCertificateVerifier) SetAugmentationAlertOnSignatureWithoutCertificates(augmentationAlertOnSignatureWithoutCertificates alert.StatusAlert) {
	v.augmentationAlertOnSignatureWithoutCertificates = augmentationAlertOnSignatureWithoutCertificates
}

// AugmentationAlertOnSignatureWithoutCertificates returns the defined augmentation behaviour
// for a signature without certificates. Port of
// getAugmentationAlertOnSignatureWithoutCertificates().
func (v *CommonCertificateVerifier) AugmentationAlertOnSignatureWithoutCertificates() alert.StatusAlert {
	return v.augmentationAlertOnSignatureWithoutCertificates
}

// SetAugmentationAlertOnSelfSignedCertificateChains allows to change the augmentation
// behaviour for a signature containing only self-signed certificate chains. Port of
// setAugmentationAlertOnSelfSignedCertificateChains(...).
func (v *CommonCertificateVerifier) SetAugmentationAlertOnSelfSignedCertificateChains(augmentationAlertOnSelfSignedCertificateChains alert.StatusAlert) {
	v.augmentationAlertOnSelfSignedCertificateChains = augmentationAlertOnSelfSignedCertificateChains
}

// AugmentationAlertOnSelfSignedCertificateChains returns the defined augmentation behaviour
// for a signature containing only self-signed certificate chains. Port of
// getAugmentationAlertOnSelfSignedCertificateChains().
func (v *CommonCertificateVerifier) AugmentationAlertOnSelfSignedCertificateChains() alert.StatusAlert {
	return v.augmentationAlertOnSelfSignedCertificateChains
}

// SetCheckRevocationForUntrustedChains allows enabling of revocation checking for untrusted
// certificate chains. Port of setCheckRevocationForUntrustedChains(boolean).
func (v *CommonCertificateVerifier) SetCheckRevocationForUntrustedChains(enable bool) {
	v.checkRevocationForUntrustedChains = enable
}

// IsCheckRevocationForUntrustedChains returns true if revocation check is enabled for
// untrusted certificate chains. Port of isCheckRevocationForUntrustedChains().
func (v *CommonCertificateVerifier) IsCheckRevocationForUntrustedChains() bool {
	return v.checkRevocationForUntrustedChains
}

// compile-time assertion: a CommonCertificateVerifier is a CertificateVerifier.
var _ CertificateVerifier = (*CommonCertificateVerifier)(nil)
