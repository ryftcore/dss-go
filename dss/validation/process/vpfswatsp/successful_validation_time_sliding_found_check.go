// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/pcv/checks/SuccessfulValidationTimeSlidingFoundCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SuccessfulValidationTimeSlidingFoundCheck verifies whether a successful
// Validation Time Sliding process was found.
type SuccessfulValidationTimeSlidingFoundCheck struct {
	*process.ChainItemBase[*jaxb.XmlPCV]

	// vts is the best successful VTS result; nil is Java's null.
	vts *jaxb.XmlVTS
}

// NewSuccessfulValidationTimeSlidingFoundCheck is the default constructor. Port
// of SuccessfulValidationTimeSlidingFoundCheck(Provider, XmlPCV, XmlVTS, LevelRule).
func NewSuccessfulValidationTimeSlidingFoundCheck(i18nProvider *i18n.Provider,
	result *process.Result[*jaxb.XmlPCV], vts *jaxb.XmlVTS,
	constraint policy.LevelRule) *SuccessfulValidationTimeSlidingFoundCheck {
	c := &SuccessfulValidationTimeSlidingFoundCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		vts:           vts,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *SuccessfulValidationTimeSlidingFoundCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeVTS
}

// Process performs the check. Port of process().
func (c *SuccessfulValidationTimeSlidingFoundCheck) Process() bool {
	return c.vts != nil && c.IsValid(&c.vts.XmlConstraintsConclusionContent)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SuccessfulValidationTimeSlidingFoundCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagPCVICCSVTSF
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SuccessfulValidationTimeSlidingFoundCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagPCVICCSVTSFANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SuccessfulValidationTimeSlidingFoundCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *SuccessfulValidationTimeSlidingFoundCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoPOE
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(), whose null result leaves the element absent.
//
// Java passes getTrustAnchor() straight into MessageFormat; the generated Go
// member is a *string, whose nil would be the literal "null" there - but
// ValidationTimeSliding#addAdditionalInfo always sets it, so the nil case is
// unreachable and renders as the empty string here.
func (c *SuccessfulValidationTimeSlidingFoundCheck) BuildAdditionalInfo() *string {
	if c.vts != nil {
		trustAnchor := ""
		if c.vts.TrustAnchor != nil {
			trustAnchor = *c.vts.TrustAnchor
		}
		// XmlVTS#getControlTime(): the generated member is a *XSDateTime, whose
		// nil is Java's null.
		var controlTime *time.Time
		if c.vts.ControlTime != nil {
			t := c.vts.ControlTime.Time()
			controlTime = &t
		}
		message := c.I18nProvider.GetMessage(i18n.MessageTagControlTimeWithTrustAnchor, trustAnchor,
			process.GetFormattedDate(controlTime))
		return &message
	}
	return nil
}
