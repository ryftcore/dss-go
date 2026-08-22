# xcv (X.509 certificate validation) XCVA + XCVB oracle corpora

All `oracle/*.jsonl` files are pure Java dumps, produced by the drivers in `gen/`.

| file | rows | driver | input |
| --- | --- | --- | --- |
| `oracle/xcva_blocks.jsonl` | 1067 | `XcvaOracle.java` | `X509CertificateValidation` / `CertificateRevocationSelector` / `RevocationAcceptanceChecker` over every signature, time-stamp, revocation and used certificate of the marshal-parity diagnostic-data corpus and of the twelve synthetic dumps in `dd/`, validation time 2024-01-01T00:00:00Z |
| `oracle/xcva_direct.jsonl` | 1341 | `XcvaDirectOracle.java` | all 12 `rac/checks` classes driven alone at `Level.FAIL` over every (certificate, certificate revocation data) pair of the same dumps |
| `oracle/xcv_direct.jsonl` | ~33800 | `XcvOracle.java` | all 75 instantiable `checks` / `sub/checks` / `sub/checks/pseudo` / `rfc/checks` classes driven alone at `Level.FAIL`, over every used certificate (and, for the revocation-facing checks, every certificate-revocation pair) of the marshal-parity corpus plus `dd/`, plus a handful of hand-built certificates/revocations (see "XCVB direct corpus" below) |

The real inputs are the marshal-parity corpus in `dss/diagnostic/jaxb/testdata/oracle`
plus `dd/`; the Go tests (`../xcva_blocks_oracle_test.go`,
`../xcva_direct_oracle_test.go`, sharing `../xcva_oracle_test.go`; and
`../xcv_direct_oracle_test.go` for the `xcv_direct.jsonl` row above, reusing
`xcvaConclusion`/`xcvaConstraint`/`xcvaMessage` and the corpus/policy loaders
from `xcva_oracle_test.go`) read the very same files, so neither side gets a
private fixture.

The `XmlAOV` the blocks consume is built by the `aov` package; the block corpus
feeds them a PASSED one, the way the `sav` corpus does.

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

is ever reached. `gen/XcvaSyntheticDumps.java` writes the twelve dumps in `dd/` -
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

`policy/constraint-certhash.xml` is the default policy with the two OCSP certHash
constraints (`OCSPCertHashPresent`, `OCSPCertHashMatch`) raised to `FAIL`; the
block corpus runs the `dd/` dumps under it too. The default ETSI policy leaves
both undefined, so `RevocationAcceptanceChecker`'s certHash items are level-less
and record no constraint at all - which leaves the "wire the match check only when
the extension is present" gate in its `initChain()` unobservable (an audit
mutation inverting that gate survived the corpus until this variant existed).

Four of the `dd/` dumps exist for the same reason - a branch of the crs/rac
selection logic that no other input reaches, each one found by an audit mutation
that survived without it:

- `crs-two-acceptable.xml`: a certificate with TWO acceptable revocations at
  different production dates, so `CertificateRevocationSelector`'s "keep the
  latest acceptable one" comparison has to choose between candidates (reversing
  it, or dropping its null seed, changes the outcome).
