// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAAdministrativeIssuanceDatePresentCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAAAdministrativeIssuanceDatePresentCheck verifies whether the EAA contains
// an administrative issuance date.
type EAAAdministrativeIssuanceDatePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAAdministrativeIssuanceDatePresentCheck is the default constructor.
func NewEAAAdministrativeIssuanceDatePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAAAdministrativeIssuanceDatePresentCheck {
	c := &EAAAdministrativeIssuanceDatePresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAAdministrativeIssuanceDatePresentCheck) Process() bool {
	return c.eaa.AdministrativeIssuanceDate() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAAdministrativeIssuanceDatePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAAIDPresent
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAAdministrativeIssuanceDatePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAAIDPresentANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAAdministrativeIssuanceDatePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAAdministrativeIssuanceDatePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
