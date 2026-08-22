// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/QualifiedCheck.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QualifiedCheck checks whether the certificate is qualified at validation
// time.
type QualifiedCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// qualifiedStatus is the certificate qualification status.
	qualifiedStatus enumerations.CertificateQualifiedStatus

	// validationTime is the validation time type.
	validationTime enumerations.ValidationTime
}

// NewQualifiedCheck is the default constructor. Port of
// QualifiedCheck(I18nProvider, XmlValidationCertificateQualification, CertificateQualifiedStatus, ValidationTime, LevelRule).
func NewQualifiedCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationCertificateQualification],
	qualifiedStatus enumerations.CertificateQualifiedStatus, validationTime enumerations.ValidationTime,
	constraint policy.LevelRule) *QualifiedCheck {
	c := &QualifiedCheck{
		ChainItemBase:   process.NewChainItemBase(i18nProvider, result, constraint),
		qualifiedStatus: qualifiedStatus,
		validationTime:  validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QualifiedCheck) Process() bool {
	return enumerations.CertificateQualifiedStatusIsQC(c.qualifiedStatus)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QualifiedCheck) MessageTag() i18n.MessageTag {
	switch c.validationTime {
	case enumerations.ValidationTimeBESTSignatureTime:
		return i18n.MessageTag_QUAL_QC_AT_ST
	case enumerations.ValidationTimeCertificateIssuanceTime:
		return i18n.MessageTag_QUAL_QC_AT_CC
	case enumerations.ValidationTimeValidationTime:
		return i18n.MessageTag_QUAL_QC_AT_VT
	default:
		panic(fmt.Sprintf("Unsupported time %s", c.validationTime))
	}
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QualifiedCheck) ErrorMessageTag() i18n.MessageTag {
	switch c.validationTime {
	case enumerations.ValidationTimeBESTSignatureTime:
		return i18n.MessageTag_QUAL_QC_AT_ST_ANS
	case enumerations.ValidationTimeCertificateIssuanceTime:
		return i18n.MessageTag_QUAL_QC_AT_CC_ANS
	case enumerations.ValidationTimeValidationTime:
		return i18n.MessageTag_QUAL_QC_AT_VT_ANS
	default:
		panic(fmt.Sprintf("Unsupported time %s", c.validationTime))
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QualifiedCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QualifiedCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
