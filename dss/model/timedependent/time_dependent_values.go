// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/timedependent/TimeDependentValues.java (DSS 6.5.RC1).
package timedependent

import (
	"iter"
	"strings"
	"time"
)

// Values is an immutable list of time-dependent values, with the latest value
// first.
//
// java.io.Serializable is dropped silently (no Go counterpart). The "protected final List<T>
// list" field Java subclasses (MutableTimeDependentValues) mutate directly is kept unexported
// here: MutableValues lives in this same package and reaches it directly, which
// is the Go counterpart of a protected field accessed from a subclass.
type Values[T TimeDependent] struct {
	list []T
}

// NewValues is the empty list of values.
func NewValues[T TimeDependent]() *Values[T] {
	return &Values[T]{}
}

// NewValuesFrom is the copy constructor.
func NewValuesFrom[T TimeDependent](srcList []T) *Values[T] {
	v := &Values[T]{list: make([]T, 0, len(srcList))}
	v.list = append(v.list, srcList...)
	return v
}

// Iterator returns a range-over-func iterator on the immutable list, the Go counterpart of
// Java's Iterable<T>#iterator().
func (v *Values[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, x := range v.list {
			if !yield(x) {
				return
			}
		}
	}
}

// Latest gets the latest time dependent value, or the zero value of T when the list is empty
// (Java's null). Port of getLatest().
func (v *Values[T]) Latest() T {
	if len(v.list) == 0 {
		var zero T
		return zero
	}
	return v.list[0]
}

// Current gets the value with the date d if present, or the zero value of T otherwise
// (Java's null). Port of getCurrent(Date).
func (v *Values[T]) Current(d time.Time) T {
	for _, x := range v.list {
		if !x.StartDate().After(d) {
			endDate := x.EndDate()
			if endDate.IsZero() || endDate.After(d) {
				return x
			}
		}
	}
	var zero T
	return zero
}

// After gets a list of time dependent values occurred after notBefore. Port of
// getAfter(Date).
func (v *Values[T]) After(notBefore time.Time) []T {
	result := make([]T, 0)
	for _, x := range v.list {
		endDate := x.EndDate()
		if endDate.IsZero() || !endDate.Before(notBefore) {
			result = append(result, x)
		}
	}
	return result
}

// String renders the list the way java.util.List#toString() would, joining each element's
// String() (or default formatting when T does not implement fmt.Stringer) with ", " between
// square brackets. Port of toString().
func (v *Values[T]) String() string {
	parts := make([]string, len(v.list))
	for i, x := range v.list {
		parts[i] = timeDependentValueString(x)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// timeDependentValueString renders a single element the way Java's String.valueOf would,
// preferring a Stringer implementation when present.
func timeDependentValueString(x any) string {
	if s, ok := x.(interface{ String() string }); ok {
		return s.String()
	}
	return "null"
}
