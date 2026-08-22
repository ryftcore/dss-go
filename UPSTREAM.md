# Upstream tracking

esig is a Go port of [esig/dss](https://github.com/esig/dss), the EU
"DSS - Digital Signature Services" Java library. This file records the
current upstream baseline and how to move that baseline forward.

## Current baseline

- **Repository**: https://github.com/esig/dss
- **Version**: 6.5.RC1
- **Commit**: `4c2129862948bfd53ca1455832260aa17e183cf8`

Every ported Go file carries a `// Ported from <upstream-path> (DSS
6.5.RC1).` header identifying the exact upstream source file and
version it was derived from — see `NOTICE` and `dss/PORTING.md` for how
that attribution is structured, and `docs/migrating-from-java/index.md`
for the full Maven-module → Go-package mapping.

esig does not track upstream continuously; it was ported from this one
pinned commit. Rebasing onto a newer upstream release is a deliberate,
occasional exercise, not an automated sync.

## What "compatibility" means here

esig's compatibility contract with upstream is **interoperability**,
not a literal API mirror: byte/spec-level signature output, validation
verdicts (Indication/SubIndication), and report/diagnostic-data schemas
match upstream Java DSS for the same inputs, verified by cross-
validating in both directions against upstream-derived test vectors and
by differential test harnesses that replay upstream oracle output
against the Go implementation. See `docs/compatibility/methodology.md`
and `docs/compatibility/numbers.md` for the current, verified scope of
that claim, and `docs/compatibility/known-gaps.md` for the deviations
that are tracked rather than fixed.

## Rebasing onto a future DSS release

There is no single automated command for this; it is a structured
manual process. As a rough recipe, based on how this port was built and
kept in sync during development:

1. **Diff upstream.** Check out (or diff) the new upstream tag/commit
   against the commit this port is pinned to (above). Focus the diff on
   the Maven modules listed in `docs/migrating-from-java/index.md`'s
   module-mapping table — that table is also the map from an upstream
   module to the Go package(s) that need the corresponding update.

2. **Locate the Go twin of each changed Java file.** Every ported file's
   `// Ported from ...` header names its upstream source path — grep for
   the changed upstream file's path across `dss/` to find the Go file(s)
   derived from it. A file with no Go twin is either newly added
   upstream (needs a first port) or in a module this port deliberately
   does not cover (see `docs/compatibility/known-gaps.md`).

3. **Port the change**, following the conventions in `dss/PORTING.md`
   (file/package layout, enum representation, error mapping, streaming
   document model, ASN.1/crypto conventions, test-vector policy). Keep
   the `Ported from` header's version string updated to the new upstream
   version once the file matches it.

4. **Re-run the oracle / cross-validation harnesses.** Compatibility in
   this project is verified, not assumed. Test vectors and byte-parity
   fixtures generated against the OLD upstream version need
   regenerating against the NEW one wherever the underlying algorithm or
   output changed. Oracle generators live alongside the fixtures they
   produce, typically under a package's `testdata/oracle/gen/` directory
   (Java programs invoked against the upstream classpath) — see
   `dss/harness/README.md` for one fully worked example (the
   document-level end-to-end oracle covering all six signature-format
   families) including the exact classpath/`javac`/`java` invocation
   pattern to follow for other packages' generators.

5. **Run the parity gates.** After regenerating fixtures and porting the
   change, the full gate must be green: `go build ./...`, `go vet
   ./...`, `go test ./... -count=1 -timeout 40m`, and `go build -tags
   eaa ./...` from `dss/` (see `CONTRIBUTING.md`). Do not weaken or
   delete an oracle/parity test to make it pass — a mismatch means
   either the port needs a fix or the fixture/oracle needs a documented
   reason it's expected to differ (following the pattern already used
   for the tracked gaps in `docs/compatibility/known-gaps.md`).

6. **Update the baseline.** Once the rebase is complete and the gate is
   green, update the version/commit at the top of this file, in `NOTICE`,
   and in `docs/compatibility/numbers.md`'s "Upstream baseline" table,
   and record what changed (new modules ported, gaps closed, new gaps
   discovered) as a `CHANGELOG.md` entry.
