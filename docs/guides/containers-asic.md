# ASiC containers

**ASiC** — Associated Signature Container — is what you use when the thing being
signed is not one file. A submission of "the form, three attachments and a
covering letter", signed as a unit, is an ASiC container.

Underneath it is a ZIP file with rules. That is genuinely all it is, and it is
worth knowing, because it means `unzip` is a legitimate debugging tool.

## S or E?

| | **ASiC-S** (simple) | **ASiC-E** (extended) |
|---|---|---|
| Signed data objects | exactly one | many |
| Manifest describing what is signed | no | yes |
| Use it for | wrapping a single file tidily, with its signature and any time-stamps | a multi-file submission signed as a whole |
| Constant | `dss.ContainerASiCS` | `dss.ContainerASiCE` |
| CLI `-format` | `asics` | `asice` |

Inside either, the signatures are **CAdES** or **XAdES** — chosen with
`dss.FormatASiCWithCAdES` / `dss.FormatASiCWithXAdES`, or the CLI's
`-asic-format cades|xades`. Which one is usually dictated by whoever receives
the container; if nobody has said, XAdES is the more common default in EU
schemes, and it is what the CLI uses when `-asic-format` is omitted.

## Signing several documents into one container

```go
docs := []dss.Document{
	dss.NewDocument("payload.bin", []byte("the payload being signed")),
	dss.NewDocument("metadata.json", []byte(`{"amount":1250.00,"currency":"EUR"}`)),
}

container, err := dss.SignMultiple(docs, signer, dss.SignOptions{
	Format: dss.FormatASiCWithXAdES,
	Level:  dss.LevelB,
})
```

`SignMultiple` is the entry point; `Sign` is the same thing for one document.
Only the two ASiC formats accept more than one — anything else returns
`ErrMultipleDocuments`. `ContainerType` defaults to `ContainerASiCE` for several
documents and `ContainerASiCS` for one, so the common cases need no explicit
setting.

**The document names matter.** They become the entry names inside the ZIP, they
appear in the manifest, and they are what the reports show. `NewDocument`'s
first argument is not decoration.

## What comes out

Opening the result as an ordinary ZIP archive — this is the actual output of
[`examples/06-asice-container`](https://github.com/utain/esig/tree/main/dss/examples/06-asice-container):

```console
container: container-signed-xades-baseline-b.sce
entries:
  mimetype (31 bytes)
  payload.bin (24 bytes)
  metadata.json (35 bytes)
  META-INF/signatures001.xml (4202 bytes)
  META-INF/manifest.xml (481 bytes)
verdict: XAdES-BASELINE-B, indication=TOTAL_PASSED
```

Three rules are visible there:

1. **`mimetype` comes first and is stored uncompressed.** This is what lets a
   reader identify the container type by looking at a fixed byte offset,
   without unzipping. Get it wrong and tools reject the container.
2. **Signatures live under `META-INF/`.** `signatures001.xml` for XAdES,
   `signature001.p7s` for CAdES; numbered, so several signers can each add
   their own.
3. **`META-INF/manifest.xml`** binds the signature to what it covers — the
   piece that makes ASiC-E more than a zip with a signature dropped in.

## Validating

Exactly like anything else — format detection recognises the container:

```go
reports, err := dss.Validate(container, dss.ValidateOptions{
	TrustedCertificates: anchors,
})
```

```sh
esig validate submission.asice -trust ca.cer
esig inspect submission.asice
```

You do **not** need `-detached`: the signed files are inside the container, and
the validator finds them there.

The reports carry container-specific facts worth reading: the container type,
the manifest contents, and — importantly — **whether every file in the container
is actually signed**. A container can legitimately hold unsigned entries. The
default policy raises that as a warning (`AllFilesSigned`); whether it should be
a hard failure for you is a policy decision. See
[A custom validation policy](custom-validation-policy.md).

## Adding a signature to an existing container

Sign the container again. Each signer's signature lands as its own numbered
entry under `META-INF/`, and validation reports one verdict per signature.

## Extending

Same as any other format — no key needed:

```sh
esig extend submission.asice -format asice -asic-format xades -level LTA \
    -tsa https://tsa.example.org/tsa
```

For long-term retention this is the normal shape: the container arrives at B or
T, gets extended to LT on ingest while the revocation data is still available,
and is re-timestamped to LTA on a schedule. See
[Signature levels](../concepts/signature-levels.md).

## Things worth knowing

**File extensions.** `.asice` / `.sce` for extended, `.asics` / `.scs` for
simple. The library names its output from the format and level; `-out` overrides.

**Zip-bomb defences are on.** The extractor enforces limits on entry counts and
expansion ratios. A hostile container is a real attack surface for anything that
unzips untrusted input, and the port's guards for this are mutation-tested.

**Containers can be merged.** Two containers covering related material can be
combined, with the merge rules the standard defines. That path lives in the
`asic` package rather than the facade — see
[pkg.go.dev](https://pkg.go.dev/github.com/utain/esig/dss/asic).

**ASiC-S with a time-stamp only.** A container can carry a time-stamp over its
single data object rather than a signature — proving existence without claiming
authorship. Also an `asic` package path.

## Next

- [Signature formats](../concepts/signature-formats.md) — why ASiC rather than
  a bare signature.
- [Signature levels](../concepts/signature-levels.md) — how long it has to last.
