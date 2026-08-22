// Ported from dss-validation/.../validation/process/bbb/fc/checks/ZipCommentPresentCheck.java (DSS 6.5.RC1).
package fc

import (
	"strings"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ZipCommentPresentCheck checks if the zip comment is present.
type ZipCommentPresentCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	zipComment string
}

// NewZipCommentPresentCheck is the default constructor.
func NewZipCommentPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	zipComment string, constraint policy.LevelRule) *ZipCommentPresentCheck {
	c := &ZipCommentPresentCheck{zipComment: zipComment}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *ZipCommentPresentCheck) Process() bool {
	return strings.TrimSpace(c.zipComment) != ""
}

// MessageTag returns the constraint message i18n key.
func (c *ZipCommentPresentCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_ITZCP }

// ErrorMessageTag returns the error message i18n key.
func (c *ZipCommentPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ITZCP_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *ZipCommentPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *ZipCommentPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
