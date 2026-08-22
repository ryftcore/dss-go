# Executor / report-builder oracle corpus

Three pure-Java dumps, all produced by the generators in `oracle/gen/` against
DSS 6.5.RC1 (`dss-validation`'s `target/classes` plus the module's runtime
classpath).

| file | rows | generator | input |
| --- | --- | --- | --- |
| `oracle/reports.jsonl` | 50 | `oracle/gen/ReportsOracle.java` | every diagnostic-data document of the marshal-parity corpus at `dss/diagnostic/jaxb/testdata/oracle`, minus the four `model-*.xml` schema-coverage fixtures |
| `oracle/full_corpus.jsonl` | 273 | `oracle/gen/FullCorpusOracle.java` | EVERY file under upstream's `dss-validation/src/test/resources/diag-data/**`, vendored at `oracle/full-corpus/` - see below |
| `oracle/full_corpus_reports.jsonl` | 252 | `oracle/gen/FullCorpusReportsOracle.java` | the same corpus, digests of every marshalled report - see below |
| `oracle/hash_order.tsv` | 8 | `oracle/gen/HashOrder.java` | four hand-written id/name sequences |

## `oracle/full_corpus.jsonl` and `oracle/full-corpus/`

Executor-level verdict parity, checked by `full_corpus_oracle_test.go`'s
`TestFullCorpusExecutorOracle` against
literally every file (273, no exclusions) under upstream's diag-data test
resource tree - the real corpus `DefaultSignatureProcessExecutorTest` and its
neighbors draw from, spanning every JIRA-ticket regression fixture plus the
categorized `cert-validation/`, `qwac-validation/`, `policy/`, `crypto-suite/`,
`timestamp-validation/`, `eaa-validation/`, `er-validation/`, `sig-qualification/`
and `mra/` subdirectories - vendored byte-identically under `oracle/full-corpus/`
with paths preserved.

