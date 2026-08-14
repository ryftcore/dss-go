// Ported from dss-enumerations/.../ValidationTime.java (DSS 6.5.RC1).
package enumerations

// ValidationTime defines the used validation time types.
type ValidationTime string

const (
	// ValidationTime_CERTIFICATE_ISSUANCE_TIME is the not-before date of
	// the certificate.
	ValidationTime_CERTIFICATE_ISSUANCE_TIME ValidationTime = "CERTIFICATE_ISSUANCE_TIME"
	// ValidationTime_BEST_SIGNATURE_TIME is the lowest time at which
	// there exists a POE for the signature.
	ValidationTime_BEST_SIGNATURE_TIME ValidationTime = "BEST_SIGNATURE_TIME"
	// ValidationTime_VALIDATION_TIME is the date of the validation.
	ValidationTime_VALIDATION_TIME ValidationTime = "VALIDATION_TIME"
	// ValidationTime_TIMESTAMP_GENERATION_TIME is the time of the
	// timestamp generation.
	ValidationTime_TIMESTAMP_GENERATION_TIME ValidationTime = "TIMESTAMP_GENERATION_TIME"
	// ValidationTime_TIMESTAMP_POE_TIME is the time when the timestamp is
	// proved to exist.
	ValidationTime_TIMESTAMP_POE_TIME ValidationTime = "TIMESTAMP_POE_TIME"
)

// ValidationTimeValues returns all constants in declaration order.
func ValidationTimeValues() []ValidationTime {
	return []ValidationTime{
		ValidationTime_CERTIFICATE_ISSUANCE_TIME,
		ValidationTime_BEST_SIGNATURE_TIME,
		ValidationTime_VALIDATION_TIME,
		ValidationTime_TIMESTAMP_GENERATION_TIME,
		ValidationTime_TIMESTAMP_POE_TIME,
	}
}
