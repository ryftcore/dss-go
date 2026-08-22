# aov (algorithm-obsolescence validation) oracle corpora

All `oracle/*.jsonl` files are pure Java dumps, produced by the driver in `gen/`.

| file | rows | driver | input |
| --- | --- | --- | --- |
| `oracle/aov_blocks.jsonl` | 1328 | `AovOracle.java` | every one of upstream's ten concrete `AlgorithmObsolescenceValidation` subclasses over every signature, time-stamp, revocation, evidence record and used certificate of the marshal-parity diagnostic-data corpus and of the twelve synthetic dumps in `../../xcv/testdata/dd/`, default ETSI policy, validation time 2024-01-01T00:00:00Z |
| `oracle/aov_cc.jsonl` | 150 | `AovOracle.java` | `SignatureAlgorithmCryptographicChecker` and `DigestAlgorithmCryptographicChecker` over an algorithm x key-length x validation-time matrix |
| `oracle/aov_direct.jsonl` | 19 | `AovOracle.java` | the four `aov/checks` result checks driven alone through a one-item chain at `Level.FAIL`, over hand-built `XmlAOV` / `XmlCC` shapes |

The real inputs are the marshal-parity corpus in `dss/diagnostic/jaxb/testdata/oracle`
plus XCVA's `../../xcv/testdata/dd/`; the Go tests (`../aov_blocks_oracle_test.go`,
`../aov_cc_oracle_test.go` and `../aov_direct_oracle_test.go`, sharing
`../aov_oracle_test.go`) read the very same files, so neither side gets a private
fixture. The four `model-*.xml` schema-coverage fixtures are excluded, for the
reason given in `../../fc/testdata/README.md`.

A block row carries the whole `XmlAOV`: its constraint sequence, its conclusion,
and all four `XmlCryptographicValidation` members (signature value, signed
attributes, digest matchers, and the certificate-chain list), so the Go replay
compares the produced tree and not only its verdict.

## Why the cc matrix probes boundary instants

`CryptographicSuiteUtils` compares the validation time against an algorithm's
expiration date, and the reliability rules against its key size. A matrix that
only sampled a validation time far from either boundary would not tell a `<` from
a `<=`. Each (algorithm, key length) pair is therefore driven at four instants:
the fixed validation time, and the exact millisecond of the expiration date the
policy names for it, one millisecond before it and one millisecond after it. The
key lengths include one below the policy's threshold, one at it, one above it, an
unparseable one and a null one.

There is deliberately no empty-but-non-null key length in the matrix:
`TokenProxy#getKeyLengthUsedToSignThisToken` is nullable and the Go wrapper
collapses that null to `""`, so `""` is the port's spelling of null and a
non-null empty Java key length is unreachable from any caller.

Each of the six wired `cc/checks` classes reaches both its OK and its NOT OK
branch in this corpus (they raise five distinct constraint tags -
`SignatureAlgorithmAtValidationTimeCheck` and `DigestAlgorithmAtValidationTimeCheck`
share `ASCCM_AR`); `../aov_cc_oracle_test.go` fails if any of the five becomes
one-sided.

## Why aov_direct.jsonl exists

`AlgorithmObsolescenceValidationCheck`, `AlgorithmObsolescenceValidationCheckWithId`,
`SignatureAlgorithmCryptographicCheckerResultCheck` and
`DigestAlgorithmCryptographicCheckerResultCheck` are consumed by `fc` / `sav` /
`xcv`, not by the AOV blocks themselves (the blocks wire the first two only
indirectly), so the block corpus cannot reach every branch of their `process()`
and `buildAdditionalInfo()`. They are driven alone over the `XmlAOV` shapes
`passed` / `passed-with-algo` / `error` / `warning` / `info` and the `XmlCC`
shapes `passed` / `passed-no-keylength` / `failed`, which is every branch each of
them takes.

## Known mappings

- A one-item chain defines no title `MessageTag`, so Java leaves the `Title`
  attribute null while the generated Go model carries `Title` as a plain string.
  `../aov_direct_oracle_test.go` normalises that one field and nothing else (same
  as the `fc`/`sav` and XCVA direct corpora).
- Two block rows deviate from Java in the ORDER of their constraint list, and in
  nothing else:
  `SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation#buildSignedAttributesValidationChain`
  groups the signing-certificate references into a `java.util.HashMap` keyed by
  certificate id and iterates it in the JDK's bucket order, which is deterministic
  but is not the references' document order; the Go port iterates in first-seen
  (document) order instead, per PORTING.md's rule for an order-sensitive upstream
  map. `../aov_blocks_oracle_test.go` names the two rows explicitly and asserts
  that their constraint lists are permutations of each other and that everything
  else - conclusion, messages, all four cryptographic-validation members - is
  equal; a row that starts or stops deviating fails the test.

## Regenerating

From a `dss-upstream` checkout with the modules built:

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle <dss-repo>/validation/process/bbb/aov/testdata/gen/AovOracle.java
    java  -cp "$CP:/tmp/oracle" AovOracle <dss-repo>/diagnostic/jaxb/testdata/oracle <dss-repo>

`AovOracle` reads XCVA's `../../xcv/testdata/dd/*.xml` (already committed; no
`XcvaSyntheticDumps` run needed) alongside the marshal-parity corpus.
