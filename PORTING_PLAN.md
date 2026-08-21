# DSS Java → Go Port: Tech-Lead Roadmap

## Upstream baseline

- **Source**: [esig/dss](https://github.com/esig/dss), version **6.5.RC1**, commit `4c2129862948bfd53ca1455832260aa17e183cf8`
- **Scale**: 117 Maven modules, ~420K lines of main-source Java (2,683 files), ~3,800 test files
- A read-only reference clone lives outside this repo during development sessions; it is never committed here.

## What "100% compatibility" means

API-identical Java-in-Go is an anti-goal. The compatibility contract is **interoperability**, verified continuously against upstream artifacts:

1. **Byte/spec-level outputs** — signatures produced by the Go port validate in Java DSS and vice versa; canonicalization, digests, ASN.1 structures, and serialized enum values (names, OIDs, URIs) are exact.
2. **Verdict parity** — validation of the same document with the same policy yields the same conclusion (indication/sub-indication) as Java DSS.
3. **Schema parity** — diagnostic data, simple/detailed/ETSI validation reports conform to the same upstream XSD/JSON schemas.
4. **Test-vector parity** — upstream `src/test/resources` signed files, certificates, and policies are reused as golden vectors. Small per-package fixtures live in-module under `dss/<package>/testdata/`; heavier oracle corpora live in the repo-root `corpus/` tree (outside the module, so `go get` stays lean) and are located at test time via `dss/internal/corpustest`, which skips the affected tests gracefully when `corpus/` is absent (e.g. a bare module checkout).

## Module mapping (Maven → Go package under `dss/`)

| Java module(s) | Go package | Notes |
|---|---|---|
| dss-enumerations | `enumerations` | Typed string constants; value == Java `name()`; OID/URI/registry lookups preserved exactly |
| dss-alert | `alert` | Alert/status abstractions |
| dss-utils (+ apache-commons/guava impls) | `utils` | Single stdlib-based implementation of the `IUtils` surface |
| dss-model | `model` | `DSSDocument` becomes `io.Reader`-based; digest/signature value objects |
| dss-spi, dss-crl-parser*, dss-token, dss-document | `spi`, `crlparser`, `token`, `document` | `crypto/x509` + `golang.org/x/crypto/cryptobyte` replace BouncyCastle |
| dss-cms, dss-cms-object, dss-cms-stream | `cms` | Native CMS SignedData build/parse |
| dss-cades | `cades` | B/T/LT/LTA |
| dss-xml-common, dss-xml-utils | `xmldsig`, `xmlc14n` | Byte-exact C14N 1.0/1.1/exclusive is the hard core |
| dss-xades | `xades` | B/T/LT/LTA |
| dss-pades (+ native PDF engine replacing pdfbox/openpdf) | `pades`, `pdf` | Incremental updates, ByteRange, DSS dictionary, doc timestamps |
| dss-jades | `jades` | ETSI TS 119 182 JOSE |
| dss-asic-* | `asic`, `asic/cades`, `asic/xades` | ASiC-S/E containers |
| dss-validation (+ *-report-jaxb, dss-policy-*) | `validation`, `validation/diagnostic`, `validation/reports`, `validation/policy` | EN 319 102-1 process; JAXB → `encoding/xml` against same XSDs |
| dss-tsl-validation, dss-validation-job | `tsl`, `tsljob` | EU LOTL/TL, pivots, refresh job |
| dss-i18n | `i18n` | Message bundles for report generation |

**Deferred (not in initial scope)**: REST/SOAP remote services and clients, EAA/mdoc modules, dss-cookbook, coverage/BOM modules, evidence-record modules (revisit after Phase 9).

## Phases (tracked in the session task list)

| Phase | Content | Size (main LOC) | Exit criteria |
|---|---|---|---|
| 0 | Scaffolding, roadmap, conventions | — | `go build ./...` green on empty module; docs merged |
| 1a | enumerations, alert, utils | ~15K | Registry unit tests pass; OID/URI fidelity review |
| 1b | model | ~25K | Value-object semantics tests |
| 2 | spi, token, crl-parser, document | ~46K | Parse upstream certs/CRLs/OCSP vectors |
| 3 | cms, cades | ~13K | Go-signed CAdES validates in Java DSS; upstream vectors parse |
| 4 | xmldsig/c14n, xades | ~32K | Byte-exact c14n vs upstream vectors; XAdES cross-validation |
| 5 | pades + pdf engine | ~21K+ | PAdES cross-validation on upstream PDFs |
| 6 | jades | ~14K | JAdES cross-validation |
| 7 | asic | ~12K | Container cross-validation |
| 8 | validation engine + reports | ~77K (6 sub-batches) | Verdict parity on upstream validation test corpus |
| 9 | tsl + validation job | ~12K | Parse production EU LOTL; equivalence tests |
| — | Compatibility harness | continuous | CI job cross-validating both directions |

## Working rules

- Mechanical porting (enums, DTOs, JAXB models, builders, utils) is high-volume, pattern-following work; the crypto-sensitive core (ASN.1/CMS, C14N, XML-DSig, PDF ByteRange, signature builders) gets the most scrutiny, as do fidelity reviews of OID tables, enum values and verdict mapping.
- Every batch ends with: `go build ./... && go vet ./... && go test ./...` green, an adversarial fidelity review, then commit + push.

## Progress (actuals)

| Phase | Status | Notes |
|---|---|---|
| 0, 1a, 1b | ✅ done | enumerations/alert/utils/model — audited vs OpenJDK 21 oracle |
| 2a | ✅ done | spi core + crlparser; native BER/DER engine + OCSP; audited vs maven-built upstream DSS |
| CMS core | ✅ done | internal/asn1ber + internal/cmscore (RFC 5652/3161), BC-oracle byte-identical builds |
| 2b | ✅ done | spi/validation (AdvancedSignature 86/86 methods), TimestampToken, token, document, CMS sources |
| 3 | ✅ done | cms + cades; CAdES-B signed attrs byte-identical to BC (36 KATs); **cross-validation GREEN both directions**: upstream Java DSS validates Go-signed CAdES-B/T; Go validates 21/21 upstream fixtures strictly. PFX parser + determinism + test backfill follow-ups closed |
| 4a-c | ✅ done | internal/xmldom + all 7 c14n variants (1829 Santuario byte-exact KATs); xpath10 subset (Xalan-parity); xml/common+utils; serializer Transformer-parity; internal/xmldsig (3421 oracle rows over 195 fixtures) |
| 4d | ✅ done | dss-xades ported; **XAdES cross-validation GREEN both directions** (60 fixtures/659 reference validations upstream→Go; Go-signed B/T validated by Java DSS); byte-parity on SignedInfo/SignedProperties (164 subtests). Audit fixes: v1-CRL parsing, Brainpool curves (internal/eccurve), manifest dependent validations |
| 5a | ✅ done | internal/pdf: native reader (xref tables/streams/hybrid, objstm, filters, RC4/AES encryption, revision extraction) + strictly-incremental writer (/Sig, ByteRange, /DSS+/VRI, DocTimeStamp, re-encryption). 267/267 corpus parity vs pdfbox 3.0.7 |
| 5b | ✅ done | dss-pades ported (110 files + gap-ported eu.europa.esig.dss.pdf package); **PAdES cross-validation GREEN both directions** (36 fixtures upstream→Go; Go-signed B/T/LT/LTA validated by Java DSS incl. /DSS semantic equivalence with Java-extended documents). Broad differential run over the full 248-PDF upstream corpus: **shadow-attack/modification-detection parity 0 mismatches across 298 signatures**. Audit fixes incl. CMS message-digest gate (security: forged /ByteRange no longer verifies), PDFDocEncoding, negative-serial + large-exponent (RFC 8017 manual path) certificate tolerance, virtual-dispatch corrections in timestamp sources |
| 6 | ✅ done | dss-jades + specs-jades ported on native internal/jose (jose4j replacement; ordered-JSON writer differential-fuzzed vs jose4j 0.9.6 — 0 byte diffs over 3000 documents in both parse modes incl. Java HashMap bucket-order emulation); **JAdES cross-validation GREEN both directions** (30 fixtures + sign-KAT goldens byte-identical to Java oracle; Go-signed B/T validated live by upstream). Broad differential run over the entire 63-file upstream corpus: 62/63 byte-identical (the one difference is an identical rejection with different wording). Audit fixes: signature-timestamp KeyInfo-cert coverage, content-timestamp DSS-Id attribute-order identity |
| 7 | ✅ done | dss-asic-{common,cades,xades} ported (123 files, zero-collision flatten into asic + asic/cades + asic/xades); **container cross-validation GREEN both directions**: Go-built ASiC-S/E × CAdES/XAdES at B/T validated live by Java DSS (ALL OK 8/8, mimetype first+STORED verified); Java-built containers extended by Go and re-validated by Java (extension round-trip); 189-fixture broad corpus — extraction/manifests/merge/per-signature analysis all 0 mismatches (16290 merge pairs); zip-bomb guards mutation-tested. Audit: un-gated the XAdES per-signature analyzer (gate was not import-forced), DSSUtilsEncodeURI now java.net.URI-exact (809 KATs); 28-method virtual-dispatch sweep clean |
| 8a | ✅ done | validation-model foundation: dss/i18n (1034 MessageTags, all formatting byte-identical to java.text.MessageFormat), dss/policy{,/jaxb,/crypto/json,/crypto/xml}, dss/diagnostic{,/jaxb}. **Marshal parity proven**: 889/889 diagnostic dumps + 32/32 policies + 480 crypto-suite keys byte-identical through unmarshal→remarshal; XSD completeness 100% (476 elements/74 attrs/160 types); 97,603-value wrapper differential — 0 real divergences. Audit caught and reversed a doctored oracle (porter had edited 4 policy reference files instead of fixing the Go marshaller) |
| 8b | ✅ done | report models: dss/simplereport, dss/simplecertificatereport, dss/detailedreport, dss/validationreport (ETSI TS 119 102-2) + validation/reports assembly. Marshal parity via PRODUCTION facade paths: 30+14+16+18 oracles all Java fixed points, raw-byte KATs (no normalization); XSD completeness 4/4 schemas zero dropped fields; DetailedReport wrapper 82/82 accessors, 0 behavioral mismatches over 2998-key differential. Audit recalibrated detailedreport oracles off a wrong (non-production) JAXB path. Deferred to 8c: 9 validation/reports/diagnostic builder files (unreferenced as of 8b) |
| 8c | ✅ done | check framework + BBB part 1: Chain/ChainItem conclusion machinery (14/14 mutation kills), isc/vci/cv/fc/sav check families (126/126 manifest files), diagnostic builders. **Check-semantics parity over 3,782 Java-oracle rows** — every fc/sav check has OK + NOT-OK rows; builder parity 45/49 certificate documents byte-identical. Audit fixes incl. critical ASiC enum-spelling failure (every ASiC signature failed FC), 17 spurious extension descriptions, map-order nondeterminism, field-lock check inversion |
| 8d | ✅ done | XCV certificate-chain validation (98 files: X509CertificateValidation + crs/rac revocation selection + 83 checks), AOV algorithm-obsolescence (31 files + cc crypto checkers), BasicBuildingBlocks moved to pkg blocks (cycle break, un-tagged, 474-row dispatcher KAT), process/eaa ported (40 files, mechanical 1:1 audit; behavioral corpus deferred to EAA phase). **Parity over ~38K oracle rows** (xcv 36K, aov 1.5K, blocks 474; all provenance-proven; expiry instants pinned at ±1ms). Audit fixes incl. Critical: AOV InitChain pointed at the LAST check so every AOV block ran exactly one check (557/1290 rows wrong, 19 conclusions flipped); plural MessageTag positions; null-vs-empty KeyLength/AdditionalInfo semantics |
| 8e | ✅ done | LTV/archival processes + qualification (240 files): vpfbs/vpftsp/vpfltvd/vpfswatsp (POE machinery: 61-row oracle incl. all comparator tie-breaker pairs + revoked-revived/timestamp-chain synthetics)/vpftspwatsp, EAA validation blocks, and the full qualification tree (trust-service filters, QSCD/type strategies, MRA, QWAC) with a 9,946-row oracle built by the audit. Audit Critical fix: nil-vs-empty sentinel made qualification PANIC on ordinary input (all trust services filtered out). Package splits for Go import cycles documented (eaa/checks, vpfswatsp/evidencerecord, vpfltvdsig). Three documented hash-order carve-outs (multiset-compared). Open test gaps (vpftsp/vpftspwatsp/vpfltvdsig/eaa blocks/evidencerecord + 2 self-asserting tests) BOUND to the 8f end-to-end harness |
| 8f | ✅ done | **PHASE 8 COMPLETE.** Executor + report builders + validation root/leaves; all 27 phase8-gated files un-gated (zero build tags remain); **end-to-end verdict parity GREEN with zero tolerances**: 273/273 executor-corpus rows (every BBB indication/sub-indication/qualification), 60/60 documents across six formats under the strict Id-inclusive comparator (all divergence allowlists deleted), **1,258/1,258 report hashes byte-identical** (simple/detailed/ETSI VR/cert flavors over the full corpus). Audit fixed 10 defects incl. pointer-identity attribute order (the tracked latent bug — at its spi source), CMS encapContentInfo in archive-timestamp-v2 imprints (intact LTA failed as HASH_FAILURE), three java.util.HashMap order emulations, IDREF subclass resolution, BigInteger octal parsing, JAXB error tolerance. Booked follow-ups: X500 toString rendering (D11), diagnostic-data cosmetic drift (D12–D14), QWAC-cert + EAA-presentation executor coverage |

Known accepted gaps (tracked): JKS/PKCS#11/Windows-MY/macOS-Keychain key stores unsupported (no Go counterpart under the dependency policy; PKCS#12 itself is fully supported for RSA/EC/Ed25519/DSA since the native `internal/pfx` RFC 7292 reader landed — the earlier "Ed25519/DSA unloadable" x/crypto limitation is closed, pinned by `token.TestKeyStoreLegacyKeyTypeFixtures`), no concrete evidence-record (RFC 4998 ERS / RFC 6283 XMLERS) analyzer — the framework and report plumbing are ported but nothing registers an `EvidenceRecordAnalyzerFactory`, no online CRL/OCSP/TSP HTTP sources (dss-service is out of scope; AIA fetching and the TSL downloader do use the native HTTP data loader), MD2/WHIRLPOOL digests, streaming CMS in-memory-only, visible-signature rasterization (pdf/visible drawers) returns not-supported, 6 CMS/CAdES-layer parity gaps surfaced by the 248-PDF PAdES differential run (booked against Phase 3: BadEncodedCMS signature-count, pdf-eof cert extraction, 4 legacy PKCS#7 reference-data cases, PLAIN-ECDSA BSI TR-03111, wrong-digest-algo, empty-vs-absent /Reason optionality — none is a false-accept; see dss/pades/testdata/broadgen/README.md), AOV signing-cert-ref grouping iterates first-seen order where Java uses java.util.HashMap bucket order (deterministic-order sanction; pinned by permutation-allowlist test; 2/1328 oracle rows, conclusions unaffected — revisit if 8f end-to-end parity trips), and process/eaa awaits a behavioral corpus (mechanical 1:1 verified; EAA phase post-9).

| 9 | ✅ done | **CORE PORT COMPLETE.** dss-tsl-validation + dss-validation-job + specs-trusted-list (227 manifest entries): TL/LOTL/pivot/MRA parsing, conditions, TL signature validation, the full validation-job framework with its cache state machine. **EU LOTL cross-validation GREEN A–D**: parse parity over 37 fixtures byte-exact (incl. pivot chains, 4 real-MRA LOTLs, v5/v6 TLs, Brainpool TL); TL signature verdicts 8/8; offline TLValidationJob parity (cache states, alert firings, full TrustedListsCertificateSource content); end-to-end trust hand-off into the Phase 8 engine incl. complete XmlTrustedList diagnostic blocks. Audit Critical: an XML subtree-decode bug silently discarded ALL otherCriteriaList conditions — empty ALL-composites matched unconditionally, granting TL qualifiers certificates must not receive; plus trust-source determinism (frozen spi/tsl), pointer-identity cert-token map keys, untested sync strategies (now a 16-row truth table). Mutation: 41/41 real mutants killed |

## Pacing

- One bounded batch at a time, each ending green, committed and pushed, so that an interrupted stretch of work loses nothing.
- Scrutiny is spent where it pays: the crypto-critical few percent and the review passes.
