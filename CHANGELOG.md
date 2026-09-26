# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
While the major version is 0, a minor release may contain source-incompatible
API changes; each one is listed under **Changed**. Release procedure:
`RELEASING.md`.

## [Unreleased]

## [0.1.0] - 2026-09-26

First tagged release. Initial public port of [esig/dss](https://github.com/esig/dss) (upstream
baseline: version 6.5.RC1, commit
`4c2129862948bfd53ca1455832260aa17e183cf8`) to Go, module
`github.com/ryftcore/dss-go/dss`. See `UPSTREAM.md` for the baseline pin,
the documentation site's compatibility pages for the verification detail
behind every claim below, and `dss/PORTING.md` for the porting
conventions applied throughout.

### Added

- **Signature formats — sign, extend, and validate** — ported, with each
  format family cross-validated against upstream Java DSS 6.5.RC1 in
  both directions. The two directions do not have the same reach, and
  the difference is worth stating plainly:
  - *Upstream Java-signed fixtures validated by this port*: at every
    level the fixtures cover, including LT and LTA.
  - *Go-signed documents validated by a running Java DSS*: at levels
    **B and T** for CAdES, XAdES, JAdES and ASiC, and at **B, T, LT and
    LTA** for PAdES (which additionally checks `/DSS` equivalence with
    Java-extended documents). ASiC also covers Java-built containers
    extended by this port and re-validated by Java.

  Fixture counts and results, and what runs in `go test` versus what was
  a one-time sweep, are separated in the repository's compatibility
  documentation. The formats:
  - **CAdES** (CMS Advanced Electronic Signatures), levels B/T/LT/LTA.
  - **XAdES**, levels B/T/LT/LTA, including a byte-exact native XML
    canonicalization (C14N 1.0/1.1/exclusive) and XML-DSig
    implementation.
  - **PAdES** (PDF signatures), levels B/T/LT/LTA, on a native PDF
    reader/incremental-update writer (no cgo, no external PDF library)
    covering xref tables/streams/hybrid files, object streams, standard
    filters, RC4/AES encryption, `/DSS` and `/VRI` dictionaries, and
    document timestamps.
  - **JAdES** (JSON Advanced Electronic Signatures / ETSI TS 119 182),
    levels B/T/LT/LTA, on a native JOSE/ordered-JSON implementation.
  - **ASiC** containers (ASiC-S and ASiC-E) wrapping CAdES and XAdES
    signatures, including container build/extract/merge and manifest
    handling.
- **Validation engine** implementing the EN 319 102-1 signature
  validation process (building blocks, basic/long-term-validation-data
  processes, POE handling, algorithm-obsolescence checks) with
  end-to-end Indication/SubIndication verdicts, and eIDAS qualification
  determination (trust-service filters, QSCD/type strategies, MRA,
  QWAC).
- **Reports**: SimpleReport, SimpleCertificateReport, DetailedReport,
  and ETSI TS 119 102-2 ValidationReport, plus the underlying
  diagnostic-data model — all XML-schema-conformant and marshal-parity
  tested against upstream JAXB output.
- **Trusted lists / LOTL**: TSL/LOTL/pivot parsing (including MRA
  extensions), trusted-list signature validation, and the
  `TLValidationJob` framework (cache state machine, alert firing),
  cross-validated against production EU LOTL data.
- **PKCS#12 (PFX) key store support** via a native parser
  (`internal/pfx`) built on this port's own BER/DER engine, covering
  RSA, EC, Ed25519, and DSA private keys and multi-entry key stores —
  see "Known limitations" below for key-store formats this does *not*
  cover.
- A compatibility test harness with byte-parity and verdict-parity
  oracles generated from upstream Java DSS, exercised per-package and,
  for the full signature-format matrix end to end, in
  `dss/harness` — see the compatibility documentation for the scope and
  scale verified (fixture counts, corpus sizes and parity results are
  recorded there; this file does not restate them, to avoid the two
  drifting apart).
- **Go 1.27 is the minimum supported release** (`go 1.27.0` in
  `dss/go.mod`). Earlier toolchains cannot build the module: the public
  API relies on generic methods, which Go 1.27 is the first release to
  allow.
