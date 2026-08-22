// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/checks/AlgorithmObsolescenceValidationCheck.java (DSS 6.5.RC1).
//
// This is a deliberately minimal slice of Java's
// eu.europa.esig.dss.validation.process.bbb.aov package tree: only
// AlgorithmObsolescenceValidationCheck is ported in this file, as the sole
// consumer-side check that reads an already-built XmlAOV result, which is
// what AbstractAcceptanceValidation.cryptographic() (bbb/sav) needs. It has
// no dependency on the rest of the Java aov tree (AlgorithmObsolescenceValidation
// and its Certificate/Digest/Signature/Timestamp/EvidenceRecord/EAA
// subclasses, which *produce* an XmlAOV by orchestrating bbb/xcv and the
// aov/cc cryptographic-checker family, or the aov/cc package itself) - those
// live alongside AlgorithmObsolescenceValidationCheckWithId in this same
// package.
package aov

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AlgorithmObsolescenceValidationCheck verifies result of the
// bbb.aov.AlgorithmObsolescenceValidation process.
type AlgorithmObsolescenceValidationCheck[T any] struct {
	*process.ChainItemBase[T]

	// aovResult is the result of the AOV validation process.
	aovResult *jaxb.XmlAOV

	// validationDate is the check execution time.
	validationDate time.Time

	// position is the validating constraint position.
	position i18n.MessageTag

	// blockType is the type of the validation block.
	blockType jaxb.XmlBlockType
}

// NewAlgorithmObsolescenceValidationCheck is the convenience constructor,
// defaulting the block type to XmlBlockTypeAOV. Port of
// AlgorithmObsolescenceValidationCheck(I18nProvider, T, XmlAOV, Date, MessageTag, String).
func NewAlgorithmObsolescenceValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	aovResult *jaxb.XmlAOV, validationDate time.Time, position i18n.MessageTag,
	tokenId string) *AlgorithmObsolescenceValidationCheck[T] {
	return NewAlgorithmObsolescenceValidationCheckWithBlockType(i18nProvider, result, aovResult, validationDate,
		position, jaxb.XmlBlockTypeAOV, tokenId)
}

// NewAlgorithmObsolescenceValidationCheckWithBlockType is the full constructor.
// Port of AlgorithmObsolescenceValidationCheck(I18nProvider, T, XmlAOV, Date,
// MessageTag, XmlBlockType, String).
func NewAlgorithmObsolescenceValidationCheckWithBlockType[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	aovResult *jaxb.XmlAOV, validationDate time.Time, position i18n.MessageTag,
	blockType jaxb.XmlBlockType, tokenId string) *AlgorithmObsolescenceValidationCheck[T] {
	c := &AlgorithmObsolescenceValidationCheck[T]{
		ChainItemBase:  process.NewChainItemBaseWithId(i18nProvider, result, getLevelRule(aovResult), tokenId),
		aovResult:      aovResult,
		validationDate: validationDate,
		position:       position,
		blockType:      blockType,
	}
	c.InitChainItem(c)
	return c
}

// getLevelRule ports the private static getLevelRule(XmlAOV).
func getLevelRule(aovResult *jaxb.XmlAOV) policy.LevelRule {
	conclusion := aovResult.Conclusion
	if utils.IsCollectionNotEmpty(conclusion.Errors) {
		return process.GetLevelRule(enumerations.LevelFail)
	} else if utils.IsCollectionNotEmpty(conclusion.Warnings) {
		return process.GetLevelRule(enumerations.LevelWarn)
	} else if utils.IsCollectionNotEmpty(conclusion.Infos) {
		return process.GetLevelRule(enumerations.LevelInform)
	}
	return process.GetLevelRule(enumerations.LevelFail) // default
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *AlgorithmObsolescenceValidationCheck[T]) BlockType() jaxb.XmlBlockType {
	return c.blockType
}

// Process performs the check. Port of process().
func (c *AlgorithmObsolescenceValidationCheck[T]) Process() bool {
	return c.IsValid(&c.aovResult.XmlConstraintsConclusionContent)
}

// IsValidConclusion checks if the conclusion has a passed status, and neither
// warnings nor infos. Port of isValidConclusion(XmlConclusion).
func (c *AlgorithmObsolescenceValidationCheck[T]) IsValidConclusion(conclusion *jaxb.XmlConclusion) bool {
	return c.ChainItemBase.IsValidConclusion(conclusion) &&
		utils.IsCollectionEmpty(conclusion.Warnings) && utils.IsCollectionEmpty(conclusion.Infos)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *AlgorithmObsolescenceValidationCheck[T]) BuildAdditionalInfo() *string {
	dateTime := process.GetFormattedDate(&c.validationDate)
	if c.Process() {
		cryptographicValidation := process.GetPrimaryCryptographicValidation(c.aovResult)
		if cryptographicValidation != nil {
			algorithm := cryptographicValidation.Algorithm
			var message string
			if algorithm.KeyLength != nil && *algorithm.KeyLength != "" {
				message = c.I18nProvider.GetMessage(i18n.MessageTagCryptographicCheckSuccessKeySize,
					algorithm.Name, *algorithm.KeyLength, dateTime)
			} else {
				message = c.I18nProvider.GetMessage(i18n.MessageTagCryptographicCheckSuccess, algorithm.Name, dateTime)
			}
			return &message
		}
		return nil
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagCryptographicCheckFailure, c.errorMessage(), dateTime)
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AlgorithmObsolescenceValidationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.aovResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AlgorithmObsolescenceValidationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.aovResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.aovResult.Conclusion.SubIndication.SubIndication()
}

// BuildConstraintMessage builds a constraint message. Port of
// buildConstraintMessage().
func (c *AlgorithmObsolescenceValidationCheck[T]) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTagACCM, c.position)
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *AlgorithmObsolescenceValidationCheck[T]) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTagACCMANS, c.position)
}

// errorMessage returns the first error/warning/info message value, or the
// empty string if the check succeeded. Port of getErrorMessage().
func (c *AlgorithmObsolescenceValidationCheck[T]) errorMessage() string {
	conclusion := c.aovResult.Conclusion
	if utils.IsCollectionNotEmpty(conclusion.Errors) {
		return conclusion.Errors[0].Value
	}
	if utils.IsCollectionNotEmpty(conclusion.Warnings) {
		return conclusion.Warnings[0].Value
	}
	if utils.IsCollectionNotEmpty(conclusion.Infos) {
		return conclusion.Infos[0].Value
	}
	return ""
}

// PreviousErrors returns a list of previous errors occurred in the chain. Port
// of getPreviousErrors().
func (c *AlgorithmObsolescenceValidationCheck[T]) PreviousErrors() []*jaxb.XmlMessage {
	return c.aovResult.Conclusion.Errors
}
