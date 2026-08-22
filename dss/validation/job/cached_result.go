// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/CachedResult.java (DSS 6.5.RC1).
package job

// CachedResult is used to define a cached result for a single job. Java's empty marker
// interface becomes Go's empty interface; every cached result type (DownloadResult,
// ParsingResult, ValidationResult) trivially satisfies it.
type CachedResult any