- **Generic-method API finalization.** Java generic instance methods that
  the port had been carrying as package-level generic functions — a
  pre-1.27 language limitation, not a design choice — are now methods on
  their receiver, matching the Java shape: `CertificateWrapper`'s
  `CertificateExtensionForOid[T]` and
  `SignedDocumentDiagnosticDataBuilder`'s
  `CreateOrphanTokenFromRevocationIdentifier[R]` stay exported, while the
  ports of Java `protected` helpers whose callers are all in-package
  became unexported methods. Generic methods declared on Java
  *interfaces* still erase to the constraint's base type: Go forbids type
  parameters on interface methods. See `dss/PORTING.md`.
- **Password-protected (encrypted) PDFs through the facade and the CLI.**
  `dss.SignOptions`, `dss.ExtendOptions` and `dss.ValidateOptions` gained a
  `PasswordProtection` field — the facade's name for upstream's
  `PAdESSignatureParameters.setPasswordProtection` /
  `PDFDocumentValidator.setPasswordProtection` — so an encrypted PDF can
  be signed, extended and validated without dropping to the `pades`
  package; the signed document stays encrypted under the same password,
  as upstream's does. Setting it for anything but a PDF returns the new
  `dss.ErrPasswordProtectionNotApplicable`. The `esig` CLI exposes it as
  `-pdf-pass env:VARNAME` on `sign`, `extend`, `validate` and `inspect`,
  with the same env-var-only rule as `-p12-pass`. Encrypting a PDF that
  is not yet encrypted is, as upstream, not something the library does.
  Making the password a public contract also fixed how the native PDF
  engine hashes it: `internal/pdf` now treats the password as UTF-8 text
  and reproduces pdfbox's charset step (ISO-8859-1 for `/R 2-4`, UTF-8
  for `/R 5`, SASLprep then UTF-8 for `/R 6` — `SaslPrep` is ported
  table for table), where it previously hashed the raw bytes, so a
  non-ASCII password such as `café` on an RC4 or AES-128 document that
  Java DSS opens was refused as invalid. A password SASLprep prohibits
  is the new `pdf.ErrProhibitedPassword`, distinct from
  `ErrInvalidPassword` as upstream's exception is. Pinned by
  pdfbox-generated goldens under `internal/pdf/testdata/password/`.

### Fixed

- `esig <subcommand> -h` exits 0, like `esig -h`, instead of 2.
- The facade no longer panics on nil input.
  - `LoadCertificateBytes(nil)` returns an error.
  - `NewDocument(name, nil)` is an empty document.
  - `Sign` with a zero `Signer` returns a clear error.
  - `(*Signer)(nil).Close()` is a no-op.

### Changed

- **Exported type and function names no longer repeat their package name.**
  `alert.AlertHandler` is now `alert.Handler`, `asic.ASiCContent` is
  `asic.Content`, `xades.XAdESSignature` is `xades.Signature`,
  `cades.NewCAdESService` is `cades.NewService`. 403 types and functions
  across 29 packages (248 types, 155 functions), plus 313 constructors,
  `Abstract...` bases, `...Overrides` hooks and same-family constants that
  follow them.

  The rule, which is the one `revive`'s `exported` check applies: if the
  package name is a case-insensitive prefix of an exported identifier and the
  next character is `_` or upper-case, the prefix is dropped. **Put the
  package name back and you have the Java class name** — the mapping in
  `docs/migrating-from-java/` is updated, and `dss/PORTING.md` records the
  rule for future ports.

  Two kinds of name are deliberately left alone. A name identical to its
  package name keeps it. So does a name whose next character is a **digit**,
  which is why the XAdES version helpers are unchanged:
  `xades.XAdES111XSDUtils`, `xades.XAdES122XSDUtils`,
  `xades.XAdES319132XSDUtils`, `xades.Xades141Namespace`, plus
  `evidencerecord.EvidenceRecordsValidationBlock` and seven
  `timestamp.Timestamp*` methods (methods are never written `pkg.Name`, so
  they never stuttered). Renaming those would mangle the version numbers.

  Identifiers whose name merely *contains* a renamed type also keep their
  names, because shortening them would lose meaning:
  `lote.LoLoTEIdentifier` (a list *of* LoTE), the
  `validation.SignatureValidationContext` family, and
  `pades.CMSForPAdESBaselineRequirementsChecker`.

  **Serialized output is unchanged.** Enum values, XML element and attribute
  names, i18n keys, OIDs, MIME types and JSON keys are byte-identical
  (verified by diffing all 53,707 Go string literals). One internal detail
  moved to keep that promise: `model/eaa/claim` rendered
  `AbstractClaim#toString()` from the Go type name via reflection, and now
  prepends the dropped `Claim` prefix explicitly, so claim rendering is
  character-for-character what it was.

  Source-incompatible. It lands before the first tagged release, so no
  deprecated aliases are provided; update call sites by deleting the package
  name from the identifier.

