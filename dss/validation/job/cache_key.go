// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/CacheKey.java (DSS 6.5.RC1).
package job

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/spi"
)

// CacheKey defines a key for a cache record.
//
// Java's CacheKey overrides equals()/hashCode() on the normalized key string so two
// instances built from the same URL compare equal and hash identically when used as a
// HashMap key. CacheKey is therefore ported as a value type (a single comparable string
// field) rather than a pointer, so plain Go map[CacheKey]... equality (==) already matches
// Java's equals()/hashCode() contract, per the wrapper-equality porting rule.
type CacheKey struct {
	// key is the key of the entry.
	key string
}

// NewCacheKey creates a CacheKey from a document url string. Panics if url is empty,
// mirroring Java's Objects.requireNonNull("URL cannot be null.") — Go strings cannot be
// null, so the empty string is the port of the missing-value case.
func NewCacheKey(url string) CacheKey {
	if url == "" {
		panic("URL cannot be null.")
	}
	return CacheKey{key: spi.DSSUtilsNormalizedString(url)}
}

// Key returns the encoded key. Port of getKey().
func (k CacheKey) Key() string {
	return k.key
}

// String implements fmt.Stringer. Port of toString().
func (k CacheKey) String() string {
	return fmt.Sprintf("CacheKey with the key [%s]", k.key)
}

// Equals reports whether other is a CacheKey with the same key, mirroring Java's
// equals(Object). Provided for parity with the Java API; plain == is equivalent and
// idiomatic in Go.
func (k CacheKey) Equals(other CacheKey) bool {
	return k == other
}
