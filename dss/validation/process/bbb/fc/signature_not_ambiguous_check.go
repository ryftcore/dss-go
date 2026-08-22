// Ported from dss-validation/.../validation/process/bbb/fc/checks/SignatureNotAmbiguousCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignatureNotAmbiguousCheck checks if the signature can be identified.
type SignatureNotAmbiguousCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	signature *diagnostic.SignatureWrapper
}

// NewSignatureNotAmbiguousCheck is the default constructor.
func NewSignatureNotAmbiguousCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *SignatureNotAmbiguousCheck {
	c := &SignatureNotAmbiguousCheck{signature: signature}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *SignatureNotAmbiguousCheck) Process() bool { return !c.signature.IsSignatureDuplicated() }

// MessageTag returns the constraint message i18n key.
func (c *SignatureNotAmbiguousCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_ISD }

// ErrorMessageTag returns the error message i18n key.
func (c *SignatureNotAmbiguousCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISD_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *SignatureNotAmbiguousCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *SignatureNotAmbiguousCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
