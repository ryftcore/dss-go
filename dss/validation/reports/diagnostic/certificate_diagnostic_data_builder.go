// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/CertificateDiagnosticDataBuilder.java (DSS 6.5.RC1).
package diagnostic

import (
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
)

// CertificateDiagnosticDataBuilder builds the DiagnosticData for a CertificateToken validation.
type CertificateDiagnosticDataBuilder struct {
	DiagnosticDataBuilder
}

// NewCertificateDiagnosticDataBuilder is the port of the default constructor.
func NewCertificateDiagnosticDataBuilder() *CertificateDiagnosticDataBuilder {
	b := &CertificateDiagnosticDataBuilder{
		DiagnosticDataBuilder: *NewDiagnosticDataBuilder(),
	}
	b.InitDiagnosticDataBuilder(b)
	return b
}

// Build builds the XmlDiagnosticData. Port of the public @Override build().
func (b *CertificateDiagnosticDataBuilder) Build() *jaxb.XmlDiagnosticData {
	diagnosticData := b.DiagnosticDataBuilder.Build()

	diagnosticData.OrphanTokens = b.BuildXmlOrphanTokens()

	return diagnosticData
}

// UsedCertificates re-declares the fluent setter with the CertificateDiagnosticDataBuilder
// return type. Port of the covariant-return @Override usedCertificates(Set<CertificateToken>).
func (b *CertificateDiagnosticDataBuilder) UsedCertificates(usedCertificates []*model.CertificateToken) *CertificateDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.UsedCertificates(usedCertificates)
	return b
}

// UsedRevocations re-declares the fluent setter with the CertificateDiagnosticDataBuilder
// return type. Port of the covariant-return @Override usedRevocations(Set<RevocationToken<?>>).
func (b *CertificateDiagnosticDataBuilder) UsedRevocations(usedRevocations []validation.AnyRevocationToken) *CertificateDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.UsedRevocations(usedRevocations)
	return b
}

// AllCertificateSources re-declares the fluent setter with the CertificateDiagnosticDataBuilder
// return type. Port of the covariant-return @Override allCertificateSources(ListCertificateSource).
func (b *CertificateDiagnosticDataBuilder) AllCertificateSources(trustedCertSources *spi.ListCertificateSource) *CertificateDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.AllCertificateSources(trustedCertSources)
	return b
}

// ValidationDate re-declares the fluent setter with the CertificateDiagnosticDataBuilder return
// type. Port of the covariant-return @Override validationDate(Date).
func (b *CertificateDiagnosticDataBuilder) ValidationDate(validationDate time.Time) *CertificateDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.ValidationDate(validationDate)
	return b
}

// TokenExtractionStrategy re-declares the fluent setter with the
// CertificateDiagnosticDataBuilder return type. Port of the covariant-return @Override
// tokenExtractionStrategy(TokenExtractionStrategy).
func (b *CertificateDiagnosticDataBuilder) TokenExtractionStrategy(tokenExtractionStrategy enumerations.TokenExtractionStrategy) *CertificateDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.TokenExtractionStrategy(tokenExtractionStrategy)
	return b
}

// DefaultDigestAlgorithm re-declares the fluent setter with the CertificateDiagnosticDataBuilder
// return type. Port of the covariant-return @Override defaultDigestAlgorithm(DigestAlgorithm).
func (b *CertificateDiagnosticDataBuilder) DefaultDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *CertificateDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.DefaultDigestAlgorithm(digestAlgorithm)
	return b
}
