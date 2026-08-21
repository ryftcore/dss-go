// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/checks/QEAACheck.java (DSS 6.5.RC1).
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

// QEAACheck checks whether the EAA is qualified.
type QEAACheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationEAAQualificationProcess]

	// trustServicesAtTime is the list of TrustServices declaring EAA/Q
	// status for the certificate.
	trustServicesAtTime []*diagnostic.TrustServiceWrapper
}

// NewQEAACheck is the default constructor. Port of
// QEAACheck(I18nProvider, XmlValidationEAAQualificationProcess, List, LevelRule).
func NewQEAACheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationEAAQualificationProcess],
	trustServicesAtTime []*diagnostic.TrustServiceWrapper, constraint policy.LevelRule) *QEAACheck {
	c := &QEAACheck{
		ChainItemBase:       process.NewChainItemBase(i18nProvider, result, constraint),
		trustServicesAtTime: trustServicesAtTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QEAACheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustServicesAtTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QEAACheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_QEAA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QEAACheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_QEAA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QEAACheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QEAACheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
