// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/CertificateTypeCheck.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateTypeCheck checks if the certificate type has been identified at
// the given time.
type CertificateTypeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// certType is the CertificateType in question.
	certType enumerations.CertificateType

	// validationTime is the used validation time.
	validationTime enumerations.ValidationTime
}

// NewCertificateTypeCheck is the default constructor. Port of
// CertificateTypeCheck(Provider, XmlValidationCertificateQualification, CertificateType, ValidationTime, LevelRule).
func NewCertificateTypeCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlValidationCertificateQualification],
	certType enumerations.CertificateType, validationTime enumerations.ValidationTime,
	constraint policy.LevelRule) *CertificateTypeCheck {
	c := &CertificateTypeCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		certType:       certType,
		validationTime: validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateTypeCheck) Process() bool {
	return enumerations.CertificateTypeUnknown != c.certType
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateTypeCheck) MessageTag() i18n.MessageTag {
	switch c.validationTime {
	case enumerations.ValidationTimeBESTSignatureTime:
		return i18n.MessageTagQualCertTypeAtST
	case enumerations.ValidationTimeCertificateIssuanceTime:
		return i18n.MessageTagQualCertTypeAtCC
	case enumerations.ValidationTimeValidationTime:
		return i18n.MessageTagQualCertTypeAtVT
	default:
		panic(fmt.Sprintf("Unsupported time %s", c.validationTime))
	}
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateTypeCheck) ErrorMessageTag() i18n.MessageTag {
	switch c.validationTime {
	case enumerations.ValidationTimeBESTSignatureTime:
		return i18n.MessageTagQualCertTypeAtSTANS
	case enumerations.ValidationTimeCertificateIssuanceTime:
		return i18n.MessageTagQualCertTypeAtCCANS
	case enumerations.ValidationTimeValidationTime:
		return i18n.MessageTagQualCertTypeAtVTANS
	default:
		panic(fmt.Sprintf("Unsupported time %s", c.validationTime))
	}
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *CertificateTypeCheck) BuildAdditionalInfo() *string {
	if enumerations.CertificateTypeUnknown != c.certType {
		message := c.I18nProvider.GetMessage(i18n.MessageTagCertificateType, c.certType.Label())
		return &message
	}
	return nil
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateTypeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *CertificateTypeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
