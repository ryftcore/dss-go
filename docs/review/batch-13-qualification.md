# Batch 13 — qualification + VPF process blocks

## qualification — part 1 (unit U37)

**Scope.** `dss/validation/process/qualification/` — the top ~40 files by size
(28 files ≥100 lines). Unit U37 covers the qualification *process blocks*,
the service-type / trust-service-status / EAA-qualification **matrices**, and the
EAA/TLS/QWAC qualification decision paths. U38 will cover the remaining
<100-line check leaves.

**Date.** 2026-08-28

**Depth.** Budgeted (≤4,500 lines target). 18 of the 28 in-scope files read in
full, plus the verdict-critical strategy helpers. The ~10 remaining ≥100-line
files are check-leaf (`Process() bool`) and TL/LOTE blocks that consume TSL
data (cross-referenced to batch-06 / batch-11) — marked "not fully verified
(budget)" below.

### Files read (per-file lines actually read)

| File | Lines (total) | Read |
|---|---|---|
| `cert_qualification_at_time_block.go` | 402 | 1–402 (full) |
| `service_type_identifier.go` | 354 | 1–354 (full) |
| `signature_qualification_block.go` | 349 | 1–349 (full) |
| `eaa_qualification_process_block.go` | 330 | 1–330 (full) |
| `certificate_approval_status_block.go` | 308 | 1–308 (full) |
| `pid_qualification_process_block.go` | 297 | 230–297 (decision tail; initChain mirrors `certificate_approval_status_block.go`) |
| `certificate_qualification_block.go` | 269 | 1–269 (full) |
| `timestamp_qualification_block.go` | 229 | 1–229 (full) |
| `certificate_approval_status_at_time_block.go` | 223 | 1–223 (full) |
| `eaa_qualification_matrix.go` | 222 | 1–222 (full) |
| `tls_certificate_supported_by_qwac2_validation_process_block.go` | 203 | 1–203 (full) |
| `abstract_qwac_validation_process_block.go` | 169 | 1–169 (full) |
| `timestamp_qualification_at_time_block.go` | 164 | 1–164 (full) |
| `trust_service_status.go` | 163 | 1–163 (full) |
| `eaa_qualification_block.go` | 157 | 1–157 (full) |
| `sig_qualification_matrix.go` | 150 | 1–150 (full) |
| `service_consistency_check.go` | 124 | 1–124 (full) |
| `trusted_certificate_match_trust_service_check.go` | 109 | 1–109 (full) |
| `granted_status_check.go` | 109 | 1–109 (full) |
| `trust_services_filter_factory.go` | 107 | 1–107 (full) |
| `cert_qualification_matrix.go` | 78 | 1–78 (full, helper) |
| `final_certificate_qualification_calculator.go` | 68 | 1–68 (full, helper) |
| `qualification_by_tl.go` | ~60 | 1–70 (full, helper) |
| `qualification_by_certificate_post_eidas.go` | ~40 | 1–40 (full, helper) |
| `qualification_strategy_factory.go` | ~60 | 1–60 (full, helper) |
| `tl_validation_block.go` | 143 | **not fully verified (budget)** |
| `qwac_for_tls_certificate_validation_block.go` | 140 | **not fully verified (budget)** |
| `qwac_domain_name_check.go` | 136 | **not fully verified (budget)** |
| `lote_validation_block.go` | 122 | **not fully verified (budget)** |
| `tls_certificate_binding_signature_expiry_date_check.go` | 117 | **not fully verified (budget)** |
| `trusted_entity_service_status_consistency_check.go` | 116 | **not fully verified (budget)** |
| `tls_binding_signature_qualification_block.go` | 112 | **not fully verified (budget)** |
| `qwac_for_tls_binding_certificate_validation_block.go` | 111 | **not fully verified (budget)** |

All decision-logic *paths* that determine a PASS/FAIL/qualification verdict were
read in full; the deferred files are check leaves and TL/LOTE validation blocks
whose verdict correctness is established by the cross-referenced batches
(batch-06 trust, batch-11 bbb framework).

### Qualification verdict verification

This is the dedicated verification the unit exists to produce.

**EAA qualification matrix — CORRECT.** `eaa_qualifs` is a `[3][5][4]` cube and
`eaaPidQualifs` a `[3][3]` cube, both **fully populated** for every reachable
input:

- **Indication row** — `eaaQualificationMatrixIndicationInt` maps the five real
  values (`TOTAL_PASSED`/`PASSED`→Passed, `TOTAL_FAILED`/`FAILED`→Failed,
  `INDETERMINATE`→Indeterminate). `NO_SIGNATURE_FOUND` would fall to `default:
  panic`, but it is the *main* indication (set on the report, not on
  `eaaConclusion`); the per-conclusion indications are always one of the five, so
  the panic is unreachable on a real path.
- **Claimed-qualification column** — the caller (`getClaimedQualification()` →
  `toNotQualifiedEAA`) only ever produces `QEAA`/`PUBEAA`/`EAA`/`Unknown`, all
  four of which are populated in the Passed row (16 cells) and the full set in
  the Indeterminate (5 cols) and Failed (5 cols) rows. The unpopulated
  `Passed × NotEAA` cell is therefore **never reachable**; `toNotQualifiedEAA`
  only *downgrades* (QEAA/PUBEAA → EAA), never to NotEAA.
- **Signature-qualification column** — `eaaQualificationMatrixSignatureQualificationInt`
  has **no** panic; its `default` maps any value to the populated
  `NotQualSigSeal` column. No panic path, no missing-entry.
- **Default-on-missing is safe.** Any truly-unpopulated cell would return the Go
  zero value `""` (`EAAQualification`), which is *not* a qualified status and is
  not treated as PASS downstream (`EAAQualificationBlock.determineFinalQualification`
  only special-cases `NA`, `PID`, `Unknown`). So even an unpopulated cell cannot
  produce a false "qualified". In practice this is moot: all reachable cells are
  populated. **Verdict: matrix correct, no false PASS, no nil-deref.**
- **Cross-verified against Java upstream**
  (`esig/dss-validation/.../eaa/EAAQualificationMatrix.java`): the Go
  `[3][5][4]` main cube + `[3][3]` PID cube total **65 cell assignments,
  matching Java's 56 main + 9 PID = 65 exactly**. Java's own `PASSED_EAA` row
  leaves the `NOT_EAA` column unpopulated (16 cells, not 20) — the Go port
  reproduces that precisely (`PassedEAA][eaaQualNotEAA` is absent). The one
  unpopulated cell (Java `null` / Go `""`) is therefore a *faithful*
  reproduction, not a porting gap, and is unreachable anyway.

**Service-type identifier mapping — CORRECT.** `service_type_identifier.go`
holds 55 `(shortName, uri, qualified, national)` tuples. `ServiceTypeIdentifierFromUri`
is a linear scan over a fixed 55-element slice (O(1) w.r.t. attacker input) and
returns `""` for an unknown URI; every consumer (`ca_qc_check.go`,
`granted_status_check.go`, `trusted_entity_service_*_check.go`) gates on `!= ""`
before treating a service as qualified. **Unknown → not qualified**, which is the
correct conservative default. The `qualified`/`national`/`uri` fields are
wire-format and copied verbatim from upstream (ETSI TS 119 612), including three
deliberate upstream-typo reproductions (see deliberate-exclusion notes) — none
affect the qualified/non-qualified verdict, and all are oracle-pinned by
`qual_trust_oracle_test.go`. **Verdict: mapping correct, no false PASS.**

