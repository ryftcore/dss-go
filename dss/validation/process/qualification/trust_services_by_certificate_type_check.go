// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/TrustServicesByCertificateTypeCheck.java (DSS 6.5.RC1).
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

// TrustServicesByCertificateTypeCheck checks if a trust service
// corresponding to the certificate type has been found.
type TrustServicesByCertificateTypeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// trustServices is the list of TrustServiceWrappers at control time.
	trustServices []*diagnostic.TrustServiceWrapper
}

// NewTrustServicesByCertificateTypeCheck is the default constructor. Port of
// TrustServicesByCertificateTypeCheck(I18nProvider, XmlValidationCertificateQualification, List, LevelRule).
func NewTrustServicesByCertificateTypeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateQualification], trustServices []*diagnostic.TrustServiceWrapper,
	constraint policy.LevelRule) *TrustServicesByCertificateTypeCheck {
	c := &TrustServicesByCertificateTypeCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		trustServices: trustServices,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustServicesByCertificateTypeCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustServices)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustServicesByCertificateTypeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_TS_CERT_TYPE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustServicesByCertificateTypeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_TS_CERT_TYPE_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustServicesByCertificateTypeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustServicesByCertificateTypeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
