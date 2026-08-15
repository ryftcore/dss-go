# Simple-certificate-report marshal-parity oracles

Ground truth for `dss/simplecertificatereport/jaxb`: every file under
`oracle/` was produced by the JAXB reference implementation of upstream DSS
6.5.RC1, and the KAT in `../xml_kat_test.go` requires each of them to
unmarshal into the Go model and marshal back byte for byte.

- 7 files named after a certificate fixture (`*.cer.xml`) —
  `SimpleCertificateReportOracle.java` runs upstream certificate validation
  (`CertificateValidator`, default validation policy, pinned validation
  time 2025-06-15T12:00:00Z) over an X.509 certificate fixture drawn from
  upstream's own `dss-certificate-validation-common`/`dss-spi`/`dss-model`/
  `dss-cades` test resources and dumps
  `CertificateReports.getXmlSimpleReport()`.
- `simple-cert-report*.xml` (7 files) — hand/IDE-formatted
  dss-simple-certificate-report-jaxb unit-test fixtures
  (`src/test/resources/simple-cert-report*.xml`), re-marshalled through
  `SimpleCertificateReportFacade` (unmarshal, then marshal with
  `validateXml=false`) so they become valid parity bytes regardless of
  their original formatting. These are the fixtures that exercise QWAC
  profiles/TLS binding signatures, PID certificate-approval-status entries,
  and (`simple-cert-report2.xml`) a certificate with no notBefore/notAfter
  at all - the finding that pinned notBefore/notAfter as *XSDateTime rather
  than XSDateTime in jaxb_chain_item.go, since the schema's "required"
  (no minOccurs="0") is not enforced by the RI's unmarshaller/marshaller.
- `xsd/SimpleCertificateReport.xsd`, `xsd/bindings.xml` — byte-identical
  copies of `dss-simple-certificate-report-jaxb/src/main/resources/xsd`,
  swept against the model by `../jaxb_schema_test.go`.
- `xslt/html/simple-certificate-report-bootstrap4.xslt`,
  `xslt/pdf/simple-certificate-report.xslt` — copied for the record only;
  XSLT rendering is deferred, no XSLT engine is ported (see
  `simplecertificatereport/simple_certificate_report_xml_definer.go`).

## Regenerating

Build upstream DSS with maven, then compose a classpath of the installed
`eu.europa.ec.joinup.sd-dss` jars and their dependencies (no pdfbox needed -
certificate validation never touches a PDF document):

    javac -cp "$CP" -d /tmp/oracle SimpleCertificateReportOracle.java ReserializeSimpleCertificateReport.java
    java -cp "$CP:/tmp/oracle" SimpleCertificateReportOracle oracle <certFile>...
    java -cp "$CP:/tmp/oracle" ReserializeSimpleCertificateReport oracle <upstream-test-resource>...

`DSS_SIMPLE_CERTIFICATE_REPORT_ORACLE_DIR` replays the KAT over a larger
local corpus without committing it.
