// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/X509CertificateValidationResultCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// X509CertificateValidationResultCheck verifies if the X.509 Certificate
// Validation as per clause 5.2.6 succeeded.
type X509CertificateValidationResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// xmlXCV is the token's X509CertificateValidation result.
	xmlXCV *jaxb.XmlXCV
}

// NewX509CertificateValidationResultCheck is the default constructor. Port of
// X509CertificateValidationResultCheck(I18nProvider, T, XmlXCV, TokenProxy, LevelRule).
func NewX509CertificateValidationResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	xmlXCV *jaxb.XmlXCV, token diagnostic.TokenProxy, constraint policy.LevelRule) *X509CertificateValidationResultCheck[T] {
	c := &X509CertificateValidationResultCheck[T]{
		// X509 Certificate Validation building block suffix ("-XCV"), a
		// per-class private constant in Java; inlined here since Go
		// package-level constants share one namespace.
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, token.Id()+"-XCV"),
		xmlXCV:        xmlXCV,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *X509CertificateValidationResultCheck[T]) Process() bool {
	return c.xmlXCV != nil && c.IsValid(&c.xmlXCV.XmlConstraintsConclusionContent)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *X509CertificateValidationResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.xmlXCV.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *X509CertificateValidationResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.xmlXCV.Conclusion.SubIndication == nil {
		return ""
	}
	return c.xmlXCV.Conclusion.SubIndication.SubIndication()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *X509CertificateValidationResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_IXCVRC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *X509CertificateValidationResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_IXCVRC_ANS
}
