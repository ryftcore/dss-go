# Batch 02 review — XML chain (unit U06)

- **Unit:** U06 — `internal/xmldom`
- **Scope:** the full `dss/internal/xmldom` package — XML DOM model, parser, serializer,
  encoding transcode, ID index (29 files, 6528 lines).
- **Date:** 2026-08-23
- **Depth:** deep (line-level; **every file read in full** — see Files read)
- **Context read (allowed only):** `dss/PORTING.md` conventions, `docs/compatibility/known-gaps.md`.
  Deliberate conventions treated as **not** findings: 1:1 Java/Xerces/OpenJDK port shape,
  `panic` standing in for unchecked exceptions at the mutation API, staticcheck-style
  categories disabled, stdlib-only, `internal/` frozen.

**Headline (SEC lens):** this package is the *deliberately hardened* member of the XML chain.
The attack surface the review was primed to find is **closed by construction**:
DOCTYPE is rejected by default (no internal → no entity expansion / billion-laughs, no
external → no SSRF, no DTD fetching at all), the parser is **iterative with a depth limit of
500** (no stack exhaustion), there are **no regular expressions** (no ReDoS), **no goroutines,
no network, no `unsafe`**, and the serializer **escapes `& < >` (and `"` in attributes)** in
both text and attribute values. The one genuinely exploitable defect found is a
**CPU DoS**, not a memory blow-up: an O(n²) attribute-value scanner (X06-PERF-001).

---

## Files read (29/29, all FULL — none truncated)

| File | Lines | Read |
|---|---|---|
| `serialize.go` | 856 | full (1–320, 320–620, 620–856) |
| `node_test.go` | 508 | full |
| `parse.go` | 419 | full |
| `node.go` | 348 | full |
| `stdlib_pin_test.go` | 313 | full (hostile-input baseline) |
| `ids_test.go` | 312 | full (wrapping-attack baseline) |
| `navigate.go` | 303 | full |
| `parse_test.go` | 295 | full |
| `attvalue.go` | 295 | full (O(n²) site) |
| `encoding.go` | 276 | full |
| `differential_test.go` | 242 | full (mutation-sweep baseline) |
| `serialize_encoding.go` | 241 | full |
| `encoding_test.go` | 221 | full |
| `serialize_test.go` | 207 | full |
| `owner_test.go` | 203 | full |
| `attvalue_test.go` | 187 | full |
| `ids.go` | 186 | full |
| `roundtrip_test.go` | 180 | full |
| `nilsafe_test.go` | 175 | full |
| `fuzz_test.go` | 157 | full (hostile-input baseline) |
| `reject_test.go` | 148 | full (rejection baseline) |
| `text_test.go` | 118 | full |
| `codepage.go` | 78 | full |
| `attrs.go` | 60 | full |
| `nodeset.go` | 48 | full |
| `doc.go` | 47 | full |
| `errors.go` | 39 | full |
| `name.go` | 37 | full |
| `helpers_test.go` | 29 | full |

No gap was left unread in any file. The test files (`fuzz_test.go`, `stdlib_pin_test.go`,
`differential_test.go`, `reject_test.go`, `ids_test.go`) were read first, as they define the
hostile-input coverage baseline the SEC findings are measured against.

---

## Findings

Counts: **1 High, 1 Medium, 3 Low, 3 Info** (8 total). No Critical. No false-accept /
exploitable verification finding in this package.

### X06-PERF-001

- **Severity:** High
- **Category:** PERF-BIGO
- **Location:** `internal/xmldom/attvalue.go:195` (`normalizeAttValue`, the `&` branch),
  mirror at `attvalue.go:282` (`checkTextRefs`)
- **Evidence:**
  ```go
  case c == '&':
      j := strings.IndexByte(string(raw[i:]), ';')   // O(remaining) copy+scan, per reference
      if j < 0 { return "", errRefNoSemi }
      r, err := decodeRef(string(raw[i+1 : i+j]))
      ...
      i += j + 1
  ```
- **Impact:** Each `&` reference in an attribute value allocates `raw[i:]` (a full copy of the
  remaining bytes) and scans it for `;`. For a value with *k* references spread across *n*
  bytes this is O(n²) time and O(n²) transient allocation. It is **attacker-controlled input**
  on the **parse hot path of a signature validator**: `a="&amp;&amp;&amp;…"` (a few KB) is O(n²)
  in that one value. This is the one place a large but *legal* document can be turned into a
  CPU DoS — distinct from the billion-laughs *memory* attack, which is blocked by construction
  (no DTD / no entity expansion at all; see the clean-SEC section below). `checkTextRefs` has
  the same `string(s[i:])` per `&` but is only entered when
  the run contains U+FFFD, so `normalizeAttValue` is the primary exposure.
- **Recommendation:** Replace the per-reference `strings.IndexByte(string(raw[i:]), ';')` with a
  single forward scan that records the next `;` (or use `bytes.IndexByte(raw[i:], ';')` to
  avoid the `string(...)` copy, then advance). Linear in value length; no behavior change.
  Keep the `decodeRef` validation as-is (it is correct).

### X06-PERF-002

- **Severity:** Medium
- **Category:** PERF-MEM
- **Location:** `internal/xmldom/parse.go:53` (`ParseReader`), `parse.go:62` (`io.ReadAll(r)`)
- **Evidence:**
  ```go
  if o.MaxBytes > 0 {
      src, err = io.ReadAll(io.LimitReader(r, o.MaxBytes+1))
  } else {
      src, err = io.ReadAll(r)        // no bound by default
  }
  ```
- **Impact:** `ParseReader` with nil/`MaxBytes==0` options buffers the **entire stream into
  memory with no cap** (`io.ReadAll(r)`). A memory-exhaustion DoS on attacker-controlled input:
  the only protection is that a *caller* sets `MaxBytes`, but the default is unlimited. The
  in-module callers (fuzz cap 1 MiB) are careful, but `ParseReader` is an exported entry point
  whose secure posture depends on every caller remembering to set `MaxBytes`. This is a real
  hardening gap on the hostile-input surface, hence Medium rather than Low.
- **Recommendation:** Give `ParseReader` a documented, enforced default cap (or a dedicated
  bounded-reader helper) so the secure posture does not depend on caller discipline. At minimum,
  a doc comment on `ParseReader` warning that nil options = unbounded, matching how `Parse`
  documents `MaxBytes`, and a note in `docs/compatibility/known-gaps.md` if it is accepted.

### X06-PERF-003

- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `internal/xmldom/serialize.go:318` (`domAttrOrder`)
- **Evidence:**
  ```go
  out := make([]*Node, 0, len(attrs))
  for _, a := range attrs {
      i := sort.Search(len(out), func(i int) bool { return compareUTF16(out[i].Name.QName(), q) >= 0 })
      out = append(out, nil)
      copy(out[i+1:], out[i:])   // O(n) memmove, once per attribute
      out[i] = a
  }
  ```
- **Impact:** Insertion-sort: O(n²) attribute moves per element, where *n* = attribute count.
  Attribute counts are small (tens at most) in XAdES/PAdES, so the constant is negligible;
  reported for completeness because it is the one non-linear step in the serializer and it is
  on the digest-relevant serialize path.
- **Recommendation:** Leave as-is unless a wide-attribute document ever appears; if so, switch
  to `sort.Slice` over `[]*Node` keyed by `compareUTF16(QName)`. Not worth a change today.

### X06-STD-001

- **Severity:** Low
- **Category:** STD
- **Location:** `internal/xmldom/node.go:61` (exported `Attrs` field), contract at `node.go:42`;
  `bump()` at `node.go:344`
- **Evidence:**
  ```go
  type Node struct {
      ...
      Attrs []*Node   // exported, mutable; "Callers must not reorder or mutate it directly"
      ...
  }
  // bump invalidates the owning document's derived indexes.
  func (n *Node) bump() { ... }   // only the documented mutators call it
  ```
- **Impact:** The ID index is invalidated lazily via a generation counter that **only the
  documented mutators** (`AppendChild`/`InsertBefore`/`RemoveChild`/`ReplaceChild`/`SetAttr`/
  `RemoveAttr`/`SetTextContent`) bump. A caller who mutates `n.Attrs` directly (the field is
  exported and the invariant is carried only by a comment) silently defeats the
  `ElementByID`/`IDAttrs`/`DuplicateIDs` invalidation contract, producing **stale ID
  resolutions** — a correctness landmine in exactly the `#id` dereference path the package
  documents as wrapping-attack-sensitive (`ids.go`). This is a latent footgun, not a
  reachable bug in the current codebase.
- **Recommendation:** Either make `Attrs` unexported with accessors, or (to keep the port's
  shape) add a one-line doc note that direct mutation of `Attrs` bypasses the ID-index
  generation counter and is unsafe, so the contract is at least discoverable next to the field.
  Low priority — the current in-tree call sites all go through the mutators.

### X06-STD-002

- **Severity:** Low
- **Category:** STD
- **Location:** `internal/xmldom/serialize_encoding.go:158` and `:165`
  (`writeUTF16BE` / `writeUTF16LE`)
- **Evidence:**
  ```go
  func writeUTF16BE(dst []byte, r rune) []byte {
      for _, u := range utf16.Encode([]rune{r}) {   // allocates []rune{r} per character
          dst = append(dst, byte(u>>8), byte(u))
      }
      return dst
  }
  ```
- **Impact:** Every character write allocates a fresh `[]rune{r}` (and, for a supplementary
  character, a two-unit slice) on the serializer's hottest loop. Pure allocation overhead on the
  UTF-16 output path; no correctness impact. UTF-8/Latin-1 paths do not have this.
- **Recommendation:** Use `utf16.EncodeRune(r)` (returns the two `uint16`s with no slice
  allocation) instead of `utf16.Encode([]rune{r})`. Trivial, local, zero behavior change.

---

### X06-INFO-001 — deliberate JDK bug-compatibility in the serializer (cite, not a defect)

- **Severity:** Info
- **Category:** SEC
- **Location:** `internal/xmldom/serialize_encoding.go` (`writeXalanUTF8`, ~line 96;
  `writeTruncatingASCII`, ~line 135)
- **Evidence:**
  ```go
  // writeXalanUTF8 ... for U+40000 up ... the three high bits of the code point are dropped.
  // U+10FFFF comes out F0 8F BF BF instead of F4 8F BF BF, which is not even valid UTF-8.
  ```
- **Note:** The package *deliberately reproduces* OpenJDK's UTF-8 encoder bug (astral code
  points mis-encoded) and the `US-ASCII` truncating writer, because those exact bytes are what
  `ds:Reference DigestValue` is computed over. This is a documented, goldens-pinned divergence
  (see `serialize.go` header + `xml/utils/testdata/serialize/goldens.txt`). It is the correct
  porting choice for byte-parity and is **not a finding**; listed here so the SEC lens's
  "malformed UTF-8 output" concern is explicitly answered: it is intentional and load-bearing.

