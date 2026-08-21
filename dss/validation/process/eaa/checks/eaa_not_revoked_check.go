// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAANotRevokedCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAANotRevokedCheck checks whether the corresponding status declares that
// the EAA is not revoked.
type EAANotRevokedCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationWrapper
}

// NewEAANotRevokedCheck is the default constructor.
func NewEAANotRevokedCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaStatusToken *diagnostic.EAARevocationWrapper, constraint policy.LevelRule) *EAANotRevokedCheck {
	c := &EAANotRevokedCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaaStatusToken: eaaStatusToken,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAANotRevokedCheck) Process() bool {
	return c.eaaStatusToken == nil || enumerations.EAAStatus_INVALID != c.eaaStatusToken.Status()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAANotRevokedCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_NOT_REV
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAANotRevokedCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_NOT_REV_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAANotRevokedCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAANotRevokedCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_REVOKED
}
