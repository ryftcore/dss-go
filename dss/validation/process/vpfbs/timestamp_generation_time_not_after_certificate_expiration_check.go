// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/TimestampGenerationTimeNotAfterCertificateExpirationCheck.java (DSS 6.5.RC1).
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

// TimestampGenerationTimeNotAfterCertificateExpirationCheck verifies if the
// generation time of a content timestamp is not after the certificate's
// expiration time.
type TimestampGenerationTimeNotAfterCertificateExpirationCheck[T any] struct {
	*process.ChainItemBase[T]

	// contentTimestamp is the content timestamp.
	contentTimestamp *diagnostic.TimestampWrapper

	// signingCertificateNotAfter is the signing certificate's notAfter time;
	// nil is Java's null.
	signingCertificateNotAfter *time.Time
}

// NewTimestampGenerationTimeNotAfterCertificateExpirationCheck is the default
// constructor. Port of
// TimestampGenerationTimeNotAfterCertificateExpirationCheck(I18nProvider, T, TimestampWrapper, Date, LevelRule).
func NewTimestampGenerationTimeNotAfterCertificateExpirationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	contentTimestamp *diagnostic.TimestampWrapper, signingCertificateNotAfter *time.Time,
	constraint policy.LevelRule) *TimestampGenerationTimeNotAfterCertificateExpirationCheck[T] {
	c := &TimestampGenerationTimeNotAfterCertificateExpirationCheck[T]{
		ChainItemBase:              process.NewChainItemBaseWithId(i18nProvider, result, constraint, contentTimestamp.Id()),
		contentTimestamp:           contentTimestamp,
		signingCertificateNotAfter: signingCertificateNotAfter,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *TimestampGenerationTimeNotAfterCertificateExpirationCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeCNTTSTBBB
}

// Process performs the check. Port of process().
func (c *TimestampGenerationTimeNotAfterCertificateExpirationCheck[T]) Process() bool {
	productionTime := c.contentTimestamp.ProductionTime()
	return productionTime != nil && c.signingCertificateNotAfter != nil && !productionTime.After(*c.signingCertificateNotAfter)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TimestampGenerationTimeNotAfterCertificateExpirationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *TimestampGenerationTimeNotAfterCertificateExpirationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationExpired
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TimestampGenerationTimeNotAfterCertificateExpirationCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVICTGTNASCET
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TimestampGenerationTimeNotAfterCertificateExpirationCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVICTGTNASCETANS
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *TimestampGenerationTimeNotAfterCertificateExpirationCheck[T]) BuildAdditionalInfo() *string {
	tstGenerationTime := " ? "
	if productionTime := c.contentTimestamp.ProductionTime(); productionTime != nil {
		tstGenerationTime = process.GetFormattedDate(productionTime)
	}
	certificateNotAfter := " ? "
	if c.signingCertificateNotAfter != nil {
		certificateNotAfter = process.GetFormattedDate(c.signingCertificateNotAfter)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagTimestampAndRevocationTime, c.contentTimestamp.Id(), tstGenerationTime, certificateNotAfter)
	return &message
}
