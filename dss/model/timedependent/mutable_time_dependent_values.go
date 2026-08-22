// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/timedependent/MutableTimeDependentValues.java (DSS 6.5.RC1).
package timedependent

import "reflect"

// MutableValues is a mutable list of time-dependent values.
type MutableValues[T TimeDependent] struct {
	Values[T]
}

// NewMutableTimeDependentValues is the empty constructor.
func NewMutableTimeDependentValues[T TimeDependent]() *MutableValues[T] {
	return &MutableValues[T]{}
}

// NewMutableTimeDependentValuesFrom is the default constructor from a source list.
func NewMutableTimeDependentValuesFrom[T TimeDependent](srcList []T) *MutableValues[T] {
	return &MutableValues[T]{Values: *NewValuesFrom(srcList)}
}

// Clear clears the current list.
//
// Java declares this method synchronized; the Go port is not goroutine-safe, matching the
// rest of the value objects in dss-model.
func (v *MutableValues[T]) Clear() {
	v.list = v.list[:0]
}

// AddOldest adds the value only if it is the oldest in the current list.
//
// Panics with the Java message when x is a nil pointer/interface (Java
// Objects.requireNonNull(x, "Cannot add null")); panics with Java's IllegalArgumentException
// message when x's end date overlaps an existing entry's start date.
//
// Java declares this method synchronized; the Go port is not goroutine-safe, matching the
// rest of the value objects in dss-model.
func (v *MutableValues[T]) AddOldest(x T) {
	if timeDependentIsNil(x) {
		panic("Cannot add null")
	}
	if len(v.list) != 0 {
		endDate := x.EndDate()
		if !endDate.IsZero() {
			for _, y := range v.list {
				if !y.StartDate().IsZero() && endDate.After(y.StartDate()) {
					panic("Cannot add overlapping item")
				}
			}
		}
	}
	v.list = append(v.list, x)
}

// List gets the current list. Port of getList().
func (v *MutableValues[T]) List() []T {
	return v.list
}

// timeDependentIsNil reports whether the generic value x holds a nil pointer or interface.
// TimeDependent is always implemented by a pointer or interface type in this codebase; Go
// generics offer no direct nil comparison against a type parameter, so reflection stands in
// for Java's simple x == null check.
func timeDependentIsNil[T TimeDependent](x T) bool {
	rv := reflect.ValueOf(x)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return rv.IsNil()
	default:
		return false
	}
}
