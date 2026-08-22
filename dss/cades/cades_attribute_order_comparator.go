// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESAttributeOrderComparator.java (DSS 6.5.RC1).
package cades

// AttributeOrderComparator compares the order - and only the original order - of
// CAdESAttributes from an AttributeTable. Port of the class CAdESAttributeOrderComparator,
// implementing Comparator<Attribute>.
//
// AttributeCompare returns:
//   - -1 if attributeOne has original order before attributeTwo
//   - 0 if attributes have the same order (should not happen)
//   - 1 if attributeOne has original order after attributeTwo
func AttributeCompare(attributeOne, attributeTwo *Attribute) int {
	if attributeOne.Order() != nil && attributeTwo.Order() != nil {
		if *attributeOne.Order() < *attributeTwo.Order() {
			return -1
		} else if *attributeOne.Order() > *attributeTwo.Order() {
			return 1
		}
	}
	return 0
}

// CAdESAttributeOrderComparator is the port of the class of the same name, provided as a
// sort.Interface-compatible less-function wrapper for the common case of sorting a slice by
// original attribute order.
type AttributeOrderComparator struct{}

// NewCAdESAttributeOrderComparator is the port of the default constructor.
func NewCAdESAttributeOrderComparator() AttributeOrderComparator {
	return AttributeOrderComparator{}
}

// Compare is the port of compare(CAdESAttribute, CAdESAttribute).
func (AttributeOrderComparator) Compare(attributeOne, attributeTwo *Attribute) int {
	return AttributeCompare(attributeOne, attributeTwo)
}
