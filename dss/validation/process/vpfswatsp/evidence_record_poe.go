// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/EvidenceRecordPOE.java (DSS 6.5.RC1).
//
// See poe.go for the POE hierarchy note.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
)

// EvidenceRecordPOE is a POE provided by an evidence record.
type EvidenceRecordPOE struct {
	*POEBase

	// evidenceRecord is the evidence record.
	evidenceRecord *diagnostic.EvidenceRecordWrapper
}

// NewEvidenceRecordPOE is the constructor to instantiate POE by an evidence
// record. Port of EvidenceRecordPOE(EvidenceRecordWrapper).
func NewEvidenceRecordPOE(evidenceRecord *diagnostic.EvidenceRecordWrapper) *EvidenceRecordPOE {
	return &EvidenceRecordPOE{
		POEBase:        NewPOE(evidenceRecordPOETime(evidenceRecord)),
		evidenceRecord: evidenceRecord,
	}
}

// evidenceRecordPOETime ports the private static getPOETime(EvidenceRecordWrapper).
// The two requireNonNull messages are upstream's, typo included; the production
// time of the first time-stamp being null trips POE(Date)'s own requireNonNull
// in Java, raised here with that constructor's message (see NewPOE).
func evidenceRecordPOETime(evidenceRecord *diagnostic.EvidenceRecordWrapper) time.Time {
	if evidenceRecord == nil {
		panic("The evidenceRecord must be defined!")
	}
	firstTimestamp := evidenceRecord.FirstTimestamp()
	if firstTimestamp == nil {
		panic("EvidenceRecord shall have at leats one time-stamp!")
	}
	productionTime := firstTimestamp.ProductionTime()
	if productionTime == nil {
		panic("The controlTime must be defined!")
	}
	return *productionTime
}

// POEProviderId returns the evidence record's Id. Port of the overridden
// getPOEProviderId().
func (p *EvidenceRecordPOE) POEProviderId() *string {
	id := p.evidenceRecord.Id()
	return &id
}

// POEObjects returns the objects covered by the evidence record. Port of the
// overridden getPOEObjects().
func (p *EvidenceRecordPOE) POEObjects() []*diagnosticjaxb.XmlTimestampedObject {
	return p.evidenceRecord.CoveredObjects()
}

// IsTokenProvided returns whether the POE is provided by a token. Port of the
// overridden isTokenProvided().
func (p *EvidenceRecordPOE) IsTokenProvided() bool {
	return true
}
