// Ported from dss-simple-certificate-report-jaxb/target/generated-sources/xjc/eu/europa/esig/dss/simplecertificatereport/jaxb
// (DSS 6.5.RC1), the JAXB model xjc generates from SimpleCertificateReport.xsd.

// Package jaxb is the schema-shaped model of a DSS simple certificate
// report.
//
// Every generated JAXB class keeps its name (XmlChainItem stays
// XmlChainItem) and its property order, and every property carries the xml
// struct tag that reproduces its JAXB annotation, so that marshalling with
// Marshal produces the bytes the JAXB reference implementation produces for
// the same tree. Per the phase-8a generated-JAXB rule the classes are
// grouped into files by schema area rather than one file per class.
//
// Two generated files have no Go counterpart of their own, by design (see
// dss/diagnostic/jaxb's and dss/simplereport/jaxb's doc.go, which document
// the same choice for those packages):
//
//   - package-info.java carries @XmlSchema(namespace=...,
//     elementFormDefault=QUALIFIED). There is no package-level annotation
//     in Go, so the namespace is the Namespace constant and the qualified
//     element form is what Marshal applies in xml.go.
//   - ObjectFactory.java is xjc's no-arg createXxx() factories plus the
//     @XmlElementDecl naming SimpleCertificateReport as the document
//     element. Go composite literals replace the factories, and the
//     document element is bound by XmlSimpleCertificateReport's XMLName
//     field (see jaxb_root.go), so porting the class would add nothing
//     callers can use.
//
// The parity contract is pinned by the round-trip KAT in xml_kat_test.go
// over testdata/oracle: simple-certificate-report dumps produced by
// upstream validation of real signed certificate fixtures
// (testdata/SimpleCertificateReportOracle.java). Each must unmarshal into
// this model and marshal back byte for byte. The XSD-completeness sweep in
// xml_schema_test.go checks the model against SimpleCertificateReport.xsd
// element by element and attribute by attribute.
package jaxb
