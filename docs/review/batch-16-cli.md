# Batch 16 — CLI / examples / harness

## cmd + examples + harness + corpus (unit U46)

**Scope.** `cmd/esig` (11 non-test files / 1,542 L; 12 incl. tests / 2,235 L) + `examples` (10 files / 845 L) + `harness` (5 files / 1,921 L, all `_test.go`) + 3 corpus generators (`corpus/pades/.../gobroad_main.go` 245 L, `corpus/jades/.../gobroad_main.go` 395 L, `corpus/spi/.../synthetic_certificates.go` 357 L) + the supporting network layer both network paths actually use: `dss/spi/client/http/{native_http_data_loader,native_http_data_loader_call,max_size_input_stream}.go`
**Date.** 2026-08-29
**Depth.**
- `cmd/esig/` — **ALL 11 non-test files read in full**, including the RFC 3161 HTTP client `tsa.go` (priority file) and the `tl refresh` loader wrapper `tl.go`; both test files read in full (`main_test.go`, `e2e_test.go` grep-mapped for password patterns).
- `examples/` — **ALL 10 files read in full** (incl. `internal/fixtures`).
- `harness/` — `document_level_oracle_test.go` read in full (first 120 L + comparison loop); remaining 4 files grep-mapped (structure, fixture sourcing, `exec` surface).
- Corpus generators — pades generator read in full; jades generator header + `detachedContent` map + diff against pades; `synthetic_certificates.go` first 150 L.
- `spi/client/http` loader/call/`MaxSizeInputStream` read in full (the actual bytes-on-wire path for `tl refresh` and example 07); cross-checked against Java upstream `NativeHTTPDataLoader.java`.

### Findings

**T46-SEC-001 — Medium — Unbounded HTTP response body (DoS) — `cmd/esig/tl.go:160-164` (with `dss/spi/client/http/native_http_data_loader_call.go:140-146`) —**
```go
l := &httpFileLoader{}
l.native.SetConnectTimeout(tlFetchTimeoutMillis)
l.native.SetReadTimeout(tlFetchTimeoutMillis)
// SetMaxInputSize is never called -> maxInputSize stays 0
```
`newHTTPFileLoader` sets connect/read timeouts but never `SetMaxInputSize`, so `MaxSizeInputStream` is bypassed and `Call()` reads the response body with an unbounded `io.ReadAll` (`utils.ToByteArray`). **Impact:** `esig tl refresh` follows child-TL and pivot-chain URLs *taken from the downloaded LOTL's contents* (network-controlled URL list); a hostile or compromised LOTL/TL endpoint (or a mistyped `-lotl`) can stream arbitrarily large bodies — the 30 s read timeout bounds time but not size, so a fast connection can exhaust memory. **Recommendation:** call `l.native.SetMaxInputSize(...)` in `newHTTPFileLoader` (a LOTL/TL XML is at most a few MB; 10–50 MiB is generous). Note: this is parity-neutral hardening — upstream Java also defaults `maxInputSize = 0` and no production caller caps it either (verified against `NativeHTTPDataLoader.java`), so adopting a cap is a sanctioned, CLI-only divergence under the upstream-tracking rule, not a fix of a port bug.

**T46-SEC-002 — Medium — No HTTPS enforcement for the `-tsa` time-stamp URL — `cmd/esig/sign.go:121-175`, `cmd/esig/extend.go:93-116`, `cmd/esig/tsa.go:38-39` —**
```go
opts.TSPSource = newHTTPTSPSource(tsaURL)   // tsaURL taken verbatim from -tsa
...
return &httpTSPSource{url: url, client: &http.Client{Timeout: 30 * time.Second}}
```
`-tsa http://...` is accepted with no scheme check (test data itself uses `-tsa http://x`). **Impact:** over plaintext HTTP a MITM can answer with a forged `TimeStampResp`/`TimeStampToken`; at level T/LT/LTA the token is *embedded* into the signature without being anchored to any trust anchor at signing time, so a substituted token is indistinguishable from a genuine one to later verifiers unless they independently anchor the TSA. **Recommendation:** reject (or require an explicit opt-in flag for) non-`https://` `-tsa` URLs in `cmdSign`/`cmdExtend`. The `-lotl` flag has the same property, but its default is the https EU endpoint and the fetched LOTL is signature-checked, so it is lower risk; a scheme check there is a bonus, not the core ask.

**T46-STD-001 — Low — Ignored/short `Read` in example — `examples/05-jades-json-payload/main.go:51-54` —**
```go
n, _ := body.Read(buf)
compact := string(buf[:n])
```
A single `Read` may legally return a short read (and the error is discarded); the printed JWS could be truncated with no indication. **Impact:** cosmetic for this in-memory stream, but it is exactly the pattern users copy from examples. **Recommendation:** `io.ReadAll(body)` (or `io.ReadFull`) and handle the error.

