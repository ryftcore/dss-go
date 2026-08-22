// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/status/EAARevocationExpirationTimeCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAARevocationExpirationTimeCheck verifies whether the EAA revocation token
// contains an expiration time.
type EAARevocationExpirationTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper
}

// NewEAARevocationExpirationTimeCheck is the default constructor.
func NewEAARevocationExpirationTimeCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper, constraint policy.LevelRule) *EAARevocationExpirationTimeCheck {
	c := &EAARevocationExpirationTimeCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaaStatusToken: eaaStatusToken,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAARevocationExpirationTimeCheck) Process() bool {
	return c.eaaStatusToken.ExpirationTime() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAARevocationExpirationTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevExp
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationExpirationTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevExpANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationExpirationTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationExpirationTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