- **Exported constants renamed from Java enum style to Go MixedCaps.** The
  port originally spelled every enum constant `<TypeName>_<JAVA_NAME>`
  (e.g. `enumerations.DigestAlgorithm_SHA256`,
  `i18n.MessageTag_BBB_XCV_CCCBB`); they are now idiomatic Go
  (`enumerations.DigestAlgorithmSHA256`, `i18n.MessageTagBBBXCVCCCBB`).
  2,835 identifiers across 23 packages — every exported constant and
  package-level variable whose name contained an underscore.

  **Serialized values are unchanged.** Only Go identifiers moved: every
  enum *value*, XML element and attribute name, i18n message key, OID,
  URI, MIME type and JSON key is byte-identical (verified by diffing all
  53,687 Go string literals before and after). Reports, signatures and
  trusted-list parsing produce exactly the same output.

  This is a source-incompatible change. It lands before the first tagged
  release, so no deprecated aliases are provided; update call sites
  mechanically. `dss/PORTING.md`'s "Enums" section documents the new
  naming rule, and `docs/migrating-from-java/` shows the Java → Go
  mapping.

  `Test*`/`Example*` function names keep their underscores (idiomatic Go).

### Security

- **`esig`'s RFC 3161 client now checks the token it receives against its
  request.** The message-imprint algorithm and digest must match, and the
  nonce must be echoed (RFC 3161 §2.4.2) — the checks BouncyCastle's
  `TimeStampResponse.validate` performs for upstream's `OnlineTSPSource`.
  Previously a TSA, or anyone on the path, could return a token for other
  data or replay an old one and it was embedded in the signature. A
  password in the `-tsa` URL no longer appears in error messages, and
  `esig tl` downloads are capped at 64 MiB.
