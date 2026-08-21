// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationCertHashPresenceCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationCertHashPresenceCheck checks if the revocation's certHash is
// present.
type RevocationCertHashPresenceCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper
}

// NewRevocationCertHashPresenceCheck is the default constructor. Port of
// RevocationCertHashPresenceCheck(I18nProvider, XmlRAC, RevocationWrapper, LevelRule).
func NewRevocationCertHashPresenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRAC],
	revocationData *diagnostic.RevocationWrapper, constraint policy.LevelRule) *RevocationCertHashPresenceCheck {
	c := &RevocationCertHashPresenceCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationCertHashPresenceCheck) Process() bool {
	return c.revocationData.IsCertHashExtensionPresent()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationCertHashPresenceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_CERT_HASH_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationCertHashPresenceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_CERT_HASH_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationCertHashPresenceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *RevocationCertHashPresenceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}
