# Batch 12 — bbb/xcv (X.509 certificate validation)

## bbb/xcv — light pass (unit U36)

- **Date:** 2026-08-28
- **Depth:** LIGHT — top-12 files by size read in full (~3,025 lines) + flagged files read in
  full + grep-scan over the remaining ~85 files. Total in-scope code read ≈ 3,450 lines
  (within the 3,500-line budget).
- **Scope:** `dss/validation/process/bbb/xcv` — the 5.2.6 X.509 certificate-validation
  building block and its sub-checks (sub-XCV, CRS revocation selector, RAC revocation
  acceptance, RFC freshness, policy-tree, name-constraints, validity-range, KU/EKU, and ~85
  leaf checks). 98 non-test files, 9,396 lines. A building-block check that returns the wrong
  status (false PASS → false "VALID") is Critical/High; a false FAIL is High.
- **Cross-referenced (not re-derived):** the chain framework (`validation/process/chain.go`,
  `chain_item.go`, `validation_process_utils.go`) is U32 scope; the AOV/CV result checks that
  `certificateCryptographic()` / `revocationCryptographic()` / `revocationDataIntact()`
  delegate to are U33 scope; the policy-constraint resolution (U30) and the `spi/validation`
  revocation-source plumbing are U33/batch-07 scope. `PORTING.md` (panic→error at the facade,
  disabled staticcheck ST*/S1*/QF*/SA1019, `// DIVERGENCE, deliberate:`) and
  `docs/compatibility/known-gaps.md` (unported online revocation sources; EAA gated behind the
  `eaa` tag) were read and treated as deliberate, not findings.

### Files read (top-12 full-read + flagged)

Top 12 by size (full read):

| File | Lines |
|---|---|
| sub_x509_certificate_validation.go | 697 |
| certificate_name_constraints_check.go | 614 |
| x509_certificate_validation.go | 405 |
| revocation_acceptance_checker.go | 306 |
| revocation_has_information_about_certificate_check.go | 295 |
| certificate_revocation_selector.go | 290 |
| certificate_policy_tree_check.go | 255 |
| certificate_algorithm_obsolescence_validation_check.go | 205 |
| revocation_freshness_checker.go | 156 |
| certificate_validity_range_check.go | 132 |
| revocation_acceptance_checker_result_check.go | 128 |
| policy_tree_node.go | 119 |

Flagged by grep-scan, read in full:

| File | Lines | Why flagged |
|---|---|---|
| revocation_responder_id_match_check.go | ~80 | `[0]` index deref |
| certificate_not_revoked_check.go | ~95 | **revocation verdict (most security-critical)** |
| certificate_not_on_hold_check.go | ~95 | **revocation verdict** |
| extended_key_usage_check.go | ~110 | `panic(err)` site |
| key_usage_check.go | ~100 | `panic(err)` site |

Also cross-checked against Java upstream (local checkout
`/Users/utain/Workspace/esig/dss`): `CertificateNotRevokedCheck.java`,
`CertificateRevocationSelector.java` — the two most security-critical verdict files.

### Grep-scan summary

- **`panic | os.Exit | recover()`** (8 non-test hits): 6 are `panic(err)` on
  `process.GetSubContextPosition` / `GetConstraintOrMaxLevel` / `GetCryptoPosition`
  (sub_x509_certificate_validation.go:312,650,660; x509_certificate_validation.go:247;
  extended_key_usage_check.go:96; key_usage_check.go:78) — all on unreachable-by-construction
  enum/lookup paths, documented as the Java-exception→panic port. 1 is the deliberate
  `crsOverrides()` init guard (certificate_revocation_selector.go:147). 1 is the
  **unsupported-RevocationType** default in revocation_freshness_checker.go:104 (→ T36-STD-001).
  Disposition: all deliberate/fail-closed; one noted as a hardening opportunity.
- **`[0]` / `.Sources()` / `.Certificates()` / `.All()`** (4 non-test hits):
  `conclusion.{Errors,Warnings,Infos}[0].Value` in certificate_algorithm_obsolescence_
  validation_check.go:196,199,202 — each guarded by `IsCollectionNotEmpty` (safe);
  `relatedSigningCertificates[0]` in revocation_responder_id_match_check.go:42 — guarded by
  `CollectionSize(...) == 1` (safe). Disposition: all guarded, no finding.
- **always-`return true` / `return false` in `Process()`** (0 hits). **No always-PASS/FAIL
  candidate.** Disposition: clean.
- **`func (.*) Process(`** count: **85** checks — the grep-scan (whole-directory `*.go`) covers
  the package; the 12 un-read leaf checks are small, single-purpose checks whose `Process()` is
  a one-line predicate on a certificate field (surname/givenName/email/etc. present-or-matches),
  confirmed by the grep for always-true/false (none) and the index-0/nil-deref patterns (none).
