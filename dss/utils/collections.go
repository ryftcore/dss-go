// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1),
// collection/map-related methods, matching
// org.apache.commons.collections4.CollectionUtils / MapUtils semantics.
//
// Java's Collection<T>/Map<K,V> raw-typed helpers become Go generics over
// slices and maps.

package utils

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
