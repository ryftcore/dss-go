// Ported from dss-enumerations/.../ValidationTime.java (DSS 6.5.RC1).
package enumerations

// ValidationTime defines the used validation time types.
type ValidationTime string

const (
	// ValidationTimeCertificateIssuanceTime is the not-before date of
	// the certificate.
	ValidationTimeCertificateIssuanceTime ValidationTime = "CERTIFICATE_ISSUANCE_TIME"
	// ValidationTimeBESTSignatureTime is the lowest time at which
	// there exists a POE for the signature.
	ValidationTimeBESTSignatureTime ValidationTime = "BEST_SIGNATURE_TIME"
	// ValidationTimeValidationTime is the date of the validation.
	ValidationTimeValidationTime ValidationTime = "VALIDATION_TIME"
	// ValidationTimeTimestampGenerationTime is the time of the
	// timestamp generation.
	ValidationTimeTimestampGenerationTime ValidationTime = "TIMESTAMP_GENERATION_TIME"
	// ValidationTimeTimestampPOETime is the time when the timestamp is
	// proved to exist.
	ValidationTimeTimestampPOETime ValidationTime = "TIMESTAMP_POE_TIME"
)

// ValidationTimeValues returns all constants in declaration order.
func ValidationTimeValues() []ValidationTime {
	return []ValidationTime{
		ValidationTimeCertificateIssuanceTime,
		ValidationTimeBESTSignatureTime,
		ValidationTimeValidationTime,
		ValidationTimeTimestampGenerationTime,
		ValidationTimeTimestampPOETime,
	}
}
