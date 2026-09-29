# Review — Batch 06 (Trust & Revocation)

- Scope: trust/revocation packages of the Go port (`dss/crlparser`, and subsequent
  trust units).
- Date: 2026-08-23
- Depth: deep
- Compatibility contract: interop with Java DSS 6.5.RC1 (byte-level + verdict parity).
  `dss/PORTING.md` conventions and `docs/compatibility/known-gaps.md` are treated as
  deliberate, not findings.

## crlparser (unit U21)

CRL parsing and revocation-status lookup: X.509 CRL parse (`crypto/x509`), CRL signature /
issuer / key-usage validity, and per-serial revocation entry matching.

### Files read

In-scope (read fully):
- `crlparser/crl_utils_x509crl_impl.go` (492)
- `crlparser/crl_validity.go` (414)
- `crlparser/abstract_crl_utils.go` (272)
- `crlparser/crl_utils.go` (43)
- `crlparser/icrl_utils.go` (37)
- `crlparser/x509_crl_validity.go` (29)
- `crlparser/pem_to_der_converter.go` (25)
- `crlparser/crl_binary.go` (21)
- `crlparser/doc.go` (12)

Test files (skimmed for coverage only): `crl_utils_test.go`, `crl_validity_test.go`,
`abstract_crl_utils_test.go`, `crl_utils_x509crl_impl_test.go`, `crl_indirect_kat_test.go`.

Cross-checked against upstream Java: `dss-crl-parser/.../AbstractCRLUtils.java`,
`CRLValidity.java`, `PemToDerConverter.java`, and
`dss-crl-parser-x509crl/.../CRLUtilsX509CRLImpl.java`, plus the Go 1.27
`crypto/x509.ParseRevocationList` (parser.go:1190) and the `tbsCertificateList` struct.

### Findings

**STD (idioms, error handling, defer/cleanup, context, sync, doc comments, test quality)**

C21-STD-001 — Info — STD — `crlparser/doc.go:3`
- Evidence: package doc claims "a **streaming** parser for X.509 Certificate Revocation
  Lists that **avoids loading a large CRL fully into memory**", but the implementation
  (`crlUtilsParseRevocationList` → `x509.ParseRevocationList`, and
  `crlUtilsX509CRLImplRevocationInfo`) always decodes the whole CRL into
  `[]x509.RevocationListEntry` and retains it on `CRLValidity.x509CRL`.
- Impact: Documentation overstates the memory profile of the implementation; a reader
  optimising for huge CRLs would expect a constant-memory scan that does not exist. No
  functional defect.
