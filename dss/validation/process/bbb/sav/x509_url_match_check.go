// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/X509UrlMatchCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// X509UrlMatchCheck verifies if the application was able to derive
// signing-certificate using the value of the 'x5u' (X.509 URL) header parameter
// of the protected header of the signature.
type X509UrlMatchCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// signature is the signature to verify.
	signature *diagnostic.SignatureWrapper
}

// NewX509UrlMatchCheck is the default constructor. Port of
// X509UrlMatchCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewX509UrlMatchCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *X509UrlMatchCheck {
	c := &X509UrlMatchCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *X509UrlMatchCheck) Process() bool {
	signingCertificate := c.signature.SigningCertificate()
	if signingCertificate == nil {
		return false
	}
	x509UrlCertificates := c.signature.FoundCertificates().RelatedCertificatesByRefOrigin(enumerations.CertificateRefOrigin_X509_URL)
	for _, r := range x509UrlCertificates {
		if signingCertificate.Id() == r.Id() {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *X509UrlMatchCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISAX509UA
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *X509UrlMatchCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISAX509UA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *X509UrlMatchCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *X509UrlMatchCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
