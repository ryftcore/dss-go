// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/ContentTimestampBasicValidationCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ContentTimestampBasicValidationCheck verifies validity of a content
// timestamp.
type ContentTimestampBasicValidationCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// timestamp is the timestamp to check.
	timestamp *diagnostic.TimestampWrapper

	// timestampValidationResult is the timestamp validation result.
	timestampValidationResult *jaxb.XmlConclusion
}

// NewContentTimestampBasicValidationCheck is the default constructor. Port of
// ContentTimestampBasicValidationCheck(I18nProvider, XmlSAV, TimestampWrapper, XmlConclusion, LevelRule).
func NewContentTimestampBasicValidationCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	timestamp *diagnostic.TimestampWrapper, timestampValidationResult *jaxb.XmlConclusion,
	constraint policy.LevelRule) *ContentTimestampBasicValidationCheck {
	c := &ContentTimestampBasicValidationCheck{
		ChainItemBase:             process.NewChainItemBaseWithId(i18nProvider, result, constraint, timestamp.Id()),
		timestamp:                 timestamp,
		timestampValidationResult: timestampValidationResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *ContentTimestampBasicValidationCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeTSTBBB
}

// Process performs the check. Port of process().
func (c *ContentTimestampBasicValidationCheck) Process() bool {
	return c.IsValidConclusion(c.timestampValidationResult)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ContentTimestampBasicValidationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ICTVS
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ContentTimestampBasicValidationCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ICTVS_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ContentTimestampBasicValidationCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ContentTimestampBasicValidationCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo().
func (c *ContentTimestampBasicValidationCheck) BuildAdditionalInfo() *string {
	date := process.GetFormattedDate(c.timestamp.ProductionTime())
	timestampTypeMessageTag, err := process.GetTimestampTypeMessageTag(c.timestamp.Type())
	if err != nil {
		panic(err)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_TIMESTAMP_VALIDATION, timestampTypeMessageTag, c.timestamp.Id(), date)
	return &message
}
