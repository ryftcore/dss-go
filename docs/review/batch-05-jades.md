# Batch 05 — JAdES core (U17a)

- **Unit:** U17a
- **Scope:** `dss/jades` core — the JWS/JWT signature object, the B-B signing header builder, and the
  JSON/base64/digest utility layer.
  Focus files: `jades_signature.go` (1354), `jades_level_baseline_b.go` (1119),
  `dss_json_utils.go` (972) — 3445 lines total.
- **Date:** 2026-08-23
- **Depth:** deep-budgeted (hard read budget ~8k lines incl. grep; actual ≈ 3.4k focus + ≈ 1.6k
  supporting cross-check reads + signature greps, well under budget)
- **Context:** Go port of Java DSS 6.5.RC1. `dss/PORTING.md` read first — 1:1 Java port, staticcheck
  style rules disabled, SA1019 allowed, byte-ident enum values, stdlib-first, no cgo.
  `docs/compatibility/known-gaps.md` read — listed gaps not re-reported (notably: MD2/Whirlpool,
  streaming-CMS, evidence-record analyzers unregistered). `internal/jose` (SEC-001 stale validSignature
  cache, SEC-002 HS* dead code, PERF-001 unbounded header, SEC-003 lenient JWK base64) already reviewed —
  **not re-reported**; this unit focuses on the jades layer that sits on top of it. The level-check /
  cert-selection / timestamp-validation logic (B-LT/B-LTA, `jades_level_baseline_lt.go`,
  `jades_timestamp_source.go`, `jades_certificate_source.go`) is U17b and is **out of scope** here; the
  `jades_level_baseline_b.go` file in scope is the **signing-side** protected-header builder, not the
  level-validation side.

## Files read

| File | Coverage |
|---|---|
| `jades/jades_signature.go` (1354) | **FULL** (1–250, 250–560, 560–860, 860–1130, 1130–1354) + all 70 `func` signatures grepped. Covers `NewSignature`, `SignatureAlgorithm`, `CheckSignatureIntegrity`, `ReferenceValidations`, `signingInputReferenceValidation`, `jadesSignaturePayload`, the detached-content payload builders, `referenceValidationsByUriHashMechanism`, `isDocumentDigestMatch`, `counterSignatureReferenceValidation`, `keyBindingSignatureReferenceValidation`, `DataFoundUpToLevel`, `OriginalDocuments`. |
| `jades/jades_level_baseline_b.go` (1119) | **FULL** (1–250, 250–700, 700–1119) + all 44 `func` signatures grepped. Covers `NewLevelBaselineB`, `SignedProperties` ordering, all `incorporate*` methods (alg/cty/kid/x5u/x5t#S256/x5c/typ/exp/b64/iat/sigT/srCms/sigPl/srAts/adoTst/sigPId/sigD/crit), the three sigD mechanism builders, `signedDataDigests`, `assertPayloadEncodingValid`, `PayloadBytes`. |
| `jades/dss_json_utils.go` (972) | **FULL** (1–250, 250–503, 503–668, 668–935, 935–972) + all 62 `func` signatures grepped. Covers base64url/url-safe/utf8 predicates, crit-header sets, `DSSJsonUtilsDigest`, `OidObject`, `TstContainer`, `DocumentOctets`, `IsJsonDocument`, `ExtractJAdESCounterSignature`, `ValidateAgainstJAdESSchema`, `ParseEtsiUComponent`, `SigningInputBytes`, the `GetAs*`/`To*` coercion helpers, `ParseBase64UrlEncoded`/`ParseJSONString`/`ParseJSONStringToMap`. |
| `jades/jades_signature_integrity_validator.go` (60) | **FULL** — the `Verify` override that calls `jose.JWS.VerifySignature`. |
| `jades/jades_timestamp_message_digest_builder.go` 120–175 | **REGION** — the one *non*-recovered caller of the panicking detached-doc resolvers. |
| `jades/jws.go` 40–60, 185–230 | **REGION** — `ProtectedHeaderValueAs*` (map reads, no re-parse) and the `CheckCritOverride` install. |
| Supporting cross-check (for the verify boundary) | `spi/signature_integrity_validator.go` 59–160 (`Validate`/`isSignatureIntact`), `internal/jose/jws.go` 60–330 (`VerifySignature`/`SigningInputBytes`/`resolveAlgorithm`/`CheckCrit`), `internal/jose/sigalg.go` 1–160 (alg table + key-type mismatch → error), `internal/jose/number.go` 1–120 (`Int64` overflow-safe), `internal/jose/parser.go` 1–110 (JSON state machine), `internal/jose/base64url.go` 80–150 (lenient decode), `internal/jose/headers.go` 28–130, `model/reference_validation.go` (`IsIntact` defaults false), `enumerations/signature_algorithm.go` 862–905, `spi/dss_utils.go` 1470–1520 (EDDSA-by-length), `format.go` 1–80 + `recovered()` (facade panic→error). |
| Tests (coverage only, not read fully) | `dss_json_utils_kat_test.go` (13 tests), `jades_sign_kat_test.go` (3 tests), `jades_root_chunk_test.go` (16 tests) — confirm KAT coverage of the protected-header bytes, serialization, scalar coercion, crit sets, and the sign path. |

## Security posture (top priority) — VERIFIED CLEAN on the attack surface

This was the primary target. The headline is that **no exploitable defect and no false-accept was
found** in the JWS verification / digest / JSON / base64 surface of the three in-scope files. Specifically
verified:

- **No false-accept in the verify path.** `jades.SignatureIntegrityValidator.Verify` →
  `jose.JWS.VerifySignature` → `resolveAlgorithm`. `alg=none` is rejected in `resolveAlgorithm`
  (`internal/jose/jws.go`, `name == AlgorithmNone` → `ErrUnsupportedAlgorithm`), so a `none` JWS can
  never verify. A **key-type / key-algorithm mismatch returns an error, not a boolean** (`sigalg.go`:
  `hmacAlgorithm.verify` does `secret, ok := key.([]byte)`; `rsaAlgorithm.verify` does
  `pub, ok := key.(*rsa.PublicKey`), each mismatch → `return false, fmt.Errorf(...)`). That error
  propagates up: `spi.SignatureIntegrityValidatorBase.isSignatureIntact` returns `(false, err)` →
  `Validate` records the error and does **not** return that candidate → `certificateValidity` stays `nil`
  → `signingInputReferenceValidation` sets `SetIntact(certificateValidity != nil)` = **false**.
  `ReferenceValidation.IsIntact()` defaults to `false` (zero value), so every early-return / failure
  path in `signingInputReferenceValidation` (empty header, unknown alg, digest error) is **fail-closed**.
- **No alg-confusion / RS256→HS256 key swap.** An attacker who sets `alg=HS256` while the supplied
  "key" is actually an RSA public key hits `hmacAlgorithm.verify`'s `key.([]byte)` assertion, which
  fails (a `crypto.PublicKey` is never a `[]byte`) → error → not intact. Conversely a `RS256` claim
  with an HMAC-style secret fails the `*rsa.PublicKey` assertion. Both directions are blocked by the
  Go type assertion, not by a shared secret. Clean.
