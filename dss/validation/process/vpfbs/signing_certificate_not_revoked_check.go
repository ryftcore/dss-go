// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/SigningCertificateNotRevokedCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SigningCertificateNotRevokedCheck verifies if the X.509 Certificate
// Validation as per clause 5.2.6 did not return INDETERMINATE/REVOKED_NO_POE
// indication.
type SigningCertificateNotRevokedCheck[T any] struct {
	*process.ChainItemBase[T]

	// xmlXCV is the token's X509CertificateValidation result.
	xmlXCV *jaxb.XmlXCV
}

// NewSigningCertificateNotRevokedCheck is the default constructor. Port of
// SigningCertificateNotRevokedCheck(I18nProvider, T, XmlXCV, TokenProxy, LevelRule).
func NewSigningCertificateNotRevokedCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	xmlXCV *jaxb.XmlXCV, token diagnostic.TokenProxy, constraint policy.LevelRule) *SigningCertificateNotRevokedCheck[T] {
	c := &SigningCertificateNotRevokedCheck[T]{
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
func (c *SigningCertificateNotRevokedCheck[T]) Process() bool {
	if c.xmlXCV == nil {
		return false
	}
	conclusion := c.xmlXCV.Conclusion
	var subIndication enumerations.SubIndication
	if conclusion.SubIndication != nil {
		subIndication = conclusion.SubIndication.SubIndication()
	}
	return !(enumerations.IndicationIndeterminate == conclusion.Indication.Indication() &&
		enumerations.SubIndicationRevokedNoPOE == subIndication)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SigningCertificateNotRevokedCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *SigningCertificateNotRevokedCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationRevokedNoPOE
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SigningCertificateNotRevokedCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_ISCRAVTC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SigningCertificateNotRevokedCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_ISCRAVTC_ANS
}
