// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/POEComparator.java (DSS 6.5.RC1).
//
// See poe.go for the POE hierarchy note. java.io.Serializable and the
// serialVersionUID have no Go counterpart and are dropped.
package vpfswatsp

// POEComparator compares two POE instances, by its production time, origin and
// covered context.
//
// The comparator returns the following values:
//
//	-1 if the poe1 is preferred over poe2
//	 0 if the POEs are equal
//	 1 if the poe2 is preferred over poe1
type POEComparator struct{}

// NewPOEComparator is the default constructor.
func NewPOEComparator() *POEComparator {
	return &POEComparator{}
}

// Compare compares the two POEs. Port of compare(POE, POE).
func (c *POEComparator) Compare(poe1 POE, poe2 POE) int {
	result := c.compareByTime(poe1, poe2)
	if result == 0 {
		result = c.compareByType(poe1, poe2)
	}
	if result == 0 {
		result = c.compareByTimestampType(poe1, poe2)
	}
	if result == 0 {
		result = c.compareByTimestampedReferences(poe1, poe2)
	}
	return result
}

// compareByTime ports the private compareByTime(POE, POE).
func (c *POEComparator) compareByTime(poe1 POE, poe2 POE) int {
	return poe1.Time().Compare(poe2.Time())
}

// compareByType ports the private compareByType(POE, POE).
func (c *POEComparator) compareByType(poe1 POE, poe2 POE) int {
	// POE defined by a timestamp is preferred over a POE defined by a control time
	if poe1.IsTokenProvided() && !poe2.IsTokenProvided() {
		return -1
	} else if !poe1.IsTokenProvided() && poe2.IsTokenProvided() {
		return 1
	}
	return 0
}

// compareByTimestampType ports the private compareByTimestampType(POE, POE).
//
// Java's "instanceof TimestampPOE" is a type assertion here; the null-guarded
// TimestampType comparison keeps its guard, the empty TimestampType standing in
// for Java's null (the wrapper returns it for a time-stamp with no type).
func (c *POEComparator) compareByTimestampType(poe1 POE, poe2 POE) int {
	if timestampPOE1, ok1 := poe1.(*TimestampPOE); ok1 {
		if timestampPOE2, ok2 := poe2.(*TimestampPOE); ok2 {
			poe1TstType := timestampPOE1.TimestampType()
			poe2TstType := timestampPOE2.TimestampType()
			if poe1TstType != "" && poe2TstType != "" {
				return poe1TstType.Compare(poe2TstType)
			}
		}
	}
	return 0
}

// compareByTimestampedReferences ports the private
// compareByTimestampedReferences(POE, POE).
//
// Java guards both lists against null; neither getPOEObjects() implementation
// can return null (the base answers Collections.emptyList(), the two subclasses
// a wrapper's list), so the guard is always taken and only that branch is
// ported - a Go nil slice being the empty list, not a missing one.
func (c *POEComparator) compareByTimestampedReferences(poe1 POE, poe2 POE) int {
	poe1References := poe1.POEObjects()
	poe2References := poe2.POEObjects()
	if len(poe1References) < len(poe2References) {
		return -1
	} else if len(poe1References) > len(poe2References) {
		return 1
	}
	return 0
}

// Before checks if the poe1 is before the poe2. Port of before(POE, POE).
func (c *POEComparator) Before(poe1 POE, poe2 POE) bool {
	return c.Compare(poe1, poe2) < 0
}
