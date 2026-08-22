// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/AllCertificatesInPathReferencedCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AllCertificatesInPathReferencedCheck checks if all certificates in the path
// have the corresponding signing certificate references.
type AllCertificatesInPathReferencedCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// token is the token to check.
	token diagnostic.TokenProxy
}

// NewAllCertificatesInPathReferencedCheck is the default constructor. Port of
// AllCertificatesInPathReferencedCheck(I18nProvider, XmlSAV, TokenProxy, LevelRule).
func NewAllCertificatesInPathReferencedCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	token diagnostic.TokenProxy, constraint policy.LevelRule) *AllCertificatesInPathReferencedCheck {
	c := &AllCertificatesInPathReferencedCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AllCertificatesInPathReferencedCheck) Process() bool {
	relatedSigningCertificates := c.token.FoundCertificates().RelatedCertificatesByRefOrigin(enumerations.CertificateRefOriginSigningCertificate)
	signingCertificateIds := make(map[string]struct{}, len(relatedSigningCertificates))
	for _, cert := range relatedSigningCertificates {
		signingCertificateIds[cert.Id()] = struct{}{}
	}

	for _, certificate := range c.token.CertificateChain() {
		if _, ok := signingCertificateIds[certificate.Id()]; !ok {
			// certificate in the certificate path is not covered by a signing certificate reference
			return false
		}
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AllCertificatesInPathReferencedCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ACPCCRSCA
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *AllCertificatesInPathReferencedCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ACPCCRSCA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AllCertificatesInPathReferencedCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *AllCertificatesInPathReferencedCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}
