// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/ValidationTimeAtCertificateValidityRangeCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ValidationTimeAtCertificateValidityRangeCheck verifies if the result of
// X509CertificateValidation is not indication INDETERMINATE with the
// sub-indication OUT_OF_BOUNDS_NO_POE or OUT_OF_BOUNDS_NOT_REVOKED.
type ValidationTimeAtCertificateValidityRangeCheck[T any] struct {
	*process.ChainItemBase[T]

	// xmlXCV is the token's X509CertificateValidation result.
	xmlXCV *jaxb.XmlXCV
}

// NewValidationTimeAtCertificateValidityRangeCheck is the default constructor.
// Port of
// ValidationTimeAtCertificateValidityRangeCheck(Provider, T, XmlXCV, TokenProxy, LevelRule).
func NewValidationTimeAtCertificateValidityRangeCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	xmlXCV *jaxb.XmlXCV, token diagnostic.TokenProxy, constraint policy.LevelRule) *ValidationTimeAtCertificateValidityRangeCheck[T] {
	c := &ValidationTimeAtCertificateValidityRangeCheck[T]{
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
func (c *ValidationTimeAtCertificateValidityRangeCheck[T]) Process() bool {
	if c.xmlXCV == nil {
		return false
	}
	conclusion := c.xmlXCV.Conclusion
	var subIndication enumerations.SubIndication
	if conclusion.SubIndication != nil {
		subIndication = conclusion.SubIndication.SubIndication()
	}
	return !(enumerations.IndicationIndeterminate == conclusion.Indication.Indication() &&
		(enumerations.SubIndicationOutOfBoundsNoPOE == subIndication ||
			enumerations.SubIndicationOutOfBoundsNotRevoked == subIndication))
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ValidationTimeAtCertificateValidityRangeCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.xmlXCV.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *ValidationTimeAtCertificateValidityRangeCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.xmlXCV.Conclusion.SubIndication == nil {
		return ""
	}
	return c.xmlXCV.Conclusion.SubIndication.SubIndication()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ValidationTimeAtCertificateValidityRangeCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVIVTAVRSC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ValidationTimeAtCertificateValidityRangeCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVIVTAVRSCANS
}
