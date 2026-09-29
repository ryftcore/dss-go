# Review — batch 01a: cms / asn1 / jose

## internal/jose

Scope: all `.go` files in `dss/internal/jose` (13 non-test files). Port of org.jose4j 0.9.6, the
JOSE (RFC 7515/7518) machinery that `dss-jades` uses to parse and verify JWS signatures.
Deliberate conventions from `dss/PORTING.md` (byte-exact JSON escaping, lenient commons-codec
base64, Java `HashMap` iteration order, 1:1 Java port, SA1019 legacy crypto allowed) are
excluded from findings.

### Tool log

```
$ cd dss && go vet ./internal/jose/...
VET_EXIT=0
$ go version
go version go1.27.0 darwin/arm64
```

`go vet` reports no issues.

---

### Findings

#### SEC-001 — `validSignature` cache is not invalidated when the payload or header changes

| | |
|---|---|
| **Severity** | High |
| **Category** | SEC |
| **Location** | `dss/internal/jose/jws.go:285-306` (`VerifySignature`), `dss/internal/jose/jws.go:102-106` (`SetPayload`), `dss/internal/jose/jws.go:115` (`SetPayloadBytes`), `dss/internal/jose/jws.go:127-131` (`SetEncodedPayload`), `dss/internal/jose/headers.go:86-98` (`SetFullHeaderAsJSONString`), `dss/internal/jose/headers.go:104-117` (`SetEncodedHeader`) |
| **Evidence** | `VerifySignature` computes `alg.verify(j.signature, j.key, j.SigningInputBytes())` and caches the result in `j.validSignature` (a `*bool`), then returns the cached value on every subsequent call as long as it is non-nil. The only two methods that clear the cache are `SetKey` (line 215) and `SetSignature` (line 150). Every method that changes the *signing input* leaves the cache untouched: `SetPayload` (clears `encodedPayload`/`hasEncodedPay` only), `SetPayloadBytes`, `SetEncodedPayload`, and both header setters `SetFullHeaderAsJSONString` / `SetEncodedHeader` (which clear `headerValid`/`encodedValid` only). `SigningInputBytes()` reads the *current* `payloadBytes` and `EncodedHeader()`, so the cached verdict is over a different input than the one that is now present. |
| **Impact** | A caller that calls `VerifySignature()` before the payload or protected header is finalised, then changes either, and then calls `VerifySignature()` again, receives the *stale* verdict — a `true` for a signature that no longer matches the payload. This is a textbook stale-cache crypto bug: the signature is over `header || '.' || payload`, and both `header` and `payload` can be mutated after the verdict is cached. In the current `dss/jades` flow the order is always set-header → set-payload → set-key → verify, so the bug is not triggered today; but `jose.JWS` is a reusable public type and the `Set*` methods are all exported. Any code that reuses a `JWS` across payloads (a detached-signature loop, a test harness, or a future DSS change) will silently get a wrong answer. The risk is a false-positive "valid" verdict, which is the dangerous direction. |
| **Recommendation** | Add `j.validSignature = nil` to `SetPayload`, `SetPayloadBytes`, `SetEncodedPayload`, `SetFullHeaderAsJSONString`, and `SetEncodedHeader`. This is a five-line change that closes the gap. Add a comment to `validSignature` stating that *every* setter that touches the signing input (payload, header, key, signature) must clear it. |

#### SEC-002 — HMAC algorithms are unreachable dead code

| | |
|---|---|
| **Severity** | Medium |
| **Category** | SEC |
| **Location** | `dss/internal/jose/sigalg.go:110-124` (`hmacAlgorithm`), `dss/internal/jose/jws.go:213-217` (`SetKey`) |
| **Evidence** | `SetKey(key crypto.PublicKey)` is the only way to install a key. Go's `crypto.PublicKey` is a sealed interface (unexported `public()` method); `[]byte` does not satisfy it, so `SetKey([]byte("secret"))` is a compile error. `hmacAlgorithm.validateVerificationKey` and `hmacAlgorithm.verify` both do `key.([]byte)`, which can never succeed for any value that could have been set through `SetKey`. The `jwk` parser (`jwk.go`) only handles `kty` RSA / EC / OKP, never `oct`. Yet `lookupJWSAlgorithm` (`sigalg.go:61-89`) returns `hmacAlgorithm{…}` for `HS256`, `HS384`, `HS512`. |
| **Impact** | A hostile JWS with `alg: "HS256"` is resolved to `hmacAlgorithm`, then rejected with `ErrKeyMismatch` ("HMAC needs a secret key, got *rsa.PublicKey") rather than `ErrUnsupportedAlgorithm`. The error misleads a developer debugging interop. More importantly, the three HS\* entries in the switch create a false sense of HMAC support: a future reader may assume HMAC verification works and skip an independent check. If the sealed-interface restriction were ever relaxed (e.g. a wrapper type that embeds `[]byte` and implements `public()`), the HMAC path would become live with a raw-secret key, which is a different threat model than the public-key path. |
| **Recommendation** | Either (a) remove `HS256/HS384/HS512` from `lookupJWSAlgorithm` and let them fall through to the `default` → `ErrUnsupportedAlgorithm` branch, with a comment explaining why; or (b) introduce a `SecretKey` type that wraps `[]byte` and implements `crypto.PublicKey` (via an in-package `public()` method), widen `SetKey` to accept `any`, and route the HMAC path through it. Option (a) is lower risk and matches the package's stated scope ("only what dss-jades actually consumes"). |

