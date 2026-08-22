// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/POEExtraction.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening and POE-hierarchy notes.
//
// Go has no method overloading, so the three extractPOE overloads are named
// after what they extract from: ExtractPOE(TimestampWrapper),
// ExtractEvidenceRecordPOE(EvidenceRecordWrapper) and
// ExtractTimestampedObjectsPOE(List, Date).
//
// The poeMap is a java.util.HashMap upstream, but it is only ever read by key
// (get(tokenId)) - never iterated - so its iteration order cannot reach any
// output and a plain Go map is a faithful stand-in.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// POEExtraction performs 5.6.2.3 POE extraction.
//
// 5.6.2.3.1 Description
//
// This building block derives POEs from a given time-stamp. Assumptions:
//   - The time-stamp validation has returned PASSED.
//   - The cryptographic hash function used in the time-stamp
//     (messageImprint.hashAlgorithm) is considered reliable at current time or,
//     if this is not the case, a PoE for that time-stamp exists for a time when
//     the hash function has still been considered reliable.
//
// In the simple case, a time-stamp gives a POE for each data item protected by
// the time-stamp at the generation date/time of the token.
//
// EXAMPLE: A time-stamp on the signature value gives a POE of the signature
// value at the generation date/time of the time-stamp.
//
// A time-stamp can also give an indirect POE when it is computed on the hash
// value of some data instead of the data itself. A POE for DATA at T1 can be
// derived from the time-stamp:
//   - If there is a POE for h(DATA) at a date T1, where h is a cryptographic
//     hash function and DATA is some data (e.g. a certificate),
//   - if h is asserted in the cryptographic constraints to be trusted until at
//     least a date T after T1; and
//   - if there is a POE for DATA at a date T after T1.
type POEExtraction struct {
	// poeMap is the map of proofs of existence by token ids.
	poeMap map[string][]POE
}

// NewPOEExtraction is the default constructor initializing an empty map.
func NewPOEExtraction() *POEExtraction {
	return &POEExtraction{poeMap: make(map[string][]POE)}
}

// Init instantiates a controlTime POE for all used tokens. Port of
// init(Data, Date).
func (p *POEExtraction) Init(diagnosticData *diagnostic.Data, controlTime time.Time) {

	controlTimePoe := NewPOE(controlTime)

	for _, signature := range diagnosticData.AllSignatures() {
		p.addPOE(signature.Id(), controlTimePoe)
	}
	for _, timestamp := range diagnosticData.TimestampList() {
		p.addPOE(timestamp.Id(), controlTimePoe)
	}
	for _, evidenceRecord := range diagnosticData.EvidenceRecords() {
		p.addPOE(evidenceRecord.Id(), controlTimePoe)
	}
	for _, eaa := range diagnosticData.EAAs() {
		p.addPOE(eaa.Id(), controlTimePoe)
	}
	for _, certificate := range diagnosticData.UsedCertificates() {
		p.addPOE(certificate.Id(), controlTimePoe)
	}
	for _, revocation := range diagnosticData.AllRevocationData() {
		p.addPOE(revocation.Id(), controlTimePoe)
	}
	for _, signerData := range diagnosticData.AllSignerDocuments() {
		p.addPOE(signerData.Id(), controlTimePoe)
	}
	for _, orphanCertificate := range diagnosticData.AllOrphanCertificateObjects() {
		p.addPOE(orphanCertificate.Id(), controlTimePoe)
	}
	for _, orphanCertificateRef := range diagnosticData.AllOrphanCertificateReferences() {
		p.addPOE(orphanCertificateRef.Id(), controlTimePoe)
	}
	for _, orphanRevocation := range diagnosticData.AllOrphanRevocationObjects() {
		p.addPOE(orphanRevocation.Id(), controlTimePoe)
	}
	for _, orphanRevocationRef := range diagnosticData.AllOrphanRevocationReferences() {
		p.addPOE(orphanRevocationRef.Id(), controlTimePoe)
	}

}

// CollectAllPOE extracts all POEs from the provided collection of timestamps.
// Port of collectAllPOE(Collection).
func (p *POEExtraction) CollectAllPOE(timestamps []*diagnostic.TimestampWrapper) {
	for _, timestamp := range timestamps {
		p.ExtractPOE(timestamp)
	}
}

