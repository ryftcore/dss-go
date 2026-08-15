// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/checks/EAACategoryForEAAPresenceCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAACategoryForEAAPresenceCheck verifies whether the EAA payload contains
// an indication that the attestation has been issued as an EU
// non-qualified electronic attestation of attributes.
type EAACategoryForEAAPresenceCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationEAAQualificationProcess]

	// eaa is the EAA presentation to be checked.
	eaa *diagnostic.EAAWrapper
}

// NewEAACategoryForEAAPresenceCheck is the default constructor. Port of
// EAACategoryForEAAPresenceCheck(I18nProvider, XmlValidationEAAQualificationProcess, EAAWrapper, LevelRule).
func NewEAACategoryForEAAPresenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationEAAQualificationProcess],
	eaa *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAACategoryForEAAPresenceCheck {
	c := &EAACategoryForEAAPresenceCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaa,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAACategoryForEAAPresenceCheck) Process() bool {
	category := c.eaa.EAACategory()
	if category == "" {
		return false
	}
	for _, v := range enumerations.EAACategoryValues() {
		if v.URN() == category {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAACategoryForEAAPresenceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CAT_EAA
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *EAACategoryForEAAPresenceCheck) BuildErrorMessage() *jaxb.XmlMessage {
	category := c.eaa.EAACategory()
	if category == "" {
		return c.BuildXmlMessage(i18n.MessageTag_EAA_CAT_EAA_ANS_1)
	}
	return c.BuildXmlMessage(i18n.MessageTag_EAA_CAT_EAA_ANS_2, category)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAACategoryForEAAPresenceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *EAACategoryForEAAPresenceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
