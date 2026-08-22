# FAQ

## Is this the same as Java DSS?

It is a **port** of it — the same algorithms, checks, reports and enumeration
values, written in Go against upstream 6.5.RC1, class by class, with each file
naming its Java original in a `Ported from` header. (The engines DSS gets from
BouncyCastle, PDFBox, jose4j and Santuario have no Java DSS class to mirror;
those were written to the specifications and pinned to the real libraries'
output byte for byte.) It is not the
same API: Java DSS's shape is not idiomatic Go, and reproducing it exactly is an
explicit anti-goal. What is reproduced exactly is the *output*. See
[Methodology](compatibility/methodology.md).

## Is it affiliated with the European Commission or the DSS project?

**No.** It is an independent derivative work, distributed under the same
LGPL-2.1 licence as the original, with attribution in `NOTICE` and in every
ported file's header.

## Can I use this in production?

That is your call, and it should be an informed one. In its favour: parity is
verified against recorded Java DSS output at several layers, on every build,
with no divergence allowlists in the end-to-end comparisons; the
[numbers](compatibility/numbers.md) name the tests and the commands so you can
re-run them yourself. Against it: it is a young port of a large standards stack,
and the [known gaps](compatibility/known-gaps.md) page is not short.

For anything with legal or financial weight, validate the specific flows you
depend on against Java DSS with your own documents — the method is in
[Verifying Java DSS interop](guides/verifying-java-dss-interop.md).

## What are the dependencies?

Two direct — `golang.org/x/crypto` and `golang.org/x/text` — plus
`golang.org/x/sys` indirectly. No cgo. Everything else, including the CMS and
ASN.1 engines, XML canonicalization, the PDF engine and the JOSE serializer, is
code in the module. `cat dss/go.mod` is the authority.

## Why does validation return `INDETERMINATE` on a signature I know is good?

Almost always because you did not tell it what to trust
(`NO_CERTIFICATE_CHAIN_FOUND`), or because revocation data is missing or stale.
`INDETERMINATE` means "I cannot conclude with what I have", not "this is bad".
The sub-indication says what is missing. See
[How validation works](concepts/validation.md).

## Why is the qualification `NA`?

Because you supplied no trusted-list information. `NA` means *not determined*,
not *not qualified*. Qualification comes from the EU trusted lists; a plain
trust store cannot produce one. See
[Validate against the EU trusted lists](guides/validate-against-eu-trusted-lists.md).

## Which signature format should I use?

Decided mostly by what you are signing, and often prescribed by whoever receives
it. There is a decision table in
[Signature formats](concepts/signature-formats.md).

## Which level should I use?

B if it is checked immediately. T if *when* it was signed matters. LT if it
will be checked after the certificate expires. LTA if it must survive a decade
and someone will maintain it. See
[Signature levels](concepts/signature-levels.md).

## Can I sign with a smart card or an HSM?

Not through PKCS#11 — that has no Go counterpart and is not supported. But the
library never touches your private key: it hands a digest to a
`token.SignatureTokenConnection` and asks for a signature. Implement that
interface over whatever driver, service or SDK you have and everything else
works unchanged.

## Can I use a JKS key store, the Windows store, or the macOS Keychain?

No. Each is reached in Java through a JCA provider with no Go counterpart, and
the port's dependency policy allows nothing outside the standard library and
`golang.org/x/…`. For JKS, convert it: `keytool -importkeystore -srcstoretype
JKS -deststoretype PKCS12`.

## Does the library make network requests?

Not unless you ask. It ships no HTTP time-stamping or revocation client at all
— you supply those. Trusted-list downloading and AIA certificate fetching do
exist and are explicitly opt-in. See
[Timestamps and TSAs](guides/timestamps-and-tsas.md).

## How do I time-stamp without a TSA of my own?

Point the CLI at any public RFC 3161 endpoint with `-tsa <url>` — it ships a
client. From Go, implement `TSPSource` (the CLI's implementation is a small,
readable, dependency-free example), or use the library's local key-based source
for tests.

## Can I add a visible signature to a PDF?

Not today. The signature field and its rectangle are created; nothing is painted
into them, because the native PDF engine has no rasteriser. See
[Known gaps](compatibility/known-gaps.md).

## Why is the API reference not on this site?

Because Go already has a good one. Everything exported is documented on
[pkg.go.dev](https://pkg.go.dev/github.com/ryftcore/dss-go/dss), generated from the
source, so it cannot drift. This site explains *concepts* and *decisions*;
duplicating godoc here would only create a second version to go stale.

## What does LGPL-2.1 mean for my program?

The library is LGPL-2.1, the same as upstream — it must stay that way, because
this is a derivative work. In broad terms the LGPL is designed so that a program
can *use* the library without itself becoming LGPL, subject to conditions about
relinking, which are stricter in the static-linking case that Go always uses.

**This is a summary, not legal advice, and it is not a licence grant.** Read
`LICENSE`, and if your distribution model matters to you, ask a lawyer.

## How do I keep up with upstream DSS releases?

`UPSTREAM.md` documents the baseline pin and the procedure: diff the upstream
tags, locate the Go counterparts through the `Ported from` headers, and re-run
the oracle harness and parity gates.

## How do I report a bug?

The [issue tracker](https://github.com/ryftcore/dss-go/issues), which has templates
for bugs, features, and specifically for **interop mismatches** against Java
DSS. For a suspected security vulnerability, do **not** open a public issue —
follow `SECURITY.md`.

## How do I contribute?

`CONTRIBUTING.md` has the development setup, the gate commands, and the porting
conventions. The one rule worth knowing before you start: **a change that
diverges from upstream Java behaviour needs a documented reason.** Interop is
the product.

## Why is the module `github.com/ryftcore/dss-go/dss` and not `github.com/ryftcore/dss-go`?

Because the Go module lives in the `dss/` subdirectory, so that the heavy test
corpora can sit beside it at the repository root — outside the module, and
therefore outside what `go get` downloads. See
[Methodology](compatibility/methodology.md).

## Why do some tests skip when I run `go test` after `go get`?

They do not, by design — the module ships small fixtures so every package has
real tests. What skips is anything needing the heavy `corpus/` tree, which is
not in the module zip. Clone the repository if you want to run everything; the
CI test job fails the build if any test skips for want of `corpus/` in a full
checkout.

A full checkout does still skip a small, fixed set of tests for reasons that
have nothing to do with `corpus/`: those that cross-validate against a **built
upstream Java DSS checkout** (they need `DSS_UPSTREAM_HOME` plus `java`,
`javac` and `mvn`), two gated on their own opt-in environment variables, and one
canonicalization fixture that is rejected by design. `README.md` lists them.
