package diagnostic

import (
	"fmt"
	"testing"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
)

// TestXmlTrustedEntityBuilderGetLangAndValuesIsDeterministic covers T31-PERF-BIGO-002: the
// trusted-entity builder ranged its map[string][]string directly, so the emitted
// XmlLangAndValue order (multi-language names / trade names / service names) differed between
// runs. Languages are now emitted sorted, each language's values in their given order, like the
// sibling XmlTrustServiceProviderBuilder.
func TestXmlTrustedEntityBuilderGetLangAndValuesIsDeterministic(t *testing.T) {
	names := map[string][]string{
		"fr": {"Fournisseur A", "Fournisseur B"},
		"de": {"Anbieter"},
		"en": {"Provider B", "Provider A"},
		"bg": {"Dostavchik"},
		"nl": {"Aanbieder"},
	}
	want := []string{
		"bg|Dostavchik",
		"de|Anbieter",
		"en|Provider B",
		"en|Provider A",
		"fr|Fournisseur A",
		"fr|Fournisseur B",
		"nl|Aanbieder",
	}

	b := &XmlTrustedEntityBuilder{}
	// Go randomizes map iteration order per range, so a handful of repetitions would have
	// caught the old behaviour with overwhelming probability.
	for i := 0; i < 50; i++ {
		if got := flattenLangAndValues(b.getLangAndValues(names)); !equalStrings(got, want) {
			t.Fatalf("run %d: got %v, want %v", i, got, want)
		}
	}

	if got := b.getLangAndValues(nil); got != nil {
		t.Errorf("a nil map must yield nil, got %v", got)
	}
	if got := b.getLangAndValues(map[string][]string{}); got != nil {
		t.Errorf("an empty map must yield nil, got %v", got)
	}
}

func flattenLangAndValues(items []*jaxb.XmlLangAndValue) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		lang := "<nil>"
		if item.Lang != nil {
			lang = *item.Lang
		}
		out = append(out, fmt.Sprintf("%s|%s", lang, item.Value))
	}
	return out
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
