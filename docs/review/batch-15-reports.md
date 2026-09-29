# Batch 15 — Reports

## detailedreport + diagnostic (unit U43)

**Scope.** detailedreport (16/3894L full-read) + diagnostic (78/17505L LIGHT pass)
**Date.** 2026-08-29
**Depth.**
- `detailedreport/` — ALL 16 non-test files read in full, including the whole `jaxb/`
  marshal path (`xml.go`, `jaxb_report.go`, `jaxb_content_model.go`, `jaxb_qualification.go`,
  `jaxb_adapters.go`, `jaxb_bbb.go`, `jaxb_common.go`, `jaxb_process.go`, `jaxb_model.go`,
  `doc.go`) and the verdict-navigation facade (`detailed_report.go`,
  `detailed_report_message_collector.go`).
- `diagnostic/` — LIGHT. Top files by size read: `diagnostic_data.go` (full),
  `certificate_wrapper.go` (L1–800), `signature_wrapper.go` (L1–250),
  `timestamp_wrapper.go` (L1–250), `abstract_token_proxy.go` (full),
  `abstract_signature_wrapper.go` (L40–110), `revocation_wrapper.go` (full),
  `certificate_revocation_wrapper.go` (full), `found_revocations_proxy.go` (L1–120),
  `credential_subject_claim_wrapper.go` (L1–120), plus `jaxb/xml.go` (full: simple types,
  `Link`/IDREF, `CertificateExtensionsWrapper`, `Marshal`/`jaxbCanonical` escape path).
  Remaining files (eaa_payload_proxy, eaa_wrapper, claim_wrapper, pdf_revision_wrapper,
  evidence_record_wrapper, orphan_* wrappers, the `jaxb/jaxb_*` model files) covered by
  grep-scan for verdict/panic/nil/sync patterns + Java-parity spot-checks. Marked
  "not fully verified (budget)" where noted.
- Java cross-checks run against `/Users/utain/Workspace/esig/dss`
  (`dss-detailed-report-jaxb`, `dss-diagnostic-jaxb`).

### Findings

No Critical or High findings. The package is a data-object/report layer: it carries and
navigates verdicts produced by the bbb/process layer (batch-11), and its own
verdict-derivation (`IsValidCertificate`, `IsSigningCertificateIdentified`, the level
`Is*Valid`/`IsThere*` predicates, `FinalConclusion`, `HighestConclusion`) is either
fail-closed (absent data ⇒ `false`/empty/nil) or a line-for-line Java port. No method flips
FAIL→VALID or drops a failure.

**T43-M-001 — Low — Concurrency — `detailedreport/detailed_report.go:966` (and
`diagnostic/found_revocations_proxy.go:25`) — unsynchronized lazy init of shared mutable state.**
Evidence:
```go
func (r *DetailedReport) MessageCollector() *MessageCollector {
    if r.messageCollector == nil {          // read-then-write, no lock
        r.messageCollector = newMessageCollector(r)
    }
    return r.messageCollector
}
```
and `FoundRevocationsProxy.ensureFoundRevocations()` does the same
(`if p.foundRevocations == nil { p.foundRevocations = &jaxb.XmlFoundRevocations{} }`).
Impact: a `DetailedReport`/`Data` object shared by concurrent readers that both take the
first-call path races on the collector/proxy field (undefined value, possible duplicate
allocation). Single-goroutine use (the norm — one report per validation job) is safe.
Parity note: Java's `DetailedReport.getMessageCollector()` is identically unsynchronized,
so this is NOT a behavioural divergence, only a latent Go data-race.
Recommendation: build the `MessageCollector` eagerly in `NewDetailedReport`, or guard the
two lazy-init sites with a `sync.Once`. Low because the report object is effectively
single-consumer in the validated call graph.

**T43-L-001 — Low — Correctness (parity) — `detailedreport/detailed_report.go:142` —
off-by-one bounds check.**
Evidence:
```go
bbbs := r.jaxbDetailedReport.BasicBuildingBlocks
if bbbs != nil && len(bbbs) >= index {     // should be len(bbbs) > index
    bbb := r.jaxbDetailedReport.BasicBuildingBlocks[index]
```
Impact: when `index == len(bbbs)` the check passes and `bbbs[index]` is the nil sentinel
that `append` never stores, so the method returns `""` (the `bbb != nil` guard) rather than
the correct id. No panic, no wrong verdict. This is a 1:1 copy of the Java
`getBasicBuildingBlocksSignatureId` (`bbbs.size() >= index`), so it is a faithful port of an
upstream defect, not a Go regression; no in-repo caller passes an out-of-range index.
Recommendation: tighten to `len(bbbs) > index` and record the deliberate-parity note if
kept as-is.