#### PERF-001 — No input-size bound on the JSON header parse path

| | |
|---|---|
| **Severity** | Medium |
| **Category** | PERF |
| **Location** | `dss/internal/jose/headers.go:104-117` (`SetEncodedHeader`), `dss/internal/jose/parser.go:44-54` (`ParseJSON`), `dss/internal/jose/base64url.go:90-134` (`Base64URLDecode`) |
| **Evidence** | `SetEncodedHeader` calls `Base64URLDecodeToUTF8String(encodedHeader)`, which allocates a `[]byte` of up to `3/4 · len(encodedHeader)` plus 3, then `string(b)`. The result is handed to `ParseJSON`, which tokenises the entire string and builds an in-memory tree (`Object` + `[]any`) with no size cap. `ParseJSON` has no `maxBytes` parameter; `jsonLexer` holds the full input string for the lifetime of the parse. `Base64URLDecode` pre-allocates `len(encoded)*3/4+3` with no upper bound. |
| **Impact** | `dss` is a validation library that processes untrusted input by design. A JWS whose first compact part (the base64url-encoded protected header) is several megabytes — or, in a pathologically constructed document, hundreds of megabytes — causes a proportional memory allocation with no ceiling. If the JWS is embedded in a larger container (CMS, PDF) that has its own size limit, the risk is bounded; but `Headers.SetEncodedHeader` is a public method on an internal package and could be reached through any code path that feeds attacker-controlled text into it. There is no `io.LimitReader` analogue, no `maxHeaderBytes` constant, and no early rejection. |
| **Recommendation** | Add a package-level constant (e.g. `maxEncodedHeaderLen = 1 << 20`, 1 MiB — generous for any real JOSE header) and reject in `SetEncodedHeader` and `SetFullHeaderAsJSONString` when the input exceeds it. This is a one-line guard at the public boundary; the internal parser can remain unbounded for tests that need large fixtures. |

#### SEC-003 — Lenient base64url decode in JWK coordinate parsing

| | |
|---|---|
| **Severity** | Low |
| **Category** | SEC |
| **Location** | `dss/internal/jose/jwk.go:134-142` (`jwkBigInt`), `dss/internal/jose/base64url.go:90-134` (`Base64URLDecode`) |
| **Evidence** | `jwkBigInt` calls `Base64URLDecode(encoded)` and passes the result directly to `big.Int.SetBytes`. `Base64URLDecode` silently skips any byte not in the 64-character alphabet (`base64DecodeTable[c] < 0` → `continue`). A JWK member `x: "AAAA!!!!BBBB"` therefore decodes to the same `[]byte` as `x: "AAAABBBB"`. RFC 7518 §2 requires the value to be a base64url-encoded byte string; non-alphabet characters are not valid. The leniency is a deliberate port of commons-codec's `Base64.decode`, and `DSSJsonUtilsIsBase64UrlEncoded` (`jades/dss_json_utils.go:155-163`) compensates by doing a separate per-byte alphabet check — but that check is on the *JAdES payload* path, not on the *JWK* path. |
| **Impact** | A hostile JWS whose `jwk` header contains interspersed garbage in the coordinate fields (`x`, `y`, `n`, `e`) is accepted by the Go parser where a strict RFC 7518 implementation would reject it. The decoded numeric value is the same as if the garbage were absent, so the resulting public key is well-formed; the risk is a divergence from strict validators, not a cryptographic break. |
| **Recommendation** | In `jwkBigInt`, validate that every byte of `encoded` is in the base64url alphabet (reuse `IsBase64URLCharacter`) before decoding, and return an error otherwise. This is a three-line addition that closes the gap between the JWK path and the JAdES payload path. |

#### PERF-002 — `[]string` and `[]*Object` to `[]any` allocation on every JSON array render

