// Ported from dss-validation/.../validation/process/bbb/fc/checks/FieldMDPCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// FieldMDPCheck verifies a signature according to given permissions for the document in /FieldMDP.
type FieldMDPCheck struct {
	AbstractPdfLockDictionaryCheck
}

// NewFieldMDPCheck is the default constructor.
func NewFieldMDPCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *FieldMDPCheck {
	c := &FieldMDPCheck{}
	c.InitAbstractPdfLockDictionaryCheck(i18nProvider, result, pdfRevision, pdfRevision.FieldMDP(), constraint)
	c.InitChainItem(c)
	return c
}

// MessageTag returns the constraint message i18n key.
func (c *FieldMDPCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_ISVAFMDPD }

// ErrorMessageTag returns the error message i18n key.
func (c *FieldMDPCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISVAFMDPD_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *FieldMDPCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *FieldMDPCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
