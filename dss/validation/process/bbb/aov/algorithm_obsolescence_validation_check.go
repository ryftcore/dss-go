// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/checks/AlgorithmObsolescenceValidationCheck.java (DSS 6.5.RC1).
//
// Integration note (phase 8c pass): this is a deliberately minimal slice of
// Java's eu.europa.esig.dss.validation.process.bbb.aov package tree. Only
// AlgorithmObsolescenceValidationCheck is ported here - the sole consumer-side
// check that reads an already-built XmlAOV result, which is exactly what
// AbstractAcceptanceValidation.cryptographic() (bbb/sav) needs and is the only
// call site anywhere in the phase 8c chunks. It has no dependency on the rest
// of the Java aov tree: not on AlgorithmObsolescenceValidation.java or its
// Certificate/Digest/Signature/Timestamp/EvidenceRecord/EAA subclasses (the
// classes that *produce* an XmlAOV by orchestrating bbb/xcv + the aov/cc
// cryptographic-checker family), and not on the aov/cc package itself. Those
// were genuinely out of phase 8c scope and were added by phase 8d to this same
// package - along with AlgorithmObsolescenceValidationCheckWithId.java, the
// sibling class in the same Java package, which likewise had no caller in
// phase 8c. (The dispatcher this paragraph called "FRAME's
// basic_building_blocks.go, gated behind //go:build phase8d" is now the
// un-tagged dss/validation/process/blocks package; no phase8d build tag
// remains anywhere in the tree.)
//
// Every dependency this file does use - jaxb.XmlAOV/XmlConclusion/
// XmlCryptographicValidation/XmlCryptographicAlgorithm/XmlMessage/XmlBlockType,
// process.ChainItemBase/Result, process.GetFormattedDate/
// GetPrimaryCryptographicValidation/GetLevelRule, and the ACCM/ACCM_ANS/
// CRYPTOGRAPHIC_CHECK_* message tags - already existed untagged in the tree
// before this pass; nothing new was invented to support it.
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
// defaulting the block type to XmlBlockType_AOV. Port of
// AlgorithmObsolescenceValidationCheck(I18nProvider, T, XmlAOV, Date, MessageTag, String).
func NewAlgorithmObsolescenceValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	aovResult *jaxb.XmlAOV, validationDate time.Time, position i18n.MessageTag,
	tokenId string) *AlgorithmObsolescenceValidationCheck[T] {
	return NewAlgorithmObsolescenceValidationCheckWithBlockType(i18nProvider, result, aovResult, validationDate,
		position, jaxb.XmlBlockType_AOV, tokenId)
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
		return process.GetLevelRule(enumerations.Level_FAIL)
	} else if utils.IsCollectionNotEmpty(conclusion.Warnings) {
		return process.GetLevelRule(enumerations.Level_WARN)
	} else if utils.IsCollectionNotEmpty(conclusion.Infos) {
		return process.GetLevelRule(enumerations.Level_INFORM)
	}
	return process.GetLevelRule(enumerations.Level_FAIL) // default
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
				message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_KEY_SIZE,
					algorithm.Name, *algorithm.KeyLength, dateTime)
			} else {
				message = c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS, algorithm.Name, dateTime)
			}
			return &message
		}
		return nil
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE, c.getErrorMessage(), dateTime)
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
	return c.BuildXmlMessage(i18n.MessageTag_ACCM, c.position)
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *AlgorithmObsolescenceValidationCheck[T]) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTag_ACCM_ANS, c.position)
}

// getErrorMessage returns the first error/warning/info message value, or the
// empty string if the check succeeded. Port of getErrorMessage().
func (c *AlgorithmObsolescenceValidationCheck[T]) getErrorMessage() string {
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
