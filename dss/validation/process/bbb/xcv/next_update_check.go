// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rfc/checks/NextUpdateCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// NextUpdateCheck checks if the nextUpdate is present.
type NextUpdateCheck struct {
	*process.ChainItemBase[*jaxb.XmlRFC]

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper
}

// NewNextUpdateCheck is the default constructor. Port of
// NextUpdateCheck(I18nProvider, XmlRFC, RevocationWrapper, LevelRule).
func NewNextUpdateCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRFC],
	revocationData *diagnostic.RevocationWrapper, constraint policy.LevelRule) *NextUpdateCheck {
	c := &NextUpdateCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *NextUpdateCheck) Process() bool {
	if c.revocationData != nil {
		return c.revocationData.NextUpdate() != nil
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *NextUpdateCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_RFC_NUP
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *NextUpdateCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_RFC_NUP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *NextUpdateCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *NextUpdateCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_TRY_LATER
}
