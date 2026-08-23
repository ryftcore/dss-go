# PDF engine binding design: `internal/pdf`

*Design record for the native PDF engine. Companion document:
`dss/internal/xmldom/DESIGN.md` (the same exercise for the XML stack).*

| Upstream artefact this package replaces | Where it lives upstream |
| --- | --- |
| `org.apache.pdfbox:pdfbox:3.0.7` (parser, `COSDocument`/`COSWriter`, security handlers, filters) | maven central |
| `com.github.librepdf:openpdf:1.3.43` | second SPI implementation upstream — **not** ported (§0.3) |
| `dss-pades-pdfbox` (`PdfBoxDocumentReader`, `PdfBoxSignatureService`, `PdfBoxDict`, `PdfBoxArray`, `PdfBoxObjectKey`, `PdfBoxUtils`) | the only consumer whose calls define our scope |
| `dss-pades` `eu.europa.esig.dss.pdf.*` SPI (`PdfDocumentReader`, `PdfDict`, `PdfArray`, `PdfObject`) | the interface `internal/pdf` must be able to satisfy |

Per `PORTING.md`: pdfbox has no DSS Java class to mirror, so — exactly like `internal/asn1ber`,
`internal/cmscore` and `internal/xmldom` — the machinery lives under `internal/` and states its
provenance in `doc.go`. It never imports a DSS package. Beyond the standard library it imports
exactly one module, `golang.org/x/text`, and from one file: `saslprep.go` needs NFKC normalisation
(`x/text/unicode/norm`) and Unicode bidirectional classes (`x/text/unicode/bidi`) to reproduce what
pdfbox's `SaslPrep` gets from `java.text.Normalizer` and `Character.getDirectionality` (§2.6).
`golang.org/x/text` is already a direct dependency of the module and is on `CONTRIBUTING.md`'s
allowed list; no other file here reaches outside the standard library.

---

## 0. Scope, layering and non-goals

### 0.1 What DSS actually asks a PDF library to do

Read every call site in `dss-pades-pdfbox` and the answer is small and sharp. `internal/pdf` is
**not** a general PDF library. The complete list of capabilities the port needs:

**Reading** (`PdfBoxDocumentReader` + `AbstractPDFSignatureService.getRevisions`):

1. Load a document, optionally with a password (`Loader.loadPDF(bytes, password)`), and report
   `isEncrypted()`.
2. `getSignatureFields()` — walk `/Root /AcroForm /Fields` (recursively through `/Kids`), keep only
   `/FT /Sig` fields, and for each report its fully-qualified partial name, its widget annotation,
   and the **object number** of its `/V` value (`sigDictObject.getKey().getNumber()` — used to
   deduplicate fields that share one signature dictionary).
3. Read the signature dictionary entries `/Type /Filter /SubFilter /Name /M /Reason /Location
   /ContactInfo /Contents /ByteRange /Reference` (and through `/Reference`, `/TransformMethod`,
   `/TransformParams /P /V`, `/Data`) — `PdfSigDictWrapperFactory`.
4. Read the catalog: `/Version /Extensions /Perms (/DocMDP /UR /UR3) /AcroForm /DSS`.
5. Read the `/DSS` dictionary and its `/Certs /CRLs /OCSPs` **stream** arrays and `/VRI` sub-dictionaries
   (`/Cert /CRL /OCSP /TU /TS`), including the object key of each array element
   (`DSSDictionaryExtractionUtils`, `SingleDssDict`, `PdfVriDict`).
6. Page tree: page count, `/MediaBox`, `/Rotate`, `/Annots` with each annotation's `/Rect`, `/T`,
   `/V`, `/F` flags (`getPdfAnnotations`).
7. Encryption permission bits: `canModify`, `canModifyAnnotations`, `canFillInForm`,
   `isOwnerPermission` (`PdfPermissionsChecker`).
8. Header version and catalog `/Version`.
9. Generic object-graph traversal with identity by object key, and **raw** (still-encoded) stream
   bytes plus raw stream size, for `DefaultPdfObjectModificationsFinder.compareDictStreams` — it
   compares *encoded* bytes, never decoded ones. Decoded bytes are needed only for `/DSS` token
   streams and for `PdfObjectModificationsFilter.isStreamFill` (a non-empty test).
10. Byte offsets of every incremental revision — but note **DSS does not walk `/Prev` to do this**:
    `PAdESUtils.extractRevisions` scans raw bytes for `%%EOF` + its EOL. We must reproduce that
    scanner exactly, because the resulting lengths feed
    `PAdESUtils.getPreviousRevision(byteRange, revisions)`.

**Writing** (`PdfBoxSignatureService`):

11. Incremental update only — `pdDocument.saveIncremental(os)`. Prior bytes are copied verbatim.
12. Create a signature dictionary with a placeholder `/Contents` and a placeholder `/ByteRange`,
    then patch both in place after the increment is laid out.
13. Create a `/Sig` field + widget annotation into `/AcroForm /Fields` and the page's `/Annots`, or
    fill an existing empty field named by `/T`.
14. Write `/DSS` (+ optional `/VRI`) into the catalog, reusing already-present token streams by
    object key where possible.
15. Write a `DocTimeStamp` signature dictionary (same machinery, `/Type /DocTimeStamp`,
    `/SubFilter /ETSI.RFC3161`, no `/M`, no `/Name`).
16. Add developer-extension dictionaries into `/Root /Extensions` (pure dict manipulation).
17. Re-encrypt the increment when the source document is encrypted, with a caller-supplied
    deterministic CSPRNG (`SecureRandomProvider` exists upstream precisely for reproducible output).

That is the whole contract. Everything else pdfbox does is out of scope.

### 0.2 Non-goals — stated so nobody "helpfully" adds them

* **No rasterization.** `generateImageScreenshot`, `generateImageScreenshotWithoutAnnotations`,
  `PdfBoxScreenshotBuilder`, `PDFRenderer`, and the whole `pdf/visible` drawer stack are not ported.
  Consequence, to be documented in the `pades` layer: `DefaultPdfDifferencesFinder.getVisualDifferences`
  yields no `VISUAL_DIFFERENCE` modifications. `getAnnotationOverlaps` and `getPagesDifferences` are
  purely geometric (`/Rect`, `/MediaBox`, page count) and **are** fully supported.
* **No content-stream interpretation.** No text extraction, no fonts, no graphics operators.
  Image filters (`DCTDecode`, `CCITTFaxDecode`, `JPXDecode`, `JBIG2Decode`) are never decoded —
  see §2.5.
* **No full-rewrite serializer.** `internal/pdf` can only *append*. There is no code path that
  re-emits an object that already exists in the input at its original offset.
* **No encryption of a document that is not encrypted.** `StandardSecurityHandler
  .prepareDocumentForEncryption` / `StandardProtectionPolicy` are not ported: the writer only
  re-encrypts under a handler the input already carries (§2.6, §3.2 R20), which is also all DSS
  itself does (`PAdESSignatureParameters.setPasswordProtection` opens a protected document; nothing
  in DSS protects one). Encrypting would need the full rewrite above. Consequently only
  `SaslPrep.saslPrepQuery` is ported, not `saslPrepStored`.
* **No PDF/A, no linearization, no tagged-PDF, no object-stream *writing*** (§3.2 R15).
* **No openpdf port.** Upstream ships two interchangeable SPI backends selected by
  `ServiceLoaderPdfObjFactory`; the Go port has exactly one native backend, so `IPdfObjFactory`
  collapses to a constructor. Nothing in `ITextDocumentReader` requires a capability
  `PdfBoxDocumentReader` does not (the two implement the identical `PdfDocumentReader` interface,
  verified method-for-method).
* **No `PdfMemoryUsageSetting` modes.** Go works from an `io.ReaderAt`; `MEMORY_FULL`,
  `MEMORY_BUFFERED` and `FILE` all map onto "whatever `io.ReaderAt` you hand us". The upstream enum
  is ported in `pades` as a no-op knob for API compatibility.

### 0.3 Layering

```
pades  (ported Java classes: PAdESUtils, PdfSigDictWrapper, SingleDssDict, ByteRange, …)
  │  imports
  ▼
internal/pdf      ← this document. stdlib, plus golang.org/x/text in saslprep.go.
```

`internal/pdf` knows nothing about CMS, certificates, OCSP, or ETSI. It hands `[]byte` up. The
`pades` layer keeps upstream's class names and does the ETSI semantics. Concretely:
`internal/pdf` parses `/ByteRange` into `[]int64`; `pades.ByteRange` (the port of
`eu.europa.esig.dss.pades.validation.ByteRange`) owns `validate()` and `getLength()`.

The writer depends on the reader only through the `*Document` surface in §4.3.

---

## 1. Evidence: what the upstream corpus actually contains

Everything in §2 is scoped from measurement, not from the PDF spec's table of contents. The survey
program is `internal/pdf/testdata/gen/PdfOracle.java` (§5).

*Reproduction artefacts from the run that produced this table (not committed):*
`PdfCorpusSurvey.java` (the prototype oracle), `pdfprobe-pom.xml` (the two-dependency pom that
fetches pdfbox 3.0.7 and its sources), `pdf_corpus_survey.txt` (the full 248-line per-file dump plus
the summary quoted below), and the extracted pdfbox 3.0.7 sources cited throughout
§2.7 and §3.2. The numbers below are its output over
**all 248 PDFs** under `dss-pades/src/test/resources` with pdfbox 3.0.7.

```
total=248 loaded=246 failed=2 encrypted=8 withSignatureFields=207
withObjStm=106 hybridXref=19 maxEofRevisions=54 maxObjects=2927
```