### X06-INFO-002 — "last ID wins" is a deliberate parity choice in a wrapping-attack-sensitive path

- **Severity:** Info
- **Category:** SEC
- **Location:** `internal/xmldom/ids.go` (`register`, the `d.ids[attr.Value] = elem` line)
- **Evidence:**
  ```go
  // LAST registration wins ... A same-document "#id" reference therefore dereferences to the
  // last such element. That is exactly the shape of an XML signature wrapping attack ...
  ```
- **Note:** `ElementByID` resolves a duplicated `Id` to the **last** element in document order,
  mirroring Xerces/OpenJDK `CoreDocumentImpl.putIdentifier` (a `HashMap.put`). The package
  documents that this is the *correct* parity behavior and that upstream DSS ships
  `xades-with-manifest-with-duplicated-reference.xml` precisely to exercise it. The actual
  *decision* (whether a duplicate is fatal) is made by the caller — `DuplicateIDs()` reports the
  condition, `ids.go:54` documents that registration never fails. So the wrapping-attack verdict
  is owned by `validation/`/`xmldsig`, not this package. Correctly delegated; **not a finding**
  here. (Cross-check that the DSS caller actually calls `DuplicateIDs()` and rejects is in the
  scope of the `xmldsig`/`xades` units, U09/U14.)

### X06-INFO-003 — no `recover()` in the parse/serialize hot path (acceptable; panic surface is narrow)

- **Severity:** Info
- **Category:** SEC
- **Location:** `internal/xmldom/parse.go`, `serialize.go` (package-wide: no `recover()` outside
  tests)
- **Evidence:** grep for `recover()` in non-test files → none. The only in-package `panic`s are
  on the **mutation** API (`AppendChild`/`RemoveChild`/`SetAttr`/`RegisterIDs` on a non-document,
  etc.), which require a caller to pass a *nil or misparented pointer* — programmer error, not
  hostile *data*.
- **Note:** Parse and Serialize are the hostile-data entry points and they are **panic-free by
  design**: every rejection path returns `*SyntaxError` or an `error`; the only `fail` in the
  serializer is on a hand-built tree with an invalid character, which `Bytes`/`Serialize` surface
  as a returned error (see X06-STD-003). The iterative `parser` loop (`parse.go:103`) and the
  depth guard (`parse.go:155`) mean a hostile *document* cannot overflow the stack. The `FuzzParse`
  gate (`fuzz_test.go`) and the `parseNoPanic` recovery in `differential_test.go` both assert
  "never panic on arbitrary input". This is a clean result, recorded so the "panic on hostile
  input" lens is explicitly closed.

---

## Fuzz / hostile-input coverage assessment

**Existing coverage is strong.** The package ships a real fuzz target *and* a deterministic
mutation floor:

- `FuzzParse` (`fuzz_test.go`) — arbitrary `[]byte`, asserts (a) no panic, (b) structural
  invariants on every parsed tree, (c) every accessor is nil-safe, and (d) Serialize→Parse→
  Serialize idempotence + tree stability. 34 seeds including UTF-16 BOMs, DOCTYPE, entities,
  CDATA splits, namespace edge cases, astral chars, ISO-8859-1.
- `TestParseVsEncodingXMLOnMutants` (`differential_test.go`) — 10,000 deterministic mutants,
  `recover()`-guarded, asserting the parser is *strictly stricter than* `encoding/xml` (containment
  property) and never panics.
- `reject_test.go`, `stdlib_pin_test.go`, `attvalue_test.go`, `encoding_test.go` — pin every
  added well-formedness check, the clause 3.3.3 attribute normalization, and the encoding table.

**Fuzz gaps (hostile-input classes NOT covered by the current fuzz target) — for the follow-up
that adds targets:**

1. **Serializer fuzz target.** `FuzzParse` only fuzzes *parse*; the **serializer is not fuzzed at
   all** as an entry point. There is no `Fuzz(n *Node / []byte tree → Bytes)` that drives
   `Serialize` with arbitrary *hand-built* trees (the `isInvalidXMLChar` path at `serialize.go:844`,
   the `pushNamespace`/`patchName`/`namespaceURIOf` "undeclared prefix" error at `serialize.go:473`,
   and the CDATA-split state machine are only reachable by building nodes by hand). A
   `FuzzSerialize` over a generated `[]Node` tree would close the one un-fuzzed entry point.
2. **`MaxBytes` / large-input DoS is not fuzz-asserted.** The fuzz cap is 1 MiB and the
   mutation sweep caps at 64 KiB, so the O(n²) path in X06-PERF-001 and the unbounded
   `ParseReader` (X06-PERF-002) are exercised only for *crash*, never for *time/alloc*. A
   targeted property test (e.g., an N-KB all-`&amp;` attribute must parse in < T) would pin the
   complexity.
3. **Deep-nesting is tested at exactly the boundary (500/501) but not fuzzed for the *serializer's*
   recursion.** `TestDeepNestingRoundTrip` covers 200; `Walk`/`Clone`/`serialize.walk` are
   recursive and are bounded by the same 500 only *if the tree came from Parse*. A hand-built
   10⁵-deep tree (legal via `AppendChild`, no depth guard on mutation) would recurse in
   `Walk`/`Clone`/`serialize.walk`. See X06-STD-001 / open question.
4. **UTF-16/UTF-8 BOM + encoding combos** are well covered by unit tests but not part of the
   open-ended fuzz seed set beyond a few seeds; a seed table of BOM×declared-encoding pairs as
   fuzz seeds would harden the `decodeSource` branch matrix.

---

## Tool log (run from `dss/`)

| Tool | Command | Result |
|---|---|---|
| gofmt | `gofmt -l internal/xmldom` | clean (no files printed) |
| go vet | `go vet ./internal/xmldom/...` | clean (exit 0) |
| golangci-lint | `golangci-lint run --config=../.github/.golangci.yml ./internal/xmldom/...` | `0 issues.` |
| go test | `go test ./internal/xmldom/ -count=1 -timeout 300s` | `ok … 0.645s` (all unit tests pass) |
| fuzz | `go test ./internal/xmldom/ -run NONE -fuzz FuzzParse -fuzztime 15s` | `PASS` — 670,201 execs, 0 crashes, 181 interesting |
| (none unavailable) | — | all five tools present and run |

No tool was unavailable. Nothing in the package is skipped by `corpustest` (the package has its
own `testdata` and no `corpus/`-gated tests in scope here).

---

## Open questions / observations

- **Recursion on hand-built deep trees (Medium potential, currently Low):** Parse enforces
  `MaxDepth=500`, so a *parsed* tree cannot exceed it. But `AppendChild`/`InsertBefore` carry no
  depth guard, and `Walk`, `Clone(deep)`, `TextContent`, `serialize.walk`, and
  `NodeSet.AddSubtree` are all recursive. A caller who builds (or imports) a 100k-deep element
  chain and then walks/serializes it can overflow the goroutine stack. This is out of scope for
  "hostile *document* input" (Parse blocks it) and is a *mutation-API* concern, so it is noted
  rather than raised as a finding; if DSS ever builds deep chains programmatically (it does not
  today) this would become a real DoS. Cheap mitigation if ever needed: an iterative
  work-stack version of `Walk`/`Clone`/`serialize.walk`.
- **`ParseReader` unbounded by default (X06-PERF-002):** confirm every DSS call site sets
  `MaxBytes` (or feeds `[]byte` via `Parse`) before this is treated as closed. The in-module
  callers do; the public API does not enforce it.
- **Serializer output is intentionally not re-parseable in the astral-UTF-8 case (X06-INFO-001):**
  the mis-encoded bytes are the digest input, so downstream code must never re-`Parse` a
  `UTF-8`-declared document that contains astral characters and expect valid UTF-8. This is a
  property of the ported serializer, not a bug — but it is a sharp edge for any future tooling
  that treats serializer output as a general XML stream.
- **Cross-check (outside this unit):** the two SEC-relevant *decisions* this package defers —
  "is a duplicate `Id` fatal?" and "is a namespace-confused attribute a verification error?" —
  are owned by `validation/` and `internal/xmldsig`. Verify in U09 that `DuplicateIDs()` is
  consulted and that the `register()` last-wins behavior is what the DSS wrapping-attack test
  expects (the `ids.go` comment says yes, and the corpus file is pinned).
- **`bump()` is called on the parent, not the subtree:** mutation of a node's own attributes via
  `SetAttr`/`RemoveAttr` bumps the *document* generation (correct), but `SetTextContent` detaches
  children via direct link-clearing rather than `RemoveChild`; it still calls `bump()`, so the
  ID index is invalidated. Verified by `TestIDIndexIsInvalidatedByMutation/SetTextContent`. No
  action.

---

**No findings (explicitly clean per lens):**
- **SEC — XXE / SSRF / billion-laughs / ReDoS / namespace-confusion / serializer escape gaps:**
  all **clean**. DOCTYPE is rejected by default (`parse.go:directive`); no DTD, no external
  entities, no parameter entities, no entity expansion at all → billion-laughs is structurally
  impossible; no `net/http`, no I/O of any kind → SSRF impossible; no `regexp` package used
  anywhere in non-test code → ReDoS impossible; the serializer escapes `& < >` (and `"` in
  attribute values) and `>`/`]]>` are legal in attribute values and handled (`serialize_test.go`,
  `reject_test.go`); the parser rejects duplicate expanded names and undeclared prefixes
  (`buildAttrs`, `declareNamespaces`).
- **PERF-LEAK:** **no findings.** No goroutines, no channels, no timers, no open handles, no
  `io.ReadCloser` obtained without a matching read-to-completion, no grow-only unbounded state in
  the hot path (the ID index is rebuilt, not accumulated).
- **STD (error-handling, defer, context, doc comments, test quality):** the package is
  exemplary — every exported symbol has a doc comment, errors are typed (`*SyntaxError` with
  line/col/offset), `ParseReader` propagates read errors, the mutation API panics *only* on
  programmer error (nil/misparented) and is covered by `TestMutationPanics`, nil-receiver safety
  is covered by `TestNilReceiversAreSafe`, and the test suite is table-driven, deterministic,
  and oracle-pinned. `Serialize` correctly returns `Bytes`'s error and writes nothing on
  failure (`serialize.go:63-70`), so there is no silent-write path. The two Low STD findings
  above (X06-STD-001, X06-STD-002) are the only nitpicks.

---

# Batch 02 review — XML chain (unit U07)

