// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/ServiceConsistencyCheck.java (DSS 6.5.RC1).
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

// ServiceConsistencyCheck checks if the Trusted Service is consistent.
type ServiceConsistencyCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// trustService is the Trusted Service to check.
	trustService *diagnostic.TrustServiceWrapper

	// errorMessage is the internal cached error message, if applicable.
	errorMessage i18n.MessageTag
}

// NewServiceConsistencyCheck is the default constructor. Port of
// ServiceConsistencyCheck(Provider, XmlValidationCertificateQualification, TrustServiceWrapper, LevelRule).
func NewServiceConsistencyCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlValidationCertificateQualification],
	trustService *diagnostic.TrustServiceWrapper, constraint policy.LevelRule) *ServiceConsistencyCheck {
	c := &ServiceConsistencyCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		trustService:  trustService,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ServiceConsistencyCheck) Process() bool {

	if c.trustService == nil {

		c.errorMessage = i18n.MessageTagQualTLServCONSANS0
		return false

	}

	if !TrustServiceCheckerIsQCStatementConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTagQualTLServCONSANS1
		return false
	}

	if !TrustServiceCheckerIsQSCDConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTagQualTLServCONSANS3
		return false
	}

	if !TrustServiceCheckerIsQSCDStatusAsInCertConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTagQualTLServCONSANS3A
		return false
	}

	if !TrustServiceCheckerIsPostEIDASQSCDConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTagQualTLServCONSANS3B
		return false
	}

	if !TrustServiceCheckerIsQualifiersListKnownConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTagQualTLServCONSANS3C
		return false
	}

	if !TrustServiceCheckerIsUsageConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTagQualTLServCONSANS4
		return false
	}

	if !TrustServiceCheckerIsPreEIDASStatusConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTagQualTLServCONSANS5
		return false
	}

	if !TrustServiceCheckerIsPreEIDASQualifierAndAdditionalServiceInfoConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTagQualTLServCONSANS6
		return false
	}

	if !TrustServiceCheckerIsQualifierAndAdditionalServiceInfoConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTagQualTLServCONSANS7
		return false
	}

	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ServiceConsistencyCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualTLServCONS
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *ServiceConsistencyCheck) ErrorMessageTag() i18n.MessageTag {
	return c.errorMessage
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *ServiceConsistencyCheck) BuildAdditionalInfo() *string {
	if c.trustService != nil && utils.IsCollectionNotEmpty(c.trustService.ServiceNames) {
		message := c.I18nProvider.GetMessage(i18n.MessageTagTrustServiceName, c.trustService.ServiceNames[0])
		return &message
	}
	return nil
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ServiceConsistencyCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *ServiceConsistencyCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
