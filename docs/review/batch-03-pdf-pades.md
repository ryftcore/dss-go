# Review — unit U11a: internal/pdf (core parsing half)

**Batch:** 03 (internal/pdf, depth = **deep-budgeted**)
**Unit:** U11a — **core parsing half** of the native PDF engine (parser, lexer,
xref, object streams, filters, object model, revisions, DSS dict, SASLprep,
PDFDocEncoding, errors, doc).
**Scope (non-test, in-scope):** `document.go` 1267 · `xref.go` 711 ·
`lexer.go` 564 · `filter.go` 562 · `parser.go` 493 · `object.go` 353 ·
`objstm.go` 166 · `revision.go` 160 · `saslprep.go` 186 · `pdfdocencoding.go`
120 · `dss.go` 119 · `errors.go` 96 · `doc.go` 50.
**Skimmed (quality/coverage only):** `fuzz_test.go` (full — hostile baseline),
`mustreject_test.go` (full — hostile baseline), and `codeql_hardening_test.go`,
`parser_test.go`, `xref_test.go`, `lexer_test.go`, `filter_test.go`,
`document_test.go`, `kat_test.go`, `roundtrip_test.go`, `objectmodel_test.go`,
`revision_test.go`, `objstm_test.go`, `dss_test.go`, `corpuswrite_test.go`,
`zz_dump_writer_test.go` (signatures/structure only).
**Date:** 2026-08-23
**Depth:** deep, **budgeted** (the prior full-package attempt overflowed
context; this pass obeys the read budget below).
**Out of scope (not re-reported):** `crypt.go` (RC4/AES, password check) and
`writer.go` / `incremental.go` / `sign.go` / `repair.go` / `xrefwrite.go` are
the *writer/encryption* half and are not part of U11a. Deliberate conventions
from `dss/PORTING.md` and `internal/pdf/DESIGN.md` are not findings: 1:1
pdfbox-3.0.7 port, the enumerated leniency rules (H/O/S/T/X), the deliberate
*refusals* where pdfbox recovers, byte-identical enum/wire values, `SA1019`
legacy crypto allowed, the single allowed `golang.org/x/text` import
(`saslprep.go`). Known gaps in `docs/compatibility/known-gaps.md` (visible
signature appearances not supported; in-memory-only streaming) are not
re-reported.

---

## Read budget (per file — full vs budgeted)

Files > 400 lines were read as: first 250 lines → **all** `func` signatures
grep'd → **only** the security-relevant regions (parse, xref, offset, length,
stream, decompress, filter, flate, objstm, trailer, indirect, ref, resolve)
read in full. Files ≤ 400 lines were read **fully**.

| File | Lines | First 250 | Func map | Security regions read | Net coverage |
|---|---|---|---|---|---|
| `document.go` | 1267 | ✅ | ✅ (68) | `Resolve`/`object`/`loadAt`/`checkPagesDictionary`/`enqueueKids`/`walkField`/`appendSignatureDictionary` (342–530, 724–833, 1083–1231) | **budgeted** (~40% read; accessors `Get*`/`parseDate`/`Page*`/`Annotations`/`AcroForm`/`collectSignatures` signature-only) |
| `xref.go` | 711 | ✅ | ✅ (23) | `buildXRef`/`walkChain`/`parseXrefTableAt`/`parseXrefStreamAt`/`decodeXRefStream`/`mergeSections`/`checkXrefOffsets`/`findObjectKey`/`readSubsectionHeader` | **full** (every fn is xref/offset/trailer = security-relevant) |
| `lexer.go` | 564 | ✅ | ✅ (26) | `readLiteralString`/`checkForEndOfString`/`readHexString`/`decodeHexDigits`/`readNumber`/`parseLenientReal` (272–564) | **budgeted** (~70%; `isWhitespace`/`skipSpaces`/`next`/`peek` read in first-250 block) |
| `filter.go` | 562 | ✅ | ✅ (20) | `Decode`/`decodeOne`/`FlateDecode`/`ApplyPredictor`/`tiffPredictor`/`ascii85Decode`/`runLengthDecode`/`lzwDecode`/`streamFilters` (1–250, 273–562) | **full** (every fn is a decompressor = security-relevant) |
| `parser.go` | 493 | ✅ | ✅ (19) | `parseDirObject`/`objectFromToken`/`maybeRef`/`parseStreamBody`/`validStreamLength`/`scanForEndstream`/`parseIndirectAt`/`parseTrailerDictAt` (1–250, 236–493) | **full** (every fn is parse/stream/ref = security-relevant) |
| `object.go` | 353 | — | — | **full** | full |
| `objstm.go` | 166 | — | — | **full** | full |
| `revision.go` | 160 | — | — | **full** | full |
| `saslprep.go` | 186 | — | — | **full** | full |
| `pdfdocencoding.go` | 120 | — | — | **full** | full |
| `dss.go` | 119 | — | — | **full** | full |
| `errors.go` | 96 | — | — | **full** | full |
| `doc.go` | 50 | — | — | **full** | full |
| `fuzz_test.go` | 411 | — | — | **full** (baseline) | full |
| `mustreject_test.go` | 191 | — | — | **full** (baseline) | full |
| `repair.go` | 317 | — | — | `bf()`/`findXRefNear`/`rebuildXRefFromBruteForce`/`registerObjectStreams` read (in-scope `xref.go` calls them; read to answer the O(n²) question) | budgeted (supporting) |

**Full-read coverage:** 9/13 in-scope files read in full; `xref.go`,
`filter.go`, `parser.go` effectively full (every function is a
security-relevant parse/decode routine); `document.go`, `lexer.go` budgeted to
the hostile-input regions. Zero security-relevant function left unexamined at
least at signature level.

---

## Tool log

Run from `dss/` (module root). Go `go1.27.0 darwin/arm64`.

```
$ gofmt -l internal/pdf
(empty — clean)

$ go vet ./internal/pdf/...
VET_OK            # exit 0, no diagnostics

$ golangci-lint run --config=../.github/.golangci.yml ./internal/pdf/...
0 issues.         # config passed explicitly, not auto-discovered

$ go test ./internal/pdf/ -run '^$' -fuzz FuzzOpenBytes -fuzztime 15s
fuzz: elapsed: 15s, execs: 1018415, new interesting: 189 (total: 319)
PASS              # ok ... internal/pdf  16.7s   (no crash/hang on ~1M execs)
```

All available tools ran and passed. **No tool was unavailable.** Fuzz target
`FuzzOpenBytes` exercised live; `FuzzScanRevisions` exists and is covered by
the mutation sweeps. Note: the fuzz target self-caps input at **1 MiB**
(`if len(data) > 1<<20 { t.Skip }`) and sets `MaxStreamSize: 4 << 20` — see
the fuzz-coverage gap section, because those caps mask the decompression-bomb
class.

---

## Findings

### Summary table

| ID | Severity | Category | Location | Title |
|---|---|---|---|---|
| P11A-SEC-001 | **High** | PERF-MEM / SEC | `filter.go:138-152` | `FlateDecode` inflates into an unbounded `bytes.Buffer` — decompression bomb |
| P11A-SEC-002 | **High** | SEC / PERF-MEM | `filter.go:183-210` | `ApplyPredictor` `make([]byte, rowLen)` panics on hostile `/DecodeParms` |
| P11A-SEC-003 | **High** | SEC | `objstm.go:31-44` | `parseObjStmHeader` `make([]int64, 0, n)` panics on hostile `/N` |
| P11A-PERF-004 | Medium | PERF-MEM | `filter.go:368-523` | LZW / ASCII85 / RunLength decode have no output-size cap |
| P11A-STD-005 | Low | STD | `filter.go:161-181` | `FlateEncode` panics on a `flate.NewWriter` error |
| P11A-STD-006 | Low | STD / PERF | `object.go:175-190` | `Dict.Delete` O(n) splice + full reindex |
| P11A-STD-007 | Low | STD | `objstm.go:166` | `errObjStmCycle` is dead code (defined, never returned) |
| P11A-STD-008 | Info | STD | package-wide | No `context.Context` propagation for a cancellable decode |
| P11A-PERF-005 | Info | PERF-BIGO | `xref.go:511-564` + `repair.go` | Brute-force recovery path is O(filesize + objects×objstm) — not O(n²) on the happy path |

**Counts by severity:** High 3 · Medium 1 · Low 3 · Info 2 = **9 findings**.
**Counts by category:** SEC 3 · PERF 3 · STD 3.

---

### SEC lens

#### P11A-SEC-001 — `FlateDecode` inflates into an unbounded buffer (decompression bomb)

| | |
|---|---|
| **Severity** | **High** |
| **Category** | PERF-MEM / SEC (memory-exhaustion DoS on hostile input in the prod read path) |
| **Location** | `internal/pdf/filter.go:138-152` |

```go
r := flate.NewReader(bytes.NewReader(raw[2:]))
defer r.Close()
var buf bytes.Buffer
_, err := io.Copy(&buf, r)      // no cap on buf growth
...
return buf.Bytes()
```

`raw` (the *encoded* stream) is bounded by `MaxStreamSize` (default
**512 MiB**, `document.go:36`) — but the *decoded* size is never bounded.
`StreamData` (`document.go:724`) calls `Decode` with no size limit, and the
PAdES layer calls `StreamData` on the streams it validates
(`pades/native_pdf_dict.go:214`, `native_pdf_array.go:68`). A small
Flate-encoded stream expands by a large, unbounded factor.

**Impact.** A network-facing signature validator that reads a hostile PDF can
be made to allocate many GB for a few-MB stream (flate expansion of
highly-repetitive data routinely exceeds 100×; an isolated probe inflated
~1.0 MB of `A`-run data to 1.07 GB, a 1031× ratio). This is a memory-exhaustion
DoS in the production validation path — not a false-accept, so High rather than
Critical. The Go fuzzer does **not** catch it: `FuzzOpenBytes` caps input at
1 MiB and there is no assertion on decoded size, so a bomb neither crashes nor
trips a check (see the fuzz-coverage gaps).

**Recommendation.** Wrap the inflate reader in a size-limited copy (e.g. a
counting `io.Reader`, or `io.CopyN` against `len(raw) * maxRatio` where
`maxRatio` is a small constant such as 1024, or a dedicated
`MaxDecodedStreamSize` option defaulting to a few × `MaxStreamSize`). On
overflow, return `ErrLimitExceeded` / a `*FilterError` rather than allocating.
Pin the chosen ratio in `DESIGN.md` and add a KAT seed of a hand-built bomb.

---

#### P11A-SEC-002 — `ApplyPredictor` panics (`makeslice: len out of range`) on hostile `/DecodeParms`

| | |
|---|---|
| **Severity** | **High** |
| **Category** | SEC (panic on hostile input in a production read path) |
| **Location** | `internal/pdf/filter.go:183-210` (`rowLen := (colors*bpc*columns+7)/8`; `prev := make([]byte, rowLen)`; `cur := make([]byte, rowLen)`) |

```go
rowLen := (colors*bpc*columns + 7) / 8
if rowLen <= 0 { return data }
...
prev := make([]byte, rowLen)   // rowLen is attacker-controlled (see below)
cur  := make([]byte, rowLen)
```

`colors`, `bpc`, `columns` come from the stream's `/DecodeParms` via
`intFromDict` (`filter.go:121`), i.e. from `Integer` (int64) values an
attacker writes into a hostile PDF. With `/Predictor 12` and `/Columns 1<<60`
(1·1·2^60 = 2^60, no int64 overflow), `rowLen = 2^57` and
`make([]byte, 2^57)` panics with `makeslice: len out of range`.

**Empirically confirmed** (isolated probe of the exact expression):
`make([]byte, (1*1*(1<<57)+7)/8)` → `PANIC: runtime error: makeslice: len out
of range`.

**Impact.** A panic in the decode path is a crash. `ApplyPredictor` is reached
for **every** `FlateDecode` and `LZWDecode` stream that carries
`/Predictor > 1` (`decodeOne`, `filter.go:72-96`) — a shape an attacker can
force. It is a hard DoS in the production validation path. Not caught by the
general fuzzer (a random 1 MiB blob almost never yields a *consistent*
Flate stream + well-formed `/DecodeParms` dict with overflow-inducing values).

**Recommendation.** Validate `colors`, `bpc`, `columns` before computing
`rowLen` (bound each to a sane maximum, e.g. reject `columns` large enough that
`colors*bpc*columns` overflows, or check the product in `int64` against
`len(data)`), or compute `rowLen` with overflow checks and return a
`*FilterError`/`ErrLimitExceeded` instead of calling `make`. Add a KAT seed of
a Flate stream with a hostile `/DecodeParms`.

---

#### P11A-SEC-003 — `parseObjStmHeader` panics (`makeslice: cap out of range`) on hostile `/N`

| | |
|---|---|
| **Severity** | **High** |
| **Category** | SEC (panic on hostile input in a production read path) |
| **Location** | `internal/pdf/objstm.go:31-44` |

