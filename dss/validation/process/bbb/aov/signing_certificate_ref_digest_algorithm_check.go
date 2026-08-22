// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/checks/SigningCertificateRefDigestAlgorithmCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: see abstract_cryptographic_check.go.
//
// Java's isValid(XmlConstraintsConclusion) is itself a protected overridable
// method with a base body of isValidConclusion(getConclusion()); this class
// overrides it directly, adding an allConstraintsValid() check. The Go port's
// ChainItemOverrides interface only exposes the finer-grained IsValidConclusion
// hook (see process/chain_item.go), so IsValid itself cannot be
// virtually overridden the way Java allows. Since every internal call to
// isValid(...) in this Java class is made on `this` (never dispatched from a
// base-class self-call this port would otherwise miss), the extra
// allConstraintsValid check is folded directly into the two call sites
// (Process and BuildAdditionalInfo) via the unexported validResult helper
// below, rather than attempting to override process.ChainItemBase.IsValid.
package aov

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SigningCertificateRefDigestAlgorithmCheck verifies whether a used
// DigestAlgorithm for a signing-certificate-reference signing-attribute is
// reliable and acceptable at validation time.
type SigningCertificateRefDigestAlgorithmCheck[T any] struct {
	*process.ChainItemBase[T]

	// certificateRefs is the certificate references being validated.
	certificateRefs []*diagnostic.CertificateRefWrapper

	// certificateId is the Id of the certificate being validated.
	certificateId string

	// validationDate is the validation time.
	validationDate time.Time

	// context is the validation context.
	context enumerations.Context

	// subContext is the validation SubContext.
	subContext enumerations.SubContext

	// validationPolicy is the validation policy.
	validationPolicy policy.ValidationPolicy

	// constraint defines the check level.
	constraint policy.LevelRule

	// cryptographicValidationResult is the final validation result.
	cryptographicValidationResult *jaxb.XmlCC
}

// NewSigningCertificateRefDigestAlgorithmCheck is the default constructor.
func NewSigningCertificateRefDigestAlgorithmCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	validationDate time.Time, certificateRefs []*diagnostic.CertificateRefWrapper, certificateId string,
	context enumerations.Context, subContext enumerations.SubContext, validationPolicy policy.ValidationPolicy,
	constraint policy.LevelRule) *SigningCertificateRefDigestAlgorithmCheck[T] {
	c := &SigningCertificateRefDigestAlgorithmCheck[T]{
		certificateRefs:  certificateRefs,
		certificateId:    certificateId,
		validationDate:   validationDate,
		context:          context,
		subContext:       subContext,
		validationPolicy: validationPolicy,
		constraint:       constraint,
	}
	c.ChainItemBase = process.NewChainItemBaseWithId(i18nProvider, result, constraint, certificateId)
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of the overridden
// getBlockType().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeAOVXCV
}

// Process performs the check. Port of the overridden process().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) Process() bool {
	return c.validResult(c.validateCertReferences())
}

// validResult ports the this-dispatched Java isValid(XmlConstraintsConclusion)
// override: super.isValid(constraintConclusion) && allConstraintsValid(constraintConclusion).
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) validResult(ccResult *jaxb.XmlCC) bool {
	if ccResult == nil {
		return false
	}
	return c.IsValid(&ccResult.XmlConstraintsConclusionContent) && allCCConstraintsValid(ccResult)
}

// validateCertReferences performs validation of the signing certificate
// references' digest algorithms. Port of the protected validateCertReferences().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) validateCertReferences() *jaxb.XmlCC {
	if c.cryptographicValidationResult != nil {
		return c.cryptographicValidationResult
	}

	for _, certificateRefWrapper := range c.certificateRefs {
		digestAlgorithm := certificateRefWrapper.DigestMethod()
		if digestAlgorithm == "" {
			continue
		}

		dacResult := c.signingCertificateDigestCryptographicCheckResult(certificateRefWrapper)

		// overwrite only if previous checks are secure
		if c.cryptographicValidationResult == nil || !c.validResult(c.cryptographicValidationResult) {
			c.cryptographicValidationResult = dacResult
		}

		if c.validResult(c.cryptographicValidationResult) {
			break
		}
	}
	return c.cryptographicValidationResult
}

// signingCertificateDigestCryptographicCheckResult ports the private get
// signingCertificateDigestCryptographicCheckResult(CertificateRefWrapper).
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) signingCertificateDigestCryptographicCheckResult(
	certificateRef *diagnostic.CertificateRefWrapper) *jaxb.XmlCC {
	certificateConstraint := c.validationPolicy.CertificateCryptographicConstraint(c.context, c.subContext)
	dac := NewDigestAlgorithmCryptographicChecker(c.I18nProvider, certificateRef.DigestMethod(),
		c.validationDate, i18n.MessageTagACCMPosSigCertRef, certificateConstraint)
	return dac.Execute()
}

