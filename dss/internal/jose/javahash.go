// Ported from java.util.HashMap's iteration order (OpenJDK 21) and java.lang.String.hashCode.
// There is no jose4j class behind this file; it exists because DSS's own JsonObject, in the
// eu.europa.esig.dss.jades.JsonObject no-argument constructor, wraps a plain HashMap, and the
// resulting member order reaches serialized JAdES bytes.
package jose

import "sort"

// JavaStringHashCode is java.lang.String.hashCode(): s[0]*31^(n-1) + s[1]*31^(n-2) + ... + s[n-1],
// evaluated in 32-bit two's-complement arithmetic over UTF-16 code units.
//
// Go strings are UTF-8, so the code units have to be recovered: a rune below U+10000 is one unit,
// anything above is the surrogate pair Java would hold. Header names are ASCII in practice, but
// getting this wrong for a non-ASCII key would silently reorder a signed header.
func JavaStringHashCode(s string) int32 {
	var h int32
	for _, r := range s {
		if r < 0x10000 {
			h = 31*h + int32(r)
			continue
		}
		r -= 0x10000
		hi := int32(0xD800 + (r >> 10))
		lo := int32(0xDC00 + (r & 0x3FF))
		h = 31*h + hi
		h = 31*h + lo
	}
	return h
}

// javaHashMapSpread is HashMap.hash(Object): h ^ (h >>> 16), which mixes the high bits down so
// that they take part in the (n-1)-masked bucket index.
func javaHashMapSpread(key string) uint32 {
	h := uint32(JavaStringHashCode(key))
	return h ^ (h >> 16)
}

// JavaHashMapOrder returns keys reordered the way iterating a java.util.HashMap would yield them,
// given that they were put in the order supplied and that none was ever removed.
//
// Why this is needed at all: JsonObject's no-argument constructor is `map = new HashMap<>()`, and
// three DSS call sites build multi-member objects through it - LevelBaselineLT.getRVals
// ({crlVals, ocspVals}), LevelBaselineLT.getTstVd ({xVals, rVals}) and, indirectly,
// LevelBaselineB. Those objects are serialized into 'etsiU' components, which an archive
// timestamp then covers. "Whatever order Go's map gives" would produce bytes upstream never
// produces, so the order is reproduced instead of approximated.
//
// The model: a table of capacity cap (16, doubling when size exceeds 0.75*cap), each key in
// bucket spread(key) & (cap-1), iteration walking buckets in index order and, within a bucket,
// following the linked list. New entries are appended at the tail of their bucket, and a resize
// splits each bucket into a "lo" list (same index) and a "hi" list (index + oldCap) while
// preserving relative order - so relative order within a bucket is always insertion order, and
// the whole thing reduces to a stable sort by final bucket index.
//
// Deliberately not modelled: treeification. HashMap converts a bucket to a red-black tree once it
// holds 8 entries AND the table has at least 64 slots, at which point iteration order within that
// bucket becomes hash/comparison order rather than insertion order. Reaching it needs at least 64
// keys with a heavy collision pattern; no DSS JSON object comes within an order of magnitude of
// that (the largest is two members). If one ever does, this returns the untreeified order, which
// is why callers that care use insertion-ordered objects instead.
func JavaHashMapOrder(keys []string) []string {
	n := len(keys)
	if n < 2 {
		out := make([]string, n)
		copy(out, keys)
		return out
	}

	// Table capacity after n puts: start at 16, double whenever size exceeds 0.75*capacity.
	// (HashMap resizes when ++size > threshold, so capacity 16 holds 12 entries.)
	capacity := 16
	for n > capacity*3/4 {
		capacity <<= 1
	}
	mask := uint32(capacity - 1)

	type entry struct {
		key    string
		bucket uint32
		order  int
	}
	entries := make([]entry, n)
	for i, k := range keys {
		entries[i] = entry{key: k, bucket: javaHashMapSpread(k) & mask, order: i}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].bucket < entries[j].bucket
	})

	out := make([]string, n)
	for i, e := range entries {
		out[i] = e.key
	}
	return out
}
