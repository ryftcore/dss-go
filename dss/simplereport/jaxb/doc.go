// Ported from dss-simple-report-jaxb/target/generated-sources/xjc/eu/europa/esig/dss/simplereport/jaxb
// (DSS 6.5.RC1), the JAXB model xjc generates from SimpleReport.xsd.

// Package jaxb is the schema-shaped model of a DSS simple report.
//
// Every generated JAXB class keeps its name (XmlSignature stays
// XmlSignature) and its property order, and every property carries the xml
// struct tag that reproduces its JAXB annotation, so that marshalling with
// Marshal produces the bytes the JAXB reference implementation produces for
// the same tree. The classes are grouped into files by schema area rather
// than one file per class.
//
// Two generated files have no Go counterpart of their own, by design (see
// dss/diagnostic/jaxb's doc.go, which documents the same choice for that
// package):
//
//   - package-info.java carries @XmlSchema(namespace=...,
//     elementFormDefault=QUALIFIED). There is no package-level annotation in
//     Go, so the namespace is the Namespace constant and the qualified
//     element form is what Marshal applies in xml.go.
//   - ObjectFactory.java is xjc's no-arg createXxx() factories plus the
//     @XmlElementDecl naming SimpleReport as the document element. Go
//     composite literals replace the factories, and the document element is
//     bound by XmlSimpleReport's own MarshalXML/UnmarshalXML (see
//     jaxb_root.go), so porting the class would add nothing callers can use.
//
// The parity contract is pinned by the round-trip KAT in xml_kat_test.go
// over testdata/oracle: simple-report dumps produced by upstream validation
// of real signed fixtures (testdata/SimpleReportOracle.java). Each must
// unmarshal into this model and marshal back byte for byte. The
// XSD-completeness sweep in xml_schema_test.go checks the model against
// SimpleReport.xsd element by element and attribute by attribute.
//
// Deferred: Facade.generateHtmlReport/generatePdfReport apply
// dss-simple-report-jaxb's xslt/html and xslt/pdf stylesheets to a
// simple-report document. The stylesheets are copied to testdata/xslt for
// the record, but no XSLT engine is ported and nothing in scope consumes
// them (same deferral as dss/diagnostic/jaxb's SVG rendering).
package jaxb
