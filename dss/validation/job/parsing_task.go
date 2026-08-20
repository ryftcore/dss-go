// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/parsing/ParsingTask.java (DSS 6.5.RC1).
package job

// ParsingTask performs a parsing job. See DownloadTask for the (T, error) return-shape
// rationale.
type ParsingTask interface {
	// Get performs the parsing and returns its result. Port of get() (java.util.function.
	// Supplier<ParsingResult>).
	Get() (ParsingResult, error)
}
