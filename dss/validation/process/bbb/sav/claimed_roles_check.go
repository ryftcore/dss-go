// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/ClaimedRolesCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// ClaimedRolesCheck checks if the claimed roles are acceptable.
type ClaimedRolesCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewClaimedRolesCheck is the default constructor. Port of
// ClaimedRolesCheck(I18nProvider, XmlSAV, SignatureWrapper, MultiValuesRule).
func NewClaimedRolesCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.MultiValuesRule) *ClaimedRolesCheck {
	c := &ClaimedRolesCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		signature:                    signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ClaimedRolesCheck) Process() bool {
	claimedRoles := c.signature.SignerRoleDetails(c.signature.ClaimedRoles())
	return c.ProcessValuesCheck(claimedRoles)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ClaimedRolesCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ICRM
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ClaimedRolesCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ICRM_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ClaimedRolesCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ClaimedRolesCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
