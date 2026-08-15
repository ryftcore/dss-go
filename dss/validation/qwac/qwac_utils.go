// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/qwac/QWACUtils.java (DSS 6.5.RC1).
//
// MINIMAL PORT: only GetIdentifiedTLSCertificates is ported. Java's
// QWACUtils also carries getTLSCertificateBindingUrl/isTLSCertificateBindingRel
// (parsing an HTTP "Link" response header via a private LinkHeaderParser
// helper class) but nothing under validation/process/qualification calls
// them, and no Go port of LinkHeaderParser exists on disk; porting them
// would mean inventing an unverified LinkHeaderParser port with no caller
// to check it against. Left unported; add both alongside a LinkHeaderParser
// port if/when a caller needs them.
package qwac

import (
	"github.com/utain/esig/dss/diagnostic"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// GetIdentifiedTLSCertificates gets TLS Binding Certificates identified from
// the binding signature. Port of getIdentifiedTLSCertificates(SignatureWrapper, List<CertificateWrapper>).
func GetIdentifiedTLSCertificates(signature *diagnostic.SignatureWrapper, certificates []*diagnostic.CertificateWrapper) []*diagnostic.CertificateWrapper {
	var result []*diagnostic.CertificateWrapper
	for _, digestMatcher := range signature.DigestMatchers() {
		if digestMatcherType(digestMatcher) == enumerations.DigestMatcherType_SIG_D_ENTRY &&
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