- **No wrong-cert false-accept (in scope).** The integrity validator iterates *all* candidates and
  returns the first one that cryptographically verifies (`spi.Validate`); candidate *selection* itself
  (x5t#S256 / x5c matching) lives in `jades_certificate_source.go` / `spi` (U17b / already-reviewed),
  not in these three files. What is in scope — that a non-matching candidate cannot be reported intact —
  is correct: only a candidate whose key actually verifies the signature is returned.
- **JSON parsing is panic-safe and error-based.** `jose.ParseJSON`/`ParseJSONAny`
  (`internal/jose/parser.go`) is a hand-written recursive-descent state machine (not `encoding/json`);
  malformed input returns a `*ParseError`, never a panic, and the ordered path rejects duplicate member
  names. Every jades JSON entry point (`DSSJsonUtilsParseEtsiUComponent`, `ValidateAgainstJAdESSchema`,
  `ParseBase64UrlEncoded`, `ParseJSONString(ToMap)`, `IsJsonDocument`) checks `err` and falls back to a
  null/empty/`false` sentinel. A hostile JWS `etsiU` array or protected header therefore degrades to
  "not found / not a map", not a crash.
- **Numeric overflow is saturated, not a panic.** `jose.Number.Int64()` (`internal/jose/number.go`)
  saturates Double at `Min/MaxInt64` and takes BigInteger low-64-bits (JVM `longValue()` parity), so a
  hostile `iat`/`exp`/`sigT`/`noticeNumbers` huge number cannot overflow-panic in
  `DSSJsonUtilsToNumericDateWithHeaderName` or `jadesSignatureSPUserNotice`.
- **base64url handling is the documented leniency, not a new bug.** `DSSJsonUtilsFromBase64Url` →
  `jose.Base64URLDecode` is the *already-reviewed* SEC-003 lenient decoder (skips non-alphabet bytes,
  optional padding). The jades layer does **not** re-introduce a new leniency: `DSSJsonUtilsIsBase64UrlEncoded`
  does the strict per-byte alphabet check the decoder lacks, exactly as the `dss_json_utils.go` comment
  documents ("the lenient decoder in internal/jose must not be replaced by a strict one"). Not a new finding.
- **Timestamp (B-LT) counter-signature *validation* is out of scope** (U17b: `jades_level_baseline_lt.go`,
  `jades_timestamp_source.go`). In scope, `jades_level_baseline_b.go` only *builds* the signed
  `adoTst` header (signing side), and `jades_signature.go`'s `TimestampSource()` is a lazy source
  accessor. The "is the counter-signature present + valid" decision is U17b territory and not re-reported.
- **No uncaught panic-DoS on hostile input.** The panics that do exist (`incorporatedPayload`,
  `payloadForHttpHeadersMechanism`, `payloadForObjectIdByUriMechanism`, `signedDocumentsForUris`,
  `SignedDocumentsByHTTPHeaderName`, `sdAlg`, `AddExternalTimestamp`, `DSSJsonUtilsOidObjectFromURI`)
  are either (a) Java-parity `IllegalArgumentException`/`UnsupportedOperationException` equivalents, (b)
  guarded by a local `defer recover()` in `jadesSignaturePayload`/`OriginalDocuments`, or (c) pre-checked
  (`assertSignaturePolicyValid` rejects empty `Id` before `OidObjectFromURI` can panic). All are reached
  under the facade `format.go`'s `recovered()` wrapper, which converts a panic into an `error`
  (`errors.As(err, *model.DSSError)`), so a hostile JWS yields an error verdict, not a process crash.
- **Signing-input bytes are RFC 7797-correct and the two variants agree.** The in-scope
  `DSSJsonUtilsSigningInputBytes` (used only for the informational DTBSR digest reference) and the
  verifying `jose.JWS.SigningInputBytes` both produce `ASCII(BASE64URL(hdr)||'.'||BASE64URL(payload))`
  for `b64=true` and `ASCII(BASE64URL(hdr))||'.'||rawPayload` for `b64=false` (raw payload never routed
  through a string for the unencoded case). No divergence between the digest-reference and the verify
  path.

## Findings

**SEC:** *no findings* — the verification / false-accept / alg-confusion / JSON / base64 / cert-selection
attack surface was verified clean (see posture above).

### J17A-PERF-001
- **Severity:** Medium
- **Category:** PERF-BIGO
- **Location:** `dss/jades/jades_signature.go:928` (`referenceValidationsByUriHashMechanism`), `:1016`
  (`detachedDocumentByDigest`), `:1097` (`isDocumentDigestMatch`)
- **Evidence:**
  ```go
  for _, signedDataName := range signedDataOrder {                 // n = sigD.pars, attacker-controlled
      ...
      detachedDocument = s.detachedDocumentByDigest(digestAlgorithm, expectedDigest, signedDataName, detachedDocuments)
      ...
  }
  // detachedDocumentByDigest:
  for _, detachedDocument := range detachedContent {               // m = detached docs
      if s.isDocumentDigestMatch(detachedDocument, ...) { return detachedDocument }
  }
  // isDocumentDigestMatch (non-digest, encoded path):
  base64UrlEncodedDocument, _ := DSSJsonUtilsToBase64UrlDocument(document)   // full-doc copy + base64
  computedDigestValue, _ = spi.DSSUtilsDigest(digestAlgorithm, []byte(base64UrlEncodedDocument)) // full-doc digest
  ```
- **Impact:** For the `OBJECT_ID_BY_URI_HASH` mechanism the resolver scans **every** detached document and
  computes a full-document base64url-encode **and** a full digest per `(sigD entry, document)` pair, only
  stopping on the first digest match. A hostile JWS can carry an arbitrarily large `sigD.pars`/`hashV`
  array (n, attacker-controlled) and, when no entry matches, forces O(n·m) full-document base64+digest
  operations over the caller's documents. `m` is caller-bounded (the documents the caller chose to
  validate), so it is O(n) with a large per-iteration constant rather than true O(n²); the CPU cost
  scales directly with hostile-JWS size. Faithful to upstream (Java's
  `getReferenceValidationsByUriHashMechanism`/`getDetachedDocumentByDigest` do the same nested scan +
  `toBase64Url`+digest), so not a divergence — but a real, avoidable CPU cost on the validation path.
- **Recommendation:** Pre-compute and cache each detached document's `base64url` form and digest once per
  `(doc, digestAlgorithm)` before the sigD loop, then compare by value (the `b64=false` path already
  digests the raw value, so only the encoded path needs the cache). Keep the first-match / name-fallback
  semantics intact.

### J17A-PERF-002
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `dss/jades/jades_signature.go:850` (`SignedDocumentsByHTTPHeaderName`), `:908`
  (`signedDocumentsForUris`), `:1031` (`jadesSignatureDetachedDocumentByName`)
- **Evidence:**
  ```go
  for _, signedDataName := range signedDataUriList {     // n names
      found := false
      for _, document := range s.DetachedContents() {    // linear scan over m docs, case-insensitive
          if utils.AreStringsEqualIgnoreCase(signedDataName, document.Name()) { found = true; ... }
      }
      if !found { panic("The detached content for a signed data with name '...
  }
  ```
- **Impact:** O(n·m) case-insensitive string comparison to bind sigD names to detached documents, for both
  the `HTTP_HEADERS` and `OBJECT_ID_BY_URI` mechanisms. Small for realistic n, m, but it re-scans the
  whole detached list per name and is O(n·m) rather than O(n+m). Faithful to upstream (Java does the same
  nested `AreStringsEqualIgnoreCase` scan), so not a divergence.
- **Recommendation:** Build a lower-cased name→`[]DSSDocument` index once (same-name docs are legal, so
  map to a slice) and resolve each sigD name in O(1); preserve the "not found" error/panic contract.

### J17A-STD-001
- **Severity:** Low
- **Category:** STD
- **Location:** `dss/jades/jades_signature.go` — `CheckSignatureIntegrity` (:608), `ReferenceValidations`
  (:646), `CertificateSource`/`CRLSource`/`OCSPSource`/`TimestampSource` (:156–191), `EtsiUHeader` (:557)
- **Evidence:**
  ```go
  func (s *Signature) CheckSignatureIntegrity() {
      if s.jadesCachedCryptoVerification != nil { return }
      verification := signature.NewCryptographicVerification()
      s.jadesCachedCryptoVerification = verification   // check-then-set, no lock
      ...
  }
  // CertificateSource / ReferenceValidations / EtsiUHeader use the same lazy check-then-set on
  // s.OfflineCertificateSource() / s.CachedReferenceValidations() / s.etsiUHeader
  ```
- **Impact:** A `*Signature` exposes several lazily-initialised, check-then-set fields with no
  synchronization. If one `*Signature` is validated from more than one goroutine concurrently, these are
  data races (torn reads / duplicate builder construction). This is parity with Java's own unsynchronised
  lazy `getXxx()` initialisation and with the pattern already logged in the xml/xpath units, and the
  facade validates a signature single-threaded per call — so it is a latent hardening gap, not a
  reachable defect in the normal path.
- **Recommendation:** Guard the lazy fields with a `sync.Once` (or a package-level note in `doc.go`) if
  concurrent access to a single `Signature` is ever supported; otherwise document that a `*Signature` is
  not safe for concurrent use, matching the Java contract.

### J17A-STD-002
- **Severity:** Info
- **Category:** STD
- **Location:** `dss/jades/jades_level_baseline_b.go:71` (`NewLevelBaselineB`)
- **Evidence:**
  ```go
  if certificateVerifier == nil {
      panic("certificateVerifier must not be null!")
  }
  // Upstream repeats the same null-check with the "signatureParameters must be defined!"
  // message against certificateVerifier (a known upstream copy-paste); reproduced verbatim.
  if certificateVerifier == nil {          // checks certificateVerifier AGAIN, never parameters
      panic("signatureParameters must be defined!")
  }
  ```
