// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CMSDocumentAnalyzer.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	mscope "github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	spiscope "github.com/ryftcore/dss-go/dss/spi/validation/scope"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CMSDocumentAnalyzer is the validation of a CMS document. Port of the class
// CMSDocumentAnalyzer, extending analyzer.DefaultDocumentAnalyzer.
type CMSDocumentAnalyzer struct {
	analyzer.DefaultDocumentAnalyzer

	// cms is the CMS to be validated.
	cms *cms.CMS
}

// newCMSDocumentAnalyzer is the port of the package-private empty constructor.
func newCMSDocumentAnalyzer() *CMSDocumentAnalyzer {
	a := &CMSDocumentAnalyzer{DefaultDocumentAnalyzer: analyzer.NewDefaultDocumentAnalyzerBase()}
	a.InitDefaultDocumentAnalyzer(a)
	return a
}

// NewCMSDocumentAnalyzer creates a CMSDocumentAnalyzer from a CMS. Port of the constructor
// CMSDocumentAnalyzer(CMS).
func NewCMSDocumentAnalyzer(cmsObj *cms.CMS) *CMSDocumentAnalyzer {
	a := newCMSDocumentAnalyzer()
	a.cms = cmsObj
	return a
}

// NewCMSDocumentAnalyzerFromDocument creates a CMSDocumentAnalyzer from a DSSDocument. Port of
// the constructor CMSDocumentAnalyzer(DSSDocument).
//
// Panics when document is nil (Objects.requireNonNull("Document to be validated cannot be
// null!")). The IllegalInputException toCMS wraps a parse failure in propagates as-is.
func NewCMSDocumentAnalyzerFromDocument(document model.DSSDocument) (*CMSDocumentAnalyzer, error) {
	if document == nil {
		panic("Document to be validated cannot be null!")
	}
	a := newCMSDocumentAnalyzer()
	a.SetDocument(document)
	parsedCMS, err := cms.CMSUtilsParseToCMS(document)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause("A CMS file is expected : "+err.Error(), err)
	}
	a.cms = parsedCMS
	return a, nil
}

// IsSupported checks if the document is supported by the current validator. Port of
// isSupported(DSSDocument).
func (a *CMSDocumentAnalyzer) IsSupported(dssDocument model.DSSDocument) bool {
	firstByte, err := spi.DSSUtilsReadFirstByte(dssDocument)
	if err != nil {
		return false
	}
	if spi.DSSASN1UtilsIsASN1SequenceTag(firstByte) {
		return !cmsDocumentAnalyzerIsTimestampToken(dssDocument) && !analyzer.EvidenceRecordAnalyzerIsSupportedDocument(dssDocument)
	}
	return false
}

// cmsDocumentAnalyzerIsTimestampToken reports whether document parses as a CMS SignedData whose
// encapsulated content type is id-ct-TSTInfo. Port of DSSUtils.isTimestampToken(DSSDocument).
//
// NOT PORTED IN spi/dss_utils.go: that file's header defers isTimestampToken(DSSDocument) to
// "the CMS phase" as a CMS-shaped method, the same way dss_asn1_utils.go defers its own
// CMS-shaped TODOs. It is completed here rather than in spi/dss_utils.go itself because cades
// imports cms and spi, while nothing in spi imports cades or cms - spi cannot import cms, which
// this check needs to parse the document.
func cmsDocumentAnalyzerIsTimestampToken(document model.DSSDocument) bool {
	parsedCMS, err := cms.CMSUtilsParseToCMS(document)
	if err != nil {
		return false
	}
	return parsedCMS.SignedContentType().Equal(cmscore.OIDCTTSTInfo)
}

// BuildSignatures builds a list of signatures to be extracted from a document. Port of the
// protected buildSignatures() override.
func (a *CMSDocumentAnalyzer) BuildSignatures() []validation.AdvancedSignature {
	signatures := make([]validation.AdvancedSignature, 0)
	if a.cms != nil {
		for _, signerInformation := range a.cms.SignerInfos() {
			cadesSignature := NewCAdESSignature(a.cms, signerInformation)
			if a.HasDocument() {
				cadesSignature.SetFilename(a.Document().Name())
			}
			cadesSignature.SetDetachedContents(a.DetachedContents())
			cadesSignature.SetContainerContents(a.ContainerContents())
			cadesSignature.SetManifestFile(a.ManifestFile())
			cadesSignature.SetSigningCertificateSource(a.SigningCertificateSource())
			cadesSignature.InitBaselineRequirementsChecker(a.CertificateVerifier())
			a.ValidateSignaturePolicy(cadesSignature)
			signatures = append(signatures, cadesSignature)
		}
	}
	return signatures
}

