// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/pid/checks/PIDIssuanceTrustedEntityServicesCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PIDIssuanceTrustedEntityServicesCheck checks if the list of extracted
// trusted entity services for PID issuance is not empty.
type PIDIssuanceTrustedEntityServicesCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationPIDQualificationProcess]

	// pidIssuanceTrustedEntityServices is the pre-filtered list of trusted
	// entity services for PID issuance.
	pidIssuanceTrustedEntityServices []*diagnostic.TrustedEntityServiceWrapper
}

// NewPIDIssuanceTrustedEntityServicesCheck is the default constructor. Port
// of
// PIDIssuanceTrustedEntityServicesCheck(I18nProvider, XmlValidationPIDQualificationProcess, List, LevelRule).
func NewPIDIssuanceTrustedEntityServicesCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationPIDQualificationProcess],
	pidIssuanceTrustedEntityServices []*diagnostic.TrustedEntityServiceWrapper,
	constraint policy.LevelRule) *PIDIssuanceTrustedEntityServicesCheck {
	c := &PIDIssuanceTrustedEntityServicesCheck{
		ChainItemBase:                    process.NewChainItemBase(i18nProvider, result, constraint),
		pidIssuanceTrustedEntityServices: pidIssuanceTrustedEntityServices,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *PIDIssuanceTrustedEntityServicesCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.pidIssuanceTrustedEntityServices)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PIDIssuanceTrustedEntityServicesCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PID_STI_PID_ISSUANCE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *PIDIssuanceTrustedEntityServicesCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PID_STI_PID_ISSUANCE_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *PIDIssuanceTrustedEntityServicesCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *PIDIssuanceTrustedEntityServicesCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
