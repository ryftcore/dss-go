# Batch 14 — XAdES

## xades core — part 1 (unit U40)

**Scope.** 50 largest non-test files in `dss/xades/` (top half by size). Largest single
file in the repo: `xades_signature_builder.go` (1747). Part 2 (bottom 50 + `xades/extension/`)
is unit U41.

**Date.** 2026-08-29

**Depth.** Budgeted. 11 files >400 lines read as header + full `func`/`panic`/`== nil`
signature map + targeted verdict/security regions (level acceptance, digest/MDP handling,
property-set inclusion); all files 150–400 lines read in full; all files <150 lines read in
full. Every verdict/level path that could flip a level or panic was cross-checked against the
upstream Java at `~/Workspace/esig` (DSS 6.5.RC1 pin), including
`XAdESBaselineRequirementsChecker.java`, `XAdESSignatureBuilder.java`, `DSSXMLUtils.java`,
`XAdESLevelC/T.java`, `EncapsulatedTimestampScopeFinder.java`, and
`SignatureTimestampSource.java`.

Counts: **0 Critical, 0 High, 2 Medium, 0 Low, 4 Info** (6 total). No exploitable /
false-accept verification finding in this package.

### Findings

### T40-SEC-001

- **Severity:** Medium
- **Category:** SEC — timestamp scope over-attribution (verification-side)
- **Location:** `xades/xades_timestamp_source.go:482-487` (`GetTimestampScopes`),
  `xades/xades_timestamp_scope_finder.go:88` (`FilterCoveredSignatureScopes`), reached from
  the base `spi/validation/timestamp/signature_timestamp_source.go:1073`
  (`getTimestampScopes`), called by `validateTimestamps` (`:1027`, `:1053`).
- **Evidence:**
  ```go
  // xades_timestamp_source.go:482 — the CORRECT XAdES scope finder
  func (s *TimestampSource) GetTimestampScopes(timestampToken *validation.TimestampToken) []mscope.SignatureScope {
      timestampScopeFinder := NewTimestampScopeFinder()
  // signature_timestamp_source.go:1073 — what the base ACTUALLY calls
  func (s *SignatureTimestampSource[AS, SA]) getTimestampScopes(...) []scope.SignatureScope {
      timestampScopeFinder := validationscope.NewEncapsulatedTimestampScopeFinder() // base, XAdES-agnostic
  ```
  The base `EncapsulatedTimestampScopeFinder.FilterCoveredSignatureScopes`
  (`spi/validation/scope/encapsulated_timestamp_scope_finder.go`) "returns all by default"
  (`return f.Signature.SignatureScopes()`), whereas `xades.TimestampScopeFinder`
  filters to the `IndividualDataObjectsTimestamp` include subset. The base's private
  `getTimestampScopes` has no override hook (it is not on `SignatureTimestampSourceOverrides`),
  so `validateTimestamps()`'s call never reaches the XAdES override.
- **Impact:** For a timestamp token with `includes` (i.e. an
  `IndividualDataObjectsTimestamp`), the validation flow attributes the timestamp to **all**
  of the signature's scopes instead of the include-filtered subset — i.e. it
  **over-claims** timestamp coverage. This is a parity gap (Java's virtual dispatch reaches
  the XAdES override), not a demonstrated false-accept: scope results feed the
  scope-consumption / `VALIDATED` roll-up already reviewed in the validation-core batches
  (batch-10/11/13), which bound the verdict. It is a real hardening/parity defect on a
  verification-ish path, hence Medium rather than Low, and it is the single most
  security-relevant item in this unit.
