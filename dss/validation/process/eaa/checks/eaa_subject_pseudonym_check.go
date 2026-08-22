// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAASubjectPseudonymCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// EAASubjectPseudonymCheck verifies whether the EAA subject pseudonym claim
// contains one of the expected values.
type EAASubjectPseudonymCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAASubjectPseudonymCheck is the default constructor.
func NewEAASubjectPseudonymCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.MultiValuesRule) *EAASubjectPseudonymCheck {
	c := &EAASubjectPseudonymCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		eaa:                          eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAASubjectPseudonymCheck) Process() bool {
	return c.ProcessValueCheck(c.eaa.HolderPseudonym())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAASubjectPseudonymCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAASubPSE
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAASubjectPseudonymCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAASubPSEANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAASubjectPseudonymCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAASubjectPseudonymCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
