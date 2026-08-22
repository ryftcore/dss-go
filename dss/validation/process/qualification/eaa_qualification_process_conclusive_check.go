// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/checks/EAAQualificationProcessConclusiveCheck.java (DSS 6.5.RC1).
//
// Java's Collection<? extends XmlConstraintsConclusion> qualificationProcesses
// -> []*jaxb.XmlConstraintsConclusionContent, the Go stand-in for the
// XmlConstraintsConclusion supertype used throughout this package (see
// chain.go's header); callers pass the embedded XmlConstraintsConclusionContent
// pointer of each concrete qualification-process result.
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAAQualificationProcessConclusiveCheck verifies whether at least one of
// the EAA qualification processes concluded with a positive status.
type EAAQualificationProcessConclusiveCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationEAAQualification]

	// qualificationProcesses is the collection of EAA qualification
	// processes.
	qualificationProcesses []*jaxb.XmlConstraintsConclusionContent
}

// NewEAAQualificationProcessConclusiveCheck is the default constructor. Port
// of
// EAAQualificationProcessConclusiveCheck(Provider, XmlValidationEAAQualification, Collection, LevelRule).
func NewEAAQualificationProcessConclusiveCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlValidationEAAQualification],
	qualificationProcesses []*jaxb.XmlConstraintsConclusionContent, constraint policy.LevelRule) *EAAQualificationProcessConclusiveCheck {
	c := &EAAQualificationProcessConclusiveCheck{
		ChainItemBase:          process.NewChainItemBase(i18nProvider, result, constraint),
		qualificationProcesses: qualificationProcesses,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAQualificationProcessConclusiveCheck) Process() bool {
	if !utils.IsCollectionNotEmpty(c.qualificationProcesses) {
		return false
	}
	for _, qualificationProcess := range c.qualificationProcesses {
		if c.IsValid(qualificationProcess) {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAQualificationProcessConclusiveCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAQualConclusive
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *EAAQualificationProcessConclusiveCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAQualConclusiveANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAQualificationProcessConclusiveCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *EAAQualificationProcessConclusiveCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
