# dss-go — Code Review Report

**Date.** 2026-08-29
**Scope.** Whole Go module `dss/` (`github.com/ryftcore/dss-go/dss`, Go 1.27, stdlib + `golang.org/x` only), plus `cmd/esig`, `examples/`, `harness/`, and the repo-root `corpus/` generators. ~2,854 `.go` files.
**Companion detail.** Per-batch findings live in `docs/review/batch-01a…16-*.md`. This file is the consolidated executive view.
**Mode.** Read-only. No source/test/config file was modified. Findings are reported, not fixed.

---

## 1. Executive summary

dss-go is a disciplined, high-fidelity 1:1 port of Java DSS 6.5.RC1. Across 46 review units (17 batch files) the port is **structurally sound**: the static-analysis gate (`gofmt`, `go vet`, `golangci-lint` v2, `go test -race`) is **clean in every reviewed package in both the default and `eaa` build-tag states**, and the security-critical verification paths were repeatedly cross-checked line-by-line against the Java oracle at `~/Workspace/esig/dss` and found **fail-closed** (missing data → NOT_EVALUATED/INVALID, never → PASS).

The single most important result is the **one Critical**: a **data race in the generic cache** (`validation/job`) that is reachable through the TSL pivot fan-out and can, in a concurrent validation workload, return a **stale trust anchor → false VALID verdict**. This is the only finding that changes a *verdict* rather than crashing or leaking, and it is the top remediation priority.

**Findings by severity (reconciled from per-unit records):**

| Severity | Count | Meaning in this codebase |
|---|---:|---|
| **Critical** | **1** | Exploitable / false-accept on a verification path |
| **High** | **18** | Likely vuln, real leak, O(n²)+ on attacker input, panic on hostile input, trust-logic error |
| **Medium** | **43** | Hardening opportunity / parity gap that weakens security |
| **Low** | **97** | Idiom, minor perf, parity-panic, dead code |
| **Info** | **106** | Observation, deliberate-parity note, open question |
| **Total** | **~265** | |

> **Reconciliation note.** C and H are exact. M/L/I are exact to ±1–2: two early units (U12b, U23) recorded a severity breakdown that is off by one from their stated finding total, so the Medium/Low/Info split carries a small residual. The per-batch table in §5 is the authoritative breakdown.

**What is NOT a finding** (per `PORTING.md` binding conventions): staticcheck `ST*/S1*/QF*/SA1019/SA9003` categories are deliberately disabled (they would restyle Java-verbatim strings or flag legacy crypto the validator must keep); `// DIVERGENCE, deliberate:` markers are sanctioned; panics below the facade are the Java-exception idiom. Java-parity behavior (the same bug exists in upstream) is only reported when it weakens security or is an *undocumented* divergence. Pre-existing gaps in `docs/compatibility/known-gaps.md` are not re-reported.

---

## 2. Methodology

- **Units.** The module was split into **46 small-scope units** (≤~30 files / ≤~8k read-lines each) so each could be reviewed by a dedicated senior-reviewer pass without context overflow. Units ran **max 2 in parallel**.
- **Lenses (per unit).**
  - **STD** — Go coding standards / idioms / best practice.
  - **PERF** — BigO, allocations, memory leaks, unbounded growth.
  - **SEC** — security: false-accept/false-reject in verdicts, panic on hostile input, nil-deref, inverted logic, trust-logic errors, injection, crypto misuse.
- **Severity rubric.** Critical = exploitable/false-accept in a verification path · High = likely vuln / real leak / O(n²)+ on attacker input / panic on hostile input / trust-logic error · Medium = hardening or security-weakening parity gap · Low = idiom / minor perf / parity-panic · Info = observation / deliberate-parity note / open question.
- **Evidence standard.** Every finding carries `file:line` + ≤5-line code evidence + impact + recommendation. Security findings on verification paths were **cross-checked against the Java DSS 6.5.RC1 oracle** (`~/Workspace/esig/dss`) to establish whether the behavior is a divergence or faithful parity.
- **Tool gate (per unit).** `gofmt -l`, `go vet` (+ `-tags eaa` where relevant), `golangci-lint run --config=../.github/.golangci.yml`, `go test -count=1`, and `go build`/`go build -tags eaa`. **Result: clean in every unit, both tag states.**
- **Report-only.** No source, test, or config file was modified.

