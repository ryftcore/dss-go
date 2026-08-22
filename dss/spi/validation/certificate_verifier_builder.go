// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/CertificateVerifierBuilder.java (DSS 6.5.RC1).
package validation

// CertificateVerifierBuilder builds a copy of CertificateVerifier.
type CertificateVerifierBuilder struct {
	// certificateVerifier is the CertificateVerifier to copy.
	certificateVerifier CertificateVerifier
}

// NewCertificateVerifierBuilder is the default constructor.
func NewCertificateVerifierBuilder(certificateVerifier CertificateVerifier) *CertificateVerifierBuilder {
	return &CertificateVerifierBuilder{certificateVerifier: certificateVerifier}
}

// BuildCompleteCopy builds a complete copy of the certificateVerifier. Port of
// buildCompleteCopy().
func (b *CertificateVerifierBuilder) BuildCompleteCopy() CertificateVerifier {
	copyVerifier := NewCommonCertificateVerifierSimple(true)
	if b.certificateVerifier != nil {
		copyVerifier.SetAIASource(b.certificateVerifier.AIASource())
		copyVerifier.SetCrlSource(b.certificateVerifier.CrlSource())
		copyVerifier.SetOcspSource(b.certificateVerifier.OcspSource())
		copyVerifier.SetRevocationDataLoadingStrategyFactory(b.certificateVerifier.RevocationDataLoadingStrategyFactory())
		copyVerifier.SetRevocationFallback(b.certificateVerifier.IsRevocationFallback())
		copyVerifier.SetRevocationDataVerifier(b.certificateVerifier.RevocationDataVerifier())
		copyVerifier.SetTimestampTokenVerifier(b.certificateVerifier.TimestampTokenVerifier())
		copyVerifier.SetTrustAnchorVerifier(b.certificateVerifier.TrustAnchorVerifier())
		copyVerifier.SetCheckRevocationForUntrustedChains(b.certificateVerifier.IsCheckRevocationForUntrustedChains())
		copyVerifier.SetAdjunctCertSourcesFromList(b.certificateVerifier.AdjunctCertSources())
		copyVerifier.SetTrustedCertSourcesFromList(b.certificateVerifier.TrustedCertSources())

		copyVerifier.SetAlertOnInvalidSignature(b.certificateVerifier.AlertOnInvalidSignature())
		copyVerifier.SetAlertOnInvalidTimestamp(b.certificateVerifier.AlertOnInvalidTimestamp())
		copyVerifier.SetAlertOnMissingRevocationData(b.certificateVerifier.AlertOnMissingRevocationData())
		copyVerifier.SetAlertOnNoRevocationAfterBestSignatureTime(b.certificateVerifier.AlertOnNoRevocationAfterBestSignatureTime())
		copyVerifier.SetAlertOnRevokedCertificate(b.certificateVerifier.AlertOnRevokedCertificate())
		copyVerifier.SetAlertOnUncoveredPOE(b.certificateVerifier.AlertOnUncoveredPOE())
		copyVerifier.SetAlertOnExpiredCertificate(b.certificateVerifier.AlertOnExpiredCertificate())
		copyVerifier.SetAlertOnNotYetValidCertificate(b.certificateVerifier.AlertOnNotYetValidCertificate())
		copyVerifier.SetAugmentationAlertOnSignatureWithoutCertificates(b.certificateVerifier.AugmentationAlertOnSignatureWithoutCertificates())
		copyVerifier.SetAugmentationAlertOnHigherSignatureLevel(b.certificateVerifier.AugmentationAlertOnHigherSignatureLevel())
		copyVerifier.SetAugmentationAlertOnSelfSignedCertificateChains(b.certificateVerifier.AugmentationAlertOnSelfSignedCertificateChains())
	}
	return copyVerifier
}