// AppendExternalEvidenceRecords appends the detached evidence record provided to the validator
// to the corresponding signatures covered by the evidence record document, then - for CMS,
// where an embedded evidence record covers every SignerInformation - links it as an external
// evidence record to every other signature of the same CMS. Port of the protected
// appendExternalEvidenceRecords(List) override.
func (a *CMSDocumentAnalyzer) AppendExternalEvidenceRecords(allSignatureList []validation.AdvancedSignature) []validation.AdvancedSignature {
	allSignatureList = a.DefaultDocumentAnalyzer.AppendExternalEvidenceRecords(allSignatureList)

	// For CMS, an embedded ER covers all SignerInformation's
	// Add as external evidence record only for
	for _, signature := range allSignatureList {
		embeddedEvidenceRecords := signature.EmbeddedEvidenceRecords()
		if utils.IsCollectionNotEmpty(embeddedEvidenceRecords) {
			for _, coveredSignature := range allSignatureList {
				if signature != coveredSignature && cmsDocumentAnalyzerSameCMS(signature, coveredSignature) {
					for _, evidenceRecord := range embeddedEvidenceRecords {
						coveredSignature.AddExternalEvidenceRecord(evidenceRecord)
					}
					for _, evidenceRecord := range embeddedEvidenceRecords {
						a.addSignatureScope(evidenceRecord, coveredSignature)
					}
				}
			}
		}
	}
	return allSignatureList
}

// cmsDocumentAnalyzerSameCMS ports the private sameCMS(AdvancedSignature, AdvancedSignature).
func cmsDocumentAnalyzerSameCMS(signatureOne, signatureTwo validation.AdvancedSignature) bool {
	cadesSignatureOne := signatureOne.(*CAdESSignature)
	cadesSignatureTwo := signatureTwo.(*CAdESSignature)
	return cadesSignatureOne.CMS() == cadesSignatureTwo.CMS()
}

// addSignatureScope ports the private addSignatureScope(EvidenceRecord, AdvancedSignature).
func (a *CMSDocumentAnalyzer) addSignatureScope(evidenceRecord validation.EvidenceRecord, signature validation.AdvancedSignature) {
	evidenceRecordScopeFinder := NewCAdESEvidenceRecordScopeFinder(evidenceRecord, signature)
	evidenceRecordScopes := evidenceRecordScopeFinder.FindEvidenceRecordScope()
	evidenceRecord.SetEvidenceRecordScopes(evidenceRecordScopes)
	evidenceRecord.SetTimestampedReferences(
		cmsDocumentAnalyzerAddTimestampedReferences(evidenceRecord.TimestampedReferences(), evidenceRecordScopes))

	timestampScopeFinder := spiscope.NewEvidenceRecordTimestampScopeFinder(evidenceRecord)
	for _, timestampToken := range evidenceRecord.Timestamps() {
		timestampScopes := timestampToken.TimestampScopes()
		for _, evidenceRecordScope := range timestampScopeFinder.FindTimestampScope(timestampToken) {
			if !cmsDocumentAnalyzerContainsSignatureScope(timestampScopes, evidenceRecordScope) {
				timestampScopes = append(timestampScopes, evidenceRecordScope)
			}
		}
		timestampToken.SetTimestampScopes(timestampScopes)
		timestampToken.SetTimestampedReferences(
			cmsDocumentAnalyzerAddTimestampedReferences(timestampToken.TimestampedReferences(), evidenceRecordScopes))
	}
}

// cmsDocumentAnalyzerAddTimestampedReferences ports the private
// addTimestampedReferences(List<TimestampedReference>, List<SignatureScope>); see
// spi/validation.EvidenceRecord.SetTimestampedReferences's doc comment on why the result has to
// be written back explicitly rather than mutated in place.
func cmsDocumentAnalyzerAddTimestampedReferences(references []*validation.TimestampedReference,
	signatureScopes []mscope.SignatureScope) []*validation.TimestampedReference {
	for _, signatureScope := range signatureScopes {
		timestampedReference := validation.NewTimestampedReference(signatureScope.DSSIDAsString(), enumerations.TimestampedObjectTypeSignedData)
		if !cmsDocumentAnalyzerContainsTimestampedReference(references, timestampedReference) {
			references = append(references, timestampedReference)
		}
	}
	return references
}

// cmsDocumentAnalyzerContainsSignatureScope reports whether scopes already holds an equal
// SignatureScope, standing in for Java's List#contains (SignatureScope#equals).
func cmsDocumentAnalyzerContainsSignatureScope(scopes []mscope.SignatureScope, target mscope.SignatureScope) bool {
	for _, s := range scopes {
		if s.Equals(target) {
			return true
		}
	}
	return false
}

// cmsDocumentAnalyzerContainsTimestampedReference reports whether references already holds an
// equal TimestampedReference, standing in for Java's List#contains
// (TimestampedReference#equals).
func cmsDocumentAnalyzerContainsTimestampedReference(references []*validation.TimestampedReference, target *validation.TimestampedReference) bool {
	for _, r := range references {
		if r.Equals(target) {
			return true
		}
	}
	return false
}

// CMS returns a CMS. Port of getCMS().
func (a *CMSDocumentAnalyzer) CMS() *cms.CMS {
	return a.cms
}

// OriginalDocumentsForSignature returns the signed document(s) without their signature(s). Port
// of the getOriginalDocuments(AdvancedSignature) override.
//
// The DSSException Java catches (logging "Cannot retrieve a list of original documents") is
// swallowed the same way here, since slf4j logging is dropped per PORTING.md.
func (a *CMSDocumentAnalyzer) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	cadesSignature := advancedSignature.(*CAdESSignature)
	document, err := cadesSignature.OriginalDocument()
	if err != nil {
		return []model.DSSDocument{}
	}
	return []model.DSSDocument{document}
}
