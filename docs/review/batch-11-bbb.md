# Batch 11 — building-block checks (bbb)

## bbb/aov + bbb/cv (unit U33)

- **Date:** 2026-08-28
- **Depth:** full-read of every in-scope file (≤3,900-line budget; 3,921 lines read
  plus the framework and policy utils that gate the verdicts)
- **Scope:** the algorithm-obsolescence-validation (AOV) building block
  (`bbb/aov`) and the certificate-validation / cryptographic-verification (CV)
  building block (`bbb/cv`). A building-block check that returns the wrong status
  (false PASS → false "VALID") is Critical/High; a false FAIL is High. The AOV
  checks delegate their actual pass/fail decision to the cryptographic-suite
  policy (`validation/policy` cryptographic-suite utils, U30 scope) — those are
  cross-referenced, not re-derived. The chain framework (`process/chain.go`,
  `chain_item.go`, `uninterrupted_chain_item.go`, `bbb/abstract_*`) is the U32
  bbb-core framework and is cross-referenced, not re-derived.
- **Lens set:** SEC, PERF-BIGO, PERF-MEM / PERF-LEAK, STD.

### Scope (files actually read)

**`dss/validation/process/bbb/aov` — 31 files, 2,703 lines — read in full (100%).**

| File | Lines | | File | Lines |
|---|---|---|---|---|
| signing_certificate_ref_digest_algorithm_check.go | 273 (full) | | digest_algorithm_cryptographic_checker_result_check.go | 46 (full) |
| digest_algorithm_obsolescence_validation.go | 227 (full) | | evidence_record_algorithm_obsolescence_validation.go | 45 (full) |
| algorithm_obsolescence_validation_check.go | 171 (full) | | eaa_algorithm_obsolescence_validation.go | 44 (full) |
| algorithm_obsolescence_validation.go | 158 (full) | | signature_algorithm_obsolescence_validation.go | 43 (full) |
| abstract_cryptographic_checker_result_check.go | 152 (full) | | revocation_data_algorithm_obsolescence_validation.go | 43 (full) |
| signature_algorithm_cryptographic_checker.go | 129 (full) | | signature_signed_data_algorithm_obsolescence_validation.go | 40 (full) |
| signature_value_and_signed_attributes_algorithm_obsolescence_validation.go | 127 (full) | | token_certificate_chain_algorithm_obsolescence_validation.go | 39 (full) |
| token_algorithm_obsolescence_validation.go | 126 (full) | | eaa_revocation_algorithm_obsolescence_validation.go | 39 (full) |
| abstract_algorithm_cryptographic_checker.go | 99 (full) | | timestamp_algorithm_obsolescence_validation.go | 28 (full) |
| signature_algorithm_cryptographic_checker_result_check.go | 94 (full) | | public_key_size_known_check.go | 48 (full) |
| digest_algorithm_cryptographic_checker.go | 93 (full) | | digest_algorithm_reliable_check.go | 58 (full) |
| signature_algorithm_at_validation_time_check.go | 82 (full) | | signature_algorithm_reliable_check.go | 55 (full) |
| digest_algorithm_at_validation_time_check.go | 78 (full) | | certificate_algorithm_obsolescence_validation.go | 51 (full) |
| digest_matcher_cryptographic_checker_result_check.go | 74 (full) | | certificate_and_chain_algorithm_obsolescence_validation.go | 51 (full) |
| public_key_size_acceptable_check.go | 63 (full) | | algorithm_obsolescence_validation_check_with_id.go | 62 (full) |
| abstract_cryptographic_check.go | 65 (full) | | | |

**`dss/validation/process/bbb/cv` — 11 files, 1,218 lines — read in full (100%).**

| File | Lines | | File | Lines |
|---|---|---|---|---|
| cryptographic_verification.go | 312 (full) | | reference_data_name_match_check.go | 77 (full) |
| evidence_record_hash_tree_renewal_timestamp_check.go | 149 (full) | | at_least_one_reference_data_object_found_check.go | 65 (full) |
| reference_data_intact_check.go | 132 (full) | | reference_data_group_check.go | 64 (full) |
| reference_data_existence_check.go | 132 (full) | | signature_intact_with_id_check.go | 40 (full) |
| manifest_entry_existence_check.go | 85 (full) | | | |
| manifest_entry_group_check.go | 82 (full) | | signature_intact_check.go | 80 (full) |

**Cross-referenced (not re-derived, U30/U32 scope):** `validation/process/chain.go`,
`chain_item.go`, `uninterrupted_chain_item.go`, `bbb/abstract_{value,multi_values,certificate}_check_item.go`,
`validation/policy/cryptographic_suite_utils.go` (the functions the AOV checks gate on),
`policy/etsi_validation_policy.go` (constraint resolution), `diagnostic/certificate_ref_wrapper.go`
(`DigestMethod()`), `utils/strings.go` (`IsStringDigits`), `utils/collections.go`
(`JavaHashMapComputeIfAbsentKeyOrder`).

### Check verdict verification (dedicated subsection)

Every check's `Process() bool` was traced to the condition it returns true for.
The chain framework (`ChainItemBase.fail`/`informOrWarn`, `chain_item.go`) calls
`overrides.Process()` and records `OK` on `true`, `NOT_OK` on `false`; the chain
default-concludes `PASSED` when no item recorded a failure (`chain.go` `Execute`).
A check therefore yields a false "VALID" only if its `Process()` returns `true` on
a condition that is actually bad, or if its `Level()` is empty so it is skipped.

**No always-PASS and no always-FAIL checks, and no inverted comparisons, exist in
either package.** The decision logic is correct for each:

AOV (all delegate to the U30-verified cryptographic-suite policy utils —
cross-referenced, not re-derived):
- `SignatureAlgorithmReliableCheck.Process` → `IsSignatureAlgorithmReliable` =
  algorithm present in the acceptable set (fail-closed: `false` when absent).
- `DigestAlgorithmReliableCheck.Process` → `IsDigestAlgorithmReliable` = digest
  present in the acceptable set (fail-closed).
- `PublicKeySizeKnownCheck.Process` → `IsStringDigits(keyLength)` = key size known
  (non-empty, all digits) — correct; `""`/non-numeric → FAIL.
- `PublicKeySizeAcceptableCheck.Process` → `IsSignatureAlgorithmWithKeySizeReliable`
  = key size within each applicable evaluation's `[Min,Max]` (`*keySize < Min` or
  `> Max` → not applicable). Comparison direction correct (`>= min`, `<= max`).
- `SignatureAlgorithmAtValidationTimeCheck.Process` → `IsSignatureAlgorithmReliableAtTime`
  = an applicable evaluation's validity window `[validityStart, validityEnd]`
  contains `validationDate` (`Before`/`After` boundaries correct).
- `DigestAlgorithmAtValidationTimeCheck.Process` → `IsDigestAlgorithmReliableAtTime`
  = same window test on the digest evaluation.
- The four result-check leaves
  (`AbstractCryptographicCheckerResultCheck`, `SignatureAlgorithm…`,
  `DigestAlgorithm…`, `DigestMatcher…`) gate on
  `IsValid(conclusion) && allCCConstraintsValid(cc)` = conclusion is
  `PASSED`/`TOTAL_PASSED` **and** every constraint is `OK`/`IGNORED` — i.e. no
  `NOT_OK`, `WARNING`, or `INFORMATION` constraint. Correct (not inverted).
- `AlgorithmObsolescenceValidationCheck.Process` →
  `IsValid(aovResult)` + `IsValidConclusion` (conclusion `PASSED`/`TOTAL_PASSED`
  and no warnings/infos). Correct.
- The `…ObsolescenceValidation` builder classes (`Token…`, `Signature…`,
  `Certificate…`, `EvidenceRecord…`, `EAA…`, `RevocationData…`, `Timestamp…`,
  `SignatureSignedData…`, `TokenCertificateChain…`, `CertificateAndChain…`,
  `EAARevocation…`) only *assemble* the chain of the above checks; they carry no
  independent verdict.

CV (all gate directly on the diagnostic `XmlDigestMatcher` / token flags — the
correct input):
- `ReferenceDataExistenceCheck.Process` → `digestMatcher.DataFound` (data present
  ⇒ PASS). Correct.
