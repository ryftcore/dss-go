// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/ValidationAlerter.java (DSS 6.5.RC1).
//
// This interface uses eu.europa.esig.dss.spi.validation.ValidationContext to perform validation
// and executes alerts based on the validation result.
package validation

import "github.com/ryftcore/dss-go/dss/model"

// ValidationAlerter is used with ValidationContext to perform validation and execute alerts
// based on the validation result.
type ValidationAlerter interface {
	// AssertAllRequiredRevocationDataPresent verifies if all processed certificates have a
	// revocation data. The behavior is configured with
	// CertificateVerifier.SetAlertOnMissingRevocationData(alert.StatusAlert). Port of
	// assertAllRequiredRevocationDataPresent().
	AssertAllRequiredRevocationDataPresent()

	// AssertAllPOECoveredByRevocationData verifies if all POE (timestamp tokens) are covered by
	// a revocation data. The behavior is configured with
	// CertificateVerifier.SetAlertOnUncoveredPOE(alert.StatusAlert). Port of
	// assertAllPOECoveredByRevocationData().
	AssertAllPOECoveredByRevocationData()

	// AssertAllTimestampsValid verifies if all processed timestamps are valid and intact. The
	// behavior is configured with CertificateVerifier.SetAlertOnInvalidTimestamp(alert.StatusAlert).
	// Port of assertAllTimestampsValid().
	AssertAllTimestampsValid()

	// AssertCertificateNotRevoked verifies if the certificate is not revoked. The behavior is
	// configured with CertificateVerifier.SetAlertOnRevokedCertificate(alert.StatusAlert). Port
	// of assertCertificateNotRevoked(CertificateToken).
	AssertCertificateNotRevoked(certificateToken *model.CertificateToken)

	// AssertAllSignatureCertificatesNotRevoked verifies recursively whether none of the
	// signature's certificate chain certificates are revoked. The behavior is configured with
	// CertificateVerifier.SetAlertOnRevokedCertificate(alert.StatusAlert). Port of
	// assertAllSignatureCertificatesNotRevoked().
	AssertAllSignatureCertificatesNotRevoked()

	// AssertAllSignatureCertificateHaveFreshRevocationData verifies whether for all signature's
	// certificate chain certificates there is a fresh revocation data, after the earliest
	// available timestamp token production time. The behavior is configured with
	// CertificateVerifier.SetAlertOnNoRevocationAfterBestSignatureTime(alert.StatusAlert). Port
	// of assertAllSignatureCertificateHaveFreshRevocationData().
	AssertAllSignatureCertificateHaveFreshRevocationData()

	// AssertAllSignaturesNotExpired verifies whether all signatures added to the
	// ValidationContext are not yet expired. The behavior is configured with
	// CertificateVerifier.SetAlertOnExpiredCertificate(alert.StatusAlert). Port of
	// assertAllSignaturesNotExpired().
	AssertAllSignaturesNotExpired()

	// AssertCertificateNotExpired verifies whether the certificate token is not yet expired.
	// The behavior is configured with CertificateVerifier.SetAlertOnExpiredCertificate(alert.StatusAlert).
	// Port of assertCertificateNotExpired(CertificateToken).
	AssertCertificateNotExpired(certificateToken *model.CertificateToken)

	// AssertAllSignaturesAreYetValid verifies whether all signatures added to the
	// ValidationContext have been produced with yet valid certificates. The behavior is
	// configured with CertificateVerifier.SetAlertOnNotYetValidCertificate(alert.StatusAlert).
	// Port of assertAllSignaturesAreYetValid().
	AssertAllSignaturesAreYetValid()

	// AssertCertificateIsYetValid verifies whether the certificate token is yet valid. The
	// behavior is configured with CertificateVerifier.SetAlertOnNotYetValidCertificate(alert.StatusAlert).
	// Port of assertCertificateIsYetValid(CertificateToken).
	AssertCertificateIsYetValid(certificateToken *model.CertificateToken)
}