- **Hostile input can no longer crash the process or exhaust its memory
  or CPU in the native engines.** Each of the following was a fatal stack
  overflow, an out-of-memory abort or a hang that `recover` cannot catch.
  Each now returns an error.
  - `internal/asn1ber`: nesting is capped at 512 levels. Tag numbers are
    bounded to 31 bits, as in BouncyCastle; a 2^64-1 tag had wrapped to -1
    and been read as a plain certificate. `ValueToString` escaping is now
    linear.
  - `internal/pfx` (PKCS#12): an IV of the wrong size used to panic.
    PBKDF2 key lengths above 32 are refused; they caused a multi-GiB
    allocation. Iteration counts must be within the JDK's
    1–5,000,000. Non-positive DSA domain parameters are refused; `P = 0`
    hung.
  - `internal/pdf`: decoded stream size, including filter chains, is
    bounded by `MaxStreamSize`, and so is the predictor row. The number of
    xref-stream rows is bounded by `MaxObjects`. A huge object-stream `/N`
    no longer panics. Two quadratic paths are now linear: lenient-real
    lexing and the scan for unterminated streams. The incremental writer
    refuses object-number overflow instead of writing negative object
    numbers.
  - `internal/xpath10`, `internal/jose`: expression and JSON nesting are
    capped at 1000.
  - `internal/xmldsig`: a cycle among Manifest references, when manifests
    are followed, returns `ErrManifestCycle`.
  - `crlparser`: `onlySomeReasons` with invalid unused bits is refused, as
    BouncyCastle does.

  Where Java would instead throw `StackOverflowError`, run out of heap or
  accept the input, the divergence is documented in
  `internal/pdf/DESIGN.md` §2.7, `internal/xmldom/DESIGN.md` §6.5 and
  `internal/jose/doc.go`.

- **An encrypted PDF's `/Encrypt /Perms` block is now verified for `/R 5`
  and `/R 6`, so a byte-edited `/P` no longer grants permissions.** For
  AES-256 revisions the permission word `/P` is not mixed into the file
  encryption key, and the native PDF reader never read `/Perms` — the
  16-byte block that carries an authenticated copy of `/P` under that key.
  Anyone able to edit bytes could therefore rewrite `/P` in a signed
  document, flip `CanCreateSignatureField()` from false to true and walk
  through `pades.PdfPermissionsChecker.CheckDocumentPermissions`, the gate
  that refuses to add a signature to a document that forbids it.

  `internal/pdf` now performs ISO 32000-2 Algorithm 13 after the file key
  is recovered, on both the owner and the user password branch: `/Perms`
  is decrypted with the file key and its embedded `/P`, its
  `/EncryptMetadata` marker and its `'a' 'd' 'b'` constant must agree with
  the dictionary. A mismatch — and an absent, non-string or wrong-length
  `/Perms` — is `pdf.ErrInvalidPassword`, which `pades` already maps onto
  upstream's `InvalidPasswordException`; the message says the password was
  accepted, so the two causes are distinguishable.

  Because the check compares `/Perms` against `/P` and `/EncryptMetadata`,
  all three are now read the way `COSDictionary.getDictionaryObject` reads
  them: an indirect reference is followed, and a `/P` written as a real is
  accepted as `COSDictionary.getInt` accepts any `COSNumber`. Previously
  `/P` and `/EncryptMetadata` were read direct-only, so an otherwise
  conforming document that wrote either indirectly was silently given the
  default.

  **This is an intentional divergence from upstream.** Apache PDFBox
  3.0.7's `StandardSecurityHandler.validatePerms` makes the same three
  comparisons but answers every failure with `LOG.warn` and loads the
  document anyway, so upstream Java DSS reports the tampered `/P`. The
  divergence is documented on `validatePerms` in
  `dss/internal/pdf/crypt.go`, in `dss/internal/pdf/DESIGN.md` §2.6 and in
  `dss/PORTING.md`. No document in the corpus is affected: the one
  AES-256 fixture validates cleanly.

- **An encrypted PDF whose cipher cannot be identified is now refused
  instead of silently decrypted with RC4, and `/V 0` is refused too.**
  `internal/pdf`'s Standard security handler set its AES flag from
  `/CF … /CFM` alone, with no `else` and no `default:`, inside `case 4,
  5:` of `setupEncryption`. A `/V 5` document whose `/CF` was missing,
  whose `/CF` had no entry under the selected filter name, or whose
  `/CFM` was `/V2`, `/None` or anything this handler does not implement,
  therefore kept `useAES` false and handed the 32 bytes unwrapped from
  `/UE` to RC4-128 — on both the read path (`decryptBytes`) and the
  write path (`encryptForWrite`, which mirrors it) — while
  `Encryption()` reported `/CFM /None`. A document naming two different
  crypt filters for streams and strings (`/StmF`/`/StrF`) was worse: only
  one of the two was ever resolved, so the other's declared cipher was
  silently overruled.

  `internal/pdf` now validates every crypt filter the document actually
  *selects* (`checkCryptFilters` in `crypt.go`, called from
  `setupEncryption` before the password check): each of `/StmF` and
  `/StrF`, when not `/Identity`, must resolve through `/CF` to an
  implemented `/CFM`; both must resolve to the *same* one; for `/V 5` the
  crypt filter must be `/AESV3`; and `/AESV3` and `/V 5` must each be
  paired with `/R 5` or `/R 6` (ISO 32000-2 defines no other pairing).
  Both halves of that last rule are load-bearing: `/AESV3` promises a
  32-byte key and only the `/UE`/`/OE` unwrap produces one, so
  constraining `/V 5` alone left the identical downgrade one `/V` value
  away — a `/V 4`/`/R 4` document naming `/CFM /AESV3` reported
  `/AESV3` with `/KeyLength 256` and applied AES-128 under an MD5 key. `/StmF` and `/StrF` are now resolved through
  an indirect reference, as `PDEncryption.getStreamFilterName`/
  `getStringFilterName` do, so a `/StmF` or `/StrF` written as `n 0 R`
  is honoured (and checked) rather than silently read as absent and
  falling back to `/Identity`. A setup that still recovers a 32-byte
  `/R 5`/`/R 6` file key without selecting AES — the `/V 4`/`/R 6`
  crossbreed `checkCryptFilters` alone cannot see, since `/CFM /V2` is a
  legitimate `/V 4` declaration on its own — is refused by an explicit
  postcondition just before the handler is published; that same
  postcondition's comment now also notes the one shape that remains and
  is parity, not a bug: `/V 4` with `/R 5`/`/R 6` and `/CFM /AESV2`
  opens reporting `{CFM: AESV2, KeyLength: 128}` while the file key's
  32-byte length still routes it through the AES-256 fast path, exactly
  as pdfbox's own `useAES && encryptionKey.length == 32` branch does.
  Separately, `computeEncryptionKey` now requires `/UE` and `/OE` to be
  exactly the 32 bytes ISO 32000-2 fixes them at — rejected with
  `ErrUnsupportedSecurityHandler`, not `ErrInvalidPassword`, since a
  malformed `/UE`/`/OE` is not a wrong password — closing a route where
  a padded entry's first two AES-CBC blocks still recovered the genuine
  key inside a longer, `len(key) == 32`-defeating slice. `/V 0` — ISO
  32000-1 Table 20's "an algorithm that is undocumented … shall not be
  used", which is also what an `/Encrypt` with no `/V` reads as — is
  refused rather than folded into the same key-length arm as `/V 1`.

  **These are intentional divergences from upstream.** Apache PDFBox
  3.0.7's `StandardSecurityHandler` reads the crypt filter as a
  hardcoded `/StdCF` and sets its AES flag only inside an
  `if (stdCryptFilterDictionary != null)` with no `else`; on three
  tampered copies of the AES-256 corpus fixture it cleared the password
  check, applied RC4, and surfaced only as an unrelated downstream parse
  error, never a security-handler verdict. `/V 0` is accepted by pdfbox
  as `/Length`/8 RC4. pdfbox derives its key and cipher from `/R`
  (`dicRevision`) rather than `/V`, so a `/V 5`/`/R 4` document is exact
  pdfbox parity too — this port refuses it anyway, because it is the
  last surviving "declares AES-256, applies AES-128" shape issue #33
  exists to close. The divergences are documented on
  `checkCryptFilters` and `computeEncryptionKey` in
  `dss/internal/pdf/crypt.go`, in `dss/internal/pdf/DESIGN.md` §2.6 and
  in `dss/PORTING.md`. No document in the corpus is affected: every
  vendored `/V 4`/`/V 5` fixture names `/StdCF` and resolves to
  `/AESV2` or `/AESV3`.

### Documentation

- **Porting-era planning documents removed.** `PORTING_PLAN.md` — the
  phase-by-phase roadmap and progress tracker for the Java→Go port — has
  been deleted. It described *how the port was built*, not how the library
  behaves, and every durable claim it carried now lives where a reader will
  look for it: the upstream baseline pin in `UPSTREAM.md` and `NOTICE`, the
  Maven-module → Go-package mapping in the "Migrating from Java DSS"
  documentation page, and the accepted gaps in the "Known gaps"
  compatibility page. All references to it — in `README.md`, `CHANGELOG.md`,
  `CONTRIBUTING.md`, `UPSTREAM.md`, `NOTICE`, the compatibility
  documentation and 28 Go source files — were repointed at those
  successors.
- **Porting narration removed from durable documents.** Phase labels,
  implementer assignments and session-diary notes were stripped from
  `dss/internal/pdf/DESIGN.md`, `dss/internal/xmldom/DESIGN.md`,
  `dss/harness/README.md` and eight test-fixture provenance READMEs. Their
  technical content — the design decisions, the measured corpus evidence,
  the pinned APIs, the oracle provenance, the regeneration recipes and the
  recorded divergences — is unchanged.
- **Attribution and licensing documents are unchanged.** `LICENSE`, `NOTICE`
  (pointer update only), `CONTRIBUTING.md`'s DCO and LGPL-2.1
  inbound = outbound terms, `SECURITY.md`, `CODE_OF_CONDUCT.md` and every
  `// Ported from` header remain exactly as they were. `dss/PORTING.md` is
  kept: it is the normative conventions document for contributed ports, not
  a porting-progress record.

### Known limitations

Tracked, accepted gaps in this initial port (the documentation site's
"Known gaps" page is the authoritative, currently-maintained list;
summarized here as of this port):

- **No online time-stamping or revocation clients in the library.**
  Upstream's `dss-service` module (`OnlineTSPSource`, `OnlineCRLSource`,
  `OnlineOCSPSource`) is out of scope, so levels T, LT and LTA require a
  `TSPSource` — and LT/LTA additionally revocation sources — that you
  supply. What *is* included: the native HTTP data loader used for AIA
  certificate fetching and for downloading trusted lists, an offline
  `KeyEntityTSPSource` issuing RFC 3161 tokens from a local key, and an
  RFC 3161-over-HTTP client inside the `esig` CLI.
