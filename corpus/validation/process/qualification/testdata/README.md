# `validation/process/qualification` oracle corpora

Four pure-Java dumps, produced by the drivers in `gen/`, which instantiate the
upstream `eu.europa.esig.dss.validation.process.qualification` classes directly.
Nothing in `oracle/` is hand-written or hand-edited; each file regenerates
byte-identical from the recipe at the bottom.

| file | rows | driver | what it pins |
| --- | --- | --- | --- |
| `oracle/qual_trust.jsonl` | 1248 | `gen/QualTrustOracle.java` | the trust-service surface |
| `oracle/qual_cert.jsonl` | 354 | `gen/QualCertOracle.java` | the certificate surface |
| `oracle/qual_block.jsonl` | 6400 | `gen/QualBlockOracle.java` | `CertQualificationAtTimeBlock`, whole result tree |
| `oracle/qual_sig_block.jsonl` | 1944 | `gen/QualSigBlockOracle.java` | `SignatureQualificationBlock`, whole result tree |

## `qual_trust.jsonl`

A designed cross-product of `TrustServiceWrapper`s - status x type x
start/end date x captured qualifiers x additional service information x country
code x TL URL x MRA fields - covering both branches of every predicate below,
plus the matrix and vocabulary cells:

- the ten `TrustServiceChecker` predicates (and the twelve
  `TrustServiceCondition` classes behind them);
- the thirteen `TrustServicesFilterFactory` filters, as accept/reject over a
  singleton list; the two date-sensitive ones at each of five probe dates;
- `QualificationByTL` / `TypeByTL` / `QSCDByTL` over every stub-strategy input;
- `CertQualificationMatrix` (the full cube), `SigQualificationMatrix` (every
  `Indication` x `CertificateQualification` cell, including the cells upstream
  throws on) and `FinalCertificateQualificationCalculator` (every ordered pair);
- `EIDASUtils` at ten probes straddling both boundaries, plus `null`;
- `TrustServiceStatus` and `ServiceTypeIdentifier` - every constant's fields and
  every URI probe, including unknown URIs and `null`.

## `qual_cert.jsonl`

A cross-product of synthetic `XmlCertificate`s - `QcCompliance` x `QcSSCD` x
`QcTypes` x `QcCClegislation` x certificate policies x `notBefore` on both sides
of the eIDAS boundary - crossed with a fixed panel of twelve trust services:

- `QCTypeIdentifiers`;
- the six `…ByCertificate{Pre,Post}EIDAS` strategies, through the factories that
  pick between them;
- the full `CertificateQualificationCalculator`, and the
  `createXxxFromCertAndTL` strategy triple, per trust service;
- `ServiceByCertificateTypeFilter`;
- `UniqueServiceFilter` over the whole panel, over every adjacent pair, over a
  singleton and over the empty list - so both its "one conclusion" and its
  "several conclusions" branches are observed.

## `qual_block.jsonl`

`CertQualificationAtTimeBlock` over 8 certificates x 80 trust-service selections
(the empty list, every singleton, every unordered pair, and the whole panel) x 5
validation dates x both `ValidationTime`s. Each row carries the whole
`XmlValidationCertificateQualification`: the **ordered** `Constraint` list with
every item's name key, status, error/warning/info keys and additional info; the
`Conclusion` with its ordered message keys; the resulting
`CertificateQualification`; and `getFilteredServices()`.

This is the block that carries the normative `TrustServiceFilter` **order** of
TS 119 615, so an out-of-order filter shows up as a differently ordered or
differently sized constraint list.

## `qual_sig_block.jsonl`

`SignatureQualificationBlock` over 4 certificates x 9 trust-service-provider
shapes x 6 `XmlTLAnalysis` panels x 3 EN 319 102-1 conclusions x 3
best-signature-times, dumping the whole `XmlValidationSignatureQualification`
including both nested `XmlValidationCertificateQualification` blocks.

## Flagged: Java hash-order in the conflict-detection message

`CertQualificationAtTimeBlock` collects the simulated per-trust-service
qualification outcomes into a `Set<CertificateQualification>` - a `HashSet` of a
plain Java enum, whose iteration order is JVM identity-hash-bucket order. The
corpus carries both `Results : [CERT_FOR_ESIG, QCERT_FOR_ESIG]` and
`Results : [QCERT_FOR_ESIG_QSCD, CERT_FOR_ESIG]`, i.e. neither insertion nor
reverse-insertion order, so it is not reproducible from Go and not guaranteed by
upstream across JVMs either. The Go port renders the same values in first-seen
order; `qual_block_oracle_test.go` therefore compares that one
`AdditionalInfo` field's bracketed value list as a set and everything else
exactly. `IsNoQualificationConflictDetectedCheck#process()` itself only reads
the set's size, which is order-independent.

`SignatureQualificationBlock` and `CertificateQualificationBlock` iterate
`HashSet<String>` LOTL/TL URL sets into report constraints; the Go port sorts
those sets, which is deterministic but not necessarily Java-bucket-identical.
The corpus only exercises one TL and one LOTL URL, so the difference is not
observable in it; a multi-TL dump would need the same set-comparison treatment.

## Regenerating

From a `dss-upstream` checkout with the modules built:

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="$PWD/dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle <dss-repo>/validation/process/qualification/testdata/gen/*.java
    for driver in QualTrustOracle QualCertOracle QualBlockOracle QualSigBlockOracle; do
      java -cp "$CP:/tmp/oracle" "$driver" <dss-repo>
    done

The drivers are independent; any of them may be run alone.
