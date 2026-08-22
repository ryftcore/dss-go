// Ported from dss-validation/.../validation/process/bbb/fc/checks/PdfAnnotationOverlapCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PdfAnnotationOverlapCheck verifies if a PDF contains overlapping annotations.
type PdfAnnotationOverlapCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfRevision *diagnostic.PDFRevisionWrapper
}

// NewPdfAnnotationOverlapCheck is the default constructor.
func NewPdfAnnotationOverlapCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *PdfAnnotationOverlapCheck {
	c := &PdfAnnotationOverlapCheck{pdfRevision: pdfRevision}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *PdfAnnotationOverlapCheck) Process() bool {
	return len(c.pdfRevision.PdfAnnotationsOverlapConcernedPages()) == 0
}

// MessageTag returns the constraint message i18n key.
func (c *PdfAnnotationOverlapCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_IAOD }

// ErrorMessageTag returns the error message i18n key.
func (c *PdfAnnotationOverlapCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IAOD_ANS
}

// BuildErrorMessage overrides the default to include the concerned pages. Port of the
// overridden protected XmlMessage buildErrorMessage().
func (c *PdfAnnotationOverlapCheck) BuildErrorMessage() *drjaxb.XmlMessage {
	pages := c.pdfRevision.PdfAnnotationsOverlapConcernedPages()
	return c.BuildXmlMessage(c.ErrorMessageTag(), bigIntSlicePrint(pages))
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *PdfAnnotationOverlapCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *PdfAnnotationOverlapCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