| Dimension | Measured distribution | Consequence |
| --- | --- | --- |
| **Stream filters** (documents using) | `FlateDecode` 240, `DCTDecode` 18, `CCITTFaxDecode` 3, `ASCII85Decode` 1 | §2.5: decode **FlateDecode only**; everything else is carried as raw bytes |
| **Predictors** | `/Predictor 12` (PNG Up) 71 objects, `/Predictor 15` (PNG Optimum) 158 objects. No `2` (TIFF), no `10/11/13/14` alone | Implement PNG predictors 10–15 (15 implies per-row tags, so 10–14 come for free); TIFF predictor 2 is implemented for completeness but has no corpus coverage |
| **xref section style** (per revision, `/Prev`-chain walk) | `table` 572, `stream` 202, unresolvable 2 | Both are mandatory. Ratio ≈ 3:1 in favour of tables |
| **Hybrid xref** (`/XRefStm` in a table trailer) | 19 documents | Mandatory |
| **Object streams** (`/Type /ObjStm`) | 106 documents | Mandatory |
| **Trailer keys seen** | `Root` 246, `Size` 246, `ID` 240, `Prev` 207, `Info` 184, `Encrypt` 9, `XRefStm` 19, `DocChecksum` 35, `AdditionalStreams` 34; plus xref-stream dict keys `Type/W/Index/Filter/Length/DecodeParms` 85–88 | `DocChecksum`/`AdditionalStreams` are real keys in these files; the writer **must drop `/DocChecksum`** (§3.2 R13) exactly as `COSWriter.doWriteTrailer` does |
| **Encryption** | `Standard V2 R3 RC4-128` ×1; `Standard V4 R4 /StdCF /AESV2` ×6; `Standard V5 R6 /StdCF /AESV3` ×1 | §2.6: implement Standard handler R2–R6. Test password for `protected/*.pdf` is a single space `" "` |
| **Header versions** | 1.2 ×1, 1.3 ×8, 1.4 ×95, 1.5 ×30, 1.6 ×84, 1.7 ×27, 2.0 ×1 | Header version drives the developer-extension logic in `AbstractPDFSignatureService` |
| **Signature `/SubFilter`** | `ETSI.CAdES.detached` 275, `ETSI.RFC3161` 74, `adbe.pkcs7.detached` 26, `adbe.pkcs7.sha1` 2, `PBAD.PAdES` 1 | The reader must not whitelist SubFilters — `PBAD.PAdES` exists and must parse |
| **Hard-failure exhibits** | `EmptyPage-corrupted.pdf`, `EmptyPage-corrupted2.pdf` (page tree root not a dictionary — pdfbox *throws*); `validation/pdf-signed-corrupted.pdf`, `validation/pades-5-signatures-and-1-document-timestamp.pdf` (broken `startxref`/`/Prev` offsets — pdfbox *recovers*) | §2.7 pins which class each belongs to |

Two corpus facts worth stating loudly:

* **54 `%%EOF` markers** in the deepest document. Revision handling must be O(n), not O(n²).
* `validation/doc-firmado-LT.pdf` — a DSS-produced file — contains literally
  `/ByteRange [0 1383 28009 307]` followed by **18 spaces** before `\n>>`, and its xref free entry
  is `0000000000 65535 f\r\n`. That is direct byte-level confirmation of §3.2 R12/R18.

---

## 2. Reader

### 2.1 Lexer (`lexer.go`)

A byte-level tokenizer over `io.ReaderAt`, never over `bufio` alone — the parser seeks constantly.

Whitespace, per ISO 32000-1 §7.2.2: `\x00 \t \n \f \r \x20`. Delimiters: `( ) < > [ ] { } / %`.
Comments (`%` to EOL) are whitespace everywhere except inside strings and stream data. `%%EOF` and
`%PDF-` are *comments to the lexer*; only the byte scanners in `revision.go` and the header parser
treat them specially.

Token kinds: integer, real, name, literal string, hex string, `[`, `]`, `<<`, `>>`, keyword
(`obj endobj stream endstream R true false null xref trailer startxref f n`).

Pinned lexer tolerances (all measured against pdfbox `BaseParser`):

| Input | Behaviour |
| --- | --- |
| `/Na#6de` | name unescapes `#XX`; a malformed `#` (fewer than two hex digits follow) is kept literally |
| `(a\)b)`, nested `(a(b)c)` | balanced-paren and backslash escapes `\n \r \t \b \f \( \) \\`, `\ddd` octal (1–3 digits, mod 256), backslash-newline line continuation |
| `<A1B>` | odd-length hex string is padded with a trailing `0` |
| `<A1 \n B2>` | whitespace inside hex strings is skipped; non-hex bytes are skipped with a warning |
| `1.` / `.5` / `--3` / `4.5.6` | parsed leniently as reals; the **literal is preserved verbatim** in `Real.Raw` (§3.2 R8) |
| `6 0 R` where `6 0` is followed by anything else | two integers, not a reference — one-token lookahead beyond the second integer is required |
| numbers exceeding `int64` | clamped and a warning recorded; never a parse error |

The lexer is allocation-conscious: names and keywords are compared against a fixed table without
allocating for the common cases.

### 2.2 Object model (`object.go`)

```
Object = Null | Bool | Integer | Real | String | Name | Array | *Dict | *Stream | Ref
```

Two decisions that matter and are not negotiable:

* **`*Dict` preserves insertion order.** pdfbox's `COSDictionary` is a `LinkedHashMap`, and
  `PdfDict.list()` (used by `SingleDssDict.extractVRIs` to enumerate VRI names) exposes that order to
  DSS. `map[Name]Object` alone is a determinism bug. `Dict` therefore carries `keys []Name` plus
  `m map[Name]int`. `Set` on an existing key updates in place and keeps its position — the same as
  `LinkedHashMap.put`.
* **`Real` keeps its source text.** `Real{Val float64, Raw string}`. Every real we parse and echo
  back is written from `Raw`; only reals we *create* go through `FormatReal` (§3.2 R8). This makes
  round-tripping a `/Rect [0.0 0.0 595.276 841.89]` byte-exact without inventing a float formatter
  that has to match Java's `Float.toString`.

`Stream` holds `Dict *Dict` plus the raw (still-filtered, still-encrypted) bytes and the source byte
range. Decoding is a `*Document` method, not a `*Stream` method, because it needs the decryption key
and the filter chain resolution.

`Ref{Num int64, Gen uint16}` is a value type; `ObjectKey` is the same shape and is what the API
exposes to `pades` (the analogue of `PdfObjectKey` / `PdfBoxObjectKey`). They are distinct types on
purpose: `Ref` is a *value inside the object graph*, `ObjectKey` *identifies a slot*.

### 2.3 Cross-reference resolution (`xref.go`)

Startup sequence, mirroring `COSParser.parseXref` but simplified to what we need:

1. Read the header (§2.7 H1–H3).
2. Read the last `min(fileLen, 2048)` bytes; find the last `%%EOF`; find the last `startxref`
   preceding it; read the offset. (pdfbox's `DEFAULT_TRAIL_BYTECOUNT` is 2048 and its
   `readTrailBytes` is settable; we pin 2048 and do not expose a knob.)
3. Walk the chain. Each hop yields an `XRefSection`:

```go
type XRefStyle uint8   // XRefTable | XRefStream
type XRefSection struct {
    Offset    int64      // byte offset of `xref` or of the `N G obj` of the xref stream
    Style     XRefStyle
    Trailer   *Dict      // trailer dict, or the xref-stream dict
    Prev      int64      // -1 when absent
    XRefStm   int64      // hybrid: /XRefStm offset, -1 when absent
    Entries   int        // entries contributed by this section
    Recovered bool       // true when the offset had to be repaired (§2.7)
}
```

4. **Cross-reference table** (`xref` keyword): subsections `first count`, then exactly-20-byte
   entries `%010d %05d [nf]` + 2-byte EOL. Leniency per `COSParser.parseXrefTable`:
   * a subsection header that does not split into exactly two integers aborts the section (warning),
   * entry lines shorter than 3 space-separated fields abort the section (warning) — this is
     pdfbox's PDFBOX-474 tolerance,
   * entries with offset `0` are skipped,
   * only entries ending in `n` are recorded; `f` entries are recorded as free,
   * the section terminates early on `t` (`trailer`) or on a delimiter byte.
5. **Cross-reference stream** (`/Type /XRef`): `/W [w1 w2 w3]` (a zero `w1` means "type 1"),
   `/Index [first count …]` defaulting to `[0 /Size]`, `/Size`. Field widths are read big-endian.
   Type 0 = free, type 1 = `(offset, gen)`, type 2 = `(objstm object number, index within it)`.
   The stream is `FlateDecode` + predictor (evidence: predictors 12 and 15 dominate).
6. **Hybrid** (`/XRefStm` in a *table* trailer): parse the referenced xref stream and merge it
   **beneath** the table — table entries win. pdfbox repairs a bad `/XRefStm` offset by brute force
   and, failing that, skips it (`LOG.error("Skipped XRef stream due to a corrupt offset")`); so do we.
7. Continue with `/Prev` until the offset repeats (cycle guard), is out of range, or is absent.
   Section cap: 512 (the corpus max is 54 revisions).
8. Earlier sections never override later ones: fill the map from newest to oldest and refuse to
   overwrite an existing key. The trailer presented to callers is the newest section's trailer,
   with `/Root` and `/Info` inherited from an older section if the newest lacks them (pdfbox does
   this in `XrefTrailerResolver.getTrailer`).
9. Then `checkXrefOffsets()` (§2.7 X1).

`XRefSections()` is exported: the writer needs the newest section's `Style` to decide table vs
stream (§3.4), and the oracle compares it per revision.

### 2.4 Object streams (`objstm.go`)

A type-2 xref entry names a container. Loading it: decode the `/ObjStm` stream, read `/N` pairs of
`(objnum, relative offset)` from the front, then parse each object at `/First + offset`.

Rules:
* Object streams are loaded lazily and cached, keyed by container object number.
* An object stream may not itself live in an object stream; a type-2 entry pointing at a container
  that is itself type-2 is a hard error (cycle).
* If the object at the declared index has a different object number than the xref claims, trust the
  **xref** and search the container's number list for the requested number (pdfbox
  `PDFObjectStreamParser` behaviour); warn.
* Streams cannot appear inside object streams; a `stream` keyword there is a parse error for that
  object only (the rest of the container still loads).
* A container that fails to decode marks all of its objects as missing and records a warning; it
  does not fail the document.

### 2.5 Filters (`filter.go`) — pinned from evidence

**We decode exactly one filter: `FlateDecode`** (plus the predictors that accompany it, plus
`LZWDecode`, `ASCIIHexDecode`, `ASCII85Decode` and `RunLengthDecode`, which are 30-line functions and
which appear in xref-stream and `/DSS` positions in the wild). Image filters are **never** decoded.

Justification, from §1 and from the call sites: the only decoded bytes DSS ever consumes are
(a) xref streams, (b) object streams, (c) `/DSS` `/Certs` `/CRLs` `/OCSPs` token streams,
(d) `/VRI /TS`, and (e) `isStreamFill`'s non-emptiness test. (a)–(d) are Flate or unfiltered in every
corpus document. Stream *comparison* in `DefaultPdfObjectModificationsFinder` is on **raw** bytes.

```go
const (
    FilterFlate     Name = "FlateDecode"
    FilterLZW       Name = "LZWDecode"
    FilterASCIIHex  Name = "ASCIIHexDecode"
    FilterASCII85   Name = "ASCII85Decode"
    FilterRunLength Name = "RunLengthDecode"
    FilterCrypt     Name = "Crypt"     // /Identity only; anything else → ErrUnsupportedFilter
)
```

Anything else — `DCTDecode`, `CCITTFaxDecode`, `JPXDecode`, `JBIG2Decode` — makes
`Document.StreamData` return `ErrUnsupportedFilter` wrapping the filter name. `RawStreamData`
always succeeds. `pades` treats `ErrUnsupportedFilter` exactly as pdfbox's callers treat an
`IOException` from `getStreamBytes`: log and continue.

**FlateDecode leniency is a hard behavioural contract, not a nicety.** `FlateFilterDecoderStream`
(pdfbox 3.0.7) does three things we must copy exactly:

