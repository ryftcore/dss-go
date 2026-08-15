//go:build phase8

// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/timestamp/ASiCWithCAdESTimestampAnalyzer.java (DSS 6.5.RC1).
//
// The Java `validation.timestamp` sub-package flattens into this Go package per the phase-7
// package layout (S7_BRIEF.md).
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation does not exist yet. Drop the tag once Phase 8 lands the package; the file may
// need adjusting to the real shape (see the BLOCKED FORWARD DEPENDENCY note below).
//
// BLOCKED FORWARD DEPENDENCY (flagged per S7_BRIEF.md's "flag needs in notes" rule): this
// class's Java base, eu.europa.esig.dss.validation.timestamp.DetachedTimestampAnalyzer, belongs
// to dss-validation, assigned to the not-yet-ported `validation` package (Phase 8). Following
// the Overrides+Init virtual-dispatch convention this porting effort uses pervasively (see
// AbstractASiCContainerAnalyzer, DefaultContainerMerger), this file assumes:
//
//	type DetachedTimestampAnalyzerOverrides interface {
//	    CreateTimestampToken() *validation.TimestampToken
//	    IsTimestampCoveredByEvidenceRecord(timestampToken *validation.TimestampToken, evidenceRecord validation.EvidenceRecord) bool
//	    GetTimestampScopes(timestampToken *validation.TimestampToken) []modelscope.SignatureScope
//	    AddReference(signatureScope modelscope.SignatureScope) bool
//	}
//	type DetachedTimestampAnalyzer struct { ... }
//	func NewDetachedTimestampAnalyzerBase() DetachedTimestampAnalyzer
//	func (a *DetachedTimestampAnalyzer) InitDetachedTimestampAnalyzer(overrides DetachedTimestampAnalyzerOverrides)
//	func (a *DetachedTimestampAnalyzer) InitFromDocument(timestamp model.DSSDocument)
//	func (a *DetachedTimestampAnalyzer) InitFromDocumentWithType(timestamp model.DSSDocument, timestampType enumerations.TimestampType)
//	func (a *DetachedTimestampAnalyzer) TimestampedData() model.DSSDocument
//	func (a *DetachedTimestampAnalyzer) SetTimestampedData(timestampedData model.DSSDocument)
//	func (a *DetachedTimestampAnalyzer) ManifestFile() *model.ManifestFile // protected field manifestFile, exported for cross-package embedding
//	func (a *DetachedTimestampAnalyzer) SetManifestFile(manifestFile *model.ManifestFile)
//	func (a *DetachedTimestampAnalyzer) SetDocument(document model.DSSDocument)
//	func (a *DetachedTimestampAnalyzer) SetCertificateVerifier(certificateVerifier validation.CertificateVerifier)
//	func (a *DetachedTimestampAnalyzer) SetDetachedEvidenceRecords(evidenceRecords []validation.EvidenceRecord)
//	func (a *DetachedTimestampAnalyzer) CreateTimestampToken() *validation.TimestampToken
//
// Revisit once Phase 8 lands the real package - only the imports/shape above need to resolve.
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/spi/validation"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// ASiCWithCAdESTimestampAnalyzer is the abstract validator for an ASiC with CAdES timestamp.
type ASiCWithCAdESTimestampAnalyzer struct {
	dssvalidation.DetachedTimestampAnalyzer

	// originalDocuments is a list of original documents present in the container.
	originalDocuments []model.DSSDocument

	// archiveDocuments is a list of package.zip embedded documents, when applicable.
	archiveDocuments []model.DSSDocument

	// archiveTimestampType defines the archive timestamp type.
	archiveTimestampType    enumerations.ArchiveTimestampType
	archiveTimestampTypeSet bool
}

var _ dssvalidation.DetachedTimestampAnalyzerOverrides = (*ASiCWithCAdESTimestampAnalyzer)(nil)

// newASiCWithCAdESTimestampAnalyzer wires the overrides registration shared by both
// constructors.
func newASiCWithCAdESTimestampAnalyzer() *ASiCWithCAdESTimestampAnalyzer {
	a := &ASiCWithCAdESTimestampAnalyzer{
		DetachedTimestampAnalyzer: dssvalidation.NewDetachedTimestampAnalyzerBase(),
	}
	a.InitDetachedTimestampAnalyzer(a)
	return a
}

