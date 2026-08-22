// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/CertificateTypeCoverageCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateTypeCoverageCheck verifies if a TrustService(s) issuing the
// certificate have been found.
type CertificateTypeCoverageCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// trustServicesAtTime is the list of TrustServices issuing the
	// certificate in question.
	trustServicesAtTime []*diagnostic.TrustServiceWrapper
}

// NewCertificateTypeCoverageCheck is the default constructor. Port of
// CertificateTypeCoverageCheck(I18nProvider, XmlValidationCertificateQualification, List, LevelRule).
func NewCertificateTypeCoverageCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateQualification], trustServicesAtTime []*diagnostic.TrustServiceWrapper,
	constraint policy.LevelRule) *CertificateTypeCoverageCheck {
	c := &CertificateTypeCoverageCheck{
		ChainItemBase:       process.NewChainItemBase(i18nProvider, result, constraint),
		trustServicesAtTime: trustServicesAtTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateTypeCoverageCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustServicesAtTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateTypeCoverageCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_CERT_TYPE_COVERAGE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateTypeCoverageCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_CERT_TYPE_COVERAGE_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateTypeCoverageCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *CertificateTypeCoverageCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
