# Porting conventions (Java DSS → Go)

Binding rules for every ported package, and for every change to one.
Reviewers reject deviations.

## Layout & naming

- One Go file per Java class: `SignatureLevel.java` → `signature_level.go`; its tests in `signature_level_test.go`.
- Java package → Go package per the mapping table in `/docs/migrating-from-java/index.md`. Keep upstream file organization recognizable so diffs against upstream stay tractable.
- Exported Go identifiers keep the Java name (minus `get`/`set` prefixes where un-idiomatic): `getSignatureAlgorithm()` → `SignatureAlgorithm()`.
- **An exported identifier never repeats its package name.** The qualifier at the call site already carries it, so the Java name loses that prefix: `ASiCContent` → `asic.Content`, `XAdESSignature` → `xades.Signature`, `AlertHandler` → `alert.Handler`, `CAdESService` → `cades.Service`. The rule is the one `revive`'s `exported` check applies: if the package name is a case-insensitive prefix of the identifier and the next rune is `_` or upper-case, drop the prefix. It leaves a name equal to the package name alone, and — deliberately — a name whose next rune is a digit, so the XAdES version numbers survive (`xades.XAdES111XSDUtils`). The Java class name is always recoverable by putting the package name back.
- Constructors and companions follow their type: `NewASiCContent` → `asic.NewContent`, `AbstractASiCContentBuilder` → `asic.AbstractContentBuilder`, `ASiCManifestElementSigReference` → `asic.ManifestElementSigReference`. Where the type name sits in the middle of another identifier, leave it: `lote.LoLoTEIdentifier` is *list of* LoTE, not a stutter.
- Never let a renamed identifier reach a wire format. If a rendered string is derived from a Go type name, pin it: `model/eaa/claim` puts the dropped `Claim` prefix back before rendering `toString()`, because that text is cross-validated against Java.
- Each file starts with a comment naming its upstream source: `// Ported from dss-enumerations/.../SignatureLevel.java (DSS 6.5.RC1).`
- Machinery that BouncyCastle provides upstream has no Java class to mirror, so it lives under `internal/` instead of being duplicated per package: `internal/asn1ber` holds the BER/DER/DL engine and the generic X.509 structures (`AlgorithmIdentifier`, `IssuerSerial`, `GeneralName`) extracted from `DSSASN1Utils`. `internal/cmscore` builds on it with the RFC 5652 CMS and RFC 3161 time-stamp structures that replace `org.bouncycastle.asn1.cms.*`, `org.bouncycastle.cms.*` and `org.bouncycastle.tsp.*`; the public `cms` package wraps it. Such a package states its provenance in `doc.go` and imports the standard library plus, at most, the other `internal/` packages it is layered on - never a DSS package, since the DSS packages are what import it. Packages that exposed those types before the extraction keep their API through type aliases.

## Enums

Java enums become typed string constants whose **value is exactly Java's `name()`** (serialization compatibility). Because Go constants share one package namespace while Java's are class-scoped, the type name is a prefix; the constant is then spelled in Go **MixedCaps**, never `Type_JAVA_NAME`:

```go
type SignatureLevel string

const (
    SignatureLevelXAdESBaselineB SignatureLevel = "XAdES_BASELINE_B"
    ...
)

type DigestAlgorithm string

const (
    DigestAlgorithmSHA256 DigestAlgorithm = "SHA256"
    ...
)
```

The **value is never touched** by this spelling rule. Enum values, XML element and attribute names, i18n message keys, OIDs, URIs, MIME types and JSON keys stay byte-identical to upstream; only the Go identifier is Go-shaped.

Naming a constant, given Java's `TYPE` + `NAME`: split `NAME` on `_`, then join the parts with no separator, rendering each part as follows.