1. It **unconditionally discards the first two bytes** and inflates with `new Inflater(true)` —
   raw DEFLATE, no zlib header validation, **no Adler-32 check**. Go equivalent:
   `flate.NewReader(bytes.NewReader(raw[2:]))`. Never `zlib.NewReader` — it would reject the
   corpus's non-conforming headers and would fail on truncated checksums.
2. Fewer than 2 bytes of input ⇒ empty output, no error.
3. On a data-format error it **returns the bytes decoded so far** and logs a warning. Go: on any
   error from the flate reader, return the accumulated buffer with a recorded `Warning`, nil error.

Predictors (`Predictor.decodePredictorRow`): `/Predictor`, `/Colors` (default 1), `/BitsPerComponent`
(default 8), `/Columns` (default 1). `1` = none. `2` = TIFF. `≥10` = PNG, where each row is prefixed
by its own filter-type byte (so a declared `15` still dispatches per row over 0–4). Row length is
`ceil(Colors*BPC*Columns/8)`; `bpp = ceil(Colors*BPC/8)`. A truncated final row is zero-padded and
warned about.

### 2.6 Encryption (`crypt.go`) — the policy, decided

The corpus forces this: 8 encrypted documents, and `validation/encrypted.pdf` is in
`AbstractPDFDocumentValidatorTest`'s standard set, i.e. **validating an encrypted signed PDF is a
supported DSS feature**, not an edge case.

**Handled.** `/Filter /Standard` only:

| `/V` | `/R` | Algorithm | Corpus |
| --- | --- | --- | --- |
| 1 | 2 | RC4 40-bit | — (implement; trivial once V2 exists) |
| 2 | 3 | RC4 40–128 bit | 1 doc |
| 4 | 4 | `/CF /StdCF /CFM` ∈ {`/V2` (RC4), `/AESV2` (AES-128-CBC)}, `/StmF`/`/StrF` ∈ {`/StdCF`, `/Identity`} | 6 docs |
| 5 | 5, 6 | `/AESV3` AES-256-CBC, SHA-256/384/512 key derivation (R6 hardened hash) | 1 doc |

**These sets are closed, and `crypt.go`'s `checkCryptFilters` is what closes them.** The table is the
whole of what the handler implements, not a list of what it has been seen to do: a `/V` outside
{1, 2, 4, 5}, a `/V 5` paired with an `/R` other than 5 or 6 (the `/V 5` row's "5, 6" is enforced, not
descriptive — `checkCryptFilters` rejects any other pairing, because `computeEncryptionKey` and the
postcondition guard below both key off `/R` alone and would otherwise derive and apply an AES-128-shaped
key to a document declaring `/AESV3`), or a crypt filter a `/V 4`/`/V 5` document actually *selects*
(named by `/StmF` or `/StrF`, and not `/Identity`) that does not resolve through `/CF` to one of the
listed `/CFM` values, is refused rather than approximated. There is no `/V 0` row because ISO 32000-1
Table 20 says `/V 0` "shall not be used"; `/V` is optional with default 0, so an `/Encrypt` that omits
it lands in the same place.

Closing those sets closes every "declares one cipher, applies a different one" *shape*, but the
`/CFM` and `KeyLength` a conforming document reports through `Encryption()` are still what the
document *declares*, not always a description of the cipher width `objectKeyFor` actually applies:
`/V 4` with `/R 5` or `/R 6` and `/CFM /AESV2`, given a genuine `/UE`/`/OE` unwrap, opens reporting
`{CFM: AESV2, KeyLength: 128}` while `objectKeyFor`'s `useAES && len(h.key) == 32` fast path
(Algorithm 1.A) applies the 32-byte file key directly, i.e. AES-256 — exact parity with
`SecurityHandler.encryptData`'s own `if (useAES && encryptionKey.length == 32)` branch
(`SecurityHandler.java:221`), not a bug. `useAES` plus the file key's length are the only fields that
settle which cipher width is actually used; see the postcondition guard's comment in `setupEncryption`
(`crypt.go`).

**Rejected**, with a typed error that `pades` maps onto `InvalidPasswordException` /
`ProtectedDocumentException`:

