# DSS Java → Go Port: Tech-Lead Roadmap

## Upstream baseline

- **Source**: [esig/dss](https://github.com/esig/dss), version **6.5.RC1**, commit `4c2129862948bfd53ca1455832260aa17e183cf8`
- **Scale**: 117 Maven modules, ~420K lines of main-source Java (2,683 files), ~3,800 test files
- A read-only reference clone lives outside this repo during development sessions (`/home/user/dss-upstream`); it is never committed here.

## What "100% compatibility" means

API-identical Java-in-Go is an anti-goal. The compatibility contract is **interoperability**, verified continuously against upstream artifacts:

1. **Byte/spec-level outputs** — signatures produced by the Go port validate in Java DSS and vice versa; canonicalization, digests, ASN.1 structures, and serialized enum values (names, OIDs, URIs) are exact.
2. **Verdict parity** — validation of the same document with the same policy yields the same conclusion (indication/sub-indication) as Java DSS.
3. **Schema parity** — diagnostic data, simple/detailed/ETSI validation reports conform to the same upstream XSD/JSON schemas.
4. **Test-vector parity** — upstream `src/test/resources` signed files, certificates, and policies are reused as golden vectors in `dss/testdata/`.

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

## Delegation matrix (sub-agent model assignment)

| Work type | Model | Rationale |
|---|---|---|
| Mechanical porting: enums, DTOs, JAXB models, builders, utils | Sonnet | High-volume, pattern-following; cheapest adequate tier |
| Crypto-sensitive core: ASN.1/CMS, C14N, XML-DSig, PDF ByteRange, signature builders | Opus | Correctness under adversarial subtlety |
| Architecture/API design, integration, plan revisions, final reviews | Fable (tech lead) | Session driver |
| Fidelity reviews (OID tables, enum values, verdict mapping) | Opus | Adversarial verification pass |
| Large-scale grep/inventory/reporting chores | Haiku/Sonnet | Trivial mechanical work |

Every batch ends with: `go build ./... && go vet ./... && go test ./...` green, an adversarial fidelity review, then commit + push.

## Progress (actuals)

| Phase | Status | Notes |
|---|---|---|
| 0, 1a, 1b | ✅ done | enumerations/alert/utils/model — audited vs OpenJDK 21 oracle |
| 2a | ✅ done | spi core + crlparser; native BER/DER engine + OCSP; audited vs maven-built upstream DSS |
| CMS core | ✅ done | internal/asn1ber + internal/cmscore (RFC 5652/3161), BC-oracle byte-identical builds |
| 2b | ✅ done | spi/validation (AdvancedSignature 86/86 methods), TimestampToken, token, document, CMS sources |
| 3 | ✅ done | cms + cades; CAdES-B signed attrs byte-identical to BC (36 KATs); **cross-validation GREEN both directions**: upstream Java DSS validates Go-signed CAdES-B/T; Go validates 21/21 upstream fixtures strictly. PFX parser + determinism + test backfill follow-ups closed |
| 4+ | pending | next: XML stack (native DOM/c14n design, then xmldsig + XAdES over two windows) |

Known accepted gaps (tracked): PKCS#12 Ed25519/DSA keystores unloadable (x/crypto limitation; native PFX parser planned), JKS/PKCS#11 unsupported, MD2/WHIRLPOOL digests, map-iteration ordering sweep pending at 12 sites, test backfill pending for validation/timestamp, document, analyzer, jdbc.

## Budget pacing (Claude Max 5-hour windows)

- One bounded batch per usage window (≈8–12 Sonnet agents + 1–2 Opus reviewers), then commit, push, and schedule a self check-in into the next window via `send_later`.
- Bulk porting always goes to Sonnet; Opus is reserved for the crypto-critical few percent and reviews.
- Progress is durable: every window ends pushed to `claude/dss-java-go-port-bte1eb`, so an interrupted window loses nothing.
