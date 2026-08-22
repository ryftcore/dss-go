// Ported from dss-validation/.../validation/process/bbb/fc/checks/PdfVisualDifferenceCheck.java (DSS 6.5.RC1).
package fc

import (
	"math/big"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PdfVisualDifferenceCheck verifies if a PDF has visual difference between revisions.
type PdfVisualDifferenceCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfRevision *diagnostic.PDFRevisionWrapper
}

// NewPdfVisualDifferenceCheck is the default constructor.
func NewPdfVisualDifferenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *PdfVisualDifferenceCheck {
	c := &PdfVisualDifferenceCheck{pdfRevision: pdfRevision}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *PdfVisualDifferenceCheck) Process() bool {
	return len(c.pdfRevision.PdfVisualDifferenceConcernedPages()) == 0
}

// MessageTag returns the constraint message i18n key.
func (c *PdfVisualDifferenceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCIVDBSFR
}

// ErrorMessageTag returns the error message i18n key.
func (c *PdfVisualDifferenceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCIVDBSFRANS
}

// BuildErrorMessage overrides the default to include the concerned pages. Port of the
// overridden protected XmlMessage buildErrorMessage().
func (c *PdfVisualDifferenceCheck) BuildErrorMessage() *drjaxb.XmlMessage {
	pages := c.pdfRevision.PdfVisualDifferenceConcernedPages()
	return c.BuildXmlMessage(c.ErrorMessageTag(), bigIntSlicePrint(pages))
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *PdfVisualDifferenceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *PdfVisualDifferenceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}

// bigIntSlicePrint mirrors Java's List<BigInteger>.toString() lexical form
// ("[1, 2, 3]"), shared by the two PDF checks that surface concerned pages in
// their error message.
func bigIntSlicePrint(values []*big.Int) string {
	s := "["
	for i, v := range values {
		if i > 0 {
			s += ", "
		}
		s += v.String()
	}
	return s + "]"
}
