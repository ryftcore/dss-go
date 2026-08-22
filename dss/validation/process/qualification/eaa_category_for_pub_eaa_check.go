// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/checks/EAACategoryForPubEAACheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAACategoryForPubEAACheck verifies whether the EAA payload contains an
// indication that the attestation has been issued as an electronic
// attestation of attributes issued by or on behalf of a public body
// responsible for an authentic source.
type EAACategoryForPubEAACheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationEAAQualificationProcess]

	// eaa is the EAA presentation to be checked.
	eaa *diagnostic.EAAWrapper
}

// NewEAACategoryForPubEAACheck is the default constructor. Port of
// EAACategoryForPubEAACheck(I18nProvider, XmlValidationEAAQualificationProcess, EAAWrapper, LevelRule).
func NewEAACategoryForPubEAACheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationEAAQualificationProcess],
	eaa *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAACategoryForPubEAACheck {
	c := &EAACategoryForPubEAACheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaa,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAACategoryForPubEAACheck) Process() bool {
	return enumerations.EAACategoryEUPubEAA.URN() == c.eaa.EAACategory()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAACategoryForPubEAACheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CAT_PUBEAA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *EAACategoryForPubEAACheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CAT_PUBEAA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAACategoryForPubEAACheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *EAACategoryForPubEAACheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