- **Impact:** The second guard is a faithful reproduction of an upstream copy-paste bug: it re-tests
  `certificateVerifier` (already proven non-nil) and never tests `parameters`. A `nil` `parameters` is
  therefore not caught here and will nil-panic later in the first `b.parameters.…` call. This is
  documented as deliberate parity (the comment says so), and the constructor is only reached from
  internal signing flows that always pass valid parameters, so it is not a reachable defect — recorded so
  it is not mistaken for a porting oversight and so the "reproduced verbatim" intent stays explicit.
- **Recommendation:** No change required for parity. If ever diverging, the second check should test
  `parameters == nil`; otherwise leave the verbatim reproduction and its comment in place.

## Open questions

1. **Concurrent `*Signature` access (J17A-STD-001).** Does any production path (e.g. a batch validator or
   the `esig` CLI) share one `*Signature` across goroutines? If yes, re-rate STD-001 to High (data race
   under concurrent use). If the facade always serializes per-signature (the apparent case), Low stands.
2. **Level acceptance (B-B vs B-LT vs B-LTA).** The lens asks about "wrong level acceptance (missing
   counter-signature/revocation data = Critical)". That decision logic lives in `jades_level_baseline_lt.go` /
   `jades_level_baseline_lta.go` / `jades_baseline_requirements_checker.go` — **U17b**, not these three
   files. `DataFoundUpToLevel()` here only *reads* the `Has*Profile()` flags. Recommend confirming in
   U17b that a B-LTA claim cannot be reported when the LTA counter-signature/revocation data is absent.
3. **`SignatureAlgorithm()` from signature length.** `jades_signature.go:94` derives the EdDSA algorithm
   from `len(signatureValue)` (64→ED25519, 114→ED448) for the metadata accessor. Verification itself uses
   the `alg` header (clean), so this only affects reported metadata, but a 64-byte non-EdDSA value would be
   mislabelled. Confirm U17b/`spi` does not feed this accessor back into a verify decision.

---

## jades timestamp/service/checker (unit U17b)

- **Unit:** U17b
- **Scope:** the JAdES `TIMESTAMP / SERVICE / CHECKER` layer — timestamp-token extraction, the signing/
  extension/counter-signing service, the signing-certificate source, the ETSI JSON-Schema engine, the
  baseline-requirements checker, the signature-parameter holder, and the B-LT extension level.
  Focus files: `jades/jades_timestamp_source.go` (772), `jades/jades_service.go` (674),
  `jades/jades_certificate_source.go` (572), `jades/specs/schema_engine.go` (441),
  `jades/jades_baseline_requirements_checker.go` (428), `jades/jades_signature_parameters.go` (422),
  `jades/jades_level_baseline_lt.go` (416) — 3725 lines total.
- **Date:** 2026-08-25
- **Depth:** every focus file **full-read** (3725 lines) + all `func` signatures grepped (231 total) +
  supporting cross-check reads well under the ~8k budget.

### Files read (full vs budgeted)

| File | Coverage | Functions covered |
|---|---|---|
| `jades/jades_timestamp_source.go` (772) | **FULL** (1–772) | all 61 `func` read in full: the `jadesSignedPropertiesAsEtsiUComponents` adapter, every `Is*` predicate, `IncorporateArchiveTimestampReferences`, the `archiveTimestampReferences`/`signatureTimestampReferences`/`signingCertificateTimestampReferences`/`encapsulatedReferencesFromTimestamps` reimplementations, the `jadesTS*` dedup helpers, `GetCertificateRefs`/`GetCRLRefs`/`GetOCSPRefs`, `GetEncapsulated{Certificate,CRL,OCSP}Identifiers` + `to*` converters, `GetTimestampMessageImprintDigestBuilderFor{Token,Algorithm}`, `GetCounterSignatures`, `Get{Signature,Archive}TimestampData`, `MakeTimestampToken(s)`, `extractTimestampTokens`/`toTimestampToken`/`extractArchiveTimestampTokens`, `GetArchiveTimestampType`, `MakeEvidenceRecords`. |
| `jades/jades_service.go` (674) | **FULL** (1–674) | all 32 `func` read in full: `NewService`, `GetContentTimestamp(ForDocuments)` + assert, `GetDataToSign(ForDocuments)` + `jadesServiceAssertMultiDocumentsAllowed`, `SignDocument(s)`, `JAdESBuilder` + `jadesServiceJWSJsonSerializationObjectToSign`/`jadesServiceContainsSignatures`, `ExtendDocument` + `jadesServiceAssertExtensionPossible`/`extensionProfile`, `TimestampDocuments`, `AddSignaturePolicyStore(WithEncoding)`, `GetDataToBeCounterSigned`/`CounterSignSignature` + `jadesServiceVerifyAndSetCounterSignatureParameters`, `IsValidSignatureValue`, `AssertSigningCertificateValid` + `jadesServiceAssertSigningCertificateValidForAlgorithm`/`jadesServiceAssertSignatureValueValid`, `MultipleDocumentsService` + adapter. |
| `jades/jades_certificate_source.go` (572) | **FULL** (1–572) | all 37 `func` read in full: `NewCertificateSource`, `KeyIdentifierCertificateRefs`/`KeyIdentifierCertificates` + `jadesCertificateSourceFindTokensFromRefs`, `extractX5T`/`extractX5TS256`/`extractX5TO`/`extractSigX5Ts`/`extractKid`/`extractX509Url`/`extractX5C`/`extractEtsiU`, `extract{CertificateValues,AttrAuthoritiesCertValues,TimestampValidationData,AnyValidationData,CompleteCertificateRefs,AttributeCertificateRefs}`, `extractCertificate{Values,Refs}FromList`/`extractX509Cert`, `ExtractCandidatesForSigningCertificate`, `resolveFromSource`/`resolveByKid`/`resolveByUri`, `extractPublicKey` + `jadesCertificateSourcePublicKey`, `checkSigningCertificateRef`/`isValid`/`signingCertificateDigest`, `OrphanCertificateRefs`, `ReferencesForCertificateToken`, `FindTokensFromCertRef`. |
| `jades/specs/schema_engine.go` (441) | **FULL** (1–441) | all 15 `func` read in full: `newSchemaSet`/`decodeJSON`/`load`/`resolveRef`/`jsonPointerGet`, `validate`/`validateObject`/`validateArray`/`validateCombinators`, `schemaSummary`/`jsonType`/`typeMatches`/`joinPath`/`errPath`/`validateJSONAgainstSchema`. |
| `jades/jades_baseline_requirements_checker.go` (428) | **FULL** (1–428) | all 13 `func` read in full: `NewBaselineRequirementsChecker`, `HasAdESProfile`, `HasBaselineBProfile`, `critRequirements` + `jadesBaselineRequirementsChecker{Equal,Contains}`, `signingTimeRequirement`, `HasBaselineTProfile`, `HasBaselineLTProfile`, `ContainsLTLevelCertificates` + `containsCertificateValues`/`containsTstOrAnyValDataCertificates`, `HasBaselineLTAProfile`. |
| `jades/jades_signature_parameters.go` (422) | **FULL** (1–422) | all 37 `func` read in full: `NewSignatureParameters`, `SetSignatureLevel`, the three `Get*TimestampParameters`, every `Is*/Set*` access pair, `SigningCertificateDigestMethod`, `JwsSerializationType`, `SigDMechanism`, `JadesSigningTimeType`, `ExpirationTime`, `IsBase64UrlEncoded{Payload,EtsiUComponents}`, `String` + `boolPtrString`, `Equals` + `boolPtrEquals`. |
| `jades/jades_level_baseline_lt.go` (416) | **FULL** (1–416) | all 20 `func` read in full: `NewLevelBaselineLT`, `ExtendSignatures`, `removeOld{Certificate,Revocation}Values`, `RemoveLastTimestampAndAnyValidationData`, `incorporateValidationDataFor{Signature,Timestamps}`, `XVals`/`jadesLevelBaselineLTX509CertObject`, `Incorporate{XVals,RVals,TstValidationData,AnyValidationData,ValidationData}`, `RVals`/`jadesLevelBaselineLT{CrlVals,OcspVals}`, `tstVd`, `jadesLevelBaselineLTExtendToLTLevelSignatures`/`jadesLevelBaselineLTLtLevelExtensionRequired`. |
| Supporting cross-checks (for the verify/level boundary) | **REGION** | `spi/certificate_validity.go` 140–240 (`IsValid` disjunction + bool defaults at 1–95), `spi/validation/baseline_requirements_checker.go` 230–400 (`Minimal{T,LT,LTA}Requirement`, `SignatureTimestampsCreatedBeforeSignCertExpiration`, `ValidationContext`), `jades/jades_signature.go` 1–80 (struct) + 405–470 (`CounterSignatures`) + 1300–1354 (`DataFoundUpToLevel`/`ValidateStructure`), `jades/jws_json_serialization_object.go` 60–110 (`IsValid`), `jades/jws_json_serialization_parser.go` 1–175 (`validateJWSStructure`→`ValidateAgainstSchema`), `jades/specs/abstract_jades_utils.go` 1–100 + `jades_utils.go` 1–70 (defs + singleton + upstream `getValidator` caching note), `jades/specs/schema/` corpus (11 files, keyword audit), `dss/document/abstract_signature_parameters.go` 44–110 (no weak-alg default). |
| Tests (coverage only, not read fully) | — | `specs/jades_utils_test.go` (schema KATs: protected/unprotected/JAdES envelope accept+reject), `jades_document_validator_smoke_test.go`, `jades_sign_kat_test.go`, `jades_{upstream,downstream}_cross_validation_test.go` — confirm the level/checker/sign path is oracle-pinned. |

