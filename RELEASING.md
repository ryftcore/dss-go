# Releasing

esig follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html);
`CHANGELOG.md` is the release notes.

## Why every release has two tags

The Go module is `github.com/ryftcore/dss-go/dss` and lives in the `dss/`
subdirectory. Go resolves a version of a module in a subdirectory only from a
tag carrying that directory as a prefix, so:

| Tag           | Consumed by                                            |
| ------------- | ------------------------------------------------------ |
| `dss/vX.Y.Z`  | `go get github.com/ryftcore/dss-go/dss@vX.Y.Z`, `go install …/cmd/esig@vX.Y.Z`, pkg.go.dev |
| `vX.Y.Z`      | `.github/workflows/release.yml` → GoReleaser → GitHub release with `esig` binaries |

Both must point at the same commit. The release workflow refuses to publish if
`dss/vX.Y.Z` is missing or points elsewhere.

## Checklist

1. On `main`, with CI green on the commit to release:
   - `CHANGELOG.md`: move the `[Unreleased]` entries under a new
     `## [X.Y.Z] - YYYY-MM-DD` heading, leave an empty `[Unreleased]`, and
     update the compare links at the bottom.
   - If a compatibility claim changed, `docs/compatibility/` is current.
2. Run the full gate from `dss/` (see `CLAUDE.md` → Commands), including
   `go test ./... -count=1 -timeout 40m` with `corpus/` present.
3. Optionally dry-run the release locally from the repo root:
   `goreleaser check && goreleaser release --snapshot --clean`.
4. Tag and push both tags in one push:

   ```sh
   git tag -s dss/vX.Y.Z -m "dss vX.Y.Z"
   git tag -s vX.Y.Z     -m "esig vX.Y.Z"
   git push origin dss/vX.Y.Z vX.Y.Z
   ```

5. Check the GitHub release: archives for linux/darwin/windows × amd64/arm64,
   `checksums.txt`, and `esig version` printing `esig vX.Y.Z`. Paste the
   CHANGELOG section into the release body.
6. Confirm the module resolves:
   `GOPROXY=https://proxy.golang.org go list -m github.com/ryftcore/dss-go/dss@vX.Y.Z`.

A tag is never moved or deleted once pushed: the Go checksum database has
recorded it. Fix a bad release with a new patch version.
