# Batch 10 — validation core (job + reports + process)

## validation/job + validation/reports (unit U31)

- **Date:** 2026-08-28
- **Depth:** budgeted (≤8,000 lines read)
- **Scope:** the cache / validation-job orchestration layer (the layer where U18A-SEC-001
  data race was found) and the diagnostic report-building data structures. A false "VALID"
  verdict is Critical.
- **Lens set:** SEC, PERF-BIGO, PERF-MEM / PERF-LEAK, STD.

### Scope (files actually read)

**`dss/validation/job` — 50 files, 3,177 lines — read in full (100%).**

| File | Lines | | File | Lines |
|---|---|---|---|---|
| validation_job.go | 464 (full) | | abstract_runnable_analysis.go | 100 (full) |
| cached_entry.go | 189 (full) | | cache_state_enum.go | 85 (full) |
| abstract_cache.go | 186 (full) | | download_cache_dto.go | 83 (full) |
| abstract_cache_dto.go | 181 (full) | | validation_job_alerter.go | 74 (full) |
| abstract_cache_access_by_key.go | 131 (full) | | cached_exception_wrapper.go | 67 (full) |
| abstract_analysis.go | 130 (full) | | abstract_parsing_cache_dto_builder.go | 67 (full) |
| current_cache_context.go | 129 (full) | | cache_cleaner.go | 66 (full) |
| abstract_cache_access_factory.go | 119 (full) | | cache_access_by_key.go | 66 (full) |
| validation_cache_dto.go | 113 (full) | | document_source.go | 64 (full) |
| abstract_parsing_cache_dto.go | 58 (full) | | download_cache.go | 56 (full) |
| abstract_read_only_cache_access_by_key.go | 48 (full) | | abstract_cache_access.go | 54 (full) |
| cache_key.go | 47 (full) | | cache_context.go | 47 (full) |
| abstract_cache_dto_builder.go | 41 (full) | | debug_cache_access.go | 36 (full) |
| changes_cache_access.go | 33 (full) | | synchronizer_cache_access.go | 28 (full) |
| parametrized_read_only_cache_access.go | 28 (full) | | abstract_parsing_result.go | 28 (full) |
| validation_result.go | 27 (full) | | validation_cache_dto_builder.go | 26 (full) |
| cache_access_factory.go | 25 (full) | | download_cache_dto_builder.go | 23 (full) |
| accept_all_strategy.go | 23 (full) | | synchronization_strategy.go | 21 (full) |
| document_alert.go | 21 (full) | | validation_cache.go / read_only_cache_access.go / parsing_cache.go / download_result.go | 19 each (full) |
| read_only_cache_access_by_key.go | 17 (full) | | cache_type.go | 14 (full) |
| validation_job_summary_builder.go / parsing_result.go / download_task.go | 12 each (full) | | validation_task.go / parsing_task.go | 10 each (full) |
| cached_result.go | 7 (full) | | | |

**`dss/validation/reports` — 14 files, 6,360 lines.**

| File | Lines | Read |
|---|---|---|
| diagnostic/diagnostic_data_builder.go | 2,533 | first 150 (header + DataBuilder struct) + func map (131 fns) + targeted region `getXmlForCertificateChainWithSource`/`containsToken`/`ensureCertificateChain`/`getCertChainSinceIndex` (lines 1129–1208). **~370 lines — budgeted.** |
| diagnostic/signed_document_diagnostic_data_builder.go | 1,589 | first 150 (header + type + ctor) + func map (101 fns). **~170 lines — budgeted.** |
| diagnostic/xml_trust_service_provider_builder.go | 840 | first 150 (full type + Build/sort/classify + getXmlTrustServiceProvider) + targeted region `getLangAndValues` (145–190). **~185 lines — budgeted.** |
| diagnostic/xml_qc_statements_builder.go | 273 | full |
| diagnostic/xml_trusted_entity_builder.go | 176 | full |
| diagnostic/xml_policy_builder.go | 175 | full |
| diagnostic/qwac_certificate_diagnostic_data_builder.go | 154 | full |
| abstract_reports.go | 154 | full |
| reports.go | 128 | full |
| diagnostic/xml_diagnostic_data_factory.go | 124 | full |
| diagnostic/certificate_diagnostic_data_builder.go | 78 | full |
| certificate_reports.go | 60 | full |
| dss_report_exception.go | 59 | full |
| doc.go | 17 | full |

