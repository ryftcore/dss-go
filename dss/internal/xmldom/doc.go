// Package xmldom is a minimal, namespace-aware XML document model sufficient for
// XML-DSig canonicalization, reference processing and XAdES construction. It replaces
// the org.w3c.dom API that upstream DSS obtains from the JDK. It is not a general
// purpose XML library: no DTDs, no schema validation, no XPath, no XSLT.
//
// # Provenance
//
// Upstream DSS gets its XML tree from a third party - the JDK's org.w3c.dom / Xerces -
// so, exactly as PORTING.md prescribes for BouncyCastle-replacement machinery, this
// package has no Java class to mirror one-to-one and therefore lives under internal/.
// It is shaped by three upstream sources, all of DSS 6.5.RC1 / Apache Santuario 3.0.6 /
// OpenJDK 21:
//
//   - eu.europa.esig.dss.xml.common.DocumentBuilderFactoryBuilder and
//     org.apache.xml.security.parser.XMLParserImpl, which is where the security posture
//     comes from: disallow-doctype-decl, no external entities, no entity expansion.
//   - eu.europa.esig.dss.xades.definition.XAdESDOMDocument.recursiveIdBrowse /
//     setIDIdentifier (and the deprecated DSSXMLUtils.setIDIdentifier), which is where
//     the ID-attribute rule in RegisterIDs comes from.
//   - org.apache.xml.security.c14n.implementations.CanonicalizerBase and AttrCompare,
//     which is why namespace declarations are modelled as ordinary attribute nodes and
//     why Name keeps the literal source prefix.
//
// # Relationship to encoding/xml
//
// The tokenizer is encoding/xml's, driven exclusively through Decoder.RawToken so that
// element and attribute names keep their literal prefixes; Decoder.Token would rewrite
// Name.Space to the resolved URI and destroy the prefix that canonicalization must emit
// verbatim. Namespace resolution, well-formedness checking beyond what RawToken does,
// and XML 1.0 clause 3.3.3 attribute-value normalization are this package's own.
//
// xml.Attr.Value is never read: encoding/xml decodes character and entity references in
// place, which destroys the literal-versus-reference distinction that clause 3.3.3
// depends on. Attribute values are recomputed from the raw start-tag bytes recovered by
// bracketing each RawToken call with Decoder.InputOffset. See attvalue.go.
//
// # Non-goals
//
// DTD and entity declarations, schema validation, XPath, XSLT, XInclude, namespace
// fixup on mutation, Document.normalize, live NodeLists, mutation events, xml:space
// interpretation and pretty-printing are all out of scope and will not be added.
//
// # Dependencies
//
// The standard library only. No DSS package is imported: the DSS packages are what
// import this one.
package xmldom
