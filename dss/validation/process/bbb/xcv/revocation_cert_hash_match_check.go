// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationCertHashMatchCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationCertHashMatchCheck checks if the revocation's certHash matches.
type RevocationCertHashMatchCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper
}

// NewRevocationCertHashMatchCheck is the default constructor. Port of
// RevocationCertHashMatchCheck(I18nProvider, XmlRAC, RevocationWrapper, LevelRule).
func NewRevocationCertHashMatchCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRAC],
	revocationData *diagnostic.RevocationWrapper, constraint policy.LevelRule) *RevocationCertHashMatchCheck {
	c := &RevocationCertHashMatchCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationCertHashMatchCheck) Process() bool {
	return c.revocationData.IsCertHashExtensionMatch()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationCertHashMatchCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_CERT_HASH_MATCH
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationCertHashMatchCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_CERT_HASH_MATCH_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationCertHashMatchCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *RevocationCertHashMatchCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
