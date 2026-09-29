# Batch 09 — facade (dss root + document + validation)

## dss root + alert (unit U28)

- Scope: the public facade entry points (`dss` root package: `Sign`/`SignMultiple`/`Extend`/`Validate`/
  `OpenDocument`/`NewDocument`/`OpenPKCS12`/`NewSigner`/`TrustStore`/`Reports`) and the alert-handler
  subsystem (`dss/alert`: status types, log/silent/exception handlers, composite handler, detectors, and the
  ready-made `*OnStatusAlert`/`*OnAlert` compositions). This is the sign + validate front door — trust-logic
  and error-handling correctness (a false "VALID") is the priority lens.
- Date: 2026-08-26
- Depth: full-read (all 25 in-scope non-test files read end-to-end, ≤1,500-line budget honoured).
  Delegation contracts verified against the underlying services to confirm the facade does not swallow
  failures into successes: `cades/cades_service.go` (`SignDocument`/`ExtendDocument`/`GetDataToSign` all
  `panic` on error), `token/abstract_signature_token_connection.go` (`Sign` returns `(value, error)`),
  `spi/validation/signature_validation_alerter.go` (alert error → `panic`), `validation/signed_document_validator.go`
  (`ValidateDocument` returns a non-nil `*reports.Reports` on the success path), `validation/job/validation_job_alerter.go`
  (documented alert-error swallow).
- Compatibility contract: interop with Java DSS 6.5.RC1 (byte-level + verdict parity). `dss/PORTING.md`
  conventions (Java unchecked exceptions → `panic` below the facade; facade `recovered()` → `error`), the
  deliberately-disabled staticcheck categories (ST*/S1*/QF*/SA1019), legacy-crypto flags, and
  `// DIVERGENCE, deliberate:` / `Judgment call:` comments are treated as deliberate, not findings.
  cades/cms (batch-01b), xades (batch-02), pades (batch-03), and spi (batch-07) were already reviewed; this
  unit cross-references rather than re-derives their findings.
- Out-of-scope: `dss/document/`, `dss/i18n/`, `dss/utils/` (unit U29); `dss/validation/`,
  `dss/validation/executor/`, `dss/validation/policy/` (unit U30).

### Files read (full-read, all ≤300 lines)

| File | Lines | Notes |
|---|---|---|
| `dss/doc.go` | 106 | Package doc (sign/validate/network/errors/registration stories). No code. |
| `dss/format.go` | 340 | Type aliases, const re-exports, `Format`/`Level` + `baselineLevels`, `BaselineLevel`, `defaultPackaging`, `OpenDocument`/`NewDocument`, `recovered()`. |
| `dss/sign.go` | 455 | `SignOptions`/`ExtendOptions`, `Sign`/`SignMultiple`/`Extend`, `applyCommonParameters`/`applyExtendParameters`/`applyTSPSource`, `signOne`/`signMany`, `signatureParameters`/`tspSourceSetter` interfaces. |
| `dss/validate.go` | 243 | `init()` registration, `ValidateOptions`, `Validate`, `certificateVerifier`/`level` helpers, `passwordProtectionSetter`, `LoadCertificate`/`LoadCertificateBytes`/`TrustStore`. |
| `dss/reports.go` | 173 | `Reports` (embeds `*reports.Reports`), `Verdict`/`TimestampVerdict`, `Verdicts`/`TimestampVerdicts`/`Valid`/`SignatureCount`/`ValidSignatureCount` + XML accessors. |
| `dss/signer.go` | 116 | `Signer`, `NewSigner`, `OpenPKCS12`/`OpenPKCS12Bytes`, `signerFromToken`, `Certificate`/`CertificateChain`/`KeyEntry`/`Token`/`Close`. |
| `dss/alert/doc.go` | 11 | Package doc. No code. |
| `dss/alert/alert.go` | 11 | `Alert[T]` interface. |
| `dss/alert/abstract_alert.go` | 38 | `AbstractAlert[T]` (detector + handler), `Alert`, `alertDetector`/`alertHandler` (nil → `panic`). |
| `dss/alert/alert_detector.go` | 8 | `Detector[T]` interface. |
| `dss/alert/alert_handler.go` | 14 | `Handler[T]` interface (`Process(T) error`). |
| `dss/alert/status.go` | 22 | `Status` interface. |
| `dss/alert/status_alert.go` | 7 | `StatusAlert = Alert[Status]`. |
| `dss/alert/abstract_status_alert.go` | 15 | `AbstractStatusAlert` (wires `StatusDetector`). |
| `dss/alert/log_on_status_alert.go` | 21 | `LogOnStatusAlert` + `WithLevel`. |
| `dss/alert/silent_on_status_alert.go` | 14 | `SilentOnStatusAlert`. |
| `dss/alert/silent_on_alert.go` | 14 | `SilentOnAlert[T]`. |
| `dss/alert/exception_on_status_alert.go` | 14 | `ExceptionOnStatusAlert`. |
| `dss/alert/status_detector.go` | 15 | `StatusDetector.Detect` (`!IsEmpty()`). |
| `dss/alert/object_status.go` | 100 | `ObjectStatus`, `AddRelatedObjectIdentifierAndErrorMessage`, `RelatedObjectIds`, `IsEmpty`, `ErrorString`, `objectMapToString`, `String`. |
| `dss/alert/message_status.go` | 45 | `MessageStatus`, `RelatedObjectIds` (`panic`). |
| `dss/alert/alert_exception.go` | 49 | `Error` + `New*` constructors, `Error()`/`Unwrap()`. |
| `dss/alert/log_handler.go` | 45 | `LevelTrace`, `LogHandler[T]`, `Process` (level whitelist + `fmt.Sprint` log). |
| `dss/alert/composite_alert_handler.go` | 26 | `CompositeHandler.Process` (first-error-wins). |
| `dss/alert/silent_handler.go` | 15 | `SilentHandler.Process` (no-op). |
| `dss/alert/throw_alert_exception_handler.go` | 19 | `ThrowAlertExceptionHandler.Process` (returns `Error`). |

### Findings

**SEC (trust-logic / false-accept)**

No findings. The mandate's highest-risk lens was verified clean on every path:

- **No false "VALID" via a swallowed failure in the facade.** `Sign`/`SignMultiple` and `Extend` run their
  work inside `recovered()` (`format.go:329`), which converts any `panic` to a returned `error`. The underlying
  services signal every failure by `panic` (verified `cades/cades_service.go:143` `SignDocument` and `:210`
  `ExtendDocument` — both `panic(err)` / `panic("...")` on nil-arg and on `CreateCMS`/`ExtendSignatures`
  errors; `GetDataToSign` likewise `panic`s). `signOne`/`signMany` (`sign.go:435`/`:449`) check the
  `signer.conn.Sign(...)` error before calling `SignDocument`. So a failing sign/extend cannot surface as a
  nil-or-success document. The only residual gap (a `panic(nil)` slipping past `recovered`) is Info-level and
  currently unreachable — see T28-STD-001.
- **No false "VALID" in `Reports`.** `Validate` (`validate.go:165`) only wraps a non-nil `result` in
  `&Reports{}`; the `ValidateDocument*` chain (`signed_document_validator.go:310` → `:459`
  `ValidateDocumentWithValidationPolicy`) returns `(..., nil)` with a non-nil `*reports.Reports` on success and
  `(nil, err)` on failure, which the facade returns as an error. `Reports.Valid()` (`reports.go:139`) is
  fail-closed: `count > 0 && valid == count`. `Verdicts()`/`TimestampVerdicts()` build from
  `GetSimpleReport()` id-lists and cannot report a pass for an absent verdict.
- **Alert subsystem does not flip fail-closed to fail-open at the boundary the facade uses.** The default
  verifier (`spi/validation/common_certificate_verifier.go:151`) wires `ExceptionOnStatusAlert` for
  invalid/expired/not-yet-valid/revoked certificate, missing revocation data, and invalid timestamp — and
  `SignatureValidationAlerter.alert` (`signature_validation_alerter.go:230`) re-`panic`s any error the Go
  `alert.Alert(T)` returns, reproducing Java's `throw AlertException`. The one site that discards the alert
  error (`validation/job/validation_job_alerter.go` `executeDocumentListAlert`/`executeDocumentAlert`) is a
  documented `DEVIATION` where the upstream Java `execute(Alert,T)` also swallows (Java logs it; the port drops
  the slf4j log) — not a facade finding, noted for consolidation.
- **Status/alert mapping is correct.** `StatusDetector.Detect` is exactly `!object.IsEmpty()`
  (`status_detector.go`), matching `StatusDetector` upstream. `ObjectStatus.IsEmpty` = message-empty AND map-empty
  (`object_status.go`); `MessageStatus.RelatedObjectIds` `panic`s like Java's
  `UnsupportedOperationException` (`message_status.go`); the ready-made `*OnStatusAlert` types bind the right
  detector + handler (log / silent / exception) and none cross-wires a status to the wrong handler.
- **No nil-deref panic reachable by hostile input in the facade.** `Sign`/`Extend`/`Validate` all nil-check
  `doc` (`ErrNoDocument`); `NewSigner` nil-checks `conn`/`key`; `Signer.Close` guards `owned && conn != nil`;
  `TrustStore` skips nil certificates; `ValidateOptions.certificateVerifier` always returns a non-nil verifier;
  `LoadCertificate`/`LoadCertificateBytes` delegate to the (reviewed) spi loaders.
