// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/SigningCertificateReferencesValidityCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// SigningCertificateReferencesValidityCheck checks if a signing certificate
// reference is present and valid (all signingCertificate references refer the
// signature certificate chain).
type SigningCertificateReferencesValidityCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// token is the token to check.
	token diagnostic.TokenProxy
}

// NewSigningCertificateReferencesValidityCheck is the default constructor. Port
// of SigningCertificateReferencesValidityCheck(I18nProvider, XmlSAV, TokenProxy, LevelRule).
func NewSigningCertificateReferencesValidityCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	token diagnostic.TokenProxy, constraint policy.LevelRule) *SigningCertificateReferencesValidityCheck {
	c := &SigningCertificateReferencesValidityCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SigningCertificateReferencesValidityCheck) Process() bool {
	foundCertificates := c.token.FoundCertificates()

	// 1) Check orphan references presence
	orphanSigningCertificateRefs := foundCertificates.OrphanCertificateRefsByRefOrigin(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)
	if utils.IsCollectionNotEmpty(orphanSigningCertificateRefs) {
		// the provided reference does not match the provided certificate chain
		return false
	}

	// 2) Check found references against the certificate chain
	relatedSigningCertificates := foundCertificates.RelatedCertificatesByRefOrigin(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)

	certificateChainIds := make(map[string]struct{})
	for _, certificate := range c.token.CertificateChain() {
		certificateChainIds[certificate.Id()] = struct{}{}
	}

	for _, signingCertificate := range relatedSigningCertificates {
		if _, ok := certificateChainIds[signingCertificate.Id()]; !ok {
			// a certificate referenced by a SigningCertificate reference is not included into the certificate chain
			return false
		}
	}

	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SigningCertificateReferencesValidityCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_DSCACRCC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SigningCertificateReferencesValidityCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_DSCACRCC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SigningCertificateReferencesValidityCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SigningCertificateReferencesValidityCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