- **Unit:** U07 — `internal/xpath10` (14 files, 2,893 lines) + `xml/utils` (14 files, 3,426 lines)
- **Scope:** the XPath 1.0 subset engine (lexer, parser, evaluator, transform widening) and the
  DSS bridge package that compiles/executes DSS's `XPathQuery` objects, plus DomUtils/XPathUtils,
  NamespaceContextMap, XMLCanonicalizer, DOMDocument and the executor layer.
- **Date:** 2026-08-24
- **Depth:** deep (line-level; see Files read)
- **Context read (allowed only):** `dss/PORTING.md`, `docs/compatibility/known-gaps.md`.
  Deliberate conventions treated as **not** findings: 1:1 Java port shape, `panic` standing in
  for Java unchecked exceptions/requireNonNull, staticcheck-style categories disabled,
  stdlib-only, `internal/` frozen, unregistered-prefix → null-namespace (pinned Java quirk).
  **Not re-reported from U06:** anything about `internal/xmldom` (node model, parser, serializer,
  ID index) — reviewed in this file's U06 section.

**Headline (SEC lens):** this is the *correctness-critical* member of the XML chain — signature
verification depends on node-set ordering, axis direction, and string/number coercion being
exactly Xalan/JDK. The engine is **spec-faithful by construction and oracle-pinned**: results
are deduped by pointer and sorted into document order via a lazily built document-order index
(`eval.go:267-323`), the attribute axis correctly skips `xmlns` declarations
(`eval.go:161-170`), axes are implemented in the right direction (parent/ancestor walk
`n.Parent` chain, descendant walks children — `eval.go:183-213`), string/boolean coercion
follows clause 3.4/4.2/4.3 (`eval.go:355-407`, `transform.go:254-271`), and `equals` is
EXISTENTIAL over node-sets as clause 3.4 requires (`eval.go:366-379`). The only regular
expression in either package is the anchored-linear `\s` split in `dom_utils.go:58` (no ReDoS).
There are **no Critical or High findings**. The engine's one structural risk — recursive
descent over an attacker-sized document — is **bounded by xmldom's parser depth cap of 500**
(X07-INFO-001). The real exposure the review was primed to find — an **XPath-injection surface**
where a document-derived `Id` value is interpolated into a string literal **without escaping** —
lives in `xml/common` (out of scope here; tracked as open question OQ-3 for unit U10), and is
**parity with Java** (upstream `XPathQueryAttributeParameter` interpolates identically), so it
is not reported as a new finding.

---

## Files read (28/28; 1 budgeted-truncated)

| Package | File | Lines | Read |
|---|---|---|---|
| xpath10 | `doc.go` | 92 | full |
| xpath10 | `xpath10.go` | 84 | full |
| xpath10 | `ast.go` | 112 | full |
| xpath10 | `errors.go` | 48 | full |
| xpath10 | `namespace.go` | 91 | full |
| xpath10 | `lexer.go` | 237 | full (no regex; linear scan) |
| xpath10 | `parser.go` | 469 | **truncated** — lines 1-300 + grep of every `func` + the security-relevant parse regions (300-469: steps, node tests, axes) |
| xpath10 | `eval.go` | 407 | full (evaluator, axes, ordering, coercion) |
| xpath10 | `transform.go` | 272 | full (transform widening, `id()`, `NamespaceContextOf`) |
| xpath10 | `adversarial_test.go` | 184 | full (hostile-input baseline) |
| xpath10 | `xpath10_test.go` | 388 | full |
| xpath10 | `kat_test.go` | 304 | full |
| xpath10 | `transform_kat_test.go` | 299 | full |
| xpath10 | `namespace_test.go` | 111 | full |
| xml/utils | `doc.go` | 58 | full |
| xml/utils | `dom_utils.go` | 865 | **truncated** — lines 1-300 + grep of all 51 `func` + security regions (Id/XPointer/GetDate/GetNodeBytes/javaString/deepcopy: 300-865) |
| xml/utils | `serialize_oracle_test.go` | 510 | full (byte-parity gate) |
| xml/utils | `dom_utils_test.go` | 359 | full |
| xml/utils | `dom_utils_oracle_test.go` | 229 | full (XPointer/date oracle) |
| xml/utils | `xml_canonicalizer.go` | 182 | full |
| xml/utils | `dom_document.go` | 181 | full |
| xml/utils | `xpath_utils.go` | 167 | full (global registry) |
| xml/utils | `namespace_context_map.go` | 118 | full |
| xml/utils | `dom_document_test.go` | 92 | full |
| xml/utils | `namespace_context_map_test.go` | 87 | full |
| xml/utils | `xml_canonicalizer_test.go` | 86 | full |
| xml/utils | `xpath_query_executor_loader.go` | 67 | full |
| xml/utils | `java_xml_xpath_query_executor.go` | 60 | full |
| xml/utils | `native_dom_xpath_query_executor.go` | 45 | full |
| xml/utils | `santuario_initializer.go` | 32 | full |
| xml/utils | `dss_xml_error_listener.go` | 29 | full |
| xml/utils | `xpath_string_executor.go` | 20 | full |
| xml/utils | `xpath_query_executor.go` | 20 | full |
| xml/utils | `abstract_xpath_query_executor.go` | 14 | full |

Cross-package reads (allowed, for SEC context, not findings): `internal/xmldsig/transform_xpath.go`,
`transform_xpath2.go`, `signature.go` (confirm `RegisterIDs` precondition is satisfied in the
production path via `xades/xades_signature.go:1086` → `RecursiveIdBrowse` → `RegisterIDs`, so
`id()` resolution is **not** a finding); `xades/dss_xml_utils.go` (concurrency of the global
registry); `xml/common/xpath_query_attribute_parameter.go` + `xpath_query_identifier_parameter.go`
+ `abstract_xpath_query*.go` (the unescaped-literal surface — tracked as OQ-3, not a finding).

---

## Findings

Counts: **0 Critical, 0 High, 0 Medium, 3 Low, 3 Info** (6 total).
No false-accept / exploitable verification finding in either package. The three Lows are
concurrency/robustness hardening in `xml/utils`; the three Infos are structural notes on the
engine. No STD/PERF finding in `internal/xpath10` itself — it is clean on every lens.

### X07-STD-001

- **Severity:** Low
- **Category:** STD
- **Location:** `xml/utils/xpath_utils.go:19-23` (globals) + `namespace_context_map.go:33-39`
  (`RegisterNamespace`), exercised through `xades/dss_xml_utils.go:127-137`,
  `asic/asic_manifest_parser.go:16-19`, `tsl/..._with_sha2_predicate.go:35`
- **Evidence:**
  ```go
  xPathUtilsNamespacePrefixMapper = NewNamespaceContextMap()   // package global
  func (m *NamespaceContextMap) RegisterNamespace(prefix, namespace string) bool {
      _, existed := m.prefixMap[prefix]
      m.prefixMap[prefix] = namespace        // unsynchronized write
      m.createNamespace(prefix, namespace)   // + nested map write
      return !existed
  }
  ```
- **Impact:** The prefix registry is a process-global mutable map. `RegisterNamespace` (writes)
  and `PrefixMap`/`NamespaceURI` (reads, via `namespaceContextMapToXPath10` at
  `namespace_context_map.go:113-118`) are called from multiple format packages. Upstream
  `NamespaceContextMap` guards this with `synchronized`. A caller registering a namespace
  concurrently with another goroutine compiling an XPath expression is a data race (undefined
  behavior; `go test -race` would catch it if exercised). In the typical single-threaded
  validator flow this never fires, and the writes are idempotent re-registrations of fixed URIs,
  so it is a hardening gap, not an exploitable defect.
- **Recommendation:** Add a `sync.RWMutex` to `NamespaceContextMap` (write-lock in
  `RegisterNamespace`, read-lock in `PrefixMap`/`NamespaceURI`/`Prefix`/`Prefixes`), mirroring
  the `sync.RWMutex` that `xml_canonicalizer.go:53` already correctly applies to its own
  registry.

### X07-STD-002

- **Severity:** Low
- **Category:** STD
- **Location:** `xml/utils/xpath_query_executor_loader.go:31-37` and `:56-62`
  (`GetXPathQueryExecutor` / `GetXPathStringExecutor`)
- **Evidence:**
  ```go
  func (l *XPathQueryExecutorLoader) GetXPathQueryExecutor() XPathQueryExecutor {
      if l.xPathQueryExecutor == nil {
          l.xPathQueryExecutor = l.loadXPathQueryExecutor()   // unsynchronized lazy init
      }
      return l.xPathQueryExecutor
  }
  ```
- **Impact:** Textbook unsynchronized lazy-init on a shared loader instance (the global
  `xPathUtilsExecutorLoader` at `xpath_utils.go:22`). Two goroutines calling
  `GetXPathQueryExecutor` first can both observe `nil` and both write the field — a data race.
  Same class as X07-STD-001; also low-likelihood in practice because the executor is
  stateless-after-construction and the write is a benign double-assignment.
- **Recommendation:** Either initialize `xPathQueryExecutor`/`xPathStringExecutor` at
  construction (`NewXPathQueryExecutorLoader`) — the loader's only other path is
  `Set*`, so lazy init buys nothing — or guard the two getters with a `sync.Once`/mutex.

### X07-STD-003

- **Severity:** Low
- **Category:** STD
- **Location:** `xml/utils/xml_canonicalizer.go:125-131` (`CanonicalizeBytesTo`)
- **Evidence:**
  ```go
  func (c *XMLCanonicalizer) CanonicalizeBytesTo(toCanonicalizeBytes []byte, w io.Writer) error {
      out, err := c.CanonicalizeBytes(toCanonicalizeBytes)
      if err != nil { return err }
      _, err = w.Write(out)      // partial write is not detected
      return err
  }
  ```
- **Impact:** `w.Write` is allowed to return a short `n < len(out)` with `err == nil`; the
  function then reports success on a truncated canonicalized output. `CanonicalizeNodeTo`
  (`xml_canonicalizer.go:163-169`) avoids this by delegating to `xmlc14n.Canonicalize`'s
  streaming writer, but the `...BytesTo` overload does not. Callers that write canonical bytes
  into a pipe/socket could accept a silently-truncated canonical form.
- **Recommendation:** Wrap `w` in `io.WriteFull` (Go 1.27) or loop until `n == len(out)`, and
  return an error on a short write.

### X07-INFO-001

- **Severity:** Info
- **Category:** PERF-BIGO
- **Location:** `internal/xpath10/eval.go:209-213` (`appendDescendants`),
  `eval.go:298-316` (`buildOrder`/`walk`), `parser.go:44,208,337` (recursive descent)
- **Evidence:** the evaluator and the parser both recurse over tree depth (a `//x` step and a
  predicate each add one frame per level); there is no explicit depth limit in `xpath10`