// NewASiCWithCAdESTimestampAnalyzer is the default constructor. Ports
// ASiCWithCAdESTimestampAnalyzer(DSSDocument).
func NewASiCWithCAdESTimestampAnalyzer(timestamp model.DSSDocument) *ASiCWithCAdESTimestampAnalyzer {
	a := newASiCWithCAdESTimestampAnalyzer()
	a.InitFromDocument(timestamp)
	return a
}

// NewASiCWithCAdESTimestampAnalyzerWithType is the default constructor with a timestamp type.
// Ports ASiCWithCAdESTimestampAnalyzer(DSSDocument, TimestampType).
func NewASiCWithCAdESTimestampAnalyzerWithType(timestamp model.DSSDocument, tstType enumerations.TimestampType) *ASiCWithCAdESTimestampAnalyzer {
	a := newASiCWithCAdESTimestampAnalyzer()
	a.InitFromDocumentWithType(timestamp, tstType)
	return a
}

// GetCoveredManifest returns the covered ManifestFile. Ports getCoveredManifest().
func (a *ASiCWithCAdESTimestampAnalyzer) GetCoveredManifest() *model.ManifestFile {
	return a.ManifestFile()
}

// SetOriginalDocuments sets the original documents present in the ASiC container. Ports
// setOriginalDocuments(List).
func (a *ASiCWithCAdESTimestampAnalyzer) SetOriginalDocuments(originalDocuments []model.DSSDocument) {
	a.originalDocuments = originalDocuments
}

// SetArchiveDocuments sets the document embedded inside package.zip, when applicable. Ports
// setArchiveDocuments(List).
func (a *ASiCWithCAdESTimestampAnalyzer) SetArchiveDocuments(archiveDocuments []model.DSSDocument) {
	a.archiveDocuments = archiveDocuments
}

// SetArchiveTimestampType sets the archive timestamp type. Ports
// setArchiveTimestampType(ArchiveTimestampType).
func (a *ASiCWithCAdESTimestampAnalyzer) SetArchiveTimestampType(archiveTimestampType enumerations.ArchiveTimestampType) {
	a.archiveTimestampType = archiveTimestampType
	a.archiveTimestampTypeSet = true
}

// CreateTimestampToken ports the @Override protected createTimestampToken().
func (a *ASiCWithCAdESTimestampAnalyzer) CreateTimestampToken() *validation.TimestampToken {
	timestamp := a.DetachedTimestampAnalyzer.CreateTimestampToken()
	if a.ManifestFile() != nil {
		timestamp.SetManifestFile(a.ManifestFile())
	}
	if a.archiveTimestampTypeSet {
		timestamp.SetArchiveTimestampType(a.archiveTimestampType)
	}
	return timestamp
}

// IsTimestampCoveredByEvidenceRecord ports the @Override protected
// isTimestampCoveredByEvidenceRecord(TimestampToken, EvidenceRecord).
func (a *ASiCWithCAdESTimestampAnalyzer) IsTimestampCoveredByEvidenceRecord(timestampToken *validation.TimestampToken, evidenceRecord validation.EvidenceRecord) bool {
	erManifestFile := evidenceRecord.ManifestFile()
	if erManifestFile == nil {
		// detached ER, covers all content
		return true
	}
	for _, entry := range erManifestFile.Entries() {
		if timestampToken.Filename() != "" && timestampToken.Filename() == entry.Uri() {
			return true
		}
	}
	return false
}

// GetTimestampScopes ports the @Override protected getTimestampScopes(TimestampToken).
func (a *ASiCWithCAdESTimestampAnalyzer) GetTimestampScopes(timestampToken *validation.TimestampToken) []mscope.SignatureScope {
	timestampScopeFinder := NewASiCWithCAdESTimestampScopeFinder()
	timestampScopeFinder.SetContainerDocuments(a.originalDocuments)
	timestampScopeFinder.SetArchiveDocuments(a.archiveDocuments)
	timestampScopeFinder.SetTimestampedData(a.TimestampedData())
	return timestampScopeFinder.FindTimestampScope(timestampToken)
}

// AddReference ports the @Override protected addReference(SignatureScope).
//
// Cross-chunk assumption (ZIPCORE): ASiCUtilsIsSignature/ASiCUtilsIsTimestamp/
// ASiCUtilsIsEvidenceRecord take a filename string.
func (a *ASiCWithCAdESTimestampAnalyzer) AddReference(signatureScope mscope.SignatureScope) bool {
	fileName := signatureScope.DocumentName()
	return fileName == "" || (!asic.ASiCUtilsIsSignature(fileName) && !asic.ASiCUtilsIsTimestamp(fileName) && !asic.ASiCUtilsIsEvidenceRecord(fileName))
}