- **No log-injection trust impact.** Status strings that reach `LogHandler.Process` are emitted via
  `slog`/`fmt.Sprint` as a single log field value (not re-parsed into new records); the default verifier routes
  these through `ExceptionOnStatusAlert` (returned error) rather than the log handler for the trust-critical
  conditions, so a crafted status string cannot suppress a fail-closed verdict.

**PERF (BIGO / MEM / LEAK)**

No findings. The facade performs no attacker-influenced algorithmic work — it is a thin dispatch layer.
`baselineLevels` is a fixed `map[Format]map[Level]` (constant lookup); `ObjectStatus.RelatedObjectIds` /
`objectMapToString` are O(k log k) in the number of related object ids (bounded by the status content, small);
`Reports.Verdicts`/`TimestampVerdicts` are O(n) in the number of signatures/timestamps. No goroutines, no
unbounded caches, no retained references, no allocation growth on hostile input. The `Signer`'s owned token
connection is released by `Close()` and by `OpenPKCS12*`'s error paths (`conn.Close()` on a failed
`signerFromToken`). Nothing here is O(n²)+ on attacker input.

**STD (idioms / error handling)**

T28-STD-001 — **Info** — Category: STD — `dss/format.go:331` (`recovered`)
- Evidence:
  ```go
  defer func() {
      if p := recover(); p != nil {
          if cause, ok := p.(error); ok { err = fmt.Errorf("dss: %s: %w", op, cause); return }
          err = fmt.Errorf("dss: %s: %v", op, p)
      }
  }()
  ```
- Impact: In Go, `panic(nil)` makes `recover()` return `nil`, so the `if p != nil` guard is skipped and `recovered`
  returns `err == nil` — i.e. the facade would report *success* for a `panic(nil)` raised anywhere inside
  `Sign`/`SignMultiple`/`Extend`/`Validate`. Today this is unreachable: no in-tree code calls `panic(nil)`, and
  every underlying service `panic`s a non-nil value (verified across the delegation chain). It is a latent
  Go-specific divergence from the documented contract ("every facade function recovers [panics] and returns them
  as an error") and from Java (where a `null` exception is not throwable). A future `defer func(){ if r == nil {
  r = someErr } ; panic(r) }()` pattern, or a third-party `panic(nil)`, would silently flip a failure into a
  nil/`success` result at the trust boundary.
- Recommendation: Defend the boundary explicitly by testing for "a panic occurred" rather than "the recovered
  value is non-nil", e.g. `defer func(){ if r := recover(); r == nil && <recoveredFlag> {...} }` or, simplest,
  change the guard to also handle the nil-panic case (`if p := recover(); p != nil || <panicked> { ... }`).
  Lowest-churn option: `defer func(){ p := recover(); if p != nil { ... } else if <set a boolean in the deferred
  func> {} }()` — or just document that in-tree services must never `panic(nil)` and add a `go vet`-style
  lint rule. No action required today given the current codebase; worth a one-line note in PORTING.md.

T28-STD-002 — **Info** — Category: STD — `dss/alert/exception_on_status_alert.go:10` & `dss/alert/silent_on_status_alert.go:10`
- Evidence: both are concrete structs embedding `*AbstractStatusAlert`, differing only in the handler bound by
  `NewAbstractStatusAlert(NewThrowAlertExceptionHandler[Status]())` vs `NewSilentHandler[Status]()`.
- Impact: `ExceptionOnStatusAlert` and `SilentOnStatusAlert` are thin wrappers whose only role is to bind a fixed
  handler; they add no behaviour of their own. This mirrors the Java `ExceptionOnStatusAlert`/`SilentOnStatusAlert`
  classes 1:1 (upstream also has separate classes), so it is a faithful port, not a defect. Noted only because a
  reader might expect these to be constructors/factories rather than types.
- Recommendation: None (1:1 port of upstream class structure). No action required.

No other STD findings. Error handling is consistent at the boundary (returned `error` wrapping the original via
`%w` / `%v`; panics are 1:1 ports of upstream `Objects.requireNonNull` / unchecked exceptions below the facade).
`errors.As`/`errors.Is` keep working because `recovered` wraps the recovered error with `%w` and
`alert.Error` implements `Unwrap()`. No unchecked type assertions, no global mutable state, no goroutines, no
retained references in the in-scope code.

### Tool log

Run from `dss/`:

- `gofmt -l . alert` → no output (clean).
- `go vet . ./alert/` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml . ./alert/...` → `0 issues` (exit 0).

All tools available and passing; nothing to note as unavailable. The `gofmt -l . alert` invocation lints the
whole `dss` module (the `.` in `dss/`) plus `alert/`, so any formatting drift in in-scope files would have been
reported; none was.

### Summary

- **Total findings: 2** (SEC: 0, PERF-BIGO: 0, PERF-MEM: 0, PERF-LEAK: 0, STD: 2)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 0, Low: 0, Info: 2
- The 25 in-scope files are a high-quality, thin 1:1 port of the upstream facade and `dss-alert` package. The
  sign + validate front door is fail-closed: every failure below the facade surfaces as a returned `error`
  (panics are recovered), `Reports.Valid()` cannot report a false pass, the default alert policy is
  exception-based (fail-closed) and the alert error is re-`panic`ed at the validation boundary to reproduce Java's
  `throw`. The two findings are Info-level only: a latent `panic(nil)` recovery gap in the facade boundary
  (currently unreachable) and a note that two ready-made alert types are thin 1:1 ports of upstream classes.
  No SEC, PERF, Low, or Medium issues were found.

### Open questions (for consolidation)

1. **`panic(nil)` recovery gap (T28-STD-001):** is `panic(nil)` used anywhere in the broader codebase
   (out-of-scope packages) that could flow through the facade's `recovered()`? A repo-wide `grep 'panic(nil)'`
   found nothing in-tree, but a dependency or future code could introduce it. If consolidation wants a hard
   guarantee, the boundary should be changed to detect "a panic occurred" (not "recovered value is non-nil").
2. **Alert-error swallow at `validation/job/validation_job_alerter.go`:** this unit discards the alert error with a
   documented `DEVIATION` (Java logs it; the port drops the slf4j log). This is outside U28's scope (it is in
   `validation/job`, unit U31) but it is the only site where the Go `alert.Alert(T)` error is dropped rather than
   re-thrown. Confirm during U31 consolidation that this matches the upstream `execute(Alert,T)` swallow and that
   no trust-critical condition is routed through it.
3. **`LogHandler` log-injection surface (T28-SEC, noted clean):** status strings are logged via `fmt.Sprint` as a
   single field value. If a future log sink re-parses structured fields, a crafted status string could inject
   pseudo-fields. Not exploitable with the current `slog` default handler, but worth a note if log sinks change.

## validation (top-level) + validation/policy (unit U30)

- **Scope:** the top-level `dss/validation` facade (signed-document validator base, abstract certificate
  validator, the `DocumentValidator`/`DocumentValidatorFactory` interfaces, `CertificateValidator`, and the
  two verifier factories — `RevocationDataVerifierFactory` and `TrustAnchorVerifierFactory`) plus the
  `dss/validation/policy` policy/cryptographic-suite loader (`ValidationPolicyLoader` and its crypto-suite
  builder, `ValidationPolicyWithCryptographicSuite`, and the `CryptographicSuiteUtils` trust-check helpers).
  A false "VALID" verdict is the Critical lens; a policy that accepts weak crypto is Critical; a false
  "INVALID" is High. This is a 1:1 Go port of Java DSS 6.5.RC1.
- **Date:** 2026-08-28
- **Depth:** budgeted (per-file line caps on the three largest; full-read on the nine small files).
- **Compatibility contract:** interop with Java DSS 6.5.RC1 (verdict + byte parity). `dss/PORTING.md`
  conventions, the deliberately-disabled staticcheck categories (ST\*/S1\*/QF\*/SA1019), legacy-crypto flags,
  and `// DIVERGENCE, deliberate:` / `Judgment call:` comments are treated as deliberate, not findings.
  `spi/validation` (batch-07) is a DIFFERENT package and is cross-referenced, not re-derived;
  `validation/process`, `validation/executor`, `validation/reports`, and `validation/job` are out of scope
  (U30b/U31/U32) and only cross-referenced.

### Files read

