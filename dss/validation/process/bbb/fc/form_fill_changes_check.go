// Ported from dss-validation/.../validation/process/bbb/fc/checks/FormFillChangesCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// FormFillChangesCheck checks whether a document contains form fill or signing modifications
// occurred after the signature revision.
type FormFillChangesCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfRevision *diagnostic.PDFRevisionWrapper
}

// NewFormFillChangesCheck is the default constructor.
func NewFormFillChangesCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *FormFillChangesCheck {
	c := &FormFillChangesCheck{pdfRevision: pdfRevision}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *FormFillChangesCheck) Process() bool {
	return len(c.pdfRevision.PdfSignatureOrFormFillChanges()) == 0
}

// MessageTag returns the constraint message i18n key.
func (c *FormFillChangesCheck) MessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCDSCNFFSM }

// ErrorMessageTag returns the error message i18n key.
func (c *FormFillChangesCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCDSCNFFSMANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *FormFillChangesCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *FormFillChangesCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
