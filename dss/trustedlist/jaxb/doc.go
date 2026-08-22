// Ported from specs-trusted-list/target/generated-sources/xjc/eu/europa/esig/trustedlist/jaxb
// (DSS 6.5.RC1), the JAXB model xjc generates from ts_119612v020401_xsd.xsd,
// ts_119612v020401_sie_xsd.xsd, ts_119612v020401_additionaltypes_xsd.xsd and
// mra_schema_v2_19612v020401.xsd, across the four Java packages
// eu.europa.esig.trustedlist.jaxb.{tsl,ecc,tslx,mra}.

// Package jaxb is the schema-shaped model of an ETSI TS 119 612 Trusted List
// / List of Trusted Lists (LOTL) document, plus its MRA (Mutual Recognition
// Agreement, ETSI TS 119 612 Annex B / the EC's mra_schema_v2) extension.
//
// Every generated JAXB class keeps its name (TrustStatusListType stays
// TrustStatusListType) and its exact property order, and every property
// carries the xml struct tag that reproduces its JAXB annotation, so that
// marshalling with Marshal/MarshalMRA reproduces the bytes the JAXB
// reference implementation produces for the same tree. The classes are
// grouped into files by schema area rather than one file per class:
//
//   - jaxb_common.go: the schema's shared simple types (Base64Binary,
//     BigInteger) and its xs:any/xs:anyType extensibility machinery
//     (AnyType, ExtensionType, the wildcard dispatch registry, and the
//     raw-capture stand-ins for the handful of properties whose Java type
//     is generated from a schema outside this manifest - see that file's
//     header).
//   - jaxb_tsl_common.go, jaxb_tsl_service.go, jaxb_tsl_tsp.go: the "tsl"
//     package (ts_119612v020401_xsd.xsd + its sie/additionaltypes-adjacent
//     leaf types), split by schema area - scheme-level types, TSP service
//     types, and TSP/scheme-information types respectively.
//   - jaxb_tsl_root.go: TrustStatusListType, the document element, plus the
//     Marshal/Unmarshal entry points and the root xmlns:* handling.
//   - jaxb_ecc.go: the "ecc" package (Qualified Certificate qualifiers
//     extension, ts_119612v020401_sie_xsd.xsd).
//   - jaxb_tslx.go: the "tslx" package (additional types extension,
//     ts_119612v020401_additionaltypes_xsd.xsd).
//   - jaxb_mra.go: the "mra" package (mra_schema_v2_19612v020401.xsd).
//   - jaxb_model.go: the modelTypes registry every other file's types are
//     listed in exactly once, and the self-closing-tag content model it
//     derives (see that file's header; same technique as
//     dss/diagnostic/jaxb/jaxb_content_model.go).
//
// Two generated files per Java package have no Go counterpart of their own,
// by design, for the same reasons dss/diagnostic/jaxb's doc.go gives:
// package-info.java (the namespace constants below replace its
// @XmlSchema) and ObjectFactory.java (composite literals replace its
// no-arg createXxx() factories; TrustStatusListType's XMLName field and
// Marshal/Unmarshal bind the document element ObjectFactory's
// @XmlElementDecl named). The four packages' XmlAdapter classes (ecc's
// Adapter1/Adapter2, mra's Adapter1/Adapter2/Adapter3) also have no
// counterpart: what they adapt (dss/enumerations' KeyUsageBit and Assert,
// and MRAStatus/MRAEquivalenceContext via dss/trustedlist's
// MRAStatusParser/MRAEquivalenceContextParser) is folded into the bound
// field's own MarshalText/UnmarshalText, the same way dss/diagnostic/jaxb
// folds CollapsedStringAdapter into CollapsedString.
//
// # Namespaces
//
// The four schemas' target namespaces, matching their package-info.java's
// @XmlSchema(namespace=...):
package jaxb

const (
	// NamespaceTSL is ts_119612v020401_xsd.xsd's namespace (package tsl);
	// the document's own, unprefixed, default namespace.
	NamespaceTSL = "http://uri.etsi.org/02231/v2#"
	// NamespaceECC is ts_119612v020401_sie_xsd.xsd's namespace (package ecc,
	// Qualified Certificate qualifiers).
	NamespaceECC = "http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/#"
	// NamespaceTSLX is ts_119612v020401_additionaltypes_xsd.xsd's namespace
	// (package tslx, "additional types").
	NamespaceTSLX = "http://uri.etsi.org/02231/v2/additionaltypes#"
	// NamespaceMRA is mra_schema_v2_19612v020401.xsd's namespace (package
	// mra).
	NamespaceMRA = "http://ec.europa.eu/tools/lotl/mra/schema/v2#"
	// NamespaceDSig is the XMLDSig namespace ds:Signature/ds:KeyValue
	// resolve to (specs-xmldsig, outside this manifest - see
	// jaxb_common.go's header).
	NamespaceDSig = "http://www.w3.org/2000/09/xmldsig#"
	// NamespaceXAdES132 is the XAdES 1.3.2 namespace (specs-xades, outside
	// this manifest - see jaxb_common.go's header).
	NamespaceXAdES132 = "http://uri.etsi.org/01903/v1.3.2#"
	// NamespaceXAdES141 is the XAdES 1.4.1 namespace: registered in every
	// TrustedListUtils/MRAUtils JAXBContext (so the RI always pre-declares
	// it at the document root, see jaxb_tsl_root.go) but never bound by any
	// property in this schema.
	NamespaceXAdES141 = "http://uri.etsi.org/01903/v1.4.1#"
)