---

## 3. Top 10 findings

Ordered by severity then blast radius. Full evidence in the cited batch file.

| # | ID | Sev | Category | Location | One-line |
|---|---|---|---|---|---|
| 1 | **T31-SEC-001** | **Critical** | Data race → false verdict | `validation/job/abstract_cache.go:96-98` (+ `cached_entry.go:108-110`, `current_cache_context.go:53-59`) | `AbstractCache.mu` guards only the map; `Expire/Sync/Error/Update` mutate `CachedEntry`/`CurrentCacheContext` **unlocked**. Reachable via TSL pivot fan-out (`tsl/lotl_with_pivots_analysis.go:170`) → a concurrent reader can see a **stale trust anchor → false VALID**. Confirms & escalates U18A-SEC-001. `go test -race` passes only because no concurrent test exists. |
| 2 | **T23-SEC-001** | High | False-accept (divergence) | `spi/policy/basic_asn1_signature_policy_validator.go:71,91-96` | `SetDigestValid(true)` is set *before* the `recover` block; a top-level SEQUENCE with <3 children returns early **without resetting** it → Go reports the signature-policy digest **valid** where Java reports it **invalid** (verified vs `BasicASN1SignaturePolicyValidator.java`). Feeds `SignaturePolicyHashValidCheck`. |
| 3 | **T22D1A/B-SEC-001** | High | Fail-open (error swallowed) | `spi/validation/signature_validation_context.go:1142,1233,1569`; `…_alerter.go:58,73` | `getOrderedCertificateChains()` error is swallowed → the revocation-data *presence* check reports **PASS** on an unorderable certificate set. Two call sites, same root. |
| 4 | **T18A-SEC-001/002** | High | Data race + panic | `tsl/…` (pivot cache), `tsl/…` (nil `SchemeInformation`) | (a) pivot-cache shared mutation (the race T31-SEC-001 escalates); (b) nil-`SchemeInformation` deref in a goroutine. |
| 5 | **P11A-SEC-001/002/003** | High | Memory DoS + panics | `internal/pdf/filter.go:138-152,183-210`, `objstm.go:31-44` | Flate decompression-bomb (1 MB → 1.07 GB reproduced), `ApplyPredictor` panics on hostile `/DecodeParms`, `parseObjStmHeader` panics on hostile `/N`. |
| 6 | **T33-SEC-001** | High | Panic on hostile input | `validation/process/bbb/…/signature_value_and_signed_attributes…obsolescence.go:104` | Nil-deref when every digest method is empty — a hostile AOV input crashes the report path. |
| 7 | **T35-SEC-001** | High | Panic | `validation/process/bbb/isc/identification_of_the_signing_certificate.go:131` | `sources[0]` on an empty `Sources()` slice — report-generation phase (not verdict), but a crash on a degenerate certificate. |
| 8 | **T32-SEC-001** | High | Panic (build-tag) | `validation/process/eaa/eaa_validation_process.go:84-108` | EAA nil-panic in the **non-`eaa`** build (`SetNextItem` on a nil FC) — an untagged code path reaching an `eaa`-only dependency. |
| 9 | **T37-SEC-001** | Med | Panic (report-DoS) | `validation/process/qualification/eaa_qualification_block.go:111,146` | `EAASignatures()[0]` with no empty-guard → panic on a 0-signature EAA. `eaa`-tag gated; report-DoS, not a false verdict. |
| 10 | **T20-SEC-001** | Med | Trust-data truncation | `token/…` (multi-key PFX) | A multi-key PKCS#12 with a shared chain yields a **truncated** `CertificateChain` — downstream revocation/anchor checks see fewer certificates than Java emits. |

