// Ported from dss-validation/.../validation/process/bbb/fc/checks/AcceptableZipCommentCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// AcceptableZipCommentCheck checks if the zip comment is acceptable.
type AcceptableZipCommentCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*drjaxb.XmlFC]

	zipComment string
}

// NewAcceptableZipCommentCheck is the default constructor.
func NewAcceptableZipCommentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	zipComment string, constraint policy.MultiValuesRule) *AcceptableZipCommentCheck {
	c := &AcceptableZipCommentCheck{zipComment: zipComment}
	c.AbstractMultiValuesCheckItem = bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *AcceptableZipCommentCheck) Process() bool { return c.ProcessValueCheck(c.zipComment) }

// MessageTag returns the constraint message i18n key.
func (c *AcceptableZipCommentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ITEZCF
}

// ErrorMessageTag returns the error message i18n key.
func (c *AcceptableZipCommentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ITEZCF_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *AcceptableZipCommentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *AcceptableZipCommentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
