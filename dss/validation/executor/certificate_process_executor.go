// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/certificate/CertificateProcessExecutor.java
// (DSS 6.5.RC1).
//
// Java places this interface in the sub-package
// eu.europa.esig.dss.validation.executor.certificate; the whole executor tree
// is flattened into one Go package (see the batch manifest), which needs no
// renaming - every type name in the tree is already unique.

package executor

import (
	"github.com/ryftcore/dss-go/dss/validation/reports"
)

// CertificateProcessExecutor processes a certificate validation. Port of the
// CertificateProcessExecutor interface (extends
// ProcessExecutor<CertificateReports>).
type CertificateProcessExecutor interface {
	ProcessExecutor[*reports.CertificateReports]

	// SetCertificateId allows specifying the target certificate present in the
	// Diagnostic Data to be verified. Port of setCertificateId(String).
	SetCertificateId(certificateId string)
}
