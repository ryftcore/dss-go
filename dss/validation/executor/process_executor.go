// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/ProcessExecutor.java
// (DSS 6.5.RC1).

package executor

import (
	"time"

	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/model/policy"
)

// ProcessExecutor allows to define how the validation process should be
// carried out. Port of the ProcessExecutor<R extends AbstractReports>
// interface; the Java type parameter bound is expressed by the concrete
// report container types (reports.Reports, reports.CertificateReports)
// that instantiate it.
type ProcessExecutor[R any] interface {

	// SetCurrentTime allows to set the time that is used during the
	// validation process execution. Port of setCurrentTime(Date).
	SetCurrentTime(currentDate time.Time)

	// CurrentTime returns the validation time. Port of getCurrentTime().
	CurrentTime() time.Time

	// SetDiagnosticData allows to set the XmlDiagnosticData that is used
	// during the validation process execution. Port of
	// setDiagnosticData(XmlDiagnosticData).
	SetDiagnosticData(diagnosticData *diagnosticjaxb.XmlDiagnosticData)

	// SetValidationPolicy allows to set the validation policy that is used
	// during the validation process execution. Port of
	// setValidationPolicy(ValidationPolicy).
	SetValidationPolicy(validationPolicy policy.ValidationPolicy)

	// ValidationPolicy returns the used validation policy. Port of
	// getValidationPolicy().
	ValidationPolicy() policy.ValidationPolicy

	// SetLocale allows to set a language setting for generated Reports.
	// Port of setLocale(Locale); Java's Locale is represented by its
	// language tag, matching i18n.NewI18nProviderForLocale.
	SetLocale(locale string)

	// Execute allows to run the validation process. Port of execute().
	Execute() R
}
