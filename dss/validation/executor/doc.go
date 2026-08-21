// Package executor ports dss-validation's
// eu.europa.esig.dss.validation.executor package: the process executors
// that run the EN 319 102-1 validation process (dss/validation/process) over
// already-built diagnostic data and turn the result into the three report
// flavors consumers read - SimpleReport, DetailedReport and the ETSI
// Validation Report.
//
// # Main entry types
//
// ProcessExecutor is the generic executor interface; DocumentProcessExecutor
// and CertificateProcessExecutor specialize it for signed documents and bare
// certificates. DefaultSignatureProcessExecutor and
// DefaultCertificateProcessExecutor are the concrete implementations
// consumers construct via New*, configured with a diagnostic data source, a
// validation policy and the current validation time. DetailedReportBuilder,
// SimpleReportBuilder and ETSIValidationReportBuilder assemble the
// corresponding report from the process's conclusions; the "ForCertificate",
// "ForQWAC" and "ForEAAPresentation" variants adapt the same machinery to
// certificate-only, website-authentication and electronic-attestation
// validation.
//
// Ported from dss-validation's executor package (DSS 6.5.RC1); every file
// names its own upstream source in its header.
package executor