**T46-SEC-003 — Low — Unsanitized user input in log lines (log injection) — `cmd/esig/common.go:226`, `cmd/esig/validate.go:180`, `cmd/esig/tl.go:131` (pattern) —**
```go
return nil, fmt.Errorf("opening %s: %w", path, err)
```
File paths, TSA URLs and LOTL URLs are interpolated into stderr/stdout with `%s`/`%v` and no newline sanitization. **Impact:** a filename or URL containing `\n` can forge additional log lines in CI/scrollback (cosmetic; no secret is involved — passwords never reach these paths, see verification notes). **Recommendation:** optional; sanitize (`strings.ReplaceAll(s, "\n", "\\n")`) in the CLI's error helpers if CI-log integrity matters.

**T46-SEC-004 — Info — TSA client redirect behavior — `cmd/esig/tsa.go:39` —**
Default `http.Client` redirect policy (≤10 hops). Go drops the POST body and converts to GET on 301/302/303, so the message imprint is never re-sent to a redirect target; a redirect could still downgrade to `http://`. **Impact:** minimal — the imprint is a hash of content the signer already holds, and the operator chose the initial URL; relevant only in combination with T46-SEC-002. **Recommendation:** if the scheme check lands, also pin the redirect policy (`CheckRedirect` enforcing https) for defense in depth.

**T46-STD-002 — Info — Live-Java exec surface has no timeout — `dss/pades/pades_downstream_cross_validation_test.go:69-96` (pattern in all six `*_downstream_cross_validation_test.go`) —**
```go
javac := exec.Command("javac", javacArgs...)
java := exec.Command("java", javaArgs...)
```
All argument-vector (no shell), `cmd.Dir` fixed to `testdata/crossgen`, inputs are `t.TempDir()` + operator-controlled `DSS_UPSTREAM_HOME` — **no injection surface**. But no `exec.CommandContext`/timeout, so a wedged `javac`/`java`/`go run` wedges `go test`. **Impact:** test-only, opt-in (`DSS_UPSTREAM_HOME`), skipped under `-short`. **Recommendation:** wrap in `exec.CommandContext` with a generous timeout (e.g. 5–10 min).

**T46-SEC-005 — Info — Example credentials are shared test fixtures — `examples/01-sign-pdf-pades/main.go:30`, `examples/02-validate-pdf/main.go:21` (all examples) —**
```go
signer, err := dss.OpenPKCS12(fixtures.Path("signer_rsa.p12"), "testpassword")
```
`testpassword` is the password of the vendored `dss/testdata/signer_rsa.p12` self-signed test key, shared verbatim with the library's own tests and documented in `examples/internal/fixtures` ("a self-signed RSA test key store"). **Impact:** none — no real credential; acceptable and clearly test data. **Recommendation:** none.

**T46-SEC-006 — Info — Corpus generators make no network calls — `corpus/pades/testdata/broadgen/gobroad_main.go:47-49`, `corpus/jades/testdata/broadgen/gobroad_main.go:34-39`, `corpus/spi/testdata/certificate_extensions/generator/synthetic_certificates.go:75-80` —**
The two `gobroad` generators walk *local* fixture directories and dump JSON via the port's own analyzers; the jades one's `detachedContent` entries that look like URLs (`https://nowina.lu/...`) are detached-content *references* only — no corpus Go file imports `net/http` (grep-verified). `synthetic_certificates.go` (`//go:build ignore`) generates synthetic test certificates with a fresh local RSA-2048 key. **Impact:** none. **Recommendation:** none.

### Security verification notes

1. **RFC 3161 HTTP client (`cmd/esig/tsa.go`) — SAFE, with the two findings above.**
   - **TLS:** default `crypto/tls` verification; `InsecureSkipVerify` appears *nowhere* in the module (grep-verified, whole `dss/`). No root-CA pinning (system pool) — correct for a TSA client that must work against arbitrary public TSAs.
   - **Timeouts:** `http.Client{Timeout: 30s}` bounds dial + body read as a whole; no unbounded-wait path.
   - **Size cap:** present — `io.ReadAll(io.LimitReader(body, 1<<20))` (1 MiB); an oversized body is truncated and then fails the strict `asn1.Unmarshal` (fails closed, no silent acceptance).
   - **Redirects:** default policy, body not re-sent on POST→GET conversion (see T46-SEC-004).
   - **Scheme:** not enforced (T46-SEC-002) — the one substantive gap in this client.
   - **SSRF:** not applicable in a meaningful sense — the CLI is local and the operator supplies the URL; no network-controlled URL feeds this client.
   - Nonce: 64-bit from `crypto/rand` — correct per RFC 3161.
2. **File handling — SAFE.**
   - **Passwords:** `resolvePassword` (`common.go:40`) accepts **only** `env:VARNAME`; literal passwords are refused, and the rejection message deliberately never echoes the typed value (pinned by `TestPasswordMustBeEnvForm`). No password is ever written to a file or temp file.
   - **Temp files:** the one `os.CreateTemp` in the module (`document/temp_file_resources_handler.go:44`) creates with Go's 0600 temp-file perms; the re-open is explicit 0600. No world-readable temp key material anywhere.
   - **Output writes:** `Save` → `os.Create` (0666 & ~umask → 0644 under default umask) — all outputs are public artifacts (signed documents, report XML, trust-anchor PEMs), so 0644 is appropriate; nothing secret is persisted.
   - **Symlinks/TOCTOU:** standard `os.Open`/`os.Create` semantics on operator-supplied paths; no check-then-act race, no privilege boundary crossed. Default output name (`signed.Name()`, cwd-relative) is derived from the user's own input path — traversal would require the user to supply it.
