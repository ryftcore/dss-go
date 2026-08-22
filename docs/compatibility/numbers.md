# The numbers

Every figure on this page names the test that produces it and the command that
reproduces it. Run the command from a **full checkout** — one that includes the
repository-root `corpus/` tree — with `cd dss` first.

If a number here and a test disagree, the test is right and this page is a bug.
Please [report it](https://github.com/ryftcore/dss-go/issues).

## Upstream baseline

| | |
|---|---|
| Upstream project | [esig/dss](https://github.com/esig/dss) |
| Version | **6.5.RC1** |
| Commit | `4c2129862948bfd53ca1455832260aa17e183cf8` |

Recorded in `UPSTREAM.md` and `NOTICE`. Every ported source file names the
Java file it came from in its header comment.

## The library itself

| Measure | Value | How to check |
|---|---|---|
| Direct runtime dependencies | **2** — `golang.org/x/crypto`, `golang.org/x/text` | `cat dss/go.mod` |
| Indirect runtime dependencies | **1** — `golang.org/x/sys` | `cat dss/go.mod` |
| Exported packages | **115** | `go list ./... \| grep -v '/internal\|/cmd\|/examples\|/harness' \| wc -l` |
| Test files | **573** | `find dss -name '*_test.go' \| wc -l` |
| cgo | none | — |

No BouncyCastle equivalent, no XML library, no PDF library, no JOSE library.
Every one of those is code in the module: the CMS and ASN.1 engines, the seven
canonicalization algorithms, the PDF reader and incremental writer, the JOSE
serializer.

## Parity measurements

Each row is a test in `go test ./... -count=1`.

| What is compared | Scale | Test | Command |
|---|---|---|---|
| **End-to-end document validation**, both engines starting from the *original* signed document — diagnostic-data core, every building-block conclusion, every top-level indication and sub-indication, every qualification, **identifiers included, no tolerance** | **60 documents**, 10 in each of the six format families (CAdES, XAdES, PAdES, JAdES, ASiC-CAdES, ASiC-XAdES) | `TestDocumentLevelOracle` | `go test ./harness/ -run TestDocumentLevelOracle -count=1` |
| **Verdicts through the EN 319 102-1 engine** from already-built diagnostic data — every building block's indication, sub-indication and qualification | **273 documents** | `TestFullCorpusExecutorOracle` | `go test ./validation/executor/ -run TestFullCorpusExecutorOracle -count=1` |
| **Rendered report bytes**, compared as digests against Java's | **1,258 report digests** over **252 documents** — SimpleReport, DetailedReport and ETSI validation report for each, plus the two certificate-report flavours for 251 of them. A default build compares 1,255 of the 1,258: three DetailedReport digests belong to documents whose reports carry the `eaa`-gated blocks, and are checked only under `-tags eaa` | `TestFullCorpusReportsByteParity` | `go test ./validation/executor/ -run TestFullCorpusReportsByteParity -count=1` |
| **XML canonicalization**, byte-exact against Apache Santuario's output through the DSS call path, across all seven registered algorithms | **1,829 known-answer tests** | `TestC14nGolden` | `go test ./internal/xmlc14n/ -run TestC14nGolden -count=1` |
| **XMLDSig reference processing and signature verification** against a Java-generated oracle | **3,421 known answers** over **195 fixtures** (the test logs both counts) | `TestKnownAnswers` | `go test ./internal/xmldsig/ -run TestKnownAnswers -count=1 -v` |
| **PDF parsing**, against an oracle generated with pdfbox — file sizes, revision boundaries, object structure, encryption | **30 documents**, each byte-pinned by a SHA-256 manifest | `TestOracleCorpus`, `TestOracleManifest` | `go test ./internal/pdf/ -run 'TestOracle' -count=1` |
| **Trusted-list and LOTL parsing**, byte-exact — including pivot chains, real MRA lists, and v5 and v6 trusted lists | **37 fixtures** | `TestTSLParsingOracle` | `go test ./harness/ -run TestTSLParsingOracle -count=1` |
| **Trusted-list signature verdicts**, and the offline validation-job cache state machine | — | `TestTLValidationJobOracle`, `TestCertificateQualificationOracle` | `go test ./harness/ -count=1` |

!!! note "How these were checked for this page"
    Each figure above was produced by running the named command and counting
    the passing subtests, not taken from a previous document. Counts of
    "documents" are subtest counts; counts of "digests" are fields in the
    committed oracle file.

## Live cross-validation against a running Java DSS

The measurements above compare against **recorded** Java output. Separately,
each format family was cross-validated against a **live**
Java DSS 6.5.RC1 build in both directions — documents signed here validated by
Java, documents signed by Java validated here, plus broad differential sweeps
over entire upstream corpora.

Those runs are not part of `go test` and are not reproducible from this
repository alone — see the warning below. What *is* reproducible is
everything in the table above.

!!! warning "They are not reproducible from this repository alone"
    Re-running live cross-validation needs a JDK and a built upstream DSS
    classpath, neither of which is committed here. The generator programs are —
    under the oracle `gen/` directories, with the classpath recipe in
    `dss/harness/README.md` — so the runs are repeatable by anyone who sets that
    up, but they are not part of `go test`, and no CI job asserts them.

    Everything in the table above **is** part of `go test`. That is the
    difference this page is careful about.

## What these numbers do not say

- They do not say the port is bug-free. They say that on this material, at these
  layers, it agrees with the reference implementation.
- They do not cover behaviour outside the tested corpora. The
  [known gaps](known-gaps.md) list what is known to be outside them.
- They say nothing about *your* documents. If yours are unusual, test them —
  [Verifying Java DSS interop](../guides/verifying-java-dss-interop.md) is the
  method.

## Next

- [Methodology](methodology.md) — how the oracles are built and why.
- [Known gaps](known-gaps.md) — the honest other half of this page.
