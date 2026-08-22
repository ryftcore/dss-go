// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateAlgorithmObsolescenceValidationCheck.java (DSS 6.5.RC1).
//
// Java's class extends AlgorithmObsolescenceValidationCheckWithId, itself
// extending aov.checks.AlgorithmObsolescenceValidationCheck, overriding only
// getBlockType() and process(); buildAdditionalInfo() is inherited from
// WithId, which wraps the base class's buildAdditionalInfo() with the token
// id, virtual-dispatching into this class's process() along the way.
//
// The Go aov package's AlgorithmObsolescenceValidationCheck self-registers
// with process.ChainItemBase.InitChainItem in its own constructor (it is a
// leaf, not designed for further Go-embedding subclassing the way
// bbb.Abstract*CheckItem is): embedding it here would make its own Process /
// BuildAdditionalInfo win the dispatch, never this type's overrides. This
// file instead replicates the base class's and WithId's method bodies
// directly (mirroring the Java call graph, including the getPrimaryCryptographicValidation/
// getErrorMessage calls that read the *overall* aovResult rather than the
// per-certificate one - matching upstream's buildAdditionalInfo() exactly),
// duplicating aov.AlgorithmObsolescenceValidationCheck's private
// getLevelRule/getErrorMessage the same way every AbstractCertificateCheckItem-style
// Go type duplicates state its ChainItemBase already privately holds.
package xcv

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateAlgorithmObsolescenceValidationCheck checks validity of the
// bbb.aov.AlgorithmObsolescenceValidation process performed on a certificate
// token alone.
type CertificateAlgorithmObsolescenceValidationCheck[T any] struct {
	*process.ChainItemBase[T]

	// aovResult is the overall result of the AOV validation process.
	aovResult *jaxb.XmlAOV

	// validationDate is the check execution time.
	validationDate time.Time

	// position is the validating constraint position.
	position i18n.MessageTag

	// tokenID is the identifier of the certificate to be validated.
	tokenID string
}

// NewCertificateAlgorithmObsolescenceValidationCheck is the default
// constructor. Port of
// CertificateAlgorithmObsolescenceValidationCheck(I18nProvider, T, XmlAOV, Date, MessageTag, String).
func NewCertificateAlgorithmObsolescenceValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	aovResult *jaxb.XmlAOV, validationDate time.Time, position i18n.MessageTag,
	certificateId string) *CertificateAlgorithmObsolescenceValidationCheck[T] {
	c := &CertificateAlgorithmObsolescenceValidationCheck[T]{
		ChainItemBase:  process.NewChainItemBaseWithId(i18nProvider, result, certificateAOVLevelRule(aovResult), certificateId),
		aovResult:      aovResult,
		validationDate: validationDate,
		position:       position,
		tokenID:        certificateId,
	}
	c.InitChainItem(c)
	return c
}

// certificateAOVLevelRule ports the private static
// AlgorithmObsolescenceValidationCheck#getLevelRule(XmlAOV).
func certificateAOVLevelRule(aovResult *jaxb.XmlAOV) policy.LevelRule {
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
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeAOVXCV
}

// Process performs the check. Port of process().
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) Process() bool {
	cryptographicValidation := c.getCertificateCryptographicValidation()
	return cryptographicValidation != nil && isValidAOVCryptographicConclusion(cryptographicValidation.Conclusion)
}

// getCertificateCryptographicValidation ports the private
// getCertificateCryptographicValidation().
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) getCertificateCryptographicValidation() *jaxb.XmlCryptographicValidation {
	if c.aovResult != nil && c.aovResult.CertificateChainCryptographicValidation != nil &&
		utils.IsCollectionNotEmpty(c.aovResult.CertificateChainCryptographicValidation.CertificateCryptographicValidation) {
		for _, cryptographicValidation := range c.aovResult.CertificateChainCryptographicValidation.CertificateCryptographicValidation {
			if cryptographicValidation.TokenId != nil && c.tokenID == *cryptographicValidation.TokenId {
				return cryptographicValidation
			}
		}
	}
	return nil
}

// isValidAOVCryptographicConclusion ports
// AlgorithmObsolescenceValidationCheck#isValidConclusion(XmlConclusion): a
// passed conclusion carrying neither warnings nor infos.
func isValidAOVCryptographicConclusion(conclusion *jaxb.XmlConclusion) bool {
	return conclusion != nil &&
		(enumerations.IndicationPassed == conclusion.Indication.Indication() ||
			enumerations.IndicationTotalPassed == conclusion.Indication.Indication()) &&
		utils.IsCollectionEmpty(conclusion.Warnings) && utils.IsCollectionEmpty(conclusion.Infos)
}

// BuildConstraintMessage builds a constraint message. Port of
// AlgorithmObsolescenceValidationCheck#buildConstraintMessage().
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTagACCM, c.position)
}

// BuildErrorMessage builds an error message. Port of
// AlgorithmObsolescenceValidationCheck#buildErrorMessage().
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTagACCMANS, c.position)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// AlgorithmObsolescenceValidationCheck#getFailedIndicationForConclusion().
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.aovResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of AlgorithmObsolescenceValidationCheck#getFailedSubIndicationForConclusion().
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.aovResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.aovResult.Conclusion.SubIndication.SubIndication()
}

// PreviousErrors returns a list of previous errors occurred in the chain.
// Port of AlgorithmObsolescenceValidationCheck#getPreviousErrors().
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) PreviousErrors() []*jaxb.XmlMessage {
	return c.aovResult.Conclusion.Errors
}

// BuildAdditionalInfo builds an additional information. Port of
// AlgorithmObsolescenceValidationCheckWithId#buildAdditionalInfo(), wrapping
// the base class's buildAdditionalInfo() with the token id.
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) BuildAdditionalInfo() *string {
	base := c.baseAdditionalInfo()
	var baseStr string
	if base != nil {
		baseStr = *base
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagACCMDescWithIDResult, baseStr, c.tokenID)
	return &message
}

// baseAdditionalInfo ports AlgorithmObsolescenceValidationCheck#buildAdditionalInfo().
// NOTE: like upstream, this reads the *overall* aovResult (via
// process.GetPrimaryCryptographicValidation), not the per-certificate
// cryptographic validation used by Process() - only the branch taken (via
// c.Process(), which dispatches to this type's override) differs from the
// base's own behavior.
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) baseAdditionalInfo() *string {
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
	message := c.I18nProvider.GetMessage(i18n.MessageTagCryptographicCheckFailure, c.overallErrorMessage(), dateTime)
	return &message
}

// overallErrorMessage ports AlgorithmObsolescenceValidationCheck#getErrorMessage().
func (c *CertificateAlgorithmObsolescenceValidationCheck[T]) overallErrorMessage() string {
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
