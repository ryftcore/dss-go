// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAExpirationPresentCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAAExpirationPresentCheck verifies whether the EAA contains the expiration.
type EAAExpirationPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAExpirationPresentCheck is the default constructor.
func NewEAAExpirationPresentCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAAExpirationPresentCheck {
	c := &EAAExpirationPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAExpirationPresentCheck) Process() bool {
	return c.eaa.EAAExpiration() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAExpirationPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAExpPresent
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAExpirationPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAExpPresentANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAExpirationPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAExpirationPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
