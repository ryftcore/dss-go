// Ported from dss-simple-report-jaxb/src/main/java/eu/europa/esig/dss/simplereport/SimpleReportXmlDefinedUtils.java
// (DSS 6.5.RC1).
//
// DEFERRED: Java's SimpleReportXmlDefinedUtils is a singleton holding
// pluggable eu.europa.esig.dss.xml.common.TransformerFactoryBuilder/
// SchemaFactoryBuilder instances that SimpleReportXmlDefiner uses to build
// its (also deferred, see simple_report_xml_definer.go) secure
// TransformerFactory/SchemaFactory. Since neither XSD validation nor XSLT
// rendering is implemented in this port, there is nothing for a Go
// equivalent to configure; this file exists only so every hand-written
// manifest file keeps its own Go file per PORTING.md's one-file-per-class
// rule; SimpleReportXmlDefinedUtils has no Go declarations of its own.
package simplereport
