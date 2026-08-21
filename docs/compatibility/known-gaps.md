# Known gaps

The other half of [The numbers](numbers.md). This page lists what the port does
**not** do, and where it knowingly differs from Java DSS.

It is here because in a signature library an overclaim in the documentation is
as much a defect as a bug in the code. If you find something missing from this
list that should be on it, that is a bug worth
[reporting](https://github.com/ryftcore/dss-go/issues).

Every entry below was checked against the code in this repository, not copied
forward from an earlier document.

## Not ported at all

Whole upstream modules that are out of scope. These are absences, not
divergences — nothing here behaves differently from Java, it simply is not
present.

| Missing | What it means for you |
|---|---|
| **`dss-service` online sources** — `OnlineTSPSource`, `OnlineCRLSource`, `OnlineOCSPSource` | The library never fetches a time-stamp or revocation data by itself. You supply a `TSPSource` and revocation sources. The `esig` CLI ships its own RFC 3161-over-HTTP client for `-tsa`. Downloading trusted lists and fetching issuer certificates over AIA **do** work — those loaders live in `dss-spi` upstream and were ported. See [Timestamps and TSAs](../guides/timestamps-and-tsas.md). |
| **REST and SOAP remote services**, and their clients | No remote-signing server or client. Sign in-process. |
| **Evidence-record analyzers** — RFC 4998 ERS and RFC 6283 XMLERS | The framework, the interfaces and the report plumbing are ported, but nothing registers an `EvidenceRecordAnalyzerFactory`, so an evidence record cannot actually be parsed from a document today. |
| **`dss-cookbook`**, coverage and BOM modules | Documentation and build-tooling modules with no runtime role. This site and `dss/examples/` serve the cookbook's purpose. |

## Supported with limits

Present, but narrower than upstream.

### Key stores

- **PKCS#12 is fully supported**, including RSA, ECDSA, **Ed25519 and DSA**
  keys, through a native RFC 7292 reader written for this port and pinned by
  tests in `dss/token/pkcs12_signature_token_test.go`.
- **JKS is not supported.** Java's proprietary key store format has no Go
  implementation in the standard library or under `golang.org/x/…`, and the
  port's dependency policy allows nothing else. Every `JKSSignatureToken`
  constructor returns a clear "not supported in the Go port" error.
- **PKCS#11 is not supported.** Hardware tokens are reached in Java through the
  JCA `SunPKCS11` provider, which has no Go counterpart. Same treatment: the
  constructors return an explicit error rather than pretending.
- **The OS key stores are not supported** — the Windows certificate store
  (`Windows-MY`) and the macOS Keychain (`KeychainStore`) are JCA providers with
  no Go counterpart, and their tokens report that explicitly.

You can still drive a smart card or an HSM: implement the port's
`token.SignatureTokenConnection` interface over whatever driver you have. The
library never handles your private key directly — it asks the token to sign.

### Digest algorithms

**MD2 and WHIRLPOOL** are recognised as enumeration values (so documents naming
them parse, and OID and URI lookups work) but cannot be computed: neither has a
Go implementation in the standard library or under `golang.org/x/…`. They are
the only two.

Everything else upstream supports is present, **including the SHAKE
extendable-output functions** — `spi.DSSUtilsDigest` computes SHAKE-128 and
SHAKE-256 through `golang.org/x/crypto/sha3`, at the same output lengths
BouncyCastle's `SHAKEDigest` uses, and SHAKE256-512 is additionally available as
a `hash.Hash`. This mirrors upstream's own availability exactly: Java DSS also
reaches SHAKE-128 and SHAKE-256 directly rather than through
`MessageDigest.getInstance`, because BouncyCastle registers no JCA
`MessageDigest` under those two names either.

### PDF

**Visible signature appearances are not supported.** The signature field, its
`/Rect` and an empty appearance stream are produced correctly; painting an
image or a text block into it is not, because the native PDF engine has no
rasteriser, font engine or content-stream interpreter. Requests to render an
appearance return an explicit error rather than silently producing a blank one.

Everything else in the PDF path is present: incremental updates, byte ranges,
the DSS dictionary and VRI, document time-stamps, RC4 and AES encryption,
cross-reference tables and streams, object streams, revision extraction, and
modification-detection.

### CMS

**Streaming is in-memory only.** Upstream offers two interchangeable CMS
implementations — an in-memory one and a streaming one for very large
documents. This port has a single native implementation that follows the
in-memory behaviour. The streaming API surface is accepted for parity but not
consulted. Practically: signing a multi-gigabyte file will use memory
proportional to the file.

### EAA (electronic attestation of attributes)

The EAA validation blocks are ported and **gated behind the `eaa` build tag**,
off by default. They were verified as a mechanical one-for-one port of the Java
classes, but a behavioural oracle corpus for them has not been built, so they
do not carry the parity evidence every other part of the validation engine
does. Treat them as unproven until that corpus exists — it is a tracked item.

## Known divergences from Java DSS

Places where both implementations do something and the something differs.

### Six CAdES/CMS-layer parity gaps

A differential sweep over the **entire** upstream PAdES corpus — 248 PDFs, 298
signatures — found six divergences, all in the CMS/CAdES layer beneath PAdES,
all tracked as follow-ups:

1. Signature count on a deliberately badly-encoded CMS structure.
2. Certificate extraction from a PDF with unusual end-of-file handling.
3. Four legacy PKCS#7 reference-data cases.
4. PLAIN-ECDSA (BSI TR-03111) signature verification.
5. A wrong-digest-algorithm case.
6. Present-but-empty `/Reason` and `/Location` optionality.

**None of these is a false accept** — no case where this port declares valid
something Java declares invalid. Everything else in that sweep matched exactly:
verdicts, modification-detection classification, byte ranges, coverage. The
details, and the runner that found them, are documented in the repository at
`dss/pades/testdata/broadgen/README.md`.

### Deliberate ordering choices

A handful of places where upstream's output order comes from a
`java.util.HashMap` or `HashSet` iteration rather than from anything in a
specification. Most of these are reproduced exactly, including Java's bucket
ordering, because the reports are compared byte for byte. Where a
carve-out remains it is compared as a multiset — same contents, order not
asserted — and pinned by a test so it cannot quietly widen.

If you are diffing report XML against Java's output and see identical content
in a different order within one list, that is this. It never changes a
conclusion.

### Go toolchain sensitivity

Two behaviors follow the Go release you build with, not this port's own code.

- **Container framing bytes.** ASiC containers and PAdES incremental updates are
  byte-stable for a given Go toolchain, but the same input can produce different
  bytes across Go releases, because `compress/flate`'s encoder changes between
  them (Go 1.27 is such a release). Signatures stay valid: they cover entry
  contents and `ByteRange`s, never the container framing.
- **Unicode-category-dependent behavior.** The DN `CANONICAL` form (`x/text`
  casing plus NFKD), `IsStringDigits`, and the XPath 1.0 name lexer read the
  building toolchain's Unicode tables — Unicode 17 on Go 1.27 — while Java
  follows the JDK's, Unicode 15 on JDK 21. No corpus fixture is affected today.

## Not gaps, but frequently mistaken for them

- **`qualification=NA`.** Not "unqualified" — *not determined*. You did not
  supply trusted-list information. See
  [Trust, eIDAS and qualification](../concepts/trust-and-eidas.md).
- **`INDETERMINATE` on a signature you know is fine.** Almost always a missing
  trust anchor or missing revocation data. See
  [How validation works](../concepts/validation.md).
- **The library not reaching the network.** Deliberate, and documented in the
  package doc. Nothing here makes a request you did not ask for.
- **`SignOptions.Packaging` ignored for PDFs and ASiC.** Those formats have no
  packaging choice to make.

## Where these are tracked

- [`PORTING_PLAN.md`](https://github.com/ryftcore/dss-go/blob/main/PORTING_PLAN.md) —
  the historical record of the port, phase by phase, including the accepted-gaps
  list this page is derived from and verified against.
- [`UPSTREAM.md`](https://github.com/ryftcore/dss-go/blob/main/UPSTREAM.md) — the
  baseline pin and the procedure for moving to a newer DSS release.
- The repository's issue tracker, which has a dedicated **interop mismatch**
  template.

## Next

- [Methodology](methodology.md) — how the parity that *is* claimed is
  established.
- [The numbers](numbers.md) — what is measured, and how to reproduce it.