* `/Filter` other than `/Standard` (public-key / PKCS#7 handlers) → `ErrUnsupportedSecurityHandler`.
* Wrong or missing password → `ErrInvalidPassword`. The password is tried as the owner password
  first, then as the user password, exactly as `StandardSecurityHandler.prepareForDecryption` does
  (`isOwnerPassword` before `isUserPassword`); which one matched determines
  `Permissions.OwnerAccess`.
* For `/R 5`/`/R 6`, an `/Encrypt /Perms` that is absent, is not a 16-byte string, or does not
  decrypt under the recovered file key to the `'a' 'd' 'b'` marker, the dictionary `/P` and the
  dictionary `/EncryptMetadata` → `ErrInvalidPassword` (ISO 32000-2 Algorithm 13, `validatePerms`).
  The check runs on both password branches, immediately before the handler is installed, so it
  cannot be bypassed. **This is a deliberate divergence**: pdfbox's
  `StandardSecurityHandler.validatePerms` makes the same three comparisons but answers each failure
  with `LOG.warn` and loads the document anyway, driving `AccessPermission` from the unauthenticated
  `/P`. For `/R 5`/`/R 6` `/P` is not mixed into the file key, so `/Perms` is the only thing that
  authenticates it and a warning leaves `PdfPermissionsChecker`'s `CanCreateSignatureField` gate
  defeated by a one-byte edit. See the `// DIVERGENCE, deliberate:` note on `validatePerms` in
  `crypt.go`.
* `/EncryptMetadata false` is honoured (metadata streams left in the clear).
* **The password is text, and the bytes hashed are pdfbox's bytes.** `Options.Password` is the
  password as UTF-8 — the Go shape of the Java `String` that `Loader.loadPDF(bytes, password)`
  receives — and `crypt.go`'s `passwordBytes` reproduces the charset step of
  `StandardSecurityHandler.prepareForDecryption`: `String.getBytes(ISO_8859_1)` for `/R 2`–`/R 4`
  (a code point above U+00FF becomes one `'?'`, Java's replacement byte; so does a malformed UTF-8
  sequence, which is what Java makes of an unpaired surrogate), `String.getBytes(UTF_8)` for `/R 5`,
  and `SaslPrep.saslPrepQuery` then UTF-8 for `/R 6` (PDFBOX-4155). `saslprep.go` ports pdfbox's
  `SaslPrep` table for table, including its `(char)` truncation of supplementary code points in two
  of the prohibition checks, because the byte string that has to match is the one pdfbox hashes. A
  password SASLprep prohibits (a control character, private-use or non-character code point, mixed
  bidirectional text) is `ErrProhibitedPassword`, the Go shape of the `IllegalArgumentException`
  pdfbox lets escape there — not `ErrInvalidPassword`, since no retry can help. The encryption side
  (`prepareDocumentForEncryption`) encodes the same way, with SASLprep's stored-string profile
  (`saslPrepStored`, which differs only in rejecting unassigned code points), so a document pdfbox
  protects with `café` opens here with `café`; `testdata/password/` holds pdfbox-generated goldens
  for `/R 3`, `/R 4` and `/R 6` and `password_kat_test.go` pins them (§5.5). One residue cannot be
  closed: the bidi step classifies code points by x/text's Unicode tables where pdfbox uses the
  JDK's, so a code point assigned in one Unicode version and not the other may be judged
  differently; unassigned code points are treated as Java does (no directionality). Without this
  step a non-ASCII password on the commonest encryption in the wild (`/V 4 /R 4`) opened in Java
  DSS and was `ErrInvalidPassword` here.
* **A crypt filter a `/V 4` or `/V 5` document selects and that this handler cannot identify** →
  `ErrUnsupportedSecurityHandler`. "Selects" means named by `/StmF` or `/StrF` and not `/Identity`.
  Each such name must resolve through `/CF` to a `/CFM` of `/V2`, `/AESV2` or `/AESV3`; both selected
  filters must resolve to the **same** one (the handler carries one key and one AES flag, so it
  cannot honour two — ISO 32000-1 itself permits distinct filters for streams and strings, this is a
  limitation of *this* implementation, not a pdfbox divergence); and for `/V 5` that one must be
  `/AESV3`. Both filters `/Identity` selects nothing and stays legal. **Deliberate divergence** —
  pdfbox keeps `useAES` false and RC4s the 32 bytes unwrapped from `/UE`, surfacing only as an
  unrelated parse error downstream; see `checkCryptFilters`'s `// DIVERGENCE, deliberate:` note in
  `crypt.go`.
* **`/AESV3` under an `/R` other than 5 or 6, and `/V 5` under an `/R` other than 5 or 6** →
  `ErrUnsupportedSecurityHandler`. ISO 32000-2 defines no other pairing. Both halves are needed:
  `/AESV3` promises a 32-byte key, and the only derivation that produces one is the `/UE`/`/OE`
  unwrap, which `computeEncryptionKey` selects on `/R` alone. Constraining `/V 5` alone would leave
  the identical downgrade one `/V` value away — a `/V 4` `/R 4` document naming `/CFM /AESV3` reports
  `/AESV3` with `/KeyLength 256` and applies AES-128 under an MD5 key, on read and on write.
  Not a new divergence in itself: such a document takes `computeEncryptionKey`'s `/R`-keyed
  `computeKeyRev234` (MD5) path exactly as pdfbox's `dicRevision`-keyed `prepareForDecryption` does,
  deriving and applying the same short key both implementations would — but it is the last surviving
  "declares AES-256, applies AES-128" shape issue #33 exists to close, and this port refuses it
  rather than reproduce it. See `checkCryptFilters`'s `// DIVERGENCE, deliberate:` note.
* **`/V 0`, including an `/Encrypt` with no `/V` at all** → `ErrUnsupportedSecurityHandler`.
  **Deliberate divergence** — pdfbox accepts it as `/Length`/8 RC4 (except for `dicLength`, where it
  uses 5 bytes only for `/V 1`, so pdfbox's own `/V 0` key length already disagreed with this port's
  pre-existing `case 0, 1: keyLenBytes = 5`; rejecting `/V 0` retires that silent disagreement too).
  See the `// DIVERGENCE, deliberate:` note opening the `default:` arm of `setupEncryption` (`crypt.go`).
* **`/UE` or `/OE` that is not exactly 32 bytes**, for `/R 5`/`/R 6` → `ErrUnsupportedSecurityHandler`
  (a malformed `/UE`/`/OE` is not a wrong password, so it is not `ErrInvalidPassword`). ISO 32000-2
  §8.7.4.1 fixes both at exactly 32 bytes (the wrapped AES-256 file key, two whole AES blocks).
  **Deliberate divergence** — pdfbox hands `fileKeyEnc` straight to `Cipher.getInstance("AES/CBC/
  NoPadding").doFinal(fileKeyEnc)`, which is equally lenient only over lengths that are already a
  whole multiple of the AES block size (any other length throws `IllegalBlockSizeException`, rethrown
  as `IOException`); this port cannot afford even that narrower leniency, because AES-CBC decrypts
  each block independent of the ones after it, so a padded `/UE` still recovers the genuine 32-byte
  key in its first two blocks, defeating every `len(key) == 32` check downstream by the padding
  length alone. See `computeEncryptionKey`'s `// DIVERGENCE, deliberate:` note in `crypt.go`.

Never-encrypted objects, matching `SecurityHandler.decrypt`:

* the `/Encrypt` dictionary itself,
* strings inside the trailer's `/ID` array,
* **the `/Contents` string of a dictionary whose `/Type` is `/Sig` or `/DocTimeStamp`** — pdfbox
  guards this explicitly, and getting it wrong silently corrupts every signature in an encrypted
  document,
* cross-reference streams (never encrypted, by spec),
* streams with `/Type /XRef` or a `/Crypt /Identity` filter.

Permissions are decoded from `/P` into a `Permissions` bitfield exposing exactly the four predicates
`PdfPermissionsChecker` needs.

The writer's re-encryption is §3.2 R20.

### 2.7 Lenient parsing — the exact tolerances, enumerated from pdfbox 3.0.7 sources

pdfbox runs `isLenient = true` by default and `Loader.loadPDF` never turns it off, so **lenient is
our only mode**. Every rule below was read out of `pdfbox-src/org/apache/pdfbox/pdfparser/`; the
citation is the source line, and each becomes a Go unit test with a hand-built input.

**Header** (`COSParser.parseHeader`, l.1620–1690):
* H1. Garbage before `%PDF-` is trimmed: the header search scans the first 1024 bytes for the marker
  and the file origin is moved there. *All subsequent offsets remain relative to byte 0 of the file*,
  which is what makes the `/ByteRange` arithmetic work.
* H2. `%PDF-` with fewer than 3 following bytes ⇒ version defaults to **1.4**.
* H3. Garbage after `%PDF-X.Y` on the same line is discarded. An unparseable version ⇒ **1.7**.

**Trailer / startxref** (`COSParser.parseTrailer`, l.500–560):
* T1. A missing `%%EOF` is tolerated; the search for `startxref` then runs to end-of-file.
* T2. Only the last 2048 bytes are searched.
* T3. A missing or unparseable `startxref` triggers full brute-force reconstruction (X2).
* T4. `LOG.warn("Expected trailer object at offset …")` — a `/Prev` that does not land on `xref` or
  on an `N G obj` whose object is `/Type /XRef` is repaired by X1.

**xref offsets** (`COSParser.checkXrefOffsets`/`validateXrefOffsets`, l.1195–1290) — this is the
single most important repair path, and it is what saves `validation/pdf-signed-corrupted.pdf`:
* X1. Every xref entry is validated by seeking to its offset and reading the `N G obj` header.
  If the *generation* differs, the key is corrected. If the *object number* differs or nothing
  parses, validation fails for the whole table.
* X2. On failure, a **brute-force scan** of the entire file for `\d+ \d+ obj` rebuilds the offset map
  from scratch and **replaces** the xref table wholesale. Last occurrence of a given object number
  wins (later revisions are later in the file). This is `BruteForceParser.getBFCOSObjectOffsets`.
* X3. A bad `startxref` or `/Prev` offset is repaired by searching the brute-force map for an `xref`
  keyword or an `/Type /XRef` object near the claimed offset; failing that the section is dropped
  (`return 0`), which ends the chain without failing the document.
* X4. A bad `/XRefStm` is repaired the same way, then skipped.

**Objects** (`COSParser.parseFileObject`, l.780–960):
* O1. `LOG.warn("Object (n:g) at offset … does not match")` — a mismatched object header is a
  non-fatal warning if the object number matches after generation correction; otherwise the object
  resolves to `null`.
* O2. `endobj` missing ⇒ tolerated.
* O3. A dangling reference (`R` to a number with no xref entry, or to a free entry) resolves to
  `null`, never an error. **`Document.Resolve` therefore returns `Null{}`, not an error, for a
  dangling reference** — this matters because `visitFromDictionary` skips nil-valued entries.

**Streams** (`COSParser.parseCOSStream`, l.860–960):
* S1. A missing or non-integer `/Length`, or a `/Length` that does not land on `endstream`, falls
  back to scanning forward for `endstream`; the dictionary's `/Length` is then **rewritten in memory**
  to the discovered length.
* S2. `endobj` where `endstream` was expected ⇒ warning, rewind, continue.
* S3. Extra bytes between the stream data and `endstream` ⇒ warning, rewind by the excess.
* S4. `stream` may be followed by `\r\n` or `\n`; a lone `\r` is accepted with a warning (spec
  forbids it).
* S5. `/Length 0` on a stream that clearly has data is treated as S1.

**Pages** (`COSParser`, l.1440): a `/Kids` entry that resolves to null is removed with a warning.

**The one *parsing* defect pdfbox does *not* recover from**, and neither do we: a `/Root` that is
missing or whose `/Pages` is not a dictionary. (§2.6 lists the encryption inputs this port refuses
where pdfbox does not; those are policy, not parsing.) `EmptyPage-corrupted.pdf` and
`EmptyPage-corrupted2.pdf` throw
`IOException: Page tree root must be a dictionary`. Our `Open` returns `ErrBrokenCatalog` for these
two, and the KAT asserts *failure on both sides* — a Go-side success where Java fails is a test
failure, exactly as in the `xmldom` oracle contract.

Every repair records a `Warning{Code, Offset, Message}`. `Document.Warnings()` is compared against
the oracle's warning-class list (not its English text).

**Resource guards** (not pdfbox behaviour — our own, because we are the network-facing side):
`MaxObjects` (default 5·10⁵), `MaxDepth` (default 512, guards recursive `/Kids` and self-referential
dictionaries), `MaxStreamSize` (default 512 MiB), and a global visited-set on `Resolve` so a cyclic
graph cannot spin. Exceeding a guard is a hard error, never a silent truncation.

### 2.8 Revisions and `/ByteRange` (`revision.go`)

Two *different* notions of "revision" live here and must not be confused.

**(a) The `%%EOF` scan** — the port of `PAdESUtils.extractRevisions`, and the only one DSS actually
uses to find previous revisions. Byte-for-byte semantics, which are subtle enough to spell out:

Scan forward one byte at a time, accumulating into a line buffer. When the buffer equals exactly
`%%EOF` (5 bytes), a revision boundary is recorded at `position` (1-based count of bytes consumed),
then **one lookahead byte** is read: if it is `\n`, the boundary advances by 1; if it is `\r`, the
boundary advances by 1 and a *second* lookahead byte is read, advancing by 1 more if it is `\n`.
Then the buffer resets. The buffer also resets on any line-break byte or as soon as it exceeds 5
bytes. The result is `[]Revision{{End: boundary}}`, and each `Revision` denotes bytes `[0, End)`.

This is deliberately *not* a `/Prev`-chain walk, and it deliberately counts `%%EOF` occurrences that
appear inside object data. Reproducing the quirk is required: `PAdESUtils.getPreviousRevision`
picks the candidate whose length is the largest below `byteRange[0]+byteRange[1]`, and a different
revision list changes which document DSS reports as the signed original.

**(b) `XRefSections()`** — the `/Prev` chain from §2.3, used by the writer to match the previous
revision's xref style and by the oracle to report per-revision style. Not used for coverage.

`/ByteRange` handling in `internal/pdf` is deliberately thin: `SignatureDictionary.ByteRange` is
`[]int64` exactly as written in the file (length is *not* forced to 4 — `PdfSigDictWrapperFactory`
reads whatever size is there and `ByteRange.validate()` upstream is what rejects it). The helpers
`SignedRanges` and `ContentsRange` derive the covered spans:

```
signed  = [br[0], br[0]+br[1]) ∪ [br[2], br[2]+br[3])
contents(hex, incl. < >) = [br[0]+br[1], br[2])
```

`Document.SignatureCoversWholeDocument` reproduces `PdfBoxDocumentReader.isSignatureCoversWholeDocument`
including its arithmetic, which is *not* the obvious one — it computes
`(br[1]-br[0]) + (br[2]-br[1]-br[0]) + br[3]` and compares to the file length. Port the formula as
written; do not "fix" it, or documents upstream reports as fully covered will stop matching.

### 2.9 AcroForm, fields and annotations (`document.go`)

`SignatureFields()` reproduces `PDDocument.getSignatureFields()`:

1. `/Root /AcroForm /Fields`; recurse into `/Kids` when a node has no `/FT` of its own (inheritable
   attributes `/FT /Ff /V /DA` propagate from parent to kid).
2. Keep nodes whose effective `/FT` is `/Sig`.
3. Fully-qualified name = `/T` values of the ancestor chain joined with `.`, skipping nodes without
   `/T` (pdfbox `PDField.getFullyQualifiedName`).
4. Widgets: the field dict itself if it carries `/Subtype /Widget` (merged field+widget, the common
   case), else its `/Kids`.
5. `Value` is `/V` resolved to a dict; `ValueKey` is the *object key of the reference*, which is what
   deduplicates two fields pointing at one signature dictionary.
6. A `/Sig` field whose `/V` is absent or not a dictionary is reported with `Value == nil` — that is
   an *empty* signature field, which `getAvailableSignatureFields` returns and
   `findExistingSignatureField` fills.

Page geometry: `Page(i)` is **1-based** (matching DSS's `ImageUtils.DEFAULT_FIRST_PAGE == 1`).
`/MediaBox` and `/Resources` are inheritable through `/Parent`. `/Rotate` is normalised into
`{0,90,180,270}` by `((r % 360) + 360) % 360` rounded to the nearest multiple of 90.

---

## 3. Writer

### 3.1 The one invariant

> **The output is the input's bytes, unchanged, followed by an appended increment.**

Not "logically equivalent". Byte-identical prefix. The writer takes the original bytes and an
`*Updater`, emits `original || increment`, and every offset it computes is absolute in the
concatenation. There is no code path that seeks backwards into the original. This is what makes
prior signatures survive, and it is checked by a test that asserts
`bytes.Equal(out[:len(in)], in)` on every writer KAT.

### 3.2 Byte-determinism rules — **PINNED**

These exist so that two runs, two machines and two Go versions produce identical bytes, so KAT
goldens are stable, and so the `/ByteRange` arithmetic is decidable in advance. Rules R1–R14 also
reproduce pdfbox's `COSWriter` formatting, because that formatting lands *inside the signed byte
range* and matching it makes divergence easy to spot in a diff. R15 is a deliberate deviation.

**R1 — EOL.** The end-of-line byte is `\n` (0x0A). Exceptions, and only these: each xref-table entry
ends `\r\n`; the `stream` keyword is followed by `\r\n`; stream data is followed by `\r\n` before
`endstream`.

**R2 — EOL suppression.** The output writer tracks an `onNewLine` flag. `writeEOL()` emits nothing
if the last byte written was already `\n` or `\r`. This is `COSStandardOutputStream.writeEOL` and it
is *load-bearing*: without it, every dictionary gains a blank line and no golden matches.

**R3 — Dictionary.**
```
<<\n
/Key<SP>value<EOL>
… (insertion order; entries whose value is nil are skipped entirely)
>><EOL>
```

**R4 — Array.** `[` then items separated by `<SP>`, except that after every 10th item the separator
is `<EOL>` instead; then `]` then `<EOL>`. (`visitFromArray`: `count % 10 == 0`.)

**R5 — Name.** `/` then, for each byte of the UTF-8 name: emit it verbatim if it is in
`[A-Za-z0-9+\-_@*$;.]`, else `#` + two **uppercase** hex digits. Note this is stricter than the PDF
spec (PDFBOX-2073) — `#` is emitted for `!`, `,`, `~`, `'` and so on. Copy the set exactly.

**R6 — String.** Hex form `<` + uppercase hex + `>` if any byte is ≥ 0x80, or is 0x0D or 0x0A, or if
`ForceHex` is set. Otherwise literal form `(` … `)` escaping **only** `(`, `)` and `\` with a
backslash. No octal escapes are ever emitted.

**R7 — Integer.** `strconv.FormatInt`, base 10.

**R8 — Real.** If `Real.Raw != ""`, write `Raw` verbatim. Otherwise `FormatReal(v float64) string`,
pinned to Java `Float.toString` semantics: shortest decimal that round-trips as a `float32`, always
containing a `.` (so `1` → `"1.0"`), and if that form would use exponent notation, the
`BigDecimal.stripTrailingZeros().toPlainString()` expansion instead. `FormatReal` gets its own
exhaustive table test (0, ±1, 0.5, 1e-7, 1e8, 595.276, 841.89, 1e20, -0.0, NaN→"0.0", ±Inf→"0.0").

**R9 — Reference.** `<num><SP><gen><SP>R`.

**R10 — Indirect object.** `<num><SP><gen><SP>obj<EOL>` body `<EOL>endobj<EOL>`.

**R11 — Stream.** The stream dictionary per R3 (with `/Length` set to the raw byte count as a direct
integer — never an indirect reference), then `stream\r\n`, then the raw bytes verbatim, then
`\r\nendstream<EOL>`.

**R12 — xref table.** `xref<EOL>`, then per contiguous subsection `<first><SP><count><EOL>`, then per
entry exactly 20 bytes: `%010d<SP>%05d<SP>[n|f]\r\n`. Subsections are computed by grouping sorted
object numbers into maximal contiguous runs (`COSWriter.getXRefRanges`). An incremental update
**always** prepends the free head entry for object 0: `0000000000<SP>65535<SP>f\r\n` as its own
`0 1` subsection when object 0 is not otherwise present.

**R13 — Trailer.** `trailer<EOL>` then the trailer dictionary, which is the previous trailer with:
`/Prev` = the previous revision's `startxref` value; `/Size` = highest object number written or
carried, plus 1; `/DocChecksum` **removed**; `/XRefStm` **removed** whenever we emit a table;
`/ID` forced to a *direct* array. `/Encrypt`, `/Root` and `/Info` are carried unchanged.

**R14 — Tail.** `startxref<EOL><offset><EOL>%%EOF<EOL>`.

**R15 — Object write order: ascending object number.** *This is a deliberate deviation.* pdfbox
writes objects in the BFS order of its `objectsToWrite` deque, which depends on `HashSet` iteration
and is not reproducible. We sort. Offsets in the xref are computed from actual write positions, so
correctness is unaffected and nothing about validity changes. **We never emit an object stream on
write** — every new object is a plain `N G obj` — which keeps the increment human-diffable and
removes an entire class of nondeterminism.

**R16 — Object-number allocation.** New numbers start at `Document.HighestObjectNumber() + 1` (the
maximum over *all* xref sections, including free entries and `/Size - 1`) and increase by one per
`Updater.Add` call, in call order. Generation is always 0. Free-list reuse is not implemented — a
reused number in an incremental update is a well-known interoperability hazard.

**R17 — `/Contents` placeholder.** `<` + `2 × ContentSize` ASCII `'0'` + `>`. Default
`ContentSize = 9472` (`0x2500`), matching both `PAdESSignatureParameters.signatureSize` and
`SignatureOptions.DEFAULT_SIGNATURE_SIZE`. The real CMS is written as **uppercase** hex
(pdfbox `Hex.getBytes`) starting at the byte after `<`; the remaining reserved bytes stay `'0'`.
*Note for the `pades` layer:* `PAdESUtils.replaceSignature` (the cached-to-be-signed path) uses
`Utils.toHex`, which is **lowercase** in both upstream `dss-utils` implementations. Both are legal
PDF hex strings and both parse. This proves upstream itself is not byte-stable across its two
signing paths, which is why §3.6 is worded as it is.

**R18 — `/ByteRange` placeholder and patch.** Written as `[0 1000000000 1000000000 1000000000]`
(pdfbox `RESERVE_BYTE_RANGE`). Record `byteRangeOffset` = position of the `[` **plus 1**, and
`byteRangeLength` = 35 (34 content bytes + the `]`). After the increment is fully laid out, format
`fmt.Sprintf("0 %d %d %d]", before, afterOffset, afterLength)` and copy it over those 35 bytes,
padding the remainder with `0x20`. If the formatted string exceeds 35 bytes, fail loudly. The four
values are
`before = sigContentsOffset`, `afterOffset = sigContentsOffset + sigContentsLength`,
`afterLength = totalLen - afterOffset`, where `sigContentsOffset` is the absolute offset of the `<`.
The corpus confirms the result shape exactly: `[0 1383 28009 307]` + 18 spaces.

**R19 — Only one signature per increment.** `Updater.AddSignature` returns an error on the second
call, matching `PDDocument.addSignature`'s `IllegalStateException`. Multiple signatures are multiple
increments.

**R20 — Encryption on write.** If the source document is encrypted, every string and stream in the
increment is encrypted with the **document's existing handler and key**, with the exceptions listed
in §2.6 (notably the signature `/Contents`). The handler is never upgraded, downgraded or
re-keyed — `/Encrypt` is carried by reference. AES-CBC initialisation vectors come from a
caller-supplied `io.Reader` (`Options.Random`, defaulting to `crypto/rand`), so the `pades` layer can
inject the deterministic `DSSSecureRandomProvider` seed and get reproducible output.

**R21 — Nothing else may vary.** No timestamps, no PRNG, no map iteration, no locale, no
`time.Now()`. `/ID` on an increment: the first element is carried from the source; the second is
`SignatureOptions.DocumentID` when supplied (the port of `generateDocumentId`, which is itself
deterministic from the deterministic-id + file size), else the first element repeated.

### 3.3 Signature placement (`sign.go`)

`AddSignature` reproduces `PDDocument.addSignature` + `PdfBoxSignatureService.createSignatureDictionary`:

1. Reserve `/Contents` (R17) and the placeholder `/ByteRange` (R18) **before** anything is laid out.
2. Build the signature dictionary. Key order is fixed and is the order pdfbox produces:
   `/Type /Filter /SubFilter /Name /Location /Reason /ContactInfo /M /Prop_Build /Reference
   /Contents /ByteRange`. `/M` and `/Name` are omitted for `/DocTimeStamp`.
3. Field: if `FieldID` names an existing field, it must exist, must be a `/Sig` field, and must be
   unsigned — otherwise an error whose text matches upstream's (`"The signature field '%s' can not
   be signed since its already signed."`). Set its `/V` and mark it updated. Otherwise create a new
   merged field+widget dictionary: `/FT /Sig`, `/Type /Annot`, `/Subtype /Widget`, `/F 4` (Print),
   `/T`, `/Rect`, `/P`, `/V`, and `/AP << /N <appearance stream> >>` when an appearance is supplied
   (invisible signatures get `/Rect [0 0 0 0]` and no `/AP`).
4. AcroForm: create `/AcroForm` if absent; append to `/Fields`; set `/SigFlags 3`
   (`SignaturesExist|AppendOnly`); mark `/AcroForm` and `/Fields` updated. `/AcroForm` is written
   direct inside the catalog when it was direct, indirect when it was indirect.
5. Page: append the widget to the page's `/Annots`, and **re-emit `/Annots` as a direct array** in
   the updated page object — pdfbox does this deliberately (`page.setAnnotations`) because an
   indirect `/Annots` that is not itself rewritten makes Adobe Reader report the document as
   modified.
6. `/Lock` on an existing field ⇒ build the `FieldMDP` `/Reference` array (`/Type /SigRef`,
   `/TransformMethod /FieldMDP`, `/TransformParams` = a copy of the `/Lock` dict with
   `/Type /TransformParams` and `/V /1.2`, direct; `/Data` = the catalog).
7. `DocMDP` (permission 1/2/3), only when no filled signature already exists in the document: build
   the `/Reference` array with `/TransformMethod /DocMDP`, and set `/Root /Perms /DocMDP` to the
   signature dictionary.
8. Mark the catalog updated.

Objects marked updated form the increment's write set, together with everything reachable from them
that is itself new. Reachability stops at objects that already exist unchanged in the input.

### 3.4 xref emission (`xrefwrite.go`)

Style selection reproduces `COSWriter.doWriteXRefInc` and is not a free choice:

```
if previousSectionStyle == XRefTable  ||  document.HasHybridXRef() {
        emit an xref TABLE + trailer      // hybrid always degrades to a table
} else {
        emit an xref STREAM               // /Type /XRef, /W [1 <n> 2], /Index, /Size, /Prev
}
```

For the stream form: `/Size` = highest object number + 2 (the extra one is the xref stream object
itself, which pdfbox allocates last — `pdfxRefStream.setSize(number + 2)`), `/Filter /FlateDecode`
with `/DecodeParms << /Predictor 12 /Columns <w1+w2+w3> >>`, `/W [1 <bytes needed for the largest
offset> 2]`, `/Index` listing the contiguous runs actually written, `/Prev` = previous `startxref`.
The xref stream object carries the trailer keys (`/Root /Info /ID /Encrypt`) and is **not** encrypted.

### 3.5 `/DSS`, `/VRI` and DocTimeStamp (`dss.go`)

`SetDSSDictionary` writes `/Root /DSS` and marks the catalog updated. Rules:

* Each token is a new plain stream object containing the DER bytes **unfiltered** (no `/Filter`) —
  this is what `PdfBoxDocumentReader.createCOSStream` produces, and it keeps the increment readable.
* A token whose `TokenRef.Key` is non-zero is *not* re-serialised; the existing object is referenced.
  This is `getPdfObjectForToken`'s dedup and it is why `internal/pdf` must expose object keys for
  array elements.
* `/DSS` is `<< /Certs [refs] /CRLs [refs] /OCSPs [refs] /VRI << … >> >>`; empty arrays are omitted.
* `/VRI` keys are the uppercase base-16 SHA-1 of the signature (computed by `pades`, passed in as a
  string); each value is a direct dictionary `<< /Cert [] /CRL [] /OCSP [] /TU (D:…) /TS <stream> >>`.
* Ordering inside `/Certs`/`/CRLs`/`/OCSPs` and inside `/VRI` is **the caller's slice order**.
  `internal/pdf` never sorts and never deduplicates — dedup is upstream's `indexOf` check and belongs
  in `pades`, where the token identity is known.

A DocTimeStamp is `AddSignature` with `Type: "DocTimeStamp"`, `Filter: "Adobe.PPKLite"`,
`SubFilter: "ETSI.RFC3161"`, no `/M`, no `/Name`, no `/Reference`.

### 3.6 The output contract — decided, and this is the answer to "do we match pdfbox byte-for-byte?"

**Our target is output that upstream DSS validates, not output byte-identical to pdfbox's.**

Reasoning, from evidence rather than taste:

1. Nothing in the PDF or PAdES specification, and nothing in DSS's validation path, reads our
   formatting. Validation recomputes the digest over `/ByteRange`, so what matters is that the range
   is *correct and self-consistent*, not how the surrounding dictionary was spaced.
2. Byte-identity with pdfbox is **not achievable even in principle**: pdfbox's object write order
   comes from `HashSet`/deque iteration (R15), and upstream DSS's own two signing paths already
   disagree on `/Contents` hex casing (R17). Chasing identity would mean reproducing a bug.
3. Byte-identity is also not *useful*: it would pin an artefact of pdfbox 3.0.7 that will move on the
   next upstream bump, exactly the trap the `xmldom` design avoided with Santuario's
   `circumventBug2650`.

So the three assertions are split across three KAT families, and each one is explicit about what it
proves (§5.3). What we *do* commit to, and test:

* **Determinism** — same inputs ⇒ same bytes, always. Golden files.
* **Prefix preservation** — `out[:len(in)] == in`. Always.
* **Structural parity** — pdfbox parses our output and reports the same revision count, xref style,
  signature inventory and `/ByteRange` values that we do.
* **Cryptographic parity** — upstream DSS validates our signatures to the same
  `SignatureLevel`/`Indication`/`SubIndication` as it does for its own.
* **Formatting parity where it is free** — R1–R14 match pdfbox because there is no cost to matching.

---

## 4. The exported Go API — **PINNED**

This surface is pinned. It does not change without an amendment to this
file. Every exported symbol below is final: name, signature, semantics.

### 4.1 `object.go` — object model *(written first, jointly reviewed, then frozen)*

```go
// Package pdf implements the subset of ISO 32000-1/2 that DSS's PAdES support needs.
package pdf

// Object is the sum type of PDF object kinds: Null, Bool, Integer, Real, String,
// Name, Array, *Dict, *Stream, Ref.
type Object interface{ pdfObject() }

type Null struct{}
type Bool bool
type Integer int64

// Real carries the literal it was parsed from; Raw is "" for values built in memory.
type Real struct {
    Val float64
    Raw string
}

// String is a PDF string. Hex records the source form and forces hex on output.
type String struct {
    Bytes []byte
    Hex   bool
}

type Name string
type Array []Object

// Ref is an indirect reference appearing as a value inside the object graph.
type Ref struct {
    Num int64
    Gen uint16
}

// ObjectKey identifies an indirect object slot. It is the analogue of upstream's
// PdfObjectKey / PdfBoxObjectKey. The zero value means "no key".
type ObjectKey struct {
    Num int64
    Gen uint16
}

func (k ObjectKey) IsZero() bool
func (k ObjectKey) String() string // "12 0"
func (r Ref) Key() ObjectKey

// Dict is an insertion-ordered PDF dictionary. Order is part of the contract:
// upstream's COSDictionary is a LinkedHashMap and PdfDict.list() exposes that order.
type Dict struct{ /* unexported */ }

func NewDict() *Dict
func DictOf(kv ...any) *Dict // alternating Name, Object; panics on a malformed pair

func (d *Dict) Len() int
func (d *Dict) Keys() []Name          // insertion order; a fresh slice
func (d *Dict) Has(key Name) bool
func (d *Dict) GetRaw(key Name) Object // unresolved; nil when absent
func (d *Dict) Set(key Name, v Object) // in-place for an existing key, append otherwise
func (d *Dict) Delete(key Name)
func (d *Dict) Clone() *Dict           // shallow: values are shared
func (d *Dict) String() string

// Stream is a PDF stream: its dictionary plus its still-encoded, still-encrypted bytes.
type Stream struct {
    Dict *Dict
    Raw  []byte
    // Offset and Length locate Raw in the source file; both are 0 for streams
    // built in memory.
    Offset int64
    Length int64
}

func NewStream(d *Dict, raw []byte) *Stream

// Rect is a PDF rectangle, normalised so Min <= Max on both axes.
type Rect struct{ MinX, MinY, MaxX, MaxY float64 }

func (r Rect) Width() float64
func (r Rect) Height() float64
func (r Rect) Array() Array
func RectFromArray(a Array) (Rect, bool)
```

### 4.2 `errors.go` — error taxonomy *(frozen with `object.go`)*

```go
package pdf

import "errors"

var (
    ErrNotPDF                    = errors.New("pdf: missing %PDF- header")
    ErrBrokenCatalog             = errors.New("pdf: page tree root must be a dictionary")
    ErrInvalidPassword           = errors.New("pdf: invalid password")
    ErrProhibitedPassword        = errors.New("pdf: password contains characters SASLprep prohibits")
    ErrUnsupportedSecurityHandler = errors.New("pdf: unsupported security handler")
    ErrUnsupportedFilter         = errors.New("pdf: unsupported stream filter")
    ErrLimitExceeded             = errors.New("pdf: resource limit exceeded")
    ErrNoSuchObject              = errors.New("pdf: object not found")
    ErrSignatureAlreadyAdded     = errors.New("pdf: only one signature may be added per increment")
    ErrContentsTooLarge          = errors.New("pdf: CMS does not fit the reserved /Contents space")
    ErrByteRangeTooLarge         = errors.New("pdf: /ByteRange does not fit the reserved space")
)

// FilterError names the filter that could not be decoded. errors.Is(err, ErrUnsupportedFilter).
type FilterError struct{ Filter Name }

func (e *FilterError) Error() string
func (e *FilterError) Unwrap() error

// WarningCode classifies a recovered defect. Stable across releases; the oracle
// compares codes, never message text.
type WarningCode string

const (
    WarnHeaderGarbage      WarningCode = "header-garbage"
    WarnHeaderVersion      WarningCode = "header-version-default"
    WarnMissingEOF         WarningCode = "missing-eof"
    WarnXRefOffsetRepaired WarningCode = "xref-offset-repaired"
    WarnXRefBruteForce     WarningCode = "xref-brute-force"
    WarnXRefEntryInvalid   WarningCode = "xref-entry-invalid"
    WarnXRefStmSkipped     WarningCode = "xrefstm-skipped"
    WarnObjectHeaderFixed  WarningCode = "object-header-fixed"
    WarnDanglingReference  WarningCode = "dangling-reference"
    WarnStreamLengthFixed  WarningCode = "stream-length-fixed"
    WarnStreamEndFixed     WarningCode = "stream-end-fixed"
    WarnFlateTruncated     WarningCode = "flate-truncated"
    WarnPredictorTruncated WarningCode = "predictor-truncated"
    WarnObjStmBroken       WarningCode = "objstm-broken"
    WarnKidRemoved         WarningCode = "kid-removed"
)

type Warning struct {
    Code    WarningCode
    Offset  int64
    Message string
}
```

### 4.3 `document.go` — the reader API

```go
package pdf

import (
    "io"
    "time"
)

type Options struct {
    // Password is the document's password as UTF-8 text, tried as the owner
    // password, then as the user password; hashed as pdfbox hashes it (§2.6).
    Password []byte
    // Random supplies AES initialisation vectors on write. nil means crypto/rand.
    Random io.Reader
    MaxObjects    int   // default 500000
    MaxDepth      int   // default 512
    MaxStreamSize int64 // default 512 << 20
}

type Document struct{ /* unexported */ }

func Open(r io.ReaderAt, size int64, opts *Options) (*Document, error)
func OpenBytes(b []byte, opts *Options) (*Document, error)

func (d *Document) Close() error
func (d *Document) Size() int64
func (d *Document) Bytes() ([]byte, error) // the whole source; used by the writer

// --- versions, trailer, catalog ---
func (d *Document) HeaderVersion() float32   // from %PDF-x.y
func (d *Document) Version() float32         // catalog /Version when present, else header
func (d *Document) Trailer() *Dict
func (d *Document) Catalog() (*Dict, error)
func (d *Document) Info() (*Dict, error)     // nil, nil when absent
func (d *Document) ID() [2][]byte

// --- object access ---
func (d *Document) Resolve(o Object) Object            // follows Ref chains; Null{} for dangling
func (d *Document) Object(k ObjectKey) (Object, error) // ErrNoSuchObject when absent
func (d *Document) ObjectKeys() []ObjectKey            // ascending Num, then Gen
func (d *Document) HighestObjectNumber() int64
func (d *Document) Warnings() []Warning

// --- typed dictionary accessors; each resolves indirect references ---
// The bool result is false when the key is absent or has the wrong type; that is
// never an error, matching pdfbox's null-returning getters.
func (d *Document) Get(dict *Dict, key Name) Object
func (d *Document) GetDict(dict *Dict, key Name) (*Dict, bool)
func (d *Document) GetArray(dict *Dict, key Name) (Array, bool)
func (d *Document) GetStream(dict *Dict, key Name) (*Stream, bool)
func (d *Document) GetName(dict *Dict, key Name) (Name, bool)
func (d *Document) GetString(dict *Dict, key Name) ([]byte, bool)
func (d *Document) GetInt(dict *Dict, key Name) (int64, bool)
func (d *Document) GetReal(dict *Dict, key Name) (float64, bool)
func (d *Document) GetBool(dict *Dict, key Name) (bool, bool)
// GetDate parses a PDF date string "D:YYYYMMDDHHmmSSOHH'mm'" leniently: any
// truncation from the right is accepted, as is a missing "D:" prefix.
func (d *Document) GetDate(dict *Dict, key Name) (time.Time, bool)
// RefAt returns the key of an indirect reference stored at key, if it is one.
// This is upstream's PdfDict.getObjectKey.
func (d *Document) RefAt(dict *Dict, key Name) (ObjectKey, bool)
// IndexRef is the array-element form: upstream's PdfArray.getObjectKey(i).
func (d *Document) IndexRef(a Array, i int) (ObjectKey, bool)

// --- streams ---
func (d *Document) StreamData(s *Stream) ([]byte, error) // decrypted + fully decoded
func (d *Document) RawStreamData(s *Stream) ([]byte, error) // decrypted, still encoded
func (d *Document) RawStreamSize(s *Stream) int64           // -1 when s is nil

// --- encryption / permissions ---
func (d *Document) IsEncrypted() bool
func (d *Document) Encryption() *Encryption // nil when not encrypted
func (d *Document) Permissions() Permissions

type Encryption struct {
    Handler   Name // always "Standard" (others are rejected at Open)
    V, R      int
    KeyLength int
    StmF      Name // "StdCF" | "Identity"
    StrF      Name
    CFM       Name // "V2" | "AESV2" | "AESV3"; "None" only when both filters are /Identity
}

type Permissions struct {
    Raw                int32
    OwnerAccess        bool // the owner password matched
    CanModify          bool // bit 4
    CanModifyAnnots    bool // bit 6
    CanFillInForm      bool // bit 9
    CanPrint           bool // bit 3
    CanExtract         bool // bit 5
    CanAssemble        bool // bit 11
    CanPrintFaithful   bool // bit 12
}

// --- pages and annotations (1-based page numbers) ---
func (d *Document) NumberOfPages() int
func (d *Document) Page(page int) (*Dict, ObjectKey, error)
func (d *Document) PageBox(page int) (Rect, error)   // /MediaBox, inherited
func (d *Document) PageRotation(page int) int        // normalised to 0/90/180/270
func (d *Document) Annotations(page int) ([]Annotation, error)

type Annotation struct {
    Key    ObjectKey
    Dict   *Dict
    Rect   Rect
    Name   string // /T
    Signed bool   // /V present
    Hidden bool   // /F bit 2
    NoRotate bool // /F bit 5
}

// --- AcroForm and signatures ---
func (d *Document) AcroForm() (*Dict, bool)
func (d *Document) SignatureFields() ([]SignatureField, error)
func (d *Document) SignatureDictionaries() ([]SignatureDictionary, error)

type SignatureField struct {
    Key      ObjectKey
    Dict     *Dict
    Name     string     // fully-qualified /T, dot-joined
    Value    *Dict      // /V, nil for an empty field
    ValueKey ObjectKey  // key of the /V reference; zero when /V is direct or absent
    Widgets  []*Dict
    WidgetKeys []ObjectKey
    Page     int        // 0 when the widget is not on any page
    Rect     Rect
    Lock     *Dict      // /Lock, nil when absent
}

// SignatureDictionary is the raw /V content. All ETSI semantics live in `pades`.
type SignatureDictionary struct {
    Key       ObjectKey
    Dict      *Dict
    Type      Name   // "Sig" | "DocTimeStamp" | ""
    Filter    Name
    SubFilter Name
    Contents  []byte  // decoded string bytes, i.e. the DER CMS
    ByteRange []int64 // verbatim, length not forced to 4
    Fields    []int   // indices into the SignatureFields slice that point here
}

func (d *Document) SignatureCoversWholeDocument(sd SignatureDictionary) bool

// --- revisions ---
// Revision is one %%EOF-delimited prefix, per PAdESUtils.extractRevisions.
type Revision struct {
    Index int
    End   int64 // the revision is bytes [0, End)
}

// ScanRevisions reproduces PAdESUtils.extractRevisions byte for byte, including
// its %%EOF + EOL lookahead. It does not walk /Prev.
func ScanRevisions(r io.Reader) ([]Revision, error)

func (d *Document) Revisions() []Revision
func (d *Document) XRefSections() []XRefSection // newest first
func (d *Document) StartXref() int64
func (d *Document) HasHybridXRef() bool

type XRefStyle uint8

const (
    XRefTable XRefStyle = iota + 1
    XRefStream
)

func (s XRefStyle) String() string

type XRefSection struct {
    Offset    int64
    Style     XRefStyle
    Trailer   *Dict
    Prev      int64 // -1 when absent
    XRefStm   int64 // -1 when absent
    Entries   int
    Recovered bool
}

// --- /ByteRange helpers ---
// SignedRanges returns the two covered spans as [start,end) pairs.
func SignedRanges(br []int64) ([2][2]int64, error)
// ContentsRange returns the span of the /Contents hex string including its < >.
func ContentsRange(br []int64) ([2]int64, error)
```

### 4.4 `filter.go` — decoding, exported for tests and for the writer

```go
package pdf

// Decode applies the filter chain named by /Filter with the parameters in
// /DecodeParms (or the /F, /DP abbreviations) to raw. Unsupported filters yield
// a *FilterError. Any recovered defect appends to warn.
func Decode(raw []byte, filters []Name, parms []*Dict, warn *[]Warning) ([]byte, error)

// FlateDecode reproduces pdfbox FlateFilterDecoderStream exactly: it discards the
// first two bytes, inflates as raw DEFLATE with no checksum validation, and on a
// corrupt stream returns the bytes decoded so far with a WarnFlateTruncated
// warning and a nil error.
func FlateDecode(raw []byte, warn *[]Warning) []byte

// FlateEncode produces zlib-wrapped DEFLATE at the fixed compression level 6, so
// output is byte-stable across Go versions that keep the flate algorithm stable.
func FlateEncode(data []byte) []byte

// ApplyPredictor undoes /Predictor: 1 = none, 2 = TIFF, 10..15 = PNG (per-row tag).
func ApplyPredictor(data []byte, predictor, colors, bpc, columns int, warn *[]Warning) []byte
```

### 4.5 `writer.go` — serialization primitives

```go
package pdf

import "io"

// Writer emits PDF syntax under the byte-determinism rules R1..R14 of DESIGN.md.
// It tracks column state so writeEOL is suppressed at a line start (R2).
type Writer struct{ /* unexported */ }

func NewWriter(w io.Writer) *Writer

func (w *Writer) Pos() int64
func (w *Writer) Err() error

func (w *Writer) WriteObject(o Object) error          // dispatches on kind
func (w *Writer) WriteIndirect(k ObjectKey, o Object) error // "N G obj … endobj"
func (w *Writer) WriteEOL() error                     // suppressed at a line start
func (w *Writer) WriteRaw(b []byte) error

// FormatReal renders v with Java Float.toString semantics (R8).
func FormatReal(v float64) string

// EncodeName renders a name under R5; EncodeString renders a string under R6.
func EncodeName(n Name) []byte
func EncodeString(s String) []byte
```

### 4.6 `incremental.go` + `sign.go` + `dss.go` — the writer API

```go
package pdf

import (
    "io"
    "time"
)

// Updater accumulates an incremental update over a Document. The Document's bytes
// are never modified; Write emits original||increment.
type Updater struct{ /* unexported */ }

func NewUpdater(d *Document) (*Updater, error)

// Alloc reserves the next object number (R16) without writing anything.
func (u *Updater) Alloc() ObjectKey
// Add allocates a key and schedules obj to be written.
func (u *Updater) Add(obj Object) ObjectKey
// Put schedules obj to be written under an existing key, replacing that object in
// this revision. It never touches the original bytes.
func (u *Updater) Put(k ObjectKey, obj Object)
// Catalog returns a mutable clone of the catalog, already scheduled for writing.
// Repeated calls return the same instance.
func (u *Updater) Catalog() *Dict
// Update returns a mutable clone of an existing object, already scheduled.
func (u *Updater) Update(k ObjectKey) (Object, error)
// Scheduled reports the keys that will be written, ascending (R15).
func (u *Updater) Scheduled() []ObjectKey

type SignatureOptions struct {
    Type        Name      // "Sig" (default) or "DocTimeStamp"
    Filter      Name      // default "Adobe.PPKLite"
    SubFilter   Name      // e.g. "ETSI.CAdES.detached", "ETSI.RFC3161"
    ContentSize int       // reserved /Contents bytes; default 9472
    SignerName  string    // /Name
    Reason      string
    Location    string
    ContactInfo string
    SigningTime time.Time // /M; zero value omits the key
    AppName     string    // /Prop_Build /App /Name
    FieldID     string    // fill this existing empty field; "" creates a new one
    Page        int       // 1-based; used only when creating a field
    Rect        Rect      // zero Rect creates an invisible field
    Appearance  *Stream   // /AP /N; nil for an invisible field
    DocMDP      int       // 1..3; 0 = none
    Lock        *Dict     // /Lock of the target field, drives FieldMDP
    DocumentID  []byte    // second element of /ID
}

type Placeholder struct {
    SigKey   ObjectKey
    FieldKey ObjectKey
}

// AddSignature installs a signature dictionary with placeholder /Contents and
// /ByteRange. It returns ErrSignatureAlreadyAdded on a second call (R19).
func (u *Updater) AddSignature(opts SignatureOptions) (*Placeholder, error)

type TokenRef struct {
    Key  ObjectKey // non-zero reuses an existing object; zero writes a new stream
    Data []byte    // DER; ignored when Key is non-zero
}

type VRIEntry struct {
    Name  string // uppercase base-16 SHA-1 of the signature
    Certs []TokenRef
    CRLs  []TokenRef
    OCSPs []TokenRef
    TU    time.Time
    TS    []byte
}

type DSSDictionary struct {
    Certs []TokenRef
    CRLs  []TokenRef
    OCSPs []TokenRef
    VRI   []VRIEntry // omitted from output when empty
}

// SetDSSDictionary writes /Root /DSS. Order is the caller's; nothing is sorted or
// deduplicated here.
func (u *Updater) SetDSSDictionary(dss DSSDictionary) error

// Result is the laid-out increment, before the CMS is inserted.
type Result struct {
    Bytes          []byte    // original || increment
    OriginalLength int64
    ByteRange      [4]int64  // zero-valued when the update carries no signature
    ContentsOffset int64     // absolute offset of the '<'
    ContentsLength int64     // reserved bytes including '<' and '>'
    StartXref      int64
    XRefStyle      XRefStyle
    HighestObjectNumber int64
}

// Write lays out the increment and patches /ByteRange (R18). /Contents still holds
// the placeholder.
func (u *Updater) Write() (*Result, error)

// SignedData is the byte stream the CMS must be computed over: exactly the two
// spans named by ByteRange.
func (r *Result) SignedData() io.Reader

// InsertContents writes cms as uppercase hex into the reserved /Contents span
// (R17). It returns ErrContentsTooLarge when cms does not fit.
func (r *Result) InsertContents(cms []byte) error

// ReplaceContents is the port of PAdESUtils.replaceSignature for the cached
// to-be-signed path: it finds the single all-zero hex placeholder in doc and
// substitutes cms. It errors when zero or more than one placeholder is present.
func ReplaceContents(doc []byte, cms []byte) ([]byte, error)
```

---

## 5. Oracle strategy and KATs

### 5.1 Shape

Same contract as the BouncyCastle and Santuario oracles used in phases 2–4, restated because it is
what makes the goldens trustworthy:

* The oracle is a Java program run **by hand**, never from `go test`. Go tests never invoke Java.
* Goldens are checked in. A golden that changes in a PR is a red flag and must be justified in the
  PR description.
* A Java-side failure is recorded in the golden as `!ERROR <SimpleClassName>` and the Go test asserts
  a corresponding Go-side failure. A Go success where Java failed is a **test failure**.
* `manifest.txt` records the SHA-256 of every golden and of every corpus file.

### 5.2 The oracle program

`internal/pdf/testdata/gen/PdfOracle.java`, header comment carrying the exact, proven command:

```
M2=$HOME/.m2/repository
CP=$M2/org/apache/pdfbox/pdfbox/3.0.7/pdfbox-3.0.7.jar\
:$M2/org/apache/pdfbox/pdfbox-io/3.0.7/pdfbox-io-3.0.7.jar\
:$M2/org/apache/pdfbox/fontbox/3.0.7/fontbox-3.0.7.jar\
:$M2/commons-logging/commons-logging/1.3.5/commons-logging-1.3.5.jar
javac -nowarn -cp "$CP" -d /tmp/pdforacle PdfOracle.java
java -Dorg.slf4j.simpleLogger.defaultLogLevel=off -cp "$CP:/tmp/pdforacle" \
     PdfOracle <path to your upstream DSS checkout>/dss-pades/src/test/resources <golden dir> <manifest>
```

The pdfbox 3.0.7 jars are fetched once with a two-dependency `pom.xml` and
`mvn dependency:build-classpath`; the sources jar (`-Dartifact=…:jar:sources`) is the citation
source for every rule in §2.7 and §3.2.

**Per-PDF dump fields** — one record per corpus file, tab-separated, stable field order:

| Field | Source | Which Go assertion it feeds |
| --- | --- | --- |
| `path`, `size` | filesystem | corpus identity |
| `eofRevisions` | raw `%%EOF` count | `ScanRevisions` length |
| `revisionEnds` | the `PAdESUtils.extractRevisions` boundaries | `ScanRevisions` exact offsets |
| `xrefChain` | `/Prev` walk, `table`/`stream`/`?` per hop | `XRefSections()[i].Style` |
| `hybrid` | `/XRefStm` present | `HasHybridXRef()` |
| `headerVersion`, `catalogVersion` | `COSDocument.getVersion()`, `PDDocument.getVersion()` | `HeaderVersion()`, `Version()` |
| `objects` | size of `COSDocument.getXrefTable()` | `len(ObjectKeys())` |
| `objectKeys` | sorted `num:gen` list | `ObjectKeys()` — the strongest single check |
| `trailer` | sorted trailer key names | `Trailer().Keys()` as a set |
| `filters` | sorted set of `/Filter` names over all streams | `Decode` coverage |
| `predictors` | `/DecodeParms /Predictor` values seen | `ApplyPredictor` coverage |
| `ObjStm` | any `/Type /ObjStm` | object-stream loading |
| `encrypted` | `Filter,V,R,Length,CFM` | `Encryption()` |
| `password` | which of `{none," "}` opened it | `Options.Password` |
| `permissions` | the four `AccessPermission` predicates + `isOwnerPermission` | `Permissions()` |
| `pages`, `pageBoxes`, `pageRotations` | `PDPage` | `NumberOfPages`, `PageBox`, `PageRotation` |
| `annots` | per page: `rect,name,signed` | `Annotations()` |
| `sigs` | per field: `name:Type,Filter,SubFilter,BR=[…],ContentsLen,ContentsSHA256` | `SignatureFields`, `SignatureDictionary` |
| `coversWholeDocument` | `PdfBoxDocumentReader.isSignatureCoversWholeDocument` | `SignatureCoversWholeDocument` |
| `dss` | `/DSS` present + `Certs/CRLs/OCSPs` counts + each element's object key + SHA-256 of its decoded stream | `/DSS` reading, `IndexRef`, `StreamData` |
| `vri` | `/VRI` key names in order + per-entry counts | `SingleDssDict.extractVRIs` parity — depends on dictionary order |
| `error` | exception class + message when pdfbox throws | negative assertions |

A prototype of this program already exists and has been run over the full corpus; §1's table is its
output. Productionising it means adding `revisionEnds`, `objectKeys`, `permissions`, page/annotation
geometry, `dss`/`vri` and the SHA-256 digests.

### 5.3 The three KAT families, and exactly what each asserts

**KAT-A — reader parity (`TestOracleCorpus`).** Corpus = all 248 upstream PDFs, copied into
`internal/pdf/testdata/corpus/` mirroring their upstream paths per `PORTING.md`. For each file, every
oracle field in §5.2 is compared. **Asserts**: our parse agrees with pdfbox on structure, on repair
outcomes, and on every byte we hand upward. **Does not assert** anything about output bytes.

**KAT-B — writer determinism (`TestWriterGolden`).** ~18 hand-authored scenarios, each a
`(input PDF, Updater script)` pair with a checked-in golden of the **increment only**:

| # | Scenario | What it pins |
| --- | --- | --- |
| 1 | invisible signature onto a table-xref 1.4 doc | R1–R14, R17–R18, table emission |
| 2 | invisible signature onto an xref-stream 1.6 doc | xref-stream emission, `/W`, `/Index`, predictor 12 |
| 3 | invisible signature onto a hybrid-xref doc | style degradation to a table (§3.4) |
| 4 | visible signature with an `/AP` stream, page 2 | widget, `/Rect`, `/Annots` direct re-emission |
| 5 | fill an existing empty field by `/T` | the fill-in branch, no new field |
| 6 | second signature on an already-signed doc | `/Prev` chain, `/Size` growth, prefix preservation |
| 7 | DocTimeStamp | `/Type /DocTimeStamp`, `/ETSI.RFC3161`, no `/M` |
| 8 | `/DSS` with 3 certs, 1 CRL, 2 OCSPs, no `/VRI` | §3.5 stream emission |
| 9 | `/DSS` + `/VRI` with two entries, one reusing an existing cert object | `TokenRef.Key` reuse |
| 10 | DocMDP = 2 | `/Perms /DocMDP`, `/Reference` |
| 11 | FieldMDP from a `/Lock` | `/TransformParams` copy semantics |
| 12 | developer extension into an empty `/Extensions` | catalog mutation via `Updater.Catalog` |
| 13 | AES-128 (`AESV2`) encrypted source, fixed IV reader | R20, `/Contents` left unencrypted |
| 14 | AES-256 (`AESV3`) encrypted source, fixed IV reader | R20 for V5/R6 |
| 15 | RC4-128 encrypted source | R20 for V2/R3 |
| 16 | a name needing `#` escapes and a string needing hex | R5, R6 |
| 17 | an array of 25 elements | R4's 10-per-line rule |
| 18 | `/ByteRange` values wide enough to nearly fill 35 bytes | R18's padding and its overflow error |

Each runs twice in-process and asserts identical bytes; each asserts `out[:len(in)] == in`.
**Asserts**: our writer is deterministic and formatted per §3.2. **Does not assert** identity with
pdfbox output — that is deliberate, per §3.6, and the test file says so in a comment.

**KAT-C — interop, both directions (`TestInterop`, build tag `interop`).**
* *Go → Java*: every KAT-B scenario's output is fed to the oracle; the oracle must parse it, report
  the expected revision count, xref style, signature inventory and `/ByteRange`, and — for the signed
  scenarios driven from the `pades` layer — upstream DSS must validate
  it to the expected `SignatureLevel` and `Indication`.
* *Java → Go*: every one of the 199 corpus PDFs that carries at least one filled signature
  dictionary (207 have signature fields; 8 of those hold only empty fields) is re-validated by the Go port and each
  signature's recomputed digest over `SignedRanges` must match what the CMS declares. That is 206
  independently-derived expected answers for free, and it is the strongest evidence the reader's
  offset arithmetic is right.

**Asserts**: real-world validity. This is the family that decides whether the port ships.

### 5.4 Unit tests that do not need an oracle

Every §2.7 tolerance gets a hand-built minimal input (usually under 400 bytes) and a Go test named
for its rule ID: `TestLenient_H2_NoVersion`, `TestLenient_X2_BruteForce`,
`TestLenient_S1_MissingLength`, and so on. These are cheap, they document the behaviour better than
prose, and they are the regression net when pdfbox is bumped. `FormatReal`, `EncodeName`,
`EncodeString`, `ApplyPredictor` and `FlateDecode` each get an exhaustive table test.

### 5.5 The password-charset goldens

`testdata/password/` is a second, smaller set of goldens with the same contract as §5.1 — produced
by a Java program run by hand, never from `go test` — but made by pdfbox itself rather than read
from the DSS corpus: `testdata/password/PasswordFixtures.java` (pdfbox 3.0.6; the charset code is
the same in 3.0.7) protects one empty page under `StandardProtectionPolicy` with non-ASCII
passwords for `/R 3`, `/R 4` and `/R 6`, and the PDFs it writes are the goldens.
`password_kat_test.go`'s `TestPasswordCharsetFixtures` pins them (they are in-module and not in
`manifest.txt`, since they are the fixtures, not derived records); its README lists each file's
parameters and passwords. A regenerated golden changes bytes (pdfbox draws fresh `/ID` and salts)
and needs the same PR justification as any other.

---

## 6. Review checklist

1. Does `Dict` preserve insertion order through `Set` on an existing key? (§2.2)
2. Does the writer ever seek backwards into the original bytes? It must not. (§3.1)
3. Is `writeEOL` suppressed at a line start? (R2 — the single most common cause of golden drift)
4. Are xref-table entries exactly 20 bytes, `\r\n`-terminated, with the object-0 free head? (R12)
5. Is `/DocChecksum` dropped and `/ID` forced direct in the trailer? (R13)
6. Is the object write order ascending object number, and is that deviation commented in the code
   with a pointer to R15?
7. Is `FlateDecode` skipping two bytes and inflating raw, with no Adler-32 check and partial output
   on corruption? (§2.5 — do not "fix" this to use `zlib`)
8. Does `Resolve` return `Null{}` — not an error — for a dangling reference? (§2.7 O3)
9. Is the signature `/Contents` string excluded from encryption on both read and write? (§2.6, R20)
10. Does `ScanRevisions` reproduce the `%%EOF` + `\r`/`\n`/`\r\n` lookahead exactly, including
    counting `%%EOF` sequences that occur inside object data? (§2.8)
11. Is `SignatureCoversWholeDocument` the upstream formula, not the obvious one? (§2.8)
12. Does anything in the package call `time.Now`, `math/rand`, or iterate a `map` into output? (R21)
13. Does `internal/pdf` import any DSS package? It must not. (§0.3)
14. Does every source file carry the provenance header `PORTING.md` requires?

---

## Encryption on write

R20 (append to encrypted documents) requires an encryption-on-write path.
The implemented surface is `writerEncryptHook` in
`incremental.go`: the Updater re-encrypts new/updated strings and streams with
the document's existing security handler (AESV2, AESV3, RC4-128) before
serialization, reusing the reader's decryption state. Verified empirically
against pdfbox on all three variants (see corpuswrite_test.go). Public API is
unchanged — encryption engages automatically for encrypted source documents.