### Security posture (top priority) — VERIFIED CLEAN on the attack surface

Every SEC lens the unit asked for was checked end-to-end; **no Critical, no High, and no exploitable
false-accept was found**:

- **B-LT / B-LTA level (counter-sig present+valid, correct digest over the JWS).** `MinimalLTRequirement`
  (frozen base, `baseline_requirements_checker.go:274`) is the actual LT gate: it requires non-all-self-signed
  certs **and** non-empty CRL+OCSP **and** `isAllRevocationDataPresent()` (or LT/XL validation data for
  self-signed chains). `HasBaselineLTProfile` adds the ref/`sigRTst`/`rfsTst` absence checks. **A counter-
  signature is NOT a B-LT requirement** (it is the B-**LTA** gate — `MinimalLTARequirement` requires
  `ArchiveTimestamps()` non-empty), so there is no missing-counter-sig check to flag. The counter-sig digest
  over the JWS value is enforced by the already-reviewed `internal/jose` (`VerifySignature`/`alg=none` reject/
  key-type mismatch→`IsIntact=false`), and `GetCounterSignatures` (`timestamp_source.go:617`) only *reports*
  the token — it is not itself a verify gate. No false-accept.
- **Certificate source (wrong cert = false-accept, chain building, revocation attachment).** `checkSigning-
  CertificateRef` (`certificate_source.go:429`) selects the signing candidate via `CertificateValidity.IsValid()`
  = `IsDigestEqual() || (IsDistinguishedNameEqual() && IsSerialNumberEqual()) || IsResponderIdMatch()`. All
  those booleans **default to `false`** (`certificate_validity.go:13–35`, no constructor sets them), so an
  unmatched candidate is **not** selected — the path is **fail-closed**. Selection is a *preference*; the real
  gate is the cryptographic `IsValidSignatureValue` → base `Verify` (already reviewed). Revocation/cert values
  are attached by `extractEtsiU`/`extractX5C`/`extractX509Url` with a `DER`-only encoding guard and `otherCert`
  skip — faithful to Java. No false-accept.
- **Schema engine (ETSI U / schema bypassable → level mis-claim).** `schema_engine.go` implements exactly the
  draft-07 subset the corpus uses. I audited all 11 schema files: the corpus uses **none** of `minLength`/
  `maxLength`/`minimum`/`maximum`/`uniqueItems`/`patternProperties`/`dependencies`/`if-then-else`/`exclusive*`/
  `const`/`multipleOf`, and `type` is always a bare string (never a multi-type array). The only `$ref` siblings
  present are `$schema`/`definitions` (annotation-only, harmless under `$ref`-replaces-siblings). `format`
  (1 file) and `contentEncoding` (5 files) are draft-07 *annotation* keywords — non-enforcement is the correct
  draft-07 behaviour and is documented. `additionalProperties:false` (32×) is enforced. So a hostile/invalid JWS
  **cannot** pass a schema it should fail. No bypass.
- **Signature parameters (unsafe weak-alg default).** `NewSignatureParameters` sets `signingCertificateDigest-
  Method: SHA512`, `jadesSigningTimeType: IAT`; the **signature** algorithm is never defaulted here — it is
  derived from the signing certificate's key (`document/abstract_signature_parameters.go` `SetSigningCertificate`
  → `documentEncryptionAlgorithmForKey`). No weak/HMAC default is exposed; the HS* dead code in `internal/jose`
  (U01 SEC-002) is not reachable from these parameters. The service additionally enforces the RFC 7518 ECDSA
  curve/digest pairing (`jadesServiceAssertSigningCertificateValidForAlgorithm`) and signature-bitlength
  (`jadesServiceAssertSignatureValueValid`). Clean.
- **Baseline requirements checker (B-B vs B-LT vs B-LTA enforcement).** `critRequirements` rejects non-array/
  empty/`crit`-with-RFC-reserved/unknown/unused critical headers (fail-closed); `signingTimeRequirement`
  enforces the 2025-07-15 `iat`-vs-`sigT` transition; `HasBaselineBProfile` enforces `alg`+exactly-one cert
  header+`sigPSt` cardinality+requirement (b); `HasBaselineTProfile` enforces single `sigTst`+pre-expiry
  (reqs c,d); `HasBaselineLTProfile`/`HasBaselineLTAProfile` per above. All faithful to `JAdESBaseline-
  RequirementsChecker.java`. No per-level mis-enforcement.
- **Service signing flow (byte-exact JWS before signing) + panics on hostile input.** `GetDataToSign` builds
  the `ToBeSigned` once via the builder; `IsValidSignatureValue` re-checks the value. Panics are Java-`throw`
  propagations at the service boundary (facade `recovered()` turns them into errors) — the documented
  convention, not a hostile-input DoS. `GetCounterSignatures`' `.(*Signature)` type-assert is safe:
  `Signature.CounterSignatures()` (`jades_signature.go:412`) only ever appends `*Signature` from
  `DSSJsonUtilsExtractJAdESCounterSignature`. No panic-on-hostile-input in a prod verify path.

### Findings

**Critical:** *no findings.*  **High:** *no findings.*

**PERF:**
- **J17B-PERF-001** (Medium) — schema engine re-parses all schema docs per call (see below).
- **J17B-PERF-002** (Low) — `jadesTSAddReferences` O(n²) dedup.

**SEC:** *no findings* (posture above).

**STD:**
- **J17B-STD-001** (Medium) — `jsonType`/`typeMatches` deviate from draft-07 integer semantics.
- **J17B-STD-002** (Low) — unsynchronized lazy caches.
- **J17B-STD-003** (Info) — `extractX5T` dead no-op body (parity with upstream; slf4j dropped).

### J17B-PERF-001
- **Severity:** Medium
- **Category:** PERF-LEAK (redundant re-parse / missing cache)
- **Location:** `dss/jades/specs/schema_engine.go:428` (`validateJSONAgainstSchema`), `:47`
  (`newSchemaSet`), `:146`/`:210` (`validate`/`validateObject`)
- **Evidence:**
  ```go
  // validateJSONAgainstSchema: builds a FRESH set (empty docs cache) on EVERY call
  set := newSchemaSet(files, defs)          // newSchemaSet: docs: map[string]any{}
  root, err := set.load(rootURI)            // load: fs.ReadFile + decodeJSON per uncached URI
  ```
- **Impact:** `validateAgainstSchema` is called once per schema (JWS envelope, protected header, unprotected
  header, and per base64 `etsiU` component) and once per signature in a multi-signature JSON-serialization
  document. Because each call constructs a new `schemaSet` with an empty `docs` cache, the ~11 embedded schema
  documents are re-read and re-decoded, and every `pattern` is re-compiled via `regexp.Compile`
  (`schema_engine.go:167`), on **every** call — and re-validated across the whole batch. Upstream's
  `JSONSchemaAbstractUtils#getValidator` (whose "…getValidator caching machinery" this file's header cites as
  folded in, `abstract_jades_utils.go:8`) caches the compiled validator across calls; that cache is **not**
  reproduced here. Bounded (schemas are small and embedded), so Low per single signature, but a real redundant
  cost on multi-signature/batch validation and a genuine deviation from the upstream caching contract.
