# POE core (EN 319 102-1 5.6.2.3 / 5.6.2.4) oracle corpus

`oracle/poe.jsonl` is a pure Java dump, produced by `gen/POEOracle.java`, which
drives the upstream classes

- `eu.europa.esig.dss.validation.process.vpfswatsp.POE` (`../poe.go`)
- `…vpfswatsp.TimestampPOE` (`../timestamp_poe.go`)
- `…vpfswatsp.EvidenceRecordPOE` (`../evidence_record_poe.go`)
- `…vpfswatsp.POEComparator` (`../poe_comparator.go`)
- `…vpfswatsp.POEExtraction` (`../poe_extraction.go`)

| file | rows | input |
| --- | --- | --- |
| `oracle/poe.jsonl` | 61 | the 54 diagnostic-data dumps of the marshal-parity corpus in `dss/diagnostic/jaxb/testdata/oracle`, plus this package's 7 synthetic dumps in `dd/`; control time 2024-01-01T00:00:00Z |

The real inputs are the marshal-parity corpus; the Go test
(`../poe_oracle_test.go`) reads the very same files, so neither side gets a
private fixture.

A row carries, per dump:

- **`tokens`** – for every token id `init()` seeds, the lowest POE after
  `init()`, after `collectAllPOE(getTimestampList())` and after
  `extractPOE(evidenceRecord)` for each evidence record, each as
  `{time, providerId, tokenProvided, objects}`; the `isPOEExists()` answer at
  four probe times (2000-01-01, 2016-06-01, 2020-01-01, the control time); and,
  for a used certificate, the `isPOEExistInRange(notBefore, notAfter)` answer.
- **`erExtraction`** – `"ok"` or `"throws"` per evidence record. Three dumps of
  the corpus carry an evidence record whose time-stamp reference never resolved,
  on which upstream's `EvidenceRecordPOE` raises a `NullPointerException` out of
  `getFirstTimestamp()`; the Go replay asserts the equivalent panic rather than
  hiding the case.
- **`compare`** – `POEComparator#compare` and `#before` over *every ordered
  pair* of the POEs the dump can build (the control-time POE, one `TimestampPOE`
  per time-stamp, one `EvidenceRecordPOE` per evidence record), so all four
  tie-breakers - `compareByTime`, `compareByType`, `compareByTimestampType`,
  `compareByTimestampedReferences` - are observed directly.
- **`signaturePOE`** – the lowest POE of each signature after
  `addSignaturePOE()`, which is what
  `ValidationProcessForSignaturesWithArchivalData` step 4) does.

## Why `dd/` exists

The marshal-parity corpus carries no evidence record that resolves, no two
time-stamps sharing a production time, no broken message imprint and no orphan
token, so `POEComparator`'s three tie-breakers past `compareByTime`,
`POEExtraction#extractPOE`'s message-imprint gate, and four of the eleven lists
`init()` walks are dead weight in it. `gen/POESyntheticDumps.java` adds:

| dump | what it turns on |
| --- | --- |
| `poe-timestamp-chain.xml` | three archive time-stamps, each covering the previous: POE propagation along a chain |
| `poe-evidence-record.xml` | an evidence record with two time-stamps covering signature + certificate + revocation: `EvidenceRecordPOE` takes the FIRST time-stamp's production time |
| `poe-revoked-revived.xml` | a certificate revoked in 2017 with a 2016 archive time-stamp over it and its revocation: the "revoked/expired revived by POE" shape of 5.6.2.4 step 3) |
| `poe-tie-timestamp-type.xml` | three time-stamps with the same production time and different types: `compareByTimestampType` |
| `poe-tie-covered-count.xml` | two time-stamps with the same production time and type, covering 1 vs 3 objects: `compareByTimestampedReferences` |
| `poe-broken-message-imprint.xml` | message-imprint not found / not intact / valid, plus a time-stamp covering nothing: `extractPOE`'s two guards |
| `poe-orphans-and-signer-data.xml` | orphan certificate token and signer data: the `getAllOrphan…` / `getAllSignerDocuments` arms of `init()` |

## Flagged: Java hash-order in `init()`'s inputs

`DiagnosticData#getAllSignatures()` and `#getAllRevocationData()` return
`java.util.HashSet`s, so two of the eleven lists `POEExtraction#init` walks are
iterated in String-hash order upstream. The Go wrappers return deterministic
slices instead, so the *order* differs. It is invisible: `init()` only seeds a
map keyed by token id, and every reader looks the id up. The Go test therefore
compares the seeded ids as a set, and matches `signaturePOE` rows by id rather
than by position.

## Not covered

`POEExtraction#extractPOE(List<XmlTimestampedObject>, Date)` is exercised only
indirectly (it is what `EvidenceRecordTimestampsValidationBlock#getPoe` calls);
its own corpus belongs with that block, which lives in
`../evidencerecord` and cannot be built until the `qualification` package lands
(it reaches `vpftsp`, which imports `qualification`).

## Regenerating

From a `dss-upstream` checkout with the modules built:

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle <dss-repo>/validation/process/vpfswatsp/testdata/gen/*.java
    java  -cp "$CP:/tmp/oracle" POESyntheticDumps <dss-repo>
    java  -cp "$CP:/tmp/oracle" POEOracle <dss-repo>/diagnostic/jaxb/testdata/oracle <dss-repo>

`POESyntheticDumps` must run first: `POEOracle` reads `dd/`.
