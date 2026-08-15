// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/IsAbleToSelectOneTrustService.java (DSS 6.5.RC1).
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

// IsAbleToSelectOneTrustService checks whether the validator was able to
// select one TrustService (in condition that there is no conflict with
// other TrustServices).
type IsAbleToSelectOneTrustService struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// trustServicesAtTime is the list of selected TrustServices.
	trustServicesAtTime []*diagnostic.TrustServiceWrapper
}

// NewIsAbleToSelectOneTrustService is the default constructor. Port of
// IsAbleToSelectOneTrustService(I18nProvider, XmlValidationCertificateQualification, List, LevelRule).
func NewIsAbleToSelectOneTrustService(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateQualification], trustServicesAtTime []*diagnostic.TrustServiceWrapper,
	constraint policy.LevelRule) *IsAbleToSelectOneTrustService {
	c := &IsAbleToSelectOneTrustService{
		ChainItemBase:       process.NewChainItemBase(i18nProvider, result, constraint),
		trustServicesAtTime: trustServicesAtTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *IsAbleToSelectOneTrustService) Process() bool {
	return utils.CollectionSize(c.trustServicesAtTime) == 1
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *IsAbleToSelectOneTrustService) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_ONLY_ONE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *IsAbleToSelectOneTrustService) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_ONLY_ONE_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *IsAbleToSelectOneTrustService) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *IsAbleToSelectOneTrustService) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
