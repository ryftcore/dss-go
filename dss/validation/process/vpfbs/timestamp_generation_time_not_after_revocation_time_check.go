// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/TimestampGenerationTimeNotAfterRevocationTimeCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TimestampGenerationTimeNotAfterRevocationTimeCheck checks if the generation
// time of a content timestamp is not after the revocation time of a
// signature's signing certificate.
type TimestampGenerationTimeNotAfterRevocationTimeCheck[T any] struct {
	*process.ChainItemBase[T]

	// contentTimestamp is the content timestamp.
	contentTimestamp *diagnostic.TimestampWrapper

	// signingCertificateRevocationTime is the revocation time of the signing
	// certificate; nil is Java's null.
	signingCertificateRevocationTime *time.Time
}

// NewTimestampGenerationTimeNotAfterRevocationTimeCheck is the default
// constructor. Port of
// TimestampGenerationTimeNotAfterRevocationTimeCheck(I18nProvider, T, TimestampWrapper, Date, LevelRule).
func NewTimestampGenerationTimeNotAfterRevocationTimeCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	contentTimestamp *diagnostic.TimestampWrapper, signingCertificateRevocationTime *time.Time,
	constraint policy.LevelRule) *TimestampGenerationTimeNotAfterRevocationTimeCheck[T] {
	c := &TimestampGenerationTimeNotAfterRevocationTimeCheck[T]{
		ChainItemBase:                    process.NewChainItemBaseWithId(i18nProvider, result, constraint, contentTimestamp.Id()),
		contentTimestamp:                 contentTimestamp,
		signingCertificateRevocationTime: signingCertificateRevocationTime,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *TimestampGenerationTimeNotAfterRevocationTimeCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeCNTTSTBBB
}

// Process performs the check. Port of process().
func (c *TimestampGenerationTimeNotAfterRevocationTimeCheck[T]) Process() bool {
	productionTime := c.contentTimestamp.ProductionTime()
	return productionTime != nil && c.signingCertificateRevocationTime != nil &&
		!productionTime.After(*c.signingCertificateRevocationTime)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TimestampGenerationTimeNotAfterRevocationTimeCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *TimestampGenerationTimeNotAfterRevocationTimeCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationRevoked
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TimestampGenerationTimeNotAfterRevocationTimeCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_ICTGTNASCRT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TimestampGenerationTimeNotAfterRevocationTimeCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_ICTGTNASCRT_ANS
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *TimestampGenerationTimeNotAfterRevocationTimeCheck[T]) BuildAdditionalInfo() *string {
	tstGenerationTime := " ? "
	if productionTime := c.contentTimestamp.ProductionTime(); productionTime != nil {
		tstGenerationTime = process.GetFormattedDate(productionTime)
	}
	revocationTime := " ? "
	if c.signingCertificateRevocationTime != nil {
		revocationTime = process.GetFormattedDate(c.signingCertificateRevocationTime)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_TIMESTAMP_AND_REVOCATION_TIME, c.contentTimestamp.Id(), tstGenerationTime, revocationTime)
	return &message
}
