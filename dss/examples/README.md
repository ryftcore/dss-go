# Examples

Nine small, runnable programs, each teaching one story with the root
[`dss`](https://pkg.go.dev/github.com/ryftcore/dss-go/dss) facade. Every example
is a self-contained `package main` you can read top to bottom; run it with:

```sh
go run ./examples/01-sign-pdf-pades
```

from anywhere - the working directory does not matter, each example locates
its fixtures relative to its own source file.

All but one example run entirely offline, against the same tiny, self-signed
test fixtures the facade's own tests use (`dss/testdata`, see that
directory's `README.md`): a self-signed RSA key, an EC key acting as a local
RFC 3161 time-stamp authority, and a minimal sample PDF. None of them are
trust anchors of anything real; they exist to make a runnable example
possible without a network call or a real CA.

| Example | Story |
|---|---|
| [`01-sign-pdf-pades`](01-sign-pdf-pades) | Sign a PDF at PAdES-BASELINE-B, then raise a second copy to B-T with a time-stamp. Notes where the port's local test TSA stands in for a real HTTP time-stamping authority. |
| [`02-validate-pdf`](02-validate-pdf) | Validate the same signature twice - without a trust anchor (correctly INDETERMINATE) and with one (TOTAL_PASSED) - to show what `ValidateOptions` actually controls. |
| [`03-xades-enveloped-invoice`](03-xades-enveloped-invoice) | Sign an XML invoice with an enveloped XAdES-BASELINE-B signature, the usual shape for XML business documents. |
| [`04-cades-detached`](04-cades-detached) | Produce a detached CAdES-BASELINE-B signature over a binary payload, and show what validation looks like with and without the original content. |
| [`05-jades-json-payload`](05-jades-json-payload) | Sign a JSON payload with a compact JAdES-BASELINE-B signature - the JOSE-family format for JSON APIs. |
| [`06-asice-container`](06-asice-container) | Sign two documents into one ASiC-E container and list what is actually inside it (an ASiC container is a zip file). |
| [`07-validate-eu-trusted-lists`](07-validate-eu-trusted-lists) | Run the TSL validation job against the real EU List Of Trusted Lists - the piece that turns "the signature verifies" into "the signature is eIDAS qualified". Needs the network; degrades gracefully without it. |
| [`08-custom-policy`](08-custom-policy) | Load the library's default validation policy, relax one constraint, and see the verdict change - `ValidateOptions.Policy` takes any policy document, not just the built-in one. |
| [`09-render-reports`](09-render-reports) | Walk through all four reports `Validate` returns: SimpleReport, DetailedReport, diagnostic data, and the ETSI TS 119 102-2 validation report. |

## What these are not

These are teaching examples for the facade, not a tour of the full port. The
facade itself covers the common paths only (see the `dss` package doc's
"What the facade deliberately leaves out" section) - counter-signatures,
visible PDF signature appearances, custom certificate/revocation sources,
XAdES references and transforms, and more are reached through the
underlying packages the facade delegates to
([`cades`](https://pkg.go.dev/github.com/ryftcore/dss-go/dss/cades),
[`xades`](https://pkg.go.dev/github.com/ryftcore/dss-go/dss/xades),
[`pades`](https://pkg.go.dev/github.com/ryftcore/dss-go/dss/pades),
[`jades`](https://pkg.go.dev/github.com/ryftcore/dss-go/dss/jades),
[`asic`](https://pkg.go.dev/github.com/ryftcore/dss-go/dss/asic),
[`validation`](https://pkg.go.dev/github.com/ryftcore/dss-go/dss/validation),
[`tsl`](https://pkg.go.dev/github.com/ryftcore/dss-go/dss/tsl)), documented on
their own.

For the CLI built on this same facade, see `dss/cmd/esig`.
