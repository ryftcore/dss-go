// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/validation/ValidationTask.java (DSS 6.5.RC1).
package job

// ValidationTask performs a signature validation job. See DownloadTask for the (T, error)
// return-shape rationale.
type ValidationTask interface {
	// Get performs the validation and returns its result. Port of get() (java.util.function.
	// Supplier<ValidationResult>).
	Get() (ValidationResult, error)
}
