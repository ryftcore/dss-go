// Ported from dss-validation/.../validation/process/bbb/fc/checks/AnnotationChangesCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AnnotationChangesCheck checks whether a document contains annotation creation, modification or
// deletion changes occurred after the signature revision.
type AnnotationChangesCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfRevision *diagnostic.PDFRevisionWrapper
}

// NewAnnotationChangesCheck is the default constructor.
func NewAnnotationChangesCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *AnnotationChangesCheck {
	c := &AnnotationChangesCheck{pdfRevision: pdfRevision}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *AnnotationChangesCheck) Process() bool {
	return len(c.pdfRevision.PdfAnnotationChanges()) == 0
}

// MessageTag returns the constraint message i18n key.
func (c *AnnotationChangesCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCDSCNACMDM
}

// ErrorMessageTag returns the error message i18n key.
func (c *AnnotationChangesCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCDSCNACMDMANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *AnnotationChangesCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *AnnotationChangesCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