- A part that is already mixed-case is kept verbatim, with its first rune upper-cased: `XAdES`, `CAdES`, `PAdES`, `JAdES`, `ASiC`, `AdES`, `signingCertificate` → `SigningCertificate`.
- A purely numeric part is kept verbatim and simply concatenated: `XAdES_141` → `XAdES141`, `P_256` → `P256`.
- An all-caps part that is a **crypto or standards initialism** stays upper-case: `SHA1 SHA3 SHA256 SHA384 SHA512 SHAKE256 MD5 RIPEMD160 RSA RSASSA PSS ECDSA DSA EDDSA HMAC MGF1 ASN1 DER PEM BER CER CMS CRL OCSP TSA TSL TSP TST TL LOTL MRA PKI PKCS PKCS7 PKCS12 SPKI X509 X400 X448 ED25519 XML XMLDSIG XPATH XSD JSON JWS JWT JOSE PDF PNG SVG JPEG MIME HTML HTTP HTTPS URI URL URN UID UUID UTC ID DN CA CN AIA OID EU ETSI IEC ISO RFC W3C QC QSCD QTSP QTSA QWAC PSD2 EEA LDAP FTP TLS SSL GCM MAC B64 UTF8 SECP` — plus the ETSI validation-process block codes (`BBB XCV CV SAV FC ICS PSV VCI ADEST …`).
- Any other all-caps part is Titlecased word by word: `BASELINE` → `Baseline`, `FOR` → `For`, `SIGNING` → `Signing`, `COMMONNAME` → `CommonName`, `ESIGNATURES` → `ESignatures`.
- An all-caps part that is an opaque domain code with no word structure stays upper-case: `MessageTag_BBB_XCV_CCCBB` → `MessageTagBBBXCVCCCBB`.
- In the `OID_*` families, the lowercase ASN.1 arc labels are Titlecased, not initialism-uppercased, because the alternative is unreadable: `OID_id_aa_ets_archiveTimestampV3` → `OIDIdAaEtsArchiveTimestampV3`.

`Test*`/`Example*` function names keep their underscores — that is idiomatic Go (sub-test and example association), not Java naming.

- Instance methods → methods on the type; multi-field enums (OID, URI, code…) back their methods with package-level lookup tables defined next to the constants.
- Static factories keep their contract: `valueOf` → `SignatureLevelValueOf(string) (SignatureLevel, error)`; `forOid`/`forUri`/`forName` → `SignatureAlgorithmForOID(...)` etc. Java's thrown `IllegalArgumentException` becomes a returned error.
- Marker interfaces (`OidBasedEnum`, `UriBasedEnum`, …) → Go interfaces in `enumerations`.
- **Never invent, abbreviate, or "fix" an OID, a URI, or an enum VALUE.** Copy them from upstream verbatim. This is about the wire: the Go *identifier* is MixedCaps per the rule above, but the string it is bound to is upstream's, byte for byte.

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
- The module requires Go 1.27, so a Java **generic instance method** ports to a Go generic method on the receiver (`func (w *CertificateWrapper) CertificateExtensionForOid[T …](oid string) T`); Java `protected` maps to an unexported method when every caller is in-package, otherwise it stays exported (Go has no `protected`).
- A **generic method declared on an interface** cannot be ported: Go forbids type parameters on interface methods, Go 1.27 included. Such methods erase to the constraint's base type (`<T extends AdvancedSignature> … (Collection<T>)` → `([]AdvancedSignature)`) and the erasure is permanent.
- A **stateless Java "protected helpers for subclasses" class** (no fields, every method pure w.r.t. `this`) ports to package-level functions, generic or not — not to methods on an empty receiver. `spi/validation/timestamp` is the reference case.

## Tests

- Port test *vectors*, not JUnit code. Upstream resources are copied into `testdata/` mirroring their upstream path.
- Every registry-like table (OIDs, URIs, algorithm mappings) gets an exhaustive table test.
- Gate for every batch: `go build ./... && go vet ./... && go test ./...`.
- Keep in-module `testdata/` small: a heavy fixture (an oracle dump, a large corpus, a bulky KAT) belongs in the repo-root `corpus/` tree instead, mirroring the package's path under `dss/`. Reach it at test time with `internal/corpustest.Path`/`RootPath`, which resolve into `corpus/` and skip the test gracefully when it is absent (a bare module checkout has no `corpus/`); keep a small representative subset directly under the package's own `testdata/` so `go test` still exercises real code from the module zip alone.

## Dependency policy

Stdlib-first. Allowed: `golang.org/x/…`. Anything else needs tech-lead sign-off recorded in this file.
