// Ported from dss-validation/.../validation/process/bbb/fc/checks/SignerInformationStoreCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignerInformationStoreCheck checks if only one SignerInformationStore entry is present for a PAdES signature.
type SignerInformationStoreCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	signature *diagnostic.SignatureWrapper
}

// NewSignerInformationStoreCheck is the default constructor.
func NewSignerInformationStoreCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *SignerInformationStoreCheck {
	c := &SignerInformationStoreCheck{signature: signature}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *SignerInformationStoreCheck) Process() bool {
	store := c.signature.SignatureInformationStore()
	if len(store) > 0 {
		return len(store) == 1
	}
	return false
}

// MessageTag returns the constraint message i18n key.
func (c *SignerInformationStoreCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IOSIP
}

// ErrorMessageTag returns the error message i18n key.
func (c *SignerInformationStoreCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IOSIP_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *SignerInformationStoreCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *SignerInformationStoreCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
