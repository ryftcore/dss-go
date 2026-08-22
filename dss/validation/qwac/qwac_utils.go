// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/qwac/QWACUtils.java (DSS 6.5.RC1).
//
// Completed in phase 8f (S8F_BRIEF.md): GetTLSCertificateBindingUrl/isTLSCertificateBindingRel
// (parsing an HTTP "Link" response header) were previously left unported for lack of a
// LinkHeaderParser port to check them against (see link_header_parser.go, now ported alongside
// this completion).
package qwac

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
)

// headerLink represents a "Link" response header.
const headerLink = "Link"

// relationType represents a "rel" (relation type) attribute of the "Link" response header.
const relationType = "rel"

// tlsCertificateBinding represents the "tls-certificate-binding" value of the "rel" attribute
// identifying a Link to a TLS/SSL binding signature file.
const tlsCertificateBinding = "tls-certificate-binding"

// GetTLSCertificateBindingUrl loops over headers and returns the tls-certificate-binding URL
// value from a "Link" header if found. If no matching value is found, the function returns an
// empty string. Port of getTLSCertificateBindingUrl(Map<String, List<String>>).
//
// Java logs and swallows any parse error per Link header value candidate (LOG.debug); dropped
// per PORTING.md's "slf4j dropped unless load-bearing" - the loop simply moves on to the next
// candidate value on a parse error, exactly as Java's catch block did.
func GetTLSCertificateBindingUrl(headers map[string][]string) string {
	if headers == nil {
		return ""
	}

	headerValues, ok := headers[headerLink]
	if !ok {
		return ""
	}

	parser := NewLinkHeaderParser()
	for _, linkHeaderValue := range headerValues {
		linkHeader, err := parser.Parse(linkHeaderValue)
		if err != nil {
			continue
		}
		if utils.IsCollectionNotEmpty(linkHeader) {
			for _, singleHeaderValue := range linkHeader {
				if isTLSCertificateBindingRel(singleHeaderValue) {
					return singleHeaderValue.URL()
				}
			}
		}
	}
	return ""
}

// isTLSCertificateBindingRel checks whether the "Link" header attributes contain a rel value of
// tls-certificate-binding. Port of isTLSCertificateBindingRel(LinkHeaderParser.LinkHeader).
func isTLSCertificateBindingRel(linkHeader *LinkHeader) bool {
	rel, ok := linkHeader.Attributes()[relationType]
	return ok && rel != nil && *rel == tlsCertificateBinding
}

// GetIdentifiedTLSCertificates gets TLS Binding Certificates identified from
// the binding signature. Port of getIdentifiedTLSCertificates(SignatureWrapper, List<CertificateWrapper>).
func GetIdentifiedTLSCertificates(signature *diagnostic.SignatureWrapper, certificates []*diagnostic.CertificateWrapper) []*diagnostic.CertificateWrapper {
	var result []*diagnostic.CertificateWrapper
	for _, digestMatcher := range signature.DigestMatchers() {
		if digestMatcherType(digestMatcher) == enumerations.DigestMatcherTypeSigDEntry &&
			digestMatcher.DataFound && digestMatcher.DataIntact && digestMatcher.DocumentName != nil {
			for _, certificate := range certificates {
				if *digestMatcher.DocumentName == certificate.Id() {
					result = append(result, certificate)
				}
			}
		}
	}
	return result
}

// digestMatcherType reads XmlDigestMatcher#getType(), tolerating the nil the
// generated model uses for Java's null (mirrors bbb/cv's digestMatcherType helper).
func digestMatcherType(digestMatcher *diagnosticjaxb.XmlDigestMatcher) enumerations.DigestMatcherType {
	if digestMatcher.Type == nil {
		return ""
	}
	return digestMatcher.Type.DigestMatcherType()
}
