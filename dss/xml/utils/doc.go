// Package utils ports dss-xml-utils (13 Java source files, root package plus the xpath
// subpackage) into a single Go package, as instructed for this phase. It is DSS's bridge
// between internal/xmldom + internal/xmlc14n + internal/xpath10 and the rest of DSS:
// DOMDocument adapts an *xmldom.Node to model.DSSDocument, DomUtils/XPathUtils are the query
// and construction surface XAdES (a later phase) builds on, and XMLCanonicalizer is a thin,
// stateless-per-call wrapper over internal/xmlc14n.
//
// Import as xmlutils from any package that also imports github.com/utain/esig/dss/utils
// (same bare package name "utils", different import path) to avoid a collision; within this
// package's own files no alias is needed.
//
// # Provenance
//
// DSS 6.5.RC1, package eu.europa.esig.dss.xml.utils and its subpackage xpath.
//
// # Consolidated deviations and assumptions (flagged for integrator)
//
// See each file's header for full detail; summarized here for a single point of review:
//
//   - GetSecureDocumentBuilderFactory/GetSecureTransformerFactory/GetSecureTransformer are
//     omitted: internal/xmldom.Parse and internal/xmldom.Serialize already apply the
//     equivalent secure posture internally, with no factory/transformer object to return.
//   - DomUtilsWriteDocumentTo/DomUtilsSerializeNode/DOMDocument's serialized bytes do not
//     echo the source document's original xml declaration encoding, because internal/xmldom
//     does not retain it after parsing (XML_DESIGN.md §1.5 discards the declaration once
//     read). XML_DESIGN.md §1.8 explicitly does not require byte-parity here.
//   - NativeDOMXPathQueryExecutor and JavaXmlXPathQueryExecutor share one xpath10-backed
//     implementation per this phase's task brief, even though xml/common/doc.go anticipated
//     NativeDOMXPathQueryExecutor being a from-scratch MatchNode-chain walker (that logic
//     remains fully available, unused, in package xml/common). See
//     native_dom_xpath_query_executor.go's header for the full reasoning; this is the single
//     most consequential design decision in this port and should be reconciled explicitly.
//   - XPathQueryExecutorLoader has no ServiceLoader/plugin-discovery equivalent; it is
//     hardwired to JavaXmlXPathQueryExecutor, matching upstream's practical (not merely
//     theoretical) default.
//   - DomUtilsGetDate's xsd:dateTime parsing (RFC3339, optionally with fractional seconds,
//     falling back to a no-offset/UTC form) has not been checked against a Java
//     DatatypeFactory oracle; flagged as the one XML_DESIGN.md "must be exact" item this
//     package could not KAT-verify without dss-xades fixtures.
//   - DomUtilsAddNamespaceAttribute resolves an ambiguity in Java's own
//     Element#setAttribute("xmlns:"+prefix, uri) (which, per the DOM Level 2 contract,
//     creates an attribute with a null namespaceURI - not a namespace-aware declaration) in
//     favour of creating a genuine xmldom namespace-declaration attribute, since its call
//     sites need it to function as one. See dom_utils.go's header comment on the function.
//   - NamespaceContextMap.Prefixes returns an empty slice instead of reproducing Java's
//     NullPointerException for an unregistered URI (dead code path upstream); Prefix/Prefixes
//     pick deterministically (lexicographically smallest / sorted) rather than reproducing
//     Java's HashSet-iteration-order-dependent answer, matching the same choice
//     internal/xpath10.NamespaceContext already made for the same upstream method.
package utils
