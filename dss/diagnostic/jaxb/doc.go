// Ported from dss-diagnostic-jaxb/target/generated-sources/xjc/eu/europa/esig/dss/diagnostic/jaxb
// (DSS 6.5.RC1), the JAXB model xjc generates from DiagnosticData.xsd.

// Package jaxb is the schema-shaped model of a DSS diagnostic-data report.
//
// Every generated JAXB class keeps its name (XmlCertificate stays
// XmlCertificate) and its property order, and every property carries the xml
// struct tag that reproduces its JAXB annotation, so that marshalling with
// Marshal produces the bytes the JAXB reference implementation produces for the
// same tree. Per the phase-8a generated-JAXB rule the classes are grouped into
// files by schema area rather than one file per class; jaxb_model.go lists them
// all, and the schema sweep in jaxb_schema_test.go checks the model against
// DiagnosticData.xsd complexType by complexType.
//
// Two generated files have no Go counterpart of their own, by design:
//
//   - package-info.java carries @XmlSchema(namespace=...,
//     elementFormDefault=QUALIFIED). There is no package-level annotation in Go,
//     so the namespace is the Namespace constant and the qualified element form
//     is what Marshal applies in xml.go.
//   - ObjectFactory.java is xjc's 146 no-arg createXxx() factories plus the
//     @XmlElementDecl naming DiagnosticData as the document element. Go composite
//     literals replace the factories, and the document element is bound by
//     XmlDiagnosticData's XMLName field, so porting the class would add nothing
//     callers can use.
//
// The parity contract is pinned by the round-trip KAT in xml_kat_test.go over
// testdata/oracle: diagnostic-data dumps produced by upstream validation of real
// signed fixtures (testdata/DiagnosticDataOracle.java), plus dumps in which
// every property of every reachable class is set
// (testdata/DiagnosticDataFillOracle.java). Each must unmarshal into this model
// and marshal back byte for byte. TestOracleCorpusExercisesModel additionally
// requires that corpus to reach every element and attribute name the model
// binds, so no binding is pinned by the schema sweep alone.
//
// Deferred: DiagnosticDataFacade.generateSVG applies
// dss-diagnostic-jaxb's xslt/svg/diagnostic-data.xslt to a diagnostic-data
// document. The stylesheet is copied to testdata/xslt/svg for the record, but no
// XSLT engine is ported and nothing in scope consumes it.
package jaxb
