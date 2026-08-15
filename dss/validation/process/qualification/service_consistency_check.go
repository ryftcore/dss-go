// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/ServiceConsistencyCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
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
// ServiceConsistencyCheck(I18nProvider, XmlValidationCertificateQualification, TrustServiceWrapper, LevelRule).
func NewServiceConsistencyCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationCertificateQualification],
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

		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS0
		return false

	}

	if !TrustServiceCheckerIsQCStatementConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS1
		return false
	}

	if !TrustServiceCheckerIsQSCDConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS3
		return false
	}

	if !TrustServiceCheckerIsQSCDStatusAsInCertConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS3A
		return false
	}

	if !TrustServiceCheckerIsPostEIDASQSCDConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS3B
		return false
	}

	if !TrustServiceCheckerIsQualifiersListKnownConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS3C
		return false
	}

	if !TrustServiceCheckerIsUsageConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS4
		return false
	}

	if !TrustServiceCheckerIsPreEIDASStatusConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS5
		return false
	}

	if !TrustServiceCheckerIsPreEIDASQualifierAndAdditionalServiceInfoConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS6
		return false
	}

	if !TrustServiceCheckerIsQualifierAndAdditionalServiceInfoConsistent(c.trustService) {
		c.errorMessage = i18n.MessageTag_QUAL_TL_SERV_CONS_ANS7
		return false
	}

	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ServiceConsistencyCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_TL_SERV_CONS
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *ServiceConsistencyCheck) ErrorMessageTag() i18n.MessageTag {
	return c.errorMessage
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *ServiceConsistencyCheck) BuildAdditionalInfo() *string {
	if c.trustService != nil && utils.IsCollectionNotEmpty(c.trustService.ServiceNames) {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_TRUST_SERVICE_NAME, c.trustService.ServiceNames[0])
		return &message
	}
	return nil
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ServiceConsistencyCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *ServiceConsistencyCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
