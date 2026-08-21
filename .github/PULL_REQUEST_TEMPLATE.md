## What

<!-- What does this change do, in one or two sentences? -->

## Why

<!-- Motivation: bug, upstream parity gap, new feature, docs fix, etc. -->

## Upstream alignment

- [ ] This change matches Java DSS's behavior for the equivalent code path, OR
- [ ] This change intentionally diverges from Java DSS — the reason is
      documented in a code comment (see CONTRIBUTING.md's upstream-tracking
      rule).

## Checklist

- [ ] `gofmt -w` run on changed files
- [ ] `cd dss && go build ./... && go build -tags eaa ./... && go vet ./... && go vet -tags eaa ./...` passes
- [ ] `cd dss && go test ./... -count=1` passes (no test skipped or weakened to make this pass; new tests added for new behavior)
- [ ] Existing `// Ported from` attribution headers are unchanged
- [ ] Docs/comments updated if public API or behavior changed
- [ ] I have read CONTRIBUTING.md and agree to its sign-off / LGPL-2.1 inbound=outbound terms