```go
n := dictInt(dict, "N", 0)
first := dictInt(dict, "First", 0)
if n < 0 || first < 0 || first > int64(len(data)) {   // n is NOT bounded above
    return nil, nil, ...
}
...
nums := make([]int64, 0, n)      // n may be a clamped MaxInt64 -> panic
offs := make([]int64, 0, n)
```

`/N` is an `Integer` (int64); a value beyond int64 is **clamped to
`MaxInt64`** by the lexer (`lexer.go:440-450`, `WarnNumberClamped`), so `n` can
be `MaxInt64`. The range check bounds `first` and `n < 0` but never `n` against
any sane maximum. `make([]int64, 0, MaxInt64)` then panics.

**Empirically confirmed:** `make([]int64, 0, 1<<62)` →
`PANIC: runtime error: makeslice: cap out of range`.

**Impact.** A hostile PDF with a type-2 xref entry pointing at an `/ObjStm`
whose `/N` is huge (and `/First 0`) panics `loadObjStm` → `parseObjStmHeader`.
This is in the production object-resolution path (`object` → `objectFromStm`,
`document.go:448`). Hard DoS crash.

**Recommendation.** Bound `n` (e.g. `n > len(data)` or a fixed maximum such as
the number of representable offsets) in the same guard that bounds `first`, and
return the existing "bad /N or /First" error. Add a KAT seed of an ObjStm with
a hostile `/N`.

---

### PERF lens

#### P11A-PERF-004 — LZW / ASCII85 / RunLength decoders have no output-size cap

| | |
|---|---|
| **Severity** | Medium |
| **Category** | PERF-MEM |
| **Location** | `internal/pdf/filter.go:368-523` (`ascii85Decode`, `runLengthDecode`, `lzwDecode`) |

These decoders append to `out` with no bound and are not subject to the
`FlateDecode` treatment above. `ascii85Decode` emits ~4 bytes per 5 input
bytes (~0.8× input — not a bomb, listed for completeness); `runLengthDecode`
expands each run byte up to 256× (bounded by input × 256); `lzwDecode` grows
the table to 4096 entries and resets, but the *decoded* output is unbounded
relative to input.

**Impact.** Bounded by the (capped) encoded input, so far less severe than the
flate bomb (P11A-SEC-001), but the same class: decoded size is never checked
against a limit before being handed up. Medium.

**Recommendation.** Apply the same decoded-size cap from P11A-SEC-001
uniformly in `Decode`/`decodeOne` rather than per-filter, so every decoder is
bounded.

---

#### P11A-PERF-005 — xref parsing is not O(n²); the brute-force *recovery* path is O(filesize + objects×objstm)

| | |
|---|---|
| **Severity** | Info (observation — answers the "O(n²) over object count?" question) |
| **Category** | PERF-BIGO |
| **Location** | `internal/pdf/xref.go:511-564` (`checkXrefOffsets`), `repair.go:29-150` (`bf`), `repair.go:186-230` (`registerObjectStreams`) |

The **happy path** is linear: `parseXrefTableAt` reads each entry once,
`decodeXRefStream` is O(entries), `mergeSections` is O(total entries), and
`checkXrefOffsets` does one bounded `findObjectKey` (a few tokens + a bounded
backward digit scan) per type-1 entry. `sortedXRefKeys` (O(n log n)) is called
a bounded number of times, not per-object. **There is no O(n²) over object
count on the normal path.**

The **recovery path** (only when `startxref`/`/Root`/a type-1 header fails)
is `rebuildFromBruteForce` → `bf()` (a single O(filesize) byte scan, **cached**)
→ `registerObjectStreams` (calls `loadAt` + `objStmNumbers` for **every**
object, i.e. O(objects × objstm-decode)). That is the expensive shape, but it
is a deliberate last resort mirroring pdfbox's `BruteForceParser`, it is
cached, and it is pinned by `repair_test.go`. No action required; recorded so
the O(n²) question has a documented answer.

---

### STD lens

#### P11A-STD-005 — `FlateEncode` panics on a `flate.NewWriter` error

| | |
|---|---|
| **Severity** | Low |
| **Category** | STD |
| **Location** | `internal/pdf/filter.go:161-181` |

`w, err := flate.NewWriter(&buf, 6); if err != nil { panic(...) }`. Level 6 is
always valid, so the branch is effectively dead — but a `panic` on a
constructor-error path is not idiomatic and would mask a real error if the
level constant ever changed. Recommend returning an error (or a
`*FilterError`) consistent with the rest of the package's error taxonomy.

#### P11A-STD-006 — `Dict.Delete` is O(n) splice + full reindex

| | |
|---|---|
| **Severity** | Low |
| **Category** | STD / PERF |
| **Location** | `internal/pdf/object.go:175-190` |

```go
d.keys = append(d.keys[:i], d.keys[i+1:]...)
d.vals = append(d.vals[:i], d.vals[i+1:]...)
delete(d.m, key)
for j := i; j < len(d.keys); j++ { d.m[d.keys[j]] = j }
```

Correct (preserves insertion order, which is the contract), but O(n) on both
the slice shift and the reindex. `Dict` is insertion-ordered by design, so a
lazy-delete (tombstone) or a two-pass rebuild would be the only real fix; not
worth the complexity unless a hot loop calls `Delete` repeatedly. Noted for
completeness.

#### P11A-STD-007 — `errObjStmCycle` is dead code

| | |
|---|---|
| **Severity** | Low |
| **Category** | STD |
| **Location** | `internal/pdf/objstm.go:166` |

`var errObjStmCycle = errors.New("pdf: object stream cycle")` is defined but
never returned or referenced anywhere in the package (grep confirms a single
hit). Object-stream cycles are actually guarded by the `d.objStms` cache
("cache first: a cycle must not recurse forever", `objstm.go:62`). Either wire
the sentinel in or delete it.

#### P11A-STD-008 — No `context.Context` propagation for a cancellable decode

| | |
|---|---|
| **Severity** | Info |
| **Category** | STD |
| **Location** | package-wide |

`Open`, `OpenBytes`, `StreamData`, `Decode`, and the decoders take no
`context.Context`. Combined with the unbounded decodes above
(P11A-SEC-001/002/004), a caller cannot cancel a long/bad decode. Low priority
given the library is in-memory and single-call, but a `ctx` on `StreamData`/
`Decode` would be the natural companion to the size caps.

---

## Lenses with no findings (explicit)

- **Circular-reference / loop detection (SEC):** *no findings.* Indirect-object
  loops are handled: `Document.Resolve` bounds ref chains by `MaxDepth`
  (`document.go:414-437`), `Document.object` uses a `loading` map to prevent
  re-entrant `/Length` cycles (`document.go:455-460`), object-stream loading
  caches before decoding ("a cycle must not recurse forever",
  `objstm.go:62`), and the `/Prev` xref chain has both a `seen` set and a
  `maxSections = 512` cap (`xref.go:139-183`). Object-streams-inside-object-
  streams are refused (`objstm.go:59-61`).
- **SASLprep correctness (SEC):** *no findings.* `saslprep.go` faithfully
  reproduces pdfbox's `SaslPrep.saslPrepQuery` including its deliberate
  departs-from-RFC-3454 truncation of supplementary code points
  (`saslProhibited`, `saslprep.go:73-88`); the NFKC step uses `x/text` (the
  one allowed module), and the Unicode-version skew is documented in the file
  header and `known-gaps.md`. `ErrProhibitedPassword` is correctly
  distinguished from `ErrInvalidPassword` (`errors.go:15-24`).
- **Lexer hostile tokenization (SEC):** *no findings.* Unterminated literal/hex
  strings return what they have (`lexer.go:339-347`, `382-395`); huge
  integers are **clamped** to ±`MaxInt64` with `WarnNumberClamped`, never a
  parse error (`lexer.go:440-450`); the lenient-float repairs are bounded and
  the "longest parseable prefix" fallback is linear in the literal length
  (`lexer.go:475-505`). Integer overflow on offset/length is therefore not a
  panic (it clamps) — the *consequence* of a clamped value flowing into
  `make(...)` is P11A-SEC-003, reported under that site.
- **objstm recursion / count mismatch (SEC):** the recursion guard and
  out-of-range offset checks are correct (`objstm.go:79-118`); the only defect
  is the `/N` cap panic, P11A-SEC-003.
- **xref malformed-entry handling (SEC):** *no findings.* Negative/zero
  offsets rejected (`xref.go:212-224`), generation bounded to 0..65535
  (`xref.go:656-659`, `isValidGeneration`), `/W` array bounded
  (`xref.go:365-386`, total ≤ 20), truncated xref streams keep partial output
  (`xref.go:414-417`), `/Prev` cycle + section cap (`xref.go:139-183`).

---

## Fuzz-coverage gaps

The hostile-input baseline (`fuzz_test.go`, `mustreject_test.go`) is strong:
`FuzzOpenBytes` walks every reader surface the PAdES layer can reach, and two
deterministic 10k/12k-mutant sweeps (`TestMutantsNoPanic`,
`TestMutantsReaderWriterNoPanic`) push corruption through reader **+** a full
signing increment with a recover() net. **However**, three gaps leave the
High findings above undetected:

1. **Decompression-bomb class is invisible to the fuzzer.** `FuzzOpenBytes`
   self-caps input at 1 MiB (`if len(data) > 1<<20 { t.Skip }`) and sets
   `MaxStreamSize: 4 << 20`, and there is **no assertion on decoded size**. A
   flate/LZW/predictor bomb neither crashes (the panics above need a *specific*
   dict shape, not random bytes) nor trips a check, so P11A-SEC-001/002 and
   P11A-PERF-004 are not caught. **Gap: no fuzz seed or test that asserts a
   decoded-output bound / a hand-built bomb.**
2. **No targeted hostile-`/DecodeParms` seed.** The predictor panic
   (P11A-SEC-002) needs a *consistent* Flate/LZW stream + a well-formed
   `/DecodeParms` dict with overflow-inducing `/Columns`. Random mutation
   almost never constructs that shape. **Gap: no KAT/fuzz seed of a
   Flate stream with a malicious `/DecodeParms`.**
3. **No targeted hostile-ObjStm `/N` seed.** The `/N` cap panic
   (P11A-SEC-003) needs a type-2 xref entry → `/ObjStm` with a huge `/N` and
   a valid `/First 0`. **Gap: no KAT/fuzz seed of an ObjStm with a malicious
   `/N`.**
4. **Brute-force recovery path is under-exercised at scale.** The
   `TestBruteForceScanIsBoundedByTheLastEOF` and `TestBruteForce*` tests pin
   behaviour but there is no fuzz target that specifically stresses
   `rebuildFromBruteForce`/`registerObjectStreams` over many objects with
   object streams (the P11A-PERF-005 shape).

---

## Open questions

- **Q1 (P11A-SEC-001/004):** What decoded-size policy does the project want?
  A per-filter cap, a shared `MaxDecodedStreamSize` option, or a ratio bound
  (`len(decoded) ≤ k · len(encoded)`)? This should be a named `DESIGN.md`
  entry and an `Options` field, consistent with the existing
  `MaxObjects`/`MaxDepth`/`MaxStreamSize` guards.
- **Q2 (P11A-SEC-002/003):** Should the two `make(...)` panic sites be fixed by
  clamping the attacker dict values to a documented maximum (recommended), or
  by a shared "decode budget" that both the predictor and the ObjStm header
  respect?
- **Q3:** Is `errObjStmCycle` (P11A-STD-007) intended for a future error path
  (e.g. a `/Type /ObjStm` whose own container is compressed)? If not, delete
  it; if so, wire it in.
- **Q4 (fuzz gaps 1–3):** Can targeted bomb/`/DecodeParms`/`/N` seeds be added
  to `fuzzSeeds` without bloating the corpus (they are tiny hand-built blobs,
  so likely yes) — and should the 1 MiB input cap be raised for the
  decompression-bomb class, or is a dedicated size-bounded assertion the
  cleaner fix?
- **Q5 (P11A-SEC-001):** Confirm the production blast radius — `StreamData` is
  called on the `/VRI /TS`, `/DSS token`, xref- and object-streams
  (per `doc.go`/`filter.go` header). Is any *arbitrary* document stream ever
  decoded during validation, or strictly that set? Either way the attacker
  controls which stream is which, but the exact set determines how many
  independent bombs a single document can carry.

---

*Scope note: this unit (U11a) covers the **core parsing half** only. The
writer/encryption half (`crypt.go`, `writer.go`, `incremental.go`, `sign.go`,
`repair.go`, `xrefwrite.go`) and its own findings (e.g. the non-constant-time
password compare and the AES key-size silent failure previously noted under
U11) are out of scope here and are not re-reported.*

---

## internal/pdf crypt+sign+writer+repair (unit U11b)

