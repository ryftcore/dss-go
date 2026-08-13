// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/EvidenceRecordTimestampScopeFinder.java (DSS 6.5.RC1).
package scope

import (
	mscope "github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/spi/validation"
)

// EvidenceRecordTimestampScopeFinder finds timestamped scopes for evidence record
// time-stamps.
type EvidenceRecordTimestampScopeFinder struct {
	EvidenceRecordScopeFinder
}

// NewEvidenceRecordTimestampScopeFinder is the default constructor. Port of
// EvidenceRecordTimestampScopeFinder(EvidenceRecord).
func NewEvidenceRecordTimestampScopeFinder(evidenceRecord validation.EvidenceRecord) *EvidenceRecordTimestampScopeFinder {
	return &EvidenceRecordTimestampScopeFinder{EvidenceRecordScopeFinder: *NewEvidenceRecordScopeFinder(evidenceRecord)}
}

// FindTimestampScope returns a timestamp scope for the given TimestampToken. Port of
// findTimestampScope(TimestampToken).
func (f *EvidenceRecordTimestampScopeFinder) FindTimestampScope(timestampToken *validation.TimestampToken) []mscope.SignatureScope {
	if timestampToken.IsMessageImprintDataIntact() && (!f.EvidenceRecord.IsEmbedded() || timestampToken.AreReferenceValidationsValid()) {
		return f.EvidenceRecord.EvidenceRecordScopes()
	}
	return []mscope.SignatureScope{}
}

// compile-time interface assertion.
var _ TimestampScopeFinder = (*EvidenceRecordTimestampScopeFinder)(nil)
