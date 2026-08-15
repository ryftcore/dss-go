// Ported from dss-validation/.../validation/process/bbb/fc/checks/AcceptableMimetypeFileContentCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// AcceptableMimetypeFileContentCheck checks if the mimetype file is acceptable.
type AcceptableMimetypeFileContentCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*drjaxb.XmlFC]

	mimetypeFileContent string
}

// NewAcceptableMimetypeFileContentCheck is the default constructor.
func NewAcceptableMimetypeFileContentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	mimetypeFileContent string, constraint policy.MultiValuesRule) *AcceptableMimetypeFileContentCheck {
	c := &AcceptableMimetypeFileContentCheck{mimetypeFileContent: mimetypeFileContent}
	c.AbstractMultiValuesCheckItem = bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *AcceptableMimetypeFileContentCheck) Process() bool {
	return c.ProcessValueCheck(c.mimetypeFileContent)
}

// MessageTag returns the constraint message i18n key.
func (c *AcceptableMimetypeFileContentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IEMCF
}

// ErrorMessageTag returns the error message i18n key.
func (c *AcceptableMimetypeFileContentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IEMCF_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *AcceptableMimetypeFileContentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *AcceptableMimetypeFileContentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