**Batch:** 03 · **Date:** 2026-08-23 · **Depth:** deep, budgeted.
**Scope (in-scope, non-test):** `crypt.go` (1051), `writer.go` (552), `sign.go` (539), `incremental.go` (511), `repair.go` (317), `xrefwrite.go` (286).
**SKIM (coverage only):** `cryptfilter_test.go`, `cryptaes256_test.go`, `sign_test.go`, `incremental_test.go`, `writer_test.go`, `writerfixture_test.go`, `xrefwrite_test.go`, `repair_test.go`.
**READ FULLY (define the crypto baseline):** `password_kat_test.go`, `crypt_test.go`.

This is the companion to U11a above, focused on the CRYPTO/SIGN half. It deliberately does **not** re-report the U11a findings (e.g. the non-constant-time password compare and the AES key-size silent failure) or the known gap (visible signature appearances). It cross-checks every security-relevant function against the actual pdfbox 3.0.7 source found on this host (`/tmp/pb307src/.../StandardSecurityHandler.java`, `SecurityHandler.java`) rather than trusting the port's doc comments.

### Files read (full vs budgeted)

| File | Lines | Read |
|---|---|---|
| `crypt.go` | 1051 | **full** (1–250 + 250–398 + 398–760 + 760–1051; every function read) |
| `writer.go` | 552 | **full** (1–250 + 250–552) |
| `sign.go` | 539 | **full** (1–250 + 250–539) |
| `incremental.go` | 511 | **full** (1–250 + 250–511) |
| `repair.go` | 317 | **full** (≤400, read end-to-end) |
| `xrefwrite.go` | 286 | **full** (≤400, read end-to-end) |
| `password_kat_test.go` | — | **full** (baseline) |
| `crypt_test.go` | 306 | **full** (baseline) |
| `cryptaes256_test.go`, `cryptfilter_test.go`, `sign_test.go`, `incremental_test.go`, `writer_test.go`, `xrefwrite_test.go`, `repair_test.go` | — | **skim** (function inventory + target grep only) |

**Every function in the six in-scope files was read in full** — none signature-only. The budgeted-read strategy (first 250 lines → grep all `func ` → targeted region reads) resolved to a full read because the security-relevant functions (`setupEncryption`, `validatePerms`, `checkCryptFilters`, `isUserPassword`, `isOwnerPassword`, `userPasswordFromOwner`, `computeRC4Key`, `computeUserEntry`, `computeKeyRev234`, `computeEncryptionKey`, `hash2AOr256`, `computeHash2B`, `rc4Apply`, `objectKeyFor`, `decryptBytes`, `decryptValue/Dict/Stream`, `encryptForWrite`, `Write`/`WriteIndirect`/`writeStream`/`EncodeName`/`EncodeString`, `patchByteRange`, `SignedData`, `InsertContents`, `ReplaceContents`, `writeXRefTable`/`writeXRefEntry`, `buildXRefStream`, the repair `bf()` family) all sat inside the 250-line windows or were reached by the grep map.

### Verified-correct (no finding) — the high-risk claims, checked against upstream

