// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1),
// collection/map-related methods, matching
// org.apache.commons.collections4.CollectionUtils / MapUtils semantics.
//
// Java's Collection<T>/Map<K,V> raw-typed helpers become Go generics over
// slices and maps.

package utils

import "sort"

// IsCollectionEmpty checks if the collection is nil or empty.
func IsCollectionEmpty[T any](collection []T) bool {
	return len(collection) == 0
}

// IsCollectionNotEmpty checks if the collection is not nil nor empty.
func IsCollectionNotEmpty[T any](collection []T) bool {
	return len(collection) != 0
}

// CollectionSize gets the size of the collection.
func CollectionSize[T any](collection []T) int {
	return len(collection)
}

// IsMapEmpty checks if the map is nil or empty.
func IsMapEmpty[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

// IsMapNotEmpty checks if the map is not nil nor empty.
func IsMapNotEmpty[K comparable, V any](m map[K]V) bool {
	return len(m) != 0
}

// MapSize gets the size of the map.
func MapSize[K comparable, V any](m map[K]V) int {
	return len(m)
}

// ReverseList creates a reversed copy of the list.
func ReverseList[T any](list []T) []T {
	reversed := make([]T, len(list))
	for i, v := range list {
		reversed[len(list)-1-i] = v
	}
	return reversed
}

// ContainsAny returns whether superCollection contains any element of
// subCollection.
// Ex. {'A', 'B', 'C'}, {'B', 'C', 'D'} = TRUE
func ContainsAny[T comparable](superCollection, subCollection []T) bool {
	if len(superCollection) == 0 || len(subCollection) == 0 {
		return false
	}
	set := make(map[T]struct{}, len(superCollection))
	for _, v := range superCollection {
		set[v] = struct{}{}
	}
	for _, v := range subCollection {
		if _, ok := set[v]; ok {
			return true
		}
	}
	return false
}

// JavaHashMapStringKeyOrder reorders keys into the iteration order a
// java.util.HashMap<String, ?> (equivalently a java.util.HashSet<String>) built
// by inserting them in the given sequence yields for its entrySet()/keySet().
//
// This is the deliberate EXCEPTION to OrderedMap's rule (see ordered_map.go):
// insertion order is the right substitute wherever Java's own callers only rely
// on "arbitrary but stable", but where the order reaches an OBSERVABLE,
// byte-compared output - a marshalled report element sequence, or a "last entry
// wins" accumulation whose winner is reported - only Java's actual order gives
// parity. The order is fully determined, not JVM-dependent:
// HashMap.hash(key) = h ^ (h >>> 16) over String#hashCode, bucket index
// (n-1) & hash, buckets walked in ascending index and, inside a bucket, in
// insertion order (a resize's split preserves relative order, and treeification
// needs 8 entries in one bucket). n is the table size the map has grown to for
// the entry count, from HashMap's default capacity 16 and load factor 0.75.
//
// dss/validation/executor's JavaHashSetStringOrder is this function; the
// oracle-pinned executor test testdata/oracle/hash_order.tsv covers it.
func JavaHashMapStringKeyOrder(keys []string) []string {
	if len(keys) < 2 {
		return keys
	}
	n := 16
	for len(keys) > n*3/4 {
		n *= 2
	}
	indexed := make([]struct {
		key    string
		bucket int
	}, len(keys))
	for i, key := range keys {
		h := JavaStringHashCode(key)
		spread := h ^ int32(uint32(h)>>16)
		indexed[i].key = key
		indexed[i].bucket = int(uint32(spread) & uint32(n-1))
	}
	sort.SliceStable(indexed, func(i, j int) bool { return indexed[i].bucket < indexed[j].bucket })
	result := make([]string, len(keys))
	for i := range indexed {
		result[i] = indexed[i].key
	}
	return result
}

// JavaStringHashCode reproduces java.lang.String#hashCode for an ASCII string,
// where each byte is one UTF-16 code unit.
func JavaStringHashCode(s string) int32 {
	var h int32
	for i := 0; i < len(s); i++ {
		h = 31*h + int32(s[i])
	}
	return h
}
