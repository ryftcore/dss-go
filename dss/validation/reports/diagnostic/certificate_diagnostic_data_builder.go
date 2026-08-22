// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/CertificateDiagnosticDataBuilder.java (DSS 6.5.RC1).
package diagnostic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// CertificateDiagnosticDataBuilder builds the Data for a CertificateToken validation.
type CertificateDiagnosticDataBuilder struct {
	DataBuilder
}

// NewCertificateDiagnosticDataBuilder is the port of the default constructor.
func NewCertificateDiagnosticDataBuilder() *CertificateDiagnosticDataBuilder {
	b := &CertificateDiagnosticDataBuilder{
		DataBuilder: *NewDataBuilder(),
	}
	b.InitDiagnosticDataBuilder(b)
	return b
}

// Build builds the XmlDiagnosticData. Port of the public @Override build().
func (b *CertificateDiagnosticDataBuilder) Build() *jaxb.XmlDiagnosticData {
	diagnosticData := b.DataBuilder.Build()

	diagnosticData.OrphanTokens = b.BuildXmlOrphanTokens()

	return diagnosticData
}

// UsedCertificates re-declares the fluent setter with the CertificateDiagnosticDataBuilder
// return type. Port of the covariant-return @Override usedCertificates(Set<CertificateToken>).
func (b *CertificateDiagnosticDataBuilder) UsedCertificates(usedCertificates []*model.CertificateToken) *CertificateDiagnosticDataBuilder {
	b.DataBuilder.UsedCertificates(usedCertificates)
	return b
}

// UsedRevocations re-declares the fluent setter with the CertificateDiagnosticDataBuilder
// return type. Port of the covariant-return @Override usedRevocations(Set<RevocationToken<?>>).
func (b *CertificateDiagnosticDataBuilder) UsedRevocations(usedRevocations []validation.AnyRevocationToken) *CertificateDiagnosticDataBuilder {
	b.DataBuilder.UsedRevocations(usedRevocations)
	return b
}

// AllCertificateSources re-declares the fluent setter with the CertificateDiagnosticDataBuilder
// return type. Port of the covariant-return @Override allCertificateSources(ListCertificateSource).
func (b *CertificateDiagnosticDataBuilder) AllCertificateSources(trustedCertSources *spi.ListCertificateSource) *CertificateDiagnosticDataBuilder {
	b.DataBuilder.AllCertificateSources(trustedCertSources)
	return b
}

// ValidationDate re-declares the fluent setter with the CertificateDiagnosticDataBuilder return
// type. Port of the covariant-return @Override validationDate(Date).
func (b *CertificateDiagnosticDataBuilder) ValidationDate(validationDate time.Time) *CertificateDiagnosticDataBuilder {
	b.DataBuilder.ValidationDate(validationDate)
	return b
}

// TokenExtractionStrategy re-declares the fluent setter with the
// CertificateDiagnosticDataBuilder return type. Port of the covariant-return @Override
// tokenExtractionStrategy(TokenExtractionStrategy).
func (b *CertificateDiagnosticDataBuilder) TokenExtractionStrategy(tokenExtractionStrategy enumerations.TokenExtractionStrategy) *CertificateDiagnosticDataBuilder {
	b.DataBuilder.TokenExtractionStrategy(tokenExtractionStrategy)
	return b
}

// DefaultDigestAlgorithm re-declares the fluent setter with the CertificateDiagnosticDataBuilder
// return type. Port of the covariant-return @Override defaultDigestAlgorithm(DigestAlgorithm).
func (b *CertificateDiagnosticDataBuilder) DefaultDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *CertificateDiagnosticDataBuilder {
	b.DataBuilder.DefaultDigestAlgorithm(digestAlgorithm)
	return b
}