- **RC4/AES key derivation, byte-for-byte.** `computeRC4Key` (50-round MD5 over `digest[:keyLen]` — the *documented* pdfbox `length`-not-`keyLen` quirk), `computeKeyRev234` (Algorithm 2 incl. the `/R 4 && !encMeta` 4-byte `0xFF` salt), `computeUserEntry` (20-round iterated RC4, `buf[:16]‖padding[:16]`), and `userPasswordFromOwner` (19→0 XOR-ladder) all match `StandardSecurityHandler` verbatim. The `clampKeyLen`/`clampDictInt` guards close the hostile `/Length` → negative/250 MB digest-slice panic **without** changing any well-formed derivation (the default 40-bit = 5-byte path is untouched).
- **`objectKeyFor` byte layout = upstream `calcFinalKey`.** `crypt.go:895–919` writes `num` low-3 (LE) then `gen` low-2 (LE) exactly as `SecurityHandler.calcFinalKey:256–260`, adds the `sAlT` suffix only for AES, truncates to `min(len(key)+5, 16)`, and fast-paths the 32-byte AES-256 file key (Algorithm 1.A). A byte-order slip here would silently mis-decrypt high-numbered objects; there is none.
- **`/UE`/`/OE` exactness.** `computeEncryptionKey` requires exactly 32 bytes (stricter than pdfbox, which accepts any whole block count). The 32-byte requirement is sound: AES-CBC decrypts block-0 independent of later blocks, so a padded `/UE` would still expose the genuine key — the port's hard check makes the `len(h.key)==32 ⇔ R5/R6 unwrap` invariant used by `objectKeyFor` and the postcondition guard sound rather than padding-bypassable.
- **`computeHash2B` termination (the classic R6 infinite-loop bug).** `crypt.go:782` `for round:=0; round<64 || int(e[len(e)-1]) > round-32; round++` — Go's `int(byte)` is *unsigned* (0–255), byte-for-byte the same comparison as Java's `(e[e.length-1] & 0xFF) > round-32` (`StandardSecurityHandler.java:1139`). The loop is provably bounded (round strictly increases; the `||` branch only extends it while the final byte keeps falling below the moving `round-32` threshold). Not a finding.
- **R5/R6 `/Perms` binding (`validatePerms`).** Checks the `'adb'` marker, the little-endian `/P`, and the `/EncryptMetadata` flag — exactly the three fields pdfbox checks. Refusing on mismatch (vs pdfbox's `LOG.warn` + load) is the deliberate `PORTING.md` divergence that defeats the one-byte `/P` edit; correct.
- **Byte-range exclusion (the critical false-accept risk) — correct.** `patchByteRange` (`incremental.go:472–501`) sets `ByteRange = [0, before, before+contentsLen, total-(before+contentsLen)]` and `SignedData()` (`incremental.go:503–511`) returns exactly `Bytes[0:before] ∪ Bytes[after:total]` — i.e. it excludes precisely the `/Contents` `<…>` span the CMS is later written into, and the remainder (trailer/xref/`%%EOF`) is *included* in the signed data. `InsertContents`/`ReplaceContents` bounds are checked against the reserved span, so no out-of-bounds write. No false-accept in the signing path.
- **xref offset tracking — correct.** `Write` (`incremental.go:286–…`) records `w.Pos()` *before* each `WriteIndirect` and uses those captured positions in the xref table/stream, so the emitted offsets equal the actual byte positions in `original‖increment`. `freeHeadEntry` (gen 65535, next-free 0), `xrefRanges` contiguity, and `buildXRefStream`'s `/W [1 n 2]` + `/Predictor 12` all match the reader's `ApplyPredictor` inverse (pinned by `TestPNGUpEncodeRoundTrip`).
- **incremental old-content preservation — correct.** The `§3.1` invariant holds: `orig` is copied verbatim once, at the end; nothing seeks backwards. `Update` refuses to rewrite a stream of an encrypted document (double-encrypt guard), which is the right failure.
- **`repair.go` hostile-input posture — sound.** The brute-force scan is bounded by the last `%%EOF`, parses `num`/`gen` with overflow-checked `ParseInt`, bounds `gen ≤ 65535`, and `searchForTrailer`/`searchForTrailerItems` use depth-limited `parseDictBody` (`parser.go:167` checks `p.depth >= p.maxDepth`) and a 32-token look-ahead in `looksLikeXRefStream`. Failures degrade to `ErrNotPDF` / a warning, not a crash.
- **RC4/AES read-path cipher correctness.** `rc4Apply` is symmetric (RC4 self-inverse), and `decryptBytes`' AES arm (IV = first 16 bytes, `AES/CBC`, lenient PKCS#5 strip) matches pdfbox's `createCipher` (`AES/CBC/PKCS5Padding`) swallowing `BadPaddingException`. The `len(data) < aes.BlockSize → []byte{}` and `len(body)==0 → []byte{}` guards prevent OOB slicing on a truncated stream.
- **`writer.go` encoding determinism.** `EncodeName` (R5 restricted charset), `EncodeString` (R6 hex/literal, octal-free), and `FormatReal`/`javaFloatToString`/`plainDigits`/`plainFromScientific` (Java `Float.toString` + `stripTrailingZeros().toPlainString()`) are pure functions of their input — no map order, clock, or PRNG. `WriteEOL` suppression (R2) and `writeArray`'s every-10th-EOL (R4) match `COSWriter`.

### KAT-coverage notes

The crypto baseline is genuinely strong where it is covered:
- **RC4 (`/V 1`–`/V 4`):** `crypt_test.go` builds a *real* `/Filter /Standard /V 1 /R 2` document, opens it, decrypts strings and streams, and pins the never-decrypt `/Contents` rule (incl. the PDFBOX-4466 `/ByteRange`-without-`/Type` case). `TestRC4KnownAnswer` pins the RFC 6229 vector. `TestObjectKeyDerivation` pins Algorithm 1 lengths and the AES-256 fast path.
- **AES-256 (`/V 5 /R 6`):** `cryptaes256_test.go` is a thorough acceptance-boundary + rejection suite (unidentifiable/disagreeing crypt filters, `/Perms` missing/too-short/too-long, `/P` mismatch, `/UE`/`/OE` length, `/Perms` literal vs indirect, `/P` as a real) built from `buildAES256Doc`.
- **Charset step:** `password_kat_test.go` pins `passwordBytes` against ISO-8859-1/UTF-8/SASLprep and a set of *pdfbox-produced* fixtures (`testdata/password/*.pdf`), asserting the wrong byte-forms are *refused*.
- **`computeHash2B`:** `TestComputeHash2BIsDeterministic` pins termination + 32-byte output + salt-sensitivity.

**The one real gap (P11B-STD-001):** the **`/R 5` (AES-256 with SHA-256, not SHA-512-384/512)** AES-256 path has **zero** independent known-answer coverage anywhere — no builder, no fixture, no corpus file. `hash2AOr256`'s `r==5` arm (plain `SHA-256(pw‖salt‖userKey)`) is only ever exercised with `rev=6`. So a regression in the `r==5` branch — e.g. a wrong salt slice, a copy-paste into the R6 path, or a missing `truncate127` — would not be caught by any test, and it would change the derived key for a class of real-world AES-256 PDFs (R5 is the more common of the two in the wild). This is a coverage gap, not a confirmed bug: the code path is present and reads correctly, and the R5/R6 `computeEncryptionKey` unwrap is shared and *is* exercised (via the R6 builder).

### Findings

#### P11B-STD-001 — `/R 5` (AES-256/SHA-256) key-derivation path has no known-answer pin

| | |
|---|---|
| **Severity** | Low (coverage gap, not a live bug) |
| **Category** | STD |
| **Location** | `dss/internal/pdf/crypt.go:750` (`hash2AOr256` `r==5` arm), `:766` (`adjustUserKey`), `:682` (`computeEncryptionKey` R5 unwrap); absent from `cryptaes256_test.go`, `cryptfilter_test.go`, `password_kat_test.go`, and the `corpus/` oracle. |
| **Evidence** | `hash2AOr256`'s `r==5` arm is `sha256.Sum256(pw‖salt‖userKey)`; every test caller passes `rev = 6` (`cryptaes256_test.go:75 const rev = 6`; `:97,100,104,107`). `grep -rn "rev = 5\|/V 5 /R 5\|aes256_r5" internal/pdf/*_test.go corpus/` → no builder, no fixture. |
| **Impact** | The R5 AES-256 derivation (SHA-256, no `computeHash2B`) is one of two real-world AES-256 shapes and the more common one, yet no test asserts it against an independent vector. A regression that swaps the R5/R6 salt/userKey handling, or drops `truncate127`, would pass the whole suite and silently corrupt key recovery for R5 documents (false-reject: a valid password fails to open the file). The shared R5/R6 `/UE` unwrap is covered, so this is narrowly the *password-hash* half. |
| **Recommendation** | Add one R5 KAT: a `buildAES256Doc` variant with `const rev = 5` (the builder already threads `rev` through `hash2AOr256`), plus ideally a pdfbox-produced `aes256_r5_*.pdf` alongside the existing `aes256_r6_utf8.pdf` (`testdata/password/PasswordFixtures.java` already has the R6 generator). Assert it opens with the user password, `enc.R == 5`, and a wrong password is refused. |

#### P11B-STD-002 — `Writer.WriteRaw` ignores a short `io.Writer` write

| | |
|---|---|
| **Severity** | Low |
| **Category** | STD |
| **Location** | `dss/internal/pdf/writer.go:79–92` (`WriteRaw`). |
| **Evidence** | ```go n, err := w.w.Write(b); w.pos += int64(n); if err != nil { w.err = err }; return w.err ``` — the `n < len(b)` partial-write case is neither checked nor retried, and the bytes actually written are not validated. |
| **Impact** | `io.Writer` permits a partial write (`0 < n < len(b)`) to signal "wrote some, not all." Every real `io.Writer` this is used with (`bytes.Buffer`, `os.File`) writes all-or-nothing, so the bug is latent. But a custom `io.Writer` (a bounded/buffered sink, or a writer that fills a fixed chunk) that returns a short count with `nil` error would cause `pos` to under-advance *and* drop bytes, producing a silently corrupted PDF with wrong xref offsets — the "wrong offset = unopenable PDF" class this file's doc comments are careful about. The contract is stated ("every write is a no-op after [an] error") but a partial write is not an error. |
| **Recommendation** | Loop until all of `b` is written (or `n==0` on a non-nil error), mirroring `io.WriteFull` semantics: `for len(b) > 0 && err == nil { n, err = w.w.Write(b); b = b[n:]; w.pos += int64(n) }`. Cheap, and it makes the "offsets match actual byte positions" invariant hold for *any* `io.Writer`, not just all-or-nothing ones. |

#### P11B-SEC-001 — R5/R6 (AES-256) password check is non-constant-time (refines U11a's password-compare finding on the AES arms)

| | |
|---|---|
| **Severity** | Low |
| **Category** | SEC |
| **Location** | `dss/internal/pdf/crypt.go:584` (`isUserPassword` R5/R6 `bytes.Equal(hash, u[:32])`), `:593` (`isOwnerPassword` R5/R6 `bytes.Equal(hash, o[:32])`). |
| **Evidence** | `hash := hash2AOr256(truncate127(pw), u[32:40], nil, r); return bytes.Equal(hash, u[:32])` — `bytes.Equal` is a short-circuiting byte loop. |
| **Impact** | A timing side-channel on the R5/R6 (AES-256) password check. This is the *same class* as U11a's non-constant-time password-compare finding, but on the **AES-256 arms** (`crypt.go:584,593`), which that finding scoped to the MD5-based R2–R4 arms and did not enumerate. Practical exposure is negligible — the file is public and the derived key is brute-forceable offline — but the textbook-correct primitive for comparing a derived hash is `hmac.Equal`. |
| **Recommendation** | `return hmac.Equal(hash, u[:32])` (and the owner-password site). One-word swap, no behaviour change for a correct password. `hmac.Equal` is constant-time over the *shorter* operand, so comparing the 32-byte `hash` to `u[:32]` is the right shape. If the security owner signs off on `hmac.Equal` generally, fold this into U11a's password-compare fix. |

**Counts by severity:** Low 3 (P11B-STD-001, P11B-STD-002, P11B-SEC-001). **No Critical, High, Medium, or Info findings** in this unit's in-scope set.
**Counts by category:** STD 2 · SEC 1 · PERF 0.

> The two PERF hot paths the brief asked about are clean on inspection: the RC4/AES key stretching (`computeRC4Key`/`computeKeyRev234` 50-round MD5, `computeHash2B` ≤64-round AES/SHA loop) runs **once per open, per password** (owner then user), not per object — per-object work is only the cheap `objectKeyFor` MD5 (or the AES-256 zero-copy fast path) and the single `decryptBytes`. There is no per-string re-derivation and no unbounded buffer growth in these six files (the unbounded-allocation exposure lives in `filter.go`/`lexer.go`, covered under U11a).

### Tool log (U11b)

Run from `dss/` (module root). Go `go1.27.0 darwin/arm64`. All four tools **available**.

```
$ gofmt -l internal/pdf
(empty — clean)

$ go vet ./internal/pdf/...
VET_OK            # exit 0, no diagnostics

$ golangci-lint run --config=../.github/.golangci.yml ./internal/pdf/...
0 issues.         # config passed explicitly, not auto-discovered

$ go test ./internal/pdf/ -run 'TestCrypt|TestEncryption|TestPassword|TestComputeHash2B|TestRC4|TestObjectKey|TestAES|TestXRefStream|TestUnencrypted|TestTruncateOrPad|TestPermissionBits|TestSignature|TestByteRange|TestInsertContents|TestReplaceContents|TestXRef|TestLenient|TestBruteForce|TestPNGUp|TestFormatReal|TestWriteStream|TestDeterminism|TestPrefixPreserved|TestStartXref|TestTrailer|TestFormatDate' -count=1
ok  github.com/ryftcore/dss-go/dss/internal/pdf  0.732s
```

The crypto/signing/writer KAT subset (all of `crypt_test.go`, `password_kat_test.go`, `cryptaes256_test.go`, `cryptfilter_test.go`, plus the `sign`/`incremental`/`writer`/`xrefwrite`/`repair` test functions) is **green** on this host.

### Open questions (U11b)

1. **Should R5 AES-256 get a dedicated KAT (P11B-STD-001)?** Strongly yes — it is the more common of the two AES-256 revisions and the only derivation path in this file with no independent vector. Confirm with the security owner whether an R5 fixture can be produced from the existing `PasswordFixtures.java` (it already emits the R6 one).
2. **Is the `Writer` ever pointed at a non-buffered `io.Writer` in the `pades`/`asic` layers (P11B-STD-002)?** If so the short-write guard is worth adding now; if every caller is `bytes.Buffer`-backed it stays latent. Needs a call-site check outside this package (out of scope here).
3. **Does any corpus document exercise `/R 5`?** The `kat_test.go` oracle sweep (`corpus.tsv`) is the only place it could be pinned today; if no corpus file is R5, the R5 gap (P11B-STD-001) is fully untested and the recommendation is the highest-value single test addition in this unit.
4. **`hmac.Equal` for the AES-256 password check (P11B-SEC-001):** same approval question as U11a's password-compare finding, but on the R5/R6 arms. If the security owner signs off on `hmac.Equal` generally, fold P11B-SEC-001 into that fix.

---

## pades mid layer (unit U13)

Scope: the DSS-dictionary / revocation-source / external-CMS layer of `dss/pades` — 21 files: `pades_utils.go`, `pdf_dss_dict.go`, `pades_service.go`, `native_pdf_document_reader.go`, `pdf_document_analyzer.go`, `native_pdf_dict.go`, `signature_image_parameters.go`, `pdf_dss_dict_crl_source.go`, `pdf_dss_dict_ocsp_source.go`, `pdf_signature_dictionary.go`, `pades_with_external_cms_service.go`, `pades_diagnostic_data_builder.go`, `pades_crl_source.go`, `pdf_composite_dss_dict_crl_source.go`, `pdf_composite_dss_dict_ocsp_source.go`, `pades_ocsp_source.go`, `signature_image_text_parameters.go`, `pades_timestamp_parameters.go`, `external_cms_service.go`, `default_pdf_differences_finder.go`, `pdf_byte_range_document.go`.

### Files read (full vs budgeted)

- **Budgeted (read in full for completeness):** `pades_utils.go` (417 — over the 400-line budget by 17, read whole rather than 250+grep because the two I/O state machines `UtilsReplaceSignature`/`UtilsExtractRevisions` are the security-relevant regions and are contiguous), `pdf_dss_dict.go` (395 — ≤400 so full read per budget).
- **Full read (≤400 lines):** all other 19 target files.
- **Supporting files read to verify findings** (not in the 21, for cross-check only): `pdf_validation_data_container.go`, `pdf_cms_crl_source.go`, `pdf_cms_ocsp_source.go`, `pdf_composite_dss_dict_certificate_source.go`, `byte_range_input_stream.go`, `pades_signature.go`, `pdf_signature_revision.go`, `spi/offline_revocation_source.go`, `spi/dss_asn1_utils.go`, `spi/dss_utils.go`, `utils/codec.go`, `cades/cades_signature.go`, `spi/validation/analyzer/default_document_analyzer.go`.
- **Upstream Java checked against** (repo at `~/Workspace/esig/dss`, DSS 6.5.RC1): `PAdESUtils.java`, `PAdESCRLSource.java`, `PAdESOCSPSource.java`, `PdfDssDictCRLSource.java`, `AbstractPdfDssDict.java`, `SingleDssDict.java`, `PdfVriDict.java`, `DSSDictionaryExtractionUtils.java`, `PAdESService.java`, `PAdESWithExternalCMSService.java`, `PAdESTimestampParameters.java`, `PdfSigDictWrapper.java`.

### unsafe.Pointer assessment — `pdf_document_analyzer.go:354`

```go
func (a *PDFDocumentAnalyzer) ValidationData(signatures []validation.AdvancedSignature,
	detachedTimestamps []*validation.TimestampToken) *PdfValidationDataContainer {
	container, err := a.GetValidationDataWithTimestamps(signatures, detachedTimestamps)
	if err != nil { panic(err) }
	return (*PdfValidationDataContainer)(unsafe.Pointer(container))
}
```

**Verdict: sound, legal, and necessary. No action required (recorded as P13-STD-003).**

- **Same underlying type / layout invariant holds.** `PdfValidationDataContainer` embeds `validation.DataContainer` as its **first** field (`pdf_validation_data_container.go:17`), so by the Go spec the embedded field is at offset 0 and `&pvc.DataContainer == (*PdfValidationDataContainer)(unsafe.Pointer(&pvc.DataContainer))`. The producer is single and fixed: `GetValidationDataWithTimestamps` (`analyzer/default_document_analyzer.go:523`) calls `overrides.InstantiateValidationDataContainer()`, whose PDF override (`pdf_document_analyzer.go:319`) returns `&NewPdfValidationDataContainer(a.DssRevisions()).DataContainer`. So the dynamic type behind the `*validation.DataContainer` is **always** a `*PdfValidationDataContainer` — the cast never reinterprets a different concrete type.
- **No aliasing / GC hazard.** It is one 1:1 pointer reinterpretation with no escaping, no lifetime extension, no pointer arithmetic. The `unsafe.Pointer` round-trip (`*A → unsafe.Pointer → *B`) is legal Go for struct types per the spec's pointer-conversion rules.
- **Necessity.** Go has no covariant return types and the public `analyzer.DocumentAnalyzer` interface pins `GetValidationDataWithTimestamps` to return `*validation.DataContainer`; Java's unchecked cast has no direct Go equivalent. A type assertion `container.(*PdfValidationDataContainer)` would be cleaner **if** the base returned the concrete pointer, but it returns the base type, so the only options are (a) this cast, or (b) adding a `PdfValidationDataContainer() *PdfValidationDataContainer` override method to the analyzer. (b) is the lower-risk refactor if a future change ever makes the layout invariant non-obvious — noted as an open question, not a defect.

### Findings

#### P13-SEC-001 — VRI extraction panics on a non-dictionary `/VRI <key>` entry (divergence from Java)

| | |
|---|---|
| **Severity** | **High** |
| **Category** | SEC (panic on hostile input in the production validation path; also a Java-parity divergence) |
| **Location** | `pdf_dss_dict.go:254` (`NewPdfVriDict(name, vriDict.AsDict(name))`) → `NewPdfVriDict:278` → `newDssDictExtraction:66` → `dssDictionaryExtractionUtilsGetCertsFromArray:86` (`dict.AsArray(arrayName)` on a nil `PdfDict`) |

```go
// singleDssDictExtractVRIs — pdf_dss_dict.go:243
for _, name := range names {
    if !spi.DSSUtilsIsSHA1Digest(name) { continue }
    result = append(result, NewPdfVriDict(name, vriDict.AsDict(name))) // AsDict → nil when the value is not a dict
}
// ...
// dssDictionaryExtractionUtilsGetCertsFromArray — pdf_dss_dict.go:86
certsArray := dict.AsArray(arrayName)   // dict is a nil PdfDict here → nil-pointer panic
```

`AsDict` returns `nil` when the entry under a key is absent **or of another type** (`native_pdf_dict.go:88-104`). A PDF whose `/DSS /VRI` maps a valid 40-hex SHA-1 key to a **non-dictionary** value (a number, name, string, or `null`) passes the `DSSUtilsIsSHA1Digest(name)` gate, then `NewPdfVriDict(name, nil)` → `newDssDictExtraction(nil, …)` → `dict.AsArray(…)` on a nil interface → **panic**.

**Impact.** This is in the revision/DSS-dictionary extraction path (`GetRevisions` → `SingleDssDictExtract` → `newSingleDssDict`), i.e. the production validation path. A crafted PDF makes validation abort with a panic. The top-level facade (`dss/format.go` `recovered()`) recovers it into an error, so it is a validation-failure / robustness defect rather than a process crash — but it is a **divergence from upstream Java**, where `SingleDssDict.extractVRIs` wraps the whole loop in `catch (Exception e)` and logs + continues (`SingleDssDict.java:89-95`). The code comment at `pdf_dss_dict.go:241` ("Upstream logs and swallows any exception raised while walking the /VRI dictionary; this stays silent") is **incorrect** — it does not stay silent, it panics. Not a false-accept, so High (panic-on-hostile-input in the prod path) rather than Critical.

**Recommendation.** Guard the non-dict case the way Java's `catch` does: in `singleDssDictExtractVRIs`, skip (or log-and-continue) when `vriDict.AsDict(name)` is `nil` before calling `NewPdfVriDict`; and/or make `NewPdfVriDict`/`newDssDictExtraction` tolerate a `nil` `PdfDict` by returning empty token maps (mirroring `dssDictionaryExtractionUtilsGetCertsFromArray`'s `certsArray == nil` early-return). Correct the `:241` comment. Add a KAT/oracle fixture: a PDF with `/DSS << /VRI << <40-hex> 42 >> >>` asserting it validates (with zero VRIs) rather than panicking.

