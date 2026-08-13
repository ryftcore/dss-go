// Package xpath10 evaluates the subset of XPath 1.0 that DSS uses, over internal/xmldom
// trees. It replaces javax.xml.xpath, which upstream DSS obtains from the JDK. It is not a
// general purpose XPath engine: anything outside the subset is refused with an
// *UnsupportedError naming the construct, never approximated and never silently ignored.
//
// # Provenance
//
// Upstream DSS gets its XPath engine from a third party - the JDK's Xalan, reached through
// javax.xml.xpath in eu.europa.esig.dss.xml.utils.xpath.JavaXmlXPathQueryExecutor - so,
// exactly as PORTING.md prescribes for the BouncyCastle-replacement packages, this package has
// no Java class to mirror one-to-one and therefore lives under internal/. It imports the
// standard library and internal/xmldom, nothing else.
//
// Three upstream sources shape it, all of DSS 6.5.RC1:
//
//   - eu.europa.esig.dss.xml.common.xpath.XPathQueryBuilder and the XPathQueryItem /
//     XPathQueryParameter hierarchy, which is what actually writes DSS's XPath expressions.
//     Every expression it can produce is in testdata/expressions.txt.
//   - eu.europa.esig.dss.xml.utils.NamespaceContextMap, whose semantics NamespaceContext
//     reproduces - in particular that an unregistered prefix resolves to the null namespace
//     instead of raising.
//   - eu.europa.esig.dss.xml.utils.xpath.XPathUtils, whose getNodeList / getNode / getElement
//     contracts are mirrored by Select and SelectOne.
//
// # The subset
//
// Determined by testdata/expressions.txt, the deduplicated inventory of every expression the
// upstream tree can hand to that engine, and no wider: every construct below is exercised by
// at least one inventory expression against the fixtures, and every one of those answers is
// pinned node-for-node against Xalan by TestKnownAnswers.
//
//	construct                        written as              comes from
//	----------------------------------------------------------------------------------------
//	relative location path           ds:Signature            XPathQueryElementItem
//	self step abbreviation           ./ds:Object             FromCurrentPositionXPathQuery
//	absolute path                    /*                      setXPathLocationString
//	absolute descendant-or-self      //ds:Signature          AllXPathQuery
//	relative descendant-or-self      .//ds:Reference         AllFromCurrentPositionXPathQuery
//	"//" as a step separator         //*[..]//*[..]          setXPathLocationString
//	child axis, by name              child::node()           DomUtils.isNotEmpty
//	parent axis                      parent::xades132:Cs     XPathQueryNotChildOfParameter
//	self axis, by name               self::text()            DomUtils.isNotEmpty
//	attribute step                   @URI, @*                XPathQueryAttributeItem
//	element wildcard                 *                       XPathQueryAnyItem
//	node() and text() tests          node(), text()          DomUtils.isNotEmpty
//	predicate                        [ .. ]                  every parameter
//	equality                         @Id='x'                 XPathQueryAttributeParameter
//	disjunction                      .. or ..                XPathQueryIdentifierParameter
//	not()                            not(parent::x)          XPathQueryNotChildOfParameter
//	local-name()                     local-name()='Id'       XPathQueryAttributeParameter
//	string literal                   'Id'                    every parameter
//
// Everything else in XPath 1.0 is an *UnsupportedError that names it: the "and", "!=",
// relational, arithmetic and union operators, variable references, parenthesized and filter
// expressions, number literals and therefore positional predicates, the ".." abbreviation, the
// "prefix:*" node test, comment() and processing-instruction(), the ancestor, descendant,
// following, preceding, sibling and namespace axes, the "attribute::" and
// "descendant-or-self::" spellings of the two axes reachable only through "@" and "//", and
// every core function other than not() and local-name().
//
// Compile also refuses a top-level expression that is not a location path. DSS evaluates every
// expression it compiles as a NODESET, so Evaluate returning a node-set is the whole contract;
// a boolean or string evaluation would be API that no known answer covers.
//
// The XML-DSig transform expressions - XPathTransform, XPath2FilterTransform and the
// enveloped-signature transform's not(ancestor-or-self::ds:Signature) - are deliberately NOT
// in that inventory. Santuario evaluates them inside the Reference/Transform pipeline, never
// through JavaXmlXPathQueryExecutor, and XPath Filter 2.0 is a different evaluation model on
// top of a different engine. They belong to the xmldsig package. A package that needs them
// must extend this subset deliberately, with its own known answers; until then it gets an
// *UnsupportedError rather than a plausible guess.
//
// # Namespaces and node identity
//
// Prefixes are resolved once, at Compile time, against the NamespaceContext passed in; that is
// where javax.xml.xpath resolves them too. An unprefixed name is in no namespace, per XPath
// 1.0 clause 2.3, and never in the default namespace of the document - so "//placeOfSignature"
// matches only an element outside any namespace.
//
// A namespace declaration is an attribute node in internal/xmldom, mirroring org.w3c.dom, but
// it is not one in XPath: the attribute axis skips every xmlns declaration, and the namespace
// axis is not implemented. Results are node-sets in document order without duplicates, where
// an element's attributes sort after the element and before its children.
package xpath10
