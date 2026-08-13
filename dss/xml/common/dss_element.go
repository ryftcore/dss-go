// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/DSSElement.java (DSS 6.5.RC1).
package common

// DSSElement is the XML element interface.
type DSSElement interface {
	// TagName returns the element tag name. Ports getTagName().
	TagName() string

	// Namespace returns the namespace, or nil when the element carries none. Ports
	// getNamespace().
	Namespace() *DSSNamespace

	// URI returns the namespace URI, or "" when Namespace() is nil. Ports getURI(), whose
	// Java returns null in that case.
	//
	// Judgment call: xpath_query_element_item.go's element matcher checks
	// "element.getURI() == null || element.getURI().equals(node.getNamespaceURI())"; with
	// URI() == "" standing in for null, the first disjunct is unconditionally true whenever
	// URI() == "", so the overall condition holds regardless of the node's own namespace -
	// exactly Java's behaviour for a nil DSSNamespace. This only diverges from Java for a
	// DSSElement explicitly constructed with a non-nil DSSNamespace whose URI is itself "",
	// which no definition in this codebase does (every DSSNamespace here carries a real URI).
	URI() string

	// IsSameTagName reports whether value equals this element's tag name. Ports
	// isSameTagName(String).
	IsSameTagName(value string) bool
}

// dssElementFromDefinition is the anonymous DSSElement the Java static factory
// DSSElement.fromDefinition builds.
type dssElementFromDefinition struct {
	localName string
	namespace *DSSNamespace
}

// TagName implements DSSElement.
func (e *dssElementFromDefinition) TagName() string {
	return e.localName
}

// Namespace implements DSSElement.
func (e *dssElementFromDefinition) Namespace() *DSSNamespace {
	return e.namespace
}

// URI implements DSSElement.
func (e *dssElementFromDefinition) URI() string {
	if e.namespace != nil {
		return e.namespace.Uri()
	}
	return ""
}

// IsSameTagName implements DSSElement.
func (e *dssElementFromDefinition) IsSameTagName(value string) bool {
	return e.localName == value
}

// DSSElementFromDefinition creates a DSSElement from the given definition. Ports the static
// factory method DSSElement.fromDefinition(String, DSSNamespace).
func DSSElementFromDefinition(localName string, namespace *DSSNamespace) DSSElement {
	return &dssElementFromDefinition{localName: localName, namespace: namespace}
}
