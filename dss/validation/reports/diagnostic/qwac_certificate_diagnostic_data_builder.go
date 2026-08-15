// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/QWACCertificateDiagnosticDataBuilder.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

// QWACCertificateDiagnosticDataBuilder builds a Diagnostic Data report for a QWAC certificate
// validation.
type QWACCertificateDiagnosticDataBuilder struct {
	SignedDocumentDiagnosticDataBuilder

	// websiteUrl is the URL used to establish a remote connection to verify a TLS/SSL
	// certificate for.
	websiteUrl string

	// tlsCertificate is the TLS/SSL certificate returned by a remote server during the TLS/SSL
	// handshake.
	tlsCertificate *model.CertificateToken

	// tlsCertificateBindingUrl is the TLS Certificate Binding URL, when present (under the
	// 'Link' response header).
	tlsCertificateBindingUrl string

	// tlsCertificateBindingSignature is the TLS Certificate Binding Signature, when present.
	tlsCertificateBindingSignature validation.AdvancedSignature

	// signatureDiagnosticDataBuilder is the builder used to build a signature object.
	signatureDiagnosticDataBuilder *SignedDocumentDiagnosticDataBuilder
}

// NewQWACCertificateDiagnosticDataBuilder is the port of the default constructor.
func NewQWACCertificateDiagnosticDataBuilder() *QWACCertificateDiagnosticDataBuilder {
	b := &QWACCertificateDiagnosticDataBuilder{
		SignedDocumentDiagnosticDataBuilder: *NewSignedDocumentDiagnosticDataBuilder(),
	}
	b.InitSignedDocumentDiagnosticDataBuilder(b)
	return b
}

// WebsiteUrl sets the website URL used to establish a TLS connection. Port of websiteUrl(String).
func (b *QWACCertificateDiagnosticDataBuilder) WebsiteUrl(websiteUrl string) *QWACCertificateDiagnosticDataBuilder {
	b.websiteUrl = websiteUrl
	return b
}

// TLSCertificate sets the TLS/SSL certificate obtained during the handshake. Port of
// tlsCertificate(CertificateToken).
func (b *QWACCertificateDiagnosticDataBuilder) TLSCertificate(tlsCertificate *model.CertificateToken) *QWACCertificateDiagnosticDataBuilder {
	b.tlsCertificate = tlsCertificate
	return b
}

// TLSCertificateBindingUrl sets the TLS Certificate Binding URL, when present (under the
// 'Link' response header). Port of tlsCertificateBindingUrl(String).
func (b *QWACCertificateDiagnosticDataBuilder) TLSCertificateBindingUrl(tlsCertificateBindingUrl string) *QWACCertificateDiagnosticDataBuilder {
	b.tlsCertificateBindingUrl = tlsCertificateBindingUrl
	return b
}

// TLSCertificateBindingSignature sets the TLS Certificate Binding signature, when present.
// Port of tlsCertificateBindingSignature(AdvancedSignature).
func (b *QWACCertificateDiagnosticDataBuilder) TLSCertificateBindingSignature(tlsCertificateBindingSignature validation.AdvancedSignature) *QWACCertificateDiagnosticDataBuilder {
	b.tlsCertificateBindingSignature = tlsCertificateBindingSignature
	return b
}

// SetSignatureDiagnosticDataBuilder sets a builder for a signature object. Port of
// setSignatureDiagnosticDataBuilder(SignedDocumentDiagnosticDataBuilder).
func (b *QWACCertificateDiagnosticDataBuilder) SetSignatureDiagnosticDataBuilder(signatureDiagnosticDataBuilder *SignedDocumentDiagnosticDataBuilder) *QWACCertificateDiagnosticDataBuilder {
	b.signatureDiagnosticDataBuilder = signatureDiagnosticDataBuilder
	return b
}

