// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/EvidenceRecordProductionComparator.java (DSS 6.5.RC1).
//
// See time_stamp_token_production_comparator.go's header - UnsignedAttributes needs this
// sibling comparator too.
package cades

import "github.com/ryftcore/dss-go/dss/internal/asn1ber"

// EvidenceRecordProductionComparator compares production time of RFC 4998 EvidenceRecords,
// checking their generation time. Port of the class EvidenceRecordProductionComparator,
// implementing Comparator<org.bouncycastle.asn1.tsp.EvidenceRecord>.
//
// The Java type carries the raw ASN.1 EvidenceRecord; this port has no Go type for it yet (see
// cades_attribute.go), so Compare takes the DER encoding directly - what
// Attribute.ToEvidenceRecord's *asn1ber.Element.Encoded() hands out.
//
// Compare returns:
//   - -1 if evidenceRecordOne was created before evidenceRecordTwo
//   - 0 if the EvidenceRecords were created at the same time
//   - 1 if evidenceRecordOne was created after evidenceRecordTwo
type EvidenceRecordProductionComparator struct{}

// NewEvidenceRecordProductionComparator is the port of the default constructor.
func NewEvidenceRecordProductionComparator() EvidenceRecordProductionComparator {
	return EvidenceRecordProductionComparator{}
}

// Compare is the port of compare(EvidenceRecord, EvidenceRecord).
func (c EvidenceRecordProductionComparator) Compare(evidenceRecordOne, evidenceRecordTwo *asn1ber.Element) int {
	return c.compareByGenerationTime(evidenceRecordOne, evidenceRecordTwo)
}

func (c EvidenceRecordProductionComparator) compareByGenerationTime(evidenceRecordOne, evidenceRecordTwo *asn1ber.Element) int {
	t1, err1 := UtilsEvidenceRecordGenerationTime(evidenceRecordOne.Encoded())
	t2, err2 := UtilsEvidenceRecordGenerationTime(evidenceRecordTwo.Encoded())
	if err1 != nil || err2 != nil {
		// Upstream's getEvidenceRecordGenerationTime raises a DSSException that would propagate
		// out of compare(); there is no such channel through a Comparator here, so an unreadable
		// evidence record's generation time compares as equal (0), leaving the pair's order
		// unaffected the way the mergesort in UnsignedAttributes.sortAttributes only ever
		// swaps on a strictly positive result.
		return 0
	}
	switch {
	case t1.Before(t2):
		return -1
	case t1.After(t2):
		return 1
	default:
		return 0
	}
}