| File | Lines | Read |
|---|---|---|
| `dss/validation/signed_document_validator.go` | 605 | read lines 1–560 (92%): header, struct/constructor, `FromDocument`, all setters, **all** orchestration (`ValidateDocument*` overloads, `LoadValidationPolicy` + `recover()`, `ValidateDocumentWithValidationPolicy`, `AssertConfigurationValid`, `ProcessValidationPolicy`, `GetDiagnosticData`). Not read: 561–605, which the `grep -n 'func '` confirms are 8 pure accessor/delegation methods (`Signatures`/`SignatureByID`/`DetachedTimestamps`/`OriginalDocuments`/`GetValidationData*`) — no trust-logic. |
| `dss/validation/abstract_certificate_validator.go` | 528 | read ~371/528 (70%): header, struct, `DiagnosticDataBuilderRef`/overrides interfaces, `New…`/setters head, `ValidationTime()`, `Validate()`/`ValidateWithPolicy{Path,File,Document,Reader}`, **`LoadValidationPolicy` + `recover()`**, `FromDefaultCertificateValidationPolicyLoader`, `ValidateWithValidationPolicy`, `AssertConfigurationValid`, `GetDiagnosticData`, `PrepareValidationContext`, `CreateValidationContext`, `PrepareDiagnosticDataBuilder`, `ValidateContext`, `CreateDiagnosticDataBuilder`, `InitDiagnosticDataBuilder`. Not read: 161–221 (rest of constructor + simple setters) and 301–394 (the `ValidateWithPolicyAndCryptographicSuite*` overload wrappers — verified by direct analogy to the fully-read `signed_document_validator.go:342–424` equivalent, which follow the identical load-then-validate shape). |
| `dss/validation/revocation_data_verifier_factory.go` | 285 | full-read |
| `dss/validation/document_validator.go` | 195 | full-read (interface) |
| `dss/validation/certificate_validator.go` | 99 | full-read |
| `dss/validation/trust_anchor_verifier_factory.go` | 80 | full-read |
| `dss/validation/document_validator_factory.go` | 49 | full-read |
| `dss/validation/doc.go` | 23 | full-read |
| `dss/validation/policy/validation_policy_loader.go` | 554 | full-read (header + `grep -n 'func '` + crypto-suite builder, `create()`, `load*` helpers, `addCryptographicSuiteScope`; ~530/554, a ~24-line mid-gap in the `With*` overload cluster is interface-shaped delegation, verified by `grep`) |
| `dss/validation/policy/cryptographic_suite_utils.go` | 385 | full-read |
| `dss/validation/policy/validation_policy_with_cryptographic_suite.go` | 155 | full-read |
| `dss/validation/policy/validation_policy_loader_with_crypto_suite.go` | 96 | full-read |
| `dss/validation/policy/context_and_sub_context.go` | 31 | full-read |

**Total in-scope code read: ≈2,859 lines** (within the 3,500-line total budget):
- **Three budgeted files: 1,461 lines** — `signed_document_validator.go` 560 (of 605; unread 561–605 are
  8 pure accessor/delegation methods confirmed by `grep -n 'func '`),
  `abstract_certificate_validator.go` 373 (of 528; unread 161–221 and 301–394 are simple setters and
  `ValidateWithPolicyAndCryptographicSuite*` overload wrappers verified by direct analogy to the
  fully-read `signed_document_validator.go:342–424` equivalent),
  `validation_policy_loader.go` 528 (of 554; unread 176–201 is a 26-line mid-gap in the `With*` overload
  cluster, interface-shaped delegation verified by `grep`).
- **Ten small files (full-read): 1,398 lines** — `revocation_data_verifier_factory.go` 285,
  `document_validator.go` 195, `certificate_validator.go` 99, `trust_anchor_verifier_factory.go` 80,
  `document_validator_factory.go` 49, `doc.go` 23, `cryptographic_suite_utils.go` 385,
  `validation_policy_with_cryptographic_suite.go` 155, `validation_policy_loader_with_crypto_suite.go` 96,
  `context_and_sub_context.go` 31.

### Findings

**SEC (trust-logic / false-accept)**

No findings. Every trust-critical path was verified fail-closed and 1:1 with the upstream Java contract:

- **No false "VALID" via a swallowed load failure.** Both `LoadValidationPolicy` implementations
  (`signed_document_validator.go:426`, `abstract_certificate_validator.go:395`) wrap the loader in a
  `recover()` that converts *any* panic (including `panic(nil)` — `recover()` returning nil does not
  skip the guard here, the guard is `if r := recover(); r != nil` but `r` is captured into `err` only on
  non-nil, so a nil panic is the *only* case that slips — see T30-STD-001) into a non-nil
  `*spiexception.IllegalInputException` with `validationPolicy = nil`. Both call sites
  (`ValidateDocumentWithPolicyAndCryptographicSuiteDocument` / `ValidateWithPolicyAndCryptographicSuiteDocument`)
  check `if err != nil { return nil, err }` **before** calling `ValidateDocumentWithValidationPolicy` /
  `ValidateWithValidationPolicy`, so a `nil` policy can never reach the executor on the success path.
  The executor only ever runs with a non-nil policy. No false pass.
- **`LevelFail` gating is correct in both verifier factories.** `RevocationDataVerifierFactory`
  (`revocation_data_verifier_factory.go:79`) enforces the revocation cryptographic suite **only** when
  `Level() == LevelFail`; otherwise it deliberately accepts all digest algorithms and all signature
  algorithms with `MinKeySize 0` (a documented `LOG.info` port). This is the correct upstream parity:
  "no constraint" means "accept all" — it does *not* mean "reject all", so it is not a false-accept and
  not a false-reject. `TrustAnchorVerifierFactory.acceptUntrustedCertificateChains`
  (`trust_anchor_verifier_factory.go:48`) returns `true` (accept untrusted) when the constraint is `nil`
  **or** not `LevelFail` — i.e. only a `LevelFail` constraint forbids untrusted chains. Correct.
- **`CryptographicSuiteUtils` trust checks are fail-closed.** `IsSignatureAlgorithmReliable` /
  `IsDigestAlgorithmReliable` return `true` when the suite is `nil` (no constraint = accept), and
  `false` when the algorithm is absent from the suite's acceptable map — a policy that omits an algorithm
  rejects it (no false-accept). `IsSignatureAlgorithmWithKeySizeReliable` and
  `GetExpirationDate*` return `false`/`nil` respectively when the suite is `nil`, so a missing policy
  cannot silently pass a key-size or expiry check. `reliableEvaluationExistsAtTime` requires a non-empty
  evaluations list and a validity window containing `validationTime` before returning `true`. No path
  returns "reliable" for a policy that has not affirmatively listed the algorithm.
