// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/pid/checks/PIDProviderListCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PIDProviderListCheck verifies whether the List of Trusted Entities is of
// PID Providers list type.
type PIDProviderListCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationPIDQualificationProcess]

	// listTypeUri is the List Type URI.
	listTypeUri string
}

// NewPIDProviderListCheck is the default constructor. Port of
// PIDProviderListCheck(I18nProvider, XmlValidationPIDQualificationProcess, String, LevelRule).
func NewPIDProviderListCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationPIDQualificationProcess],
	listTypeUri string, constraint policy.LevelRule) *PIDProviderListCheck {
	c := &PIDProviderListCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		listTypeUri:   listTypeUri,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *PIDProviderListCheck) Process() bool {
	listType := enumerations.ListTypeFromURI(c.listTypeUri)
	return enumerations.LoTETypeEnumEUPIDProvidersList == listType
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *PIDProviderListCheck) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTag_CERTIFICATE_USAGE_LIST_TYPE, c.listTypeUri)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PIDProviderListCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PID_LOTE_TYPE_PID_PROVIDERS
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *PIDProviderListCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PID_LOTE_TYPE_PID_PROVIDERS_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *PIDProviderListCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *PIDProviderListCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
