// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/SignatureAcceptanceValidationNoCryptoResultCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignatureAcceptanceValidationNoCryptoResultCheck verifies if the format
// Signature Acceptance Validation process as per clause 5.2.8 succeeded, but
// skips the cryptographic check failures, if any.
type SignatureAcceptanceValidationNoCryptoResultCheck[T any] struct {
	*SignatureAcceptanceValidationResultCheck[T]
}

// NewSignatureAcceptanceValidationNoCryptoResultCheck is the default
// constructor. Port of
// SignatureAcceptanceValidationNoCryptoResultCheck(I18nProvider, T, XmlSAV, TokenProxy, LevelRule).
//
// The constructor re-registers the overrides with the outer type, so that the
// base's self-calls reach this class' Process rather than the one inherited
// from SignatureAcceptanceValidationResultCheck.
func NewSignatureAcceptanceValidationNoCryptoResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	xmlSAV *jaxb.XmlSAV, token diagnostic.TokenProxy, constraint policy.LevelRule) *SignatureAcceptanceValidationNoCryptoResultCheck[T] {
	c := &SignatureAcceptanceValidationNoCryptoResultCheck[T]{
		SignatureAcceptanceValidationResultCheck: NewSignatureAcceptanceValidationResultCheck(i18nProvider, result, xmlSAV, token, constraint),
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignatureAcceptanceValidationNoCryptoResultCheck[T]) Process() bool {
	return c.XmlSAV != nil && (c.IsValid(&c.XmlSAV.XmlConstraintsConclusionContent) || c.isCryptoFailure(c.XmlSAV))
}

// isCryptoFailure ports the private isCryptoFailure(XmlSAV).
func (c *SignatureAcceptanceValidationNoCryptoResultCheck[T]) isCryptoFailure(xmlSAV *jaxb.XmlSAV) bool {
	var subIndication enumerations.SubIndication
	if xmlSAV.Conclusion.SubIndication != nil {
		subIndication = xmlSAV.Conclusion.SubIndication.SubIndication()
	}
	return enumerations.Indication_INDETERMINATE == xmlSAV.Conclusion.Indication.Indication() &&
		enumerations.SubIndication_CRYPTO_CONSTRAINTS_FAILURE_NO_POE == subIndication
}