- `crs-mixed-acceptance.xml`: a certificate with one acceptable and one
  unacceptable revocation, the only state in which the selector's two overridden
  message collectors branch (`collectMessages` drops the RAC constraint's
  messages, `collectAdditionalMessages` collects only the valid RAC's).
- `rac-anchor-mid-chain.xml`: a revocation whose own certificate chain carries a
  certificate behind its trust anchor, so that
  `RevocationAcceptanceChecker#initChain()`'s `break` at the first trust anchor
  differs from merely skipping the entry.
- `rac-passed.xml` (pre-existing) supplies the certHash-present shape the
  `certhash` policy needs.

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
  The direct test normalises that one field and nothing else (same as the `fc`/`sav`
  direct corpora).
- `XmlRAC/RevocationThisUpdate` and `XmlRAC/RevocationProductionDate` are required
  elements in `DetailedReport.xsd`, so the generated Go model carries them as plain
  (non-pointer) `XSDateTime` members where the JAXB class has nullable `Date`s. The
  zero time stands in for Java's null, and the block corpus renders a zero one as
  `null`. A revocation without a production date is therefore marshalled by the Go
  detailed report as `0001-01-01T00:00:00Z` where Java omits the element - a
  deviation of the frozen `detailedreport/jaxb` package, not of this one.

## XCVB direct corpus (`xcv_direct.jsonl`)

`gen/XcvOracle.java` covers the 83-file
manifest ported into this package's `checks` / `sub.checks` / `sub.checks.pseudo`
/ `rfc.checks` half (everything **except** the `X509CertificateValidation` /
`crs` / `rac` root that XCVA owns and that `XcvaOracle`/`XcvaDirectOracle`
already cover above). It follows the same "drive every check alone through a
one-item chain at `Level.FAIL`" pattern as `FcSavDirectOracle` and
`XcvaDirectOracle`.

Of the 83 ported files, 75 are directly instantiable `ChainItem`s and get a row
here; the other 8 are exercised only indirectly:

- `SubX509CertificateValidation` (the per-certificate `Chain` that wires most of
  these checks together) and `RevocationFreshnessChecker` (the `Chain` wiring
  `NextUpdateCheck`/`RevocationDataFreshCheck`/`RevocationDataFreshCheckWithNullConstraint`)
  are `Chain`s, not `ChainItem`s; their own constituent checks are all driven
  directly here, and `SubX509CertificateValidation` itself is exercised at block
  level by XCVA's `xcva_blocks.jsonl` (`X509CertificateValidation` calls it for
  every certificate of the chain).
