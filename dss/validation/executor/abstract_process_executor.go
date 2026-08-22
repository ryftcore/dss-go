// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/AbstractProcessExecutor.java
// (DSS 6.5.RC1).
//
// Java's Objects.requireNonNull(...) guards in assertConfigurationValid()
// signal a programming error (a NullPointerException escapes execute()), so
// they are ported as panics, matching the convention already used across the
// ported validation process tree (see process/chain_item.go).

package executor

import (
	"time"

	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
)

// AbstractProcessExecutor is the abstract validation process executor.
// Port of AbstractProcessExecutor; the protected fields are exported so
// that the concrete executors embedding it can read them the way the Java
// subclasses do.
type AbstractProcessExecutor struct {
	// CurrentTimeValue is the validation time. Port of the protected
	// currentTime field.
	CurrentTimeValue time.Time

	// Policy is the validation policy. Port of the protected policy field.
	Policy policy.ValidationPolicy

	// JaxbDiagnosticData is the DiagnosticData. Port of the protected
	// jaxbDiagnosticData field.
	JaxbDiagnosticData *diagnosticjaxb.XmlDiagnosticData

	// i18nProvider is the i18n provider. Port of the private i18nProvider
	// field.
	i18nProvider *i18n.Provider
}

// NewAbstractProcessExecutor is the default constructor instantiating the
// object with null values and current time. Port of the protected
// AbstractProcessExecutor() constructor, which relies on the
// "protected Date currentTime = new Date();" field initializer.
func NewAbstractProcessExecutor() AbstractProcessExecutor {
	return AbstractProcessExecutor{CurrentTimeValue: time.Now()}
}

// SetCurrentTime sets the validation time. Port of setCurrentTime(Date).
func (e *AbstractProcessExecutor) SetCurrentTime(currentTime time.Time) {
	e.CurrentTimeValue = currentTime
}

// CurrentTime gets the validation time. Port of getCurrentTime().
func (e *AbstractProcessExecutor) CurrentTime() time.Time {
	return e.CurrentTimeValue
}

// SetDiagnosticData sets the DiagnosticData. Port of
// setDiagnosticData(XmlDiagnosticData).
func (e *AbstractProcessExecutor) SetDiagnosticData(diagnosticData *diagnosticjaxb.XmlDiagnosticData) {
	e.JaxbDiagnosticData = diagnosticData
}

// ValidationPolicy gets the validation policy. Port of
// getValidationPolicy().
func (e *AbstractProcessExecutor) ValidationPolicy() policy.ValidationPolicy {
	return e.Policy
}

// SetValidationPolicy sets the validation policy. Port of
// setValidationPolicy(ValidationPolicy).
func (e *AbstractProcessExecutor) SetValidationPolicy(validationPolicy policy.ValidationPolicy) {
	e.Policy = validationPolicy
}

// SetLocale sets the locale to use to generate messages. Port of
// setLocale(Locale).
//
// Java's Objects.requireNonNull("Locale cannot be null!") has no Go
// counterpart here: a Locale becomes a language tag string, and the empty
// string is NOT its null - i18n.NewProviderForLocale documents "" as
// exactly Java's Locale.getDefault(), which is what both
// SignedDocumentValidator and AbstractCertificateValidator initialise their
// locale field to and pass straight through to this setter.
func (e *AbstractProcessExecutor) SetLocale(locale string) {
	e.i18nProvider = i18n.NewProviderForLocale(locale)
}

// Provider gets the i18nProvider, instantiating it with the default
// locale on first use. Port of the protected getI18nProvider().
func (e *AbstractProcessExecutor) I18nProvider() *i18n.Provider {
	if e.i18nProvider == nil {
		e.i18nProvider = i18n.NewProvider()
	}
	return e.i18nProvider
}

// AssertConfigurationValid checks if the configuration is valid. Port of
// the protected assertConfigurationValid().
func (e *AbstractProcessExecutor) AssertConfigurationValid() {
	if e.JaxbDiagnosticData == nil {
		panic("The diagnostic data is missing")
	}
	if e.Policy == nil {
		panic("The validation policy is missing")
	}
	if e.CurrentTimeValue.IsZero() {
		panic("The current time is missing")
	}
}
