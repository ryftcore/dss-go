# Batch 04 — ASiC core (U14a)

- **Unit:** U14a
- **Scope:** `dss/asic` core — ASiC container handling, custom ZIP reader, and zip-bomb guards.
  Focus files: `secure_container_handler.go`, `asic_utils.go`,
  `abstract_asic_container_analyzer.go`, `default_container_merger.go`,
  `abstract_asic_container_evidence_record_builder.go`, `dss_zip_entry.go`, `asic_content.go`.
- **Date:** 2026-08-23
- **Depth:** deep-budgeted (hard read budget ~10k lines incl. grep; actual ≈ 5.0k read + signature greps)
- **Context:** Go port of Java DSS 6.5.RC1. `dss/PORTING.md` read first — 1:1 Java port,
  staticcheck style rules disabled, SA1019 allowed, byte-identical enum values, stdlib-first,
  no cgo. `docs/compatibility/known-gaps.md` read — listed gaps not re-reported (notably:
  evidence-record analyzers are *not* registered, so ER parsing is a known gap, not a finding).
  `cades` (C04A/B), `cmscore`, `internal/*` already reviewed — not re-reported.

## Files read

| File | Coverage |
|---|---|
| `secure_container_handler.go` (1164) | **FULL** (1–1164) + all 44 `func` signatures grepped. Covers the custom `zipLocalHeaderReader`/`zipEntryStream` engine, `SecureCopy`/`SecureSkip` bomb guards, `BuildZip` write path. |
| `asic_utils.go` (972) | **FULL** (1–972) + all 60 `func` signatures grepped. Covers zip-comment scan, mimetype/container-type detection, manifest-name predicates. |
| `dss_zip_entry.go` (304) | **FULL** (1–304). Extra-field (NTFS/EXTT/ZIP64) time decoding. |
| `asic_content.go` (274) | **FULL** (1–274). Document grouping + mutation idiom. |
| `abstract_asic_container_evidence_record_builder.go` (389) | **FULL** (1–389). ER incorporation + digest-algorithm/coverage asserts. |
| `default_container_merger.go` (444) | **FULL** (1–444). Merge + `CreateZipArchiveAt` write path. |
| `abstract_asic_container_analyzer.go` (506) | **BUDGETED** — lines 100–506 read; lines 1–99 (type + `Overrides` interface decls) **signature-only** via the 37 `func` grep. All functions covered via read-or-signature. |
| `secure_container_handler_thresholds_test.go` | **SKIMMED** (1–200) for coverage only — confirms the four guards are load-bearing (mutation-tested via setter *and* builder). |
| Supporting (for cross-check) | `file_archive_entry.go` 1–200 (lazy read path), `default_asic_container_extractor.go` 1–160, `asic_manifest_parser.go` 1–120, `utils/io.go` 40–120, `spi/dss_utils.go` 830–870, `PORTING.md` 1–120, `known-gaps.md` 1–60. |

## Security posture (top priority) — VERIFIED CLEAN on the attack surface

The zip-bomb guards and the custom ZIP reader were the primary target. Findings below are the
residue; the headline is that **no exploitable defect was found** in the asic core read/merge/verify
paths. Specifically verified:

- **Bomb guards enforced on ALL read paths — no bypass.** `extractZipEntries` (→`SecureSkip`) and
  `getCurrentEntryDocument` (→`SecureCopy`) both flow through `assertExtractEntryLengthValid`
  (`byteCounter > threshold && byteCounter > allowedSize`, `allowedSize = containerSize * ratio`)
  and `assertCollectionSizeValid` (`size > maxAllowedFilesAmount`). `getNextValidEntry` is bounded
  by `maxMalformedFiles`. The **lazy** `FileArchiveEntry.OpenStream` read is *not* independently
  guarded, but it is only reachable after a container already walked+vetted through the guarded
  path, so it is **not** a bypass (it re-reads bytes already counted). `BuildZip`'s write path uses
  `SecureCopy(…, -1)` (guard disabled) but its inputs are caller-supplied documents, not a hostile
  ZIP — correct.
- **No path traversal / no symlink following.** The core never extracts to disk by entry name;
  entry names stay in-memory `DSSDocument` names. `FileArchiveEntry.Save(filePath)` takes a
  *caller-supplied* absolute path, not an entry name — not an entry-name traversal sink.
- **No OOM/panic on hostile ZIP input.** Every structural read (`readAt`) bounds-checks
  `off < 0 || off+len > r.size`. ZIP64 sizes are resolved from the extra field (bounds-checked).
  `int64` size conversions plus `io.NewSectionReader(…, n ≤ reader.size)` keep decompression
  bounded by the *real* file size, so a forged `uncompressedSize` cannot drive allocation.
  `flate` errors are returned, not panicked. `time.Date` (in `zipDosToTime`) normalizes, not
  panics, out-of-range MS-DOS fields. `getNextValidEntry` cannot loop forever (bounded by
  `maxMalformedFiles`). Truncated headers / malformed central directory → clean error or skipped
  entry, never a panic.
- **CRC enforced, not skippable.** STORED entries check CRC at EOF in `Read`; DEFLATED entries
  check it in `finish()` (after the data descriptor is consumed). The `archive/zip` (central-dir)
  path also validates. Note the Go reader is *stricter* than Java: it validates CRC/size while
  `SecureSkip` walks, where Java's `skip` does not — a deliberate hardening consistent with
  PORTING.md's "refuse what upstream only warns about" (allowed, not a divergence defect).
- **No XML entity-expansion in the manifest parser.** `asic_manifest_parser.go` delegates all
  DOM building to `internal/xmldom` (reviewed clean) and does no custom entity handling.
- **Merger preserves bytes.** `Merge`→`CreateZipArchiveAt`→`BuildZip` streams each entry through
  `SecureCopy` into the `zip.Writer`; it does not re-encode or corrupt. Same-name conflicts with
  differing SHA-256 are rejected (integrity check), matching upstream.
- **Evidence-record builder** validates digest-algorithm agreement and that every signed document
  is covered (via `referenceValidation.IsIntact()`); the actual digest computation is delegated to
  the (known-gap, unregistered) ER analyzers — not re-reported as a finding.

## Findings

### A14A-SEC-001
- **Severity:** Low
- **Category:** SEC
- **Location:** `dss/asic/asic_utils.go:854` (`UtilsZipCommentFromArchiveContainer`)
- **Evidence:**
  ```go
  for ii := length - 22; ii >= 0; ii-- {
      isMagicStart := true
      for jj := 0; jj < len(asicUtilsMagicDir); jj++ {   // asicUtilsMagicDir = {0x50,0x4b,0x05,0x06}
          if buffer[ii+jj] != asicUtilsMagicDir[jj] { isMagicStart = false; break }
      }
      if isMagicStart { realLen := length - ii - 22; ... return string(buffer[ii+22 : ii+22+realLen]) ... }
  }
  ```
- **Impact:** The ZIP comment is recovered by a backward scan for the EOCD magic `PK\x05\x06`. A
  crafted container whose *comment* contains those 4 bytes (and that omits/neutralises the `mimetype`
  document) can cause the scan to latch at the wrong offset, yielding a truncated or wrong comment.
  The comment feeds `asicUtilsContainerTypeFromZipComment` → container-type routing (ASiC-S vs ASiC-E).
  Impact is limited to mis-routing the analyzer, not a false-accept of a signature, because the
  `mimetype` document is the primary (step-1) type signal and the manifest/signature digests are
  still verified. Faithful to upstream (Java performs the same backward scan), so not a divergence.
- **Recommendation:** Hardening only — anchor the scan to the last EOCD signature (scan strictly
  backward from the end and stop at the *first* match) or validate that the found `realLen` equals
  the EOCD's own `commentLen` field before accepting, so a comment-embedded magic cannot win.
  Low priority given limited impact and parity with Java.

### A14A-PERF-001
- **Severity:** Medium
- **Category:** PERF-BIGO
- **Location:** `dss/asic/asic_utils.go:401` (`UtilsIsASiC`), `:552` (`UtilsContainerType`), `default_asic_container_extractor.go:97` (`Extract`)
- **Evidence:**
  ```go
  // UtilsIsASiC        -> ZipUtilsInstance().ExtractEntryNames(doc)        // full walk + SecureSkip (inflates all DEFLATED)
  // UtilsContainerType -> ZipUtilsInstance().ExtractEntryNames(archiveContainer) // second walk
  //                       + asicUtilsMimetypeDocument -> ExtractContainerContent // third walk
  // DefaultContainerExtractor.Extract -> ExtractContainerContent          // (re)walk
  ```