---

#### P13-STD-002 — `TimestampParameters.String()` prints the PDF password bytes

| | |
|---|---|
| **Severity** | **Low** |
| **Category** | STD (credential-in-string; upstream-faithful) |
| **Location** | `pades_timestamp_parameters.go:195-199` (`passwordProtection=%v`, `p.passwordProtection`) |

```go
return fmt.Sprintf("PAdESTimestampParameters [pdfSignatureCache=%v, ... passwordProtection=%v] %s",
    ..., p.passwordProtection, p.TimestampParameters.String())
```

`passwordProtection` is a `[]byte`; `%v` renders it as the raw byte values (e.g. `[104 101 108 108 111]`), so any code path that formats this value with `%v` (a `fmt.Stringer`, so any log/`Sprintf` of the parameters) emits the document password in cleartext.

**Impact.** Latent credential leak. **Upstream-faithful** — Java's `PAdESTimestampParameters.toString()` does `", passwordProtection=" + Arrays.toString(passwordProtection)` (`PAdESTimestampParameters.java:240`), which also prints the chars — so this is a deliberate 1:1 port, not a regression. Rated Low because it is only reachable if a caller ever stringifies the parameters, and it matches Java.

**Recommendation.** If the team wants to harden beyond parity, redact the field in `String()` (e.g. `passwordProtection=<set>` / `<unset>`) and add a `// DIVERGENCE, deliberate:` note plus a `DESIGN.md` entry per the upstream-tracking rule. Otherwise leave as-is and record it as an accepted upstream behavior.

---

#### P13-STD-003 — `unsafe.Pointer` cast at `pdf_document_analyzer.go:354` (assessment, no defect)

| | |
|---|---|
| **Severity** | **Info** |
| **Category** | STD (assessment) |
| **Location** | `pdf_document_analyzer.go:354` |

Assessed in the dedicated section above. The cast is legal, the layout invariant (first embedded field) is guaranteed by the Go spec, the producer is single and fixed so the dynamic type is always `*PdfValidationDataContainer`, and there is no aliasing/GC hazard. **No action required.** See the `unsafe.Pointer assessment` section.

### Lenses with no findings (explicit)

- **PERF — no findings.** The brief's two hot spots are clean:
  - *DSS-dict archival is not unbounded.* `dssDictExtraction`'s three token maps (`certMap`/`crlMap`/`ocspMap`) are bounded by the number of entries in the `/DSS` `/Certs`/`/CRLs`/`/OCSPs` arrays and the per-VRI arrays — one entry per PDF object, no recursive/expanding structure. `RevocationInfoArchival` (`UtilsRevocationInfoArchival`) parses the ADBE attribute once; its CRL/OCSP values are stored as their parsed binaries, not re-expanded. No unbounded buffer growth in these 21 files.
  - *Repeated re-parsing is bounded.* The composite sources (`pdf_composite_dss_dict_{crl,ocsp,certificate}_source.go`) and the per-signature sources cache their maps (`crlMap`/`ocspMap`/`certMap`) on first `CrlMap()`/`OcspMap()`/`CertificateMap()` call; the `filter*FromKeys` / `RevocationTokenIDs` / merge helpers (`padesCRLSourceMergeBinaryOrigins`, etc.) are linear-in-tokens scans, bounded by the (small) number of embedded tokens, not by attacker-controlled document size. `SortedKeys` sorts a key set of object-ids per call — O(k log k), k = tokens, fine.
  - `pades_utils.go` `UtilsExtractRevisions` / `UtilsReplaceSignature` do byte-at-a-time I/O, but the scratch buffer (`tempLine` / `temp`) is **bounded** (reset on every line-break or once it exceeds 5 / `len(signature)` bytes), so no O(n) buffer growth and no per-byte unbounded allocation. Faithful to Java's `BufferedInputStream` + `ByteArrayOutputStream` loop.
- **SEC (revocation-source correctness) — no false-accept found.** The CRL/OCSP sources compose the CMS (`PdfCmsCRLSource`/`PdfCmsOCSPSource`) and DSS-dict (`PdfDssDictCRLSource`/`PdfDssDictOCSPSource`) sources and `RevocationTokens` merges both. I verified the one suspicious asymmetry — `ADBERevocationValuesBinaries` reading from the **CMS** source while `ADBERevocationValuesTokens` reads from the **DSS-dict** source (`pades_crl_source.go:153-160`, `pades_ocsp_source.go:143-150`) — is **faithful to upstream Java** (`PAdESCRLSource.java:107-113`, `PAdESOCSPSource.java:109-115` do exactly this), so it is not a port bug. The `revocationDataOrigins*` nil-dereference of `s.dssDictionary` (`pdf_dss_dict_crl_source.go:222,258`, `pdf_dss_dict_ocsp_source.go:215,251`) is **unreachable** when `dssDictionary` is nil: the `filteredBinaries`/`filteredTokens` they iterate are filtered against `CrlMap()`/`OcspMap()`, which are empty when the dict is nil, so `revocationDataOrigins*` is never called — the "As upstream, dssDictionary is dereferenced without a nil check here" comments are accurate. `removeReferenceData` (`pdf_signature_dictionary.go:227`) mutates in place via `modifications[:0]`, matching Java's `removeIf` — not a bug.

### Tool log (U13)

Run from `dss/` (module root). Go `go1.27.0 darwin/arm64`. All three tools **available**.

```
$ gofmt -l pades
(empty — clean)

$ go vet ./pades/...
VET_OK            # exit 0, no diagnostics

$ golangci-lint run --config=../.github/.golangci.yml ./pades/...
0 issues.         # config passed explicitly, not auto-discovered
```

### Open questions (U13)

