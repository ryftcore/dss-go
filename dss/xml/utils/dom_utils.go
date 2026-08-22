// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/DomUtils.java (DSS 6.5.RC1).
//
// NAMING: see xpath_utils.go's file header for why every exported function here is prefixed
// DomUtils (mirrors package spi's DSSUtilsXxx/DSSASN1UtilsXxx precedent for merged Java
// "Utils" classes). Java overloads (buildDOM(), buildDOM(String), buildDOM(byte[]),
// buildDOM(InputStream), buildDOM(DSSDocument)) become distinctly-suffixed Go funcs, since Go
// has no overloading.
//
// SKIPPED (no Go equivalent, flagged for integrator): getSecureDocumentBuilderFactory,
// getSecureTransformerFactory and getSecureTransformer return javax.xml.parsers /
// javax.xml.transform factory/transformer objects. internal/xmldom.Parse already applies the
// equivalent secure posture internally (no DOCTYPE, no external entities) with no factory
// object to obtain first, and internal/xmldom.Serialize has no pluggable Transformer to
// configure either. There is nothing for a Go port of these three methods to usefully
// return, so they are omitted rather than stubbed.
package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xpath10"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// Constants mirroring DomUtils's private static fields.
const (
	domUtilsHash         = "#"
	domUtilsXNSOpen      = "xmlns("
	domUtilsXPOpen       = "xpointer("
	domUtilsXPWithIDOpen = "#xpointer(id("
	domUtilsXPRoot       = "#xpointer(/)"
	domUtilsIDAttribute  = "Id"
)

// domUtilsXmlPreamble is the starting binaries of an XML file.
var domUtilsXmlPreamble = []byte{'<'}

// domUtilsXmlWithBomPreamble is the starting binaries of an XML file with a UTF-8 BOM. Java
// writes this as signed bytes {-17, -69, -65, '<'}, i.e. the unsigned bytes 0xEF, 0xBB, 0xBF
// (the UTF-8 BOM) followed by '<'.
var domUtilsXmlWithBomPreamble = []byte{0xEF, 0xBB, 0xBF, '<'}

// domUtilsWhitespaceSplit ports Java's String.split("\\s") for IsXPointerQuery.
var domUtilsWhitespaceSplit = regexp.MustCompile(`\s`)

// DomUtilsRegisterNamespace registers a namespace and associated prefix. Deprecated: use
// XPathUtilsRegisterNamespace. Ports registerNamespace(DSSNamespace).
func DomUtilsRegisterNamespace(namespace *common.DSSNamespace) bool {
	return XPathUtilsRegisterNamespace(namespace)
}

// DomUtilsStartsWithXmlPreamble checks if byteArray starts with an XML preamble '<',
// with or without a BOM. NOTE: does not check XML-conformity of the whole file; call
// DomUtilsIsDOMBytes for a deep check. Ports startsWithXmlPreamble(byte[]).
func DomUtilsStartsWithXmlPreamble(byteArray []byte) bool {
	return utils.StartsWith(byteArray, domUtilsXmlPreamble) || utils.StartsWith(byteArray, domUtilsXmlWithBomPreamble)
}

// DomUtilsStartsWithXmlPreambleDocument checks if document starts with an XML preamble '<',
// with or without a BOM. NOTE: does not check XML-conformity of the whole file; call
// DomUtilsIsDOM for a deep check. Ports startsWithXmlPreamble(DSSDocument).
func DomUtilsStartsWithXmlPreambleDocument(document model.DSSDocument) (bool, error) {
	ok, err := domUtilsStreamStartsWith(document, domUtilsXmlPreamble)
	if err != nil {
		return false, model.NewDSSErrorMessageCause("Cannot read a sequence of bytes from the InputStream.", err)
	}
	if ok {
		return true, nil
	}
	ok, err = domUtilsStreamStartsWith(document, domUtilsXmlWithBomPreamble)
	if err != nil {
		return false, model.NewDSSErrorMessageCause("Cannot read a sequence of bytes from the InputStream.", err)
	}
	return ok, nil
}

// domUtilsStreamStartsWith ports the private startsWith(InputStream, byte[]) helper.
func domUtilsStreamStartsWith(document model.DSSDocument, preamble []byte) (bool, error) {
	r, err := document.OpenStream()
	if err != nil {
		return false, err
	}
	defer r.Close()
	return utils.StartsWithStream(r, preamble)
}

// DomUtilsBuildDOMEmpty creates a new empty Document. Ports buildDOM().
func DomUtilsBuildDOMEmpty() *xmldom.Node {
	return xmldom.NewDocument()
}

// DomUtilsBuildDOMFromString returns the Document created based on the XML string. Ports
// buildDOM(String).
func DomUtilsBuildDOMFromString(xmlString string) (*xmldom.Node, error) {
	return DomUtilsBuildDOMFromBytes([]byte(xmlString))
}