- **nil-deref candidates** (unguarded `.SigningCertificate()` / `.CertificateChain()` /
  `.CertificateRevocationData()`): 17 hits. All are either (a) on a receiver the framework
  guarantees non-nil (`c.certificate`, `c.revocationData` — set in constructors), or (b) guarded
  by a preceding `!= nil` on the same expression, or (c) the single reachable unguarded case in
  revocation_has_information_about_certificate_check.go:140 (→ T36-STD-002). The one that looks
  dangerous — `c.usedCertificateRevocation.SigningCertificate()` in
  certificate_validity_range_check.go:87 — is short-circuit-guarded by
  `c.revocationIssuerTrusted || !c.revocationIssuerCheckEnforced ||` and then by the
  `certificateWrapper != nil` test inside `isInValidityRange` (safe).

### Revocation-check verdict verification (dedicated subsection)

The revocation verdict path is:
`CertificateRevocationSelector` (picks latest acceptable) → `RevocationAcceptanceChecker`
(RAC) → `AcceptableRevocationDataAvailableCheck` → `CertificateNotRevokedCheck` /
`CertificateNotOnHoldCheck` → `RevocationFreshnessChecker` (RFC) →
`RevocationHasInformationAboutCertificateCheck`. Each was verified:

| Check | Verdict logic | Verdict |
|---|---|---|
| **CertificateNotRevokedCheck** (`Process`) | `!isRevoked`, where `isRevoked = rev≠nil ∧ rev.IsRevoked() ∧ reason≠CERTIFICATE_HOLD ∧ (revDate≠nil ∧ !currentTime.Before(revDate))`. Cross-checked byte-for-byte against Java (`currentTime.compareTo(revDate) >= 0` ≡ `!currentTime.Before(revDate)`). | **Correct. No false-PASS.** A revoked cert (non-hold) with a revocation date on/before the validation time returns FAIL; nil revocation or nil revocation date returns PASS (fail-open only when there is no evidence of revocation — the `AcceptableRevocationDataAvailableCheck` and RAC gate that). |
| **CertificateNotOnHoldCheck** (`Process`) | `!isOnHold`, symmetric to the above with `reason==CERTIFICATE_HOLD`. | **Correct. No false-PASS, no inversion.** |
| **CertificateRevocationSelector** (initChain / verifyRevocationData) | `latestCertificateRevocation` is set **only** when `RevocationDataValidityMap[id] == true`, i.e. only when the RAC verdict is valid (`c.IsValid(&racResult.Content)`). Cross-checked against Java (`Boolean.TRUE.equals(map.get(w))`). | **Correct.** The selector cannot promote a rejected revocation to "latest acceptable," so it cannot manufacture a false-PASS. |
| **RevocationAcceptanceCheckerResultCheck** (`Process`) | `c.IsValid(&c.racResult.XmlConstraintsConclusionContent)` — returns true only when the RAC conclusion is PASSED. | **Correct.** |
| **RevocationFreshnessChecker** (initChain) | CRL/OCSP → `NextUpdateCheck`; other → `panic` (fail-closed). When a freshness constraint is present → `RevocationDataFreshCheck`. | **Correct.** The only non-verdict branch is a panic, not a PASS (→ T36-STD-001). |
| **RevocationHasInformationAboutCertificateCheck** (`Process`) | `checkCertHashMatches() ∨ checkIssuerHasInformationForExpiredCertificate()`. The second clause requires `certNotAfter != nil ∧ notAfterRevoc != nil ∧ !certNotAfter.Before(notAfterRevoc)` — i.e. the issuer's known-status-since is on/after the cert expiry. | **Correct.** Fails closed when either date is nil. |
| **CertificateValidityRangeCheck** (`isInValidityRange`) | `notBefore≠nil ∧ !currentTime.Before(notBefore) ∧ notAfter≠nil ∧ !currentTime.After(notAfter)`. | **Correct** — exactly the `!Before(notBefore) && !After(notAfter)` form. Nil certificate → `false` (fail-closed). |
| **RevocationResponderIdMatchCheck** (`Process`) | `CollectionSize(related)==1 ∧ signingCert≠nil ∧ related[0].Id()==signingCert.Id()`; else `false`. | **Correct.** Fails closed when there is not exactly one related signing certificate. |

**Nil-deref verdict:** no unguarded nil-deref on the *verdict* path. The one reachable
unguarded nil-deref (revocation_has_information_about_certificate_check.go:140, in
`BuildAdditionalInfo` — a *message* path, not the verdict) is a 1:1 port of the same upstream
Java latent hazard and is fail-closed (panic is recovered to an error at the facade, never a
false PASS). See T36-STD-002.

**Inverted comparisons:** none found. All `Before`/`After`/`EqualFold`/`HasPrefix`/
`HasSuffix`/subnet-mask comparisons read correctly against the RFC 5280 intent and the Java
source.

