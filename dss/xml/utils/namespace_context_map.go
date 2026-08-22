// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/NamespaceContextMap.java (DSS 6.5.RC1).
package utils

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/internal/xpath10"
)

// NamespaceContextMap manages a prefix <-> namespace-URI registry, as used by XPath queries.
//
// DEVIATION: upstream implements javax.xml.namespace.NamespaceContext so it can be installed
// directly on a javax.xml.xpath.XPath. Go has no such interface for internal/xpath10 (or
// anything else) to consume, so NamespaceURI/Prefix/Prefixes are plain methods rather than an
// interface implementation; namespaceContextMapToXPath10 (in xpath_utils.go) is the actual
// bridge xpath10.Compile is given.
type NamespaceContextMap struct {
	prefixMap    map[string]string
	namespaceMap map[string]map[string]struct{}
}

// NewNamespaceContextMap creates an empty NamespaceContextMap. Ports the default constructor.
func NewNamespaceContextMap() *NamespaceContextMap {
	return &NamespaceContextMap{
		prefixMap:    make(map[string]string),
		namespaceMap: make(map[string]map[string]struct{}),
	}
}

// RegisterNamespace registers prefix for namespace, replacing any existing binding for
// prefix, and reports whether prefix was not already registered. Ports
// registerNamespace(String, String).
func (m *NamespaceContextMap) RegisterNamespace(prefix, namespace string) bool {
	_, existed := m.prefixMap[prefix]
	m.prefixMap[prefix] = namespace
	m.createNamespace(prefix, namespace)
	return !existed
}

// createNamespace ports the private createNamespace(String, String) helper.
func (m *NamespaceContextMap) createNamespace(prefix, namespace string) {
	set, ok := m.namespaceMap[namespace]
	if !ok {
		set = make(map[string]struct{})
		m.namespaceMap[namespace] = set
	}
	set[prefix] = struct{}{}
}

// PrefixMap returns a copy of the prefix-to-URI bindings, so a caller cannot mutate the
// registry through it. Ports getPrefixMap().
func (m *NamespaceContextMap) PrefixMap() map[string]string {
	out := make(map[string]string, len(m.prefixMap))
	for k, v := range m.prefixMap {
		out[k] = v
	}
	return out
}

// NamespaceURI resolves prefix, returning "" (XMLConstants.NULL_NS_URI) when prefix is not
// registered. Ports getNamespaceURI(String). Java's Objects-non-null check on the argument
// has no Go analogue, since a Go string is never nil.
func (m *NamespaceContextMap) NamespaceURI(prefix string) string {
	return m.prefixMap[prefix]
}

// Prefix returns a prefix registered for namespaceURI and reports whether one exists. Ports
// getPrefix(String).
//
// JUDGMENT CALL: Java returns an arbitrary member of a HashSet (JVM-hash-order dependent,
// already unspecified upstream); this returns the lexicographically smallest registered
// prefix instead, for determinism - the same choice internal/xpath10.NamespaceContext.Prefix
// makes for the same upstream method. Nothing in this package calls Prefix/Prefixes during
// XPath evaluation - xpath10 resolves prefixes from its own NamespaceContext snapshot taken
// at Compile time (see namespaceContextMapToXPath10) - so this exists only to keep the
// type's public surface complete.
func (m *NamespaceContextMap) Prefix(namespaceURI string) (string, bool) {
	set, ok := m.namespaceMap[namespaceURI]
	if !ok || len(set) == 0 {
		return "", false
	}
	best, found := "", false
	for prefix := range set {
		if !found || prefix < best {
			best, found = prefix, true
		}
	}
	return best, found
}

// Prefixes returns every prefix registered for namespaceURI, sorted for determinism (see
// Prefix). Ports getPrefixes(String).
//
// DEVIATION: Java's getPrefixes throws NullPointerException when namespaceURI was never
// registered ("Set<String> set = namespaceMap.get(namespaceURI); return set.iterator();" -
// no null check). That NPE is not deliberate upstream behaviour worth reproducing - nothing
// in DSS ever calls getPrefixes with an unregistered URI - so this returns an empty, non-nil
// slice instead of panicking.
func (m *NamespaceContextMap) Prefixes(namespaceURI string) []string {
	set := m.namespaceMap[namespaceURI]
	out := make([]string, 0, len(set))
	for prefix := range set {
		out = append(out, prefix)
	}
	sort.Strings(out)
	return out
}

// namespaceContextMapToXPath10 adapts m to internal/xpath10's namespace-context type, which
// is structurally identical (map[string]string) but kept distinct so that xpath10 stays
// dependency-free of this package. A nil m compiles an expression with no bound prefixes,
// matching xpath10.Compile's documented nil behaviour.
func namespaceContextMapToXPath10(m *NamespaceContextMap) xpath10.NamespaceContext {
	if m == nil {
		return nil
	}
	return xpath10.NamespaceContext(m.PrefixMap())
}