- **Impact:** A single ASiC open/validation can walk the container 2–3×, and each walk
  `SecureSkip`s (i.e. fully inflates) every DEFLATED entry to reach the next header. That is O(k·n)
  CPU for k walks over n inflated bytes, re-doing work already performed. Faithful to upstream
  (Java's `ASiCUtils.isASiC`/`getContainerType` make the same repeated `getEntryNames` calls), so
  not a divergence — but a real, avoidable performance cost on large containers.
- **Recommendation:** Cache one validated entry-inventory (names + sizes + vetting result) per
  container document and reuse it across `UtilsIsASiC`/`UtilsContainerType`/`Extract`, so the ZIP
  is walked and inflated once. Keep the per-call bomb-guard semantics intact (reset counters per
  walk).

### A14A-PERF-002
- **Severity:** Medium
- **Category:** PERF-MEM
- **Location:** `dss/asic/default_container_merger.go:214` (`Merge`→`CreateZipArchiveAt`), `secure_container_handler.go:415` (`CreateZipArchive`→`InstantiateResourcesHandler`)
- **Evidence:**
  ```go
  func (h *SecureContainerHandler) InstantiateResourcesHandler() (resources.DSSResourcesHandler, error) {
      return h.resourcesHandlerBuilder.CreateResourcesHandler(), nil   // default = InMemoryResourcesHandler
  }
  // CreateZipArchive: zos := zip.NewWriter(outputStream); ... dssResourcesHandler.WriteToDSSDocument()
  ```
- **Impact:** The merged (and created) output container is fully materialized in memory by the
  default `InMemoryResourcesHandler`. Merging several large ASiC containers therefore holds the
  entire concatenated output in RAM. Bounded by input size and configurable (a file-backed
  `DSSResourcesHandlerBuilder` via `SetResourcesHandlerBuilder` avoids it), so not a leak — but a
  memory ceiling for large merges under the default configuration.
- **Recommendation:** Document the default in-memory output and, for large-merge use cases, steer
  callers to a file-backed `DSSResourcesHandlerBuilder`. No code change required for parity.

### A14A-STD-001
- **Severity:** Low
- **Category:** STD
- **Location:** `dss/asic/secure_container_handler.go:127` (`byteCounter`, `malformedFilesCounter`), `:617` (`resetCounters`), `:694`/`:719` (`SecureCopy`/`SecureSkip`)
- **Evidence:**
  ```go
  type SecureContainerHandler struct {
      ...
      byteCounter int64            // NOTE: shall be reset on every use
      malformedFilesCounter int    // NOTE: shall be reset on every use
  }
  func (h *SecureContainerHandler) SecureCopy(is io.Reader, os io.Writer, allowedSize int64) error {
      ... h.byteCounter += int64(nRead) ...
  ```
- **Impact:** The two guard counters are shared mutable instance state with no synchronization and a
  "reset on every use" contract enforced only by convention. If a single handler instance is shared
  across goroutines, `SecureCopy`/`SecureSkip`/`resetCounters` race (undefined behavior, would trip
  `-race`), and interleaved resets could in principle weaken the bomb guard. Not reachable via the
  library's own per-operation builder usage, and faithful to Java (instance fields, same convention)
  — so a robustness/idiom note rather than a defect.
- **Recommendation:** Either make the handler explicitly non-shareable (doc + `//go:build`-free
  guard), or move the counters onto the per-call walk (e.g. a small `walkState` passed through
  `extractZipEntries`/`getNextDocument`) so there is no shared mutable state to race.

### A14A-PERF-003
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `dss/asic/secure_container_handler.go:391` (`extractCommentsInto`), `:288` (`isInFileProcessingSupported`)
- **Evidence:**
  ```go
  for _, zipEntry := range zipEntries {
      for _, zipFileEntry := range readCloser.File {          // O(n*m) name match
          if zipFileEntry.Name == zipEntry.Name { zipEntry.Comment = zipFileEntry.Comment; break }
      }
  }
  // isInFileProcessingSupported: zip.OpenReader(fileDocument.Path())  // #1 open
  // extractCommentsInto:         zip.OpenReader(fileDocument.Path())  // #2 open
  ```
- **Impact:** When `extractComments` is enabled (default **off**), entry comments are matched by a
  nested O(n·m) name scan (n,m ≤ `maxAllowedFilesAmount` = 1000) and the archive is opened a second
  time via `archive/zip`. Bounded and off by default, so low impact; noted for completeness.
- **Recommendation:** Index `readCloser.File` by name into a `map[string]*zip.FileHeader` for O(1)
  comment lookup; optionally reuse a single `zip.OpenReader` handle. Low priority.

## Lens summaries

- **STD:** One Low (A14A-STD-001, shared mutable guard state). Otherwise idiomatic; doc comments are
  thorough and consistently note the Java provenance and deliberate divergences. Error handling
  returns `error`/`panic` per PORTING.md; no leaky deferred-close patterns in the read path.
- **PERF:** Two Medium (A14A-PERF-001 repeated walk/inflate, A14A-PERF-002 in-memory output) and one
  Low (A14A-PERF-003). The streaming design is sound (no full-container buffering of *inputs*; the
  merger streams entries); the costs above are inherent to the faithful, self-contained utility
  decomposition and the default in-memory resources handler.
- **SEC:** No Critical/High. One Low hardening note (A14A-SEC-001). The core read/merge/verify
  surface is well-hardened: bomb guards on all read paths, bounded decompression, enforced CRC, no
  traversal/symlink/entity-expansion, and no panic on hostile ZIP input.

## Tool log

Run from `dss/` (Go module root):

| Tool | Command | Result |
|---|---|---|
| gofmt | `gofmt -l asic` | clean (no output) |
| go vet | `go vet ./asic/...` | clean (no output) |
| golangci-lint | `golangci-lint run --config=../.github/.golangci.yml ./asic/...` | `0 issues.` (exit 0) |

All three tools available and clean. (No tools unavailable.)

## Open questions

1. **Concurrent handler use.** Is a single `SecureContainerHandler` instance ever shared across
   goroutines in the facade/CLI? If not, A14A-STD-001 stays a convention note; if yes, it should be
   re-rated to High (data race on the guard counters). Worth confirming against the `dss` facade
   and `cmd/esig` call sites (out of this unit's scope).
2. **ER analyzer wiring.** The evidence-record builder's digest-intact checks depend on
   `analyzer.EvidenceRecordAnalyzerFromDocument`, which per `known-gaps.md` has no registered
   factory today (so `getEvidenceRecord` panics). Once ER analyzers are ported, re-verify
   `assertEvidenceRecordValid`/`assertSignedDataCovered` against a real parsed ER.
3. **A14A-PERF-001 caching.** Confirm the 2–3× walk count against a concrete end-to-end
   sign/validate trace before committing to a shared entry-inventory cache, so the cache invalidation
   semantics (and per-call bomb-guard resets) stay correct.

---

## asic/cades (unit U15)

- **Unit:** U15
- **Scope:** `dss/asic/cades` — ASiC-C / ASiC-S / ASiC-E with CAdES: container merger,
  analyzer/validator, signature extension (LTA), filename factories, data-to-sign helpers,
  manifest builders, evidence-record builder, timestamp service/analyzer, and the
  `testdata/crossgen` generator + `extension/*`.
- **Date:** 2026-08-25
- **Depth:** deep-budgeted (hard read budget ~10k lines incl. grep; actual ≈ 6.1k read +
  signature greps for the one >400-line file). 43 non-test Go files (5,574 lines) + 2
  `extension/` files + `testdata/crossgen/main.go` (220).
- **Context:** Go port of Java DSS 6.5.RC1. `dss/PORTING.md` read first — 1:1 Java port,
  staticcheck style rules disabled, SA1019 allowed, byte-identical enum values, stdlib-first,
  no cgo. `docs/compatibility/known-gaps.md` read — listed gaps not re-reported. **Already
  reviewed (not re-reported):** asic core (A14A-*), `cades` core (C04A/B-*), `cmscore`
  (U02), `internal/*`. This unit covers only the `asic/cades` orchestration layer; the
  security-critical primitives it calls into (CMS parse/re-serialize, level checks, manifest
  digesting, ZIP read/merge) live in those already-reviewed packages and are not re-assessed.

## Files read

All 43 non-test files listed by `find asic/cades -name '*.go' -not -name '*_test.go' | xargs wc -l | sort -rn`,
plus `extension/*` and `testdata/crossgen/main.go`. 42 files read **FULL**; the single >400-line
file read **BUDGETED** (first 250 lines + all 36 `func` signatures grepped + every security-relevant
region read in full). Every function is covered via read-or-signature.

| File (lines) | Coverage |
|---|---|
| `asic_with_cades_service.go` (740) | **BUDGETED** — 1–250 + 227–740 (sign/extend/counter-sign/evidence-record/policy-store/adapter) read; all 36 `func` signatures grepped. All functions covered via read-or-signature. |
| `asice_with_cades_container_merger.go` (393) | **FULL** — ASiC-E merge guards, manifest/ER dedup, CMS merge orchestration. |
| `abstract_asic_with_cades_container_merger.go` (232) | **FULL** — `MergeCmsSignatures` (delegates to reviewed `cms` layer) + cert/CRL/OCSP store dedup. |
| `asics_with_cades_container_merger.go` (228) | **FULL** — ASiC-S merge guards (single-sig/ER rules, root-signer-name check). |
| `asic_container_with_cades_analyzer.go` (349) | **FULL** — CMS/timestamp analyzer wiring, manifest validation, SPP-propagation gap. |
| `simple_asic_with_cades_filename_factory.go` (331) | **FULL** — name templates + EN 319 162-1 conformance panics. |
| `asic_with_cades_signature_extension.go` (250) | **FULL** — LT/LTA extension dispatch, `ExtensionRequired`/`AssertExtend*`/`IsCoveredBy*`. |
| `asic_with_cades_level_baseline_lta.go` (249) | **FULL** — ASiC-E LTA: archive-manifest build + timestamp + last-TST re-extension. |
| `asic_with_cades_container_merger_factory.go` (140) | **FULL** — ASiC-S/E selection, type-mix reject. |
| `asic_with_cades_format_detector.go` (111) | **FULL** — CAdES vs XAdES vs OpenDocument routing. |
| `asic_with_cades_timestamp_analyzer.go` (130) | **FULL** — scope/token/evidence-record-coverage. |
| `asic_with_cades_data_to_sign_helper_builder.go` (123) | **FULL** — from-archive/from-files manifest dispatch. |
| `asic_with_cades_timestamp_validator.go` (116) | **FULL** — leaf validator embedding note (non-promotable analyzer). |
| `asic_container_with_cades_validator.go` (100) | **FULL** — validator leaf + naming note. |
| `asic_with_cades_timestamp_parameters.go` (91) | **FULL** — `zipCreationDate` default `time.Now()`. |
| `asic_with_cades_signature_parameters.go` (90) | **FULL** — CAdES-form level assert, `ZipCreationDate`. |
| `extension/asic_with_cades_document_extender.go` (89) | **FULL** — `document` augmentation extender. |
| `asic_with_cades_timestamp_service.go` (85) | **FULL** — TSA `TimeStampResponse` + manifest attach. |
| `asic_with_cades_signature_extension_helper.go` (83) | **FULL** — counter-sign helper, `GetDocumentAnalyzer` panic. |
| `asice_with_cades_archive_manifest_builder.go` (76) | **FULL** — root-file ref-equality, TST SigReference. |
| `asice_with_cades_manifest_builder.go` (69) | **FULL** — signed-documents-only filter. |
| `asic_with_cades_container_extractor.go` (64) | **FULL** — per-entry-type allow predicates (delegated to reviewed `asic` layer). |
| `data_to_sign_asics_with_cades_from_archive.go` (55) | **FULL** — single-embedded-sig select. |
| `asic_with_cades_container_evidence_record_builder.go` (55) | **FULL** — ER filename convention assert. |
| `asic_with_cades_diagnostic_data_builder.go` (54) | **FULL** — ATSHI XML build. |
| `asic_with_cades_timestamp_manifest_builder.go` (48) | **FULL** — TST SigReference leaf. |
| `asic_with_cades_signature_manifest_builder.go` (48) | **FULL** — PKCS7 SigReference leaf. |
| `asic_container_with_cades_validator_factory.go` (48) | **FULL** — factory + `init()` registration. |
| `asic_container_with_cades_analyzer_factory.go` (47) | **FULL** — factory + `init()` registration. |
| `asic_with_cades_lta_validation_phase8.go` (46) | **FULL** — re-validate + exclude-dup validation-data. |
| `asic_with_cades_container_extractor_factory.go` (46) | **FULL** — factory + `init()` registration. |
| `abstract_asic_with_cades_data_to_sign_helper_builder.go` (46) | **FULL** — `IsASiCArchive` + `GetDataPackageName`. |
| `asic_with_cades_timestamp_data_to_sign_helper_builder.go` (42) | **FULL** — timestamp manifest builder dispatch. |
| `extension/asic_with_cades_document_extender_factory.go` (40) | **FULL** — `init()` registration. |
| `data_to_sign_asice_with_cades_helper.go` (40) | **FULL** — cached `ToBeSigned`, empty detached. |
| `asic_with_cades_utils.go` (40) | **FULL** — `getSignedDocument` (ASiC-S vs ASiC-E). |
| `asic_with_cades_signature_data_to_sign_helper_builder.go` (39) | **FULL** — signature manifest builder dispatch. |
| `data_to_sign_asics_with_cades_from_files.go` (34) | **FULL** — `SignedDocuments()[0]` select. |
| `asic_with_cades_asic_content_builder.go` (34) | **FULL** — content builder leaf. |
| `get_data_to_sign_asic_with_cades_helper.go` (30) | **FULL** — `GetDataToSignASiCWithCAdESHelper` interface. |
| `asic_with_cades_filename_factory.go` (30) | **FULL** — `ASiCWithCAdESFilenameFactory` interface. |
| `asic_with_cades_service_lta_phase8.go` (26) | **FULL** — LTA phase-8 split-out (re-validate path). |
| `asic_with_cades_common_parameters.go` (25) | **FULL** — `ASiCWithCAdESCommonParameters` interface. |
| `abstract_get_data_to_sign_asics_with_cades.go` (20) | **FULL** — shared ASiC-S helper base. |
| `testdata/crossgen/main.go` (220) | **FULL** (test-only generator — see A15-SEC-002). |
| Test files (4, 972 lines) | **SKIMMED** for coverage: `manifest_kat_test.go` (byte-exact vs Java oracle), `asic_with_cades_signature_test.go` (helper/dispatch/override KATs), `asic_downstream_cross_validation_test.go` (live-Java round trip), `asic_container_with_cades_validator_smoke_test.go`. |

## Security posture — VERIFIED CLEAN on the asic/cades attack surface

The security-critical work (CMS parse/re-serialize, level-evidence checks, manifest digesting,
ZIP read/merge) is delegated to already-reviewed packages (`cms`, `cmscore`, `cades`, `internal/*`,
asic core). Verified in this layer specifically:

- **CMS merge does not corrupt or false-reject.** `MergeCmsSignatures`
  (`abstract_asic_with_cades_container_merger.go:66`) orchestrates the merge but every
  re-serialization call (`cms.UtilsReplaceSigners`, `UtilsReplaceCertificatesAndCRLs`,
  `UtilsPopulateDigestAlgorithmSet`, `UtilsWriteToDSSDocument`) is in the already-reviewed `cms`
  layer. The merger only gathers SignerInfos/certs/CRLs/OCSP and dedups (order-preserving
  `appendUniqueBytes`); the byte-exact CMS output is the reviewed layer's contract, pinned by
  `manifest_kat_test.go` (byte-exact vs Java oracle) and the live-Java round trip.
- **Level checks (B-LT/LTA) require correct evidence.** `SignDocumentMultiple`
  (`asic_with_cades_service.go:163`) downgrades to LT *before* signing an ASiC-E LTA signature
  (the archive timestamp is a container-level artifact, not CMS-embedded), then restores LTA and
  runs the LTA extension (archive manifest + TSA token) — matching Java and the ASiC spec.
  `ASiCWithCAdESLevelBaselineLTA.Extend` (`asic_with_cades_level_baseline_lta.go:52`) builds the
  archive manifest, requests a real `TimeStampResponse`, and (when extending) re-validates and
  re-embeds revocation data via `extendLastArchiveTimestampWithValidationData`
  (`asic_with_cades_lta_validation_phase8.go`). `AssertExtendSignaturePossible`/
  `AssertExtendTimestampPossible` correctly *refuse* to extend a signature/timestamp already
  covered by a manifest (no double-covering).
- **Signature extension preserves the existing signature.** `extendSignatureDocument`
  (`asic_with_cades_signature_extension.go:143`) captures the detached contents *and* the
  original filename before delegating to the reviewed `cades.SignatureExtender`, and restores the
  name afterwards — the CMS is augmented in place (counter-signature), not rebuilt, so the existing
  signature value is intact.
- **Analyzer extraction is correct (no false-accept/reject).** `GetSignatureAnalyzers`
  (`asic_container_with_cades_analyzer.go:93`) wires each CMS to its *linked* manifest
  (`ManifestParserGetLinkedManifest`) and the signed document (`ASiCWithCAdESUtilsGetSignedDocument`
  — ASiC-S single root / ASiC-E linked-manifest). `OriginalDocumentsForSignature` returns exactly the
  manifest-covered documents. Timestamp scoping (`ASiCWithCAdESTimestampScopeFinder`) only accepts a
  manifest entry when `IsIntact()` and the URI matches a real container document.
- **Filename factory collisions are rejected, not silently accepted.** `SimpleASiCWithCAdESFilenameFactory`
  asserts the EN 319 162-1 template (`META-INF/signature*.p7s`, `timestamp*.tst`, etc.) and panics on
  a name already present (`AssertFilenameValid`) or on the reserved default archive-manifest name.
  Predictable names are a documented, expected property of the ASiC spec (not a vuln).
- **No panic on hostile input in the prod path** beyond Java-parity `IllegalInputException`
  panics (container type, missing manifest, empty lists) — all faithful to upstream's
  `Objects.requireNonNull`/`IllegalInputException`/`UnsupportedOperationException` contract, which
  the facade converts to errors (see `format.go` `recovered()`).
- **No network.** The only `os/exec` reference in the package is in a test
  (`asic_downstream_cross_validation_test.go`) that drives the live-Java oracle. No
  `http.Get`/`Dial`/`url.Parse` in any non-test file.
- **`AddSignaturePolicyStore` / counter-sign / evidence-record paths** all validate the container
  is ASiC, assert the target signature exists, and refuse to add to a signature already covered by
  a manifest (`AssertAddSignaturePolicyStorePossible`).

## Findings

### A15-SEC-001
- **Severity:** Info
- **Category:** SEC
- **Location:** `dss/asic/cades/asic_container_with_cades_analyzer.go:15` (header note), `:93` (`GetSignatureAnalyzers`)
- **Evidence:**
  ```go
  // ...so this propagation is dropped here: each nested CMSDocumentAnalyzer instead lazily
  // instantiates its own default SignaturePolicyProvider. This only differs observably when a
  // caller has set a *custom* SignaturePolicyProvider on the outer analyzer via
  // SetSignaturePolicyProvider.
  // TODO: add an exported SignaturePolicyProvider accessor on analyzer.DefaultDocumentAnalyzer.
  ...
  cadesValidator.SetCertificateVerifier(a.CertificateVerifier())
  cadesValidator.SetContainerContents(a.GetArchiveDocuments())
  // (no SetSignaturePolicyProvider(...) — the Java `this.getSignaturePolicyProvider()` forward is dropped)
  ```
- **Impact:** Java's `ASiCContainerWithCAdESAnalyzer.getSignatureAnalyzers()` forwards the outer
  analyzer's `getSignaturePolicyProvider()` into each nested `CMSDocumentAnalyzer`. The Go port drops
  that forward, so a *custom* `SignaturePolicyProvider` set on the outer analyzer via
  `SetSignaturePolicyProvider` never reaches the per-CMS analyzers — each silently falls back to its
  own default. If a caller relied on a strict custom provider to *tighten* a CMS verdict, the nested
  analyzer's lenient default could yield a different (potentially more accepting) result. Default
  configuration is unaffected (default provider = default behaviour). Already documented in the file
  header with a `TODO`; the gap is narrow and requires an unusual caller setup.
- **Recommendation:** Add a `SignaturePolicyProvider` accessor to `analyzer.DefaultDocumentAnalyzer`
  and forward it in `GetSignatureAnalyzers` (resolves the `TODO`), and record this as a known gap in
  `docs/compatibility/known-gaps.md` until it is closed so it is discoverable.

### A15-STD-001
- **Severity:** Low
- **Category:** STD
- **Location:** `dss/asic/cades/asic_container_with_cades_analyzer.go:199` (`BuildDetachedTimestamps`)
- **Evidence:**
  ```go
  detachedTimestampSource := timestampsrc.NewDetachedTimestampSource()
  for _, timestampAnalyzer := range a.GetTimestampAnalyzers() {
      _ = detachedTimestampSource.AddExternalTimestamp(timestampAnalyzer.Timestamp()) // error discarded
  }
  return detachedTimestampSource.DetachedTimestamps()
  ```
- **Impact:** `DetachedTimestampSource.AddExternalTimestamp` returns an `error`
  (`detached_timestamp_source.go:57`), but it is discarded (`_ =`). If adding one timestamp token
  fails (e.g. an unparseable/malformed token that can't be registered in the source), that timestamp
  is silently dropped from `BuildDetachedTimestamps()`'s result with no alert or log — the report's
  detached-timestamp list can be incomplete relative to the container's actual timestamps. The sibling
  `Signature.AddExternalTimestamp` call at `:247` returns nothing (no error), so only this one is a
  discard. Faithful-enough to Java's void/boolean idiom, but the Go error is real and unhandled.
- **Recommendation:** Surface the error — either panic (Java-parity unchecked) or, better, emit an
  `Alert`/collect it so the report reflects that a timestamp failed to attach. At minimum add a
  comment explaining the intentional discard.

### A15-PERF-001
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `dss/asic/cades/asice_with_cades_container_merger.go:51-77` (`doesNotContain*Document`), `asic_with_cades_container_merger_factory.go:49-60` (`CreateFromDocuments`)
- **Evidence:**
  ```go
  // IsSupportedDocument(ASiCE) -> base walk (IsSupportedZip)
  //   + doesNotContainSignaturesDocument   -> ExtractEntryNames   // full ZIP walk #2
  //   + doesNotContainTimestampsDocument   -> ExtractEntryNames   // full ZIP walk #3
  //   + doesNotContainEvidenceRecordsDocument -> ExtractEntryNames // full ZIP walk #4
  // factory CreateFromDocuments per container:
  //   newASiCSWithCAdESContainerMerger().IsSupportedDocument(c)   // walk
  //   newASiCEWithCAdESContainerMerger().IsSupportedDocument(c)   // walk + 3 walks
  ```
- **Impact:** Selecting a merger for `n` containers walks each container up to ~5× (format-detector
  walk + the ASiC-E merger's three `doesNotContain*` `ExtractEntryNames` walks + the ASiC-S merger's
  walk). Each `ExtractEntryNames` walks + inflates the DEFLATED entries (see A14A-PERF-001's
  `SecureSkip` cost). The `doesNotContain*` checks could share a single `ExtractEntryNames` result.
  Faithful to Java's structure (Java makes the same per-check `getEntryNames` calls), so not a
  divergence, but a real, avoidable repeated-I/O cost during container-merger selection on large
  ASiC-E containers.
- **Recommendation:** In `ASiCEWithCAdESContainerMerger.IsSupportedDocument`, call
  `ExtractEntryNames` once and evaluate all three `doesNotContain*` predicates against that single
  name list; and/or have the factory cache a container's `ExtractEntryNames` result across the
  ASiC-S/ASiC-E support probes.

### A15-SEC-002
- **Severity:** Info
- **Category:** SEC
- **Location:** `dss/asic/cades/testdata/crossgen/main.go:73,79,181,207` + git-tracked `signer_rsa.p12` (2.5 KB), `tsa_ec.p12` (993 B)
- **Evidence:**
  ```go
  signerEntry, err := loadKeyEntry(filepath.Join(selfDir, "signer_rsa.p12"), "testpassword")
  ... validation.NewKeyEntityTSPSourceFromKeyStorePath(filepath.Join(selfDir, "tsa_ec.p12"), "PKCS12",
      "testpassword", "", "testpassword") ... tspSource.SetTsaPolicy("1.2.3.4.5.6.7.8.9")
  ```
- **Impact:** Two private-key PKCS#12 test fixtures (RSA signer, EC self-TSA) and the literal
  password `testpassword` are committed to git under `testdata/`. This is test-only: the keys are
  disposable cross-validation fixtures (the same two stores `cades/testdata/crossgen` uses), the
  policy OID is an explicitly-unregistered test placeholder, and the harness only *builds/validates*
  locally — there is no production secret, no network, and no real trust anchor. Matches upstream
  DSS's own `KeyEntityTSPSource` test fixtures. No exploitable impact.
- **Recommendation:** None required (test-only, expected pattern). If the project later adds a
  secret-scanning gate, allow-list `*/testdata/crossgen/*.p12` + the `testpassword` literal to
  avoid a false positive, and keep the key generation in a checked-in script (already the case via
  the sibling `cades` crossgen).

## Lens summaries

- **STD:** One Low (A15-STD-001, discarded `AddExternalTimestamp` error). Otherwise idiomatic:
  doc comments are thorough and consistently carry the `// Ported from` provenance, the Java→Go
  overloading/embedding rationale (the `Overrides`-registration idiom is applied consistently and
  documented per type), and deliberate divergences are flagged. `defer`-based `try/finally`
  restoration in `SignDocumentMultiple` correctly mirrors Java. No leaky deferred-close; no
  goroutines/concurrency in this package (the A14A shared-handler concurrency question does not
  apply here — `asic/cades` types are not shared across goroutines).
- **PERF:** One Low (A15-PERF-001, repeated ZIP walks in merger selection). The ASiC-E manifest/
  evidence-record dedup loops (`ensure*DocumentsValid`) are O(k²·n²) by name+digest over the small,
  caller-bounded merge set — faithful to Java and on a rare, deliberate merge path, so not a
  standalone finding. No full-container materialization in this layer (streaming/merge delegated to
  reviewed asic-core + `cms`); allocations are proportional to the (small) document set.
- **SEC:** No Critical/High/Medium. Two Info (A15-SEC-001 documented SPP-propagation drop —
  recommend a known-gaps entry; A15-SEC-002 test-only committed fixtures). The merge/level/
  extension/analyzer surfaces are faithful and delegate the byte-exact crypto to already-reviewed
  packages; no false-accept, no false-reject, no panic on hostile input beyond Java-parity
  `IllegalInputException`s, no network, no cross-layer invariant violation found.

## Tool log

Run from `dss/` (Go module root):

| Tool | Command | Result |
|---|---|---|
| gofmt | `gofmt -l asic/cades/` | clean (no output) |
| go vet | `go vet ./asic/cades/...` | clean (no output) |
| golangci-lint | `golangci-lint run --config=../.github/.golangci.yml ./asic/cades/...` | `0 issues.` (exit 0) |

All three tools available and clean. (No tools unavailable.)

## Open questions

1. **SPP propagation (A15-SEC-001).** Confirm no current production caller of the ASiC-C/AdES
   analyzer sets a *custom* `SignaturePolicyProvider` (i.e. the narrow impact is theoretical today).
   If one exists, re-rate to Medium/Low and prioritize the accessor fix.
2. **`AddExternalTimestamp` failure mode (A15-STD-001).** Verify when
   `DetachedTimestampSource.AddExternalTimestamp` can actually return a non-nil error (it is
   accumulating already-parsed tokens); if it is effectively unreachable, this can be re-rated to
   Info with a clarifying comment.
3. **ASiC-S LTA evidence.** For ASiC-S with CAdES at LTA level (no container archive manifest), the
   LTA evidence is the CMS-embedded timestamp + revocation data (not a container artifact). This was
   verified by construction (LT extension + TSA token) but is not exercised by a committed ASiC-S LTA
   oracle fixture in this package — worth confirming against the `harness`/`corpus` ASiC-S LTA case
   before treating ASiC-S LTA as fully pinned.

---

## asic remaining (unit U14b)

- **Unit:** U14b
- **Scope:** `dss/asic` top-level package — every remaining non-test `.go` file **not** covered by
  U14a. Explicitly excluded (already reviewed, not re-read/re-reported):
  `secure_container_handler.go`, `asic_utils.go`, `abstract_asic_container_analyzer.go`,
  `default_container_merger.go`, `abstract_asic_container_evidence_record_builder.go`,
  `dss_zip_entry.go`, `asic_content.go`.
- **Date:** 2026-08-25
- **Depth:** all 40 in-scope files read **FULL** (each ≤274 lines; the largest,
  `abstract_asic_manifest_builder.go`, is 251). Test files skims only (for coverage).
  Cross-checked the security-relevant helpers they call into: `xml/utils/dom_utils.go`,
  `utils/codec.go`, `spi/dss_utils.go` (URI encode/decode), `validation/reports/diagnostic`
  builder-init contract.
- **Context:** Go port of Java DSS 6.5.RC1. `dss/PORTING.md` read first — 1:1 Java port,
  staticcheck style rules disabled, SA1019 allowed, byte-identic enum values, stdlib-first, no cgo.
  `docs/compatibility/known-gaps.md` read — listed gaps not re-reported (notably: evidence-record
  analyzers are *not* registered, so ER *parsing* is a known gap — but the digest builders that
  compute per-data-object hashes over entry contents are in scope here and were assessed).
  asic core (A14A-*) and `cades`/`asic/cades` (C04A/B-*, A15-*) already reviewed — not re-reported.

## Files read

All 40 in-scope non-test files, each **FULL**:

| File (lines) | Coverage |
|---|---|
| `abstract_asic_manifest_builder.go` (251) | **FULL** — `Overrides`/`Init` dispatch, `Build`, root-element/`SigReference`/`DataObjectReference` DOM assembly, digest-value base64, filter null-panic. |
| `abstract_asic_signature_service.go` (222) | **FULL** — generic base, single-doc→multi-doc wrappers, `BuildASiCContainer(At)`→`CreateZipArchiveAt`, `Assert*` panics. |
| `file_archive_entry.go` (212) | **FULL** — lazy `OpenStream` (central-dir lookup by name), `zipFileEntryInputStream`, `DigestValue` cache, `Save` (caller path). |
| `default_asic_container_extractor.go` (210) | **FULL** — `Extract`/`zipParsing` grouping, `containerDocuments` (ASiC-S nested), allow-predicates via `Overrides`. |
| `asic_signature_extension_helper.go` (200) | **FULL** — `ExtractSignatureDocument`, `containsSignatureWithId` (recover-swallow), `CheckSignatureExtensionPossible`. |
| `container_entry_document.go` (197) | **FULL** — wrapper doc, `OpenStream`→content, `DigestValue` cache, `Equals` (ref-identity note). |
| `abstract_asic_container_validator.go` (190) | **FULL** — two-tier `Overrides` dispatch, document-list accessors, `InitializeDiagnosticDataBuilder`. |
| `abstract_asic_filename_factory.go` (188) | **FULL** — `NextAvailableDocumentName`/`documentNameRecursively`, `WithMetaInfFolder`, `Assert*`, `Valid*Filename`. |
| `asic_content_document_filter.go` (176) | **FULL** — 9 category flags + include/exclude, `Filter`/`filterDocuments`. |
| `asic_manifest_parser.go` (173) | **FULL** — `GetManifestFile`/`GetLinkedManifest`, root-element (double-parse), entry/digest/mimetype parsing. |
| `zip_utils.go` (147) | **FULL** — singleton (mutex-guarded), `Extract*`/`CreateZip*` delegation, write-side deviations note. |
| `secure_container_handler_builder.go` (139) | **FULL** — default bomb-guard thresholds, `Build` wiring. |
| `abstract_asic_content_builder.go` (132) | **FULL** — `Build`/`fromZipArchive`/`fromFiles`, name-defaulting, type-mismatch panics. |
| `asic_evidence_record_digest_builder.go` (125) | **FULL** — `BuildDigestGroup`/`AssertConfigurationValid`/`GetDocumentListToComputeDigest`. |
| `zip_content_evidence_record_digest_builder.go` (107) | **FULL** — base digest builder, `ComputeDigestForDocuments`. |
| `asic_evidence_record_manifest_builder.go` (91) | **FULL** — ER manifest leaf, `InitDefaultAsicContentDocumentFilter`=`ArchiveDocumentsFilter`. |
| `asic_content_document_filter_factory.go` (89) | **FULL** — `Empty`/`SignedDocumentsOnly`/`Archive`/`AllSupported`/`All`/`AllowedFilenames` filters. |
| `abstract_asic_data_to_sign_helper_builder.go` (88) | **FULL** — `GetASiCSSignedDocument`/`CreatePackageZip` (multi-file→package ZIP). |
| `asic_container_diagnostic_data_builder.go` (85) | **FULL** — `ContainerInfo` + `getXmlContainerInfo`/`getXmlManifests`. |
| `asic_manifest_validator.go` (70) | **FULL** — `ValidateEntries` (strict name match, digest compare, graceful on missing docs). |
| `asic_manifest_path.go` (30), `asic_manifest_element.go` (60), `asic_manifest_attribute.go` (36), `asic_manifest_namespace.go` (24) | **FULL** — XSD enum/attribute/path/namespace wire tables (values verified byte-exact vs ASiC XSD). |
| `asic_container_merger.go` (38), `asic_container_merger_factory.go` (46), `asic_container_extractor.go` (20), `asic_container_extractor_factory.go` (40), `zip_container_handler.go` (28), `zip_container_handler_builder.go` (17), `dss_zip_entry_document.go` (15), `asic_format_detector.go` (40), `asic_filename_factory.go` (35), `asic_evidence_record_filename_factory.go` (24), `evidence_record_incorporation_service.go` (27), `serializable_asic_container_evidence_record_parameters.go` (16) | **FULL** — interface/registry declarations + `init()`-free registration surface. |
| `abstract_get_data_to_sign_asics.go` (24), `abstract_get_data_to_sign_helper.go` (27), `asic_parameters.go` (108), `asic_container_evidence_record_parameters.go` (94) | **FULL** — value structs + accessors/`Equals`/`String`. |
| Test files (13, 3,634 lines) | **SKIMMED** for coverage: `abstract_asic_filename_factory_test.go`, `default_asic_container_extractor_test.go`, `asic_content_test.go`, `asic_parameters_test.go`, `container_entry_document_test.go`, `zip_utils_test.go`, `dss_zip_entry_test.go`, `secure_container_handler_test.go`, `secure_container_handler_thresholds_test.go`, `asic_utils_test.go`, `broad_corpus_cross_validation_test.go`, `container_merge_cross_validation_test.go`, `abstract_asic_container_analyzer_test.go`. |

## Security posture — VERIFIED CLEAN on the in-scope attack surface

The in-scope files are the orchestration/model layer over the already-hardened core. Verified:

- **Manifest parsing uses the reviewed-clean `internal/xmldom`, no custom entity handling.**
  `asic_manifest_parser.go` builds all DOM through
  `xmlutils.DomUtilsBuildDOMFromDocument` → `DomUtilsBuildDOMFromReader` → `internal/xmldom`
  (reviewed clean, U02-xml-chain). No entity expansion, no XXE, no custom parser. XPath lookups go
  through the reviewed `xmlc14n`/`xpath10` path helpers. **Confirmed clean — not a finding.**
- **No path traversal / no symlink following.** Entry names are never used to build a write path.
  `FileArchiveEntry.Save(filePath)` and `ContainerEntryDocument.Save(filePath)` take a
  *caller-supplied* absolute path; `OpenStream` reads the entry's in-memory content stream (central-dir
  lookup by name), never extracting to a name-derived location. `zip.OpenReader` on a fixed path does
  not follow attacker entry names to disk. Faithful to the U14a conclusion — **not a finding.**
- **Filename-factory collisions are rejected, not silently accepted.** `AssertFilenameValid` panics
  with `IllegalInputException` when the candidate name already exists in the container
  (`IsAvailableName` scan). `ValidDataPackageFilename`/`ValidEvidenceRecordManifestFilename` enforce
  the EN 319 162-1 shape. Predictable names are a spec-documented property, not a vuln (matches A15's
  note on the CAdES factory). **Not a finding.**
- **Content filter is correct for every shipped filter.** All six factory filters
  (`asic_content_document_filter_factory.go`) behave correctly: they either enable a category set
  (no include list) or set only `includedFilenames` (no category enabled), so no document that should
  be covered is dropped. See A14B-SEC-001 for the one Java-parity semantic quirk (include list only
  restricts when the category flag is off) — not triggered by any shipped filter.
- **Evidence-record digests are over the exact entry contents (correct for ASiC ERS/XMLERS).**
  `ComputeDigestForDocuments` delegates per-document hashing to
  `DataObjectDigestBuilderFactory.CreateWithAlgorithm(doc, algo).Build()`, which hashes
  `doc.OpenStream()` content — the entry bytes, not the ZIP framing. The ER *parsing* side remains a
  known gap (unregistered analyzers), not re-reported.
- **Extractor and signature-service delegate to the reviewed, bomb-guarded core.**
  `DefaultContainerExtractor.Extract` → `ZipUtilsInstance().ExtractContainerContent` →
  `SecureContainerHandler.ExtractContainerContent` (U14a: guards on all read paths, bounded
  decompression, enforced CRC, no panic). `BuildASiCContainerAt` → `CreateZipArchiveAt` →
  `BuildZip` (U14a: byte-exact STORED/DEFLATED write). The signing flow's byte-exact CMS work is in
  the reviewed `cades`/`cms` layers (C04A/B, A15). **Not findings here.**
- **No panic on hostile input beyond the documented Java-parity contract.** Malformed/missing manifest
  → `ManifestParserGetManifestFile` returns `nil` (no panic). Missing entry → `FileArchiveEntry.OpenStream`
  returns an error (no panic). Empty signed-doc set → `ManifestValidator.ValidateEntries` returns
  entries as-is (no panic). The `panic(err)`/`panic(IllegalInputException)` sites (extract, zip
  creation, `requireOverrides`) are the PORTING.md unchecked-exception idiom that the `dss` facade's
  `recovered()` converts to errors. **Not a finding.**
- **No network, no `os/exec`, no goroutines** in any in-scope non-test file. The only concurrency is
  the mutex-guarded `ZipUtils` singleton (correctly synchronized; `zipContainerHandler()` reads the
  builder under the same lock — noted and clean).

## Findings

### A14B-SEC-001
- **Severity:** Info
- **Category:** SEC
- **Location:** `dss/asic/asic_content_document_filter.go:154` (`filterDocuments`)
- **Evidence:**
  ```go
  if utils.IsCollectionNotEmpty(f.includedFilenames) {
      for _, d := range documents { if slices.Contains(f.includedFilenames, d.Name()) { result = append(result, d) } }
  }
  if !formatSupported { return result }
  if utils.IsCollectionEmpty(f.excludedFilenames) { return documents }   // <- returns ALL, include list ignored
  ```
- **Impact:** `includedFilenames` only acts as a *restriction* when the category's `formatSupported`
  flag is false. When a category is enabled (`formatSupported==true`) and `excludedFilenames` is empty,
  `filterDocuments` returns **all** of that category's documents, ignoring the include list. So a
  hand-built filter that sets both `includedFilenames` and, say, `SetSignedDocuments(true)` would
  return every signed document, not just the allowed ones. This is a faithful 1:1 port of Java's
  `ASiCContentDocumentFilter.filterDocuments` (identical control flow), so it is **not** a divergence.
  None of the six shipped factory filters combine an include list with an enabled category, so no
  missing-entry / over-inclusion occurs in practice. Impact is a maintainability trap only.
- **Recommendation:** No code change required for parity. Add a doc note on `SetIncludedFilenames`
  (and/or `Filter`) that the include list is only consultable on categories whose flag is off, so a
  future caller does not combine the two expecting intersection semantics.

### A14B-PERF-001
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `dss/asic/asic_manifest_parser.go:44` (`GetManifestRootElement`), `:29` (`GetManifestFile`), `:40` (`GetLinkedManifest`)
- **Evidence:**
  ```go
  func asicManifestParserGetManifestRootElement(manifestDocument model.DSSDocument) *xmldom.Node {
      if !xmlutils.DomUtilsIsDOM(manifestDocument) { return nil }   // full XML parse #1 (result discarded)
      manifestDom, err := xmlutils.DomUtilsBuildDOMFromDocument(manifestDocument)  // full XML parse #2
      ...
  }
  // GetLinkedManifest: for _, manifest := range manifestDocuments { GetManifestRootElement(manifest) ... }
  ```
- **Impact:** `DomUtilsIsDOM` parses the entire DOM and discards it, then
  `DomUtilsBuildDOMFromDocument` parses the same document again — two full XML parses per manifest.
  `GetLinkedManifest` repeats this once per manifest in the container (early-exits on a match, so
  ~M/2 manifests × 2 parses on average), and `GetManifestFile` does it once more. For a
  `FileDocument`-backed container each parse also re-opens the archive via `FileArchiveEntry.OpenStream`
  (central-dir read + entry scan), so it is 2 archive opens + 2 XML parses per manifest. Manifests are
  small (KB), so the cost is low; it is a real, avoidable repeated-parse cost on the validation hot path
  (the analyzer calls these per signature). Faithful to Java (`isDOM` + `buildDOM` both parse), so not a
  divergence.
- **Recommendation:** Parse once and reuse — have `GetManifestRootElement` call
  `DomUtilsBuildDOMFromDocument` directly and treat a parse error as "not a manifest" (dropping the
  separate `DomUtilsIsDOM` pre-check), or memoize the parsed root per document. Low priority.

### A14B-PERF-002
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `dss/asic/abstract_asic_filename_factory.go:85` (`documentNameRecursively`)
- **Evidence:**
  ```go
  func (f *AbstractFilenameFactory) documentNameRecursively(nameTemplate string, existingDocumentNames []string) string {
      number := len(existingDocumentNames) + 1
      candidateName := strings.ReplaceAll(nameTemplate, "001", candidateSuffix)
      if !f.abstractASiCFilenameFactoryOverrides().IsAvailableName(candidateName, existingDocumentNames) {
          grown := make([]string, 0, len(existingDocumentNames)+1)
          grown = append(grown, existingDocumentNames...)   // O(n) copy each level
          grown = append(grown, candidateName)
          return f.documentNameRecursively(nameTemplate, grown)
      }
      return candidateName
  }
  ```
- **Impact:** On a name collision the function copies the whole (growing) name slice and recurses,
  giving O(n·k) work for k consecutive collisions (n = existing document count, attacker-influenced via
  container entry names). In practice k is tiny: `NextAvailableDocumentName` starts at `len+1`, which is
  only a collision if the numbering has a gap, so recursion depth is 0–2 for real containers. Bounded and
  low-impact; noted because the lens asked about allocations-in-loops. Faithful to Java's recursive
  `getDocumentNameRecursively`, so not a divergence.
- **Recommendation:** Rewrite as a `for` loop that increments `number` until `IsAvailableName` passes,
  avoiding the per-level slice copy. Cosmetic/defensive; low priority.

### A14B-STD-001
- **Severity:** Info
- **Category:** STD
- **Location:** `dss/asic/container_entry_document.go:192` (`Equals`), `asic_container_evidence_record_parameters.go:90` (`Equals`), `file_archive_entry.go:205` (`Equals`)
- **Evidence:**
  ```go
  // ContainerEntryDocument.Equals
  return d.content == other.content && d.zipEntry.Equals(other.zipEntry)   // ref-identity for DSSDocument
  // ContainerEvidenceRecordParameters.Equals
  return p.asicEvidenceRecordManifest == other.asicEvidenceRecordManifest // ref-identity
  ```
- **Impact:** `model.DSSDocument` declares no `Equals` in the interface (each implementation carries its
  own), so these `Equals` methods compare the wrapped document by **reference identity**, not by value.
  Two logically-equal-but-distinct document instances compare unequal. Each site is documented in-file
  with the rationale and a cross-reference, and `equals()` is only used for cache/dedup bookkeeping (not
  a security decision), so there is no functional defect. Noted for completeness and to flag that a
  future value-equality change to `DSSDocument` would silently change these comparisons.
- **Recommendation:** None required (documented, faithful). If `model.DSSDocument` ever gains a value
  `Equals`, revisit these three sites to compare by value.

## Lens summaries

- **STD:** One Info (A14B-STD-001, documented reference-identity `Equals`). Otherwise idiomatic and
  consistently documented: every file carries its `// Ported from` provenance, the `Overrides`+`Init`
  virtual-dispatch idiom is applied and explained per type (matching the U14a/U15 precedent), error
  handling follows the `error`/`panic` PORTING.md contract, and `defer`-based cleanup is correct
  (e.g. `zipFileEntryInputStream.Close` closes entry then archive; `CreateZipArchive` closes the
  resources handler). No goroutines; the single mutex-guarded `ZipUtils` singleton is correctly
  synchronized.
- **PERF:** Two Low (A14B-PERF-001 double XML parse per manifest; A14B-PERF-002 recursive slice-copy in
  the filename factory). Both are faithful to Java and low-impact; no O(n²)+ on unbounded attacker
  input, no leak, no repeated full-container materialization in this layer (delegated to the reviewed
  core). The repeated-walk/inflate cost of the extractor is already covered by A14A-PERF-001 and is not
  re-reported.
- **SEC:** No Critical/High/Medium. One Info (A14B-SEC-001, filter include/exclude quirk — not triggered
  by any shipped filter). The manifest parser, entry-name handling, filename collisions, content filter,
  evidence-record digests, extractor, and signature-service delegation are all verified clean or
  correctly delegated to the reviewed core; no false-accept, no path traversal, no XML injection, no
  panic on hostile input beyond the documented Java-parity contract.

## Tool log

Run from `dss/` (Go module root):

| Tool | Command | Result |
|---|---|---|
| gofmt | `gofmt -l asic` | clean (no output) |
| go vet | `go vet ./asic/` | clean (no output) |
| golangci-lint | `golangci-lint run --config=../.github/.golangci.yml ./asic/...` | `0 issues.` (exit 0) |

All three tools available and clean. (No tools unavailable.)

## Open questions

1. **Filter include/exclude semantics (A14B-SEC-001).** Confirm no current caller (in `asic/cades`,
   `asic/xades`, or the `dss` facade) builds a `ContentDocumentFilter` that sets both
   `includedFilenames` and an enabled category flag. If one does, re-rate to Low and fix the doc
   note; if not, keep as Info.
2. **Double XML parse (A14B-PERF-001).** Confirm the per-manifest parse count against a concrete
   multi-manifest validate trace (signature + archive + ER + ER-manifest) before deciding whether the
   memoization is worth the `Overrides`-dispatch complication; manifests are small, so this may stay
   a non-issue.
3. **ER digest builders once analyzers are ported.** The evidence-record digest builders here compute
   per-data-object hashes correctly, but the ER *parsing/verification* side is a known gap
   (unregistered analyzers). Once `dss/evidencerecord` is ported and an
   `EvidenceRecordAnalyzerFactory` is registered, re-verify `BuildDigestGroup` output against a real
   parsed ER (the `evidence_record_incorporation_service.go` local interface copy is flagged for
   replacement at that time).

---

## asic/xades (unit U16)

- **Unit:** U16
- **Scope:** `dss/asic/xades` — ASiC with XAdES: container orchestration over XAdES (service,
  ASiC-S/E mergers, analyzer/validator, filename factories, data-to-sign helpers, manifest
  parser/builder, evidence-record builder, format detector, extractor, and the
  `testdata/crossgen` generator + `extension/*`).
- **Date:** 2026-08-25
- **Depth:** deep-budgeted. 34 non-test Go files (3,757 lines) + 2 `extension/` files +
  `testdata/crossgen/main.go` (179). Two >400-line files read **BUDGETED** (first 250 lines + all
  `func` signatures grepped + every security-relevant region read in full); all others read **FULL**.
  Every function is covered via read-or-signature.
- **Context:** Go port of Java DSS 6.5.RC1. `dss/PORTING.md` read first — 1:1 Java port,
  staticcheck style rules disabled, SA1019 allowed, byte-identical enum values, stdlib-first, no
  cgo. `docs/compatibility/known-gaps.md` read — listed gaps not re-reported (evidence-record
  analyzers unregistered is a *known gap*, not re-reported). **Already reviewed (not
  re-reported):** asic core (A14A-*), asic/cades (A15-*), `xades` core (batch-14), the XML chain
  (`xml/utils`, `internal/xmldom`, `internal/xpath10` — X06–X10-*). This unit covers only the
  `asic/xades` orchestration layer; the security-critical primitives it calls into (XAdES
  sign/extend/counter-sign, XML parse/adopt/serialize, manifest digesting, ZIP read/merge) live in
  those already-reviewed packages and are not re-assessed.

## Files read

All 34 non-test files listed by `find asic/xades -name '*.go' -not -name '*_test.go' | xargs wc -l | sort -rn`,
plus `extension/*` and `testdata/crossgen/main.go`. 32 files read **FULL**; the two >400-line
files read **BUDGETED** (first 250 lines + all `func` signatures grepped + sign/verify/merge/
manifest/digest/level/timestamp/cert/chain/policy/container/xml regions read in full). Every
function is covered via read-or-signature.

| File (lines) | Coverage |
|---|---|
| `asic_with_xades_service.go` (580) | **BUDGETED** — 1–250 + 250–580 (sign/extend/counter-sign/evidence-record/policy-store/`Multiple` adapter) read; all 36 `func` signatures grepped. All functions covered via read-or-signature. |
| `asice_with_xades_container_merger.go` (452) | **BUDGETED** — 1–250 + 250–452 (merge guards, covered-doc names, `createNewManifest`, `ensure*NamesDiffer`, helpers) read; all 28 `func` signatures grepped. All functions covered via read-or-signature. |
| `asics_with_xades_container_merger.go` (367) | **FULL** — ASiC-S merge guards, ID-dup checks, `assertSameRootElement`, `getMergedSignaturesXml` adopt-children. |
| `simple_asic_with_xades_filename_factory.go` (178) | **FULL** — name templates + EN 319 162-1 conformance panics. |
| `asic_container_with_xades_analyzer.go` (152) | **FULL** — signature-analyzer wiring, manifest descriptions, SPP-propagation gap. |
| `asice_with_xades_manifest_parser.go` (148) | **FULL** — delegated XML parse, qualified-attr lookup (documented prior-bug fix). |
| `asic_with_xades_container_merger_factory.go` (143) | **FULL** — ASiC-S/E selection, type-mix reject, `init()` registration. |
| `asice_with_xades_manifest_builder.go` (129) | **FULL** — ODF manifest.xml build, entries-vs-documents conflict. |
| `asic_with_xades_data_to_sign_helper_builder.go` (120) | **FULL** — from-archive/from-files/OpenDocument dispatch. |
| `default_asic_with_xades_filename_factory.go` (103) | **FULL** — default names + `NextAvailableDocumentName`. |
| `asic_container_with_xades_validator.go` (99) | **FULL** — validator leaf + naming note. |
| `extension/asic_with_xades_document_extender.go` (91) | **FULL** — `document` augmentation extender. |
| `asic_with_xades_container_extractor.go` (82) | **FULL** — per-entry-type allow predicates (delegated to reviewed `asic` layer). |
| `asic_with_xades_signature_extension_helper.go` (79) | **FULL** — counter-sign helper, `GetDocumentAnalyzer` panic. |
| `asic_with_xades_container_evidence_record_builder.go` (70) | **FULL** — ER filename-convention assert. |
| `asic_with_xades_signature_parameters.go` (67) | **FULL** — params + `Equals`/`String`. |
| `asic_with_xades_format_detector.go` (66) | **FULL** — XAdES vs CAdES vs OpenDocument routing. |
| `abstract_asic_with_xades_container_merger.go` (57) | **FULL** — merger base + filename factory field. |
| `open_document_support_utils.go` (55) | **FULL** — ODF `external-data/` coverage filter. |
| `asic_container_with_xades_analyzer_factory.go` (54) | **FULL** — factory + `init()` registration. |
| `manifest_element.go` (52) | **FULL** — KAT-verified element enum. |
| `asic_container_with_xades_validator_factory.go` (50) | **FULL** — factory + `init()` registration. |
| `asic_with_xades_container_extractor_factory.go` (49) | **FULL** — factory + `init()` registration. |
| `manifest_path.go` (46) | **FULL** — KAT-verified XPath/attr-path helpers. |
| `manifest_attribute.go` (40) | **FULL** — KAT-verified attribute enum. |
| `data_to_sign_open_document_helper.go` (40) | **FULL** — OpenDocument `ToBeSigned`. |
| `asic_with_xades_asic_content_builder.go` (36) | **FULL** — content-builder leaf. |
| `data_to_sign_asics_with_xades_helper.go` (34) | **FULL** — ASiC-S `ToBeSigned`. |
| `data_to_sign_asice_with_xades_helper.go` (34) | **FULL** — ASiC-E `ToBeSigned`. |
| `get_data_to_sign_asic_with_xades_helper.go` (25) | **FULL** — helper interface. |
| `asic_with_xades_filename_factory.go` (20) | **FULL** — factory interface. |
| `extension/asic_with_xades_document_extender_factory.go` (42) | **FULL** — `init()` registration. |
| `manifest_namespace.go` (18) | **FULL** — KAT-verified namespace. |
| `testdata/crossgen/main.go` (179) | **FULL** (test-only generator — see A16-SEC-001). |
| Test files (4, 560 lines) | **SKIMMED** for coverage: `asice_with_xades_manifest_builder_test.go` (KAT), `manifest_definitions_test.go` (KAT namespace/element/attribute/path), `asic_downstream_cross_validation_test.go` (live-Java round trip), `asic_container_with_xades_validator_smoke_test.go`. |

## Security posture — VERIFIED CLEAN on the asic/xades attack surface

The lens targets (byte-exact container merge, manifest parse, data-to-sign extraction, signature
extension, analyzer wiring) were each verified. **No exploitable defect was found** in the
`asic/xades` orchestration layer. Specifically verified:

- **Container merger does not corrupt or re-serialize an existing signature in a way that breaks
  it.** The ASiC-S merge (`getMergedSignaturesXml`) adopts child `Signature` elements via
  `xmlutils.DomUtilsAdoptChildren` → `xmldom.Node.Import` (deep `Clone`) and serializes with the
  reviewed `internal/xmldom` serializer, which re-emits namespace declarations from the cloned
  `Name.Space`/`Prefix` (see `serialize.go` prefix-mapping pass). Namespace URIs survive, so the
  merged `signatures.xml` still validates — no false-reject from a lost `xmlns`. The ASiC-E merge
  does *not* re-serialize signature files at all; it only rebuilds `manifest.xml` and renames
  colliding entries.
- **Manifest parsing delegates to the reviewed-clean XML chain — no bespoke XML handling.**
  `asice_with_xades_manifest_parser.go` uses `xmlutils.DomUtilsBuildDOMFromDocument` /
  `XPathUtilsGetNodeList` / `DomUtilsBrowseRecursivelyForNamespaceWithUri` (the same `xml/utils`
  helpers `asic` core uses), plus one local qualified-attribute reader
  (`asiceWithXAdESManifestParserAttrValue`) whose correctness is *documented* as a prior-bug fix
  (qualified-name vs expanded-name mismatch) with an explanatory comment — the right shape.
- **Data-to-sign extraction is correct.** `DataToSignASiCE/SWithXAdESHelper.ToBeSigned()` returns
  `asicContent.SignedDocuments()`; the OpenDocument helper returns the ODF `external-data/`-filtered
  coverage. No wrong-extraction (false-accept/reject) path found.
- **Signature extension (counter-sign / policy-store / evidence-record) does not break the
  existing signature.** Each path extracts the single target signature document by id
  (`ExtractSignatureDocument`), hands it to the reviewed `xades` service, and replaces it in the
  content by name — the untouched signatures are carried through verbatim by the container builder.
  `GetXAdESService()` builds a *fresh* service per call, so no cross-call state reuse.
- **Analyzer wiring is correct.** `GetSignatureAnalyzers` builds one `XMLDocumentAnalyzer` per
  signature document and sets the right `DetachedContents`/`ContainerContents` per container type
  (OpenDocument coverage / ASiC-S signed+archive / else all-docs). `GetManifestFilesDescriptions`
  pairs each signature doc with each manifest doc (Java: "All signatures use the same file:
  manifest.xml").
- **No panic on hostile input in a reachable production path.** Malformed container / missing
  manifest / truncated XML surface as returned `error`s from the reviewed `xmlutils`/`xmldom`
  chain (re-raised as `panic` only inside `recovered`-wrapped facade calls, which convert them to
  `error`). The merger's `IllegalInputException` panics are deliberate reject-on-invalid-input,
  mirroring Java.
- **Filename factories are predictable and collision-free** (`NextAvailableDocumentName` is
  KAT-pinned in `asic/abstract_asic_filename_factory_test.go`); no path-traversal sink (entry names
  stay in-memory, `Save` takes a caller-supplied absolute path).
- **No network, no crypto of its own** — all signing/digesting/TSAs delegate to already-reviewed
  `xades`/`spi`/`token`.

## Findings

### A16-SEC-001
- **Severity:** Info
- **Category:** SEC
- **Location:** `dss/asic/xades/extension/asic_with_xades_document_extender.go` (no SPP forward), `asic_container_with_xades_analyzer.go:82-110` (`GetSignatureAnalyzers`)
- **Evidence:**
  ```go
  // TODO: add an exported SignaturePolicyProvider accessor on analyzer.DefaultDocumentAnalyzer.
  func (a *ASiCContainerWithXAdESAnalyzer) GetSignatureAnalyzers() []analyzer.DocumentAnalyzer {
      ...
          documentAnalyzer.SetCertificateVerifier(a.CertificateVerifier())
          // (no SetSignaturePolicyProvider(...) — the Java getSignaturePolicyProvider() forward is dropped)
          ...
  }
  ```
- **Impact:** Java's `ASiCContainerWithXAdESAnalyzer.getSignatureAnalyzers()` forwards the outer
  analyzer's `getSignaturePolicyProvider()` into each nested `XMLDocumentAnalyzer`. The Go port drops
  that forward, so a *custom* `SignaturePolicyProvider` set on the outer analyzer via
  `SetSignaturePolicyProvider` never reaches the per-signature analyzers — each silently falls back
  to its own default. If a caller relied on a strict custom provider to *tighten* an XAdES verdict,
  the nested analyzer's lenient default could yield a different (potentially more accepting) result.
  Default configuration is unaffected. This is the **same** gap as A15-SEC-001 (asic/cades),
  independently present in this package. Already documented in the file header with a `TODO`.
- **Recommendation:** Add a `SignaturePolicyProvider` accessor to `analyzer.DefaultDocumentAnalyzer`
  and forward it in both `asic/cades` and `asic/xades` `GetSignatureAnalyzers` (resolves both `TODO`s
  at once); until then record it once in `docs/compatibility/known-gaps.md` so it is discoverable.

### A16-SEC-002
- **Severity:** Info
- **Category:** SEC
- **Location:** `dss/asic/xades/testdata/crossgen/main.go:61,67,74,169` + git-tracked `signer_rsa.p12`, `tsa_ec.p12`
- **Evidence:**
  ```go
  signerEntry, err := loadKeyEntry(filepath.Join(selfDir, "signer_rsa.p12"), "testpassword")
  ... validation.NewKeyEntityTSPSourceFromKeyStorePath(filepath.Join(selfDir, "tsa_ec.p12"), "PKCS12",
      "testpassword", "", "testpassword") ... tspSource.SetTsaPolicy("1.2.3.4.5.6.7.8.9")
  ```
- **Impact:** Two private-key PKCS#12 test fixtures (RSA signer, EC self-TSA) and the literal
  password `testpassword` are committed under `testdata/`. Test-only: the keys are disposable
  cross-validation fixtures (the same two stores `cades/testdata/crossgen` uses), the policy OID is an
  explicitly-unregistered test placeholder, and the harness only builds/validates locally — no
  production secret, no network, no real trust anchor. Matches upstream DSS's own `KeyEntityTSPSource`
  test fixtures. No exploitable impact. (Same pattern as A15-SEC-002.)
- **Recommendation:** None required (test-only, expected pattern). If a secret-scanning gate is added,
  allow-list `*/testdata/crossgen/*.p12` + the `testpassword` literal to avoid a false positive.

### A16-STD-001
- **Severity:** Low
- **Category:** STD
- **Location:** `dss/asic/xades/asice_with_xades_manifest_parser.go:64-104` (`getEntries`)
- **Evidence:**
  ```go
  result = []*model.ManifestEntry{}
  defer func() { if recover() != nil { result = []*model.ManifestEntry{} } }()
  manifestDom, err := xmlutils.DomUtilsBuildDOMFromDocument(p.manifestDocument)
  if err != nil { return []*model.ManifestEntry{} }   // malformed/undecodable manifest -> empty, no error
  ...
  if !xmlutils.DomUtilsIsDOM(p.manifestDocument) { return []*model.ManifestEntry{} }
  ```
- **Impact:** `getEntries()` swallows *every* failure to build the manifest DOM (parse error,
  undecodable bytes) and a panic alike, returning an empty entry list rather than an error. Mirrors
  Java's broad `catch (Exception e) { LOG.warn(...); return emptyList; }`, so it is faithful — but a
  *corrupt/truncated* manifest silently yields an empty `ManifestFile`, which a downstream consumer
  (e.g. a report's manifest description, or the merger's `getManifestFileEntries`) can read as
  "no entries" rather than "unreadable manifest." In a validation context this is a soft false
  negative (missing coverage) rather than a hard one. Java-parity, so not a divergence defect.
- **Recommendation:** Keep the Java-parity swallow, but at minimum emit an `Alert` (the `alert`
  package is already a dependency of the analyzer path) or a `DIVERGENCE`-free comment noting the
  empty-on-error contract, so a caller that surfaces `ManifestFile` to a report knows an empty
  list can mean "unreadable," not just "no entries."

### A16-PERF-001
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `dss/asic/xades/asics_with_xades_container_merger.go:307-319` (`checkDuplicatesPresent`), `asice_with_xades_container_merger.go:304-316` (`getConflictingDocumentNames`)
- **Evidence:**
  ```go
  func (m *ASiCSWithXAdESContainerMerger) checkDuplicatesPresent(strs []string) bool {
      for _, s := range strs {
          count := 0
          for _, other := range strs { if other == s { count++ } }   // O(n^2)
          if count > 1 { return true }
      }
  ```
- **Impact:** Duplicate detection is O(n²) over the collected id/name lists. Called 3× for ASiC-S
  (`checkNoCommonIds*`) and for the ASiC-E conflict scan. On the merge path `n` is the total number
  of signature/manifest identifiers across the small, caller-bounded set of containers being merged —
  so `n` is tiny in practice and this is not a realistic cost. Faithful to Java's structure (Java
  likewise scans linearly). Not attacker-driven to any large `n` in a supported flow.
- **Recommendation:** Optional: replace the double loop with a `map[string]int` count (O(n)) for
  readability; no behavioural change. Skip if the Java-verbatim mapping is preferred here.

### A16-PERF-002
- **Severity:** Low
- **Category:** PERF-MEM
- **Location:** `dss/asic/xades/asice_with_xades_container_merger.go:66-81` (`doesNotContain*Document`), `asic_with_xades_container_merger_factory.go:53-92` (`CreateFromDocuments`)
- **Evidence:**
  ```go
  func (m *ASiCEWithXAdESContainerMerger) IsSupportedDocument(container model.DSSDocument) bool {
      ...
      return !isASiCS || (m.doesNotContainSignaturesDocument(container) &&    // ExtractEntryNames walk #1
          m.doesNotContainEvidenceRecordsDocument(container))                  // ExtractEntryNames walk #2
  }
  // factory, per container:
  //   newASiCSWithXAdESContainerMerger().IsSupportedDocument(c)  // walk
  //   newASiCEWithXAdESContainerMerger().IsSupportedDocument(c)  // walk + 2 walks
  ```
- **Impact:** Selecting a merger for `n` containers walks each container up to ~4× (format-detector
  `ExtractEntryNames` + the ASiC-E merger's two `doesNotContain*` `ExtractEntryNames` + the ASiC-S
  merger's walk), and the factory re-probes each container against both mergers. Each
  `ExtractEntryNames` walks + inflates DEFLATED entries (see A14A-PERF-001's `SecureSkip` cost). The
  two `doesNotContain*` checks could share a single `ExtractEntryNames` result. Faithful to Java's
  structure (Java makes the same per-check `getEntryNames` calls), so not a divergence, but a real,
  avoidable repeated-I/O cost during merger selection on large ASiC-E containers. (Same pattern as
  A15-PERF-001 in asic/cades.)
- **Recommendation:** In `ASiCEWithXAdESContainerMerger.IsSupportedDocument`, call
  `ExtractEntryNames` once and evaluate both `doesNotContain*` predicates against that single name
  list; and/or have the factory cache a container's `ExtractEntryNames` result across the ASiC-S/ASiC-E
  support probes.

## Lens summaries

- **STD:** One Low (A16-STD-001, manifest parser swallows errors → empty entry list). Otherwise
  idiomatic: doc comments are thorough and consistently carry the `// Ported from` provenance, the
  Java→Go overloading/embedding rationale (the `Overrides`-registration idiom is applied
  consistently and documented per type, including the `Multiple`-adapter trick that lets one Go type
  satisfy both `SignatureService` and `MultipleDocumentsSignatureService`), and deliberate divergences
  are flagged (e.g. the manifest qualified-attribute lookup fix). The
  `getSignedDataObjectIds` inner `if referenceURI == ""` is *dead* (guarded by an outer
  `referenceURI != ""`) but is ported verbatim with an explanatory comment — not a finding. No
  leaky deferred-close; no goroutines/concurrency in this package.
- **PERF:** Two Low (A16-PERF-001 O(n²) duplicate scans, A16-PERF-002 repeated ZIP walks in merger
  selection) — both faithful to Java, both on rare/deliberate merge paths with small `n`. No
  full-container materialization in this layer (streaming/merge delegated to reviewed asic-core);
  allocations are proportional to the (small) document set.
- **SEC:** No Critical/High/Medium. Two Info (A16-SEC-001 documented SPP-propagation drop — same as
  A15-SEC-001; A16-SEC-002 test-only committed fixtures — same as A15-SEC-002). The merge/level/
  extension/analyzer/data-to-sign surfaces are faithful and delegate the byte-exact crypto/XML to
  already-reviewed packages; no false-accept, no false-reject, no panic on hostile input beyond
  Java-parity `IllegalInputException`s, no network, no cross-layer invariant violation found.

## Tool log

Run from `dss/` (Go module root):

| Tool | Command | Result |
|---|---|---|
| gofmt | `gofmt -l asic/xades/` | clean (no output) |
| go vet | `go vet ./asic/xades/...` | clean (`vet-ok`) |
| golangci-lint | `golangci-lint run --config=../.github/.golangci.yml ./asic/xades/...` | `0 issues.` (exit 0) |

All three tools available and clean. (No tools unavailable.)

## Open questions

1. **SPP propagation (A16-SEC-001).** Same as A15-SEC-001 — resolve with a single
   `SignaturePolicyProvider` accessor on `analyzer.DefaultDocumentAnalyzer` shared by both the
   cades and xades analyzers, and add one `known-gaps.md` entry covering both packages.
2. **Manifest empty-on-error (A16-STD-001).** Confirm no production caller of
   `ASiCEWithXAdESManifestParser.Manifest()` relies on a non-empty `ManifestFile` as a
   well-formedness signal (i.e. the empty-list-on-corrupt-manifest path is not load-bearing for a
   verdict). If it is, re-rate to Medium and surface the error/Alert.
3. **ASiC-S XAdES LTA evidence.** For ASiC-S with XAdES at LTA level there is no container archive
   manifest; the LTA evidence is the XAdES-embedded timestamp + revocation data. This is exercised by
   the `testdata/crossgen` B/T pair (baseline B and T generated and handed to `CrossGenValidator.java`),
   but an LTA-level ASiC-S XAdES container is *not* among the committed crossgen cases (only B and T) —
   worth confirming against the `harness`/`corpus` ASiC-S XAdES LTA case before treating that level as
   fully pinned in this package.
