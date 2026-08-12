# Porting conventions (Java DSS → Go)

Binding rules for every ported package. Reviewers reject deviations.

## Layout & naming

- One Go file per Java class: `SignatureLevel.java` → `signature_level.go`; its tests in `signature_level_test.go`.
- Java package → Go package per the mapping table in `/PORTING_PLAN.md`. Keep upstream file organization recognizable so diffs against upstream stay tractable.
- Exported Go identifiers keep the Java name (minus `get`/`set` prefixes where un-idiomatic): `getSignatureAlgorithm()` → `SignatureAlgorithm()`.
- Each file starts with a comment naming its upstream source: `// Ported from dss-enumerations/.../SignatureLevel.java (DSS 6.5.RC1).`

## Enums

Java enums become typed string constants whose **value is exactly Java's `name()`** (serialization compatibility):

```go
type SignatureLevel string

const (
    XAdES_BASELINE_B SignatureLevel = "XAdES_BASELINE_B"
    ...
)
```

- Instance methods → methods on the type; multi-field enums (OID, URI, code…) back their methods with package-level lookup tables defined next to the constants.
- Static factories keep their contract: `valueOf` → `SignatureLevelValueOf(string) (SignatureLevel, error)`; `forOid`/`forUri`/`forName` → `SignatureAlgorithmForOID(...)` etc. Java's thrown `IllegalArgumentException` becomes a returned error.
- Marker interfaces (`OidBasedEnum`, `UriBasedEnum`, …) → Go interfaces in `enumerations`.
- **Never invent, abbreviate, or "fix" an OID, URI, or enum name.** Copy them from upstream verbatim.

## Errors, exceptions, alerts

- Checked/runtime exceptions → `error` returns. Exception classes that carry meaning (`DSSException`, `IllegalInputException`, …) → error types in the owning package, matched with `errors.As`/`errors.Is`.
- Java `throw` in constructors → constructor funcs returning `(T, error)`.
- `dss-alert` handlers keep their semantics: alerts receive a status and decide to log/throw; in Go they receive the status and may return an error.

## Streams & documents

- `DSSDocument` → interface with `OpenStream() (io.ReadCloser, error)`, `Name() string`, `MimeType() MimeType`, `Digest(DigestAlgorithm) (string, error)` — implementations: `InMemoryDocument`, `FileDocument`, `DigestDocument`.
- `InputStream`/`OutputStream` plumbing → `io.Reader`/`io.Writer`; no whole-file `[]byte` slurping where upstream streams.

## Crypto

- BouncyCastle → Go stdlib (`crypto/*`, `encoding/asn1`) plus `golang.org/x/crypto/cryptobyte` for ASN.1 that `encoding/asn1` can't round-trip byte-exactly. No cgo, no unmaintained third-party crypto.
- DER output must be byte-exact where upstream's is (signed attributes, c14n, etc.). When Java produces BER quirks we normalize only where the spec allows, and the deviation is documented in the file header.

## Collections & generics

- `List<T>`→`[]T`, `Map<K,V>`→`map[K]V` (order-sensitive upstream iteration → slice of pairs or explicit sort), `Set<T>`→`map[T]struct{}` behind small helpers in `utils`.

## Tests

- Port test *vectors*, not JUnit code. Upstream resources are copied into `testdata/` mirroring their upstream path.
- Every registry-like table (OIDs, URIs, algorithm mappings) gets an exhaustive table test.
- Gate for every batch: `go build ./... && go vet ./... && go test ./...`.

## Dependency policy

Stdlib-first. Allowed: `golang.org/x/…`. Anything else needs tech-lead sign-off recorded in this file.
