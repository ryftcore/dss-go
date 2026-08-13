// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/evidencerecord/EvidenceRecordAnalyzer.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.validation.analyzer.evidencerecord lands in this same Go package
// (dss/spi/validation/analyzer) per S2B_BRIEF.md's package layout table.
package analyzer

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/spi/validation"
)

// EvidenceRecordAnalyzer performs validation of an evidence record document.
type EvidenceRecordAnalyzer interface {
	DocumentAnalyzer

	// EvidenceRecord returns a single EvidenceRecord to be validated. Port of
	// getEvidenceRecord().
	EvidenceRecord() validation.EvidenceRecord

	// EvidenceRecordType returns the type of the evidence record supported by the current
	// validator. Port of getEvidenceRecordType().
	EvidenceRecordType() enumerations.EvidenceRecordTypeEnum

	// SetEvidenceRecordOrigin sets the origin of the extracted evidence record.
	// Default: EvidenceRecordOrigin.EXTERNAL. Port of setEvidenceRecordOrigin(EvidenceRecordOrigin).
	SetEvidenceRecordOrigin(origin enumerations.EvidenceRecordOrigin)

	// SetEvidenceRecordIncorporationType sets the incorporation type of the evidence record
	// within a signature's unsigned attributes.
	// NOTE: only used for attached CAdES evidence records. Port of
	// setEvidenceRecordIncorporationType(EvidenceRecordIncorporationType).
	SetEvidenceRecordIncorporationType(evidenceRecordIncorporationType enumerations.EvidenceRecordIncorporationType)

	// SetEmbeddedEvidenceRecordHelper sets a helper for processing and validation of the
	// embedded evidence record type. Port of
	// setEmbeddedEvidenceRecordHelper(EmbeddedEvidenceRecordHelper).
	SetEmbeddedEvidenceRecordHelper(embeddedEvidenceRecordHelper validation.EmbeddedEvidenceRecordHelper)
}