- **Recommendation:** Hoist a shared, read-only `schemaSet` (or at least the parsed `docs` map + compiled
  `*regexp.Regexp`s) to the singleton utils (`JAdESUtils`/`JAdESProtectedHeaderUtils`/`JAdESUnprotectedHeader-
  Utils` already use `sync.Once`) and pass it into `validate`, so schemas are parsed once per process. Keep the
  per-call instance fresh if state must remain caller-local, but share the immutable parsed artifacts.

### J17B-STD-001
- **Severity:** Medium
- **Category:** STD (draft-07 semantics deviation → false-REJECT, not false-accept)
- **Location:** `dss/jades/specs/schema_engine.go:372` (`jsonType`), `:396` (`typeMatches`)
- **Evidence:**
  ```go
  case json.Number:
      if strings.ContainsAny(string(v), ".eE") { return "number" }  // "1.0" -> "number"
      return "integer"
  // typeMatches:
  case "integer": return actual == "integer"                        // "1.0" fails an integer schema
  ```
- **Impact:** Draft-07 defines `integer` as "any number with a zero fractional part," so `1.0` is a valid
  `integer`. The corpus declares `type:integer` for the JWT `iat` claim (`schema/rfcs/rfc7519.json:7`) and for
  `noticeNumbers` items (`19182-jsonSchema.json:197`). A conforming producer that serializes `iat` as `1.0`
  (legal JSON, legal per draft-07) is therefore **rejected** by this engine even though upstream
  (jsonsKema/everit) accepts it. This is a **false-reject** (over-strict), so it cannot cause a level
  mis-claim or false-accept — but it is a concrete conformance deviation on the `iat` path that the
  baseline-requirements checker's `signingTimeRequirement` depends on.
- **Recommendation:** In `jsonType`/`typeMatches`, treat a `json.Number` as `integer` when its value has no
  fractional part (e.g. parse to a big/float and check `== float64(int64(...))`, or reject only when a
  non-zero fractional component is present), matching draft-07. Add a KAT with `iat: 1.0` to pin it.

### J17B-PERF-002
- **Severity:** Low
- **Category:** PERF-BIGO (bounded)
- **Location:** `dss/jades/jades_timestamp_source.go:351` (`jadesTSAddReferences`)
- **Evidence:**
  ```go
  for _, candidate := range referencesToAdd {
      found := false
      for _, existing := range *referenceList { if existing.Equals(candidate) { found = true; break } }
      if !found { *referenceList = append(*referenceList, candidate) }
  }
  ```
- **Impact:** O(n·m) linear-scan dedup. Called from the `signatureTimestampReferences`/`archiveTimestamp-
  References`/`IncorporateArchiveTimestampReferences` paths on each signature. `n`/`m` are the timestamped-
  reference lists (certificate refs + signing-cert refs + key-info refs), typically single digits, so the
  quadratic cost is negligible in practice. Faithful to Java's `addReferences`/`addReference` (which likewise
  linear-scan `List#contains`). Reported for completeness; not a real-world cost.
- **Recommendation:** None required. If the lists ever grow (many certs/refs), replace the scan with a
  `map[string]struct{}` keyed on `TimestampedReference` identity.

### J17B-STD-002
- **Severity:** Low
- **Category:** STD (concurrency — lazy caches)
- **Location:** `dss/jades/jades_timestamp_source.go:62` (`components`), `:90`
  (`timestampAttributeMap`); `jades/spi/validation/baseline_requirements_checker.go` (`validationContext`)
- **Evidence:**
  ```go
  type jadesSignedPropertiesAsEtsiUComponents struct { ... components []*EtsiUComponent } // lazily filled in Attributes()
  func (s *TimestampSource) ... { s.timestampAttributeMap[timestampToken] = signatureAttribute }  // written lazily
  ```
- **Impact:** Several lazy caches (`components`, `timestampAttributeMap`, the base `validationContext`) are
  populated on first use with **no** mutex. Under concurrent access to the same `*Signature`/source this is a
  data race. This is the same class of issue as U17a **J17A-STD-001** (unsynchronized lazy fields on
  `*Signature`); the facade appears to serialize per-signature validation, in which case it is not exploitable.
  Listed here because the timestamp source adds its own two lazy fields on top of the one already noted.
- **Recommendation:** Same as J17A-STD-001 — confirm the facade never shares one `*Signature` across
  goroutines; if it can, add a `sync.Once`/mutex guard around the lazy population.

### J17B-STD-003
- **Severity:** Info
- **Category:** STD
- **Location:** `dss/jades/jades_certificate_source.go:118` (`extractX5T`)
- **Evidence:**
  ```go
  func (s *CertificateSource) extractX5T() {
      base64UrlSHA1Certificate := s.jws.ProtectedHeaderValueAsString(jose.HeaderX509CertificateThumbprint)
      if utils.IsStringNotEmpty(base64UrlSHA1Certificate) { //nolint:staticcheck // ... Java's body is `LOG.warn(...)` only.
          // Upstream builds a Digest(SHA1, ...) purely to log "...not supported by the JAdES standard";
          // with slf4j dropped per PORTING.md there is nothing left for this branch to do.
      }
  }
  ```
- **Impact:** The method is a deliberate no-op: Java's `extractX5T` only *logs* that `x5t` (SHA-1 thumbprint)
  is "not supported by the JAdES standard" and adds nothing, so with slf4j dropped the branch body is empty.
  Correct and faithful; recorded only because an empty `if` body is an unusual idiom and could look like a
  missing implementation. `x5t` (SHA-1) is genuinely not a JAdES signing-cert reference — the supported ones
  are `x5t#S256`/`x5t#o`/`sigX5ts` (`extractX5TS256`/`extractX5TO`/`extractSigX5Ts`).
- **Recommendation:** None (parity with upstream; the `//nolint:staticcheck` + comment documents intent).

### Tool log

```text
$ (dss/) gofmt -l jades jades/specs              # -> (no output)
$ (dss/) go vet ./jades/...                       # -> (no output)
$ (dss/) golangci-lint run --config=../.github/.golangci.yml ./jades/...
           # golangci-lint present (/opt/homebrew/bin/golangci-lint): "0 issues."
```
All three checks clean. No unavailable tools.

### Open questions

1. **`iat` as `1.0` (J17B-STD-001).** Confirm whether any oracle fixture serializes `iat` with a fractional
   form (`1.0`); if so, J17B-STD-001 is a live false-reject against upstream and should be rated High.
   If no fixture does, it stays a latent conformance deviation (Medium).
2. **Schema cache (J17B-PERF-001).** Confirm the batch validator / `esig` CLI does not validate many
   signatures per invocation such that the per-call schema re-parse (J17B-PERF-001) becomes a measurable cost;
   if it does, re-rate Medium→High.
3. **Concurrent `*Signature` access (J17B-STD-002 / J17A-STD-001).** Same as U17a OQ-1 — if a production path
   shares one `*Signature`/`TimestampSource` across goroutines, re-rate the lazy-cache findings to High
   (data race).

## jades tail + specs + extension (unit U17c)

- **Unit:** U17c
- **Scope:** the **remaining** non-test Go files of `dss/jades`, `dss/jades/specs`, and `dss/jades/extension` —
  i.e. every file NOT already covered by U17a (`jades_signature.go`, `jades_level_baseline_b.go`,
  `dss_json_utils.go`, `jades_timestamp_source.go`, `jades_service.go`, `jades_certificate_source.go`,
  `specs/schema_engine.go`, `jades_baseline_requirements_checker.go`, `jades_signature_parameters.go`,
  `jades_level_baseline_lt.go`). That leaves **60 files** (all ≤377 lines) plus the `crossgen` testdata
  generator (232).
- **Date:** 2026-08-25
- **Depth:** all 60 in-scope files **full-read** (1:1) + `crossgen/main.go` full-read + supporting cross-check
  reads (`internal/jose/jws.go` verify-cache invalidation, `utils/codec.go` base64 padding,
  `utils/strings.go` `Split`/`Trim`, `jades_signature.go` verify-path ordering 600–770/1050–1240) + a
  `/tmp` repro confirming the `getDigest` split behaviour. Well under the ~8k-line budget.
- **Context:** Go port of Java DSS 6.5.RC1. `dss/PORTING.md` read first (1:1 port, staticcheck style rules
  disabled, SA1019 allowed, byte-ident enum values, stdlib-first, no cgo). `docs/compatibility/known-gaps.md`
  read (not re-reporting listed gaps). `internal/jose` (incl. SEC-001 stale `validSignature` cache) and the
  jades core/timestamp layer (J17A/B-*) already reviewed — **not re-reported**; this unit focuses on the tail.

### Files read (all FULL unless noted)

