// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESUnsignedAttributes.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
)

// CAdESUnsignedAttributes represents the CAdES Unsigned attributes. Port of the class
// UnsignedAttributes, extending SigProperties.
type UnsignedAttributes struct {
	SigProperties
}

// newCAdESUnsignedAttributes is the port of the package-private
// UnsignedAttributes(ASN1Set) constructor.
func newUnsignedAttributes(attributeTable cmscore.Attributes, exists bool) *UnsignedAttributes {
	return &UnsignedAttributes{SigProperties: newSigProperties(attributeTable, exists)}
}

// UnsignedAttributesBuild builds the UnsignedAttributes from a SignerInfo.
// Extraction from SignerInfo allows keeping the actual order. Port of the static
// build(SignerInformation).
func UnsignedAttributesBuild(signerInformation *cmscore.SignerInfo) *UnsignedAttributes {
	// Extraction from SignerInfo allows to keep actual order
	return newUnsignedAttributes(signerInformation.UnsignedAttributes, signerInformation.HasUnsignedAttributes())
}

// Attributes overrides SigProperties.Attributes(): multiple timestamps need to be sorted
// in CAdES by their production date. Port of getAttributes().
func (u *UnsignedAttributes) Attributes() []*Attribute {
	attributes := u.SigProperties.Attributes()
	return cadesUnsignedAttributesSort(attributes)
}

// cadesUnsignedAttributesSort ports the private sortAttributes(List<CAdESAttribute>): a
// bubble sort, kept as one rather than a stable library sort because
// cadesAttributeTimeStampCompare is not a total order (many pairs compare equal, e.g. any two
// non-timestamp, non-evidence-record attributes) and Java relies on the bubble sort's specific
// adjacent-swap pattern to leave those pairs in their original relative order.
func cadesUnsignedAttributesSort(attributes []*Attribute) []*Attribute {
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
// CAdESAttributeTimeStampComparator.compare(Attribute, Attribute).
func cadesAttributeTimeStampCompare(o1, o2 *Attribute) int {
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
func cadesAttributeCompareByType(attributeOne, attributeTwo *Attribute) int {
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
// compareByTimeStampToken(Attribute, Attribute).
func cadesAttributeCompareByTimeStampToken(attributeOne, attributeTwo *Attribute) int {
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
// compareByTimestampType(Attribute, Attribute).
func cadesAttributeCompareByTimestampType(attributeOne, attributeTwo *Attribute) int {
	timestampTypeOne := attributeOne.TimestampTokenType()
	timestampTypeTwo := attributeTwo.TimestampTokenType()
	if timestampTypeOne != "" && timestampTypeTwo != "" {
		return timestampTypeOne.Compare(timestampTypeTwo)
	}
	return 0
}

// cadesAttributeCompareByEvidenceRecord ports the private
// compareByEvidenceRecord(Attribute, Attribute).
func cadesAttributeCompareByEvidenceRecord(attributeOne, attributeTwo *Attribute) int {
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
