# esig — digital signatures in pure Go

[![CI](https://github.com/ryftcore/dss-go/actions/workflows/ci.yml/badge.svg)](https://github.com/ryftcore/dss-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ryftcore/dss-go/dss.svg)](https://pkg.go.dev/github.com/ryftcore/dss-go/dss)
[![Go 1.27+](https://img.shields.io/badge/Go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![License: LGPL-2.1](https://img.shields.io/badge/License-LGPL--2.1-blue.svg)](LICENSE)

**esig** creates and validates AdES digital signatures — CAdES, XAdES, PAdES, JAdES
and ASiC containers — and runs the full ETSI EN 319 102-1 validation process,
including eIDAS qualification against the EU trusted lists.

It is a Go port of [esig/dss](https://github.com/esig/dss), the European
Commission's Digital Signature Services library, and it is built to
*interoperate* with it rather than to imitate its API: signatures produced here
validate in Java DSS, signatures produced there validate here, and validating
the same document yields the same indication, sub-indication and qualification —
checked continuously against a Java oracle corpus with no tolerance allowances.
[Compatibility](#compatibility) says exactly what was checked, in which
direction, at which levels.

- **Pure Go.** No cgo, no BouncyCastle, no PDFBox, no jose4j — every one of those
  is replaced by code in this repository.
- **Two dependencies.** `golang.org/x/crypto` and `golang.org/x/text`
  (`golang.org/x/sys` arrives transitively). Nothing else is compiled into your
  binary.
- **A library, a facade and a CLI.** ~314,000 lines of Go across 2,274 source
  files, behind a small top-level API — plus an `esig` command for the shell.

> **Status.** This project has not been released, tagged or independently
> security-audited, and it has no production track record. Read
> [What is not here](#what-is-not-here) before you depend on it, and read
> [SECURITY.md](SECURITY.md) before you report a problem with it.

## Install

```sh
go get github.com/ryftcore/dss-go/dss
```

Go 1.27 or newer. The CLI, once a version is tagged:

```sh
go install github.com/ryftcore/dss-go/dss/cmd/esig@latest
```

Until then, build it from a checkout — `cd dss && make cli-build` puts it in
`dss/bin/`.

## Quick start

Sign an XML document and validate the result — the whole round trip through the
top-level facade:

```go
package main

import (
	"fmt"
	"log"

	"github.com/ryftcore/dss-go/dss"
)

func main() {
	// A Signer is a key store connection plus the key entry to sign with;
	// the library never sees the private key, it asks the token to sign.
	signer, err := dss.OpenPKCS12("keystore.p12", "password")
	if err != nil {
		log.Fatal(err)
	}
	defer signer.Close()

	doc := dss.NewDocument("invoice.xml", []byte("<invoice><total>42</total></invoice>"))

	// Sign: an enveloped XAdES signature at baseline level B.
	signed, err := dss.Sign(doc, signer, dss.SignOptions{
		Format: dss.FormatXAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Validate: the full EN 319 102-1 process, against your trust anchors.
	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: []*dss.CertificateToken{signer.Certificate()},
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range reports.Verdicts() {
		fmt.Println(v.ID, v.Indication, v.SubIndication, v.SignatureLevel)
	}
}
```

```
S-5F32BDAC…20A0D995 TOTAL_PASSED  XAdES-BASELINE-B
```

`dss.Validate` also returns the four DSS reports themselves — SimpleReport,
DetailedReport, the ETSI TS 119 102-2 validation report and the diagnostic data —
as XML identical to what Java DSS emits.

The root [`dss`](https://pkg.go.dev/github.com/ryftcore/dss-go/dss) package is a thin
facade over the ported packages, covering the common paths and holding no logic
of its own. Everything it deliberately leaves out — counter-signatures, XAdES
references and transforms, signature policy stores, custom revocation sources,
trusted-list refresh jobs — is reached through the packages underneath, which are
fully exported and documented on [pkg.go.dev](https://pkg.go.dev/github.com/ryftcore/dss-go/dss).

Nine runnable, commented programs live in
[`dss/examples/`](dss/examples/README.md) — one story each, all but the
trusted-list one entirely offline.

## Command line

```sh
esig sign contract.pdf -format pades -level B \
    -p12 keystore.p12 -p12-pass env:P12_PASS -out contract-signed.pdf

esig validate contract-signed.pdf -trust ca.cer
```

```
S-CC2F8D5D…93A496C89  TOTAL_PASSED  level=PAdES-BASELINE-B  qualification=NA  signed-by="…"
```

| Command | What it does |
|---|---|
| `esig validate` | Validate a signed document; human summary by default, or `-format simple\|detailed\|diagnostic\|etsi-vr` for the DSS report XML |
| `esig sign` | Sign a document — `-format cades\|xades\|pades\|jades\|asice\|asics`, `-level B\|T\|LT\|LTA` |
| `esig extend` | Raise an existing signature to a higher level |
| `esig inspect` | Signature, time-stamp and certificate summary from the diagnostic data |
| `esig report` | Render a previously saved SimpleReport XML |
| `esig tl refresh` | Refresh a local trusted-list cache from the EU LOTL |
| `esig version` | Build information |

Exit codes: `0` success, `1` a signature did not reach TOTAL_PASSED, `2` usage
error, `3` runtime error. Every subcommand is a living example of the facade
above. Unlike the library, the CLI ships its own RFC 3161 HTTP client, so
`-tsa <url>` talks to a real time-stamping authority.

## Feature matrix

### Signature formats

| Format | Sign | Extend | Validate | Levels | Variants |
|---|:--:|:--:|:--:|---|---|
| **CAdES** (EN 319 122) | ✅ | ✅ | ✅ | B, T, LT, LTA | enveloping, detached; counter-signatures |
| **XAdES** (EN 319 132) | ✅ | ✅ | ✅ | B, T, LT, LTA | enveloped, enveloping, detached, internally-detached; counter-signatures |
| **PAdES** (EN 319 142) | ✅ | ✅ | ✅ | B, T, LT, LTA | incremental updates, `/DSS` + `/VRI`, document time-stamps; **invisible signatures only** — see gaps |
| **JAdES** (TS 119 182) | ✅ | ✅ | ✅ | B, T, LT, LTA | compact, JSON and flattened-JSON serialization |
| **ASiC-S / ASiC-E with CAdES** (EN 319 162) | ✅ | ✅ | ✅ | B, T, LT, LTA | container merge, manifests |
| **ASiC-S / ASiC-E with XAdES** (EN 319 162) | ✅ | ✅ | ✅ | B, T, LT, LTA | container merge, manifests, OpenDocument |

Levels above B need a time-stamp source you supply; LT and LTA additionally need
revocation data. The library ships no HTTP TSA/CRL/OCSP clients — see
[What is not here](#what-is-not-here).

### Validation and trust

| Capability | Status |
|---|---|
| EN 319 102-1 validation process (basic, time-stamp, long-term, archival) | ✅ |
| SimpleReport / DetailedReport / ETSI TS 119 102-2 report / diagnostic data | ✅ byte-identical XML to Java DSS |
| Certificate validation, revocation selection, algorithm obsolescence | ✅ |
| eIDAS qualification (signature, time-stamp, certificate; QSCD, MRA, QWAC) | ✅ |
| EU trusted lists: LOTL, TL, pivots, MRA, the refresh job and its cache | ✅ |
| Custom validation policy and ETSI TS 119 312 cryptographic suite | ✅ |
| Signature policies, signature scopes, modification detection (PDF shadow attacks) | ✅ |
| Evidence records (RFC 4998 ERS / RFC 6283 XMLERS) | ❌ framework only — no parser, see gaps |
| Electronic attestations of attributes (EAA) | ⚠️ behind `-tags eaa`, unverified — see gaps |

### Keys and tokens

| Key store | Status |
|---|---|
| PKCS#12 (`.p12` / `.pfx`) — RSA, EC, Ed25519, DSA | ✅ native RFC 7292 reader |
| Any custom `token.SignatureTokenConnection` (HSM, smart card, remote signing service) | ✅ |
| JKS, PKCS#11, Windows certificate store, macOS Keychain | ❌ no Go counterpart under this project's dependency policy |

## What is not here

Honest gaps, all tracked on the documentation site's
[Known gaps](https://ryftcore.github.io/dss-go/compatibility/known-gaps/) page:

- **No online revocation or time-stamping clients.** Upstream's `dss-service`
  module (`OnlineTSPSource`, `OnlineCRLSource`, `OnlineOCSPSource`) is out of
  scope, so levels T/LT/LTA need a `TSPSource` and revocation sources you
  provide. What *is* included: a native HTTP data loader used for AIA
  certificate fetching and for downloading trusted lists, an offline
  `KeyEntityTSPSource` that issues RFC 3161 tokens from a local key, and an RFC
  3161 HTTP client inside the `esig` CLI.
- **No evidence-record parser.** The validation framework, report plumbing and
  ASiC container hooks for evidence records are ported, but no concrete
  RFC 4998 / RFC 6283 analyzer is registered, so documents carrying evidence
  records are not validated against them.
- **EAA is unverified.** The electronic-attestation-of-attributes validation
  blocks are ported and build under `-tags eaa`, but they have been checked only
  by a mechanical 1:1 review against the Java source — no behavioural oracle
  corpus exists for them yet. Treat that build tag as experimental.
- **No visible PDF signature rendering.** Signature-appearance geometry is
  ported, but the rasterization stack (`eu.europa.esig.dss.pdf.visible`) is not:
  the port produces invisible signatures, and comparing rendered page images is
  not available.
- **No REST/SOAP remote services**, no cookbook or BOM modules.
- **MD2 and WHIRLPOOL digests are unavailable** (no Go implementation); every
  other digest upstream supports, including the SHAKE functions, is present.
- **CMS streaming is in-memory only.**
- **Six known CMS/CAdES-layer parity gaps** surfaced by the 248-PDF differential
  run — BadEncodedCMS signature counting, PDF-EOF certificate extraction, four
  legacy PKCS#7 reference-data cases, PLAIN-ECDSA (BSI TR-03111), a
  wrong-digest-algorithm case, and empty-versus-absent `/Reason` and
  `/Location` optionality.
  None of them is a false accept; all are enumerated in
  `dss/pades/testdata/broadgen/README.md`.
- **Open cosmetic follow-ups**: X.500 distinguished-name `toString` rendering,
  some diagnostic-data drift that no verdict depends on, and executor-corpus
  coverage for the QWAC-certificate and EAA-presentation paths.
- **One deterministic-order sanction**: certificate-reference grouping in the
  algorithm-obsolescence process iterates in first-seen order where Java iterates
  in `java.util.HashMap` bucket order. Pinned by a permutation allow-list test;
  affects 2 of 1,328 oracle rows and no conclusion.

## Compatibility

Compatibility here means interoperability, not API similarity, and it is
measured against upstream Java DSS **6.5.RC1** (commit
`4c2129862948bfd53ca1455832260aa17e183cf8`). Most of what follows is a test in
this repository, run on every build against committed Java oracle data; the rest
is a one-time sweep over an upstream corpus too large to vendor. The two are
marked apart below, because the difference is the whole point of the claim.

Cross-validation runs in both directions, per format family. The left column is
live: those tests drive a *running* Java DSS and skip unless `DSS_UPSTREAM_HOME`
points at a built upstream checkout — so they are part of `go test`, but only
for someone who has set that up. The "broad run" figures in the right column are
one-time sweeps, not repeatable from this repository alone; the fixture
counts beside them are committed and run on every build.

| Format | Go-signed, validated by a live Java DSS | Upstream fixtures validated by Go |
|---|---|---|
| CAdES | B, T | 21/21 fixtures, strict |
| XAdES | B, T | 60 fixtures / 659 reference validations |
| PAdES | B, T, LT, LTA (incl. `/DSS` equivalence with Java-extended files) | 36 fixtures; broad run over 248 PDFs / 298 signatures, 0 modification-detection mismatches |
| JAdES | B, T | 30 fixtures; broad run over the full 63-file corpus, 62 byte-identical (the remaining one is an identical rejection worded differently) |
| ASiC (CAdES + XAdES) | B, T, ASiC-S and ASiC-E, plus Java→Go→Java extension round trips | 189-fixture broad corpus, 0 mismatches |

And layer by layer, against Java oracles. The **Runs** column matters: ✅ means
the comparison is a test in `go test ./...` from a full checkout, re-run on
every build against a committed Java oracle; ⏸ means a one-time sweep over an
upstream corpus too large to vendor, repeatable only with a JDK and an
out-of-tree checkout.

| Layer | Comparison | Result | Runs |
|---|---|---|:--:|
| XML canonicalization | 1,829 Santuario known-answer tests, all 7 C14N variants | byte-exact | ✅ |
| XMLDSig | 3,421 oracle rows over 195 fixtures | match | ✅ |
| Validation process | 273-row executor corpus (every building-block indication, sub-indication and qualification) | 273/273 | ✅ |
| End-to-end documents | 60 signed documents, 10 per format family, both engines starting from the original file | 0 mismatches, no divergence allow-lists | ✅ |
| Reports | SimpleReport, DetailedReport, ETSI VR and certificate reports over 252 documents — 1,258 digests in the committed oracle | byte-identical | ✅ |
| EU trusted lists | 37 LOTL/TL fixtures incl. pivot chains, real MRAs, v5/v6 TLs | byte-exact parse parity; 8/8 TL signature verdicts | ✅ |
| PDF engine | 30-document committed oracle vs PDFBox 3.0.7 | 30/30 | ✅ |
| PDF engine, full sweep | the whole 267-document upstream corpus vs PDFBox 3.0.7 (`PDF_CORPUS_DIR`) | 267/267 | ⏸ |
| JSON/JOSE writer | 3,000 documents differential-fuzzed vs jose4j 0.9.6 (the committed oracle keeps 196 rows of it) | 0 byte differences | ⏸ |
| Diagnostic data & policies | 889 diagnostic dumps, 32 policies, 480 crypto-suite keys (the committed oracle keeps 55 dumps) | byte-identical round trips | ⏸ |

The parity gates carry **no tolerance allowances**: the divergence allow-lists
that existed during porting have been deleted, and the comparators are
identifier-inclusive. The methodology, the oracle-generation
recipes and the full numbers are on the documentation site's compatibility pages.

## Documentation

| | |
|---|---|
| **Guides and concepts** | [ryftcore.github.io/dss-go](https://ryftcore.github.io/dss-go/) — what the signature formats and levels mean, how validation reaches a verdict, trusted lists and eIDAS, and how to do the common tasks. Built from `docs/` by the Pages workflow |
| **Compatibility** | [Methodology, numbers and known gaps](https://ryftcore.github.io/dss-go/compatibility/methodology/) |
| **API reference** | [pkg.go.dev/github.com/ryftcore/dss-go/dss](https://pkg.go.dev/github.com/ryftcore/dss-go/dss) |
| **Examples** | [`dss/examples/`](dss/examples/README.md) |
| **Porting conventions** | [dss/PORTING.md](dss/PORTING.md) — the rules contributed ports follow |
| **Upstream tracking** | [UPSTREAM.md](UPSTREAM.md) — the pinned baseline and how to rebase onto a newer DSS |
| **Contributing** | [CONTRIBUTING.md](CONTRIBUTING.md) · [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) |
| **Security** | [SECURITY.md](SECURITY.md) — report vulnerabilities privately, never as a public issue |
| **Changes** | [CHANGELOG.md](CHANGELOG.md) |

## Repository layout

```
dss/           the Go module — github.com/ryftcore/dss-go/dss
  cmd/esig/      the CLI
  examples/      nine runnable example programs
  internal/      engines with no Java class to mirror (ASN.1/BER, CMS, PDF, JOSE, C14N, PFX)
corpus/        heavy oracle fixtures, deliberately OUTSIDE the module
docs/          documentation site sources (MkDocs Material)
```

`corpus/` sits outside the module so that `go get` does not ship ~180 MB of
oracle data; the in-module `testdata/` trees keep small representative fixtures
(about 2 MB in total) so that a bare module checkout still runs meaningful tests.
Tests that need the heavy corpus resolve it through `dss/internal/corpustest` and
skip themselves gracefully when it is absent.

## Development

```sh
cd dss
make check      # go build ./... && go vet ./... && go test ./... -timeout 40m
make race       # the race detector over the full suite
make cli-build  # build ./bin/esig
```

The full gate additionally covers the `eaa` build tag and formatting:

```sh
go build -tags eaa ./... && go vet -tags eaa ./... && gofmt -l .
```

A full checkout has `corpus/` present, so no `corpustest`-gated test skips, and
CI fails the build if one does. Eleven tests still skip, for reasons unrelated
to `corpus/`:

- **Eight need a built upstream Java DSS checkout** — the six
  `TestDownstreamCrossValidation` cases (`cades`, `xades`, `jades`, `pades`,
  `asic/cades`, `asic/xades`) plus `TestCAdESLevelBaselineBBytesMatchUpstream`
  and `asic/cades`'s `TestRoundTripJavaBuiltExtendedByGo`. They skip unless
  `DSS_UPSTREAM_HOME` points at a checkout whose Maven modules are built, with
  `java`/`javac`/`mvn` on `PATH`.
- **Two are gated on their own opt-in environment variables** for out-of-tree
  PDF comparisons: `TestWriterOverCorpus` (`PDF_CORPUS_DIR`) and
  `TestDumpWriterOutputs` (`PDFOUT`).
- **One skips by design**: `TestDOMRoundTrip/ns-relative.xml`, a fixture
  canonicalization is supposed to reject.

Contribution rules, including the upstream-tracking rule for any behaviour that
diverges from Java DSS, are in [CONTRIBUTING.md](CONTRIBUTING.md).

## Attribution and license

This project is a derivative work: a Go port of the Java **DSS — Digital
Signature Services** library published by the European Commission at
<https://github.com/esig/dss>, copyright © 2015 European Commission and DSS
contributors, baseline version 6.5.RC1
(`4c2129862948bfd53ca1455832260aa17e183cf8`).

It is licensed under the **GNU Lesser General Public License, version 2.1**, the
same license as upstream. The full text is in [LICENSE](LICENSE); the
project-level attribution is in [NOTICE](NOTICE), and every Go file that ports a
specific Java class carries a `// Ported from …` header naming the exact
upstream source file it derives from (packages replacing machinery that has no
Java class to mirror — the ASN.1, CMS, PDF, JOSE and canonicalization engines —
state their provenance in their `doc.go` instead). Those headers are the
file-level attribution record and must not be removed.

**This project is not affiliated with or endorsed by the European Commission or
the DSS project.**

### A note on linking

Go links statically, which is the case LGPL-2.1 section 6 addresses. In outline:
you may distribute a program that *uses* this library under terms of your own
choosing, provided you permit modification of the library and reverse
engineering for debugging such modifications, and you satisfy one of the
options in section 6 — for a statically linked binary, option 6(a) is the usual
one, and requires supplying whatever a recipient needs to relink the program
against a modified version of the library. Changes to the library itself stay
under LGPL-2.1 (inbound equals outbound). That is a summary, not the license and
not legal advice: read [LICENSE](LICENSE), and take your own counsel.
