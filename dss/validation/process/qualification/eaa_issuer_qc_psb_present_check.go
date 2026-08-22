// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/checks/EAAIssuerQcPSBPresentCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAAIssuerQcPSBPresentCheck verifies presence of a QcPSB QcStatement.
type EAAIssuerQcPSBPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationEAAQualificationProcess]

	// signingCertificate is the signing-certificate of the EAA signature.
	signingCertificate *diagnostic.CertificateWrapper
}

// NewEAAIssuerQcPSBPresentCheck is the default constructor. Port of
// EAAIssuerQcPSBPresentCheck(I18nProvider, XmlValidationEAAQualificationProcess, CertificateWrapper, LevelRule).
func NewEAAIssuerQcPSBPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationEAAQualificationProcess],
	signingCertificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *EAAIssuerQcPSBPresentCheck {
	c := &EAAIssuerQcPSBPresentCheck{
		ChainItemBase:      process.NewChainItemBase(i18nProvider, result, constraint),
		signingCertificate: signingCertificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process(). TODO: check country ?
func (c *EAAIssuerQcPSBPresentCheck) Process() bool {
	return c.signingCertificate.QcPSB() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAIssuerQcPSBPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAQCPSB
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *EAAIssuerQcPSBPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAQCPSBANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAIssuerQcPSBPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *EAAIssuerQcPSBPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
