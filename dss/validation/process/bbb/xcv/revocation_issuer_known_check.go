// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationIssuerKnownCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationIssuerKnownCheck verifies whether the signing certificate of the
// revocation data has been found.
type RevocationIssuerKnownCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper
}

// NewRevocationIssuerKnownCheck is the default constructor. Port of
// RevocationIssuerKnownCheck(I18nProvider, XmlRAC, RevocationWrapper, LevelRule).
func NewRevocationIssuerKnownCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRAC],
	revocationData *diagnostic.RevocationWrapper, constraint policy.LevelRule) *RevocationIssuerKnownCheck {
	c := &RevocationIssuerKnownCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationIssuerKnownCheck) Process() bool {
	return c.revocationData.SigningCertificate() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationIssuerKnownCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_ISSUER_KNOWN
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationIssuerKnownCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_ISSUER_KNOWN_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationIssuerKnownCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *RevocationIssuerKnownCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}
