// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAClaimsCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// EAAClaimsCheck verifies whether the EAA contains all the specified claims,
// either as part of the original payload or through the provided selective
// disclosures.
type EAAClaimsCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAClaimsCheck is the default constructor.
func NewEAAClaimsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.MultiValuesRule) *EAAClaimsCheck {
	c := &EAAClaimsCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		eaa:                          eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAClaimsCheck) Process() bool {
	claimNames := c.eaa.AllEAAPayloadClaimNames()
	return c.ProcessValuesForEachExpectedCheck(claimNames)
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *EAAClaimsCheck) BuildAdditionalInfo() *string {
	claimNames := c.eaa.AllEAAPayloadClaimNames()
	notPresentClaims := make([]string, 0)
	for _, v := range c.Values() {
		if !process.ProcessValueCheck(v, claimNames) {
			notPresentClaims = append(notPresentClaims, v)
		}
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_EAA_CLAIMS_INFO, utils.JoinStrings(notPresentClaims, ", "))
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAClaimsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CLAIMS
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAClaimsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CLAIMS_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAClaimsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAClaimsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
