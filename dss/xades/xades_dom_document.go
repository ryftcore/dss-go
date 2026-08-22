// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/dom/XAdESDOMDocument.java (DSS 6.5.RC1).
//
// Santuario/ID-attribute replacement: Java's recursiveIdBrowse()/setIDIdentifier(Element) walk
// the whole tree once, marking each element's first case-insensitively-"Id"-named attribute as a
// w3c DOM ID attribute (Element#setIdAttribute), so that later same-document URI dereferencing
// (getElementById) works without a DTD/schema. internal/xmldom's *xmldom.Node#RegisterIDs is
// exactly that walk, ported once and shared by every caller that needs it (its own doc comment
// says so verbatim: "reproducing DOMDocument.recursiveIdBrowse"), so RecursiveIdBrowse below
// delegates to it instead of re-implementing the walk; Java's two protected helper methods
// (recursiveIdBrowse(Element), setIDIdentifier(Element)) have no separate Go counterpart as a
// result - nothing in this port subclasses DOMDocument to override them, and Go has no
// virtual dispatch across embedding for such an override to reach anyway.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// DOMDocument represents a wrapper of an *xmldom.Node Document, containing utility methods
// as well as common cached values.
type DOMDocument struct {
	*xmlutils.DOMDocument

	// xadesPathsHolders contains the list of XAdESPaths adapted to the specific signature
	// schema.
	xadesPathsHolders []definition.XAdESPath

	// signatureNodes caches the list of ds:Signature nodes, except counter signatures.
	signatureNodes    []*xmldom.Node
	signatureNodesSet bool
}

// NewDOMDocumentDefault is the default constructor instantiating a XAdES DOM document
// using XAdES 1.3.2 namespace paths. Ports the single-argument XAdESDOMDocument(Document)
// constructor.
func NewDOMDocumentDefault(document *xmldom.Node) *DOMDocument {
	return NewDOMDocument(document, []definition.XAdESPath{&definition.XAdES132Path{}})
}

// NewDOMDocument is the constructor with provided XAdES Path holders. The method
// instantiates a new XAdES Path list based on the provided one. Ports the two-argument
// DOMDocument(Document, List<XAdESPath>) constructor.
func NewDOMDocument(document *xmldom.Node, xadesPathsHolders []definition.XAdESPath) *DOMDocument {
	return NewDOMDocumentWithName(document, "", xadesPathsHolders)
}

// NewDOMDocumentWithName is the constructor with provided XAdES Path holders and document
// name. The method instantiates a new XAdES Path list based on the provided one. Panics with the
// Java message when xadesPathsHolders is nil (Objects.requireNonNull upstream). Ports the
// three-argument DOMDocument(Document, String, List<XAdESPath>) constructor.
func NewDOMDocumentWithName(document *xmldom.Node, name string, xadesPathsHolders []definition.XAdESPath) *DOMDocument {
	if xadesPathsHolders == nil {
		panic("XAdES Path holders cannot be null!")
	}
	holders := make([]definition.XAdESPath, len(xadesPathsHolders))
	copy(holders, xadesPathsHolders)
	return &DOMDocument{
		DOMDocument:       xmlutils.NewDOMDocumentWithName(document, name),
		xadesPathsHolders: holders,
	}
}

// Document gets the DOM Document. Ports getDocument().
func (d *DOMDocument) Document() *xmldom.Node {
	return d.Node()
}

// XAdESPathHolders gets a list of registered XAdES Path holders. Ports getXAdESPathHolders().
func (d *DOMDocument) XAdESPathHolders() []definition.XAdESPath {
	return d.xadesPathsHolders
}

// AddXAdESPathHolder appends p to the registered XAdES Path holders. Go accommodation for
// Java's `getXAdESPathHolders().add(p)` call sites (XAdESSignature.registerXAdESPaths,
// XMLDocumentAnalyzer), which rely on the returned List being the same mutable field instance -
// a guarantee a `[]definition.XAdESPath` value returned by XAdESPathHolders cannot offer in Go.
// There is no such method on the Java class itself.
func (d *DOMDocument) AddXAdESPathHolder(p definition.XAdESPath) {
	d.xadesPathsHolders = append(d.xadesPathsHolders, p)
}

// ClearXAdESPathHolders removes all elements from the registered XAdES Path holders. Go
// accommodation for Java's `getXAdESPathHolders().clear()` call site
// (XMLDocumentAnalyzer.clearQueryHolders), the same mutable-list-instance accommodation
// AddXAdESPathHolder documents above. There is no such method on the Java class itself.
func (d *DOMDocument) ClearXAdESPathHolders() {
	d.xadesPathsHolders = nil
}

// SignatureNodes gets a node list containing all ds:Signature elements present within the
// document, with exception to counter signatures. Ports getSignatureNodes().
func (d *DOMDocument) SignatureNodes() []*xmldom.Node {
	if !d.signatureNodesSet {
		nodes, err := DSSXMLUtilsGetAllSignaturesExceptCounterSignatures(d.Document())
		if err == nil {
			d.signatureNodes = nodes
		}
		d.signatureNodesSet = true
	}
	return d.signatureNodes
}

// RecursiveIdBrowse indexes ID attributes over the whole document: an ID attribute can only be
// dereferenced if it is declared in the validation context, since the attribute carries no
// attached type information; this process attaches it without needing a DTD/XML schema. See the
// file header for why this delegates to *xmldom.Node#RegisterIDs rather than re-implementing the
// walk. Ports recursiveIdBrowse().
func (d *DOMDocument) RecursiveIdBrowse() {
	d.Document().RegisterIDs()
}
