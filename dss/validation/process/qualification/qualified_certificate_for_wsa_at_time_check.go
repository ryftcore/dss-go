// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/QualifiedCertificateForWSAAtTimeCheck.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QualifiedCertificateForWSAAtTimeCheck verifies whether the certificate is
// a Qualified Certificate for WebSiteAuthentication at the given time.
type QualifiedCertificateForWSAAtTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// certificateQualification is the certificate qualification validation
	// result at time.
	certificateQualification *jaxb.XmlValidationCertificateQualification
}

// NewQualifiedCertificateForWSAAtTimeCheck is the default constructor. Port
// of
// QualifiedCertificateForWSAAtTimeCheck(I18nProvider, XmlValidationQWACProcess, XmlValidationCertificateQualification, LevelRule).
func NewQualifiedCertificateForWSAAtTimeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	certificateQualification *jaxb.XmlValidationCertificateQualification, constraint policy.LevelRule) *QualifiedCertificateForWSAAtTimeCheck {
	c := &QualifiedCertificateForWSAAtTimeCheck{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		certificateQualification: certificateQualification,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QualifiedCertificateForWSAAtTimeCheck) Process() bool {
	return c.certificateQualification.CertificateQualification != nil &&
		enumerations.CertificateQualification_QCERT_FOR_WSA == c.certificateQualification.CertificateQualification.CertificateQualification()
}

// BuildConstraintMessage builds a constraint message. Port of
// buildConstraintMessage().
func (c *QualifiedCertificateForWSAAtTimeCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTag_QWAC_IS_WSA_AT_TIME, c.validationTimeMessageTag())
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *QualifiedCertificateForWSAAtTimeCheck) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTag_QWAC_IS_WSA_AT_TIME_ANS, c.validationTimeMessageTag())
}

// validationTimeMessageTag ports the ValidationProcessUtils.getValidationTimeMessageTag
// call shared by BuildConstraintMessage/BuildErrorMessage above.
func (c *QualifiedCertificateForWSAAtTimeCheck) validationTimeMessageTag() i18n.MessageTag {
	if c.certificateQualification.ValidationTime == nil {
		return ""
	}
	tag, err := process.GetValidationTimeMessageTag(c.certificateQualification.ValidationTime.ValidationTime())
	if err != nil {
		panic(err)
	}
	return tag
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *QualifiedCertificateForWSAAtTimeCheck) BuildAdditionalInfo() *string {
	var dateTime time.Time
	if c.certificateQualification.DateTime != nil {
		dateTime = time.Time(*c.certificateQualification.DateTime)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_VALIDATION_TIME, process.GetFormattedDate(&dateTime))
	return &message
}

// MessageTag/ErrorMessageTag are not overridden: Java's ChainItem source
// does not override getMessageTag()/getErrorMessageTag() either (only
// buildConstraintMessage()/buildErrorMessage(), used directly instead), so
// the ChainItemBase defaults (both "") are inherited unchanged via
// embedding.

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QualifiedCertificateForWSAAtTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QualifiedCertificateForWSAAtTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
