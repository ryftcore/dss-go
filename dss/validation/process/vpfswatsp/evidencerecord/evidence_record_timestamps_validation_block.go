// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/evidencerecord/EvidenceRecordTimestampsValidationBlock.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
//
// The base class lives in package vpftsp
// (eu.europa.esig.dss.validation.process.vpftsp.TimestampsValidationBlock). It
// is subclassed here through the InitTimestampsValidationBlockStateWithoutPOE
// / InitTimestampsValidationBlock pair (the state-then-register shape this
// port uses everywhere a Java class is designed for extension); the overrides
// interface routes getTimestamps() and getPoe() back to this type.
//
// Java's override of getTimestamps() reads the base's protected `timestamps`
// field, which holds exactly what this class's constructor passed up -
// evidenceRecord.getTimestampList(). The wrapper is kept here and re-read
// instead, so that the override does not depend on the base's field.
package evidencerecord

import (
	"sort"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfswatsp"
	"github.com/ryftcore/dss-go/dss/validation/process/vpftsp"
)

// TimestampsValidationBlock verifies a time-stamp of an Evidence
// Record.
type TimestampsValidationBlock struct {
	*vpftsp.TimestampsValidationBlock

	// evidenceRecord is the evidence record whose time-stamps are validated;
	// see the file header.
	evidenceRecord *diagnostic.EvidenceRecordWrapper
}

// NewEvidenceRecordTimestampsValidationBlock is the default constructor. Port
// of TimestampsValidationBlock(Provider, EvidenceRecordWrapper, Data, ValidationPolicy, Date, Map, List, ValidationLevel).
func NewEvidenceRecordTimestampsValidationBlock(i18nProvider *i18n.Provider,
	evidenceRecord *diagnostic.EvidenceRecordWrapper, diagnosticData *diagnostic.Data,
	validationPolicy policy.ValidationPolicy, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tlAnalysis []*jaxb.XmlTLAnalysis,
	validationLevel enumerations.ValidationLevel) *TimestampsValidationBlock {
	b := &TimestampsValidationBlock{
		TimestampsValidationBlock: &vpftsp.TimestampsValidationBlock{},
		evidenceRecord:            evidenceRecord,
	}
	b.InitTimestampsValidationBlockStateWithoutPOE(i18nProvider, evidenceRecord.TimestampList(), diagnosticData,
		validationPolicy, currentTime, bbbs, tlAnalysis, validationLevel)
	b.InitTimestampsValidationBlock(b)
	return b
}

// Timestamps returns a list of time-stamp tokens to be validated. Port of the
// overridden getTimestamps(): evidence record time-stamps are validated in the
// order of their appearance.
//
// Java sorts a copy with Comparator.comparing(TimestampWrapper::getProductionTime),
// which is a stable sort and raises a NullPointerException on a time-stamp
// without a production time; sort.SliceStable is the same stable sort, and a
// missing production time dereferences here in its place.
func (b *TimestampsValidationBlock) Timestamps() []*diagnostic.TimestampWrapper {
	timestamps := b.evidenceRecord.TimestampList()
	timestampList := make([]*diagnostic.TimestampWrapper, len(timestamps))
	copy(timestampList, timestamps)
	sort.SliceStable(timestampList, func(i, j int) bool {
		return timestampList[i].ProductionTime().Before(*timestampList[j].ProductionTime())
	})
	return timestampList
}

// Poe returns the POE container to be used for the given timestamp. Port of the
// overridden getPoe(TimestampWrapper).
func (b *TimestampsValidationBlock) Poe(timestamp *diagnostic.TimestampWrapper) *vpfswatsp.POEExtraction {
	poe := b.TimestampsValidationBlock.Poe(timestamp)
	/*
	 * i) Before validating a time-stamp the process shall extract POEs (as per clause 5.6.2.3) of the
	 * time-stamp within the next Archive timestamp and initialize the set of temporary POEs with the
	 * extracted POEs.
	 */
	nextTimestamp := b.nextTimestamp(timestamp)
	if nextTimestamp != nil {
		// skip message-imprint check
		poe.ExtractTimestampedObjectsPOE(nextTimestamp.TimestampedObjects(), nextTimestamp.ProductionTime())
	}
	return poe
}

// nextTimestamp ports the private getNextTimestamp(TimestampWrapper): the
// iterator walk returns the entry following the first one whose id matches,
// or null when the match is the last entry (or there is none).
func (b *TimestampsValidationBlock) nextTimestamp(
	currentTimestamp *diagnostic.TimestampWrapper) *diagnostic.TimestampWrapper {
	evidenceRecordTimestamps := b.Timestamps()
	for i, timestampWrapper := range evidenceRecordTimestamps {
		if currentTimestamp.Id() == timestampWrapper.Id() && i+1 < len(evidenceRecordTimestamps) {
			return evidenceRecordTimestamps[i+1]
		}
	}
	return nil
}
