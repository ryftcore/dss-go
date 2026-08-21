// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/download/DownloadTask.java (DSS 6.5.RC1).
package job

// DownloadTask performs a download job. Java's Supplier<DownloadResult>.get() (a checked-
// exception-free functional interface) becomes a Go func returning (DownloadResult, error):
// callers in this package (see AbstractAnalysis.download) always expect a possible error
// from a download task, so the error is threaded explicitly rather than panicking.
type DownloadTask interface {
	// Get performs the download and returns its result. Port of get() (java.util.function.
	// Supplier<DownloadResult>).
	Get() (DownloadResult, error)
}
