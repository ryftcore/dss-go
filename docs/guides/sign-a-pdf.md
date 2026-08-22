# Sign a PDF

The single most common task, end to end: a PDF in, a PAdES-signed PDF out,
verified.

## What you need first

A signing key and its certificate, in a **PKCS#12** file. If you have a
`.p12` or `.pfx` and a password, you are ready. If your certificate came as
separate PEM files, `openssl pkcs12 -export` will bundle them.

!!! danger "Never put the password on the command line"
    It ends up in your shell history and in the process table. The CLI takes
    the *name of an environment variable*, not the password itself:
    `-p12-pass env:P12PASS`.

## From Go

```go
package main

import (
	"log"

	"github.com/ryftcore/dss-go/dss"
)

func main() {
	signer, err := dss.OpenPKCS12("keystore.p12", "password")
	if err != nil {
		log.Fatal(err)
	}
	defer signer.Close()

	doc, err := dss.OpenDocument("contract.pdf")
	if err != nil {
		log.Fatal(err)
	}

	signed, err := dss.Sign(doc, signer, dss.SignOptions{
		Format: dss.FormatPAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := signed.Save("contract-signed.pdf"); err != nil {
		log.Fatal(err)
	}
}
```

`OpenDocument` reads lazily, so a large PDF is not slurped into memory. The
returned `signed` document is not written anywhere until you `Save` it (or read
it as a stream).

That is level **B**. It proves who signed and that nothing changed since — for
as long as the certificate remains valid and checkable. For anything that has
to outlive that, add a time-stamp.

## From the CLI

```sh
export P12PASS='…'

esig sign contract.pdf \
    -format pades -level B \
    -p12 keystore.p12 -p12-pass env:P12PASS \
    -out contract-signed.pdf
```

Without `-out`, the output is named after the input the way the library names
it — `contract-signed-pades-baseline-b.pdf`.

## Adding a time-stamp (level T)

```go
signed, err := dss.Sign(doc, signer, dss.SignOptions{
	Format:    dss.FormatPAdES,
	Level:     dss.LevelT,
	TSPSource: tsa, // see the timestamps guide
})
```

```sh
esig sign contract.pdf -format pades -level T \
    -p12 keystore.p12 -p12-pass env:P12PASS \
    -tsa https://tsa.example.org/tsa
```

The library itself ships **no HTTP time-stamping client** — the CLI has one of
its own. From Go you either supply a `TSPSource` implementation or use the
port's local, key-based one. Both routes, and why it is this way, are in
[Timestamps and TSAs](timestamps-and-tsas.md).

## Check what you produced

Always. Signing succeeding is not the same as the result validating.

```sh
esig validate contract-signed.pdf -trust your-ca.cer
```

```console
S-C2AE57DE…  TOTAL_PASSED  level=PAdES-BASELINE-B  qualification=NA  signed-by="…"
```

`qualification=NA` here is expected — you anchored the chain with your own CA,
which says nothing about eIDAS qualification. See
[Trust, eIDAS and qualification](../concepts/trust-and-eidas.md).

Or from Go:

```go
reports, err := dss.Validate(signed, dss.ValidateOptions{
	TrustedCertificates: signer.CertificateChain(),
})
if err != nil {
	log.Fatal(err)
}
log.Println(reports.Valid(), reports.Verdicts()[0].Indication)
```

(Trusting the signer's own chain is only sensible in a test. In production the
trust anchor is a CA you decided to believe, or the EU trusted lists.)

## Signing a PDF that is already signed

Sign it again. PAdES adds an **incremental update**: the previous revision, and
the previous signature, stay byte-for-byte intact, and the new signature covers
the file including them.

```go
twiceSigned, err := dss.Sign(signed, secondSigner, dss.SignOptions{
	Format: dss.FormatPAdES,
	Level:  dss.LevelB,
})
```

Validation then reports both signatures, each with its own scope — the first
covering the original revision, the second covering everything. That
`SignatureScope` field is how a reader tells "signed the whole thing" from
"signed an earlier version of the thing".

## Things that will surprise you

**`SignOptions.Packaging` is ignored for PDFs.** A PDF signature is always
embedded; there is no enveloped/detached choice to make.

**No visible signature appearance.** The port creates the signature field and
its rectangle but paints nothing into it — no image, no text block. There is no
rasteriser in the native PDF engine. If your requirement includes a visible
stamp on the page, this port cannot produce it today; see
[Known gaps](../compatibility/known-gaps.md).

**An encrypted PDF is fine to read.** The native engine handles RC4 and AES
encryption and re-encrypts on write.

**Level LT needs revocation data, which means network sources you supply.**
Signing at LT requires actually fetching CRLs or OCSP responses to embed. The
library ships no HTTP client for that; set the sources you want on a
`CertificateVerifier` and pass it in `SignOptions.CertificateVerifier`.

## Next

- [Timestamps and TSAs](timestamps-and-tsas.md) — everything above level B.
- [Signature levels](../concepts/signature-levels.md) — how far up to go.
- [How validation works](../concepts/validation.md) — reading the result.