| File | Lines | Coverage |
|---|---|---|
| `jades/jades_timestamp_message_digest_builder.go` | 377 | FULL — all 5 `*MessageDigest()` impls + `writeSignedDataBinaries`/`writeSigDReferencedOctets`/`etsiUComponentValue`/`jadesTMDBCanonicalizedValue` (documented upstream canonicalization TODO). |
| `jades/jades_header_parameter_names.go` | 322 | FULL — all `JAdESHeaderParameterNames*` constants + the `TstVD`/`TstVd` and `SigPSt`/`SigPst` alias pairs (see J17C-STD-003). |
| `jades/jades_level_baseline_lta.go` | 281 | FULL — `ExtendSignatures`, `ltaIncorporate{ValidationDataForTimestamps,AnyValidationData}`, `incorporateArcTst`/`archiveTimestamp`, the three `assert*/check*` guards. |
| `jades/jades_counter_signature_builder.go` | 256 | FULL — `GetSignatureValueToBeSigned`, `BuildEmbeddedCounterSignature`, `updateMasterSignatureRecursively`, `extractSignatureById`/`signatureOrItsCounterSignature` (recursive cSig walk). |
| `jades/jades_level_baseline_t.go` | 240 | FULL — `ExtendSignaturesDocument`/`ExtendSignatures` (sigTst over `GetSignatureTimestampData`), `extendToTLevelSignatures`, `jadesLevelBaselineTTimestampSource`. |
| `jades/jades_signature_policy_store_builder.go` | 239 | FULL — `AddSignaturePolicyStore(ForSignature)`, `AddSignaturePolicyStoreIfDigestMatch`, `CheckDigest` (digest-match gate, fail-closed). |
| `jades/http_headers_payload_builder.go` | 237 | FULL — `Build` (draft-cavage canonicalization), `httpHeadersPayloadBuilderHTTPHeaderWithName`/`AsHTTPHeader`/`ToHTTPHeaders`, `assertHttpHeadersConfigurationIsValid`/`checkIfDigestHTTPHeaderDocument`. |
| `jades/jws_converter.go` | 235 | FULL — all 4 `JWSConverterFrom*` + `jwsConverter{AssertConvertPossible,ToClearJSONIncorporation,ToBase64UrlIncorporation,AssertComponentSupportsConversion}` (timestamp-header exclusion list). |
| `jades/jws.go` | 214 | FULL — `NewJWS` (+`CheckCritOverride`), `SetPayloadOctets`, `SignedPayload`, `DecodedPayload` (**cached**, see J17C-STD-002), all `ProtectedHeaderValueAs*`, `CheckCrit`. |
| `jades/jades_signature_scope_finder.go` | 196 | FULL — `FindSignatureScope`, `getSignatureScopeFromOriginalDocument(s)`, `getHttpHeader(Signature/MessageBody)SignatureScope`, `getDigest` (**see J17C-STD-001**). |
| `jades/jades_etsi_u_header.go` | 183 | FULL — `Attributes` (lazy, nil-safe), `Add/Remove/ReplaceComponent`/`RemoveLastComponent` (all error-returning, panic-free), `getEtsiUToEdit`/`setEtsiU`/`clearCachedAttributes`. |
| `jades/jades_serialization_builder.go` | 180 | FULL — `New(FromSignature)`, `Build` (parallel-signing append), `assertB64ConfigurationConsistent`, `jadesSerializationBuilderExtractDocumentToBeSigned`. |
| `jades/jws_json_serialization_parser.go` | 159 | FULL — `Parse`, `IsSupported`, `extractSignatures`/`extractAndAddJWSSignature` (b64-gated payload decode), `validateJWSStructure`. |
| `jades/jades_ocsp_source.go` | 155 | FULL — `NewOCSPSource` + `extractEtsiU` + every `extract*`/`add` (fail-soft, logged-and-discarded on parse error). |
| `jades/jades_crl_source.go` | 152 | FULL — `NewCRLSource` + `extractEtsiU` + every `extract*`/`add`. |
| `jades/jws_compact_serialization_parser.go` | 132 | FULL — `Parse` (bufio, 16 MB cap), `IsSupported` (byte-walk + `ending` line-break trim). |
| `jades/abstract_jws_document_analyzer.go` | 131 | FULL — `Init*`, `OriginalDocumentsForSignature`, `GetDefaultSignaturePolicyValidator`, the `jwsDocumentAnalyzerBase` recovery shim. |
| `jades/jades_revocation_ref_extraction_utils.go` | 127 | FULL — `CreateOCSPRef`/`CreateCRLRef` + `ResponderId` (documented tolerant-cast deviation). |
| `jades/jws_json_serialization_generator.go` | 126 | FULL — `Generate` + `build(Flat)JWSJsonSerialization`/`signatureJSONMap` (insertion-ordered). |
| `jades/abstract_jades_builder.go` | 120 | FULL — `InitAbstractJAdESBuilder`, `BuildDataToBeSigned`, `IncorporateHeader`/`IncorporatePayload`. |
| `jades/jades_extension_builder.go` | 112 | FULL — `AssertEtsiUComponentsConsistent(WithEncoding)` (tri-state `*bool`), `AssertJWSJsonSerializationObjectValid`/`AssertJSONSerializationObjectMayBeExtended`. |
| `jades/json_object.go` | 109 | FULL — `New(FromMap)JsonObject`, `ContainsValue` (documented `==`-panic-on-uncomparable caveat), all Map-facade methods. |
| `jades/jades_compact_builder.go` | 103 | FULL — `Build` (detached-asymmetric payload), `AssertConfigurationValidity` (B-only). |
| `jades/jwt_payload.go` | 102 | FULL — `JWTPayload` claim accessors (all delegate to `DSSJsonUtilsGetAs*`, array-`aud`→"" reproduced). |
| `jades/jws_json_serialization_object.go` | 97 | FULL — container + `AddSignature` + `IsValid`. |
| `jades/jws_serialization_analyzer_validator.go` | 95 | FULL — `BuildSignatures`/`BuildJwsJsonSerializationObject` (panic-on-unsupported, facade-recovered). |
| `jades/extension/jades_document_extender.go` | 92 | FULL — `CreateSignatureService`, `IsSupported*`, `EmptySignatureParameters`, `SignatureForm`. |
| `jades/jws_compact_document_analyzer.go` | 85 | FULL — `IsSupported`/`BuildSignatures`/`BuildJwsJsonSerializationObject`. |
| `jades/specs/abstract_jades_utils.go` | 83 | FULL — `rfcDefinitions`/`schemaDefinitions`/`validateAgainstSchema` (folded 3-level Java hierarchy). |
| `jades/http_header.go` | 83 | FULL — `HTTPHeader` DSSDocument pseudo-doc (unsupported ops panic, Java-parity). |
| `jades/jades_signature_identifier_builder.go` | 81 | FULL — `CounterSignaturePosition`/`SignaturePosition` (documented `hashCode`→full-equality deviation). |
| `jades/etsi_u_component.go` | 81 | FULL — `EtsiUComponentBuild(FromValue)`/`Component`/`IsBase64UrlEncoded`. |
| `jades/http_header_digest.go` | 76 | FULL — `NewHTTPHeaderDigest` + `httpHeaderDigestBuildInstanceDigestValue` (RFC 3230, `utils.ToBase64` **padded** — see J17C-STD-001). |
| `jades/jades_attribute.go` | 73 | FULL — `Attribute` + `Identifier` (lazy) + `Equals`. |
| `jades/jades_attribute_identifier.go` | 70 | FULL — `AttributeIdentifierBuild(WithOrder)` + `jadesAttributeIdentifierWriteChars` (UTF-16 2-byte big-endian). |
| `jades/abstract_jws_document_validator.go` | 70 | FULL — `newAbstractJWSDocumentValidator`, `InitializeDiagnosticDataBuilder`, `DocumentAnalyzer` (narrowing assert). |
| `jades/jades_timestamp_parameters.go` | 65 | FULL — `TimestampParameters` + `SetCanonicalizationMethod` (panic, upstream unimplemented). |
| `jades/jades_counter_signature_parameters.go` | 61 | FULL — `CounterSignatureParameters` + accessors + `Equals`. |
| `jades/jades_signed_properties.go` | 60 | FULL — `SignedProperties` + `Attributes`/`mapKeyValues` (TODO: re-parse). |
| `jades/jws_serialization_document_validator.go` | 59 | FULL — leaf validator + `DocumentAnalyzer` narrowing. |
| `jades/specs/jades_utils.go` | 58 | FULL — `JAdESUtils` singleton + `SchemaURI`/`SchemaDefinitions`/`ValidateAgainstSchema`. |
| `jades/jws_document_analyzer_factory.go` | 53 | FULL — `IsSupported`/`Create` + `init()` registration. |
| `jades/jades_document_validator_factory.go` | 53 | FULL — `IsSupported`/`Create` + `init()` registration. |
| `jades/specs/jades_unprotected_header_utils.go` | 50 | FULL — `JAdESUnprotectedHeaderUtils` singleton + schema defs. |
| `jades/specs/jades_protected_header_utils.go` | 50 | FULL — `JAdESProtectedHeaderUtils` singleton + schema defs. |
| `jades/jwt_claim_names.go` | 50 | FULL — RFC 7519 registered claim-name constants. |
| `jades/jades_signature_integrity_validator.go` | 50 | FULL — `Verify` (delegates to `jose.JWS.VerifySignature`, `(bool,error)` split). |
| `jades/extension/jades_document_extender_factory.go` | 40 | FULL — factory + `init()` registration. |
| `jades/http_header_signature_scope.go` | 39 | FULL — `HTTPHeaderSignatureScope` scope type. |
| `jades/jws_compact_document_validator.go` | 37 | FULL — leaf validator + `DocumentAnalyzer` narrowing. |
| `jades/jades_diagnostic_data_builder.go` | 37 | FULL — `BuildDetachedXmlSignature` (adds JWS serialization type + expiration). |
| `jades/jades_signing_time_type.go` | 34 | FULL — `SigningTimeType` enum (IAT/SIG_T/NONE) + `Values()`. |
| `jades/jades_level_baseline_extension.go` | 32 | FULL — `LevelBaselineExtension` interface. |
| `jades/jades_certificate_ref_extraction_utils.go` | 31 | FULL — `CreateCertificateRef` (kid/issuerSerial). |
| `jades/jades_builder.go` | 28 | FULL — `Builder` interface. |
| `jades/jws_constants.go` | 26 | FULL — RFC 7515 member-name constants. |
| `jades/http_header_message_body_signature_scope.go` | 29 | FULL — `HTTPHeaderMessageBodySignatureScope` (overrides `Description`). |
| `jades/specs/schema_fs.go` | 12 | FULL — `//go:embed schema`. |
| `jades/specs/doc.go` | 11 | FULL — package doc. |
| `jades/testdata/crossgen/main.go` | 232 | **FULL** (test-only generator — see J17C-INFO-001). |

