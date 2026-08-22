// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateIssuedToLegalPersonCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateIssuedToLegalPersonCheck checks if the certificate has been
// issued to a legal person.
type CertificateIssuedToLegalPersonCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateIssuedToLegalPersonCheck is the default constructor. Port of
// CertificateIssuedToLegalPersonCheck(Provider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewCertificateIssuedToLegalPersonCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificateIssuedToLegalPersonCheck {
	c := &CertificateIssuedToLegalPersonCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
//
// This check only uses the certificate (not the TL).
func (c *CertificateIssuedToLegalPersonCheck) Process() bool {
	return process.IsLegal(c.certificate)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateIssuedToLegalPersonCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCIITLP
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateIssuedToLegalPersonCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCIITLPANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateIssuedToLegalPersonCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateIssuedToLegalPersonCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
