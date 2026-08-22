// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/status/EAARevocationSubjectMatchCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAARevocationSubjectMatchCheck verifies whether the EAA revocation token's
// subject matches the subject of the related EAA.
type EAARevocationSubjectMatchCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper
}

// NewEAARevocationSubjectMatchCheck is the default constructor.
func NewEAARevocationSubjectMatchCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper, constraint policy.LevelRule) *EAARevocationSubjectMatchCheck {
	c := &EAARevocationSubjectMatchCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaaStatusToken: eaaStatusToken,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAARevocationSubjectMatchCheck) Process() bool {
	return c.eaaStatusToken.SubjectMatch()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAARevocationSubjectMatchCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevSubMatch
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationSubjectMatchCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevSubMatchANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationSubjectMatchCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationSubjectMatchCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
