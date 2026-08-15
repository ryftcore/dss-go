// Ported from dss-validation/.../validation/process/bbb/fc/checks/PdfPageDifferenceCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// PdfPageDifferenceCheck verifies if a PDF contains difference between page amount in different revisions.
type PdfPageDifferenceCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfRevision *diagnostic.PDFRevisionWrapper
}

// NewPdfPageDifferenceCheck is the default constructor.
func NewPdfPageDifferenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *PdfPageDifferenceCheck {
	c := &PdfPageDifferenceCheck{pdfRevision: pdfRevision}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *PdfPageDifferenceCheck) Process() bool {
	return len(c.pdfRevision.PdfPageDifferenceConcernedPages()) == 0
}

// MessageTag returns the constraint message i18n key.
func (c *PdfPageDifferenceCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_DSFREAP }

// ErrorMessageTag returns the error message i18n key.
func (c *PdfPageDifferenceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_DSFREAP_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *PdfPageDifferenceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *PdfPageDifferenceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