- **No evidence-record analyzer.** The evidence-record framework,
  interfaces, report plumbing and ASiC container hooks are ported, but
  nothing registers an `EvidenceRecordAnalyzerFactory`, so an RFC 4998
  ERS or RFC 6283 XMLERS record cannot actually be parsed from a
  document. (Listed again under "Deferred" for the upstream modules that
  were not ported at all.)
- **JKS, PKCS#11 and the OS key stores are not supported.** JKS is a
  proprietary Java binary format with no Go implementation available
  (stdlib or `golang.org/x/...`); PKCS#11 hardware-token access has no
  Go counterpart to the JCA `SunPKCS11` provider this port could wrap;
  the Windows certificate store (`Windows-MY`) and the macOS Keychain
  are likewise JCA providers with no Go counterpart. All of them report
  a clear "not supported in the Go port" error rather than failing
  silently. A hardware token can still be driven by implementing
  `token.SignatureTokenConnection`.
- **MD2 and WHIRLPOOL digest algorithms** are recognized as enum values
  (for parsing/round-tripping data that references them) but are not
  implemented for actual digesting, matching their status as
  deprecated/rarely-available algorithms.
- **Streaming CMS is in-memory-only**: this port does not yet stream
  CMS content the way upstream's `InputStream`-based API can for very
  large payloads.