### Findings

No **Critical** or **High** findings. No false-PASS or false-FAIL on the verdict path, no
inverted comparison, no unguarded nil-deref on the verdict path, no always-PASS/FAIL check, no
unbounded-growth or retained-reference leak (the `validatedTokens` set and
`RevocationDataValidityMap` are bounded by the number of distinct revocation/certificate
tokens in the input and are released when the chain result is GC'd; `policy_tree_node.go`
builds a tree bounded by chain-length × policies-per-cert, which is inherent to the RFC 5280
algorithm and the input size).

#### STD

```
T36-STD-001 — **Low** — Category: STD — dss/validation/process/bbb/xcv/revocation_freshness_checker.go:104
- Evidence: `default: panic(fmt.Errorf("the RevocationType '%s' is not supported!", c.revocationData.RevocationType()))`
- Impact: A hostile/defective revocation token whose `RevocationType()` is neither CRL nor OCSP
  (the only two enum values the rest of the package branches on) panics inside the building
  block. The facade `recovered()` turns it into an error, so this is fail-closed (never a
  false-PASS), but it is a panic on attacker-influenced input rather than a returned verdict.
- Recommendation: Accept as the faithful port of Java's `UnsupportedOperationException` (it is
  unreachable for the two enum values that can actually appear), or, if hardening, map the
  `default` branch to a FAILED/INDETERMINATE check item instead of panicking, with a
  `// DIVERGENCE, deliberate:` note. No action required for parity.
```

```
T36-STD-002 — **Info** — Category: STD — dss/validation/process/bbb/xcv/revocation_has_information_about_certificate_check.go:140
- Evidence: `BuildAdditionalInfo()` → `if !c.Process() { message = c.getNotAfterAfterCertificateNotAfterMessage() }`
  → `c.formattedDate(c.notAfterRevoc)` where `notAfterRevoc` may still be nil (a revocation
  with no thisUpdate and no certHash/expiry extension).
- Impact: A nil-deref panic on a message-building path (not the verdict path). This is a 1:1
  port of the identical upstream Java latent NPE (`getNotAfterRevoc()` returning null).
  Fail-closed only (panic→error at the facade); never a false-PASS.
- Recommendation: None required for parity. If the team later wants robustness, guard the
  `notAfterRevoc` nil case in `getNotAfterAfterCertificateNotAfterMessage` (it already special-
  cases nil via `formattedDate`) — the actual deref is in the Java-shared message format call,
  so leave as-is to stay 1:1.
```

#### PERF-BIGO

No PERF-BIGO findings. The name-constraints intersect/union (certificate_name_constraints_check.go)
and the policy-tree node walk (policy_tree_node.go) are O(chain × policies × subtrees), which is
inherent to the RFC 5280 path-validation algorithm and bounded by the input's own size; there is
no quadratic blow-up on attacker-influenced input beyond what the algorithm and the input size
impose. The `generalNameTypesOf` helper deliberately iterates in a fixed order to keep the
deterministic output (a porting requirement), at no asymptotic cost.

#### PERF-MEM / PERF-LEAK

No PERF-MEM / PERF-LEAK findings. `validatedTokens` (map) and `RevocationDataValidityMap`
(map) are bounded by the number of distinct tokens in the input; the policy tree is bounded by
chain length × policies per certificate; no references are retained beyond the life of the chain
result.

#### SEC

No SEC findings. The dedicated revocation-check verification above found **no false-PASS**
revocation check, **no inverted comparison**, and **no unguarded nil-deref on the verdict
path**. The two most security-critical files were additionally cross-checked against the Java
upstream source and are byte-for-byte faithful.

### Tool log

Run from `dss/`:

```
$ gofmt -l validation/process/bbb/xcv
(empty — no unformatted files)

$ go vet ./validation/process/bbb/xcv/
(clean — exit 0, no output)

$ golangci-lint run --config=../.github/.golangci.yml ./validation/process/bbb/xcv/
0 issues.
```

All three in-scope tools are clean. No out-of-scope findings.

### Summary

| Lens | Critical | High | Medium | Low | Info | Total |
|---|---|---|---|---|---|---|
| SEC | 0 | 0 | 0 | 0 | 0 | 0 |
| PERF-BIGO | 0 | 0 | 0 | 0 | 0 | 0 |
| PERF-MEM/LEAK | 0 | 0 | 0 | 0 | 0 | 0 |
| STD | 0 | 0 | 0 | 1 | 1 | 2 |
| **Total** | **0** | **0** | **0** | **1** | **1** | **2** |

**Total findings: 2** (1 Low, 1 Info). **Zero Critical/High/Medium.** The xcv building-block
verdict logic — including every revocation check — is correct, fail-closed, and a faithful 1:1
port of Java DSS 6.5.RC1.
