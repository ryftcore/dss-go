// Ported from dss-validation/.../validation/process/bbb/fc/checks/FormatCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// FormatCheck checks if the signature format is acceptable.
type FormatCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*drjaxb.XmlFC]

	signature *diagnostic.SignatureWrapper
}

// NewFormatCheck is the default constructor.
func NewFormatCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	signature *diagnostic.SignatureWrapper, constraint policy.MultiValuesRule) *FormatCheck {
	c := &FormatCheck{signature: signature}
	c.AbstractMultiValuesCheckItem = bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *FormatCheck) Process() bool {
	// Java calls SignatureLevel#toString(), which is overridden to replace '_'
	// with '-' ("XAdES-BASELINE-B"); the policy's AcceptableFormats ids carry
	// that dashed spelling. A plain string(...) conversion would yield the enum
	// name and only ever match the "*" wildcard.
	return c.ProcessValueCheck(c.signature.SignatureFormat().String())
}

// MessageTag returns the constraint message i18n key.
func (c *FormatCheck) MessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCIEFF }

// ErrorMessageTag returns the error message i18n key.
func (c *FormatCheck) ErrorMessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCIEFFANS }

// FailedIndicationForConclusion returns the Indication on failure.
func (c *FormatCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *FormatCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
