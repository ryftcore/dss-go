# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project intends to adopt [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
once tagged releases begin (see `SECURITY.md`'s supported-versions note —
there is no tagged release yet).

## [Unreleased]

Initial public port of [esig/dss](https://github.com/esig/dss) (upstream
baseline: version 6.5.RC1, commit
`4c2129862948bfd53ca1455832260aa17e183cf8`) to Go, module
`github.com/utain/esig/dss`. See `PORTING_PLAN.md` for the full
phase-by-phase porting record, module mapping, and verification detail
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

  Per-phase scope, fixture counts and results are in `PORTING_PLAN.md`;
  what runs in `go test` versus what was a one-time sweep is separated
  in the repository's compatibility documentation. The formats:
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
  `dss/harness` — see `PORTING_PLAN.md`'s progress table for the scope
  and scale verified in each phase (fixture counts, corpus sizes, and
  parity results are recorded there per phase; this file does not
  restate per-phase numbers to avoid the two drifting apart).
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

### Known limitations

Tracked, accepted gaps in this initial port (see `PORTING_PLAN.md`'s
"Known accepted gaps" for the authoritative, currently-maintained list;
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

Not ported in this initial release — see `PORTING_PLAN.md`'s module
mapping for the full list: upstream's REST/SOAP remote services and
clients, `dss-cookbook`, coverage/BOM modules, and evidence-record
modules.

[Unreleased]: https://github.com/utain/esig/compare/main...HEAD
