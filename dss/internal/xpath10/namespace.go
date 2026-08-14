package xpath10

import "sort"

// NamespaceContext is the prefix-to-URI map an expression is compiled against. It is the port
// of eu.europa.esig.dss.xml.utils.NamespaceContextMap (DSS 6.5.RC1), which is the
// javax.xml.namespace.NamespaceContext that XPathUtils installs on every XPath it builds.
//
// A plain map is the whole state: the Java class also keeps a reverse URI-to-prefixes index,
// but that index exists only to serve getPrefix / getPrefixes, which an XPath evaluator never
// calls, and rebuilding it here would buy nothing but a second thing to keep consistent.
//
// The nil map is usable: every prefix resolves to the null namespace.
type NamespaceContext map[string]string

// defaultNamespaceKey is the sentinel NamespaceContextOf stores the innermost default namespace
// declaration (xmlns="u", no prefix) under. It is a control character, so it can never collide
// with a real XML prefix (NCName forbids it), and it is queried only by parser.go's deliberate
// leniency for a QName whose prefix component is present but empty (a stray leading colon, e.g.
// "descendant:::Signature") - never by ordinary prefix resolution. A genuinely unprefixed node
// test (no colon at all) never even reaches NamespaceURI: per XPath 1.0 clause 2.3 it is
// resolved to the null namespace directly, bypassing the map entirely, so storing the default
// namespace here cannot change that pinned behaviour. See parseNodeTest's stray-colon case for
// why the distinction matters.
const defaultNamespaceKey = "\x00"

// NamespaceURI resolves prefix. An unregistered prefix resolves to the null namespace, "",
// rather than reporting an error - NamespaceContextMap.getNamespaceURI returns
// XMLConstants.NULL_NS_URI for a prefix it does not know, so a typo in a prefix silently turns
// a namespaced node test into an unnamespaced one. That is load-bearing behaviour, not an
// accident to be fixed: expressions compiled before a namespace is registered must keep
// matching what Xalan matched.
func (c NamespaceContext) NamespaceURI(prefix string) string {
	return c[prefix]
}

// RegisterNamespace binds prefix to uri, replacing any existing binding, and reports whether
// the prefix was new - the return value of NamespaceContextMap.registerNamespace.
//
// It panics on a nil receiver, since a nil map cannot be written; build one with make or a
// composite literal first.
func (c NamespaceContext) RegisterNamespace(prefix, uri string) bool {
	if c == nil {
		panic("xpath10: RegisterNamespace on a nil NamespaceContext")
	}
	_, existed := c[prefix]
	c[prefix] = uri
	return !existed
}

// Prefix returns a prefix bound to uri and reports whether one exists.
//
// NamespaceContextMap.getPrefix returns an arbitrary member of a HashSet, so upstream's answer
// is unspecified when a URI carries several prefixes; this returns the lexicographically
// smallest so that callers and golden files do not depend on map iteration order. Nothing in
// DSS calls it during evaluation - it exists for the JAXP interface - so pinning the choice
// costs no compatibility.
func (c NamespaceContext) Prefix(uri string) (string, bool) {
	best, found := "", false
	for prefix, u := range c {
		if u != uri {
			continue
		}
		if !found || prefix < best {
			best, found = prefix, true
		}
	}
	return best, found
}

// Prefixes returns every prefix bound to uri, sorted, for the same reason Prefix sorts.
func (c NamespaceContext) Prefixes(uri string) []string {
	var out []string
	for prefix, u := range c {
		if u == uri {
			out = append(out, prefix)
		}
	}
	sort.Strings(out)
	return out
}

// PrefixMap returns a copy of the bindings, matching NamespaceContextMap.getPrefixMap, which
// hands out a copy so a caller cannot mutate the registry through it.
func (c NamespaceContext) PrefixMap() map[string]string {
	out := make(map[string]string, len(c))
	for k, v := range c {
		out[k] = v
	}
	return out
}
