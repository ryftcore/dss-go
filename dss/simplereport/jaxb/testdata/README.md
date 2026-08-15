# Simple-report marshal-parity oracles

Ground truth for `dss/simplereport/jaxb`: every file under `oracle/` was
produced by the JAXB reference implementation of upstream DSS 6.5.RC1, and
the KAT in `../xml_kat_test.go` requires each of them to unmarshal into the
Go model and marshal back byte for byte.

- 18 files named after a signed fixture (`*.p7m.xml`, `*.xml.xml`,
  `*.json.xml`, `*.asic{s,e}.xml`) — `SimpleReportOracle.java` runs upstream
  validation (default validation policy, pinned validation time
  2025-06-15T12:00:00Z) over a CAdES/XAdES/JAdES/ASiC fixture from this
  repo's own `cades/`, `xades/`, `jades/`, `asic/` testdata trees and dumps
  `Reports.getXmlSimpleReport()`.
- `sr-*.xml`, `sr1.xml`, `sr2.xml` (12 files) — hand/IDE-formatted
  dss-simple-report-jaxb unit-test fixtures (`src/test/resources/sr*.xml`),
  re-marshalled through `SimpleReportFacade` (unmarshal, then marshal with
  `validateXml=false`) so they become valid parity bytes regardless of their
  original formatting. These are the only fixtures in the corpus that
  exercise EAA/mdoc/PID tokens, the EAALevel/EAASignature/KeyBindingSignature/
  EAAPayload elements, PDFAInfo, Semantic, and a sunset-dated trust anchor -
  none of upstream's own signed-fixture corpus produces those.
- `xsd/SimpleReport.xsd`, `xsd/bindings.xml` — byte-identical copies of
  `dss-simple-report-jaxb/src/main/resources/xsd`, swept against the model
  by `../jaxb_schema_test.go`.
- `xslt/html/simple-report-bootstrap4.xslt`,
  `xslt/pdf/simple-report.xslt` — copied for the record only; XSLT rendering
  is deferred, no XSLT engine is ported (see `simplereport/simple_report_xml_definer.go`).

## Regenerating

Build upstream DSS with maven, then compose a classpath of the installed
`eu.europa.ec.joinup.sd-dss` jars, their dependencies, and
`dss-pades-pdfbox/target/classes` plus the pdfbox jars (the PAdES validator
is loaded through `ServiceLoader`, so the implementation module has to be
on the classpath even to validate a non-PDF document):

    javac -cp "$CP" -d /tmp/oracle SimpleReportOracle.java ReserializeSimpleReport.java
    java -cp "$CP:/tmp/oracle" SimpleReportOracle oracle <fixture>...
    java -cp "$CP:/tmp/oracle" ReserializeSimpleReport oracle <upstream-test-resource>...

`DSS_SIMPLE_REPORT_ORACLE_DIR` replays the KAT over a larger local corpus
without committing it.
