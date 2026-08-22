// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAASupportedClaimsCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// EAASupportedClaimsCheck verifies whether the EAA contains only supported
// claims.
type EAASupportedClaimsCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAASupportedClaimsCheck is the default constructor.
func NewEAASupportedClaimsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.MultiValuesRule) *EAASupportedClaimsCheck {
	c := &EAASupportedClaimsCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		eaa:                          eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAASupportedClaimsCheck) Process() bool {
	claimNames := c.eaa.AllEAAPayloadClaimNames()
	return c.ProcessAllValuesCheck(claimNames)
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *EAASupportedClaimsCheck) BuildAdditionalInfo() *string {
	unsupportedClaims := make([]string, 0)
	for _, cl := range c.eaa.AllEAAPayloadClaimNames() {
		if !c.ProcessValueCheck(cl) {
			unsupportedClaims = append(unsupportedClaims, cl)
		}
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_EAA_UNSUPPORTED_CLAIMS, utils.JoinStrings(unsupportedClaims, ", "))
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAASupportedClaimsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_SUPPORTED_CLAIMS
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAASupportedClaimsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_SUPPORTED_CLAIMS_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAASupportedClaimsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAASupportedClaimsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
