// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/TrustedCertificateMatchTrustServiceCheck.java (DSS 6.5.RC1).
package qualification

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TrustedCertificateMatchTrustServiceCheck checks if the
// ServiceDigitalIdentifier of the TrustService matches the TrustService name.
type TrustedCertificateMatchTrustServiceCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// trustService is the Trusted Service to check.
	trustService *diagnostic.TrustServiceWrapper

	// errorMessage is the internal cached error status message.
	errorMessage i18n.MessageTag
}

// NewTrustedCertificateMatchTrustServiceCheck is the default constructor.
// Port of
// TrustedCertificateMatchTrustServiceCheck(I18nProvider, XmlValidationCertificateQualification, TrustServiceWrapper, LevelRule).
func NewTrustedCertificateMatchTrustServiceCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateQualification], trustService *diagnostic.TrustServiceWrapper,
	constraint policy.LevelRule) *TrustedCertificateMatchTrustServiceCheck {
	c := &TrustedCertificateMatchTrustServiceCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		trustService:  trustService,
		errorMessage:  i18n.MessageTag_EMPTY,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustedCertificateMatchTrustServiceCheck) Process() bool {
	trustedCert := c.trustService.ServiceDigitalIdentifier
	if trustedCert == nil {
		c.errorMessage = i18n.MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS0
		return false
	}

	organizationName := trustedCert.OrganizationName()
	if utils.IsStringBlank(organizationName) {
		c.errorMessage = i18n.MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS1
		return false
	}

	if !c.isMatch(trustedCert) {
		c.errorMessage = i18n.MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS2
		return false
	}

	return true
}

// isMatch ports the private isMatch(CertificateWrapper).
func (c *TrustedCertificateMatchTrustServiceCheck) isMatch(trustedCert *diagnostic.CertificateWrapper) bool {
	candidates := []string{trustedCert.OrganizationName(), trustedCert.CommonName(), trustedCert.OrganizationalUnit(),
		trustedCert.CertificateDN()}

	var possibleMatchers []string
	possibleMatchers = append(possibleMatchers, c.trustService.TspNames...)
	possibleMatchers = append(possibleMatchers, c.trustService.TspTradeNames...)
	possibleMatchers = append(possibleMatchers, c.trustService.ServiceNames...)

	for _, candidate := range candidates {
		if utils.IsStringBlank(candidate) {
			continue
		}
		for _, matcher := range possibleMatchers {
			if strings.EqualFold(candidate, matcher) {
				return true
			}
		}
	}

	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustedCertificateMatchTrustServiceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustedCertificateMatchTrustServiceCheck) ErrorMessageTag() i18n.MessageTag {
	return c.errorMessage
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustedCertificateMatchTrustServiceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustedCertificateMatchTrustServiceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
