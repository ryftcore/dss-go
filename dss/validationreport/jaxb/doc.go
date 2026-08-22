// Ported from specs-validation-report/target/generated-sources/xjc/eu/europa/esig/validationreport/jaxb
// (DSS 6.5.RC1), the JAXB model xjc generates from 1910202xmlSchema.xsd (ETSI TS 119 102-2).

// Package jaxb is the schema-shaped model of an ETSI TS 119 102-2 Validation
// Report.
//
// Every generated JAXB class keeps its name (ValidationReportType stays
// ValidationReportType) and its property order, and every property carries
// the xml struct tag that reproduces its JAXB annotation, so that
// marshalling with Marshal produces the bytes the JAXB reference
// implementation produces for the same tree. The classes are grouped into
// files by schema area rather than one file per class.
//
// Two generated files have no Go counterpart of their own, by design (see
// dss/diagnostic/jaxb's doc.go, which documents the same choice):
//
//   - package-info.java carries @XmlSchema(namespace=...,
//     elementFormDefault=QUALIFIED). There is no package-level annotation
//     in Go, so the namespace is the Namespace constant and the qualified
//     element form is what Marshal applies in xml.go.
//   - ObjectFactory.java is xjc's no-arg createXxx() factories plus the
//     @XmlElementDecl naming ValidationReport as the document element. Go
//     composite literals replace the factories, and the document element
//     is bound by ValidationReportType's own XMLName field, so porting the
//     class would add nothing callers can use.
//
// # Cross-namespace types
//
// 1910202xmlSchema.xsd imports the XMLDSig, XAdES 1.3.2 and ETSI TS 119 612
// (trusted-list) namespaces for a handful of properties: ds:Signature,
// ds:SignatureValue, ds:DigestMethod/DigestValue, XAdES's
// DigestAlgAndValueType and SignaturePolicyIdentifierType, and the
// trusted-list DigitalIdentityType/TSPInformationType. None of those
// modules' own generated-JAXB classes has been ported to Go (they belong to
// specs-xmldsig, specs-xades and specs-trusted-list), so this package
// cannot import a matching Go type for them. Two shapes are used instead, both documented in
// jaxb_crossns.go:
//
//   - The small, low-cardinality XMLDSig types this schema actually
//     dereferences (ds:DigestMethod, ds:DigestValue via
//     DigestAlgAndValueType, ds:SignatureValue) are modelled directly from
//     the well-known xmldsig-core schema, reusing
//     dss/internal/xmldsig.NamespaceDSig for the namespace literal's
//     provenance (verified equal in xml_kat_test.go; struct tags need a
//     compile-time literal, so the constant cannot be referenced from the
//     tag itself).
//   - The large ones this schema treats as opaque payloads from this
//     package's point of view - ds:Signature (SignatureType),
//     XAdES's SignaturePolicyIdentifierType, and the trusted-list
//     DigitalIdentityType/TSPInformationType - are modelled as raw-XML
//     capture stand-ins (innerxml in, innerxml out), which guarantees
//     marshal-parity byte-for-byte for whatever content a real report
//     carries there without requiring a full port of those schemas. If
//     specs-xmldsig-jaxb/specs-xades-jaxb/specs-trusted-list-jaxb are ever
//     ported, the real generated types should replace these stand-ins; the
//     field names and positions are already correct and would not move.
//
// # xs:anyType / xs:any content (AnyType, TypedDataType.Value, ...)
//
// The schema uses AnyType (mixed, lax xs:any) and bare Object-typed
// properties (VOReferenceType.Any, TypedDataType.Value,
// IndividualValidationConstraintReportType.Indications,
// ValidationObjectRepresentationType's "direct" choice member) as
// extensibility points nothing in this package interprets. They are
// captured verbatim with `xml:",innerxml"` (see AnyType and RawContent in
// jaxb_common.go) rather than modelled as a DOM tree: unmarshal captures
// the exact bytes and marshal re-emits them unchanged, which is a stronger
// parity guarantee than reconstructing the content from a parsed
// representation would be, and is sufficient for the marshal-parity KAT
// this package is pinned to.
//
// The parity contract is pinned by the round-trip KAT in xml_kat_test.go
// over testdata/oracle: validation-report dumps produced by upstream
// validation of real signed fixtures (testdata/gen/ValidationReportOracle.java).
// Each must unmarshal into this model and marshal back byte for byte. The
// XSD-completeness sweep in xml_schema_test.go checks the model against
// 1910202xmlSchema.xsd element by element and attribute by attribute.
package jaxb