3. **Harness Java invocation — NO INJECTION.**
   - `dss/harness/` (the unit's "harness" package) contains **no `os/exec` at all** — it compares against *pre-generated* Java oracle dumps (JSON/JSONL) checked into `testdata/`. Zero process-execution surface.
   - The live-Java path is the per-package `*_downstream_cross_validation_test.go` tests: all `exec.Command` calls use argument vectors (never a shell), fixed `cmd.Dir`, `t.TempDir()` outputs, operator-supplied `DSS_UPSTREAM_HOME`; gated by `DSS_UPSTREAM_HOME` + tool detection + `-short`. Only gap: no exec timeout (T46-STD-002).
4. **Examples credentials — NO HARDCODED CREDENTIALS** (T46-SEC-005): only the shared, documented test-fixture password; no private keys or real secrets in source.
5. **Corpus generators — what they do:** (a) pades `gobroad` — walks a local PDF corpus, dumps per-signature facts (byte ranges, MDP, modification detection, digests) as JSON for parity diffing against `BroadOracle.java`; (b) jades `gobroad` — same for the JAdES corpus, with a transcribed `detachedContent` map mirroring upstream test classes; (c) `synthetic_certificates` — synthesizes edge-case X.509 certs (SAN types, policy constraints, name constraints) for the SPI known-answer corpus. **No network calls in any of them** (T46-SEC-006).
6. **Second network path (`esig tl refresh` + example 07):** both go through `spi/client/http.NativeHTTPDataLoader`. Timeouts are set in both (30 s CLI / 5 s example), TLS is default-verified, redirects are the default GET-following policy (acceptable for document fetches). The **size cap is missing in both** — the CLI instance is the finding (T46-SEC-001); the example's is benign (5 s read timeout bounds it, and it fetches a single well-known endpoint).
7. **CLI hygiene:** documented exit-code contract (0/1/2/3) consistently implemented and test-pinned; `flag.ContinueOnError` with all diagnostics on stderr; no `os.Exit` inside `run()` (testable); `go run`-able examples resolve fixtures via `runtime.Caller`, not cwd.

### Open questions

1. **T46-SEC-002 decision:** should `-tsa` hard-reject `http://`, or accept it with a loud stderr warning? Hard-reject is cleaner; a warning preserves parity-friendly flexibility for testing against a local plain-HTTP TSA (the e2e tests use `http://x`). Needs a maintainer call — it is CLI-only, so no upstream-parity cost.
2. **T46-SEC-001 cap value:** what size ceiling is sane for a LOTL + 27 member-state TLs fetched in one refresh? (A single EU TL is a few MB; 50 MiB per-document seems safe.) Also whether to apply it in `NewNativeHTTPDataLoader`'s *default* (library-wide) or only in the CLI wrapper — the latter keeps the 1:1 port untouched.
3. Whether the `tl refresh` child-TL URL-following should log each URL it fetches (it currently does not) — useful for auditing which endpoints a refresh actually touched; not a security requirement, but cheap.

### Tool log

Run from `dss/` unless noted:
1. `find cmd examples harness -name '*.go' | wc -l` + `find corpus -name '*.go'` — scope enumeration (12 + 10 + 5 files; 3 corpus generators).
2. `grep -rn 'InsecureSkipVerify|exec.Command|os.CreateTemp|ioutil.TempFile|http.Get|os.OpenFile|0o6[0-6][0-6]|symlink|Lstat|password' cmd/ examples/ harness/` — security-pattern sweep (results: env-only password handling; two 0o644 report/cache writes; no skip-verify; no exec in harness).
3. `grep -rn 'DSS_UPSTREAM_HOME' dss/` — located the live-Java exec surface (6 `*_downstream_cross_validation_test.go` + crossgen generators).
4. `grep -rn 'InsecureSkipVerify' dss/**` — **zero matches** in the whole module.
5. `grep -rn 'os/exec|exec.Command' dss/harness/**` — **zero matches** in the harness package.
6. `grep -rn 'http.|net/http' corpus/**` — only Java oracle sources and XML fixtures; no corpus Go generator imports `net/http`.
7. Java cross-check at `/Users/utain/Workspace/esig/dss`: `NativeHTTPDataLoader.java` — `maxInputSize` defaults to 0 (no cap) and the only `setMaxInputSize` callers are its own unit tests → confirms T46-SEC-001 is a faithful port of upstream default behavior, i.e. a parity-neutral hardening opportunity.
8. `gofmt -l cmd/ examples/ harness/` — **clean** (no output).
9. `go vet ./cmd/... ./examples/... ./harness/...` — **clean**.
10. `golangci-lint run --config=../.github/.golangci.yml ./cmd/... ./examples/... ./harness/...` — **0 issues**.
11. `go build ./cmd/... ./examples/... ./harness/...` — **clean**; `go run ./cmd/esig --help` — prints usage, exits 0.
