// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAIdentifierPresentCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAAIdentifierPresentCheck verifies whether the EAA contains the identifier.
type EAAIdentifierPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAIdentifierPresentCheck is the default constructor.
func NewEAAIdentifierPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAAIdentifierPresentCheck {
	c := &EAAIdentifierPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAIdentifierPresentCheck) Process() bool {
	switch c.eaa.EAAType() {
	case enumerations.EAAType_SD_JWT_VC:
		return c.eaa.EAAIdentifier() != ""
	case enumerations.EAAType_ISO_IEC_MDOC:
		return c.eaa.DocumentNumber() != ""
	default:
		// Other EAA types not supported.
		return false
	}
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAIdentifierPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_IDENTIFIER_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAIdentifierPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_IDENTIFIER_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAIdentifierPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAIdentifierPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
