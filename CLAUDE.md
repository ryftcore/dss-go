# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

**esig** (`github.com/ryftcore/dss-go/dss`) is a pure-Go, file-by-file port of the EU
Java library [esig/dss](https://github.com/esig/dss), pinned to **DSS 6.5.RC1**
(`4c2129862948bfd53ca1455832260aa17e183cf8`). It signs and validates CAdES, XAdES,
PAdES, JAdES and ASiC documents and runs the full ETSI EN 319 102-1 validation
process. The compatibility contract is *interoperability with Java DSS* — byte-level
signature output, validation verdicts and report XML must match upstream for the same
input — verified continuously against committed Java oracle fixtures.

This shapes almost every rule below: "idiomatic Go" is deliberately constrained by
"still maps 1:1 to the Java source".

## Commands

**The Go module root is `dss/`, not the repo root.** Run every `go`/`make` command from
inside `dss/`. Go 1.27+ (the API uses generic methods).

```sh
cd dss
make build        # go build ./...
make vet          # go vet ./...
make test         # go test ./... -timeout 40m
make race         # go test ./... -race -timeout 60m   (asic merge tests exceed 10m under -race)
make check        # build + vet + test
make cli-build    # builds ./bin/esig

# Full CI gate (what a PR must pass):
gofmt -l .                                   # must print nothing
go vet ./... && go vet -tags eaa ./...
go build ./... && go build -tags eaa ./...
go test ./... -count=1 -timeout 40m
golangci-lint run --config=../.github/.golangci.yml ./...   # config is NOT auto-discovered; pass it explicitly

# Single package / single test:
go test ./cades/ -count=1
go test ./internal/pdf/ -run 'TestCrypt' -count=1 -v
go run ./examples/01-sign-pdf-pades
```

CI splits `go test` into three package groups (`core`, `formats`, `validation` — see
`.github/workflows/ci.yml`) purely for parallelism; if you add a top-level package, add
it to one of those groups.

`-tags eaa` builds the electronic-attestation-of-attributes code paths. A change must
not break either the default or the `eaa` build.

### Tests that skip, and why

- `corpus/` at the **repo root** (outside the module, ~180 MB) holds heavy oracle
  fixtures, mirroring the package layout under `dss/` (`corpus/validation/executor/testdata/…`).
  Tests reach it via `internal/corpustest.Path`/`RootPath`, which `t.Skip` when `corpus/` is
  absent (a bare `go get`) but **fail** when `corpus/` exists and a fixture is missing.
  In a full checkout no corpustest skip may fire — CI greps the log for
  `corpustest: corpus/ not found` and fails the build if it appears.
- Eight tests (`TestDownstreamCrossValidation/*`, `TestCAdESLevelBaselineBBytesMatchUpstream`,
  `TestRoundTripJavaBuiltExtendedByGo`) drive a live Java DSS and skip unless
  `DSS_UPSTREAM_HOME` points at a built upstream checkout with `java`/`javac`/`mvn` on PATH.
- `TestWriterOverCorpus` (`PDF_CORPUS_DIR`) and `TestDumpWriterOutputs` (`PDFOUT`) are opt-in.

Never weaken or delete an existing test (especially an `*_oracle_test.go`) to make it
pass. A parity mismatch means the port needs a fix or the divergence needs a documented
reason.

## Architecture

### Layers (top to bottom)

```
dss/  (root package)          thin facade: Sign / Extend / Validate / OpenPKCS12. No logic of its own.
  cmd/esig/                   CLI built on the facade; ships its own RFC 3161 HTTP client.
  examples/                   nine runnable programs on the facade.
cades xades pades jades asic/{cades,xades}   one package per format family: Service, SignatureParameters,
                                            the format's AdvancedSignature + document-validator factory.
validation/                   EN 319 102-1 process: executor/ (process executors + report builders),
                              process/{bbb,vpf*,qualification} (building blocks), policy/ (policy loader),
                              reports/, timestamp/, job/ …
spi/                          dss-spi: certificate/revocation sources, tokens, DSSUtils, x509, tsl, validation
                              (CertificateVerifier, AdvancedSignature, TimestampToken, TSPSource).
model enumerations utils alert i18n   leaf model/enum/utility packages.
policy diagnostic simplereport detailedreport validationreport simplecertificatereport   JAXB-schema ports.
tsl trustedlist               TL/LOTL validation job and trusted-list parsing.
internal/                     engines replacing third-party Java libs with NO Java class to mirror.
```

`internal/` packages replace BouncyCastle / PDFBox / Santuario / jose4j / Xerces:
`asn1ber` (BER/DER engine), `cmscore` (RFC 5652 CMS + RFC 3161), `pdf` (ISO 32000 reader and
incremental writer), `xmldom`, `xmlc14n` (7 canonicalization variants), `xmldsig`,
`xpath10`, `jose`, `pfx` (PKCS#12 reader), `eccurve`. Rules for them: they import the stdlib
and at most other `internal/` packages — **never a DSS package** — and state their
provenance in `doc.go` instead of a `// Ported from` header. An allowed `golang.org/x/…`
module is permitted only where the stdlib has no counterpart to the Java API being mirrored,
and `doc.go` plus the DESIGN entry must name it (today: `internal/pdf/saslprep.go` on
`x/text/unicode/{norm,bidi}`). `internal/pdf/DESIGN.md` and
`internal/xmldom/DESIGN.md` are the design records; deliberate divergences from upstream get
an entry there.

### Cross-cutting mechanisms

- **Registration replaces `java.util.ServiceLoader`.** Each format package's
  `*_document_validator_factory.go` has an `init()` that calls
  `validation.RegisterDocumentValidatorFactory`; the root `dss` package blank-imports all
  six format packages (`format.go`) and registers the ETSI policy and XML/JSON crypto
  suites (`validate.go`). Registration **order matters** (first factory = default). Code that
  uses the lower packages without importing `dss` must register explicitly.
- **Panics → errors at the facade boundary.** Ported services follow Java and raise
  unchecked exceptions as panics. `format.go`'s `recovered()` turns them into errors that
  still satisfy `errors.As(err, *model.DSSError)`. Below the facade, callers recover
  themselves.
- **Facade types are type aliases** of the underlying types, so the two levels mix with no
  conversion.
- **No network by default.** Upstream's `OnlineTSPSource`/`OnlineCRLSource`/`OnlineOCSPSource`
  are not ported; T/LT/LTA need a caller-supplied `TSPSource` and revocation sources.
  `spi/validation.KeyEntityTSPSource` issues RFC 3161 tokens from a local key for tests.
- `harness/` is a test-only package (no non-test `.go`) holding the cross-format
  document-level oracle; it lives outside the format packages to avoid import cycles.

### Finding the Go twin of a Java class

Every ported file starts with `// Ported from <upstream-path> (DSS 6.5.RC1).` —
`grep -r "SomeClass.java" dss/` finds it. The Java→Go naming rules and the
Maven-module→Go-package table are in `docs/migrating-from-java/index.md`.

## Porting conventions (binding — see `dss/PORTING.md` for the full text)

- One Go file per Java class (`SignatureLevel.java` → `signature_level.go`), tests alongside.
- **Exported identifiers never repeat the package name**: `CAdESService` → `cades.Service`,
  `AlertHandler` → `alert.Handler`, `NewASiCContent` → `asic.NewContent`. Exceptions: a
  name equal to the package name, and a name whose next rune is a digit
  (`xades.XAdES111XSDUtils`). `revive`'s `exported` stutter findings must be eliminated.
- **Enums** are typed string constants whose *value* is exactly Java's `name()`; the Go
  identifier is `TypeMixedCaps` (`SignatureLevelXAdESBaselineB = "XAdES_BASELINE_B"`,
  `DigestAlgorithmSHA256`), never `Type_JAVA_NAME`. Crypto/standards initialisms stay
  upper-case; other all-caps parts are Titlecased.
- **Never invent, abbreviate or "fix" an OID, URI, enum value, XML name, i18n key or JSON
  key.** They are wire format and are copied verbatim from upstream.
- Java exceptions → returned `error` (types like `model.DSSError`); throwing constructors →
  `New…() (T, error)`. Static utility classes → prefixed package functions
  (`DSSUtils.loadCertificate` → `spi.DSSUtilsLoadCertificate`).
- Order-sensitive Java iteration (HashMap etc.) → slice of pairs or explicit sort;
  `*_determinism_test.go` files pin this.
- **Upstream-tracking rule**: any observable divergence from Java DSS (different verdict,
  bytes, or error) needs a `// DIVERGENCE, deliberate:` comment naming the upstream method
  and why, plus a PORTING.md / DESIGN.md entry. Silent drift is a bug even when the new
  behaviour is "better". Security-motivated refusals of input upstream merely warns about
  are allowed under this rule (reference: `internal/pdf/crypt.go`).
- **Never remove or rewrite a `// Ported from` header**, even when refactoring the file —
  it is the LGPL attribution record.
- **Dependencies**: stdlib + `golang.org/x/…` only. Anything else needs tech-lead sign-off
  recorded in PORTING.md.
- Port test *vectors*, not JUnit code; upstream resources go in `testdata/` mirroring the
  upstream path. Heavy fixtures go under repo-root `corpus/`, with a small representative
  subset kept in-module. Test-file suffixes in use: `_kat_test.go` (known-answer),
  `_oracle_test.go` (Java-oracle comparison), `_smoke_test.go`, `_determinism_test.go`.

## Lint policy

`.github/.golangci.yml` is a v2 config. Blocking: `govet` (minus `inline`) and
`staticcheck` with `ST*`, `S1*`, `QF*`, `SA1019`, `SA9003` **disabled on purpose** — they
would restyle Java-verbatim error strings, "simplify" code away from its 1:1 Java mapping,
or flag the legacy algorithms (SHA-1, MD5, PKCS#1 v1.5) a validation library must keep.
Do not "fix" those categories. `revive exported` runs as advisory only.

## Git / PR conventions

- Commits carry a `Signed-off-by:` trailer (DCO) — use `git commit -s`.
- PRs follow `.github/PULL_REQUEST_TEMPLATE.md`, including the "Upstream alignment"
  checkbox (matches Java DSS, or documented divergence).
- `CHANGELOG.md` is maintained; behaviour changes and closed gaps get an entry.
- Rebasing onto a newer upstream DSS is a manual procedure described in `UPSTREAM.md`.

## Documentation site

`docs/` is an MkDocs Material site (`mkdocs.yml`) published by the Pages workflow. The
compatibility claims (`docs/compatibility/{methodology,numbers,known-gaps}.md`) are the
authoritative statement of what is verified; update them when a gap closes or opens.