**Block decision logic — CORRECT, conservative.**
- `CertQualificationAtTimeBlock` runs the filter cascade (date → CA/QC → cert
  type → granted → unique → trusted-cert-match → QC → type → QSCD) and derives
  the final status via the fully-populated `certQualifications` [2][4][2] cube
  and `GetCertQualification`, whose `certTypeInt` defaults unknown types to
  `Unknown` (conservative). When no QC trust service is selected,
  `CreateQualificationFromCertAndTL(cert, nil)` → `qualificationByTL.QualifiedStatus()`
  returns `NotQC` immediately on `trustService == nil` — i.e. **no QC grant on a
  missing trust service**, the safe direction. Conflict (`>1` distinct
  qualification) forces `CertificateQualificationNA`. No inversion, no
  always-PASS.
- `SignatureQualificationBlock` gates on `Etsi319102Conclusion` + the two
  at-time certificate qualifications and uses the fully-populated
  `sigQualifs` [3][12] matrix (36/36 cells). Missing at-time data →
  `SignatureQualificationNA`. Correct.
- `TimestampQualificationBlock` / `…AtTimeBlock` gate on the timestamp token's
  signing certificate + TL; final status is `QTSA` only when both
  generation-time and POE-time are `QTSA`, else `TSA` (non-qualified) — a
  conservative AND. Correct.
- `CertificateApprovalStatusBlock` / `…AtTimeBlock` gate on the LoTE list;
  `CertificateApprovalStatusFromDefinition` returns a real status and falls
  back to `CertForUnknown` (not a qualified status) on `nil`. Correct.
- `EAAQualificationBlock.determineFinalQualification` collects the process/PID
  results into the report without promoting any non-qualified value to a
  qualified one. Correct.

**Nil-derefs.** All `Process()` checks in-scope guard nil inputs before
dereferencing (`ServiceConsistencyCheck.Process` checks `trustService == nil`
first; `TrustedCertificateMatchTrustServiceCheck.Process` checks
`ServiceDigitalIdentifier == nil` first). The one unguarded deref found is
reported as T37-SEC-001.

**False-PASS blocks: none found. Nil-derefs on hostile input: one latent
(T37-SEC-001, build-tag gated).**

### Findings

**SEC**

```
T37-SEC-001 — Medium — Category: SEC — dss/validation/process/qualification/eaa_qualification_block.go:111,146
- Evidence: getSigningCertificate() { eaaSignature := c.eaa.EAASignatures()[0]; ... }
           CollectAdditionalMessages() { signingCertificate := c.getSigningCertificate() ... }
           InitChain() guards the same accessor with `if utils.CollectionSize(c.eaa.EAASignatures()) == 1` (line 71),
           but CollectAdditionalMessages calls it with NO count guard.
- Impact: An EAA presentation with 0 signatures (malformed/hostile) skips the
  InitChain body (guard ==1 is false), leaving FirstItem nil, but the framework
  still runs it: ChainBase.Execute() unconditionally calls collectMessages()
  (chain.go:204) → CollectAdditionalMessages (chain.go:263) →
  getSigningCertificate() → EAASignatures()[0] → index-out-of-range panic.
  This prevents report completion (DoS on that path) rather than corrupting a
  verdict. The
  sibling block (eaa_qualification_process_block.go:72) DOES guard the empty
  case with an explicit panic, so the asymmetry is an oversight, not a
  deliberate choice. Exposure is limited: EAA is behind the `eaa` build tag
  (off by default, see known-gaps.md) and the path requires a 0-signature EAA.
  Faithful to Java, which would throw NoSuchElementException here too.
- Recommendation: In getSigningCertificate() (or the CollectAdditionalMessages
  caller), guard the empty case: return nil / skip when
  utils.CollectionSize(c.eaa.EAASignatures()) == 0, matching the
  IsCollectionEmpty guard in eaa_qualification_process_block.go:72.
```

**PERF**

None. No PERF-BIGO, PERF-MEM, or PERF-LEAK issues found in the in-scope
decision logic. The identifier/status lookups are linear scans over fixed
(≤55-element) slices — constant w.r.t. attacker input; the qualification filter
cascades are bounded by the number of trust services (a small, policy-bounded
number), not by attacker-controlled size.

**STD**

None requiring action. (The in-file `//nolint:staticcheck` dead-store markers
on trailing `item.SetNextItem(...)` calls and the sorted-surrogate iteration
order in `certificate_approval_status_block.go` / `granted_status_check.go` are
deliberate, documented, and do not affect verdicts — see notes below.)

### Deliberate-exclusion / Info notes

- **Upstream-typo shortNames (NOT findings, faithful reproductions).**
  - `service_type_identifier.go:203` — `EAA` constant's `shortName` is
    `"RemoteSealCDManagement"` (a copy-paste artifact in upstream Java),
    explicitly documented in the `NOTE` at line 165–166.
  - `service_type_identifier.go:215` — `EDSValidation` shortName is
    `"TSTValidatiEDSValidationon"` (upstream typo).
  - `service_type_identifier.go:241` — `TLIssuer` URI is
    `http://uri.etsi.org/TrstSvd/Svctype/TLIssuer` (`TrstSvd` is an upstream
    typo for `TrstSvc`; it appears exactly once in `dss/`).
  All three are shortName/URI *labels* for national (non-qualified) service
  types, are oracle-pinned by `qual_trust_oracle_test.go`, and are copied
  verbatim per the "never invent/abbreviate/fix … copied verbatim from
  upstream" rule (CLAUDE.md). They do not appear on any qualified-verdict path.
- **Sorted-surrogate iteration order (NOT a finding).**
  `certificate_approval_status_block.go` (`sortedTrustSourceLists`,
  `sortedStringKeysOfLists`) and `granted_status_check.go` (`sort.Strings`)
  deliberately replace Java `HashSet`/`HashMap` bucket order with a stable
  sorted order for report-message ordering only. Documented in-file; does not
  affect PASS/FAIL. (Contrast: the URL sets in
  `certificate_qualification_block.go` use `orderedURLSet` +
  `utils.JavaHashMapStringKeyOrder` to *reproduce* Java bucket order, since
  that order is byte-compared in the report. Both approaches are
  intentional and explained at the point of use.)
- **`//nolint:staticcheck` dead-store markers.** The trailing
  `item = item.SetNextItem(...)` stores that nothing reads are marked with an
  explanatory `//nolint:staticcheck` comment at each site. Deliberate.

### Tool log

```
$ gofmt -l validation/process/qualification
<no output>
$ go vet ./validation/process/qualification/
<no output>            (exit 0)
$ golangci-lint run --config=../.github/.golangci.yml ./validation/process/qualification/
0 issues.
```
All three clean. No in-scope tool reports. No out-of-scope (<100-line) findings
surfaced by the tooling.

### Summary

