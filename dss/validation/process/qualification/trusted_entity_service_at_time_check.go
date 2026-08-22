// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/TrustedEntityServiceAtTimeCheck.java (DSS 6.5.RC1).
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

// TrustedEntityServiceAtTimeCheck verifies whether the filtered trusted
// entity services exist at the given time.
type TrustedEntityServiceAtTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateApprovalStatus]

	// trustedServicesAtTime is the list of TrustedEntityServiceWrappers at
	// control time.
	trustedServicesAtTime []*diagnostic.TrustedEntityServiceWrapper

	// validationTime is the validation time type.
	validationTime enumerations.ValidationTime
}

// NewTrustedEntityServiceAtTimeCheck is the default constructor. Port of
// TrustedEntityServiceAtTimeCheck(I18nProvider, XmlValidationCertificateApprovalStatus, List, ValidationTime, LevelRule).
func NewTrustedEntityServiceAtTimeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateApprovalStatus], trustedServicesAtTime []*diagnostic.TrustedEntityServiceWrapper,
	validationTime enumerations.ValidationTime, constraint policy.LevelRule) *TrustedEntityServiceAtTimeCheck {
	c := &TrustedEntityServiceAtTimeCheck{
		ChainItemBase:         process.NewChainItemBase(i18nProvider, result, constraint),
		trustedServicesAtTime: trustedServicesAtTime,
		validationTime:        validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustedEntityServiceAtTimeCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustedServicesAtTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustedEntityServiceAtTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_HAS_ATTIME
}

// BuildConstraintMessage builds a constraint message. Port of
// buildConstraintMessage().
func (c *TrustedEntityServiceAtTimeCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	tag, err := process.GetValidationTimeMessageTag(c.validationTime)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(c.MessageTag(), tag)
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustedEntityServiceAtTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_HAS_ATTIME_ANS
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *TrustedEntityServiceAtTimeCheck) BuildErrorMessage() *jaxb.XmlMessage {
	tag, err := process.GetValidationTimeMessageTag(c.validationTime)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(c.ErrorMessageTag(), tag)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustedEntityServiceAtTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustedEntityServiceAtTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
