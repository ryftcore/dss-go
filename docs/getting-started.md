# Getting started

Twenty minutes from nothing to a signed, validated document — first from Go,
then from the command line.

## What you need

- **Go 1.27 or later.** The module pins `go 1.27.0` in its `go.mod`.
- **A signing key in a PKCS#12 file** (`.p12` or `.pfx`) with its password.
  For a first experiment a self-signed one is fine — the repository ships one
  under `dss/testdata/` that all the examples use. For anything real you want a
  certificate from a certificate authority, and for an eIDAS *qualified*
  signature, one issued on a qualified signature creation device. See
  [Trust, eIDAS and qualification](concepts/trust-and-eidas.md).

!!! note "PKCS#12 only"
    JKS, PKCS#11 (smart cards and HSMs through a driver), the Windows
    certificate store and the macOS Keychain are all reached in Java through
    JCA providers that have no Go counterpart, and are not supported here.
    PKCS#12 is, fully, with RSA, ECDSA, Ed25519 and DSA keys.

    A hardware token can still be used: implement the port's
    `token.SignatureTokenConnection` interface over whatever driver you have.
    The library never touches your private key directly — it asks the token to
    sign a digest.

## Install

```sh
go get github.com/utain/esig/dss
```

Two runtime dependencies come with it — `golang.org/x/crypto` and
`golang.org/x/text` (plus `golang.org/x/sys`, indirectly). Nothing else: no
BouncyCastle equivalent, no XML library, no PDF library, no cgo. The port
replaces every one of those with code in the module.

## Sign and validate, in twenty lines

```go
package main

import (
	"fmt"
	"log"

	"github.com/utain/esig/dss"
)

func main() {
	signer, err := dss.OpenPKCS12("keystore.p12", "password") // (1)!
	if err != nil {
		log.Fatal(err)
	}
	defer signer.Close()

	doc := dss.NewDocument("invoice.xml", []byte("<invoice><total>42</total></invoice>"))

	signed, err := dss.Sign(doc, signer, dss.SignOptions{ // (2)!
		Format: dss.FormatXAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatal(err)
	}

	reports, err := dss.Validate(signed, dss.ValidateOptions{ // (3)!
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range reports.Verdicts() { // (4)!
		fmt.Println(v.ID, v.Indication, v.SubIndication, v.SignatureLevel, v.SignedBy)
	}
}
```

1. `Signer` pairs an open key store with the one key entry to sign with. Use
   `dss.NewSigner` to drive any other token implementation.
2. `Format` picks the signature family, `Level` picks how long it has to stay
   verifiable. Everything else in `SignOptions` has a documented default —
   SHA-256, enveloped for XAdES, the signing time is now.
3. Validation needs to know what to trust. Without a trust anchor the answer is
   correctly <span class="verdict indeterminate">INDETERMINATE</span>: the
   chain reaches nothing you have said you believe in. Here we trust the
   signer's own chain, which is only sensible in an example.
4. One `Verdict` per signature found. `reports.Valid()` is the
   single-boolean answer if that is all you need.

That prints one line per signature, like this:

```console
S-562166A7…F79A1BC TOTAL_PASSED  XAdES-BASELINE-B Go Port Test RSA
```

The sub-indication is empty because there is nothing to explain. Drop the
`TrustedCertificates` option and the same line reads `INDETERMINATE
NO_CERTIFICATE_CHAIN_FOUND` — which is the correct answer, not a failure. [How
validation works](concepts/validation.md) explains why.

!!! tip "Where the runnable versions live"
    `dss/examples/` holds nine complete, commented programs — one per story,
    each runnable offline against the repository's small test fixtures:

    ```sh
    go run ./examples/01-sign-pdf-pades
    ```

    And the root package's own `Example*` functions on
    [pkg.go.dev](https://pkg.go.dev/github.com/utain/esig/dss#pkg-examples) are
    compiled and run by the test suite, so they are always current.

## Where the options actually matter

The facade takes two option structs, and almost everything in them has a
sensible default. The fields worth knowing on day one:

| Option | Why you would touch it |
|---|---|
| `SignOptions.Format` | Which signature family. Decided by what you are signing — see [Signature formats](concepts/signature-formats.md). |
| `SignOptions.Level` | `B`, `T`, `LT`, `LTA`. Anything above `B` needs a `TSPSource`. See [Signature levels](concepts/signature-levels.md). |
| `SignOptions.Packaging` | Enveloped, enveloping or detached — CAdES, XAdES and JAdES only. |
| `SignOptions.TSPSource` | Where time-stamps come from. See [Timestamps and TSAs](guides/timestamps-and-tsas.md). |
| `ValidateOptions.TrustedCertificates` | Trust anchors, as plain certificates. |
| `ValidateOptions.TrustedCertificateSources` | Whole trust stores — this is where the EU trusted lists go, and the only way to get an eIDAS qualification verdict. |
| `ValidateOptions.DetachedContents` | The original file, when the signature is detached from it. |
| `ValidateOptions.ValidationTime` | Ask "was this valid last March?" rather than "is it valid now?". |
| `ValidateOptions.Policy` | Your own constraint set instead of the bundled ETSI one. See [A custom validation policy](guides/custom-validation-policy.md). |

The full field-by-field reference is the godoc:
[`SignOptions`](https://pkg.go.dev/github.com/utain/esig/dss#SignOptions),
[`ValidateOptions`](https://pkg.go.dev/github.com/utain/esig/dss#ValidateOptions).

## The CLI

Same library, no Go code.

```sh
go install github.com/utain/esig/dss/cmd/esig@latest
```

Sign a PDF, then validate what came out:

```sh
export P12PASS=…                     # the password, never on the command line

esig sign sample.pdf \
    -format pades -level B \
    -p12 keystore.p12 -p12-pass env:P12PASS
```

```console
$ esig validate sample-signed-pades-baseline-b.pdf -trust signer.cer
S-C2AE57DE…1A7BD6  TOTAL_PASSED  level=PAdES-BASELINE-B  qualification=NA  signed-by="Go Port Test RSA"
```

Without the trust anchor, the same command tells you exactly what is missing
and exits `1`:

```console
$ esig validate sample-signed-pades-baseline-b.pdf
S-C2AE57DE…1A7BD6  INDETERMINATE (NO_CERTIFICATE_CHAIN_FOUND)  level=PAdES-BASELINE-B  qualification=NA  signed-by="Go Port Test RSA"
    error: The certificate chain for signature is not trusted, it does not contain a trust anchor.
```

The subcommands:

| Command | What it does |
|---|---|
| `esig sign` | Sign a document at a chosen format and level. |
| `esig extend` | Raise an existing signature to a higher level. No key needed. |
| `esig validate` | Validate, and print a summary or render one of the four reports. |
| `esig inspect` | Print what is *in* the document — signatures, time-stamps, certificates — without judging it. |
| `esig report` | Re-render a previously saved SimpleReport XML as the human summary. |
| `esig tl refresh` | Download the EU List of Trusted Lists into a local cache for `validate -tl-cache`. |
| `esig version` | Build information. |

Exit codes are meant for scripts: **0** every signature passed, **1** at least
one did not (or there were none), **2** you got the flags wrong, **3** the
validation could not run at all. `esig <command> -h` prints the full flag list.

## Next

- [Signature formats](concepts/signature-formats.md) — pick the right one.
- [Sign a PDF](guides/sign-a-pdf.md) — the most common single task.
- [Validate against the EU trusted lists](guides/validate-against-eu-trusted-lists.md) —
  the piece that produces a qualification verdict.
