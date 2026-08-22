// Ported from dss-validation/.../validation/process/bbb/fc/checks/PdfSignatureDictionaryCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PdfSignatureDictionaryCheck verifies whether the corresponding signature dictionary is consistent
// across PDF revisions.
type PdfSignatureDictionaryCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfRevision *diagnostic.PDFRevisionWrapper
}

// NewPdfSignatureDictionaryCheck is the default constructor.
func NewPdfSignatureDictionaryCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *PdfSignatureDictionaryCheck {
	c := &PdfSignatureDictionaryCheck{pdfRevision: pdfRevision}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *PdfSignatureDictionaryCheck) Process() bool {
	return c.pdfRevision.IsPdfSignatureDictionaryConsistent()
}

// MessageTag returns the constraint message i18n key.
func (c *PdfSignatureDictionaryCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISDC
}

// ErrorMessageTag returns the error message i18n key.
func (c *PdfSignatureDictionaryCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISDC_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *PdfSignatureDictionaryCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *PdfSignatureDictionaryCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
