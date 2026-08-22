# Broad differential corpus runner (PAdES)

Differential harness over the **entire** upstream `dss-pades/src/test/resources`
corpus (~248 PDFs, 298 signatures). Unlike the committed 36-fixture
cross-validation suite, this walks every upstream PDF
through both implementations and diffs the dumped JSON — it is the tool that
surfaced the dispatch, DSS-dict-equality, PDFDocEncoding, and large-exponent
defects the sampled suite could not see.

- `BroadOracle.java` — Java side; compile against a maven-built upstream DSS
  (same classpath recipe as `../gen/CrossValidationOracle.java`), run with the
  corpus directory as argument, produces `broad-java.json`.
- `gobroad_main.go` — Go counterpart (kept under `testdata/` so the go tool
  ignores it; build it from a scratch directory with
  `go mod edit -replace`, or `go run` it with the repo module on the path),
  produces `broad-go.json` in the same shape.
- Diff the two JSON files; every difference is a parity defect in one side.

Known accepted differences, all booked as CMS-layer follow-ups (see
`docs/compatibility/known-gaps.md`): BadEncodedCMS signature-count divergence,
pdf-eof certificate extraction, 4 legacy PKCS#7 reference-data cases,
PLAIN-ECDSA (BSI TR-03111) verification, wrong-digest-algo, and
present-but-empty `/Reason`/`/Location` optionality. Everything else matches
exactly — verdicts, modification-detection buckets, byte ranges, coverage.

TODO (tracked): promote this into a maintained, regularly-run harness once the
CMS-layer follow-ups land.
