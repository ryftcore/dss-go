// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/XmlDiagnosticDataFactory.java (DSS 6.5.RC1).
package diagnostic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// XmlDiagnosticDataFactory creates an XmlDiagnosticData.
type XmlDiagnosticDataFactory struct {
	// diagnosticDataBuilder is the corresponding diagnostic data builder.
	diagnosticDataBuilder *SignedDocumentDiagnosticDataBuilder

	// document is the document to be validated.
	document model.DSSDocument

	// validationTime is the validation time.
	validationTime time.Time

	// validationContext is the current validation context.
	validationContext validation.ValidationContext

	// defaultDigestAlgorithm is the used default digest algorithm for tokens definition.
	defaultDigestAlgorithm enumerations.DigestAlgorithm

	// tokenExtractionStrategy is the token extraction strategy to be used (i.e. binaries vs
	// digest).
	tokenExtractionStrategy enumerations.TokenExtractionStrategy

	// tokenIdentifierProvider is the class to compute identifiers for tokens to be returned in
	// the reports.
	tokenIdentifierProvider model.TokenIdentifierProvider
}

// NewXmlDiagnosticDataFactory is the port of the default constructor: panics (Java
// requireNonNull) if diagnosticDataBuilder is nil.
func NewXmlDiagnosticDataFactory(diagnosticDataBuilder *SignedDocumentDiagnosticDataBuilder) *XmlDiagnosticDataFactory {
	if diagnosticDataBuilder == nil {
		panic("SignedDocumentDiagnosticDataBuilder is null!")
	}
	return &XmlDiagnosticDataFactory{diagnosticDataBuilder: diagnosticDataBuilder}
}

// SetDocument sets the original document to be validated. Port of setDocument(DSSDocument).
func (f *XmlDiagnosticDataFactory) SetDocument(document model.DSSDocument) *XmlDiagnosticDataFactory {
	f.document = document
	return f
}

// SetValidationTime sets the validation time. Port of setValidationTime(Date).
func (f *XmlDiagnosticDataFactory) SetValidationTime(validationTime time.Time) *XmlDiagnosticDataFactory {
	f.validationTime = validationTime
	return f
}

// SetValidationContext sets the validation context. Port of setValidationContext(ValidationContext).
func (f *XmlDiagnosticDataFactory) SetValidationContext(validationContext validation.ValidationContext) *XmlDiagnosticDataFactory {
	f.validationContext = validationContext
	return f
}

// SetDefaultDigestAlgorithm sets the digest algorithm to be used to compute references to the
// data objects. Port of setDefaultDigestAlgorithm(DigestAlgorithm).
func (f *XmlDiagnosticDataFactory) SetDefaultDigestAlgorithm(defaultDigestAlgorithm enumerations.DigestAlgorithm) *XmlDiagnosticDataFactory {
	f.defaultDigestAlgorithm = defaultDigestAlgorithm
	return f
}

// SetTokenExtractionStrategy sets the token extraction strategy. Port of
// setTokenExtractionStrategy(TokenExtractionStrategy).
func (f *XmlDiagnosticDataFactory) SetTokenExtractionStrategy(tokenExtractionStrategy enumerations.TokenExtractionStrategy) *XmlDiagnosticDataFactory {
	f.tokenExtractionStrategy = tokenExtractionStrategy
	return f
}

// SetTokenIdentifierProvider sets the token identifier provider. Port of
// setTokenIdentifierProvider(TokenIdentifierProvider).
func (f *XmlDiagnosticDataFactory) SetTokenIdentifierProvider(tokenIdentifierProvider model.TokenIdentifierProvider) *XmlDiagnosticDataFactory {
	f.tokenIdentifierProvider = tokenIdentifierProvider
	return f
}

// Create creates an XmlDiagnosticData. Port of create().
//
// Dispatches through the registered SignedDocumentDiagnosticDataBuilderOverrides rather than
// calling Build() directly on the base-typed f.diagnosticDataBuilder: Java's virtual dispatch
// reaches ASiCContainerDiagnosticDataBuilder.build()'s ContainerInfo-adding @Override through
// this same call in XmlDiagnosticDataFactory.create(), but f.diagnosticDataBuilder is statically
// typed *SignedDocumentDiagnosticDataBuilder (SignedDocumentValidatorOverrides.
// InitializeDiagnosticDataBuilder's declared return type - see validation/signed_document_validator.go),
// so a direct .Build() call always resolves to this package's own base implementation and never
// reaches an embedding package's override. Added to SignedDocumentDiagnosticDataBuilderOverrides
// during phase 8f un-gating for the same reason as BuildXmlOrphanTokens - see that method's doc
// comment. Every existing implementer (CAdES, XAdES, JAdES, QWAC's signature-diagnostic-data use)
// keeps its current behavior: none define their own Build, so embedding promotes this package's
// concrete implementation for them unchanged.
func (f *XmlDiagnosticDataFactory) Create() *jaxb.XmlDiagnosticData {
	builder := f.InitBuilder()
	return builder.signedDocumentDiagnosticDataBuilderOverrides().Build()
}

// InitBuilder instantiates the Diagnostic Data builder with the validation context. Port of the
// protected initBuilder().
func (f *XmlDiagnosticDataFactory) InitBuilder() *SignedDocumentDiagnosticDataBuilder {
	f.diagnosticDataBuilder.Document(f.document)
	f.diagnosticDataBuilder.ValidationDate(f.validationTime)
	f.diagnosticDataBuilder.FoundSignatures(f.validationContext.GetProcessedSignatures())
	f.diagnosticDataBuilder.UsedTimestamps(f.validationContext.GetProcessedTimestamps())
	f.diagnosticDataBuilder.FoundEvidenceRecords(f.validationContext.GetProcessedEvidenceRecords())
	f.diagnosticDataBuilder.AllCertificateSources(f.validationContext.GetAllCertificateSources())
	f.diagnosticDataBuilder.DocumentCertificateSource(f.validationContext.GetDocumentCertificateSource())
	f.diagnosticDataBuilder.DocumentCRLSource(f.validationContext.GetDocumentCRLSource())
	f.diagnosticDataBuilder.DocumentOCSPSource(f.validationContext.GetDocumentOCSPSource())
	f.diagnosticDataBuilder.UsedCertificates(f.validationContext.GetProcessedCertificates())
	f.diagnosticDataBuilder.UsedRevocations(f.validationContext.GetProcessedRevocations())
	f.diagnosticDataBuilder.DefaultDigestAlgorithm(f.defaultDigestAlgorithm)
	f.diagnosticDataBuilder.TokenExtractionStrategy(f.tokenExtractionStrategy)
	f.diagnosticDataBuilder.TokenIdentifierProvider(f.tokenIdentifierProvider)
	return f.diagnosticDataBuilder
}
