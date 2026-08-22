// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/QWACDomainNameCheck.java (DSS 6.5.RC1).
package qualification

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QWACDomainNameCheck verifies whether the website domain name in question
// appears in the QWAC's subject alternative name(s).
type QWACDomainNameCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// certificate is the certificate to be validated.
	certificate *diagnostic.CertificateWrapper

	// websiteUrl is the URL of the website.
	websiteUrl string
}

// NewQWACDomainNameCheck is the default constructor. Port of
// QWACDomainNameCheck(I18nProvider, XmlValidationQWACProcess, CertificateWrapper, String, LevelRule).
func NewQWACDomainNameCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	certificate *diagnostic.CertificateWrapper, websiteUrl string, constraint policy.LevelRule) *QWACDomainNameCheck {
	c := &QWACDomainNameCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
		websiteUrl:    websiteUrl,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QWACDomainNameCheck) Process() bool {
	host := spi.DSSUtilsHost(c.websiteUrl)

	for _, generalName := range c.certificate.SubjectAlternativeNames() {
		if generalName.Type == nil {
			continue
		}
		switch generalName.Type.GeneralNameType() {
		case enumerations.GeneralNameTypeDNSName:
			if c.matchesDNSName(host, generalName.Value) {
				return true
			}
		case enumerations.GeneralNameTypeIPAddress:
			if c.matchesIPAddress(host, generalName.Value) {
				return true
			}
		default:
			// not supported
		}
	}
	return false
}

// matchesDNSName verifies if the given domain name matches the SAN DNS
// name, according to CAB Forum BR section 3.2.2.6 and RFC 6125 section
// 6.4.3 (wildcard matching). Port of the private matchesDNSName(String, String).
func (c *QWACDomainNameCheck) matchesDNSName(hostname, subAltName string) bool {
	subAltName = strings.ToLower(subAltName)
	hostname = strings.ToLower(hostname)

	if subAltName == hostname {
		return true
	}

	/*
	 * If a client matches the reference identifier against a presented
	 * identifier whose DNS domain name portion contains the wildcard
	 * character '*', the following rules apply:
	 *
	 * 1. The client SHOULD NOT attempt to match a presented identifier in
	 * which the wildcard character comprises a label other than the
	 * left-most label (e.g., do not match bar.*.example.net).
	 *
	 * 2. If the wildcard character is the only character of the left-most
	 * label in the presented identifier, the client SHOULD NOT compare
	 * against anything but the left-most label of the reference identifier
	 * (e.g., *.example.com would match foo.example.com but not
	 * bar.foo.example.com or example.com).
	 */

	// NOTE: Only Full Qualified Domain Names are considered (FQDN)

	// Wildcard match
	if strings.HasPrefix(subAltName, "*.") && strings.Contains(hostname, ".") {
		suffix := subAltName[2:]
		return strings.HasSuffix(hostname, suffix) && c.countParts(hostname) == c.countParts(suffix)+1
	}

	return false
}

// countParts ports the private countParts(String).
func (c *QWACDomainNameCheck) countParts(domain string) int {
	return len(strings.Split(domain, "."))
}

// matchesIPAddress verifies if the given domain name matches one the IP
// Address, according to CAB Forum BR section 7.1.2.7.12. Port of the private
// matchesIPAddress(String, String): must be a complete match, no wildcards
// supported.
func (c *QWACDomainNameCheck) matchesIPAddress(hostname, subAltName string) bool {
	return hostname == subAltName
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QWACDomainNameCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QWAC_DOMAIN_NAME
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QWACDomainNameCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QWAC_DOMAIN_NAME_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QWACDomainNameCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QWACDomainNameCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
