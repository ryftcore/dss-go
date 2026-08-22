// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/DocumentValidator.java
// (DSS 6.5.RC1).
//
// Java overloads validateDocument() on the policy source (URL, classpath
// resource path, File, DSSDocument, InputStream, ValidationPolicy) and again
// on the (policy, cryptographic suite) pair. Go has no overloading, so each
// overload keeps its own name, suffixed by the source it reads - matching the
// naming already used by dss/validation/policy.ValidationPolicyLoader
// (FromValidationPolicyDocument / ...Reader / ...File / ...Path).
//
// Java's classpath-resource overloads (String policyResourcePath) have no Go
// counterpart: Go has no classpath, so - like
// dss/policy/validation_policy_facade.go's GetValidationPolicyFromPath - the
// path is read straight from the filesystem.
//
// The java.net.URL overloads (validateDocument(URL[, URL])) are not ported:
// they only wrap URL.openStream() around the InputStream overload, and the
// ported dss/validation/policy.ValidationPolicyLoader likewise has no URL
// entry point. Callers open the stream themselves and use the ...Reader
// overload.

package validation

import (
	"io"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	spi "github.com/ryftcore/dss-go/dss/spi"
	spipolicy "github.com/ryftcore/dss-go/dss/spi/policy"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
	spiexecutor "github.com/ryftcore/dss-go/dss/spi/validation/executor"
	"github.com/ryftcore/dss-go/dss/validation/executor"
	"github.com/ryftcore/dss-go/dss/validation/reports"
)

