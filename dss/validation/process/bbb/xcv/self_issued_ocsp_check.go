// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/SelfIssuedOCSPCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SelfIssuedOCSPCheck checks if the certificate in question is not present in
// the OCSP's certificate chain.
type SelfIssuedOCSPCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// certificateWrapper is the certificate in question.
	certificateWrapper *diagnostic.CertificateWrapper

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper
}

// NewSelfIssuedOCSPCheck is the default constructor. Port of
// SelfIssuedOCSPCheck(I18nProvider, XmlRAC, CertificateWrapper, RevocationWrapper, LevelRule).
func NewSelfIssuedOCSPCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRAC],
	certificateWrapper *diagnostic.CertificateWrapper, revocationData *diagnostic.RevocationWrapper,
	constraint policy.LevelRule) *SelfIssuedOCSPCheck {
	c := &SelfIssuedOCSPCheck{
		ChainItemBase:      process.NewChainItemBase(i18nProvider, result, constraint),
		certificateWrapper: certificateWrapper,
		revocationData:     revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SelfIssuedOCSPCheck) Process() bool {
	for _, certificate := range c.revocationData.CertificateChain() {
		if c.certificateWrapper.Id() == certificate.Id() {
			return false
		}
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SelfIssuedOCSPCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_SELF_ISSUED_OCSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SelfIssuedOCSPCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_SELF_ISSUED_OCSP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SelfIssuedOCSPCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SelfIssuedOCSPCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}