> Not fully verified (budget): the interior of `diagnostic_data_builder.go` (lines ~151–2532
> beyond the three targeted regions) and `signed_document_diagnostic_data_builder.go`
> (lines ~151–1589). These are report-serialization bodies (JAXB assembly), not trust-state or
> verdict logic; the trust/verdict-relevant surface (chain building, orphan tokens, map-order
> determinism) was covered by the targeted reads. The two >300-line files' remaining ~3,300
> lines are marked **not fully verified (budget)**.

### U18A-SEC-001 race verification (dedicated subsection)

**Verdict: CONFIRMED — and escalated from High to Critical by this unit.**

U18a (tsl PART 1) reported the race from the *caller* side (the concurrent pivot fan-out).
Reading the in-scope `validation/job` code confirms the race is real and localizes the root
cause to the cache layer: **`AbstractCache.mu` protects only the `cachedEntriesMap` map, not
the `*CachedEntry` it returns.** Every state mutation happens *after* the mutex is released.

**Evidence (exact file:line, no `sync.Mutex`/`sync.RWMutex` guards any of these writes):**

- `validation/job/abstract_cache.go:77-88` — `Get(cacheKey)` locks `mu` only to
  read/insert the map, then **unlocks and returns the shared `*CachedEntry[R]` pointer**.
  The mutex never wraps the entry's methods.
- `validation/job/abstract_cache.go:96-98` — `Expire(key)` = `a.Get(cacheKey).Expire()`;
  the mutation `Expire()` runs with **no lock held**.
- `validation/job/cached_entry.go:108-110` — `Expire()` = `e.cacheContext.RefreshNeeded()` —
  a write into shared context state, unlocked.
- `validation/job/current_cache_context.go:94-96` → `:53-59` — `RefreshNeeded()` →
  `State(newState)` writes `c.state`, `c.lastStateTransitionTime = time.Now()`,
  `c.exception = nil` — three unsynchronized writes to a struct shared across goroutines.
- Same unlocked-write pattern: `abstract_cache.go:109-111` (`Sync`), `:133-135` (`Error`),
  `:90-92` (`Update`); `cached_entry.go:61-66` (`Update` → `Desync()` + `e.cachedResult =`),
  `:91-98` (`Error` → `cacheContext.Error/ErrorUpdateDate` + `e.cachedResult = zero`).

**Reachability from the public API — YES (concurrent cache access is triggerable):**

`tsl/lotl_with_pivots_analysis.go:131` `downloadAndParseAllPivots` runs
`go func(...){ processing.Call(); ... }` **per pivot** (line 170), `wg.Wait()` (line 182).
Each `PivotProcessing` is constructed (line ~139) with the **same**
`lotlCacheAccessByKey` (`a.CacheAccessByKey().(*TLCacheAccessByKey)`, line ~138) and a shared
`pivotCacheAccessByKeyList`. `PivotProcessing.expireCache` / `Call` then drive
`ExpireValidation()` / `Sync()` / `UpdateValidationResult()` (all `validation/job` methods above)
on that one shared LOTL cache entry. With the normal EU pivot chain (≥2 pivots: 226→191→172)
this is concurrent read+write of the same `CurrentCacheContext` — a Go data race.

