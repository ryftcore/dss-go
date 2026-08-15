// Ported from dss-validation/.../validation/process/bbb/fc/checks/ByteRangeCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// ByteRangeCheck verifies the applicability of the /ByteRange field extracted from a corresponding PDF revision.
type ByteRangeCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfRevision *diagnostic.PDFRevisionWrapper
}

// NewByteRangeCheck is the default constructor.
func NewByteRangeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *ByteRangeCheck {
	c := &ByteRangeCheck{pdfRevision: pdfRevision}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *ByteRangeCheck) Process() bool { return c.pdfRevision.IsSignatureByteRangeValid() }

// MessageTag returns the constraint message i18n key.
func (c *ByteRangeCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_IBRV }

// ErrorMessageTag returns the error message i18n key.
func (c *ByteRangeCheck) ErrorMessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_IBRV_ANS }

// FailedIndicationForConclusion returns the Indication on failure.
func (c *ByteRangeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *ByteRangeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
