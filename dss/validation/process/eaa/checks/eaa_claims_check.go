// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAClaimsCheck.java (DSS 6.5.RC1).
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

// EAAClaimsCheck verifies whether the EAA contains all the specified claims,
// either as part of the original payload or through the provided selective
// disclosures.
type EAAClaimsCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAClaimsCheck is the default constructor.
func NewEAAClaimsCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
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
		if !process.ValueCheck(v, claimNames) {
			notPresentClaims = append(notPresentClaims, v)
		}
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagEAAClaimsInfo, utils.JoinStrings(notPresentClaims, ", "))
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAClaimsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAClaims
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAClaimsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAClaimsANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAClaimsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAClaimsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
