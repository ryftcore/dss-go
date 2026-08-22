// Ported from dss-validation/.../validation/process/bbb/fc/checks/AcceptableMimetypeFileContentCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// AcceptableMimetypeFileContentCheck checks if the mimetype file is acceptable.
type AcceptableMimetypeFileContentCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*drjaxb.XmlFC]

	mimetypeFileContent string
}

// NewAcceptableMimetypeFileContentCheck is the default constructor.
func NewAcceptableMimetypeFileContentCheck(i18nProvider *i18n.Provider, result *process.Result[*drjaxb.XmlFC],
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
	return i18n.MessageTagBBBFCIEMCF
}

// ErrorMessageTag returns the error message i18n key.
func (c *AcceptableMimetypeFileContentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCIEMCFANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *AcceptableMimetypeFileContentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *AcceptableMimetypeFileContentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