- **Impact:** A maliciously deep XML document would push Go's stack deep. **This is bounded, not
  unbounded:** `internal/xmldom.Parse` is iterative with a hard depth cap of 500 (U06 headline),
  so the maximum recursion depth in `xpath10` is ~500 levels — a few hundred KB of stack, well
  within Go's per-goroutine limit and not reachable to exhaustion. Recorded so the bound is
  attributed to the right package and not re-investigated. No action.
- **Recommendation:** None. If a future change ever lets `xpath10` evaluate a tree built without
  going through `xmldom.Parse` (e.g. a hand-assembled `*xmldom.Node` chain), revisit and add a
  depth guard.

### X07-INFO-002

- **Severity:** Info
- **Category:** PERF-BIGO
- **Location:** `internal/xpath10/eval.go:366-379` (`equals`, node-set × node-set case)
- **Evidence:**
  ```go
  case lok && rok:
      for _, a := range lset { for _, b := range rset {
          if stringValue(a) == stringValue(b) { return true }
      } }
  ```
  where `stringValue` (`eval.go:340-349`) for an element is `n.TextContent()` — a full subtree walk
- **Impact:** A `node-set = node-set` comparison where both sides are large element sets is
  O(|L|·|R|·subtree). This is spec-correct (clause 3.4 existential semantics) and **not reachable
  through DSS**: every expression in `testdata/expressions.txt` and
  `testdata/transform-expressions.txt` compares a node-set against a *literal* or a *single
  function result* (`@*[local-name()='Id']='x'`, `name()='ds:Signature'`), never node-set ×
  node-set. Recorded for completeness; no corpus expression reaches the quadratic path.
- **Recommendation:** None.

### X07-INFO-003

- **Severity:** Info
- **Category:** PERF-LEAK
- **Location:** `internal/xpath10/xpath10.go:44-51` (`Evaluate`), `transform.go:69-75`
  (`EvaluateBoolean`), `eval.go:31-37` (`evaluator`)
- **Evidence:**
  ```go
  func (e *Expr) Evaluate(ctx *xmldom.Node) ([]*xmldom.Node, error) {
      ...
      ev := &evaluator{}                     // fresh per call; order map built lazily
      return ev.sortUnique(toNodeSet(ev.eval(e.root, ctx))), nil
  }
  ```
- **Impact:** Every `Evaluate`/`EvaluateBoolean` allocates a new `evaluator` and, when the result
  needs ordering, a fresh `order map[*xmldom.Node]int` sized to the whole document
  (`eval.go:298-316`). A `ds:XPath` transform calls `EvaluateBoolean` once per candidate node
  (`xmldsig/transform_xpath.go:71-83`), so a large document re-walks the tree for each node. The
  engine is **correct and leaks nothing** (all state is call-scoped, no goroutines, no handles),
  but the per-call document walk is the hot-path cost. It is a deliberate consequence of the
  per-node boolean evaluation model upstream Santuario uses, and `Expr` is documented as
  compile-once/evaluate-many — so this is an inherent cost of the ported model, not a bug.
- **Recommendation:** If profiling ever shows the transform filter dominating, hoist a shared
  document-order index out of `evaluator` into `Expr` (build once per document) — a
  pure optimization, behavior-preserving.

---

## No findings (explicitly clean per lens)

- **SEC — node-set ordering:** **clean.** Results are deduped by pointer and sorted into document
  order by a lazily built document-order index that numbers element→attributes→children
  (`eval.go:245-323`), matching XPath 1.0 clause 5 and Xalan's DTM. Pinned by
  `TestResultIsDeduplicatedAndOrdered` and the 12 000-answer `TestKnownAnswers`/
  `TestTransformKnownAnswers` oracle (order-compared, not set-compared). No false-accept path.
- **SEC — axis direction:** **clean.** `parent`/`ancestor` walk `n.Parent` (nearest-first,
  `eval.go:183-207`); `child`/`descendant` walk `FirstChild`/`NextSibling`; the attribute axis
  skips `xmlns` declarations (`eval.go:161-170`); `matches`/`isPrincipal` correctly gate name
  tests on node kind (`eval.go:215-237`). No axis silently returns the wrong direction.
- **SEC — coercion / string-value:** **clean.** `toBool`/`toString`/`stringValue`/`equals`
  implement clauses 3.4/4.2/4.3/5 for the three types the subset produces (`eval.go:327-407`,
  `transform.go:254-271`); `equals` is existential over node-sets (`eval.go:366-388`). The only
  numeric path (`number()`) is refused at parse time, so there is no number-coercion divergence.
- **SEC — `id()`/`idref()`:** **clean.** `evalID` (`transform.go:227-247`) consults
  `doc.ElementByID`, the index `xmldom.RegisterIDs` builds; the production XAdES path calls
  `RecursiveIdBrowse()` before `NewXMLSignature` (`xades/xades_signature.go:1086`), so the index
  is populated. No confusion between the `id()` function and `getElementById` (the latter goes
  through the `local-name()` predicate, a different, literal-based path).
- **SEC — namespace axis / `xml:base`:** **clean / N/A.** The `namespace::` axis is refused by
  name (`parser.go:455-465`), so there is no namespace-node identity to get wrong; `xml:base`
  is not consulted by this engine (base-URI resolution lives in `xmldsig`/`xmlc14n`, out of
  scope). `NamespaceContextOf` (`transform.go:84-110`) reproduces `DOMNamespaceContext`
  innermost-wins and binds `xml` always; the `defaultNamespaceKey` leniency for a stray leading
  colon is a *documented, oracle-pinned* transform-grammar accommodation, not a bypass.
- **SEC — panics on hostile input:** **clean.** The lexer is a linear byte scan (no regex, no
  backtracking → no ReDoS); the parser is bounded recursive descent over a token list; the
  evaluator returns errors (`*SyntaxError`/`*UnsupportedError`/`*EvalError`) rather than panicking
  on malformed or unsupported input. The only deliberate `panic`s are the
  `requireNonNull`-style precondition failures (`namespace.go:46`, `dom_utils.go:114,145`) and the
  unreachable `unknown AST node` guard (`eval.go:67`) — both programmer-error, both documented.
  Hostile-input coverage is pinned by `adversarial_test.go` (deep 60-level tree, prefix
  re-binding, default-namespace separation) and the transform adversarial KAT.
- **SEC — ReDoS:** **clean.** The only regex in either package is the anchored-linear `\s`
  whitespace split at `xml/utils/dom_utils.go:58`; `domUtilsJavaSplit` trims the trailing run it
  introduces (load-bearing for the XPointer regression, pinned by
  `dom_utils_oracle_test.go`). No other `regexp` use.
- **PERF-LEAK:** **no findings** beyond X07-INFO-003 (deliberate model cost). No goroutines, no
  channels, no timers, no open handles, no grow-only unbounded state. `XMLCanonicalizer`'s
  registry is correctly mutex-protected; `DOMDocument`'s caches (`bytesData`, `digestMap`) are
  bounded by the document and its digest set.
- **STD (error handling, defer, context, doc comments, test quality) — `internal/xpath10`:**
  **clean.** Typed errors with offset + the offending expression; `Evaluate(nil)` returns an
  `*EvalError`; every exported symbol documented; test suite table-driven, deterministic, and
  oracle-pinned against `javax.xml.xpath`/JDK XPath. (The three Low STD findings above are all in
  `xml/utils` concurrency/robustness, not in the engine.)

---

## Tool log (run from `dss/`)

| Command | Result |
|---|---|
| `gofmt -l internal/xpath10 xml/utils` | clean (no output) |
| `go vet ./internal/xpath10/... ./xml/utils/...` | clean |
| `golangci-lint run --config=../.github/.golangci.yml ./internal/xpath10/... ./xml/utils/...` | `0 issues.` |
| `go test ./internal/xpath10/... ./xml/utils/... -count=1 -timeout 10m` | `ok ... xpath10 2.068s` / `ok ... xml/utils 0.848s` |

All tools available and clean. No `eaa` build-tag surface in either package.

---

## Open questions / observations

- **OQ-1 (cross-unit, `xml/utils` → `xades`/`asic`):** X07-STD-001/002 are latent data races on
  the global namespace registry and executor loader. Confirm with the integrator whether the
  library is ever exercised from multiple goroutines (e.g. parallel signature validation). If
  yes, these two become the top remediation items in `xml/utils`; if single-threaded is a
  guaranteed usage contract, document it and close as Info.
- **OQ-2 (parity, `xml/utils`):** `DOMDocument` (`dom_document.go`) carries its own copy of
  `CommonDocument`'s digest/serialization plumbing because `model`'s is unexported
  (flagged in the file header). If more `DSSDocument` implementations accumulate outside
  `model`, consider exporting a helper surface rather than growing copies.