- **Visible signature rasterization is not implemented**: PAdES visible
  signature appearance drawing (`pdf`/visible-signature drawers) is not
  supported; non-visible signing and validation are unaffected.
- **A handful of narrow, documented parity gaps** surfaced by broad
  differential corpus runs against upstream (booked, tracked by name,
  and confirmed not to be false-accepts): a small number of CMS/CAdES
  edge cases found by a 248-PDF PAdES differential run (see
  `dss/pades/testdata/broadgen/README.md`), and a deterministic-ordering
  difference in one algorithm-obsolescence (AOV) grouping step where
  this port iterates first-seen order instead of replicating Java
  `HashMap` bucket order (does not change any validation conclusion).
- **`process/eaa`** (electronic attestation of attributes validation)
  is ported and mechanically verified 1:1 against upstream but does not
  yet have a behavioral oracle corpus; it builds only under the `eaa`
  build tag pending that follow-up work.

### Deferred (out of initial scope)

Not ported in this initial release — see the "Not ported at all" table
on the documentation site's "Known gaps" page: upstream's REST/SOAP
remote services and clients, `dss-cookbook`, coverage/BOM modules, and
evidence-record modules.

[Unreleased]: https://github.com/ryftcore/dss-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/ryftcore/dss-go/releases/tag/v0.1.0
