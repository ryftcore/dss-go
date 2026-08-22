// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/SelfIssuedOCSPCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
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
// SelfIssuedOCSPCheck(Provider, XmlRAC, CertificateWrapper, RevocationWrapper, LevelRule).
func NewSelfIssuedOCSPCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlRAC],
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
	return i18n.MessageTagBBBXCVRevocSelfIssuedOCSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SelfIssuedOCSPCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVRevocSelfIssuedOCSPANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SelfIssuedOCSPCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SelfIssuedOCSPCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