- Recommendation: Soften the doc.go wording (e.g. "parses X.509 CRLs and answers
  is-revoked lookups") or implement a true entry-streaming lookup. Low priority; cosmetic.

No other STD findings. Error handling is correct (`error` returns at the facade boundary;
the two in-package `panic`s — `NewCRLValidity` on nil and `crlUtilsX509CRLImplRevocationInfo`
on an unparseable CRL — are 1:1 ports of upstream `Objects.requireNonNull` /
unchecked-exception semantics, not reachable by hostile input on the normal path). No
goroutines, so no sync issues. Tests are thorough (KAT vectors, reasonCode-absent-vs-zero
disambiguation, indirect-CRL certificateIssuer carry-forward, PEM↔DER parity, hostile-input
error paths, and the v1 splice indirectly via the minimal-CRL fixtures).

**PERF (CRL parsing cost, allocations)**

C21-PERF-BIGO-001 — Low — PERF-BIGO — `crlparser/crl_utils_x509crl_impl.go:288` (loop at :302–304)
- Evidence: `crlUtilsX509CRLImplRevocationInfo` does a linear `for i := range
  revocationList.RevokedCertificateEntries { … entry.SerialNumber.Cmp(serialNumber) == 0 }`
  over all entries per call (O(n) per serial lookup); no index by serial number.
- Impact: For a CRL with N entries, one revocation query costs O(N); validating K
  certificates against the same CRL costs O(K·N). This is parity with BouncyCastle's
  `X509CRLObject#getRevokedCertificate` (which also scans the entry collection), so it is
  not a divergence — but it is a genuine cost on large CRLs. Not exploitable: the caller
  (`spi.CRLToken.setRevocationStatus`) issues one lookup per certificate being validated.
- Recommendation: Optional hardening — build a `map[*big.Int]*CRLEntry` (or sort by serial)
  once per CRL if large-CRL + many-certificate workloads matter. Document as an accepted
  parity cost if left as-is.

C21-PERF-MEM-002 — Info — PERF-MEM — `crlparser/crl_utils_x509crl_impl.go:150` (splicing)
- Evidence: `crlUtilsSpliceV2VersionForParsing` builds a throwaway patched DER copy
  (rebuilt `patchedTBS` + `patchedList` via `cryptobyte.Builder`) for every CRL that
  `ParseRevocationList` rejects, before the splice is applied.
- Impact: One extra O(size-of-CRL) allocation for the v1 (versionless) CRL shape only; the
  patched copy is discarded after parsing and `Raw`/`RawTBSRevocationList` are restored to
  the caller's bytes. Negligible; noted for completeness.
- Recommendation: None required.

No unbounded-memory or O(n²) attacker-input findings. Allocations in the hot path are bounded
by CRL size and match the "parse once, retain" design.

**SEC (top priority)**

No findings. Each mandated lens was verified clean:

- **CRL validity window / false-accept:** `CRLValidity.IsValid`
  (`crl_validity.go`) = `issuerX509PrincipalMatches && signatureIntact && crlSignKeyUsage &&
  !IsUnknownCriticalExtension()`, byte-for-byte the upstream `CRLValidity.isValid()`.
  Notably it does **not** enforce the thisUpdate/nextUpdate window — this is correct
  parity: upstream enforces freshness in the validation layer (see
  `dss/spi/validation/revocation_data_verifier.go` `isRevocationThisUpdateAfterValidationTime`
  and `signature_validation_context.go` `earliestNextUpdate`), not in `CRLValidity`. So there
  is no false-accept of an out-of-window CRL and no false-reject of a current one within this
  package's contract.
- **Revocation entry matching (serial comparison):** `crl_utils_x509crl_impl.go:304` uses
  `entry.SerialNumber.Cmp(serialNumber) == 0` — exact integer equality, no false match and no
  false miss. A non-matching serial returns `nil` → `CRLToken` sets `GOOD`; a match sets
  `REVOKED` with date + reason. Indirect-CRL `certificateIssuer` inheritance
  (`newCRLEntry`) is RFC 5280 §5.3.3-correct and KAT-pinned (`crl_indirect_kat_test.go`).
- **pem_to_der_converter (hostile PEM):** `pem.Decode` never panics on malformed input and
  returns a nil block → `model.NewDSSError("Unable to read PEM Object")`; base64 decoding is
  done inside `pem.Decode` and bounds-checked. No panic path.
- **crl_utils_x509crl_impl / crypto/x509 usage:** `crlUtilsCheckSignatureValue` verifies the
  CRL signature via `x509.Certificate.CheckSignature(revocationList.SignatureAlgorithm,
  revocationList.RawTBSRevocationList, revocationList.Signature)`. The v1 splice restores
  `RawTBSRevocationList` to the caller's original bytes before this call, so the signature is
  verified over the genuine TBS — no bypass. Issuer/subject matching goes through
  `model.X500Principal.Equals` (canonical-name comparison), the same primitive upstream uses.
  `NotBefore/NotAfter` are not used for the validity decision (window is a validation-layer
  concern, as above). MD5-signed CRLs are refused by `crypto/x509`, which is a deliberate
  documented hardening over JCA.
- **Reason codes:** `newCRLEntry` disambiguates "reasonCode extension absent" (nil) from
  "present with value 0 / Unspecified" (0) by inspecting the entry's Extensions — resolving
  `crypto/x509.ReasonCode`'s documented zero-ambiguity; KAT-pinned
  (`TestNewCRLEntry_ReasonCodeAmbiguity`). Revocation is recorded for any present entry
  (reason is metadata), matching upstream.
- **Panics on hostile CRL input:** No panic path reachable via `CRLUtilsBuildCRLBinary` /
  `CRLUtilsBuildCRLValidity` with hostile DER/PEM. Both `ParseRevocationList` and the splice
  use bounds-checked `cryptobyte.String` reads (no unbounded recursion / expansion). The only
  in-scope `panic` (`crlUtilsX509CRLImplRevocationInfo`) fires only when a `CRLValidity`
  carries no cached CRL **and** re-parse fails — an internal invariant, not hostile input.
- **Time handling:** Dates are carried as absolute `time.Time` (UTC); comparisons in the
  validation layer use `time.Time.After/Before/Equal` (instant-based), so there is no
  timezone/clock-skew bug introduced by this package. GeneralizedTime-only decoding for
  `expiredCertsOnCRL`/`invalidityDate` matches upstream's `ASN1GeneralizedTime` requirement.

### Tool log

Run from `dss/`:

- `gofmt -l crlparser` → no output (clean).
- `go vet ./crlparser/...` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./crlparser/...` → `0 issues`.
- `go test ./crlparser/... -count=1` → `ok  github.com/ryftcore/dss-go/dss/crlparser  0.689s`.

All tools available and passing; nothing to note as unavailable.

### Open questions

1. `CRLValidity.IsValid` deliberately omits the thisUpdate/nextUpdate window check (parity
   with upstream). Confirmed the window is enforced by `spi/validation`
   (`RevocationDataVerifier` / `SignatureValidationContext`). No action in this package, but
   the invariant "CRL consumers must check freshness downstream" is implicit and worth a one-
   line note in `doc.go` if a future refactor adds a convenience `IsCurrent()` helper.
2. The v1-CRL splice (`crlUtilsSpliceV2VersionForParsing`) is a sound and necessary
   deviation (Go 1.27 `ParseRevocationList` rejects a versionless TBSCertList). It is
   documented and effectively covered, but there is no dedicated unit test that feeds a
   genuine versionless CRL and asserts the splice engages and the signature still verifies.
   Consider adding one if a v1 fixture (e.g. the Belgium Root CA CRLs referenced in the
   header) becomes available in `testdata/`.

## tsl PART 1 (unit U18a)

- Scope: `dss/tsl` — TSL/LoTL (List of Trusted Lists) parsing & validation layer,
  **PART 1** = the first 48 files of `find tsl -name '*.go' -not -name '*_test.go' | xargs wc -l | sort -rn`
  (all in-scope files ≤369 lines, read in full). PART 2 (the other ~47 files) is a separate unit.
- Date: 2026-08-23
- Depth: deep
- Compatibility contract: interop with Java DSS 6.5.RC1. `dss/PORTING.md` conventions and
  `docs/compatibility/known-gaps.md` treated as deliberate, not findings.

### Files read (48, in-scope — all read fully)

`sha2_file_cache_data_loader.go` (369), `trusted_list_certificate_source_synchronizer.go` (309),
`tl_validation_job_summary_builder.go` (280), `abstract_parsing_task.go` (259),
`trust_service_equivalence_converter.go` (257), `criteria_list_converter.go` (253),
`tl_parsing_cache_dto.go` (240), `trust_service_provider_builder.go` (239), `tl_validation_job.go` (212),
`tl_parsing_cache_dto_builder.go` (193), `document_with_sha2.go` (193),
`lotl_with_pivots_analysis.go` (192), `default_trusted_list_with_sha2_predicate.go` (189),
`trust_service_converter.go` (172), `lotl_parsing_task.go` (171), `composite_condition.go` (164),
`tl_parsing_task.go` (163), `tl_validator_task.go` (162), `trust_service_provider_converter.go` (156),
`other_tsl_pointer_converter.go` (147), `mra_path.go` (143),
`cert_subject_dn_attribute_condition.go` (139), `extended_key_usage_condition.go` (131),
`abstract_tl_parsing_result.go` (121), `abstract_trusted_list_with_sha2_predicate.go` (120),
`qc_statement_condition.go` (115), `tl_source.go` (113), `mra_element.go` (110),
`pivot_processing.go` (108), `lotl_source.go` (106), `key_usage_condition.go` (100),
`expiration_and_signature_check_strategy.go` (96), `lotl_change_applier.go` (93),
`policy_id_condition.go` (92), `xml_download_task.go` (85), `tl_validation_result.go` (76),
`lotl_signing_certificates_announcement_scheme_information_uri.go` (74),
`abstract_runnable_tl_analysis.go` (69), `lotl_parsing_result.go` (68),
`abstract_other_tsl_pointer_predicate.go` (68), `tl_source_builder.go` (57),
`tl_read_only_cache_access.go` (54), `international_names_type_converter.go` (52),
`tl_predicate_factory.go` (50), `pivot_processing_result_from_cache_access_builder.go` (49),
`granted_trust_service.go` (49), `xml_download_result.go` (47),
`official_registration_identifier_predicate.go` (47).

Test files (coverage only, not read fully): `*_test.go` in `tsl/`, plus the cross-format oracle
`harness/tsl_parsing_oracle_test.go` and the Java oracle generator
`corpus/harness/testdata/oracle/tsl/gen/TSLParsingOracle.java`.

Cross-checked against upstream Java where the port is not self-evident: `TrustServiceEquivalenceConverter`
MRA condition/replacement direction, `AbstractParsingTask` nil-`SchemeInformation` behaviour, the
`validation/job` cache-entry state machine, and `trustedlist/jaxb` / `internal/xmldom` XML handling.

### Findings

**SEC (top priority)**

T18A-SEC-001 — **High** — Category: SEC — `tsl/lotl_with_pivots_analysis.go:165` (goroutine fan-out) & `tsl/pivot_processing.go:98` (`expireCache`)
- Evidence: `downloadAndParseAllPivots` runs one goroutine per pivot needing a refresh
  (`lotl_with_pivots_analysis.go:166-181`, `mu` guards only the `processingResults` map). Each
  `PivotProcessing` is handed the SAME `lotlCacheAccessByKey` (`lotl_with_pivots_analysis.go:139`) and a
  shared `preceedingPivotCacheAccessByKeyList`. `PivotProcessing.expireCache`
  (`pivot_processing.go:98-106`) then calls `p.lotlCacheAccess.ExpireValidation()` and
  `pivotCacheAccessByKey.ExpireValidation()` on those shared entries.
- Root cause (out of this unit's scope, in `dss/validation/job`): `AbstractCache.Expire(key)` is
  `Get(key).Expire()` (`abstract_cache.go:96-98`) — `Get` holds `a.mu` only to fetch the
  `*CachedEntry`, then `entry.Expire()` → `CurrentCacheContext.RefreshNeeded()` writes
  `c.state`/`c.lastStateTransitionTime`/`c.exception` (`current_cache_context.go:60-71,92-94`) with
  **no lock**. `CachedEntry`/`CurrentCacheContext` carry no mutex.
- Impact: With ≥2 pivots (the normal EU pivot chain 226→191→172), the goroutines concurrently
  read+write the same shared LOTL/pivot cache-entry state. A Go data race = undefined behaviour:
  lost state transitions (a pivot/LOTL left "fresh" when it should have been expired → stale trust
  state), torn `c.state` interface value / `c.exception` pointer, or a crash. This is the concurrent
  path Java guards with its own synchronization; the Go port's map-mutex does not cover the entry
  transition. Trust-relevant (pivot validation state feeds the accepted-signer set).
- Recommendation: Make `CachedEntry`/`CurrentCacheContext` state transitions thread-safe in
  `dss/validation/job` (per-entry lock, or move the transition under `AbstractCache.mu`), and/or add a
  `recover()` + shared-state guard in the pivot goroutines. Flag for the `validation/job` unit; the
  tsl-side trigger is the concurrent fan-out above. Add a `-race` test that drives ≥2 pivots.

T18A-SEC-002 — **High** — Category: SEC — `tsl/lotl_with_pivots_analysis.go:170` (unrecovered `go func`) & `tsl/lotl_parsing_task.go:100` (`extractOtherTSLPointers`)
- Evidence: The pivot goroutines call `processing.Call()` with **no `recover()`**
  (`lotl_with_pivots_analysis.go:169-179`). `Call()` → `parsing()` → `LOTLParsingTask.Get()` →
  `parseSchemeInformation(result, jaxbObject.SchemeInformation)`; `extractOtherTSLPointers` dereferences
  `schemeInformation.PointersToOtherTSL` with no nil guard (`lotl_parsing_task.go:100`). The file's own
  NOTE confirms a TL missing `<SchemeInformation>` raises a null-deref (parity with the upstream NPE).
- Impact: A hostile/malformed **pivot** LOTL that omits `<SchemeInformation>` panics inside the pivot
  goroutine. Unlike the main analysis path — which is contained by `AbstractRunnableAnalysis.Run`'s
  `defer func(){ _ = recover() }()` (`validation/job/abstract_runnable_analysis.go:97`) — the pivot path
  has no recover, so the panic is not caught and **crashes the whole process** (an un-recovered panic in
  any goroutine is fatal). DoS via a single malicious pivot in the LOTL's pivot chain; the main path's
  identical input is safely contained, so this is an inconsistency, not an upstream-parity behaviour.
- Recommendation: Wrap the pivot goroutine body in `defer func(){ if r := recover(); r != nil { ... } }()`
  (recording a `ParsingError` on the pivot cache access instead of letting it escape), matching the
  main path's recover contract. Add a hostile-pivot fixture (LOTL with a pivot whose body lacks
  `<SchemeInformation>`) asserting the job survives.

Verified clean on the other mandated SEC lenses (each confirmed, not assumed):

- **TSL's own signature validation:** `TLValidatorTask.Get` (`tl_validator_task.go:70-118`) runs the full
  frozen XAdES validator over the TL document with a `CommonCertificateVerifierSimple(true)` whose trusted
  set is exactly the TL's announced SDIs (`buildTrustedCertificateSource`), `ValidateDocumentWithValidationPolicy`
  over the embedded `tsl-constraint.xml`. `fillResult` requires exactly one signature and surfaces
  `Indication`/`SubIndication`; `ValidationCacheDTO.IsValid()` = `TOTAL_PASSED`
  (`validation/job/validation_cache_dto.go:98`). So an invalid TSL signature ⇒ not `TOTAL_PASSED` ⇒ not
  treated as valid. No false-accept of an untrusted trust anchor.
- **Pivot chain:** `LOTLWithPivotsAnalysis.currentCertificateSourceFromPivots` walks the pivots in
  **reversed** order and only promotes a pivot's announced certs into the current signer set when that
  pivot's cached validation is `IsResultExist() && IsValid()` (`lotl_with_pivots_analysis.go:85-101`). A
  broken pivot (validation absent/invalid) simply does not advance the signer set — no silent trust-chain
  skip. (The concurrency hazards of this path are T18A-SEC-001/002.)
- **MRA equivalence:** `trust_service_equivalence_converter.go` + `validation/reports/diagnostic`
  `xml_trust_service_provider_builder.go` — see Open question 1; the consumer enacts a replacement only when
  `Condition.Check(cert)` holds, and `MRACertificateEquivalenceApplied` fails the check when any equivalence is
  `!Enacted` (conservative). No unconditional trust grant.
- **Expiration/parsing-error detection:** `DefaultTrustedListWithSha2Predicate.Test` forces a refresh when the
  sha2 digest mismatches, `NextUpdate` is reached, or the cache expired (`default_trusted_list_with_sha2_predicate.go:110-160`).
  `ExpirationAndSignatureCheckStrategy.isSyncSupported` refuses to sync expired or non-`TOTAL_PASSED` lists by
  default (`expiration_and_signature_check_strategy.go:78-95`). `TLParsingCacheDTOBuilder`/`AbstractCacheDTO`
  correctly propagate `IsResultExist`/`IsError`/`IsDesynchronized` (`validation/job/abstract_cache_dto.go:140-167`).
  No missed-expiration accept.
- **Criteria list:** `CriteriaListConverter.Apply` builds `CompositeCondition`s; `CompositeCondition.Check`
  treats an unknown indicator as `false` (conservative, `composite_condition.go:73-95`); `CertSubjectDNAttribute`,
  `ExtendedKeyUsage`, `QCStatement`, `KeyUsage`, `PolicyId` conditions all default to "does not match" on absent
  data. No criterion that silently grants trust.
- **TSP / service matching:** `TLParsingTask.filter` and the `GrantedTrustService` / `OfficialRegistrationIdentifier`
  / `TradeName` predicates are *removal* filters (drop non-conforming entries) — conservative. `GrantedTrustService`
  returning true for a historically-granted status is parity with upstream. No false-accept.
- **SHA2 cache poisoning:** `Sha2FileCacheDataLoader.GetDocumentRefresh` only returns a cached document when
  `DefaultTrustedListWithSha2Predicate` accepts it (digest match + not expired); a mismatched/absent sha2 forces a
  re-download (`sha2_file_cache_data_loader.go:184-259`). The published `.sha2` digest is the integrity reference —
  consistent with ETSI TS 119 612 §6.1; no attacker-plantable cache entry bypasses it. `IsHexEncoded`/`FromHex`
  (`utils/codec.go`) are bounds-checked (`hex.DecodeString`).
- **XML parsing:** TSL/LOTL bodies parse via `trustedlist.NewFacade().Unmarshal` → `trustedlist/jaxb`
  (not this package's own XML handling); the only DOM read here is `NextUpdate` via the reviewed
  `xml/utils` + `internal/xmldom` (`abstract_trusted_list_with_sha2_predicate.go:98-118`).
  `criteria_list_converter.go` re-decodes the opaque `xs:any` token stream with Go's stdlib `encoding/xml`
  over a finite captured stream (`tokenSliceReader`/`decodeTokenSubtree`) — no unbounded recursion.
- **Panics on hostile TSL input:** All constructor `panic`s are `Objects.requireNonNull` parity. The one hostile-
  input-reachable panic (nil `SchemeInformation`, T18A-SEC-002) is on the pivot path; on the *main* path it is
  contained by `Run()`'s recover and does not yield a false-accept (a TL with no parsing result simply
  contributes no trust properties).
- **Concurrency / goroutine leaks:** No goroutine leaks — every pivot goroutine does `defer wg.Done()` and the
  fan-out `wg.Wait()`s before returning (`lotl_with_pivots_analysis.go:166-181`). The only sync defect is the
  shared-state data race in T18A-SEC-001.

**PERF (parsing cost, cache growth, allocations, leaks)**

T18A-PERF-BIGO-001 — **Low** — Category: PERF-BIGO — `tsl/trust_service_equivalence_converter.go:218` (`appendUnique`)
- Evidence: `appendUnique(list, values...)` linearly scans `list` for each appended `v`
  (`for _, existing := range list`), called per child condition in `populateFromChild` — O(n²) in the number of
  distinct QC-statement OIDs.
- Impact: Negligible in practice — the lists are the small set of QC-statement/QC-type OIDs from an MRA
  `CertificateContentReferencesEquivalence` (single digits), not attacker-scalable. Not a real cost on realistic
  input; noted for completeness.
- Recommendation: None required; use a `map[string]struct{}` accumulator if the lists ever grow.

No unbounded-memory, O(n²)-on-attacker-input, or goroutine-leak findings in `tsl` PART 1. The parse-once/retain
model (parsing result cached once, `TLParsingCacheDTO` immutable after build) means repeated `TSL` XML re-parsing
is bounded; the shared `validation/job` caches that hold the retained DTOs are the container, out of this unit's
scope.

**STD (idioms, error handling, defer/cleanup, context, sync, doc comments, test quality)**

T18A-STD-001 — **Info** — Category: STD — `tsl/tl_validation_job.go:133` (`Summary`)
- Evidence: `func (j *TLValidationJob) Summary()` does an **unchecked** type assertion
  `adapter := j.GetSummary().(*tlValidationJobSummaryAdapter)` (`tl_validation_job.go:134`).
- Impact: Panics if `GetSummary()` ever returns a different concrete summary type. Today the builder always
  returns the adapter (via `BuildTyped`/`tlValidationJobSummaryAdapter`), so this is an internal invariant and not
  reachable by hostile input — but it is a latent panic and a brittle coupling to a private type name.
- Recommendation: Use the comma-ok form and return a wrapped `error` (or a defensive default) instead of an
  unchecked assertion, so a future change to the summary type degrades gracefully.

Otherwise clean. Error handling is idiomatic (`error` returns at the facade boundary; `defer`/`Close()` used
correctly in `document_with_sha2.go` `WriteTo`/`Save` and `abstract_parsing_task.go` `readDocument`); no goroutine
leaks; doc comments are consistently excellent (every deliberate Java→Go mapping and DEVIATION is spelled out). The
`Init…`/`Overrides` self-registration pattern (e.g. `sha2_file_cache_data_loader.go`, `abstract_parsing_task.go`,
`default_trusted_list_with_sha2_predicate.go`) is a clean, well-documented stand-in for Java virtual dispatch.

### Tool log

Run from `dss/`:

- `gofmt -l tsl` → no output (clean).
- `go vet ./tsl/...` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./tsl/...` → `0 issues`.

All tools available and passing; nothing to note as unavailable.

### Open questions

1. **MRA condition/replacement direction (verify against upstream Java).**
   `trust_service_equivalence_converter.go:120-142` maps `Condition = Apply(PointedParty)` and
   `ContentReplacement = qCStatementOids(PointingParty)`; the consumer
   (`validation/reports/diagnostic/xml_trust_service_provider_builder.go:~605`) enacts the replacement when
   `Condition.Check(cert)` holds — i.e. "if the cert matches the **Pointed** party, treat it with the **Pointing**
   party's OIDs." This is the *opposite* of the reading one might expect from ETSI MRA "equivalence" wording
   (Pointing-party cert granted Pointed-party equivalence), but it is a 1:1 port and I could not locate the
   upstream `TrustServiceEquivalenceConverter.java` in the workspace to confirm which party Java assigns to
   `getCondition()` vs `getContentReplacement()`. **Action:** confirm against the Java source; if Java assigns the
   same mapping, this is correct parity (no change); if Java is inverted, this is a trust-bypass divergence to fix.
2. **T18A-SEC-001 root cause lives in `dss/validation/job`** (out of this unit). The tsl-side trigger is the
   concurrent pivot fan-out. The fix is a `validation/job` concern (thread-safe `CachedEntry` transitions) — flag
   for that unit and add a `-race` multi-pivot test here.
3. **Panic-on-nil-`SchemeInformation` parity.** The Go port deliberately reproduces the upstream NPE for a TL
   missing `<SchemeInformation>` (`lotl_parsing_task.go` NOTE). On the main path this is contained and safe; on the
   pivot path it is the DoS in T18A-SEC-002. Confirm the intended behaviour for a *pivot* LOTL missing the element
   (error the pivot, not crash) and pin it with a fixture.

## trustedlist + jaxb (unit U19)

- Scope: `dss/trustedlist` (TSL/MRA facade + parsers) and `dss/trustedlist/jaxb` — the
  ETSI TS 119 612 Trusted List / LOTL JAXB-equivalent model types and the
  MRA (Mutual Recognition Agreement) extension types, plus the `Unmarshal`/`Marshal`
  entry points. This is the *structural* model layer: it parses a TSL document into a
  Go tree and writes it back. Trust-logic (signature validation, expiry, pivot chains,
  criterion enforcement) lives in `dss/tsl` (U18a/U18b/U21) and is out of scope here.
- Date: 2026-08-26
- Depth: deep
- Compatibility contract: interop with Java DSS 6.5.RC1. `dss/PORTING.md` conventions
  (in particular: staticcheck `ST*`/`S1*`/`QF*`/`SA1019` and legacy-crypto flags are
  deliberately disabled; JAXB porting follows the JAXB-schema port convention; no
  runtime XSD validation is performed) and `docs/compatibility/known-gaps.md` are
  treated as deliberate, not findings.

### Files read

In-scope (read fully):
- `trustedlist/trusted_list_utils.go` (39)
- `trustedlist/trusted_list_facade.go` (30)
- `trustedlist/mra_facade.go` (30)
- `trustedlist/mra_equivalence_context_parser.go` (27)
- `trustedlist/mra_status_parser.go` (26)
- `trustedlist/mra_utils.go` (24)
- `trustedlist/doc.go` (17)
- `trustedlist/jaxb/jaxb_model.go` (326)
- `trustedlist/jaxb/jaxb_tsl_root.go` (322)
- `trustedlist/jaxb/jaxb_mra.go` (281)
- `trustedlist/jaxb/jaxb_tsl_common.go` (242)
- `trustedlist/jaxb/jaxb_ecc.go` (129)
- `trustedlist/jaxb/doc.go` (85)
- `trustedlist/jaxb/jaxb_tsl_tsp.go` (54)
- `trustedlist/jaxb/jaxb_tsl_service.go` (54)
- `trustedlist/jaxb/jaxb_tslx.go` (36)

Read with budget (626 lines): `trustedlist/jaxb/jaxb_common.go` — read the first 300
lines, then `grep -n 'func '` for all signatures, then read the `MarshalXML`/
`UnmarshalXML`/type-assertion regions (lines 300–626) in full. Net: every function in
the file was read; only the long cross-schema-stand-in header comment (lines 1–70) was
skipped as it is provenance documentation, not logic.

Test files (coverage only, not read fully): `trustedlist/facade_test.go`,
`trustedlist/jaxb/xml_test.go`, `trustedlist/jaxb/xml_kat_test.go`.

Cross-checked against upstream Java and the Go stdlib where the port is not
self-evident: `TrustedListFacade.java`/`MRAFacade.java`/`MRAUtils.java`/
`MRAStatusParser.java`/`MRAEquivalenceContextParser.java` (specs-trusted-list), the
generated JAXB classes for `AnyType`/`ExtensionType`/`AdditionalInformationType`/
`BigInteger`, and the Go 1.27 `encoding/xml` decoder (`read.go:323`,
`maxUnmarshalDepth=10000`, `errUnmarshalDepth`).

### Findings

**SEC (top priority)**

T19-SEC-001 — **Medium** — Category: SEC — `trustedlist/jaxb/jaxb_common.go:604-610`
(`ExtensionType.UnmarshalXML`)
- Evidence: `x.Critical = a.Value == "true" || a.Value == "1"` — the `Critical`
  attribute is read from `start.Attr` with no presence check; a missing attribute
  leaves `x.Critical` at its zero value `false`. The XSD declares it
  `<xsd:attribute name="Critical" type="xsd:boolean" use="required"/>`.
- Impact: `Critical` is trust-relevant — it is consumed downstream at
  `tsl/trust_service_converter.go:94` →
  `model.ConditionForQualifiers.Critical`, the flag that decides whether a QC
  qualifier condition is *mandatory* (critical) or *informational*. A hostile TSL
  that omits the required `Critical` on a `SchemeExtensions/Extension` carrying a QC
  qualifier would be parsed as `Critical=false` (non-critical) rather than rejected,
  silently downgrading the qualifier from mandatory to informational. This is not a
  standalone false-accept: the TSL's own signature is validated by `TLValidatorTask`
  (U18a scope) before its tree is consulted, so the attacker must already possess a
  trusted anchor. It is a lenient-parse gap on a trust-relevant field, consistent with
  the documented "no runtime XSD validation" stance but worth tightening.
- Recommendation: In `ExtensionType.UnmarshalXML`, error (or set a sentinel) when the
  required `Critical` attribute is absent, rather than defaulting to `false`; or add a
  `// DIVERGENCE, deliberate:` note + DESIGN entry recording that a missing required
  attribute is tolerated as `false` and why that is safe given the upstream signature
  validation. Flag for cross-check with the U18a consumer.

T19-SEC-002 — **Info** — Category: SEC — `trustedlist/jaxb/jaxb_tsl_common.go:142-169`
(`AdditionalInformationType.UnmarshalXML`)
- Evidence: the `default:` branch of the child-element switch calls `d.Skip()` for any
  child that is neither `TextualInformation` nor `OtherInformation`, silently
  discarding it. The XSD declares a strict
  `<xsd:choice maxOccurs="unbounded">` of exactly those two elements.
- Impact: A hostile TSL can embed arbitrary (potentially large) subtrees in
  `OtherTSLPointer.AdditionalInformation` and have them silently dropped. Not a DoS —
  `d.Skip()` is linear in the discarded subtree's size, bounded by the input — and not
  a false-accept (the discarded content is a free-form annotation, not trust-relevant
  per the schema). Consistent with the documented no-XSD-validation stance.
- Recommendation: None required; if XSD strictness is ever added, this `default:`
  should become an error rather than a `Skip`.

T19-SEC-003 — **Info** — Category: SEC — `trustedlist/jaxb/jaxb_common.go:499-534`
(`unmarshalAnyContent`) & `:184-226` (`foreignContent.unmarshal`)
- Evidence: the recursive decode cycle
  `AnyType → OtherInformation → AdditionalServiceInformation → OtherInformation → …`
  has **no depth limit in the ported code**. Empirically verified (throwaway probe,
  deleted after): a 1 000-deep chain parses; a 10 000-deep chain returns
  `exceeded max depth`; 100 000-deep likewise. The bound is Go 1.27's stdlib
  `encoding/xml` constant `maxUnmarshalDepth = 10000`
  (`GOROOT/src/encoding/xml/read.go:323`, `errUnmarshalDepth`), not a library guard.
- Impact: No DoS — the recursion is bounded by the toolchain and fails cleanly (error,
  not panic/OOM). The observation is that the "bounded recursion" property is
  *inherited* from a stdlib `const` that is not a documented API (it is already
  wasm-tuned to 5000) — a Go release that changes it would change this library's max
  parseable nesting. The module pins `go 1.27.0` in `go.mod`, so this is currently
  stable; it is an undocumented dependency, not a defect.
- Recommendation: Optionally add an explicit depth guard in `unmarshalAnyContent` /
  `foreignContent.unmarshal` (e.g. a counter capped at a documented value) so the bound
  is owned by the library rather than inherited; or record the stdlib dependency in the
  package `doc.go`. Low priority; the current behaviour is safe.

**PERF**

T19-PERF-BIGO-001 — **Low** — Category: PERF-BIGO — `trustedlist/jaxb/jaxb_common.go:415-433`
(`wildcardByName` / `wildcardPrefix`)
- Evidence: both functions do a **linear scan** of the 43-entry `wildcardElements`
  table (`for _, w := range wildcardElements { if w.namespace == space && … }`) on
  every recognized wildcard child during `unmarshalAnyContent` (`:505`) and
  `marshalAnyItems` (`:545`).
- Impact: For a TSL with N wildcard-dispatched children, dispatch is O(N·43). 43 is a
  small constant, so this is O(N) in practice — negligible on realistic input (the EU
  LOTL carries a few dozen such children). Not exploitable; noted for completeness.
- Recommendation: Optional — a `map[[2]string]wildcardEntry` (namespace+local → entry)
  would make dispatch O(1) per child; not worth it at 43 entries.

T19-PERF-MEM-001 — **Info** — Category: PERF-MEM — `trustedlist/jaxb/jaxb_common.go:184-226`
(`foreignContent.unmarshal`)
- Evidence: every `RawWildcardElement` / `ObjectIdentifierType` / `dsigSignature` /
  `dsigKeyValue` / `xadesAnyType` captured via `foreignContent` retains **every** token
  of its subtree in `c.Tokens []xml.Token` (each copied via `t.Copy()` /
  `xml.CopyToken`) for the lifetime of the parsed tree.
- Impact: On a large TSL (the EU LOTL is ~1 MB and carries many `Qualifications` /
  `OtherCriteriaList` / `ObjectIdentifierType` subtrees), the in-memory tree is a
  superset of the input — the dominant memory cost of `Unmarshal`. Bounded by input
  size, no leak, no unbounded growth: the tree is freed when the caller drops the
  `*TrustStatusListType`. The parse-once/retain model is by design (the tree is handed
  to the `tsl` parsing tasks and the validation-report builders).
- Recommendation: None required. If very large TSLs become a concern, a streaming
  decode that does not retain the foreign-content token streams would reduce the peak,
  at the cost of the byte-parity round-trip that the KAT corpus pins.

No PERF-LEAK findings: the parse-once/retain model means no retained references escape
the tree; there are no goroutines, no caches, and no global mutable state in this unit.

**STD**

T19-STD-001 — **Low** — Category: STD — `trustedlist/jaxb/jaxb_common.go:141-147`
(`BigInteger.UnmarshalText`)
- Evidence: `b.Int.SetString(s, 10)` accepts a leading `-` and `+`; the XSD types bound
  to `BigInteger` are `xs:integer` / `xs:positiveInteger` /
  `xs:nonNegativeInteger` (per the type's own doc comment), which forbid negatives for
  `TSLVersionIdentifier`, `TSLSequenceNumber`, `HistoricalInformationPeriod`,
  `MRADepth`, `technicalType`, `version`.
- Impact: A hostile TSL carrying a negative `TSLVersionIdentifier` /
  `TSLSequenceNumber` / `HistoricalInformationPeriod` would parse successfully and be
  carried into the trust model as a negative `*big.Int`. Schema-invalid but accepted.
  Not a standalone false-accept (the TSL's signature is validated elsewhere, U18a), but
  a lenient-parse gap on a trust-relevant numeric field — the sequence-number ordering
  logic in `tsl` (U18a/U21) is the likely consumer and may or may not re-validate sign.
- Recommendation: In `BigInteger.UnmarshalText`, reject values that do not satisfy the
  bound of the specific `xs:*Integer` type the field is bound to (or at minimum reject
  a leading `-`/`+` for the `positiveInteger`/`nonNegativeInteger` positions); or add a
  `// DIVERGENCE, deliberate:` note. Flag for cross-check with the U18a/U21 sequence-
  number consumer.

T19-STD-002 — **Info** — Category: STD — `trustedlist/trusted_list_utils.go:25-40` &
`trustedlist/mra_utils.go:15-24`
- Evidence: five XSD files (~64 KB total: `ts_119612v020401_xsd.xsd` 28 KB,
  `mra_schema_v2.xsd` 12 KB, `mra_schema_v2_19612v020401.xsd` 12 KB,
  `ts_119612v020401_sie_xsd.xsd` 8 KB, `ts_119612v020401_additionaltypes_xsd.xsd` 4 KB)
  are `go:embed`ded into the binary. The file headers explicitly state "no runtime XSD
  validation is performed" and "no JAXBContext is built from it"; the only consumer is
  `TestEmbeddedSchemas`.
- Impact: The embedded bytes inflate the module binary by ~64 KB for
  provenance/documentation only. Deliberate and documented ("exist for
  provenance/documentation and for any future validator"). No functional defect.
- Recommendation: None required. If binary size matters, the XSDs could move to
  `testdata/` (test-only) with the provenance note kept in the file header; the
  current choice is defensible.

T19-STD-003 — **Info** — Category: STD — `trustedlist/jaxb/jaxb_mra.go:39-66` &
`trustedlist/jaxb/jaxb_ecc.go:39-66`
- Evidence: the four enum-attribute adapters (`mraStatusAttr`,
  `mraEquivalenceContextAttr`, `keyUsageBitAttr`, `assertAttr`) each do a **linear
  scan** of their enum's value list on every `UnmarshalXMLAttr`. The lists are small
  (2–7 values), so this is O(1) in practice. The scan is also duplicated on the
  `dss/trustedlist` side (`MRAStatusParser`/`MRAEquivalenceContextParser`) per the
  `jaxb_mra.go` header — a deliberate, documented duplication to avoid an import cycle
  (this package sits below `dss/trustedlist`, which wraps it into the Facade API).
- Impact: Negligible cost; the duplication is a documented, justified consequence of
  the package layering. No functional defect.
- Recommendation: None required. If the enum lists ever grow, a `map[string]T` lookup
  would make each adapter O(1); not worth it at 2–7 values.

Verified clean (not findings):
- `Base64Binary.UnmarshalText` (`jaxb_common.go:85-98`) strips only space/tab/CR/LF
  before a strict `base64.StdEncoding.DecodeString` — correct for `xs:base64Binary`;
  no padding or non-alphabet injection path.
- `BigInteger.UnmarshalText` uses base 10 (not base 0) to avoid the octal-prefix trap
  of `big.Int.UnmarshalText` — correct for `xs:integer` (the sign gap is T19-STD-001).
- All `New…` constructors (`NewBigInteger`, `NewBigIntegerFromInt64`) are nil-safe;
  `BigInt()` returns nil for a nil receiver.
- No unchecked type assertions in this unit (unlike `tsl/tl_validation_job.go:134`,
  which U18a flags as T18A-STD-001). The `AnyItem.Elem` field is `any` and is
  dispatched by `ElemName`, not by a type assertion.
- No goroutines, no global mutable state, no `sync` primitives — the unit is a pure
  value-type model with deterministic `Marshal`/`Unmarshal`.

### Tool log

Run from `dss/`:

- `gofmt -l trustedlist` → no output (clean).
- `go vet ./trustedlist/...` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./trustedlist/...` → `0 issues`.

All tools available and passing; nothing to note as unavailable.

### Open questions

1. **T19-SEC-001 (missing `Critical`) — cross-check with the U18a consumer.**
   `tsl/trust_service_converter.go:94` passes `extensionType.Critical` to
   `toConditionForQualificationsType`, which feeds
   `model.ConditionForQualifiers.Critical` — the mandatory/informational flag on a QC
   qualifier condition. Confirm whether the U18a/U21 consumer re-validates the
   attribute's *presence* (not just its value) before enacting the condition; if it
   does, T19-SEC-001 is a non-issue; if it does not, the lenient default is a
   trust-relevant gap to tighten.
2. **T19-STD-001 (negative `BigInteger`) — cross-check with the U18a/U21 sequence-
   number consumer.** `TSLSchemeInformationType.TSLVersionIdentifier` /
   `TSLSequenceNumber` / `HistoricalInformationPeriod` are `*BigInteger` and feed the
   sequence/ordering logic in `tsl`. Confirm whether that logic re-validates sign; if
   it does, T19-STD-001 is a non-issue; if it does not, a negative version/sequence
   number in a hostile TSL would be carried into the trust model.
3. **T19-SEC-003 (stdlib depth bound) — toolchain dependency.** The recursion bound is
   Go 1.27's `maxUnmarshalDepth = 10000` (a `const`, not a documented API). The module
   pins `go 1.27.0` in `go.mod`, so this is currently stable. If the project ever
   supports a range of Go versions, add an explicit depth guard in
   `unmarshalAnyContent` / `foreignContent.unmarshal` so the bound is owned by the
   library rather than inherited.

### Summary

- Total findings: **8** (1 Medium, 2 Low, 5 Info). No Critical, no High.
- SEC: 3 (1 Medium, 2 Info) — one lenient-parse gap on a trust-relevant required
  attribute (`Critical`), one strict-choice tolerance, one toolchain-inherited depth
  bound (all consistent with the documented no-XSD-validation stance; none a
  standalone false-accept).
- PERF: 2 (1 Low, 1 Info) — one small O(N·43) dispatch, one bounded retain-by-design
  memory cost. No PERF-LEAK.
- STD: 3 (1 Low, 2 Info) — one sign gap on a trust-relevant `xs:*Integer` field, one
  deliberate ~64 KB binary bloat, one documented O(1) enum-scan duplication.
- Tool log: clean (gofmt, go vet, golangci-lint all pass with zero findings in-scope).

## tsl PART 2 (unit U18b)

- Scope: `dss/tsl` — TSL/LoTL parsing & validation layer, **PART 2** = the last 25 non-test files of
  `dss/tsl` (all ≤28 lines, read in full). PART 1 (the other 48 files) is the `U18a` section above; its findings
  are not duplicated here.
- Date: 2026-08-26
- Depth: deep
- Compatibility contract: interop with Java DSS 6.5.RC1. `dss/PORTING.md` conventions and
  `docs/compatibility/known-gaps.md` treated as deliberate, not findings. staticcheck `ST*`/`S1*`/`QF*`/`SA1019`
  and legacy-crypto flags are deliberately disabled (PORTING.md) and were not reported.
- Upstream cross-checked against `~/Workspace/esig/dss` (DSS 6.5.RC1): `TrustServiceEquivalenceConverter.java`
  (MRA direction, the PART-1 open question), `SchemeInformationURIByLang.java`, `PivotSchemeInformationURI.java`,
  `NonEmptyTSPInformation.java`, `NonEmptyTrustService.java`, `NonEmptyServiceInformation.java`,
  `TradeNamePredicate.java`, `EUTLOtherTSLPointer.java`, `LogTLSignatureErrorAlertHandler.java`,
  `MRANamespace.java`, and the Java consumer `XmlTrustServiceProviderBuilder.java:684-700`.

### Files read (25, in-scope — all read fully)

`tl_parsing_result.go` (28), `other_tsl_pointer_predicate.go` (28), `tl_cache_access_by_key.go` (27),
`pivot_scheme_information_uri.go` (27), `scheme_information_uri_by_lang.go` (26),
`log_tl_signature_error_alert_handler.go` (25), `log_tl_parsing_error_alert_handler.go` (25), `doc.go` (25),
`cache_cleaner.go` (25), `tl_parsing_error_detection.go` (24), `lotl_alert.go` (24), `tl_alert.go` (23),
`non_empty_trust_service.go` (23), `trade_name_predicate.go` (20), `non_empty_tsp_information.go` (19),
`non_empty_service_information.go` (19), `eu_lotl_other_tsl_pointer.go` (16), `xml_other_tsl_pointer.go` (15),
`pdf_other_tsl_pointer.go` (15), `eu_tl_other_tsl_pointer.go` (15), `trust_anchor_period_predicate.go` (14),
`official_journal_scheme_information_uri.go` (13), `trust_service_provider_predicate.go` (12),
`trust_service_predicate.go` (12), `mra_namespace.go` (9).

Supporting (read for context, ≤ budget): `validation/job/abstract_cache_access_by_key.go` (constructor + field
layout, ~90 lines), `xml/common/dss_namespace.go` (immutability of the `MRANamespaceNS` global), and the two base
predicates `type_other_tsl_pointer.go` / `mimetype_other_tsl_pointer.go` (embedded by value by the EU/XML/PDF
pointer types). `trust_service_equivalence_converter.go:100-160` re-read for the MRA cross-check (PART-1 file, not
re-reported).

### Findings

**SEC (top priority)**

T18B-SEC-001 — **Info** — Category: SEC — `tsl/scheme_information_uri_by_lang.go:24`
- Evidence: `func (p *SchemeInformationURIByLang) Test(schemeInformationURI *jaxb.NonEmptyMultiLangURIType) bool
  { return p.lang == schemeInformationURI.Lang }` — dereferences `schemeInformationURI` with **no nil guard**, so a
  nil argument panics. Sibling `pivot_scheme_information_uri.go:24` guards the same shape
  (`if t != nil && t.Value != ""`).
- Impact: A nil `NonEmptyMultiLangURIType` would panic (DoS) rather than return `false`. Verified this is a
  faithful 1:1 port of upstream: Java `SchemeInformationURIByLang.test` is `lang.equals(schemeInformationURI.getLang())`
  (NPEs on null, no guard) whereas Java `PivotSchemeInformationURI.test` is `if (t != null && t.getValue() != null)`
  (guarded). The asymmetry is a reproduction of upstream, not a porting defect. In practice the argument comes from a
  present `<SchemeInformationURI>` list element (non-nil after unmarshal), and any parse-path panic is contained by
  `AbstractRunnableAnalysis.Run`'s `recover` (PART-1). This is a *filter* predicate, so the worst case is a contained
  panic or a dropped match — **never a false "trusted" verdict** (a predicate that misses an entry removes it, it does
  not grant it).
- Recommendation: None required for parity. If defensive hardening is later desired (beyond upstream behaviour), it
  must be recorded as a `// DIVERGENCE, deliberate:` + DESIGN entry, since upstream NPEs here.

All other in-scope files verified clean on the mandated SEC lenses:

- **No false-accept in any predicate.** Every filter in this part is a *removal* filter — it returns `false` to drop
  non-conforming entries: `NonEmptyTrustService` (`t.TSPServices != nil && IsCollectionNotEmpty(...TSPService)`),
  `NonEmptyTSPInformation` (`t.TSPInformation != nil`), `NonEmptyServiceInformation`
  (`ServiceInformation != nil && StatusStartingTime != nil`), `TradeNamePredicate`
  (`t != "" && !registrationIdentifier.Test(t)`), the `EULOTLOtherTSLPointer`/`EUTLOtherTSLPointer` (exact `TSLType`
  URI equality), and the `XMLOtherTSLPointer`/`PDFOtherTSLPointer` (exact `MimeType` equality, case-insensitive).
  A bug in any of these causes a false-*reject* (a service/pointer dropped), not a false-*accept*. Wire values
  (the ETSI `TSLType` URIs and MIME types) are byte-verbatim from upstream — confirmed against the Java constants
  (`EUTLOtherTSLPointer.EXPECTED_EU_TL_TYPE`, `XMLOtherTSLPointer.EXPECTED_MIMETYPE`,
  `PDFOtherTSLPointer.EXPECTED_MIMETYPE`, `EULOTLOtherTSLPointer.EXPECTED_EU_LOTL_TYPE`).
- **`other_tsl_pointer_predicate.go`** — `otherTSLPointerPredicateAnd.Test` is `a.first.Test(o) && a.second.Test(o)`,
  a correct short-circuit port of Java `Predicate#and`. No logic defect.
- **`tl_parsing_error_detection.go:23`** — `Detect` returns true only on `IsError()` or non-empty
  `StructureValidationMessages()`; it *surfaces* an error condition (conservative), never suppresses one.
- **`tl_parsing_result.go`** — pure data holder (getter/setter for `trustServiceProviders`); no trust logic.
- **`mra_namespace.go:9`** — `var MRANamespaceNS = common.NewDSSNamespace(...)` is global state, but `DSSNamespace`
  exposes only getters (`Uri`/`Prefix`/`IsSameUri`/`String`) with no setters, so the global is effectively immutable.
  Clean port of Java `MRANamespace.NS` (`public static final`). No mutable global-state hazard.
- **Alert handlers** (`log_tl_signature_error_alert_handler.go`, `log_tl_parsing_error_alert_handler.go`) — log-and-
  return-`nil`; match Java `void process(TLInfo) { LOG.warn(...) }`. They never alter a verdict.
- **Cache/alert wrappers** (`tl_cache_access_by_key.go`, `lotl_alert.go`, `tl_alert.go`, `cache_cleaner.go`) — thin
  embeds of `validation/job` types with no added logic. `TLCacheAccessByKey` embeds
  `job.AbstractCacheAccessByKey[...]` **by value**, the same documented pattern `tl_parsing_cache_dto.go` established
  in PART-1; the outer type is always used by pointer and the inner state lives in shared `*Cache` pointers, so no
  copy-escape. The underlying cache-entry concurrency hazard is PART-1's `T18A-SEC-001` (a `validation/job` concern),
  not re-reported here.

**PERF (parsing cost, allocations, leaks)**

No findings in this lens. All 25 files are O(1) value types, predicate `Test` methods, or thin wrappers: no loops,
no allocations on the hot path, no goroutines, no caches created or retained, no `defer`/stream plumbing, and no
attacker-scalable input. The only retained state is the `TLParsingResult.trustServiceProviders` slice (bounded by the
TL's own content, held once) and the immutable `MRANamespaceNS` constant. No unbounded-memory, O(n²)-on-attacker-input,
or goroutine-leak concerns.

**STD (idioms, error handling, defer/cleanup, context, sync, doc comments)**

T18B-STD-001 — **Info** — Category: STD — `tsl/tl_parsing_result.go:26`
- Evidence: `func (r *TLParsingResult) TrustServiceProviders() []*tslmodel.TrustServiceProvider { return
  r.trustServiceProviders }` — returns the backing slice directly, so a caller can mutate the internal list through
  the returned reference (mutable escape).
- Impact: Inherent to the 1:1 Java contract — upstream `getTrustServiceProviders()` returns the raw `List` (also
  mutable-escaping). In this codebase the consumer (the qualification/report builders) treats it read-only, so no
  behavioural defect. Noted only as the one place in PART 2 where internal mutable state is exposed by reference.
- Recommendation: None required for parity; a defensive `append([]..., r.trustServiceProviders...)` copy would
  diverge from upstream's raw-`List` return and is not recommended without a DESIGN entry.

Otherwise clean. Error handling is idiomatic for the types present: the two log alert handlers correctly return `nil`
for a log-only `Process`; every `New…` constructor that guards a required argument panics with the Java
`Objects.requireNonNull` message (e.g. `NewSchemeInformationURIByLang` on empty `lang`, matching upstream
`requireNonNull`), which is the PORTING-mandated form and not reachable by hostile input on the normal path. Doc
comments are consistently precise (every Java→Go mapping and the deprecated-`CacheCleaner` provenance is spelled
out). No context, sync, or defer issues — none of the 25 files uses goroutines or `context`.

### Cross-checks

- **MRA condition/replacement direction — VERDICT: confirmed-correct (parity), NOT inverted.** This resolves
  PART-1 open question 1. Go `certificateEquivalence` (`trust_service_equivalence_converter.go:140-148`):
  `expected := CertificateContentDeclarationPointedParty`; `substitute :=
  CertificateContentDeclarationPointingParty`; `condition := Apply(substitute)`; then
  `SetCondition(Apply(expected))` and `SetContentReplacement(qCStatementOids(condition))`. Java
  `TrustServiceEquivalenceConverter.getCertificateEquivalence` does the identical mapping:
  `expected = getCertificateContentDeclarationPointedParty()`; `substitute =
  getCertificateContentDeclarationPointingParty()`; `condition = criteriaConverter.apply(substitute)`;
  `equiv.setCondition(criteriaConverter.apply(expected))`;
  `equiv.setContentReplacement(getQCStatementOids(condition))`. Both consumers gate identically — Go
  `xml_trust_service_provider_builder.go:616` `if condition.Check(certToken)` and Java
  `XmlTrustServiceProviderBuilder.java:684-688` `if (condition.check(certToken))` — then apply
  `contentReplacement`. So the semantic is, in both: "if the cert matches the **Pointed** party's criteria, replace
  its QC-statement OIDs with the **Pointing** party's." The mapping the PART-1 reviewer found counter-intuitive is
  exactly what upstream Java assigns. **No trust-bypass divergence; the port is byte-for-byte faithful. No change
  required.**
- **`SchemeInformationURIByLang` nil-guard asymmetry — confirmed parity.** Java `SchemeInformationURIByLang.test`
  has no null guard (NPEs on null); Java `PivotSchemeInformationURI.test` guards. The Go port reproduces both exactly
  (see T18B-SEC-001). Not a divergence.
- **Wire values — confirmed verbatim.** The ETSI `TSLType` URIs and MIME types in the four pointer predicates match
  the Java `EXPECTED_*` constants byte-for-byte (checked `EUTLOtherTSLPointer`, `EULOTLOtherTSLPointer`,
  `XMLOtherTSLPointer`, `PDFOtherTSLPointer`).
- **`MRANamespaceNS`** URI/prefix match Java `MRANamespace.NS` verbatim.

### Tool log

Run from `dss/`:

- `gofmt -l tsl` → no output (clean).
- `go vet ./tsl/` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./tsl/...` → `0 issues.`

All tools available and passing; nothing in-scope was reported, so no deliberate-exclusion notes were required.

### Summary counts

Total findings: **2** — Critical 0 / High 0 / Medium 0 / Low 0 / **Info 2**
(SEC 1 · PERF 0 · STD 1). No false-accept, no leak, no O(n²)-on-attacker-input, no trust-logic error in PART 2.
The MRA direction (PART-1 open question 1) is resolved **confirmed-correct**.

## token (unit U20)

- Scope: `dss/token` — signature-token abstractions: HSM/PKCS#11/PKCS#12/JKS/Apple/MSCAPI
  key access, digest-info encoding, password callbacks, and key-entry selection predicates.
  A 1:1 port of Java DSS 6.5.RC1 `dss-token`; most classes are thin adapters over crypto
  APIs, so the focus is the shared plumbing (sign dispatch, PFX→keyStore correlation,
  `DigestInfoEncoder`, predicate wiring).
- Date: 2026-08-26
- Depth: deep
- Compatibility contract: interop with Java DSS 6.5.RC1. `dss/PORTING.md` conventions and
  `docs/compatibility/known-gaps.md` are treated as deliberate, not findings. staticcheck
  `ST*`/`S1*`/`QF*`/`SA1019` and legacy-crypto flags (SHA-1, MD5, PKCS#1 v1.5, DSA) are
  deliberately disabled for a signing/validation library (PORTING.md) and were not reported.

### Files read (22, in-scope — all read fully, every file ≤305 lines)

`key_store_signature_token_connection.go` (305), `abstract_signature_token_connection.go`
(257), `pkcs11_signature_token.go` (183), `digest_info_encoder.go` (148),
`abstract_key_store_token_connection.go` (145), `ks_private_key_entry.go` (133),
`extended_key_usage_key_entry_predicate.go` (76), `jks_signature_token.go` (64),
`pkcs12_signature_token.go` (60), `signature_token_connection.go` (44),
`apple_signature_token.go` (44), `password_protection.go` (42), `mscapi_signature_token.go`
(41), `prefilled_password_callback.go` (33), `key_usage_key_entry_predicate.go` (30),
`valid_at_time_key_entry_predicate.go` (24), `dss_private_key_entry.go` (21),
`dss_private_key_access_entry.go` (20), `sun_pkcs11_initializer.go` (18),
`password_input_callback.go` (12), `doc.go` (12), `all_key_entry_predicate.go` (11),
`dss_key_entry_predicate.go` (10).

Cross-checked against the Go 1.27 stdlib and upstream Java where the port is not self-
evident: `crypto/rsa` `Sign`/`SignPKCS1v15` (`GOROOT/src/crypto/rsa/fips.go:335`, the
`hash == crypto.Hash(0)` → "hashed is signed directly" raw path), `crypto/ecdsa`
`Signer`/`SignerOpts` pre-hashed contract, the `enumerations.SignatureAlgorithm` →
`(EncryptionAlgorithm, DigestAlgorithm)` table, `enumerations.DigestAlgorithm.SaltLength`,
`spi.CertificateExtensionsUtilsExtendedKeyUsage` nil-return semantics, `internal/pfx.Store`
shape, and upstream `dss-token/.../digest/DigestInfoEncoder.java` `isEncoded` (byte-for-byte
parity) plus `AppleSignatureToken.java:60` / `MSCAPISignatureToken.java:57` (the `"nimp"`
literal). Test files (coverage only): `pkcs12_signature_token_test.go`,
`digest_info_encoder_test.go`, `key_entry_predicate_test.go`, `unsupported_platforms_test.go`.

### Findings

**SEC (top priority)**

T20-SEC-001 — **Medium** — Category: SEC — `dss/token/key_store_signature_token_connection.go:201`
(`used := make([]bool, len(certs))`; leaf claim at `:229`; chain walk at `:247–248`)
- Evidence: a single shared `used []bool` is allocated once, *before* the
  `for i, keyBag := range store.PrivateKeys` loop (`:201`). After each key's leaf is chosen,
  `used[leafIndex] = true` (`:229`), and the chain walk marks every appended cert
  `used[nextIndex] = true` (`:248`). The *next* key's chain walk then skips every cert
  already marked `used` — including shared intermediates/roots — because it only looks at
  `!used[ci]` (`:238`).
- Impact: For a PFX holding **≥2 private keys whose chains share an intermediate or root**
  (a common real-world shape: two keys from the same CA), the 2nd-and-later key's
  `CertificateChain()` is **truncated to the leaf** — the shared chain certs are "consumed"
  by the first key and become invisible to the rest. `KSPrivateKeyEntry.CertificateChain()`
  is a wire-relevant value (it feeds the certificate chain embedded in CAdES/PAdES/XAdES
  output), so an incomplete chain yields a signature missing its intermediate/root. The
  leaf and private key per entry are still correct, so this is **not** a wrong-key or
  false-accept — it is a wrong/incomplete *chain* selection. Secondary: the "first unclaimed
  cert" leaf fallback (`:218–224`) is ambiguous for a multi-key PFX whose bags omit
  `localKeyId`, and can claim another key's leaf as this key's. Neither path is exercised by
  the in-module fixtures (all single-key: `rsa_test.p12`, `ec_test.p12`, `user_a_rsa.p12`),
  so it is untested.
- Recommendation: Track leaf-claiming separately from chain membership. Allocate a
  `claimedLeaf` set and mark only the leaf per key; walk each key's chain independently
  against the *full* certificate list (allowing the same intermediate/root to appear in
  multiple chains). Add a multi-key, shared-intermediate PFX fixture asserting every key's
  `CertificateChain()` is the full leaf→root chain. The file header calls this an
  approximation of `java.security.KeyStore`'s PKCS#12 provider — verify the expected
  per-key chain shape against upstream if parity matters. (Upgrade to High if a multi-key
  shared-chain PFX is a supported signing workload, since `CertificateChain()` is
  wire-relevant.)

**PERF (algorithmic complexity, memory, leaks)**

No findings in this lens.

- **PERF-BIGO:** no O(n²)-on-attacker-input. `pkcs12BuildKeyStore`'s correlation is
  O(K·C) (K keys × C certs); `keyStore.entry`/`aliases` are linear in the entry count. All
  are bounded by the PFX's own size, and a PFX is caller-supplied secret material, not
  hostile network input — not attacker-scalable.
- **PERF-MEM / PERF-LEAK:** the `keyStore` retains decrypted private keys
  (`crypto.Signer`) for the lifetime of the connection — **expected** for an in-memory
  signing token (the key must persist to sign repeatedly), and it is not copied
  unnecessarily (`KSPrivateKeyEntry.PrivateKey()` returns the stored signer;
  `NewCertificateToken` wraps the parsed `*x509.Certificate` without re-encoding the key).
  `PasswordProtection.Destroy()` zeroes the password bytes on `Close()`. No goroutines, no
  unbounded growth, no leak.

**STD (idioms, error handling, cleanup, doc comments)**

T20-STD-001 — **Low** — Category: STD — `dss/token/pkcs11_signature_token.go:135,165,172`
(`pkcs11SignatureTokenBuildConfig`, `pkcs11SignatureTokenEscapePath`,
`pkcs11SignatureTokenRandomUUID`)
- Evidence: these three free functions port `Pkcs11SignatureToken`'s `buildConfig` /
  `escapePath`, but every `Pkcs11SignatureToken` constructor and `KeyStore()` return
  `pkcs11SignatureTokenUnsupported()`, so no live path calls them. The only consumer is
  `TestPkcs11SignatureTokenBuildConfig` (`unsupported_platforms_test.go:51`).
- Impact: Dead code retained for upstream fidelity (documented in the file header: "still
  ported below as free functions for fidelity, even though nothing in this port can consume
  their output"). No functional or security impact; a minor maintainability cost (a
  `rand.Read`-based v4-UUID generator and a path-escaper that nothing reachable can use).
- Recommendation: Acceptable as-is with the existing fidelity note. If dead-code reduction is
  preferred, move the three functions into the test file as their sole consumer and drop
  them from the non-test build. No action required.

Verified clean on the mandated SEC/STD lenses (each confirmed, not assumed):

- **RSA pre-DigestInfo path — correct, not double-wrapped.** `prepareMessage`
  (`abstract_signature_token_connection.go`) hashes then `DigestInfoEncoderEncode`s for
  plain RSA, and `abstractSignatureTokenConnectionSign` signs with `crypto.Hash(0)`.
  Verified against `GOROOT/src/crypto/rsa/fips.go:345` ("If hash is zero, hashed is signed
  directly. This isn't advisable except for interoperability") — so the already-wrapped
  DigestInfo is signed verbatim, exactly what JCA `<digest>withRSA` builds internally. No
  double DigestInfo, no false signature.
- **ECDSA pre-hashed path — correct, no double-hash.** Every non-raw ECDSA/PLAIN_ECDSA
  digest (SHA1/224/256/384/512, SHA3-*, RIPEMD160) resolves in the `PSSHashes` map or the
  explicit RIPEMD160 branch, so the named hash is length-check-only and never re-hashes
  `preparedInput`; `ECDSA_RAW` passes `nil` opts. `ED25519`/`ED448` bypass the hash
  entirely. Matches JCA `NONEwithECDSA`.
- **RSASSA-PSS — correct.** `prepareMessage` does **not** pre-wrap for PSS (raw digest
  passed), and `Sign` uses `rsa.PSSOptions{Hash: <map>, SaltLength:
  pssDigestAlgorithm.SaltLength()}`, the same two values upstream `initParameters` threads.
  Unsupported digests (e.g. RIPEMD160) are refused, not silently signed.
- **`DigestInfoEncoder` — byte-for-byte parity + hardening.** `Encode` builds the identical
  DER `SEQUENCE{ AlgorithmIdentifier{ OID, NULL }, OCTET STRING }` as upstream. `IsEncoded`
  is a statement-for-statement port of Java `isEncoded` (single-byte-length assumption and
  all included — confirmed against `DigestInfoEncoder.java:158–231`), with an added Go
  `defer recover()` that maps a Go index panic to the same `false` Java's `catch(Exception)`
  yields — a strict improvement, not a divergence.
- **`"nimp"` password (Apple/MSCAPI) — parity, not a secret leak.** Confirmed verbatim in
  upstream `AppleSignatureToken.java:60` / `MSCAPISignatureToken.java:57
  (new PasswordProtection("nimp".toCharArray()))`. Both tokens' `KeyStore()` always error,
  so the literal is never used to unlock anything. Deliberate fidelity; not reported.
- **EKU predicate panic (`extended_key_usage_key_entry_predicate.go:71–72`) — deliberate.**
  Fires only when a cert *carries* an EKU extension that fails to parse
  (`CertificateExtensionsUtilsExtendedKeyUsage` returns nil only on `!parsed`/`!ok`; an
  *absent* extension returns an empty-but-valid list, so well-formed certs never hit it).
  Documented in the file header as preserving upstream's "exception surfaces from a
  predicate call" behaviour; signing-side (not a verification false-accept). Not a finding.
- **Constructor `panic`s** (`NewKSPrivateKeyEntry`, `NewExtendedKeyUsageKeyEntryPredicate`,
  `SetKeyEntryPredicate`, `ensureDigestUniform` guards, etc.) are `Objects.requireNonNull`
  parity per PORTING.md and not reachable by hostile input on the normal path.
- **No unchecked type assertions** in-scope (`key.(crypto.Signer)` / `key.(*dsa.PrivateKey)`
  both use the comma-ok form; `dssPrivateKeyEntry` is a type switch, not an assertion).
  **No goroutines, no global mutable state, no `sync`** in this package.
- **Password handling** — `PasswordProtection.Destroy()` zeroes and nils the `[]byte`;
  `keyStoreSignatureTokenConnectionPasswordString` copies to a `string` scoped to the
  `pfx.Load` call, so nothing retains the secret past the load. `PrefilledPasswordCallback`
  delegates correctly. No retention/leak.

### Tool log

Run from `dss/`:

- `gofmt -l token` → no output (clean).
- `go vet ./token/` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./token/...` → `0 issues`.
- `go test ./token/ -count=1` → `ok  github.com/ryftcore/dss-go/dss/token  0.782s`.

All tools available and passing; nothing in-scope was reported, so no deliberate-exclusion
notes were required.

### Open questions

1. **T20-SEC-001 (multi-key chain truncation) — verify the expected shape against
   upstream.** The file header frames `pkcs12BuildKeyStore` as "the closest a Go program
   can get to what `java.security.KeyStore`'s PKCS#12 provider does internally." Confirm
   whether upstream produces a *full* per-key chain (leaf→root) for every key entry even
   when chains are shared; if it does, the Go truncation is a divergence to fix (T20-SEC-
   001); if upstream also shares/consumes chain certs, document it as a `// DIVERGENCE,
   deliberate:` parity note. Add a multi-key, shared-intermediate PFX fixture either way.
2. **DSA `crypto/dsa` deprecation.** `pkcs12Signer` wraps `*dsa.PrivateKey` and
   `ksPrivateKeyEntryJCEName` accepts it (both `//nolint:staticcheck`). Deliberate (legacy
   key stores) and consistent with PORTING.md's "keep legacy algorithms a validation
   library must retain." Noted so the consolidation step does not mistake it for a
   staticcheck finding.

### Summary counts

Total findings: **2** — Critical 0 / High 0 / Medium **1** / Low **1** / Info 0.
SEC: 1 (1 Medium) — one wrong/incomplete *chain* selection in multi-key PFX correlation
(not a false-accept or wrong-key). PERF: 0 — no O(n²)-on-attacker-input, no leak (key
retention is by-design for an in-memory signing token). STD: 1 (1 Low) — dead PKCS#11
config helpers retained for fidelity. Tool log clean (gofmt, go vet, golangci-lint all
pass with zero in-scope findings; `go test ./token/` passes).