1. **P13-SEC-001 reachability in the corpus.** Is there any oracle/corpus PDF with a malformed `/VRI <key> → non-dict` entry? If not, the High finding is currently untested by the oracle sweep and the KAT fixture in its recommendation is the highest-value addition in this unit. Confirm whether `harness/` or `corpus/pades` exercises it.
2. **P13-SEC-001 fix shape.** Guard at `singleDssDictExtractVRIs` (skip nil sub-dict, matching Java's `catch`) vs. tolerating a nil `PdfDict` in `NewPdfVriDict`/`newDssDictExtraction`. The former is the closer port of `SingleDssDict.extractVRIs`; the latter is more robust if other callers can pass a nil dict. Decide and record in `DESIGN.md`.
3. **P13-STD-003 refactor option.** If the team wants to remove the project's only `unsafe.Pointer`, adding a `PdfValidationDataContainer() *PdfValidationDataContainer` override to the analyzer (so the concrete type is returned directly) would replace the cast with a plain method call. Low priority — the current cast is sound — but worth tracking if `spi/validation/analyzer` is ever refactored to return a concrete/parameterized container.
4. **P13-STD-002 approval.** Does the security owner sign off on printing the PDF password in `String()` (upstream-faithful) as acceptable, or should it be redacted with a `// DIVERGENCE, deliberate:` entry? Same class of question as the U11b password-compare findings.

---

## pades core A (unit U12a)

Scope: the PAdES core signing/level-check/timestamp layer — 3 files:
`native_pdf_signature_service.go` (1499), `pades_baseline_requirements_checker.go` (729),
`pades_timestamp_source.go` (504). Plus the supporting files they load through, read to
verify findings: `byte_range.go`, `byte_range_input_stream.go`, `pades_utils.go`,
`native_pdf_document_reader.go`, `pdf_sig_dict_wrapper_factory.go`, `pdf_signature_dictionary.go`,
`pdf_signature_dictionary_comparator.go`, `pdf_doc_timestamp_revision.go`,
`pdf_signature_revision.go`, `pdf_cms_revision.go`, `pdf_timestamp_token.go`,
`native_pdf_dict.go`, `pdf_document_analyzer.go` (caller), `dss/format.go` (recover boundary),
`spi/validation/baseline_requirements_checker.go`, `spi/validation/timestamp/signature_timestamp_source.go`.

**Out of scope / not re-reported:** `internal/pdf` (U11a/U11b), `cades` (U04a/b), and the
`pades` files in U12b/U13a/U13b/U13c (modif finder, signature, service, revocation sources,
etc.). Known gaps in `docs/compatibility/known-gaps.md` (visible-signature appearances not
supported; the six CAdES/CMS-layer parity gaps in `pades/testdata/broadgen/README.md`) are
not re-reported. Deliberate `PORTING.md` conventions are not findings: 1:1 Java port,
panic→error at the facade boundary, byte-identical enum/wire values, dropped slf4j logging,
`SA1019`/staticcheck-style rules disabled.

**Upstream Java checked against** (repo at `~/Workspace/esig/dss`, DSS 6.5.RC1):
`AbstractPDFSignatureService.java`, `PdfBoxSignatureService.java`,
`PAdESBaselineRequirementsChecker.java`, `PAdESTimestampSource.java`,
`PdfSignatureDictionaryComparator.java`, `PdfDocTimestampRevision.java`,
`PdfBoxDict.java` (`match`), `AbstractTimestampSource.java` (`addReferences`),
`DigestAlgorithm.java`.

### Files read (full vs budgeted)

All three target files were read **in full** (the piecewise reads below tile 1→end with no
gaps), with **every** `func` signature grep'd first and then the security-relevant regions
(level checks, byte-range, digest, sign, timestamp, revision, cert/chain) read. This exceeds
the 250+grep budget; nothing was left signature-only.

| File | Lines | Read | Func map | Notes |
|---|---|---|---|---|
| `native_pdf_signature_service.go` | 1499 | **full** (1–250, 250–420, 420–526, 526–616, 616–643, 643–925, 925–1208, 1208–1499) | ✅ (69) | every fn read; security regions = `buildRevision`/digest, `GetRevisions`/`ValidateByteRange`/`isSignedContentComplete`, `signatureOptions`, `BuildDSSDictionary`, `Analyze*Modifications`, `Sign`/`SignDocument` |
| `pades_baseline_requirements_checker.go` | 729 | **full** (1–250, 250–560, 560–729) | ✅ (40) | every fn read; level checks B/BES/T/LT/LTA/LTV/EPES/PKCS7 + `coversLTLevelData`/`coversRevocationTokens`/`coversOwnRevocationData` |
| `pades_timestamp_source.go` | 504 | **full** (1–250, 250–504) | ✅ (26) | every fn read; `populateAndValidateDocumentTimestamps`, `GetSignatureTimestampReferences`, `GetAdbeRevocationInfoArchivalReferences`, `padesTS*` helpers |

**Test coverage (skimmed, not read fully):** `pdf_document_validator_smoke_test.go` (full
pipeline against `testdata/upstream` fixtures, permissive verifier), `pades_upstream_cross_validation_test.go`
/ `pades_downstream_cross_validation_test.go` (live Java oracle, skip-gated on
`DSS_UPSTREAM_HOME`), `byte_range_validate_test.go`, `pdf_number_test.go`,
`native_pdf_dict_match_test.go`. 11 top-level test funcs total. The level-checker and
timestamp-source have **no dedicated unit test** in-package — they are exercised end-to-end
only via the smoke test and the cross-validation oracle.

### Findings

#### P12A-SEC-001 — `GetRevisions` drops Java's per-iteration catch; one malformed document timestamp aborts the whole document

| | |
|---|---|
| **Severity** | **High** |
| **Category** | SEC (panic on hostile input in the production validation path + Java-parity divergence; false-reject/DoS, not a false-accept) |
| **Location** | `native_pdf_signature_service.go:712` (`GetRevisions` per-iteration body, no recover) → `pdf_doc_timestamp_revision.go:47-53` (the escaping panic) |

```go
// native_pdf_signature_service.go:712 — per-iteration body of GetRevisions, NO surrounding recover
if s.IsDocTimestamp(signatureDictionary) {
    newRevision = NewPdfDocTimestampRevision(signatureDictionary, fields, signedContent, ...)
}
// pdf_doc_timestamp_revision.go:51 — the panic that escapes GetRevisions
if _, err := timestampToken.MatchDataDocument(revision.SignedData()); err != nil {
    panic(fmt.Sprintf("Unable to create a PdfDocTimestampRevision : %s", err.Error()))
}
```

Java's `AbstractPDFSignatureService.getRevisions` wraps the **entire per-dictionary body**
(byte-range check, `isSignedContentComplete`, `verifyPdfSignatureDictionary`, the
`new PdfDocTimestampRevision(...)`/`new PdfSignatureRevision(...)` construction, and the
previous-revision read) in `} catch (Exception e) { LOG.warn("Unable to parse signature {}
. Reason : {}", ...); }` — i.e. a failure on one revision is logged and the loop **continues
to the next signature**. The Go `GetRevisions` has **no** such per-iteration `recover`: the
only `recover()` in the function is inside `ValidateByteRange` (`:850`), which protects only
the `byteRange.Validate()`/hex-extraction path. So a document timestamp whose RFC 3161 token
does not match its own signed data (a malformed/corrupt TST — or an attacker-crafted one)
makes `NewPdfDocTimestampRevision` panic, and that panic escapes `GetRevisions` →
`PDFDocumentAnalyzer.Revisions()` (no recover, `pdf_document_analyzer.go:312`) → the whole
`Validate` call.

**Impact.** (1) **Parity divergence:** Java isolates the damage to the single bad revision and
still validates the document's *other* signatures; Go aborts the entire document with an
error, so no signature in it is validated. (2) **DoS vector:** an attacker can embed one bad
document timestamp in a multi-signature PDF to make the whole document's validation fail. It
is **not** a false-accept (fails safe), and the top-level facade `recovered()`
(`format.go:329`) bounds the panic to a clean returned error rather than a process crash.
Rated High to stay consistent with the sibling P13-SEC-001 (same pattern: panic in a prod
path that Java catches-and-continues, facade-bounded, not a false-accept) and with the
"panic on hostile input in prod path" rubric band.

**Recommendation.** Add the missing per-iteration boundary in `GetRevisions`: wrap the body of
the `for _, sigDictEntry := range sigDictionaries` loop in a `defer`-based `recover()`
(log-and-continue on the Java `catch` sites) — or at minimum guard the two revision
constructors — so a single malformed revision is skipped and the remaining signatures are
still validated, matching Java. Add an oracle/KAT fixture: a PDF with one good signature and
one document-timestamp whose TST message-imprint does not match its signed data, asserting the
good signature still validates rather than the whole document erroring.

---

#### P12A-PERF-001 — process-wide `pdfTimestampTokenRegistry` grows unboundedly and pins the full original PDF bytes

| | |
|---|---|
| **Severity** | **Medium** |
| **Category** | PERF-LEAK |
| **Location** | `pdf_timestamp_token.go:37-39` (registry), `:60-62` (insert), `:76-78` (lookup) — supporting file for `pades_timestamp_source.go` (`PdfTimestampTokenOf`) and `native_pdf_signature_service.go` (`NewPdfDocTimestampRevision`) |

```go
// pdf_timestamp_token.go:37 — process-wide map, populated per NewPdfTimestampToken, NEVER removed
var pdfTimestampTokenRegistry = map[*validation.TimestampToken]*PdfTimestampToken{}
// :60-62 — the only mutation besides lookup; grep confirms no delete(...) anywhere in the package
pdfTimestampTokenRegistryMu.Lock()
pdfTimestampTokenRegistry[base] = token
pdfTimestampTokenRegistryMu.Unlock()
```

`NewPdfTimestampToken` is called for **every** document-timestamp revision, both from
`NewPdfDocTimestampRevision` (in `GetRevisions`, once per `IsDocTimestamp` revision) and from
`PDFDocumentAnalyzer.createPdfTimestampToken`. Each call allocates a **fresh** base
`*validation.TimestampToken` (the map key) and stores a `*PdfTimestampToken` that embeds
`pdfRevision *PdfDocTimestampRevision`, whose `pdfCMSRevisionBase.signedContent` is a
`PdfByteRangeDocument` wrapping the **entire original `model.DSSDocument`** (not just the
byte-range slice — `pdf_byte_range_document.go:54` stores the full `pdfDocument`). Because the
base key is freshly allocated each call and entries are never deleted, the registry grows by
one full-document pin per document-timestamp per validation, with no release path.

**Impact.** In a long-lived validation service (or a single document with many document
timestamps), the full bytes of every validated PDF stay alive in memory indefinitely — an
unbounded leak. This is a **Go-specific structural workaround** (the registry stands in for
Java's `instanceof PdfTimestampToken`, which needs no side table) and is not present in the
upstream Java. Not a correctness bug, so Medium.

**Recommendation.** Scope the registry to the validation: either (a) key the map by the
`PdfDocTimestampRevision`/signature under validation and clear the entries for a document when
its analyzer is done (add a `release`/`Clear` called at the end of the validation pipeline), or
(b) avoid the global side table entirely by having `PdfTimestampTokenOf` recover the wrapper
through a field on the `TimestampToken` itself (or by returning the wrapper from the
constructors that already hold it) instead of a process-wide lookup. If the table must stay
global, bound it (e.g. a per-signature weak reference or an explicit eviction on
`TimestampToken` GC) and record the choice in `DESIGN.md`.

---

#### P12A-PERF-002 — `BuildDSSDictionary`'s per-token linear dedup makes DSS-dict accumulation O(tokens²)

| | |
|---|---|
| **Severity** | **Low** |
| **Category** | PERF-BIGO |
| **Location** | `native_pdf_signature_service.go:1099` (`nativePDFSignatureServiceContainsToken`) as called from `BuildDSSDictionary` (`:932-1065`) |

```go
// native_pdf_signature_service.go:1099 — linear scan per token
func nativePDFSignatureServiceContainsToken(tokens []pdf.TokenRef, token pdf.TokenRef) bool {
    for _, candidate := range tokens { ... }
}
// BuildDSSDictionary:971-999 — called once per cert/CRL/OCSP appended → O(n²) overall
if !nativePDFSignatureServiceContainsToken(vriEntry.Certs, tokenRef) {
    vriEntry.Certs = append(vriEntry.Certs, tokenRef)
    if !nativePDFSignatureServiceContainsToken(dssDictionary.Certs, tokenRef) { ... }
}
```

`BuildDSSDictionary` loops over each signature's certificate/CRL/OCSP tokens and, for each,
runs two O(n) `nativePDFSignatureServiceContainsToken` scans over the growing `vriEntry.*`
and `dssDictionary.*` slices — O(n²) in the number of tokens.

**Impact.** Minor. This is the **signing/extension** path (`addDssDictionary`), not the
validation hot path, and n (the number of embedded validation tokens in a `/DSS` dictionary)
is small in practice and bounded by the document's own data, not by an unbounded attacker
dimension. It is also **upstream-faithful** — Java's `PdfBoxSignatureService` uses
`COSArray#indexOf(...)` (also linear) for the same dedup. Low / observation.

**Recommendation.** Optional: replace the linear `ContainsToken` scans with the existing
`knownObjects map[string]pdf.TokenRef` (already built in `BuildDSSDictionary`) as the
membership test, making the accumulation O(n). Only worth doing if a profiling pass shows the
DSS-dict build is hot; otherwise leave as the faithful port.

### Lenses with no findings (explicit)

- **SEC (level checks) — no false-accept found.** Every level predicate was checked against
  the Java oracle and is a faithful 1:1 port, with the higher-level evidence actually required:
  - `HasBaselineBProfile` / `HasExtendedBESProfile` — all SPO entries (M, Contents, Filter,
    ByteRange, SubFilter), CMS `id-data`, requirement (d)/(m)/(e)/(j) Reason-vs-commitment/
    signature-policy exclusions, and the `ETSI.CAdES.detached` SubFilter check match
    `PAdESBaselineRequirementsChecker.java` exactly.
  - `HasBaselineLTProfile` (`hasLTProfile`) = `MinimalLTRequirement()` (revocation-data
    presence via `spi/validation`) **and** (all-self-signed **or** DSS dictionary present) —
    so LT is never accepted without revocation evidence. `HasBaselineLTAProfile`
    (`isBaselineLTATimestamp`) additionally requires `containsRFC3161SubFilter &&
    coversLTLevelData`. `HasPKCS7LTAProfile` requires only `coversLTLevelData` — **faithful**:
    Java's `hasPKCS7LTAProfile` also omits the RFC3161-subfilter test that PAdES-LTA carries,
    so the PAdES-LTA vs PKCS7-LTA asymmetry is preserved, not a port bug.
  - `coversLTLevelData`/`coversRevocationTokens`/`coversDSSCertificateTokens`/
    `coversOwnRevocationData`/`coversTimestampTokens`/`coversToken` — line-for-line faithful to
    the Java private methods (including the "at least one revocation token per related
    certificate must be covered" rule and the `checkAllRequiredRevocationDataPresent` fallback).
  - Byte-range math: `isSignedContentComplete` (`(firstEnd-firstStart) + secondEnd`) and
    `IsSignatureCoversWholeDocument` are byte-identical to Java's formulas; `ByteRange.Validate`
    enforces start-0 / non-negative / ordered parts; the SIWA check
    (`isContentValueEqualsByteRangeExtraction` → `bytes.Equal(cms, getSignatureValue(...))`)
    matches Java and is wrapped in the `ValidateByteRange` `recover`.
  - Revision ordering: `PdfSignatureDictionaryComparator.Compare` (envelop / equal / strange
    byte-range + `SigningDate` tiebreak) and the stable insertion sort that preserves Java's
    `Stream.sorted()` stability are faithful; `getRevisions`'s `containsDSSRevisions` /
    `getPreviousDssDictAndUpdateIfNeeded` logic matches.
  - Timestamp population: `populateAndValidateDocumentTimestamps`'s `PdfDocTimestampRevision` /
    `PdfDocDssRevision` / `PdfSignatureRevision` switch, the `signatureRevisionReached` /
    `dssRevisionReached` flags, the `ArchiveTimestampTypePAdES` assignment, the VRI-token dedup
    by identity, and the VRI `matchData` pass all mirror `PAdESTimestampSource.makeTimestampTokensFromUnsignedAttributes`
    + `validateTimestamps`. The two "concrete-but-overridable" hooks that are unreachable from
    the base's internal dispatch (`getSignatureTimestampReferences`, `getTimestampScopes`) are
    **already documented as narrow GAPs in the file header** (and match the three GAPs in
    `cades_timestamp_source.go`'s header) — not re-reported.
  - `nativePdfObjectEquals`/`nativePdfIsNull` (developer-extension `match`) faithfully reproduce
    `PdfBoxDict.match`'s `targetObject != null && !targetObject.equals(currentObject)` guard.

- **PERF (hot-path BigO / repeated re-parsing) — clean beyond P12A-PERF-002.** The level checks
  are linear in the (small) number of CMS attributes / timestamps / revocation tokens, with no
  attacker-controlled super-linear dimension (the one O(n²) is the signing-path DSS-dict dedup,
  P12A-PERF-002, upstream-faithful and small-n). `padesTSAddReferences`/`padesTSAddReference`
  dedup is O(n²) but bounded by the small reference set and mirrors Java's
  `DSSUtils.enrichCollection`. Digest computation in `buildRevision` streams
  (`io.Copy(digest, result.SignedData())`) rather than slurping. No unbounded buffer growth or
  per-byte unbounded allocation in these three files.

- **STD — no findings.** Doc comments are thorough and accurate (port provenance, Java method
  mapping, and the structural-deviation rationale for the checker and timestamp source are all
  correct and match the code). Panic→error is applied consistently at the facade boundary, and
  the in-package panics carry the Java messages. The one `unsafe`-free, `sync`-correct
  concurrency touch is the mutex-guarded `pdfTimestampTokenRegistry` (its *lifetime* is the
  P12A-PERF-001 issue, not a data race). `defer`/cleanup is correct: every
  `NewNativePdfDocumentReader`/`NewNativePdfDocumentReader`-style reader is paired with
  `defer func() { _ = reader.Close() }()`; `AddDssDictionary`, `SignDocument`, `buildRevision`,
  `GetRevisions`, `GetAvailableSignatureFields`, `AddNewSignatureField`, `Analyze*Modifications`
  all close their readers.

### Tool log (U12a)

Run from `dss/` (module root). Go `go1.27.0 darwin/arm64`. All three tools **available**.

```
$ gofmt -l pades
(empty — clean)

$ go vet ./pades/...
VET_OK            # exit 0, no diagnostics

$ golangci-lint run --config=../.github/.golangci.yml ./pades/...
0 issues.         # config passed explicitly, not auto-discovered
```

### Open questions (U12a)

1. **P12A-SEC-001 fix shape & scope.** Add the per-iteration `recover()` in `GetRevisions`
   (closest port of Java's `catch`) vs. guarding just the two revision constructors vs. giving
   `NewPdfDocTimestampRevision`/`NewPdfSignatureRevision` error-returning variants. The first is
   the most faithful and also covers any *other* future panic in the per-iteration body
   (`CheckConsistency`, `IsSignatureCoversWholeDocument`, the previous-revision reader). Decide
   and record in `DESIGN.md`. Confirm the smoke/oracle corpus has no multi-signature PDF with a
   corrupt TST (the recommended KAT would be the first coverage of this path).
2. **P12A-PERF-001 registry lifetime.** Should the `pdfTimestampTokenRegistry` be per-validation
   (cleared at pipeline end) or replaced by a `TimestampToken`-owned back-pointer? Note this
   registry is shared with the U13/U12b callers (`PDFDocumentAnalyzer.createPdfTimestampToken`,
   `AnalyzePdfModifications`'s `instanceof` recovery), so the fix should be coordinated across
   the pades unit, not just U12a.
3. **Level-checker test gap.** `pades_baseline_requirements_checker.go` and
   `pades_timestamp_source.go` have no dedicated in-package unit tests — they are only exercised
   via the end-to-end smoke test and the (skip-gated) Java cross-validation oracle. A focused
   table test pinning each `Has*Profile`/`coversLTLevelData` branch (especially the PAdES-LTA vs
   PKCS7-LTA asymmetry and the `coversOwnRevocationData` fallback) would protect the SEC-clean
   result above against regression. Worth adding even if the oracle is absent locally.

---

## pades modification-detection + signature (unit U12b)

**Batch:** 03 (pades, modification-detection + signature layer)
**Unit:** U12b — the four PAdES classes that find and categorize object
modifications between revisions, plus the PAdES `AdvancedSignature` and its
`SignatureParameters`.
**Scope (non-test):** `default_pdf_object_modifications_finder.go` 462 ·
`pdf_object_modifications_filter.go` 459 · `pades_signature.go` 440 ·
`pades_signature_parameters.go` 439.
**Skimmed (coverage only, not fully read):** `pades_upstream_cross_validation_test.go`
(the oracle that pins modification-detection verdicts + the four category counts),
and the signatures of `object_modification.go` / `pdf_object_tree.go` /
`pdf_object_key.go` / `native_pdf_object_key.go` / `native_pdf_dict.go` /
`pades_timestamp_parameters.go` / `utils/io.go` (supporting types the four files
call into).
**Date:** 2026-08-25
**Depth:** **full-read** — all four in-scope files read end-to-end (each was read as
first-250 + func-map + remaining lines, which together covered 100 % of every
function; no function is signature-only).
**Out of scope (not re-reported):** `internal/pdf` (U11a/U11b) and `cades` (U04)
parity findings; the six CAdES/CMS-layer gaps in `known-gaps.md` /
`pades/testdata/broadgen/README.md`; the DSS-dict/revocation sources and service/
reader/analyzer (U13a/b/c, sibling units already in this file). Deliberate
conventions from `dss/PORTING.md` are not findings: 1:1 Java port, byte-identical
enum/wire values, stdlib-first, panics-at-facade-boundary, `SA1019`/`staticcheck`
style rules disabled.

---

### Read budget (per file — full vs budgeted)

All four files were read in full (first-250 block → **all** `func` signatures
grep'd → remaining security-relevant regions read). Because each file is 439–462
lines and every function is a compare/filter/sign/param function (i.e. all
security-relevant), the net coverage is **100 % / full** for all four.

| File | Lines | Func map | Net coverage |
|---|---|---|---|
| `default_pdf_object_modifications_finder.go` | 462 | ✅ (24) | **full** — `Find`/`FindInDicts`/`compareDictsRecursively`/`compareObjectsRecursively`/`compareSimpleObjects`/`compareNumericValues`/`compareArraysRecursively`/`isProcessedReference`/`addProcessedReference`/`compareDictStreams`/`rawStream*`/`objectModificationSet.Add`/`pdfObjectTreeReferenceSet` all read |
| `pdf_object_modifications_filter.go` | 459 | ✅ (43) | **full** — `Filter`/`skipChange`/`isExtensionChange` + all 18 sub-predicates/`isSignatureOrFormFillChange` + all 9 sub-predicates/`isAnnotationChange`/`checkRecursivelyForNewSignatureCreation`/`isStreamFill`/`isFontCreationChange`/`getParentKey`/`isOneOf` all read |
| `pades_signature.go` | 440 | ✅ (39) | **full** — `NewSignature`/`CertificateSource`/`CRLSource`/`OCSPSource`/`Complete*Source`/`TimestampSource`/`DocumentTimestamps`/`VRITimestamps`/`FindSignatureScopes`/`SigningTime`/`SignerDocumentContent`/`BuildSignatureDigestReference`/`DataFoundUpToLevel`/`Has*Profile`/`DssDictionary`/`VRIKey`/`VRICreationTime`/`AddExternalTimestamp` all read |
| `pades_signature_parameters.go` | 439 | ✅ (42) | **full** — `NewSignatureParameters`/`SetSignatureLevel`/all `Get*`/`Set*`/`SigningDate`/`DeterministicId`/`SigningTimeZone`/`IsIncludeVRIDictionary`/`getPAdESContext`/`*TimestampParameters`/`newTimestampParametersFromCAdES`/`PdfSignatureCache`/`Reinit`/`String`/`Equals`/`signatureImageParametersEqualPointers` all read |

Supporting types read to answer the security questions (budgeted, targeted):
`object_modification.go` (`Equals`), `pdf_object_tree.go` (`Copy`/`Equals`/
`AddReference`/`IsProcessedReference`/`ChainDeepness`), `pdf_object_key.go` +
`native_pdf_object_key.go` (map-key comparability), `native_pdf_dict.go`
(`RawStreamSize`/`Parent`/`AsDict`/`newNativePdfDictWithParent`),
`pades_timestamp_parameters.go` (`newTimestampParametersFromCAdES`),
`utils/io.go` (`CompareInputStreams`), `spi/validation/default_advanced_signature.go`
(`Complete*Source` base).

---

### Findings

#### SEC — clean (no findings)

I treated this as the top-priority lens and verified the three false-accept /
bypass questions against the upstream Java source (`DefaultPdfObjectModificationsFinder.java`,
`PdfObjectModificationsFilter.java`) line-by-line:

- **No missing-modification / false-accept.** The finder walks from the catalog
  dict and compares every signed vs final key; the indirect-reference dedup
  (`processedObjects` keyed by `(name, objectKey)` + per-tree `refChain`) is
  faithful to Java and, on first encounter, still *performs* the comparison
  before recording the ref — so a modified object that is also deduped is still
  compared, never silently skipped. Type-mismatch (dict↔simple) falls to the
  `default` branch and records a `Modify` (not dropped). Recursion is bounded by
  `maximumObjectVerificationDeepness` (500), so no unbounded stack on a
  deep-nested hostile tree. This is oracle-pinned: the 248-PDF sweep in
  `broadgen/README.md` reports "modification-detection classification" matched
  Java exactly (none of the six known gaps is a false accept).
- **No filter bypass.** The filter only *categorizes*; its only silent-drop is
  `skipChange` (delete `/AP`, modify `/F`, `/T`, `/Itext`), which is byte-for-byte
  Java's `skipChange`. The broad "allow" predicates (`isDSSDictionaryChange`,
  `isExtensionsChange`, `isMetaDataChange`, …) mirror Java exactly, and
  `checkRecursivelyForNewSignatureCreation` walks the `Parent()` chain without a
  cycle guard — but the chain is exactly the object tree the finder already walked
  (≤ 500 deep), so it is bounded and cannot infinite-recurse. `isOtherAnnotChange`'s
  deletion fall-through (returns `true` when the `/V` sub-dict is absent) is a
  faithful copy of Java's identical fall-through — not a bypass.
- **No unsafe signing defaults.** `signatureFilter`/`signatureSubFilter` default to
  `Adobe.PPKLite` / `ETSI.CAdES.detached` (the correct PAdES CAdES-detached pair);
  digest defaults come from the CAdES base (SHA-2 family); `SetSignatureLevel`
  panics on a non-PAdES level (Java's `IllegalArgumentException`). `passwordProtection`
  is opt-in (nil default) and is a *signing* parameter, not a validation one.
- **No byte-range arithmetic** in these four files (that is `internal/pdf`, U11b).
- **No re-introduced unbounded access.** `compareDictStreams` uses the
  already-bounded `internal/pdf` `RawStreamSize`/`CreateRawInputStream`
  (U11a/U11b-reviewed); the finder adds no new unbounded reads.

#### PERF — 2 findings

**P12B-PERF-001** — Severity **Medium** — Category **PERF-BIGO**
Location: `pades/default_pdf_object_modifications_finder.go:419-427`
(`objectModificationSet.Add`)
Evidence:
```go
func (s *objectModificationSet) Add(objectModification ObjectModification) {
	for _, existing := range s.items {          // linear scan over ALL prior entries
		if existing.Equals(objectModification) {
			return
		}
	}
	s.items = append(s.items, objectModification)
}
```
Impact: the modification set is built on the validation hot path (every revised-PDF
validation). Each `Add` scans all previously-added modifications and calls
`ObjectModification.Equals` (which compares the full key/ref chains), so a report
of `M` modifications costs O(M²). This is a **complexity divergence from upstream**:
Java's `LinkedHashSet<ObjectModification>` gives O(1)-amortized insert *and*
insertion order; the Go port keeps the ordering (slice) but regressed the insert
from a hash probe to a linear scan. Bounded by the number of *distinct* modified
object-tree paths (small on a validly-extended PDF, larger on a genuinely tampered
one), but it is slower than the Java it ports with no functional reason.
Recommendation: keep the ordering slice for iteration but back dedup with an
`map[string]struct{}` keyed on a stable identity (e.g. the `PdfObjectTree` string +
action type) for O(1) membership — matching Java's `LinkedHashSet`. Low-risk, no
behaviour change. (If left as-is, it is a faithful-enough-but-simpler stand-in, not
a correctness bug.)

**P12B-PERF-002** — Severity **Low** — Category **PERF-BIGO**
Location: `pades/default_pdf_object_modifications_finder.go:285-315`
(`compareArraysRecursively`)
Evidence:
```go
if objectKey != nil {
	for j := 0; j < secondArray.Size(); j++ { // linear scan for each i
		if objectKey == secondArray.ObjectKey(j) {
			finalRevObject = secondArray.Object(j)
		}
	}
}
```
Impact: for arrays whose elements are indirect object references (e.g. `/Kids`,
`/B`, `/Annots`), matching the first array against the second is O(n·m) on array
element count. This is **faithful to upstream Java** (`DefaultPdfObjectModificationsFinder.compareArraysRecursively`
has the identical nested loop) and is bounded by a *single* array's size (not the
whole document), so it is a perf observation rather than a port regression. Only
relevant for PDFs carrying a large array of indirect refs.
Recommendation: optional — index `secondArray.ObjectKey(j)` into a
`map[PdfObjectKey]int` once per array (O(m) build + O(1) lookup) → O(n+m). Because
it changes no verdict, it would need a `// DIVERGENCE, deliberate:` note (upstream
tracks no such optimization) before adopting.

#### STD — clean (no findings)

- **Error handling** is faithful: error returns below the facade, panics only at
  the documented boundary (`SetPdfObjectModificationsFilter(nil)`,
  `NewSignature`/`CertificateSource`/`VRIKey`/`BuildSignatureDigestReference`
  digest errors) — all mirror Java's thrown `IllegalArgumentException` /
  `IOException`/`RuntimeException`.
- **Doc comments** are thorough — every function carries a `Port of #…` line, and
  the non-obvious Go adaptations are explained at length (the `NewSignature`
  eager-source warm-up that compensates for Go's lack of Java virtual dispatch,
  and the `SignatureParameters` context/timestamp-field split that keeps
  `&parameters.SignatureParameters` sharing intact). Both are deliberate, correct,
  and well-justified — not findings.
- **No shared-state mutation:** `CompleteCertificateSource/CRLSource/OCSPSource`
  call the CAdES base (which builds a *fresh* `ListCertificateSource` per call) and
  then `AddAll` the PAdES DSS source — no cross-call cache corruption.
- **Test coverage** is oracle-grade: `pades_upstream_cross_validation_test.go`
  pins `pdfModificationsDetected` and the four category counts
  (`secure`/`formFill`/`annot`/`undefined`) against the Java oracle across the
  corpus. (No dedicated KAT unit test for the finder/filter, but oracle coverage is
  the project's standard and is the stronger guarantee here.)
- **Info (cross-reference, not a new defect):** `SignatureParameters.String()`
  (`pades_signature_parameters.go:396`) renders the `passwordProtection` byte
  slice in cleartext via `%v`, faithful to Java's `Arrays.toString(passwordProtection)`.
  This is the *same* pattern already recorded as **P13-STD-002** in the U13 section
  above (sibling unit, same file family); listed here only for completeness of the
  U12b record, not re-counted as a separate finding.

### Tool log (U12b)

Run from `dss/` (module root). Go `go1.27.0 darwin/arm64`. All three tools
**available**.

```
$ gofmt -l pades
(empty — clean)

$ go vet ./pades/...
VET_OK            # exit 0, no diagnostics

$ golangci-lint run --config=../.github/.golangci.yml ./pades/...
0 issues.         # config passed explicitly, not auto-discovered
```

(No fuzz target exists for the finder/filter in this package; the modification-detection
path is instead oracle-verified by the cross-validation test noted above.)

### Open questions (U12b)

1. **P12B-PERF-001 — worth fixing?** The O(M²) dedup is the one place the Go code is
   measurably slower than the Java it ports, with no functional benefit. Should the
   ordered-slice + map-membership hybrid be adopted (matching Java's `LinkedHashSet`
   semantics exactly)? Low risk; no verdict change.
2. **P12B-PERF-002 — adopt the map index?** It would require a
   `// DIVERGENCE, deliberate:` entry (upstream has no such optimization). Given it is
   faithful and bounded by single-array size, recommend leaving as-is unless a corpus
   PDF with a very large indirect array surfaces.
3. **`NewSignature` warm-up cost.** The constructor eagerly builds and caches the
   PAdES certificate/CRL/OCSP sources to win the lazy-cache race (documented). This is
   correct but means every `NewSignature` does non-trivial work up front, whereas Java
   would build lazily. Acceptable for correctness; noted in case a future refactor of the
   `spi/validation/timestamp` generic dispatch removes the need for the warm-up.
4. **P12B cross-ref to P13-STD-002.** If the security owner redacts the password from
   `String()` under P13-STD-002, the same fix should be applied here in
   `SignatureParameters.String()` for consistency (single `// DIVERGENCE, deliberate:`
   entry covering both).

---

## pades tail + alerts + exception (unit U13c)

**Unit:** U13c — the remaining **tail** of `dss/pades` (non-test `.go` files not yet
reviewed in the earlier U11/U12/P11/P12/P13 units), plus the two leaf subpackages
`dss/pades/alerts` and `dss/pades/exception`, plus the test-only generator
`dss/pades/testdata/crossgen/main.go`.
**In scope (all ≤187 lines, all read in full):** every non-test `.go` under
`pades/`, `pades/alerts/`, `pades/exception/` **except** the 29 already-reviewed
files (native_pdf_signature_service, pades_baseline_requirements_checker,
pades_timestamp_source, default_pdf_object_modifications_finder,
pdf_object_modifications_filter, pades_signature, pades_signature_parameters,
pades_utils, pdf_dss_dict, pades_service, native_pdf_document_reader,
pdf_document_analyzer, native_pdf_dict, signature_image_parameters,
pdf_dss_dict_crl_source, pdf_dss_dict_ocsp_source, pdf_signature_dictionary,
pades_with_external_cms_service, pades_diagnostic_data_builder, pades_crl_source,
pdf_composite_dss_dict_crl_source, pdf_composite_dss_dict_ocsp_source,
pades_ocsp_source, signature_image_text_parameters, pades_timestamp_parameters,
external_cms_service, default_pdf_differences_finder, pdf_byte_range_document).
That leaves **82** in-scope files (79 in `pades/` — 78 listed plus
`cms_for_pades_builder_helper.go` — 1 in `pades/alerts/`, 2 in
`pades/exception/`) plus `pades/testdata/crossgen/main.go` (187 lines).
**All 83 files read in full.** Test files skimmed for coverage only (names, not bodies):
`byte_range_validate_test.go`, `native_pdf_dict_match_test.go`,
`pades_downstream_cross_validation_test.go`,
`pades_upstream_cross_validation_test.go`, `pdf_document_validator_smoke_test.go`,
`pdf_number_test.go`.
**Deliberate conventions (not findings):** 1:1 Java DSS 6.5.RC1 port; staticcheck style
rules disabled on purpose; `SA1019` legacy crypto allowed; byte-ident enum/wire values;
stdlib-first; no cgo; panic→error at the facade boundary. `docs/compatibility/known-gaps.md`
items not re-reported: visible-signature rendering non-goal, in-memory-only streaming, the
six CAdES/CMS-layer parity gaps, Go-toolchain `compress/flate` + Unicode-table sensitivity.
**Upstream cross-check:** two look-suspicious sites were diffed against
`/Users/utain/Workspace/esig/dss` (Java DSS 6.5.RC1) and confirmed faithful — **not**
findings: (a) `byte_range.go` `Length()`/`Validate()` and `byte_range_input_stream.go`
`skip()` are byte-identical to `ByteRange.java` / `ByteRangeInputStream.java` (including the
mixed start/length reading of the 4-tuple — that is upstream's own convention, not a port
bug); (b) `pdf_permissions_checker.go` `isSignatureFieldCreationForbidden`'s INCLUDE/EXCLUDE
fall-through is line-for-line the upstream `switch`+`break` (an INCLUDE with an empty field-id
genuinely returns `false` upstream too — creating a *new* field is allowed when the field id
is empty).

### Findings — SEC

**P13C-SEC-001** — Severity **Low** — Category **SEC**
Location: `pades/cms_for_pades_baseline_requirements_checker.go:105`
Evidence:
```go
if !signature.CMS().IsDetachedSignature() { return false }
return c.HasBaselineBProfile()   // adds CAdES requirement (k); upstream calls cmsBaselineBRequirements()
```
Impact: `IsValidForPAdESBaselineBProfile` rejects a CMS that Java's
`isValidForPAdESBaselineBProfile` (which calls the protected `cmsBaselineBRequirements()`
directly) would accept — specifically a CMS failing only CAdES requirement (k): a
`signature-policy-store` present without a `signature-policy-identifier` defining
`sigPolicyHash`. Direction is fail-closed (false-reject, never a false-accept), so no
security weakening; the cost is a spurious "not PAdES-B compliant" verdict on an edge CMS
Java would pass. Already documented in the file header as a DEVIATION (the cades
`cmsBaselineBRequirements()` is unexported in Go, so only the stricter
`HasBaselineBProfile()` is reachable).
Recommendation: leave as-is (stricter and safer, rationale documented). If strict parity is
later required, export a `cades` accessor for the (k)-free `cmsBaselineBRequirements()` and
call it here, and add a `// DIVERGENCE, deliberate:` entry referencing this finding.

### Findings — PERF

**P13C-PERF-001** — Severity **Medium** — Category **PERF-LEAK**
Location: `pades/pdf_timestamp_token.go:39,61`
Evidence:
```go
var pdfTimestampTokenRegistry = map[*validation.TimestampToken]*PdfTimestampToken{}
func NewPdfTimestampToken(rev *PdfDocTimestampRevision) (*PdfTimestampToken, error) {
	...
	pdfTimestampTokenRegistryMu.Lock()
	pdfTimestampTokenRegistry[base] = token   // insert; grep confirms no delete()/eviction anywhere
	pdfTimestampTokenRegistryMu.Unlock()
```
Impact: a process-wide map that is **write-only** — no `delete()` or eviction exists.
Every `NewPdfTimestampToken` (i.e. every document time-stamp revision extracted while
validating/analysing a PDF) inserts one entry holding a `*PdfTimestampToken` (its
`*validation.TimestampToken` + DER-encoded token binary) and the enclosing
`*PdfDocTimestampRevision` (signature dict, fields, signed content). In a long-running
process validating many PDFs (DSS-as-a-service), the map grows without bound: an unbounded
memory leak retaining full token binaries and certificate/CRL/OCSP data after callers drop
their references. This is a Go-side construction (standing in for Java's
`instanceof PdfTimestampToken`) with no upstream analogue, so it is not a faithful-port
behaviour.
Recommendation: make the mapping reclaimable — e.g. a Go 1.24+ `weak`-map keyed on the
`*validation.TimestampToken`; or scope the registry to a validation job and clear it at job
end; or let the token carry a back-pointer so no process-wide registry is needed. Verify
with a heap-diff leak check over a batch of `NewPDFDocumentValidator(...).Validate()` calls.

### Findings — STD

**P13C-STD-001** — Severity **Low** — Category **STD**
Location: `pades/dss_java_font.go:61-63` (`Name`)
Evidence:
```go
// Name gets the name of the font. Port of #getName.
func (f *DSSJavaFont) Name() string {
	return f.javaFont.Name   // f.javaFont is *NativeJavaFont; nil when built via NewDSSJavaFont(nil)
```
Impact: `NewDSSJavaFont(nil)` (legal — the constructor guards `if javaFont != nil` only for
`size`) yields a `*DSSJavaFont` whose `javaFont` is nil; a later `Name()` dereferences nil
and panics. Java's `DSSJavaFont(Font)` likewise NPEs on `new DSSJavaFont(null).getName()`,
so this is **upstream-faithful** and only reachable by a caller that explicitly passes a nil
font (a font is not attacker-controlled input and the native engine renders nothing, so blast
radius is small). Noted because the sibling constructors are otherwise defensive.
Recommendation: optional hardening — guard `Name()`/`SetSize()` on a nil `javaFont`
(return `""` / no-op) or document the nil contract on `NewDSSJavaFont`. No parity change needed.

### Findings — Info

**P13C-STD-002 (Info)** — `pades/pdf_revision_timestamp_source.go:109-123`
`pdfRevisionTSAddReferences` is an O(n²) linear-scan dedup, the same pattern already
recorded as **P12A-PERF-002** (`padesTSAddReferences`), bounded by the (small)
timestamp-reference set and mirroring Java's `AbstractTimestampSource.addReferences`.
Cross-reference only — not a new defect.

**P13C-SEC-002 (Info)** — `pades/testdata/crossgen/main.go:58,64,172`
The cross-validation generator loads two **test** PKCS#12 keystores with the literal
password `"testpassword"`. Keys are strong (verified: signer RSA **2048-bit**, TSA EC
**P-256**), the TSA policy is an unregistered test OID `1.2.3.4.5.6.7.8.9`, and the
artifacts only feed the Java oracle `CrossGenValidator` for interop proof — no production
path, no real credential. Info / expected for test-only fixtures.

### Lens summaries

- **STD:** no blocking findings. Doc comments are uniformly thorough (port provenance,
  `Port of #…` mapping, and the Go-vs-Java structural deviations are all justified — e.g. the
  `LevelBaselineTOverrides` dispatch shim in `pades_level_baseline_t.go`, the
  `CMSForPAdESBuilderHelper` covariant-setter mirroring in `cms_for_pades_builder_helper.go`,
  the `pdfCMSRevisionBase` shared-state split in `pdf_cms_revision.go`). The only mutex in
  scope (the token registry's) is held correctly. One Low hardening (P13C-STD-001), one Info.
- **PERF:** one Medium (P13C-PERF-001, the write-only token registry) and one Info cross-ref
  (P13C-STD-002). No O(n²) on attacker-controlled input in the tail; the
  `PdfObjectTree.IsProcessedReference` linear scan (`pdf_object_tree.go:95`) is bounded by
  object-tree depth and is the faithful port of upstream's `isProcessedReference`.
- **SEC:** no Critical, no High. The two highest-risk lens targets were checked and are clean /
  upstream-faithful: byte-range arithmetic (`byte_range.go`, `byte_range_input_stream.go`)
  matches Java exactly; the permission checkers (`pdf_permissions_checker.go`,
  `sig_field_permissions.go`) match Java exactly, including the INCLUDE/EXCLUDE fall-through
  and the `default: panic` on an unknown `PdfLockAction`. `pdf_signature_cache.go` is a plain
  DTO (no staleness/leak). The identifier builders derive IDs from field names / VRI names /
  revision binaries — no weak crypto. `revocation_info_archival.go` (ASN.1 EXPLICIT-tag parse)
  is well-guarded (rejects >3 children, wrong tag numbers, non-single-child tags) — no
  over-read. The one real acceptance divergence (P13C-SEC-001) is fail-closed and documented.

### Tool log (U13c)

Run from `dss/` (module root). Go `go1.27.0 darwin/arm64`. All three tools **available**.

```
$ gofmt -l pades pades/alerts pades/exception
(empty — clean)

$ go vet ./pades/...
VET_OK            # exit 0, no diagnostics

$ golangci-lint run --config=../.github/.golangci.yml ./pades/...
0 issues.         # config passed explicitly, not auto-discovered
```

(No fuzz target exists for the tail in this package; the byte-range / permissions /
certificate-source paths are oracle-verified by
`pades_upstream_cross_validation_test.go` / `pades_downstream_cross_validation_test.go`.)

### Open questions (U13c)

1. **P13C-PERF-001 — the real leak.** Confirm the registry is intended to be
   process-lifetime. If a DSS server validates unbounded PDFs, is a `weak`-map or a
   job-scoped registry acceptable, or is there a reason the token must stay reachable for the
   life of the process (a caller holds the bare `*validation.TimestampToken` and expects
   `PdfTimestampTokenOf` to resolve it later)? That determines the fix shape.
2. **P13C-SEC-001 — parity vs. safety.** The (k)-stricter check is documented but not in
   `known-gaps.md`. Should it be added there (it is the only PAdES-B *acceptance* divergence
   in the tail), or is the fail-closed direction acceptable to leave only in the file header?
3. **Overall:** this unit found **no Critical/High** and no new false-accept in the
   verification path — consistent with the U11/U12/P13 sweeps. The tail is dominated by data
   carriers (revisions, scopes, fields, fonts, parameters) and SPI seams whose behaviour was
   already exercised by the oracle tests; the one genuine defect is the Go-only token registry
   (P13C-PERF-001).
