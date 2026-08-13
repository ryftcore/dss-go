// OrderedMap is a small helper for the "slice + index map" pattern PORTING.md's Collections
// section calls for wherever upstream Java iterates a HashMap/HashSet-backed collection: Java's
// HashMap iteration order is arbitrary but stable within a JVM run; Go's map iteration order is
// randomized on every run. A bare Go map is therefore the wrong port whenever the iteration
// order of a map-backed field can leak into an observable, ordered result (a returned slice, a
// formatted string, a "first match wins" search). OrderedMap keeps O(1) key lookup (the backing
// map) alongside deterministic, insertion-ordered iteration (the parallel key slice), so a
// range over it always produces the same sequence for the same sequence of insertions -
// matching the "arbitrary but stable" contract Java's callers already depend on, without
// claiming to reproduce Java's specific (and unspecified) HashMap bucket order.
package utils

// OrderedMap is an insertion-ordered map: Set/Get/Delete are O(1) amortized (Delete is O(n) in
// the rare case, since it must also splice the key out of the order slice); iteration follows
// insertion order, and re-Set-ing an existing key updates its value in place without moving it.
// The zero value is not ready to use; construct with NewOrderedMap.
type OrderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]V
}

// NewOrderedMap creates an empty OrderedMap.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{values: make(map[K]V)}
}

// Set inserts or updates the value for key. A new key is appended to the iteration order; an
// existing key keeps its current position.
func (m *OrderedMap[K, V]) Set(key K, value V) {
	if m.values == nil {
		m.values = make(map[K]V)
	}
	if _, found := m.values[key]; !found {
		m.keys = append(m.keys, key)
	}
	m.values[key] = value
}

// Get returns the value for key and whether it was found. A nil *OrderedMap reads as empty,
// mirroring Go's built-in map (a nil map is readable, just not writable).
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	v, found := m.values[key]
	return v, found
}

// Delete removes key, if present, and drops it from the iteration order.
func (m *OrderedMap[K, V]) Delete(key K) {
	if m == nil {
		return
	}
	if _, found := m.values[key]; !found {
		return
	}
	delete(m.values, key)
	for i, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
}

// Len returns the number of entries. A nil *OrderedMap reads as empty.
func (m *OrderedMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.keys)
}

// Keys returns the keys in insertion order. The returned slice is a copy; mutating it does not
// affect the OrderedMap. A nil *OrderedMap reads as empty.
func (m *OrderedMap[K, V]) Keys() []K {
	if m == nil {
		return nil
	}
	result := make([]K, len(m.keys))
	copy(result, m.keys)
	return result
}

// Values returns the values in insertion (key) order. A nil *OrderedMap reads as empty.
func (m *OrderedMap[K, V]) Values() []V {
	if m == nil {
		return nil
	}
	result := make([]V, 0, len(m.keys))
	for _, k := range m.keys {
		result = append(result, m.values[k])
	}
	return result
}

// Range calls f for every entry in insertion order, stopping early if f returns false. A nil
// *OrderedMap ranges zero times.
func (m *OrderedMap[K, V]) Range(f func(key K, value V) bool) {
	if m == nil {
		return
	}
	for _, k := range m.keys {
		if !f(k, m.values[k]) {
			return
		}
	}
}

// Reset clears the map back to empty, keeping it ready to use.
func (m *OrderedMap[K, V]) Reset() {
	m.keys = nil
	m.values = make(map[K]V)
}
