# Batch 01b review — cades / cms

- **Scope:** public API layer of `dss/cades` (unit U04, written separately) and `dss/cms` (unit U05).
- **Date:** 2026-08-23
- **Depth:** deep
- **Baseline:** Go port of Java DSS 6.5.RC1. `dss/PORTING.md` conventions (1:1 Java port,
  staticcheck style rules intentionally disabled, SA1019 legacy crypto allowed, byte-identical enum
  values, stdlib-first, no cgo) are **deliberate** and are not findings. `docs/compatibility/known-gaps.md`
  entries are not re-reported (cited as Info where relevant).

---

## dss/cms (unit U05)

Public CMS API wrapping `internal/cmscore` (the RFC 5652 / RFC 3161 engine). `cmscore` was reviewed
separately; this unit covers only the API layer — parsing entry points, builders, the generator, the
attribute-assembly helpers, the digest-provider operators, and the ETSI `SignerAttributeV2`
ASN.1 values. There is **no signature verification** in this package: it builds and re-shapes CMS
SignedData, it never decides a verdict, so there is no false-accept surface here (verdicts live in
`validation/` + `cmscore`).

### Files read

| File | Lines | Read |
|---|---|---|
| `cms/cms_utils.go` | 414 | truncated >400: lines 1–300 + 300–414 (full coverage, two reads) |
| `cms/cms_signer_info_generator_builder.go` | 286 | full |
| `cms/cms_builder.go` | 252 | full |
| `cms/abstract_cms_generator.go` | 238 | full |
| `cms/custom_content_signer.go` | 189 | full |
| `cms/signer_attribute_v2.go` | 149 | full |
| `cms/certified_attributes_v2.go` | 96 | full |
| `cms/cms_signed_document.go` | 83 | full |
| `cms/cms_signed_attribute_table_generator.go` | 81 | full |
| `cms/signed_assertion.go` | 80 | full |
| `cms/doc.go` | 73 | full |
| `cms/signed_assertions.go` | 58 | full |
| `cms/cms_generator.go` | 57 | full |
| `cms/custom_content_signer_builder.go` | 53 | full |
| `cms/precomputed_digest_calculator_provider.go` | 45 | full |
| `cms/custom_message_digest_calculator_provider.go` | 45 | full |
| `cms/digest_calculator_provider.go` | 20 | full |
| `cms/cms.go` | 100 | full |
| `cms/cms_integration_test.go` | 189 | full (quality/coverage skim) |
| `cms/custom_content_signer_kat_test.go` | 199 | full (quality/coverage skim) |
| `cms/asn1_types_test.go` | 120 | full (quality/coverage skim) |
| `cms/cms_signed_attribute_table_generator_test.go` | 78 | full (quality/coverage skim) |

No file was read with an omitted gap; the single >400-line file (`cms_utils.go`) was read in its
entirety across two contiguous reads.

### Findings

**Lens: STD** — 1 finding (Low).

**C05-STD-001**
- **Severity:** Low
- **Category:** STD
- **Location:** `cms/custom_content_signer.go:169` (`ContentSigner` interface, `OutputStream()`),
  consumer at `cms/cms_signer_info_generator_builder.go:129`
- **Evidence:**
  ```go
  type ContentSigner interface {
      AlgorithmIdentifier() *asn1ber.AlgorithmIdentifier
      OutputStream() *bytes.Buffer   // concrete type in a public interface
      Signature() []byte
  }
  ```
- **Impact:** The public interface exposes the concrete `*bytes.Buffer` rather than an
  `io.Writer`/reader abstraction. Callers building a custom `ContentSigner` (the interface is exported
  for exactly that) are pinned to `bytes.Buffer` semantics and to the package's choice of a mutable
  buffer, which the generator then reads with `Write(...).Bytes()`. It works, but it is the one place
  the otherwise interface-clean API leaks a concrete, stateful type across the public boundary.
- **Recommendation:** Keep the buffer (the two-step "data-to-sign" flow needs the accumulated bytes),
  but consider returning a narrow reader (`io.Reader` over `Bytes()`) or documenting that the buffer
  is the stable contract. Low priority — this mirrors the port's deliberate "replace BC's
  `OutputStream` with `*bytes.Buffer`" choice noted in the file header, so treat as an optional
  idiom tidy, not a defect.

**Lens: PERF-BIGO** — 1 finding (Low).

**C05-PERF-001**
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `cms/cms_builder.go:197` (`containsCertificateToken`), `cms/cms_builder.go:270`
  (`uniqueByteSet.add`), `cms/abstract_cms_generator.go:259` (`containsAlgorithmIdentifier`)
- **Evidence:**
  ```go
  func containsCertificateToken(certificates []*model.CertificateToken, candidate *model.CertificateToken) bool {
      for _, certificate := range certificates { if certificate.Equals(candidate) { return true } }
      return false
  }
  // add() likewise scans s.values linearly; called once per cert/CRL/OCSP in ExtendCMSSignedData
  ```
- **Impact:** De-duplication is a linear scan per element, so `Builder.certificateStore` /
  `addJcaCertStoreCertificates` / `ExtendCMSSignedData` are O(n²) in the number of certificates, CRLs,
  or OCSP responses. At realistic chain sizes (single digits to low tens) this is negligible, but it is
  an avoidable quadratic that would hurt if a caller ever passed a large certificate set.