Supporting cross-checks (REGION, for the two findings + boundary): `internal/jose/jws.go` 140–200
(`SetSignature`→`validSignature=nil`, `SetCompactSerializationParts`, `SigningInputBytes`), `utils/codec.go`
70–130 (`ToBase64` = **StdEncoding, padded**; `FromBase64` lenient), `utils/strings.go` 100–160 (`Trim`/
`LowerCase`), `jades/jades_signature.go` 600–770 (verify-path order: `jadesSignaturePayload` sets payload
**before** `keyBindingSignatureReferenceValidation` reads it) + 1050–1240 (`sdHash`/`sdAlg` via
`DecodedPayload`), `jades/jades_root_chunk_test.go` 22–43 (alias-pinning test).

### Security posture (top priority) — VERIFIED CLEAN on the attack surface

Every SEC lens the unit asked for was checked end-to-end; **no Critical and no exploitable false-accept was
found** in the in-scope files:

- **`jws.go` stale-cache (the SEC-001 class, explicitly requested).** The *verify* cache
  (`internal/jose.JWS.validSignature`) IS correctly invalidated — `SetSignature` sets `validSignature = nil`
  (`internal/jose/jws.go:150`), so the already-reviewed SEC-001 pattern does not recur on the hot verify path.
  However the **wrapper** `jades.JWS` carries a *sibling* cache, `decodedPayload`, that is **never
  invalidated** by any payload setter — see **J17C-STD-002**. It is not currently reachable as a bug (payload
  is set at parse time before first read), but it is the same class and the one real cache-hygiene gap in the
  unit.
- **JWS↔JSON round-trip (`jws_converter.go` / `jws_json_serialization_parser.go`).** Lossless for the
  supported cases: the converter refuses (error, not silent) any mixed-form `etsiU` array
  (`jwsConverterAssertConvertPossible`) and any timestamp-bearing component whose re-encoding would invalidate
  a covering timestamp (`jwsConverterAssertComponentSupportsConversion` against the `arcTst`/`rfsTst`/`sigRTst`
  list). The parser installs `protected` **before** the payload so `b64` gates the decode order
  (load-bearing, documented). No false-reject / corruption path found.
- **HTTP header construction (`http_headers_payload_builder.go`).** This builds a **signed canonicalization
  string** (draft-cavage-http-signatures-10), not an HTTP response — so the CRLF-injection lens does not
  apply: a value with `\r\n` becomes signed bytes, not an injected header line. Values are caller-controlled
  (signing path) or detached-content-controlled, and repeated-header joining is encounter-ordered with `, `.
  `assertHttpHeadersConfigurationIsValid` enforces exactly-one `Digest` header. Clean.
- **ETSI U header parsing (`jades_etsi_u_header.go`).** Malformed component → `EtsiUComponentBuild` returns
  nil (skipped), not a panic; every mutator returns an error (nil today); `getEtsiUToEdit` is nil-safe.
  `ReplaceComponent`'s `Equals`→`Identifier`→`jose.JSON(value)` is JSON-safe because `value` comes from a
  parsed component. Panic-free.
- **Counter-signature (B-LT) — `jades_counter_signature_builder.go`.** The digest-over-exact-JWS-bytes work is
  the TSA's / `jades_timestamp_source.go`'s (U17b, already reviewed); this builder only *embeds* the cSig and
  re-serializes the nested chain innermost-out (`updateMasterSignatureRecursively`), preserving nesting and
  the ordered `cSig` members. `extractSignatureById`/`signatureOrItsCounterSignature` is a correct recursive
  cSig walk with a timestamp-before-extend guard and a compact-extension refusal. No wrong-match.
- **Timestamp message-imprint (`jades_timestamp_message_digest_builder.go`).** Digests delegate to
  `spi.DSSUtilsDigest`/`DSSMessageDigestCalculator`; `archiveTimestampMessageDigest` assembles
  payload/headers ‖ '.' ‖ encodedHeader ‖ '.' ‖ sigValue ‖ '.' ‖ etsiU (per ETSI TS 119 182-1 5.3.6.3.1).
  The only approximation is `jadesTMDBCanonicalizedValue`, a **documented upstream TODO** (canonicalization
  unsupported), reproduced verbatim — not a porting shortcut. Clean.
- **Scope finder — `jades_signature_scope_finder.go`.** Signature *matching* is by ID (`extractSignatureById`)
  and by original-document type; the fail-soft `getDigest` (see **J17C-STD-001**) drops a scope rather than
  mis-matching one. The crypto verdict is driven by `ReferenceValidations`/`CheckSignatureIntegrity` (U17a,
  clean), not by the scope, so this is a reporting gap, not a false-accept.
- **Policy store — `jades_signature_policy_store_builder.go`.** `CheckDigest` is fail-closed: no
  `SignaturePolicy` → not added; no expected digest → not added; content absent → digest "not checked" (upstream
  behaviour); otherwise the computed digest must `Equals` the expected one or the store is not added.
  `spDocSpecification.Id()` is routed through `DSSJsonUtilsOidObjectFromURI` (URI→OID, already-reviewed
  helper). Clean.
- **Serialization — `jades_serialization_builder.go` / `jws_json_serialization_generator.go`.** Insertion-
  ordered `jose.Object` throughout; `assertB64ConfigurationConsistent` rejects a mixed `b64` population; the
  generator enforces exactly-one signature for flattened form. Byte-exact by construction (ordered map +
  single `ToJSONString`). Clean.
- **Level evidence (`jades_level_baseline_t.go` / `jades_level_baseline_lta.go`).** Correct evidence
  requirements: T-level adds a `sigTst` over `GetSignatureTimestampData`; LTA-level runs the LT chain first,
  then (only when a signature already `HasLTAProfile`) removes the old timestamp + incorporates
  revocation/timestamp/any validation data per the `ValidationDataEncapsulationStrategy`, then adds `arcTst`
  over `GetArchiveTimestampData`. The LTA guards (`checkArchiveTimestampParameters`,
  `assertDetachedDocumentsContainBinaries`, `checkEtsiUContentUnicity`) are all fail-closed errors. The
  *level-accept* verdict itself is U17b's checker (already reviewed). Clean.