// ExtractPOE extracts POE for all covered objects from a single timestamp
// wrapper. Port of extractPOE(TimestampWrapper).
func (p *POEExtraction) ExtractPOE(timestamp *diagnostic.TimestampWrapper) {
	/*
	 * 5.6.2.3.4 Processing (5.6.2.3 POE extraction)
	 *
	 * 1) The building block shall determine the set S of references to objects and
	 * objects that are part of the signature and are protected by the time-stamp.
	 */
	if timestamp.IsMessageImprintDataFound() && timestamp.IsMessageImprintDataIntact() {
		timestampedObjects := timestamp.TimestampedObjects()
		if utils.IsCollectionNotEmpty(timestampedObjects) {
			poe := NewTimestampPOE(timestamp)
			for _, xmlTimestampedObject := range timestampedObjects {
				p.addPOE(timestampedObjectId(xmlTimestampedObject), poe)
			}
		}
	}
}

// ExtractEvidenceRecordPOE extracts POE for all objects covered by an evidence
// record. Port of extractPOE(EvidenceRecordWrapper).
func (p *POEExtraction) ExtractEvidenceRecordPOE(evidenceRecord *diagnostic.EvidenceRecordWrapper) {
	coveredObjects := evidenceRecord.CoveredObjects()
	if utils.IsCollectionNotEmpty(coveredObjects) {
		poe := NewEvidenceRecordPOE(evidenceRecord)
		for _, xmlTimestampedObject := range coveredObjects {
			p.addPOE(timestampedObjectId(xmlTimestampedObject), poe)
		}
	}
}

// ExtractTimestampedObjectsPOE extracts POE for given timestamped objects at
// the provided poeTime. Port of extractPOE(List, Date): the nullable Java Date
// keeps its nullability, the method being a no-op for a null one.
func (p *POEExtraction) ExtractTimestampedObjectsPOE(
	timestampedObjects []*diagnosticjaxb.XmlTimestampedObject, poeTime *time.Time) {
	if utils.IsCollectionNotEmpty(timestampedObjects) && poeTime != nil {
		poe := NewPOE(*poeTime)
		for _, xmlTimestampedObject := range timestampedObjects {
			p.addPOE(timestampedObjectId(xmlTimestampedObject), poe)
		}
	}
}

// addPOE ports the private addPOE(String, POE). Java's null guard on the POE
// becomes a nil-interface guard.
func (p *POEExtraction) addPOE(tokenId string, proofOfExistence POE) {
	if proofOfExistence != nil {
		p.poeMap[tokenId] = append(p.poeMap[tokenId], proofOfExistence)
	}
}

// AddSignaturePOE adds a specific POE for a signature wrapper. Port of
// addSignaturePOE(SignatureWrapper, POE).
func (p *POEExtraction) AddSignaturePOE(signature *diagnostic.SignatureWrapper, proofOfExistence POE) {
	if signature != nil {
		p.addPOE(signature.Id(), proofOfExistence)
	}
}

// IsPOEExists returns true if there is a POE exists for a given id at (or
// before) the control time. Port of isPOEExists(String, Date).
func (p *POEExtraction) IsPOEExists(tokenId string, controlTime time.Time) bool {
	poes, ok := p.poeMap[tokenId]
	if ok {
		for _, poe := range poes {
			if poe.Time().Compare(controlTime) <= 0 {
				return true
			}
		}
	}
	return false
}

// IsPOEExistInRange checks if a POE exists for the token with the given Id
// within the validity range between notBefore and notAfter inclusively. Port of
// isPOEExistInRange(String, Date, Date).
//
// The two bounds stay nullable, the callers passing a certificate's notBefore /
// notAfter straight through; a nil bound dereferences here exactly where Java's
// Date#compareTo raises a NullPointerException.
func (p *POEExtraction) IsPOEExistInRange(tokenId string, notBefore, notAfter *time.Time) bool {
	poes, ok := p.poeMap[tokenId]
	if ok {
		for _, poe := range poes {
			if poe.Time().Compare(*notBefore) >= 0 && poe.Time().Compare(*notAfter) <= 0 {
				return true
			}
		}
	}
	return false
}

// GetLowestPOETime returns the lowest POE time for the requested token. Port of
// getLowestPOETime(String).
//
// Java dereferences getLowestPOE(tokenId) unguarded, so a token with no POE at
// all (init(controlTime) not executed beforehand) raises a
// NullPointerException; the nil POE interface panics here in its place.
func (p *POEExtraction) GetLowestPOETime(tokenId string) time.Time {
	return p.GetLowestPOE(tokenId).Time()
}

// GetLowestPOE returns the lowest POE for the requested token. Port of
// getLowestPOE(String).
//
// NOTE: can return NULL (nil) if POE is not found (Init must be executed
// before).
func (p *POEExtraction) GetLowestPOE(tokenId string) POE {
	var lowestPOE POE
	poes, ok := p.poeMap[tokenId]
	if ok {
		comparator := NewPOEComparator()
		for _, poe := range poes {
			if lowestPOE == nil || comparator.Before(poe, lowestPOE) {
				lowestPOE = poe
			}
		}
	}
	return lowestPOE
}