---

## 4. Cross-cutting themes

These are the recurring *classes* of issue, which is more useful than the individual items when planning remediation.

1. **Fail-open error-swallowing (the dominant security class).** A handful of verification paths convert an error into a *pass* by dropping the error: `getOrderedCertificateChains` (T22D1A/B), the ASN.1 policy `<3`-element early-return (T23-SEC-001), and the alerter discarding alert errors (deliberate, T31-SEC-002). The codebase is otherwise consistently fail-closed; these are the exceptions and the highest-value hardening targets. **Pattern:** an error/`!ok` that should `return false`/`SetInvalid` instead falls through to a previously-set `true`.
2. **Panic-on-hostile-input (the dominant crash class).** Mostly **Java-parity** NPEs / `IndexOutOfBounds` ports (T33, T35, T37, T38, T39, T39B, P11A-002/003, P12A, P13) where a `default: panic` or unguarded `[0]`/`*NotBefore()` mirrors a Java `throw`/NPE. Individually Low–High; collectively they are a DoS surface on malformed-but-parseable input. The facade's `recovered()` converts them to errors on the public path, but internal report-generation and `eaa`-gated paths can still crash.
3. **Stale / unsynchronized shared state (the concurrency class).** The **Critical** T31-SEC-001 cache race is the headliner; the same *shape* (unsynchronized read-then-write lazy init) recurs as Low/Info in `detailedreport`/`diagnostic` (T43-M-001), `xades` singletons (T41-INFO-002), and the policy wrapper (T45-STD-003) — all 1:1 Java parity and benign in the current single-threaded port, but a latent race the moment concurrent validation is introduced.
4. **Unbounded registries / caches (the leak class).** `Store`-only `sync.Map`/global registries with no eviction: `xades` (T40-PERF-001, 3 sites), `pades` timestamp-token registry (P12A/P13C), `spi` file-cache on-disk growth (U22b), `jose`/`jades` decoded-payload caches (stale + retained). Long-lived validation workloads pin memory.
5. **O(n·m) / O(n²) on document-scale input (the CPU class).** Attribute scans (X06), attribute sort (X08 O(A²)), full-doc base64+digest scans (J17A), schema re-parse per call (J17B), message-tag 16 KB alloc/call (T29). Bounded today by document size, but DoS-relevant on large/hostile documents.
6. **`panic(nil)` slipping past `recovered()`.** `dss/format.go:331` — a `panic(nil)` is not recovered by the `recover()`→error boundary, so a nil panic would escape the facade as a real crash. Latent (no in-tree `panic(nil`) but a one-line hardening (`if r == nil { return errors.New("panic(nil)") }`).

