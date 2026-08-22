// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/status/EAARevocationIssuanceTimeCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAARevocationIssuanceTimeCheck verifies whether the EAA revocation token
// contains an issuance time.
type EAARevocationIssuanceTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper
}

// NewEAARevocationIssuanceTimeCheck is the default constructor.
func NewEAARevocationIssuanceTimeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper, constraint policy.LevelRule) *EAARevocationIssuanceTimeCheck {
	c := &EAARevocationIssuanceTimeCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaaStatusToken: eaaStatusToken,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAARevocationIssuanceTimeCheck) Process() bool {
	return c.eaaStatusToken.IssuedAt() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAARevocationIssuanceTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevISS
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationIssuanceTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevISSANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationIssuanceTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationIssuanceTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