- `PolicyTreeNode` is a plain RFC 5280 policy-tree data structure with no
  `ChainItem` surface of its own; it is exercised through
  `CertificatePolicyTreeCheck`, which is driven directly (including a
  synthetic anyPolicy/qualifier scenario reaching its node-deletion branch via
  the real corpus's varied certificate-policy chains).
- `PseudoStrategy` (interface), `PseudoAttributeStrategy`, `PseudoGermanyStrategy`
  and `JoinedPseudoStrategy` are exercised through `PseudoUsageCheck`, driven
  directly over the real corpus (which contains no `:PN`-suffixed German common
  names, so `PseudoGermanyStrategy`'s own branch is a Go-side unit-test
  responsibility, not this oracle's).
- `AbstractRevocationFreshCheck` is `abstract`; both of its concrete
  subclasses (`RevocationDataFreshCheck`, `RevocationDataFreshCheckWithNullConstraint`)
  are driven directly.

`RevocationDataFreshCheck` is driven twice per (certificate, revocation) pair -
a loose `DurationRule` then a tight one, to reach both its OK and NOT OK branch
over the very same revocation - so it is the one check for which
`(file, token, check)` is not a unique key in this file: two consecutive lines
share it. `../xcv_direct_oracle_test.go` queues rows per key instead of
keying on the last one written, and replays them in the same loose-then-tight
order, so both survive the KAT.

`../xcv_direct_oracle_test.go` buckets the "every check has an OK and a NOT OK
row" invariant by check *class*, stripping the `-any`/`-none`/`-low`/`-high`/
`-SIGNING_CERT`/`-CA_CERTIFICATE` variant suffixes `check()` carries in this
file (e.g. `CommonNameCheck-any` and `CommonNameCheck-none` both count toward
`CommonNameCheck`) - a variant is deliberately one-sided (`-any` always OK,
`-none` always NOT OK), only the class as a whole is guaranteed both.

Writing this replay surfaced two genuine bugs in the ported checks (both
fixed, not just worked around in the test): `CertificateForbiddenExtensionsCheck`
and `CertificateSupportedCriticalExtensionsCheck` passed a raw `[]string` to
`BuildXmlMessage`, which the frozen `i18n` package's `messageFormatArgString`
falls back to `fmt.Sprint` for (rendering `[a b c]`, space-separated) where
Java's `MessageFormat` calls the `List`'s own `toString()` (`[a, b, c]`,
comma-separated) - both now pre-render the list themselves before the call.
`ExtendedKeyUsageCheck.BuildAdditionalInfo` substituted `""` for a `null`
`XmlOID` description, where Java's `Arrays.toString` prints the literal text
`"null"` for a null element; the `[]string` `Process()` feeds
`ProcessValuesCheck` is unchanged (an empty string is not a value any real
policy constraint accepts, so this could not affect a verdict), only the
message rendering was fixed. The `i18n.messageFormatArgString` gap itself is
not fixed here - any other call site in the tree that passes a raw
`[]string`/`[]T` to `GetMessage`/`BuildXmlMessage` has the same latent bug.
Open follow-up.

### Why some rows are synthetic

The real corpus - even widened with XCVA's `dd/*.xml` dumps, which this
generator also replays (as `dd/<name>` file entries) since several of them
(trust anchors, sunset dates) exercise the same "no corpus certificate is
trusted" gap XCVB's checks hit too - still cannot reach every branch:

- **No corpus certificate carries a QC/PSD2 statement or a pseudonym**, so the
  nine `CertificateQc*`/`CertificatePS2DQc*` checks and `PseudonymCheck` only
  ever see the "absent" `NOT OK` branch under an `ANY` (`"*"`) rule. A single
  hand-built `XmlCertificate` (`C-SYNTH-QC`) carrying every one of those
  statements at once supplies the missing `OK` row for all of them.
- **No corpus certificate's DN falls outside a name constraint its issuer
  actually declares** (`CertificateNameConstraintsCheck`'s only failure
  branch): a two-certificate synthetic chain (`C-SYNTH-NC-CA` restricting
  `directoryName` to `O=Allowed`, `C-SYNTH-NC-LEAF` with DN `O=Excluded`)
  supplies it.
- **No corpus certificate declares `noRevAvail` while still publishing a
  conflicting OCSP access point** (`NoRevAvailCheck`'s only failure branch,
  RFC 9608 §3): `C-SYNTH-NORA-VIOLATION` supplies it.
- **Every corpus certificate carries a serial number**
  (`SerialNumberCheck`'s only failure branch): `C-SYNTH-NO-SERIAL` (serial
  number left unset) supplies it.
- **No corpus certificate is associated with a `TrustServiceProvider`**
  (`TrustServiceStatusCheck`/`TrustServiceTypeIdentifierCheck`'s only "match"
  branch): `C-SYNTH-TRUST-SERVICE` (one `TrustService` with a status/type an
  `ANY` rule matches, `StartDate` before the fixed usage time) supplies it.
- **No corpus revocation carries reason `CERTIFICATE_HOLD`**
  (`CertificateNotOnHoldCheck`'s only failure branch): `R-SYNTH-ON-HOLD`
  supplies it.
- The `XmlAOV`/`XmlCRS`/`XmlRFC`/`XmlSubXCV` "result wrapper" checks
  (`CertificateAlgorithmObsolescenceValidationCheck`,
  `CertificateRevocationSelectorResultCheck`,
  `RevocationFreshnessCheckerResultCheck`, `CheckSubXCVResult`) read an
  already-built result object rather than a diagnostic-data token; each is
  driven over hand-built result shapes keyed to a real certificate id, at
  every `Indication` the check branches on (`PASSED` / `INDETERMINATE`, plus
  `passed-with-algo`/`error`/`warning` for the AOV shape, mirroring
  `FcSavDirectOracle`'s `aovOfShape`).

Every one of the 75 checks has both an `OK` and a `NOT OK` row; the Go test
fails if either property is lost for any of them (same rule as XCVA's corpus
above and the `sav`/`fc` corpora).

### Regenerating `xcv_direct.jsonl`

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle <dss-repo>/validation/process/bbb/xcv/testdata/gen/XcvOracle.java
    java  -cp "$CP:/tmp/oracle" XcvOracle <dss-repo>/diagnostic/jaxb/testdata/oracle <dss-repo>

`XcvOracle` reads XCVA's `dd/*.xml` (already committed; no `XcvaSyntheticDumps`
run needed) alongside the marshal-parity corpus.
