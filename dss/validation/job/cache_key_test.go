package job

import "testing"

func TestCacheKey_EqualityByNormalizedURL(t *testing.T) {
	k1 := NewCacheKey("https://example.org/tsl.xml")
	k2 := NewCacheKey("https://example.org/tsl.xml")
	if k1 != k2 {
		t.Fatalf("two CacheKeys built from the same URL should compare equal: %v != %v", k1, k2)
	}
	if !k1.Equals(k2) {
		t.Fatal("Equals() should agree with ==")
	}

	k3 := NewCacheKey("https://example.org/other.xml")
	if k1 == k3 {
		t.Fatal("CacheKeys built from different URLs should not compare equal")
	}
}

func TestCacheKey_UsableAsMapKey(t *testing.T) {
	m := map[CacheKey]int{}
	m[NewCacheKey("https://example.org/tsl.xml")] = 1
	m[NewCacheKey("https://example.org/tsl.xml")] = 2
	if len(m) != 1 {
		t.Fatalf("expected a single map entry for two CacheKeys from the same URL, got %d", len(m))
	}
	if got := m[NewCacheKey("https://example.org/tsl.xml")]; got != 2 {
		t.Fatalf("m[key] = %d, want 2", got)
	}
}

func TestCacheKey_PanicsOnEmptyURL(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic constructing a CacheKey from an empty URL")
		}
	}()
	NewCacheKey("")
}
