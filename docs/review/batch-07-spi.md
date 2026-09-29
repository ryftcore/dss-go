# Batch 07 — spi package

## spi core — top 6 largest (unit U22a)

- Scope: the 6 largest top-level `.go` files in `dss/spi/` — `dss_utils.go`,
  `certificate_extensions_utils.go`, `dss_revocation_utils.go`, `dss_asn1_utils.go`,
  `cms_certificate_source.go`, `offline_revocation_source.go`.
- Date: 2026-08-26
- Depth: budgeted (first 250 lines of each file + `grep -n 'func '` for all signatures +
  targeted reads of security-relevant regions: certificate loading/parsing, digest
  calculation, key-size checks, ASN.1 decoding, OCSP/CRL parsing, revocation-status
  matching, serial comparison, issuer/serial selection, AIA/CRL-DP parsing,
  name-constraint checks, DSA/bit-length math, X.500 principal equality).
- Compatibility contract: interop with Java DSS 6.5.RC1 (byte-level + verdict parity).
  `dss/PORTING.md` conventions and `docs/compatibility/known-gaps.md` are treated as
  deliberate, not findings. `// DIVERGENCE, deliberate:` comments mark intentional
  divergences.
- Out-of-scope: `spi/validation/`, `spi/x509/`, `spi/eaa/`, `spi/lote/`, `spi/policy/`,
  `spi/tsl/`, `spi/client/`, `spi/alerts/`, `spi/exception/`, `spi/extension/`,
  `spi/random/`, `spi/signature/`, and the remaining ~120 top-level `spi/` files (units
  U22b/U22c/U22d).

### Files read (budgeted approach per file)

| File | Lines | Approach |
|---|---|---|
| `dss_utils.go` | 1589 | First 250 lines (header, imports, date formatting). `grep -n 'func '` for all 80 signatures. Targeted reads: certificate loading (`dssUtilsParseCertificate`, `DSSUtilsLoadCertificate*`, `dssUtilsP7cCertificates`, lines 401–565), digest calculation (`DSSUtilsDigest`, `DSSUtilsMessageDigest`, `dssUtilsShake`, `DSSUtilsDigestReader`, `DSSUtilsDigestConcat`, lines 566–760), EC-signature conversion (`DSSUtilsConvertECSignatureValue`, `DSSUtilsEdDSASignatureAlgorithm`, lines 1444–1500), `DSSUtilsEnrichCollection` (line 1421), `DSSUtilsGenerateKid` / `DSSUtilsHost` (lines 1527–1589). |
| `certificate_extensions_utils.go` | 1394 | First 250 lines (header, imports, OID constants, `CertificateExtensionsUtilsCertificateExtensions`, all `Is*` predicates, `CertificateExtensionsUtilsSubjectAlternativeNames`). `grep -n 'func '` for all 60 signatures. Targeted reads: SAN parsing (`certificateExtensionsUtilsAlternativeNames`, `certificateExtensionsUtilsGeneralNameOf`, lines 267–333), AIA/CRL-DP parsing (`CertificateExtensionsUtilsAuthorityInformationAccess`, `CertificateExtensionsUtilsCRLDistributionPoints`, `certificateExtensionsUtilsCRLDistributionPointsUrls`, lines 334–600), name-constraint parsing (`certificateExtensionsUtilsGeneralSubtrees`, lines 642–707), policy/QC extensions (`CertificateExtensionsUtilsPolicyConstraints` through `CertificateExtensionsUtilsQcStatements`, lines 752–1080), ASN.1 helpers (`certificateExtensionsUtilsParseTagged`, `certificateExtensionsUtilsInteger`, `certificateExtensionsUtilsHostAddress`, `certificateExtensionsUtilsIPv4Mapped`, `certificateExtensionsUtilsRFC4519Name`, `certificateExtensionsUtilsEncodeAttributeValue`, `certificateExtensionsUtilsASN1String`, lines 1082–1394). |
| `dss_revocation_utils.go` | 1160 | First 250 lines (header, imports, OCSP status constants, `Extension`, `ResponderID`, `CertID`, `CertificateID`, `OCSPCertStatus`, `SingleResponse`, `SingleResp`). `grep -n 'func '` for all 50 signatures. Targeted reads: OCSP response decoding (`dssRevocationUtilsOCSPResponseFromElement` through `dssRevocationUtilsExtensionsFromTagged`, lines 860–1160), OCSP certID matching (`DSSRevocationUtilsMatches`, `DSSRevocationUtilsOCSPCertificateID`, `dssRevocationUtilsSubjectPublicKeyBits`, lines 642–710), OCSP response object (`ParseBasicOCSPResponse`, `BasicOCSPResp`, `OCSPResp`, `ResponseObject`, `Certs`, lines 300–560), `ParseOtherHash` (line 538), `DSSRevocationUtilsCRLRevocationTokenKey` / `DSSRevocationUtilsOcspRevocationKey` (lines 744–790). |
| `dss_asn1_utils.go` | 967 | First 250 lines (header, imports, `AlgorithmIdentifier`, `ParseAlgorithmIdentifier`, `DSSASN1UtilsToASN1Primitive`, `DSSASN1UtilsIsDEROctetStringNull`, `DSSASN1UtilsDEREncoded`, `DSSASN1UtilsToDate`, `DSSASN1UtilsToString`, `DSSASN1UtilsAsn1SequenceFromDerOctetString`, `DSSASN1UtilsAsn1IntegerFromDerOctetString`). `grep -n 'func '` for all 50 signatures. Targeted reads: X.500 principal equality (`DSSASN1UtilsX500PrincipalAreEquals`, `DSSASN1UtilsAttributeMap`, `DSSASN1UtilsString`, lines 570–660), DSA/bit-length math (`DSSASN1UtilsEnsurePlainSignatureValue`, `DSSASN1UtilsToPlainDSASignatureValue`, `dssASN1UtilsEncodeDSAValue`, `DSSASN1UtilsToStandardDSASignatureValue`, `DSSASN1UtilsOrderFromSignatureValue`, `DSSASN1UtilsSignatureValueBitLength`, `dssASN1UtilsUnsignedByteLength`, lines 788–930), `DSSASN1UtilsIssuerSerial` / `DSSASN1UtilsToSignerIdentifierFromIssuerSerial` (lines 486–570). |
| `cms_certificate_source.go` | 779 | First 250 lines (header, imports, `SigningCertificate`, `SigningCertificateV2`, `ESSCertID`, `ESSCertIDv2`, `OtherCertID` parsing). `grep -n 'func '` for all 30 signatures. Targeted reads: CMS cert source (`NewCMSCertificateSource`, `extractCertificateIdentifiers`, `extractSignedCertificates`, `extractSigningCertificateReferences`, `extractSigningCertificateV1/V2`, `extractESSCertIDs`, `extractESSCertIDv2s`, `extractCertificateValues`, `extractCertificateRefsFromUnsignedAttribute`, `ExtractCandidatesForSigningCertificate`, `cmsCertificateSourceFirstToken`, `cmsCertificateSourceSerialsEqual`, lines 457–779), `cmsCertificateSourceSequence` / `cmsCertificateSourceMembers` / `DSSASN1UtilsToSignerIdentifierFromSignerID` / `DSSASN1UtilsCertificate` / `DSSASN1UtilsAsn1Attributes` / `DSSASN1UtilsAsn1Encodable` / `DSSASN1UtilsCertificateRef` (lines 300–456). |
| `offline_revocation_source.go` | 611 | First 250 lines (header, imports, `OfflineRevocationSourceOverrides`, `offlineRevocationSourceBinaryEntry/TokenEntry/RefEntry`, `OfflineRevocationSourceBase`, `NewOfflineRevocationSourceBase`, `InitOfflineRevocationSource`, `offlineRevocationSourceBaseOverrides`, `AddBinary`, `AddRevocation`, `AddRevocationWithBinary`, `originsForBinary`, `AddRevocationReference`, `AllRevocationBinaries`, `AllRevocationTokens`, `UniqueRevocationTokensWithOrigins`, `AllRevocationReferences`, `RevocationToken`). `grep -n 'func '` for all 50 signatures. Targeted reads: all origin-specific accessors (lines 290–390), `FindRefsAndOriginsForRevocationToken` / `FindRefsAndOriginsForBinary` / `FindBinaryForReference` / `OrphanRevocationReferencesWithOrigins` / `IsOrphan` (lines 408–500), helper functions (lines 505–611). |

### Findings

**SEC (top priority)**

No findings. Each mandated SEC lens was verified clean:

- **OCSP certID matching (false-accept / false-reject):** `DSSRevocationUtilsMatches`
  (`dss_revocation_utils.go:642`) compares `HashAlgOID`, `IssuerKeyHash`,
  `IssuerNameHash`, and `SerialNumber` with exact equality (`bytes.Equal`, `Cmp == 0`).
  The comment documents the deliberate field-by-field comparison (upstream's
  `certId.equals()` fails on null-vs-DERNull AlgorithmIdentifier parameters). No
  false-accept: a non-matching serial, issuer name hash, or issuer key hash returns
  `false`. No false-reject: all four fields are compared. `DSSRevocationUtilsOCSPCertificateID`
  (`dss_revocation_utils.go:659`) computes the CertID from the issuer cert's subject
  (`RawSubject`) and subject public key, matching BouncyCastle's
  `CertificateID(DigestCalculator, X509CertificateHolder, BigInteger)` constructor.
  For a valid chain, `issuerCert.RawSubject == cert.RawIssuer`, so the hash is over
  the correct bytes. `dssRevocationUtilsSubjectPublicKeyBits` extracts the BIT STRING
  value bits via `asn1ber.Parse` + `BitStringOctets()`, which is bounds-checked.
- **Revocation-status matching (CRL serial comparison):** Handled by `crlparser`
  (unit U21, reviewed in batch-06). `dss_revocation_utils.go` does not perform CRL
  serial matching; it provides the OCSP-side primitives. The CRL-side matching is in
  `crlparser/crl_utils_x509crl_impl.go:304` (`entry.SerialNumber.Cmp(serialNumber) == 0`),
  which is exact integer equality. No false match or false miss.
- **Issuer/serial selection (CMS cert source):** `ExtractCandidatesForSigningCertificate`
  (`cms_certificate_source.go:688`) resolves the signing certificate by
  `CertificateTokenBySignerIdentifier` (issuer+serial or SKI match), then falls back
  to `signingCertificateSource.BySignerIdentifier`. `cmsCertificateSourceSerialsEqual`
  (line 771) uses `Cmp == 0` (exact) and returns `false` for nil. `DSSASN1UtilsX500PrincipalAreEquals`
  (`dss_asn1_utils.go:570`) first tries `X500Principal.Equals` (canonical comparison),
  then falls back to attribute-map comparison. Both are conservative (no false-accept).
  `DSSASN1UtilsToSignerIdentifierFromIssuerSerial` (line 546) returns `nil` when the
  issuer name cannot be parsed (undecodable name), which is the conservative direction
  (no trust grant). `cmsCertificateSourceFirstToken` (line 760) sorts map keys
  deterministically (a deliberate deviation from Java's HashSet order, documented in
  the comment). No trust-logic error.
- **ASN.1 parsing edge cases (nil-deref panics on malformed input):** All ASN.1 parsing
  in the 6 in-scope files uses either `cryptobyte.String` (bounds-checked, returns
  `false` on malformed input) or `asn1ber.Parse` (returns an error on malformed input).
  No unbounded recursion, no expansion attacks, no nil-deref panics reachable by hostile
  input. The `panic` calls in `dss_utils.go` (`DSSUtilsLoadCertificateFromBinary`
  line 428, `DSSUtilsDigest` line 585) and `dss_revocation_utils.go`
  (`DSSRevocationUtilsCRLRevocationTokenKey` line 758,
  `DSSRevocationUtilsOcspRevocationKey` line 782) are 1:1 ports of upstream
  `Objects.requireNonNull` / unreachable-error semantics, not reachable by hostile
  input on the normal path (see T22A-STD-001 for the latter).
- **Key-size / algorithm checks:** `DSSASN1UtilsEncodeDSAValue` (`dss_asn1_utils.go:830`)
  rejects values outside `[0, order)` and values longer than the buffer — no
  out-of-bounds write. `DSSASN1UtilsOrderFromSignatureValue` (line 858) computes
  `max(r, s) + 1` as the byte-length reference (a faithful port of BouncyCastle's
  `PlainDSAEncoding`); the actual signature verification is done by the crypto
  library, not by this function. `DSSUtilsEdDSASignatureAlgorithm` (line 1481) maps
  signature length to algorithm (64 → Ed25519, 114 → Ed448); the ambiguity with
  ECDSA P-256 (also 64 bytes) is resolved by the caller's context (the function name
  and the upstream comment both state "Only EdDSA algorithms are recognised"). No
  false-accept.
