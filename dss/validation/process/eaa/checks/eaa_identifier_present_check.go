// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAIdentifierPresentCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAAIdentifierPresentCheck verifies whether the EAA contains the identifier.
type EAAIdentifierPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAIdentifierPresentCheck is the default constructor.
func NewEAAIdentifierPresentCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAAIdentifierPresentCheck {
	c := &EAAIdentifierPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAIdentifierPresentCheck) Process() bool {
	switch c.eaa.EAAType() {
	case enumerations.EAATypeSDJWTVC:
		return c.eaa.EAAIdentifier() != ""
	case enumerations.EAATypeISOIECMDoc:
		return c.eaa.DocumentNumber() != ""
	default:
		// Other EAA types not supported.
		return false
	}
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAIdentifierPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAIdentifierPresent
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAIdentifierPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAIdentifierPresentANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAIdentifierPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAIdentifierPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
