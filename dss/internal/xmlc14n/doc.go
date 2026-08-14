// Package xmlc14n canonicalizes XML documents and document subsets per Canonical XML 1.0
// (http://www.w3.org/TR/2001/REC-xml-c14n-20010315), Canonical XML 1.1
// (http://www.w3.org/2006/12/xml-c14n11) and Exclusive XML Canonicalization 1.0
// (http://www.w3.org/2001/10/xml-exc-c14n#), each with and without comments, plus the
// Santuario physical method (http://santuario.apache.org/c14n/physical).
//
// # Provenance
//
// Upstream DSS obtains canonicalization from Apache Santuario (xmlsec) through
// eu.europa.esig.dss.xml.utils.XMLCanonicalizer, so there is no Java class to mirror
// one-to-one and, per PORTING.md, this machinery lives under internal/ instead of being
// duplicated per package. The algorithms are ports of Apache Santuario xmlsec 3.0.6:
//
//	CanonicalizerBase              -> walk.go, escape.go
//	Canonicalizer20010315          -> emit_inclusive.go
//	Canonicalizer11_*              -> emit_inclusive.go, xmlattrs.go, xmlbase.go
//	Canonicalizer20010315Excl      -> emit_exclusive.go
//	CanonicalizerPhysical          -> emit_physical.go
//	NameSpaceSymbTable             -> nsstack.go
//	XmlAttrStack                   -> xmlattrs.go, xmlbase.go
//	AttrCompare                    -> attrsort.go
//	C14nHelper, InclusiveNamespaces -> algorithm.go
//
// The registry, the two defaults and the call paths follow XMLCanonicalizer (DSS 6.5.RC1).
// This package imports only the standard library and internal/xmldom; it never imports a DSS
// package, because the DSS packages are what import it.
//
// # Bug-for-bug
//
// Interoperability with Java DSS is the contract, so Santuario behaviour is reproduced even
// where it is arguably wrong: the epilog-drop quirk for a childless document element (walk.go),
// the relative-namespace check firing only for declarations physically present on the element
// being emitted (emit_inclusive.go, emit_exclusive.go), Java's UTF-16 code-unit ordering of
// attribute sort keys (attrsort.go), and the whole C14N 1.1 xml:base surface (xmlbase.go) -
// java.net.URI's decode-then-requote round trip, Santuario's own removeDotSegments, and the
// two unchecked exceptions it raises on an opaque or authority-only base, which abort the
// canonicalization rather than being swallowed. Each site carries a comment and a
// known-answer test.
//
// # State
//
// Santuario's canonicalizers are single-use: Canonicalizer20010315.firstCall is never reset, so
// a reused instance silently drops inherited namespaces and xml:* attributes (SANTUARIO-463).
// Every DSS call site builds one inline and therefore always gets the first-call behaviour.
// This package is stateless per call instead: all mutable state is created inside Canonicalize
// and dies with it, so the functions here are safe for concurrent use.
//
// # Coverage
//
// All seven algorithms XMLCanonicalizer registers canonicalize: Canonical XML 1.0 and 1.1 and
// Exclusive XML Canonicalization 1.0, each with and without comments, plus the Santuario
// physical method (which has no separate #WithComments URI because it always keeps comments).
// Both document-subset forms are implemented; the physical method has none, matching
// CanonicalizerPhysical, and Canonicalize reports ErrPhysicalNodeSet for it.
package xmlc14n
