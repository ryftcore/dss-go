// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/DocumentProcessExecutor.java
// (DSS 6.5.RC1).

package executor

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/validation/reports"
)

// DocumentProcessExecutor processes a document validation. Port of the
// DocumentProcessExecutor interface (extends ProcessExecutor<Reports>).
type DocumentProcessExecutor interface {
	ProcessExecutor[*reports.Reports]

	// SetValidationLevel allows to set the validation level that is used
	// during the validation process execution. Port of
	// setValidationLevel(ValidationLevel).
	SetValidationLevel(validationLevel enumerations.ValidationLevel)

	// SetEnableEtsiValidationReport specifies if the ETSI Validation Report
	// must be created. Port of setEnableEtsiValidationReport(boolean).
	SetEnableEtsiValidationReport(enableEtsiValidationReport bool)

	// SetIncludeSemantics allows to enable/disable the semantics inclusion
	// in the reports (Indication / SubIndication meanings). Disabled by
	// default. Port of setIncludeSemantics(boolean).
	SetIncludeSemantics(includeSemantics bool)
}