// Build builds the XmlDiagnosticData. Port of the public @Override build().
func (b *QWACCertificateDiagnosticDataBuilder) Build() *jaxb.XmlDiagnosticData {
	xmlDiagnosticData := b.SignedDocumentDiagnosticDataBuilder.Build()
	xmlDiagnosticData.ConnectionInfo = b.buildXmlConnectionInfo()
	return xmlDiagnosticData
}

func (b *QWACCertificateDiagnosticDataBuilder) buildXmlConnectionInfo() *jaxb.XmlConnectionInfo {
	xmlConnectionInfo := &jaxb.XmlConnectionInfo{}
	if b.websiteUrl != "" {
		url := b.websiteUrl
		xmlConnectionInfo.Url = &url
	}
	if b.tlsCertificate != nil {
		xmlTLSCertificate := &jaxb.XmlTLSCertificate{}
		xmlTLSCertificate.Certificate = b.xmlCertsMap[b.tlsCertificate.DSSIDAsString()]
		xmlConnectionInfo.TLSCertificate = xmlTLSCertificate
	}
	if b.tlsCertificateBindingUrl != "" {
		url := b.tlsCertificateBindingUrl
		xmlConnectionInfo.TLSCertificateBindingUrl = &url
	}
	if b.tlsCertificateBindingSignature != nil {
		xmlTLSCertificateBindingSignature := &jaxb.XmlTLSCertificateBindingSignature{}
		xmlTLSCertificateBindingSignature.Signature = b.xmlSignaturesMap[b.tlsCertificateBindingSignature.ID()]
		xmlConnectionInfo.TLSCertificateBindingSignature = xmlTLSCertificateBindingSignature
	}
	return xmlConnectionInfo
}

// BuildDetachedXmlSignature overrides the base to delegate to signatureDiagnosticDataBuilder
// and to identify TLS certificates in the resulting signature. Port of the public @Override
// buildDetachedXmlSignature(AdvancedSignature).
func (b *QWACCertificateDiagnosticDataBuilder) BuildDetachedXmlSignature(sig validation.AdvancedSignature) *jaxb.XmlSignature {
	xmlSignature := b.signatureDiagnosticDataBuilder.BuildDetachedXmlSignature(sig)
	b.identifyTLSCertificates(xmlSignature)
	return xmlSignature
}

func (b *QWACCertificateDiagnosticDataBuilder) identifyTLSCertificates(xmlSignature *jaxb.XmlSignature) {
	for _, digestMatcher := range xmlSignature.DigestMatchers.All() {
		if digestMatcher.Type != nil && enumerations.DigestMatcherType_SIG_D_ENTRY == digestMatcher.Type.DigestMatcherType() &&
			digestMatcher.DataFound && digestMatcher.DataIntact {
			tlsCertificate := b.getMatchingTLSCertificate(digestMatcher)
			if tlsCertificate != nil {
				digestMatcher.DataObjectReferences = &jaxb.DataObjectReferencesWrapper{
					Items: []string{b.identifierProvider.IDAsString(tlsCertificate)},
				}
			}
		}
	}
}

func (b *QWACCertificateDiagnosticDataBuilder) getMatchingTLSCertificate(digestMatcher *jaxb.XmlDigestMatcher) *model.CertificateToken {
	if digestMatcher.DigestMethod == nil || digestMatcher.DigestValue == nil {
		return nil
	}
	algo := digestMatcher.DigestMethod.DigestAlgorithm()
	for _, certificate := range b.usedCertificates {
		digest, err := certificate.Digest(algo)
		if err != nil {
			continue
		}
		if bytesEqual(*digestMatcher.DigestValue, digest) {
			return certificate
		}
	}
	return nil
}

// AssertConfigurationValid overrides the base to also require websiteUrl. Port of the
// protected @Override assertConfigurationValid().
func (b *QWACCertificateDiagnosticDataBuilder) AssertConfigurationValid() {
	if b.websiteUrl == "" {
		panic("websiteUrl shall be provided!")
	}
}
