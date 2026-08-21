// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampedReference.java (DSS 6.5.RC1).
package validation

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// TimestampedReference stocks the timestamped reference, which is composed of: the timestamp
// reference category (TimestampedObjectType) and the object id in the case where the reference
// applies to the signature.
type TimestampedReference struct {
	// objectId is the id of the timestamped object.
	objectId string
	// category is the timestamped object type.
	category enumerations.TimestampedObjectType
}

// NewTimestampedReference builds a reference to the object of the given id and type.
// Port of the TimestampedReference(String, TimestampedObjectType) constructor.
func NewTimestampedReference(objectId string, category enumerations.TimestampedObjectType) *TimestampedReference {
	return &TimestampedReference{objectId: objectId, category: category}
}

// Category gets the timestamped object type. Port of getCategory().
func (r *TimestampedReference) Category() enumerations.TimestampedObjectType {
	return r.category
}

// ObjectId gets the timestamped object Id. Port of getObjectId().
func (r *TimestampedReference) ObjectId() string {
	return r.objectId
}

// Equals reports whether both references carry the same object id and category.
// Port of equals(Object).
//
// NOTE: hashCode() has no Go counterpart; upstream needs it only to key the JDK hash
// collections, which this port replaces with slices keyed on Equals.
func (r *TimestampedReference) Equals(other *TimestampedReference) bool {
	if r == other {
		return true
	}
	if other == nil {
		return false
	}
	if r.category != other.category {
		return false
	}
	return r.objectId == other.objectId
}

// String renders the reference. Port of toString().
func (r *TimestampedReference) String() string {
	return fmt.Sprintf("TimestampedReference with Id [%s] and type [%s]", r.ObjectId(), r.Category())
}