// BuildOfflineCopy builds a copy of the certificateVerifier by skipping the data sources, but
// keeping alerts. Port of buildOfflineCopy().
func (b *CertificateVerifierBuilder) BuildOfflineCopy() CertificateVerifier {
	offlineCertificateVerifier := NewCommonCertificateVerifierSimple(true)
	if b.certificateVerifier != nil {
		offlineCertificateVerifier.SetAdjunctCertSourcesFromList(b.certificateVerifier.AdjunctCertSources())
		offlineCertificateVerifier.SetTrustedCertSourcesFromList(b.certificateVerifier.TrustedCertSources())
		offlineCertificateVerifier.SetRevocationDataVerifier(b.certificateVerifier.RevocationDataVerifier())
		offlineCertificateVerifier.SetTimestampTokenVerifier(b.certificateVerifier.TimestampTokenVerifier())
		offlineCertificateVerifier.SetTrustAnchorVerifier(getTrustAnchorVerifierOfflineCopy(b.certificateVerifier.TrustAnchorVerifier()))

		// keep alerting
		offlineCertificateVerifier.SetAlertOnInvalidSignature(b.certificateVerifier.AlertOnInvalidSignature())
		offlineCertificateVerifier.SetAlertOnInvalidTimestamp(b.certificateVerifier.AlertOnInvalidTimestamp())
		offlineCertificateVerifier.SetAlertOnMissingRevocationData(b.certificateVerifier.AlertOnMissingRevocationData())
		offlineCertificateVerifier.SetAlertOnNoRevocationAfterBestSignatureTime(b.certificateVerifier.AlertOnNoRevocationAfterBestSignatureTime())
		offlineCertificateVerifier.SetAlertOnRevokedCertificate(b.certificateVerifier.AlertOnRevokedCertificate())
		offlineCertificateVerifier.SetAlertOnUncoveredPOE(b.certificateVerifier.AlertOnUncoveredPOE())
		offlineCertificateVerifier.SetAlertOnExpiredCertificate(b.certificateVerifier.AlertOnExpiredCertificate())
		offlineCertificateVerifier.SetAlertOnNotYetValidCertificate(b.certificateVerifier.AlertOnNotYetValidCertificate())
		offlineCertificateVerifier.SetAugmentationAlertOnSignatureWithoutCertificates(b.certificateVerifier.AugmentationAlertOnSignatureWithoutCertificates())
		offlineCertificateVerifier.SetAugmentationAlertOnHigherSignatureLevel(b.certificateVerifier.AugmentationAlertOnHigherSignatureLevel())
		offlineCertificateVerifier.SetAugmentationAlertOnSelfSignedCertificateChains(b.certificateVerifier.AugmentationAlertOnSelfSignedCertificateChains())
	}
	return offlineCertificateVerifier
}

// BuildOfflineAndSilentCopy builds a copy of the certificateVerifier by skipping the data
// sources and disabling alerts. Port of buildOfflineAndSilentCopy().
func (b *CertificateVerifierBuilder) BuildOfflineAndSilentCopy() CertificateVerifier {
	offlineCertificateVerifier := NewCommonCertificateVerifierSimple(true)
	if b.certificateVerifier != nil {
		offlineCertificateVerifier.SetAdjunctCertSourcesFromList(b.certificateVerifier.AdjunctCertSources())
		offlineCertificateVerifier.SetTrustedCertSourcesFromList(b.certificateVerifier.TrustedCertSources())
		offlineCertificateVerifier.SetRevocationDataVerifier(b.certificateVerifier.RevocationDataVerifier())
		offlineCertificateVerifier.SetTimestampTokenVerifier(b.certificateVerifier.TimestampTokenVerifier())
		offlineCertificateVerifier.SetTrustAnchorVerifier(getTrustAnchorVerifierOfflineCopy(b.certificateVerifier.TrustAnchorVerifier()))
	}
	// disable alerting
	offlineCertificateVerifier.SetAlertOnInvalidSignature(nil)
	offlineCertificateVerifier.SetAlertOnInvalidTimestamp(nil)
	offlineCertificateVerifier.SetAlertOnMissingRevocationData(nil)
	offlineCertificateVerifier.SetAlertOnNoRevocationAfterBestSignatureTime(nil)
	offlineCertificateVerifier.SetAlertOnRevokedCertificate(nil)
	offlineCertificateVerifier.SetAlertOnUncoveredPOE(nil)
	offlineCertificateVerifier.SetAlertOnExpiredCertificate(nil)
	offlineCertificateVerifier.SetAlertOnNotYetValidCertificate(nil)
	offlineCertificateVerifier.SetAugmentationAlertOnSignatureWithoutCertificates(nil)
	offlineCertificateVerifier.SetAugmentationAlertOnHigherSignatureLevel(nil)
	offlineCertificateVerifier.SetAugmentationAlertOnSelfSignedCertificateChains(nil)
	return offlineCertificateVerifier
}

// getTrustAnchorVerifierOfflineCopy is the port of the private
// getTrustAnchorVerifierOfflineCopy(TrustAnchorVerifier) method.
func getTrustAnchorVerifierOfflineCopy(originalTrustAnchorVerifier *TrustAnchorVerifier) *TrustAnchorVerifier {
	trustAnchorVerifier := NewEmptyTrustAnchorVerifier()
	trustAnchorVerifier.SetUseSunsetDate(false) // set to FALSE for offline processing
	if originalTrustAnchorVerifier != nil {
		trustAnchorVerifier.SetTrustedCertificateSource(originalTrustAnchorVerifier.TrustedCertificateSource())
		trustAnchorVerifier.SetAcceptRevocationUntrustedCertificateChains(originalTrustAnchorVerifier.IsAcceptRevocationUntrustedCertificateChains())
		trustAnchorVerifier.SetAcceptTimestampUntrustedCertificateChains(originalTrustAnchorVerifier.IsAcceptTimestampUntrustedCertificateChains())
	}
	return trustAnchorVerifier
}

// BuildCompleteCopyForValidation builds a local copy of a CertificateVerifier used by a
// signature validation process. Port of buildCompleteCopyForValidation().
func (b *CertificateVerifierBuilder) BuildCompleteCopyForValidation() CertificateVerifier {
	copyVerifier := b.BuildCompleteCopy()
	copyVerifier.SetRevocationFallback(true)
	return copyVerifier
}
