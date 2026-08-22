// Ported from dss-validation/.../validation/process/bbb/fc/checks/PDFAComplianceCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PDFAComplianceCheck checks whether the input document is compliant with the determined PDF/A format.
type PDFAComplianceCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfaCompliant bool
}

// NewPDFAComplianceCheck is the default constructor.
func NewPDFAComplianceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfaCompliant bool, constraint policy.LevelRule) *PDFAComplianceCheck {
	c := &PDFAComplianceCheck{pdfaCompliant: pdfaCompliant}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *PDFAComplianceCheck) Process() bool { return c.pdfaCompliant }

// MessageTag returns the constraint message i18n key.
func (c *PDFAComplianceCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_IDPDFAC }

// ErrorMessageTag returns the error message i18n key.
func (c *PDFAComplianceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IDPDFAC_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *PDFAComplianceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *PDFAComplianceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