**T43-L-002 — Low — Robustness (parity) — `diagnostic/diagnostic_data.go:271` — nil-deref on
hostile input.**
Evidence:
```go
func (d *Data) TimestampSigningCertificateId(timestampId string) string {
    return d.timestampByIdNullSafe(timestampId).SigningCertificate().Id() // .Id() on nil
}
```
`SigningCertificate()` (abstract_token_proxy.go:163) returns `nil` when the timestamp has
no `SigningCertificate`; `.Id()` then panics. Java's `getTimestampSigningCertificateId`
throws the equivalent NPE, so this is parity, and the diagnostic model always carries a
signing certificate for a real timestamp (the validator uses the null-safe
`SigningCertificateId` path). Impact: panic (not a wrong verdict) only on a malformed
timestamp block. Recommendation: guard the nil `SigningCertificate()` before `.Id()`.

**T43-L-003 — Low — Robustness (parity) — `detailedreport/detailed_report.go:939` —
`HighestConclusion` dereferences a possibly-nil process block.**
Evidence:
```go
xmlSignature := r.XmlSignatureById(signatureId)
if xmlSignature.ValidationProcessArchivalData != nil { ... }
else if xmlSignature.ValidationProcessLongTermData != nil { ... }
return &xmlSignature.ValidationProcessBasicSignature...XmlConstraintsConclusionContent
```
The `else` branch takes the address of `ValidationProcessBasicSignature` with no nil check;
a signature carrying neither ArchivalData nor LongTermData (malformed) panics. Java's
`getHighestConclusion` makes the same unchecked call (NPE). `XmlSignatureById` guarantees a
non-nil `XmlSignature`, so the reachable case is a structurally-invalid signature. Parity,
no wrong verdict. Recommendation: nil-guard the `else` branch or document the precondition.

### Verdict-path verification notes

