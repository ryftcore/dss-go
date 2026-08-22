// Ported from dss-validation/.../validation/process/bbb/fc/checks/UndefinedChangesCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// UndefinedChangesCheck checks whether a document contains undefined object modifications.
type UndefinedChangesCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfRevision *diagnostic.PDFRevisionWrapper
}

// NewUndefinedChangesCheck is the default constructor.
func NewUndefinedChangesCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *UndefinedChangesCheck {
	c := &UndefinedChangesCheck{pdfRevision: pdfRevision}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *UndefinedChangesCheck) Process() bool {
	return len(c.pdfRevision.PdfUndefinedChanges()) == 0
}

// MessageTag returns the constraint message i18n key.
func (c *UndefinedChangesCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_DSCNUOM }

// ErrorMessageTag returns the error message i18n key.
func (c *UndefinedChangesCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_DSCNUOM_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *UndefinedChangesCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *UndefinedChangesCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