// DomUtilsBuildDOMFromBytes returns the Document created based on the byte array. Panics if
// data is nil (Java Objects.requireNonNull("bytes is required")). Ports buildDOM(byte[]).
func DomUtilsBuildDOMFromBytes(data []byte) (*xmldom.Node, error) {
	if data == nil {
		panic("bytes is required")
	}
	doc, err := xmldom.Parse(data, nil)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to parse content (XML expected) : %s", err.Error()), err)
	}
	return doc, nil
}

// DomUtilsBuildDOMFromReader returns the Document created based on the XML io.Reader. NOTE:
// closes r after reading. Ports buildDOM(InputStream).
func DomUtilsBuildDOMFromReader(r io.Reader) (*xmldom.Node, error) {
	data, err := io.ReadAll(r)
	if closer, ok := r.(io.Closer); ok {
		_ = closer.Close()
	}
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("An error occurred while reading InputStream : %s", err.Error()), err)
	}
	doc, err := xmldom.Parse(data, nil)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to parse content (XML expected) : %s", err.Error()), err)
	}
	return doc, nil
}

// DomUtilsBuildDOMFromDocument returns the Document created based on the model.DSSDocument.
// Panics if document is nil (Java Objects.requireNonNull("The document is null")). Ports
// buildDOM(DSSDocument).
func DomUtilsBuildDOMFromDocument(document model.DSSDocument) (*xmldom.Node, error) {
	if document == nil {
		panic("The document is null")
	}
	r, err := document.OpenStream()
	if err != nil {
		return nil, err
	}
	return DomUtilsBuildDOMFromReader(r)
}

// DomUtilsIsDOMBytes returns true if the binaries contain a Document. Ports isDOM(byte[]).
func DomUtilsIsDOMBytes(data []byte) bool {
	if !DomUtilsStartsWithXmlPreamble(data) {
		return false
	}
	_, err := DomUtilsBuildDOMFromBytes(data)
	return err == nil
}

// DomUtilsIsDOM returns true if the provided document is a valid XML Document. Any panic
// (mirroring Java's broad catch(Exception e)) is treated as "not a DOM". Ports
// isDOM(DSSDocument).
func DomUtilsIsDOM(document model.DSSDocument) (result bool) {
	defer func() {
		if recover() != nil {
			result = false
		}
	}()
	ok, err := DomUtilsStartsWithXmlPreambleDocument(document)
	if err != nil || !ok {
		return false
	}
	_, err = DomUtilsBuildDOMFromDocument(document)
	return err == nil
}

// DomUtilsSetAttributeNS adds an attribute with the namespace and the value. Ports
// setAttributeNS(Element, DSSNamespace, DSSAttribute, String).
func DomUtilsSetAttributeNS(element *xmldom.Node, namespace *common.DSSNamespace, attribute common.DSSAttribute, value string) {
	// Java always prepends "prefix:" unconditionally here (unlike DomUtilsCreateElementNS,
	// which guards on a non-empty prefix); xmldom.Name.QName() only emits "prefix:local" when
	// Prefix != "", so an empty-prefix DSSNamespace (never used by any caller in practice)
	// cannot reproduce Java's malformed ":AttrName" qualified name here - a harmless,
	// structurally-forced deviation.
	name := xmldom.Name{Space: namespace.Uri(), Local: attribute.AttributeName(), Prefix: namespace.Prefix()}
	element.SetAttr(name, value)
}

// DomUtilsAddElement creates and adds a new XML Element. Ports
// addElement(Document, Element, DSSNamespace, DSSElement).
func DomUtilsAddElement(document, parentDom *xmldom.Node, namespace *common.DSSNamespace, element common.DSSElement) *xmldom.Node {
	dom := DomUtilsCreateElementNS(document, namespace, element)
	parentDom.AppendChild(dom)
	return dom
}

// DomUtilsAdoptChildren adopts all children of toBeAdopted, excluding the node itself, into
// parentElement. Ports adoptChildren(Element, Node).
func DomUtilsAdoptChildren(parentElement, toBeAdopted *xmldom.Node) []*xmldom.Node {
	var adopted []*xmldom.Node
	for c := toBeAdopted.FirstChild; c != nil; c = c.NextSibling {
		imported := parentElement.Import(c, true)
		parentElement.AppendChild(imported)
		adopted = append(adopted, imported)
	}
	return adopted
}

// DomUtilsCreateXPathExpression creates a new instance of *xpath10.Expr with the given xpath
// expression. Deprecated: use XPathUtils methods. Ports createXPathExpression(String).
func DomUtilsCreateXPathExpression(xpathString string) (*xpath10.Expr, error) {
	expr, err := xpath10.Compile(xpathString, namespaceContextMapToXPath10(XPathUtilsGetNamespaceContextMap()))
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to create an XPath expression : %s", err.Error()), err)
	}
	return expr, nil
}