- **Recommendation:** Add `GetTimestampScopes(*TimestampToken)` to
  `SignatureTimestampSourceOverrides` and route `validateTimestamps`'s call through it
  (the file header at `xades_timestamp_scope_finder.go:3-30` and
  `xades_timestamp_source.go:18-33` already describe this exact fix as "GAP flagged for
  integrator"). Alternatively, document the over-attribution in
  `docs/compatibility/known-gaps.md` if it is accepted (it is not currently listed there).

### T40-PERF-001

- **Severity:** Medium
- **Category:** PERF — unbounded retention of per-signature objects
- **Location:** three package-level `sync.Map` registries, all `Store`-only with **no
  `Delete`/eviction**: `xades/xades_reference_validation.go:35` (+`Store` `:66`),
  `xades/xades_signature_policy.go:22` (+`Store` `:50`,`:58`),
  `xades/xades_signature_builder.go:1723` (+`Store` `:1734`).
- **Evidence:**
  ```go
  var xadesReferenceValidationRegistry sync.Map // map[*model.ReferenceValidation]*ReferenceValidation
  ...
  xadesReferenceValidationRegistry.Store(&v.ReferenceValidation, v)  // xades_reference_validation.go:66
  // no .Delete anywhere in dss/xades (verified by grep: only .Load/.Store)
  ```
  `NewReferenceValidation` is called once per `ds:Reference` inside
  `Signature.ReferenceValidations()` (`xades_signature.go:803`), which is invoked for every
  signature in every document of every validation run; `SignaturePolicy` is built per
  `SignaturePolicyIdentifier` read.
- **Impact:** Every XAdES signature ever validated in-process leaves its
  `*ReferenceValidation` (holding a `*xmldsig.Reference` → DOM node tree) and
  `*SignaturePolicy` (holding `*xmldom.Node` transforms) pinned for the lifetime of the
  process. A long-lived server that validates N signatures with R references each holds
  O(N·R) reference objects + O(N) policy objects that are never released — a steady
  memory-growth / DoS vector on a verification workload, not a correctness bug.
- **Recommendation:** Since the registry exists purely to recover the concrete type from an
  embedded base pointer (Go has no covariant return through embedding), bound it: (a) attach
  the concrete value to the base struct directly (e.g. a typed field on
  `model.ReferenceValidation`) instead of a side map, or (b) add a `Delete` when the
  enclosing `Signature` is released / `ReferenceValidations()` result is replaced, or (c) if
  the leak is accepted, record it in `known-gaps.md`. Note this is a systemic consequence of
  the port's no-`instanceof` workaround (same pattern as batch-05 JAdES policy recovery); a
  single structural fix in `model`/`spi` would close it across formats.

### Info

### T40-INFO-001 — policy-transforms registry is an `instanceof` workaround

- **Category:** Info — deliberate-parity note
- **Location:** `xades/xades_signature_builder.go:1718-1745`
  (`SignatureBuilderRegisterPolicyTransforms` / `SignatureBuilderPolicyTransforms`)
- **Evidence:**
  ```go
  // Java: `if (signaturePolicy instanceof XmlPolicyWithTransforms)` — Go cannot recover the
  // outer value from *model.Policy (XmlPolicyWithTransforms embeds model.Policy by value).
  var xadesSignatureBuilderPolicyTransformsRegistry sync.Map // map[*model.Policy][]DSSTransform
  ```
- **Impact:** Correct behaviour, but a global registry keyed by `*model.Policy` that
  `incorporatePolicy` consults; a caller that sets a plain `*model.Policy` (not via
  `RegisterPolicyTransforms`) silently gets no `ds:Transforms` written — a behavioural
  footgun, not a bug. Flagged "FOR THE INTEGRATOR" in the file header; the clean fix is a
  change to the frozen `model.BLevelParameters`.
- **Recommendation:** None required for this unit; track with T40-PERF-001.

### T40-INFO-002 — substring ID matching in indented-signature replacement

- **Category:** Info — deliberate-parity note
- **Location:** `xades/dss_xml_utils.go:254` (`DSSXMLUtilsGetDocWithIndentedSignature`)
- **Evidence:**
  ```go
  if utils.IsStringNotEmpty(signatureAttrIdValue) && strings.Contains(signatureAttrIdValue, signatureId) {
  ```
- **Impact:** Uses substring (`strings.Contains`) rather than equality to match a signature
  `Id` while replacing its indented copy. Cross-checked 1:1 against Java
  `DSSXMLUtils.getDocWithIndentedSignature` (`DSSXMLUtils.java:290`,
  `signatureAttrIdValue.contains(signatureId)`) — **parity, not a divergence**. In theory two
  signatures whose Ids are substrings of one another could both match; Java has the identical
  property, so no finding.
- **Recommendation:** None.

### T40-INFO-003 — BASELINE-B `SignaturePolicyStore` requirement (m) is log-only

- **Category:** Info — deliberate-parity note
- **Location:** `xades/xades_baseline_requirements_checker.go:135-143`
  (`HasBaselineBProfile`)
- **Evidence:**
  ```go
  if signaturePolicyStoreAmount == 1 {
      if !b.IsSignaturePolicyIdentifierHashPresent() {
          // Upstream logs "SignaturePolicyStore shall not be present ... (requirement (m))!".
      }   // ← no `return false` here
  } else if signaturePolicyStoreAmount > 1 {
      return false
  }
  ```
- **Impact:** When a `SignaturePolicyStore` is present without a `SigPolicyHash`, the check
  only logs and does **not** fail the level. Verified **1:1** against Java
  `XAdESBaselineRequirementsChecker.java:149-156` (the `== 1 && !hashPresent` branch is
  `LOG.warn(...)` only; the `> 1` branch `return false`). This is a level-downgrading risk in
  principle, but it is Java-parity, so per the unit rule it is **not a finding** — recorded
  here so it is not mistaken for a Go-port error. (Contrast: the EPES branch at
  `:395-410` correctly does `return false`, matching its Java twin.)
- **Recommendation:** None (parity). If the team wants stricter BASELINE-B acceptance than
  Java, that is an upstream-alignment decision, not a port defect.

### T40-INFO-004 — `sync.Map` concurrency posture of the three registries

- **Category:** Info — observation
- **Location:** see T40-PERF-001
- **Evidence:** All three registries are `sync.Map` (concurrency-safe), but unlike
  `dssXMLUtilsMu` (which guards the two plain-map transform sets in `dss_xml_utils.go:77`),
  they are write-once-per-key with no lock coordination against the `Signature`'s own lazy
  caches (`CachedReferenceValidations`, `santuarioSignature`, `timestampAttributeMap`).
- **Impact:** No data race in the single-threaded signing/validation flow that the port
  targets; `sync.Map` is used defensively. The only shared mutable state in the service/builder
  is per-instance (the `Service` embeds `AbstractSignatureService` with no package-level
  mutable field; builders are per-call). **No goroutines are spawned anywhere in `dss/xades/`**
  (verified by grep).
- **Recommendation:** None.

### Verdict-path verification notes

**Level acceptance (SPECIAL FOCUS #1) — verified fail-closed.** Every `Has*Profile` /
`HasExtended*Profile` method in `xades_baseline_requirements_checker.go` was read in full and
cross-checked against Java. The acceptance logic is uniformly **fail-closed**:
- Cardinality checks use `!= 1` / `> 1` / `== 0` guards that `return false` on violation
  (e.g. `HasBaselineBProfile` `:55-200`, `HasExtendedBESProfile` `:335-385`).
- The one log-only branch (BASELINE-B `SignaturePolicyStore` requirement (m)) is confirmed
  Java-parity (T40-INFO-003), **not** a Go divergence.
- `HasExtendedCProfile` correctly requires `CompleteCertificateRefs(V2)` (`:409-420`) and
  `CompleteRevocationRefs` (`:431-439`) with the self-signed exemption; `HasExtendedXProfile`
  requires a `RefsOnlyTimestamp(V2)`/`SigAndRefsTimestamp(V2)` (`:444-457`);
  `HasExtendedERSProfile` requires `SealingEvidenceRecords >= 1` (`:523-535`).
- `DataFoundUpToLevel` (`xades_signature.go:1229+`) is a pure dispatch over the above and
  cannot promote a signature past what the checks return.
- **No inverted `!` and no always-true branch found** in any level-acceptance path. The
  requirement-check functions that *could* wrongly return "satisfied" were specifically
  audited; all are fail-closed.

**Builder chain (SPECIAL FOCUS #2) — verified clean.**
`xades_signature_builder.go` (`Build()` `:233-302`) drives a fixed sequence and
**propagates every error**: `AssertSignaturePossible` → `ensureConfigurationValidity` →
`IncorporateFiles` → `IncorporateSignedInfo` → `IncorporateKeyInfo` → `IncorporateObjects` →
reference incorporation → canonicalization. No swallowed error on the signing path.
`IncorporateSignedProperties` → `IncorporateSignedSignatureProperties` /
`incorporateSignedDataObjectProperties` each `return` the first non-nil error, so a
required property is either written or the build fails (no silent drop). `AddContentTimestamp`
(`:1625-1676`) fails on an empty canonicalization method rather than emitting an empty
timestamp. `incorporateSigningCertificateV1/V2` propagate `IncorporateCert` errors.
The `Objects.requireNonNull`-equivalent panics (`DataObjectFormat cannot be null!` `:1194`,
`CommitmentTypeQualifier cannot be null!` `:1441`, `CommitmentType cannot be null!` `:1480`)
are the sanctioned Java idiom (T40-INFO-004 note; PORTING.md). MDP/digest builders
(`xades_timestamp_message_digest_builder.go`) guard `references == 0` and `includes`
referencedData before computing (`:190-280`) and `defer recover()` to the documented
"empty digest" collapse — nil-safe.

**Timestamp source (SPECIAL FOCUS #3) — trust handling verified; one parity gap.**
`xades_timestamp_source.go` is construction-side (it builds the message-imprint for a
`TSPSource` to timestamp) and does **not** itself trust/verify a TSP response — the
`TSPSource` is caller-supplied (`KeyEntityTSPSource` for tests; no online source is ported,
see `known-gaps.md` "Not ported at all"). The token is built from the caller's
`TimeStampResponse` bytes (`all_data_objects_time_stamp_builder.go:142-149`) and no
network trust decision is made here. The **only** verification-side trust issue found is the
scope-finder dispatch gap (T40-SEC-001), which is over-attribution (over-claim), not a
trust-bypass. `MakeTimestampToken` (singular) `panic`s by design (Java throws
`UnsupportedOperationException`); `MakeTimestampTokens` (plural) is the real path.

**`dss_xml_utils.go` vs batch-02 X10-SEC-001 (SPECIAL FOCUS #4) — cross-referenced, not re-reported.**
This file *consumes* the X10-SEC-001 pattern via `xmlutils.XPathUtilsGetElementByIdWithQuery`
(e.g. `DSSXMLUtilsGetObjectById` `:691`, `DSSXMLUtilsIsKeyInfoReference` `:648`,
`DSSXMLUtilsIsSignaturePropertiesReference` `:657`, and the checker's
`getSignedPropertiesById` `xades_signature.go:1037`). It does **not** re-implement the
single-quote-interpolation, so per the unit rule it is "same class as X10-SEC-001" and is not
re-reported. `DSSXMLUtilsApplyTransforms` (`:985+`) and
`DSSXMLUtilsGetReferenceOriginalContentBytes` (`:496`) propagate transform/dereference errors
correctly.

**Concurrency (SPECIAL FOCUS #5) — no shared mutable state, no goroutines.**
- No `go func`/goroutines anywhere in `dss/xades/` (verified by grep).
- Only package-level mutable state is the three `sync.Map` registries (T40-PERF-001) and the
  two transform sets guarded by `dssXMLUtilsMu` (`dss_xml_utils.go:77`). Both are
  concurrency-safe; `sync.Map` is used defensively.
- `Service` (`xades_service.go`) and `AbstractSignatureBuilder` are per-instance; their
  mutable fields (`Built`, `Documents`, the cached DOMs, `statusMap`
  `xades_baseline_requirements_checker.go:700`) are per-object and not shared across goroutines.
  No lock is needed for the port's single-threaded use, and none is missing for the
  `sync.Map` paths.

### Open questions

1. **T40-PERF-001 structural fix:** the three `sync.Map` registries are a systemic consequence
   of Go's lack of covariant-return-through-embedding (the port's documented no-`instanceof`
   workaround). A single change in the frozen `model`/`spi` packages (attach the concrete
   value to the base struct instead of a side map) would close the leak across XAdES, CAdES,
   and JAdES simultaneously. **Question for the integrator:** is the leak accepted (and
   documented in `known-gaps.md`), or is the structural fix in scope for a follow-up unit?
2. **T40-SEC-001 severity arbitration:** the over-attribution is bounded by the
   scope-consumption roll-up reviewed in the validation-core batches, so it is rated Medium.
   If a downstream consumer trusts `TimestampToken.TimestampScopes()` directly (outside the
   executor), the effective severity is higher. **Question:** is there any consumer path that
   reads `GetTimestampScopes` output without going through `validateTimestamps`'s roll-up?
3. **`xadesSignatureBuilderPolicyTransformsRegistry` footgun** (T40-INFO-001): a plain
   `*model.Policy` silently omits `ds:Transforms`. **Question:** should `SetSignaturePolicy`
   on a plain policy with transforms emit a warning, or is the register-first API contract
   sufficient?

### Tool log (run from `dss/`)

- `gofmt -l xades/` → **clean** (no output).
- `go vet ./xades/...` → **clean** (exit 0, no diagnostics).
- `golangci-lint run --config=../.github/.golangci.yml ./xades/...` → **0 issues** (exit 0).
- `go test ./xades/... -count=1 -timeout 10m` → **pass**: `xades` 2.394s,
  `xades/definition` 1.254s; `xades/extension` and `xades/tsl` have no test files.
- `grep`/`find` freely (level-acceptance sweep, `panic(` sweep, `sync.Map`/registry
  eviction sweep, goroutine sweep).
- Java cross-checks at `/Users/utain/Workspace/esig/dss`:
  `find … -name 'XAdESBaselineRequirementsChecker.java'`,
  `sed -n` on `XAdESSignatureBuilder.java` (`incorporateSigningTime`), `DSSXMLUtils.java`
  (`getDocWithIndentedSignature`), `SignatureTimestampSource.java` (`getTimestampScopes`),
  `EncapsulatedTimestampScopeFinder` (base `FilterCoveredSignatureScopes`).

## xades/definition + xades/tsl (unit U42)

**Scope.** xades/definition (19/5807L) + xades/tsl (6/765L)
**Date.** 2026-08-29
**Depth.** Full. Every file in both packages was read in full (all ≤934 lines; the three
large element files read as header + complete constant/tag-map/method-body maps, and every
wire value programmatically extracted and diffed against the Java twin). `xades/tsl/` all six
files read end-to-end. Every structural element (element sets, unsupported-method sets, XPath
chains, tag/attribute names, namespaces, path chains, builder surface) was cross-checked 1:1
against the upstream Java at `/Users/utain/Workspace/esig/dss` (DSS 6.5.RC1 pin) using
mechanical extraction + set/diff, not visual inspection.

**What this unit actually is.** Contrary to the unit brief's framing, neither subpackage
contains *level-requirement logic*, *property-set inclusion logic*, or *trust-service-type →
XAdES-property mapping*:
- `xades/definition/` is **pure XML vocabulary**: element/attribute enums (with wire tag
  names), XPath path builders, and namespace constants, per schema version (1.1.1/1.2.2/1.3.2/
  1.4.1 + TL + EvidenceRecord). It has no verdict/level/trust logic at all. The level
  requirement definitions (SPECIAL FOCUS #1) live in `xades/xades_baseline_requirements_checker.go`
  and were already audited in U40 (fail-closed verified, `T40-INFO-003`). The property-set
  inclusion definitions live in the `xades` core builder/checker (U40), not here.
- `xades/tsl/` is **TL conformance + TL-signing parameter builders** (`TLStructureVerifier`,
  `XAdESTrustedListUtils`, abstract + V5/V6 builders). It does **not** map a TSL
  trust-service-type URI to an XAdES property (that is the `tsl`/`trustedlist` package,
  batch-06, out of scope here). What it *does* do on the trust side — the
  `SigningCertificate`/`SigningCertificateV2` (and ProductionPlace/SignerRole V1-vs-V2)
  conformance check in `validateSignatureElement` — is the security-relevant part and is
  verified below.

Counts: **0 Critical, 0 High, 0 Medium, 0 Low, 2 Info** (2 total). No exploitable /
false-accept verification finding in this unit. No level-downgrade, no false-qualification,
no race, no leak.

### Findings

### Info

### T42-INFO-001 — XSD structural validation is a permanent no-op (documented stance, re-confirmed)

- **Category:** Info — deliberate-parity / known-limitation note (not re-reported as a new gap)
- **Location:** `xades/tsl/tl_structure_verifier.go:131-150`
  (`validateAgainstXSD` + `TrustedListXSDUtils`/`TrustedList211XSDUtils`), reached from
  `validateTrustedListV5`/`validateTrustedListV6`, via `xades.DSSXMLUtilsValidateAgainstXSD`
  (`xades/dss_xml_utils.go:443-445`, `return nil`).
- **Evidence:**
  ```go
  // dss_xml_utils.go:443 — the XSD hook both TL versions call
  func DSSXMLUtilsValidateAgainstXSD(xsdUtils StructureXSDUtils, source *xmldom.Node) []string {
      return nil   // always "no errors"
  }
  // tl_structure_verifier.go:144-150 — opaque placeholders, never dereferenced
  var trustedListXSDUtils xades.StructureXSDUtils = struct{}{}
  ```
- **Impact:** `TLStructureVerifier`'s XSD conformance step never rejects a malformed TL — only
  the *semantic* checks (namespace root, `ds:Signature` presence, `QualifyingProperties` →
  `SignedProperties` → `SignedSignatureProperties`, and the `SigningCertificate(V2)` V1/V2
  match) run. This is the port's documented "no runtime XSD validation" stance (see
  batch-06-trust U19: "no runtime XSD validation is performed … treated as deliberate, not
  findings") and is consistent with the XAdES `StructureValidator` no-op in U40's package. It is
  **not** a false-accept *of a signature*: this verifier is on the *signing* side (it guards
  "may I build a TL signature over this document"), not the validation side that issues a
  trust verdict. Recorded here so the XAdES-TL consumer's reliance on it is explicit and is not
  mistaken for a Go-port defect.
- **Recommendation:** None required for this unit (consistent with the accepted stance). If the
  team later wants to surface "TL XSD conformance is not enforced in the Go port" as a user-facing
  limitation, add a line to `docs/compatibility/known-gaps.md` — it is currently only in the
  review batches and code comments.

### T42-INFO-002 — two doc-comment typos copied verbatim from Java (parity, not a defect)

- **Category:** Info — deliberate-parity note
- **Location:** `xades/definition/xades_element.go:83` and `:161`
- **Evidence:**
  ```go
  // :83  ElementCommitmentTypeQualifiers gets the "CommitmentTypeQualifies" element.  (missing 'r')
  // :161 ElementIssuer gets the "IssueTime" element.                                   (copy-paste of IssueTime's doc)
  ```
- **Impact:** The *identical* typos exist in the Java Javadoc
  (`XAdESElement.java:192` `"CommitmentTypeQualifies"`, `:367/:374` both `getElementIssueTime`
  and `getElementIssuer` documented as `Gets "IssueTime" element`). These are comments only — the
  method names, the tag-name maps, and the path/element wiring are all correct and were
  verified 1:1 (see notes). Per the unit rule, Java-parity is **not a finding**; recorded so it is
  not mistaken for a Go error and so a future "fix" does not accidentally rename the method (which
  would *break* parity with the `XAdESElement` interface).
- **Recommendation:** None.

### Verdict-path verification notes

**SPECIAL FOCUS #1 — level requirement definitions.** *Not present in this unit.*
`xades/definition/` contains no `Has*Profile`/level logic (grep for
`baseline|profile|level|requirement|dataobject` over the package, excluding comments/vocabulary,
returns nothing). The level requirement definitions were audited in U40
(`xades_baseline_requirements_checker.go`) and confirmed **fail-closed**: cardinality guards
(`!= 1`/`> 1`/`== 0` → `return false`), the one log-only BASELINE-B `SignaturePolicyStore`
branch confirmed Java-parity (T40-INFO-003), and `HasExtendedC/X/ERS` all require their
respective V2/revocation/evidence elements. **No incomplete or inverted level-requirement
definition found in this unit** — there is none here to be incomplete.

**SPECIAL FOCUS #2 — TSL mapping.** *No trust-service-type → XAdES-property mapping in this unit.*
`xades/tsl/` is the XAdES-side *consumer* (TL conformance + signing builders), not the TST→property
mapper (that is batch-06's `tsl`/`trustedlist`). The security-relevant trust logic that *is* here —
`TLStructureVerifier.validateSignatureElement`'s V1/V2 `SigningCertificate`/`ProductionPlace`/
`SignerRole` check — was read in full and cross-checked 1:1 against Java `TLStructureVerifier.java`:
- The `v2Expected ^ doesMatch(elem, SIGNING_CERTIFICATE_V2)` XOR is ported exactly (Go `!=`, same
  truth table: a V1 element is rejected when V2 is expected and vice-versa; both present → the
  non-matching one is rejected). **No inverted condition.**
- `doesMatch` = `IsSameTagName(localName) && URI()==namespaceURI`, identical to Java, so an
  element from a *wrong namespace* is not matched (no cross-namespace false-qualification).
- `getChildElement` returns the child only when `len==1` (mirrors Java's
  `collectionSize==1`), so a *duplicated* `SignedProperties`/`SignedSignatureProperties` is
  rejected rather than silently taking the first — fail-closed, matching Java.
- `verifySignatureElementPresence` enforces "no `ds:Signature` when signing / one present when
  validating" identically. **No false-qualification path found.**

**SPECIAL FOCUS #3 — property-set / definition inclusion correctness + nil-safety.** Verified:
- **Element sets:** for each version the Go `Element*` method set, the Java `getElement*` set,
  and the **unsupported-method (panic) set** all match exactly — 110 names each for 1.1.1/1.2.2/
  1.3.2; unsupported counts 20 (1.1.1) / 17 (1.2.2) / 0 (1.3.2) identical to Java; the 1.4.1 set
  (17) and TL set (70) and EvidenceRecord set (3) all match. No method returns a *different*
  element than its name (the only "miswire" candidates are Go case-normalisation of the constant
  identifier — `TimeStamp`→`Timestamp`, `Id`→`ID`, `MimeType`→`MIMEType` — which is the wire tag
  preserved correctly, not a rewire).
- **XPath chains:** all 119 `XPathQuery` methods per version (1.1.1/1.2.2/1.3.2) were extracted and
  the element-chain (and the `nil`/`null` returns) diffed against Java — **all identical**.
- **Wire names:** element tag-name maps (95/94/110 per version), attribute maps (6/6/7/3), and the
  `TrustedListPath` chains (7) and all namespace URIs/prefixes (`xades111/122/132/141`, `xadesen`,
  `dsig-filter2`, `tl`) are byte-identical to Java, *including* Java's own irregular tags
  (`SchemeExtension→"SchemeExtensions"`, `TSLAddress→"TSPAddress"`, `DateTime→"dateTime"`,
  `DigitalID→"DigitalId"`, `SPURI`, `XAdESTimeStamp`).
- **nil-safety:** `Validate`/`ValidateDocument`/`Init…`/`XAdESTrustedListUtils…` guard their
  document args with the sanctioned `Objects.requireNonNull`-equivalent `panic` (Java throws
  `NullPointerException`; PORTING.md idiom) and return clean error strings otherwise.
  `XAdES*ElementFromTagName` iterates a fixed `…ElementValues()` **slice** (deterministic, no
  Go map-iteration-order leak), not a map. `GetTSLVersionIdentifier` handles the non-digit /
  absent case by returning `nil` (no panic) before `strconv.Atoi`. `validateNamespace` nil-checks
  via the `Name.Local`/`Name.Space` access on an already-DOM-built root (no nil deref).

**SPECIAL FOCUS #4 — concurrency / shared mutable state.** Verified clean:
- **No goroutines, no `sync.*`, no mutex** anywhere in `xades/definition/` or `xades/tsl/`
  (grep). The only package-level state is *immutable* `var` blocks of constants/enum values
  (`XAdESNamespace*`, `TrustedListNamespaceNS`, the tag-name maps, `trustedListXSDUtils`) — all
  read-only after init, so safe to share across goroutines.
- The two `init()`s that call `XPathUtilsRegisterNamespace(TrustedListNamespaceNS)`
  (`tl_structure_verifier.go:30`, `xades_trusted_list_utils.go:19`) are idempotent by design
  (`RegisterNamespace` replaces an existing prefix — see `xml/utils/xpath_utils.go:44-58`) and
  run once at program start before any goroutine; no data race.
- `TLStructureVerifier` and the builders are per-instance (setters mutate `v`/`b` on the
  concrete pointer); nothing is shared across goroutines, matching Java's non-thread-safe
  builder semantics. **No race, no leak** in this unit (contrast with U40's `sync.Map`
  registries, which are in `dss/xades/` proper and already reported as T40-PERF-001).

**Builders (signing-side) — verified 1:1.** `AbstractTrustedListSignatureParametersBuilder.Build`
sets `ENVELOPED` + `XAdES_BASELINE_B` + `En319132(isEn319132())` + the enveloped-signature
`DSSReference` (uri `""`, contents = TL document, digest = `referenceDigestAlgorithm` default
`SHA512`, transforms = `EnvelopedSignature` + `Canonicalization(EXCLUSIVE)`), exactly matching
Java. V5 (`isEn319132=false`, target `5`) and V6 (`isEn319132=true`, target `6`) match Java's
`TrustedListV5/V6…Builder`. `AssertConfigurationIsValid` → `XAdESTrustedListUtilsValidateUnsignedTrustedList`
→ version-match + `TLStructureVerifier(signingMode=true)` — identical control flow to Java
(`assertConfigurationIsValid`/`validateUnsignedTrustedList`). The base
`document.AbstractSignatureParametersBuilder.Build` does **not** call `AssertConfigurationIsValid`
in either Go or Java (it is a separate caller-invoked step), so there is no missing virtual-dispatch
hook here (unlike U40's T40-SEC-001).

### Open questions

1. **XSD no-op documentation (T42-INFO-001):** the "no runtime XSD validation" stance is stated in
   the code comments and in the batch-06 review, but not in `docs/compatibility/known-gaps.md`.
   Should it be added there as a user-facing limitation (the XAdES-TL conformance path and the
   XAdES `StructureValidator` both rely on it), or is the existing "not ported / deliberate"
   treatment sufficient? (Integration decision, not a code defect.)
2. **Scope framing for future units:** the unit brief assumed `xades/definition/` holds
   "level-requirement / property-set definitions". It does not — it is pure XML vocabulary, and the
   level/property-set logic it feeds is in `dss/xades/` (U40). If a future unit is scoped by the
   *name* `definition/` expecting verdict logic, it will find none; the real level logic is in
   `xades_baseline_requirements_checker.go`. Flagging to prevent a duplicate or a missed audit.

### Tool log (run from `dss/`)

- `gofmt -l xades/definition/ xades/tsl/` → **clean** (no output).
- `go vet ./xades/definition/... ./xades/tsl/...` → **clean** (exit 0, no diagnostics).
- `golangci-lint run --config=../.github/.golangci.yml ./xades/definition/... ./xades/tsl/...`
  → **0 issues** (exit 0).
- `go test ./xades/definition/... ./xades/tsl/... -count=1` → **pass**:
  `xades/definition` 0.603s; `xades/tsl` has no test files.
- `grep`/`find` freely: level/property-set/TST-mapping scope check (none), `panic(`/`== nil`/
  `go func`/`sync.`/mutable-`var` sweep (only sanctioned panics + immutable `var` blocks + 2
  idempotent `init()`s), `FromTagName` iteration (slice, deterministic).
- Mechanical 1:1 diff vs Java at `/Users/utain/Workspace/esig/dss` (element sets +
  unsupported-method sets, XPath chains, tag/attribute/namespace wire names, `TrustedListPath`
  chains, `TLStructureVerifier`, `AbstractTrustedListSignatureParametersBuilder`,
  `XAdESTrustedListUtils`, V5/V6 builders) via `python3` extraction + set/diff.
- `find /Users/utain/Workspace/esig/dss -path '*xades*definition*' -name '*.java'` and
  `-path '*xades*tsl*'` to locate all Java twins (all located and cross-checked).

## xades core — part 2 + extension (unit U41)

**Scope.** 52 smallest non-test files in `dss/xades/` (13–107 lines; the bottom half of the
100 core files, complementing U40's top 50) + all of `dss/xades/extension/` (3 files).
**Date.** 2026-08-29
**Depth.** All 55 in-scope files read in full — the largest is 107 lines, so none required
budgeting. Supporting code read as needed for verdict/security context: `xml/utils/xpath_utils.go`
(`XPathUtilsGetElementByIdWithQuery`'s `defer recover()` + `DomUtilsGetId` normalization),
`internal/xmldom/navigate.go` (`Elements()` nil-safety), `internal/xmldsig/sigalg.go` +
`enumerations/signature_algorithm.go` (RIPEMD160 signature-method URI dispatch), and
`xades/tsl/tl_structure_verifier.go` + `xades/dss_xml_utils.go:443` (`DSSXMLUtilsValidateAgainstXSD`)
to assess the XSD no-op's blast radius. Every file with a verdict/trust/panic path was
cross-checked against the upstream Java at `~/Workspace/esig` (DSS 6.5.RC1 pin), including
`XAdESTimestampScopeFinder.java`, `XAdESAttributeIdentifier.java` (`getOrder`/`writeInt`),
`EnvelopingSignatureBuilder.java` (`getUri().substring(1)`), `XAdESStructureValidatorFactory.java`
(`getInstance` lazy-init), and `EncapsulatedTimestampScopeFinder.java` (`findTimestampScope`).

Counts: **0 Critical, 0 High, 0 Medium, 0 Low, 3 Info** (3 total). No exploitable /
false-accept verification finding in this unit; no false level-acceptance; no panic on
hostile (attacker-supplied) input. Two items are cross-referenced to U40 (T40-SEC-001,
T40-PERF-001) rather than re-reported, and one is cross-referenced to batch-06-trust.

### Findings

### Info

### T41-INFO-001 — XSD structure validation is a permanent no-op (not in `known-gaps.md`)

- **Category:** Info — deliberate-parity / documentation gap
- **Location:** `xades/xades_structure_validator.go` (`Validate()` `:62-66` always
  `return utils.IsCollectionEmpty(nil)` → `true`), driven by `xades/dss_xml_utils.go:443`
  (`DSSXMLUtilsValidateAgainstXSD` unconditionally `return nil`). Consumed at
  `xades/xades_signature.go:1330-1334` (`Signature.ValidateStructure`) and, trust-relevant,
  at `xades/tsl/tl_structure_verifier.go:117,130,138-140` (TL V5/V6 structure conformity,
  feeding `dss/tsl/abstract_parsing_task.go:244-256` → `StructureValidationMessages`).
- **Evidence:**
  ```go
  func DSSXMLUtilsValidateAgainstXSD(xsdUtils StructureXSDUtils, source *xmldom.Node) []string {
      return nil   // always: "the structure is always considered valid" (header, dss_xml_utils.go:439-444)
  }
  func (v *StructureValidator) Validate() bool {
      v.errors = DSSXMLUtilsValidateAgainstXSD(v.getUtils(v.xadesPath), v.signatureElement)
      v.errorsSet = true
      return utils.IsCollectionEmpty(v.errors)   // → always true
  }
  ```
- **Impact:** XAdES/TL structure validation (`validateStructure()` / TL `verifyTLVersionConformity`)
  reports the structure as valid for **any** input, because no XML-Schema validator package
  exists in the port. This is a report-fidelity gap, **not a core-verdict false-accept**: the
  structural-validation result only populates the `structuralValidation` report field
  (`validation/reports/diagnostic/diagnostic_data_builder.go:598,654,869,921`) and TL
  `structureValidationMessages`; the cryptographic verdict path is independent (cross-checked
  against Java, where `getStructureValidationMessages` is likewise only read by
  `DiagnosticDataBuilder`). Already surfaced by **batch-06-trust** (the "documented
  'no runtime XSD validation' stance" and T19-SEC-001) and documented in the three code
  headers — but **not** in `docs/compatibility/known-gaps.md` (verified by grep). Hence it is
  an *undocumented* divergence, which the porting rules treat as reportable.
- **Recommendation:** Add a one-line entry to `known-gaps.md` ("Supported with limits" →
  new "XML structure validation" bullet) stating that `StructureValidator`/`TLStructureVerifier`
  always report the structure as valid because no XSD validator is ported. No code change is
  in scope here (the underlying `ValidateAgainstXSD` is frozen-package territory and the
  gap is deliberate).

### T41-INFO-002 — `StructureValidatorFactory` lazy singleton is unsynchronized

- **Category:** Info — observation (Java-parity)
- **Location:** `xades/xades_structure_validator_factory.go:22-27`
- **Evidence:**
  ```go
  var xadesStructureValidatorFactorySingleton *StructureValidatorFactory
  func StructureValidatorFactoryGetInstance() *StructureValidatorFactory {
      if xadesStructureValidatorFactorySingleton == nil {
          xadesStructureValidatorFactorySingleton = &StructureValidatorFactory{}
      }
      return xadesStructureValidatorFactorySingleton
  }
  ```
- **Impact:** A naive check-then-set on a package-level pointer is a data race (UB) under
  concurrent `ValidateStructure()` calls. Two mitigations make it benign in practice: (a)
  it is **1:1 Java-parity** — upstream `XAdESStructureValidatorFactory.getInstance()`
  (`:45-51`) has the identical unsynchronized `if (singleton == null) singleton = new …`;
  and (b) `StructureValidatorFactory` is **stateless** (no fields), so a torn double-init
  yields two equivalent instances with no observable difference. U40 verified no goroutines
  are spawned in `dss/xades/`, so the port's single-threaded use never triggers it. `go vet`
  and `golangci-lint` do not flag it (a race requires concurrent access).
- **Recommendation:** None for this unit. If the library ever grows a concurrent API,
  replace the lazy-init with `sync.Once` (or pre-allocate at `init()`).

### T41-INFO-003 — `SetDisableXSWProtection` is a public, security-relevant opt-out knob

- **Category:** Info — observation (Java-parity)
- **Location:** `xades/xml_document_validator.go:67-70`
- **Evidence:**
  ```go
  // SetDisableXSWProtection: NOT RECOMMENDED - this parameter allows disabling protection against
  // XML Signature wrapping attacks (XSW). It disables the research by XPath expression for defined
  // Type attributes. Port of setDisableXSWProtection(boolean).
  func (v *XMLDocumentValidator) SetDisableXSWProtection(disableXSWProtection bool) {
      v.DocumentAnalyzer().SetDisableXSWProtection(disableXSWProtection)
  }
  ```
- **Impact:** A caller can explicitly disable XSW (signature-wrapping) protection, widening
  the class of documents that validate. This is a **faithful port** of the upstream
  `setDisableXSWProtection` API (documented "NOT RECOMMENDED" upstream) and is opt-in —
  the default is protected — so it is not a defect. Recorded only because it is a
  security-relevant surface in an otherwise boilerplate file, so a future reader does not
  mistake it for a porting error or "fix" it away.
- **Recommendation:** None. Keep the "NOT RECOMMENDED" doc line.

### Verdict-path verification notes

**Level checkers (SPECIAL FOCUS #1) — no false-accept in this unit.** The level-*acceptance*
checker (`xades_baseline_requirements_checker.go`, `Has*Profile` / `DataFoundUpToLevel`) sits in
U40's top-50 and was already verified **fail-closed** there (no inverted `!`, no always-true
branch). The only level-related file in *this* bottom-50 set is
`xades/xades_level_baseline_b.go`, which is the **construction** side — `LevelBaselineB`, the
B-level *signing profile* (`GetDataToSign` / `SignDocument` / `SignDocuments`), not an acceptance
checker. It is clean: `GetDataToSignForDocuments` and `SignDocuments` both propagate the
builder factory / `builder.SignDocument` errors (`return nil, err`), and the re-read of
`parameters.GetContext().Builder()` at `:101-103` is a documented no-op kept verbatim for
parity. There is no "requirement satisfied" return path in this file to get wrong.

**Property builders / digest builders (SPECIAL FOCUS #2) — nil-safe, no dropped errors.**
- `data_object_format_builder.go` `Build()` `:33-58`: correctly implements EN 319 132-1 §6.3(k) —
  skips counter-signature references and (via the caller) the SignedProperties element;
  `toDataObjectFormat` `:59-70` `panic`s on a nil reference (Java `Objects.requireNonNull`) and
  falls back to `binary` MIME when the content has none. Clean.
- `xades_attribute_identifier.go` `AttributeIdentifierBuild` `:42-66`: `getOrder` `:73-84` is a
  linear sibling scan (O(n), fine); the order is encoded with
  `binary.BigEndian.PutUint32` — verified **1:1** against Java `dos.writeInt(getOrder(node))`
  (`XAdESAttributeIdentifier.java:56-58`, 4-byte big-endian). The serialize-error `panic` is
  documented Java-parity. Clean.
- `xades_sig_properties.go` / `xades_unsigned_sig_properties.go` /
  `xades_signed_data_object_properties.go`: `IsExist()` guards nil; `Attributes()` guards
  `signaturePropertiesDom == nil`; the `*Dom` getters swallow the XPath "not found" error and
  return nil (Java returns null on the same catch) — a missing property correctly reads as
  *absent*, never *present*. Clean.
- `xades_timestamp_parameters.go`, `xades_evidence_record_incorporation_parameters.go`,
  `xml_policy_with_transforms.go`: setters `panic` on the Java `IllegalArgumentException` /
  `Objects.requireNonNull` cases (empty canonicalization, nil/foreign-URI ERS namespace);
  `XmlPolicyWithTransforms.IsEmpty/Equals` are faithful overrides (digest is user-computed per
  the documented NOTE, not a builder duty here). Clean.
- `reference_id_provider.go`: deterministic ID; `SetReferenceIdPrefix` `panic`s on blank
  (parity); `ReferenceId` is a monotonic counter. Clean.
- `xades_signature_identifier_builder.go`: `CounterSignaturePosition` / `SignaturePosition`
  are linear scans (O(n) over signature/counter-signature count — fine); the
  `masterSignature.(*Signature)` type-assert at `:42` is safe because `CounterSignatures()`
  yields XAdES `*Signature` values on this code path. Clean.

**`extension/` subpackage (SPECIAL FOCUS #3) — what it is + inclusion correctness.**
`dss/xades/extension/` is the XAdES **extend-to-a-higher-level** entry point (adding a content
timestamp, revocation data, or an archive timestamp to an *existing* XAdES signature), i.e.
Java's `XAdESDocumentExtender` + its `SignedDocumentExtenderFactory`, not the B-LT/B-LTA
*properties* themselves (those live in U40's `xades_signature_builder.go` / timestamp builders).
- `extension/xades_document_extender.go`: `CreateSignatureService()` `panic`s when
  `CertificateVerifier == nil` (Java `Objects.requireNonNull`, `:47`); `IsSupportedParameters` /
  `IsSupportedService` are type-asserts against the concrete XAdES types (correct dispatch);
  `SignatureForm()` returns `XAdES`. The five override hooks registered via
  `InitAbstractDocumentExtender` match Java's `AbstractDocumentExtender`. Clean.
- `extension/xades_document_extender_factory.go`: `init()` registers the factory into
  `document.RegisterSignedDocumentExtenderFactory` — the exact stand-in for upstream's
  `META-INF/services/…SignedDocumentExtenderFactory` ServiceLoader entry; `Create` returns a
  fresh extender. Clean.
- No property-inclusion logic lives here (the extender delegates to `xades.NewService`), so
  there is no dropped-property risk in the subpackage itself.

**Signature-construction transforms (SEC) — injection surface is guarded or static.**
- `xpath_transform.go` / `xpath_enveloped_signature_transform.go` /
  `xpath2_filter_enveloped_signature_transform.go`: the XPath filter expressions are built from
  the *namespace's own prefix* (`xmlDSigNamespace.Prefix()`) plus fixed constants — static, not
  attacker-supplied. Actual execution routes through `internal/xmldsig`, whose
  `EnforcedResolverFragment` applies the `checkValueForXPathInjection` / `xpathCharFilter` guard
  (see `enforced_resolver_fragment.go`'s header and `internal/xmldsig/resolver.go:130-160`).
  The `XPathEnvelopedSignatureTransform` "cannot be used with parallel signatures" warning is
  carried over verbatim. Clean.
- `xades_signature_integrity_validator.go` `Verify()` `:31-44`: returns `(false, *DSSError)`
  when `CheckSignatureValue` errors — i.e. a verification *failure* is never conflated with a
  valid signature; only a clean `true` is returned for a good signature. Correct.
- `counter_signature_resolver.go` `Resolve()` `:38-47`: returns an **error** (not a silent
  `nil`/accept) when the counter-signed `ds:SignatureValue` is not found by URI; `createFromNode`
  `panic`s only on a DOM serialize failure (documented Java-parity); `CanResolve` gates on
  `resolveNode != nil`. The lookup goes through `XPathUtilsGetElementByIdWithQuery`, which
  normalizes the ID via `DomUtilsGetId` and `defer recover()`s to `nil` — no XPath injection, no
  panic on a missing node. Clean.
- `sp_doc_digest_as_in_specification_transform.go` `PerformTransform` always returns the
  verbatim Java `IllegalArgumentException` — the transform is a digest-rule marker, and it is
  correctly refused when misused for reference processing. `base64_transform.go` /
  `enveloped_signature_transform.go` are documented identity/no-op `PerformTransform` (parity).
- `xslt_transform.go`: writing the `ds:Transform` is byte-pinned (KAT `tf-xslt`); **executing**
  it is a deliberate `xmldsig.ErrForbiddenTransform` refusal — the file header records the Java
  oracle (`gen/RefsOracle.java` "apply-xslt") transforming the fixture, so this is a *documented*
  divergence owned by frozen `internal/xmldsig`, not a silent drop. Not a finding.
- `dss_object.go` / `xml_element_signature_scope.go` / `xml_root_signature_scope.go` /
  `xml_full_signature_scope.go` / `manifest_entry_signature_scope.go` / `xpointer_signature_scope.go`:
  pure getters / scope descriptions. `XPointerSignatureScope` `Type()` correctly returns `Full`
  only for a root XPointer and `Partial` otherwise; `ManifestEntrySignatureScope` returns `Full`.
  `reference_output_type.go` / `dss_transform.go` / `signature_builder.go` /
  `xades_signature_profile.go` / `digest_document_xml_signature_input.go`: interfaces / constants.
  `signature_ecdsaripemd160.go` / `signature_rsaripemd160_at.go`: URI-constant stubs — both URIs
  (`…/xmldsig-more/rsa-ripemd160` and `…/07/05/…ecdsa-ripemd160`) are confirmed present in the
  frozen `enumerations/signature_algorithm.go` URI table, so these non-standard methods still
  verify (no false-negative). `enforced_resolver_fragment.go`: comment-only, the real guard lives
  in `internal/xmldsig`. Clean.

**Cross-referenced, not re-reported.**
- `xades_timestamp_scope_finder.go` — verified **1:1** against Java
  `XAdESTimestampScopeFinder.java` (same `isCollectionNotEmpty(includes)` gate, same
  `endsWithIgnoreCase(uri, documentName)` match, same `isContentTimestampedReference`). The
  re-implemented `FindTimestampScope` faithfully mirrors the Java base
  `EncapsulatedTimestampScopeFinder.findTimestampScope` (`isMessageImprintDataIntact()` guard).
  The **reachability** gap (the base `validateTimestamps` never dispatches to this override) is
  **U40's T40-SEC-001** — not re-reported here.
- `xades_signature_policy.go` and `xades_reference_validation.go` — the two `sync.Map`
  recovery registries are **U40's T40-PERF-001** (unbounded retention). Not re-reported.

### Open questions

1. **T41-INFO-001 → known-gaps.md.** The XSD structure-validation no-op is a genuine
   *undocumented* divergence (verified absent from `known-gaps.md`) and is trust-adjacent
   (TL structure conformity). **Question for the orchestrator:** add the `known-gaps.md`
   "XML structure validation" bullet now (documentation-only, no code), or fold it into the
   batch-06-trust follow-up that already surfaced it? Either way it should not stay
   documented only in code headers.
2. **T41-INFO-002 concurrency posture.** The unsynchronized lazy singleton is benign today
   (stateless struct, single-threaded port) but is real UB under concurrency. **Question:**
   is there any roadmap item that would add concurrent validation (e.g. a batch/signature
   fan-out), which would justify a `sync.Once` hardening pass across the three
   `…GetInstance`/singleton sites now rather than later?

### Tool log (run from `dss/`)

- `gofmt -l xades/` → **clean** (no output, exit 0).
- `go vet ./xades/...` → **clean** (exit 0, no diagnostics).
- `golangci-lint run --config=../.github/.golangci.yml ./xades/...` → **0 issues** (exit 0).
- `go test ./xades/... -count=1 -timeout 10m` → **pass**: `xades` 2.611s,
  `xades/definition` 1.491s; `xades/extension` and `xades/tsl` have no test files.
- `grep`/`find` sweeps: `panic(` enumeration across all in-scope files (12 sites, each mapped
  to its Java `Objects.requireNonNull` / `IllegalArgumentException` / `UnsupportedOperationException`
  origin); `sync.Map` registry sweep; `StructureValidator`/`ValidateAgainstXSD`/
  `StructureValidationMessages` blast-radius sweep; `CounterSignatureResolver`/
  `EnforcedResolverFragment` resolver-wiring sweep; RIPEMD160 signature-method URI sweep in
  `enumerations/`.
- Java cross-checks at `/Users/utain/Workspace/esig/dss`: `XAdESTimestampScopeFinder.java`
  (full `filterCoveredSignatureScopes` + `isContentTimestampedReference`), `EncapsulatedTimestampScopeFinder.java`
  (`findTimestampScope`), `XAdESAttributeIdentifier.java` (`getOrder`/`writeInt`),
  `EnvelopingSignatureBuilder.java` (`getUri().substring(1)`), `XAdESStructureValidatorFactory.java`
  (`getInstance`), and `DiagnosticDataBuilder.java` (structure-validation gating).
