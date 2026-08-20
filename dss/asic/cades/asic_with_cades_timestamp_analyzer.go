// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/timestamp/ASiCWithCAdESTimestampAnalyzer.java (DSS 6.5.RC1).
//
// The Java `validation.timestamp` sub-package flattens into this Go package per the phase-7
// package layout (S7_BRIEF.md).
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/spi/validation"
	dsstimestamp "github.com/utain/esig/dss/validation/timestamp"
)

// ASiCWithCAdESTimestampAnalyzer is the abstract validator for an ASiC with CAdES timestamp.
type ASiCWithCAdESTimestampAnalyzer struct {
	dsstimestamp.DetachedTimestampAnalyzer

	// originalDocuments is a list of original documents present in the container.
	originalDocuments []model.DSSDocument

	// archiveDocuments is a list of package.zip embedded documents, when applicable.
	archiveDocuments []model.DSSDocument

	// archiveTimestampType defines the archive timestamp type.
	archiveTimestampType    enumerations.ArchiveTimestampType
	archiveTimestampTypeSet bool
}

var _ dsstimestamp.DetachedTimestampAnalyzerOverrides = (*ASiCWithCAdESTimestampAnalyzer)(nil)

// newASiCWithCAdESTimestampAnalyzer wires the overrides registration shared by both
// constructors.
func newASiCWithCAdESTimestampAnalyzer() *ASiCWithCAdESTimestampAnalyzer {
	a := &ASiCWithCAdESTimestampAnalyzer{
		DetachedTimestampAnalyzer: dsstimestamp.NewDetachedTimestampAnalyzerBase(),
	}
	a.InitDetachedTimestampAnalyzer(a)
	return a
}

// NewASiCWithCAdESTimestampAnalyzer is the default constructor. Ports
// ASiCWithCAdESTimestampAnalyzer(DSSDocument), which delegates to the single-argument
// DetachedTimestampAnalyzer(DSSDocument) constructor and so keeps its
// TimestampType_CONTENT_TIMESTAMP default.
func NewASiCWithCAdESTimestampAnalyzer(timestamp model.DSSDocument) *ASiCWithCAdESTimestampAnalyzer {
	return NewASiCWithCAdESTimestampAnalyzerWithType(timestamp, enumerations.TimestampType_CONTENT_TIMESTAMP)
}

// NewASiCWithCAdESTimestampAnalyzerWithType is the default constructor with a timestamp type.
// Ports ASiCWithCAdESTimestampAnalyzer(DSSDocument, TimestampType).
func NewASiCWithCAdESTimestampAnalyzerWithType(timestamp model.DSSDocument, tstType enumerations.TimestampType) *ASiCWithCAdESTimestampAnalyzer {
	a := newASiCWithCAdESTimestampAnalyzer()
	a.SetDocument(timestamp)
	a.SetTimestampType(tstType)
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
func (a *ASiCWithCAdESTimestampAnalyzer) CreateTimestampToken() (*validation.TimestampToken, error) {
	timestamp, err := a.DetachedTimestampAnalyzer.CreateTimestampToken()
	if err != nil {
		return nil, err
	}
	if a.ManifestFile() != nil {
		timestamp.SetManifestFile(a.ManifestFile())
	}
	if a.archiveTimestampTypeSet {
		timestamp.SetArchiveTimestampType(a.archiveTimestampType)
	}
	return timestamp, nil
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