// DomUtilsGetValue returns the String value corresponding to the XPath query, evaluated
// directly as an XPath STRING result (first node's string value in document order, no
// "more than one result" check). Deprecated: use XPathUtilsGetValue. Ports
// getValue(Node, String).
//
// RECONCILED against the finished internal/xpath10 (its Status handoff: green, gate clean,
// 96.7% coverage): the package's public surface is intentionally Expr.Evaluate (a node-set
// result) plus Compile/Select/SelectOne - no separate STRING/BOOLEAN evaluation mode, since
// DSS never performs one (see internal/xpath10/doc.go). XPath 1.0's string(node-set)
// conversion is exactly "the string-value of the first node in the node-set, in document
// order, or \"\" if the node-set is empty" (clause 4.2), which is what javax.xml.xpath's
// XPathConstants.STRING evaluation performs internally for a node-set-typed expression - so
// evaluating as a node-set and taking the first result's TextContent() (itself XPath
// string-value for every node kind, per internal/xmldom's TextContent doc comment) reproduces
// getValue's semantics exactly, with no dedicated string-evaluation API needed.
func DomUtilsGetValue(xmlNode *xmldom.Node, xPathString string) (string, error) {
	expr, err := DomUtilsCreateXPathExpression(xPathString)
	if err != nil {
		return "", err
	}
	nodes, err := expr.Evaluate(xmlNode)
	if err != nil {
		return "", model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to extract value of the node. Reason : %s", err.Error()), err)
	}
	if len(nodes) == 0 {
		return "", nil
	}
	return utils.Trim(nodes[0].TextContent()), nil
}

// DomUtilsGetNodeList returns the nodes corresponding to the XPath query. Deprecated: use
// XPathUtilsGetNodeList. Ports getNodeList(Node, String).
func DomUtilsGetNodeList(xmlNode *xmldom.Node, xPathString string) ([]*xmldom.Node, error) {
	expr, err := DomUtilsCreateXPathExpression(xPathString)
	if err != nil {
		return nil, err
	}
	nodes, err := expr.Evaluate(xmlNode)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to find a NodeList by the given xPathString '%s'. Reason : %s", xPathString, err.Error()), err)
	}
	return nodes, nil
}

