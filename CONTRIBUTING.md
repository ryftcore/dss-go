# Contributing to esig

Thanks for your interest in contributing. esig is a Go port of
[esig/dss](https://github.com/esig/dss) (the EU "DSS" digital-signature
library), so contributions here have one extra constraint that most Go
projects don't: behavior is expected to track a specific upstream Java
codebase. Please read the "Porting conventions" section below before
sending non-trivial changes.

## Development setup

- **Go 1.27 or later** (the module currently pins `go 1.27.0` in
  `dss/go.mod`; keep the module's `go` directive and this requirement in
  sync if you bump it). The public API uses generic methods, which Go 1.27
  is the first release to allow.
- No other tooling is required to build and test the library itself.
  Docs-site and CI tooling (MkDocs, golangci-lint, goreleaser, etc.) are
  only needed if you're touching `docs/` or `.github/`.

Clone and build:

```sh
git clone https://github.com/ryftcore/dss-go
cd dss-go/dss
go build ./...
```

The Go module root is `dss/` (module path `github.com/ryftcore/dss-go/dss`),
not the repository root — run `go`/`make` commands from inside `dss/`.

### Make targets

`dss/Makefile` provides:

```sh
make build   # go build ./...
make vet     # go vet ./...
make test    # go test ./... -timeout 40m
make race    # go test ./... -race -timeout 60m
make check   # build + vet + test
```

### Full gate

Before sending a PR, run the same checks CI runs:

```sh
cd dss
go build ./...
go vet ./...
go test ./... -count=1 -timeout 40m
go build -tags eaa ./...
```

`-tags eaa` builds the electronic-attestation-of-attributes (EAA)
validation code paths, which are gated behind a build tag (see the
documentation site's "Known gaps" page for their current status) — a PR
must not break either the default build or the `eaa`-tagged build.

If your change touches a package with a large fixture corpus, see
`dss/internal/corpustest` and the top-level `corpus/` directory: heavy
test fixtures live outside the Go module so `go get` stays small, and
tests that need them skip gracefully when `corpus/` isn't present (e.g.
in a module downloaded via `go get` rather than a full git checkout).
Run the full checkout gate above (with `corpus/` present) so nothing is
silently skipped.

Never weaken or delete an existing test to make it pass — if a test
seems wrong, say so in the PR description and propose the fix
separately.

### Formatting

Run `gofmt -w` on every changed `.go` file; CI checks that `gofmt -l`
reports nothing.

## Porting conventions

Because this is a port, "the right way to write this in idiomatic Go"
is deliberately constrained by "does this still map back to the
upstream Java source." The binding rules — file/package layout, enum
representation, error/exception mapping, document/stream modeling,
crypto/ASN.1 conventions, and the test-vector policy — are documented in
[`dss/PORTING.md`](dss/PORTING.md). Read it before porting new code or
modifying ported code.

Key points worth calling out here:

- Every ported file keeps a `// Ported from <upstream-path> (DSS
  <version>).` header. This is load-bearing attribution (see `NOTICE`)
  — never remove or rewrite it, including when you refactor the file's
  contents.
- Enum values are exactly the upstream Java `name()` string, and
  OIDs/URIs/registry values are copied verbatim from upstream. Never
  "fix" or normalize them based on what looks more correct.
- **Upstream-tracking rule**: any change whose observable behavior
  diverges from what upstream Java DSS does for the same input (a
  different validation verdict, a different byte-level signature
  output, a different error) needs a comment and, ideally, a
  `dss/PORTING.md` note explaining *why* — a known upstream bug you're
  intentionally not reproducing, a Go-standard-library limitation, an
  accepted gap, etc. Silent behavioral drift from
  upstream is treated as a bug, not a feature, even when the new
  behavior is arguably "better."
- New third-party dependencies are effectively out of scope: the
  library's dependency set is deliberately limited to
  `golang.org/x/crypto` and `golang.org/x/text` (with `golang.org/x/sys`
  as an indirect dependency of those). If you believe a new dependency
  is genuinely necessary, open an issue to discuss it before writing
  code — see the dependency policy at the end of `dss/PORTING.md`.
- Interoperability (byte-level output parity and validation-verdict
  parity against upstream Java DSS) is the project's compatibility
  contract, not API-identical mirroring of the Java class structure. See
  the documentation site's compatibility pages for the precise claims and
  how they're verified.

## Reporting bugs / interoperability mismatches

Use the issue templates under `.github/ISSUE_TEMPLATE/`. For a result
that disagrees with upstream Java DSS (a different validation verdict,
a signature one side accepts and the other rejects, etc.), use the
"interop mismatch" template and include the Java DSS version you
compared against plus the differing report output — that's the fastest
path to a fix.

## Sign-off (DCO)

By submitting a contribution, you certify that you wrote it (or
otherwise have the right to submit it) under the project's license, per
the [Developer Certificate of Origin](https://developercertificate.org/).
Please add a `Signed-off-by` trailer to your commits:

```
Signed-off-by: Your Name <you@example.com>
```

`git commit -s` adds this automatically.

## License of contributions

esig is licensed under LGPL-2.1 (see `LICENSE`). This project uses
**inbound = outbound** licensing: by submitting a contribution, you
agree it is licensed under LGPL-2.1, the same as the rest of the
project, with no additional restrictions.

## Code of conduct

Participation in this project is governed by our
[Code of Conduct](CODE_OF_CONDUCT.md).
