// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificatePolicyQualifiedIdsCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificatePolicyQualifiedIdsCheck checks if the certificate policies
// contain a Qualified identifier(s).
type CertificatePolicyQualifiedIdsCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificatePolicyQualifiedIdsCheck is the default constructor. Port of
// CertificatePolicyQualifiedIdsCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewCertificatePolicyQualifiedIdsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificatePolicyQualifiedIdsCheck {
	c := &CertificatePolicyQualifiedIdsCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
//
// This check only uses the certificate (not the TL).
func (c *CertificatePolicyQualifiedIdsCheck) Process() bool {
	isQCP := process.IsQCP(c.certificate)
	isQCPPlus := process.IsQCPPlus(c.certificate)
	return isQCP || isQCPPlus
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificatePolicyQualifiedIdsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCIQC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificatePolicyQualifiedIdsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCIQCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificatePolicyQualifiedIdsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificatePolicyQualifiedIdsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
