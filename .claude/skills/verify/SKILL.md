---
name: verify
description: Build and drive the esig CLI to observe a dss-go change at runtime.
---

# Verifying dss-go changes

- Module root is `dss/`. `go.mod` needs Go 1.27; the system Go auto-downloads the toolchain on first build (~1 min).
- Build the CLI: `cd dss && go build -o <scratch>/esig ./cmd/esig`. For before/after comparisons, build `origin/main` in a `git worktree` the same way.
- Fixtures: `dss/testdata/sample.pdf`, `signer_rsa.p12` and `tsa_ec.p12` (password `testpassword`). Pass passwords as `-p12-pass env:VAR`.
- Sign: `esig sign testdata/sample.pdf -format pades -level B -p12 testdata/signer_rsa.p12 -p12-pass env:PW -out o.pdf`.
  - `-format xades` needs XML input. `-format jades -level T` fails by default (compact serialization).
- T-level needs a TSA. Copy `testTSA` from `cmd/esig/e2e_test.go` and `tsaRequest`/`tsaResponse` from `cmd/esig/tsa.go` into a throwaway `package main` under `dss/` (so it can import `internal/cmscore`). Wrap them in an `http.ListenAndServe` handler. Delete the throwaway package afterwards.
- `esig validate` returns 1 for INDETERMINATE: the test chain is untrusted, so that is expected.
- Hostile input: craft files with Python and run `validate`/`inspect`/`sign` under a Python `subprocess` wrapper that sets `RLIMIT_AS` and a timeout. Read peak RSS from `getrusage(RUSAGE_CHILDREN)`, because there is no `/usr/bin/time`.
- Release: `go install github.com/goreleaser/goreleaser/v2@latest`. Then run `goreleaser check` and `goreleaser release --snapshot --clean --skip=publish` from the repo root. `dist/` is gitignored.
