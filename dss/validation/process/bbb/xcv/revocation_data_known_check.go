// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationDataKnownCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationDataKnownCheck checks if the revocation status is known.
type RevocationDataKnownCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// revocationData is the revocation data to check.
	revocationData *diagnostic.CertificateRevocationWrapper
}

// NewRevocationDataKnownCheck is the default constructor. Port of
// RevocationDataKnownCheck(I18nProvider, XmlRAC, CertificateRevocationWrapper, LevelRule).
func NewRevocationDataKnownCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRAC],
	revocationData *diagnostic.CertificateRevocationWrapper, constraint policy.LevelRule) *RevocationDataKnownCheck {
	c := &RevocationDataKnownCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationDataKnownCheck) Process() bool {
	return c.revocationData.IsKnown()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationDataKnownCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCUKN
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationDataKnownCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCUKN_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationDataKnownCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *RevocationDataKnownCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
