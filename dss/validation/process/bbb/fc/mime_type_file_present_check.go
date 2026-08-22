// Ported from dss-validation/.../validation/process/bbb/fc/checks/MimeTypeFilePresentCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// MimeTypeFilePresentCheck checks if a mimetype file is present.
type MimeTypeFilePresentCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	mimetypePresent bool
}

// NewMimeTypeFilePresentCheck is the default constructor.
func NewMimeTypeFilePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	mimetypePresent bool, constraint policy.LevelRule) *MimeTypeFilePresentCheck {
	c := &MimeTypeFilePresentCheck{mimetypePresent: mimetypePresent}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *MimeTypeFilePresentCheck) Process() bool { return c.mimetypePresent }

// MessageTag returns the constraint message i18n key.
func (c *MimeTypeFilePresentCheck) MessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCITMFP }

// ErrorMessageTag returns the error message i18n key.
func (c *MimeTypeFilePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCITMFPANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *MimeTypeFilePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *MimeTypeFilePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
