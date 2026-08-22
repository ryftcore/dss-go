// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/checks/AlgorithmObsolescenceValidationCheckWithId.java (DSS 6.5.RC1).
//
// This class only has one Java caller (CertificateAlgorithmObsolescenceValidationCheck,
// in eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks), which the xcv
// chunk replicates directly rather than embedding this type (see the header
// of bbb/xcv/certificate_algorithm_obsolescence_validation_check.go: embedding
// a leaf that self-registers with InitChainItem would make the base's own
// Process/BuildAdditionalInfo win the dispatch). It is ported here regardless,
// as a manifest file, following the same embed-and-re-register technique used
// by cv.SignatureIntactWithIdCheck for an analogous WithId variant.
package aov

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AlgorithmObsolescenceValidationCheckWithId verifies result of the
// bbb.aov.AlgorithmObsolescenceValidation process. This class also allows
// providing an identifier of a token to be validated.
type AlgorithmObsolescenceValidationCheckWithId[T any] struct {
	*AlgorithmObsolescenceValidationCheck[T]

	// tokenId is the identifier of a token to be validated.
	tokenId string
}

// NewAlgorithmObsolescenceValidationCheckWithId is the default constructor.
//
// The constructor re-registers the overrides with the outer type, so that the
// base's self-calls reach this class' BuildAdditionalInfo rather than the one
// inherited from AlgorithmObsolescenceValidationCheck.
func NewAlgorithmObsolescenceValidationCheckWithId[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	aovResult *jaxb.XmlAOV, validationDate time.Time, position i18n.MessageTag,
	tokenId string) *AlgorithmObsolescenceValidationCheckWithId[T] {
	c := &AlgorithmObsolescenceValidationCheckWithId[T]{
		AlgorithmObsolescenceValidationCheck: NewAlgorithmObsolescenceValidationCheck(i18nProvider, result, aovResult, validationDate, position, tokenId),
		tokenId:                              tokenId,
	}
	c.InitChainItem(c)
	return c
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo(), wrapping the base class's
// buildAdditionalInfo() with the token id.
func (c *AlgorithmObsolescenceValidationCheckWithId[T]) BuildAdditionalInfo() *string {
	base := c.AlgorithmObsolescenceValidationCheck.BuildAdditionalInfo()
	// Java hands super.buildAdditionalInfo()'s result straight to
	// MessageFormat, which renders a null argument as the literal text "null"
	// (the reachable case: a passing check whose XmlAOV carries no primary
	// cryptographic validation). Substituting "" here would drop that word.
	baseStr := "null"
	if base != nil {
		baseStr = *base
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagACCMDescWithIDResult, baseStr, c.tokenId)
	return &message
}
