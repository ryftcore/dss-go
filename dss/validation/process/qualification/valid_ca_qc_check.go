// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/ValidCAQCCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ValidCAQCCheck checks if an acceptable Trust Service for a qualified
// certificate issuance found.
type ValidCAQCCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// trustService is the TrustServiceWrapper.
	trustService *diagnostic.TrustServiceWrapper
}

// NewValidCAQCCheck is the default constructor. Port of
// ValidCAQCCheck(I18nProvider, XmlValidationCertificateQualification, TrustServiceWrapper, LevelRule).
func NewValidCAQCCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationCertificateQualification],
	trustService *diagnostic.TrustServiceWrapper, constraint policy.LevelRule) *ValidCAQCCheck {
	c := &ValidCAQCCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		trustService:  trustService,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ValidCAQCCheck) Process() bool {
	return c.trustService != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ValidCAQCCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasValidCAQC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *ValidCAQCCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasValidCAQCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ValidCAQCCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *ValidCAQCCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
