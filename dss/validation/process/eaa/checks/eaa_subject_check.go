// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAASubjectCheck.java (DSS 6.5.RC1).
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

// EAASubjectCheck verifies whether the EAA was issued to an expected subject.
type EAASubjectCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAASubjectCheck is the default constructor.
func NewEAASubjectCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.MultiValuesRule) *EAASubjectCheck {
	c := &EAASubjectCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		eaa:                          eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAASubjectCheck) Process() bool {
	return c.ProcessValueCheck(c.eaa.EAASubject())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAASubjectCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAASub
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAASubjectCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAASubANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAASubjectCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAASubjectCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
