// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAAdministrativePeriodNotExpiredCheck.java (DSS 6.5.RC1).
package checks

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAAAdministrativePeriodNotExpiredCheck verifies whether the validation time
// is within the EAA administrative validity period range.
type EAAAdministrativePeriodNotExpiredCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper

	// validationTime is the EAA validation time.
	validationTime time.Time
}

// NewEAAAdministrativePeriodNotExpiredCheck is the default constructor.
func NewEAAAdministrativePeriodNotExpiredCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, validationTime time.Time, constraint policy.LevelRule) *EAAAdministrativePeriodNotExpiredCheck {
	c := &EAAAdministrativePeriodNotExpiredCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:            eaaWrapper,
		validationTime: validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAAdministrativePeriodNotExpiredCheck) Process() bool {
	return c.notAdministrativePeriodBefore() && c.notAdministrativePeriodAtOrAfter()
}

// notAdministrativePeriodBefore ports the private
// notAdministrativePeriodBefore().
//
// Same logic is applied as for IETF RFC 7519 "nbf" mutatis mutandis.
func (c *EAAAdministrativePeriodNotExpiredCheck) notAdministrativePeriodBefore() bool {
	return c.eaa.AdministrativeIssuanceDate() != nil && !c.validationTime.Before(*c.eaa.AdministrativeIssuanceDate())
}

// notAdministrativePeriodAtOrAfter ports the private
// notAdministrativePeriodAtOrAfter().
//
// Same logic is applied as for IETF RFC 7519 "exp" mutatis mutandis.
func (c *EAAAdministrativePeriodNotExpiredCheck) notAdministrativePeriodAtOrAfter() bool {
	return c.eaa.AdministrativeExpirationDate() != nil && c.validationTime.Before(*c.eaa.AdministrativeExpirationDate())
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *EAAAdministrativePeriodNotExpiredCheck) BuildAdditionalInfo() *string {
	if !c.notAdministrativePeriodBefore() || !c.notAdministrativePeriodAtOrAfter() {
		message := c.I18nProvider.GetMessage(i18n.MessageTagEAAVTIAVRValidity,
			process.GetFormattedDate(&c.validationTime), process.GetFormattedDate(c.eaa.AdministrativeIssuanceDate()),
			process.GetFormattedDate(c.eaa.AdministrativeExpirationDate()))
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAAdministrativePeriodNotExpiredCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAVTIAVR
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAAdministrativePeriodNotExpiredCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAVTIAVRANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAAdministrativePeriodNotExpiredCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAAdministrativePeriodNotExpiredCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationOutOfBoundsNoPOE
}
