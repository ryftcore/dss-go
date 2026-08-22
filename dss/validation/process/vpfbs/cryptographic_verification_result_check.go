// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/CryptographicVerificationResultCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CryptographicVerificationResultCheck verifies if the format Cryptographic
// Verification process as per clause 5.2.7 succeeded.
type CryptographicVerificationResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// xmlCV is the Cryptographic Verification result.
	xmlCV *jaxb.XmlCV
}

// NewCryptographicVerificationResultCheck is the default constructor. Port of
// CryptographicVerificationResultCheck(I18nProvider, T, XmlCV, TokenProxy, LevelRule).
func NewCryptographicVerificationResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	xmlCV *jaxb.XmlCV, token diagnostic.TokenProxy, constraint policy.LevelRule) *CryptographicVerificationResultCheck[T] {
	c := &CryptographicVerificationResultCheck[T]{
		// Cryptographic Verification building block suffix ("-CV"), a per-class
		// private constant in Java; inlined here since Go package-level constants
		// share one namespace.
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, token.Id()+"-CV"),
		xmlCV:         xmlCV,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CryptographicVerificationResultCheck[T]) Process() bool {
	return c.xmlCV != nil && c.IsValid(&c.xmlCV.XmlConstraintsConclusionContent)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CryptographicVerificationResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.xmlCV.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CryptographicVerificationResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.xmlCV.Conclusion.SubIndication == nil {
		return ""
	}
	return c.xmlCV.Conclusion.SubIndication.SubIndication()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CryptographicVerificationResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_ICVRC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *CryptographicVerificationResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_ICVRC_ANS
}
