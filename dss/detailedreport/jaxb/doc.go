// Ported from dss-detailed-report-jaxb/target/generated-sources/xjc/eu/europa/esig/dss/detailedreport/jaxb
// (DSS 6.5.RC1), the JAXB model xjc generates from DetailedReport.xsd.

// Package jaxb is the schema-shaped model of a DSS detailed report.
//
// Every generated JAXB class keeps its name (XmlSignature stays XmlSignature)
// and its property order, and every property carries the xml struct tag that
// reproduces its JAXB annotation, so that marshalling with Marshal produces the
// bytes the JAXB reference implementation produces for the same tree. The
// classes are grouped into files by schema area rather than one file per class;
// jaxb_model.go lists them all, and the schema sweep in jaxb_schema_test.go
// checks the model against DetailedReport.xsd complexType by complexType.
//
// Two generated files have no Go counterpart of their own, by design, exactly
// as in dss/diagnostic/jaxb:
//
//   - package-info.java carries @XmlSchema(namespace=...,
//     elementFormDefault=QUALIFIED). There is no package-level annotation in Go,
//     so the namespace is the Namespace constant and the qualified element form
//     is what Marshal applies in xml.go.
//   - ObjectFactory.java is xjc's single createDetailedReport() factory plus the
//     @XmlElementDecl naming DetailedReport as the document element. Marshal
//     binds the document element directly, so porting the class would add
//     nothing callers can use.
//
// The parity contract is pinned by the round-trip KAT in xml_kat_test.go over
// testdata/oracle: DetailedReport dumps produced by the JAXB reference
// implementation (org.glassfish.jaxb:jaxb-runtime 3.0.2, the same version xjc
// generated the classes with) re-marshalling the module's own test fixtures
// (dss-detailed-report-jaxb/src/test/resources/dr*.xml) - see
// testdata/oracle/README.md for the exact regeneration recipe. Each dump must
// unmarshal into this model and marshal back byte for byte.
// TestOracleCorpusExercisesModel additionally requires that corpus to reach
// every element and attribute name the model binds, so no binding is pinned by
// the schema sweep alone.
//
// Deferred: nothing in this schema requires an XSLT engine or XSD validator (see
// the detailedreport package's doc.go for the DetailedReportXmlDefiner/
// DetailedReportFacade deferrals, which mirror dss/diagnostic's).
package jaxb