- **OQ-3 (SEC, cross-unit → U10 `xml/common`):** the `Id`-lookup path interpolates a
  document-derived value into an XPath string literal **without escaping the XPath quote
  character**: `xml/common/xpath_query_attribute_parameter.go:107` builds
  `@*[local-name()='Id']='<value>'` and `dom_utils.go:475-477`
  (`DomUtilsGetXPathByIdAttribute`) does the same for the three Id spellings. An `Id` attribute
  whose value contains `'` can therefore break out of the literal and change (or fail to
  compile) the expression. This is **parity with Java** (upstream interpolates identically, so
  it is not a new finding under PORTING.md's upstream-tracking rule), but it is a genuine
  hardening opportunity that should be evaluated in the `xml/common` review (unit U10) — e.g. by
  escaping `'` as `''` or, better, by having `xpath10` compare against a parameter rather than a
  re-lexed literal. Tracked here because it is the single most security-relevant string this
  unit hands to the engine.
- **OQ-4 (deliberate, `xml/utils`):** `DomUtilsGetDate` (`dom_utils.go:368-385`) intentionally
  parses only `xsd:dateTime` (rejecting Java's seven other lexical forms and the comma
  fractional separator). The rejection of the comma is the *safe* direction (Java accepts a
  form Go refuses) and is oracle-pinned; the seven-form gap is schema-invalid at every DSS call
  site and JVM-timezone-dependent upstream, so it is a deliberate, documented non-port. No action.


---

## internal/xmlc14n (unit U08)

- **Unit:** U08 — `internal/xmlc14n`
- **Scope:** the full `dss/internal/xmlc14n` canonicalizer — traversal, namespace symbol
  table, `xml:*` attribute stack, C14N 1.0/1.1/exclusive/physical emitters, attribute
  ordering, escaping, and the C14N 1.1 `xml:base` join (13 source + 8 test files, 3740
  lines).
- **Date:** 2026-08-24
- **Depth:** deep (line-level; **every file read in full** — see Files read).
- **Context read (allowed only):** `dss/PORTING.md` conventions,
  `docs/compatibility/known-gaps.md`. Deliberate conventions treated as **not** findings:
  1:1 Santuario port shape, **bug-for-bug fidelity to Santuario** (the epilog-drop quirk,
  the relative-namespace check scoping, the UTF-16 attribute order, the whole C14N 1.1
  `xml:base` decode-then-requote surface, `removeMappingIfRender` always returning
  `false`), `SA1019`/staticcheck-style categories disabled, stdlib-only, `internal/`
  frozen. The `xmlbase.go` `ErrXMLBaseUnjoinable` aborts (opaque/authority-only base) are
  **deliberate parity** with Santuario's unchecked exceptions, not defects.

**Headline (SEC lens):** this is the **verification boundary** — the canonical bytes are the
digest input, so any divergence from Santuario is a signature/verification mismatch. On
correctness the package is **exemplary and clean**: no findings. Every algorithm surface is
pinned against Java goldens generated by `testdata/gen/C14nOracle.java` (`kat_test.go`,
including a golden-`sha256` self-check so a corrupt regeneration can't hide behind a matching
Go bug), a `FuzzCanonicalize` round-trip harness, an end-to-end `CanonicalizeBytes` call-path
test, and table-driven unit tests for the comparator, the namespace stack, the `xml:*`
stack, `joinURI`, and the escape writers. The two things that most often break a canonicalizer
— attribute ordering (Java UTF-16 code-unit order, astral vs U+FFFD pinned) and namespace
declaration suppression (superfluous re-declaration, apex default-namespace undeclaration) —
are both KAT-locked. The findings below are **performance and idiom**, not correctness.

---

### Files read (21/21, all FULL — none truncated)

The only file over 400 lines (`xmlbase.go`, 536) was read in full as two contiguous chunks
(1–300, 300–536). No file was read by signature+region only.

| File | Lines | Read |
|---|---|---|
| `xmlbase.go` | 536 | full (1–300, 300–536) |
| `walk.go` | 365 | full |
| `nsstack.go` | 241 | full |
| `c14n.go` | 182 | full |
| `xmlattrs.go` | 175 | full |
| `attrsort.go` | 151 | full |
| `algorithm.go` | 150 | full |
| `emit_exclusive.go` | 148 | full |
| `emit_inclusive.go` | 124 | full |
| `escape.go` | 129 | full |
| `filter.go` | 96 | full |
| `doc.go` | 55 | full |
| `emit_physical.go` | 35 | full |
| `kat_test.go` | 343 | full (golden KAT harness) |
| `nsstack_test.go` | 207 | full |
| `c14n_test.go` | 178 | full (API/concurrency) |
| `xmlbase_test.go` | 162 | full (defensive `xml:base` paths) |
| `algorithm_test.go` | 144 | full (registry/defaults/PrefixList) |
| `attrsort_test.go` | 131 | full (ordering + UTF-16) |
| `escape_test.go` | 93 | full (escape tables, invalid UTF-8) |
| `xmlattrs_test.go` | 95 | full (`xml:*` stack) |

No gap was left unread in any file.

---

### Findings

Counts: **0 Critical, 1 High, 1 Medium, 2 Low, 1 Info** (5 total). **No SEC findings** — the
c14n correctness surface is clean and oracle-pinned. All findings are performance/idiom.

### X08-PERF-001

- **Severity:** High
- **Category:** PERF-BIGO
- **Location:** `internal/xmlc14n/attrsort.go:130` (`attrSet.add`), consumed by every
  emitter's attribute loop (`emit_inclusive.go`, `emit_exclusive.go`, `emit_physical.go`)
- **Evidence:**
  ```go
  func (s *attrSet) add(a outAttr) {
      i := sort.Search(len(s.items), func(i int) bool { return attrCompare(s.items[i], a) >= 0 })
      if i < len(s.items) && attrCompare(s.items[i], a) == 0 {
          return // TreeSet.add: an equal element is already present.
      }
      s.items = append(s.items, outAttr{})
      copy(s.items[i+1:], s.items[i:])   // O(n) shift on every insert
      s.items[i] = a
  }
  ```
- **Impact:** Each `add` is O(log n) search + O(n) element shift, so building the attribute
  set for an element with **A** attributes is **O(A²)**. The parser imposes **no
  per-element attribute-count cap** — only `MaxDepth=500` (nesting, not attributes) and
  `MaxBytes` (default **unlimited**) — so a *flat* element carrying ~10⁵ attributes (a few
  hundred KB of input, trivially submittable) costs ~10¹⁰ element-copies to canonicalize.
  This is a CPU DoS **on the verification boundary**: any attacker-supplied document that is
  opened/verified is canonicalized, and DSS's own signing path canonicalizes the same way.
  Notably this is a *regression* versus upstream: Santuario collects into a `TreeSet<Attr>`,
  which is O(A log A). The output bytes are identical, so this is purely a complexity gap the
  port introduced by substituting a slice + manual shift for the balanced tree.
- **Recommendation:** Replace the per-insert shift with build-then-sort: `append` all
  attributes (O(A)), `sort.Slice` once with `attrCompare` (O(A log A)), then a single linear
  dedup pass for the `TreeSet.add` "first wins" rule. This restores upstream's O(A log A)
  without changing any output bytes. If the slice is kept for shape, at minimum document the
  O(A²) bound and add a `MaxAttributes` guard at the parser seam.

### X08-PERF-002

- **Severity:** Medium
- **Category:** PERF-MEM
- **Location:** `internal/xmlc14n/escape.go:116` (`runes`), called by all three hot writers
  at `escape.go:36`, `escape.go:59`, `escape.go:81`
- **Evidence:**
  ```go
  func writeTextEscaped(w *bufio.Writer, s string) {
      for _, r := range runes(s) {        // runes() materialises a []rune
          switch r { ... }
  }
  // runes decodes s into a heap slice before the writer can touch it
  func runes(s string) []rune {
      out := make([]rune, 0, len(s))
      for i := 0; i < len(s); {
          r, size := utf8.DecodeRuneInString(s[i:])
          if r == utf8.RuneError && size == 1 { out = append(out, invalidRune); i++; continue }
          out = append(out, r); i += size
      }
      return out
  }
  ```
- **Impact:** Every text node, attribute value, comment, and PI value in the document is
  decoded into a **fresh `[]rune`** in the hottest inner loop — one heap allocation per node
  plus a full extra decode pass over the content. Since `rune` is 4 bytes, a single large
  text node of N bytes temporarily allocates ~4N bytes, and a document with many small text
  nodes adds proportional GC pressure. The only reason the slice exists is to tell an
  invalid byte (`RuneError`, size 1) apart from a well-formed U+FFFD, which is a two-line
  inline check, not a reason to materialize the whole sequence.
- **Recommendation:** Range over the bytes directly
  (`for i:=0; i<len(s); { r,size:=utf8.DecodeRuneInString(s[i:]); if r==utf8.RuneError && size==1 { w.WriteByte('?'); i++; continue }; writeRune(w,r); i+=size }`)
  and drop `runes`. Same output, no per-node allocation, no extra pass.

### X08-PERF-003

- **Severity:** Low
- **Category:** PERF-BIGO
- **Location:** `internal/xmlc14n/xmlbase.go:459`, `:475` (`output.String()` in-loop),
  `:515` (`dotDotOutput` copies the buffer)
- **Evidence:**
  ```go
  case strings.HasPrefix(input, "../"):
      input = input[3:]
      if output.String() != "/" {   // full-buffer copy, once per `..` segment
          output.WriteString("../")
  }
  ...
  func dotDotOutput(output *strings.Builder, input string) string {
      s := output.String()          // full-buffer copy on every 2C branch
      ...
  }
  ```
- **Impact:** `strings.Builder.String()` copies the entire accumulated buffer, and it is
  called on the `../`/`..` branches inside the segment loop and again in `dotDotOutput`. A
  base URI with many dot-segments is therefore O(len²). This sits on the C14N 1.1 `xml:base`
  join path only, and `xml:base` values are short URIs in practice, so the practical cost is
  bounded — hence Low rather than High. It is a deviation from the "avoid repeated string
  building in hot loops" goal and is avoidable.
- **Recommendation:** Replace the `output.String() != "/"` guard with a `rootOnly bool` set
  when the builder holds exactly `/`, and in `dotDotOutput` operate on the builder's
  contents via a tracked length / last-segment rather than re-materializing the whole string.
  No output change.

### X08-STD-001

- **Severity:** Low
- **Category:** STD
- **Location:** `internal/xmlc14n/filter.go:92` (`noteFilterErr`), `c14n.go:84` (`bufio`
  wrapping), `c14n.go:118` (return path)
- **Evidence:**
  ```go
  // noteFilterErr records the first filter failure ... the error is made sticky instead
  // and Canonicalize reports it in place of the (discarded) output.
  func (e *engine) noteFilterErr(err error) { if e.filterErr == nil { e.filterErr = err } }
  ...
  bw := bufio.NewWriter(w)
  ...
  if e.filterErr != nil {
      return e.filterErr           // bw is never Flush()'d on this path
  }
  return bw.Flush()
  ```
- **Impact:** On a filter error the traversal **continues** and the output is discarded, but
  `bufio.Writer` auto-flushes to the caller's `io.Writer` every time its 4 KiB buffer fills.
  So a caller who passes a **streaming** `io.Writer` (not `*bytes.Buffer`) can receive up to
  ~4 KiB of the would-be-dropped canonical bytes *before* the error is returned, and must
  then detect the error to know that stream is garbage. All in-tree callers use the
  `CanonicalizeToBytes`/`CanonicalizeNode` form (in-memory, discarded on error), so this is
  latent, but it is a sharp edge on the exported `Canonicalize(alg, in, io.Writer)` API that
  is only guarded by the caller checking `err`.
- **Recommendation:** Either document on `Canonicalize` that a caller using a streaming
  writer must discard any output already written when the returned error is non-nil, or have
  the error path write to a throwaway sink once a filter error is noted so no partial bytes
  reach the caller. Low priority — no in-tree caller is affected.

### X08-INFO-001 — correctness is oracle-pinned; the "verification boundary" lens is clean

- **Severity:** Info
- **Category:** SEC
- **Location:** `kat_test.go` (golden KAT + `sha256` self-check), `attrsort.go` (UTF-16
  order), `nsstack.go` (declaration suppression), `walk.go` (epilog quirk), `xmlbase.go`
  (C14N 1.1 `xml:base`), `escape.go` (escape tables)
- **Evidence:**
  ```go
  // kat_test.go — a corrupt regeneration must not hide behind a matching Go bug
  if sum := hex.EncodeToString(sha256Sum(want)); sum != k.sha256 { t.Fatalf(...) }
  // assertJavaError — a Go SUCCESS on a Java failure is the interesting direction and
  // must never pass:
  if err == nil { t.Fatalf("Canonicalize succeeded where Java raised %s", javaClass) }
  ```
- **Impact (observation):** The lenses most likely to find an exploitable defect —
  semantically-equal-but-textually-different XML, conflicting/undeclared namespaces,
  namespace-prefix confusion, exclusive-c14n inference scope, `xml:base` confusion,
  attribute canonical ordering, entity/CDATA/comment/PI handling, and escaping — are all
  either **KAT-pinned against Java goldens** or **structurally impossible** in this
  package (no entity expansion: the `xmldom` parser rejects DOCTYPE; the canonicalizer
  emits only the DOM's already-expanded text). `assertJavaError` even fails a Go
  *success* where Java *failed*, which is the correct direction for a false-accept check.
  The bug-for-bug reproductions (epilog drop, relative-namespace scoping, UTF-16 order,
  `xml:base` decode-then-requote, the two `ErrXMLBaseUnjoinable` aborts) are **deliberate
  parity** and are each individually KAT-locked, so they are not findings. No SEC action.

---

### Tool log (run from `dss/`)

| Tool | Command | Result |
|---|---|---|
| gofmt | `gofmt -l internal/xmlc14n` | clean (no files printed, exit 0) |
| go vet | `go vet ./internal/xmlc14n/...` | clean (exit 0) |
| golangci-lint | `golangci-lint run --config=../.github/.golangci.yml ./internal/xmlc14n/...` | `0 issues.` (exit 0) |
| go test | `go test ./internal/xmlc14n/... -count=1` | `ok … 0.996s` (all tests pass) |
| (none unavailable) | — | all four tools present and run; `fuzz` not run (KAT + `FuzzCanonicalize` round-trip already committed and passing) |

No tool was unavailable. The KAT corpus is reached via `corpustest.Path` (repo-root `corpus/`);
it is present in this checkout, so no `corpustest` skip fired.

---

### Open questions / observations

- **X08-PERF-001 is the one to close before this is treated as done.** It is the only
  finding that turns a few hundred KB of *legal* input into a multi-second stall on the
  verification path, and it is a genuine regression against Santuario's `TreeSet`. The
  build-then-sort fix is byte-preserving and small.
- **Confirm the caller contract for `Canonicalize` (X08-STD-001).** Grep the `xmldsig`
  consumer (U09) to confirm it always routes through the `*ToBytes` helpers and never passes
  a long-lived streaming writer to `Canonicalize` while ignoring the returned error. If it
  does, promote X08-STD-001 to Medium.
- **`compareUTF16` allocation on the non-ASCII path** (`attrsort.go:112-113`, two
  `utf16.Encode([]rune(...))` per comparison that escapes the ASCII fast path) is a minor
  per-comparison allocation inside the O(A log A) sort. It is only reached for astral/
  non-ASCII attribute keys (rare in practice) and is dominated by X08-PERF-001, so it is
  noted rather than raised as a separate finding.
- **Cross-check (U09, `xmldsig` consumer):** this package hands back canonical *bytes*; it does
  no dedup, ID resolution, or duplicate-`Id` rejection itself. Verify in U09 that the consumer
  is the one enforcing duplicate-`Id`/reference semantics on the *input* DOM (not expecting the
  canonicalizer to), consistent with the U06 note that `DuplicateIDs()` is consulted upstream
  of canonicalization.
- **`errURISyntax` vs `ErrXMLBaseUnjoinable` split is correct and load-bearing.** The former
  (caught, previous value kept) and the latter (unchecked, aborts the run) map 1:1 to
  Santuario's `catch (URISyntaxException)` vs the two propagated unchecked exceptions, and both
  directions are pinned (`TestRemoveDotSegmentsEmptyPathIsUnjoinable`,
  `TestUnjoinableXMLBasePropagatesOutOfCanonicalize`). No action — flagged only so a future
  "simplification" of `xmlbase.go` doesn't collapse the two error classes.

---

**No findings (explicitly clean per lens):**
- **SEC (c14n correctness as the verification boundary):** **no findings.** Attribute
  ordering (Java UTF-16 code-unit order, astral vs U+FFFD), namespace declaration
  suppression, exclusive-c14n inference scope, `xml:base` join, escaping, and comment/PI/CDATA
  handling are all KAT-pinned against Java goldens or structurally safe; the
  bug-for-bug reproductions are deliberate parity, not defects. No false-accept / exploitable
  verification finding.
- **PERF-LEAK:** **no findings.** No goroutines, channels, timers, or open handles; all mutable
  state (`engine`, `nsStack`, `xmlAttrStack`) is created inside `Canonicalize` and dies with
  it, so there is no grow-only state and the functions are safe for concurrent use
  (`TestCanonicalizeIsStatelessPerCall`).
- **STD (error handling, defer, context, doc comments, test quality):** the package is
  exemplary — every exported symbol is documented, errors are typed and matched with
  `errors.Is`/`errors.As`, `CanonicalizeToBytes` returns `nil` on error (no partial output),
  the concurrency test pins statelessness, and the test suite is table-driven and
  oracle-pinned. The single Low STD finding (X08-STD-001) is a latent streaming-writer edge on
  an API no in-tree caller uses in that mode.

---

## internal/xmldsig (unit U09)

- **Unit:** U09 — `internal/xmldsig`
- **Scope:** the full `dss/internal/xmldsig` XML-DSig reference-processing engine — `Data`
  union, `Reference`/`Manifest`/`SignedInfo`, `XMLSignature`, the transform pipeline
  (enveloped-signature, base64, C14N, ds:XPath, XPath Filter 2.0, XSLT refusal), URI
  resolvers (fragment / XPointer / detached), and the `sigalg.go` signature-verify engine
  (16 source + 3 test files, 3858 lines).
- **Date:** 2026-08-24
- **Depth:** deep (line-level; 17/19 files read **in full**, 2 budget-truncated — see Files read).
- **Context read (allowed only):** `dss/PORTING.md` conventions, `docs/compatibility/known-gaps.md`,
  and the sibling review sections U06–U08 in this file (their findings are **not** re-reported).
  Deliberate conventions treated as **not** findings: 1:1 Santuario/DSS port shape, **bug-for-bug
  fidelity** (the base64 MIME-decoder quirks, the `length()==1` InclusiveNamespaces guard,
  the physical-c14n *not* being a transform, the XSLT refusal, the ED25519-only EdDSA branch,
  the "last registration wins" ID rule), `secureValidation` OFF by default, `SA1019`/staticcheck-style
  categories disabled, stdlib-only, `internal/` frozen.

**Headline (SEC lens):** this is the **XML-DSig verification boundary** — every `ds:Reference`
digest and the `ds:SignatureValue` over `ds:SignedInfo` are computed here, so any divergence from
Santuario is a signature/verification mismatch or a false-accept. On correctness the package is
**clean**: **no SEC findings.** The attack surfaces the review was primed for are all closed or
deliberately refused — the **enveloped-signature** transform excludes the `ds:Signature` element
correctly (filter prunes the whole subtree; `searchSignatureElement` errors if the transform is
detached, so it can never silently include it), **base64** decodes with a faithful Java MIME-decoder
port that *skips* in-band whitespace/junk but **rejects data after padding and a malformed final
quantum** (pinned by `TestDecodeBase64MatchesTheJavaMimeDecoder`), **XPath**/C14N transforms
dispatch on the correct discriminator state, and `sigalg.go` keeps Santuario's
`(false "signature is wrong")` vs `(error "algorithm/key unanswerable")` split — a wrong-length RSA
value and a key-type mismatch are reported as **errors, not as "invalid signature"**, which is the
correct direction (a validator that folds them into `false` has silently decided the document was
signed and merely tampered with). `Reference.Verify` compares with `subtle.ConstantTimeCompare`, and
there is **no `panic()`** anywhere in the package (hostile input becomes returned errors). The one
genuinely subtle risk — the **XML-signature wrapping attack** — is *not* defended by this package's
default path (`protectAgainstWrappingAttack` is gated behind `secureValidation`, which DSS leaves
off, exactly as Santuario/DSS do upstream); it is instead defended at the DSS layer by
`DSSXMLUtilsIsDuplicateIdsDetected` / `DSSXMLUtilsIsReferencedContentAmbiguous` over
`xmldom.DuplicateIDs()`, with the `xmldom` ID index resolving **last-wins** (Xerces behavior,
deliberate). That cross-layer defense is what makes a duplicate-`Id` wrapping document fail rather
than verify — and it is the thing a future refactor must not accidentally remove (see Open questions).
The two findings below are **performance and idiom**, not correctness.

---

### Files read (17/19 full, 2 budget-truncated)

The two files over the 400-line budget (`kat_test.go` 631, `xmldsig_test.go` 489) were read by the
budget rule: first 300 lines in full, plus `grep func ` signatures for the remainder, plus a
security-term grep (`verify|sign|transform|envelope|exclude|digest|canonical|base64|alg|key|cert|
attack|wrap|duplicate|false`) over the full file to confirm the unread tail is the KAT
manifest/fixture harness and table-driven resolver/transform tests, not new logic.

| File | Lines | Read |
|---|---|---|
| `kat_test.go` | 631 | **truncated** — 1–300 full + `func` signatures (301–631) + security-term grep |
| `xmldsig_test.go` | 489 | **truncated** — 1–300 full + `func` signatures (301–489) + security-term grep |
| `data.go` | 305 | full |
| `xmlutils.go` | 264 | full |
| `resolver_detached.go` | 236 | full |
| `resolver.go` | 233 | full |
| `manifest.go` | 231 | full |
| `transform_xpath2.go` | 213 | full |
| `reference.go` | 212 | full |
| `sigalg.go` | 210 | full |
| `transform.go` | 171 | full |
| `signature.go` | 136 | full |
| `signedinfo.go` | 109 | full |
| `transform_xpath.go` | 81 | full |
| `doc.go` | 81 | full |
| `transform_enveloped.go` | 74 | full |
| `internal_test.go` | 74 | full (base64 + stringFromNode KAT) |
| `transform_base64.go` | 65 | full |
| `transform_c14n.go` | 43 | full |

---

### Findings

Counts: **0 Critical, 0 High, 0 Medium, 2 Low, 0 Info** (2 findings + 3 observations).
**No SEC findings** — the verification-boundary surfaces are clean and oracle-pinned.

### X09-PERF-001

- **Severity:** Low
- **Category:** PERF-MEM
- **Location:** `internal/xmldsig/reference.go:156` (`ReferencedBytes`), `reference.go:195`
  (`Verify`), `reference.go:133` (`ContentsAfterTransformation`); cache field `reference.go:34`
  (`transformsOutput`)
- **Evidence:**
  ```go
  func (r *Reference) ReferencedBytes() ([]byte, error) {
      out, err := r.ContentsAfterTransformation()   // re-resolves + re-runs the whole transform chain
      if err != nil { return nil, err }
      return out.Bytes()
  }
  func (r *Reference) Verify() (bool, error) {
      ...
      got, err := r.CalculateDigest()               // also re-resolves + re-runs the chain + re-canonicalizes
      ...
  }
  // ContentsAfterTransformation sets r.transformsOutput = out, but nothing reads it back:
  func (r *Reference) TransformsOutput() *Data { return r.transformsOutput }
  ```
- **Impact:** `ContentsBeforeTransformation` → `PerformTransforms` → `Data.Bytes()` is the expensive
  path (URI dereference, then the transform chain, then canonicalization). `ReferencedBytes` and
  `CalculateDigest` each re-run it from scratch even though the immediately preceding call already
  produced and cached the result in `r.transformsOutput`. In the XAdES validation flow a reference
  that is (a) verified — `Manifest.VerifyReferences` → `Reference.Verify` → `CalculateDigest` — and
  (b) reported — `xades.DSSXMLUtilsGetReferenceOriginalContentBytes` → `Reference.ReferencedBytes` —
  therefore **dereferences, transforms, and canonicalizes the same content twice**, and the
  `followManifests` branch (`manifest.go:203`) runs it a third time. This mirrors Santuario
  (which likewise does not memoize), so it is not a parity defect, and the signed content is
  normally small; the cost is a redundant canonicalization pass over the referenced node set per
  reference-per-report, i.e. wasted CPU and a transient second copy of the transform output in
  memory. Low rather than Medium because the per-reference signed content in XAdES is typically
  bounded (the signed element or a `ds:Object`), not multi-megabyte.
- **Recommendation:** Make `ReferencedBytes` and `CalculateDigest` reuse `r.transformsOutput` when
  it is non-nil and was produced for the same `ds:Transforms` element (clear it if the reference
  element or chain changes). If that state-tracking is judged fragile against the
  "reads the element fresh on every call" contract documented at `reference.go:30`, at minimum
  document on `ReferencedBytes`/`Verify` that a prior call's output is deliberately discarded and
  re-derived, so a caller can batch report+verify over one `ContentsAfterTransformation` result.

### X09-STD-001

- **Severity:** Low
- **Category:** STD
- **Location:** `internal/xmldsig/kat_test.go:185`
- **Evidence:**
  ```go
  var checked int
  for _, fixture := range order {
      fixture := fixture          // pre-Go-1.22 loopvar shim
      t.Run(fixture, func(t *testing.T) {
          st := loadFixture(t, fixture, detached[fixture])
  ```
- **Impact:** The `fixture := fixture` self-assignment is a workaround for the pre-Go-1.22 shared
  loop-variable semantics. The module is **Go 1.27** (`dss/go.mod`), where each `for` iteration
  already gets a fresh variable, so the shim is a no-op. It is harmless and idiomatic-enough, but it
  is dead weight that a reader may mistake for load-bearing, and it will be flagged (or look wrong)
  to anyone reasoning about loop capture. Pure idiom — no behavior change either way.
- **Recommendation:** Delete the `fixture := fixture` line; the closure already captures the
  per-iteration variable under Go 1.22+ semantics. (No test change needed; the loop body is
  otherwise unchanged.)

---

### Observations (not raised as findings)

- **ED448 is deliberately refused.** `sigalg.go` accepts only `SignatureAlgorithmED25519` in the
  EdDSA branch and returns `ErrUnsupportedSignatureAlgorithm` for `ED448` — the correct, *safe*
  direction (an "unanswerable" error, never a `false` that could be read as "signature invalid").
  This follows the stdlib-first / `golang.org/x`-only dependency policy (no Ed448 in
  `crypto/ed25519`; `golang.org/x/crypto` is not used here). It is consistent with, though not
  explicitly named in, `known-gaps.md`'s digest-algorithm coverage note. No action.
- **`rooted()` is O(|set|·depth) on the non-document-order path** (`transform_xpath2.go:194`,
  `IsNodeInclude`). The XPath Filter 2.0 filter prefers the O(1) `IsNodeIncludeDO` counter path
  during a single document-order canonicalization and only falls back to `rooted()` for the
  out-of-order cases (attributes, end tags) — the same shape as Santuario's `XPath2NodeFilter`.
  Bounded by the selected set size and tree depth; not a finding.
- **`Data.Bytes()` caches and `IsElement` deliberately dropped `hasOctets`** (`data.go:96-118`,
  `data.go:186-204`). The FIX note at `data.go:110` records that requiring `!hasOctets` previously
  broke any identity transform followed by another (the exact chain the enveloped transform's
  javadoc prescribes), and the fix was verified against the Java oracle (55/57 → 56/57 KAT). This
  is a *corrected* parity bug, not a new one — the one KAT case Java itself mis-handles remains
  documented. No action.

---

### Tool log (run from `dss/`)

| Tool | Command | Result |
|---|---|---|
| gofmt | `gofmt -l internal/xmldsig` | clean (no files printed, exit 0) |
| go vet | `go vet ./internal/xmldsig/...` | clean (exit 0) |
| golangci-lint | `golangci-lint run --config=../.github/.golangci.yml ./internal/xmldsig/...` | `0 issues.` (exit 0) |
| go test | `go test ./internal/xmldsig/... -count=1` | `ok … 1.008s` (all tests pass) |
| (none unavailable) | — | all four tools present and run; KAT corpus reached via `corpustest.Path` (repo-root `corpus/`), present in this checkout, so no `corpustest` skip fired |

No tool was unavailable.

---

### No findings (explicitly clean per lens)

- **SEC (the verification boundary):** **no findings.** Enveloped-signature exclusion
  (subtree-pruning filter + `searchSignatureElement` detachment guard), base64 decode strictness
  (skips junk, rejects data-after-padding and bad final quantum), XPath/XPath2/C14N transform
  discriminator dispatch, reference URI resolution (empty URI = whole document; `#id` = last-wins
  element; XPointer forms; detached digest-over-name preference), the transform chain application
  order, `ds:SignedInfo` c14n caching (no re-encoding drift), and the `sigalg.go` RSA/PSS/ECDSA/DSA/
  Ed25519 verify with correct key-type and malformed-value classification are all correct and
  KAT-pinned against Santuario/Java goldens (`kat_test.go`, `internal_test.go`). Digest comparison
  is constant-time. There is **no `false`-accept and no `panic` on hostile input**. The bug-for-bug
  reproductions (MIME base64 quirks, `length()==1` InclusiveNamespaces, physical-c14n-not-a-transform,
  XSLT refusal, ED25519-only EdDSA) are **deliberate parity/safety** and individually pinned.
- **PERF-LEAK / PERF-BIGO:** **no findings.** No goroutines, channels, `unsafe`, or open handles
  (grep-confirmed); all mutable state is per-call and dies with the call, so the package is safe
  for concurrent use and has no grow-only state. The one complexity observation is X09-PERF-001
  (a redundant canonicalization pass, Low).
- **STD:** the package is exemplary — every exported symbol is documented with its Santuario
  counterpart, errors are typed and matched with `errors.Is`/`errors.As`, the `(bool, error)`
  split is load-bearing and documented, and the test suite is table-driven and oracle-pinned.
  The single Low STD finding (X09-STD-001) is a redundant Go-1.22+ loopvar shim in the KAT
  harness — no behavioral impact.

---

### Open questions / observations

- **The wrapping-attack defense lives OUTSIDE this package — confirm it is always on.**
  `protectAgainstWrappingAttack` (`resolver.go:108`, `:203`) is gated behind `ctx.SecureValidation`,
  which DSS sets to **false** (`xades/xades_signature.go:1089` passes `nil` `ManifestOptions`). So
  a `#id` reference with a duplicate `Id` resolves (last-wins) *here* without complaint, exactly as
  upstream Santuario/DSS do. What actually stops the XML-signature wrapping attack is the DSS-layer
  check — `DSSXMLUtilsIsDuplicateIdsDetected` (`xades/dss_xml_utils.go:486`) and
  `DSSXMLUtilsIsReferencedContentAmbiguous` (`xades/dss_xml_utils.go:1077`), both over
  `xmldom.DuplicateIDs()`. **Confirm every verify path (XAdES, and any PAdES/ASiC path that
  resolves `#id`) actually consults these before trusting a reference.** If a future refactor routes
  a `#id` reference through `xmldsig` without the DSS-layer duplicate check, the wrapping attack
  becomes exploitable (attacker's second element wins resolution and its digest can be made to
  match). This is the single most important cross-layer invariant in the XML chain.
- **X09-PERF-001 is the only actionable item and it is Low.** Closing it (reuse
  `r.transformsOutput`) is a small, byte-preserving change; do it only if profiling shows the
  double canonicalization matters, otherwise it is optional hardening.
- **X09-STD-001 is a one-line cleanup** (drop the `fixture := fixture` shim) — safe under
  Go 1.27, no test change needed.
- **Cross-check with U08:** U08's open question about the `Canonicalize` streaming-writer contract
  (X08-STD-001) is **answered** by this review — every `xmldsig` caller routes through the
  in-memory `xmlc14n.Canonicalize(alg, in, *bytes.Buffer)` form (`data.go:201`, `signedinfo.go`,
  `transform_c14n.go`); no `xmldsig` path passes a long-lived streaming writer, so X08-STD-001
  stays Low and does not need promotion.

---

## xml/common (unit U10)

- **Unit:** U10 — `dss/xml/common` (47 Go files, 3,339 lines)
- **Scope:** the XML definer framework shared by every report "jaxb" package — the
  `DSSNamespace`/`DSSElement`/`DSSAttribute` definition model, the `XMLDSig*` concrete
  vocabulary, the `XPathQuery` chain items/parameters that build XPath *strings*, the
  `AbstractPath`/`XPathQueryBuilder`/`XPathExpressionBuilder` factories, and the JAXP
  security-plumbing builders (`DocumentBuilderFactoryBuilder` + three documented stubs).
- **Date:** 2026-08-24
- **Depth:** deep (line-level; see Files read)
- **Context read (allowed only):** `dss/PORTING.md`, `docs/compatibility/known-gaps.md`, and
  the U07 section of this file (which deferred the unescaped-literal surface to this unit as
  OQ-3). Deliberate conventions treated as **not** findings: 1:1 Java port shape, `panic`
  standing in for Java unchecked exceptions/`requireNonNull`, staticcheck-style categories
  disabled, stdlib-only, `internal/` frozen, documented stubs (`SchemaFactory`/`Validator`/
  `TransformerFactory`) as explicitly documented no-ops, verbatim upstream `toString()` quirk.
- **Cross-unit note resolved:** this unit closes **OQ-3** left by U07 — the unescaped
  single-quote interpolation in the Id/attribute-lookup XPath literal. See **X10-SEC-001**.

### Files read (47/47 — all full)

No file exceeds 300 lines (largest: `xpath_query_builder.go`, 222), so every file was
read in full; none was sampled or skimmed. 38 source + 9 test.

| File | Lines | Read |
|---|---|---|
| `doc.go` | 49 | full |
| `abstract_configurator.go` | 132 | full |
| `abstract_configurator_test.go` | 80 | full |
| `abstract_factory_builder.go` | 31 | full |
| `abstract_path.go` | 46 | full |
| `abstract_xpath_query.go` | 93 | full |
| `abstract_xpath_query_item.go` | 59 | full |
| `abstract_xpath_query_parameter.go` | 24 | full |
| `all_from_current_position_xpath_query.go` | 19 | full |
| `all_xpath_query.go` | 18 | full |
| `definition_test.go` | 71 | full |
| `dss_attribute.go` | 25 | full |
| `dss_element.go` | 64 | full |
| `dss_error_handler.go` | 81 | full |
| `dss_error_handler_alert.go` | 90 | full |
| `dss_error_handler_alert_test.go` | 96 | full |
| `dss_namespace.go` | 37 | full |
| `document_builder_factory_builder.go` | 131 | full |
| `document_builder_factory_builder_test.go` | 39 | full |
| `from_current_position_xpath_query.go` | 18 | full |
| `schema_factory_builder.go` | 118 | full |
| `security_configuration_exception.go` | 30 | full |
| `transformer_factory_builder.go` | 101 | full |
| `validator_configurator.go` | 128 | full |
| `validator_configurator_test.go` | 70 | full |
| `xsd_validation_exception.go` | 39 | full |
| `xml_definer_utils.go` | 104 | full |
| `xmldsig_attribute.go` | 36 | full |
| `xmldsig_element.go` | 128 | full |
| `xmldsig_namespace.go` | 9 | full |
| `xmldsig_path.go` | 103 | full |
| `xmldsig_path_test.go` | 64 | full |
| `xpath_expression_builder.go` | 145 | full |
| `xpath_expression_builder_test.go` | 70 | full |
| `xpath_matchnode_test.go` | 143 | full |
| `xpath_query.go` | 28 | full |
| `xpath_query_any_item.go` | 44 | full |
| `xpath_query_attribute_item.go` | 61 | full |
| `xpath_query_attribute_parameter.go` | 115 | full |
| `xpath_query_builder.go` | 222 | full |
| `xpath_query_builder_test.go` | 112 | full |
| `xpath_query_element_item.go` | 61 | full |
| `xpath_query_end_item.go` | 65 | full |
| `xpath_query_identifier_parameter.go` | 23 | full |
| `xpath_query_item.go` | 37 | full |
| `xpath_query_not_child_of_parameter.go` | 73 | full |
| `xpath_query_parameter.go` | 7 | full |

Cross-package reads (allowed, for SEC context, not findings): `xml/utils/xpath_utils.go`
(`XPathUtilsGetElementById[WithQuery]` — the Id-lookup entry point), `xml/utils/dom_utils.go`
(`DomUtilsGetId`/`DomUtilsGetXPointerId` — where the attacker-derived value originates),
`xml/utils/java_xml_xpath_query_executor.go` (proves the built string is compiled+evaluated
by a real engine), `internal/xpath10/{doc,parser,eval}.go` (confirms the evaluable subset:
`or`, `=`, `local-name()`, string literals), `xades/reference_verifier.go` + `reference_processor.go`
(verification-path consumers of the Id lookup), `internal/xmldsig/resolver.go` (the
*separate* `EnforcedResolverFragment` char-filter guard — confirms it does **not** cover this
path).

### Tool log

Run from `dss/`:

- `gofmt -l xml/common` → **no output** (clean).
- `go vet ./xml/common/...` → **OK** (no diagnostics).
- `golangci-lint run --config=../.github/.golangci.yml ./xml/common/...` → **`0 issues.`**

No tool was unavailable.

### Findings

Counts: **0 Critical, 0 High, 1 Medium, 0 Low, 0 Info** (1 total).
The package is clean on the STD and PERF lenses (explicit "no findings" below); the single
finding is the confirmed special-attention XPath-interpolation surface, reported as a
hardening gap and noted explicitly as **Java parity** (closing U07's OQ-3).

#### X10-SEC-001

- **Severity:** Medium
- **Category:** SEC
- **Location:** `xml/common/xpath_query_attribute_parameter.go:94-104`
  (`xPathQueryAppendAttributeCondition`), reached through
  `xml/common/xpath_query_identifier_parameter.go:14` (`NewXPathQueryIdentifierParameter`) and
  `xml/common/xpath_query_builder.go` (`IdValue`), consumed by
  `xml/utils/xpath_utils.go:146-164` (`XPathUtilsGetElementById[WithQuery]`).
- **Evidence:**
  ```go
  sb.WriteString("@*[local-name()='")   // xpath_query_attribute_parameter.go:98
  sb.WriteString(attrName)
  sb.WriteString("']")
  sb.WriteString("='")
  sb.WriteString(attributeValue)        // raw value — single quotes NOT escaped
  ```
  The built string is not display-only: `JavaXmlXPathQueryExecutor` compiles+evaluates it with
  `internal/xpath10`, whose subset supports exactly the breakout constructs — `or`
  (`parser.go:77` `parseOr`), `=` (`eval.go:357` `equals`), `local-name()` and string literals.
- **Impact:** A document-derived value (an `Id`/attribute value taken off a `Reference/@URI`
  under validation, via `DomUtilsGetId`) is interpolated into a single-quoted XPath literal with
  no escaping. Because the result is evaluated by a real engine that accepts `or`-disjunctions,
  a value containing a `'` can break out of the literal and alter the predicate (e.g. append a
  disjunct that matches a chosen node), which feeds the `XPathUtilsGetElementById(...) == nil`
  validity checks in the verify path (`xades/reference_verifier.go:67`). **Java parity:** upstream
  `XPathQueryAttributeParameter`/`XPathQueryIdentifierParameter` build the identical
  unescaped `...='value'` literal, so this is a faithful 1:1 reproduction, not a port
  divergence — it is a *pre-existing upstream hardening gap*, not a new defect introduced by the
  Go port. Bounded in practice: number literals are refused by `xpath10` (so a clean
  always-true `1=1` is not reachable), a meaningful exploit needs a known target id, and the
  primary anti-wrapping defense is the separate DSS-layer duplicate-`Id` check (see U09's open
  question). Notably, this DSS `DomUtils`/`XPathUtils` Id-lookup path is **not** covered by the
  `EnforcedResolverFragment` `xpathCharFilter` guard, which only protects the separate
  `internal/xmldsig` reference-resolver path.
- **Recommendation:** Do not "fix" by escaping silently — that would be a divergence from Java
  and breaks the parity contract unless done with the `// DIVERGENCE, deliberate:` ceremony +
  a PORTING.md entry. The correct action is to **record it as a known hardening gap** in
  `docs/compatibility/known-gaps.md` (it is absent today), naming the exact
  `QueryString()` surface, the Java-parity status, and the fact that the DSS Id-lookup path
  lacks the char-filter the `xmldsig` resolver path has. If a security team later wants the
  literal escaped (`'` → `''`), apply the divergence ceremony and add a KAT pinning the new
  string so the byte-parity gate is updated deliberately.

### Lens summaries

- **STD — no findings.** Builder-based string construction throughout (`strings.Builder`, no
  concat-in-loop); consistent `model.DSSError`/`errors.As` handling; deterministic iteration
  (`featureNames`/`attributeNames` slices guard map order); `sync.Once` singleton (strictly
  safer than Java's non-thread-safe lazy init); panics are only `requireNonNull`/
  unsupported-builder-state guards (documented, unreachable via hostile XML — the `process`/
  `MatchNode` methods return `bool`, not panic, on node-kind mismatch). Tests are strong:
  KAT-pinned `QueryString()` against a Java oracle, `MatchNode` against real `*xmldom.Node`
  trees, panic paths, and alert/exception semantics all covered. No `context.Context`/`defer`/
  goroutine/leak surface in a pure definition+query-building package.
- **PERF — no findings.** Element traversal is linear over a node's attributes
  (`xpath_query_attribute_parameter.go:58` `for _, attributeNode := range node.Attrs`);
  `QueryString()` builds O(path-depth + params) with a builder; `XPathQueryBuilderFromXPathQuery`
  uses type-switch (not reflection); no O(n²) over element count, no hot-loop allocations, no
  unbounded recursion.
- **SEC — one finding (X10-SEC-001, Medium, Java-parity hardening gap).** No XML-injection
  surface (this package builds query strings and matches nodes; it does not serialize XML —
  that is U06/U07 `xmldom`/`xml/utils` territory). Namespace matching in
  `XPathQueryElementItem.process` (`URI() == "" || URI() == node.Name.Space`) is a documented
  Java-parity judgment call (`dss_element.go` `URI()` doc) and is inert for every real
  `DSSElement` here (all carry `XMLDSigNS` with a real URI). XXE is a non-issue: the
  `DocumentBuilderFactoryBuilder` maps onto `xmldom.ParseOptions`, which has no DTD/entity layer
  at all (documented no-op table). No `unsafe`.

### Open questions / observations

- **X10-SEC-001 is the one actionable item and it is a *documentation* action, not a code
  action.** The defensible move under the project's parity contract is to add it to
  `docs/compatibility/known-gaps.md` (explicitly noting Java parity and that the DSS Id-lookup
  path is outside the `EnforcedResolverFragment` char-filter), not to change `QueryString()`.
  Any escaping change must go through the divergence ceremony.
- **Cross-check with U07/U09:** this closes U07's **OQ-3** (the surface U07 deferred here as
  "parity with Java … not reported as a new finding") — confirmed real, confirmed Java parity,
  and confirmed the Go `xpath10` engine is permissive enough to evaluate a `or`-based breakout
  (so it is not neutralized by the engine's subset restriction). It is *complementary* to, not a
  duplicate of, U09's wrapping-attack invariant: the duplicate-`Id` DSS-layer check is the
  primary defense; this is the secondary, unhardened Id-lookup path.
- **`document_builder_factory_builder.go` is quietly the *strongest* security posture in the
  chain** — because `internal/xmldom` has no DTD/entity layer, the JAXP hardening features are
  no-ops on a parser that simply cannot be XXE'd. Worth a line in the compatibility docs so
  readers don't assume a DTD path exists.
- **The three stub builders (`SchemaFactory`/`TransformerFactory`/`Validator`) are
  correctly-inert** — every `setSecurityFeature`/`setSecurityAttribute` is a documented no-op
  and `Build()` always succeeds, matching the "XSD/XSLT not ported" stance in known-gaps. No
  dead-but-dangerous code; the feature/attribute machinery is preserved for a future phase.