- **Recommendation:** If a large-set path is ever expected, back the de-dup with a `map[string]struct{}`
  keyed on the DER bytes (or `CertificateToken.Equals`'s underlying value). Not worth changing today
  given the small n; note it so a future "large CMS" streaming phase (see known-gaps) doesn't inherit
  it.

**Lens: PERF-MEM** — 1 finding (Low).

**C05-PERF-002**
- **Severity:** Low
- **Category:** PERF-MEM
- **Location:** `cms/cms_utils.go:254` (`UtilsWriteSignedDataDigestAlgorithmsEncoded`),
  `cms/cms_utils.go:313` (`UtilsWriteSignedDataSignerInfosEncoded`)
- **Evidence:**
  ```go
  element := cms.core.SignedData().DigestAlgorithmsElement()
  if element != nil { w.Write(element.DEREncoded()) } else {
      for _, a := range algorithms { encodings[i] = a.DER() }   // re-serialize every call
      w.Write(cmscoreDERSetOf(encodings))
  }
  // SignerInfos path: signer.DER() re-serializes the whole SignerInfo on every call
  ```
- **Impact:** `SignedData` caches the pre-encoded `digestAlgorithmsElement` / `signerInfosElement`
  **only when the CMS was parsed** (see `internal/cmscore/signed_data.go:150-153`). For a CMS that was
  **built** in this process (the normal CAdES/PAdES signing path), both are nil, so each call to these
  two `write*` helpers re-DERs every `AlgorithmIdentifier` and — more expensively — every full
  `SignerInfo` (its signed attributes, signature, etc.) from scratch. These are used for
  evidence-record / archive-time-stamp-v2 message-imprint computation; if a caller invokes them more
  than once per signature (or over many signatures in a batch) without caching, the per-call cost is
  O(CMS size) re-serialization with no memoization at the API layer.
- **Recommendation:** Memoize the built SET encoding (e.g., a `sync.OnceValue`-style cache or a field
  on the `CMS` wrapper) so repeated `write*` calls over the same in-memory CMS don't re-serialize. Low
  priority because the known call sites are once-per-signature; flag it for the streaming/batch phase.

**Lens: PERF-LEAK** — **no findings.** No goroutines, no timers, no open handles in the API layer.
Every `io.ReadCloser` obtained (`spi.DSSUtilsToByteArrayOfDocument`) is `defer stream.Close()`'d.
`*bytes.Buffer` instances in `CustomContentSigner` are short-lived and GC-eligible; no unbounded
growth.

**Lens: SEC** — 2 findings (both Medium). No exploitable / false-accept finding (this package
performs no verification; see top note).

**C05-SEC-001**
- **Severity:** Medium
- **Category:** SEC
- **Location:** `cms/custom_content_signer.go:40` (`signatureAlgorithmIdentifierHex` table),
  `cms/custom_content_signer_builder.go:28` (`Build`), `cms/cms_builder.go:117` (`CreateCMS`)
- **Evidence:**
  ```go
  "MD2withRSA":   "300d06092a864886f70d0101020500",
  "MD5withRSA":   "300d06092a864886f70d0101040500",
  "RIPEMD160withRSA": "300a06062b24030301020500",
  "SHA1withRSA":  "300d06092a864886f70d0101050500",
  "SHA1withECDSA"/"SHA1withDSA"/"SHA1withPLAIN-ECDSA": ...
  // NewCustomContentSigner(algorithmIdentifier string) accepts any of these with no rejection
  ```
- **Impact:** The public API lets a caller select legacy/weak signature algorithms (MD2, MD5,
  RIPEMD160, and every SHA-1 family) at generation time with no pinning or policy gate. A caller that
  trusts a configuration value or enum, or that is driven by an upstream request, will emit a CAdES/PAdES
  whose digest/signature algorithm the validation engine will later downgrade — but nothing at the
  *generation* boundary refuses it. This is parity with upstream (the BouncyCastle 1.84 finder resolves
  the same names, and the table is KAT-pinned to it), so it is a deliberate porting choice rather than a
  bug; it is reported here because the SEC lens asks whether the API exposes unsafe defaults without
  pinning — and it does, by design.
- **Recommendation:** Keep the table (byte-parity is the contract), but add either (a) a documented
  policy hook on `Builder` / `CustomContentSignerBuilder` (e.g., `SetSignatureAlgorithmPolicy`) so
  deployments can pin to a minimum, or (b) an explicit doc comment on `NewCustomContentSigner` and
  `Build` warning that legacy algorithms are accepted and will fail validation. Low-to-Medium hardening;
  do not remove table entries (would break byte-exactness and the KAT).

**C05-SEC-002**
- **Severity:** Medium
- **Category:** SEC
- **Location:** `cms/precomputed_digest_calculator_provider.go:34` (`Digest`),
  `cms/cms_signer_info_generator_builder.go:243` (`messageDigest`), sink at
  `cms/cms_signed_attribute_table_generator.go:44`
- **Evidence:**
  ```go
  // PrecomputedDigestCalculatorProvider.Digest
  digestAlgorithm, err := enumerations.DigestAlgorithmForOID(...); if err != nil { return []byte{}, nil }
  value, err := p.digestDocument.DigestValue(digestAlgorithm);   if err != nil { return []byte{}, nil }
  // SignerInfoGeneratorBuilder.messageDigest (DigestDocument, digest not configured branch)
  if err != nil { return []byte{}, nil }   // and provider.Digest(...) can also return []byte{}, nil
  // cmsSignedAttributeTableGenerate then:
  table = append(table, cmscore.NewAttribute(cmscore.OIDMessageDigest,
      asn1ber.WriteTLV(asn1ber.TagOctetString, messageDigest)))   // -> 04 00 (empty OCTET STRING)
  ```
- **Impact:** When the precomputed `DigestDocument` does not carry a value for the resolved digest
  algorithm (or the algorithm can't be named), the port returns an **empty** digest with a **nil error**.
  The generator then emits a CMS whose `message-digest` signed attribute is an empty OCTET STRING
  (`04 00`). The caller receives a fully-formed `*CMS` and no error, so it cannot distinguish "signed
  successfully" from "signed over an empty digest." This is **not a false accept** (a verifier computing
  the real content digest will reject it), but it is a silent failure at the API boundary: the library
  hands back an invalid CMS as if it were valid. It is reachable via the port's own documented fallback
  (the `PrecomputedDigestCalculatorProvider` branch), not only via caller error.
- **Recommendation:** Propagate the "no digest value for algorithm" condition as an error (or, if
  matching Java's warn-and-`EMPTY_BYTE_ARRAY` is required, keep it but add a returned alert/error and a
  `// DIVERGENCE, deliberate:` note per PORTING.md) so the caller cannot unknowingly produce an invalid
  CMS. Medium hardening; the fix is a behaviour change from upstream and should be recorded accordingly.

### Tool log (run from `dss/`)

| Tool | Command | Result |
|---|---|---|
| gofmt | `gofmt -l cms` | clean (no files printed) |
| go vet | `go vet ./cms/...` | clean |
| golangci-lint | `golangci-lint run --config=../.github/.golangci.yml ./cms/...` | `0 issues.` |
| (none unavailable) | — | all three tools present and run |

### Open questions / observations (Info)

- **In-memory only (cite, not re-report):** `doc.go` and `docs/compatibility/known-gaps.md` both state
  CMS is in-memory only and the streaming / `DSSResourcesHandlerBuilder` surface is accepted for parity
  but never consulted (`UtilsResourcesHandlerBuilder` is a pass-through). Signing a multi-GB document
  uses memory proportional to the document. This is a known, documented gap — **Info** here, not a
  finding.
- **`CMS.Core()` exposure:** `cms/cms.go:36` exports the internal `*cmscore.CMS`. This is a deliberate,
  documented seam for same-phase callers (`cades`, `spi/validation`) to avoid re-parsing, with no Java
  counterpart. Acceptable; noted so a future reviewer doesn't flag it as an encapsulation leak.
- **`Generate()` / `Build*` panics on nil:** `abstract_cms_generator.go:132`,
  `custom_content_signer_builder.go:28/41`, `cms_signed_document.go:40` all `panic` on nil programmer
  inputs, matching Java's `Objects.requireNonNull`. This is the port's documented convention (Java NPE →
  Go panic) and is **not** a finding; it only matters if a hostile *data* path can reach a nil, which the
  documented call sites do not.
- **Test quality:** strong. KAT test (`custom_content_signer_kat_test.go`) pins the algorithm table and
  the `CMSAlgorithmProtection` / `DERTaggedObject` tag shapes against a real BouncyCastle 1.84 oracle
  (`testdata/signature-algorithm-identifiers.txt`), and guards against the oracle losing a section.
  Round-trip tests cover all three `SignerAttributeV2` alternatives and the `OtherAttributeCertificate`
  drop path. Integration test exercises the full two-step (data-to-sign → real signature) flow and
  verifies the RSA signature over the signed-attributes DER SET. Coverage of the
  `UtilsWriteSignedData*Encoded` `write*` helpers and of `UtilsReplaceSigners` /
  `UtilsReplaceCertificatesAndCRLs` / `UtilsPopulateDigestAlgorithmSet` round-trips is thin in this
  package (likely covered by `cades`'s oracle tests) — worth a cross-check that those three
  re-shaping `Utils*` functions have at least one in-module round-trip assertion.
- **No false-accept / exploitable finding in this unit.** The package builds and re-shapes CMS; it does
  not verify. Verdict correctness is owned by `validation/` + `internal/cmscore` (reviewed separately).

---
---

# Batch 01b review — cades (unit U04a)

- **Unit:** U04a (CAdES signing / levels / timestamps)
- **Scope (exactly 8 files):**
  - `cades/cades_signature.go` (1655)
  - `cades/cades_timestamp_source.go` (990)
  - `cades/cades_level_baseline_lta_timestamp_extractor.go` (939)
  - `cades/cades_level_baseline_b.go` (933)
  - `cades/cades_utils.go` (714)
  - `cades/cades_baseline_requirements_checker.go` (645)
  - `cades/cades_timestamp_message_digest_builder.go` (557)
  - `cades/cades_service.go` (538)
- **Date:** 2026-08-23
- **Depth:** deep
- **Context read (allowed only):** `dss/PORTING.md`, `docs/compatibility/known-gaps.md`,
  `dss/pades/testdata/broadgen/README.md`. Plus, to grade one byte-exactness question
  (not an extra file under review): `cades/testdata/attribute-table-order-oracle.txt`,
  `cades/cades_level_baseline_lta_timestamp_extractor_kat_test.go`, and a run of the LTA KAT tests.

Deliberate conventions treated as **not** findings (per `PORTING.md`): 1:1 Java port shape,
`panic` standing in for unchecked exceptions at interface boundaries, staticcheck-style
categories disabled, SA1019 / legacy crypto (SHA-1, PSS-SHA1 default) allowed where upstream
carries them, byte-identical enum/OID values, stdlib + `golang.org/x` only, `internal/` frozen.

---

## Coverage notes (functions read fully vs grep-signature only)

Read budget: first 300 lines of each file read fully, then every `func` signature grep-mapped,
then the security-relevant bodies read in full.

- **Read fully (body + behavior):** essentially all of `cades_signature.go` (integrity check
  `CheckSignatureIntegrity`, `signedContentForIntegrityCheck`, `ReferenceValidations*`,
  `messageDigestReferenceValidation`, `verifyDigestAlgorithm`, `contentReferenceValidation`,
  `DataToBeSignedRepresentation`, policy/qualifiers, signer-attr v1/v2, alg accessors,
  `DataFoundUpToLevel`, counter-signature walk); all of `cades_timestamp_source.go` (every `Is*`
  OID predicate, `MakeTimestampToken`, `IncorporateArchiveTimestampReferences`,
  `archiveTimestampReferences`, `signedData*CRL/OCSP*`, `GetSignatureSignedDataReferences`,
  `GetCounterSignature*`, the `cadesTS*` helpers); all of
  `cades_level_baseline_lta_timestamp_extractor.go` (`AtsHashIndex`, `VerifiedAtsHashIndex`,
  `verified*HashIndex`, `findCRLMatches`, `unsignedAttributesHashIndex`,
  `ArchiveTimestampV3MessageImprint`, `writeSignedFields`, and the `cadesLTAAttributeHashtable`
  ordering core); all of `cades_level_baseline_b.go` (`SignedAttributes`/`AddSignedAttributes*`,
  policy/qualifiers, `AddSigningCertificateAttribute`, `AddMimeType`, `Time`, `DisplayText`);
  all of `cades_utils.go`; all of `cades_baseline_requirements_checker.go` (every `Has*Profile`,
  `cmsBaselineBRequirements`, `cmsExtendedBESRequirements`, `ContainsLTLevelCertificates`,
  `isSigningCertificateAttributeValid`); all of `cades_timestamp_message_digest_builder.go`
  (every `*MessageDigest`, `writeSignerInfoBytes`, `filterUnauthenticatedAttributes`); all of
  `cades_service.go`.
- **Grep-signature only (trivial OID-equality / pass-through one-liners, each inside an already-read
  300-line window, bodies confirmed to be single `Equal(...)`/assignment returns):** the `Is*`
  attribute predicates in `cades_timestamp_source.go:149-289` and the `New*`/`Set*` accessors in
  `cades_service.go:72-90`. No behavior was judged from signature alone.

**Security-relevant verification performed:** the `ats-hash-index` unsigned-attrs ordering
(`cadesLTAAttributeTableOrder`) was confirmed by running
`go test ./cades/ -run TestCadesLTAAttributeTableOrder` → **PASS** against the committed
`attribute-table-order-oracle.txt` (header: "BouncyCastle 1.84 … Generated by
gen/AtsHashIndexOracle.java"). The `gen/AtsHashIndexOracle.java` generator is **not present in the
repo** (see C04A-SEC-001).

---

## Findings

### SEC

**C04A-SEC-001**
- **Severity:** Medium
- **Category:** SEC
- **Location:** `cades/cades_level_baseline_lta_timestamp_extractor.go:896` (`cadesLTAOIDHashCode`)
- **Evidence:**
  ```go
  func cadesLTAOIDHashCode(oid asn1.ObjectIdentifier) int32 {
      contents := cadesLTAOIDContentOctets(oid)
      index := len(contents)
      hashCode := uint32(index + 1)
      for index > 0 {
          index--
          hashCode *= 257
          hashCode ^= uint32(int32(int8(contents[index])))   // reverse iter, ×257, sign-XOR
      }
  }
  // comment: "…org.bouncycastle.util.Arrays.hashCode … while (--i >= 0) { hc *= 257; hc ^= data[i]; }"
  ```
- **Impact:** `cadesLTAOIDHashCode` implements a reverse-iteration, ×257, sign-XOR hash and
  labels it in-comment as BouncyCastle's `Arrays.hashCode` / `ASN1ObjectIdentifier#hashCode()`.
  The canonical Java array hash is forward ×31. This value seeds
  `cadesLTAAttributeHashtable.index()`, which fixes the `unsignedAttrsHashIndex` order that
  `AtsHashIndex()` / `VerifiedAtsHashIndex()` write into the **archive-timestamp-v3 bytes a TSA
  signs**. The KAT `TestCadesLTAAttributeTableOrder` passes, so the port is *self-consistent*, but
  the oracle fixture's generator (`AtsHashIndexOracle.java`) is not committed, so the claim that
  this ordering matches real BouncyCastle 1.84 cannot be re-verified from the repo. If the
  algorithm is not the real one, archive timestamps *created* by this port would carry a different
  `unsignedAttrsHashIndex` order than Java DSS (a cross-implementation interop break of the
  byte-exactness contract). **Not a false-accept**: `verifiedUnsignedAttributesHashIndex` matches a
  multiset of digests (order-independent), so verification of a Java-made token still passes.
- **Recommendation:** Commit `gen/AtsHashIndexOracle.java` and regenerate the fixture under a
  committed, re-runnable oracle so the ordering pin is reproducible; confirm against the exact
  upstream BouncyCastle version whether `AttributeTable#toASN1EncodableVector()` really yields this
  order — fix `cadesLTAOIDHashCode` and/or the comment if it does not. Escalate to High only if a
  real-BC run shows divergence.

**C04A-SEC-002**
- **Severity:** Info
- **Category:** SEC
- **Location:** `cades/cades_baseline_requirements_checker.go:614` (`isSigningCertificateAttributeValid`)
- **Evidence:**
  ```go
  func (b *BaselineRequirementsChecker) isSigningCertificateAttributeValid(signerInformation *cmscore.SignerInfo) bool {
      certificateRefs := b.Signature().CertificateSource().SigningCertificateRefs()
      if utils.IsCollectionNotEmpty(certificateRefs) {
          ...
          if certDigest.Algorithm() != "" {
              digestAlgorithm := certDigest.Algorithm()
              if enumerations.DigestAlgorithmSHA1 == digestAlgorithm { ... } else { ... }
          }   // else: check silently skipped
      }
      return true
  }
  ```
- **Impact:** When the signing-certificate ref's digest algorithm is empty/undetermined, the
  req (h)/(i) SHA1-vs-SHA256 attribute check is **waived** (falls through to `return true`). This
  is a **level-classification** nuance (B-B vs BES), not a cryptographic verification gate; it
  mirrors upstream's `if (digestAlgorithm != null)` guard. No false-accept of a signature.
- **Recommendation:** None required for parity; optional hardening: report an undetermined digest
  algorithm as an alert rather than a silent pass (would need a `// DIVERGENCE, deliberate:` + DESIGN
  entry per the upstream-tracking rule).

**C04A-SEC-003**
- **Severity:** Info
- **Category:** SEC
- **Location:** `cades/cades_signature.go:1261` (`contentReferenceValidation`)
- **Evidence:**
  ```go
  func (s *Signature) contentReferenceValidation(_ model.DSSDocument, _ *cmscore.SignerInfo) *model.ReferenceValidation {
      contentValidation := model.NewReferenceValidation()
      contentValidation.SetType(enumerations.DigestMatcherTypeContentDigest)
      return contentValidation   // Found stays false
  }
  ```
- **Impact:** The RFC 5652 "message-digest absent" content-digest path always reports *not found*
  (no BC `SignerInformation#getContentDigest()` equivalent). Reachable only for a non-ETSI CMS
  signature with no signed attributes; the outcome matches Java's failure branch, so it is a
  false-**reject**-leaning narrowing, not a false-accept. Documented in the file-header DEVIATION
  note.
- **Recommendation:** None (documented, deliberate, parity-safe). Cited as Info.

**C04A-SEC-004**
- **Severity:** Info
- **Category:** SEC
- **Location:** `cades/cades_service.go:1576` (guard) · weak-alg gate summary (`cades_utils.go:196`, `cades_signature.go:949`)
- **Evidence:**
  ```go
  cades_utils.go:196:   if digestAlgorithm == enumerations.DigestAlgorithmSHA1 { // ESSCertID v1, only when caller chose SHA-1
  cades_signature.go:949: return spi.DSSASN1UtilsAlgorithmIdentifierForDigest(enumerations.DigestAlgorithmSHA1) // PSS ASN.1 DEFAULT, read-side only
  cades_service.go:1576:  if !timestamp.IsProcessed() { panic("Timestamp token must be validated first !") }
  ```
- **Impact:** Legacy/weak algorithm selection is **caller-gated only**: SHA-1 appears solely in
  (a) the ESSCertID-v1 branch taken when the caller selected SHA-1, and (b) the read-side PSS hash
  default (correct per RFC 4055, and only *reports* the digest of an existing signature). No
  non-legacy path in these 8 files auto-selects a weak digest for a new signature, and
  `AddExternalTimestamp` refuses an unvalidated token. Gate is sound.
- **Recommendation:** None (SA1019 / legacy-compat is explicitly allowed by `PORTING.md`).

> **STD lens (idioms / error handling / defer / context / sync / doc comments):** no blocking
> findings. `panic`-as-exception is the deliberate, documented porting convention (not an idiom
> defect); `archiveTimestampDataV2Try` closes its digest writer via `defer`
> (`cades_timestamp_message_digest_builder.go:379`); there is no `context.Context` threading
> because upstream Java has none (parity, not a gap). The panic-on-hostile-input concern is folded
> into Open Questions rather than scored.

> **PERF-LEAK lens:** no findings. No unbounded growth and no unclosed resources in scope (the
> sole stream in these files is closed by `defer`).

### PERF

**C04A-PERF-001**
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `cades/cades_timestamp_source.go:949` (`cadesTSAddReferences`)
- **Evidence:**
  ```go
  func cadesTSAddReferences(referenceList *[]*validation.TimestampedReference, referencesToAdd []*validation.TimestampedReference) {
      for _, candidate := range referencesToAdd {
          found := false
          for _, existing := range *referenceList { if existing.Equals(candidate) { found = true; break } }
          if !found { *referenceList = append(*referenceList, candidate) }
      }
  }
  ```
- **Impact:** Dedup is O(n·m) with a linear `Equals` scan, invoked per unsigned attribute and per
  previous timestamp inside `IncorporateArchiveTimestampReferences` / `unsignedAttributesReferences`
  / `encapsulatedReferencesFromTimestamps`. Bounded by a single SignerInfo's attribute + timestamp
  counts (a handful in practice), so real-world cost is small; a hostile CMS with many unsigned
  attributes would push it toward quadratic.
- **Recommendation:** Optional: dedup via a map keyed on the reference's identity fields (type + ID)
  instead of a linear scan. Low priority.

**C04A-PERF-002**
- **Severity:** Low
- **Category:** PERF-MEM
- **Location:** `cades/cades_level_baseline_lta_timestamp_extractor.go:911` (`cadesLTAOIDContentOctets`)
- **Evidence:**
  ```go
  func cadesLTAOIDContentOctets(oid asn1.ObjectIdentifier) []byte {
      element, _, err := asn1ber.Parse(asn1ber.EncodeOID(oid))   // encode + parse round-trip
      if err != nil { return nil }
      return element.Content()
  }
  ```
- **Impact:** `EncodeOID` → `Parse` → `Content()` is an encode+parse round-trip just to strip
  tag+length; it runs once per OID inside `cadesLTAOIDHashCode`, which is invoked per attribute in
  the ≥2-attribute ordering path. Minor avoidable allocation/parse per OID.
- **Recommendation:** Derive the content octets from the already-encoded OID directly (skip the
  re-parse), or cache per-OID content. Low priority.

**C04A-PERF-003**
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `cades/cades_baseline_requirements_checker.go:221` (`ContainsLTLevelCertificates`)
- **Evidence:**
  ```go
  func (b *BaselineRequirementsChecker) ContainsLTLevelCertificates() bool {
      var signedDataCertificates []*model.CertificateToken
      for _, certificateDER := range cmsDoc.Certificates() {
          certificateToken, err := spi.DSSASN1UtilsCertificate(certificateDER)  // re-parse every cert, per call
          ...
      }
      for _, timestampCertificateSource := range timestampCertificateSources {
          if !cadesBaselineContainsAnyCertificate(signedDataCertificates, ...) { return false } // O(n·m) Equals
      }
  }
  ```
- **Impact:** Re-parses every `SignedData.certificates` DER into a fresh `[]*CertificateToken` on
  each invocation and matches with O(n·m) `Equals`; the result is not cached, so repeated level
  checks re-do the parse. Bounded by cert count, but wasteful when the checker is consulted more
  than once.
- **Recommendation:** Cache the parsed signed-data cert list on the checker (or the signature) and
  match by identifier. Low priority.

---

## Tool log (run from `dss/`)

| Command | Result |
|---|---|
| `gofmt -l cades` | **no output** (clean) |
| `go vet ./cades/...` | exit `0` (clean) |
| `golangci-lint run --config=../.github/.golangci.yml ./cades/...` | **0 issues** (tool was available; config passed explicitly) |
| `go test ./cades/ -run TestCadesLTAAttributeTableOrder -count=1 -v` | **PASS** (0.00s) — used to grade C04A-SEC-001 |

No unavailable tools. No tests weakened or deleted.

---

## Open questions

1. **C04A-SEC-001 provenance:** Was `attribute-table-order-oracle.txt` actually generated by real
   BouncyCastle 1.84 `AttributeTable#toASN1EncodableVector()`, or by a hand-written Java replica of
   `cadesLTAOIDHashCode`? The generator (`AtsHashIndexOracle.java`) is not in the repo, so this
   cannot be settled from the checkout. Needs a real-BC cross-check (or a committed, re-runnable
   oracle) before the byte-exactness claim is considered independently verified.
2. **Panic on hostile input (Low, not scored):** `UtilsEvidenceRecordIncorporationType`
   (`cades_utils.go:568`) `panic`s on an OID that is neither internal nor external evidence record.
   It is reached only after the base layer's `IsEvidenceRecord()` classifies the attribute as an
   evidence record, and it mirrors Java's `UnsupportedOperationException`. If `IsEvidenceRecord()`
   can be satisfied by an OID outside that fixed set, a hostile CMS could reach the panic — worth a
   one-line confirmation that the gate's OID set is exactly {internal, external}.
3. **`ReferenceValidationsForSignerInformation` caching (`cades_signature.go:1090`):** the cache is
   populated on first call and reused for every subsequent signer-info query. This is correct today
   because `CheckSignatureIntegrity` is always the first caller with the (possibly recreated) signer
   info. Confirm no future caller can prime the cache with a different `signerInformationToCheck`
   before `CheckSignatureIntegrity` runs.

---

## dss/cades remaining + extension (unit U04b)

- **Unit:** U04b (CAdES remaining files + `cades/extension`)
- **Scope:** all `*.go` in `dss/cades/` not covered by U04a (the 8 largest files), plus all of
  `dss/cades/extension/`, plus `testdata/crossgen/main.go` and `testdata/bytecmp/main.go`.
- **Date:** 2026-08-23

### Files read

**Non-test (read fully):**

| File | Lines |
|---|---|
| `cades/doc.go` | 14 |
| `cades/cades_attribute.go` | 154 |
| `cades/cades_attribute_identifier.go` | 41 |
| `cades/cades_attribute_order_comparator.go` | 36 |
| `cades/cades_certificate_source.go` | 25 |
| `cades/cades_crl_source.go` | 24 |
| `cades/cades_ocsp_source.go` | 24 |
| `cades/cades_counter_signature_builder.go` | 338 |
| `cades/cades_counter_signature_parameters.go` | 54 |
| `cades/cades_diagnostic_data_builder.go` | 78 |
| `cades/cades_embedded_evidence_record_builder.go` | 347 |
| `cades/cades_embedded_evidence_record_helper.go` | 108 |
| `cades/cades_evidence_record_digest_builder.go` | 279 |
| `cades/cades_evidence_record_incorporation_parameters.go` | 17 |
| `cades/cades_evidence_record_scope_finder.go` | 82 |
| `cades/cades_level_baseline_lt.go` | 279 |
| `cades/cades_level_baseline_lta.go` | 169 |
| `cades/cades_level_baseline_t.go` | 163 |
| `cades/cades_sig_properties.go` | 43 |
| `cades/cades_signature_extension.go` | 226 |
| `cades/cades_signature_identifier_builder.go` | 108 |
| `cades/cades_signature_integrity_validator.go` | 107 |
| `cades/cades_signature_parameters.go` | 193 |
| `cades/cades_signature_policy_store_builder.go` | 316 |
| `cades/cades_signature_scope_finder.go` | 139 |
| `cades/cades_signed_attributes.go` | 22 |
| `cades/cades_timestamp_parameters.go` | 33 |
| `cades/cades_unsigned_attributes.go` | 128 |
| `cades/cms_document_analyzer.go` | 224 |
| `cades/cms_document_analyzer_factory.go` | 39 |
| `cades/cms_document_validator.go` | 81 |
| `cades/cms_document_validator_factory.go` | 39 |
| `cades/cms_for_cades_builder_helper.go` | 216 |
| `cades/evidence_record_production_comparator.go` | 52 |
| `cades/time_stamp_token_production_comparator.go` | 104 |
| `extension/cades_document_extender.go` | 84 |
| `extension/cades_document_extender_factory.go` | 40 |
| `extension/doc.go` | 9 |
| `testdata/crossgen/main.go` | 194 |
| `testdata/bytecmp/main.go` | 92 |

**Tests (skipped for quality/coverage):** `cades_byte_exactness_test.go`,
`cades_level_baseline_b_kat_test.go`, `cades_level_baseline_lta_timestamp_extractor_kat_test.go`,
`cades_downstream_cross_validation_test.go`, `cades_upstream_cross_validation_test.go`,
`cms_document_validator_smoke_test.go`.

**Adjacent files (read for context, not re-reported):** `cades_level_baseline_b.go`
(`AddMimeType` at 678-693, `counterSignature` field at 102-108), `cades_service.go`
(`CAdESCounterSignatureBuilder` at 492-497), `document/abstract_signature_parameters.go`
(`Equals` at 207-226), `cms/cms_utils.go` (`UtilsReplaceUnsignedAttributes` at 365-383),
`spi/validation/scope/evidence_record_scope_finder.go`, `spi/validation/identifier/
signature_attribute_identifier.go`, `model/identifier.go` (`IdentifierBase.Equals` at 123-139).

### Findings

#### SEC

**C04B-SEC-001**

- **Severity:** High
- **Category:** SEC
- **Location:** `cades/cades_counter_signature_builder.go:331-337` (`InitCMSBuilderHelper`);
  root cause at `cades/cms_for_cades_builder_helper.go:165` (`InitCAdESProfile` reads
  `h.counterSignature`, always `false`); symptom at `cades/cades_level_baseline_b.go:680-681`
  (adjacent to U04a scope — `AddMimeType` does not skip the attribute)
- **Evidence:**
  ```go
  // cades_counter_signature_builder.go:331-337 — SetCounterSignature is never called
  func (b *CounterSignatureBuilder) InitCMSBuilderHelper(contentToSign model.DSSDocument,
      signatureParameters *SignatureParameters, contentSigner cms.ContentSigner) *CMSForCAdESBuilderHelper {
      helper := NewCMSForCAdESBuilderHelper(contentToSign, signatureParameters, contentSigner)
      helper.SetTrustedCertificateSource(b.certificateVerifier.TrustedCertSources())
      return helper   // ← helper.counterSignature stays false
  }
  ```
  ```go
  // cms_for_cades_builder_helper.go:164-166 — profile gets counterSignature=false
  func (h *CMSForCAdESBuilderHelper) InitCAdESProfile() *LevelBaselineB {
      profile := NewLevelBaselineBWithDocument(h.DocumentToSign)
      profile.SetCounterSignature(h.counterSignature)   // false
      return profile
  }
  ```
  ```go
  // cades_level_baseline_b.go:680-681 — mime-type attribute IS added (Java skips it)
  if b.counterSignature {  // always false in counter-signature path
      return signedAttributes, nil
  }
  ```
- **Impact:** Upstream Java `CAdESLevelBaselineB#addMimeType` checks
  `parameters instanceof CAdESCounterSignatureParameters` and returns the signed attributes
  unchanged (no `mimeType` attribute). The Go port's `CounterSignatureBuilder` never sets the
  `counterSignature` flag on the `CMSForCAdESBuilderHelper`, so `LevelBaselineB.AddMimeType`
  falls through and adds a `mimeType = application/octet-stream` signed attribute to every
  counter-signature. This is a byte-level divergence from upstream (the counter-signature
  carries an extra signed attribute Java does not produce) and a semantic error: a
  counter-signature signs a `SignerInfo`, not a content document, so a `mimeType` attribute
  is inappropriate. The KAT test
  `TestCAdESLevelBaselineBCounterSignatureSuppressesMimeType`
  (`cades_level_baseline_b_kat_test.go:477-489`) tests the mechanism in isolation by calling
  `profile.SetCounterSignature(true)` directly, which passes, but does not exercise the
  integration path (`CounterSignatureBuilder` → `CMSForCAdESBuilderHelper` → `LevelBaselineB`),
  so CI does not catch the regression.
- **Recommendation:** In `CounterSignatureBuilder.GenerateCounterSignature`
  (`cades_counter_signature_builder.go:240`), call
  `cmsBuilderHelper.SetCounterSignature(true)` after `InitCMSBuilderHelper` returns, so the
  flag is set before `CreateSignerInfoGenerator` → `CAdESProfile` → `InitCAdESProfile` reads
  it. Add an integration test that calls `CounterSignatureBuilder.AddCounterSignature` (or
  `GenerateCounterSignature`) and asserts the resulting `SignerInfo` signed attributes do not
  contain `OIDIdAaEtsMimeType`.

#### STD

**C04B-STD-001**

- **Severity:** Low
- **Category:** STD
- **Location:** `cades/cades_embedded_evidence_record_builder.go:148-149`
  (`getUnsignedAttributeToEmbed`)
- **Evidence:**
  ```go
  return nil, exception.NewIllegalInputException(fmt.Sprintf(
      "Unable to embed the parallel evidence record. Expected type '%s', obtained type '%s'.",
      lastUnsignedAttribute.ASN1Oid().String(), expectedEvidenceRecordAttributeType.String()))
  ```
- **Impact:** The format string reads "Expected type '%s', obtained type '%s'" but the
  arguments are passed in reverse order: `lastUnsignedAttribute.ASN1Oid()` (the obtained/actual
  OID) is printed in the "Expected" slot and `expectedEvidenceRecordAttributeType` (the expected
  OID) is printed in the "obtained" slot. Upstream Java has the correct order
  (`expectedEvidenceRecordAttributeType` first, `lastUnsignedAttribute.getASN1Oid()` second).
  A user reading the error message would see the two OIDs swapped, making diagnosis harder.
- **Recommendation:** Swap the two `fmt.Sprintf` arguments to
  `expectedEvidenceRecordAttributeType.String(), lastUnsignedAttribute.ASN1Oid().String()`.

**C04B-STD-002**

- **Severity:** Info
- **Category:** STD
- **Location:** `cades/testdata/crossgen/main.go:53-56`, `cades/testdata/bytecmp/main.go:44`
- **Evidence:**
  ```go
  // crossgen/main.go:53
  signerEntry, err := loadKeyEntry(filepath.Join(selfDir, "signer_rsa.p12"), "testpassword")
  // crossgen/main.go:55
  filepath.Join(selfDir, "tsa_ec.p12"), "PKCS12", "testpassword", "", "testpassword")
  // crossgen/main.go:63
  tspSource.SetTsaPolicy("1.2.3.4.5.6.7.8.9")
  ```
- **Impact:** The test-data generators use a hardcoded password `"testpassword"` for both the
  signer and TSA PKCS#12 key stores, and an unregistered test TSA policy OID
  `"1.2.3.4.5.6.7.8.9"`. These are test-only fixtures (committed to the repo, used exclusively
  by the cross-validation harness), and the comments document them as such. No security impact —
  the key stores are test keys with no private value, and the OID is explicitly noted as an
  "unregistered test OID" matching DSS's own `KeyEntityTSPSource` unit tests. Recorded for
  completeness.
- **Recommendation:** None.

### Lens summary

| Lens | Findings |
|---|---|
| **STD** | C04B-STD-001 (Low), C04B-STD-002 (Info). No error-handling, defer, context, sync, or doc-comment issues found. The `panic` sites in factory `Create` methods and `Objects.requireNonNull` mirrors are deliberate Java-1:1 ports. The `cadesLikeSignature` interface in `cades_diagnostic_data_builder.go` is a well-documented workaround for Go's lack of covariance. |
| **PERF** | No findings. The O(n²) bubble sort in `cadesUnsignedAttributesSort` is deliberate (documented: the comparator is not a total order, and unsigned attribute tables are small, typically 5-15 entries). The O(n²) nested loop in `CMSDocumentAnalyzer.AppendExternalEvidenceRecords` is bounded by the number of signatures in one CMS, which is small. The `cadesLTAAttributeTableOrder` hashtable reproduction is O(n) amortised. No memory-growth or allocation issues found. |
| **SEC** | C04B-SEC-001 (High). Counter-signature mime-type suppression is silently lost in the integration path. Scope-finder logic (`SignatureScopeFinder`, `EvidenceRecordScopeFinder`) correctly uses pointer identity for CMS comparison and does not falsely match signatures. Evidence-record digest handling (`EvidenceRecordDigestBuilder`) correctly strips the last ER attribute before computing the message imprint. Document validator false-accept surface is in `cades_signature.go` (U04a scope), not in the files reviewed here. |

### Tool log (run from `dss/`)

| Command | Result |
|---|---|
| `gofmt -l cades cades/extension` | **no output** (clean) |
| `go vet ./cades/...` | exit `0` (clean) |
| `golangci-lint run --config=../.github/.golangci.yml ./cades/...` | **0 issues** (clean) |

No unavailable tools. No tests weakened or deleted.

### Open questions

1. **C04B-SEC-001 integration test:** The KAT test
   `TestCAdESLevelBaselineBCounterSignatureSuppressesMimeType` tests `SetCounterSignature(true)`
   in isolation but the `CounterSignatureBuilder` path never calls it. After the fix, an
   integration test through `CounterSignatureBuilder.GenerateCounterSignature` (or
   `AddCounterSignature`) should assert the absence of `OIDIdAaEtsMimeType` in the
   counter-signer's signed attributes. Does the existing `counterSig.p7m` upstream fixture
   (used by `cms_document_validator_smoke_test.go`) carry a counter-signature that was produced
   by Java (and therefore lacks the `mimeType` attribute)? If so, it could serve as a
   read-side cross-check.
2. **`EmbeddedEvidenceRecordBuilder.NewEmbeddedEvidenceRecordBuilder`** calls
   `validation.NewCertificateVerifierBuilder(certificateVerifier).BuildOfflineCopy()`
   (`cades_embedded_evidence_record_builder.go:38`), while
   `CounterSignatureBuilder.NewCounterSignatureBuilder` stores the verifier directly
   (`cades_counter_signature_builder.go:47`). The `BuildOfflineCopy` in the evidence-record
   builder is a defensive copy that prevents the caller's verifier from being mutated during
   validation. The counter-signature builder does not make this copy. Confirm this is
   intentional (the counter-signature builder does not call `Validate()` on the verifier) or
   whether it is a subtle asymmetry worth aligning.
