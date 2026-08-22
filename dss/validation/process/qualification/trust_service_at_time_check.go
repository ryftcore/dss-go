// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/TrustServiceAtTimeCheck.java (DSS 6.5.RC1).
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

// TrustServiceAtTimeCheck checks if a corresponding Trust Service found
// valid at control time.
type TrustServiceAtTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// trustServicesAtTime is the list of TrustServiceWrappers at control time.
	trustServicesAtTime []*diagnostic.TrustServiceWrapper

	// validationTime is the validation time type.
	validationTime enumerations.ValidationTime
}

// NewTrustServiceAtTimeCheck is the default constructor. Port of
// TrustServiceAtTimeCheck(Provider, XmlValidationCertificateQualification, List, ValidationTime, LevelRule).
func NewTrustServiceAtTimeCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlValidationCertificateQualification],
	trustServicesAtTime []*diagnostic.TrustServiceWrapper, validationTime enumerations.ValidationTime,
	constraint policy.LevelRule) *TrustServiceAtTimeCheck {
	c := &TrustServiceAtTimeCheck{
		ChainItemBase:       process.NewChainItemBase(i18nProvider, result, constraint),
		trustServicesAtTime: trustServicesAtTime,
		validationTime:      validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustServiceAtTimeCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustServicesAtTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustServiceAtTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasAtTime
}

// BuildConstraintMessage builds a constraint message. Port of
// buildConstraintMessage().
func (c *TrustServiceAtTimeCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	tag, err := process.GetValidationTimeMessageTag(c.validationTime)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(c.MessageTag(), tag)
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustServiceAtTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasAtTimeANS
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *TrustServiceAtTimeCheck) BuildErrorMessage() *jaxb.XmlMessage {
	tag, err := process.GetValidationTimeMessageTag(c.validationTime)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(c.ErrorMessageTag(), tag)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustServiceAtTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustServiceAtTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