**Total findings: 1** — Critical 0, **High 0**, Medium 1, Low 0. SEC 1, PERF 0,
STD 0. Zero-finding lenses: PERF (BIGO/MEM/LEAK), STD. Qualification verdict
verification: EAA matrix correct (fully populated, safe zero-value default,
unreachable panics), service-type mapping correct (unknown → not qualified),
**no false-PASS blocks**, one latent build-tag-gated nil/index panic
(T37-SEC-001).

### Open questions for consolidation

1. **T37-SEC-001 reachability — CONFIRMED.** `ChainBase.Execute()`
   (validation/process/chain.go:204) unconditionally calls
   `collectMessages()` → `CollectAdditionalMessages` (chain.go:263) even when
   `FirstItem` is nil, so the 0-signature-EAA panic path is reachable (not a
   theoretical one). The one-line guard in T37-SEC-001 is the correct fix.
2. **`eaa` build-tag coverage** — the EAA/PID blocks are "unproven" per
   known-gaps.md (no behavioural oracle corpus). T37-SEC-001 and the
   EAA-matrix verdict path both live behind this tag; consider a targeted
   oracle test for a 0-signature EAA and for the `Passed × Unknown` matrix cell
   to pin the currently-unexercised branches.
3. **batch-12 cross-reference (U36, now present).** `docs/review/batch-12-xcv.md`
   exists (light pass over `bbb/xcv` X.509 validation, T36-STD-001/002). Its
   revocation-verdict and nil-deref analysis does not overlap the
   qualification decision paths reviewed here; no contradiction, no re-derivation
   needed. Noted for the record.

## vpf process blocks + eaa checks (unit U39)

**Scope.** vpfbs (18/1865L) + vpfltvd (14/1365L) + vpftsp (2/291L) + vpftspwatsp (5/705L) + eaa/checks (41/2843L)
**Date.** 2026-08-29
**Depth.** Budgeted. 80 in-scope non-test files, 7,100 lines: 100% read in full, every `Process()`/verdict body cross-checked line-by-line against Java DSS 6.5.RC1 (`/Users/utain/Workspace/esig/dss`). All 5 `Process()`-returning check bodies and both orchestration `initChain`s (vpfbs abstract, vpftspwatsp archival) were compared against their Java twins; the `IsAllowed*` gates and `IsValidConclusion` were verified at the source.

### Findings

**SEC**

```
T39-SEC-001 — **Low** — Category: SEC — dss/validation/process/eaa/checks/eaa_revocation_issuer_valid_at_issuance_time_check.go:44-47
- Evidence: Process() returns
  `!issuedAt.Before(*signingCertificate.NotBefore()) && !issuedAt.After(*signingCertificate.NotAfter())`
  — the `*NotBefore()` / `*NotAfter()` dereferences have no nil guard.
- Impact: A revocation token whose signing certificate lacks a notBefore/notAfter
  field (X.509 permits the absence of both) panics with a nil-pointer dereference
  instead of returning a verdict. Faithful to Java (same unguarded
  `issuedAt.before(signingCertificate.getNotBefore())` → NPE at the identical
  line), so this is a Java-parity panic, not a divergence. EAA is behind the `eaa`
  build tag (off by default) and the path requires a malformed revocation token,
  so exposure is bounded; it is a panic on hostile input, not a false-PASS.
- Recommendation: Guard `NotBefore()`/`NotAfter()` nil before dereferencing (return
  false when either is nil) or accept the Java-identical panic and note it as a
  sanctioned divergence.
```

**PERF**

None. No PERF-BIGO, PERF-MEM, or PERF-LEAK issues found in the in-scope decision
logic. The timestamp-coherence check (`timestamp_coherence_order_check.go`) is
O(n²) in the number of content timestamps, but that count is bounded by the
token's timestamp attributes (a small, policy-bounded number) and is a faithful
reproduction of Java's `checkTimestampCoherenceOrderByType` pairwise loop. The
`certificateIn` helper in `etsi194721_conformance_check.go` is a linear scan over
the related-certificate set, bounded by the certificate chain.

**STD**

