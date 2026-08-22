// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/DocumentAnalyzer.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.validation.analyzer lands in this Go package;
// spi.validation.analyzer.evidencerecord (EvidenceRecordAnalyzer, EvidenceRecordAnalyzerFactory)
// flattens into it too.
package analyzer

import (
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/policy"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/executor"
)

// DocumentAnalyzer performs processing of a signature document, including extraction of
// signature and timestamp tokens, cryptographic validation, certificate chain building and
// revocation data validation. The class works exclusively with Go objects, and does not
// include ETSI EN 319 102-1 validation process, nor JAXB objects.
type DocumentAnalyzer interface {
	// Document gets the document to be validated. Port of getDocument().
	Document() model.DSSDocument

	// Signatures retrieves the signatures found in the document. Port of getSignatures().
	Signatures() []validation.AdvancedSignature

	// DetachedTimestamps retrieves the detached timestamps found in the document. Port of
	// getDetachedTimestamps().
	DetachedTimestamps() []*validation.TimestampToken

	// DetachedEvidenceRecords retrieves the detached evidence records found in the document.
	// Port of getDetachedEvidenceRecords().
	DetachedEvidenceRecords() []validation.EvidenceRecord

	// SetCertificateVerifier provides a CertificateVerifier to be used during the validation
	// process. Port of setCertificateVerifier(CertificateVerifier).
	SetCertificateVerifier(certificateVerifier validation.CertificateVerifier)

	// SetValidationContextExecutor sets the ValidationContextExecutor for validation of the
	// prepared Context. Default: executor.DefaultValidationContextExecutor (performs
	// basic validation of tokens, including certificate chain building and revocation data
	// extraction, without processing of validity checks). Port of
	// setValidationContextExecutor(ValidationContextExecutor).
	SetValidationContextExecutor(validationContextExecutor executor.ValidationContextExecutor)

	// TokenIdentifierProvider gets the TokenIdentifierProvider. Port of
	// getTokenIdentifierProvider().
	TokenIdentifierProvider() model.TokenIdentifierProvider

	// SetTokenIdentifierProvider sets the TokenIdentifierProvider. Port of
	// setTokenIdentifierProvider(TokenIdentifierProvider).
	SetTokenIdentifierProvider(tokenIdentifierProvider model.TokenIdentifierProvider)

	// ValidationTime returns the document validation time. Port of getValidationTime().
	ValidationTime() time.Time

	// SetValidationTime allows defining a custom validation time. Port of
	// setValidationTime(Date).
	SetValidationTime(validationTime time.Time)

	// SetDetachedContents sets the list of DSSDocument containing the original contents to
	// sign, for detached signature scenarios. Port of setDetachedContents(List).
	SetDetachedContents(detachedContent []model.DSSDocument)

	// SetDetachedEvidenceRecordDocuments sets a list of DSSDocument containing the evidence
	// record documents covering the signature document. Port of
	// setDetachedEvidenceRecordDocuments(List).
	SetDetachedEvidenceRecordDocuments(detachedEvidenceRecordDocuments []model.DSSDocument)

	// SetContainerContents sets the list of DSSDocument containing the original container
	// content for ASiC-S signatures. Port of setContainerContents(List).
	SetContainerContents(archiveContents []model.DSSDocument)

	// SetManifestFile sets a related ManifestFile to the document to be validated. Port of
	// setManifestFile(ManifestFile).
	SetManifestFile(manifestFile *model.ManifestFile)

	// IsSupported checks if the document is supported by the current validator. Port of
	// isSupported(DSSDocument).
	IsSupported(dssDocument model.DSSDocument) bool

	// SetSigningCertificateSource sets a certificate source which allows finding the signing
	// certificate by kid or certificate's digest. Port of
	// setSigningCertificateSource(CertificateSource).
	SetSigningCertificateSource(certificateSource spi.CertificateSource)

	// SetSignaturePolicyProvider allows setting a provider for Signature policies. Port of
	// setSignaturePolicyProvider(SignaturePolicyProvider).
	SetSignaturePolicyProvider(signaturePolicyProvider *policy.SignaturePolicyProvider)

	// SetSignaturePolicyValidatorLoader sets a loader for a SignaturePolicyValidator. Port of
	// setSignaturePolicyValidatorLoader(SignaturePolicyValidatorLoader).
	SetSignaturePolicyValidatorLoader(signaturePolicyValidatorLoader policy.SignaturePolicyValidatorLoader)

	// OriginalDocuments returns the signed document(s) without their signature(s), looked up by
	// the DSS ID of the signature. Port of the getOriginalDocuments(String) overload.
	OriginalDocuments(signatureId string) []model.DSSDocument

	// OriginalDocumentsForSignature returns the signed document(s) without their signature(s).
	// Port of the getOriginalDocuments(AdvancedSignature) overload.
	OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument

	// GetValidationData extracts a validation data for the provided collection of signatures.
	// Port of the getValidationData(Collection) overload.
	GetValidationData(signatures []validation.AdvancedSignature) (*validation.DataContainer, error)

	// GetValidationDataWithTimestamps extracts a validation data for the provided collection of
	// signatures and/or timestamps. Port of the getValidationData(Collection, Collection)
	// overload.
	GetValidationDataWithTimestamps(signatures []validation.AdvancedSignature,
		detachedTimestamps []*validation.TimestampToken) (*validation.DataContainer, error)

	// Validate performs validation of the document. Port of validate().
	Validate() validation.Context
}