// BuildConstraintMessage builds a constraint message. Port of the overridden
// buildConstraintMessage().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTagACCM, i18n.MessageTagACCMPosSigCertRef)
}

// BuildErrorMessage builds an error message. Port of the overridden
// buildErrorMessage().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) BuildErrorMessage() *jaxb.XmlMessage {
	return c.extractXmlMessage()
}

// ErrorMessage gets error message. Port of the protected getErrorMessage().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) ErrorMessage() string {
	errorMessage := c.BuildErrorMessage()
	if errorMessage != nil {
		return errorMessage.Value
	}
	return utils.EmptyString
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port
// of the overridden getFailedIndicationForConclusion().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	if c.cryptographicValidationResult != nil {
		return c.cryptographicValidationResult.Conclusion.Indication.Indication()
	}
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the overridden getFailedSubIndicationForConclusion().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.cryptographicValidationResult != nil {
		if c.cryptographicValidationResult.Conclusion.SubIndication == nil {
			return ""
		}
		return c.cryptographicValidationResult.Conclusion.SubIndication.SubIndication()
	}
	return enumerations.SubIndicationCryptoConstraintsFailure
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) BuildAdditionalInfo() *string {
	dateTime := process.GetFormattedDate(&c.validationDate)
	if c.validResult(c.cryptographicValidationResult) {
		message := c.I18nProvider.GetMessage(i18n.MessageTagCryptographicCheckSuccessDMWithID,
			c.cryptographicValidationResult.CryptographicValidation.Algorithm.Name, dateTime,
			i18n.MessageTagACCMPosSigCertRef, c.certificateId)
		return &message
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagCryptographicCheckFailureWithID,
		c.ErrorMessage(), dateTime, i18n.MessageTagACCMPosSigCertRef, c.certificateId)
	return &message
}

// PreviousErrors returns a list of previous errors occurred in the chain.
// Port of the overridden getPreviousErrors().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) PreviousErrors() []*jaxb.XmlMessage {
	if c.cryptographicValidationResult != nil {
		return c.cryptographicValidationResult.Conclusion.Errors
	}
	return nil
}

// Level returns an execution Level of the chain item. Port of the overridden
// getLevel().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) Level() enumerations.Level {
	if c.constraint != nil {
		currentConstraintLevel := c.constraint.Level()
		subProcessLevel := c.subProcessLevel()
		return getLowestCryptoLevel(currentConstraintLevel, subProcessLevel)
	}
	return ""
}

// subProcessLevel ports the private getSubProcessLevel().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) subProcessLevel() enumerations.Level {
	conclusion := c.CryptographicValidationResult().Conclusion
	if conclusion != nil {
		if utils.IsCollectionNotEmpty(conclusion.Errors) {
			return enumerations.LevelFail
		} else if utils.IsCollectionNotEmpty(conclusion.Warnings) {
			return enumerations.LevelWarn
		} else if utils.IsCollectionNotEmpty(conclusion.Infos) {
			return enumerations.LevelInform
		}
	}
	return ""
}

// getLowestCryptoLevel ports the private static getLowestLevel(Level, Level).
func getLowestCryptoLevel(currentLevel, subProcessLevel enumerations.Level) enumerations.Level {
	if currentLevel == "" {
		return subProcessLevel
	} else if subProcessLevel == "" {
		return currentLevel
	} else if enumerations.LevelInform == currentLevel || enumerations.LevelInform == subProcessLevel {
		return enumerations.LevelInform
	} else if enumerations.LevelWarn == currentLevel || enumerations.LevelWarn == subProcessLevel {
		return enumerations.LevelWarn
	} else if enumerations.LevelFail == currentLevel || enumerations.LevelFail == subProcessLevel {
		return enumerations.LevelFail
	}
	return currentLevel
}

// extractXmlMessage ports the private extractXmlMessage().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) extractXmlMessage() *jaxb.XmlMessage {
	conclusion := c.cryptographicValidationResult.Conclusion
	if conclusion != nil {
		// Collects messages from higher levels only
		var messages []*jaxb.XmlMessage
		switch c.Level() {
		case enumerations.LevelInform:
			messages = append(messages, conclusion.Infos...)
			messages = append(messages, conclusion.Warnings...)
			messages = append(messages, conclusion.Errors...)
		case enumerations.LevelWarn:
			messages = append(messages, conclusion.Warnings...)
			messages = append(messages, conclusion.Errors...)
		case enumerations.LevelFail:
			messages = append(messages, conclusion.Errors...)
		}
		if utils.IsCollectionNotEmpty(messages) {
			return messages[0] // take the first one
		}
	}
	return nil
}

// CryptographicValidationResult gets the final cryptographic validation
// result. Port of the public getCryptographicValidationResult().
func (c *SigningCertificateRefDigestAlgorithmCheck[T]) CryptographicValidationResult() *jaxb.XmlCC {
	return c.validateCertReferences()
}