// DomUtilsGetNode returns the Node corresponding to the XPath query, or nil if none match.
// Deprecated: use XPathUtilsGetNode. Ports getNode(Node, String).
func DomUtilsGetNode(xmlNode *xmldom.Node, xPathString string) (*xmldom.Node, error) {
	list, err := DomUtilsGetNodeList(xmlNode, xPathString)
	if err != nil {
		return nil, err
	}
	if len(list) > 1 {
		return nil, model.NewDSSError("More than one result for XPath: " + xPathString)
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

// DomUtilsGetElement returns the Element corresponding to the XPath query. Deprecated: use
// XPathUtilsGetElement. Ports getElement(Node, String).
func DomUtilsGetElement(xmlNode *xmldom.Node, xPathString string) (*xmldom.Node, error) {
	return DomUtilsGetNode(xmlNode, xPathString)
}

// DomUtilsIsNotEmpty returns true if the xpath query contains something. Deprecated. Ports
// isNotEmpty(Node, String).
func DomUtilsIsNotEmpty(xmlNode *xmldom.Node, xPathString string) (bool, error) {
	n, err := DomUtilsGetNodesAmount(xmlNode, xPathString+"/child::node()[not(self::text())]")
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// DomUtilsGetNodesAmount returns the amount of nodes matching xPathString. Deprecated: use
// XPathUtilsGetNodesAmount. Ports getNodesAmount(Node, String).
func DomUtilsGetNodesAmount(xmlNode *xmldom.Node, xPathString string) (int, error) {
	list, err := DomUtilsGetNodeList(xmlNode, xPathString)
	if err != nil {
		return 0, err
	}
	return len(list), nil
}

// DomUtilsAddTextElement creates and adds a new XML Element with a text value. Ports
// addTextElement(Document, Element, DSSNamespace, DSSElement, String).
func DomUtilsAddTextElement(document, parentDom *xmldom.Node, namespace *common.DSSNamespace, element common.DSSElement, value string) *xmldom.Node {
	dom := DomUtilsCreateElementNS(document, namespace, element)
	parentDom.AppendChild(dom)
	valueNode := xmldom.NewText(value)
	dom.AppendChild(valueNode)
	return dom
}

// DomUtilsSetTextNode sets a text node to the given DOM element. Ports
// setTextNode(Document, Element, String).
func DomUtilsSetTextNode(document, parentDom *xmldom.Node, text string) {
	textNode := xmldom.NewText(text)
	parentDom.AppendChild(textNode)
}

// DomUtilsCreateXMLGregorianCalendar converts date to its canonical xsd:dateTime lexical
// representation (UTC, no fractional seconds - e.g. "2013-11-23T11:22:52Z"), i.e. the
// XMLGregorianCalendar.toXMLFormat() string every upstream call site immediately extracts
// (grep over dss-xades confirms createXMLGregorianCalendar(date).toXMLFormat() is the only
// use pattern). Returns "" for the zero Time (Java null). Ports
// createXMLGregorianCalendar(Date), folding in the always-following .toXMLFormat() call since
// Go has no XMLGregorianCalendar type to return instead.
func DomUtilsCreateXMLGregorianCalendar(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.UTC().Format("2006-01-02T15:04:05Z")
}

// DomUtilsGetDate converts text (an XML xsd:dateTime representation) to a time.Time. Returns
// the zero Time on failure (Java: LOG.warn then return null; logging dropped per PORTING.md).
//
// ASSUMPTION (flagged for integrator): upstream delegates
// to javax.xml.datatype.DatatypeFactory#newXMLGregorianCalendar(String), which lenently
// accepts any of the eight W3C XML Schema date/time lexical forms. This port only handles
// xsd:dateTime, with or without a fractional-second component, with a "Z" or "+HH:MM"/"-HH:MM"
// offset (time.RFC3339's layout already accepts an optional fractional-second field per the
// time package's documented parsing behaviour), or with no offset at all (treated as UTC,
// since XMLGregorianCalendar's own no-offset handling is JVM-default-timezone-dependent and
// therefore not a fixed target to match). This is every form DSS itself ever produces via
// DomUtilsCreateXMLGregorianCalendar and every form actually exercised by
// XAdESRevocationRefExtractionUtils/XAdESSignature callers.
//
// The assumption above has been run against a DatatypeFactory oracle. It holds for
// xsd:dateTime, with one correction applied here and one gap left open:
//
//   - CORRECTED: time.Parse accepts a COMMA as the fractional-second separator, so
//     "2024-01-15T10:30:00,5Z" parsed to 10:30:00.5 while DatatypeFactory rejects it - the
//     port was accepting a lexical form Java refuses. XSD allows only ".", and no comma can
//     occur anywhere else in any of the eight lexical forms, so one rejection covers it.
//   - OPEN, reported not fixed: the seven non-dateTime forms (xsd:date "2024-01-15",
//     gYearMonth, gYear, time, gDay, gMonth, gMonthDay) parse in Java and return the zero
//     Time here, as do "24:00:00" hour rollover, leap second ":60", negative and >4-digit
//     years. All are schema-invalid at every DSS call site (SigningTime, ProducedAt,
//     IssueTime are xsd:dateTime), and Java's answer for each is JVM-default-timezone
//     dependent - the same reason the no-offset case above is not matched either - so there
//     is no fixed oracle to port. Closing it is a deliberate decision, not a transcription.
func DomUtilsGetDate(text string) time.Time {
	if strings.ContainsRune(text, ',') {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339, text); err == nil {
		return t
	}
	if t, err := time.Parse("2006-01-02T15:04:05", text); err == nil {
		return t.UTC()
	}
	return time.Time{}
}

// DomUtilsGetChildrenNames returns the list of children's names for the element matched by
// xPathString. Deprecated: use XPathUtilsGetChildrenNames. Ports
// getChildrenNames(Node, String), which - unlike the non-deprecated XPathUtils version -
// unconditionally appends every child's local name, including "" for Text/CDATA/Comment
// children (Java's Node.getLocalName() returns null there; List.add(null) is legal in Java
// and has no error-free Go analogue, so "" stands in for null here).
func DomUtilsGetChildrenNames(xmlNode *xmldom.Node, xPathString string) ([]string, error) {
	var childrenNames []string
	element, err := DomUtilsGetElement(xmlNode, xPathString)
	if err != nil {
		return nil, err
	}
	if element != nil {
		for c := element.FirstChild; c != nil; c = c.NextSibling {
			if c.Kind != xmldom.Element {
				// Node.getLocalName() is non-null only for Element and Attr; every other
				// node kind reports null, which this port spells "". A ProcInst keeps its
				// target in Name.Local, so reading Name.Local unconditionally would report
				// the target where Java reports null - the only kind where the two differ,
				// since Text, CDATA and Comment all carry the zero Name already.
				childrenNames = append(childrenNames, "")
				continue
			}
			childrenNames = append(childrenNames, c.Name.Local)
		}
	}
	return childrenNames, nil
}

// DomUtilsWriteDocumentTo writes the Document content to w. Ports
// writeDocumentTo(Document, OutputStream).
//
// The bytes are byte-for-byte the identity Transformer's, the source document's own
// declared encoding included: internal/xmldom keeps the XML declaration's encoding and
// standalone on the Document and Serialize re-emits both, and writes the markup in that
// encoding. TestSerializeAgainstJavaTransformerOracle is the gate.
func DomUtilsWriteDocumentTo(dom *xmldom.Node, w io.Writer) error {
	if err := dom.Serialize(w, nil); err != nil {
		return model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to store a DOM document to OutputStream : %s", err.Error()), err)
	}
	return nil
}

// DomUtilsCreateDssDocumentFromDomDocument creates a new model.DSSDocument with the
// Document's content and the given name. Ports
// createDssDocumentFromDomDocument(Document, String).
func DomUtilsCreateDssDocumentFromDomDocument(document *xmldom.Node, name string) (model.DSSDocument, error) {
	var buf bytes.Buffer
	if err := DomUtilsWriteDocumentTo(document, &buf); err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to create a DSSDocument from DOM document : %s", err.Error()), err)
	}
	return model.NewInMemoryDocumentWithMimeType(buf.Bytes(), name, enumerations.MimeTypeEnumXML), nil
}

// DomUtilsXmlToString converts an XML Node to a string. Ports xmlToString(Node).
func DomUtilsXmlToString(node *xmldom.Node) (string, error) {
	b, err := DomUtilsSerializeNode(node)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// DomUtilsSerializeNode performs the serialization of the given node. See
// DomUtilsWriteDocumentTo's deviation note regarding the xmlEncoding echo. Ports
// serializeNode(Node) (the public byte[]-returning overload; the private
// serializeNode(Node, Result) helper is folded into this and into DomUtilsWriteDocumentTo).
func DomUtilsSerializeNode(xmlNode *xmldom.Node) ([]byte, error) {
	b, err := xmlNode.Bytes(nil)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("An error occurred during a node serialization.", err)
	}
	return b, nil
}

// DomUtilsGetId gets the Id value from the given URI reference, e.g. "#signature" ->
// "signature". Ports getId(String).
func DomUtilsGetId(uri string) string {
	id := uri
	if DomUtilsStartsFromHash(uri) {
		if DomUtilsIsXPointerQuery(uri) {
			if xpointerID, ok := DomUtilsGetXPointerId(uri); ok {
				id = xpointerID
			}
		} else {
			id = id[1:]
		}
	}
	return id
}

// DomUtilsGetXPathByIdAttribute returns a case-insensitive xPath expression matching an Id
// attribute equal to the id extracted from uri. Deprecated. Ports
// getXPathByIdAttribute(String).
func DomUtilsGetXPathByIdAttribute(uri string) string {
	id := DomUtilsGetId(uri)
	return "[@*[local-name()='Id']='" + id + "' or @*[local-name()='id']='" + id + "' or @*[local-name()='ID']='" + id + "']"
}

// DomUtilsGetElementById extracts an element from the given document node with the given Id,
// namespace independently. Deprecated: use XPathUtilsGetElementById. Ports
// getElementById(Node, String).
func DomUtilsGetElementById(node *xmldom.Node, id string) *xmldom.Node {
	return XPathUtilsGetElementById(node, id)
}

// DomUtilsStartsFromHash returns true if uri starts with the "#" character. Ports
// startsFromHash(String).
func DomUtilsStartsFromHash(uri string) bool {
	return utils.IsStringNotBlank(uri) && strings.HasPrefix(uri, domUtilsHash)
}

// DomUtilsIsElementReference returns true if uri refers to an element in the signature. Ports
// isElementReference(String).
func DomUtilsIsElementReference(uri string) bool {
	return DomUtilsStartsFromHash(uri) && !DomUtilsIsXPointerQuery(uri)
}

// DomUtilsToElementReference translates uri to a local element reference. E.g.: "r-123-id" ->
// "#r-123-id"; "sample.xml" -> "#sample.xml"; "#r-xades-enveloped" -> "#r-xades-enveloped".
// Ports toElementReference(String).
func DomUtilsToElementReference(uri string) string {
	if !DomUtilsStartsFromHash(uri) {
		uri = domUtilsHash + uri
	}
	return uri
}

// DomUtilsIsXPointerQuery indicates if uriValue is an XPointer query. Ports
// isXPointerQuery(String).
func DomUtilsIsXPointerQuery(uriValue string) bool {
	if utils.IsStringBlank(uriValue) {
		return false
	}
	uri := domUtilsDecodeUrlSilently(uriValue)
	if DomUtilsStartsFromHash(uri) {
		uri = uri[1:]
	}
	parts := domUtilsJavaSplit(uri)
	if len(parts) == 0 {
		// Java reaches parts[0] on an empty array and throws ArrayIndexOutOfBoundsException
		// (reachable: uriValue "#  " is not blank, so the guard above lets it through, and
		// splitting "  " on whitespace yields nothing but trailing empties). An uncaught
		// runtime exception on a Reference/@URI an attacker controls is not behaviour worth
		// reproducing as a panic; "not an XPointer query" is the answer every DSS caller of
		// this predicate would have wanted. See DomUtilsGetXPointerId for the same choice.
		return false
	}
	ii := 0
	for ; ii < len(parts)-1; ii++ {
		if !strings.HasSuffix(parts[ii], ")") || !strings.HasPrefix(parts[ii], domUtilsXNSOpen) {
			return false
		}
	}
	if !strings.HasSuffix(parts[ii], ")") || !strings.HasPrefix(parts[ii], domUtilsXPOpen) {
		return false
	}
	return true
}

// domUtilsJavaSplit splits s on whitespace with java.lang.String#split(String) semantics.
//
// regexp.Split(s, -1) is NOT that function: Java's split with the default limit of zero
// removes trailing empty strings from the result, Go's keeps them. The difference is
// load-bearing here - "#xpointer(id('x')) " with one trailing space splits to one part in
// Java and to two in Go, so the loop over the leading parts runs in Go and not in Java, and
// isXPointerQuery answers false where Java answers true. That flips isElementReference and
// changes getId, which is how a Reference/@URI is resolved to the element it signs.
// Empty strings BETWEEN parts are kept by both, so only the trailing run is trimmed.
func domUtilsJavaSplit(s string) []string {
	parts := domUtilsWhitespaceSplit.Split(s, -1)
	end := len(parts)
	for end > 0 && parts[end-1] == "" {
		end--
	}
	return parts[:end]
}

// domUtilsDecodeUrlSilently ports the private decodeUrlSilently(String) helper: URL-decodes
// uriValue, returning it unchanged on any error.
func domUtilsDecodeUrlSilently(uriValue string) string {
	decoded, err := url.QueryUnescape(uriValue)
	if err != nil {
		return uriValue
	}
	return decoded
}

// DomUtilsGetXPointerId extracts the xpointerId to search for from uri, and reports whether
// uri matched the "#xpointer(id(...))" form (Java: non-null return). Ports getXPointerId(String).
// See eu.europa.esig.dss.xml.utils.xpath's ResolverXPointer counterpart in Apache Santuario.
//
// DIVERGENCE, deliberate: Java indexes the delimited payload without checking that it holds
// both delimiters, so "#xpointer(id())" (empty payload) and "#xpointer(id(\"))" (a single
// quote character, read as both the opening and the closing delimiter) each throw an
// uncaught StringIndexOutOfBoundsException out of getXPointerId - and therefore out of
// getId, getXPathByIdAttribute and every Reference/@URI resolution above them. uri comes
// straight off a signature under validation, so this port answers "not an xpointer id"
// instead of crashing. A payload needs at least two characters to carry two delimiters,
// which is what idLen >= 1 requires.
func DomUtilsGetXPointerId(uri string) (string, bool) {
	if strings.HasPrefix(uri, domUtilsXPWithIDOpen) && strings.HasSuffix(uri, "))") {
		idPlusDelim := uri[len(domUtilsXPWithIDOpen) : len(uri)-2]
		idLen := len(idPlusDelim) - 1
		if idLen >= 1 &&
			((idPlusDelim[0] == '"' && idPlusDelim[idLen] == '"') ||
				(idPlusDelim[0] == '\'' && idPlusDelim[idLen] == '\'')) {
			return idPlusDelim[1:idLen], true
		}
	}
	return "", false
}

// DomUtilsIsRootXPointer checks if uri refers to the document root. Ports
// isRootXPointer(String).
func DomUtilsIsRootXPointer(uri string) bool {
	return domUtilsXPRoot == uri
}

// DomUtilsCreateElementNS creates an element with the given namespace. The documentDom
// parameter is unused: xmldom elements are not document-scoped until inserted into a tree
// (kept for signature fidelity with Java's Document#createElementNS receiver). Ports
// createElementNS(Document, DSSNamespace, DSSElement).
func DomUtilsCreateElementNS(documentDom *xmldom.Node, namespace *common.DSSNamespace, element common.DSSElement) *xmldom.Node {
	_ = documentDom
	name := xmldom.Name{Space: namespace.Uri(), Local: element.TagName()}
	if utils.IsStringNotEmpty(namespace.Prefix()) {
		name.Prefix = namespace.Prefix()
	}
	return xmldom.NewElement(name)
}

// DomUtilsAddNamespaceAttribute adds a namespace declaration attribute to element.
//
// JUDGMENT CALL (flagged for integrator): Java builds this via the non-namespace-aware
// Element#setAttribute("xmlns:"+prefix, uri), which - per the DOM Level 2 contract - creates
// an attribute node with a null namespaceURI, not one recognized as a real namespace
// declaration by namespace-aware processing (Java DOM's own semantics require
// setAttributeNS(XMLNS_URI, ...) for that). Since this method's only call sites
// (dss-xades: XAdESSignatureBuilder/XAdESBuilder) use it to declare the very ds:/xades:
// prefixes those signature elements are built with - i.e., it must functionally act as a
// real namespace declaration for serialization and canonicalization to be correct - this port
// creates a genuine xmldom namespace-declaration attribute (Space=XMLNSNamespace,
// Prefix="xmlns"), resolving the ambiguity in favour of runtime correctness rather than a
// literal quirk-for-quirk port. Reconcile against the XAdES KAT.
// Ports addNamespaceAttribute(Element, DSSNamespace).
func DomUtilsAddNamespaceAttribute(element *xmldom.Node, namespace *common.DSSNamespace) {
	name := xmldom.Name{Space: xmldom.XMLNSNamespace, Local: namespace.Prefix(), Prefix: "xmlns"}
	if namespace.Prefix() == "" {
		// An empty prefix has to become the DEFAULT declaration, xmlns="uri", which xmldom
		// spells {XMLNSNamespace, "xmlns", ""}. Concatenating "xmlns:" with "" the way Java
		// does and handing the result on produces the attribute name "xmlns:", whose empty
		// local part is not an XML Name: the element serialized to xmlns:="uri", markup no
		// parser will read back. Java never emits that either - Xerces' serializer applies
		// namespace fixup and writes xmlns="uri" - so the default declaration is both the
		// well-formed answer and the one that matches the Java oracle.
		name = xmldom.Name{Space: xmldom.XMLNSNamespace, Local: "xmlns"}
	}
	element.SetAttr(name, namespace.Uri())
}

// DomUtilsExcludeComments returns a Document with comments excluded. NOTE: the method
// modifies the original node in place before round-tripping it through serialization and
// re-parsing, matching Java's documented "workaround to handle the transforms correctly
// (clone does not work)". Ports excludeComments(Node).
func DomUtilsExcludeComments(node *xmldom.Node) (*xmldom.Node, error) {
	domUtilsExcludeCommentsRecursively(node)
	data, err := DomUtilsSerializeNode(node)
	if err != nil {
		return nil, err
	}
	return DomUtilsBuildDOMFromBytes(data)
}

// domUtilsExcludeCommentsRecursively ports the private excludeCommentsRecursively(Node)
// helper.
func domUtilsExcludeCommentsRecursively(node *xmldom.Node) {
	child := node.FirstChild
	for child != nil {
		next := child.NextSibling
		if child.Kind == xmldom.Comment {
			node.RemoveChild(child)
		} else if child.FirstChild != nil {
			domUtilsExcludeCommentsRecursively(child)
		}
		child = next
	}
}

// DomUtilsBrowseRecursivelyForNamespaceWithUri browses through element looking for a
// namespace with the target uri and returns the matching DSSNamespace if found, else nil.
// Ports browseRecursivelyForNamespaceWithUri(Element, String).
func DomUtilsBrowseRecursivelyForNamespaceWithUri(element *xmldom.Node, uri string) *common.DSSNamespace {
	if element == nil {
		return nil
	}
	if element.Name.Space == uri {
		return common.NewDSSNamespace(element.Name.Space, element.Name.Prefix)
	}
	for c := element.FirstChild; c != nil; c = c.NextSibling {
		if c.Kind == xmldom.Element {
			if ns := DomUtilsBrowseRecursivelyForNamespaceWithUri(c, uri); ns != nil {
				return ns
			}
		}
	}
	return nil
}

// DomUtilsGetNodeBytes returns the bytes of the given node: for Element/Document/Comment
// nodes, its serialization with a leading XML declaration stripped; for Text nodes, its
// base64-decoded content (falling back to the raw text bytes only in Java - see below); nil
// for anything else. Ports getNodeBytes(Node).
//
// NOTE: unlike DomUtilsSerializeNode, this only ever handles the Text node kind (not CDATA):
// Java's switch matches Node.TEXT_NODE, and CDATASection nodes report
// getNodeType()==CDATA_SECTION_NODE, so they fall to the unhandled default (nil) exactly as
// Attribute and ProcInst nodes do. The base64-decode-failure fallback
// (catch(Exception) -> textContent.getBytes()) is unreachable in this port: utils.FromBase64
// mirrors commons-codec's lenient decoder, which never errors (see utils/codec.go).
//
// The new String(bytes)/str.getBytes() pair around the declaration-stripping is NOT a
// no-op the way it looks: both use the platform default charset, which is UTF-8 since
// JEP 400, so a serialization written in any other encoding - which happens whenever the
// source document declared one - is decoded lossily and re-encoded, replacing every
// malformed byte sequence with U+FFFD. That is what the digest of a no-transform
// ds:Reference is then computed over, so it is reproduced here rather than skipped.
func DomUtilsGetNodeBytes(node *xmldom.Node) ([]byte, error) {
	switch node.Kind {
	case xmldom.Element, xmldom.Document, xmldom.Comment:
		b, err := DomUtilsSerializeNode(node)
		if err != nil {
			return nil, err
		}
		str := javaString(b)
		if strings.HasPrefix(str, "<?") {
			if idx := strings.Index(str, "?>"); idx >= 0 {
				str = str[idx+2:]
			}
		}
		return []byte(str), nil
	case xmldom.Text:
		return utils.FromBase64(node.Value), nil
	default:
		return nil, nil
	}
}

// javaString is new String(bytes) under the UTF-8 default charset: a decode whose
// malformed input is replaced with U+FFFD. The result is re-encoded to UTF-8 by the
// caller's []byte conversion, so what actually has to match Java is how many replacement
// characters a bad sequence collapses to, and that is decided by
// java.nio.charset.UTF_8$Decoder's malformed-input LENGTHS, not by "one per byte":
//
//   - a lead byte followed by too few continuation bytes consumes the whole partial
//     sequence and yields one U+FFFD (E9 A9 -> one, F0 9F 98 -> one);
//   - a lead byte whose SECOND byte is not a legal continuation for it consumes one byte
//     only, so the rest are re-examined and can produce their own (E0 80 80 -> three);
//   - a well-formed encoding of a surrogate or of a code point above U+10FFFF is rejected
//     as a whole (ED A0 80 -> one).
//
// Go's utf8.DecodeRune always reports a length of 1, which would give three replacements
// for E9 A9 21 where Java gives two. The rows named latin1-utf8-decoder-* in
// testdata/serialize pin every branch below.
func javaString(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	var sb strings.Builder
	sb.Grow(len(b))
	for i := 0; i < len(b); {
		r, size := utf8.DecodeRune(b[i:])
		if r != utf8.RuneError || size > 1 {
			sb.Write(b[i : i+size])
			i += size
			continue
		}
		sb.WriteRune(utf8.RuneError)
		i += javaMalformedLength(b[i:])
	}
	return sb.String()
}

// javaMalformedLength returns how many bytes UTF_8$Decoder consumes for the malformed
// sequence starting at b[0]. It mirrors that decoder's malformedN: the second byte decides
// whether the sequence is rejected one byte at a time or as a unit.
func javaMalformedLength(b []byte) int {
	b1 := b[0]
	switch {
	case b1&0xE0 == 0xC0: // 2-byte lead
		// Both malformed cases (an overlong C0/C1 lead, or a bad continuation) consume
		// one byte.
		return 1
	case b1&0xF0 == 0xE0: // 3-byte lead
		if len(b) < 2 || (b1 == 0xE0 && b[1]&0xE0 == 0x80) || !isContinuation(b[1]) {
			return 1
		}
		if len(b) < 3 || !isContinuation(b[2]) {
			return 2
		}
		// Well formed as far as UTF-8 goes; it is the surrogate range that rejects it.
		return 3
	case b1&0xF8 == 0xF0: // 4-byte lead
		if len(b) < 2 || b1 > 0xF4 ||
			(b1 == 0xF0 && (b[1] < 0x90 || b[1] > 0xBF)) ||
			(b1 == 0xF4 && b[1]&0xF0 != 0x80) || !isContinuation(b[1]) {
			return 1
		}
		if len(b) < 3 || !isContinuation(b[2]) {
			return 2
		}
		if len(b) < 4 || !isContinuation(b[3]) {
			return 3
		}
		return 4
	}
	// A continuation byte on its own, or a 0xF8-0xFF byte that leads nothing.
	return 1
}

func isContinuation(b byte) bool { return b&0xC0 == 0x80 }

// DomUtilsCreateDeepCopy creates a deep copy of the Document for the given element, and
// returns the corresponding element from the copied Document. This addresses a
// canonicalization issue on document copies, reported in SANTUARIO-139. Ports
// createDeepCopy(Element).
func DomUtilsCreateDeepCopy(element *xmldom.Node) *xmldom.Node {
	if element == nil {
		return nil
	}

	originalID := element.AttrValue("", domUtilsIDAttribute)
	blank := utils.IsStringBlank(originalID)
	var id string
	if blank {
		// A temporary identifier to retrieve the element from the document copy after; the
		// Id attribute is then removed from both the original element and the obtained copy.
		id = domUtilsRandomUUID()
	} else {
		id = originalID
	}
	element.SetAttr(xmldom.Name{Local: domUtilsIDAttribute}, id)

	var copyElement *xmldom.Node
	defer func() {
		if blank {
			domUtilsRemoveAttribute(element, domUtilsIDAttribute)
			domUtilsRemoveAttribute(copyElement, domUtilsIDAttribute)
		}
	}()

	originalRoot := element.Document().DocumentElement()
	documentCopy := xmldom.NewDocument()
	copiedRoot := documentCopy.Import(originalRoot, true)
	documentCopy.AppendChild(copiedRoot)

	copyElement = XPathUtilsGetElementById(documentCopy, id)
	return copyElement
}

// domUtilsRemoveAttribute ports the private removeAttribute(Element, String) helper.
func domUtilsRemoveAttribute(element *xmldom.Node, attributeName string) {
	if element != nil {
		element.RemoveAttr("", attributeName)
	}
}

// domUtilsRandomUUID returns a random RFC 4122 version 4 UUID string, standing in for
// java.util.UUID.randomUUID() (stdlib-only per PORTING.md's dependency policy), matching the
// same approach already used in token/pkcs11_signature_token.go.
func domUtilsRandomUUID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// DomUtilsGetCurrentNamespaces returns the stored namespace definitions. Deprecated: use
// XPathUtilsGetNamespaceContextMap().PrefixMap(). Ports getCurrentNamespaces().
func DomUtilsGetCurrentNamespaces() map[string]string {
	return XPathUtilsGetNamespaceContextMap().PrefixMap()
}