// DocumentValidator is the interface of a document validator. Port of the
// DocumentValidator interface (extends
// ProcessExecutorProvider<DocumentProcessExecutor>).
type DocumentValidator interface {
	executor.ProcessExecutorProvider[executor.DocumentProcessExecutor]

	// Signatures retrieves the signatures found in the document. Port of
	// getSignatures().
	Signatures() []spivalidation.AdvancedSignature

	// DetachedTimestamps retrieves the detached timestamps found in the
	// document. Port of getDetachedTimestamps().
	DetachedTimestamps() []*spivalidation.TimestampToken

	// DetachedEvidenceRecords retrieves the detached evidence records found
	// in the document. Port of getDetachedEvidenceRecords().
	DetachedEvidenceRecords() []spivalidation.EvidenceRecord

	// SetCertificateVerifier provides the external sources of certificates
	// and revocation data needed by the validation process. Port of
	// setCertificateVerifier(CertificateVerifier).
	SetCertificateVerifier(certificateVerifier spivalidation.CertificateVerifier)

	// SetValidationContextExecutor sets the ValidationContextExecutor for
	// validation of the prepared ValidationContext. Port of
	// setValidationContextExecutor(ValidationContextExecutor).
	SetValidationContextExecutor(validationContextExecutor spiexecutor.ValidationContextExecutor)

	// SetDefaultDigestAlgorithm sets the DigestAlgorithm used for tokens'
	// digest calculation. Port of setDefaultDigestAlgorithm(DigestAlgorithm).
	SetDefaultDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm)

	// SetTokenExtractionStrategy sets the TokenExtractionStrategy. Port of
	// setTokenExtractionStrategy(TokenExtractionStrategy).
	SetTokenExtractionStrategy(tokenExtractionStrategy enumerations.TokenExtractionStrategy)

	// SetTokenIdentifierProvider sets the TokenIdentifierProvider. Port of
	// setTokenIdentifierProvider(TokenIdentifierProvider).
	SetTokenIdentifierProvider(tokenIdentifierProvider model.TokenIdentifierProvider)

	// SetIncludeSemantics allows including the semantics for Indication /
	// SubIndication. Port of setIncludeSemantics(boolean).
	SetIncludeSemantics(include bool)

	// SetValidationTime allows defining a custom validation time. Port of
	// setValidationTime(Date).
	SetValidationTime(validationTime time.Time)

	// SetDetachedContents sets the original contents for detached signature
	// scenarios. Port of setDetachedContents(List).
	SetDetachedContents(detachedContent []model.DSSDocument)

	// SetDetachedEvidenceRecordDocuments sets the evidence record documents
	// covering the signature document. Port of
	// setDetachedEvidenceRecordDocuments(List).
	SetDetachedEvidenceRecordDocuments(detachedEvidenceRecordDocuments []model.DSSDocument)

	// SetContainerContents sets the original container content for ASiC-S
	// signatures. Port of setContainerContents(List).
	SetContainerContents(archiveContents []model.DSSDocument)

	// SetManifestFile sets a related ManifestFile to the document to be
	// validated. Port of setManifestFile(ManifestFile).
	SetManifestFile(manifestFile *model.ManifestFile)

	// SetSigningCertificateSource sets a certificate source which allows
	// finding the signing certificate by kid or certificate's digest. Port of
	// setSigningCertificateSource(CertificateSource).
	SetSigningCertificateSource(certificateSource spi.CertificateSource)

	// SetValidationLevel sets the expected validation level. Port of
	// setValidationLevel(ValidationLevel).
	SetValidationLevel(validationLevel enumerations.ValidationLevel)

	// SetEnableEtsiValidationReport specifies if the ETSI Validation Report
	// shall be produced. Port of setEnableEtsiValidationReport(boolean).
	SetEnableEtsiValidationReport(enableEtsiValidationReport bool)

	// SetSignaturePolicyProvider allows setting a provider for signature
	// policies. Port of setSignaturePolicyProvider(SignaturePolicyProvider).
	SetSignaturePolicyProvider(signaturePolicyProvider *spipolicy.SignaturePolicyProvider)

	// SetSignaturePolicyValidatorLoader sets a loader for a
	// SignaturePolicyValidator. Port of
	// setSignaturePolicyValidatorLoader(SignaturePolicyValidatorLoader).
	SetSignaturePolicyValidatorLoader(signaturePolicyValidatorLoader spipolicy.SignaturePolicyValidatorLoader)

	// ValidateDocument validates the document and all its signatures using
	// the default validation policy. Port of the no-arg validateDocument().
	ValidateDocument() (*reports.Reports, error)

	// ValidateDocumentWithPolicyPath validates the document with the policy
	// read from policyResourcePath. Port of the validateDocument(String)
	// overload; see the file header on the classpath deviation.
	ValidateDocumentWithPolicyPath(policyResourcePath string) (*reports.Reports, error)

	// ValidateDocumentWithPolicyFile validates the document with the policy
	// read from policyFile. Port of the validateDocument(File) overload.
	ValidateDocumentWithPolicyFile(policyFile string) (*reports.Reports, error)

	// ValidateDocumentWithPolicyDocument validates the document with the
	// policy carried by policyDocument. Port of the
	// validateDocument(DSSDocument) overload.
	ValidateDocumentWithPolicyDocument(policyDocument model.DSSDocument) (*reports.Reports, error)

	// ValidateDocumentWithPolicyReader validates the document with the policy
	// read from policyDataStream. Port of the validateDocument(InputStream)
	// overload.
	ValidateDocumentWithPolicyReader(policyDataStream io.Reader) (*reports.Reports, error)

	// ValidateDocumentWithValidationPolicy validates the document with the
	// given ValidationPolicy. Port of the validateDocument(ValidationPolicy)
	// overload.
	ValidateDocumentWithValidationPolicy(validationPolicy modelpolicy.ValidationPolicy) (*reports.Reports, error)

	// ValidateDocumentWithPolicyAndCryptographicSuitePath validates the
	// document with the policy and cryptographic suite read from the given
	// paths. Port of the validateDocument(String, String) overload.
	ValidateDocumentWithPolicyAndCryptographicSuitePath(policyResourcePath, cryptographicSuitePath string) (*reports.Reports, error)

	// ValidateDocumentWithPolicyAndCryptographicSuiteFile validates the
	// document with the policy and cryptographic suite read from the given
	// files. Port of the validateDocument(File, File) overload.
	ValidateDocumentWithPolicyAndCryptographicSuiteFile(policyFile, cryptographicSuiteFile string) (*reports.Reports, error)

	// ValidateDocumentWithPolicyAndCryptographicSuiteDocument validates the
	// document with the policy and cryptographic suite carried by the given
	// documents. Port of the validateDocument(DSSDocument, DSSDocument)
	// overload.
	ValidateDocumentWithPolicyAndCryptographicSuiteDocument(policyDocument, cryptographicSuiteDocument model.DSSDocument) (*reports.Reports, error)

	// ValidateDocumentWithPolicyAndCryptographicSuiteReader validates the
	// document with the policy and cryptographic suite read from the given
	// streams. Port of the validateDocument(InputStream, InputStream)
	// overload.
	ValidateDocumentWithPolicyAndCryptographicSuiteReader(policyDataStream, cryptographicSuiteStream io.Reader) (*reports.Reports, error)

	// OriginalDocuments returns the signed document(s) without their
	// signature(s), looked up by signature id. Port of the
	// getOriginalDocuments(String) overload.
	OriginalDocuments(signatureId string) []model.DSSDocument

	// OriginalDocumentsForSignature returns the signed document(s) without
	// their signature(s). Port of the getOriginalDocuments(AdvancedSignature)
	// overload.
	OriginalDocumentsForSignature(advancedSignature spivalidation.AdvancedSignature) []model.DSSDocument

	// GetValidationData extracts the validation data for the provided
	// signatures. Port of the getValidationData(Collection) overload.
	GetValidationData(signatures []spivalidation.AdvancedSignature) (*spivalidation.ValidationDataContainer, error)

	// GetValidationDataWithTimestamps extracts the validation data for the
	// provided signatures and detached timestamps. Port of the
	// getValidationData(Collection, Collection) overload.
	GetValidationDataWithTimestamps(signatures []spivalidation.AdvancedSignature,
		detachedTimestamps []*spivalidation.TimestampToken) (*spivalidation.ValidationDataContainer, error)
}
