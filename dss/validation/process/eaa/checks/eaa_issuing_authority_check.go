// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAIssuingAuthorityCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// EAAIssuingAuthorityCheck verifies whether the EAA issuing authority claim
// contains one of the expected values.
type EAAIssuingAuthorityCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAIssuingAuthorityCheck is the default constructor.
func NewEAAIssuingAuthorityCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.MultiValuesRule) *EAAIssuingAuthorityCheck {
	c := &EAAIssuingAuthorityCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		eaa:                          eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAIssuingAuthorityCheck) Process() bool {
	return c.ProcessValueCheck(c.eaa.DocumentIssuingAuthority())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAIssuingAuthorityCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_ISS_AUTH
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAIssuingAuthorityCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_ISS_AUTH_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAIssuingAuthorityCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAIssuingAuthorityCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
