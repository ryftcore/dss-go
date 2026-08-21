// Package reports ports dss-validation's
// eu.europa.esig.dss.validation.reports package: the top-level report
// containers consumers receive back from a validation call, bundling the
// diagnostic data, simple report, detailed report and (where applicable) the
// ETSI Validation Report produced by validation/executor.
//
// # Main entry types
//
// Reports wraps the four report flavors for a signed-document validation;
// CertificateReports wraps the diagnostic data, simple and detailed report
// for a bare-certificate validation. AbstractReportsBase implements the
// XML-marshaling and typed-accessor plumbing both embed. DSSReportException
// reports a malformed or inconsistent report.
//
// Ported from dss-validation's reports package (DSS 6.5.RC1); every file
// names its own upstream source in its header.
package reports
