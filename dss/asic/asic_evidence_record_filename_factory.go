// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ASiCEvidenceRecordFilenameFactory.java
// (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/enumerations"

// EvidenceRecordFilenameFactory creates a new evidence record's filename for the current
// container type and Content.
type EvidenceRecordFilenameFactory interface {
	// EvidenceRecordFilename returns a filename for an evidence record file to be created.
	// Port of getEvidenceRecordFilename(ASiCContent, EvidenceRecordTypeEnum).
	EvidenceRecordFilename(asicContent *Content, evidenceRecordType enumerations.EvidenceRecordTypeEnum) string

	// EvidenceRecordManifestFilename returns a filename for an evidence record's ASIC manifest
	// file to be created. Port of getEvidenceRecordManifestFilename(ASiCContent).
	EvidenceRecordManifestFilename(asicContent *Content) string
}
