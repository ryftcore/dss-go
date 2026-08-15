# xcv (X.509 certificate validation) XCVA oracle corpora

Both `oracle/*.jsonl` files are pure Java dumps, produced by the drivers in `gen/`.

| file | rows | driver | input |
| --- | --- | --- | --- |
| `oracle/xcva_blocks.jsonl` | 916 | `XcvaOracle.java` | `X509CertificateValidation` / `CertificateRevocationSelector` / `RevocationAcceptanceChecker` over every signature, time-stamp, revocation and used certificate of the marshal-parity diagnostic-data corpus and of the nine synthetic dumps in `dd/`, validation time 2024-01-01T00:00:00Z |
| `oracle/xcva_direct.jsonl` | 1279 | `XcvaDirectOracle.java` | all 12 `rac/checks` classes driven alone at `Level.FAIL` over every (certificate, certificate revocation data) pair of the same dumps |

The real inputs are the marshal-parity corpus in `dss/diagnostic/jaxb/testdata/oracle`
plus `dd/`; the Go tests (`../xcva_blocks_oracle_test.go`,
`../xcva_direct_oracle_test.go`, sharing `../xcva_oracle_test.go`) read the very same
files, so neither side gets a private fixture.

The `XmlAOV` the blocks consume is a phase 8d product; the block corpus feeds them a
PASSED one, the way the phase 8c sav corpus does.

## Why `dd/` exists

Not one certificate of the 50-dump marshal-parity corpus is trusted
(`grep -c 'Trusted' | <Trusted>false</Trusted>` throughout), and the default ETSI
policy runs `ProspectiveCertificateChainCheck` at `Level.FAIL`. Over that corpus
alone `X509CertificateValidation` therefore stops on its very first chain item for
every token, and none of

- the trust-anchor walk and its `SubContext` progression,
- the sunset-date items (`CertificateValidationBeforeSunsetDateWithIdCheck`,
  `ProspectiveCertificateChainAtValidationTimeCheck`) and the `SUB_XCV_TA` message
  filtering of `collectMessages` / `collectAdditionalMessages`,
- the `SubX509CertificateValidation` loop over the CA certificates,
- the `ValidationModel` lastDate progression,
- an `XmlRAC` that concludes PASSED

is ever reached. `gen/XcvaSyntheticDumps.java` writes the nine dumps in `dd/` -
trusted root, trusted signing certificate, past and future sunset dates, a trust
start date, a fully untrusted chain, a trusted root with an OCSP response, a
PASSED-RAC shape, and the revocation shapes (`certHash` present and matching or
mismatching, a resolving responder-id reference, an absent `thisUpdate`, an
`expiredCertsOnCRL` both before and after `thisUpdate`, an `archiveCutoff`) the
corpus has none of. They are marshalled with the upstream `DiagnosticDataFacade`
and committed, so the Go replay reads them unchanged.

`policy/constraint-chain.xml` is the default ETSI policy with `Model` switched to
`CHAIN` and every `ProspectiveCertificateChain` dropped to `WARN`; the block corpus
runs the `dd/` dumps under it as well as under the default policy, which is what
exercises the CA-certificate loop without an early return and the per-model
`lastDate` progression. Each row names the policy it was produced under.

The four `model-*.xml` schema-coverage fixtures are excluded from both corpora, for
the reason given in `../../fc/testdata/README.md`.

Every one of the 12 `rac/checks` classes has both an `OK` and a `NOT OK` row, and
the block corpus carries both `PASSED` and `INDETERMINATE` conclusions for each of
`XCV`, `CRS` and `RAC`; both tests fail if either property is lost.

## Regenerating

From a `dss-upstream` checkout with the modules built:

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle <dss-repo>/validation/process/bbb/xcv/testdata/gen/*.java
    java  -cp "$CP:/tmp/oracle" XcvaSyntheticDumps <dss-repo>
    java  -cp "$CP:/tmp/oracle" XcvaOracle       <dss-repo>/diagnostic/jaxb/testdata/oracle <dss-repo>
    java  -cp "$CP:/tmp/oracle" XcvaDirectOracle <dss-repo>/diagnostic/jaxb/testdata/oracle <dss-repo>

`XcvaSyntheticDumps` must run first: the two oracles read `dd/`.

## Known mappings

- A single-item chain defines no title `MessageTag`, so Java leaves the `Title`
  attribute null while the generated Go model carries `Title` as a plain string.
  The direct test normalises that one field and nothing else (same as the phase 8c
  direct corpora).
- `XmlRAC/RevocationThisUpdate` and `XmlRAC/RevocationProductionDate` are required
  elements in `DetailedReport.xsd`, so the generated Go model carries them as plain
  (non-pointer) `XSDateTime` members where the JAXB class has nullable `Date`s. The
  zero time stands in for Java's null, and the block corpus renders a zero one as
  `null`. A revocation without a production date is therefore marshalled by the Go
  detailed report as `0001-01-01T00:00:00Z` where Java omits the element - a
  deviation of the frozen `detailedreport/jaxb` package, not of this one.
