// Ported from dss-validation/.../validation/process/bbb/fc/checks/PDFAProfileCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// PDFAProfileCheck checks whether a determined PDF/A profile of the input document is acceptable.
type PDFAProfileCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*drjaxb.XmlFC]

	pdfaProfile string
}

// NewPDFAProfileCheck is the default constructor.
func NewPDFAProfileCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfaProfile string, constraint policy.MultiValuesRule) *PDFAProfileCheck {
	c := &PDFAProfileCheck{pdfaProfile: pdfaProfile}
	c.AbstractMultiValuesCheckItem = bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *PDFAProfileCheck) Process() bool { return c.ProcessValueCheck(c.pdfaProfile) }

// MessageTag returns the constraint message i18n key.
func (c *PDFAProfileCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_DDAPDFAF }

// ErrorMessageTag returns the error message i18n key.
func (c *PDFAProfileCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_DDAPDFAF_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *PDFAProfileCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *PDFAProfileCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