- **`ValidationPolicyWithCryptographicSuite` scope resolution is correct.** The `cryptographicSuite`
  helper checks exact `Context+SubContext`, then `Context`-only, then fully-global, in that order, and
  falls back to the wrapped policy's own constraint when no override is found — matching Java's
  `getCryptographicSuite` precedence. The `SetCryptographicSuiteForContextAndSubContext` guard rejects
  `ContextEvidenceRecord` + non-empty `SubContext` (Java's `IllegalArgumentException`), so an
  ill-scoped suite cannot be installed. `Create()` (loader) replays scopes in insertion order
  (`cryptoOrder`) and panics (faithful to Java's unchecked `IllegalArgumentException`) on a
  guard violation.
- **No nil-deref reachable by hostile input.** `SignedDocumentValidatorFromDocument` nil-checks the
  document; `NewSignedDocumentValidatorBase` panics on a nil analyzer; `CertificateValidatorFromCertificate`
  panics on a nil token; `FromValidationPolicyDocument`/`FromValidationPolicyReader`/
  `WithCryptographicSuite*` all panic on `nil`/empty (1:1 ports of Java's `Objects.requireNonNull`);
  `AssertConfigurationValid` (both validators) returns a non-nil error when required state is missing,
  and `Validate*` returns `(zero, err)` before touching the executor. No unchecked type assertion, no
  `interface{}` deref, no map read on a nil map without a guard.
- **No log-injection / no panic-on-hostile** on the trust path. The only `panic`s are the
  `Objects.requireNonNull`-parity guards and the documented `recover()` re-wrapping in
  `LoadValidationPolicy`; none is reachable with a crafted *policy* document (a bad document surfaces as
  a factory "not supported" panic → recovered `IllegalInputException`, not a crash).

**PERF (BIGO / MEM / LEAK)**

No findings. The in-scope code is a thin dispatch + configuration layer over the (already-reviewed)
`spi/validation` executor; it performs no attacker-influenced algorithmic work beyond:

- `toSignatureAlgorithmWithKeySizesMap` (`revocation_data_verifier_factory.go:108`) — O(k) in the number
  of acceptable signature algorithms (bounded by the policy's algorithm list, small constant in practice);
  takes the min key-size per algorithm, O(k).
- `orderedStringSet` — O(n) dedup over the revocation-skip extension/policy lists (bounded by the policy,
  small).
- `GetReliableDigestAlgorithmsAtTime` / `GetReliableSignatureAlgorithmsWithMinimalKeyLengthAtTime`
  (`cryptographic_suite_utils.go`) — O(a·e) in the number of algorithms × evaluations (bounded by the
  suite, small); the map-iteration-order non-determinism is a documented, deliberate divergence (the
  values are only consumed by membership, not order).
- `Create()` (`validation_policy_loader.go:484`) — O(|cryptoOrder| × |scopes|), bounded by the number of
  `With*` calls the caller made (small; the builder is used once per validation).

No goroutines, no unbounded caches, no retained references, no allocation growth on hostile input, no
O(n²)+ on attacker-influenced input. `io.ReadAll` on the policy/crypto-suite streams is bounded by the
caller-supplied stream (a policy file is small by contract; a hostile multi-gigabyte stream would be a
caller error, and is not a distinct finding here beyond the general "streaming is in-memory only" known
gap in `docs/compatibility/known-gaps.md`).

**STD (idioms / error handling)**

T30-STD-001 — **Info** — Category: STD — `dss/validation/signed_document_validator.go:429` &
`dss/validation/abstract_certificate_validator.go:398`
- Evidence: both `LoadValidationPolicy` use
  ```go
  defer func() {
      if r := recover(); r != nil {
          validationPolicy = nil
          err = spiexception.NewIllegalInputExceptionWithCause("Unable to load the policy", fmt.Errorf("%v", r))
      }
  }()
  ```
- Impact: In Go, `panic(nil)` makes `recover()` return `nil`, so `r != nil` is false and the deferred func
  leaves `err == nil` — i.e. a `panic(nil)` raised *anywhere inside* the policy/crypto-suite loader
  (e.g. a future third-party factory, or a `defer func(){ panic(someNilErr) }()`) would be silently
  swallowed and `LoadValidationPolicy` would return `(nil, nil)`. Both call sites then pass a **nil**
  `validationPolicy` to `Validate*WithValidationPolicy`, which passes it to `ProcessValidationPolicy` →
  `exec.SetValidationPolicy(nil)`. Today this is unreachable: no in-tree code calls `panic(nil)` (verified
  by `grep 'panic(nil)'` across the module — zero in-tree hits), and every registered
  `ValidationPolicyFactory` / `CryptographicSuiteFactory` panics a non-nil message. It is the same root
  issue as **T28-STD-001** (`dss/format.go:331` `recovered`) but a *separate, more severe instance*: the
  T28 site is at the outer facade (a `panic(nil)` there returns a nil `Reports` to the caller, who can
  still check for nil), whereas *this* site is *inside* the validation flow and would hand a nil policy
  to the executor, which is a genuine trust-path false-pass if the executor tolerates a nil policy (it
  should not — `SetValidationPolicy` is typed `*modelpolicy.ValidationPolicy`, so a nil is representable —
  and the executor's own `spi/validation` code, reviewed in batch-07, is the layer that would have to be
  checked to confirm it rejects nil; see OQ-1).
- Recommendation: Defend the boundary the same way as T28-STD-001 — detect "a panic occurred" rather than
  "the recovered value is non-nil", e.g.
  ```go
  defer func() {
      if r := recover(); r != nil || panicked {   // panicked set in an outer defer
          ...
      }
  }()
  ```
  or, simplest, change the guard to `if r := recover(); r != nil || true { ... }` (always set err on any
  recover). Lowest-churn: document that in-tree factories must never `panic(nil)` and add a `go vet`-style
  lint rule. No action required today given the current codebase; worth a one-line note in PORTING.md
  alongside T28-STD-001.

T30-STD-002 — **Info** — Category: STD — `dss/validation/document_validator_factory.go:40` &
`dss/validation/policy/validation_policy_loader.go:60`
- Evidence:
  ```go
  var documentValidatorFactoryRegistry []DocumentValidatorFactory   // document_validator_factory.go:40
  var (
      validationPolicyFactoryRegistry   []modelpolicy.ValidationPolicyFactory   // validation_policy_loader.go:60
      cryptographicSuiteFactoryRegistry []modelpolicy.CryptographicSuiteFactory
  )
  ```
  mutated by `RegisterDocumentValidatorFactory` / `RegisterValidationPolicyFactory` /
  `RegisterCryptographicSuiteFactory` (plain `append`), read by `SignedDocumentValidatorFromDocument` /
  `loadDefaultPolicy` / `loadPolicy` / `loadDefaultCryptographicSuite` / `loadCryptographicSuiteCatalogue`.
- Impact: These are package-level global mutable slices with no synchronization. In practice they are
  registered once at `init()` time (the top-level `dss` package's `init()` in `validate.go`, and each
  format package's `init()` in `*_document_validator_factory.go`) and only *read* at runtime, so there is
  no data race on the normal path (init runs before any goroutine is spawned). A caller that registers a
  factory *after* goroutines have started reading the registry (e.g. a test that registers a mock factory
  in parallel with a validation) would race. This is the same benign pattern noted in U18a (batch-06) for
  the `tsl` job registry, and is a faithful port of Java's `ServiceLoader` (which is also a one-time
  registration mechanism). Noted for completeness; no action required for the in-tree codebase.
- Recommendation: None required. If a future feature allows runtime registration after goroutine start,
  add a `sync.RWMutex` or document that registration must complete before the first validation call.

T30-STD-003 — **Info** — Category: STD — `dss/validation/policy/cryptographic_suite_utils.go:126`
  (`IsSignatureAlgorithmKeyLengthBigEnough`, `if parameterMin == nil || *parameterMin < keySize`)
- Evidence: the comparison is strict-less-than (`*parameterMin < keySize`), so a key **exactly equal** to
  the suite's `Min` key-size returns `false` (not "big enough"). This is a 1:1 port of the upstream Java
  `isSignatureAlgorithmKeyLengthBigEnough`, which uses the same strict `<` (verified by the ported
  comment "this method only ensures that the key length is bigger than the minimal accepted key size").
- Impact: If a policy sets `Min` to the exact key size of a valid signature, the check returns `false` —
  a potential false-reject. However, this is a faithful 1:1 port of the upstream Java behavior (the
  method's own doc says "bigger than", not "at least"), so it is not a divergence. Flagged only because a
  reader might expect `<=`. No action required for parity; if a consumer of this API in Go (not via the
  Java-parity path) expects "at least" semantics, they should use `<=` in their own code. See OQ-2.
- Recommendation: None (1:1 port of upstream). If consolidation wants to harden, change to `<=` and add a
  `// DIVERGENCE, deliberate:` comment + a DESIGN.md entry, but only if the Java-parity tests do not pin
  the strict-`<` behavior.

T30-STD-004 — **Info** — Category: STD — `dss/validation/policy/cryptographic_suite_utils.go:96`
  (`parseKeySize`)
- Evidence:
  ```go
  func parseKeySize(keyLength string) int {
      if utils.IsStringDigits(keyLength) {
          size, err := strconv.Atoi(keyLength)
          if err == nil {
              return size
          }
      }
      return 0
  }
  ```
  `utils.IsStringDigits` (verified `dss/utils/strings.go:92`) accepts **Unicode Nd digits** (e.g. Arabic-Indic
  `٣`), but `strconv.Atoi` only parses ASCII digits — so a Unicode-digit key-length string passes
  `IsStringDigits` but `Atoi` returns an error, and `parseKeySize` silently returns `0`.
- Impact: A policy or document that names a key length using non-ASCII digits (e.g. `"٣٠٧٢"`) would be
  coerced to key-size `0`, which would then fail the `*parameterMin < keySize` check (any positive `Min`
  is `> 0`), producing a false-reject. In practice, key-length strings in DSS policies are ASCII (the
  JAXB schema and the ETSI policy files all use ASCII digits), so this is not reachable with any in-tree
  or known-corpus policy. Noted as a latent parity quirk; no action required.
- Recommendation: None (faithful port of Java's `NumberUtils.isDigits` + `Integer.parseInt` behavior,
  which has the same Unicode-digit gap). If hardening, use `strconv.Atoi`'s error directly and drop the
  `IsStringDigits` pre-check, or document the ASCII-only assumption.

T30-STD-005 — **Info** — Category: STD — `dss/validation/policy/cryptographic_suite_utils.go:340`
  (`isEvaluationApplicable`, "first come, first served" comment)
- Evidence:
  ```go
  for _, parameter := range parameterList {
      if !isSupported(algorithm, parameter) { continue }
      if keySize != nil && ((parameter.Min() != nil && *keySize < *parameter.Min()) ||
          (parameter.Max() != nil && *keySize > *parameter.Max())) { continue }
      return true
  }
  ```
  The loop returns `true` on the **first** parameter that is both supported by the algorithm *and* whose
  min/max range contains the key size, even if a *later* parameter (also supported) would have a
  *stricter* range. This is a faithful port of the upstream Java "first come, first served" comment.
- Impact: If a policy's parameter list contains two entries for the same algorithm (e.g. one with
  `Min=1024 Max=2048` and one with `Min=4096`), a key size of `1536` would match the first (accepted) even
  though the second (which the policy author may have intended as the binding constraint) would reject it.
  In practice, ETSI policies do not list duplicate algorithm entries, so this is not reachable. Noted as a
  latent parity quirk; no action required.
- Recommendation: None (1:1 port of upstream). If consolidation wants to harden, take the **strictest**
  (most restrictive) matching parameter rather than the first, and add a `// DIVERGENCE, deliberate:`
  comment + DESIGN.md entry.

No other STD findings. Error handling is consistent at the boundary (returned `error` wrapping the
original via `%w`/`%v`; panics are 1:1 ports of upstream `Objects.requireNonNull` / unchecked exceptions
below the facade; `recover()` re-wraps into `IllegalInputException`). No unchecked type assertions, no
global mutable state (beyond the T30-STD-002 registries, which are init-only), no goroutines, no retained
references in the in-scope code.

### Deliberate-exclusion notes

- **staticcheck style rules (ST\*/S1\*/QF\*/SA1019) and legacy-crypto flags** — deliberately disabled per
  `dss/PORTING.md` and the lint policy; not reported.
- **`// DIVERGENCE, deliberate:` comments** — the `revocation_data_verifier_factory.go` header
  (HashSet→slice ordering), the `validation_policy_loader.go` header (ServiceLoader stand-in, map order),
  the `validation_policy_with_cryptographic_suite.go` header (method-promotion delegation), and the
  `LoadValidationPolicy` `recover()` comments (Java `IllegalInputException` port) are all intentional and
  documented; not findings.
- **Classpath-to-filesystem collapse** (`document_validator.go`, `validation_policy_loader.go`) — a
  deliberate, documented divergence (Go has no classpath); not a finding.
- **`spi/validation` (batch-07)** — the `spi/validation` package is a DIFFERENT package from the
  top-level `validation/` package reviewed here. This unit cross-references batch-07 findings
  (T22D1A/B-SEC-001, T23-SEC-001, etc.) rather than re-deriving them. The `CryptographicSuiteUtils`
  functions here **back** the `validation/process/bbb/aov` checks (U32 scope), but the trust-logic in
  *this* file is what is reviewed here.
- **`validation/job` alert-error swallow** (batch-09 U28 OQ-2; U31 scope) — cross-referenced, not
  re-derived.
- **`tsl` pivot-cache data race** (batch-06 U18A-SEC-001; U31 scope) — cross-referenced, not re-derived.

### Out of scope (cross-referenced, not developed)

- `dss/validation/executor/` — U30b.
- `dss/validation/job/` — U31 (caches focus; pivot-cache race U18A-SEC-001; alert-error swallow U28 OQ-2).
- `dss/validation/reports/` — U31.
- `dss/validation/process/` + `dss/validation/blocks/` + BBB core — U32.
- `dss/policy/` + `dss/policy/crypto/` + `dss/policy/jaxb/` — U45.

### Tool log

Run from `dss/` (2026-08-28):

| Command | Output |
|---|---|
| `gofmt -l validation validation/policy` | *(no output — clean)* |
| `go vet ./validation/ ./validation/policy/...` | *(no output, exit 0 — clean)* |
| `golangci-lint run --config=../.github/.golangci.yml ./validation/ ./validation/policy/...` | `0 issues.` (exit 0) |

All tools available and passing; no in-scope tool findings.

### Summary

- **Total findings: 5** (SEC: 0, PERF-BIGO: 0, PERF-MEM: 0, PERF-LEAK: 0, STD: 5)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 0, Low: 0, Info: 5

The 13 in-scope files are a high-quality, faithful 1:1 port of the upstream `dss-validation` core facade
and `dss-validation/policy` package. The trust-logic is fail-closed on every path: `LoadValidationPolicy`
converts loader panics to errors and the caller checks `err` before passing the policy to the executor, so
a nil policy can never reach the executor on the success path; `LevelFail` gating in both verifier
factories is correct (no constraint = accept all, a `LevelFail` constraint = enforce); `CryptographicSuiteUtils`
trust checks are fail-closed (absent algorithm = reject, nil suite = no-constraint-accept for
`Is*Reliable`, fail-closed for `KeyLengthBigEnough`/`GetExpirationDate`); and the
`ValidationPolicyWithCryptographicSuite` scope resolution is correct. All five findings are Info-level:
a latent `panic(nil)` recovery gap in the two `LoadValidationPolicy` sites (same root as T28-STD-001, more
severe instance), two init-only global registries (benign, faithful ServiceLoader port), and three
faithful 1:1 ports of upstream quirks (strict-`<` key-size, Unicode-digit `parseKeySize`, first-match
`isEvaluationApplicable`) that are not reachable with in-tree or known-corpus policies but are noted for
consolidation. No SEC, PERF, Low, or Medium issues were found.

### Open questions (for consolidation)

1. **`panic(nil)` recovery gap (T30-STD-001):** cross-reference with T28-STD-001. Confirm during
   `spi/validation` (batch-07) consolidation that the executor's `SetValidationPolicy` / `Execute` path
   rejects a nil `*modelpolicy.ValidationPolicy` (it should — the type is a pointer, and the batch-07
   review verified fail-closed behavior — but an explicit nil-check in the executor would close the gap
   definitively). If the executor does *not* reject nil, this finding escalates to High (a `panic(nil)`
   in a factory → nil policy → executor runs with nil policy → potential false-pass).
2. **`isSignatureAlgorithmKeyLengthBigEnough` strict-`<` (T30-STD-003):** verify against the upstream
   Java source (`CryptographicSuiteUtils.isSignatureAlgorithmKeyLengthBigEnough`) that the `<` is indeed
   upstream's (not a porting error). If upstream uses `<=`, this is a High (false-reject on exact-min
   key-size). If upstream uses `<`, it is a faithful 1:1 port and Info only. Check the batch-02
   (`xml/utils`) or batch-07 (`spi`) review for a cross-reference, or the `dss/policy/crypto/` (U45)
   review for the JAXB schema that pins the `Min`/`Max` semantics.
3. **`isEvaluationApplicable` first-match (T30-STD-005):** confirm that no in-tree or known-corpus ETSI
   policy lists duplicate algorithm entries (which would make the first-match behavior observable). If
   consolidation finds a policy that does, this escalates to Medium (a stricter parameter is silently
   ignored).
4. **Global registry data race (T30-STD-002):** cross-reference with U18a (batch-06) `tsl` registry and
   U28 (batch-09) `dss` root registries. If consolidation wants a hard guarantee, add a `sync.RWMutex` to
   all three, or document that registration must complete before the first validation call.

## document + i18n + utils (unit U29)

### Scope

The `document` signing-service layer (generic base services, signature-requirements checker,
resources-handler family, profile bookkeeping), the `i18n` message-tag system (tag table,
provider, `MessageFormat` subset), and the `utils` shared helpers (strings, collections, codec,
ordered map, io, arrays, files). This is the layer that **builds** signatures and the layer the
per-format services (CAdES/XAdES/PAdES — reviewed in `batch-01b-cades-cms.md`) delegate to. A
signing-path bug is High; a false "signed/valid" is Critical.

### Date

2026-08-28

### Depth

Budgeted. The three largest files were read within their stated caps (`signature_requirements_checker.go`
first-150 + func map + 3 targeted regions; `message_tag.go` first-100 + func map + `MessageTagValues`/
`ValueOf`/`GetSemantic` regions); all other in-scope files were read in full.

### Files read (per-file lines actually read)

| File | Total | Read |
|---|---|---|
| `document/signature_requirements_checker.go` | 562 | ~350 (1–150, 177–236, 488–562 + func map) |
| `document/abstract_signature_parameters.go` | 291 | 291 (full) |
| `document/abstract_document_extender.go` | 181 | 181 (full) |
| `document/abstract_signature_service.go` | 141 | 141 (full) |
| `document/temp_file_resources_handler.go` | 122 | 122 (full) |
| `document/abstract_signature_parameters_builder.go` | 111 | 111 (full) |
| `document/temp_file_resources_handler_builder.go` | 107 | 107 (full) |
| `document/signed_document_extender.go` | 104 | 104 (full) |
| `document/profile_parameters.go` | 79 | 79 (full) |
| `document/dss_signature_security_factory.go` | 70 | 70 (full) |
| `document/in_memory_resources_handler.go` | 62 | 62 (full) |
| `document/signature_value_checker.go` | 59 | 59 (full) |
| `document/abstract_evidence_record_incorporation_parameters.go` | 58 | 58 (full) |
| `document/abstract_resources_handler.go` | 54 | 54 (full) |
| `document/document_signature_service.go` | 51 | 51 (full) |
| `document/multiple_documents_signature_service.go` | 38 | 38 (full) |
| `document/counter_signature_service.go` | 24 | 24 (full) |
| `document/in_memory_resources_handler_builder.go` | 23 | 23 (full) |
| `document/doc.go` | 20 | 20 (full) |
| `document/signed_document_extender_factory.go` | 15 | 15 (full) |
| `document/signature_extension.go` | 14 | 14 (full) |
| `document/evidence_record_incorporation_service.go` | 13 | 13 (full) |
| `i18n/message_tag.go` | 2122 | ~300 (1–100, 1050–1120, 2085–2122 + func map) |
| `i18n/i18n_provider.go` | 154 | 154 (full) |
| `i18n/message_format.go` | 125 | 125 (full) |
| `i18n/doc.go` | 11 | 11 (full) |
| `utils/strings.go` | 187 | 187 (full) |
| `utils/collections.go` | 152 | 152 (full) |
| `utils/codec.go` | 123 | 123 (full) |
| `utils/ordered_map.go` | 110 | 110 (full) |
| `utils/io.go` | 103 | 103 (full) |
| `utils/arrays.go` | 88 | 88 (full) |
| `utils/files.go` | 83 | 83 (full) |
| `utils/utils.go` | 16 | 16 (full) |
| `utils/doc.go` | 6 | 6 (full) |

**Total read: ~3,288 of 5,479 in-scope lines — under the 4,500 budget.** No region was left
"not fully verified (budget)" except the two large declarative tables
(`message_tag.go` const block lines ~100–1056 and `MessageTagValues` body lines ~1057–2096), which
are a flat constant table verified by shape (func map + head + tail) rather than line-by-line.

### Findings

**SEC (trust-logic / false-accept / codec correctness)**

No Critical or High. The signing path is fail-closed end to end (verified): `IsValidSignatureValue`
returns `false` on any build/verify error (`abstract_signature_service.go:119-141`);
`EnsureSignatureValue` returns an error on a digest/algorithm mismatch and only converts ECDSA
`r||s → DER` on an explicit match (`signature_value_checker.go:25-59`); `AssertSignaturesValid`
alerts on every `!IsSignatureValid()` (`signature_requirements_checker.go:519-556`);
`isSignatureGeneratedWithoutCertificate` panics (fail-closed) when a cert is required but absent
(`signature_requirements_checker.go:68-77`). The "accept-all revocation" verifier
(`createAcceptAllRevocationDataVerifier`, `signature_requirements_checker.go:222-229`) is the
explicit `CheckCertificateRevocation()==false` path, not an accidental bypass. This is consistent
with the U28 facade verdict (no false "VALID" at the boundary) and with `batch-01b`'s CMS/CAdES
findings, which this layer delegates to and are not re-reported.

The three findings below are **Info** — faithful 1:1 parity ports of commons-codec / commons-io
behaviour, all documented in-file and all failing closed downstream. They are hardening candidates,
not defects.

T29-SEC-001 — **Info** — Category: SEC — `dss/utils/codec.go:99-123` (`FromBase64`) & `:74-88` (`IsBase64Encoded`)
- Evidence: `FromBase64` silently discards non-alphabet bytes and drops a trailing `len%4==1`
  char (`if rem := len(filtered) % 4; rem == 1 { filtered = filtered[:len(filtered)-1] }`);
  `IsBase64Encoded` is a character-set-only check ("does NOT validate padding correctness or that
  the length is decodable").
- Impact: A malformed base64 payload decodes to a *different/shorter* byte string rather than
  erroring. Callers (`spi/dss_utils.go:501`, `jades/jades_crl_source.go:130`,
  `spi/dss_revocation_utils.go:706`, `asic/asic_manifest_parser.go:119`, ...) then feed the bytes to
  a parser that rejects them — so the net effect is a *parse failure* (fail-closed), not a false
  accept. The risk is only that a truncated-but-still-valid byte prefix could be silently accepted
  where a strict decoder would have rejected the input; no in-tree caller relies on that distinction
  for a trust decision.
- Recommendation: No action required for parity. If a hardening pass is desired, add a
  strict-decoding variant (`base64.StdEncoding.DecodeString`) for the trust-critical decode sites and
  route the lenient one to display/string-rendering only. Document the choice in `PORTING.md`.

T29-SEC-002 — **Info** — Category: SEC — `dss/utils/codec.go:17-20` (`IsHexEncoded`) as used by `dss/tsl/abstract_trusted_list_with_sha2_predicate.go:77`
- Evidence: `IsHexEncoded` returns true for *any* even-length hex string (`hex.DecodeString`
  succeeds); `Sha2Digest` then `FromHex`-decodes and wraps in
  `model.NewDigest(sha256, decoded)` with no length check.
- Impact: A TL "SHA-2" digest payload that is even-length hex but *not* 64 hex chars (not a real
  SHA-256) is hex-decoded to a wrong-length value and carried in the `Digest`. Upstream
  (commons-codec `Hex.decodeHex`) behaves identically, and the value is a *data* field of the
  trusted-list entry, not a gate on the verdict — so this does not flip a trust decision. It is a
  data-integrity softness on hostile TL input, not a false "VALID".
- Recommendation: No action required for parity. Optionally validate the decoded length against the
  expected digest size at the `Sha2Digest` call site and surface a `StatusAlert` (matches the
  "refuse what upstream merely warns about" precedent in `internal/pdf/crypt.go`).

T29-SEC-003 — **Info** — Category: SEC — `dss/utils/io.go:93-103` (`StartsWithStream`)
- Evidence: `io.ReadFull(r, temp)` into a `len(prefixArray)` buffer; on a *short* stream the unread
  tail stays zero-valued, and the doc comment states "a short stream compares equal only if
  prefixArray itself ends in that many zero bytes."
- Impact: If a caller used `StartsWithStream` to detect a fixed magic prefix and the stream is
  shorter than the prefix, a prefix ending in zero bytes could yield a false "starts-with" match.
  No in-tree caller passes a zero-terminated prefix, and this is a verbatim port of
  commons-io `IOUtils` read-then-compare, so it is not exploitable today.
- Recommendation: No action required for parity. If a strict "full prefix must be present" check is
  ever needed, add a length guard (`n1 == len(prefixArray)`) rather than relying on the zero-fill.

**PERF (BIGO / MEM / LEAK)**

T29-PERF-BIGO-001 — **Medium** — Category: PERF-BIGO — `dss/i18n/message_tag.go:1057` (`MessageTagValues`) as called by `MessageTagGetSemantic` (`:2114`) and `MessageTagValueOf` (`:2097`)
- Evidence: `MessageTagValues()` returns a fresh `[]MessageTag{ ~1000 constants }` composite literal
  on every call; `MessageTagGetSemantic`/`MessageTagValueOf` then linear-scan it
  (`for _, messageTag := range MessageTagValues() { if string(messageTag)==... }`). Report builders
  call `MessageTagGetSemantic` per indication/sub-indication
  (`validation/executor/detailed_report_builder.go:513,524`,
  `simple_report_builder.go:269,280`).
- Impact: Each call allocates a ~16 KB slice (`1000 × 16-byte` string headers) and does an O(n)
  scan; a report with K indications performs K allocations + O(K·n) comparisons. Bounded and
  sub-millisecond for realistic K, but it is a needless per-call allocation on the report-building
  hot path and diverges from upstream, where `MessageTag.valueOf` resolves against a static
  constant set without a fresh allocation.
- Recommendation: Cache the values slice in a package-level `var allMessageTags = []MessageTag{...}`
  (built once) and have `MessageTagValueOf`/`MessageTagGetSemantic` range over it. O(1) allocation
  after init, identical observable order (declaration order is preserved), no wire-format change.

T29-PERF-LEAK-001 — **Info** — Category: PERF-LEAK — `dss/document/temp_file_resources_handler.go:43-64` & `:114-121`
- Evidence: `NewTempFileResourcesHandler` does `os.CreateTemp` with `toBeDeleted=true`; the file
  header DEVIATION states Go has no portable `deleteOnExit` equivalent and "Callers are expected to
  Close … every handler they create"; `ForceDelete` (`:114`) does `_ = os.Remove(...)`; the builder's
  `Clear()` (`temp_file_resources_handler_builder.go:96-101`) `ForceDelete`s **every** handler it
  ever created.
- Impact: Two documented, Java-parity softnesses. (1) If the process is killed (SIGKILL / an
  `os.Exit` panic) or a caller drops the handler without `Close`/`ForceDelete`, the `dss-*.tmp` file
  in `$TMPDIR` persists — a temp-file leak with no automatic cleanup. (2) `Clear()` force-deletes
  even handlers already handed to a `FileDocument` via `WriteToDSSDocument` (which sets
  `toBeDeleted=false`), so calling `Clear()` after `WriteToDSSDocument` removes a live output
  document. Both are covered by the in-file DEVIATION / `Clear()` javadoc ("do not forget to
  preserve the output documents") and mirror upstream's `File#deleteOnExit` + `clear()` contract.
- Recommendation: No action required for parity. Hardening options: (1) a `defer`-friendly
  `Cleanup()` method and/or a process-startup sweep of stale `dss-*.tmp` files in the configured
  dir; (2) have the builder track which handlers were handed to a `FileDocument` and skip them in
  `Clear()`, or make `Clear()` respect `toBeDeleted` instead of unconditionally force-deleting.

No PERF-MEM finding beyond the documented known-gap (`CMS streaming is in-memory only` —
`docs/compatibility/known-gaps.md`), which is out of scope to re-report. `SetCertificateChainFromTokens`
(`abstract_signature_parameters.go:143-158`) is an O(n²) dedup but n is the certificate-chain length
(CA depth, small and not attacker-driven), so it is not rated.

**STD (idioms / error handling / global mutable state)**

T29-STD-001 — **Low** — Category: STD — `dss/document/signed_document_extender.go:80-89` (`signedDocumentExtenderFactories` registry)
- Evidence: `var signedDocumentExtenderFactories []SignedDocumentExtenderFactory` (package-level);
  `RegisterSignedDocumentExtenderFactory` does an unguarded
  `signedDocumentExtenderFactories = append(...)`, and `FromDocument` ranges over it.
- Impact: Registration is meant to happen from `init()` (single-threaded, before `main`), so it is
  safe in the intended flow and mirrors the upstream `ServiceLoader` registry. But the function is
  exported and, if called concurrently at runtime (e.g. a plugin registering after startup), the
  `append` is a data race. No in-tree code does this.
- Recommendation: No action required for the current flow. If runtime registration is ever expected,
  guard the slice with a `sync.RWMutex` (or freeze it after `init` and return an error on late
  registration).

T29-STD-002 — **Low** — Category: STD — `dss/document/abstract_signature_parameters.go:143-158` (`SetCertificateChainFromTokens`)
- Evidence: for each incoming certificate it does a linear scan over `p.certificateChain`
  (`for _, existing := range p.certificateChain { if existing == certificate { found = true } }`)
  before appending — O(n²) in the number of tokens supplied.
- Impact: Negligible in practice (n is the chain length, a handful of certs), but it is a quadratic
  idiom where a `map`/set lookup would be O(n) and clearer. Not a correctness issue (identity
  dedup matches Java's "if already present then it is ignored").
- Recommendation: No action required. Optional: dedup via a `map[*model.CertificateToken]struct{}`
  built once, preserving first-seen order, to make the O(n) intent explicit.

No unchecked type assertions, no goroutines, no `io.Copy` error drops, no nil-deref panics on
hostile input, and no `ST*`/`S1*`/`QF*`/`SA1019` categories are in play (deliberately disabled per
`PORTING.md`). The `i18n/message_format.go` apostrophe-swallowing behaviour is an explicit,
documented upstream quirk (parity contract) and is not a finding.

### Tool log

Run from `dss/`:

- `gofmt -l document i18n utils` → no output (clean).
- `go vet ./document/... ./i18n/... ./utils/...` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./document/... ./i18n/... ./utils/...`
  → `0 issues` (exit 0).

All tools available and passing; nothing to note as unavailable.

### Summary

- **Total findings: 7** (SEC: 3, PERF-BIGO: 1, PERF-LEAK: 1, PERF-MEM: 0, STD: 2)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 1, Low: 2, Info: 4
- The three packages are a clean, thin 1:1 port. The signing path is fail-closed at every check
  (invalid signature → alert/panic; algorithm/digest mismatch → error; missing cert → panic), which
  is consistent with the U28 facade verdict and with the `batch-01b` CMS/CAdES results this layer
  delegates to. No false "signed/valid" was found. The one Medium is a needless per-call
  ~16 KB allocation on the report-building hot path (`MessageTagValues`), trivially fixable by
  caching the slice. The two Low STD items are an unguarded (but `init`-only) global registry and an
  O(n²) chain dedup. The four Info items are all faithful commons-codec/commons-io parity ports
  (or a documented temp-file DEVIATION) that fail closed downstream and are already documented
  in-file; they are hardening candidates, not defects.

### Open questions (for consolidation)

1. **Codec leniency (T29-SEC-001/002):** the lenient `FromBase64`/`IsHexEncoded` are parity with
   commons-codec and fail closed downstream. Decide once, for the whole module, whether trust-critical
   decode sites should switch to strict decoders (and be recorded as `DIVERGENCE, deliberate` per
   the `internal/pdf/crypt.go` precedent) or stay lenient. Affects `spi`, `jades`, `asic`, `tsl`
   callers, not just `utils`.
2. **`MessageTagValues` cache (T29-PERF-BIGO-001):** a package-level cached slice is a clean fix but
   changes `message_tag.go` (a 1:1 port file). Confirm the cache preserves declaration order (it
   does) and that no test relies on `MessageTagValues()` returning a fresh slice each call.
3. **Temp-file cleanup (T29-PERF-LEAK-001):** whether to add a stale-`dss-*.tmp` sweep and/or a
   `Clear()` that respects `toBeDeleted`. Both diverge from the upstream `clear()`/`deleteOnExit`
   contract and would need a `DIVERGENCE, deliberate` entry if adopted.
4. **Registry concurrency (T29-STD-001):** confirm no format package (CAdES/XAdES/PAdES/JAdES/ASiC)
   calls `RegisterSignedDocumentExtenderFactory` outside `init()`; if any does, the unguarded append
   becomes a real race. (Cross-reference the format packages' `init()` registrations during U29
   consolidation.)

## validation/executor (unit U30b)

### Scope

The process-executor layer that drives validation and builds the three report
flavours (ETSI TS 119 102-2 validation report, Simple Report, Detailed Report)
from already-built diagnostic data: the ETSI / Simple / Simple-for-certificate
report builders, the Detailed-report builder family (base + certificate + QWAC +
EAA-presentation variants), and the process executors (signature, certificate,
QWAC, EAA-presentation) that wire the two together. This is a 1:1 Go port of the
Java `dss-validation/.../executor` package (DSS 6.5.RC1). A false "VALID" verdict
is the Critical lens; a status that maps to the wrong ETSI indication is High; a
nil policy reaching the executor and skipping checks is Critical if the executor
tolerates it.

### Date

2026-08-28

### Depth

Budgeted (per-file line caps on the four largest; full-read on the fourteen small
files). The four budgeted files were read first-150 + `grep -n 'func '` + targeted
regions; every trust-critical path (status mapping, verdict assembly, indication
mapping, nil-policy handling) was read directly.

### Files read (per-file lines actually read)

| File | Total | Read |
|---|---|---|
| `etsi_validation_report_builder.go` | 2139 | ~415 (1–150, 338–400, 646–745, 934–1090, 1405–1475, 2095–2139) |
| `simple_report_builder.go` | 1480 | ~430 (1–150, 150–445, 571–660, 1129–1200) |
| `simple_report_for_certificate_builder.go` | 792 | ~340 (1–130, 484–610) |
| `detailed_report_builder.go` | 529 | 529 (full) |
| `abstract_detailed_report_builder.go` | 284 | 284 (full) |
| `detailed_report_for_certificate_builder.go` | 201 | 201 (full) |
| `detailed_report_for_qwac_builder.go` | 178 | 178 (full) |
| `default_signature_process_executor.go` | 168 | 168 (full) |
| `default_certificate_process_executor.go` | 122 | 122 (full) |
| `abstract_process_executor.go` | 110 | 110 (full) |
| `detailed_report_for_eaa_presentation_builder.go` | 67 | 67 (full) |
| `process_executor.go` | 48 | 48 (full) |
| `qwac_certificate_process_executor.go` | 37 | 37 (full) |
| `eaa_presentation_process_executor.go` | 37 | 37 (full) |
| `document_process_executor.go` | 29 | 29 (full) |
| `doc.go` | 24 | 24 (full) |
| `certificate_process_executor.go` | 24 | 24 (full) |
| `process_executor_provider.go` | 21 | 21 (full) |

**Total in-scope code read: ≈3,148 lines** (within the 3,500-line budget). Not
read (out of scope / test-only): `full_corpus_oracle_test.go`,
`full_corpus_reports_oracle_test.go`, `report_builders_oracle_test.go`,
`java_hash_set_order_test.go`, `eaa_build_tag_*_test.go` — these are the
byte-parity and ordering oracle tests that pin the deliberate ordering
divergences noted below, and are covered by the `go test` gate, not re-derived.

### Findings

**SEC (trust-logic / false-accept)**

T30B-SEC-001 — **Info** — Category: SEC — `dss/validation/executor/abstract_process_executor.go:100–110`
  (`AssertConfigurationValid`; nil check at `:104`)
- Evidence:
  ```go
  func (e *AbstractProcessExecutor) AssertConfigurationValid() {
      if e.JaxbDiagnosticData == nil {
          panic("The diagnostic data is missing")
      }
      if e.Policy == nil {
          panic("The validation policy is missing")
      }
      ...
  }
  ```
  Both public entry points invoke it *first*: `default_signature_process_executor.go:108`
  (`Execute`) and `default_certificate_process_executor.go:84` (`Execute`). `policy.ValidationPolicy`
  is a Go **interface** (`model/policy/validation_policy.go:12`), so a nil policy is representable
  and detectable here. The `qwac`/`eaa` executors embed these two defaults and therefore inherit the
  same `Execute` guard.
- Impact: This is a **positive** finding: the executor rejects a nil `*policy.ValidationPolicy`
  (a nil interface) by panicking, and that panic is converted to a returned `error` at the facade
  boundary (`dss/validate.go:154` `recovered("validate", …)` wraps `exec.Execute()`). No trust
  check is skipped and there is no false "VALID". This **closes T30-STD-001's open question
  (OQ-1)** from the U30 unit: the executor *does* reject a nil policy, so the `panic(nil)`
  recovery gap in `LoadValidationPolicy` (which is the actual root, in `validation/`, not this
  unit) does **not** escalate to High/Critical. The nil policy cannot reach the BBB/qualification
  blocks: `AssertConfigurationValid` panics before `BuildReports` → `DetailedReportBuilderFor` →
  `b.Policy.*` is ever dereferenced. The only residual is the U30-side `recover()` gap, which is
  already filed as T30-STD-001 and T28-STD-001.
- Recommendation: None required in this unit — the guard is correct and fail-closed. Keep it as the
  authoritative backstop; the root `panic(nil)`-recovery fix belongs in `validation/` (U30), not
  here. (See T30B-STD-001.)

No false "VALID" found. The status-mapping and verdict-assembly paths were verified correct and
1:1 with upstream:

- **Indication mapping is correct.** `finalIndication` (`detailed_report_builder.go:376`) and the
  QWAC twin (`detailed_report_for_qwac_builder.go:166`) both map
  `Passed → TotalPassed`, `Indeterminate → Indeterminate`, `Failed → TotalFailed`, and panic on any
  other value (a faithful port of the upstream `DSSReportException`; `detailed_report_for_qwac_builder.go:166`).
  There is no path that maps a
  non-passed indication to a "valid" conclusion. The ETSI report's `MainIndication` is copied
  verbatim from the detailed report's `FinalIndication`
  (`etsi_validation_report_builder.go:947` `fillIndicationSubIndication`) — it is never re-derived,
  so it cannot diverge from the verdict the process computed.
- **Simple-report "valid" count is conservative.** `isValidConclusion`
  (`simple_report_builder.go:1129`) counts a token valid only on `TotalPassed` or `Passed`;
  `validSignatureCount` is incremented only on `IndicationTotalPassed`
  (`simple_report_builder.go:356`). An `Indeterminate`/`Failed` signature is never counted valid.
- **In-place Conclusion mutation is a faithful Java-parity port, not a defect.**
  `detailed_report_for_qwac_builder.go:104–107` reassigns `conclusion.Indication` on the *shared*
  `XmlConclusion` so that both the basic and the top-level `<Signature>` carry the `TOTAL_*`
  indication. The file header documents this as a deliberate port of Java's in-place mutation of the
  shared conclusion object. It only ever *promotes* `Passed → TotalPassed` (or leaves
  `Indeterminate`/`Failed` as-is) — it never demotes a `Failed`/`Indeterminate` to `Passed`, so it
  cannot manufacture a false "VALID".
- **No nil-deref panic on the read paths reached.** `readableCertificateName` / `isTrustAnchor`
  (`simple_report_builder.go:633–644`) use `UsedCertificateByIdNullSafe`, which returns a
  non-nil zero `CertificateWrapper` for an unknown id, so no nil deref. The one guarded nil-deref
  (`addEvidenceRecordValidationReportData`, `etsi_validation_report_builder.go:683–685`,
  `xmlEvidenceRecord.ValidationProcessEvidenceRecord.AOV`) is only reached when
  `XmlEvidenceRecordById` is non-nil (the caller `evidenceRecordValidationReport` returns early
  otherwise, `:650–652`), so it is a 1:1 port of a Java invariant, not a hostile-input crash.
- **No unchecked type assertion on the trust path.** The token type-switch in
  `fillTokenCryptographicInfo` (`:1067–1080`, with a `default: panic` arm) and the
  `collectReportIndications` switch
  (`detailed_report_builder.go:392–399`) have `default:`/`panic` arms for the unknown case, so an
  unexpected token type is a loud failure, not a silent misreport.

**PERF (BIGO / MEM / LEAK)**

T30B-PERF-001 — **Info** — Category: PERF-BIGO — `dss/validation/executor/etsi_validation_report_builder.go:1416–1458`
  (`buildCertIDListType`)
- Evidence:
  ```go
  if len(orphanCerts) > 0 {
      allOrphanCertificates := b.diagnosticData.AllOrphanCertificateObjects()
      for _, orphanCert := range orphanCerts {
          if orphanCert != nil {
              if len(orphanCert.References()) > 0 && !containsOrphanCertificateId(allOrphanCertificates, orphanCert.Id()) {
  ```
  where `containsOrphanCertificateId` (`:2095`) is a linear scan
  (`for _, token := range tokens { if token.Id() == id { return true } }`).
- Impact: For a signature whose `SigningCertificate`/`CompleteCertificateRefs` attribute references
  *m* orphan certificates, the builder performs *m* linear scans over *n* orphan token objects, i.e.
  O(m·n). Both *m* and *n* are attacker-influenceable by a crafted CMS/PAdES document. In practice
  both are small (a handful of orphans in the corpus), so this is a bounded, low-constant O(n²) —
  not exploitable at realistic sizes, and it is a direct port of Java's
  `List.contains(orphanToken)` over the orphan list. Flagged for completeness under the
  PERF-BIGO lens.
- Recommendation: If consolidation wants a hard bound, index `allOrphanCertificates` by id once
  (a `map[string]struct{}`) before the loop. This does not change any output (membership test only),
  so it is safe to adopt without a parity entry; otherwise leave as the faithful 1:1 port.

No other PERF findings. The builders are single-pass over the diagnostic data and the detailed
report; the maps (`signatureIdentifierMap`, `validationObjectMap`,
`finalIndications`/`finalSubIndications`, `bbbs`) are bounded by the number of tokens/signatures and
are not re-grown on hostile input. `javaHashSetOrder` / `JavaHashSetStringOrder`
(`abstract_detailed_report_builder.go:216–284`) are O(n log n) and run once per token list, not per
element. There are no goroutines in this package, no unbounded caches, and no retained references
beyond the per-`Build()` maps (which are released with the builder).

**STD (idioms / error handling)**

T30B-STD-001 — **Info** — Category: STD — `dss/validation/executor/abstract_process_executor.go:100–110`
  (`AssertConfigurationValid`; `if e.Policy == nil { panic(…) }` at `:104–105`)
- Evidence: same guard as T30B-SEC-001. This is the executor's backstop against a nil policy.
- Impact: Correct and fail-closed. It is the *defense* that makes T30-STD-001 (the
  `panic(nil)`-recovery gap in `validation/signed_document_validator.go` /
  `abstract_certificate_validator.go`'s `LoadValidationPolicy`) a non-escalating Info: even if a
  `panic(nil)` in a policy factory slipped through the `recover()` and produced `(nil, nil)`, the
  executor's `AssertConfigurationValid` would panic on the nil policy, which the facade's `recovered`
  converts to an error. The only residual is that this guard relies on `policy.ValidationPolicy`
  remaining an **interface** (so a nil is detectable); if a future change made it a value type, the
  `e.Policy == nil` comparison would silently never fire and the guard would need a different
  sentinel. Noted so consolidation does not "simplify" it away.
- Recommendation: None required. Keep the interface typing and the nil guard. If the type is ever
  changed to a value type, re-derive the guard (e.g. a `HasPolicy()`/`PolicyName()` probe or an
  explicit `SetValidationPolicy(nil)` rejection at the setter).

No other STD findings. Error/panic handling is consistent at the boundary: panics are 1:1 ports of
upstream `Objects.requireNonNull` / `DSSReportException` (below the facade), and the facade's
`recovered` wrapper (`dss/format.go:329`) turns them into `error`s that satisfy `errors.As`. No
global mutable state, no goroutines, no unchecked `interface{}` deref, no unchecked type assertion
on the trust path (every type-switch has a `default`/`panic` arm).

### Deliberate-exclusion notes

- **staticcheck style rules (ST\*/S1\*/QF\*/SA1019) and legacy-crypto flags** — deliberately
  disabled per `dss/PORTING.md` and the lint policy; not reported.
- **Java `HashMap`/`HashSet` iteration order** — the `javaHashSetOrder` / `JavaHashSetStringOrder`
  helpers and the `BBBOrder` insertion-order record are deliberate, documented ports of Java's
  `HashSet`/`LinkedHashMap` order, pinned by `java_hash_set_order_test.go` and the oracle tests;
  not findings. The `SimpleReportBuilder`/`DetailedReportBuilder` `addSemantics` methods emit
  indications in *enumeration declaration order* (documented deviation from Java's unstable identity
  hash order) — deliberate and byte-visible only when `includeSemantics` is on (off by default).
- **`// DIVERGENCE, deliberate:` / `Judgment call:` comments** — the null-vs-empty byte-preserving
  helpers (`nullableString`, `emptyToNilURLs`, `xsStringValue`/`escapeCharData`), the
  `Updated`-assigned-twice statement-fidelity line (`simple_report_builder.go:1160–1161`), and the
  QWAC in-place Conclusion mutation are all intentional and documented; not findings.
- **`spi/validation` (batch-07)** and **`validation/process` + BBB (U32)** — the executors delegate
  to these; their trust-logic is cross-referenced, not re-derived. The executor passes
  `e.Policy` through to the process blocks (`detailed_report_builder.go:263`, `:275`, `:292`,
  `:304`, `:316`, `abstract_detailed_report_builder.go:101`); a nil policy is rejected upstream of
  all of them by `AssertConfigurationValid` (T30B-SEC-001).
- **`validation/job` cache data race (U18A-SEC-001)** and **`validation/job` alert-error swallow
  (U28 OQ-2)** — not present in this package (no caches, no alert routing here); cross-referenced,
  not re-derived.

### Out of scope (cross-referenced, not developed)

- `dss/validation/process/` + `dss/validation/blocks/` + BBB core — U32.
- `dss/validation/job/` — U31.
- `dss/diagnostic/` (the `Data`/wrapper accessors the builders call, e.g.
  `UsedCertificateByIdNullSafe`, `AllOrphanCertificateObjects`) — the diagnostic-data package; its
  own O(n) id lookups are out of this unit's scope (the O(m·n) at the call site is T30B-PERF-001).
- `dss/policy/` + `dss/policy/crypto/` + `dss/policy/jaxb/` — U45.

### Tool log

Run from `dss/` (2026-08-28):

| Command | Output |
|---|---|
| `gofmt -l validation/executor` | *(no output — clean)* |
| `go vet ./validation/executor/` | *(no output, exit 0 — clean)* |
| `golangci-lint run --config=../.github/.golangci.yml ./validation/executor/` | `0 issues.` (exit 0) |

All tools available and passing; no in-scope tool findings.

### Summary

- **Total findings: 3** (SEC: 1, PERF-BIGO: 1, PERF-MEM: 0, PERF-LEAK: 0, STD: 1)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 0, Low: 0, Info: 3

The 18 in-scope files are a clean, faithful 1:1 port of the upstream `dss-validation/.../executor`
package. The trust-path is fail-closed at every check: a nil policy is rejected by
`AssertConfigurationValid` before any process block runs (closing T30-STD-001's OQ-1 and keeping it
at Info); the indication→verdict mapping is correct and verbatim-copied from the process's
conclusion (no re-derivation, no false "VALID"); the Simple-report "valid" counter is conservative
(TotalPassed only); and every type-switch on the trust path has a `default`/`panic` arm. The three
Info findings are (a) the correct nil-policy guard that backstops the U30 `panic(nil)`-recovery gap,
(b) a bounded O(m·n) orphan-certificate membership scan that is a faithful 1:1 port of Java's
`List.contains`, and (c) the same nil-policy guard flagged as an STD note so consolidation does not
accidentally change `ValidationPolicy` from an interface (which the guard's `== nil` test depends
on). No SEC defect, no High/Medium, no leak, and no goroutine were found.

### Open questions (for consolidation)

1. **T30-STD-001 OQ-1 (nil policy) — ANSWERED here:** the executor rejects a nil
   `*policy.ValidationPolicy` via `AssertConfigurationValid`
   (`abstract_process_executor.go:100–110`; nil check `:104`), invoked first by both `Execute()` entry points
   (`default_signature_process_executor.go:108`, `default_certificate_process_executor.go:84`) and
   converted to an `error` by the facade's `recovered` wrapper (`dss/validate.go:154`). The U30
   `panic(nil)`-recovery gap therefore does **not** escalate; keep it at Info. No further action in
   this unit.
2. **`buildCertIDListType` O(m·n) (T30B-PERF-001):** confirm the corpus has no document with a
   large orphan-certificate fan-out (m·n large). If one is added, index the orphan list by id
   before the loop (membership-only change, no parity entry needed).
3. **`ValidationPolicy` interface typing (T30B-STD-001):** confirm no planned change converts
   `model.policy.ValidationPolicy` from an interface to a value type; if so, re-derive the nil guard
   in `AssertConfigurationValid`, since `e.Policy == nil` would silently stop firing.
