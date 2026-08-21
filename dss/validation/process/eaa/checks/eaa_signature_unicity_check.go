// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAASignatureUnicityCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAASignatureUnicityCheck verifies whether the EAA has been created with a
// single signature.
type EAASignatureUnicityCheck struct {
	*process.ChainItemBase[*jaxb.XmlFC]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAASignatureUnicityCheck is the default constructor.
func NewEAASignatureUnicityCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlFC],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAASignatureUnicityCheck {
	c := &EAASignatureUnicityCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
//
// TODO : to be implemented in TS 119 472-1. Verify the check later (ported
// as-is from Java).
func (c *EAASignatureUnicityCheck) Process() bool {
	return utils.CollectionSize(c.eaa.EAASignatures()) == 1
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAASignatureUnicityCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_SIG_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAASignatureUnicityCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_SIG_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAASignatureUnicityCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAASignatureUnicityCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