Both engines run ONE fixed executor pair uniformly over every file:
`DefaultSignatureProcessExecutor` (default policy, `ValidationLevel.ARCHIVAL_DATA`,
validation time `1700000000000`, locale `en`, ETSI report disabled - this dump
only checks verdicts, not report bytes), and, when the document carries at
least one used certificate, `DefaultCertificateProcessExecutor` against the
first one with `/policy/certificate-constraint.xml` (matching
`oracle/reports.jsonl`'s own convention). This is a uniform parity probe, not
a reproduction of each fixture's original upstream unit-test scenario -
several subdirectories are consumed upstream by OTHER executors entirely
(`DefaultCertificateProcessExecutor` directly, `DefaultTimestampProcessExecutor`,
QWAC-specific validators). Running the identical Go/Java code path on
identical input either way is what "verdict parity" means here.

21 of the 273 files are not diagnostic-data documents at all - validation-
policy and crypto-suite-constraint XML fixtures upstream's own tests pass to
the POLICY loader, not the diagnostic-data unmarshaller, living in this
resource tree only because upstream keeps every
`DefaultSignatureProcessExecutorTest` input under one directory. The Java
oracle does not special-case them: it attempts to unmarshal every file as
`XmlDiagnosticData` exactly like the Go test does, so both sides record the
SAME unmarshal failure for the SAME 21 rows - itself a parity assertion, not
a quarantine.

Each row records, for both executor runs: every `BasicBuildingBlocks`
conclusion (Id/Type/Indication/SubIndication), the final Indication/
SubIndication of every top-level Signature/Timestamp/EvidenceRecord, and every
signature's SimpleReport qualification.

## `oracle/full_corpus_reports.jsonl`

Report BYTE parity beyond `oracle/reports.jsonl`'s 50-document subset:
`full_corpus_reports_oracle_test.go`'s `TestFullCorpusReportsByteParity` runs
BOTH executors over every diagnostic-data document of `oracle/full-corpus/`
that unmarshals (252 of the 273 files - the 21 policy/crypto-suite constraint
fixtures misfiled into the tree carry no row) and compares the SHA-256 of all
five marshalled reports: 1255 report comparisons, no tolerances. Same
configuration as `oracle/reports.jsonl` (validation time `1700000000000`,
`ValidationLevel.ARCHIVAL_DATA`, locale `en`, ETSI validation report ENABLED,
`/policy/certificate-constraint.xml` for the certificate executor).

Three documents' signature `DetailedReport` carries EAA validation blocks the
`eaa` build tag gates in; the test compares their digest only under that tag,
and asserts they still DIFFER without it (`eaaOnlyDetailedReports`).

## `oracle/reports.jsonl`

One row per input document. Each row holds the SHA-256 of every report the two
executors this package ports marshal for that document:

* `simpleReport`, `detailedReport`, `etsiValidationReport` -
  `eu.europa.esig.dss.validation.executor.signature.DefaultSignatureProcessExecutor`,
  default ETSI validation policy, `ValidationLevel.ARCHIVAL_DATA`, ETSI
  validation report enabled, semantics disabled;
* `certificateId`, `simpleCertificateReport`, `certificateDetailedReport` -
  `...executor.certificate.DefaultCertificateProcessExecutor` over the SAME
  document, `/policy/certificate-constraint.xml`, targeting the first entry of
  `<UsedCertificates>`.

Both runs use validation time `1700000000000` (2023-11-14T22:13:20Z) and the
`en` locale. Comparing digests is byte-parity: equal digests mean the Go
marshalled report is byte-for-byte the Java one.

`oracle/xml/` carries the FULL Java output for five rows - one per format family
(CAdES-LTA, CAdES with an evidence record, JAdES-LTA, PAdES-LTV, XAdES with
several archive time-stamps) - so a digest mismatch on those can be diffed
rather than merely reported.

The inputs are not a private fixture: they are the very same documents the
marshal-parity corpus and the BasicBuildingBlocks corpus read.
The `model-*.xml` files are excluded for the reason given in
`../../process/bbb/fc/testdata/README.md`: they are schema-coverage dumps, not
real validation inputs.

## `oracle/hash_order.tsv`

`kind<TAB>input<TAB>expected`, the members joined by `|`. `TOKEN` rows record
the order a `java.util.HashSet` iterates keys whose `hashCode()` is
`31 + id.hashCode()` (what `AbstractTokenProxy` computes); `STRING` rows record
the same for plain `String` keys. They pin `JavaHashSetOrder` /
`JavaHashSetStringOrder`, which the report builders use wherever a Java
hash-ordered collection's iteration order reaches the marshalled report.

## Regenerating

```
CP=<dss-validation runtime classpath>:<dss-validation target/classes>
javac -cp "$CP" -d /tmp/oracle oracle/gen/ReportsOracle.java
java  -cp "$CP:/tmp/oracle" ReportsOracle \
      ../../../diagnostic/jaxb/testdata/oracle oracle oracle/xml \
      CAdESDoubleLTA.p7m cades-lt-with-er-multi-files.sce jades-lta.json \
      pades-ltv.pdf xades-with-multiple-archivetimestamps.xml

javac -d /tmp/oracle oracle/gen/HashOrder.java
java  -cp /tmp/oracle HashOrder > oracle/hash_order.tsv

# full_corpus.jsonl - only dss-validation's own classes needed, no test resources
CP=<dss-validation runtime classpath>:<dss-validation target/classes>
javac -cp "$CP" -d /tmp/oracle oracle/gen/FullCorpusOracle.java
java  -cp "$CP:/tmp/oracle" FullCorpusOracle \
      <dss-upstream>/dss-validation/src/test/resources/diag-data \
      oracle/full_corpus.jsonl

javac -cp "$CP" -d /tmp/oracle oracle/gen/FullCorpusReportsOracle.java
java  -cp "$CP:/tmp/oracle" FullCorpusReportsOracle \
      <dss-upstream>/dss-validation/src/test/resources/diag-data \
      oracle/full_corpus_reports.jsonl
```