### Findings

**SEC:** *no findings* — the verification / false-accept / alg-confusion / header-injection / scope-matching /
policy-OID attack surface was verified clean (see posture above).

### J17C-STD-001
- **Severity:** Medium (upgrade to High if any downstream consumer treats an *absent* message-body scope as a
  coverage/verdict signal — see Recommendation)
- **Category:** STD
- **Location:** `dss/jades/jades_signature_scope_finder.go:181` (`getDigest`); interacts with
  `dss/jades/http_header_digest.go:66` (`httpHeaderDigestBuildInstanceDigestValue`) and `dss/utils/codec.go:87`
  (`ToBase64` = `StdEncoding`, **padded**)
- **Evidence:**
  ```go
  valueParts := strings.Split(digestHeaderValue, "=")   // :181
  if len(valueParts) != 2 { return nil }                // :182-190
  // HTTPHeaderDigest value = "SHA-256=" + utils.ToBase64(digest)  (standard, padded base64)
  // e.g. "SHA-256=<43 chars>="  ->  strings.Split => ["SHA-256","<43>",""]  (len 3)
  ```
- **Impact:** A real RFC 3230 instance-digest value is `algo=` + **padded** standard base64. For any digest whose
  base64 carries padding (MD5, SHA-1, SHA-256, SHA-512 — all except SHA-384), `strings.Split(..., "=")` yields
  ≥3 parts, so `len != 2` and `getDigest` returns `nil`. `getHttpHeaderDigestSignatureScope` then drops the
  "message body value digest" scope, so the reported signature scope for an HTTPHeaders-mechanism JAdES
  signature silently **omits the message body it signs**. This is internally inconsistent within the port:
  `http_header_digest.go` *produces* the padded value and `jades_signature_scope_finder.go` *cannot parse* it.
  It is a porting divergence from Java (`String.split("=")` drops trailing empty strings → 2 parts), not a
  false-accept (the integrity verdict is driven separately and is clean).
- **Recommendation:** Split on the first `=` only: `valueParts := strings.SplitN(digestHeaderValue, "=", 2)`
  (and keep `len(valueParts) != 2` as the conformance check). This matches the RFC 3230 `algo=value` grammar and
  the value `HTTPHeaderDigest` itself emits. Add a KAT that parses a padded (`SHA-256=...=`) `Digest` value.
  If any consumer keys a coverage/verdict decision off the presence of the message-body scope, re-rate High.

### J17C-STD-002
- **Severity:** Low (latent; not currently reachable as a bug — see Evidence)
- **Category:** SEC (stale-cache class; same family as `internal/jose` SEC-001)
- **Location:** `dss/jades/jws.go:33` (`decodedPayload` field), `:138` (`DecodedPayload`), `:78`
  (`SetPayloadOctets` does not clear the cache)
- **Evidence:**
  ```go
  func (j *JWS) DecodedPayload() (*jose.Object, error) {
      if j.decodedPayload == nil {           // :139 — populate once, then cache
          ... j.decodedPayload = decoded ...
      }
      return j.decodedPayload, nil
  }
  func (j *JWS) SetPayloadOctets(payload []byte) {   // :78 — no `j.decodedPayload = nil`
      if j.IsRfc7797UnencodedPayload() { j.SetPayloadBytes(payload) } else { j.SetEncodedPayload(string(payload)) }
  }
  ```
- **Impact:** `decodedPayload` is cached from `UnverifiedPayload()` and is **never invalidated** by any payload
  setter (the embedded base's `SetPayloadBytes`/`SetEncodedPayload` and the wrapper's `SetPayloadOctets` all
  leave it untouched). Any mutation-then-read sequence would silently serve the *stale* parsed payload — the
  same class as `internal/jose` SEC-001, but the *verify* cache there is correctly reset (`SetSignature` clears
  `validSignature`), whereas this sibling is not. Not a live bug today: in the only call graph that reads it
  (verify path) the payload is set at parse time (`extractAndAddJWSSignature`) and, for detached signatures, by
  `jadesSignaturePayload` **before** `keyBindingSignatureReferenceValidation`→`sdHash`/`sdAlg` first read it
  (`jades_signature.go:665`, `:1199`, `:1221`), so the first (and cached) read is correct.
- **Recommendation:** Clear the cache on mutation — add `j.decodedPayload = nil` at the top of
  `SetPayloadOctets`, and (defensively) route the base `SetPayloadBytes`/`SetEncodedPayload` through it or
  document that the payload is immutable after parse. At minimum, a comment on the `decodedPayload` field
  stating the invariant, so the next mutation site doesn't quietly break it.

### J17C-INFO-001
- **Severity:** Info
- **Category:** SEC (test-only)
- **Location:** `dss/jades/testdata/crossgen/main.go:45,61` (hardcoded `testpassword`), `:67`
  (TSA policy OID `1.2.3.4.5.6.7.8.9`)
- **Evidence:**
  ```go
  signerEntry, err := loadKeyEntry(filepath.Join(selfDir, "signer_rsa.p12"), "testpassword")
  tspSource.SetTsaPolicy("1.2.3.4.5.6.7.8.9") // arbitrary unregistered test OID
  ```
- **Impact:** Test-only fixture generator (runs in CI, output handed to upstream Java DSS). The PKCS#12
  password is a throwaway test credential and the policy OID is an explicitly-unregistered placeholder — both
  correct and self-documenting (the header says so). No production path, no secret material. Noted for
  completeness only; the lens's "weak crypto / hardcoded secrets" check is satisfied (real RSA/EC keys, SHA-256,
  proper RFC 3161 via `KeyEntityTSPSource`).
- **Recommendation:** None required. (Optionally read the password from an env var if the `.p12` files are ever
  checked in unencrypted, but they are test fixtures and the value is meaningless.)

**PERF:** *no findings* — no repeated base64/JSON serialization in a hot loop and no unbounded allocation in the
in-scope files. The per-call schema re-parse and the `signedData` O(n²) sweep already logged are J17B-PERF-001
and J17A-PERF-001 respectively (already reviewed, not re-reported). `jades_signature_scope_finder.go`'s
`getHttpHeadersPayloadSignatureScope` re-serializes the headers once per scope request (not per header), which
is acceptable.

**STD (beyond the two findings above):** *clean.* Notable deliberate-port items verified as correct, not
reported: the `TstVD`/`TstVd` and `SigPSt`/`SigPst` **alias constants** (`jades_header_parameter_names.go`)
both resolve to the same wire value (`"tstVD"`/`"sigPSt"`), so read/write agree across the already-reviewed
sources (`jades_crl_source.go:74`, `jades_ocsp_source.go:75` use `…TstVD`; `jades_timestamp_source.go:210,436,515`
use `…TstVd`) and the in-scope `jades_level_baseline_lta.go`/`lt.go` (`IncorporateTstValidationData`); the
duplication is pinned by `jades_root_chunk_test.go:42-43` and documented in the file header — a code-smell
(Info, not a defect) rather than a bug. All other `Objects.requireNonNull`→panic sites are Java-parity and sit
under the facade `recovered()` wrapper.

### Tool log

Run from `dss/`:
```
$ gofmt -l jades
  # (no output) — clean
$ go vet ./jades/...
  # exit 0 — clean
$ (dss/) golangci-lint run --config=../.github/.golangci.yml ./jades/...
  # golangci-lint present (/opt/homebrew/bin/golangci-lint): "0 issues."
```
All three checks clean. No unavailable tools.

### Open questions

1. **`getDigest` reachability of a *wrong* scope (J17C-STD-001).** Confirmed the crypto verdict is independent of
   the scope (U17a, clean), so the dropped message-body scope is a reporting gap, not a false-accept. Confirm
   with the report/scope consumers (`validation/reports`, `esig` CLI) that none treat an *absent*
   `HTTPHeaderMessageBodySignatureScope` as "signature does not cover the body" — if any do, re-rate High.
2. **`decodedPayload` invariant (J17C-STD-002).** Confirm no future/`eaa`-tagged call site sets a JWS payload
   *after* a `DecodedPayload()` read (the current SD-JWT/EAA path does not). If one is added, the latent
   stale-cache becomes live — apply the `SetPayloadOctets` invalidation.
3. **Parallel-signing index (Info).** `jades_serialization_builder.go:52` dereferences
   `Signatures()[0]` when a payload is present; a degenerate payload-without-signatures object would index-out-
   of-range (recovered to an error by the facade). Confirm the analyzer never yields such an object in practice;
   if it can, add an explicit empty-check with a clear message.

