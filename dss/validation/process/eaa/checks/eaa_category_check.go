// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAACategoryCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// EAACategoryCheck verifies whether the EAA category claim contains one of the
// expected values.
type EAACategoryCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAACategoryCheck is the default constructor.
func NewEAACategoryCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.MultiValuesRule) *EAACategoryCheck {
	c := &EAACategoryCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		eaa:                          eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAACategoryCheck) Process() bool {
	return c.ProcessValueCheck(c.eaa.EAACategory())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAACategoryCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CAT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAACategoryCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CAT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAACategoryCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAACategoryCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
