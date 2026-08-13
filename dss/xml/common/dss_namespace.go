// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/DSSNamespace.java (DSS 6.5.RC1).
package common

import "fmt"

// DSSNamespace defines an XML namespace by its URI and preferred prefix.
type DSSNamespace struct {
	uri    string
	prefix string
}

// NewDSSNamespace creates a DSSNamespace. Ports DSSNamespace(String uri, String prefix).
func NewDSSNamespace(uri, prefix string) *DSSNamespace {
	return &DSSNamespace{uri: uri, prefix: prefix}
}

// Uri returns the namespace URI. Ports getUri().
func (n *DSSNamespace) Uri() string {
	return n.uri
}

// Prefix returns the namespace prefix. Ports getPrefix().
func (n *DSSNamespace) Prefix() string {
	return n.prefix
}

// IsSameUri reports whether paramUri equals this DSSNamespace's URI.
func (n *DSSNamespace) IsSameUri(paramUri string) bool {
	return n.uri == paramUri
}

// String mirrors Java's toString() exactly, including its missing closing quote after uri
// ("DSSNamespace [uri='<uri>, prefix='<prefix>]") - a verbatim, harmless, debug-only upstream
// quirk reproduced rather than fixed.
func (n *DSSNamespace) String() string {
	return fmt.Sprintf("DSSNamespace [uri='%s, prefix='%s]", n.uri, n.prefix)
}