| | |
|---|---|
| **Severity** | Low |
| **Category** | PERF |
| **Location** | `dss/internal/jose/writer.go:73-86` (`writeJSONValue`) |
| **Evidence** | ```go case []string: items := make([]any, len(v)) for i, s := range v { items[i] = s } writeJSONArray(sb, items) ``` The same pattern repeats for `[]*Object`. Every `[]string` in the JSON tree (the `crit` header, the `docRefs` array, any list of strings) triggers a heap allocation of a new `[]any` that is immediately consumed and discarded. For a header with a 10-element `crit` array, this is a 160-byte allocation per `JSON()` call. |
| **Impact** | Minor in isolation. The `JSON()` function is called on every `FullHeaderAsJSONString()` miss and on every `DSSJsonUtilsToBase64UrlObject`. In a validation run over many signatures the allocations add up, but they are small and short-lived (eligible for the allocator's size-class reuse). |
| **Recommendation** | Add a `writeJSONStringSlice(sb *strings.Builder, items []string)` helper that iterates `items` directly, and a `writeJSONObjectSlice` for `[]*Object`, and dispatch to them from `writeJSONValue` before the generic `writeJSONArray` path. Eliminates the intermediate `[]any` entirely. |

#### STD-001 — `CheckCritOverride` disables the base-library safety net

| | |
|---|---|
| **Severity** | Info |
| **Category** | STD |
| **Location** | `dss/internal/jose/jws.go:43-52` (field), `dss/jades/jws.go:45-50` (`NewJWS`) |
| **Evidence** | `jades.NewJWS()` sets `CheckCritOverride = func(*jose.JWS) error { return nil }`, which makes `JWS.CheckCrit()` a no-op for every JAdES signature. The `crit` check is then performed only by `dssJsonUtilsCriticalHeaderExceptions` in `jades/dss_json_utils.go`, which is a *negative* list (headers that must NOT appear in `crit`), not a positive allow-list. The base library's `CheckCrit` (which rejects any `crit` value the caller has not declared as known) is the positive allow-list, and it is disabled. |
| **Impact** | If a JAdES signature carries `crit: ["someUnknownHeader"]`, the base library would reject it; with the override, it is accepted. The DSS layer's negative list catches `crit: ["alg"]`, `crit: ["jwk"]`, etc., but an unknown header that is not in the negative list passes through. This is a faithful port of upstream DSS's `checkCrit()` override, so it is not a bug — but it means the Go port has a *weaker* `crit` posture than the base library was designed to have, and the only protection is the negative list. |
| **Recommendation** | No code change required (deliberate port). Consider adding a comment in `jades.NewJWS()` cross-referencing `dssJsonUtilsCriticalHeaderExceptions` so a future reader understands that the positive allow-list has been intentionally replaced by a negative deny-list. |

#### PERF-003 — `EscapeJSONString` performs two full scans when escaping is needed

| | |
|---|---|
| **Severity** | Info |
| **Category** | PERF |
| **Location** | `dss/internal/jose/writer.go:171-207` (`EscapeJSONString`, `needsJSONEscaping`) |
| **Evidence** | `needsJSONEscaping` scans the entire string rune-by-rune to determine whether any character needs escaping. If it returns `true`, `EscapeJSONString` scans the string again to build the escaped output. For strings that do need escaping (any non-ASCII character in the U+007F–U+009F or U+2000–U+20FF range, or any of the seven control characters), the total work is 2× the rune count. |
| **Impact** | Negligible for typical JOSE header names (all ASCII, no escapes → single scan). Slightly wasteful for a header value that contains a single escaped character among many (e.g. a URL with a `&`). |
| **Recommendation** | Single-pass implementation: iterate once, writing to a `strings.Builder` that is only committed if any escaping occurred (buffer in a small stack array and flush on first escape). Alternatively, accept the double-scan as a clarity trade-off; the headers are small. |

#### STD-002 — `Headers.Put` does not return the previous value, unlike `Object.Put`

| | |
|---|---|
| **Severity** | Info |
| **Category** | STD |
| **Location** | `dss/internal/jose/headers.go:52-56` (`Headers.Put`), `dss/internal/jose/object.go:70-79` (`Object.Put`) |
| **Evidence** | `Object.Put(key, value) any` returns the previous value (Java `Map.put` contract). `Headers.Put(name, value)` returns nothing. A caller who needs the old value must call `Value(name)` before `Put`, creating a TOCTOU-style read-then-write that is awkward in Go (no atomicity, though there is no concurrency here). |
| **Impact** | Minor API inconsistency within the same package. No current call site needs the previous value through `Headers.Put`. |
| **Recommendation** | Either add a return value to `Headers.Put` (breaking change for existing callers — there are none in the repo) or document the asymmetry in the `Headers.Put` comment. |

---

### Summary

| Severity | Count |
|---|---|
| Critical | 0 |
| High | 1 |
| Medium | 2 |
| Low | 2 |
| Info | 3 |
| **Total** | **8** |

### Top 5 findings

1. **SEC-001 (High)** — The `validSignature` cache is cleared only by `SetKey`/`SetSignature`, not by `SetPayload`, `SetPayloadBytes`, `SetEncodedPayload`, or the header setters. A caller that verifies, then changes the payload or header, and verifies again gets a *stale* verdict — a false-positive "valid" for a signature that no longer matches. Five-line fix: clear the cache in every signing-input setter.
2. **SEC-002 (Medium)** — `HS256/HS384/HS512` are in `lookupJWSAlgorithm` but unreachable: `[]byte` cannot satisfy the sealed `crypto.PublicKey` interface, so no HMAC key can ever be installed. Misleading error path; false sense of support.
3. **PERF-001 (Medium)** — `SetEncodedHeader` → `Base64URLDecode` → `ParseJSON` has no input-size bound. A multi-megabyte base64url header causes proportional unbounded memory allocation. Add a `maxEncodedHeaderLen` guard at the public boundary.
4. **SEC-003 (Low)** — `jwkBigInt` decodes JWK coordinates through the lenient `Base64URLDecode`, which silently skips non-alphabet bytes. A hostile JWK with `x: "AAAA!!!!"` decodes identically to `x: "AAAA"`. Validate the alphabet before decoding.
5. **PERF-002 (Low)** — `writeJSONValue` allocates an intermediate `[]any` for every `[]string` and `[]*Object` in the JSON tree before calling `writeJSONArray`. Add direct-slice render helpers.

---

## internal/asn1ber + internal/eccurve

**Scope.** Hand-rolled BER/DER/DL ASN.1 engine (`internal/asn1ber`, 10 files) and the ECC curve/cert layer it underpins (`internal/eccurve`, 5 files). Both replace BouncyCastle machinery per `PORTING.md`.

**Files read (all 15, in full).**
- `internal/asn1ber`: `doc.go`(18), `algorithm_identifier.go`(87), `element.go`(444), `encode.go`(105), `encoding_kat_test.go`(141), `issuer_serial.go`(65), `string.go`(112), `string_test.go`(53), `tags.go`(49), `time.go`(105).
- `internal/eccurve`: `brainpool.go`(197), `certificate.go`(316), `certificate_test.go`(282), `eccurve_test.go`(155), `weierstrass.go`(285).

**Budget note:** the only file over the 400-line budget is `element.go` (444); it was read in full rather than truncated, so **no file was skipped or sampled**. `eccurve/testdata/` (4 `.der` fixtures) was inspected by name/size only, not as code.

**Callers (severity grounding).** `asn1ber.Parse` is the entry for every CMS/CAdES/PAdES/XAdES signature parse (`cmscore/content_info.go:75`, `cmscore/attributes.go:48`, `cmscore/encoding.go:138`) — i.e. attacker-supplied embedded signature bytes. `eccurve.ParseCertificate` is on every cert-load path (`spi/dss_utils.go:406,475`, `spi/dss_revocation_utils.go:402`, `token/key_store_signature_token_connection.go:191`, `spi/validation/key_entity_tsp_source.go:246`). No production (non-test) code signs with the custom curves (`grep ecdsa.Sign|GenerateKey` outside `_test.go` is empty).

### Findings

**SEC lens — 3 findings.**

#### SEC-001
- **Severity:** High
- **Category:** SEC
- **Location:** `dss/internal/asn1ber/element.go:432` (definite-length child loop; indefinite at `:395`)
- **Evidence:**
  ```go
  for len(rest) > 0 {
      child, remaining, err := Parse(rest)   // recurses once per constructed level
      ...
      rest = remaining
  }
  ```
- **Impact:** `Parse` recurses once per constructed element with no depth bound. I reproduced a full process crash: a valid ~5 MB chain of 1,000,000 nested SEQUENCEs parses OK, but a 5,000,000-deep chain aborts the process with `runtime: goroutine stack exceeds 1000000000-byte limit` — a panic `recover()` cannot catch (confirmed: running `Parse` in a goroutine with a deferred `recover()` did not survive it). An attacker who supplies a deeply-nested CMS/ContentInfo inside a PDF or ASiC (attacker-controlled bytes) can crash the whole validating process with a few MB of input. This is a **port-introduced** DoS: BouncyCastle's `ASN1InputStream` parses iteratively and does not recurse per element, so Java DSS is immune.
- **Recommendation:** Bound the nesting depth in `Parse` (e.g. reject > 200 constructed levels — real X.509/CMS nesting is far below this), or convert the child loop to an explicit stack so depth is bounded by heap, not the 1 GB goroutine stack. Add a regression test feeding a ~1M-deep constructed chain and asserting a clean `error`, not a crash.

#### SEC-002
- **Severity:** Medium
- **Category:** SEC
- **Location:** `dss/internal/eccurve/weierstrass.go:222` (`ScalarMult`); invariant claimed at `weierstrass.go:15-19`
- **Evidence:**
  ```go
  // ScalarMult ... plain left-to-right double-and-add ... NOT constant-time, and only
  // ever driven by public scalars.
  func (c *weierstrassCurve) ScalarMult(bx, by *big.Int, k []byte) (...)
  ```
- **Impact:** The `math/big` math is not constant-time (branch on `b&0x80`, variable `big.Int` allocations), so a secret scalar in `k` leaks timing/alloc signals. The "only ever driven by public scalars" safety argument is a comment only — nothing enforces it — and is **contradicted by the package's own tests**, which sign with `ecdsa.SignASN1` on a Brainpool key (`eccurve_test.go:116`) and generate keys (`certificate_test.go:110`); both route secret scalars (nonce `k`, private `d`) through this `ScalarMult`. Production code does not sign with these curves (verified), so current exposure is limited, but `CurveForOID` is a public API and any consumer that signs/does ECDH on a returned curve is exposed.
- **Recommendation:** Make the "verify-only" invariant unbreakable (document + guard, and/or don't expose a sign-capable `elliptic.Curve`), or accept the non-constant-time cost explicitly for sign/ECDH. At minimum, fix the package doc so it does not assert a guarantee the in-repo tests already violate.

#### SEC-003
- **Severity:** Medium
- **Category:** SEC
- **Location:** `dss/internal/eccurve/certificate.go:55-66` (`init`)
- **Evidence:**
  ```go
  func init() {
      ...
      os.Setenv("GODEBUG", godebug+"x509negativeserial=1")
  }
  ```
- **Impact:** Importing `eccurve` flips the process-wide `crypto/x509` behavior (`x509negativeserial=1`) for **every** `x509.ParseCertificate` call in the process, not just the Brainpool/RSA-negative paths. The *intent* (accept real-world certs with unpadded negative-serial INTEGERs, matching BouncyCastle) is legitimate and fixture-backed, but the *mechanism* is a global stdlib side effect triggered by import, and it is not tagged with the `// DIVERGENCE, deliberate:` convention `PORTING.md` requires for behavioral changes. Any unrelated code in the same binary that relies on Go's default strictness silently loses it.
- **Recommendation:** Keep the accommodation but (a) tag it `// DIVERGENCE, deliberate:` naming the upstream method and motivating fixture (PAdES-LT.pdf, per the existing comment) plus a DESIGN-doc entry, or (b) scope the leniency to `ParseCertificate` (decode the serial from the already-parsed TBS) instead of mutating the global GODEBUG.

**PERF lens — 1 finding.**

#### PERF-001
- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `dss/internal/asn1ber/string.go:26-45` (escape pass in `ValueToString`)
- **Evidence:**
  ```go
  case ',', '"', '\\', '+', '=', '<', '>', ';':
      buffer = append(buffer[:index], append([]rune{'\\'}, buffer[index:]...)...)
      index++
  ```
- **Impact:** Each escape inserts a rune by copying the entire tail of `buffer`, so a value with `k` escape-worthy characters costs O(k·n) for a string of length `n`. `ValueToString` renders X.500 attribute values (RFC 4514 DNs) that can be long and comma-heavy, so a large or hostile DN degrades quadratically. (`Element.Octets`/`DERContent`/`berEncoded` likewise re-`append` accumulated slices on deep constructed OCTET/BIT strings — amortized, but re-encodes per level.)
- **Recommendation:** Build the escaped output in one pass into a pre-grown `strings.Builder` (count specials first, `Grow(n + specialCount)`), instead of in-place tail copies.

**STD lens — 5 findings.**

#### STD-001
- **Severity:** Low
- **Category:** STD
- **Location:** `dss/internal/asn1ber/encode.go:55-69` (`EncodeOID`, `EncodeInteger`)
- **Evidence:**
  ```go
  func EncodeOID(oid asn1.ObjectIdentifier) []byte {
      encoded, err := asn1.Marshal(oid)
      if err != nil { return nil }   // nil indistinguishable from a valid nil
      return encoded
  }
  ```
- **Impact:** Returning `nil` on marshal failure is indistinguishable from a caller-supplied nil and, in `AlgorithmIdentifier.DER`/`IssuerSerial.DER`, a `nil` `EncodeOID`/`EncodeInteger` result is appended as a zero-length field — a silent mis-encode rather than a loud error.
- **Recommendation:** Return `([]byte, error)` and propagate; have `AlgorithmIdentifier.DER`/`IssuerSerial.DER` return an error on a failed sub-encode.

#### STD-002
- **Severity:** Low
- **Category:** STD
- **Location:** `dss/internal/asn1ber/encode.go:74-89` (`OIDFromString`)
- **Evidence:**
  ```go
  component, err := strconv.Atoi(part)
  if err != nil { return nil, fmt.Errorf("string %s not an OID", value) }
  ```
- **Impact:** `Atoi` caps each OID arc at int32, but `asn1.ObjectIdentifier` is `[]int` (64-bit) and `asn1.Marshal` encodes larger arcs — so an otherwise-valid 64-bit arc is rejected here. The error also discards the failing part and the original `err`.
- **Recommendation:** Use `strconv.ParseInt(part, 10, 64)` and wrap the failing part: `fmt.Errorf("OID arc %q: %w", part, err)`.

#### STD-003
- **Severity:** Low
- **Category:** STD
- **Location:** `dss/internal/eccurve/certificate_test.go:104-116`, `dss/internal/eccurve/eccurve_test.go:105-120`
- **Evidence:**
  ```go
  priv := &ecdsa.PrivateKey{PublicKey: *publicKey}   // D == nil, never actually signed with
  other, err := ecdsa.GenerateKey(publicKey.Curve, newDeterministicReader())
  signature, err := ecdsa.SignASN1(newDeterministicReader(), other, digest[:])
  ```
- **Impact:** The tests sign and generate keys on the Brainpool curves (via the non-constant-time `ScalarMult`), which contradicts the package's "no key generation, signing, or ECDH" security note and makes the tests the **only** in-repo sign path. `newDeterministicReader` as a nonce source is a footgun if copy-pasted into production.
- **Recommendation:** Keep the interop coverage but add a comment that signing here is test-only and intentionally exercises the public-key path; or verify against a pre-computed KAT signature instead of signing in-test.

#### STD-004
- **Severity:** Low
- **Category:** STD
- **Location:** `dss/internal/asn1ber/string.go:90-111` (`hexUpper`, `hexLower`)
- **Evidence:**
  ```go
  const digits = "0123456789ABCDEF"
  var builder strings.Builder
  for _, b := range data { builder.WriteByte(digits[b>>4]); builder.WriteByte(digits[b&0x0F]) }
  ```
- **Impact:** Hand-rolls hex that `encoding/hex` already provides (`hex.EncodeToString` for the lowercase variant; `strings.ToUpper` for the uppercase one). Two near-identical loops.
- **Recommendation:** Replace with `hex.EncodeToString` (and `strings.ToUpper` where upper-case is required) to drop the duplicated logic.

#### STD-005
- **Severity:** Info
- **Category:** STD
- **Location:** `dss/internal/eccurve/weierstrass.go:260-264` (`hexInt`)
- **Evidence:**
  ```go
  value, ok := new(big.Int).SetString(s, 16)
  if !ok { panic("eccurve: malformed curve constant " + s) }
  ```
- **Impact:** A `panic` for a malformed curve constant is a landmine if a future contributor adds a bad hex constant — it would crash at `initRegistry` time. Safe today because `TestBrainpoolCurveConstants` re-derives every constant (a typo fails the test, not the runtime).
- **Recommendation:** Acceptable given the test guard; noted so the `panic` is not mistaken for a deliberate runtime guard. No change required.

**Positives (no finding, verified):** length/integer fields are overflow-safe (`count > 4` and `cursor+count > len(input)` guards at `element.go:413-421`; tag overflow guard at `:366`); indefinite length is rejected on primitive elements (`:389`) and bounded by the input (`:404`); `IsOnCurve` range-checks coordinates and reduces correctly (`weierstrass.go:41-59`); point-at-infinity is handled in `affineFromJacobian`/`addJacobian`/`doubleJacobian` (`:61-88`, `:113-117`); all 14 Brainpool parameters are re-derived by `TestBrainpoolCurveConstants` (G on-curve, n·G=∞, (n−1)·G=−G, `Double(G)==Add(G,G)`, 7·G two ways) and spot-checked against RFC 5639; BER/DER/DL re-encodings are pinned by a BouncyCastle-oracle KAT (`encoding_kat_test.go`).

### Tool log
(Run from `dss/`.)

| Command | Result |
|---|---|
| `gofmt -l internal/asn1ber internal/eccurve` | clean (no files listed), exit 0 |
| `go vet ./internal/asn1ber/... ./internal/eccurve/...` | clean, exit 0 |
| `golangci-lint run --config=../.github/.golangci.yml ./internal/asn1ber/... ./internal/eccurve/...` | `0 issues`, exit 0 (golangci-lint available at `/opt/homebrew/bin/golangci-lint`) |
| `go test ./internal/asn1ber/... ./internal/eccurve/...` | `ok` for both packages |
| depth probe (throwaway test, removed after use) | 1,000,000-deep nested SEQUENCEs → parse OK; 5,000,000-deep → `runtime: goroutine stack exceeds 1000000000-byte limit` (unrecoverable; not caught by `recover`) |

### Open questions
1. **Is the SEC-001 DoS acceptable for an in-process validation library?** Upstream (BouncyCastle) is iterative and immune; the recursive port is the regression. If a depth cap is added, what is the safe maximum (real X.509/CMS nesting is < ~50 levels)?
2. **Should `ParseCertificate`'s global GODEBUG mutation (SEC-003) be scoped to the function** rather than set process-wide, given it changes `crypto/x509` behavior for unrelated code in the same binary?
3. **Is the sign-path on custom curves (SEC-002) intended to stay test-only?** If so, the package API could be tightened so the non-constant-time `ScalarMult` is not reachable from production code paths.

---

## internal/cmscore

Scope: all `.go` files in `dss/internal/cmscore` (20 non-test + 6 test files, ~4735 lines). RFC 5652
(CMS) and RFC 3161 (time-stamp) engine replacing BouncyCastle's `org.bouncycastle.cms.*` /
`org.bouncycastle.tsp.*`. Parse-only for RFC 3161; DER build for SignedData. Deliberate
conventions from `dss/PORTING.md` (byte-exact preserved encodings, `expectSizedSequence`
mirroring BouncyCastle's `getInstance` guards, BER tolerance, 1:1 Java port, `nolint:staticcheck`
on the `subjectKeyIdentifier` tag) are excluded from findings.

### Files read

- Fully: `doc.go`, `oids.go`, `version.go`, `general_name.go`, `encapsulated_content_info.go`,
  `cms.go`, `content_info.go`, `certificate_set.go`, `encoding.go`, `signed_data.go`,
  `attributes.go`, `builder.go`, `revocation_info.go`, `signer_info.go`, `adversarial_test.go`,
  `attributes_test.go`, `version_test.go`, `timestamp_test.go`, `builder_bc_kat_test.go`,
  `bouncycastle_oracle_test.go`.
- Budget note (>400 lines → first 300 + grep): `timestamp.go` (459) — read L1–300 +
  L300–459 (grep confirmed no unreviewed func); `cms_kat_test.go` (619) — read L1–300 +
  L300–619. No file skipped.
- Supporting reads (to verify semantics, outside the package): `internal/asn1ber/element.go`
  (`Content`/`Octets`/`Integer`), `spi/validation/timestamp_token.go` (`timestampTokenSignedAttributeValue`,
  `timestampTokenParseCertID`), `diagnostic/signature_wrapper.go` (`MessageDigest`).
- `testdata/adversarial/` is CMS-only (adv-*/craft-*.p7s); there is no hostile TSTInfo/Accuracy
  corpus — all `.tst`/`.tsr` fixtures are well-formed OpenSSL output.

### Tool log

```
$ cd dss && gofmt -l internal/cmscore
(no output)  GOFMT_EXIT=0
$ go vet ./internal/cmscore/...
VET_EXIT=0
$ golangci-lint run --config=../.github/.golangci.yml ./internal/cmscore/...
0 issues.  LINT_EXIT=0
```

`gofmt`, `go vet`, and `golangci-lint` all clean. One semantic probe was run in a throwaway
test (removed afterwards): `big.Int.Int64()` on a 5000-bit value returns `0` (clamps, does not
panic), confirming the behaviour cited in SEC-001.

---

### Findings

#### STD lens

No findings. Error handling, doc comments, and test quality are strong: `wrapField`/`errors.As`-style
returns throughout, `nolint:staticcheck` only where it mirrors a BouncyCastle tag, and the test
suite is genuinely differential (BouncyCastle oracle `bc-oracle.txt`, byte-exact round-trips,
adversarial `craft-*` rejection). See SEC-002 for one maintainability observation.

#### PERF lens

No findings. The verify/sign hot paths are linear: `derSetOf` is O(n log n) on member count
(one `sort.SliceStable`, members are small); `Octets()` is linear; preserved encodings are
sub-slices of the input (`isSubslice` pinned by tests), so no quadratic string building and no
unbounded caches/maps/slices. No goroutines, channels, or readers in the package, so no leak or
`io.ReadAll` surface.

#### SEC lens

##### SEC-001 — `Accuracy` INTEGERs are silently clamped to 0 instead of erroring

| | |
|---|---|
| **Severity** | Medium |
| **Category** | SEC |
| **Location** | `dss/internal/cmscore/timestamp.go:312-321` (`accuracyFromElement`) |
| **Evidence** | ```go case child.IsUniversal(asn1ber.TagInteger): value := int(child.Integer().Int64()); accuracy.Seconds = &value case child.IsContextSpecific(0): value := int(child.Integer().Int64()); accuracy.Millis = &value case child.IsContextSpecific(1): value := int(child.Integer().Int64()); accuracy.Micros = &value``` `Int64()` on an out-of-range `big.Int` returns `0` (verified: `2^5000.Int64() == 0`); there is no `IsInt64()` guard. Contrast the same file's sibling fields — `TSTInfoFromElement` uses `expectSmallInteger(children[0], ...)` for `version` and `expectInteger` for `serialNumber` — and `encoding.go:expectSmallInteger`, which *errors* on `!value.IsInt64()`. |
| **Impact** | A hostile or buggy TSA that encodes `Accuracy.seconds` as a value beyond 2^63 (legal BER, no range restriction) parses to `Seconds = &0` with no error, rather than being refused. The RFC 3161 `Accuracy` is advisory and is used to widen the tolerance window when checking `genTime` against the local clock, so a clamped-to-0 accuracy tightens the window (the safe direction — a false *rejection*, not a false accept); but it is a behavioural divergence from Java BouncyCastle (which reads the full `BigInteger`) and breaks the package's own invariant that out-of-range small integers are an error. The same input is handled correctly one function over (`version`), which is the inconsistency a reviewer must catch. |
| **Recommendation** | Route all three `Accuracy` components through `expectSmallInteger(child, "Accuracy.<field>")` (or add an `IsInt64()` guard) so an out-of-range value is a parse error, consistent with `expectSmallInteger` and with `TSTInfoFromElement`. Add one crafted fixture (e.g. `craft-accuracy-hugeint.tst`) to `testdata/adversarial/` and to the BouncyCastle oracle to pin the rejection. |

##### SEC-002 — `Attributes.Get` returns the *first* match; the verify path correctly uses `GetAll`, but the asymmetry is a foot-gun

| | |
|---|---|
| **Severity** | Info |
| **Category** | STD |
| **Location** | `dss/internal/cmscore/attributes.go:120-139` (`Get` / `GetAll`); consumers `dss/spi/validation/timestamp_token.go:192-208` (`timestampTokenSignedAttributeValue`) vs `dss/cms/cms_integration_test.go:135`, `dss/cms/cms_signed_attribute_table_generator.go:46-49` |
| **Evidence** | `Attributes.Get(type)` returns the *first* attribute of a type (`for … if attribute.Type.Equal(attrType) { return attribute }`); `Attributes.GetAll(type)` returns *all*. This mirrors BouncyCastle's `AttributeTable#get` (first) vs `#getAll`. The RFC 3161 verify path is correct — `timestampTokenSignedAttributeValue` uses `GetAll` and *rejects* `len(attributes) > 1` ("MUST NOT include multiple instances of the message-digest attribute"). But other call sites use `Get` and then read `Values[0]` directly (`cms_integration_test.go:140`, `cms_certificate_source.go:396`, `timestamp_token.go:116,139`). |
| **Impact** | A document that carries two `message-digest` (or `contentType`, or `signing-certificate`) attributes is *rejected* on the RFC 3161 path (good) but, on any path that uses `Get` + `Values[0]`, the *second* (attacker-chosen) attribute is silently ignored and the *first* is trusted. No current verify path in the repo relies on the `Get`+`Values[0]` pattern for a security decision, so there is no false-accept today — but the two accessors look interchangeable and the `Get`-based pattern is the one a future maintainer is most likely to copy into a new check. |
| **Recommendation** | No code change required (faithful to BouncyCastle). Add a one-line warning to the `Attributes.Get` doc comment: "Returns the FIRST attribute of this type. In a verification context use `GetAll` and reject duplicates, as `timestampTokenSignedAttributeValue` does — a duplicate attribute is a rejection signal, not something to take the first of." Optionally add a `GetUnique(type) (Attribute, bool)` helper that errors/returns `false` on duplicates, so new verify code has a single obvious correct call. |

---

### Summary

| Severity | Count |
|---|---|
| Critical | 0 |
| High | 0 |
| Medium | 1 |
| Info | 1 |
| **Total** | **2** |

### Top findings

1. **SEC-001 (Medium)** — `accuracyFromElement` reads `Accuracy.seconds/millis/micros` with `int(child.Integer().Int64())`, which silently clamps an out-of-range BER INTEGER to `0` instead of erroring. Inconsistent with the same file's `expectSmallInteger` (used one function over for `version`) and a divergence from Java BouncyCastle. Route through `expectSmallInteger` and add a crafted fixture.
2. **SEC-002 (Info)** — `Attributes.Get` returns the *first* attribute of a type while `GetAll` returns all. The RFC 3161 verify path correctly uses `GetAll` and rejects duplicates, but other call sites use `Get` + `Values[0]`, a pattern a future maintainer could copy into a security check. Add a doc warning (and optionally a `GetUnique` helper).

### Open questions

1. **TSTInfo `version` is not validated** — `TSTInfoFromElement` reads `version` (`expectSmallInteger`) but never checks it is `v1(1)`, which RFC 3161 requires. All fixtures are v1. Is the omission deliberate (BouncyCastle also reads without enforcing) or should a non-v1 token be refused? Not reported as a finding because no BouncyCastle-oracle fixture exercises it and it is not a false-accept, but it is a gap in the "TST chain/time validation" lens worth confirming.
2. **Empty `signerInfos` set** — `SignedDataFromElement` enforces the field is *present* (errors if `index >= len(children)`) but not that it holds ≥1 member, whereas RFC 5652 mandates `SIZE (1..MAX)`. BouncyCastle's `SignedData.getInstance` accepts an empty set, so the port matches; flagging only to confirm this is the intended parity rather than an oversight.
3. **No hostile TSTInfo/Accuracy corpus** — `testdata/adversarial/` is CMS-only; every `.tst`/`.tsr` fixture is well-formed OpenSSL output. The `accuracyFromElement` (SEC-001) and `genTime`/`ordering`/`nonce` parse paths have no adversarial coverage. Adding 2–3 crafted RFC 3161 fixtures to the BouncyCastle oracle would close the gap.