**1. Verdict derivation in `detailedreport/` (Special focus #1) — CLEAN / fail-closed.**
- `FinalConclusion`/`FinalIndication`/`FinalSubIndication` only *read* the top-level
  `Conclusion` of the highest-matching token (Signature → Timestamp → EvidenceRecord → EAA →
  BBB); they never synthesize or upgrade an indication.
- `BasicBuildingBlocksIndication`/`BasicValidationIndication`/`LongTermValidationIndication`/
  `*Qualification` return `""` (zero) when the block or `Conclusion` is absent — absent data
  can never read as "PASSED".
- `HighestConclusion` (detailed_report.go:939) returns the ArchivalData → LongTermData →
  BasicSignature conclusion, exactly mirroring Java `getHighestConclusion`; the
  MessageCollector's `collectSignatureValidation` only adds a lower-level block's messages
  when its `SubIndication` equals the highest block's — a faithful copy of the Java
  `enumerate` guard, and it cannot manufacture a PASS.
- `buildFromXmlCertificateApprovalStatus` (detailed_report.go:840) falls back to
  `CertificateForUnknown` when no known `ListType`/`ServiceTypeIdentifier`/`ServiceStatus`
  matches — i.e. an unrecognised LoTE status is NOT treated as a known/approved one
  (fail-closed), and `enumerations.ListTypeFromURI` returns `nil` for an unknown URI
  (enumerations/list_type.go:14).
- **No false-VALID derivation found.** The only verdict-flip risk would be a method that
  returned "PASSED" on missing data; none do.

**2. Diagnostic data construction (Special focus #2) — CORRECT.**
- `Base64Encoded`/`PublicKey`/`Octets`/`Ski`/`KeyIdentifier` all bind to `Base64Binary`
  (`[]byte`), whose `MarshalText` is `base64.StdEncoding` (diagnostic/jaxb/xml.go:65).
  `Binaries()` (`[]byte(*w.revocation.Base64Encoded)`) therefore returns the **raw DER**, not
  a base64 text — no double-encoding, no dropped bytes. This was the primary
  "wrong revocation bytes" risk and it is clean.
- `SigningCertificate()` returns the embedded `XmlSigningCertificate.Certificate` wrapper
  (abstract_token_proxy.go:163); `CertificateRefWrapper.DigestAlgoAndValue`/`IsDigestValueMatch`
  read the matcher as stored. The `SigningCertificate` *values* the validator consumes are the
  same `XmlSigningCertificate` node, so digest/issuer/serial comparison sees identical bytes.
- Revocation-status predicates (`CertificateStatus.IsGood`/`IsRevoked`/`IsKnown`,
  `certificate_revocation_wrapper.go:58`) and `IsValidCertificate`
  (`diagnostic_data.go:295`: `signatureValid && (trusted || revocationValid)`) are fail-closed:
  a missing/`nil` revocation or unknown status yields `false`, never a spurious "good". This is
  a verbatim port of Java `DiagnosticData.isValidCertificate` (lines 552–560).
- `LatestRevocationDataForCertificate` (diagnostic_data.go:743) picks the max
  `ProductionDate`, matching Java `getLatestRevocationDataForCertificate`.

**3. XML escape in generated report/diagnostic XML (Special focus #3) — NO GAP.**
Both packages marshal via `encoding/xml` (which escapes all unsafe chars in text and
attributes to `&#…;`), then `jaxbCanonical`/`writeCharData`/`writeAttrValue` rewrite a fixed,
whitelist of numeric references to the JAXB-RI spelling (detailedreport/jaxb/xml.go:140,
diagnostic/jaxb/xml.go:875). The rewrite tables only map `&#34;`/`&#39;`/`&#x9;`/`&#xA;`/`&#xD;`
to *valid* forms (`"`, `'`, tab, newline, `&#10;`/`&#13;`) and never introduce a bare `&`.
Because the only character that can start an entity is `&`, and `&` is always emitted by the
encoder as `&amp;` (which no rewrite rule touches), there is no path to emit an unterminated or
injected entity. `writeStartTag` only rewrites inside double-quoted attribute values;
`tagEnd` is safe because the encoder guarantees no raw `>` inside a quoted value.
Round-trip byte parity is pinned by the KAT/oracle tests (all passing).

**4. Concurrency (Special focus #4).** The only shared-mutable-state sites are the two
unsynchronized lazy inits in T43-M-001. No `sync.` primitives, no goroutines (`go func`) in
either package (grep-verified). The `baseTypeView` cache in `diagnostic/jaxb/xml.go` is
a fresh per-`Link` map (diagnostic/jaxb/xml.go:684, not shared) and exists to give a single `XmlTrustSourceList` view
pointer identity so the certificate-approval-status (TS 119 602) block's set-membership test
matches Java object identity — a correctness aid, not a race.

### Open questions
- **EAA diagnostic model (eaa_payload_proxy.go / eaa_wrapper.go / claim wrappers) not fully
  read (budget).** These are the TS 119 602 / eIDAS-2 claim-attestation data objects. Their
  validation is behind the `eaa` build tag and is "unproven" per
  `docs/compatibility/known-gaps.md` (no behavioural oracle corpus) — already a tracked gap,
  so not re-reported. A dedicated pass under the `eaa` tag is the natural follow-up.
- **`CertificateApprovalStatus` match on empty URIs.** In
  `detailedreport/detailed_report_message_collector.go:440`,
  `certificateApprovalStatusProcess` requires `ListType().URI() != ""` AND
  `ServiceTypeIdentifier().URI() != ""` before matching. If a caller-supplied
  `CertificateApprovalStatus` legitimately has an empty `ServiceTypeIdentifier` URI, no
  matching message is collected (the certificate is silently treated as not-having-that-status).
  This mirrors the Java guard, but worth confirming against a real LoTE status set that the
  "empty URI ⇒ no match" semantics is intended (fail-closed, but possibly over-restrictive).
- **`XmlSemantic`/`XmlMessage` `Key` required-ness.** `XmlSemantic.Key` and
  `XmlMessage.Key` are optional (`*string`/`string` with `omitempty`); a report author could
  emit a message with an empty `Key`. Not a security issue, but the i18n lookup would get `""`.

### Tool log (run from inside `dss/`)
```
gofmt -l detailedreport/ diagnostic/                                   → (empty, clean)
go vet ./detailedreport/... ./diagnostic/...                           → exit 0
golangci-lint run --config=../.github/.golangci.yml ./detailedreport/... ./diagnostic/...
                                                                        → 0 issues
go test ./detailedreport/... ./diagnostic/... -count=1 -timeout 15m
  ok  dss/detailedreport      1.165s
  ok  dss/detailedreport/jaxb 2.290s
  ok  dss/diagnostic          1.715s
  ok  dss/diagnostic/jaxb     1.252s
Java parity greps: DetailedReport.java (getBasicBuildingBlocksSignatureId,
  getEvidenceRecordIds, getSignatureIds, getHighestConclusion, getMessageCollector),
  DiagnosticData.java (isValidCertificate, getLatestRevocationDataForCertificate),
  RevocationWrapper/CertificateWrapper/TimestampWrapper (getBinaries/base64Encoded).
```

### Summary counts
Critical 0 · High 0 · Medium 1 (Low-severity concurrency parity) · Low 3 · Info 0.

---

## simplereport + simplecertificatereport + validationreport (unit U44)

**Scope.** simplereport (13/2293L) + simplecertificatereport (12/1991L) + validationreport (20/2832L)
**Date.** 2026-08-29
**Depth.** All 45 non-test files read in full — largest: `simplereport/simple_report.go` (716L),
`simplecertificatereport/simple_certificate_report.go` (801L), `validationreport/jaxb/jaxb_enums.go`
(522L); each read whole, none truncated. Java cross-checks at `/Users/utain/Workspace/esig/dss`:
`SimpleReport.java` (isValid/getIndication), `SimpleCertificateReport.java`
(getCertificate/getFirstCertificate/toCertificateApprovalStatus/getXmlCertificateApprovalStatus),
`specs-validation-report/.../UriBasedEnumParser.java` + `xsd/bindings.xml` (URI-enum parse/print),
`SimpleReport.xsd` (TrustAnchor binding). Empirical probes: `encoding/xml` `>`-escaping in attrs vs
char data, and encoder output for empty complex vs empty text elements.

### Findings

**T44-STD-001 — Low — Dead code / misleading comment — `simplereport/jaxb/jaxb_eaa_payload.go:144`.**
Evidence: `var eaaPayloadClaimElements = []string{ "Identifier", "Issuer", … "OtherClaim" }` — a ~130-line
string slice. Header comment: *"jaxb_content_model.go consults it to mark all of them as simpleContent"*.
Grep across the whole module (incl. tests) shows **zero references**; `jaxb_content_model.go` instead uses
the flat `complexElements` map (`carriesCharData` → `!complexElements[name]`), and the EAA claim names are
not even in that map (they fall through to the `kindText` default, which is the correct answer).
Impact: ~130 lines of dead data + a comment that asserts a wiring that does not exist; a maintainer
adding/removing an EAA claim will edit a table that changes nothing, or "fix" the content model to match it.
Recommendation: either delete the slice (its stated purpose is already served by the `kindText` default in
`carriesCharData`), or actually wire it into `carriesCharData` and fix the comment. Do not keep a false claim.

**T44-INFO-001 — Info — Deliberate parity (not a defect) — `simplereport/simple_report.go:50` (IsValid).**
`IsValid(tokenID)` returns `ind == TOTAL_PASSED || ind == PASSED`. This looks like it could over-accept a
signature, but it is exactly Java's `isValid()` (`TOTAL_PASSED.equals(v) || PASSED.equals(v)`). The `PASSED`
branch is the **timestamp** result space; a *signature*'s stored `Indication` is set by
`validation/executor/simple_report_builder.go:354` to `detailedReport.FinalIndication(sigId)` and
`validSignatureCount` is only incremented on `TOTAL_PASSED` (bare `PASSED` is a timestamp value, never a
signature `FinalIndication`). So a signature whose detailed report says `TOTAL_FAILED`/`INDETERMINATE` still
yields `IsValid == false`. No inversion, no always-true, no false-VALID.

**T44-INFO-002 — Info — Deliberate parity / effectively unreachable — `simplereport/simple_report.go:504`
(GetSignatureScopes).**
`default: panic("… not supported for an object of this class")` mirrors Java's
`UnsupportedOperationException`. The `case nil: return nil` handles the not-found path, and `tokenItemFor`
only ever produces the four concrete token types, so the `default` branch is unreachable for any value
`embeddedTokenByID` can return. Defensive parity guard, not a hostile-input panic.

### Verdict-path verification notes

**Verdict summarization (focus #1) — CLEAN.**
- `simplereport.IsValid` = `TOTAL_PASSED || PASSED` — matches Java 1:1; the builder side (out of scope,
  `validation/executor/simple_report_builder.go`) only counts a signature valid on `TOTAL_PASSED`. A
  signature reported `TOTAL_FAILED`/`INDETERMINATE` → `IsValid` false. **No false-VALID.**
- `simplereport.GetIndication`/`GetSubIndication` are nil-safe (return `""` when the token id doesn't
  resolve); `GetSignatureQualification` defaults to `NA`; `GetTimestampQualification` defaults to `""`.
  No accessor is hard-wired to a passing value.
- `simplecertificatereport` exposes the per-certificate verdict as a *value*
  (`GetCertificateIndication` → `cert.Indication.Indication()`), nil-safe (`""` when unresolved). There is
  no boolean `isValid()` to invert here. **No false-VALID.**

**Report composition (focus #2) — CLEAN (pass-through container).**
- `validationreport` is a *faithful container*, not a verdict aggregator. `ValidationReportType`
  (jaxb_root.go) holds `SignatureValidationReport []*SignatureValidationReportType` (one per signature),
  `SignatureValidationObjects`, `SignatureValidator`, `Signature`; `Marshal`/`Unmarshal` round-trip every
  field in schema order. The per-signature verdict lives in
  `SignatureValidationReportType.SignatureValidationStatus` (`ValidationStatusType.MainIndication`/
  `SubIndication`); the per-constraint verdict in `IndividualValidationConstraintReportType.ValidationStatus`.
- **Nothing is dropped:** `SignatureAttributesType` (35-member choice), `ValidationObjectRepresentationType`
  (4-member) and `SignersDocumentType` (3-member) all marshal via the shared `choiceItem` machinery —
  `encodeChoiceItems` iterates every recorded item, `decodeChoiceItems` reads until the enclosing end tag;
  order preserved (the marshal-parity KAT depends on this). No sub-report or choice member is silently omitted.
- **"A FAIL in any sub-report propagates to the overall verdict" is a *builder* concern
  (`dss/validation/executor`), not this package's.** This package cannot misrepresent a sub-report's
  verdict — it stores and re-emits the indication bytes verbatim. So there is no composition error *in
  scope*; the propagation logic was already reviewed in batch-10/batch-11.

**XML escape (focus #3) — CLEAN.**
- All three packages encode with `encoding/xml` (guarantees well-formed output: `<`, `>`, `&`, `"`, `'`,
  control chars escaped), then `jaxbCanonical` post-processes only to *relax* to the JAXB-RI spelling: char
  data `&#34;`→`"`, `&#39;`→`'`; attr values `&#34;`→`&quot;`. It never introduces a raw `<` or bare `&`.
- **Verified empirically** that `encoding/xml` escapes `>` to `&gt;` in *both* attribute values and char
  data (e.g. `XPath="//a[1>0]/b"` → `XPath="//a[1&gt;0]/b"`). This is what makes the `tagEnd` scan
  (`bytes.IndexByte(in, '>')`) safe: it is only ever called on a real tag, and a `>` inside an attribute
  value is always `&gt;`, so it can never be mistaken for the tag terminator. A `>` in an `XPath`/
  `WhichDocument` value cannot break the scan.
- `jaxbRootNamespaceLast` (simplereport/simplecertificatereport only; validationreport's root carries no
  own attributes so it's intentionally absent) locates the fixed `xmlns="<Namespace>"` declaration and moves
  it after the attributes; the namespace value is a compile-time constant, never attacker input, and any `"`
  in it would be escaped by the encoder, so the `IndexByte(…,'\"')` boundary is sound.
- **No escape gap, no injection vector.**

**Concurrency (focus #4) — CLEAN.**
- No `sync.Mutex`, no `atomic`, and **no package-level mutable global written after `init()`**. The
  `init()` blocks (`buildContentModel` in jaxb_content_model.go; `registerURIEnum`/`uriToEnum` in
  jaxb_enums.go) populate `map`/`slice` once at load and are never mutated thereafter.
- Marshal paths (`Marshal`, `MarshalToWriter`, the `MarshalXML` methods) read struct fields and write to
  **per-call** `bytes.Buffer`s. A report object is safe to read/marshal from multiple goroutines; it is not
  safe to *mutate* concurrently (same as the Java JAXB objects, effectively immutable after build). No
  shared-state hazard.

**Nil-safety (focus #5) — documented parity, no new deref.**
- `NewSimpleReport(wrapped)` / `NewSimpleCertificateReport(simpleReport)` do **not** nil-check the wrapped
  pointer; calling any accessor on a `nil` wrapped value panics (Go nil-field deref) — the **same shape as
  Java's NPE** (`simpleReport.getValidationTime()` on null `wrapped`). Both file headers explicitly document
  this as a *preserved upstream NPE risk* (PORTING.md "values verbatim"), e.g. `simple_certificate_report.go`
  header and `XmlPDFAInfo.Valid`'s comment. `SimpleReport.IsPDFACompliant` likewise preserves Java's unboxed
  `pdfaInfo.isValid()` NPE (documented).
- The **facade** entry points (`Facade.Marshal`/`Unmarshal`/`MarshalToWriter`/`Unmarshal`) **do** nil-check
  their arguments (`"JAXBObject is null"` / `"InputStream is null"` / `"OutputStream is null"`), matching
  Java's `AbstractJaxbFacade`.
- Net: no *new* nil-deref is introduced beyond Java parity; the preserved NPEs are documented, not silent.
  Treated as sanctioned per the "documented parity" rule (not a finding). See Open questions re: the
  hostile-input threat model.

### Open questions

1. **`eaaPayloadClaimElements` (T44-STD-001):** confirm it is safe to delete (it is referenced nowhere,
   incl. tests), or wire it into `carriesCharData` and correct the comment. Orchestrator to pick; I did not
   modify anything.
2. **Hostile-input threat model for nil-wrapped accessors:** the preserved `NewSimpleReport(nil)`→accessor
   panic is exact Java parity and the facade guards its own entry points, so I classify it as sanctioned. If
   the review's threat model includes *attacker-supplied report construction* (not just attacker-supplied
   report *bytes*, which go through the nil-checked facade), this would rise to "panic on hostile input" —
   flagging here only for that decision.

### Tool log (run from inside `dss/`)
```
gofmt -l simplereport/ simplecertificatereport/ validationreport/            → (empty, exit 0)
go vet ./simplereport/... ./simplecertificatereport/... ./validationreport/...
                                                                               → exit 0
golangci-lint run --config=../.github/.golangci.yml ./simplereport/... \
        ./simplecertificatereport/... ./validationreport/...                  → 0 issues
go test ./simplereport/... ./simplecertificatereport/... ./validationreport/... -count=1
  ok  dss/simplereport                    1.666s
  ok  dss/simplereport/jaxb               2.241s
  ok  dss/simplecertificatereport         0.629s
  ok  dss/simplecertificatereport/jaxb    1.180s
  ok  dss/validationreport                2.622s
  ok  dss/validationreport/jaxb           3.191s
  (6/6 ok, exit 0, no FAIL/panic)
Java parity: SimpleReport.java (isValid/getIndication), SimpleCertificateReport.java
  (getCertificate/getFirstCertificate/toCertificateApprovalStatus/getXmlCertificateApprovalStatus),
  UriBasedEnumParser.java + xsd/bindings.xml (URI-enum parse/print → Go ParseMainIndication/Print:
  miss → zero value, never error — matches), SimpleReport.xsd TrustAnchor (TrustServiceName is an
  unbounded xs:string list — matches simplereport/jaxb/jaxb_common.go).
Empirical probes (throwaway /tmp programs, deleted after): encoding/xml escapes >→&gt; in both attr and
  char data (validates the tagEnd scan); encoder writes <Parent></Parent> for an empty complex and
  <Msg></Msg> for an empty text element (validates the jaxbCanonical self-closing heuristic + stack
  push/pop balance on well-formed input).
```

### Summary counts
Critical 0 · High 0 · Medium 0 · Low 1 · Info 2.

## policy (unit U45)

**Scope.** `dss/policy/` (42 non-test files, 6,118 lines): the `ValidationPolicy`
object model (`EtsiValidationPolicy`), the JAXB port of `policy.xsd`
(`policy/jaxb/`), the constraint wrappers (Level/MultiValues/Value/IntValue/
Time/CertificateValues/Cryptographic), `RuleUtils`/`DateUtils`, and the
ETSI TS 119 322 crypto-suite catalogues + factories (`policy/crypto/json/`,
`policy/crypto/xml/`). `dss/validation/policy/` was already reviewed in U30
(batch-09) and is **not** re-derived here; `dss/model/policy/` was reviewed in
U25 (batch-08) and is cross-referenced, not re-derived.
**Date.** 2026-08-29
**Depth.** Budgeted. Full-read: every file ≤300 lines (all wrappers, `date_utils.go`,
`rule_utils.go`, `validation_policy_facade.go`, `validation_policy_xml_definer.go`,
`etsi_validation_policy_factory.go`, both crypto `*_catalogue.go`/`*_factory.go`/
`*_utils.go`/`json_object.go`/`xsd_datetime.go`/`rfc3339_date_utils.go`/`xml_types.go`/
`cryptographic_suite_xml_facade.go`, and all `policy/jaxb/*.go` model+adapter+content-model
files). Budgeted (>400L): `etsi_validation_policy.go` 2520L (lines 1–250 + 414–560 +
2340–2520 + `grep -n 'func '`: the remainder is the uniform `cc := p.certificateConstraints(…);
if cc != nil { return toXxx(cc.Field) }` accessor cluster verified by grep, plus the
context-routing helpers `signatureCryptographic`/`certificateConstraints`/
`basicSignatureConstraintsByContext` read in full) and `cryptographic_constraint_wrapper.go`
580L (full-read). Total in-scope code read ≈ 5,400 lines (within the 5,500 budget).

### Findings

**T45-STD-001 — Low — Parity/robustness — `policy/crypto/json/cryptographic_suite_json_catalogue.go:233` (`asDate`), `:251` (`asDateTime`) —**
```
t, err := rfc3339GetDate(dateString)
if err != nil {
	return nil          // Java: IllegalArgumentException propagates (only caught by buildAlgorithm's try/catch)
}
```
**Impact:** A malformed (present) date in a JSON crypto-suite is degraded to "absent"
(`nil`) rather than aborting the algorithm entry as Java does. Effect is bounded to
`Metadata.PolicyIssueDate`/`NextUpdate` and `Evaluation.ValidityStart/End`; it cannot
turn a weak algorithm into an acceptable one (algorithm *identity* and `Min` key size are
parsed independently and are unaffected), and the bundled `dss-crypto-suite.json` has
well-formed dates so the path is not exercised by the shipped suite. **Recommendation:**
None required (documented divergence, no security weakening); if a JSON-suite with a bad
validity-end is ever trusted, note it would be treated as non-expiring — confirm the
consumer then fails closed on "no acceptable evaluation window."

**T45-STD-002 — Low — Parity — `policy/crypto/json/cryptographic_suite_json_factory.go:54` (`IsSupported`) —**
```
obj, err := parseJSONObject(rc)
return err == nil && obj.asObject(jsonConstraintSecuritySuitabilityPolicy) != nil
```
**Impact:** `IsSupported` is a weaker gate than Java's full ETSI TS 119 322 JSON-Schema
validation (no JSON-Schema validator ships in the Go stdlib — see the DEFERRED notes in
`cryptographic_suite_json_utils.go`/`cryptographic_suite_xml_utils.go`). A
structurally-shaped but semantically-malformed document that carries a
`SecuritySuitabilityPolicy` key returns `true` where Java returns `false`. This is a
false-**supported** (probe), not a false-accept of crypto: the subsequent
`LoadCryptographicSuiteFromReader` + consumer still enforce algorithm identity and
key-size, and the gap is explicitly documented. **Recommendation:** None required; tracked
as a deferred-dependency gap consistent with the sibling XSD/JSON-Schema stubs.

**T45-STD-003 — Info — Parity (latent concurrency) — `policy/cryptographic_constraint_wrapper.go:78` (`AcceptableDigestAlgorithms`), `:120` (`AcceptableSignatureAlgorithms`) —**
```
if w.acceptableDigestAlgorithms == nil {
	w.acceptableDigestAlgorithms = acceptable   // lazy init, no lock
}
```
**Impact:** The lazy-caching of the computed algorithm maps is not synchronised. A data
race is only reachable if a single `*CryptographicConstraintWrapper` is read from two
goroutines concurrently *and* both hit the nil-cache window. This mirrors Java's identical
unsynchronised `if (acceptableDigestAlgorithms == null)` lazy init 1:1, so it is a parity
property, not a divergence; the validation executor drives checks sequentially per policy
object, so it is not confirmed reachable in this port. **Recommendation:** None required;
note only — if the executor is ever made to evaluate one shared policy from multiple
goroutines, add a `sync.Once` (a deliberate divergence needing a `DIVERGENCE` comment).

No Critical, High or Medium findings. No false "reliable"/"acceptable" verdict, no
condition-composition inversion, no off-by-one expiry boundary, no XXE, no integer
overflow, no panic-on-hostile-policy beyond the documented Java-exception panics.

### Verdict-path verification notes

**NumericValueRule float64 vs Java — RESOLVES U25 OPEN QUESTION.** The `dss/policy`
implementation is float64-backed (`int_value_constraint_wrapper.go:28` `Value() float64`),
but it is **lossless for every reachable value**, so it does **not** diverge from Java's
exact integer semantics and is **not** a finding:

- The wire type is `xsd:int`. `policy/jaxb/jaxb_common.go:41` models it as `Value *int`
  (Go 32-bit int) and `policy.xsd` declares `<xsd:attribute type="xsd:int" name="value"/>`.
  The value domain is therefore integers in `[-2^31, 2^31-1]`.
- `float64` represents every integer with magnitude ≤ 2^53 exactly; the full `xsd:int`
  range is a strict subset, so `float64(*w.constraint.Value)` is exact for all inputs a
  policy document can carry. No 2047.999…-style mis-rounding is possible.
- Every consumer compares with `float64(actual) >= constraint.Value()`:
  `CertificateMinQcEuRetentionPeriodCheck.Process`
  (`validation/process/bbb/xcv/certificate_min_qc_eu_retention_period_check.go`) and
  `CertificateMinQcTransactionLimitCheck.Process` (`…_transaction_limit_check.go`). Java
  does `qcEuRetentionPeriod >= constraint.getValue().intValue()` (same file, Java) and
  `amount * Math.pow(10, exponent) >= constraint.getValue().intValue()`. Go's
  `float64(amount) * math.Pow(10, float64(exponent))` is the same IEEE-754 double
  arithmetic as Java's `double * Math.pow`, and `float64(int) >= float64(int)` is
  equivalent to `int >= int` for these magnitudes. Boundary at 2048/4096 and at the
  threshold value is therefore identical to Java — no below-threshold false-accept.
- Nil-safety is correct and fail-closed: `toNumericValueRule`
  (`etsi_validation_policy.go:2492`) returns a true nil interface for a nil constraint
  (avoids the typed-nil-in-interface trap), and `ChainItemBase.Execute`
  (`validation/process/chain_item.go:184`) skips the check entirely when
  `constraint == nil` (level ""), so `Process()` (which dereferences `constraint.Value()`)
  is never reached with a nil constraint. A QC statement that is *absent* fails the check
  (`return false`) exactly as Java. **Conclusion: U25 OQ resolved — float64-backed but
  exact; no parity divergence.**

**Algorithm-reliability fail-closed.** `dss/policy` produces the
`AcceptableDigestAlgorithms`/`AcceptableSignatureAlgorithms` maps
(`cryptographic_constraint_wrapper.go`). Both are built *only* from algorithms
explicitly present in the policy: Step 1 seeds the map from `AcceptableDigestAlgo.Algos` /
`AcceptableEncryptionAlgo.Algos`, and unknown names are dropped
(`cryptographicConstraintWrapperToDigestAlgorithm`/`…ToEncryptionAlgorithm` return
`ok=false` → `continue`). An algorithm absent from the policy is simply absent from the
map, so the downstream reliability check (in `validation/policy`/`model/policy`, U30/U25)
sees no entry and fails closed — there is no "reliable-by-default" path and no false
"reliable" verdict originating here. `SignatureAlgorithmGetAlgorithm` returns `""` for an
unmatched encryption/digest pair and is skipped, again excluding (not including) the
combination. `MiniPublicKeySize` floors are applied only to already-acceptable algorithms
and `floorEvaluations` *raises* a parameter `Min` (`*minSize > *parameter.Min()`), so it can
never lower a required key size — fail-safe.

**Condition composition (AND/OR).** `dss/policy` contains **no** AND/OR condition-group
combinator — it only exposes per-check `Level`/`MultiValues`/`NumericValue`/`Duration`
rules. Boolean composition of checks lives in `validation/process` (the chain-item
executor, `chain_item.go`) and in `validation/policy` (U30), both already reviewed. The
wrapper logic here is pure projection (jAXB node → model rule) with no inversion: each
accessor is `if cc != nil { return toXxx(cc.Field) }; return nil`. No Critical
"true-when-all-false" path exists in this package.

**Expiry boundaries (inclusive/exclusive).** `cryptographicConstraintWrapperBuildEvaluation`
ports Java `buildEvaluation` byte-for-byte: `validityEnd = getDate(algo, layout)`; then
`if validityEnd == nil || (forcedValidityEnd != nil && validityEnd.After(*forcedValidityEnd)) {
validityEnd = forcedValidityEnd }` — an exact copy of Java's `validityEnd.after(forcedValidityEnd)`
(strictly-later replacement, i.e. the later of the two is *not* chosen when equal — the
forced value wins on tie, identical to Java). `cryptographicConstraintWrapperAlgorithmExpirationDate`
ports `getAlgorithmExpirationDate`: returns `nil` (unbounded) if *any* evaluation has no
`validityEnd`, otherwise the latest — identical to Java's `before()` loop. Date parsing is
UTC-forced on both sides (Java `sdf.setTimeZone(UTC)`; Go `t.UTC()`), so the inclusive/exclusive
semantics of a `validityEnd` boundary match. The only asymmetry is the documented
malformed-date → nil degradation (T45-STD-001), which does not change the boundary operator.

**Loader safety (size/XXE/entity expansion).** `policy/jaxb.Unmarshal` and
`policy/jaxb.Marshal` use the Go stdlib `encoding/xml` decoder/encoder. Go's `encoding/xml`
does **not** resolve DTDs or external/internal general entities and has no entity-expansion
(billion-laughs) path, so the policy XML loader is XXE-safe by construction (this is the
batch-02 `xmldom`/`xmlc14n` chain's domain; nothing here re-introduces it). The JSON loader
(`cryptojson.parseJSONObject` → `encoding/json`) and the TS 119 322 XML loader
(`cryptoxml`, `encoding/xml.Decode`) likewise have no entity-expansion vector. The one
resource note: `ValidationPolicyFacade.Unmarshal` does `io.ReadAll(r)` and
`UnmarshalFile` does `os.ReadFile`, i.e. the whole policy document is held in memory — a
parity property with Java's in-memory JAXB unmarshal and immaterial for policy documents
(tens of KB, admin-supplied, not network-streamed). No size/depth cap is needed given the
bounded schema and the absence of any recursion-on-attacker-input path; `jaxbCanonical`'s
`tagEnd`/`writeStartTag` scans are linear over already-parsed encoder output (no
attacker-controlled nesting beyond what `encoding/xml` itself accepted).

### Open questions

1. **Shared-policy concurrency (T45-STD-003).** Confirm the `validation` executor never
   evaluates one shared `*policy.EtsiValidationPolicy`/`*CryptographicConstraintWrapper`
   from multiple goroutines; if it does (e.g. a future parallel-signature executor), the
   unsynchronised lazy cache in `cryptographic_constraint_wrapper.go` becomes a live data
   race and needs a `sync.Once` (deliberate divergence).
2. **JSON-suite validity-end trust (T45-STD-001).** If an external (non-bundled)
   TS 119 322 JSON suite with a malformed `validity.end` is ever trusted, confirm the
   consumer treats the resulting "no expiry" evaluation as fail-closed rather than
   accepting an otherwise-weak algorithm.
3. **Deferred schema validation (T45-STD-002 + the three `ValidateAgainstSchema`/`Schema()`
   stubs).** These are the documented no-Go-stdlib-validator gaps; re-rate when a
   tech-lead-signed XSD/JSON-Schema dependency is added, per PORTING.md.

### Tool log

```
cd /Users/utain/Workspace/dss-go/dss
gofmt -l policy/                     -> (empty)   gofmt-exit=0
go vet ./policy/...                  -> (clean)   vet-exit=0
golangci-lint run --config=../.github/.golangci.yml ./policy/...  -> 0 issues.
go test ./policy/... -count=1 -timeout 10m
  ok  dss/policy            0.645s
  ok  dss/policy/crypto/json  2.157s
  ok  dss/policy/crypto/xml   1.118s
  ok  dss/policy/jaxb         1.671s
  (4/4 ok, no FAIL/panic)
Java parity cross-checks (/Users/utain/Workspace/esig/dss, DSS 6.5.RC1):
  IntValueConstraint.java + policy.xsd (xsd:int value attribute → Go *int → float64 exact),
  IntValueConstraintWrapper.java (Number→float64),
  CertificateMinQcEuRetentionPeriodCheck.java + CertificateMinQcTransactionLimitCheck.java
  (>= comparison, intValue() vs float64, Math.pow(10,exp) parity, nil→false),
  CryptographicConstraintWrapper.java (getAcceptable{Digest,Signature}Algorithms lazy
  cache + Step1/2/3 build, getAlgorithmExpirationDate, getFloorEvaluations, buildEvaluation
  forced-validityEnd tie, getUsedDateFormat setTimeZone(UTC), getDate ParseException→null),
  EtsiValidationPolicy.java (initializeCryptographicSuite field-fill, getCertificateConstraints
  subContext switch, getSignatureCryptographic),
  RuleUtils.java (convertDuration null→Long.MAX_VALUE; TimeUnit ratios) +
  TimeConstraintWrapper.java (null constraint→0),
  DateUtils.java (setLenient(false)+IllegalArgumentException, matching date_utils.go),
  ChainItem.java (getLevel null→check-skip, matching chain_item.go Execute),
  CryptoSuite Json/Xml Catalogue + Factory + Utils (buildJSONAlgorithm try/catch-skip vs Go
  no-panic accessors; IsSupported schema vs structural; deferred ValidateAgainstSchema).
```

