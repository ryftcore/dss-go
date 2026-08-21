// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESUnsignedAttributes.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
)

// CAdESUnsignedAttributes represents the CAdES Unsigned attributes. Port of the class
// CAdESUnsignedAttributes, extending CAdESSigProperties.
type CAdESUnsignedAttributes struct {
	CAdESSigProperties
}

// newCAdESUnsignedAttributes is the port of the package-private
// CAdESUnsignedAttributes(ASN1Set) constructor.
func newCAdESUnsignedAttributes(attributeTable cmscore.Attributes, exists bool) *CAdESUnsignedAttributes {
	return &CAdESUnsignedAttributes{CAdESSigProperties: newCAdESSigProperties(attributeTable, exists)}
}

// CAdESUnsignedAttributesBuild builds the CAdESUnsignedAttributes from a SignerInfo.
// Extraction from SignerInfo allows keeping the actual order. Port of the static
// build(SignerInformation).
func CAdESUnsignedAttributesBuild(signerInformation *cmscore.SignerInfo) *CAdESUnsignedAttributes {
	// Extraction from SignerInfo allows to keep actual order
	return newCAdESUnsignedAttributes(signerInformation.UnsignedAttributes, signerInformation.HasUnsignedAttributes())
}

// Attributes overrides CAdESSigProperties.Attributes(): multiple timestamps need to be sorted
// in CAdES by their production date. Port of getAttributes().
func (u *CAdESUnsignedAttributes) Attributes() []*CAdESAttribute {
	attributes := u.CAdESSigProperties.Attributes()
	return cadesUnsignedAttributesSort(attributes)
}

// cadesUnsignedAttributesSort ports the private sortAttributes(List<CAdESAttribute>): a
// bubble sort, kept as one rather than a stable library sort because
// cadesAttributeTimeStampCompare is not a total order (many pairs compare equal, e.g. any two
// non-timestamp, non-evidence-record attributes) and Java relies on the bubble sort's specific
// adjacent-swap pattern to leave those pairs in their original relative order.
func cadesUnsignedAttributesSort(attributes []*CAdESAttribute) []*CAdESAttribute {
	for ii := 0; ii < len(attributes)-1; ii++ {
		for jj := 0; jj < len(attributes)-ii-1; jj++ {
			cadesAttribute := attributes[jj]
			nextCAdESAttribute := attributes[jj+1]
			if cadesAttributeTimeStampCompare(cadesAttribute, nextCAdESAttribute) > 0 {
				attributes[jj], attributes[jj+1] = attributes[jj+1], attributes[jj]
			}
		}
	}
	return attributes
}

// cadesAttributeTimeStampCompare ports the private static final class
// CAdESAttributeTimeStampComparator.compare(CAdESAttribute, CAdESAttribute).
func cadesAttributeTimeStampCompare(o1, o2 *CAdESAttribute) int {
	result := cadesAttributeCompareByType(o1, o2)
	if result == 0 {
		result = cadesAttributeCompareByTimeStampToken(o1, o2)
	}
	if result == 0 {
		result = cadesAttributeCompareByTimestampType(o1, o2)
	}
	if result == 0 {
		result = cadesAttributeCompareByEvidenceRecord(o1, o2)
	}
	return result
}

// cadesAttributeCompareByType ports the private compareByType(CAdESAttribute, CAdESAttribute):
// evidence records are always the last, timestamps are the last but before evidence records.
func cadesAttributeCompareByType(attributeOne, attributeTwo *CAdESAttribute) int {
	if !attributeOne.IsEvidenceRecord() && attributeTwo.IsEvidenceRecord() {
		return -1
	} else if attributeOne.IsEvidenceRecord() && !attributeTwo.IsEvidenceRecord() {
		return 1
	} else if !attributeOne.IsTimeStampToken() && attributeTwo.IsTimeStampToken() {
		return -1
	} else if attributeOne.IsTimeStampToken() && !attributeTwo.IsTimeStampToken() {
		return 1
	}
	return 0
}

// cadesAttributeCompareByTimeStampToken ports the private
// compareByTimeStampToken(CAdESAttribute, CAdESAttribute).
func cadesAttributeCompareByTimeStampToken(attributeOne, attributeTwo *CAdESAttribute) int {
	var current, next *cmscore.TimeStampToken
	if attributeOne.IsTimeStampToken() {
		current = attributeOne.ToTimeStampToken()
	}
	if attributeTwo.IsTimeStampToken() {
		next = attributeTwo.ToTimeStampToken()
	}
	if current != nil && next != nil {
		comparator := NewTimeStampTokenProductionComparator()
		return comparator.Compare(current, next)
	}
	return 0
}

// cadesAttributeCompareByTimestampType ports the private
// compareByTimestampType(CAdESAttribute, CAdESAttribute).
func cadesAttributeCompareByTimestampType(attributeOne, attributeTwo *CAdESAttribute) int {
	timestampTypeOne := attributeOne.TimestampTokenType()
	timestampTypeTwo := attributeTwo.TimestampTokenType()
	if timestampTypeOne != "" && timestampTypeTwo != "" {
		return timestampTypeOne.Compare(timestampTypeTwo)
	}
	return 0
}

// cadesAttributeCompareByEvidenceRecord ports the private
// compareByEvidenceRecord(CAdESAttribute, CAdESAttribute).
func cadesAttributeCompareByEvidenceRecord(attributeOne, attributeTwo *CAdESAttribute) int {
	var evidenceRecordOne, evidenceRecordTwo *asn1ber.Element
	if attributeOne.IsEvidenceRecord() {
		evidenceRecordOne = attributeOne.ToEvidenceRecord()
	}
	if attributeTwo.IsEvidenceRecord() {
		evidenceRecordTwo = attributeTwo.ToEvidenceRecord()
	}
	if evidenceRecordOne != nil && evidenceRecordTwo != nil {
		comparator := NewEvidenceRecordProductionComparator()
		return comparator.Compare(evidenceRecordOne, evidenceRecordTwo)
	}
	return 0
}