- `ReferenceDataIntactCheck.Process` → `digestMatcher.DataIntact` (integrity
  holds ⇒ PASS). Correct.
- `ReferenceDataNameMatchCheck.Process` → `uri != nil && uri == documentName`
  (reference name matches ⇒ PASS). Correct.
- `SignatureIntactCheck.Process` → `token.IsSignatureIntact()` (signature value
  verifies ⇒ PASS). Correct; `SignatureIntactWithIdCheck` inherits it.
- `ManifestEntryExistenceCheck.Process` → at least one manifest entry `DataFound`
  ⇒ PASS. Correct.
- `ManifestEntryGroupCheck.Process` → `true` unless a manifest entry is `!DataFound`
  ⇒ PASS. Correct (vacuously true when no manifest entries, but only reachable via
  `containsManifestEntries`, which requires ≥1 `DataFound` manifest entry).
- `AtLeastOneReferenceDataObjectFoundCheck.Process` → ≥1 `DataFound` ⇒ PASS. Correct.
- `ReferenceDataGroupCheck.Process` → `false` if any
  `EvidenceRecordOrphanReference` present, else `true`. Correct (orphan digests
  at the first hash-tree level are disallowed).
- `EvidenceRecordHashTreeRenewalTimestampCheck.Process` → the renewal timestamp's
  covered documents are a superset of the evidence record's covered documents.
  Correct (subset/coverage test, not inverted).

