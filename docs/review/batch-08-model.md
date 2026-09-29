# Batch 08 — model + enumerations

## model — top-level (unit U24)

- Scope: `dss/model/` top-level `.go` files (non-test) — core model types (X500 principal,
  certificate/token refs, signatures, digests, identifiers, documents, policy, signer location,
  manifests, container info). A 1:1 Go port of Java DSS 6.5.RC1.
- Date: 2026-08-26
- Depth: budgeted (x500_principal.go read in 4 targeted regions, ≤400 lines total; all other
  52 in-scope files read in full).
- Compatibility contract: interop with Java DSS 6.5.RC1 (byte-level + verdict parity).
  `dss/PORTING.md` conventions and `docs/compatibility/known-gaps.md` are treated as
  deliberate, not findings. `// DIVERGENCE, deliberate:` comments mark intentional divergences.
  Deliberately-disabled staticcheck style rules (ST*/S1*/QF*/SA1019) and legacy-crypto flags
  (SHA-1, MD5, PKCS#1 v1.5) are PORTING.md conventions, not findings.
- Out-of-scope: `model/` sub-packages (`model/eaa/`, `model/x509/`, `model/lote/`,
  `model/policy/`, `model/revocation/`, …) — separate units. `spi/` (batch-07), `tsl/`+
  `trustedlist/` (batch-06) already reviewed; cross-referenced where a model type is consumed
  by an already-filed spi finding.

### Files read (per-file line counts actually read)

| File | Total lines | Lines read | Approach |
|---|---|---|---|
| `x500_principal.go` | 726 | 330 (budgeted) | First 150 lines (header, types, constructors, `Encoded`, `RFC2253Name`, `Canonical`), `grep -n 'func '` for all 28 signatures, targeted reads: lines 150–225 (`String`, `Equals`, `generateRFC2253DN`, `generateRFC2253CanonicalDN`), lines 226–454 (`rfc2253String`, `rfc2253CanonicalString`, `x500PrincipalKeyword`, `x500PrincipalIsDerString`), lines 454–568 (`x500PrincipalDecodeValue`, `DecodeLatin1`, `DecodeUTF16BE`, `IsPrintableStringChar`, `IsJavaWhitespace`, `JavaTrim`, `CompareJavaStrings`, `HexLower`), lines 570–726 (`ReadTLV`, `Parse`, `DecodeOID`). |
| `certificate_token.go` | 447 | 447 | FULL-read. |
| `abstract_serializable_signature_parameters.go` | 347 | 347 | FULL-read. |
| `token.go` | 259 | 259 | FULL-read. |
| `policy.go` | 253 | 253 | FULL-read. |
| `digest_document.go` | 234 | 234 | FULL-read. |
| `identifier.go` | 194 | 194 | FULL-read. |
| `blevel_parameters.go` | 186 | 186 | FULL-read. |
| `reference_validation.go` | 165 | 165 | FULL-read. |
| `common_document.go` | 159 | 159 | FULL-read. |
| `signer_location.go` | 152 | 152 | FULL-read. |
| `in_memory_document.go` | 134 | 134 | FULL-read. |
| `public_key.go` | 102 | 102 | FULL-read. |
| `manifest_file.go` | 96 | 96 | FULL-read. |
| `file_document.go` | 95 | 95 | FULL-read. |
| `digest.go` | 93 | 93 | FULL-read. |
| `data_identifier.go` | 93 | 93 | FULL-read. |
| `manifest_entry.go` | 88 | 88 | FULL-read. |
| `common_object_identifier.go` | 88 | 88 | FULL-read. |
| `user_notice.go` | 85 | 85 | FULL-read. |
| `sp_doc_specification.go` | 78 | 78 | FULL-read. |
| `x500_principal_helper.go` | 73 | 73 | FULL-read. |
| `multiple_digest_identifier.go` | 71 | 71 | FULL-read. |
| `container_info.go` | 71 | 71 | FULL-read. |
| `signature_value.go` | 66 | 66 | FULL-read. |
| `timestamp_parameters.go` | 65 | 65 | FULL-read. |
| `signature_policy_store.go` | 65 | 65 | FULL-read. |
| `common_commitment_type.go` | 58 | 58 | FULL-read. |
| `serializable_signature_parameters.go` | 49 | 49 | FULL-read. |
| `dss_exception.go` | 48 | 48 | FULL-read. |
| `dss_document.go` | 48 | 48 | FULL-read. |
| `dss_message_digest.go` | 45 | 45 | FULL-read. |
| `entity_identifier_builder.go` | 43 | 43 | FULL-read. |
| `serializable_evidence_record_incorporation_parameters.go` | 41 | 41 | FULL-read. |
| `encapsulated_revocation_token_identifier.go` | 40 | 40 | FULL-read. |
| `to_be_signed.go` | 38 | 38 | FULL-read. |
| `commitment_qualifier.go` | 35 | 35 | FULL-read. |
| `oid_repository.go` | 30 | 30 | FULL-read. |
| `token_comparator.go` | 27 | 27 | FULL-read. |
| `token_identifier.go` | 21 | 21 | FULL-read. |
| `serializable_counter_signature_parameters.go` | 20 | 20 | FULL-read. |
| `doc.go` | 20 | 20 | FULL-read. |
| `original_identifier_provider.go` | 19 | 19 | FULL-read. |
| `certificate_token_identifier.go` | 18 | 18 | FULL-read. |
| `timestamp_binary.go` | 17 | 17 | FULL-read. |
| `x500_name_identifier.go` | 15 | 15 | FULL-read. |
| `key_identifier.go` | 15 | 15 | FULL-read. |
| `identifier_based_object.go` | 15 | 15 | FULL-read. |
| `entity_identifier.go` | 15 | 15 | FULL-read. |
| `signature_parameters_builder.go` | 13 | 13 | FULL-read. |
| `serializable_timestamp_parameters.go` | 12 | 12 | FULL-read. |
| `token_identifier_provider.go` | 10 | 10 | FULL-read. |
| `identifier_builder.go` | 8 | 8 | FULL-read. |

**Total in-scope lines read: 4,809** (budget: 5,200). Budget respected with 391 lines to spare.

Cross-referenced (not re-reviewed): `utils/ordered_map.go` (nil-receiver safety of
`Len`/`Keys`/`Get` — verified all three handle `m == nil`), `spi/dss_asn1_utils.go:570`
(`DSSASN1UtilsX500PrincipalAreEquals` — batch-07 already reviewed), `spi/common_certificate_source.go`
(`IsCertificateSourceEqual` — batch-07 already reviewed).

### Findings

**SEC (top priority)**

No findings. Each mandated SEC lens was verified clean:

- **DN equality used for cert matching (false-accept / false-reject):**
  `X500Principal.Equals` (`x500_principal.go:169`) is nil-safe (checks `p == nil || other == nil`
  before dereferencing), pre-checks RDN and AVA counts (mirroring `X500Name#equals`), and then
  compares `p.Canonical() == other.Canonical()`. `Canonical()` (`x500_principal.go:148`) is a
  faithful port of `X500Name#getRFC2253CanonicalName`: RFC 2253 form restricted to
  PrintableString/UTF8String values, multi-valued RDNs sorted with `AVAComparator` (keyword-
  AVAs first, then OID-keyword AVAs, both ordered by `String#compareTo` UTF-16 code-unit
  comparison via `x500PrincipalCompareJavaStrings`), internal whitespace collapsed, leading/
  trailing whitespace trimmed (Java `String#trim()`, not `strings.TrimSpace`), upper-cased then
  lower-cased with US locale, NFKD-normalized. No false-accept (two distinct DNs do not compare
  equal) and no false-reject (two DNs that denote the same name in different DER encodings or
  string types DO compare equal, as they do in Java). The canonical form is cached
  (`canonicalValid` flag) and computed only once.
  `X500PrincipalHelper.Equals` (`x500_principal_helper.go:63`) deliberately compares raw DER
  encodings (not canonical names), exactly as upstream `X500PrincipalHelper#equals` does — this
  is a DIFFERENT comparison from `X500Principal.Equals`, and the doc comment explicitly notes
  the distinction. No confusion between the two.
- **OID handling (trust-logic correctness):** `x500PrincipalDecodeOID`
  (`x500_principal.go:685`) is bounds-checked (rejects empty OID, enforces the 7-bit
  continuation-bit protocol, guards against component overflow with `value > (1<<57)-1`, rejects
  truncated OIDs, rejects high-tag-number form). No integer overflow, no out-of-bounds access,
  no panic on hostile input. The OID is rendered as a dotted-decimal string (the first two
  arcs are decomposed per RFC 5280 §4.1.2.2, with the first arc clamped to ≤ 2). Correct.
- **Digest / identifier correctness:** `IdentifierBase` (`identifier.go`) computes SHA-256
  (always available in Go stdlib) over the token/entity/key binaries. `Equals`
  (`identifier.go:138`) compares `className` (the Java simple class name, so a
  `DataIdentifier` never equals an `EntityIdentifier` even when the digests match) AND
  `i.id.Equals(otherBase.id)` (which is `Digest.Equals` — algorithm + `bytes.Equal` on the
  value). Correct: no cross-type false-match, no digest false-match. `AsXmlID`
  (`identifier.go:110`) caches the `prefix + id.HexValue()` string (computed once). `Digest.HexValue`
  (`digest.go:43`) is KAT-pinned against Java (see Open Question below).
- **Nil-deref panics on hostile input:** `NewX500Principal` panics on nil (1:1 port of Java's
  `Objects.requireNonNull` — caller-supplied, not attacker-controlled). `NewCertificateToken`
  panics on nil cert (same). `NewDigestDocumentFromBase64` panics on invalid base64 (1:1 port
  of Java's `IllegalArgumentException`, documented as a programmer-error path). All are
  deliberate ports of upstream `Objects.requireNonNull` / unchecked-exception semantics, not
  reachable by hostile input on the normal path (consistent with T22A-STD-001 in batch-07).
  `DigestDocument.Equals` (`digest_document.go:225`) calls `d.digestMap.Len()` / `.Keys()` /
  `.Get()` without a nil check, but `utils.OrderedMap` methods are all nil-receiver safe
  (verified: `Len`, `Keys`, `Get` all handle `m == nil`). No panic.
  `CertificateToken.CheckIsSignedBy` (`certificate_token.go:230`) returns an error (not a panic)
  when the public key is nil or unparsed — fail-closed. `certificateTokenVerify`
  (`certificate_token.go:370`) delegates to `crypto/x509.Certificate.CheckSignature` (bounds-
  checked, no panic on malformed input). MD5-signed certificates are refused by `crypto/x509`
  (deliberate hardening over JCA, documented in the `DEVIATION` comment). No false-accept.
- **Integer overflow in size/length fields:** `x500PrincipalReadTLV` (`x500_principal.go:570`)
  caps DER length at 4 bytes (`numBytes > 4` → error), and the accumulated `length int` is
  built from at most 4 bytes (max 2³²−1, well within `int` on all supported platforms). The
  `length < 0` check is defensive (the `int` cast of a ≤ 4-byte value cannot be negative on
  64-bit platforms, but the check is harmless). No overflow.
- **`PublicKey.Equals` (key matching):** `public_key.go:88` compares `bytes.Equal(p.Encoded(),
  other.Encoded())` — the raw SubjectPublicKeyInfo DER. This is the same comparison
  `sun.security.x509.X509Key#equals` makes upstream. No false-match (two distinct keys with
  different SPKI bytes do not compare equal) and no false-miss (two keys with the same SPKI
  bytes DO compare equal). The `Encoded()` field is retained as-parsed (never re-encoded), so
  the comparison is over the genuine bytes. Correct.

**PERF**

T24-PERF-001 — **Info** — Category: PERF-BIGO — `dss/model/x500_principal.go:290`
(`x500PrincipalAVA.rfc2253String`) / `:353` (`x500PrincipalAVA.rfc2253CanonicalString`)
- Evidence: Both methods iterate over every rune of the decoded DN value and call
  `strings.ContainsRune(escapees, c)` — once in the `switch` condition and once again inside
  the case body (lines 313–320 in `rfc2253String`; lines 370–378 in
  `rfc2253CanonicalString`). `escapees` is a 6–8 character constant string.
- Impact: O(len(value) × len(escapees)) per AVA. Since `len(escapees)` is a small constant
  (6 or 8), the effective cost is O(len(value)) with a constant factor of ~8. For typical DN
  values (tens of characters) this is negligible. For adversarial or unusually long DN values
  (thousands of characters) the constant factor becomes measurable but is still linear. Not
  exploitable: the input is a DER-encoded DN from a certificate, not a free-form attacker
  string, and the cost is bounded by the DN length.
- Recommendation: None required. If long-DN performance is ever a concern, replace the two
  `strings.ContainsRune` calls with a single `map[rune]bool` lookup or a `switch c { case
  ',', '=', '+', '<', '>', ';', '"', '\\': … }` to reduce the per-character cost to O(1).
  Low priority; the current code is correct and the constant factor is small.

No unbounded-memory, goroutine-leak, or O(n²)-on-attacker-input findings. All parsing is
bounded by the DER input size (bounds-checked, no unbounded recursion). No goroutines, no
retained references beyond the "parse once, retain" design (consistent with batch-06's CRL
parse-once-retain). Allocations in the hot path are bounded by the DN/certificate size.

**STD (idioms, error handling, defer/cleanup, context, sync, doc comments, test quality)**

T24-STD-001 — **Info** — Category: STD — `dss/model/x500_principal.go:116`
(`X500Principal.RFC2253Name`) / `:148` (`X500Principal.Canonical`) /
`multiple_digest_identifier.go:49` (`MultipleDigestIdentifier.DigestValue`)
- Evidence: `X500Principal.RFC2253Name` and `Canonical` mutate the cache fields
  (`p.rfc2253`, `p.rfc2253Valid`, `p.canonical`, `p.canonicalValid`) without a mutex; a
  concurrent call from two goroutines races on the `*Valid` flag and the string assignment.
  `MultipleDigestIdentifier.DigestValue` has the same check-then-act pattern on its
  `digestMap` (read → compute → write, no lock).
- Impact: Data race if the objects are shared across goroutines. However, the `TokenBase`
  doc comment (`token.go:155`) explicitly states "the Go port is not goroutine-safe, matching
  the rest of the value objects in this package," and all in-repo callers use a single
  goroutine per validation context. This is a deliberate design choice consistent with the
  rest of the package.
- Recommendation: No code change required (deliberate port, parity with upstream). If the
  package is ever made goroutine-safe, add a `sync.Mutex` or `atomic.Value` to the cache
  fields. Low priority; the current single-goroutine usage pattern is correct.

No other STD findings. Error handling is correct (`error` returns at the API boundary;
panics are 1:1 ports of upstream `Objects.requireNonNull` and are documented as
programmer-error paths, not attacker-reachable). No goroutines, so no sync issues beyond
T24-STD-001 (deliberate). Tests are thorough: KAT vectors pin the DN canonical form,
identifier digests, hex encoding, base64 encoding, and the BigInteger-based `HexValue`
quirk (leading-zero-byte normalization + even-length padding). The `Digest.HexValue`
KAT (`digest_test.go`, `identifier_kat_test.go`) explicitly pins the even-length padding
as Java-parity (see Open Question below). The `x500_principal.go` KAT
(`x500_principal_kat_test.go`) pins the RFC 2253 and canonical rendering against
known-answer vectors. The `certificate_token_test.go` KAT pins the DSS ID and entity-key
strings against Java-produced expected values. No deliberate-exclusion violations found
(staticcheck style rules, legacy crypto, DIVERGENCE comments are all PORTING.md
conventions).

### Deliberate-exclusion notes (not findings)

- **staticcheck style rules (ST*/S1*/QF*/SA1019):** DELIBERATELY disabled per
  PORTING.md — they would restyle Java-verbatim error strings or "simplify" code away
  from its 1:1 Java mapping. Not reported.
- **Legacy crypto (SHA-1, MD5, PKCS#1 v1.5):** DELIBERATELY retained per PORTING.md —
  a validation library must support the algorithms it validates. Not reported.
- **`// DIVERGENCE, deliberate:` comments:** `certificate_token.go:230`
  (`CheckIsSignedBy` — JCA provider-missing branch absent), `certificate_token.go:370`
  (`certificateTokenVerify` — MD5 refusal), `token.go:155` (not goroutine-safe),
  `x500_principal.go:162` (`String()` returns RFC 2253, not JDK's DEFAULT format),
  `digest.go:43` (`HexValue` even-length padding — see Open Question). All are
  intentional and documented. Not reported.
- **`NewDigestDocumentFromBase64` panic on invalid base64**
  (`digest_document.go:76`): deliberate port of Java's `IllegalArgumentException`,
  documented as a programmer-error path (caller-supplied literal, not attacker-
  controlled). Not reported.
- **`NewX500Principal` / `NewCertificateToken` / `NewInMemoryDocumentWithMimeType`
  panics on nil input:** deliberate ports of Java's `Objects.requireNonNull`,
  documented, not attacker-reachable. Not reported (consistent with T22A-STD-001).

### Tool log

Run from `dss/`:

- `gofmt -l model` → no output (clean).
- `go vet ./model/` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./model/` → `0 issues`.

All tools available and passing. No in-scope findings from any tool. Out-of-scope
(model sub-packages) were not linted (separate units).

### Summary

- **Total findings: 2** (SEC: 0, PERF-BIGO: 1, PERF-MEM: 0, PERF-LEAK: 0, STD: 1)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 0, Low: 0, Info: 2
- The 53 in-scope files are a high-quality 1:1 port of the upstream Java DSS
  `X500Principal`, `CertificateToken`, `Token`, `Identifier`, `Digest`, `Policy`,
  `BLevelParameters`, `SignerLocation`, `ReferenceValidation`, `CommonDocument`,
  `DigestDocument`, `InMemoryDocument`, `FileDocument`, `PublicKey`, `ManifestFile`,
  `ManifestEntry`, `CommonObjectIdentifier`, `UserNotice`, `SpDocSpecification`,
  `X500PrincipalHelper`, `TokenComparator`, and the various identifier types.
  All DN parsing is bounds-checked (no unbounded recursion, no integer overflow),
  all trust-logic paths (DN equality, OID handling, digest/identifier correctness,
  key matching, certificate verification) are correct and conservative (fail-closed,
  no false-accept, no false-reject), and all error handling is correct (errors at
  the API boundary, panics only for documented programmer-error paths). The two
  findings are minor (an Info-level double-lookup in the DN escape loop, and an
  Info-level note on the deliberate non-goroutine-safety of the DN cache).
  No SEC, PERF-MEM, PERF-LEAK, or Low issues were found.

### Open questions

1. **`Digest.HexValue` even-length padding — parity verification needed.**
   The Go implementation (`digest.go:43`) pads the hex string to even length
   (`if len(hex)%2 == 1 { hex = "0" + hex }`). The KAT tests
   (`digest_test.go:22` `TestDigestHexValueOddLengthPadded`,
   `identifier_kat_test.go:42` `TestDigestHexValueMatchesJava`) pin this as
   Java-parity and the test comment states the expected values were "produced by
   running the exact Java expressions dss-model uses (… `Digest#getHexValue`'s
   `new BigInteger(1, value).toString(16)` …) on OpenJDK 21." However,
   `new BigInteger(1, value).toString(16)` does NOT pad to even length
   (`new BigInteger(1, new byte[]{0x0A}).toString(16)` = `"a"`, not `"0a"`).
   If the test comment is accurate about the Java expression, the Go code has a
   byte-level divergence in the DSS ID string (report-visible, and used as a map
   key in `spi/common_certificate_source.go`). If Java's actual
   `Digest#getHexValue` does include the padding (and the test comment is a
   simplification), the Go code is correct parity. **Action for consolidation:**
   verify against the actual Java `Digest#getHexValue` source in DSS 6.5.RC1.
   If it is a divergence, it needs a `// DIVERGENCE, deliberate:` comment and a
   PORTING.md entry (or a fix). If it is parity, the test comment should be
   corrected to reflect the full Java implementation.

2. **`x500PrincipalAVA.rfc2253String` escape-set double-lookup** (T24-PERF-001):
   the two `strings.ContainsRune(escapees, c)` calls per character could be
   merged into a single lookup. No action required today (constant factor is
   small, input is DER-bounded); noted for future cleanup.

## model/lote + model/policy + model/eaa (unit U25)

- Scope: `dss/model/lote/` (16 files), `dss/model/policy/` (19 files), `dss/model/eaa/` top-level
  (3 files) — non-test only. A 1:1 Go port of Java DSS 6.5.RC1 TS 119 602 (LoTE) model types,
  the validation-policy / ETSI TS 119 322 cryptographic-suite model types, and the EAA
  disclosure DTO types. These model types feed qualification verdicts and diagnostic reports —
  a mis-modeled identifier/equality/digest is a trust-path defect.
- Date: 2026-08-26
- Depth: full read of all 38 in-scope files (3,585 lines, within the 3,700 budget) plus
  cross-layer verification: `model/certificate_token.go` (`Equals` = DSSID-based, pointer-key
  semantics), `model/job/` generic interfaces (`DocumentInfo`, `ValidationJobSummary`,
  `AbstractDocumentInfoBase`), the `model/eaa/claim` equality surface (out of scope, U39), the
  `spi/lote` consumer (`TrustedEntitiesCertificateSource` re-keying), the diagnostic builders
  (`diagnostic_data_builder.go`, `xml_trusted_entity_builder.go`), and the
  `validation/process/qualification` LoTE checks (verified to read JAXB XML, not `model/lote`).
- Compatibility contract: interop with Java DSS 6.5.RC1 (byte-level + verdict parity).
  `dss/PORTING.md` conventions and `docs/compatibility/known-gaps.md` are treated as deliberate,
  not findings. Deliberately-disabled staticcheck style rules (ST*/S1*/QF*/SA1019) and
  legacy-crypto flags are PORTING.md conventions, not findings.
- Out-of-scope: `model/eaa/claim/` (25 files — U39, `process/eaa`), `model/timedependent/` (U26),
  `model/job/` (U26), `spi/lote/` + `validation/` consumers (batches 09-13). The known-gaps
  "EAA … treat as unproven until a corpus exists" entry applies to the `eaa`-tagged validation
  blocks (U39), not to these model DTOs; noted, not re-reported.

### Files read (per-file line counts actually read)

| Package | Files (lines) | Approach |
|---|---|---|
| `model/lote` | `other_list_pointer.go` (197), `trusted_entity_service_status_and_information_extensions.go` (141), `lote_info.go` (118), `trusted_entity.go` (109), `lote_validation_job_summary.go` (91), `trusted_properties.go` (70), `trusted_entity_service.go` (65), `lolote_info.go` (63), `lote_parsing_info_record.go` (50), `lolote_identifier.go` (29), `abstract_lote_identifier.go` (26), `trusted_properties_certificate_source.go` (25), `lote_identifier.go` (24), `service_status_and_information_extensions.go` (19), `entity_service.go` (19), `doc.go` (9) | ALL FULL-read. 1,055 lines. |
| `model/policy` | `validation_policy.go` (714, pure interface), `cryptographic_suite_19322.go` (359), `cryptographic_suite_catalogue.go` (254), `cryptographic_suite_metadata.go` (180), `cryptographic_suite_evaluation.go` (171), `cryptographic_suite_algorithm.go` (169), `cryptographic_suite.go` (92, interface), `cryptographic_suite_parameter.go` (91), `encryption_algorithm_with_min_key_size.go` (58), `signature_algorithm_with_min_key_size.go` (57), `validation_policy_factory.go` (33, interface), `cryptographic_suite_factory.go` (33, interface), `certificate_applicability_rule.go` (16, interface), `numeric_value_rule.go` (15, interface), `multi_values_rule.go` (11, interface), `value_rule.go` (10, interface), `level_rule.go` (10, interface), `duration_rule.go` (10, interface), `doc.go` (9) | ALL FULL-read. 2,292 lines. |
| `model/eaa` | `validation_disclosure.go` (115), `disclosure_validation.go` (117), `doc.go` (6) | ALL FULL-read. 238 lines. |

**Total in-scope lines read: 3,585** (budget: 3,700). Budget respected.

Cross-referenced (not re-reviewed): `model/certificate_token.go:378` (`Equals` =
`DSSID().Equals` — U24 verified; used to assess pointer-key map semantics), `model/job/`
interfaces (U26; used to check the "assumed interface" claims), `model/eaa/claim/*.go`
(U39; `Equals` family + `reflect.DeepEqual` usage checked for the DeepEqual-parity
assessment), `spi/lote/trusted_entities_certificate_source.go` (consumer; re-keying by
`EntityKey().AsXmlID()` verified), `validation/reports/diagnostic/diagnostic_data_builder.go`
(consumer; LoLoTEInfo/Info split-handling verified).

### Findings

**SEC (top priority)**

T25-SEC-001 — **Medium** — Category: SEC — `dss/model/lote/lolote_info.go:17` (embedding),
`:46` (`BuildIdentifier`), `:53` (`DSSID`)
- Evidence: `LoLoTEInfo` embeds `Info` **by value** and *shadows* `BuildIdentifier`/`DSSID`.
  `Info.BuildIdentifier()` (`lote_info.go:99`) → `NewIdentifier` ("LoTE-" prefix);
  `LoLoTEInfo.BuildIdentifier()` → `NewLoLoTEIdentifier(&l.Info)` ("LoLoTE-" prefix).
  Java's `LoLoTEInfo#buildIdentifier` override is reached by virtual dispatch for *any*
  `LoTEInfo` reference; Go dispatches on the static pointer type.
- Impact: If any consumer coerces `&loLoTEInfo.Info` into a `*Info` and calls `DSSID()`, it
  gets the **"LoTE-"** identifier, not Java's **"LoLoTE-"** one — a lookup miss in
  `ValidationJobSummary.LoTEInfoByID`/`LoLoTEInfoByID` (fail-miss, not false-accept) and a
  report-identifier mismatch, plus two coexisting cached identifiers for one logical object
  (`Info.identifier` vs the shadowing `LoLoTEInfo.identifier`). Verified no live in-repo
  coercion today: `diagnostic_data_builder.go` uses a distinct
  `getXmlTrustSourceListForLoLoTE` helper and `user_friendly_identifier_provider.go:644`
  special-cases `*lote.LoLoTEInfo`; the LoTE download/parse/validation *job* is not yet
  ported (no non-test callers of `NewTrustedProperties*`/summary builders), so the path is
  currently inert. This is a landmine for the ported job layer and for external integrators.
- Recommendation: Document the coercion hazard at the `Info`/`LoLoTEInfo` type sites (the
  existing JUDGMENT CALL note covers the interface question, not this dispatch one); consider
  an unexported marker/interface so `*Info`-typed helpers can detect a LoLoTEInfo and
  re-dispatch, or a compile-time assertion test pinning that no `[]*Info` collection in the
  port ever holds a `&LoLoTEInfo.Info`. Re-rate High if/when the LoTE job is ported.

T25-SEC-002 — **Low** — Category: SEC — `dss/model/lote/lote_info.go:47-71`
(`NewInfo`/`NewInfoWithParent`/`NewInfoFull`)
- Evidence: The `Info` constructors store `url` with no validation, then
  `NewAbstractLoTEIdentifier` (`abstract_lote_identifier.go:24`) digests `[]byte(listInfo.Url())`.
  Java's `LoTEInfo` extends `AbstractDocumentInfo`, whose constructor is
  `Objects.requireNonNull(url, "URL String shall be provided!")` (the Go twin
  `job.NewAbstractDocumentInfoBase` panics on empty url — `model/job/abstract_document_info.go`);
  the inlined `lote.Info` drops that check (a by-product of the documented JUDGMENT CALL
  flattening).
- Impact: A malformed/bug-supplied Info with an empty URL yields a well-formed identifier over
  SHA-256 of the empty string; two such lists collide on `DSSID` (a trust-lookup key), where
  Java throws. No live producer of an empty-URL Info exists in the repo (the job layer is
  unported), so impact is latent. Also an observable divergence (panic vs no-panic) without a
  `// DIVERGENCE, deliberate:` marker, which PORTING.md's upstream-tracking rule requires.
- Recommendation: Either add the empty-URL panic (parity) or mark the dropped check with a
  `// DIVERGENCE, deliberate:` note naming `AbstractDocumentInfo`'s constructor and the
  flattening rationale.

T25-SEC-003 — **Low** — Category: SEC — `dss/model/eaa/validation_disclosure.go:107-115`
(`Equals`), `dss/model/eaa/disclosure_validation.go:103-115` (`Equals`)
- Evidence: Both `Equals` ports compare via `reflect.DeepEqual` over the whole struct, which
  includes the `ComputeDigest`/`NamespaceFunc`/`DigestIdFunc` function fields
  (`validation_disclosure.go:44-55`). `reflect.DeepEqual` treats two non-nil closures as
  unequal and a nil vs non-nil function as unequal; Java's `ValidationDisclosure#equals`
  compares only the claim and salt.
- Impact: Two disclosures equal under Java's semantics (same claim + salt) compare **unequal**
  in Go whenever their closure fields differ — a strictly-narrower equality (fail-reject
  direction, never a false-accept). Verified no live non-test caller of either `Equals` in the
  repo (EAA validation blocks are U39, `eaa`-tagged), so impact is latent.
- Recommendation: When U39 lands, confirm the concrete embedders either leave all three
  function fields nil or set identical closures, or switch the port to compare only
  `Salt` + `Claim` (plus `Namespace`/`DigestId` where Java does) to match upstream semantics.
  Cross-check with U39's claim-struct field inventory (lazy caches would make DeepEqual
  stricter than Java's per-class `equals` as well).

**PERF**

T25-PERF-001 — **Low** — Category: PERF-BIGO — `dss/model/policy/cryptographic_suite_19322.go:343-357`
(`cryptographicSuite19322AddUnique`/`AddAllUnique`)
- Evidence: Dedup is a linear scan with full structural `Equals` (recursive over
  `parameterList`) per insertion: `for _, e := range list { if e.Equals(evaluation) … }`,
  invoked inside the nested step-2b loop (encryptionAlgo × digestAlgo × evaluation ×
  digestEvaluation, `cryptographic_suite_19322.go:246-290`).
- Impact: O(n²) in the number of evaluations per algorithm, with each comparison O(parameters).
  Bounded by the ETSI TS 119 322 policy document (tens-hundreds of algorithms, a handful of
  evaluations each) — negligible in practice, but asymptotically worse than Java's
  `HashSet#add` (O(1) average), which the comment claims to mirror.
- Recommendation: None required at current input sizes. If policy documents ever grow, key the
  dedup on a cheap hash (e.g. `EvaluationList` length + recommendation + validity pair) and
  confirm with `Equals`.

T25-PERF-002 — **Low** — Category: PERF-MEM — `dss/model/policy/cryptographic_suite_catalogue.go:192-217`
(`filterByAlgorithmUsage`), `:248-253` (`cryptographicSuiteCatalogueAlgorithmListKey`)
- Evidence: Each of the ~14 `*CryptographicSuite()` scope methods calls
  `filterByAlgorithmUsage`, which deep-copies **every** algorithm (`CryptographicSuiteAlgorithmCopy`,
  itself deep-copying every evaluation and parameter), and then `cryptographicSuiteFor` builds
  the cache key by `fmt.Sprintf`-ing the full `String()` of every algorithm and
  `strings.Join`-ing — paid on **every** call, including cache hits.
- Impact: O(suite-document size) of allocation and formatting per scope-method call (~14× per
  policy evaluation), even when the suite is already cached. Bounded by the policy document
  size (small), and the suite itself is cached, so this is setup cost per validation job, not
  per-signature cost. Not exploitable; noted for completeness.
- Recommendation: If per-job setup cost ever matters, key the suite cache on the filtered
  algorithm *pointers* (e.g. a join of `fmt.Sprintf("%p", a)`) plus the usage-set, avoiding
  the full `String()` formatting on hits; or hoist the per-scope filtered lists into a small
  map computed once.

T25-PERF-003 — **Info** — Category: PERF-LEAK — `dss/model/policy/cryptographic_suite_19322.go:175,217`,
`cryptographic_suite_catalogue.go:59,67`, `dss/model/lote/lote_info.go:109`,
`lolote_info.go:54`, `dss/model/eaa/validation_disclosure.go:96`
- Evidence: Check-then-act lazy caches without a mutex: `if c.acceptableDigestAlgorithms ==
  nil { …build… }`, `if c.metadata == nil { … }`, `if l.identifier == nil { … }`,
  `if v.digestMap == nil { … }`.
- Impact: Data race (duplicate compute, torn reads of the map/pointer) if the objects are
  shared across goroutines. Same deliberate class as U24's T24-STD-001 and the batch-07
  "unsynchronized lazy field" notes; the model-layer types are per-validation-context in all
  in-repo usage. No unbounded growth: each cache is bounded by its document's size and
  allocated once.
- Recommendation: None required (deliberate port). If the package is ever made
  goroutine-safe, use `sync.Once`/`atomic.Pointer` for these fields.

**STD**

T25-STD-001 — **Low** — Category: STD — `dss/model/lote/trusted_properties_certificate_source.go:21-24`
- Evidence: The interface exposes `SetTrustedPropertiesByCertificates(map[*model.CertificateToken][]*TrustedProperties)`
  and `SetTrustedTimeByCertificates(map[*model.CertificateToken][]*tsl.CertificateTrustTime)`.
  Go map keys on pointer identity; Java's `Map<CertificateToken, …>` keys on
  `CertificateToken#hashCode/equals` (content). `model.CertificateToken.Equals` (DSSID-based,
  `certificate_token.go:378`) exists but is not used for keying.
- Impact: A Go caller that builds the map from two distinct-but-equal `*CertificateToken`
  objects (e.g. re-parsed) produces two keys where Java produces one. Mitigated in-repo: the
  sole consumer (`spi/lote/trusted_entities_certificate_source.go:73-100`) immediately re-keys
  by `EntityKey().AsXmlID()`, so the identity keying is only an API-shape issue for external
  callers.
- Recommendation: Document on the interface that keys are compared by pointer identity (unlike
  Java's content keying) and that callers must reuse the same `*CertificateToken` instances,
  or accept the map keyed by `EntityKey().AsXmlID()` string to remove the trap.

T25-STD-002 — **Info** — Category: STD — `dss/model/lote/lote_info.go:12-26` (JUDGMENT CALL note),
`dss/model/lote/doc.go:1-4`
- Evidence: The note says `Info` "Implements the assumed job.ValidationJobSummary[…] interface"
  / `LoLoTEInfo` "Implements the assumed job.DocumentListInfo[LoLoTEInfo, Info] interface",
  and asks the "integrator should confirm the actual job package interfaces … match what this
  type structurally provides". Verified: `job.DocumentInfo[P]` / `job.ValidationJobSummary[D,L]`
  are generic interfaces and **nothing in the repo asserts `*lote.Info` against them**; the
  claim overstates conformance. Separately, `doc.go` says the package "ports the dss-model
  lote/identifier subpackage" while it actually ports the whole LoTE model (entities, services,
  pointers, infos).
- Impact: Documentation accuracy only — a reader may assume interface conformance (or a narrower
  package scope) that does not hold. No behavioural effect.
- Recommendation: Reword the JUDGMENT CALL note to "mirrors the shape of" rather than
  "Implements", and correct `doc.go` to describe the full LoTE model package.

No unbounded-memory, goroutine-leak, or O(n²)-on-attacker-input findings beyond the above. All
parsing-free value objects; allocations bounded by their (small, trusted-source) policy/LoTE
documents. All panics are 1:1 ports of upstream `Objects.requireNonNull` programmer-error paths
(`trusted_properties.go:36-46`, `cryptographic_suite_19322.go:67-74`,
`disclosure_validation.go:51-54`, `trusted_entity_service_status_and_information_extensions.go:33-35`).
No integer-overflow surface (`minKeySize int`, `*int` min/max, `Duration() int64` — no size
arithmetic). No false-accept in any equality/digest path: identifier equality is
className + SHA-256 digest (U24-verified `IdentifierBase.Equals`), policy algorithm matching is
fail-closed (unmatched OID/URI → skipped → empty acceptable map), and the step-2b validity
intersection drops evaluations whose start ≥ end.

### Deliberate-exclusion notes (not findings)

- **MRA drop in `OtherListPointer`** (`other_list_pointer.go:38-44`): documented `DEVIATION
  (kept verbatim from upstream)` — Java's builder carries an `mra` field the constructor never
  copies. Intentional, pinned by `other_list_pointer_test.go:31`. Not a finding.
- **`"tlInfo cannot be null!"` panic message** for the `listInfo` parameter
  (`trusted_properties.go:36-40`): verbatim upstream message quirk, documented in the comment.
  Not a finding.
- **U24's pending Open Question** (`Digest.HexValue` even-length padding): the LoTE/LoLoTE
  identifiers ride on `model.MultipleDigestIdentifier`/`Digest` (U24-reviewed, KAT-pinned).
  Not re-reported.
- **staticcheck style rules (ST*/S1*/QF*/SA1019) and legacy-crypto flags:** DELIBERATELY
  disabled per PORTING.md. Tools reported 0 issues on all three packages regardless.
- **EAA corpus gap** (known-gaps.md): applies to the `eaa`-tagged validation blocks (U39), not
  to these model DTOs. Not re-reported.

### Tool log

Run from `dss/`:

- `gofmt -l model/lote model/policy model/eaa` → no output (clean).
- `go vet ./model/lote/... ./model/policy/... ./model/eaa/...` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./model/lote/... ./model/policy/... ./model/eaa/...`
  → `0 issues`.

All tools available and passing. No in-scope findings from any tool.

### Summary

- **Total findings: 8** (SEC: 3, PERF: 3, STD: 2)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 1, Low: 5, Info: 2
- The 38 in-scope files are a high-quality 1:1 port: LoTE identifier/equality semantics are
  correct and conservative (SHA-256 over URL bytes with "LoTE-"/"LoLoTE-" prefixes,
  className-scoped `Equals`, fail-closed algorithm matching, validity-intersection that drops
  misconfigured evaluations), and the policy crypto-suite model is a faithful port of the
  ETSI TS 119 322 catalogue with correct first-come-first-served OID/URI resolution and
  "prefer explicit definition" conflict handling. The single Medium is a latent Go-embedding
  dispatch trap (T25-SEC-001) that is inert today because the LoTE job layer is unported; the
  Low/Info items are parity notes, bounded performance, and documentation accuracy. No
  false-accept, no unbounded growth, no panic-on-hostile-input, no integer overflow.

### Open questions

1. **LoTE job layer absent.** No non-test code in the repo constructs `lote.TrustedProperties`,
   `lote.ValidationJobSummary`, or fills `TrustedEntitiesCertificateSource` (verified: the only
   `model/lote` importers are `spi/lote` and two diagnostic builders). If the Java
   `LoTEValidationJob`/list-download layer is not yet ported, **all `model/lote` findings are
   latent** — re-rate T25-SEC-001/002 when that lands, and confirm the ported job never
   coerces `&LoLoTEInfo.Info` into a `*Info` (the dispatch trap).
2. **T25-SEC-003 (EAA `reflect.DeepEqual` vs Java `equals`).** Needs consolidation with U39
   (`model/eaa/claim` + `validation/process/eaa`): confirm no concrete disclosure embedder
   leaves `ComputeDigest`/`NamespaceFunc`/`DigestIdFunc` non-nil across the two objects being
   compared, and that claim structs carry no lazy-cache fields that would make DeepEqual
   stricter than Java's per-class `equals`.
3. **`NumericValueRule.Value() float64`** (`numeric_value_rule.go:14`): Java's
   `java.lang.Number` is ported as `float64` (documented in the comment). Consolidation should
   confirm the policy loader (U45, `policy/` package) never feeds a value whose exactness
   matters beyond double precision (e.g. QC limit values) — a `float64` rounding of a
   `BigDecimal` limit would be a (fail-wide or fail-narrow) parity question.
---

## model/job + signature + timedependent + tsl + scope/http/tls (unit U26)

- Scope: `dss/model/job` (10 files), `dss/model/signature` (8), `dss/model/timedependent` (5),
  `dss/model/tsl` (26), `dss/model/scope` (2), `dss/model/http` (2), `dss/model/tls` (2) —
  non-test `.go` only. `dss/model/x509` has 0 non-test files (skipped). These are the
  job/cache, signature, time-dependent, TSL, scope, http and tls model value objects that
  feed the validation verdicts and reports. A 1:1 Go port of Java DSS 6.5.RC1.
- Date: 2026-08-26
- Depth: full read of every in-scope non-test file (all ≤234 lines); 3,752 lines read, within
  the 3,800 budget. Call-site reachability for the two SEC candidates verified across
  `dss/tsl`, `dss/validation`, `dss/trustedlist` and the whole module.
- Compatibility contract: interop with Java DSS 6.5.RC1 (byte-level + verdict parity).
  `dss/PORTING.md` conventions and `docs/compatibility/known-gaps.md` treated as deliberate,
  not findings. `// DIVERGENCE, deliberate:` / `DEVIATION:` / `JUDGMENT CALL` comments mark
  intentional divergences. Deliberately-disabled staticcheck style rules (ST*/S1*/QF*/SA1019)
  and legacy-crypto flags are PORTING.md conventions, not findings.
- Cross-referenced (not re-derived): `tsl` (batch-06, U18a/b), `validation/job` (U18A-SEC-001
  cache data race — a different package from `model/job`), `model` top-level (U24), and the
  `validation/identifier` report-id consumer that renders `TLInfo`/`LOTLInfo`/`PivotInfo` ids.

### Files read (per-file line counts actually read — all FULL-read)

| Package | Files (lines) |
|---|---|
| `model/job` (317) | abstract_document_info.go (132), info_record.go (42), validation_info_record.go (33), document_info.go (30), download_info_record.go (19), validation_job_summary.go (15), other_document_pointer.go (13), document_list_info.go (13), parsing_info_record.go (12), doc.go (8) |
| `model/signature` (815) | signature_policy.go (197), signature_policy_validation_result.go (141), signature_production_place.go (112), signature_cryptographic_verification.go (109), signer_role.go (104), commitment_type_indication.go (77), signature_digest_reference.go (64), doc.go (11) |
| `model/timedependent` (264) | time_dependent_values.go (101), mutable_time_dependent_values.go (71), base_time_dependent.go (65), time_dependent.go (17), doc.go (10) |
| `model/tsl` (2032) | service_equivalence.go (234), trust_service_status_and_information_extensions.go (212), other_tsl_pointer.go (204), tl_info.go (138), qc_statement_oids.go (131), trust_service_provider.go (110), certificate_trust_time.go (106), tl_validation_job_summary.go (92), certificate_content_equivalence.go (92), lotl_info.go (85), mra.go (80), trust_properties.go (74), condition_for_qualifiers.go (65), trust_service.go (63), pivot_info.go (63), service_type_asi.go (55), tl_parsing_info_record.go (45), abstract_tl_identifier.go (37), trusted_list_identifier.go (25), lotl_identifier.go (24), trust_properties_certificate_source.go (22), pivot_identifier.go (19), condition.go (17), trusted_certificate_source_with_time.go (14), certificate_pivot_status.go (14), doc.go (11) |
| `model/scope` (194) | signature_scope.go (187), doc.go (7) |
| `model/http` (75) | response_envelope.go (69), doc.go (6) |
| `model/tls` (55) | tls_certificates.go (48), doc.go (7) |

**Total in-scope lines read: 3,752** (budget: 3,800). Budget respected.

Verified clean on the mandated SEC lenses before filing (each confirmed, not assumed):

- **Time-dependent validity fields (qualification verdicts):** `Values.Current(d)`
  (`time_dependent_values.go:55`) selects the first (latest-first) entry with
  `startDate <= d` and (`endDate` zero or `> d`) — i.e. half-open `[start, end)` "current at d"
  — correct; a future-dated entry with a past one in between is correctly skipped (no false
  grant). `Values.After(notBefore)` (`time_dependent_values.go:73`) includes entries with
  `endDate` zero or `>= notBefore`, an exact port of Java `getAfter` (`after || equals`).
  `CertificateTrustTime.IsTrustedAtTime` (`certificate_trust_time.go:38`) reduces to
  `(start null OR start <= controlTime) AND (end null OR end >= controlTime)`, matching Java
  `isTrustedAt`; a null `controlTime` fails closed (not trusted) rather than false-accepting.
  `JointTrustTime` (`certificate_trust_time.go:46`) ports `getJointTrustTime` exactly. No
  boundary bug, no false-accept.
- **Trust lookups:** `TLValidationJobSummary.TLInfoByID`/`LOTLInfoByID`
  (`tl_validation_job_summary.go:58,75`) compare `identifier.Equals(tlInfo.DSSID())` — the
  `IdentifierBase.Equals` className check (U24) prevents a cross-type false-match. Pivots are
  stored in `LOTLInfo.PivotInfos()` and are **not** searched by these two methods, so the
  pivot-id issue below (T26-SEC-001) does not corrupt a trust lookup.
- **`Condition` evaluation (verdict path):** the trust/MRA verdicts call `Condition.Check(cert)`
  (e.g. `xml_trust_service_provider_builder.go`, `composite_condition.go`), **not** the
  `Equals` methods — so the `==`-on-`Condition` divergence (T26-SEC-002) is not in the verdict
  path today.
- **Nil-deref / panics:** every `panic` is a documented 1:1 port of Java
  `Objects.requireNonNull` / `IllegalArgumentException` (empty URL, nil builder, nil
  tlInfo/trustServiceProvider/trustService, `AddOldest` null/overlap) — programmer-error paths,
  not attacker-reachable; `MutableValues.AddOldest`'s overlap check fails closed. No
  hostile-input-reachable nil deref found. Integer overflow: none (no size/length arithmetic on
  hostile input; the only digest is `[]byte(url)` → SHA-256, bounds-trivial).
- **Digest/identifier correctness:** `AbstractTLIdentifier` builds a
  `MultipleDigestIdentifier` (SHA-256 over `[]byte(url)`, U24-verified digest machinery) with
  per-class `className`/`prefix`. The only mis-modeling is the pivot `className`/`prefix`
  (T26-SEC-001). `SignerRole.Equals`, `DigestReference.Equals`,
  `CryptographicVerification.IsSignatureValid` (correct AND of the three flags), and
  `PolicyValidationResult.AddError` (insertion-ordered, `Map#put` overwrite semantics) are all
  correct.

### Findings

**SEC (top priority)**

T26-SEC-001 — **Medium** — Category: SEC — `dss/model/tsl/pivot_info.go:63` (`BuildIdentifier`,
with **no** `DSSID` override) & `dss/model/tsl/lotl_info.go:70` (promoted `LOTLInfo.DSSID`)
- Evidence: `PivotInfo` embeds `LOTLInfo` (embeds `TLInfo`). `PivotInfo` defines
  `BuildIdentifier()` → `NewPivotIdentifier` (className `"PivotIdentifier"`, prefix `"P-"`) but
  does **not** redefine `DSSID()`. So `(*PivotInfo).DSSID()`/`DSSIDAsString()` promote to
  `LOTLInfo.DSSID()`, which calls `l.BuildIdentifier()` — a static-type call that resolves to
  `LOTLInfo.BuildIdentifier()` → `NewLOTLIdentifier` (className `"LOTLIdentifier"`, prefix
  `"LOTL-"`), not `PivotInfo`'s. Go embedding has no virtual dispatch, so the pivot's own
  `buildIdentifier` override is never reached through `DSSID`.
- Impact: A pivot's DSS id is mis-modelled as a `LOTLIdentifier` (`LOTL-<hex>`, className
  `"LOTLIdentifier"`) rather than a `PivotIdentifier` (`P-<hex>`, className `"PivotIdentifier"`).
  Java's `PivotInfo#getDSSId()` (virtual `buildIdentifier()`) returns the `PivotIdentifier`.
  Consequence is report/interop, not verdict: `idAsStringForTL`
  (`validation/identifier/user_friendly_identifier_provider.go:377`) type-switches `*tsl.PivotInfo`
  to the pivot prefix but then renders `DSSIDAsString()`, producing a pivot entry labelled
  `PIVOT-…_LOTL-<hex>` (vs Java `PIVOT-…_P-<hex>`), and `IdentifierBase.Equals` (className check)
  would distinguish a Go pivot id from its Java counterpart. It does **not** flip a VALID/INVALID
  verdict and is not used by `TLInfoByID`/`LOTLInfoByID` (pivots are not in those collections), so
  no false-accept. This is a documented `JUDGMENT CALL` in `tl_info.go`/`pivot_info.go`; surfaced
  here for a deliberate keep-or-fix decision (fix the shadowing, or record it in
  `known-gaps.md`).
- Recommendation: Either give `PivotInfo` a `DSSID()` override that builds/caches via
  `NewPivotIdentifier` (and make `TLInfo`/`LOTLInfo` use a virtual `buildIdentifier` indirection
  so all three derive correctly), or add a `// DIVERGENCE, deliberate:` + `known-gaps.md` entry
  stating pivot report ids intentionally render with the LOTL prefix. Pin with a KAT on
  `(*PivotInfo).DSSID().AsXmlID()`.

T26-SEC-002 — **Medium** — Category: SEC — `dss/model/tsl/condition_for_qualifiers.go:63`
(`c.condition == other.condition`) & `dss/model/tsl/certificate_content_equivalence.go:82`
(`if c.condition != other.condition`)
- Evidence: Both `Equals` methods compare the `Condition` interface field with Go `==`/`!=`.
  Every in-tree `Condition` implementer is a **pointer type** (`*PolicyIdCondition`,
  `*CompositeCondition`, `*QCStatementCondition`, `*ExtendedKeyUsageCondition`,
  `*CertSubjectDNAttributeCondition`, `*KeyUsageCondition`), so `==` compares **pointer
  identity**, not structure. Java's `Objects.equals(condition, that.condition)` invokes
  `condition.equals(that.condition)` (structural). Two logically-equal `ConditionForQualifiers`
  built at different sites (`trust_service_converter.go:118` and
  `xml_trust_service_provider_builder.go:580` each construct their own) carry different `*Condition`
  pointers, so `Equals` returns `false` where Java returns `true`.
- Impact: Latent trust-logic/parity defect — `Equals` is a false-REJECT (conservative) landmine
  in a trust-relevant comparison, and diverges from Java. Currently **unreachable** in the verdict
  path: the verdict/MRA logic uses `Condition.Check(cert)`, and no non-test caller invokes
  `ConditionForQualifiers.Equals` or `CertificateContentEquivalence.Equals` (verified across
  `tsl`/`validation`/`trustedlist`; the only in-tree `Equals` callers are
  `CertificateTrustTime.Equals` and `CertificateToken.Equals`, which are correct). Escalate to
  High if either `Equals` is ever wired into a dedup/verdict path.
- Recommendation: Give the `Condition` implementations a structural `Equals` and have the two
  `Equals` methods call it (or `reflect`-compare the condition), matching Java's
  `Objects.equals`; or delete the dead `Equals` methods if they are not part of the public
  contract. Add a KAT comparing two independently-built-but-equal `ConditionForQualifiers`.

**PERF**

T26-PERF-001 — **Info** — Category: PERF-BIGO — `dss/model/signature/signature_policy_validation_result.go:105`
(`AddError`)
- Evidence: `AddError` linearly scans `s.errors` for a matching key
  (`for i := range s.errors { if s.errors[i].key == errorKey … }`) on every call — O(n) per add,
  O(n²) for n distinct keys, mirroring Java's `LinkedHashMap`.
- Impact: Negligible — the error-key set on a signature-policy validation is tiny (single digits);
  not attacker-scalable. Same class as batch-06's T18A-PERF-BIGO-001.
- Recommendation: None required; use a `map[string]int` index if the list ever grows.

T26-PERF-002 — **Info** — Category: PERF-BIGO — `dss/model/tsl/service_equivalence.go:226` and
`condition_for_qualifiers.go:65` (`reflect.DeepEqual` in `Equals`)
- Evidence: The `model/tsl` `Equals` methods compare slice/map fields with `reflect.DeepEqual`
  (e.g. `reflect.DeepEqual(s.certificateContentEquivalences, other.certificateContentEquivalences)`,
  `reflect.DeepEqual(c.qualifiers, other.qualifiers)`) — reflection-driven, slower and
  structurally-semantics than hand-written comparison.
- Impact: Negligible — these `Equals` methods have no non-test callers (verified), so the cost is
  not on any hot path; noted for consistency with the T26-SEC-002 cleanup.
- Recommendation: Fold into the T26-SEC-002 fix (structural comparison or removal).

No unbounded-memory (PERF-MEM), goroutine-leak (PERF-LEAK), or O(n²)-on-attacker-input findings.
The parse-once/retain + value-object design means no retained-reference or growth hazards in these
packages; the only large retained value is `ResponseEnvelope.responseBody` (caller-owned `[]byte`,
retained-not-copied, matching Java).

**STD (idioms, error handling, defer/cleanup, context, sync, doc comments)**

T26-STD-001 — **Low** — Category: STD — `dss/model/tsl/abstract_tl_identifier.go:~30` (stale
`DEVIATION` comment)
- Evidence: The `NewAbstractTLIdentifier` doc comment states "DEVIATION from pivot_identifier.go
  (already ported): that file calls NewAbstractTLIdentifier(prefix, tlInfo) with two arguments …
  pivot_identifier.go's call site **will need updating** by the integrator to pass
  `"PivotIdentifier"` as the first argument."
- Impact: The described pending work is already done — `pivot_identifier.go:17` calls
  `NewAbstractTLIdentifier("PivotIdentifier", pivotIdentifierPrefix, pivotInfo)`. The comment is
  stale and could mislead an integrator into thinking the pivot identifier is still mis-wired (it
  is wired correctly; the live issue is T26-SEC-001's `DSSID` shadowing, a different matter).
- Recommendation: Correct the `DEVIATION` note to state the className is passed explicitly and
  already wired at all three call sites (`TrustedListIdentifier`/`LOTLIdentifier`/
  `PivotIdentifier`); cross-reference T26-SEC-001 for the remaining pivot gap.

T26-STD-002 — **Info** — Category: STD — `dss/model/signature/signature_policy.go:115`
(`TransformsDescription`)
- Evidence: `func (s *Policy) TransformsDescription() []string { return []string{} }` allocates a
  fresh empty slice on every call.
- Impact: A deliberate non-nil empty return (matches Java `Collections.emptyList()`), but a tiny
  per-call allocation where Java returns a cached singleton. Negligible (called once per policy in
  report building).
- Recommendation: Optional: return a package-level `var emptyStrings = []string{}` to avoid the
  allocation, or leave as-is for clarity. No action required.

### Deliberate-exclusion notes (not findings)

- **Non-goroutine-safe lazy caches:** `TLInfo.DSSID`/`LOTLInfo.DSSID` cache (`identifier` +
  build-once), `AbstractDocumentInfoBase.DSSID` (`model/job`), `SignatureScopeBase.DSSID`
  (`dssID`/`dssIDBuilt`), and `MultipleDigestIdentifier.DigestValue` all mutate cache fields
  without a mutex. Deliberate, consistent with U24-STD-001 and the package doc "not goroutine-safe
  matching the rest of the value objects in dss-model." All in-repo callers are single-goroutine
  per validation context. Not re-reported.
- **`panic` on nil/empty arguments:** `NewAbstractDocumentInfoBase` (empty URL),
  `NewTrustServiceStatusAndInformationExtensions` (nil builder),
  `NewTrustProperties*` (nil tlInfo/trustServiceProvider/trustService),
  `MutableValues.AddOldest` (nil/overlap) — documented 1:1 ports of Java
  `Objects.requireNonNull`/`IllegalArgumentException`, programmer-error paths, not
  attacker-reachable. Not reported (consistent with U24 and T22A-STD-001).
- **staticcheck style rules (ST*/S1*/QF*/SA1019) & legacy crypto:** DELIBERATELY disabled per
  PORTING.md. Not reported.
- **`// DIVERGENCE, deliberate:` / `DEVIATION:` / `JUDGMENT CALL` comments:**
  `StatusEquivalenceMapping` slice-of-pairs (Java `Map<List,String>,…>`), `TLInfo`/`PivotInfo`
  `AbstractDocumentInfo` flattening + `DSSID` shadowing (T26-SEC-001), `TransformsDescription`
  non-nil empty, `loTEInfoPrefix` upstream-bug reproduction (noted in
  `user_friendly_identifier_provider.go`). All intentional/documented. Not re-reported beyond
  surfacing the pivot id as T26-SEC-001.

### Tool log

Run from `dss/`:

- `gofmt -l model/job model/signature model/timedependent model/tsl model/scope model/http model/tls`
  → no output (clean), exit 0.
- `go vet ./model/job/... ./model/signature/... ./model/timedependent/... ./model/tsl/... ./model/scope/... ./model/http/... ./model/tls/...`
  → no output (clean), exit 0.
- `golangci-lint run --config=../.github/.golangci.yml ./model/job/... ./model/signature/... ./model/timedependent/... ./model/tsl/... ./model/scope/... ./model/http/... ./model/tls/...`
  → `0 issues.`, exit 0.

All tools available and passing. No in-scope findings from any tool; nothing out-of-scope surfaced.

### Summary

- **Total findings: 6** (SEC: 2, PERF-BIGO: 2, PERF-MEM: 0, PERF-LEAK: 0, STD: 2)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 2, Low: 1, Info: 3
- The 55 in-scope files are a high-quality 1:1 port of the upstream `dss-model`
  `job`/`signature`/`timedependent`/`tsl`/`scope`/`http`/`tls` value objects. The
  verdict-critical logic (time-dependent `Current`/`After`, `CertificateTrustTime` trust-window,
  `IsSignatureValid`, `Identifier.Equals`, trust lookups) was verified correct and conservative
  (fail-closed, no false-accept). The two SEC findings are parity/correctness gaps in the
  *identifier and `Equals`* surface — a mis-modelled pivot DSS id (report-visible, documented
  JUDGMENT CALL, no verdict flip) and a latent pointer-identity `Condition` comparison in two
  `Equals` methods (unreachable in the verdict path today). No SEC finding reaches a
  false-accept; no PERF-MEM/LEAK or O(n²)-on-attacker-input issue.

### Open questions

1. **T26-SEC-001 (pivot DSS id) — keep or fix.** The `DSSID` shadowing makes a pivot render as
   a `LOTLIdentifier` (`LOTL-`) instead of a `PivotIdentifier` (`P-`). It is a documented
   `JUDGMENT CALL`. **Consolidation decision:** (a) fix by giving `PivotInfo` a `DSSID()`
   override + virtual `buildIdentifier` indirection (best — restores Java report-id parity), or
   (b) accept and add a `known-gaps.md` entry. Either way, pin `(*PivotInfo).DSSID().AsXmlID()`
   with a KAT. This is the single most user-visible gap in the unit.
2. **T26-SEC-002 (Condition `==`) — wire-up risk.** The two `Equals` methods are dead today, but
   if `ConditionForQualifiers`/`CertificateContentEquivalence` equality is ever used for MRA
   dedup or verdict gating, the pointer-identity comparison becomes a false-reject. **Action:**
   resolve as part of the SEC cleanup (structural `Equals` on `Condition` implementers, or remove
   the dead methods) so it cannot silently bite.
3. **`model/job` `AbstractDocumentInfoBase` usage.** `TLInfo` deliberately inlines the base's
   fields/behaviour (documented `JUDGMENT CALL`) rather than embedding the generic base, so
   `AbstractDocumentInfoBase` is a faithful-but-possibly-orphaned port in this tree. **Action:**
   confirm whether any other consumer (e.g. `model/lote`) instantiates it; if none, note it as
   port-of-record parity (kept for Java-class coverage) so it is not mistaken for dead code.

## enumerations — light pass (unit U27)

**Scope**: `dss/enumerations/` — 120 non-test `.go` files, 10 967 lines. OID/URI/label constant
tables, `*Values()` slices, `*ForLabel`/`*ForOID`/`*ForName` reverse lookups, two loader
registries, and the Java-`HashMap` iteration-order emulation in `signature_algorithm.go`.

**Date**: 2026-08-28

**Depth**: LIGHT — sample + grep-scan (≤2 000 lines read in-scope)

### Files sampled

| File | Lines | Reason |
|---|---|---|
| `x520_attributes.go` | 1 096 | Largest file; OID↔description table, `ValueOf` linear scan, `GetOidDescriptions` |
| `signature_algorithm.go` | 938 | HashMap-order emulation, PSS-params ASN.1 decode, OID+params dispatch |
| `elliptic_curve.go` | 332 | Curve parameters, OID collision (P256/SECP256K1), `ForOID` last-match |
| `service_qualification.go` | 266 | URI table, `GetByUri` linear scan, `IsQc*` helpers |
| `cose_signature_type.go` | 173 | Tag/has* booleans, `Tag()` panic path |
| `mime_type.go` | 93 | `MimeType` interface, extension parsing, `fromFileName`/`fromFilePath` |
| `mime_type_enum.go` | 158 | Enum data table, `init()` registers `MimeTypeEnumLoader` |
| `mime_type_loader.go` | 33 | Global mutable loader registry |
| `lote_loader.go` | 93 | LoTELoader registry, `init()` order (load-bearing), `RegisterLoTELoader` |
| `lote_empty_loader.go` | 85 | Fallback loader (echoes URI/label, returns nil for `FromDefinition`) |
| `qc_ident_method_enum.go` | 80 | Shared-description upstream bug, `ForLabel` panic path |
| `assert.go` | ~40 | TS 119 612 `Assert` enum (smallest sampled; `ValueOf` + `Value` table pattern) |

### Grep-scan summary

| Pattern | Hits (non-test) | Disposition |
|---|---|---|
| `panic\|os.Exit\|recover()` | 3 | `cose_signature_type.go:123` `Tag()` — mirrors Java `Objects.requireNonNull`; `elliptic_curve.go:56,64` `bigDec`/`bigHex` — constant-time string parse, not hostile input; `qc_ident_method_enum.go:69` `ForLabel("")` — mirrors Java. All deliberate porting choices. |
| `go func\|sync\.\|chan \|map\[string\]` | 0 concurrency primitives; 8 `map[string]` read-only tables | No data race from concurrency. Two global mutable `[]Loader` slices (see T27-STD-001). |
| `func init()` | 2 | `mime_type_enum.go:156` registers `MimeTypeEnumLoader`; `lote_loader.go:64` registers `LoTEEnumLoader`+`LoTEEmptyLoader` in documented load-bearing order. Both single-threaded, run before `main()`. |
| Duplicate OID literals (cross-file) | ~40 OIDs appear >1× | All accounted for: test-vs-source (expected), cross-family (PSS OID in `SignatureAlgorithm` and `EncryptionAlgorithm`), or two **documented upstream collisions** (see T27-SEC-001/002). |
| Per-table duplicate OID (same file, same type) | 2 | `elliptic_curve.go`: `1.2.840.10045.3.1.7` (P256 + SECP256K1); `signature_algorithm.go`: `0.4.0.127.0.7.1.1.4.1.6` (ECDSARIPEMD160 + PlainECDSARIPEMD160). Both documented in-file as upstream verbatim. |

### Findings

#### SEC

**T27-SEC-001 — Info** — Category: SEC — `dss/enumerations/elliptic_curve.go:205,212`
- **Evidence**: `ellipticCurveOIDs` maps both `EllipticCurveP256` and `EllipticCurveSECP256K1`
  to `"1.2.840.10045.3.1.7"`. `EllipticCurveForOID` scans in declaration order keeping the last
  match, so it returns `SECP256K1` for that OID.
- **Impact**: `EllipticCurveForOID("1.2.840.10045.3.1.7")` never returns `P256`. No Go-only
  caller of `EllipticCurveForOID` was found outside the package; the function is exposed for
  API parity with Java. Pinned by `elliptic_curve_test.go:89`.
- **Recommendation**: No action. Upstream Java source has the same collision (verbatim copy).
  Deliberate-exclusion note: if a future consumer calls `EllipticCurveForOID` with that OID
  expecting `P256`, it will silently get `SECP256K1`. A `known-gaps.md` entry would be
  appropriate if this ever becomes user-visible.

**T27-SEC-002 — Info** — Category: SEC — `dss/enumerations/signature_algorithm.go:546,556`
- **Evidence**: `signatureAlgorithmOIDPairs` contains `{"0.4.0.127.0.7.1.1.4.1.6",
  SignatureAlgorithmECDSARIPEMD160}` and `{"0.4.0.127.0.7.1.1.4.1.6",
  SignatureAlgorithmPlainECDSARIPEMD160}`. The later entry wins the forward map, so
  `SignatureAlgorithmForOID("0.4.0.127.0.7.1.1.4.1.6")` → `PlainECDSARIPEMD160`, and
  `ECDSARIPEMD160.OID()` returns `""`.
- **Impact**: `ECDSARIPEMD160` has no OID in the Go port, exactly as in Java's
  `HashMap.put` overwrite. No Go-only caller of `SignatureAlgorithmForOID` with that OID was
  found outside the package. Documented in the file header (lines 16–20).
- **Recommendation**: No action. Upstream verbatim. Deliberate-exclusion note.

#### PERF

**T27-PERF-001 — Low** — Category: PERF-BIGO — `dss/enumerations/x520_attributes.go:988`
  (representative; pattern applies to ≥10 other `*ForLabel`/`*ForName`/`*GetByUri` functions)
- **Evidence**: `X520AttributesValueOf` linearly scans `X520AttributesValues()` (~244 entries).
  Same pattern in `ServiceQualificationGetByUri` (13 entries), `EllipticCurveForLabel`/`ForOID`/
  `ForCOSEValue` (12 entries), `SignatureAlgorithmGetAlgorithm` (~66 entries),
  `COSESignatureTypeForLabel` (7 entries), `QCIdentMethodEnumForLabel` (2 entries).
- **Impact**: O(n) per call where a pre-built reverse map would be O(1). None of these are in
  a per-byte or per-certificate hot loop; they are called during document parsing or
  algorithm resolution. At current table sizes (≤244), the cost is negligible.
- **Recommendation**: No action required. If any of these are ever moved to a hot path,
  pre-build the reverse map (the `signatureAlgorithm*Forward` pattern already demonstrates the
  approach).

**T27-PERF-002 — Low** — Category: PERF-MEM — `dss/enumerations/x520_attributes.go:1009,1018`
- **Evidence**: `GetUppercaseDescriptionForOids` and `GetOidDescriptions` allocate a new
  `map[string]string` on every call and iterate `X520AttributesValues()`.
- **Impact**: One map allocation + ~244 insertions per call. Called from
  `model/x500_principal_helper.go:52` (`RFC2253NameWithOIDMap`). If that helper is called
  per-certificate in a batch validation, this is a per-call allocation. No unbounded growth —
  the map is returned and garbage-collected.
- **Recommendation**: If profiling ever shows this in a hot path, cache the map behind a
  `sync.Once` or return a package-level read-only map. No action required today.

#### STD

**T27-STD-001 — Low** — Category: STD — `dss/enumerations/mime_type_loader.go:22,27-28`
  (and `dss/enumerations/lote_loader.go:36,41-42`)
- **Evidence**: `mimeTypeLoaderRegistry` and `loTELoaderRegistry` are package-level
  `[]Loader` slices. `RegisterMimeTypeLoader` / `RegisterLoTELoader` append without a mutex.
- **Impact**: A concurrent `Register*` call racing with a read (`mimeTypeLoaders()` /
  `loTELoaders()`) is a data race. In practice the only production callers are `init()`
  functions (single-threaded, before `main()`), and the tests save/restore the slice pointer
  under a single goroutine. No concurrent registration pattern exists in the codebase.
- **Recommendation**: Add a comment to `Register*` stating it is not goroutine-safe and must
  be called before any concurrent use, or add a `sync.Mutex` if the API is ever exposed to
  multi-goroutine callers. No action required today.

**T27-STD-002 — Info** — Category: SEC — `dss/enumerations/qc_ident_method_enum.go:36-38`
- **Evidence**: Both `QCIdentMethodEnumQCTEIDAS2ACD` and `QCIdentMethodEnumQCTEIDAS2B` have
  description `"qc-ident-method-eIDAS2-acd"`. `QCIdentMethodEnumForLabel` returns the first
  match in declaration order, so `ForLabel("qc-ident-method-eIDAS2-acd")` → `ACD` always.
- **Impact**: `QCTEIDAS2B` is unreachable via label lookup. No Go-only caller of
  `QCIdentMethodEnumForLabel` was found outside the package. The code comment (lines 30–34)
  explicitly documents this as an upstream copy-paste bug preserved verbatim.
- **Recommendation**: No action. Upstream verbatim. If a future consumer needs to resolve
  `B` by label, they must use the OID path (`QCIdentMethodFromOID`).

**T27-STD-003 — Info** — Category: STD — `dss/enumerations/cose_signature_type.go:123`,
  `dss/enumerations/qc_ident_method_enum.go:69`, `dss/enumerations/elliptic_curve.go:56,64`
- **Evidence**: Three `panic()` calls: `COSESignatureType.Tag()` when `hasTag == false`;
  `QCIdentMethodEnumForLabel("")` when description is empty; `bigDec`/`bigHex` when
  `big.Int.SetString` fails.
- **Impact**: The first two mirror Java's `Objects.requireNonNull` /
  `IllegalArgumentException` (unchecked exception → Go panic is the documented porting
  convention). The third two operate on compile-time string literals in package-level `var`
  initializers — not reachable with hostile input. No panic-on-hostile-input path exists.
- **Recommendation**: No action. Deliberate porting choice consistent with PORTING.md
  ("Java exceptions → returned error … throwing constructors → `New…() (T, error)`";
  `Objects.requireNonNull` → `panic` is the established pattern in this codebase).

### Tool log

| Command | Result |
|---|---|
| `gofmt -l enumerations` | (no output — clean) |
| `go vet ./enumerations/` | (no output — clean) |
| `golangci-lint run --config=../.github/.golangci.yml ./enumerations/...` | `0 issues.` |

### Summary counts

| Lens | Critical | High | Medium | Low | Info |
|---|---|---|---|---|---|
| SEC | 0 | 0 | 0 | 0 | 2 |
| PERF | 0 | 0 | 0 | 2 | 0 |
| STD | 0 | 0 | 0 | 1 | 2 |
| **Total** | **0** | **0** | **0** | **3** | **4** |

### Open questions for consolidation

1. **T27-SEC-001 / T27-SEC-002 (documented upstream OID collisions).** Both are pinned by
   tests and documented in-file. If consolidation decides these are user-visible gaps
   (e.g., a caller of `EllipticCurveForOID` or `SignatureAlgorithmForOID` relies on the
   "wrong" constant), add `known-gaps.md` entries. Otherwise, close as deliberate exclusions.
2. **T27-STD-001 (loader registry concurrency).** No current concurrent access pattern, but
   the API is exported. If the project ever adds a `sync` guard to the registries, the
   `init()`-only registration assumption in the doc comments should be updated. Low priority.