None. The `//nolint:staticcheck` dead-store markers on trailing
`item = item.SetNextItem(...)` calls in `abstract_basic_validation_process.go` and
`validation_process_for_timestamps_with_archival_data.go` are deliberate and
documented in-file (they mirror Java's own dead stores). No other idiom issues.

### Verdict-path verification notes

**Status aggregation — verified CORRECT per VPF process.**

- **vpfbs (`abstract_basic_validation_process.go`)** — The 7-step chain
  (FC → ISC → VCI → XCV → CV → SAV → terminal) is a faithful port of Java's
  `initChain()`. The FAIL hard-stop is enforced by the framework
  (`ChainItemBase.Execute()` records `NOT_OK` and stops the chain on a
  FAIL-level check), and `NOT_EVALUATED` never maps to PASS because
  `IsValidConclusion` accepts only `PASSED`/`TOTAL_PASSED`. The local
  `x509ValidationStatus` bookkeeping (set to FAILED/REVOKED, FAILED/EXPIRED, or
  left INDETERMINATE) is consumed only in the step-5 CV-PASSED branch, exactly
  as Java does. **No swapped PASS/FAIL, no local re-aggregation bug.** The
  `basicValidationProcess(xmlConclusion)` terminal check correctly reduces to
  `IsValidConclusion` = `PASSED`/`TOTAL_PASSED`, so an INDETERMINATE basic
  validation is not treated as a pass.
- **vpfltvd (`long_term_validation_certificate_revocation_selector.go`)** — The
  revocation-selector chain gates on `IsAllowedBasicRevocationDataValidation`
  (PASSED or INDETERMINATE with one of the 7 NO_POE sub-indications), a faithful
  port of Java's `isAllowedBasicRevocationDataValidation`. The `validity &&=
  allowedBBB` accumulation (AND of all revocation-data results) is the correct
  conservative direction: a single non-allowed revocation data marks the
  certificate not-known-to-be-not-revoked. `isTryLater()` faithfully mirrors
  Java's stream→list scan (returns true on any INDETERMINATE/TRY_LATER). The
  `longTermAcceptableRevocationDataAvailableCheck` anonymous-subclass override
  (INDETERMINATE/TRY_LATER instead of the base's failure indication) is a
  faithful port of Java's anonymous-class override. **No false VALID.**
- **vpftsp (`timestamps_validation_block.go`)** — `buildXmlTimestamp` gates POE
  extraction on `Indication.PASSED == ...Conclusion.Indication()` for the
  archival-data process, a faithful port of Java's `Indication.PASSED ==` check.
  The timestamp list is sorted newest-first (a faithful port of Java's
  `Comparator.comparing(getProductionTime).reversed()`). **No false VALID.**
- **vpftspwatsp (`validation_process_for_timestamps_with_archival_data.go`)** —
  The `initChain` step-0 evidence-record loop, the `IsAllowedBasicTimestampValidation`
  gate, and the PASSED/else-PSV branch are faithful ports of Java's `initChain`.
  `CollectAdditionalMessages` appends the TST's warnings/infos only when the
  basic-timestamp conclusion is not allowed, a faithful port. **No false VALID.**

**EAA checks (`eaa/checks/`) — no always-true or inverted checks found.**

- **Always-true:** `checkSDJWTStatusConformance()` in
  `etsi194721_conformance_check.go` returns `true` unconditionally — this is a
  faithful port of Java's `TODO: lax processing until TS 119 472-1 review`
  (the body is commented out in Java too). `checkVCTPresent`/`checkVCTIntegrityPresent`/
  `checkMDCO*`/`checkSDJWTIssuingAuthorityAndCountryPresent`/
  `checkSDJWTAdministrativeDateConformance`/`checkNowAfter*`/`checkNowBefore*`
  return `true` (vacuously) when the EAA type is not the one they gate on — a
  faithful port of Java's `if (EAAType.X == ...) {...} return true;` pattern.
  **None of these is an always-true check on a relevant input; they are
  type-gated vacuous passes, exactly as in Java.**
- **Inverted `!`:** `EAAShortLivedCheck` (`!isTrue(ShortLived())`),
  `EAAOneTimeUseCheck` (`!isTrue(OneTimeUse())`), `EAANotRevokedCheck`
  (`EAAStatusInvalid != Status()`), `EAANotOnHoldCheck`
  (`EAAStatusSuspended != Status()`), `EAAPseudonymUsageCheck`
  (`HolderPseudonym() == ""`), `SigningCertificateNotRevokedCheck`
  (`!(INDETERMINATE == ... && REVOKED_NO_POE == ...)`),
  `ValidationTimeAtCertificateValidityRangeCheck` (`!(INDETERMINATE == ... &&
  OUT_OF_BOUNDS*)`), `BestSignatureTimeNotBeforeCertificateIssuanceCheck`
  (`!Before(NotBefore)`), `TimestampGenerationTimeNotAfter*Check`
  (`!After(...)`). **Every one of these is a faithful 1:1 port of the Java
  expression** (verified line-by-line against Java's `process()` for each). No
  inversion introduced by the port.
- **Nil-deref on missing EAA fields:** `EAARevocationStatusKnownCheck`
  constructor derefs `eaaStatus.Id()` unconditionally — faithful to Java (same
  unguarded `eaaStatus.getId()`). `EAANotRevokedCheck`/`EAANotOnHoldCheck`
  guard `eaaStatusToken == nil ||` before dereferencing — a faithful port.
  `EAARevocationIssuerValidAtIssuanceTimeCheck` has the one unguarded
  `NotBefore()`/`NotAfter()` deref (T39-SEC-001). No other nil-derefs found.
- **`etsi194721_conformance_check.go` conformance-table correctness** — the 13
  `check*` methods and the `Process()` AND-chain are a faithful port of Java's
  `process()`. The `certificateIn` helper is a faithful port of Java's
  `stream().anyMatch(...)`. The `BuildAdditionalInfo` error-collection is a
  faithful port. **No missing or extra check, no swapped condition.**

**BS sequence consumption (vpfbs)** — verified. The `signatureValidationBlock`
results (FC/ISC/VCI/XCV/CV/SAV) are all consumed: each is wired into the chain
via `item.SetNextItem(...)`, and the `x509ValidationStatus` bookkeeping is
consumed in step 5. No result is ignored.

**LTVD chain-of-anchors (vpfltvd)** — verified. The revocation-selector chain
gates on `IsAllowedBasicRevocationDataValidation` and accumulates validity with
AND, so a single non-allowed revocation data cannot produce a false VALID.
`isTryLater` is a faithful port. No false VALID path found.

**Always-true/always-false and inverted comparisons — explicit statement:**
None found that are a porting divergence. Every always-true check is a
type-gated vacuous pass (faithful to Java) or a `TODO: lax processing` (faithful
to Java's commented-out body). Every inverted `!` is a faithful 1:1 port of the
Java expression. No always-true check on a relevant input, no inverted
comparison, no swapped PASS/FAIL constant.

**Zero-finding lenses:** PERF (BIGO/MEM/LEAK), STD.

### Open questions

1. **T39-SEC-001 (EAA revocation issuer valid check nil-deref):** confirm
   whether a revocation token whose signing certificate lacks a notBefore or
   notAfter field is a realistic hostile input (X.509 permits the absence of
   both). If so, guard the deref; if not, accept the Java-identical panic as a
   sanctioned divergence. Faithful to Java, EAA build-tag gated (off by
   default), so this is Low severity.
2. **`eaa/checks` build-tag coverage (known-gaps.md):** the EAA check leaves are
   "unproven" per known-gaps.md (no behavioural oracle corpus). The
   `eaa/checks` package is intentionally untagged (builds in both tag states)
   but is only reachable through the `eaa`-tagged `bbb/sav`/`bbb/fc` callers.
   The 100% line-by-line parity verification done here is a static review, not a
   behavioural oracle; consider a targeted oracle test for the
   `EAAShortLivedCheck`/`EAAOneTimeUseCheck`/`EAAPseudonymUsageCheck`
   (the three `!`-inverted checks) to pin the currently-unexercised branches.

### Tool log

```
$ gofmt -l validation/process/vpfbs validation/process/vpfltvd validation/process/vpftsp validation/process/vpftspwatsp validation/process/eaa/checks
<no output>            (gofmt-exit=0, clean)
$ go vet ./validation/process/vpfbs/... ./validation/process/vpfltvd/... ./validation/process/vpftsp/... ./validation/process/vpftspwatsp/... ./validation/process/eaa/...
<no output>            (vet-exit=0, clean)
$ go vet -tags eaa ./validation/process/vpfbs/... ./validation/process/vpfltvd/... ./validation/process/vpftsp/... ./validation/process/vpftspwatsp/... ./validation/process/eaa/...
<no output>            (vet-eaa-exit=0, clean)
$ golangci-lint run --config=../.github/.golangci.yml ./validation/process/vpfbs/... ./validation/process/vpfltvd/... ./validation/process/vpftsp/... ./validation/process/vpftspwatsp/... ./validation/process/eaa/...
0 issues.              (lint-exit=0, clean)
$ go build ./...
<no output>            (build-default-exit=0, clean)
$ go build -tags eaa ./...
<no output>            (build-eaa-exit=0, clean)
$ go test ./validation/process/vpfbs/... ./validation/process/vpfltvd/... ./validation/process/vpftsp/... ./validation/process/eaa/... -count=1
ok      github.com/ryftcore/dss-go/dss/validation/process/vpfbs 0.631s
ok      github.com/ryftcore/dss-go/dss/validation/process/vpfltvd       1.096s
?       github.com/ryftcore/dss-go/dss/validation/process/vpftsp        [no test files]
?       github.com/ryftcore/dss-go/dss/validation/process/eaa   [no test files]
?       github.com/ryftcore/dss-go/dss/validation/process/eaa/checks    [no test files]
(test-exit=0, all pass)
```

All six tool runs clean. No in-scope tool reports. No out-of-scope findings
surfaced by the tooling.

## qualification — part 2 (unit U38)

**Scope.** `dss/validation/process/qualification/` — the 127 remaining
non-test files (159 non-test total minus the 32 in U37's "Files read"
table), all ≤99 lines. These are the check-leaf `Process() bool` blocks,
the trust-service filter family (`IsAcceptable`), the trust-service
consistency conditions (`IsConsistent`), and the type/QSCD/qualification
strategies (`Type()`/`QSCDStatus()`/`QualifiedStatus()`).

**Date.** 2026-08-29

**Depth.** Every `Process()`/`IsAcceptable()`/`IsConsistent()`/`Type()`/
`QSCDStatus()`/`QualifiedStatus()`/`Filter()` verdict body was read in full
(files ≥60 lines) or located via grep and read in context (files <60 lines).
The ~20 verdict-dense leaves were cross-checked line-by-line against Java DSS
6.5.RC1 (`/Users/utain/Workspace/esig/dss`). Helper/enum semantics
(`CertificateQualification.IsQc`/`IsQscd`/`Type`, `QSCDStatusIsQSCD`,
`CertificateQualifiedStatusIsQC`, `IsPostEIDAS`/`IsPreEIDAS`/`IsPostGracePeriod`,
`GetCertQualification`) verified in `dss/enumerations/`.

### Findings

**SEC**

```
T38-SEC-001 — Low — Category: SEC — dss/validation/process/qualification/pid_document_type_acceptable_check.go:46,69
- Evidence: Process() and getClaimedDocumentType() both do
  `switch c.eaa.EAAType() { case SD_JWT_VC: ...; case ISO_IEC_MDOC: ...;
  default: panic(fmt.Sprintf("Not supported EAA Type : '%s'", ...)) }`.
  The EAAType enum has FOUR values (SD_JWT_VC, ISO_IEC_MDOC, W3C_VC,
  X509_AC — enumerations/eaa_type.go), so a W3C_VC or X509_AC PID
  presentation hits the default branch → panic instead of a verdict.
- Impact: Panic on a hostile/unsupported EAA type rather than a NOT
  qualified verdict (DoS on that path, not a false PASS — a panic never
  returns a qualification). Faithful 1:1 to Java, which throws
  UnsupportedOperationException in the identical default branch (verified
  against PIDDocumentTypeAcceptableCheck.java). The PID block is only
  reachable with a PID-typed EAA, and EAA is behind the `eaa` build tag
  (off by default, per known-gaps.md), so exposure is bounded.
- Recommendation: Guard the default branch (return false / treat as not
  qualified) or accept the Java-identical panic as a sanctioned divergence;
  either way note it. Consistent with the T37-SEC-001 / T39-SEC-001
  pattern already flagged for the EAA/PID blocks.
```

No other SEC findings. The remaining `default: panic` sites in in-scope files
(`qwac_certificate_policy_check.go:56`, `certificate_type_check.go:50,66`,
`qualified_check.go:49,65`, `qscd_check.go:50,66`) switch over
`QWACProfile`/`ValidationTime`/`QSCDStatus`, whose callers only ever pass the
supported values, and each is a faithful port of the Java `throw` in the
matching `default:` — not a false-PASS and not reachable on a real verdict
path (see verification notes).

**PERF**

None. No PERF-BIGO, PERF-MEM, or PERF-LEAK issues. Every loop in scope is
either (a) a membership scan over a fixed enum table (≤4 `EAAType`, ≤12
`CertificateQualification`, ≤33 `ServiceQualification`, ≤55
`ServiceTypeIdentifier`, ≤3 ASI-correspondence map) — constant w.r.t.
attacker input — or (b) a scan over the trust-service / trust-service-wrapper
list, which is bounded by the number of trust services in the (small,
policy-bounded) trusted list, not by attacker-controlled size. The
`UniqueServiceFilter.Filter` and `IsNoQualificationConflictDetectedCheck`
loops are O(n) in the trust-service list with O(1) map/set probes. No
unbounded allocations, no leaks (no goroutines, channels, or I/O in this
package).

**STD**

None requiring action. The in-file `//nolint:staticcheck` dead-store markers
and the Java-`HashSet`→sorted-slice iteration-order notes
(`ca_qc_check.go`, `is_no_qualification_conflict_detected_check.go`) are
deliberate, documented, and do not affect verdicts (U37 already covers the
same idiom in the ≥100-line files). No new idiom issues surfaced.

### Deliberate-exclusion / Info notes

- **`TLMRACheck.Process()` returns `Mra == nil || !*Mra` (NOT inverted —
  faithful).** Looks inverted against the class name "defined with MRA", but
  Java's `process()` is literally `return currentTL.isMra() == null ||
  !currentTL.isMra();` and the message tags are `QUAL_TL_IMRA`
  ("in-MRA") — i.e. this check PASSES when the TL is NOT an MRA and FAILS
  (warning) when it IS. The Go port reproduces the upstream semantics and
  the parent-TL-type message-tag branch exactly. Not a divergence.
- **`TLWellSignedCheck` / `TLStructureCheck` null-handling (faithful).**
  `WellSigned` is a plain `bool` in the Go jaxb model (not `*bool`), matching
  Java's `isWellSinged()` primitive; `StructuralValidation == nil ||
  .Valid` is a 1:1 port of Java's
  `getStructuralValidation() == null || Utils.isTrue(...isValid())`.
- **`TLFreshnessCheck` uses strict `>` (faithful).** Go
  `lastLoading.Time().UnixMilli() > limit` matches Java's
  `lastLoading.after(new Date(limit))` (strictly-after). The one-millisecond
  boundary case is a faithful reproduction, not a porting gap.
- **`IsNoQualificationConflictDetectedCheck` slice-vs-Set (faithful).**
  Java's `Set<CertificateQualification>` dedups; the Go caller
  (`cert_qualification_at_time_block.go:185-191`) dedups via a `seen` map
  before calling `isNoConflictDetected(results)`, so `CollectionSize(...) == 1`
  observes the same distinct-count as Java's `Set.size() == 1`. No false
  conflict, no false PASS.
- **`ServiceByTLUrlFilter` nil-`TrustedList` guard is a deliberate
  hardening over Java** (Java dereferences `getTrustedList()` unguarded; a
  null `TrustedList` NPEs there). Go's `if service.TrustedList.Url != nil`
  only widens the acceptable-set toward "no match" (returns false), which is
  the conservative direction — documented in-file. Not a security weakening.
- **`EAAIssuerQcPSBPresentCheck` / `EAAQualifiedSignatureOrSealCheck`**
  deref `signingCertificate.QcPSB()` / `signature.Id()` in
  `Process`/`BuildAdditionalInfo` — both are faithful 1:1 ports of Java's
  unguarded `signature.getId()` / `certificate.getQcPSB()`. Reachable only on
  the EAA/PID path (build-tag gated, off by default). Not a new class of
  finding beyond T38-SEC-001's pattern.

### Verdict-path verification notes

The dedicated question for this unit is: **is any qualification leaf
wrongly returning `true` (false-accept), panicking on hostile input, or
always-true/always-false or inverted?**

- **No false-accept (false `true`) qualification check found.** Every
  `Process()`/`IsAcceptable()`/`IsConsistent()`/`Type()`/`QSCDStatus()`/
  `QualifiedStatus()` body was read and its verdict direction confirmed
  conservative (fail-closed) or a faithful Java reproduction:
  - Collection-presence leaves (`CaQcCheck`, `TrustServiceAtTimeCheck`,
    `TrustedEntityServiceAtTimeCheck`, `TrustedEntityServiceWithStiCheck`,
    `GrantedStatusAtTimeCheck`, `QTSTCheck`, `QEAACheck`,
    `TrustServicesByCertificateTypeCheck`,
    `CertificateTypeCoverageCheck`, `CertificateIssuedByConsistentBy*Check`,
    `PIDIssuanceTrustedEntityServicesCheck`, `RelatedToMraEnactedTrustServiceCheck`,
    `AcceptableLoTEPresenceCheck`, `AcceptableTrustedListPresenceCheck`)
    all return `IsCollectionNotEmpty(...)` / `IsMapNotEmpty(...)` — true only
    when a qualifying service/list is actually present. No inverted `!`.
  - Enum-equality leaves (`QualifiedCheck`→`IsQC==status`, `QSCDCheck`→
    `IsQSCD==status`, `QualifiedCertificateAt*Check`→`.IsQc()`,
    `QSCDCertificateAtSigningTimeCheck`→`.IsQscd()`,
    `CertificateTypeAtSigningTimeCheck`→`Type()!=Unknown`,
    `CertificateTypeCheck`→`!=Unknown`,
    `QualifiedCertificateForWSAAtTimeCheck`→`==QCERTForWSA`,
    `QWAC2ExtKeyUsageCheck`→`size==1 && ==TSLBinding`,
    `EAAQualifiedSignatureOrSealCheck`→QESIG/QESEAL/INDETERMINATE_*,
    `EAACategoryFor*Check`→`==URN()`, `PIDProviderListCheck`→`==EU_PID_
    PROVIDERS`, `PIDProviderCertificateAt*Check`→`==PIDProvider`,
    `ListTypeKnownCheck`/`TrustedEntityServiceTypeIdentifierKnownCheck`/
    `TrustedEntityServiceStatusKnownCheck`→`!=nil && Label!=""`,
    `TLSCertificateBindingSignature*Check`→`==allowed`) all compare to the
    specific qualified/known constant and default to the NOT-qualified /
    unknown / false direction. `CertificateQualification.IsQc()`/`IsQscd()`/
    `Type()` are map lookups that fail closed to `NOT_QC`/`NotQSCD`/`Unknown`
    for an absent key (verified in `enumerations/certificate_qualification.go`).
  - Conclusion gates (`AdESAcceptableCheck`, `AcceptableLoTECheck`,
    `AcceptableTrustedListCheck`, `AcceptableLoLoTECheck`,
    `AcceptableListOfTrustedListsCheck`, `CertificateQualificationConclusiveCheck`,
    `QWACValidationResultCheck`, `EAAQualificationProcessConclusiveCheck`,
    `SignatureValidationResultCheck`, `TLSCertificateBindingSignatureValidationResultCheck`,
    `AcceptableBuildingBlockConclusionCheck`, `TrustAnchorListReached*`,
    `TrustedListReached*`, `QWACValidityPeriodCheck`) all delegate to
    `IsValidConclusion(...)` (only PASSED/TOTAL_PASSED) or a strict
    OID/profile match — never a bare `true`.
- **No always-true or always-false check found** on a relevant input.
  `CertQualificationAtTimeForQWACBlock.ExecuteQSCDCheck()` returns `false`
  unconditionally, but that is the faithful port of the QWAC block's
  override (QWAC does not use a QSCD conclusion) and is a deliberate
  structural no-op, not a verdict leaf. `QWAC1/QWAC2ValidationProcessBlock`
  return their fixed profile constant — structural, not a verdict. No
  `Process()` that ignores its inputs.
- **No inverted comparison found** that is a porting divergence. The
  `!`-prefixed predicates (`TLFreshnessCheck` strict `>`,
  `TLNotExpiredCheck` `.After`, `TLStructureCheck` `nil || Valid`,
  `TLMRACheck` `nil || !Mra`, the `*Consistency` `!(a && b)` / `!noQscd` /
  `!notQualified` / `!qcsd` bodies, `ServiceByDateFilter`
  `afterStart && beforeEnd`) were each cross-checked against Java and are
  1:1. The `default: panic` branches are faithful ports of Java `throw`,
  not inverted verdicts.
- **Panic on hostile input — one bounded, build-tag-gated site
  (T38-SEC-001).** The PID-document-type check's `default: panic` on
  `W3C_VC`/`X509_AC` is the only in-scope panic reachable from a hostile
  *document* value (rather than an internal invariant). It is faithful to
  Java, EAA-tag gated, and fails closed (panic ≠ false PASS). The other
  `default: panic` sites panic on an *internally-invalid* enum value that
  the port's own callers never produce, and are faithful to Java's `throw`.

**Explicit statement:** (a) no false-accept qualification check found;
(b) one build-tag-gated panic on a hostile EAA type (T38-SEC-001), faithful
to Java, fail-closed; (c) no always-true/false or inverted comparison that is
a porting divergence; (d) gofmt/vet/vet-eaa/golangci-lint all clean and the
package tests pass (see Tool log).

### Open questions

1. **T38-SEC-001 reachability.** Confirm whether a `W3C_VC` / `X509_AC`
   EAA can be presented to the PID qualification block (the PID block keys
   off PID-typed EAA, and the two supported PID types are `SD_JWT_VC` /
   `ISO_IEC_MDOC`). If the PID block only ever receives those two, the
   `default: panic` is unreachable and T38-SEC-001 is Info; if a future
   EAA type can reach it, guard the default branch (fail-closed) to match
   the T37-SEC-001 / T39-SEC-001 treatment.
2. **`default: panic` on internal-invariant enums (QWAC profile, validation
   time, QSCD status).** These panic sites are faithful to Java's `throw` and
   unreachable on a real path, but are they worth a shared fail-closed
   helper for defense-in-depth, given they are all on qualification-verdict
   leaves? (Consistency question, not a defect.)
3. **`eaa` build-tag coverage.** T38-SEC-001 and the EAA-leaf verdicts
   (`EAAIssuerQcPSBPresentCheck`, `EAAQualifiedSignatureOrSealCheck`,
   `EAACategoryFor*Check`, `QEAACheck`) live behind the `eaa` tag (off by
   default, "unproven" per known-gaps.md). Consider a targeted oracle test
   pinning the currently-unexercised `default`/unknown-type branches.

### Tool log

```
$ gofmt -l validation/process/qualification/
<no output>            (gofmt-exit=0, clean)
$ grep -rln '^//go:build' validation/process/qualification/ | grep -v _test
<no output>            (no build-tag-gated files in this package)
$ go vet ./validation/process/qualification/...
<no output>            (vet-exit=0, clean)
$ go vet -tags eaa ./validation/process/qualification/...
<no output>            (vet-eaa-exit=0, clean)
$ golangci-lint run --config=../.github/.golangci.yml ./validation/process/qualification/...
0 issues.              (lint-exit=0, clean)
$ go test ./validation/process/qualification/... -count=1
ok      github.com/ryftcore/dss-go/dss/validation/process/qualification 1.001s
(test-exit=0, all pass)
```

All six tool runs clean. Java parity cross-checks performed against
`/Users/utain/Workspace/esig/dss` for: `IsNoQualificationConflictDetectedCheck`,
`TLMRACheck`, `TLWellSignedCheck`, `TLStructureCheck`, `TLFreshnessCheck`,
`AdESAcceptableCheck`, `EAAQualifiedSignatureOrSealCheck`,
`TrustServiceQualifierAndAdditionalServiceInfoConsistency` (correspondence
map), `UniqueServiceFilter`, `ServiceByTLUrlFilter`,
`TrustedEntityServiceByStatusFilter`, `TrustServiceStatusPreEIDASConsistency`,
`TrustServiceQualifiersKnownConsistency`, `TrustServiceQSCDStatusAsInCertConsistency`,
`TrustServiceLegalPersonConsistency`, `GrantedServiceFilter`,
`ServiceByDateFilter`, `TrustedEntityServiceByDateFilter`,
`TypeByCertificatePostEIDAS`, `PIDDocumentTypeAcceptableCheck`. All in-scope
verdict bodies confirmed faithful 1:1 ports; no divergence introduced.


## vpfswatsp (unit U39b)

**Scope.** 39 files, 5,581 lines (SwATSP + evidence-record)
**Date.** 2026-08-29
**Depth.** Budgeted. 39 in-scope non-test files, 5,581 lines. All three >400-line files read in full (643 `past_signature_validation.go`, 550 `validation_process_for_signatures_with_archival_data.go`, 399 `validation_time_sliding.go`). 100% of the remaining 36 files read in full. Every `Process()`/verdict body, every `initChain`, and the POE comparator / POE extraction / evidence-record blocks were cross-checked line-by-line against Java DSS 6.5.RC1 at `/Users/utain/Workspace/esig/dss/dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/`.

### Findings

**SEC**

```
T39B-SEC-001 — **Low** — Category: SEC — dss/validation/process/vpfswatsp/poe_extraction.go:190-196
- Evidence: IsPOEExistInRange dereferences both bounds unguarded:
  `if poe.Time().Compare(*notBefore) >= 0 && poe.Time().Compare(*notAfter) <= 0`.
  Both callers (`past_signature_validation_certificate_revocation_selector.go:126`,
  `validation_process/validation_process_utils_vpfswatsp.go:110`) pass the certificate's
  `NotBefore()` / `NotAfter()` straight through, and a certificate is allowed to lack
  either (X.509 permits absence of both), so a nil bound panics on a nil-pointer
  dereference the moment the token has >=1 POE in the map.
- Impact: A revocation-data issuer certificate with no notBefore/notAfter field (a
  malformed/hostile certificate) panics instead of returning a verdict. This is a
  faithful port: Java's `POEExtraction.isPOEExistInRange` does the identical unguarded
  `compareTo` and raises a NullPointerException at the same spot (verified against
  Java POEExtraction.java). It is a panic on hostile input, not a false-PASS, and the
  path requires a malformed issuer certificate plus a POE present for it, so exposure
  is bounded. Same class as the U39 T39-SEC-001 (EAA notBefore/notAfter) finding,
  which was rated Low on identical Java-parity grounds.
- Recommendation: Guard `notBefore == nil || notAfter == nil` (return false) before the
  dereferences, or accept the Java-identical panic and record it as a sanctioned
  divergence. Note the sibling `BestSignatureTimeAfterCertificate...Check` in the same
  package already guards `notBefore != nil && ... && notAfter != nil` for the same
  reason — this range-check is the one place the guard was dropped.
```

**PERF**

None. POE extraction (`POEExtraction`) is linear per token/timestamp/evidence-record,
with POE lists bounded by the covered-object references (small, policy-bounded). The
recursive `revocationDataAlgorithmsObsolescenceValidationRecursive` in
`past_signature_validation.go` is guarded by the `checkedTokens` set (a port of Java's
`ArrayList#contains` cycle-breaker), so it cannot blow up on a hostile chain. No
O(n²)+ on attacker-controlled input, no leaks, no unbounded allocation in the
evidence-record path.

**STD**

None. The `//nolint:staticcheck` dead-store markers on the trailing
`item = item.SetNextItem(...)` calls are deliberate and documented in-file (they mirror
Java's own dead stores). `gofmt`, `go vet`, and `golangci-lint` are all clean.

### Verdict-path verification notes

**Status aggregation — verified CORRECT per VPF process.**

- **SwATSP main process (`validation_process_for_signatures_with_archival_data.go`)** —
  `initChain` is a faithful 1:1 port of Java's
  `ValidationProcessForSignaturesWithArchivalData#initChain`: the evidence-record
  loop, the `longTermValidation()` gate, the `IsAllowedValidationWithLongTermData` early
  return, the `NO_CERTIFICATE_CHAIN_FOUND_NO_POE` workaround branch, the
  newest-first timestamp sort (`sort.SliceStable` over `ProductionTime().After`, a
  faithful port of `Comparator.comparing(getProductionTime).reversed()`), the
  `tLevelTimeStamp`/`ltaLevelTimeStamp` optional checks, the step-6 lowest-POE
  recompute, and the `IsValid(LTV) → skip PSV` short-circuit all match Java line-for-line.
  The terminal `signatureIsAcceptable` dead-store `//nolint` is documented. **No
  swapped PASS/FAIL, no local re-aggregation bug.** The FAIL hard-stop is enforced by
  the shared framework (`ChainItemBase.Execute()`), which is verified in the U39
  section and not re-reported here.
- **PSV (`past_signature_validation.go`)** — the 7-step `initChain` (revocation-POE
  status, PCV-acceptable gate, POE-exists, the five `INDETERMINATE/*` branch dispatches
  on `currentConclusion`, and the terminal `pastRevocationDataValidationConclusive`)
  is a faithful port of Java's `PastSignatureValidation#initChain`. Every branch
  condition (`poeExists && INDETERMINATE == indication && <sub-indication>`) and every
  chained check (`bestSignatureTimeNotBeforeCertificateIssuance`,
  `bestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpiration`,
  `algorithmsObsolescenceValidation`, `revocationIsFresh`, `currentTimeIndicationCheck`)
  was compared against Java and matches. **No false VALID.**
- **VTS (`validation_time_sliding.go`)** — `initChain` (sunset-date control-time init,
  `reduceChainUntilFirstTrustAnchor`, reverse, the per-certificate revocation selection /
  revoked-branch / freshness-branch / crypto-notAfter-branch) is a faithful port of
  Java's `ValidationTimeSliding#initChain`. **No false VALID.**
- **PCV (`past_certificate_validation.go`)** — the trusted-certificate VTS loop,
  `bestValidationTimeSliding`, and the `controlTime != nil → cryptographicCheck` gate are
  faithful ports of Java's `PastCertificateValidation#initChain`. **No false VALID.**
- **Evidence-record (`evidencerecord/evidence_record_validation_process.go`)** —
  `initChain` (digest-matcher loop, `atLeastOneDataObjectFound`, `referenceDataGroup`,
  the signed-files-covered / signed-and-timestamped-files-covered gates, the
  `tstChainValid` timestamp walk, and the AOV selection) is a faithful port of Java's
  `EvidenceRecordValidationProcess#initChain`. **No false VALID.**

**Past-data / trust-anchor selection — verified, no false VALID (batch-07 T22D1A/B class not present).**

- The two revocation selectors (`past_signature_validation_certificate_revocation_selector.go`,
  `validation_time_sliding_certificate_revocation_selector.go`) are faithful ports of their
  Java twins. `VerifyRevocationData` in each correctly threads the `validity` Boolean
  through the `RevocationDataValidityMap` and accumulates acceptable revocations only on a
  positive POE-in-range / POE-exists test — the same conservative AND direction as the
  verified vpfltvd selector in the U39 section.
- The fail-open class flagged in batch-07 (T22D1A/T22D1B) was an *error swallow on
  `getOrderedCertificateChains`*. In this unit the analogous shared helpers
  (`GetLatestAcceptableRevocationData`, `GetAcceptableRevocationDataForPSVIfExistOrReturnAll`,
  `filterRevocationDataForPastSignatureValidation` in
  `validation_process/validation_process_utils_vpfswatsp.go`) were read and verified: they
  do NOT swallow errors — they are pure Boolean gates over `IsAllowedBasicRevocationDataValidation`,
  `IsRevocationDataAcceptable`, and POE-existence, all conservative (a non-allowed or
  POE-less revocation data is excluded, never accepted). **No error-swallowing, no
  stale/unorderable anchor set that produces a false VALID.**
- The trust-anchor comparison is identity-based (`certificate.Equals(trustedCertificate)`
  in VTS, `process.IsTrustAnchor` elsewhere), not a time-ordering that could mis-select.
  The `isCertificateSuspended` / `isTrustAnchor` / `revocationIsFresh` loops all guard
  `latestAcceptableRevocation != nil` before use. **No false VALID path found.**

**POE extraction / evidence records — allocation & path-safety verified.**

- **Unbounded allocation:** none. POE lists are bounded by covered-object references
  (small, policy-bounded). The only recursion
  (`revocationDataAlgorithmsObsolescenceValidationRecursive`) is cycle-bounded by the
  `checkedTokens` set (a faithful port of Java's `ArrayList#contains`). No host-file or
  host-memory unbounded growth on hostile evidence records.
- **Path traversal:** not applicable. The evidence-record "files covered" checks
  (`evidence_record_signed_files_covered_check.go`,
  `evidence_record_signed_and_timestamped_files_covered_check.go`) compare *document
  names* (`DocumentName`, `Filename`) purely in-memory against the evidence-record's own
  digest-matcher names — they never open, read, or resolve a filesystem path. There is no
  `os.Open`/`filepath.Join`/`ReadFile` in the unit, so no path-traversal surface.
  `signedFileCovered` correctly treats a `nil` signature `DocumentName` as "covered"
  (a faithful port of Java's `s.getDocumentName() == null ||`), and only matches
  non-nil names against non-nil ones (a faithful port of Java's `equals`).

**POEComparator / POE hierarchy — verified.** `POEComparator.Compare`
(time → type → timestamp-type → reference-count) and the `POEBase`/`TimestampPOE`/
`EvidenceRecordPOE` shadowing are faithful ports of Java's `POEComparator` and POE class
hierarchy. The `instanceof TimestampPOE` guard becomes a Go type assertion with the same
null-safe semantics. `GetLowestPOE`/`GetLowestPOETime` correctly return nil/panic-on-nil
to mirror Java's `getLowestPOE` (documented in-file). **No inversion, no swapped
preference.**

**Always-true/always-false and inverted comparisons — explicit statement:** None found
that are a porting divergence. Every always-true branch is a Java-faithful vacuous pass
or a documented `TODO : lax processing` / `TODO : this is a workaround` carried verbatim
from Java. Every inverted `!` (e.g. `!c.IsValid(...)` gates, `!containsToken(...)`) is a
faithful 1:1 port of the Java expression. No always-true check on a relevant input, no
inverted comparison, no swapped PASS/FAIL constant.

**Zero-finding lenses:** PERF (BIGO/MEM/LEAK), STD.

### Open questions

1. **T39B-SEC-001 (POE range-check nil-deref):** confirm whether a revocation-data issuer
   certificate lacking notBefore or notAfter is a realistic hostile input (X.509 permits
   the absence of both). If so, guard the deref (return false); if not, accept the
   Java-identical panic as a sanctioned divergence. Faithful to Java and the same class
   as the U39 T39-SEC-001 finding, so Low severity.
2. **`evidencerecord` subpackage has no test files** (`go test` reports `[no test files]`):
   the evidence-record process and blocks are exercised only through the cross-format
   oracle in `harness/` and the executor. The 100% line-by-line parity verification here
   is a static review, not a behavioural oracle; consider a targeted oracle test for the
   `tstChainValid` timestamp walk and the `signedFilesCovered` nil-name branch (the two
   spots with the most decision logic in the subpackage).

### Tool log

```
$ gofmt -l validation/process/vpfswatsp/
<no output>            (gofmt-exit=0, clean)
$ go vet ./validation/process/vpfswatsp/...
<no output>            (vet-exit=0, clean)
$ golangci-lint run --config=../.github/.golangci.yml ./validation/process/vpfswatsp/...
0 issues.              (lint-exit=0, clean)
$ go test ./validation/process/vpfswatsp/... -count=1
ok      github.com/ryftcore/dss-go/dss/validation/process/vpfswatsp     0.968s
?       github.com/ryftcore/dss-go/dss/validation/process/vpfswatsp/evidencerecord  [no test files]
(test-exit=0, all pass)
```

All four tool runs clean. Java parity cross-checks performed against
`/Users/utain/Workspace/esig/dss/dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/`
for: `ValidationProcessForSignaturesWithArchivalData`, `PastSignatureValidation`,
`ValidationTimeSliding`, `PastCertificateValidation`,
`PastSignatureValidationCertificateRevocationSelector`,
`ValidationTimeSlidingCertificateRevocationSelector`, `POEExtraction`, `POE`,
`POEComparator`, `TimestampPOE`, `EvidenceRecordPOE`, `EvidenceRecordValidationProcess`,
`EvidenceRecordsValidationBlock`, `EvidenceRecordTimestampsValidationBlock`,
`EvidenceRecordSignedFilesCoveredCheck`,
`EvidenceRecordSignedAndTimestampedFilesCoveredCheck`, `LongTermValidationCheck`,
`LongTermAvailabilityAndIntegrityValidationMaterialCheck`,
`AbstractPastTokenValidationCheck`, `POEExistsCheck`,
`POENotAfterCARevocationTimeCheck`, `POEExistsWithinCertificateValidityRangeCheck`,
`PastValidationAcceptableRevocationDataAvailable`, `PastCertificateValidationAcceptableCheck`,
`BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck`,
`RevocationIssuedBeforeControlTimeCheck`, `POEExistsAtOrBeforeControlTimeCheck`,
`SunsetDateCheck`, `SatisfyingRevocationDataExistsCheck`, `ControlTimeCheck`,
`SuccessfulValidationTimeSlidingFoundCheck`, `ValidationTimeSlidingCheck`,
`TimestampValidationCheck`, `CurrentTimeIndicationCheck`, `ProspectiveCertificateChainCheck`,
`PastRevocationDataValidationConclusiveCheck`, plus the shared
`ValidationProcessUtils` helpers (`getLatestAcceptableRevocationData`,
`getAcceptableRevocationDataForPSVIfExistOrReturnAll`,
`filterRevocationDataForPastSignatureValidation`). All in-scope verdict bodies and
chain logic confirmed faithful 1:1 ports; no divergence introduced.
