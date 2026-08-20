# Executor / report-builder oracle corpus

Two pure-Java dumps, both produced by the generators in `oracle/gen/` against
DSS 6.5.RC1 (`dss-validation`'s `target/classes` plus the module's runtime
classpath).

| file | rows | generator | input |
| --- | --- | --- | --- |
| `oracle/reports.jsonl` | 50 | `oracle/gen/ReportsOracle.java` | every diagnostic-data document of the marshal-parity corpus at `dss/diagnostic/jaxb/testdata/oracle`, minus the four `model-*.xml` schema-coverage fixtures |
| `oracle/hash_order.tsv` | 8 | `oracle/gen/HashOrder.java` | four hand-written id/name sequences |

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
phase-8a marshal-parity corpus and the phase-8c BasicBuildingBlocks corpus read.
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
```
