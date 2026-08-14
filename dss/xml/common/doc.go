// Package common ports dss-xml-common, flattening its six Java subpackages (root,
// definition, definition/xmldsig, exception, alert, xpath, xpath/item) into a single Go
// package, as instructed for this phase (collision-checked: snake_case-ing every one of the
// 37 upstream file names produces no clash, so no subpackage-name filename prefixing was
// needed).
//
// # What this package is
//
// It is DSS's definition/XPath framework: DSSNamespace/DSSElement/DSSAttribute describe an
// XML vocabulary (xmldsig's is the one concrete instance ported here, in
// xmldsig_element.go/xmldsig_attribute.go/xmldsig_namespace.go); AbstractPath and the
// xpath/xpath-item types build XPath expression strings and, given an *xmldom.Node, evaluate
// a single query-chain link against it without a full XPath engine (upstream: the "process"
// half of the API, consumed by dss-xml-utils's NativeDOMXPathQueryExecutor - a later phase -
// which is why this package already depends on internal/xmldom rather than leaving matchNode
// unported).
//
// # What this package is not
//
// It does not evaluate multi-step or predicate-general XPath (no axes beyond child/self, no
// arithmetic, no functions besides the fixed local-name() equality this framework itself
// generates) - upstream doesn't either; real XPath execution is javax.xml.xpath, wired in
// dss-xml-utils.
//
// # JAXP security plumbing
//
// AbstractConfigurator/AbstractFactoryBuilder and the four builders
// (DocumentBuilderFactoryBuilder/TransformerFactoryBuilder/SchemaFactoryBuilder/
// ValidatorConfigurator) port upstream's javax.xml.parsers/transform/validation
// hardening layer. Go has no JAXP: no DocumentBuilderFactory/TransformerFactory/
// SchemaFactory/Validator types, no XSLT engine, no XSD schema validator (internal/xmldom
// explicitly excludes schema validation and XSLT - see its doc.go "Non-goals"). Per this
// phase's brief, the two are handled differently:
//
//   - DocumentBuilderFactoryBuilder has a real Go counterpart: internal/xmldom.ParseOptions
//     already encodes the equivalent secure-parser posture as plain Go fields. Its file
//     documents the complete JAXP-feature-URI-to-ParseOptions-field mapping.
//   - TransformerFactoryBuilder, SchemaFactoryBuilder and ValidatorConfigurator (and the
//     Source/Schema/Validator/TransformerFactory/SchemaFactory placeholder types they
//     mention) are documented stubs: the public contract - type and method names, the
//     feature/attribute configuration machinery inherited from AbstractConfigurator - is
//     preserved so a later phase (XSD-validated XAdES manifests, if ever needed) can fill
//     in real behaviour without an API break, but no XSLT or XSD validation actually
//     happens yet. Each stub's doc comment says so explicitly.
//
// # Provenance
//
// DSS 6.5.RC1, package eu.europa.esig.dss.xml.common and its subpackages alert, definition,
// definition.xmldsig, exception, xpath, xpath.item.
package common
