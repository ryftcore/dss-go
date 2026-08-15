// Ported from dss-validation/.../validation/process/bbb/fc/checks/ContainerTypeCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// ContainerTypeCheck checks if the container type is acceptable.
type ContainerTypeCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*drjaxb.XmlFC]

	containerType enumerations.ASiCContainerType
}

// NewContainerTypeCheck is the default constructor.
func NewContainerTypeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	containerType enumerations.ASiCContainerType, constraint policy.MultiValuesRule) *ContainerTypeCheck {
	c := &ContainerTypeCheck{containerType: containerType}
	c.AbstractMultiValuesCheckItem = bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *ContainerTypeCheck) Process() bool {
	// Java calls containerType.toString(), and ASiCContainerType overrides
	// toString() to replace '_' with '-' ("ASiC-E"), which is what the policy's
	// AcceptableContainerTypes ids carry. A plain string(...) conversion would
	// yield the enum name ("ASiC_E") and never match.
	return c.ProcessValueCheck(c.containerType.String())
}

// MessageTag returns the constraint message i18n key.
func (c *ContainerTypeCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_IECTF }

// ErrorMessageTag returns the error message i18n key.
func (c *ContainerTypeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IECTF_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *ContainerTypeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *ContainerTypeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