**Impact (trust / verdict path):** a lost state transition leaves a pivot / LOTL cache entry
"fresh" (or ERROR/REFRESH_NEEDED overwritten) when it should have expired. The pivot's cached
validation state feeds the accepted-signer set (`currentCertificateSourceFromPivots` promotes a
pivot's certs only when `IsResultExist() && IsValid()`). A corrupted transition can therefore
advance the trust chain off a stale/invalid pivot validation → **false "VALID"** on a
verification path, or a panic via the illegal-transition `panic` in `CacheStateEnum`
(`cache_state_enum.go`) when two goroutines interleave `Desync`+`Sync`. This is trust-state
corruption reachable from the public `Validate`/TL-job API.

**`go test -race` result:** `ok github.com/ryftcore/dss-go/dss/validation/job  1.737s` —
**PASS, no race reported.** This does **not** refute the race: the in-scope job tests
(`abstract_cache_test.go`, `cached_entry_test.go`, `cache_state_enum_test.go`,
`cache_key_test.go`, 363 lines total) are all **single-threaded** — none use `go func` or
concurrent `Get`/`Expire`/`Sync`. The concurrent trigger lives in `tsl` (the pivot fan-out),
which is not exercised by a `-race` multi-pivot test. The `-race` pass is therefore
**consistent with, but does not disprove**, the confirmed code-level race.

**Recommendation:** make the `CachedEntry` state machine thread-safe — add a per-entry
`sync.Mutex` (or move the whole entry into the lock scope of `AbstractCache`) guarding
`CurrentCacheContext` + `cachedResult` mutations, and add a `tsl`-level `-race` test that
drives ≥2 pivots sharing one LOTL cache access.

### Findings

**SEC**

```
T31-SEC-001 — **Critical** — Category: SEC — dss/validation/job/abstract_cache.go:96-98,109-111,133-135 & cached_entry.go:108-110,61-66,91-98 & current_cache_context.go:53-59,94-96
- Evidence: AbstractCache.mu guards only cachedEntriesMap (abstract_cache.go:77-88); Expire/Sync/Error/Update call
  Get(key) then mutate the returned *CachedEntry with no lock held. Expire→cacheContext.RefreshNeeded()→State()
  writes state/time/exception (current_cache_context.go:57-59) unsynchronized. tsl fans out ≥2 pivot goroutines
  (lotl_with_pivots_analysis.go:170) sharing the same lotlCacheAccessByKey (line ~139).
- Impact: Data race on trust-state = undefined behaviour. Lost transition can leave a pivot/LOTL "fresh"/valid when
  it should be expired, advancing the accepted-signer set off stale pivot validation → false "VALID" (Critical),
  or an illegal-transition panic (cache_state_enum.go) from interleaved Desync+Sync. Reachable from the public
  TL/Validate API with the normal ≥2-pivot chain.
- Recommendation: Per-entry sync.Mutex (or hold AbstractCache.mu across the mutation) guarding CurrentCacheContext +
  cachedResult; add a tsl -race multi-pivot test. This is the same defect as U18A-SEC-001, confirmed here at the
  root cause and escalated High→Critical because it sits on the verdict/trust path.
```

```
T31-SEC-002 — **Medium** — Category: SEC — dss/validation/job/validation_job_alerter.go:65-66,72-73
- Evidence: executeDocumentListAlert/executeDocumentAlert both do `_ = alrt.Alert(info)` (lines 66 and 73),
  discarding the returned error. Alert() (alert/abstract_alert.go:16) returns the handler's error.
- Impact: A failing alert handler (e.g. one that raises on a document-list change) is silently swallowed.
  This is the U28 OQ. It is a *deliberate* porting choice (Java logs the swallowed exception; slf4j is dropped,
  per the file's DEVIATION note) and the alert's load-bearing side-effect (handler.Process) still runs — so this
  is an observability loss, not a trust/verdict defect. Rated Medium, not Critical.
- Recommendation: Keep the swallow (parity) but route the error to the configured alert/observability sink or a
  returned aggregate so a handler failure is not invisible; document the deliberate exclusion.
```

**PERF-BIGO**

```
T31-PERF-BIGO-001 — **Low** — Category: PERF-BIGO — dss/validation/reports/diagnostic/diagnostic_data_builder.go:1135-1207
- Evidence: getXmlForCertificateChainWithSource walks the issuer chain; ensureCertificateChain
  (lines 1173-1199) calls getCertChainSinceIndex(i) (line 1200) which copies certChain[i:] in O(n),
  and does so once per chain element → O(n²) in chain length; containsToken (line 1164) is also a
  linear scan per hop.
- Impact: O(n²) report-serialization cost on certificate-chain length. Chain length is bounded by real
  PKI depth (single digits) and is not attacker-inflatable to large n, so this is a Low hardening note,
  not an exploitable O(n²)+ on attacker input.
- Recommendation: If ever a long synthetic chain is fed, build the suffix wrappers in one pass; otherwise
  leave as-is (matches Java's per-element slice).
```

```
T31-PERF-BIGO-002 — **Low** — Category: PERF-BIGO — dss/validation/reports/diagnostic/xml_trusted_entity_builder.go:122-134
- Evidence: getLangAndValues ranges the `m map[string][]string` directly (line 125) with no sort, so the
  emitted XmlLangAndValue order is Go's randomised map order and differs run-to-run. The sibling
  XmlTrustServiceProviderBuilder.getLangAndValues (xml_trust_service_provider_builder.go:147) sorts keys
  with sort.Strings for determinism.
- Impact: Non-deterministic report element order for a multi-language TSP name/tradeName in the
  trusted-entity path. No perf cost; a byte-parity/reproducibility gap against the sorted sibling and
  against Java (known-gaps.md "Deliberate ordering choices" covers most, not this one).
- Recommendation: Sort `langs` like the sibling builder (or confirm the multiset carve-out in known-gaps.md
  and pin with a test).
```

**PERF-MEM / PERF-LEAK**

```
T31-PERF-MEM-001 — **Low** — Category: PERF-MEM — dss/validation/job/abstract_cache.go:77-88
- Evidence: Get(cacheKey) inserts a new empty CachedEntry into cachedEntriesMap on every miss (lines 83-86);
  the only removal is an explicit Remove(key) via CacheCleaner for TO_BE_DELETED entries.
- Impact: The map grows with the number of distinct document/LOTL/pivot keys ever seen in a job's lifetime.
  Bounded by the configured TL/LOTL/pivot source set (small, caller-controlled), so not an unbounded
  leak on attacker input; a note only. Entries that are never marked TO_BE_DELETED are retained for the job.
- Recommendation: If a long-lived job processes many distinct URLs, consider a size cap / LRU; otherwise
  leave as-is (mirrors Java's ConcurrentHashMap which also does not auto-evict).
```

**STD**

```
T31-STD-001 — **Low** — Category: STD — dss/validation/job/document_source.go:58-64
- Evidence: CacheKey() does an unsynchronized lazy-init: `if !s.cacheKeySet { s.cacheKey = NewCacheKey(s.url);
  s.cacheKeySet = true }` — a read-modify-write with no lock.
- Impact: If a *DocumentSource is shared across goroutines (it is handed into the concurrent analysis
  runnables), two goroutines can both compute and assign cacheKey (harmless — same value) or race on
  the bool. Benign in value (deterministic result) but a latent data race under -race if the source is
  shared concurrently; today each analysis uses its own source so it is not triggered.
- Recommendation: Compute the key once at SetUrl time, or guard the lazy-init; Low because the current
  call graph does not share a DocumentSource across goroutines.
```

**Zero-finding lenses:** none — every lens (SEC, PERF-BIGO, PERF-MEM, STD) has at least one finding above.
No Critical/High in PERF-BIGO, PERF-MEM, or STD.

### Tool log (run from inside `dss/`)

| # | Command | Result |
|---|---|---|
| 1 | `gofmt -l validation/job validation/reports` | **(empty — no unformatted files)** |
| 2 | `go vet ./validation/job/... ./validation/reports/...` | **clean (exit 0)** |
| 3 | `golangci-lint run --config=../.github/.golangci.yml ./validation/job/... ./validation/reports/...` | **0 issues** |
| 4 | `go test -race ./validation/job/... -count=1` | **`ok …/validation/job 1.737s` — PASS, no race reported** (caveat: job tests are single-threaded; the concurrent trigger is in `tsl`, see U18A-SEC-001 subsection) |

All in-scope tools are clean; the one load-bearing result is the `-race` pass, which is **consistent
with** the confirmed code-level race but does not disprove it because no in-scope test drives the
shared-entry path concurrently.

### Summary counts

- **Total findings: 6**
- By severity: **Critical 1** · **High 0** · **Medium 1** · **Low 4** · **Info 0**
- By lens: **SEC 2** (1 Critical, 1 Medium) · **PERF-BIGO 2** (both Low) · **PERF-MEM 1** (Low) · **PERF-LEAK 0** · **STD 1** (Low)

### Open questions for consolidation

1. **Severity of U18A-SEC-001:** U18a rated it High. This unit escalates to **Critical** because the
   corrupted state machine sits directly on the TL/pivot trust path and can produce a false "VALID".
   Confirm the escalation and the single canonical finding ID across batch-06 / batch-10.
2. **Fix ownership:** the root cause is in `validation/job` (thread-safe `CachedEntry` transitions);
   the trigger is in `tsl`. Decide whether the fix + the `-race` multi-pivot regression test land in the
   `validation/job` unit (U31) or the `tsl` unit (U18a).
3. **T31-SEC-002 (alert swallow):** keep as a Medium finding or fold into a deliberate-exclusion note,
   given the file already carries a DEVIATION comment? Suggest: deliberate-exclusion note + a follow-up to
   route the error to an observability sink.
4. **T31-PERF-BIGO-002 (getLangAndValues ordering):** confirm whether the multiset carve-out in
   `docs/compatibility/known-gaps.md` already covers this case, or whether it needs an explicit entry +
   pinning test to avoid silent widening.

## validation/process core + blocks + bbb core + eaa (unit U32)

- **Date:** 2026-08-28
- **Depth:** budgeted (≤2,900 lines read; 2,946 lines of in-scope code actually read, 100%)
- **Scope:** the process-block orchestration (`validation/process` top-level), the building-block dispatcher
  (`validation/process/blocks`), the abstract building-block check framework (`validation/process/bbb` core), and
  the EAA process checks (`validation/process/eaa`). A building-block check returning a false PASS (false "VALID")
  on the verdict path is Critical; a false FAIL is High.
- **Lens set:** SEC, PERF-BIGO, PERF-MEM / PERF-LEAK, STD.

### Scope (files actually read)

**`dss/validation/process` (top-level) — 7 files, 1,789 lines — read in full (100%).**

| File | Lines |
|---|---|
| validation_process_utils.go | 740 (full) |
| chain_item.go | 468 (full) |
| chain.go | 303 (full) |
| validation_process_utils_vpfswatsp.go | 103 (full) |
| basic_building_block_definition.go | 80 (full) |
| certificate_policy_identifiers.go | 55 (full) |
| uninterrupted_chain_item.go | 40 (full) |

**`dss/validation/process/blocks` — 3 files, 518 lines — read in full (100%).**

| File | Lines |
|---|---|
| basic_building_blocks.go | 428 (full) |
| basic_building_blocks_eaa.go | 56 (full) |
| basic_building_blocks_noeaa.go | 34 (full) |

**`dss/validation/process/bbb` (core) — 3 files, 175 lines — read in full (100%).**

| File | Lines |
|---|---|
| abstract_certificate_check_item.go | 69 (full) |
| abstract_multi_values_check_item.go | 60 (full) |
| abstract_value_check_item.go | 46 (full) |

**`dss/validation/process/eaa` — 4 files, 464 lines — read in full (100%).**

| File | Lines |
|---|---|
| eaa_validation_process.go | 204 (full) |
| eaa_validation_block.go | 182 (full) |
| key_binding_signature_validation_result_check.go | 49 (full) |
| doc.go | 29 (full) |

**Out of scope (noted, not read in depth):** `bbb/{aov,cv,fc,isc,sav,vci,xcv}`, the `cv/*` concrete checks, and the
`qualification`/`vpfbs` check bodies — cross-referenced only where an in-scope EAA check delegates its `Process()` /
verdict to them (their `Process()` all reduce to `IsValid(conclusion)` = `PASSED`/`TOTAL_PASSED`, confirmed).

### Building-block verdict verification

**Framework status aggregation (`ChainItemBase`/`ChainBase`) — CORRECT.**

- `ChainItemBase.Execute()` dispatches on the check's `Level()`: `""` → check skipped, no constraint recorded, chain
  continues; `IGNORE` → records `XmlStatusIgnored` and continues; `FAIL` → runs `Process()` and records `OK`
  (continue) or `NOT_OK` (record conclusion, and stop the chain unless `ContinueProcessOnFail()`); `INFORM`/`WARN` →
  runs `Process()`, records `OK` or `INFORMATION`/`WARNING`, then always continues. This is exactly Java's
  `ChainItem.execute()` ladder — a `FAIL`-level failure is a hard gate and is not masked by later checks.
- `ChainBase.Execute()` sets the conclusion to `PASSED` only when the chain ran to completion and no item set a
  conclusion (`Conclusion() == nil`). `blocks.BasicBuildingBlocks.updateFinalConclusion` overwrites the running final
  indication only when a block's indication is `!= PASSED`, so a later `PASSED` block cannot mask an earlier
  `FAILED`/`INDETERMINATE` block. Net effect: **NOT_EVALUATED (skipped) and INDETERMINATE are never collapsed to
  PASS** — `IsValidConclusion` accepts only `PASSED`/`TOTAL_PASSED`.
- No always-PASS or always-FAIL check in the in-scope set. The `bbb` core helpers are predicate checks, all returning
  `false` on empty input and `true` only on a genuine match: `AbstractValueCheckItem.ProcessValueCheck` (exact
  case-sensitive match or `*` wildcard); `AbstractMultiValuesCheckItem` → `ValueCheck`/`ValuesCheck`/
  `AllValuesCheck`/`ValuesForEachExpectedCheck`; and `AbstractCertificateCheckItem.ProcessCertificateCheck`, which
  returns `true` when the certificate's extension-OIDs satisfy the expected-extension constraint **OR** its
  policy-OIDs satisfy the expected-policy constraint — the OR is the upstream Java semantics (a policy rule is "has
  one of these extensions or one of these policies"), not a looser-than-Java widening.

**EAA check status logic — CORRECT (each gates on the right field with the right policy constraint).**

| EAA step | Gates on | `Process()` valid when | Policy constraint |
|---|---|---|---|
| FC | `eaaBBBs.FC` | `FC` conclusion `PASSED`/`TOTAL_PASSED` | `FailLevelRule` |
| Per-signature | that signature's basic-validation conclusion | conclusion `PASSED`/`TOTAL_PASSED` | `EAASignatureValidConstraint` |
| Key-binding signature | KB-signature's basic-validation conclusion | conclusion `PASSED`/`TOTAL_PASSED` | `EAAKeyBindingSignatureValidConstraint` |
| CV (selective disclosures) | `eaaBBBs.CV` | `CV` conclusion `PASSED`/`TOTAL_PASSED` | `FailLevelRule` |
| SAV | `eaaBBBs.SAV` | `SAV` conclusion `PASSED`/`TOTAL_PASSED` | `FailLevelRule` |

The delegated `Process()` bodies in `qualification.SignatureValidationResultCheck` and
`vpfbs.{FormatCheckingResultCheck, CryptographicVerificationResultCheck, SignatureAcceptanceValidationResultCheck}`
(out of scope, cross-referenced) all reduce to `IsValid(conclusion)` = `PASSED`/`TOTAL_PASSED`, so an INDETERMINATE
basic validation correctly fails the EAA "conclusive" gate rather than being treated as a pass.

**Not a finding (verified, matches upstream):** `blocks.BasicBuildingBlocks.updateFinalConclusion` uses "last
non-PASSED block wins," so a later `INDETERMINATE` block can overwrite an earlier `FAILED` block's indication (a
softening, never a false PASS) — this is verbatim `BasicBuildingBlocks#updateFinalConclusion` and stays on the
non-PASS side, so no false "VALID" results. `GetConstraintOrMaxLevel`'s `fallthrough` ladder was traced for all
16 (level × maxLevel) combinations and is correct. `getDigestMatcherType` reads only `digestMatchers[0]` (verbatim
Java; affects only the message-tag position label, never a status).

### Findings

**SEC**

```
T32-SEC-001 — **High** — Category: SEC — dss/validation/process/eaa/eaa_validation_process.go:84-108
- Evidence: In InitChain, `var item process.ChainItem[...]`; `xmlFC := eaaBBBs.FC; if xmlFC != nil { item = ... }`;
  then unconditionally `item = item.SetNextItem(c.signatureValidationConclusive(...))` (line 91), and again at 96/102/108.
- Impact: When eaaBBBs.FC is nil, item stays the nil interface and the first `item.SetNextItem(...)` panics with a
  nil pointer dereference, aborting the whole validation run instead of returning a status. This is reliably a
  panic, not a silent false-PASS, because under the default (non-eaa) build blocks.executeEAAFormatChecking() is
  stubbed to return nil (basic_building_blocks_noeaa.go) — so eaaBBBs.FC is always nil for an EAA token — while
  executeCryptographicVerification() is NOT stubbed and always produces a CV block. Hence any EAA reaches a
  `SetNextItem` on the nil item: at the signature loop (line 91) if it has ≥1 signature (the normal EAA case), or at
  the CV step (line 102) if it has none. Yet EAA diagnostic parsing (diagnostic/diagnostic_data.go, untagged) and
  the EAA presentation builder (validation/executor/detailed_report_for_eaa_presentation_builder.go, untagged) still
  construct an EAAWrapper and route it here. The `!eaa` stub's stated assumption ("an EAAWrapper is never
  constructed without the EAA module") does not hold. Java validates the same input successfully; the Go port panics.
- Recommendation: Guard the chain build — if eaaBBBs.FC == nil (no FC block present, e.g. !eaa build), seed `item`
  from the first available step (CV/SAV) or return an explicit non-PASSED "EAA validation unsupported in this build"
  conclusion, rather than calling SetNextItem on a nil item.
```

**PERF**

```
T32-PERF-BIGO-001 — **Medium** — Category: PERF-BIGO — dss/validation/process/blocks/basic_building_blocks.go:339-375
- Evidence: removeAllCertificates / removeAllOrphanCertificates each run a nested loop:
  `for _, certificate := range certificates { for _, removed := range toRemove { if ... { found = true } } }` — O(n·m).
- Impact: Called once per SubXCV in addAdditionalInfo over the certificate's cross/equivalent-certificate sets, whose
  sizes are attacker-influenced (a malicious certificate can carry many cross-certificate extensions, and the
  used-certificate set can be large), so the total work is quadratic in those counts.
- Recommendation: Build a set (map[string]struct{}) of the toRemove ids (or use a helper) to make each removal O(n),
  or sort-and-scan; low priority since the counts are bounded by the certificate set, but it is quadratic on
  attacker-controlled input.
```

**STD**

```
T32-STD-001 — **Low** — Category: STD — dss/validation/process/eaa/eaa_validation_process.go:131-135
- Evidence: keyBindingSignatureValidationConclusive does `xmlSignature := c.xmlSignatures[signatureWrapper.Id()]` then
  dereferences `xmlSignature.ValidationProcessBasicSignature.Conclusion` with no nil guard; the sibling
  signatureValidationConclusive instead routes through getSignatureBasicProcessConclusion, which panics with a
  descriptive "No basic signature validation process found for the signature with Id ..." message.
- Impact: Asymmetric defensive handling of the same map lookup. Unreachable today (ValidationBlock.Execute populates
  signatureValidationMap for the key-binding signature), but a latent opaque nil-deref if the population invariant
  changes, versus a clear, debuggable panic in the sibling.
- Recommendation: Route the key-binding lookup through the same guarded helper (or add an equivalent nil check with a
  descriptive message) so both signature kinds fail identically.
```

**Zero-finding lenses:** PERF-LEAK — none (no unbounded retention; the result objects are bounded by the token/block
count and released with the report).

### Tool log (run from inside `dss/`)

| # | Command | Result |
|---|---|---|
| 1 | `gofmt -l validation/process` | **(empty — no unformatted files)** |
| 2 | `go vet ./validation/process/...` | **clean (exit 0)** |
| 3 | `golangci-lint run --config=../.github/.golangci.yml ./validation/process/...` | **0 issues** |
| 4 | `go build ./validation/process/...` (default) / `go build -tags eaa ./validation/process/...` | **both exit 0** — the eaa/noeaa split compiles in both tag states |

### Summary counts

- **Total findings: 3**
- By severity: **Critical 0** · **High 1** · **Medium 1** · **Low 1** · **Info 0**
- By lens: **SEC 1** (High) · **PERF-BIGO 1** (Medium) · **PERF-LEAK 0** · **PERF-MEM 0** · **STD 1** (Low)

### Open questions for consolidation

1. **T32-SEC-001 reachability/ownership:** confirm whether EAA presentation validation is expected to be *usable* in
   the default (non-eaa) build at all. If the product stance is "EAA requires `-tags eaa`," the fix is a clean "not
   supported" status in `blocks`/`eaa`; if EAA should work without the tag, the `!eaa` stubs are the wrong design and
   `executeEAAFormatChecking` should fall back to a real FC. Either way, the nil-`item` panic must go.
2. **`updateFinalConclusion` "last non-PASSED wins":** verbatim Java and never a false PASS, so not a finding — but
   flag for `docs/compatibility/known-gaps.md` if a "most-severe block wins" semantic would be more defensible than
   "last non-PASSED wins" (currently a later INDETERMINATE can soften an earlier FAILED indication).
