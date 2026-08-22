// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/QSCDCheck.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QSCDCheck checks whether the certificate was for QSCD at validation time.
type QSCDCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// qscdStatus is the certificate QSCD status at validation time.
	qscdStatus enumerations.QSCDStatus

	// validationTime is the validation time type.
	validationTime enumerations.ValidationTime
}

// NewQSCDCheck is the default constructor. Port of
// QSCDCheck(I18nProvider, XmlValidationCertificateQualification, QSCDStatus, ValidationTime, LevelRule).
func NewQSCDCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationCertificateQualification],
	qscdStatus enumerations.QSCDStatus, validationTime enumerations.ValidationTime,
	constraint policy.LevelRule) *QSCDCheck {
	c := &QSCDCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		qscdStatus:     qscdStatus,
		validationTime: validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QSCDCheck) Process() bool {
	return enumerations.QSCDStatusIsQSCD(c.qscdStatus)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QSCDCheck) MessageTag() i18n.MessageTag {
	switch c.validationTime {
	case enumerations.ValidationTimeBESTSignatureTime:
		return i18n.MessageTagQualQSCDAtST
	case enumerations.ValidationTimeCertificateIssuanceTime:
		return i18n.MessageTagQualQSCDAtCC
	case enumerations.ValidationTimeValidationTime:
		return i18n.MessageTagQualQSCDAtVT
	default:
		panic(fmt.Sprintf("Unsupported time %s", c.validationTime))
	}
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QSCDCheck) ErrorMessageTag() i18n.MessageTag {
	switch c.validationTime {
	case enumerations.ValidationTimeBESTSignatureTime:
		return i18n.MessageTagQualQSCDAtSTANS
	case enumerations.ValidationTimeCertificateIssuanceTime:
		return i18n.MessageTagQualQSCDAtCCANS
	case enumerations.ValidationTimeValidationTime:
		return i18n.MessageTagQualQSCDAtVTANS
	default:
		panic(fmt.Sprintf("Unsupported time %s", c.validationTime))
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QSCDCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QSCDCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
