// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/TimeStampTokenProductionComparator.java (DSS 6.5.RC1).
//
// This class lives in eu.europa.esig.dss.cades, not .cades.validation, but it is ported into
// this same Go package because UnsignedAttributes needs it to sort timestamp/evidence-record
// unsigned attributes by production time, and cades_level_baseline_lt.go calls
// NewTimeStampTokenProductionComparator().
package cades

import (
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
)

// TimeStampTokenProductionComparator compares production time of TimeStampTokens, checking
// their production time and covered data. Port of the class TimeStampTokenProductionComparator,
// implementing Comparator<TimeStampToken>.
//
// Compare returns:
//   - -1 if timeStampTokenOne was created before timeStampTokenTwo
//   - 0 if the TimeStampTokens were created at the same time
//   - 1 if timeStampTokenOne was created after timeStampTokenTwo
type TimeStampTokenProductionComparator struct{}

// NewTimeStampTokenProductionComparator is the port of the default constructor.
func NewTimeStampTokenProductionComparator() TimeStampTokenProductionComparator {
	return TimeStampTokenProductionComparator{}
}

// Compare is the port of compare(TimeStampToken, TimeStampToken).
func (c TimeStampTokenProductionComparator) Compare(timeStampTokenOne, timeStampTokenTwo *cmscore.TimeStampToken) int {
	result := c.compareByGenerationTime(timeStampTokenOne, timeStampTokenTwo)
	if result == 0 {
		result = c.compareByHashTableSize(timeStampTokenOne, timeStampTokenTwo)
	}
	return result
}

// After reports whether timeStampTokenOne was created after timeStampTokenTwo. Port of
// after(TimeStampToken, TimeStampToken).
func (c TimeStampTokenProductionComparator) After(timeStampTokenOne, timeStampTokenTwo *cmscore.TimeStampToken) bool {
	return c.Compare(timeStampTokenOne, timeStampTokenTwo) > 0
}

func (c TimeStampTokenProductionComparator) compareByGenerationTime(tst1, tst2 *cmscore.TimeStampToken) int {
	t1 := spi.DSSASN1UtilsTimeStampTokenGenerationTime(tst1)
	t2 := spi.DSSASN1UtilsTimeStampTokenGenerationTime(tst2)
	switch {
	case t1.Before(t2):
		return -1
	case t1.After(t2):
		return 1
	default:
		return 0
	}
}

func (c TimeStampTokenProductionComparator) compareByHashTableSize(tst1, tst2 *cmscore.TimeStampToken) int {
	atsHashIndexOne := UtilsAtsHashIndex(timeStampTokenProductionComparatorUnsignedAttributes(tst1))
	atsHashIndexTwo := UtilsAtsHashIndex(timeStampTokenProductionComparatorUnsignedAttributes(tst2))

	if atsHashIndexOne != nil && atsHashIndexTwo != nil {
		hashTableSizeOne := timeStampTokenProductionComparatorHashTableSize(atsHashIndexOne)
		hashTableSizeTwo := timeStampTokenProductionComparatorHashTableSize(atsHashIndexTwo)
		if hashTableSizeOne < hashTableSizeTwo {
			return -1
		} else if hashTableSizeOne > hashTableSizeTwo {
			return 1
		}
	} else if atsHashIndexOne != nil {
		return 1
	} else if atsHashIndexTwo != nil {
		return -1
	}
	return 0
}

// timeStampTokenProductionComparatorUnsignedAttributes returns the unsignedAttrs of a
// TimeStampToken's single SignerInfo, nil for a token with none, mirroring
// TimeStampToken#getUnsignedAttributes() (BouncyCastle: SignerInformation#getUnsignedAttributes()).
func timeStampTokenProductionComparatorUnsignedAttributes(tst *cmscore.TimeStampToken) cmscore.Attributes {
	signers := tst.CMS().SignerInfos()
	if len(signers) == 0 {
		return nil
	}
	return signers[0].UnsignedAttributes
}

// timeStampTokenProductionComparatorHashTableSize ports the private getHashTableSize(ASN1Sequence):
// the ats-hash-index table's members are each themselves a SEQUENCE (the per-category hash
// list), and the total count is the sum of their sizes.
func timeStampTokenProductionComparatorHashTableSize(derSequence []byte) int {
	element, _, err := asn1ber.Parse(derSequence)
	if err != nil || !element.IsConstructed() {
		return 0
	}
	recordsNumber := 0
	for _, child := range element.Children() {
		if child.IsConstructed() {
			recordsNumber += len(child.Children())
		}
	}
	return recordsNumber
}
