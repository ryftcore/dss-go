// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/checks/EAACategoryForQEAACheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAACategoryForQEAACheck verifies whether the EAA payload contains an
// indication that the attestation has been issued as an EU qualified
// electronic attestation of attributes.
type EAACategoryForQEAACheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationEAAQualificationProcess]

	// eaa is the EAA presentation to be checked.
	eaa *diagnostic.EAAWrapper
}

// NewEAACategoryForQEAACheck is the default constructor. Port of
// EAACategoryForQEAACheck(I18nProvider, XmlValidationEAAQualificationProcess, EAAWrapper, LevelRule).
func NewEAACategoryForQEAACheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationEAAQualificationProcess],
	eaa *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAACategoryForQEAACheck {
	c := &EAACategoryForQEAACheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaa,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAACategoryForQEAACheck) Process() bool {
	return enumerations.EAACategory_EU_QEAA.URN() == c.eaa.EAACategory()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAACategoryForQEAACheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CAT_QEAA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *EAACategoryForQEAACheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CAT_QEAA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAACategoryForQEAACheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *EAACategoryForQEAACheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