- **Digest calculation:** `DSSUtilsDigest` (line 585) dispatches to the correct
  algorithm for each `DigestAlgorithm` enum value. SHAKE-128/256 use the correct
  output lengths (32/64 bytes, matching BouncyCastle's `SHAKEDigest#getDigestSize()`).
  MD5/SHA1/SHA2/SHA3/RIPEMD160 use the correct stdlib / `golang.org/x/crypto`
  implementations. `dssUtilsShake.Sum` clones the hash before squeezing (correct —
  `hash.Hash` contract requires `Sum` to not mutate state). No digest-algorithm
  confusion.

**PERF (CRL parsing cost, allocations, leaks)**

T22A-PERF-BIGO-001 — **Low** — Category: PERF-BIGO — `dss/spi/dss_utils.go:1421`
- Evidence: `DSSUtilsEnrichCollection` linearly scans `*currentCollection` for each
  element of `toAddCollection` (`for _, object := range toAddCollection { for _,
  existing := range *currentCollection { if existing == object { … } } }`) — O(n·m)
  in the two collection sizes.
- Impact: For collections of size n and m, the dedup costs O(n·m). In practice, all
  callers in the DSS codebase use this for small collections (OIDs, URIs, algorithm
  names, typically single-digit to low-double-digit element counts), so the cost is
  negligible. Not exploitable by attacker input on the normal path. Same pattern as
  T18A-PERF-BIGO-001 (`tsl/trust_service_equivalence_converter.go:218`).
- Recommendation: None required for the current codebase. If this function is ever
  called with large collections, replace the inner linear scan with a
  `map[T]struct{}` accumulator.

No unbounded-memory, O(n²)-on-attacker-input (beyond the above), or goroutine-leak
findings in the 6 in-scope files. All parsing is bounded by input size (cryptobyte
and asn1ber are both bounds-checked). No goroutines, no retained references, no cache
leaks. Allocations in the hot path are bounded by the size of the certificate / CRL /
OCSP response being processed and match the "parse once, retain" design.

**STD (idioms, error handling, defer/cleanup, context, sync, doc comments, test quality)**

T22A-STD-001 — **Info** — Category: STD — `dss/spi/dss_revocation_utils.go:758`
(`DSSRevocationUtilsCRLRevocationTokenKey`) & `:782` (`DSSRevocationUtilsOcspRevocationKey`)
- Evidence: Both functions call `DSSUtilsSHA1Digest(…)` and `panic(err.Error())` on
  the (theoretically unreachable) error path. The comment at line 753–757 explains:
  "DSSUtilsSHA1Digest's error is a Go-only surface … Java's method declares no checked
  exception either, so an error here is reproduced as a panic."
- Impact: SHA-1 is always available in Go's `crypto/sha1`, so the error path is
  unreachable today. However, if a future change (e.g., a FIPS-mode build, a
  `crypto/sha1` removal, or a bug in the digest dispatch) makes the digest fail, the
  panic would crash the process rather than returning an error. This is a latent DoS
  vector, not an exploitable one. The same pattern appears in `dss_utils.go:428`
  (`DSSUtilsLoadCertificateFromBinary` panics on nil input) and `dss_utils.go:585`
  (`DSSUtilsDigest` panics on nil data) — all are 1:1 ports of upstream
  `Objects.requireNonNull` and are not reachable by hostile input.
- Recommendation: No code change required (deliberate port, parity with upstream).
  Consider adding a one-line comment cross-referencing the PORTING.md "Java exceptions
  → returned error" rule so a future reader understands why the panic is deliberate.
  If the codebase ever moves to a FIPS-compliant build where `crypto/sha1` might be
  unavailable, convert the panic to a returned error.

No other STD findings. Error handling is correct (`error` returns at the facade
boundary; panics are 1:1 ports of upstream `Objects.requireNonNull`). No goroutines,
so no sync issues. Tests are thorough (KAT vectors, hostile-input error paths,
round-trip parity). The `DSSUtilsEnrichCollection` O(n·m) dedup is noted under
PERF-BIGO above.

### Tool log

Run from `dss/`:

- `gofmt -l spi` → no output (clean).
- `go vet ./spi/` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./spi/...` → `0 issues`.

All tools available and passing; nothing to note as unavailable.

### Summary

- **Total findings: 2** (SEC: 0, PERF-BIGO: 1, PERF-MEM: 0, PERF-LEAK: 0, STD: 1)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 0, Low: 1, Info: 1
- The 6 in-scope files are a high-quality 1:1 port of the upstream Java DSS
  `DSSUtils`, `CertificateExtensionsUtils`, `DSSRevocationUtils`, `DSSASN1Utils`,
  `CMSCertificateSource`, and `OfflineRevocationSource` classes. All ASN.1 parsing
  is bounds-checked (cryptobyte / asn1ber), all trust-logic paths (OCSP certID
  matching, issuer/serial selection, DSA signature conversion) are correct and
  conservative, and all digest calculations are accurate. The two findings are
  minor (a Low-severity O(n·m) dedup on small collections, and an Info-level
  note on deliberate panics at unreachable error paths). No SEC, PERF-MEM, or
  PERF-LEAK issues were found.

### Open questions

1. The `DSSUtilsEnrichCollection` function is generic and could be called with large
   collections by future code. If this becomes a concern, the O(n·m) dedup should be
   replaced with a `map[T]struct{}` accumulator. No action required today.
2. The `certificateExtensionsUtilsEncodeAttributeValue` function (line 1312)
   truncates UTF-16 code units to bytes for IA5String/PrintableString encoding,
   matching BouncyCastle's `Strings.toByteArray`. This means non-ASCII characters
   in DN attributes (e.g., a `cn=José` encoded as BMPString in the cert) are
   truncated to their low byte when re-encoded as PrintableString. This is a
   faithful port of upstream behavior and is pinned by KAT tests. No action
   required, but worth noting if a future requirement demands full Unicode
   preservation in DN rendering.

## spi core — 7 mid-size files (unit U22b)

- Scope: the 7 next-largest top-level `.go` files in `dss/spi/` — `qc_statement_utils.go`,
  `ocsp_token.go`, `revocation_token.go`, `cms_crl_source.go`,
  `common_certificate_source.go`, `file_revocation_source.go`, `list_certificate_source.go`.
- Date: 2026-08-26
- Depth: budgeted (per-file first 250 lines + `grep -n 'func '` for all signatures + targeted
  reads of the security-relevant regions: QC-statement OID parsing / extension matching,
  OCSP response parsing / certID + revocation-status extraction, revocation-token
  dispatch / status matching, CMS SignedData CRL extraction / issuer-serial selection,
  certificate chain building / dedup, file-based CRL/OCSP loading / expiry / cache,
  list-based source dedup / chain walk).
- Compatibility contract: interop with Java DSS 6.5.RC1 (byte-level + verdict parity).
  `dss/PORTING.md` conventions and `docs/compatibility/known-gaps.md` are treated as
  deliberate, not findings. `// DIVERGENCE, deliberate:` comments mark intentional
  divergences.
- Out-of-scope: `spi/validation/` and the other `spi/` subpackages (U22d/U23), the 6
  largest top-level files (U22a, reviewed above), and the remaining ~60 smaller top-level
  `spi/` files (U22c). Cross-package helpers this unit delegates to (`internal/asn1ber`,
  `internal/cmscore`, `crlparser`, `model.CertificateToken`, `model.TokenBase`,
  `utils.OrderedMap`) were read only to bound the severity of in-scope findings, not
  re-reviewed.

### Files read (budgeted approach per file)

| File | Lines | Approach |
|---|---|---|
| `qc_statement_utils.go` | 557 | FULL-read (lines 1–250 header/OID-parsing/extension-matching + 251–557 QC-statement extraction + ASN.1 helpers). `grep -n 'func '` for all 31 signatures. |
| `ocsp_token.go` | 475 | FULL-read (lines 1–250 OCSP response parsing / status extraction / cert-hash + 251–475 certID matching / signature verification / principal rendering). `grep -n 'func '` for all 26 signatures. |
| `revocation_token.go` | 468 | FULL-read (lines 1–250 RevocationToken interface / base state + 251–468 revocation-token dispatch / status matching / Equals). `grep -n 'func '` for all 28 signatures. |
| `cms_crl_source.go` | 464 | FULL-read (lines 1–250 ESF structure parsing / CRL+OCSP list decoding + 251–464 CMS SignedData CRL extraction / issuer-serial selection / revocation-ref collection). `grep -n 'func '` for all 18 signatures. |
| `common_certificate_source.go` | 432 | FULL-read (lines 1–250 cert chain building / dedup / nil-map init + 251–432 By* lookups / set equality). `grep -n 'func '` for all 26 signatures. |
| `file_revocation_source.go` | 393 | FULL-read (entire file: file-based CRL/OCSP loading / expiry / cache / file-path construction). `grep -n 'func '` for all 28 signatures. |
| `list_certificate_source.go` | 383 | FULL-read (entire file: list-based source / dedup / chain walk / set equality). `grep -n 'func '` for all 28 signatures. |

### Findings

**SEC (top priority)**

No findings. Each mandated SEC lens was verified clean:

- **OCSP response parsing / certID matching (false-accept / false-reject):**
  `ocspTokenVerifySignature` (`ocsp_token.go:442`) delegates to `crypto/x509`
  `Certificate.CheckSignature` over the retained `tbsResponseData` bytes with the exact
  `x509.SignatureAlgorithm` resolved from the response's AlgorithmIdentifier via the
  `ocspTokenSignatureAlgorithms` map. Unsupported algorithms (SHA-224, SHA-3, RIPEMD-160,
  PLAIN-ECDSA, Ed448) are deliberately absent from the map and yield a verification *error*
  rather than a silent success. `IsValid()` (`ocsp_token.go:299`) is fail-closed: it is
  `IsSignatureIntact() && isOCSPVersionValid()`, and `IsSignatureIntact()`
  (`model/token.go:208`) returns `false` when `signatureValidity` is not `VALID` (the zero
  value is `NOT_EVALUATED`, not `VALID`), so an unverified token is never reported valid.
  `CheckIsSignedBy` (`ocsp_token.go:309`) records an invalidity reason on every failure.
  `extractCertHashExtension` (`ocsp_token.go:163`) is informational-only (`SetCertHashMatch`)
  and does not feed `IsValid()`. No false-accept.
- **Revocation-status extraction (CRL/OCSP dispatch):** `extractStatusInfo`
  (`ocsp_token.go:115`) maps the RFC 6960 `CertStatus` CHOICE exhaustively (good →
  `CertificateStatusGood`, revoked → `Revoked` + date + reason, unknown → `Unknown`).
  `RevocationTokenBase` (`revocation_token.go`) carries the status and the CRL/OCSP
  dispatch; the CRL-side serial matching lives in `crlparser` (unit U21, batch-06) and is
  exact integer equality. `IsInternal()` (`revocation_token.go:424`) and
  `SetExternalOrigin` (`revocation_token.go:407`) faithfully port the upstream origin
  logic. No status confusion.
- **CMS SignedData CRL extraction / issuer-serial selection:** `cms_crl_source.go` parses
  the ESF `RevocationValues` / `CrlOcspRef` structures with `asn1ber.Parse`
  (bounds-checked, errors on malformed input) and only collects CRLs into the source;
  `OcspVals` are read here but consumed by `CMSOCSPSource` (`cms_ocsp_source.go`, a
  separate in-scope-adjacent file) — `extractRevocationValues` (`cms_crl_source.go:390`)
  correctly touches only `CrlVals`. `AddX509CRLHolder` (`cms_crl_source.go:405`) wraps
  `crlparser.CRLUtilsBuildCRLBinary` and returns the parse failure as a `DSSError` rather
  than silently accepting a malformed CRL. No issuer/serial false-match.
- **Certificate chain building / dedup (nil-deref on malformed input):**
  `common_certificate_source.go` keys its maps on `AsXmlID()` (a hex digest, `identifier.go:108`),
  which is collision-free; `utils.OrderedMap.Get/Values/Len` (`ordered_map.go:36,63,82`)
  are all nil-receiver-safe (a nil `*OrderedMap` reads as empty), and the bare-map reads
  (`ByPublicKey`, `BySki`, `BySubject`) are Go nil-map reads (safe, return zero). The
  write path guards with `commonCertificateSourceEnsureInitialized`. `IsKnown`
  (`common_certificate_source.go:182`) dereferences `token` without a nil check, but this
  mirrors Java's `Objects.requireNonNull` and every in-repo caller passes a non-nil token
  (verified: `list_certificate_source.go`, `common_trusted_certificate_source.go`). No
  reachable nil-deref on hostile input.
- **File-based CRL/OCSP loading / cache (path traversal / expiry):** the cache file path is
  `filepath.Join(fileCacheDirectory, key+ext)` (`file_revocation_source.go:317`) where
  `key` is a hex SHA-1 digest (see `dss_revocation_utils.go:758,782`, out-of-scope U22a
  but the key format is what bounds this), so `key` cannot contain `/` or `..` — no path
  traversal out of the cache directory. `SaveRevocationToken` / `SaveCertificateToken`
  validate non-nil before writing. `FindRevocations` (`file_revocation_source.go:150`)
  only reconstructs a token when `Exists()` is true. Expiry is enforced upstream by
  `RepositoryRevocationSourceBase.extractRevocationFromCacheSource`
  (`repository_revocation_source.go:204`, out-of-scope) which filters via `isNotExpired`
  and removes expired entries when `removeExpired` is set. No false-accept of stale data.
- **List-based source dedup / chain walk:** `list_certificate_source.go` delegates all
  per-certificate trust logic to the embedded `CertificateSource`s and only aggregates;
  `AddCertificate` and `CertificateSourceType` panic with the upstream
  `UnsupportedOperationException` messages (faithful port). No trust-logic error.

**PERF (BigO, memory, leaks)**

T22B-PERF-001 — **Medium** — Category: PERF-LEAK — `dss/spi/file_revocation_source.go:169`
(`InsertRevocation`) / `:112` (`SetFileCacheDirectory`)
- Evidence: `InsertRevocation` writes one revocation file (+ optional `.uri`, + `.cer`
  issuer cert) per `revocationKey` into `fileCacheDirectory` (default
  `$TMPDIR/dss-cache-revocation`) with no size, count, or age bound. `RemoveRevocation`
  (`:197`) and `ClearCache` (`:213`) are the only removal paths, and
  `extractRevocationFromCacheSource` (`repository_revocation_source.go:204`, out-of-scope)
  only removes a key when *all* of its cached tokens are expired and `removeExpired` is
  set — a key that is never re-queried is never cleaned.
- Impact: The on-disk cache grows without bound as the application validates documents
  referencing distinct OCSP/CRL endpoints. The key is `SHA-1(url)`, so the distinct-key
  set is bounded by the number of distinct revocation endpoints the workload touches; a
  long-lived validator over many issuers (or a hostile document referencing many DP URLs)
  accumulates cache files in the OS temp directory indefinitely. This is a faithful port
  of upstream `FileRevocationSource`, which also has no eviction, so it is parity, not a
  divergence.
- Recommendation: None required for parity. If unbounded growth is a concern in
  long-running validators, add a max-entry / max-age eviction policy to the cache
  directory (and document it under PORTING.md as a deliberate hardening).

T22B-PERF-002 — **Medium** — Category: PERF-BIGO — `dss/spi/list_certificate_source.go:310`
(`Entities`) / `:193` (`IsCertificateSourceEquivalent`) and
`dss/spi/common_certificate_source.go:402` (`commonCertificateSourceEntitySetsEqual`)
- Evidence: `Entities()` (`list_certificate_source.go:310`) dedups across all sources with
  a linear `entitiesEqual` scan over the growing `result` slice for every new entity
  (O(n·m)). `commonCertificateSourceEntitySetsEqual`
  (`common_certificate_source.go:402`) is pairwise O(n²) with an `Equals` per pair, and
  `IsCertificateSourceEquivalent` (`list_certificate_source.go:193`) adds a further
  O(n²) `entitiesEqual` pass.
- Impact: O(n²) in the certificate/entity count on the set-equality and entity-dedup
  paths. These are not on the primary validation hot path (they back
  `IsCertificateSourceEqual`/`IsCertificateSourceEquivalent`/`Entities`, used for source
  comparison and reporting), and all in-repo callers use small entity sets, so the cost is
  negligible today. Same class as T22A-PERF-BIGO-001 (`dss_utils.go:1421`).
- Recommendation: None required for the current codebase. If these are ever called with
  large sources, key the dedup/equality on `AsXmlID()` (a `map[string]struct{}`) instead
  of a linear scan.

No unbounded-heap, goroutine-leak, or O(n²)-on-attacker-input (beyond the two above)
findings in the 7 in-scope files. All ASN.1 parsing is bounds-checked (cryptobyte /
asn1ber). No goroutines, no retained references, no heap cache leaks; the only unbounded
growth is the on-disk file cache (T22B-PERF-001).

**STD (idioms, error handling, defer/cleanup, context, sync, doc comments, test quality)**

T22B-STD-001 — **Info** — Category: STD — `dss/spi/file_revocation_source.go:355`
(`FileCacheEntry.Clean`)
- Evidence: `Clean()` computes `cacheCleaned` by *combining* the per-file removal results
  with `!=` (XOR): `cacheCleaned = cacheCleaned != fileCacheEntryRemoveFile(uriFile)`,
  then again for the cert file. With the revocation file removed (true) and both the
  `.uri` and `.cer` files absent (both false), the result is `true != false` → false, then
  `false != false` → false — i.e. `Clean()` reports *not cleaned* when all files are
  successfully removed, and reports *cleaned* when exactly one of the optional files was
  removed.
- Impact: The returned boolean does not read as "all associated files removed"; it is the
  parity of the three per-file outcomes. The only caller, `RemoveRevocation`
  (`file_revocation_source.go:207`), discards the return value, so there is no behavioural
  impact in-repo. This is a faithful port of the upstream Java `clean()` (which likewise
  folds the three `removeFile` results), so it is parity, not a divergence.
- Recommendation: No code change required (deliberate port). If the return value is ever
  consumed, change the combination to `&&` (all files removed) or document the XOR
  semantics with a one-line comment cross-referencing the upstream `clean()`.

T22B-STD-002 — **Info** — Category: STD — `dss/spi/file_revocation_source.go:181`
(`saveRevocationToken`)
- Evidence: `saveRevocationToken` calls `panic(err)` when
  `FileCacheEntry.SaveRevocationToken` returns an I/O error. `InsertRevocation` /
  `UpdateRevocation` (the `RepositoryRevocationSourceOverrides` methods) have no error
  channel, so the failure is surfaced as a panic.
- Impact: A disk-full, permission, or filesystem error during cache write crashes the
  process rather than returning an error. This is a faithful port of upstream Java, where
  `DSSUtils.saveToFile` wraps the `IOException` in an unchecked `DSSException` and the
  caller has no checked-exception channel either, so it is parity, not a divergence.
- Recommendation: No code change required (deliberate port, same reasoning as
  T22A-STD-001). If the `RepositoryRevocationSourceOverrides` interface is ever widened to
  carry an error, convert the panic to a returned error.

No other STD findings. The `OrderedMap` nil-safety, `AsXmlID()` keying, fail-closed OCSP
`IsValid`, and bounds-checked ASN.1 parsing are all correct. Tests are thorough (KAT
vectors for every ESF structure, hostile-input error paths, determinism tests pinning the
insertion-ordered dedup). The deliberately-disabled staticcheck style rules and the
SHA-1-based cache keys (a digest, not a signature) are PORTING.md conventions, not
findings.

### Tool log

Run from `dss/`:

- `gofmt -l spi` → no output (clean).
- `go vet ./spi/` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./spi/...` → `0 issues`.

All tools available and passing; nothing in-scope (one of the 7 files) reported by a tool.

### Summary

- **Total findings: 4** (SEC: 0, PERF-BIGO: 1, PERF-MEM: 0, PERF-LEAK: 1, STD: 2)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 2, Low: 0, Info: 2
- The 7 in-scope files are a high-quality 1:1 port of the upstream Java DSS
  `QcStatementUtils`, `OCSPToken`, `RevocationToken`, `CMSCRLSource`,
  `CommonCertificateSource`, `FileRevocationSource`, and `ListCertificateSource` classes.
  All ASN.1 parsing is bounds-checked (cryptobyte / asn1ber), all trust-logic paths
  (OCSP signature verification, certID matching, revocation-status extraction, CMS CRL
  collection, certificate dedup) are correct and fail-closed, and the file cache is
  path-traversal-safe. The four findings are minor-to-moderate (a Medium unbounded
  on-disk cache, a Medium O(n²) entity-dedup, and two Info-level notes on deliberate
  parity quirk). No SEC, PERF-MEM, or Low issues were found.

### Open questions

1. T22B-STD-001: the `Clean()` XOR semantics are assumed to be a faithful port of the
   upstream Java `clean()`. Worth a one-line confirmation against the Java source at
   consolidation (the in-repo caller discards the return value, so the impact is nil
   either way).
2. T22B-PERF-001: the file cache is unbounded by design (parity with upstream). If any
   long-running validator in the supported deployment set validates documents over a large
   and growing set of revocation endpoints, a max-entry / max-age eviction policy should
   be considered and recorded under PORTING.md as a deliberate hardening.
3. The `ByCertificateDigest` method (`common_certificate_source.go:276`) recomputes
   `certificateToken.Digest(algorithm)` for every certificate on each call (no per-token
   digest cache visible in `model.CertificateToken`). This is O(total_certs) per call and
   is faithful to upstream; noted here so consolidation can decide whether it warrants a
   shared finding with the other O(n·m)/O(n²) dedup findings (T22A-PERF-BIGO-001,
   T22B-PERF-002).

## spi core — 66 smaller files (unit U22c)

- Scope: the 66 smaller top-level `.go` files in `dss/spi/` (≤383 lines each) —
  certificate/revocation references, matchers, selectors, factories, and utility
  types. The 6 largest (U22a) and the next ~13 mid-size files (U22b) are separate
  units and are **not** reviewed here. The `spi/validation/` subpackage and the other
  `spi` subpackages (`x509/`, `client/`, `eaa/`, `lote/`, `policy/`, `tsl/`, `alerts/`,
  `exception/`, `extension/`, `random/`, `signature/`) are separate units as well.
- Date: 2026-08-26
- Depth: budgeted (per-file first-100-lines + `grep -n 'func '` + security-relevant
  regions for files >100 lines; files ≤100 lines read in full). In practice every file
  in scope was read in full, because the security-relevant surface (matchers, selectors,
  candidate selection, the RSA verification fallbacks, the security factories) is
  concentrated in the smaller files.
- Compatibility contract: interop with Java DSS 6.5.RC1 (byte-level + verdict parity).
  `dss/PORTING.md` conventions and `docs/compatibility/known-gaps.md` are treated as
  deliberate, not findings. `// DIVERGENCE, deliberate:` and `DEVIATION` /
  `UPSTREAM QUIRK` comments mark intentional divergences.
- In-scope files (66): `list_certificate_source.go`, `repository_revocation_source.go`,
  `key_store_certificate_source.go`, `cms_ocsp_source.go`, `token_certificate_source.go`,
  `crl_ref.go`, `jdbc_revocation_source.go`, `rsa_large_exponent_verify.go`,
  `crl_token.go`, `dss_message_digest_calculator.go`, `certificate_ref.go`,
  `dss_signer_information_verifier_security_factory.go`, `signature_certificate_source.go`,
  `signer_identifier.go`, `ocsp_ref.go`, `certificate_reorderer.go`,
  `certificate_validity.go`, `list_revocation_source.go`, `oid.go`,
  `common_x509_url_certificate_source.go`, `signature_integrity_validator.go`,
  `equivalent_certificates_entity.go`, `dss_pk_utils.go`,
  `external_resources_ocsp_source.go`,
  `dss_content_verifier_provider_security_factory.go`, `ocsp_certificate_source.go`,
  `ocsp_token_ref_matcher.go`, `revocation_ref.go`, `candidates_for_signing_certificate.go`,
  `kid_certificate_source.go`, `external_resources_crl_source.go`,
  `certificate_token_ref_matcher.go`, `offline_crl_source.go`,
  `certificate_ref_identifier.go`, `certificate_source.go`,
  `baseline_b_certificate_selector.go`, `crl_token_ref_matcher.go`,
  `token_issuer_selector.go`, `offline_ocsp_source.go`, `composite_revocation_source.go`,
  `ocsp_response_binary.go`, `responder_id.go`, `ocsp_ref_identifier.go`,
  `common_trusted_certificate_source.go`, `ocsp_resp_status.go`,
  `dss_public_key_security_factory.go`, `dss_certificate_token_security_factory.go`,
  `dss_security_factory.go`, `dss_p7c_certificates_security_factory.go`,
  `revocation_token_ref_matcher.go`, `revocation_ref_identifier.go`,
  `web_token_payload.go`, `revocation_certificate_source.go`,
  `ocsp_signature_integrity_validator.go`, `alternate_urls_source_adapter.go`,
  `revocation_token_identifier.go`, `multiple_revocation_source.go`,
  `dss_security_provider.go`, `revocation_source_alternate_urls_support.go`, `doc.go`,
  `revocation_source.go`, `trusted_certificate_source.go`, `ocsp_source.go`,
  `crl_source.go`, `x509url_certificate_source.go`, `certificate_source_entity.go`,
  `proof_of_possession_certificate_source.go`.

### Findings

**SEC (top priority)**

T22C-SEC-001 — **Low** — Category: SEC — `dss/spi/ocsp_ref.go:157`
(`OCSPRef.String`)
- Evidence: both branches of `String()` dereference the responder unconditionally:
  `if r.responderId.X500Principal() != nil { … r.responderId.X500Principal().String() … }`
  and `… utils.ToBase64(r.responderId.Ski())`. `NewOCSPRef` (line 106) stores the
  `responderId` parameter with no nil guard.
- Impact: an `OCSPRef` built with a nil `ResponderId` panics on `String()` (Java would
  throw an NPE). `String()` is invoked from report/error rendering, so a hostile
  reference with a null responder that reaches a diagnostic path would crash the process
  rather than render. The matcher and equality paths already nil-check the responder
  (`ocspTokenRefMatcherMatchResponse`, `ocspRefResponderIdsEqual`), so this is isolated to
  the debug/render path and is not a false-accept/false-reject in the trust path.
- Recommendation: guard the two branches with `if r.responderId != nil` (render "null" when
  absent, matching the `crlRefString`/`certificateRefResponderIdString` convention already
  used in `crl_ref.go` and `certificate_ref.go`). No parity change (Java also throws here).

T22C-SEC-002 — **Info** — Category: SEC — `dss/spi/ocsp_ref.go:113`
(`NewOCSPRefFromOcspResponsesID`)
- Evidence: `digest, err := DSSRevocationUtilsDigest(ocspResponsesID.OcspRepHash)` and
  `DSSRevocationUtilsDSSResponderID(ocspResponsesID.OcspIdentifier.OcspResponderID)` run
  with no nil-check on `OcspRepHash` / `OcspIdentifier` / `OcspResponderID`.
- Impact: a by-name OCSP reference whose `OcspRepHash` (or responder) is nil reaches an
  upstream `Objects.requireNonNull`-style path that here surfaces as a nil dereference
  (panic). Reachable only from `cms_ocsp_source.go`'s `collectRevocationRefs`, where the
  `ParseCrlOcspRef` → `NewOCSPRefFromOcspResponsesID` chain is already inside an
  error-handling loop (a panic there is the port's faithful stand-in for Java's runtime
  exception). Not a trust-logic error; the trust-critical matching is the conservative
  `ocspTokenRefMatcherMatch*` functions (exact `bytes.Equal` on SKI, `Principal.Equals`
  on name, `ProducedAt` equality) — no false-accept.
- Recommendation: No action required (faithful to upstream; the panic is the documented
  port convention for Java's unchecked exceptions). If nil-responder OCSP refs are ever
  observed in the corpus, add a nil guard returning an error like the sibling `NewCRLRefFromCrlValidatedID`.

No other SEC findings. The mandated SEC lenses were verified clean across the 66 files:

- **Certificate/revocation matching (false-accept / false-reject):** All matchers are
  conservative and fail-closed. `CertificateTokenRefMatcher` (`certificate_token_ref_matcher.go`)
  uses exact `bytes.Equal` (digest) and `SignerIdentifier.IsRelatedToCertificate`
  (issuer+serial `Cmp==0` / SKI `bytes.Equal`); on a `SignerIdentifier` read error it
  returns `false` (no match) rather than surfacing — the conservative direction.
  `CRLTokenRefMatcher` / `OCSPTokenRefMatcher` require a non-nil digest (else an error),
  and the OCSP responder match is `bytes.Equal` on SKI **or** `Principal.Equals` on name
  (never a weaker comparison). `ResponderId.IsRelatedToCertificate` checks the subject
  name or the SKI — both exact. No path accepts a token on a partial/weak match.
- **Issuer/serial selection:** `TokenIssuerSelector.Issuer` filters by public key
  (`IsSignedByToken`) then by issuer-subject name, preferring an `IsValidOn(creationDate)`
  candidate and falling back to the first key-matching candidate — a correct,
  conservative issuer choice. `CertificateReorderer` / `BaselineBCertificateSelector`
  build the chain by public-key linkage and always retain the signing certificate
  (`orderedCertificates[0]`), only truncating at the trust anchor when the anchor-policy
  is enabled. `CandidatesForSigningCertificate.TheBestCandidate` prefers a
  `IsValid()` (digest/issuer-serial/responder match) candidate, else a signer-id-match,
  else the first — no trust grant without a positive check. `CRLToken.setRevocationStatus`
  refuses a CRL whose issuer is not the certificate's issuer
  (`DSSASN1UtilsX500PrincipalAreEquals`) and only marks the certificate revoked when the
  CRL entry is present — correct, no false-revoke/false-pass.
- **Nil-deref panics on malformed input:** All ASN.1 decoding in these files goes through
  `asn1ber.Parse` / `cryptobyte` (bounds-checked, error-returning) — `ParseCrlValidatedID`,
  `ParseOcspResponsesID`, `cmsOCSPSourceAsSequence`, `dssPKUtilsHasAlgorithmOID`. The
  reachable nil-deref panics are confined to the two noted above (T22C-SEC-001/-002) and to
  the deliberate `Objects.requireNonNull`-style constructor guards, which are the port's
  documented stand-in for Java's unchecked exceptions and not hostile-input vectors.
- **Key-size / algorithm checks:** `rsa_large_exponent_verify.go` is a strictly-narrower,
  never-weaker RSA fallback: it applies only to keys `crypto/rsa` refuses (`E > 2^31-1`),
  rejects `s >= n`, rebuilds the full EMSA-PKCS1v1_5 block and compares with
  `subtle.ConstantTimeCompare`, and implements the full EMSA-PSS-VERIFY (trailer, zeroed
  bits, PS/0x01 separator, recomputed H′) — a malformed padding, truncated digest, or
  trailing-byte forgery is rejected exactly as `crypto/rsa` would. `dss_pk_utils.go`
  computes key sizes without a size *check* (it reports `?`/`-1`/`0` for unknown types,
  which is the correct "cannot determine" direction, not a pass). No algorithm downgrade.
- **Security factories:** `dss_signer_information_verifier_security_factory.go` and
  `dss_content_verifier_provider_security_factory.go` verify through
  `crypto/x509.Certificate.CheckSignature` (authoritative) with fallbacks that are each
  strict (the no-NULL `DigestInfo` variant still rebuilds the full EM block; the
  large-exponent fallback is T22C-SEC-verified above). An unrecognised algorithm yields a
  `NoSuchAlgorithmException` error, never a silent success. The `buildWithProvider` nil-key
  paths return a graceful `InvalidKeyException` (documented fix, verified against
  `PAdES-LT.pdf`), not a panic.

**PERF (BIGO / MEM / LEAK)**

T22C-PERF-BIGO-001 — **Low** — Category: PERF-BIGO — `dss/spi/certificate_reorderer.go:140`
(`signingCertificates`) & `:83` (`certificateReordererBuildCertificateChainForCert`)
- Evidence: `signingCertificates` is a double loop (`for _, signer { for _, token { if
  signer.PublicKey().Equals(token.PublicKeyOfTheSigner()) } }`);
  `certificateReordererBuildCertificateChainForCert` calls the linear
  `certificateReordererContains` on every chain step.
- Impact: O(n²) in the certificate-chain size. In practice chains are small (a handful to a
  dozen certificates) and this mirrors Java's `List.contains`/nested-loop semantics exactly,
  so it is not exploitable by attacker input on the normal path. Same pattern as
  T22A-PERF-BIGO-001 and the `list_certificate_source.go` `Entities()` dedup (T22C-PERF-BIGO-002).
- Recommendation: None required. If large certificate pools are ever passed, index the
  collection by public-key hash to make the lookup O(1).

T22C-PERF-BIGO-002 — **Low** — Category: PERF-BIGO — `dss/spi/token_certificate_source.go`
(`ReferencesForCertificateToken`, `IsOrphan`, `AddCertificateWithOrigin`)
- Evidence: each does a linear scan of `certificateRefOrigins` / `certificateOrigins`
  (e.g. `for _, entry := range s.certificateRefOrigins { if s.doesCertificateReferenceMatch(…) }`).
- Impact: O(n) per lookup, O(n²) over a batch of `FindTokensFromRefs`. Bounded by the
  number of refs/tokens a single signature embeds (small); faithful to Java's
  `List`-backed collections. Not exploitable.
- Recommendation: None required (data is small and the structure mirrors Java's
  `LinkedHashMap` value-lookup).

T22C-PERF-MEM-001 — **Info** — Category: PERF-MEM — `dss/spi/common_x509_url_certificate_source.go:25`
(`mapByUri`) & `kid_certificate_source.go:22` (`mapByKid`)
- Evidence: `mapByUri map[string][]*model.CertificateToken` grows one entry per distinct
  `x5u` URI via `AddCertificateForURL`; `mapByKid map[string]*model.CertificateToken`
  likewise per KID.
- Impact: unbounded in the number of *distinct* URIs/KIDs the caller registers. Bounded by
  caller-supplied data (not by the document under validation) and released by `reset()`.
  This mirrors Java's `HashMap<String, List<CertificateToken>>`. Not a leak on the normal
  path.
- Recommendation: None required. Callers that stream many distinct `x5u` values over a long-
  lived source should call `Reset()`.

No PERF-LEAK findings. The 66 files contain no goroutines, no retained cross-request
references beyond the per-instance caches noted above, and no global mutable state
(`oid.go` holds immutable OID `var`s; the security-factory `var` instances are stateless
function holders).

**STD (idioms, error handling, sync, doc)**

T22C-STD-001 — **Low** — Category: STD — `dss/spi/dss_message_digest_calculator.go:71`
(`UpdateRange`)
- Evidence: `chunk := bytes[offset : offset+length]` — a raw slice with no bounds check.
- Impact: a caller-supplied `offset`/`length` outside `len(bytes)` panics (`slice bounds
  out of range`). Faithful to Java's `update(byte[], int, int)` throwing
  `ArrayIndexOutOfBoundsException`, and the in-package callers pass well-formed ranges
  (the `Writer` path correctly slices only the accepted `n` bytes). Programmer-error
  surface only, not reachable by hostile input.
- Recommendation: Optional — clamp or validate `offset`/`length` and return/ignore an
  out-of-range window instead of panicking, if this calculator is ever exposed to
  untrusted ranges.

T22C-STD-002 — **Info** — Category: STD — `dss/spi/crl_token_ref_matcher.go:31`,
`ocsp_token_ref_matcher.go:29`, `offline_crl_source.go:59`, `offline_ocsp_source.go:66`
- Evidence: unchecked type assertions reproduced as panics, e.g. `crlToken, ok :=
  token.(*CRLToken); if !ok { panic("ClassCastException …") }` and
  `binary.(*crlparser.CRLBinary)`.
- Impact: a type mismatch panics the process rather than returning an error. This is the
  port's documented, consistent stand-in for Java's `ClassCastException` (a programmer
  error, not a hostile-input vector), applied uniformly across the matcher/source
  dispatch. Not a correctness or trust issue.
- Recommendation: None required (deliberate port convention, consistent with U22a's
  T22A-STD-001 note on deliberate panics).

T22C-STD-003 — **Info** — Category: STD — `dss/spi/signature_certificate_source.go:156`
(`InitCandidatesList`)
- Evidence: in the `ProofOfPossessionCertificateSource` branch, when refs resolve to
  tokens the loop re-iterates `certificates` (the source's own list), not
  `certificateTokens` (the just-resolved tokens) — flagged by an in-code `NOTE`.
- Impact: an upstream quirk reproduced verbatim (per PORTING.md's fidelity rule). It does
  not change which candidates are *verified* (all candidates still flow through
  `SignatureIntegrityValidator.Validate`), so there is no trust-logic error; it is noted
  for completeness.
- Recommendation: None (documented deliberate reproduction of the upstream behaviour).

T22C-STD-004 — **Info** — Category: STD — `dss/spi/dss_pk_utils.go:89`
(`DSSPKUtilsPublicKeySize`)
- Evidence: the `ed25519.PublicKey` / `*ecdh.PublicKey` / Ed448-X448-OID branches return
  `dssPKUtilsEncodedKeySize` = `len(Encoded()) - 12` — a *byte* count — whereas every
  other branch returns a bit length.
- Impact: a consumer that treats this value uniformly as a bit length will under-report
  Ed25519 (32) / Ed448 (57) / X25519 (32) / X448 (57) sizes. This is a documented
  `UPSTREAM QUIRK` reproduced verbatim (Java's `getPublicKeySize` measures the SPKI byte
  length for these), so it is faithful to upstream and not a port defect; flagged so a
  downstream key-size policy is aware the unit differs by key type.
- Recommendation: None (faithful to upstream). If a future requirement needs a uniform
  bit-length, normalise at the policy-consumption site, not here.

No other STD findings. Error handling is consistent with the port's conventions (data-
dependent failures returned as `error`; Java's unchecked exceptions reproduced as `panic`),
there are no goroutines or shared mutable state to synchronise, and doc comments carry the
`// Ported from` attribution and the `DEVIATION`/`UPSTREAM QUIRK` markers required by
PORTING.md.

### Tool log

Run from `dss/` (2026-08-26):

- `gofmt -l spi` → no output (clean).
- `go vet ./spi/` → no output (clean).
- `golangci-lint run --config=../.github/.golangci.yml ./spi/...` → `0 issues` (clean).

All three tools available and passing; nothing to note as unavailable. Nothing the tools
report falls in the 66 in-scope files (the package as a whole, including U22a/U22b and the
subpackages, is clean).

### Summary

- **Total findings: 9** (SEC: 2, PERF-BIGO: 2, PERF-MEM: 1, PERF-LEAK: 0, STD: 4)
- **Severity breakdown:** Critical: 0, High: 0, Medium: 0, Low: 5, Info: 4
- The 66 in-scope files are a high-quality, faithful 1:1 port of the upstream Java DSS
  certificate/revocation reference, matcher, selector, factory, and utility classes.
  Every trust-logic path (certificate/revocation matching, issuer/serial selection,
  candidate selection, CRL/OCSP status, and the RSA verification fallbacks) is
  conservative and fail-closed, with exact `bytes.Equal`/`Cmp==0` comparisons and
  bounds-checked ASN.1 — no false-accept or false-reject, no exploitable panic on
  hostile input. The only security-relevant item is a Low-severity reachable nil-deref in
  `OCSPRef.String()` (T22C-SEC-001), isolated to the render path. The remaining findings
  are Low/Info notes on O(n²)-on-small-data chain building, bounded caller-supplied caches,
  and the deliberate port conventions (deliberate panics, an upstream byte-vs-bit key-size
  quirk). No PERF-LEAK or goroutine/sync issues.

### Open questions

1. `OCSPRef.String()` (T22C-SEC-001): confirm no report builder calls `String()` on an
   `OCSPRef` whose `ResponderId` is nil in a path reachable from a hostile document. If so,
   apply the nil guard (parity-neutral).
2. `DSSPKUtilsPublicKeySize` (T22C-STD-004): is the returned key size ever fed into a
   minimum-key-size policy that assumes bits? If a future policy needs a uniform unit,
   normalise at the consumer (do not change the port here).
3. The deliberate-panic convention (T22C-STD-002) is consistent with U22a (T22A-STD-001);
   confirm during consolidation that no panic in this unit is reachable from untrusted
   *data* on the normal validation path (the review found none beyond the two noted, both
   faithful to upstream's own unchecked-exception behaviour).

## spi subpackages (unit U23)

- Scope: every file under `dss/spi/` subpackages EXCEPT the `spi` top-level files and
  `spi/validation/` (units U22a-c / U22d): `spi/client/http`, `spi/client/jdbc`,
  `spi/policy`, `spi/x509/aia`, `spi/x509/evidencerecord/digest`, `spi/lote`, `spi/tsl`,
  `spi/random`, `spi/signature/resources`, `spi/extension`, `spi/alerts`, `spi/exception`,
  `spi/eaa/status`, plus the `doc.go` files of the subpackages.
- Date: 2026-08-26
- Depth: full — every in-scope non-test file read in full (all ≤ 338 lines; 50 files,
  3 389 lines total). Security-relevant regions verified against the Java upstream at
  `~/Workspace/esig/dss/dss-spi/...` (BasicASN1SignaturePolicyValidator.java, Protocol.java,
  NativeHTTPDataLoader.java) and against the verdict data-flow
  (`spi/validation/analyzer/default_document_analyzer.go` → `diagnostic/signature_wrapper.go`
  → `validation/process/bbb/vci/signature_policy_hash_valid_check.go`), the base certificate
  sources (`spi/common_certificate_source.go`, `spi/common_trusted_certificate_source.go`), and
  the random-source consumer (`internal/pdf/incremental.go`).
- Compatibility contract: interop with Java DSS 6.5.RC1 (verdict + byte parity).
  `dss/PORTING.md` and `docs/compatibility/known-gaps.md` are treated as deliberate, not
  findings. `// DIVERGENCE, deliberate:` / `DEVIATION` / `UPSTREAM QUIRK` comments mark
  intentional divergences.

### Files read (50)

| Subpackage | Files (non-test) |
|---|---|
| `spi/client/http` (11) | `native_http_data_loader.go`, `native_http_data_loader_call.go`, `advanced_memory_data_loader.go`, `advanced_data_loader.go`, `data_loader.go`, `protocol.go`, `memory_data_loader.go`, `max_size_input_stream.go`, `ignore_data_loader.go`, `dss_cache_file_loader.go`, `dss_file_loader.go` |
| `spi/client/jdbc` (4) | `jdbc_cache_connector.go`, `sql_query.go`, `sql_select_query.go`, `sql_record.go` |
| `spi/policy` (10) | `basic_asn1_signature_policy_validator.go`, `signature_policy_provider.go`, `default_signature_policy_validator_loader.go`, `non_asn1_signature_policy_validator.go`, `zero_hash_signature_policy_validator.go`, `empty_signature_policy_validator.go`, `abstract_signature_policy_validator.go`, `signature_policy_validator.go`, `signature_policy_validator_loader.go`, `doc.go` |
| `spi/x509/aia` (5) | `repository_aia_source.go`, `default_aia_source.go`, `aia_certificate_source.go`, `composite_aia_source.go`, `aia_source.go` |
| `spi/x509/evidencerecord/digest` (2) | `data_object_digest_builder.go`, `data_object_digest_builder_factory.go` |
| `spi/lote` (2) | `trusted_entities_certificate_source.go`, `doc.go` |
| `spi/tsl` (2) | `trusted_lists_certificate_source.go`, `doc.go` |
| `spi/random` (3) | `dss_secure_random_provider.go`, `secure_random_provider.go`, `doc.go` |
| `spi/signature/resources` (2) | `dss_resources_handler.go`, `dss_resources_handler_builder.go` |
| `spi/extension` (2) | `document_extender.go`, `doc.go` |
| `spi/alerts` (2) | `dss_external_resource_exception_alert.go`, `doc.go` |
| `spi/exception` (4) | `dss_data_loader_multiple_exception.go`, `dss_external_resource_exception.go`, `illegal_input_exception.go`, `doc.go` |
| `spi/eaa/status` (1) | `eaa_revocation_source.go` |

### Findings

**SEC (top priority)**

T23-SEC-001 — **High** — Category: SEC — `dss/spi/policy/basic_asn1_signature_policy_validator.go:91`
(`Validate`)
- Evidence: line 71 unconditionally runs `validationResult.SetDigestValid(true)` before the
  deferred-`recover` block. Inside that block, `if !ok || len(elements) < 3 { … return }`
  (line 91–96) exits *without* resetting `digestValid`, so a top-level SEQUENCE that has < 3
  children leaves `identified=true`, `digestValid=true`, `asn1Processable=false`, no error.
  The in-code comment asserts this "mirrors the effect" of Java's `getObjectAt` throwing.
- Impact: it does **not** mirror Java. Verified against
  `BasicASN1SignaturePolicyValidator.java`: a SEQUENCE with < 3 children makes
  `asn1Sequence.getObjectAt(0)`/`getObjectAt(2)` throw, which the Java `catch (Exception e)`
  turns into `setDigestValid(false)` + `addError`. So for the same input **Java reports the
  policy digest invalid, Go reports it valid** — a false-accept on the signature-policy-hash
  verification check. This value flows straight into the verdict: `validateSignaturePolicy`
  → `Policy.SetValidationResult` → `SignaturePolicyHashValidCheck.Process()` =
  `IsPolicyPresent() && IsPolicyDigestValid()` (`diagnostic/signature_wrapper.go:917`,
  reads `DigestAlgoAndValue.Match` = `IsDigestValid()`), i.e. the VCI `BBB_VCI_ISPM` building
  block passes in Go where it fails in Java. Reached for any ASN.1 policy (first byte `0x30`,
  `CanValidate` true) whose SEQUENCE body has fewer than three children.
- Recommendation: On the `!ok || len(elements) < 3` branch set
  `validationResult.SetDigestValid(false)` (and `SetAsn1Processable(false)`) and add the
  general error, matching Java's `catch`. Fix the misleading "mirrors the effect" comment.
  Add an oracle test pinning a `<3`-element SEQUENCE policy to `digestValid=false`. If, on the
  other hand, the Go behaviour is deemed intentional, record it as a `// DIVERGENCE,
  deliberate:` with the Java method named and the motivating input, per PORTING.md.

T23-SEC-002 — **Medium** — Category: SEC — `dss/spi/client/http/native_http_data_loader_call.go:168`
(`createClient`) & `dss/spi/client/http/max_size_input_stream.go:52`
- Evidence: `createClient` builds an `http.Transport` with only `DialContext` — no
  `ResponseHeaderTimeout`, no `CheckRedirect`/`RedirectPolicy`. The body is read via
  `utils.ToByteArray` (= `io.ReadAll`) bounded by `MaxSizeInputStream` **only when
  `c.maxInputSize > 0`** (line 142). The real validation flows use the default
  (`spi/validation/analyzer/default_document_analyzer.go:427` and
  `spi/x509/aia/default_aia_source.go:36` both call `NewNativeHTTPDataLoader()` with no
  `SetMaxInputSize`), so `maxInputSize == 0` → **no size cap and (with the default
  `readTimeout==0`) no read deadline**. Separately, `checkSize` errors only once
  `s.count > s.maxSize`, so it can return up to `len(p)` bytes past the cap before erroring.
- Impact: the fetched URLs are attacker-influenced (AIA `caIssuers`, TL, OCSP/CRL endpoints
  carried inside the document under validation). A malicious document (or a redirector, since
  Go follows up to 10 redirects by default) can point these fetches at an endpoint that
  streams unbounded / stalls, exhausting memory or holding the worker — a DoS on a
  verification network path. The library's own `MaxSizeInputStream` is the intended guard but
  is off by default. (Not a trust-logic false-accept; a resource-exhaustion / hardening gap.)
- Recommendation: Give the loader a sensible default `maxInputSize` (and a default
  `readTimeout`) for the AIA/TL/OCSP fetch paths, or apply a cap in the analyzer/AIA defaults;
  add `ResponseHeaderTimeout` and a bounded `RedirectPolicy`. Tighten `MaxSizeInputStream` to
  error as soon as `s.count >= s.maxSize` (or clamp the returned `n`) so the cap is not
  over-read by a chunk.

T23-SEC-003 — **Info** — Category: SEC — `dss/spi/random/dss_secure_random_provider.go` (whole file)
- Evidence: `DSSSecureRandomProvider.GetSecureRandom(seed)` returns a `dssFixedSecureRandom`
  that emits `digest(seed)`, `digest(digest(seed))`, … (a deterministic hash-chain, not an
  entropy source). Verified the consumer: `internal/pdf/incremental.go:88-91` uses
  `d.opts.Random` **or falls back to `rand.Reader`** (`crypto/rand`) when it is nil, and the
  only producers of a non-nil `Random` are the test/oracle harness (per this file's `doc.go`
  and the `internal/pdf` R20/R21 note).
- Impact: the "crypto/rand vs math/rand" lens is satisfied — the normal signing path derives
  AES IVs from `crypto/rand`; the deterministic source is a deliberate reproducibility seam
  mirroring upstream's `FixedSecureRandom`, not a `math/rand` misuse. No defect.
- Recommendation: None required. Guard against regressions: keep the deterministic provider
  out of any production signing/timestamping path (it must never replace `crypto/rand` for
  key/IV derivation), and consider a `// do not use for production entropy` note on
  `GetSecureRandom`.

**PERF**

T23-PERF-001 — **Low** — Category: PERF-LEAK (no pooling) —
`dss/spi/client/http/native_http_data_loader_call.go:168` (`createClient`)
- Evidence: every `Call()` constructs a fresh `net.Dialer` and `http.Transport`
  (`createClient`, lines 168–181); the `Client`/`Transport` are discarded when `Call` returns.
- Impact: no connection reuse / TLS-session resumption across the many AIA/TL/OCSP/CRL
  fetches a single validation performs, and an allocation + fresh dial per request. Not a hard
  leak (Go's `Transport` finalizer closes idle connections once the per-request `Transport` is
  GC'd), but a real efficiency cost on a fetch-heavy path.
- Recommendation: Hoist a single `*http.Transport`/`*http.Client` to the `DataLoader`
  (created once, reused, with `MaxIdleConns`/`IdleConnTimeout`) instead of per-call; pass the
  per-call `connectTimeout`/`readTimeout` via `context`/`Client.Timeout` rather than a new
  `Transport`.

(No PERF-BIGO findings: the in-scope subpackages contain no O(n²)-on-attacker-input loops —
the only nested scans are small bounded dedups/`contains` over certificate lists, already
covered as Low in U22c. No PERF-MEM unbounded-growth findings beyond the SEC-002 default-cap
gap above.)

**STD**

T23-STD-001 — **Low** — Category: STD — `dss/spi/lote/trusted_entities_certificate_source.go:79`
(`SetTrustedPropertiesByCertificates`) & `:130` (`SetTrustTimeByCertificates`)
- Evidence: both setters range the caller's map directly:
  `for certificateToken, trustPropertiesList := range trustPropertiesByCerts` (line 79) and
  `for certificateToken, certificateTrustTimes := range trustTimeByCertificate` (line 130).
  The sibling `spi/tsl/trusted_lists_certificate_source.go` does the same job but sorts the
  keys first via `trustedListsCertificateSourceSortedTokens` (line 76, 149) precisely because
  Go map order is randomized.
- Impact: `CommonCertificateSource` stores certificates in insertion order
  (`spi/common_certificate_source.go:33-36`), so LoTE's `Certificates()`/`Entities()` — and
  anything downstream that enumerates trust anchors or emits diagnostic ordering — are
  **non-deterministic run-to-run**, whereas TSL is stable. This violates PORTING.md's
  "order-sensitive upstream iteration → … explicit sort" rule and is inconsistent with the
  TSL twin. No trust-verdict impact (membership tests are order-independent), but report/diagnostic
  ordering diverges from Java's stable-per-set order and differs from run to run.
- Recommendation: Reuse the TSL `trustedListsCertificateSourceSortedTokens` approach (sort by
  `DSSIDAsString()`) in both LoTE setters so LoTE and TSL enumerate identically.

T23-STD-002 — **Low** — Category: STD — `dss/spi/lote/trusted_entities_certificate_source.go:67`
- Evidence: the LoTE `AddCertificate` panic message reads
  `panic("Cannot directly add certificate to a TrustedListsCertificateSource")` — the *TSL*
  type name, inside the *LoTE* (`TrustedEntitiesCertificateSource`) type.
- Impact: a copy-paste from the TSL source. Misleading when debugging a LoTE source that
  rejects direct `AddCertificate` (points the operator at the wrong type). No behavioural
  effect (the panic is intentional and correct).
- Recommendation: Change the message to `TrustedEntitiesCertificateSource` (or a neutral
  "this certificate source").

T23-STD-003 — **Info** — Category: STD — `dss/spi/client/jdbc/jdbc_cache_connector.go:95`
(`TableQuery`)
- Evidence: the comment at line 95 says Java "catches SQLException and returns false
  **without rolling back**; mirrored here", but the very next line is `_ = tx.Rollback()`
  (line 97).
- Impact: the comment contradicts the code. The *behaviour* (rolling back on failure) is the
  more correct choice, so this is a documentation-only defect; the reader is told one thing and
  the code does another.
- Recommendation: Fix the comment to state that the port rolls back on failure (a deliberate,
  strictly-safer improvement over Java's no-rollback), or align the code to the comment —
  prefer the former.

T23-STD-004 — **Info** — Category: STD — `dss/spi/client/jdbc/jdbc_cache_connector.go`
(`Execute`/`Select`/`TableQuery`/`ExecuteThrowable`)
- Evidence: all four swallow errors and return the Java default (`0`, `nil`, `false`)
  (e.g. `Execute` returns `0` on `Begin`/`Exec`/`Commit`/`RowsAffected` failure; `Select`
  returns `nil` on any failure). `ExecuteThrowable` is the only one that propagates `(int,
  error)`.
- Impact: faithful to Java's `JdbcCacheConnector` catch-and-return-default contract (the
  `*Throwable` variants are the error-propagating overloads). The commit-failure branches do
  not always roll back, matching Java's asymmetry. Not a defect; noted so a reader does not
  mistake the error-swallowing for a missed `errcheck`.
- Recommendation: None (deliberate Java-faithful contract; `ExecuteThrowable` already exists
  for callers that need the error).

No other STD findings. The deliberate-panic convention (Java `Objects.requireNonNull` →
`panic`) is applied uniformly and matches U22a/U22c. The AIA `compositeSourceOrderedKeys`
(sort) and `defaultSourceDedup`/`aiaCertificateSourceDedup`/`extractAIAFromCacheSource`
(map-based dedup preserving order) are correct and deterministic. `DSSDataLoaderMultipleException`
uses Go 1.20+ `Unwrap() []error` correctly and sorts URLs for deterministic `Error()` output.
No goroutines, no global mutable state, no unchecked type assertions (only the documented
`ClassCastException`-style panics already noted in U22c).

### Tool log

Run from `dss/` (2026-08-26):

1. `gofmt -l spi` → **no output** (clean).
2. `go vet ./spi/client/... ./spi/policy/... ./spi/x509/... ./spi/lote/... ./spi/tsl/...
   ./spi/random/... ./spi/signature/... ./spi/extension/... ./spi/alerts/...
   ./spi/exception/... ./spi/eaa/...` → **no output** (clean, exit 0).
3. `golangci-lint run --config=../.github/.golangci.yml ./spi/client/... ./spi/policy/...
   ./spi/x509/... ./spi/lote/... ./spi/tsl/... ./spi/random/... ./spi/signature/...
   ./spi/extension/... ./spi/alerts/... ./spi/exception/... ./spi/eaa/...` → **`0 issues.`**
   (clean).

All three tools available and passing; nothing the tools report falls in the 50 in-scope files.
No deliberate-exclusion note is required (no `ST*`/`S1*`/`QF*`/`SA1019` or legacy-crypto
findings surfaced in these subpackages).

### Summary

- **Total findings: 7** (SEC: 3, PERF: 1, STD: 3).
- **Severity breakdown:** Critical: 0, High: 1, Medium: 1, Low: 4, Info: 2 — i.e. SEC
  1 High + 1 Medium + 1 Info; PERF 1 Low; STD 2 Low + 1 Info.
- The single material item is **T23-SEC-001**: a genuine false-accept divergence from Java
  in the ASN.1 signature-policy validator (a `<3`-element SEQUENCE policy reports
  `digestValid=true` in Go but `false` in Java) that reaches the VCI policy-hash check — the
  in-code "mirrors the effect" comment is incorrect. **T23-SEC-002** is a Medium hardening
  gap (unbounded response read + no redirect/timeout policy on attacker-influenced AIA/TL/OCSP
  URLs, whose cap is off by default). The rest are Low/Info port-convention and
  documentation nits (LoTE non-deterministic enumeration vs the TSL twin, a copy-paste panic
  string, a contradictory JDBC comment, the deliberate JDBC error-swallow contract) and a
  positive confirmation that the random path is `crypto/rand` in production. No PERF-BIGO,
  no PERF-MEM (beyond the SEC-002 default-cap gap), no PERF-LEAK, and no goroutine/sync
  issues.

### Open questions

1. **T23-SEC-001 (severity):** the finding is a false-accept on the VCI `BBB_VCI_ISPM`
   check, which under this rubric's "false-accept in a verification/trust path = Critical"
   bar could be rated Critical; it is filed as **High** because it does not by itself flip an
   invalid signature to valid (other independent checks still run) and requires a crafted
   malformed policy. Please confirm the intended severity and that an oracle test is added for
   a `<3`-element SEQUENCE policy during consolidation.
2. **T23-SEC-002 (default cap):** does the project intend a default `maxInputSize` for the
   AIA/TL/OCSP fetch paths, or is "0 = unlimited" an accepted default (parity with Java's
   `maxInputSize` default of 0)? If the latter, document it in `known-gaps.md` alongside the
   redirect/timeout posture so the DoS surface is a conscious choice.

## spi/validation — verdict core (unit U22d1a)

### Scope

The 4 largest verdict-critical files in `dss/spi/validation/`: `signature_validation_context.go`
(2014), `default_advanced_signature.go` (852), `revocation_data_verifier.go` (680),
`common_certificate_verifier.go` (523). Research + report only; no source edits. All findings were
cross-checked against the Java DSS 6.5.RC1 oracle at `~/Workspace/esig/dss/dss-spi/.../spi/validation/`.

### Date

2026-08-26

### Depth

Strict-budget (per-file line caps; ~1,909 in-scope code lines + grep/tool output total, at the 1,900 target).

### Files read (line counts actually read)

| File | Lines read | Approach |
| --- | --- | --- |
| `signature_validation_context.go` | ~420 | 1–150 + `func` map + 4 verdict regions (1142–1258, 1347–1360, 1552–1582, 1583–1600, 1681–1694, 1755–1800) |
| `default_advanced_signature.go` | ~286 | 1–150 + `func` map + 3 regions (210–232, 545–600, 724–780) |
| `revocation_data_verifier.go` | 680 | full read (1–680) |
| `common_certificate_verifier.go` | 523 | full read (1–523) |

Not fully verified (budget): `signature_validation_context.go` lines outside the 5 regions above
(issuer-resolution `getIssuerWithSource`/`getIssuersFrom*`, POE bookkeeping `registerPOE*`,
`getValidationData*` recursion). The verdict-computation, status-determination, revocation-freshness
and trust-anchor call sites were all read in full.

### Findings

**SEC**

T22D1A-SEC-001 — **High** — Category: SEC — `dss/spi/validation/signature_validation_context.go:1142,1233,1569`
- Evidence:
  - `CheckAllRequiredRevocationDataPresent()` → `status, _ := c.allRequiredRevocationDataPresent(); return status.IsEmpty()`
  - `allRequiredRevocationDataPresent()` (:1151) → `orderedCertificateChains, err := c.getOrderedCertificateChains(); if err != nil { return status, err }` with `status` a fresh empty `TokenStatus`
  - `CheckAllPOECoveredByRevocationData()` (:1233) → same `status, _ := …; return status.IsEmpty()` swallow
  - `checkAtLeastOneRevocationDataPresentAfterBestSignatureTime()` (:1569) → `if err != nil { return }`
  - Java: `SignatureValidationContext.java:509` `getOrderedCertificateChains()` → `order.getOrderedCertificateChains()` which **throws** `DSSException("No signing certificate found")` / `DSSException("The certificate chain contains only bridge certificates")` (`CertificateReorderer.java:216,236`)
- Impact: **fail-OPEN divergence.** When the processed-certificate set is unorderable (all-bridge, no identifiable signing/leaf), Go's `getOrderedCertificateChains()` returns an error that the three call sites swallow, leaving the status empty → `IsEmpty()==true` → the LT "all revocation data present", "POE covered", and "fresh revocation data" requirements **report pass**. Java throws, aborting the check (fail-closed). This is a false-accept of the LT-level revocation sub-check; it can flip a document to a higher level / VALID where Java errors. Rated **High** (not Critical) consistent with the T23-SEC-001 convention: the same malformed input typically also breaks the trust-anchor, not-expired, and crypto-integrity checks, so it is not a standalone flip — but a consolidated severity decision + an oracle test for an unorderable-certs input is requested.
- Recommendation: Do not swallow. Propagate the error to the caller so the LT requirement reports *not-satisfied* (or surfaces the DSS error), matching Java's throw. Minimal: `if err != nil { return false }` at the three `Check*`/`check*` entry points, or have the `BaselineRequirementsChecker` treat a reorderer error as a failed requirement.

T22D1A-SEC-002 — **Low** — Category: SEC — `dss/spi/validation/revocation_data_verifier.go` (`isSelfIssuedRevocation`)
- Evidence:
  - `certificateChain, err := spi.NewCertificateReordererWithSigningCertificate(revocationData.IssuerCertificateToken(), revocationData.Certificates()).OrderedCertificates(); if err == nil { … return true }` then `return false`
  - Java: `isSelfIssuedRevocation` → `new CertificateReorderer(signingCert, chain).getOrderedCertificates()` **throws** `DSSException` on an unorderable chain
- Impact: On a reorderer failure Go treats the token as "not self-issued" and **continues** the `isCertificateNotRevoked` loop; Java aborts the validation. The branch taken is the conservative/natural one (a self-issued token that can't be re-ordered simply isn't matched), so there is **no false-accept** — but it is a behavioural divergence worth a parity note.
- Recommendation: No change required for correctness; add a `// DIVERGENCE, deliberate:` comment (or confirm the conservative branch is intended) so it is not mistaken for an oversight.

**PERF**

T22D1A-PERF-BIGO-001 — **Info** — Category: PERF-BIGO — `dss/spi/validation/signature_validation_context.go` / `revocation_data_verifier.go`
- Evidence: `getSignaturesIssuedByCertificateOrItsChildren`, `contains*`, and `getOrderedCertificateChains()` each re-scan the processed-certificate set; `checkAtLeastOneRevocationDataPresentAfterBestSignatureTime()` is invoked per-signature and re-calls `getOrderedCertificateChains()`.
- Impact: O(n²) on the (typically small, attacker-bounded) certificate set. Same complexity class as **T22C-PERF-BIGO-001** (already filed for this package) — not a new finding; noted so the recurrence is visible in the verdict path.
- Recommendation: No action; covered by T22C-PERF-BIGO-001.

No **PERF-MEM** / **PERF-LEAK** findings: the context's caches (`tokensToProcess`, `poeTimes`, `tokenIssuerMap`, `certificateChildrenMap`, …) are bounded by the processed token set for a single validation request and are not retained across requests; no goroutines, no global mutable state in the 4 files.

**STD**

T22D1A-STD-001 — **Low** — Category: STD — `dss/spi/validation/signature_validation_context.go:1213`
- Evidence:
  - `if freshnessStatus, ok := any(status).(*RevocationFreshnessStatus); ok { … freshnessStatus.AddTokenAndRevocationNextUpdateTime(…) }`
  - the `status` parameter is typed `*TokenStatus` and the caller passes `&status.TokenStatus` (the promoted field), so `any(status)` is always a `*TokenStatus` → the assertion **always fails** in Go
  - Java: `if (status instanceof RevocationFreshnessStatus)` is **true** on the same call path → the block runs
- Impact: `RevocationFreshnessStatus.tokenRevocationNextUpdateMap` is therefore **never populated** in Go; `MinimalNextUpdateTime()` always returns zero and `ErrorString()` omits the "nextUpdate available after …" suffix that Java appends. **No verdict impact** — the pass/fail boolean is driven by the independently-added `AddRelatedTokenAndErrorMessage` calls, not by this map; the effect is limited to the report/alert text (display parity). Consistent with the porters' own "latent no-op" note in `revocation_freshness_status.go`.
- Recommendation: Either pass the concrete `*RevocationFreshnessStatus` into `checkRevocationForCertificateChainAgainstBestSignatureTime` on the freshness path, or drop the dead block and document the deliberate drop, so the report text matches Java.

No other **STD** findings in the 4 files: nil handling is guarded (`dssPKUtilsEncodedKeySize`/`PublicKey.Encoded()`/`Key()` are nil-safe; `isAcceptableWithIssuer`/`isGoodIssuer`/`isPublicKeySizeSupported` all fail-closed on nil), no unchecked type assertions beyond the one above, no goroutines, and the `panic` sites (`maximumRevocationFreshness` unsupported-context, `SetRevocationDataLoadingStrategyFactory` nil, `IsRevocationDataSkip` ASN.1-parse) are deliberate unchecked-exception ports recovered to errors at the facade — not reported.

### Tool log (run from `dss/`)

```
$ gofmt -l spi/validation            → (no output; clean)
$ go vet ./spi/validation/           → (no output; clean)
$ golangci-lint run --config=../.github/.golangci.yml ./spi/validation/ → 0 issues.
```

All three in-scope tools are clean; no tool-derived findings. No out-of-scope tool reports.

### Summary

- **Total findings: 4** (SEC: 2, PERF: 1, STD: 1).
- **Severity breakdown:** Critical: 0, High: 1, Medium: 0, Low: 2, Info: 1 — i.e. SEC 1 High (T22D1A-SEC-001) + 1 Low (T22D1A-SEC-002); PERF 1 Info (T22D1A-PERF-BIGO-001, already-covered class); STD 1 Low (T22D1A-STD-001).
- **No Critical, no Medium, no PERF-LEAK, no PERF-MEM, no goroutine/sync issues.** The single material item is **T22D1A-SEC-001**: a fail-OPEN divergence on an unorderable certificate set (Go reports the LT revocation sub-check as *pass*; Java throws), which is the one path in these 4 files that can contribute to a false level-grant / false-VALID. Everything else is a display-parity nit, a conservative-branch parity note, or an already-filed complexity class.

### Open questions

1. **T22D1A-SEC-001 (severity):** filed **High** on the T23-SEC-001 convention (gated by other independent checks). If consolidation confirms a reachable input where the reorderer fails on the *aggregate* cert set while the signing chain independently passes crypto/trust/not-expired, this is a standalone false-VALID and should be re-rated **Critical**. Please confirm the intended severity and that an oracle test is added for an all-bridge / no-signing-cert input.
2. **T22D1A-SEC-002 / T22D1A-STD-001:** both are conservative/display-only divergences; confirm the project wants a `// DIVERGENCE, deliberate:` marker (or a fix) rather than a silent port, so they are not later mistaken for latent bugs.
3. **T23-STD-001 (LoTE order):** confirm LoTE enumeration order is not asserted anywhere in a
   byte-comparison corpus (it is order-sensitive for diagnostics only); if it is, the sort fix
   is required for parity, if not, it is a determinism improvement.

## spi/validation — mid core (unit U22d1b)

### Scope

The **mid-size core** of `spi/validation` — the document analyzer, the local TSP source, the
baseline-requirements checker, the `AdvancedSignature` / `EAAPayload` / `CertificateVerifier`
contracts, the validation-data container and record, the signature alerter, and the trust-anchor
verifier. 1:1 port of Java DSS 6.5.RC1; a false "VALID" verdict is Critical. Four of the ten files
(`advanced_signature.go`, `eaa_payload.go`, `certificate_verifier.go`, and the interface surface of the
rest) are pure interface declarations with no executable logic; the findings below concern the
logic-bearing files only.

### Date

2026-08-26

### Depth

Budgeted. The two largest files (`default_document_analyzer.go`, `key_entity_tsp_source.go`) were
read header + `func` map + targeted logic regions; the eight ≤390-line files were read in full.
Every trust-logic path that could flip a verdict was cross-checked against the upstream Java at
`~/Workspace/esig` (DSS 6.5.RC1 pin).

### Files read (line counts actually read)

| File | Read / total | Mode |
|---|---|---|
| `analyzer/default_document_analyzer.go` | 739 / 924 | budgeted (header, func map, validate/data/evidence/counter-sig regions; unread ~185 lines are trivial setters/getters) |
| `key_entity_tsp_source.go` | 699 / 699 | full |
| `baseline_requirements_checker.go` | 390 / 390 | full |
| `advanced_signature.go` | 386 / 386 | full (interface only — no logic) |
| `eaa_payload.go` | 332 / 332 | full (interface only — no logic) |
| `certificate_verifier.go` | 273 / 273 | full (interface only — no logic) |
| `validation_data_container.go` | 266 / 266 | full |
| `validation_data.go` | 256 / 256 | full |
| `signature_validation_alerter.go` | 237 / 237 | full |
| `trust_anchor_verifier.go` | 182 / 182 | full |

Total in-scope lines read: **3,737**.

### Findings

#### SEC

T22D1B-SEC-001 — **High** — Category: SEC — `dss/spi/validation/signature_validation_alerter.go:58` (and `:73`)
- Evidence:
  ```go
  status, _ := a.validationContext.allRequiredRevocationDataPresent()   // :58  (AssertAllRequiredRevocationDataPresent)
  if !status.IsEmpty() { a.populateMessage(status); a.alert(...) }        // :59-62
  status, _ := a.validationContext.allPOECoveredByRevocationData()       // :73  (AssertAllPOECoveredByRevocationData)
  ```
  `allRequiredRevocationDataPresent` (`signature_validation_context.go:1149`) and
  `allPOECoveredByRevocationData` (`:1239`) return `(status, err)`; the `err` is the
  `getOrderedCertificateChains()` failure (unorderable chain: ≥2 certs, all bridge, no signing
  cert). Java's `getOrderedCertificateChains()` **throws** `DSSException` here (fail-closed); the
  Go alerter discards the error and proceeds with an empty status.
- Impact: On an unorderable (hostile/malformed) certificate set, the discarded error leaves
  `status` empty, so `status.IsEmpty()` is `true` and **no alert is raised** — the check
  silently passes. This is on the **live main validation path**:
  `CompleteValidationContextExecutor.Validate` → `assertSignaturesValid` →
  `AssertAllRequiredRevocationDataPresent` / `AssertAllPOECoveredByRevocationData`
  (`executor/complete_validation_context_executor.go:43-44`), and it feeds the LT-level gate
  `MinimalLTRequirement` → `CheckAllRequiredRevocationDataPresent`
  (`baseline_requirements_checker.go:290,311`). A fail-closed upstream `throw` becomes a fail-open
  "pass" — a false-accept (missing revocation data / uncovered POE treated as present) in the
  revocation-presence, POE-coverage, and LT-level determination. Rated **High** (not Critical) to
  stay consistent with the sibling findings T22D1-SEC-001 and T23-SEC-001: it does not by itself
  flip an otherwise-invalid signature to valid (the cryptographic, trust-anchor and other
  independent checks still run), and it requires a crafted unorderable certificate set.
- Recommendation: Propagate the error instead of discarding it. In `AssertAllRequiredRevocationDataPresent`
  and `AssertAllPOECoveredByRevocationData`, capture `status, err := …` and, on `err != nil`,
  raise/return a fail-closed result (matching Java's thrown `DSSException`) rather than falling
  through to the `status.IsEmpty()` test. This is the **same root cause as T22D1-SEC-001**
  (the `signature_validation_context.go` `Check*` methods swallow the identical
  `getOrderedCertificateChains` error) at a different, in-scope surface — consolidate the two
  under one fix and add an oracle test that drives an unorderable chain (≥2 certs, all bridge, no
  signing cert) through `CompleteValidationContextExecutor` and asserts it does not pass.

#### PERF

No new PERF findings in this unit (0). Two O(n·m) shapes are present but are the same class
already covered by **T22C-PERF-BIGO-001** and are bounded by chain/evidence-record count (small in
practice), so they are not re-filed:
- `validation_data.go:150-157` `AddValidationData` iterates the other record's token sets and
  re-runs `containsCertificateToken` (a map `Get` + an `EntityKey` lookup) per token.
- `validation_data.go:182-190` `excludeWithEntityKey` scans all stored certificate tokens per
  excluded entity key.
Neither is attacker-influenced in a way that exceeds the existing O(n²) chain-build class; no
unbounded growth, no goroutines, no retained-reference/cache leaks in any of the ten files.

#### STD

T22D1B-STD-001 — **Info** — Category: STD — `dss/spi/validation/key_entity_tsp_source.go:266`
- Evidence:
  ```go
  // x/crypto/pkcs12 does not tell which certificate the key belongs to; the store lists the
  // key entry's own certificate first, which is also the order a time-stamp chain needs.
  return NewKeyEntityTSPSource(signer, chain[0], chain), nil   // :266  (NewKeyEntityTSPSourceFromKeyStore)
  ```
- Impact: The PKCS#12 loader assumes `chain[0]` is the key entry's own certificate. `golang.org/x/crypto/pkcs12`
  exposes no binding between the private key and its certificate, so a multi-certificate store whose
  first `CERTIFICATE` block is *not* the key's would cause the source to embed the wrong issuer
  certificate and compute the RFC 5035 `SigningCertificateV2` (ESSCertIDv2) over the wrong hash.
  This is **issuance-only** (a local TSP that *produces* tokens; it performs no verification, so it
  cannot cause a false "VALID"), caller-controlled input, and is documented in the file header's
  "Deviations" section — hence Info, not a defect.
- Recommendation: No action required for parity. If hardening is desired, either (a) require the
  store to contain exactly one key entry + its certificate, or (b) match the private key's public
  key against each candidate certificate and select the matching one before embedding.

Cleared as **parity-faithful / deliberate** (checked against Java, NOT findings):
- `trust_anchor_verifier.go` `IsTrustedAtTime` / `IsTrustedCertificateChain` /
  `isAcceptUntrustedCertificateChains` — match Java exactly; the "accept untrusted chains" branch
  is an opt-in flag defaulting to `false`, so no default false-trusted.
- `baseline_requirements_checker.go` `MinimalLTRequirement` / `MinimalTRequirement` /
  `MinimalLTARequirement` / `SignatureTimestampsCreatedBeforeSignCertExpiration` — match Java.
  (The LT gate's fail-open behaviour is T22D1B-SEC-001 / T22D1-SEC-001, not a checker defect.)
- `analyzer/default_document_analyzer.go` `GetValidationDataWithTimestamps` uses the **raw**
  `a.certificateVerifier` (not a copy) — matches Java `getValidationData`, which passes the raw
  field; `Validate()` uses the `buildCompleteCopyForValidation` copy, matching Java `validate()`.
- `validation_data_container.go` `ValidationDataForSignature` unchecked `Get` returning nil on a
  miss (Java returns `null`); the downstream `AddValidationData(nil)` nil-deref is parity-faithful
  to Java's `NullPointerException` and requires caller API misuse (querying a signature never added),
  not hostile document input.
- `validation_data.go` `AddToken` panicking on an unexpected token type — parity-faithful to Java's
  thrown exception.
- `signature_validation_alerter.go` `alert()` re-panicking `Alert(status)`'s error — the documented
  convention for Java's unchecked `AlertException` out of a `void` method (see file header).

### Tool log

Run from `dss/` (module root):

```
$ gofmt -l spi/validation
(no output)

$ go vet ./spi/validation/
(clean, exit 0)

$ golangci-lint run --config=../.github/.golangci.yml ./spi/validation/
0 issues.
```

No in-scope tool findings. The deliberately-disabled `staticcheck` `ST*`/`S1*`/`QF*`/`SA1019`
categories and legacy-algorithm flags (SHA-1/MD5/PKCS#1 v1.5 in the TSP source's accepted-algorithm
table and signing path) are not reported, per `PORTING.md` and the lint policy.

### Summary

2 findings — **1 High (SEC), 1 Info (STD)**; 0 Critical, 0 Medium, 0 Low, 0 PERF-new.
SEC: 1 (T22D1B-SEC-001, High — alerter fail-open on the live main validation path; same root as
T22D1-SEC-001). PERF: 0 new (O(n·m) shapes already covered by T22C-PERF-BIGO-001). STD: 1
(T22D1B-STD-001, Info — documented issuance-only PKCS#12 `chain[0]` assumption). All three tools
clean.

### Open questions

1. **T22D1B-SEC-001 ↔ T22D1-SEC-001 consolidation:** both are the `getOrderedCertificateChains()`
   error being swallowed at the `allRequiredRevocationDataPresent` / `allPOECoveredByRevocationData`
   level, surfaced once through the `signature_validation_context.go` `Check*` methods (T22D1) and
   once through the `signature_validation_alerter.go` `Assert*` methods (T22D1B). Confirm they are
   fixed together (a single fail-closed handling of the reorderer error) and that one oracle test
   covering an unorderable chain is added, rather than two.
2. **T22D1B-SEC-001 severity:** filed as **High** to match the T22D1-SEC-001 / T23-SEC-001
   precedent (specific-check false-accept, crafted input, other independent checks still run).
   If the project's bar for "false-accept in a verification/trust path" is Critical regardless of
   whether other checks still run, please confirm whether this and T22D1-SEC-001 should be
   re-rated Critical.

## spi/validation — scope/evidence/analyzer (unit U22d3)

### Scope

The **scope-determination, evidence-record, and analyzer subsystems** of `spi/validation` —
`scope/` (16 files: `EvidenceRecordScopeFinder`, the `TimestampScopeFinder` implementations, the
`SignatureScope` concrete types, `AbstractSignatureScopeFinder`), the evidence-record helpers
(`AbstractEmbeddedEvidenceRecordHelper`, `EmbeddedEvidenceRecordHelper`,
`AbstractSignatureEvidenceRecordDigestBuilder`, `EvidenceRecord`), the analyzer surface
(`DefaultDocumentAnalyzer`, `DocumentAnalyzer`/`DocumentAnalyzerFactory`,
`EvidenceRecordAnalyzer`/`Factory`, `eaa.PresentationAnalyzer`/`Factory`), and the EAA value
types (`EAA`, `EAAPresentation`, `EAAKeyBindingPayload`, `TimestampedReference`,
`ArchiveTimestampHashIndexStatus`). 1:1 port of Java DSS 6.5.RC1; a false "VALID" verdict is
Critical. Most of these files are pure interface declarations or small data-scope holders with no
executable trust logic; the findings below concern the two files that carry real scope/trust logic
(`evidence_record_scope_finder.go`, and the `SignatureScopeBase` base they embed).

### Date

2026-08-26

### Depth

Budgeted. `analyzer/default_document_analyzer.go` (924) was read header + `func` map + targeted
logic regions (validate / `GetValidationData` / `prepare*` / evidence-record-build / policy /
`SignatureByID`); the remaining 33 in-scope files were read in full. Every scope/trust-logic path
that could flip a verdict or panic was cross-checked against the upstream Java at `~/Workspace/esig`
(DSS 6.5.RC1 pin), including `EvidenceRecordScopeFinder.java`,
`AbstractEmbeddedEvidenceRecordHelper.java`, `DefaultDocumentAnalyzer.java`,
`DefaultSignaturePolicyValidatorLoader.java`, and `SignatureScope.java`.

### Files read (line counts actually read)

| File | Read / total | Mode |
|---|---|---|
| `analyzer/default_document_analyzer.go` | ~580 / 924 | budgeted (header, `func` map, validate/`GetValidationData`/`prepare*`/evidence-record/policy/`SignatureByID` regions; unread ~340 lines are trivial getters/setters + the `DefaultDocumentAnalyzerOverrides` interface declarations) |
| `scope/evidence_record_scope_finder.go` | 150 / 150 | full |
| `scope/signature_signature_scope.go` | 65 / 65 | full |
| `scope/signature_scope_with_transformations.go` | 59 / 59 | full |
| `scope/counter_signature_scope.go` | 56 / 56 | full |
| `scope/key_binding_signature_scope.go` | 55 / 55 | full |
| `scope/eaa_signature_scope.go` | 55 / 55 | full |
| `scope/manifest_signature_scope.go` | 52 / 52 | full |
| `scope/detached_timestamp_scope_finder.go` | 49 / 49 | full |
| `scope/encapsulated_timestamp_scope_finder.go` | 47 / 47 | full |
| `scope/abstract_signature_scope_finder.go` | 45 / 45 | full |
| `scope/container_signature_scope.go` | 40 / 40 | full |
| `scope/container_content_signature_scope.go` | 34 / 34 | full |
| `scope/digest_signature_scope.go` | 33 / 33 | full |
| `scope/evidence_record_timestamp_scope_finder.go` | 31 / 31 | full |
| `scope/evidence_record_master_signature_scope.go` | 29 / 29 | full |
| `scope/signature_scope_finder.go` | 17 / 17 | full |
| `scope/timestamp_scope_finder.go` | 14 / 14 | full |
| `abstract_embedded_evidence_record_helper.go` | 145 / 145 | full |
| `evidence_record.go` | 143 / 143 | full |
| `abstract_signature_evidence_record_digest_builder.go` | 140 / 140 | full |
| `embedded_evidence_record_helper.go` | 72 / 72 | full |
| `analyzer/document_analyzer.go` | 117 / 117 | full |
| `analyzer/document_analyzer_factory.go` | 39 / 39 | full |
| `analyzer/evidence_record_analyzer.go` | 38 / 38 | full |
| `analyzer/evidence_record_analyzer_factory.go` | 66 / 66 | full |
| `analyzer/eaa/eaa_presentation_analyzer.go` | 17 / 17 | full |
| `analyzer/eaa/eaa_presentation_analyzer_factory.go` | 76 / 76 | full |
| `archive_timestamp_hash_index_status.go` | 49 / 49 | full |
| `timestamped_reference.go` | 57 / 57 | full |
| `eaa.go` | 54 / 54 | full |
| `eaa_presentation.go` | 22 / 22 | full |
| `eaa_key_binding_payload.go` | 21 / 21 | full |
| `scope/full_signature_scope.go` (extra, in-scope dir) | 36 / 36 | full |
| `model/scope/signature_scope.go` (supporting base type, read for parity) | 220 / 220 | full (context — `SignatureScopeBase` is the shared base every scope above embeds) |

Total in-scope lines read: **~2,500** (at the 2,500 total budget cap — all 33 small files read in full plus the budgeted ~580-line slice of the 924-line `default_document_analyzer.go`). Plus ~220 lines of supporting context (`model/scope/signature_scope.go`, marked out-of-scope, read only to verify scope-equality parity) that does not count toward the in-scope budget. Unread `default_document_analyzer.go` regions are the trivial getters/setters and the `DefaultDocumentAnalyzerOverrides` interface declarations — no unverified trust-logic regions remain.

### Findings

#### SEC

T22D3-SEC-001 — **Medium** — Category: SEC — `dss/spi/validation/scope/evidence_record_scope_finder.go:40` (and `:41`)
- Evidence:
  ```go
  if f.IsSignatureEmbeddedAndValid(f.EvidenceRecord) {                                   // :40
      evidenceRecordScopes = f.enrichRecursively(evidenceRecordScopes,
          f.EvidenceRecord.MasterSignature().SignatureScopes())                          // :41
  }
  ```
  `IsSignatureEmbeddedAndValid` (`:68`) returns `true` when `IsEmbedded()` is `true` **and** some
  `ReferenceValidation` has type `EVIDENCE_RECORD_MASTER_SIGNATURE` and `IsIntact()`. It does **not**
  check that `MasterSignature()` is non-nil. `MasterSignature()` is an interface method on
  `validation.EvidenceRecord` (evidence_record.go) that is permitted to return a nil
  `AdvancedSignature`; `SignatureScopes()` on a nil interface then panics (nil pointer dereference).
- Impact: A nil-deref panic inside the embedded-evidence-record scope-determination path on a
  malformed/inconsistent `EvidenceRecord`. This is **fail-closed** (a panic aborts the whole
  validation — it cannot produce a false "VALID"), so it is not a false-accept and is not
  Critical. It mirrors upstream Java exactly (Java's `EvidenceRecordScopeFinder.findEvidenceRecordScope`
  calls `evidenceRecord.getMasterSignature().getSignatureScopes()` and throws the same
  `NullPointerException`), so it is parity, not a divergence. Reachability is currently low: no
  `EvidenceRecordAnalyzerFactory` is registered in this port (see
  `docs/compatibility/known-gaps.md` — "evidence record cannot actually be parsed from a document
  today"), so the `IsEmbedded()==true` branch is not exercised in the live flow. It becomes reachable
  the moment embedded-EAA / evidence-record parsing is wired up, and it is the one unguarded
  dereference in an otherwise nil-checked scope path.
- Recommendation: Guard the call — `masterSignature := f.EvidenceRecord.MasterSignature();
  if masterSignature != nil { … enrichRecursively(…, masterSignature.SignatureScopes()) }` — or
  document that a `MasterSignature()` returning nil while `IsEmbedded()` is `true` is a contract
  violation that callers must not produce. Add an oracle test that drives an embedded evidence record
  whose `MasterSignature()` is nil through `FindEvidenceRecordScope` and asserts a clean (non-panicking)
  result or a documented `model.DSSError`.

#### PERF

No PERF findings in this unit (0). The only loop on attacker-influenced data is the scope
deduplication in `evidence_record_scope_finder.go` (`containsSignatureScope` is O(n) per call, and
`enrichRecursively` is called per node), giving O(n·m) over the signature-scope tree — the same
small, bounded class already covered by **T22C-PERF-BIGO-001** and not re-filed. No unbounded growth,
no goroutines, no retained-reference or cache leaks in any of the 34 in-scope files.

#### STD

T22D3-STD-001 — **Low** — Category: STD — `dss/spi/validation/scope/evidence_record_scope_finder.go:136` (and `:148`)
- Evidence:
  ```go
  func documentsEqual(a, b model.DSSDocument) bool {
      switch av := a.(type) {
      case *model.InMemoryDocument: …
      case *model.DigestDocument:  …
      case *model.FileDocument:    …
      default:
          return a == b          // :148 — Go interface == (reflects DeepEqual)
      }
  }
  ```
- Impact: The `default` fallback compares two `model.DSSDocument` interface values with `==`. Go's
  `==` on interface values is legal only when the dynamic types are comparable; if a
  non-`model.*` `DSSDocument` implementation carries an incomparable field (slice, map, or func),
  `a == b` panics. In practice all live documents are `model.InMemoryDocument` / `model.DigestDocument`
  / `model.FileDocument` (each handled by an explicit case), so the fallback is effectively dead —
  hence Low. It mirrors Java's default `Object#equals` (reference equality) for custom
  `DSSDocument` subclasses, but Go's `==` is strictly narrower (it can panic where Java would merely
  return `false`), so the port is not a 1:1 behavioural match on that edge.
- Recommendation: Replace the `default` branch with a safe reference-equality check, e.g.
  `return fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b)` or an explicit `a == b` guarded by a
  `reflect.TypeOf(a) == reflect.TypeOf(b)` pre-check, so an unknown `DSSDocument` implementation
  cannot panic the scope deduplication.

### Deliberate-exclusion notes (reviewed, not findings)

- **`abstract_embedded_evidence_record_helper.go:139-143` `buildDigest` swallows the digest-builder
  error and returns an empty `model.Digest{}`.** This is a deliberate, documented port of Java's
  `buildDigest` `catch (Exception e) { LOG.warn(…); return new Digest(); }` — the file header states
  it explicitly ("preserves the swallow-and-return-empty-Digest behavior exactly"). Not re-reported.
- **`analyzer/default_document_analyzer.go:771-780` `buildEvidenceRecord` swallows the
  `EvidenceRecordAnalyzerFromDocument` error and returns `nil`.** Deliberate port of Java's
  `buildEvidenceRecord` `catch (UnsupportedOperationException/Exception) { LOG.warn(…); } return null;`
  and consistent with the known-gaps absence of evidence-record parsing. Not re-reported.
- **Global factory registries** (`documentAnalyzerFactoryRegistry`, `evidenceRecordAnalyzerFactoryRegistry`,
  `eaaPresentationAnalyzerFactoryRegistry`) are unguarded package-level `[]…` mutated by `Register*`.
  All call sites are inside `func init()` (verified across `cades/`, `xades/`, `jades/`, `pades/`,
  `asic/`, `validation/timestamp/`), so they are fully written before `main` runs and reads happen
  only at runtime — no data race, and the registration-replaces-`ServiceLoader` order-matters contract
  is documented. Deliberate, not a finding.
- **`scope/signature_signature_scope.go` `DSSID()` no-op override** and the `SignatureScopeBase.DSSID`
  lazy build that panics on a read failure: documented as parity with Java's unchecked
  `getDSSId()` and as a fail-closed (not false-accept) choice. Not re-reported.
- The deliberately-disabled `staticcheck` `ST*`/`S1*`/`QF*`/`SA1019` categories and legacy-algorithm
  flags are not reported, per `PORTING.md` and the lint policy.

### Tool log

All run from `dss/` (the Go module root), 2026-08-26.

```console
$ gofmt -l spi/validation
(no output)

$ go vet ./spi/validation/
(no output)

$ golangci-lint run --config=../.github/.golangci.yml ./spi/validation/
0 issues.
```

All three tools clean; no in-scope tool findings.

### Summary

2 findings — **0 Critical, 0 High, 1 Medium (SEC), 1 Low (STD)**; 0 PERF.
SEC: 1 (T22D3-SEC-001, Medium — unguarded `MasterSignature()` dereference in the
embedded-evidence-record scope path; parity with Java, fail-closed, currently low-reachability).
PERF: 0 (O(n·m) scope dedup already covered by T22C-PERF-BIGO-001). STD: 1
(T22D3-STD-001, Low — `documentsEqual` `==` interface fallback can panic on an unknown
`DSSDocument` implementation). All three tools clean.

### Open questions

1. **T22D3-SEC-001 severity arbitration:** filed as **Medium** (hardening opportunity; the panic is
   fail-closed — it cannot yield a false "VALID" — and the path is not reachable in the live flow
   today because no `EvidenceRecordAnalyzerFactory` is registered). If the project's bar treats any
   nil-deref panic in a trust/scope path as **High** ("panic on hostile input" in the rubric)
   regardless of current reachability, please confirm re-rating, and confirm the guard should be
   added proactively ahead of the evidence-record-parsing gap closing.
2. **T22D3-STD-001:** confirm the `default: return a == b` fallback is acceptable as-is (parity with
   Java reference equality) or should be made panic-safe per the recommendation. Low either way.
3. **Consolidation with the evidence-record gap:** T22D3-SEC-001 and the `docs/compatibility/known-gaps.md`
   "Evidence-record analyzers — not ported" entry are the same subsystem. Confirm they are tracked
   together so the guard (question 1) lands in the same change that adds the first
   `EvidenceRecordAnalyzerFactory`, rather than as two independent items.

## spi/validation — timestamp (unit U22d2)

### Scope

The **timestamp subsystem** of `spi/validation`: timestamp sources
(`signature_timestamp_source.go`, `abstract_timestamp_source.go`, `detached_timestamp_source.go`,
`timestamp_source.go`), the RFC 3161/CMS **timestamp token** and its verification
(`timestamp_token.go`), token verification/comparison/identification
(`timestamp_token_verifier.go`, `timestamp_token_comparator.go`,
`timestamp_identifier_builder.go`, `timestamp/signature_timestamp_identifier_builder.go`,
`timestamp_message_digest_builder.go`, `timestamp_token_identifier.go`), the per-token
certificate/CRL/OCSP sources (`timestamp_certificate_source.go`, `timestamp_crl_source.go`,
`timestamp_ocsp_source.go`), and the timestamp analyzer surface
(`analyzer/timestamp/timestamp_analyzer.go`, `analyzer/timestamp/timestamp_analyzer_comparator.go`).
This is the CMS/PKCS#7 (RFC 3161) timestamp path — *not* the JAdES/JSON timestamp path (U17b,
`batch-05`). 1:1 port of Java DSS 6.5.RC1; a false "VALID" timestamp verdict is Critical.

### Date

2026-08-26

### Depth

Budgeted, but the two largest files were read **fully** (not header + ≤3 regions) because they
carry the entire trust logic (`CheckIsSignedByToken` → `timestampTokenValidate` →
`timestampTokenVerifySignerInfo` in `timestamp_token.go`, and the token-creation / reference /
archive-timestamp flow in `signature_timestamp_source.go`). Every verdict-affecting line was
cross-checked against the upstream Java at `~/Workspace/esig` (DSS 6.5.RC1 pin) and against
`internal/cmscore` (the BouncyCastle replacement). The 15 smaller files were read in full.

### Files read (line counts actually read)

| File | Read / total | Mode |
|---|---|---|
| `timestamp/signature_timestamp_source.go` | 1252 / 1252 | full (exceeded the ≤400 soft budget — carries the token-creation + reference/archive-timestamp logic) |
| `timestamp_token.go` | 1199 / 1199 | full (exceeded the ≤400 soft budget — carries the RFC 3161 / CMS verification + state machine) |
| `timestamp/abstract_timestamp_source.go` | 532 / 532 | full |
| `timestamp/detached_timestamp_source.go` | 153 / 153 | full |
| `timestamp/signature_timestamp_identifier_builder.go` | 123 / 123 | full |
| `timestamp/timestamp_message_digest_builder.go` | 33 / 33 | full (interface only — no logic) |
| `timestamp_token_verifier.go` | 166 / 166 | full |
| `timestamp_token_comparator.go` | 98 / 98 | full |
| `timestamp_identifier_builder.go` | 175 / 175 | full |
| `timestamp_source.go` | 101 / 101 | full (interface only — no logic) |
| `timestamp_certificate_source.go` | 40 / 40 | full |
| `timestamp_ocsp_source.go` | 36 / 36 | full |
| `timestamp_crl_source.go` | 33 / 33 | full |
| `timestamp_include.go` | 43 / 43 | full |
| `timestamp_token_identifier.go` | 41 / 41 | full |
| `analyzer/timestamp/timestamp_analyzer.go` | 16 / 16 | full (interface only — no logic) |
| `analyzer/timestamp/timestamp_analyzer_comparator.go` | 32 / 32 | full |

Total in-scope lines read: **3,925**. The two budgeted files were read beyond their per-file caps
so the trust path could be verified end-to-end (see Depth); the remaining 15 files are all within
budget and the aggregate exceeds the 2,900-line guideline only because of those two.

### Security posture (top priority) — VERIFIED CLEAN on the verification path

No false-"VALID" timestamp verdict is reachable. Verified against the BouncyCastle/DSS upstream:

- **`IsSignatureIntact()` fails closed.** `model/token.go:103` initializes
  `signatureValidity` to `SignatureValidityNotEvaluated` and `model/token.go:208-209`
  returns `SignatureValidityValid == t.signatureValidity`, so an unverified token is
  *not* intact. `TimestampToken.IsValid()` (`timestamp_token.go:620`) additionally requires
  `IsMessageImprintDataFound()` and `IsMessageImprintDataIntact()`; both start false/nil and
  are only set by a `matchData*` call that ran `bytes.Equal` against the embedded imprint.
  `AreReferenceValidationsValid()` (`:1008`) fails closed on any found-but-not-intact
  reference validation. **All four conjuncts must be true for "VALID".**
- **Digest comparison is constant-content `bytes.Equal`** — `MatchDataSuppressingWarnings`
  (`timestamp_token.go:836`) compares the caller-supplied digest to `MessageImprint().Value()`
  with `bytes.Equal`; `timestampTokenVerifySignerInfo` (`:290`) compares the CMS message-digest
  attribute to the recomputed digest the same way. No length-prefix or early-exit weakness.
- **Certificate selection is fail-closed.** `timestampTokenValidate` (`:314`) requires the
  ESSCertID `certHash` to equal the candidate's digest **and** (when `issuerSerial` present)
  the serial + issuer name to match, then `timestampTokenValidateCertificate` (`:156`) requires
  the EKU to be *solely* `id-kp-timeStamping` and critical, then the certificate must be valid at
  `TSTInfo().GenTime`, and finally the signature must verify via
  `spi.SignerInformationVerifier`. Any failure returns an error → `SignatureValidityInvalid`.
- **The `isValidCMSSignedData` fallback** (`timestamp_token.go:728`, reached from
  `CheckIsSignedByToken` only if `isValidTimestamp` failed) still requires the cryptographic
  signature to verify against the *candidate* certificate's public key
  (`timestampTokenVerifySignerInfo` → `verifier.Verify`), and only then sets `VALID`. A
  token whose TSA cert does not match the candidate cannot reach `VALID` through either branch.
- **`TimestampTokenVerifier.IsAcceptableWithChainAt`** (`timestamp_token_verifier.go:118`) is a
  conjunction of trusted-anchor, cryptographic-validity, and chain-validity, each fail-closed
  (`containsTrustAnchor` returns false with no `TrustAnchorVerifier`; `isCryptographicallyValid`
  requires both imprint-intact and signature-intact). `isCertificateChainValid` returns `true`
  only when a `RevocationDataVerifier` is *absent* (matches Java; an absent verifier means the
  caller opted out of chain revocation checking) — not a false-trust.
- **`DigestAlgorithm()` panic** (`timestamp_token.go:876`) on an unknown message-imprint OID and
  the **`IsSignedBy`/`CheckIsSignedBy(publicKey)` panics** (`:654`, `:743`) are *documented
  deviations* (file header "Deviations") reproducing Java's unchecked `IllegalArgumentException` /
  `UnsupportedOperationException`; they are fail-closed (panic, not a silent valid) and are NOT
  findings.

### Findings

#### SEC

T22D2-SEC-001 — **Medium** — Category: SEC — `dss/spi/validation/timestamp_token.go:226`
- Evidence:
  ```go
  signingTime, err := timestampTokenSignedAttributeValue(signerInfo, cmscore.OIDSigningTime, "signing-time")
  if signingTime != nil {
      if date := spi.DSSASN1UtilsDate(signingTime.Encoded()); !date.IsZero() && !candidate.IsValidOn(date) {
          return false, timestampTokenCMSError("CMSVerifierCertificateNotValidException",
              "verifier not valid at signingTime")
      }
  }
  ```
- Impact: The signing-certificate-validity-at-`signingTime` pre-check is **skipped entirely when
  the `signingTime` attribute's date cannot be parsed** (`DSSASN1UtilsDate` returns the zero
  `time.Time`, `dss/spi/dss_asn1_utils.go:735` — mirroring Java's `getDate` returning `null`).
  A crafted token with a malformed `signingTime` therefore bypasses this specific validity
  window check. This is a **faithful port of the upstream** (Java/BouncyCastle also skip the
  check on an unparseable date) and is **fail-closed overall** — the cryptographic signature
  verification (`verifier.Verify`) and the message-imprint match still run and still gate the
  verdict, so it **cannot produce a false "VALID"** on its own. Filed **Medium** as a hardening
  opportunity (a hostile document can suppress one validity-window signal) rather than a
  trust-logic defect.
- Recommendation: No parity-mandated change (it matches upstream). If hardening is desired, treat
  a *present-but-unparseable* `signingTime` as a fail-closed error (return
  `CMSVerifierCertificateNotValidException`) rather than silently skipping the check — but only
  if the project accepts a deliberate divergence from upstream, per `PORTING.md`'s
  upstream-tracking rule (needs a `// DIVERGENCE, deliberate:` + DESIGN entry).

No other SEC findings. The verification, certificate-selection, digest, and status-mapping
paths were verified clean (see Security posture above).

#### PERF

No PERF findings in this unit (0). The O(n·m) / O(n²) shapes present are bounded by
document-internal collection sizes (timestamp count, per-token `TimestampedReferences`,
evidence-record refs) and are the same class already covered by **T22C-PERF-BIGO-001**, so they
are not re-filed:
- `abstract_timestamp_source.go:140-150` `addReferences` and `:509-517`
  `containsEqualReference` are O(n·m) linear dedup scans; bounded by per-document reference
  lists.
- `signature_timestamp_source.go:452-460` `timestampTokenSliceSortStable` is an insertion sort
  (O(n²)); `n` is the (small) number of archive/document/detached timestamps of one signature,
  and it is a faithful port of Java's `Collections.sort` (stable).
- `containsTimestampsCoveringOtherTimestamps` (`:462`) is an O(n·refs) scan; bounded.
No goroutines, no global mutable state, no unbounded growth, and no retained-reference or cache
leaks in any of the 17 files. `TimestampToken`'s lazy fields
(`candidatesForSigningCertificate`, `messageImprint`, `identifierBuilder`) are computed-once
caches, not leak-prone.

#### STD

T22D2-STD-001 — **Low** — Category: STD — `dss/spi/validation/timestamp/abstract_timestamp_source.go:442`
- Evidence:
  ```go
  // Java casts unconditionally: `(OCSPResponseBinary) token`. Every concrete
  // EncapsulatedRevocationTokenIdentifier[OCSP] in this port is a *spi.OCSPResponseBinary,
  // so this assertion always succeeds in practice; a genuine mismatch panics, exactly
  // as Java's ClassCastException would.
  ocspResponseBinary := token.(*spi.OCSPResponseBinary)
  ```
- Impact: Unchecked type assertion that **panics** on a genuine mismatch. The sibling
  `CreateReferencesForOCSPBinaries` (`:415`) deliberately uses the safe `any(x).(*T)` form to
  *skip* non-matching elements (port of Java's `instanceof`-guarded cast), whereas this one
  reproduces Java's unconditional cast. It is a **faithful port** of Java's
  `ClassCastException`, only reachable with a malformed OCSP-binary source (caller/deserialization
  error, not hostile document input in the normal flow), and the file header documents the
  choice — hence Low, not a defect.
- Recommendation: No action required for parity. If hardening is desired, use the two-value
  form and return an error (matching the `CreateReferencesForOCSPBinaries` style) instead of
  panicking — again only if a deliberate divergence from upstream is acceptable.

Cleared as **parity-faithful / deliberate** (checked against Java, NOT findings):
- `timestamp_token.go` `DigestAlgorithm()` panic, `IsSignedBy`/`CheckIsSignedBy(publicKey)`
  panics, `IsMessageImprintDataIntact` panic — all documented deviations reproducing Java's
  unchecked throws; fail-closed.
- `timestamp_token_verifier.go` `isCertificateChainValid` reading `v.revocationDataVerifier`
  (field) for its nil-check while `isTrustedTimestampToken` uses the `TrustAnchorVerifier()`
  getter — faithful to Java's `getRevocationDataVerifier()`/`isCertificateChainValid` split; the
  field read is what makes the "absent verifier ⇒ skip chain check" branch correct.
- `signature_timestamp_source.go` `NewSignatureTimestampSourceBase` / `NewTimestampIdentifierBuilderBase`
  nil-panics — port of Java `Objects.requireNonNull`.
- `signature_timestamp_source.go:44` `DocumentTimestamps()` returning `nil` (vs `[]T{}`) for the
  non-PAdES base — additive accessor; PAdES overrides it through the `overrides` interface, so no
  behavioral difference.
- `timestamp_token_comparator.go` / `analyzer/timestamp/timestamp_analyzer_comparator.go`
  ordering logic — faithful port of the Java comparators; ordering-only, never verdict-affecting.
- `timestamp_certificate_source.go` / `timestamp_crl_source.go` / `timestamp_ocsp_source.go`
  construction — thin wrappers over `spi.CMS*Source`; `DSSASN1UtilsFirstSignerInformation` returns
  nil on an empty store, but every caller has already established `len(SignerInfos()) == 1`
  (enforced in `NewTimestampTokenFromTimeStampTokenWithIdentifierBuilder`), so no nil-deref.

### Tool log

Run from `dss/` (module root):

```
$ gofmt -l spi/validation
(no output)

$ go vet ./spi/validation/
(clean, exit 0)

$ golangci-lint run --config=../.github/.golangci.yml ./spi/validation/
0 issues.
```

No in-scope tool findings. The deliberately-disabled `staticcheck` `ST*`/`S1*`/`QF*`/`SA1019`
categories and legacy-algorithm flags are not reported, per `PORTING.md` and the lint policy.

### Summary

2 findings — **0 Critical, 0 High, 1 Medium (SEC), 1 Low (STD)**; 0 PERF.
SEC: 1 (T22D2-SEC-001, Medium — signing-time validity-window check skipped on an unparseable
`signingTime`; faithful to upstream, fail-closed, no false-"VALID"). PERF: 0 new (O(n·m)/O(n²)
shapes already covered by T22C-PERF-BIGO-001). STD: 1 (T22D2-STD-001, Low — documented
unconditional OCSP type assertion). The timestamp verification path itself was verified clean —
no false-"VALID" route. All three tools clean.

### Open questions

1. **T22D2-SEC-001 severity / hardening:** filed as **Medium** (hardening opportunity; the skipped
   check cannot flip a verdict to "VALID" because signature + imprint checks still gate it).
   Confirm this rating, and whether to add the fail-closed "present-but-unparseable ⇒ error"
   behaviour as a *deliberate divergence* (requires a `// DIVERGENCE, deliberate:` + DESIGN entry
   per `PORTING.md`) or to leave it as a faithful port of upstream.
2. **T22D2-STD-001:** confirm the unconditional `token.(*spi.OCSPResponseBinary)` panic is
   acceptable as-is (parity with Java's `ClassCastException`) or should be made error-returning
   per the recommendation. Low either way.
3. **Cross-unit consolidation:** T22D2's O(n·m) dedup shapes share the T22C-PERF-BIGO-001 class;
   confirm they are tracked under that single item and not re-filed per unit.

---

## spi/validation — remainder (unit U22d4)

### Scope

The `spi/validation` files not covered by units U22d1a/d1b/d2/d3: the validation-context
interface, the revocation-data loading strategies and their factories, the TSP sources, the
executor strategy set, the signature-identifier builder machinery, the EAA revocation-token /
status-token-identifier types, the status types' remaining surface, and the misc helpers
(byte-array comparator, file-name builder, timestamp include, TLS loader, marker interfaces).
Subpackages `timestamp/`, `scope/`, `analyzer/` are excluded (reviewed in U22d2 / U22d3). This is
a 1:1 port of Java DSS 6.5.RC1; a false "VALID" verdict is Critical.

### Date

2026-08-26

### Depth

Full read of all 31 in-scope files (each ≤273 lines), plus verification of the load-bearing
supporting code they drive into: `spi/revocation_source.go`, `spi/offline_revocation_source.go`,
`spi/offline_crl_source.go`, `spi/offline_ocsp_source.go`, `spi/list_revocation_source.go`,
`spi/repository_revocation_source.go`, `model/token.go`, `model/token_identifier.go`,
`model/multiple_digest_identifier.go`, `model/dss_exception.go`,
`spi/exception/dss_external_resource_exception.go`, `alert/object_status.go`, `utils/strings.go`,
`enumerations/eaa_status.go`, `enumerations/mime_type.go`, and the revocation-strategy wiring in
`signature_validation_context.go` (`getRevocationToken`, `getRevocationDataVerifier`,
`checkRevocationForCertificateChainAgainstBestSignatureTime`). The `recover()`/`errors.As`
interplay in `checkCRL`/`checkOCSP` was verified empirically with a throwaway program using the
exact `model.DSSError` / `exception.DSSExternalResourceException` shapes — no repository change.

### Files read (31)

`validation_context.go`, `revocation_data_loading_strategy.go`,
`crl_first_revocation_data_loading_strategy.go`, `ocsp_first_revocation_data_loading_strategy.go`,
`crl_first_revocation_data_loading_strategy_factory.go`,
`ocsp_first_revocation_data_loading_strategy_factory.go`, `revocation_data_loading_strategy_factory.go`,
`composite_tsp_source.go`, `tsp_source.go`, `executor/complete_validation_context_executor.go`,
`executor/default_validation_context_executor.go`, `executor/skip_validation_context_executor.go`,
`executor/validation_context_executor.go`, `abstract_signature_identifier_builder.go`,
`signature_identifier.go`, `signature_identifier_builder.go`, `identifier/signature_attribute_identifier.go`,
`signature_properties.go`, `signature_attribute.go`, `eaa_revocation_token.go`,
`eaa_revocation_token_binary.go`, `eaa_status_token_identifier.go`, `eaa_validation_parameters.go`,
`token_status.go`, `signature_status.go`, `revocation_freshness_status.go`,
`byte_array_comparator.go`, `tls/tls_certificate_loader.go`, `file_name_builder.go`,
`timestamp_include.go`, `doc.go`.

### Findings

**SEC**

T22D4-SEC-001 — **Medium** — Category: SEC — `dss/spi/validation/revocation_data_loading_strategy.go:120-131,145-156`
- Evidence:
  - `checkCRL`/`checkOCSP` each wrap the source call in `defer func(){ if r := recover(); r != nil { var dssErr *model.DSSError; if err, ok := r.(error); ok && errors.As(err, &dssErr) { revocationToken = nil; return }; panic(r) } }()`
  - the header comment states this ports Java's `try/catch(DSSException)` around the source call
  - but `errors.As(err, &dssErr)` for `*model.DSSError` matches only the base type (and its unwrap chain), **not** the `spi/exception` DSSException *subclasses* that merely **embed** `*model.DSSError` (e.g. `exception.DSSExternalResourceException`)
  - empirically confirmed: a panic of `*DSSExternalResourceException{DSSError: &DSSError{...}}` is **not** caught (re-panics), whereas Java's `catch(DSSException)` would swallow it and fall through to the next source
- Impact: the Go `recover()` is **narrower** than the Java `catch(DSSException)` it intends to mirror — a DSSException-subclass panic propagates where Java sets the token to nil and tries the next source. Direction is conservative (fail-closed, never a false-"VALID"); a reachable instance would be a **false "INVALID"**. Currently low-reachability: the concrete CRL/OCSP revocation sources I traced panic with a plain `*model.DSSError` (caught) or a plain non-DSS error (Java would also not catch a non-DSSException), so the gap is latent rather than live on the known revocation path.
- Recommendation: Either broaden the catch to the whole DSSException family (there is no common interface — a `model.IsDSSException(err) bool` helper or an explicit type switch over the `spi/exception` family), or add a `// DIVERGENCE, deliberate:` comment (PORTING.md / DESIGN entry) recording that the Go recover is narrower than Java's `catch(DSSException)`. Re-rate **High** if consolidation finds a reachable DSSException-subclass panic on the revocation path.

T22D4-SEC-002 — **Info** — Category: SEC — `dss/spi/validation/eaa_revocation_token.go:198-200`
- Evidence: `func (t *EAARevocationToken) Encoded() []byte { return t.encoded.Binaries() }` — `t.encoded` is `*EAARevocationTokenBinary`; if nil, the promoted `Binaries()` call nil-derefs. `CreationDate()` (:145) and `CheckIsSignedBy()` (:158) likewise `panic` (deliberate — Java abstract / `UnsupportedOperationException`, documented in the file header, not findings).
- Impact: a nil-deref on a token constructed without `encoded`. Faithful to Java's `getEncoded()` NPE; the EAA block is `eaa`-build-tag gated and carries no behavioural oracle corpus (known-gaps.md), so it is unproven either way. No verdict impact on the default build.
- Recommendation: No action required for the default build; if the EAA corpus is built, confirm the concrete `EAAStatusListToken` (out of scope here) always sets `encoded` before `Encoded()` is reachable, or make `Encoded()` nil-safe and document it.

Also checked, **not** a finding: `isAcceptableToken` returns `true` when `revocationDataVerifier == nil` (:161), but the live wiring (`signature_validation_context.getRevocationToken` → `setRevocationDataVerifier(c.getRevocationDataVerifier())`, and `getRevocationDataVerifier` lazily creates a `DefaultRevocationDataVerifier`) guarantees a non-nil verifier on every real path, so the accept-all branch is dead in practice — no fail-open. `EAAStatus`'s zero value is `""` (not `"VALID"`), so an unset EAA status reads as unknown/invalid, not valid — fail-closed.

**PERF**

No PERF-BIGO, PERF-MEM, or PERF-LEAK findings in these 31 files. The only O(n²)+ shapes in the
package are the already-filed `T22C-PERF-BIGO-001` / `T22D1A-PERF-BIGO-001` (re-scan of the
processed-certificate set) — I verified none of the remainder files adds a new attacker-influenced
super-linear pass: `RelatedTokens()` sorts O(n log n) over a small related-token set;
`MinimalNextUpdateTime()` is a single O(n) scan; `WriteString`/`PositionId`/`Build` are O(len);
`CompositeTSPSource.TimeStampResponse` is O(source-count); `ByteArrayComparatorCompare` is
O(min(len₁,len₂)). All caches (`digestMap`, `tspSources`/`sourceOrder`, the status maps) are
bounded by a single validation request's token set and not retained across requests — no leak.

**STD**

T22D4-STD-001 — **Low** — Category: STD — `dss/spi/validation/doc.go:3-4`
- Evidence: `per-format AdvancedSignature implementation (Signature,\n// Signature, ...) and evidence-record validator is built from:` — the parenthetical lists "Signature" twice (a copy-paste artifact; the second was presumably meant to be a different format type).
- Impact: package doc is slightly misleading; no behavioural effect.
- Recommendation: Fix the duplicated token in the package doc (e.g. name the actual per-format types, or drop the parenthetical).

T22D4-STD-002 — **Info** — Category: STD — `dss/spi/validation/byte_array_comparator.go:14-25`
- Evidence: `ByteArrayComparatorCompare` hand-rolls an unsigned lexicographic byte compare with a length fallback for the prefix case.
- Impact: none; the logic is correct and matches BouncyCastle's `ByteArrayComparator` it is ported from.
- Recommendation: Note only — this is semantically identical to the stdlib `bytes.Compare` (unsigned, length-tie-broken), so it *could* be collapsed to `bytes.Compare(o1, o2)`. Flagging so a reviewer does not mistake it for a bug, but per PORTING.md's 1:1-mapping rule (and the provenance comment) it is a deliberate port and should **not** be "simplified" without a DESIGN note.

### Tool log (run from `dss/`)

```
$ gofmt -l spi/validation            → (no output; clean)
$ go vet ./spi/validation/           → (no output; clean)
$ golangci-lint run --config=../.github/.golangci.yml ./spi/validation/ → 0 issues.
```

All three in-scope tools are clean; no tool-derived findings, no out-of-scope tool reports.

### Summary

- **Total findings: 4** (SEC: 2, PERF: 0, STD: 2).
- **Severity breakdown:** Critical: 0, High: 0, Medium: 1, Low: 1, Info: 2 — i.e. SEC 1 Medium (T22D4-SEC-001) + 1 Info (T22D4-SEC-002); PERF none; STD 1 Low (T22D4-STD-001) + 1 Info (T22D4-STD-002).
- **No Critical, no High, no PERF findings, no false-"VALID" route** in these 31 files. The single material item is **T22D4-SEC-001**: the `checkCRL`/`checkOCSP` `recover()` catches only `*model.DSSError`, not the DSSException *subclasses* Java's `catch(DSSException)` would — a conservative (fail-closed) semantic gap that is currently latent on the known revocation path but is a real porting mismatch worth a fix or a `// DIVERGENCE` note. Everything else is a faithful-to-Java nil-deref in the unproven EAA block, a package-doc typo, or a "don't mistake for a bug" note.

### Open questions

1. **T22D4-SEC-001 (fix vs. document + severity):** the Go `recover()` in `checkCRL`/`checkOCSP` is narrower than Java's `catch(DSSException)`. Confirm whether to (a) broaden the catch to the full DSSException family (needs a family helper — there is no common interface), or (b) record it as a deliberate divergence. And confirm the **Medium** rating holds given the currently-unreachable DSSException-subclass panic on the revocation path — re-rate High if a reachable instance is found.
2. **Cross-reference:** T22D4-SEC-001 is a *different* root from the already-filed `T22D1A-SEC-001`/`T22D1B-SEC-001` (`getOrderedCertificateChains` fail-open) — this one is fail-*closed* and lives in the loading strategy, not the chain reordering. Confirm it is tracked separately and not folded into those.
3. **T22D4-SEC-002:** confirm the unproven-EAA treatment (known-gaps.md) is the right home for the `Encoded()` nil-deref rather than a standalone hardening item.
