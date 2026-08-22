// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck
// checks if the generation time of a content timestamp is not after the
// expiration time of cryptographic constraints concerned by the failure.
type TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T any] struct {
	*process.ChainItemBase[T]

	// contentTimestamp is the content timestamp.
	contentTimestamp *diagnostic.TimestampWrapper

	// cryptographicValidation is the cryptographic validation result summary.
	cryptographicValidation *jaxb.XmlCryptographicValidation
}

// NewTimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck is
// the default constructor. Port of
// TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck(Provider, T, TimestampWrapper, XmlCryptographicValidation, LevelRule).
func NewTimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T any](i18nProvider *i18n.Provider,
	result *process.Result[T], contentTimestamp *diagnostic.TimestampWrapper,
	cryptographicValidation *jaxb.XmlCryptographicValidation,
	constraint policy.LevelRule) *TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T] {
	c := &TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T]{
		ChainItemBase:           process.NewChainItemBaseWithId(i18nProvider, result, constraint, contentTimestamp.Id()),
		contentTimestamp:        contentTimestamp,
		cryptographicValidation: cryptographicValidation,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeCNTTSTBBB
}

// Process performs the check. Port of process().
func (c *TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T]) Process() bool {
	productionTime := c.contentTimestamp.ProductionTime()
	notAfter := c.cryptographicValidation.NotAfter
	return productionTime != nil && notAfter != nil && !productionTime.After(notAfter.Time())
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCryptoConstraintsFailure
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVICTGTNACCET
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVICTGTNACCETANS
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *TimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck[T]) BuildAdditionalInfo() *string {
	tstGenerationTime := " ? "
	if productionTime := c.contentTimestamp.ProductionTime(); productionTime != nil {
		tstGenerationTime = process.GetFormattedDate(productionTime)
	}
	cryptoConstraintsExpiration := " ? "
	if c.cryptographicValidation.NotAfter != nil {
		t := c.cryptographicValidation.NotAfter.Time()
		cryptoConstraintsExpiration = process.GetFormattedDate(&t)
	}

	algorithmName := "?"
	if algorithm := c.cryptographicValidation.Algorithm; algorithm != nil {
		algorithmName = algorithm.Name
		if algorithm.KeyLength != nil {
			algorithmName += fmt.Sprintf(" with keyLength '%s'", *algorithm.KeyLength)
		}
	}

	message := c.I18nProvider.GetMessage(i18n.MessageTagTimestampAndCryptoConstraintsExpiration, c.contentTimestamp.Id(),
		tstGenerationTime, algorithmName, cryptoConstraintsExpiration)
	return &message
}