**Nil-derefs on hostile input:** the AOV result-check leaves are safe because
`ccResult` is always the non-nil result of `Execute()` (the framework's
`ChainBase.Execute` guarantees a conclusion is set and `AddAdditionalInfo` sets
`CryptographicValidation` before any consumer reads it). The one genuine nil-deref
is T33-SEC-001. The CV `EvidenceRecordHashTreeRenewalTimestampCheck` carries a
deliberate `panic` (1:1 port of Java's `IllegalStateException`) — see
T33-STD-002. No other unchecked deref, unchecked type assertion, or
`interface{}` deref is reachable on the verdict path.

### Findings

**SEC**

T33-SEC-001 — **High** — Category: SEC — dss/validation/process/bbb/aov/signature_value_and_signed_attributes_algorithm_obsolescence_validation.go:104
- Evidence:
  ```go
  // signature_value_and_signed_attributes…validation.go:101-104 (builder, chain-construction time)
  cryptoValidationResult := signCertCheck.CryptographicValidationResult()   // may be nil
  if cryptographicValidation == nil || (c.isValid(cryptographicValidation) &&
      !c.IsValid(&cryptoValidationResult.XmlConstraintsConclusionContent)) {   // line 103: latent deref
      cryptographicValidation = cryptoValidationResult.CryptographicValidation   // line 104: panics (nil deref)
  ```
  where `CryptographicValidationResult()` → `validateCertReferences()` returns the
  `nil` `c.cryptographicValidationResult` when every `certificateRef` in the group
  has an empty digest method:
  ```go
  // signing_certificate_ref_digest_algorithm_check.go:102-124
  func (c *SigningCertificateRefDigestAlgorithmCheck[T]) validateCertReferences() *jaxb.XmlCC {
      ...
      for _, certificateRefWrapper := range c.certificateRefs {
          digestAlgorithm := certificateRefWrapper.DigestMethod()
          if digestAlgorithm == "" {
              continue                                   // never assigns cryptographicValidationResult
          }
          ...
      }
      return c.cryptographicValidationResult             // nil when all were skipped
  }
  ```
  (`CertificateRefWrapper.DigestMethod()` returns `""` when the certificate
  reference has no `DigestAlgoAndValue`/`DigestMethod` —
  `diagnostic/certificate_ref_wrapper.go:103-112`.)
- Impact: A signing token (CMS/XAdES/PAdES) whose `signingCertificateV2` /
  `SigningCertificate` attribute references at least one certificate but **every**
  referenced certificate lacks a digest method (a malformed/hostile CMS structure —
  `DigestMethod()` returns `""` for such a ref) panics during `Execute()` of the
  AOV chain: `buildSignedAttributesValidationChain` dereferences the `nil`
  `*jaxb.XmlCC` at line 104 (chain-construction time, before any check `Process()`
  runs). Note the *check's own verdict* is safe — `Process()` goes through
  `validResult`, which guards `ccResult == nil` and returns `false` — so this is **not**
  a false PASS and not a false FAIL; it is a **panic on hostile input**.
  - Through the public facade (`dss.Validate`) it is caught: `format.go:329
    recovered()` wraps the whole call (`validate.go:154`), so the malformed
    signature that *should* yield a verdict instead surfaces as a hard error — a
    false-"INVALID" in the user-visible sense (error where a report was expected).
  - Through the lower-level packages (a caller that uses `validation/…` without the
    facade, or any path not wrapped in `recovered()`) it is an uncaught panic →
    process crash, whereas Java DSS surfaces the same case as a caught
    `NullPointerException`.
  Reachable but narrow: requires ≥1 cert ref (so `IsSigningCertificateReferencePresent()`
  is true and the builder runs) with all refs digest-method-less.
- Recommendation: Guard the deref sites exactly the way the check's own `validResult`
  helper already does. Concretely: (a) at the builder, check
  `if cryptoValidationResult != nil` before lines 103-104 (skip storing a nil result);
  (b) make `extractXmlMessage` nil-tolerant — `if c.cryptographicValidationResult == nil
  { return nil }` at the top (line ~246); (c) optionally have `validateCertReferences`
  fall back to a synthesized non-nil `XmlCC` (conclusion `INDETERMINATE`) when all refs
  are skipped, so the report is well-formed. Add a KAT test: a signature whose
  `signingCertificateRef` entries all lack a digest method must not panic and must
  conclude `INDETERMINATE`/`FAILED` (per policy level), matching upstream.

**PERF**

T33-PERF-BIGO-001 — **Medium** — Category: PERF-BIGO — dss/validation/process/bbb/cv/evidence_record_hash_tree_renewal_timestamp_check.go:85
- Evidence:
  ```go
  for _, originalDataObject := range evidenceRecordCoveredDocuments {
      index := -1
      for i, covered := range timestampCoveredDocuments {
          if covered == originalDataObject { index = i; break }
      }
      if index < 0 { return false }
      // remove object to avoid checking duplicates
      timestampCoveredDocuments = append(timestampCoveredDocuments[:index], timestampCoveredDocuments[index+1:]...)
  }
  ```
- Impact: O(n·m) in the number of covered archive-data-objects of the evidence
  record (n) and the renewal timestamp (m); the inner slice removal is O(m) each,
  so total O(n·m). Both lists are attacker-influenced (they come from the
  evidence-record / timestamp `DigestMatcher` set in the diagnostic data). For a
  pathological evidence record with thousands of archive objects this is a
  quadratic allocation hot spot. Bounded in practice for real evidence records,
  hence Medium, not High.
- Recommendation: Replace the linear scan + slice-removal with a
  `map[string]int` (count) of `timestampCoveredDocuments`; for each
  `evidenceRecordCoveredDocuments` entry decrement the count and fail if it
  reaches 0 or is absent. O(n+m) expected, no quadratic allocation. (A faithful
  upstream-equivalent fix; verify the resulting verdict is identical — it is,
  since this is a pure multiset-subset test.)

**STD**

(The line-104 nil-deref in `signature_value_and_signed_attributes…validation.go`
is reported once under SEC — T33-SEC-001 — and is not double-counted here as an
idiom nit; the builder reading `cryptoValidationResult.CryptographicValidation`
without a local nil guard is the same root cause.)

T33-STD-002 — **Info** — Category: STD — dss/validation/process/bbb/cv/evidence_record_hash_tree_renewal_timestamp_check.go:67
- Evidence:
  ```go
  panic(fmt.Sprintf("Not found a corresponding evidence record for a time-stamp with Id '%s'",
      timestampWrapper.Id()))
  ```
- Impact: Deliberate — the comment names it a port of Java's
  `IllegalStateException` reporting a diagnostic-data inconsistency, not a
  validation outcome, and `process()` has no error channel. This is a faithful
  1:1 port and the panic is only reachable if the diagnostic data is internally
  inconsistent (a time-stamp token present without a matching evidence record),
  which is a programming/data error, not hostile document input. Noted, not a
  defect.
- Recommendation: None — keep. (Recorded here so it is not mistaken for a
  missing-guard bug in a later pass.)

No PERF-MEM / PERF-LEAK findings: the AOV/CV checks allocate only per-check report
structures (`XmlConstraint`, `XmlMessage`, small `XmlCryptographicValidation`)
and short-lived dedup slices/maps bounded by the number of digest matchers /
signing-certificate references / chain certificates — all O(input size), no
unbounded growth, no retained references outliving the result object. No PERF-BIGO
findings beyond T33-PERF-BIGO-001. No other STD findings: all type assertions use
the comma-ok form, no unchecked `interface{}` deref, all pointer members are
nil-guarded at use (the one unguarded site is reported as T33-SEC-001), status
enums are the correct generated `jaxb` values, and no `ST*`/`S1*`/`QF*`/`SA1019`-
style or legacy-crypto patterns appear (deliberately out of scope per PORTING.md).

### Tool log (run from inside `dss/`)

```
$ gofmt -l validation/process/bbb/aov validation/process/bbb/cv
(no output — clean)

$ go vet ./validation/process/bbb/aov/... ./validation/process/bbb/cv/...
(no output, exit 0)

$ golangci-lint run --config=../.github/.golangci.yml \
    ./validation/process/bbb/aov/... ./validation/process/bbb/cv/...
0 issues.
```

All three tool reports are clean → no tool-driven findings. (The one `nolint:staticcheck`
in `cv/cryptographic_verification.go` marks a deliberate dead-store port of Java's
trailing `item = item.setNextItem(...)` and is correctly suppressed.)

Out-of-scope packages (bbb/fc, bbb/isc, bbb/sav, bbb/vci, bbb/xcv,
bbb/qualification, bbb/vpf*) were not reviewed in this unit — noted and left to
their own units.

### Summary

- Total findings: **3**
- Critical: **0** — High: **1** — Medium: **1** — Low: **0** — Info: **1**
- SEC: 1 (1 High) · PERF-BIGO: 1 (1 Medium) · PERF-MEM/LEAK: 0 · STD: 1 (1 Info)
  (Low: 0 — the line-104 nil-deref is reported once under SEC, not split into a
  separate Low idiom nit)
- Check-verdict verdict: **no always-PASS, no always-FAIL, no inverted comparisons**
  in either package; the AOV checks gate on the correct inputs (algorithm OID + key
  size + validation date) and the CV checks gate on the correct inputs (the digest
  matcher's `DataFound`/`DataIntact`/`Uri`/`DocumentName` and the token's
  `IsSignatureIntact`). The single nil-deref (T33-SEC-001) is on the **chain-construction
  (builder) path** of `signingCertificateRef` digest validation — `Process()` itself is
  safe — so it is a panic/DoS on hostile input, not a false PASS or false FAIL verdict.
- Open question for consolidation: confirm T33-SEC-001's upstream parity — i.e. that
  Java DSS 6.5.RC1's `SigningCertificateRefDigestAlgorithmCheck` / caller does not
  `NullPointerException` on an all-empty-digest-method group (if it does, the Go port
  should match by converting the panic to an explicit error/`INDETERMINATE` rather
  than crashing the process; if it doesn't, this is a Go-port regression to fix).

## bbb/fc (unit U34)

- **Date:** 2026-08-28
- **Depth:** full-read of every in-scope file (≤3,200-line budget; 3,172 lines read
  across 45 files, 100%). Every check's `Process()` decision logic was traced to the
  condition it returns true for and cross-referenced against the **Java upstream source,
  available locally at `/Users/utain/Workspace/esig/dss`** (`dss-validation/.../bbb/fc/
  checks/*.java`), so verdict correctness was verified against the authoritative port
  target rather than re-derived. The chain framework (U32), the policy constraint
  resolution (U30), and the EAA FC stub (U32, T32-SEC-001) are cross-referenced, not
  re-derived.
- **Scope:** the "5.2.2 Format Checking" building block (`bbb/fc`) — signature format,
  timestamp format, file coverage (all-files-signed / signed-and-timestamped-files-
  covered), byte-range (valid / collision / all-document), signature-dictionary, page /
  annotation / visual difference, MDP / SigFieldLock / FieldMDP, filename adherence
  (signature / manifest / timestamp), container type / mimetype / zip-comment /
  manifest-present / signed-files-present, full-scope, signer-information-store, CAdES-
  V3 hash-index, PDF/A profile / compliance, elliptic-curve key-size, and the EAA +
  EAA-revocation FC (eaa-tagged).
- **Lens set:** SEC, PERF-BIGO, PERF-MEM / PERF-LEAK, STD.

### Scope (files actually read)

**`dss/validation/process/bbb/fc` — 45 files, 3,172 lines — read in full (100%).**

| File | Lines | | File | Lines |
|---|---|---|---|---|
| abstract_signature_format_checking.go | 255 (full) | | pdfa_profile_check.go | 48 (full) |
| all_files_signed_check.go | 125 (full) | | pdfa_compliance_check.go | 47 (full) |
| elliptic_curve_key_size_check.go | 121 (full) | | mime_type_file_present_check.go | 47 (full) |
| signature_format_checking.go | 111 (full) | | byte_range_check.go | 46 (full) |
| timestamp_manifest_filename_adherence_check.go | 108 (full) | | eaa_revocation_format_checking.go | 44 (full, eaa) |
| abstract_signed_and_timestamped_files_covered_check.go | 101 (full) | | field_mdp_check.go | 43 (full) |
| timestamp_format_checking.go | 97 (full) | | abstract_format_checking.go | 40 (full) |
| byte_range_collision_check.go | 91 (full) | | signed_and_timestamped_files_covered_check.go | 27 (full) |
| signed_files_present_check.go | 86 (full) | | signature_filename_adherence_check.go | 83 (full) |
| sig_field_lock_check.go | 82 (full) | | references_not_ambiguous_check.go | 80 (full) |
| timestamp_filename_adherence_check.go | 78 (full) | | pdf_visual_difference_check.go | 75 (full) |
| doc_mdp_check.go | 75 (full) | | doc.go | 68 (full) |
| abstract_pdf_lock_dictionary_check.go | 75 (full) | | filename_adherence_check.go | 65 (full) |
| signature_manifest_filename_adherence_check.go | 74 (full) | | byte_range_all_document_check.go | 62 (full) |
| eaa_format_checking.go | 74 (full, eaa) | | pdf_annotation_overlap_check.go | 57 (full) |
| manifest_file_present_check.go | 57 (full) | | signer_information_store_check.go | 56 (full) |
| cades_v3_hash_index_check.go | 56 (full) | | container_type_check.go | 54 (full) |
| pdf_signature_dictionary_check.go | 53 (full) | | full_scope_check.go | 53 (full) |
| format_check.go | 53 (full) | | annotation_changes_check.go | 53 (full) |
| acceptable_mimetype_file_content_check.go | 52 (full) | | zip_comment_present_check.go | 51 (full) |
| form_fill_changes_check.go | 51 (full) | | undefined_changes_check.go | 50 (full) |
| pdf_page_difference_check.go | 50 (full) | | acceptable_zip_comment_check.go | 50 (full) |
| signature_not_ambiguous_check.go | 48 (full) | | | |

**Cross-referenced (not re-derived):** the chain framework (`process/chain.go` `Execute`,
`process/chain_item.go` `ChainItemBase.fail`/`informOrWarn`), the bbb core
(`bbb/abstract_multi_values_check_item.go` `ProcessValueCheck`/`ValuesCheck`), the policy
constraint resolution (U30), the diagnostic-data APIs the checks gate on
(`diagnostic/diagnostic_data.go` `ManifestFileForFilename`/`ContainerType`/`ContainerInfo`,
`diagnostic/pdf_revision_wrapper.go` `SignatureByteRange`/`IsSignatureByteRangeValid`/
`DocMDPPermissions`/`FieldMDP`/`SigFieldLock`, `diagnostic/jaxb/jaxb_wrappers.go` nil-safe
`All()`), the `PdfLockAction` enum (`enumerations/pdf_lock_action.go`), the PAdES-side
field-lock semantics (`pades/pdf_permissions_checker.go` `isSignatureFieldCreationForbidden`),
and the byte-range model (`pades/byte_range.go`). The EAA FC stub is the U32 T32-SEC-001
nil-panic (cross-referenced, not re-derived).

### Check verdict verification (dedicated subsection)

Every check's `Process() bool` was traced to the condition it returns true for, and each was
verified **against the Java upstream source** (local checkout). The chain framework (U32) is
the verdict engine: `ChainItemBase` runs `Process()` and records `OK` on `true`, `NOT_OK` on
`false`; the chain concludes `PASSED` only when no item failed. A check yields a false "VALID"
only if its `Process()` returns `true` on a condition that is actually bad, or is inverted.

**Verdict: NO always-PASS and NO always-FAIL checks, NO inverted comparisons, and NO
reachable nil-deref panics on the verdict path in `bbb/fc`.** Every check gates on the correct
input and the comparison direction is correct. Per-check summary (Go ↔ Java, all faithful):

- `FormatCheck` → `ProcessValueCheck(SignatureLevel.String())` = format in the acceptable set
  (or `*` wildcard); fail-closed when absent. Faithful (Java `getSignatureFormat().toString()`).
- `SignatureNotAmbiguousCheck` → `!IsSignatureDuplicated()`. Correct.
- `ReferencesNotAmbiguousCheck` → `false` on any digest matcher with `Duplicated == true`, else
  `true`. Correct (not inverted).
- `FullScopeCheck` → `true` when **every** scope is `FULL`; `true` for empty scope list. Faithful
  (Java guards with `isCollectionNotEmpty`; a non-FULL scope fails). **Not** a vacuous-pass bug —
  empty is the faithful Java behavior.
- `SignerInformationStoreCheck` → exactly one store entry. Correct.
- `EllipticCurveKeySizeCheck` → non-ECDSA passes; ECDSA requires an authorized digest
  (SHA-256/384/512) **and** a key size matching the digest (256/384/521); missing any of the
  three known-values fails. Faithful (Java `!isECDSA || (authorized && keySizeCorresponds)`).
- `ByteRangeCheck` → `IsSignatureByteRangeValid()`. Correct.
- `ByteRangeAllDocumentCheck` → `false` if any signature/timestamp has a PDF revision with an
  invalid byte range; `true` otherwise. Correct (gates on every revision, not one).
- `ByteRangeCollisionCheck` → `false` if the current token's byte range collides with any other
  signature/timestamp's. Interval-overlap test `a1<b1 != a1<b2` is correct (not inverted); the
  `!= 4` `panic` is a faithful port of Java's `IllegalStateException` on malformed byte ranges.
- `PdfSignatureDictionaryCheck` → `IsPdfSignatureDictionaryConsistent()`. Correct.
- `PdfPageDifferenceCheck` / `PdfAnnotationOverlapCheck` / `PdfVisualDifferenceCheck` → `true`
  iff the concerned-page list is empty. Correct.
- `DocMDPCheck` / `SigFieldLockCheck` (certification-permission ladder) → `false` when a change
  class is forbidden for the declared permission (NoChange/Minimal/Changes permitted). Faithful
  (matches the `CertificationPermission` ladder); no modification detected → `true`.
- `FieldMDPCheck` / `AbstractPdfLockDictionaryCheck` (Action ladder) → `ALL` ⇒ any modified field
  fails; `EXCLUDE` ⇒ a modified field **not** in the locked list fails; `INCLUDE` ⇒ a modified
  field **in** the locked list fails. **Verified correct** against both the `PdfLockAction` enum
  doc and the PAdES-side `isSignatureFieldCreationForbidden` (all three agree: INCLUDE = listed
  fields locked, EXCLUDE = non-listed fields locked). No swap.
- `ContainerTypeCheck` / `AcceptableMimetypeFileContentCheck` / `AcceptableZipCommentCheck` /
  `PDFAProfileCheck` → `ProcessValueCheck(...)` membership in the acceptable set. Correct.
- `ZipCommentPresentCheck` → `TrimSpace(comment) != ""`. Correct.
- `MimeTypeFilePresentCheck` / `PDFAComplianceCheck` → return the boolean flag. Correct.
- `ManifestFilePresentCheck` → ASiC-E requires ≥1 manifest file; ASiC-S is optional → `true`.
  Correct.
- `SignedFilesPresentCheck` → ASiC-S exactly one root-level content file; ASiC-E ≥1 content file.
  Correct.
- `AllFilesSignedCheck` → ASiC-S exactly one content file; ASiC-E the related manifest's entries
  cover all content files (CAdES requires a manifest, XAdES falls back to FULL scopes). Faithful
  (Java `AllFilesSignedCheck.process()` matches line-for-line).
- `SignedAndTimestampedFilesCoveredCheck` / `AbstractSignedAndTimestampedFilesCoveredCheck` →
  no manifest ⇒ `true` (ASiC-S case); else the manifest's entries are recursively covered and at
  least one root-level document is covered. Faithful (Java `checkManifestFilesCoveredRecursively`
  matches, including the `rootProcess` first-pass skip and `checkedEntries` memo).
- `SignatureFilenameAdherenceCheck` / `TimestampFilenameAdherenceCheck` / `SignatureManifest-
  FilenameAdherenceCheck` / `TimestampManifestFilenameAdherenceCheck` → exact / prefix / contains /
  suffix filename rules per ASiC-S vs ASiC-E and signature form; unsupported forms/container
  types `panic`. Faithful (Java `process()` bodies match; the `panic`s are the Java
  `UnsupportedOperationException` / `IllegalArgumentException` equivalents on a diagnostic
  inconsistency, not a validation outcome).
- `CAdESV3HashIndexCheck` → `IsAtsHashIndexValid()` for a CAdES-V3 archive timestamp, `true`
  otherwise. Correct.
- EAA / EAA-revocation FC (eaa-tagged) → wire the U32/EAA-check `Process()` delegates; the EAA
  FC-stub nil-panic is T32-SEC-001 (cross-referenced, not re-derived). The `eaa`-tagged files
  only build a chain of delegates and carry no independent verdict.

**Nil-deref panics on hostile input:** none reachable on the verdict path. The jaxb collection
wrappers (`ManifestFiles`, `ContentFiles`, `Entries`, `Signatures`, `TimestampList`) all have
nil-safe `All()` (return `nil` on a nil receiver), and the diagnostic-data getters
(`ManifestFileForFilename`, `ContainerType`, `ContainerInfo`, `ZipComment`, `MimetypeFileContent`,
`IsPDFACompliant`, …) guard nil. `byteRangesCollide` panics only on a byte range whose length is
not 4, but `collideRevisions` returns `false` when either revision's `PDFRevision()` is nil and
the byte-range value is only ever populated by the reviewed diagnostic-data builder (which
produces a well-formed 4-element range from a real PDF `/ByteRange`) — so it is a faithful
upstream invariant guard, not a reachable hostile-input panic. The `panic(err)` after
`SignatureFormat().SignatureForm()` in `SignatureFilenameAdherenceCheck` / `SignatureManifest-
FilenameAdherenceCheck` fires only for a `SignatureLevel` with no form (`SignatureLevelUnknown`),
which the upstream also rejects — a faithful invariant guard, not a validation-outcome bug.

### Findings

No **SEC** findings. No **always-PASS/always-FAIL** check, no **inverted comparison**, no
reachable **nil-deref panic** on the verdict path, and no **false "VALID"** / false "INVALID"
verdict — every check is a faithful 1:1 port verified against the Java upstream.

**PERF-BIGO**

```
T34-PERF-BIGO-001 — **Low** — Category: PERF-BIGO — dss/validation/process/bbb/fc/abstract_signed_and_timestamped_files_covered_check.go:49-73
- Evidence: checkManifestFilesCoveredRecursively recurses over each manifest entry that is itself a
  manifest, with only a checkedEntries memo per root call and no bound on nesting depth or on the
  product of entries across the manifest tree:
  entryManifest := c.DiagnosticData.ManifestFileForFilename(manifestEntry)
  if entryManifest != nil && !c.checkManifestFilesCoveredRecursively(coveredEntries, entryManifest.Entries.All(), ...)
- Impact: The work is O(1) per (root, distinct-entry) pair for a *DAG* of manifests, but the memo is
  shared across the whole traversal so a deep/wide manifest nesting that is attacker-influenced (a
  crafted ASiC-E diagnostic data) can produce quadratic-to-exponential work and a deep recursion
  (unbounded stack) before the memo caps re-entry. Real evidence-record / ASiC-E manifests are small
  and shallow, so this is a Low hardening note, not an exploitable blowup.
- Recommendation: Bound recursion depth (or convert to an iterative work-list) and cap the total number
  of manifest entries visited; the verdict is a pure multiset-subset test, so an iterative rewrite is
  behavior-identical. (Faithful to Java, which has the same unbounded recursion — so this is parity,
  noted for hardening only.)
```

**STD**

```
T34-STD-001 — **Info** — Category: STD — dss/validation/process/bbb/fc/byte_range_collision_check.go:66-72
- Evidence: firstByteRangePartLength / secondByteRangePartLength truncate each *big.Int with
  .Int64():
  return byteRange[0].Int64() + byteRange[1].Int64()
- Impact: A faithful 1:1 port of Java's BigInteger.intValue() (Java truncates to int; Go to int64),
  so a byte-range part value ≥ 2^63 would truncate identically in spirit. Real PDF /ByteRange parts
  are far below 2^63 and the value is only ever set by the reviewed diagnostic-data builder from a
  real /ByteRange, so this is not reachable as a verdict change — recorded as an Info so the
  Int64() truncation is not mistaken for a missing guard in a later pass.
- Recommendation: None required for correctness (parity). If hardening, compare with .Cmp on the full
  *big.Int values instead of summing truncated Int64() parts.
```

No **PERF-MEM / PERF-LEAK** findings: the checks allocate only short-lived per-check report
structures and the `checkedEntries` memo (bounded by the manifest-entry count for one
`Process()` call, released with the result object). No unbounded growth, no retained reference
outliving the result object. No other **STD** findings: all type assertions use the comma-ok form,
no unchecked `interface{}` deref, all pointer members are nil-guarded at use (the nil-safe jaxb
`All()` and the diagnostic-data getters back this), status enums are the correct generated jaxb
values, and no `ST*`/`S1*`/`QF*`/`SA1019`-style or legacy-crypto patterns appear (deliberately out
of scope per PORTING.md). The two `//nolint:staticcheck` annotations (in `SignatureFormatChecking`
and `TimestampFormatChecking` `InitChain`) mark a deliberate dead-store port of Java's trailing
`item = item.setNextItem(...)` and are correctly suppressed.

### Tool log (run from inside `dss/`)

| # | Command | Result |
|---|---|---|
| 1 | `gofmt -l validation/process/bbb/fc` | **(empty — no unformatted files)** |
| 2 | `go vet ./validation/process/bbb/fc/` | **clean (exit 0)** |
| 3 | `golangci-lint run --config=../.github/.golangci.yml ./validation/process/bbb/fc/` | **0 issues** |
| 4 | `go build ./validation/process/bbb/fc/` (default) / `go build -tags eaa ./validation/process/bbb/fc/` | **both exit 0** — the eaa/noeaa split compiles in both tag states (golangci-lint v2 has no `-tags` CLI flag, so the eaa build is verified via `go build -tags eaa`, per CLAUDE.md) |

All in-scope tools are clean → no tool-driven findings.

### Summary

- Total findings: **2**
- Critical: **0** · High: **0** · Medium: **0** · Low: **1** · Info: **1**
- By lens: **SEC 0** · **PERF-BIGO 1** (Low) · **PERF-MEM/LEAK 0** · **STD 1** (Info)
- Check-verdict verdict: **no always-PASS, no always-FAIL, no inverted comparisons, no reachable
  nil-deref panics** — every fc check gates on the correct input and returns the correct status for
  the correct condition, verified 1:1 against the Java upstream. The package is SEC-clean.
- Zero-finding lenses: **SEC** (no false-PASS/FAIL, no inverted comparison, no reachable nil-panic),
  **PERF-MEM/PERF-LEAK** (no unbounded retention).
- Open questions for consolidation:
  1. **T34-PERF-BIGO-001 (unbounded manifest recursion):** is an iterative / depth-bounded rewrite
     wanted for hardening, or is the Java-parity unbounded recursion acceptable as-is (it is not
     reachable as a verdict change on real input)?
  2. **T34-STD-001 (`Int64()` byte-range truncation):** confirm it is acceptable parity with Java's
     `intValue()` and needs no `*big.Int.Cmp` hardening.

## bbb/isc + bbb/sav + bbb/vci (unit U35)

- **Date:** 2026-08-28
- **Depth:** full-read of every in-scope file (≤4,300-line budget; 4,271 lines read
  across 52 files, plus cross-referenced framework and diagnostic-data dependencies)
- **Scope:** the identification-of-signing-certificate (ISC), signature-acceptance-
  validation (SAV), and validation-context-initialization (VCI) building blocks.
  A building-block check that returns the wrong status (false PASS → false "VALID")
  is Critical/High; a false FAIL is High. The AOV result-check that SAV's `cryptographic()`
  step delegates to is U33 scope (cross-referenced, not re-derived). The EAA checks
  delegated from `eaa_acceptance_validation.go` live in `validation/process/eaa/checks`
  and are out of scope (cross-referenced only for the SAV wiring). The policy
  constraint resolution (U30) and the chain framework (U32) are cross-referenced,
  not re-derived.
- **Lens set:** SEC, PERF-BIGO, PERF-MEM / PERF-LEAK, STD.

### Scope (files actually read)

**`dss/validation/process/bbb/isc` — 5 files, 430 lines — read in full (100%).**

| File | Lines |
|---|---|
| identification_of_the_signing_certificate.go | 167 (full) |
| digest_value_present_check.go | 69 (full) |
| digest_value_match_check.go | 69 (full) |
| issuer_serial_match_check.go | 64 (full) |
| signing_certificate_recognition_check.go | 61 (full) |

**`dss/validation/process/bbb/sav` — 41 files, 3,400 lines — read in full (100%).**

| File | Lines | | File | Lines |
|---|---|---|---|---|
| signature_acceptance_validation.go | 347 (full) | | content_time_stamp_check.go | 64 (full) |
| eaa_acceptance_validation.go | 344 (full) | | signing_time_in_certificate_validity_range_check.go | 63 (full) |
| abstract_acceptance_validation.go | 117 (full) | | x509_url_present_check.go | 61 (full) |
| signature_type_check.go | 106 (full) | | content_identifier_check.go | 61 (full) |
| abstract_time_stamp_present_check.go | 106 (full) | | content_hints_check.go | 61 (full) |
| eaa_revocation_token_acceptance_validation.go | 105 (full) | | claimed_roles_check.go | 61 (full) |
| commitment_type_indications_check.go | 105 (full) | | certified_roles_check.go | 61 (full) |
| content_type_check.go | 104 (full) | | unicity_signing_certificate_attribute_check.go | 60 (full) |
| message_digest_or_signed_properties_check.go | 97 (full) | | tsa_general_name_value_match_check.go | 60 (full) |
| timestamp_acceptance_validation.go | 87 (full) | | tsa_general_name_order_match_check.go | 60 (full) |
| signing_certificate_references_validity_check.go | 86 (full) | | signing_certificate_attribute_present_check.go | 60 (full) |
| content_timestamp_basic_validation_check.go | 82 (full) | | key_identifier_present_check.go | 60 (full) |
| structural_validation_check.go | 77 (full) | | tsa_general_name_field_present_check.go | 59 (full) |
| counter_signature_check.go | 76 (full) | | signing_time_check.go | 59 (full) |
| abstract_time_stamp_type_check.go | 74 (full) | | signer_location_check.go | 59 (full) |
| all_certificates_in_path_referenced_check.go | 72 (full) | | tlevel_time_stamp_check.go | 50 (full) |
| x509_url_match_check.go | 71 (full) | | lta_level_time_stamp_check.go | 50 (full) |
| key_identifier_match_check.go | 65 (full) | | validation_data_time_stamp_check.go | 46 (full) |
| validation_data_refs_only_time_stamp_check.go | 46 (full) | | signature_time_stamp_check.go | 45 (full) |
| document_time_stamp_check.go | 45 (full) | | archive_time_stamp_check.go | 45 (full) |
| revocation_acceptance_validation.go | 43 (full) | | | |

**`dss/validation/process/bbb/vci` — 6 files, 441 lines — read in full (100%).**

| File | Lines |
|---|---|
| validation_context_initialization.go | 115 (full) |
| signature_policy_identifier_check.go | 88 (full) |
| signature_policy_identified_check.go | 60 (full) |
| signature_policy_hash_valid_check.go | 60 (full) |
| signature_policy_zero_hash_check.go | 59 (full) |
| signature_policy_store_check.go | 59 (full) |

**Cross-referenced (not re-derived):** `bbb/abstract_multi_values_check_item.go`
(`ProcessValueCheck`/`ProcessValuesCheck`), `process/chain.go` (`Execute`,
`CollectAllMessages`, `IsValidConclusion`), `process/chain_item.go` (`fail`,
`informOrWarn`), `process/validation_process_utils.go` (`GetCryptoPosition`),
`diagnostic/certificate_wrapper.go` (`Sources()`, `NotBefore()`, `NotAfter()`),
`diagnostic/abstract_token_proxy.go` (`CertificateChain()`),
`diagnostic/signature_wrapper.go` (`IsCounterSignature()`, `Parent()`,
`X509UrlReferences()`), `diagnostic/diagnostic_data.go` (`UsedCertificateByIdNullSafe`),
`diagnostic/jaxb/jaxb_wrappers.go` (`SourcesWrapper.All()`),
`validation/reports/diagnostic/diagnostic_data_builder.go` (`getXmlCertificateSources`),
`validation/reports/diagnostic/signed_document_diagnostic_data_builder.go`
(`attachCounterSignatures`), `enumerations/signature_level.go` (`SignatureForm()`),
`validation/process/eaa/checks/eaa_revocation_present_check.go` (`Process()`),
`validation/process/eaa/checks/acceptable_eaa_revocation_found_check.go` (`Process()`).

### Check verdict verification (dedicated subsection)

Every check's `Process() bool` was traced to the condition it returns `true` for.
The chain framework (`ChainItemBase.fail`) calls `overrides.Process()` and records
`OK` on `true`, `NOT_OK` on `false`; the chain default-concludes `PASSED` when no
item recorded a failure (`chain.go` `Execute`). A check therefore yields a false
"VALID" only if its `Process()` returns `true` on a condition that is actually bad,
or if its `Level()` is empty so it is skipped.

**No always-PASS and no always-FAIL checks, and no inverted comparisons, exist in
any of the three packages.** The decision logic is correct for each:

**ISC** (all gate on the signing-certificate reference + digest — the correct
input for 5.2.3 Identification of the signing certificate):
- `SigningCertificateRecognitionCheck.Process` → `token.SigningCertificate() != nil`
  (a signing certificate was identified ⇒ PASS). Correct.
- `DigestValuePresentCheck.Process` → at least one
  `SigningCertificateReferences` entry has `IsDigestValuePresent()` (a digest
  value is present in some reference ⇒ PASS). Correct (not inverted; the loop
  returns `true` on the first match, `false` if none).
- `DigestValueMatchCheck.Process` → at least one reference has
  `IsDigestValuePresent() && IsDigestValueMatch()` (a digest value is present
  **and** matches ⇒ PASS). Correct.
- `IssuerSerialMatchCheck.Process` →
  `signingCertificateReference != nil && IsIssuerSerialMatch()` (the first
  reference's issuer+serial matches the cert ⇒ PASS; `nil` ref ⇒ FAIL). Correct.

**SAV** (all gate on the signature type + acceptance criteria — the correct
input for 5.2.8 Signature acceptance validation):
- `StructuralValidationCheck.Process` → `IsStructuralValidationValid()`. Correct.
- `SigningCertificateAttributePresentCheck.Process` →
  `IsSigningCertificateReferencePresent()`. Correct.
- `UnicitySigningCertificateAttributeCheck.Process` →
  `IsSigningCertificateReferenceUnique()`. Correct.
- `SigningCertificateReferencesValidityCheck.Process` → no orphan refs **and**
  all related certs are in the chain (map-based O(n) lookup). Correct.
- `AllCertificatesInPathReferencedCheck.Process` → every chain cert is covered
  by a signing-certificate reference (map-based O(n) lookup). Correct.
- `KeyIdentifierPresentCheck.Process` → `KeyIdentifierReference() != nil`.
  Correct.
- `KeyIdentifierMatchCheck.Process` → if ref non-nil, `IsIssuerSerialMatch()`;
  if nil, `true` (vacuously true when no `kid` is present — correct, the check
  is only added to the chain when `KeyIdentifierReference() != nil`). Correct.
- `X509UrlPresentCheck.Process` → `IsCollectionNotEmpty(X509UrlReferences())`.
  Correct.
- `X509UrlMatchCheck.Process` → signing-cert ID is in the x5u-derived cert set.
  Correct.
- `SigningTimeCheck.Process` → `ClaimedSigningTime() != nil`. Correct.
- `SigningTimeInCertificateValidityRangeCheck.Process` →
  `claimedSigningTime != nil && signingCertificate != nil &&
  !claimed.Before(notBefore) && !claimed.After(notAfter)` (claimed time is
  within `[notBefore, notAfter]` ⇒ PASS). Correct (boundary-inclusive, both
  `!Before` and `!After` — matches Java's `!before && !after`).
- `ContentTypeCheck.Process` → `ProcessValuesCheck(contentTypes())` (the
  content-type value is in the constraint set ⇒ PASS). Correct.
- `ContentHintsCheck.Process` → `ProcessValueCheck(contentHints())`. Correct.
- `ContentIdentifierCheck.Process` → `ProcessValueCheck(contentIdentifier())`.
  Correct.
- `MessageDigestOrSignedPropertiesCheck.Process` → switch on form: XAdES →
  SignedProperties digest matcher present; CAdES/PAdES/PKCS7 → MessageDigest
  digest matcher present; JAdES/CB-AdES → `false` (skip — the check is not
  added to the chain for these forms). Correct.
- `CommitmentTypeIndicationsCheck.Process` → CTIs non-empty **and** (if
  expected values present, `containsAll(expected, present)`; if no expected
  values, `true`). Correct (not inverted; `containsAll` is O(n+m) via a
  `map[string]struct{}` — good).
- `SignerLocationCheck.Process` → `IsSignatureProductionPlacePresent()`. Correct.
- `ClaimedRolesCheck.Process` → `ProcessValuesCheck(claimedRoles)`. Correct.
- `CertifiedRolesCheck.Process` → `ProcessValuesCheck(certifiedRoles)`. Correct.
- `ContentTimeStampCheck.Process` → at least one timestamp with
  `Type().IsContentTimestamp()`. Correct.
- `CounterSignatureCheck.Process` → scans all signatures for a
  counter-signature whose `Parent().Id()` equals the current signature's ID.
  Correct (O(n) linear scan, bounded by signature count).
- `SignatureTimeStampCheck` / `ValidationDataTimeStampCheck` /
  `ValidationDataRefsOnlyTimeStampCheck` / `ArchiveTimeStampCheck` /
  `DocumentTimeStampCheck` (all via `AbstractTimeStampTypeCheck`) → at least one
  timestamp with the matching `TimestampType`. Correct.
- `TSAGeneralNameFieldPresentCheck.Process` → `IsTSAGeneralNamePresent()`. Correct.
- `TSAGeneralNameValueMatchCheck.Process` → `IsTSAGeneralNameMatch()`. Correct.
- `TSAGeneralNameOrderMatchCheck.Process` → `IsTSAGeneralNameOrderMatch()`. Correct.
- `SignatureTypeCheck.Process` → `ProcessValuesCheck(signatureTypes())`. Correct.
- `ContentTimestampBasicValidationCheck.Process` →
  `IsValidConclusion(timestampValidationResult)` (delegates to the TST BBB
  conclusion — the correct input). Correct.
- `TLevelTimeStampCheck` / `LTALevelTimeStampCheck` (via
  `AbstractTimeStampPresentCheck`) → at least one timestamp with a valid basic
  validation conclusion **or** a valid PSV conclusion. Correct (not inverted;
  the `IsValidConclusion` guard is nil-safe per `chain.go:246-248`).
- `SignatureAcceptanceValidation.timestampBasicValidationConstraintLevel` →
  if `TimestampValidConstraint` is nil **or** LTA material is present, downgrade
  to WARN (the timestamp check should not hard-fail when LTA is present —
  correct per EN 319 102-1). Correct.
- `SignatureAcceptanceValidation.CollectMessages` → skips TSTBBB-block messages
  when the timestamp-valid constraint is nil or LTA is present (avoids
  double-reporting a warning that is already downgraded). Correct.
- `EAAAcceptanceValidation.InitChain` → the `eaaRevocationPresentCheck.Process()`
  call during `InitChain()` (line 143) is a 1:1 port of Java's `initChain()`
  pattern (Java also calls `process()` as a predicate to decide whether to add
  more chain items). The `Process()` is pure (checks `EAAStatus() != nil`), so
  the double-invocation (once at construction, once at execution) is safe.
  Noted, not a defect (T35-STD-002).
- `EAAAcceptanceValidation.acceptableStatusFound` →
  `NewAcceptableEAARevocationFoundCheck` handles `nil` token: `Process()`
  returns `eaaStatusToken != nil` (nil ⇒ FAIL — correct: no acceptable status
  found). Not inverted.
- `EAAAcceptanceValidation.CollectAdditionalMessages` → the `c.bbbs[...]`
  lookups are guarded by the `InitChain` panic (line 137: `if eaaRevocationBBB
  == nil { panic(...) }`), so by the time `CollectAdditionalMessages` runs, all
  expected IDs are present in the map. Safe.
- `EAARevocationTokenAcceptanceValidation.InitChain` → straightforward chain
  of checks, each gated on the correct token field. Correct.
- `RevocationAcceptanceValidation.InitChain` → single `cryptographic()` step
  (delegates to the AOV result-check, U33 scope). Correct.

**VCI** (all gate on the validation context + policy — the correct input for
5.2.4 Validation context initialization):
- `SignaturePolicyIdentifierCheck.Process` → multi-branch: (1) "NoPolicy" in
  constraint **and** policyId empty ⇒ PASS; (2) "AnyPolicy" in constraint **and**
  policyId non-empty ⇒ PASS; (3) "ImplicitPolicy" in constraint **and** policyId
  equals "IMPLICIT_POLICY" ⇒ PASS; (4) falls through to `ProcessValueCheck`
  (policyId is in the constraint set ⇒ PASS). Correct (not inverted; the
  if-else-if chain is mutually exclusive and covers all three special policy
  types before the OID fallback).
- `SignaturePolicyIdentifiedCheck.Process` → `IsPolicyPresent() &&
  IsPolicyIdentified()`. Correct.
- `SignaturePolicyHashValidCheck.Process` → `IsPolicyPresent() &&
  IsPolicyDigestValid()`. Correct.
- `SignaturePolicyZeroHashCheck.Process` → `IsPolicyZeroHash()`. Correct.
- `SignaturePolicyStoreCheck.Process` → `IsPolicyStorePresent()`. Correct.

**Nil-derefs on hostile input:**
- **T35-SEC-001** (below): `identification_of_the_signing_certificate.go:131-135`
  — `sources[0]` on an empty slice in `AddAdditionalInfo`. The only unguarded
  deref in the three packages.
- All other `Process()` methods are nil-safe: `IsValidConclusion` guards
  `conclusion != nil` (`chain.go:246`); `ProcessValueCheck`/`ProcessValuesCheck`
  delegate to `process.ValueCheck`/`ValuesCheck` which handle empty strings
  and empty constraint sets; the EAA `AcceptableEAARevocationFoundCheck`
  guards `eaaStatusToken != nil`; `SigningTimeInCertificateValidityRangeCheck`
  guards both `claimedSigningTime != nil` and `signingCertificate != nil`
  before dereferencing `NotBefore()`/`NotAfter()`; `X509UrlMatchCheck` guards
  `signingCertificate != nil` before the loop. No other unchecked deref,
  unchecked type assertion, or `interface{}` deref is reachable on the verdict
  path.
- The `SignatureForm()` error → `panic(err)` pattern (4 sites in SAV) is a
  1:1 port of Java's `UnsupportedOperationException` on `SignatureLevelUnknown`.
  The error path is only reachable for a `SignatureLevel` value that is not one
  of the constants, which should not appear in valid diagnostic data. Noted as
  T35-STD-003.

### Findings

**SEC**

T35-SEC-001 — **High** — Category: SEC — dss/validation/process/bbb/isc/identification_of_the_signing_certificate.go:128-135
- Evidence:
  ```go
  // AddAdditionalInfo, lines 128-135
  for _, certificate := range c.token.CertificateChain() {
      chainItem := &jaxb.XmlChainItem{}
      chainItem.Id = certificate.Id()
      sources := certificate.Sources()
      if slices.Contains(sources, enumerations.CertificateSourceTypeTrustedList) {
          chainItem.Source = jaxb.CertificateSourceTypeValue(enumerations.CertificateSourceTypeTrustedList)
      } else if slices.Contains(sources, enumerations.CertificateSourceTypeTrustedStore) {
          chainItem.Source = jaxb.CertificateSourceTypeValue(enumerations.CertificateSourceTypeTrustedStore)
      } else {
          chainItem.Source = jaxb.CertificateSourceTypeValue(sources[0]) // ← panics if sources is empty
      }
  ```
  where `CertificateWrapper.Sources()` returns `nil` when the XML `<Sources>`
  element is absent or empty:
  ```go
  // diagnostic/certificate_wrapper.go:274-283
  func (w *CertificateWrapper) Sources() []enumerations.CertificateSourceType {
      values := w.certificate.Sources.All() // nil when Sources pointer is nil
      if values == nil { return nil }
      ...
  }
  ```
  and `SourcesWrapper.All()` returns `nil` for a nil receiver:
  ```go
  // diagnostic/jaxb/jaxb_wrappers.go:489-494
  func (w *SourcesWrapper) All() []CertificateSourceTypeValue {
      if w == nil { return nil }
      return w.Items // may be empty (len 0)
  }
  ```
  In the normal diagnostic-data-builder path, `getXmlCertificateSources`
  always appends `CertificateSourceTypeUnknown` when the list is empty
  (`diagnostic_data_builder.go:2419-2421`), so the normal flow is safe. But the
  `TokenProxy` interface is public — a user-supplied or externally-deserialized
  diagnostic-data XML with a chain certificate lacking `<Sources>` (or with an
  empty `<Sources>` element) produces an empty `sources` slice, and
  `sources[0]` panics with "index out of range [0] with length 0".
  `diagnostic/diagnostic_data.go:413` (`UsedCertificateByIdNullSafe`) returns
  `NewCertificateWrapper(&jaxb.XmlCertificate{})` — a bare certificate with nil
  `Sources` — confirming the type can exist in the diagnostic data.
- Impact: `AddAdditionalInfo` is called at the end of `Execute()`
  (`chain.go:197`), **after** all checks have run and the verdict has been
  recorded. The panic does not corrupt the verdict, but it prevents the report
  from being completed:
  - Through the public facade (`dss.Validate`), the panic is caught by
    `format.go:329 recovered()` and surfaced as a hard error — the user gets an
    error where a report was expected.
  - Through the lower-level packages (a caller that uses `validation/…`
    without the facade), it is an uncaught panic → process crash.
  In Java, `sources.get(0)` on an empty list throws `IndexOutOfBoundsException`,
  caught by the same facade recovery. The Go port is a faithful 1:1 port of the
  Java behaviour; the panic is **not** a false PASS or false FAIL — the verdict
  is already computed — but it is a **panic on hostile/malformed input** that
  aborts report generation.
- Recommendation: Guard the `else` branch: `if len(sources) > 0 { chainItem.Source = ... }
  else { /* leave Source unset, or use CertificateSourceTypeUnknown */ }`. This
  matches the builder's own fallback (`getXmlCertificateSources` appends
  `UNKNOWN` when empty) and prevents the index-out-of-range panic on
  externally-supplied diagnostic data. Add a KAT test: a diagnostic-data XML
  with a chain certificate lacking `<Sources>` must not panic and must produce
  a well-formed report (with the certificate's `Source` left empty or set to
  `UNKNOWN`).

**PERF**

No PERF-BIGO, PERF-MEM, or PERF-LEAK findings in this unit. All checks are
O(n) or O(1) on bounded inputs. `CommitmentTypeIndicationsCheck.containsAll`
is O(n+m) via a `map[string]struct{}`. `SigningCertificateReferencesValidityCheck`
and `AllCertificatesInPathReferencedCheck` both use map-based lookups
(O(n) construction, O(1) per lookup). `CounterSignatureCheck` is a single
O(n) linear scan over the signature list. No unbounded growth, no retained
references outliving the result object, no quadratic allocation on
attacker-influenced input.

**STD**

T35-STD-001 — **Info** — Category: STD — dss/validation/process/bbb/sav/eaa_acceptance_validation.go:137-138
- Evidence:
  ```go
  eaaRevocationBBB := c.bbbs[eaaRevocationWrapper.Id()]
  if eaaRevocationBBB == nil {
      panic("No BasicBuildingBlock found for token with Id '" + eaaRevocationWrapper.Id() + "'")
  }
  ```
- Impact: Deliberate — the comment context and the `panic` message name a
  diagnostic-data inconsistency (a revocation token present without a matching
  BBB block), not a validation outcome. This is a faithful 1:1 port of Java's
  `IllegalStateException` (or equivalent) reporting a programming/data error,
  not hostile document input. The `InitChain` loop is the guard: if any
  revocation token's BBB is missing, the chain construction panics before any
  check runs, which is correct (the chain cannot be built without the BBB
  conclusion). Noted, not a defect. (Same pattern as T33-STD-002.)
- Recommendation: None — keep.

T35-STD-002 — **Info** — Category: STD — dss/validation/process/bbb/sav/eaa_acceptance_validation.go:143
- Evidence:
  ```go
  eaaRevocationPresentCheck := c.statusPresent()
  item = item.SetNextItem(eaaRevocationPresentCheck)
  if eaaRevocationPresentCheck.Process() { // ← called during InitChain
      item = item.SetNextItem(c.statusAvailable())
      ...
  }
  ```
- Impact: `Process()` is called once during `InitChain()` (as a predicate to
  decide whether to add more chain items) and again during `Execute()` (when
  the chain runs). The `Process()` of `EAARevocationPresentCheck` is pure
  (checks `EAAStatus() != nil`), so the double-invocation is safe and
  idempotent. This is a 1:1 port of Java's `initChain()` pattern, where Java
  also calls `process()` as a predicate during chain construction. Noted,
  not a defect.
- Recommendation: None — keep (faithful port). If the Go codebase ever
  introduces side-effects into `Process()`, this pattern would need revisiting.

T35-STD-003 — **Info** — Category: STD — dss/validation/process/bbb/sav/signature_acceptance_validation.go:76 (also in `signature_type_check.go:68`, `content_type_check.go:71`, `message_digest_or_signed_properties_check.go:56`)
- Evidence:
  ```go
  signatureForm, err := c.token.SignatureFormat().SignatureForm()
  if err != nil {
      panic(err) // "the signature level 'UNKNOWN' is not supported"
  }
  ```
- Impact: `SignatureForm()` returns an error only for `SignatureLevelUnknown`
  (`enumerations/signature_level.go:210-214`), which has no associated form.
  In the normal diagnostic-data flow, the signature level is always a known
  value (set by the builder based on the actual signature format), so the error
  path is practically unreachable. This is a 1:1 port of Java's
  `UnsupportedOperationException` thrown from `getSignatureForm()` on UNKNOWN.
  Through the facade, the panic is caught; through lower-level packages, it is
  an uncaught panic. Noted, not a defect in the normal flow.
- Recommendation: None — keep (faithful port). If the diagnostic-data builder
  ever produces an UNKNOWN level, the panic would surface as an error through
  the facade, which is acceptable.

### Tool log (run from inside `dss/`)

```
$ gofmt -l validation/process/bbb/isc validation/process/bbb/sav validation/process/bbb/vci
(no output — clean)

$ go vet ./validation/process/bbb/isc/... ./validation/process/bbb/sav/... ./validation/process/bbb/vci/...
(no output, exit 0)

$ go vet -tags eaa ./validation/process/bbb/isc/... ./validation/process/bbb/sav/... ./validation/process/bbb/vci/...
(no output, exit 0) — the eaa-tagged files (eaa_acceptance_validation.go,
    eaa_revocation_token_acceptance_validation.go) also pass vet

$ golangci-lint run --config=../.github/.golangci.yml \
    ./validation/process/bbb/isc/... ./validation/process/bbb/sav/... ./validation/process/bbb/vci/...
0 issues.
```

All four tool reports are clean → no tool-driven findings. (The three
`//nolint:staticcheck` comments in `identification_of_the_signing_certificate.go:113`,
`signature_acceptance_validation.go:232`, and
`validation_context_initialization.go:87-90` mark deliberate dead-store ports of
Java's trailing `item = item.setNextItem(...)` and are correctly suppressed.)

### Summary

- Total findings: **4**
- Critical: **0** — High: **1** — Medium: **0** — Low: **0** — Info: **3**
- SEC: 1 (1 High) · PERF-BIGO: 0 · PERF-MEM/LEAK: 0 · STD: 3 (3 Info)
- Check-verdict verdict: **no always-PASS, no always-FAIL, no inverted
  comparisons** in any of the three packages. All ISC checks gate on the
  correct input (signing-certificate reference + digest), all SAV checks gate
  on the correct input (signature type + acceptance criteria + token fields),
  and all VCI checks gate on the correct input (validation context + policy).
  The single nil-deref (T35-SEC-001) is in `AddAdditionalInfo` (report-
  generation time, after the verdict is computed) — it is a panic/DoS on
  hostile input, not a false PASS or false FAIL verdict. The EAA SAV wiring
  is correct; the `Process()` double-invocation and the BBB-missing panic are
  deliberate 1:1 ports of Java patterns (T35-STD-001, T35-STD-002). The
  `SignatureForm()` panic (T35-STD-003) is a faithful port of Java's
  `UnsupportedOperationException` on UNKNOWN level, practically unreachable in
  the normal flow.
- Open questions for consolidation:
  1. **T35-SEC-001 upstream parity:** confirm that Java DSS 6.5.RC1's
     `IdentificationOfTheSigningCertificate.addAdditionalInfo()` also throws
     `IndexOutOfBoundsException` on a chain certificate with empty `Sources`
     (it should, since `sources.get(0)` on an empty `ArrayList` throws). If so,
     the Go port is faithful and the fix is a hardening improvement (guard the
     `else` branch); if Java handles it differently, the Go port diverges and
     should be aligned.
  2. **T35-STD-002 (EAA `Process()` double-invocation):** confirm that the
     EAA checks package (`validation/process/eaa/checks`) has no side-effects
     in any `Process()` method that would be triggered twice (once in
     `InitChain`, once in `Execute`). If any future EAA check introduces
     stateful `Process()`, the `InitChain` predicate pattern would need
     revisiting. (Low risk today — all current EAA `Process()` methods are
     pure field-reads.)
  3. **EAA SAV cross-reference (T32-SEC-001):** the EAA nil-panic filed in
     batch-10 (`eaa_validation_process.go:84-108`, `item.SetNextItem` on a nil
     `item` when `eaaBBBs.FC` is nil) is in the `eaa` package, not in
     `bbb/sav`. The `eaa_acceptance_validation.go` in this unit is a separate
     file (the SAV wiring for EAA tokens) and does not have the same nil-
     `item` issue — its `InitChain` always starts from a concrete
     `c.etsi194721Conformance()` item. No cross-reference needed beyond
     noting the EAA SAV wiring is correct and independent of T32-SEC-001.
