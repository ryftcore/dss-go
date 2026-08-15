# BasicBuildingBlocks (5.2 dispatcher) oracle corpus

`oracle/blocks.jsonl` is a pure Java dump, produced by `gen/BlocksOracle.java`,
which drives the upstream class
`eu.europa.esig.dss.validation.process.bbb.BasicBuildingBlocks` (ported here as
package `blocks`, see `../basic_building_blocks.go`'s header for why it cannot
live in the Go package `bbb`).

| file | rows | input |
| --- | --- | --- |
| `oracle/blocks.jsonl` | 474 | every signature (`SIGNATURE` / `COUNTER_SIGNATURE`), time-stamp (`TIMESTAMP`), revocation (`REVOCATION`) and used certificate (`CERTIFICATE`) of the marshal-parity diagnostic-data corpus, of XCVA's twelve synthetic dumps in `../../bbb/xcv/testdata/dd/`, and of this package's own `dd/` dump; default ETSI policy, validation time 2024-01-01T00:00:00Z |

The real inputs are the marshal-parity corpus in
`dss/diagnostic/jaxb/testdata/oracle` plus the two `dd/` directories; the Go test
(`../blocks_oracle_test.go`) reads the very same files, so neither side gets a
private fixture. The four `model-*.xml` schema-coverage fixtures are excluded, for
the reason given in `../../bbb/fc/testdata/README.md`.

A row carries what the dispatcher itself decides, rather than what its sub-blocks
compute (those have their own corpora, one per package):

- which of the seven sub-blocks it instantiated at all - `executeFormatChecking`'s
  three-way context switch, `executeIdentificationOfTheSigningCertificate`'s
  `CERTIFICATE`/`EAA` exclusion, `executeValidationContextInitialization`'s
  signature-only gate, `executeCryptographicVerification`'s `CERTIFICATE`
  exclusion, and the time-stamp branch that discards an `XmlFC` carrying no
  constraint at all;
- each instantiated block's title and conclusion;
- the final conclusion `updateFinalConclusion()` aggregated out of them;
- the certificate chain copied off the ISC block;
- the cross-certificate and equivalent-certificate id lists
  `addAdditionalInfo()` hangs off every `XmlSubXCV`.

Like the oracle, the Go replay shares one `bbbs` map across every token of a dump.

## Why `dd/` exists

`addAdditionalInfo(XmlXCV)` asks the diagnostic data for the cross-certificates
and equivalent certificates of every certificate in the chain. Neither the 50-dump
marshal-parity corpus nor XCVA's `dd/` dumps contain two used certificates sharing
an `EntityKey`, so both answers are empty everywhere and the whole method is dead
weight in the corpus - an audit mutation dropping its
`equivalentCertificates.removeAll(crossCertificates)` step survived until
`gen/BlocksSyntheticDumps.java` added `dd/cross-equivalent.xml`: a signing
certificate plus one further used certificate with the same `EntityKey` and a
different subject DN (cross, and therefore also equivalent - which is exactly what
the `removeAll` step has to subtract again) and one with the same DN (equivalent
but not cross).

The orphan-certificate halves of the same method
(`getOrphanCrossCertificates`/`getOrphanEquivalentCertificates`) are still not
reached by any input here.

## Not covered

The `Context.EAA` / `Context.EAA_REVOCATION` branches of every dispatch method are
unreachable from this corpus: no diagnostic-data dump carries an EAA token, and
the four constructors two of those branches call only exist under `-tags eaa`
(see `../basic_building_blocks_eaa.go`). `Context.KEY_BINDING_SIGNATURE` is
likewise unreached.

## Regenerating

From a `dss-upstream` checkout with the modules built:

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle <dss-repo>/validation/process/blocks/testdata/gen/*.java
    java  -cp "$CP:/tmp/oracle" BlocksSyntheticDumps <dss-repo>
    java  -cp "$CP:/tmp/oracle" BlocksOracle <dss-repo>/diagnostic/jaxb/testdata/oracle <dss-repo>

`BlocksSyntheticDumps` must run first: `BlocksOracle` reads `dd/`. It also reads
XCVA's `../../bbb/xcv/testdata/dd/*.xml`, which are committed already.