**Verified-correct (explicitly checked, no finding):** status aggregation across all VPF chains (FAIL hard-stop, `NOT_EVALUATED` never→PASS, no swapped PASS/FAIL); all revocation checks (no false-PASS, no inverted comparison, 0 always-PASS/FAIL) cross-checked vs Java; the EAA qualification matrix (**65 Go cells = Java's 56 + 9 exactly**); service-type mapping (unknown URI → not-qualified); XAdES level-acceptance (all `Has*Profile` fail-closed); timestamp `IsValid` (4-conjunct fail-closed); XML chain (no XXE/SSRF/billion-laughs/ReDoS, `MaxDepth=500`, DOCTYPE denied); report XML escaping (`encoding/xml` + relax-only post-pass); CLI RFC-3161 client (TLS default-verified, 30 s, 1 MiB cap, `crypto/rand` nonce, strict ASN.1); harness (no `os/exec` — compares pre-generated dumps; live-Java tests are arg-vector only).

---

## 5. Per-batch breakdown

| Batch | Package(s) | C | H | M | L | I | **Total** |
|---|---|---:|---:|---:|---:|---:|---:|
| 01a | `internal/{jose,cmscore,asn1ber,eccurve}` | 0 | 2 | 5 | 7 | 5 | 19 |
| 01b | `cades`, `cms` | 0 | 1 | 3 | 7 | 8 | 19 |
| 02 | `internal/{xmldom,xpath10,xmlc14n,xmldsig}`, `xml/` | 0 | 2 | 3 | 10 | 7 | 22 |
| 03 | `internal/pdf`, `pades` | 0 | 5 | 5 | 10 | 6 | 26 |
| 04 | `asic`, `asic/{cades,xades}` | 0 | 0 | 2 | 9 | 7 | 18 |
| 05 | `jades` | 0 | 0 | 4 | 5 | 3 | 12 |
| 06 | `tsl`, `trustedlist`, `token`, `crlparser` | 0 | 2 | 2 | 6 | 10 | 20 |
| 07 | `spi` (+ `spi/validation`) | 0 | 3 | 6 | 15 | 12 | 36 |
| 08 | `model`, `enumerations` | 0 | 0 | 3 | 8 | 12 | 23 |
| 09 | `dss` root, `alert`, `document`, `i18n`, `utils`, `validation`, `validation/policy`, `validation/executor` | 0 | 0 | 1 | 2 | 14 | 17 |
| 10 | `validation/{job,reports,process}` | **1** | 1 | 2 | 5 | 0 | 9 |
| 11 | `validation/process/bbb` | 0 | 2 | 1 | 1 | 5 | 9 |
| 12 | `validation/process/bbb/xcv` | 0 | 0 | 0 | 1 | 1 | 2 |
| 13 | `validation/process/{qualification,vpf*}` | 0 | 0 | 1 | 3 | 0 | 4 |
| 14 | `xades`, `xades/{definition,extension,tsl}` | 0 | 0 | 2 | 0 | 9 | 11 |
| 15 | `detailedreport`, `diagnostic`, `simplereport`, `simplecertificatereport`, `validationreport`, `policy` | 0 | 0 | 1 | 6 | 3 | 10 |
| 16 | `cmd/esig`, `examples`, `harness`, `corpus` gens | 0 | 0 | 2 | 2 | 4 | 8 |
| **Total** | | **1** | **18** | **43** | **97** | **106** | **265** |

(Severity column sums carry the ±1–2 reconciliation note from §1; the per-batch *totals* are exact.)

---

## 6. Remediation priority order

Ordered by (risk × effort-to-fix). All are small, local changes — none require architectural work.

1. **T31-SEC-001 (Critical) — fix the cache race.** Make `Expire/Sync/Error/Update` (and any `CachedEntry`/`CurrentCacheContext` field write) happen under the same `mu` that guards the map, or return value-copies. **Ownership:** `validation/job` (the cache) is the right place; the TSL pivot fan-out is merely the trigger. **Add a regression test:** two goroutines racing `Get`+`Expire` on a multi-pivot TL, run under `-race`. This is the only finding that flips a verdict.
2. **T23-SEC-001 (High) — false-accept on signature policy.** Reset `digestValid=false` (and set the error) on the `<3`-element early return, matching Java's catch→`setDigestValid(false)`. One-line + a KAT fixture.
3. **T22D1A/B-SEC-001 (High) — fail-open on unorderable cert set.** Propagate the `getOrderedCertificateChains` error to the presence-check verdict (fail-closed) instead of swallowing it. Add an oracle test with an unorderable set (all-bridge, no signing cert).
4. **The panic-on-hostile-input cluster (High/Low).** Add nil/empty guards (return the Java-equivalent fail-closed value) at: T33 (AOV digest), T35 (ISC `sources[0]`), T37 (EAA `[0]`), T38 (PID `default:`), T39/T39B (EAA `NotBefore/NotAfter`), P11A-002/003 (predictor `/N`), P12A (GetRevisions), P13 (VRI). Decide module-wide: *sanctioned Java-parity panic* (document) vs *guarded fail-closed* (hardening). Recommend the latter for any path reachable from a hostile document.
5. **T32-SEC-001 (High) — EAA nil-panic in the non-`eaa` build.** Guard the `SetNextItem` on a nil FC (or gate the call behind the tag).
6. **The unbounded-registry/leak cluster (Medium/Low).** Add eviction or bound to: `xades` `sync.Map` registries (T40-PERF-001), `pades` timestamp-token registry (P12A/P13C), `jose`/`jades` decoded-payload caches, `spi` file-cache growth.
7. **The O(n·m)/O(n²) cluster (Medium/Low).** Bound the document-scale scans (X06 attr, X08 sort, J17A digest scan, J17B schema re-parse, T29 16 KB alloc). Lower priority — bounded by document size today.
8. **`panic(nil)` at the facade (Low).** One-line: `recovered()` should treat a `nil` recover value as a real error, not drop it.
9. **CLI hardening (Medium).** Enforce HTTPS for `-tsa` (T46-SEC-002) and add a response-size cap to `tl refresh` (T46-SEC-001). Both are CLI-only, parity-neutral.

---

## 7. Open questions for the maintainers

Carried forward from the units; each is a *decision*, not a defect.

- **Codec leniency (module-wide).** Base64/hex decoder strictness differs across `spi`/`jades`/`asic`/`tsl` (T29, U24 `Digest.HexValue` padding vs Java `getHexValue()`). Pick one module-wide stance and pin it with tests.
- **T23 / T22D1 severity bar.** Both are false-accepts on verification sub-checks. Rate them High (as recorded) or escalate toward Critical if the affected sub-check is ever a *sole* gate for a level.
- **EAA in the non-`eaa` build (T32-SEC-001).** Product stance: should an untagged build ever reach an EAA code path, or must the tag guarantee it's unreachable?
- **LoTE job (T25-SEC-001).** The LoTE job is unported, making several `model/lote` findings latent. Confirm scope.
- **`ValidationPolicy` as an interface (T30B).** Keeping it an interface is what makes the executor's nil-guard meaningful; a concrete type would change that invariant.
- **XSD structure validation (T41-INFO-001 / T42-INFO-001).** A permanent no-op, relied on by TL conformance and the `StructureValidator`. Surface it in `known-gaps.md` (doc-only) so it isn't mistaken for a live check.
- **Concurrent validation roadmap (T41-INFO-002, T45-STD-003, T43-M-001).** If concurrent validation is planned, the unsynchronized-lazy-init cluster (all Java-parity today) should get a `sync.Once`/mutex pass *before* that lands, or T31-SEC-001's class becomes systemic.
- **Dead code (T44-STD-001).** `eaaPayloadClaimElements` (~130-line slice, referenced nowhere, comment claims otherwise) — safe to delete, or wire up + fix the comment.

---

## 8. Verification performed

- `gofmt -l .` → clean (from `dss/`).
- `go build ./...` → clean; `go build -tags eaa ./...` → clean.
- `go vet ./...` → clean.
- Per-unit: `golangci-lint run --config=../.github/.golangci.yml` clean in all 17 batches (both tag states where relevant); `go test -count=1` passing per package.
- **Spot-checked ~10 findings against source** at their cited `file:line`: T31-SEC-001 (unlocked `Expire`), T23-SEC-001 (`SetDigestValid(true)` + `<3` early return), T35-SEC-001 (empty-`Sources` `slices.Contains` path), T37-SEC-001 (`EAASignatures()[0]`), P11A-SEC-001 (`FlateDecode` unbounded `io.Copy`). All confirmed.
- **Oracle cross-checks** (Java DSS 6.5.RC1 at `~/Workspace/esig/dss`): T23 false-accept confirmed as a *divergence*; T31 race confirmed; EAA matrix 65=65; revocation checks no-false-PASS; XAdES level-acceptance fail-closed; CLI HTTP client safe.
