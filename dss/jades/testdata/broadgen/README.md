# Broad differential corpus runner (JAdES)

Differential harness over the **entire** upstream `dss-jades/src/test/resources`
corpus (63 JSON fixtures, 65 signatures), used by the Phase 6 audit. Unlike the
committed 30-fixture cross-validation suite, this walks every upstream JAdES
fixture through both implementations and diffs the dumped JSON — the same role
`../../../pades/testdata/broadgen/` plays for PAdES, where the sampled suite
missed 16 defects the broad run caught.

- `BroadOracle.java` — Java side; compile against a maven-built upstream DSS
  (same classpath recipe as `../gen/CrossValidationOracle.java`), run with the
  corpus directory and an output path as arguments.
- `gobroad_main.go` — Go counterpart of the same dump. Kept under `testdata/`
  so the go tool ignores it; build it by copying it into a throwaway package
  inside the module (`dss/zz_broad/main.go`, `go build ./zz_broad`), or from a
  scratch module with `go mod edit -replace github.com/utain/esig/dss=<repo>/dss`.
- Diff the two JSON files (any structural JSON diff works); every difference is
  a parity defect on one side.

Both sides dump, per signature: signing certificate + SHA-256, claimed signing
time, `getDataFoundUpToLevel`, counter-signature structure (recursively),
reference/signature-integrity verdicts, `getDataToBeSignedRepresentation`, the
`sigD` mechanism, the JWS serialization type, structural-error presence, the
signature algorithm, certificate extraction (source size, signing/complete/
attribute cert-ref counts, `xVals` count), revocation extraction (CRL/OCSP
binary and ref counts), signature scopes, and the full time-stamp inventory per
bucket (content / signature / X1 / X2 / archive) with each token's DSS-Id,
generation time, message-imprint found/intact/hex, token signature intactness,
and its sorted set of timestamped references.

## Result (Phase 6 audit)

62 of 63 files byte-identical. The one remaining difference is
`validation/jades-with-double-sigt.json`, where both implementations reject the
file for the same reason (duplicate `sigT` in the protected header) but word the
exception differently — the committed suite asserts that one case by
case-insensitive substring, deliberately.

Two defects were found here that no column of the sampled suite could see, both
fixed in `dss/jades/jades_timestamp_source.go`:

1. **Signature time-stamps covered one certificate too few** (14 fixtures).
   Java's `SignatureTimestampSource.makeTimestampTokensFromUnsignedAttributes`
   passes `getSignatureTimestampReferences()` into the token, and
   `JAdESTimestampSource` overrides that method to fold in
   `getKeyInfoReferences()`. The Go base's counterpart is unexported and has no
   override hook, so the KeyInfo certificate reference was silently dropped.

2. **Every content time-stamp had a different DSS-Id than upstream's**
   (3 time-stamps, all of the corpus's content time-stamps).
   `SignatureTimestampIdentifierBuilder` mixes the carrying attribute's position
   among the signature properties into the identifier, and the base finds that
   position with Go `==` (pointer identity) where Java uses
   `SignatureAttribute.equals`. JAdES's signed-properties adapter re-boxed its
   attributes on every `Attributes()` call, so the lookup always missed and the
   order came out nil.

Both are now pinned by the committed suite: `../gen/CrossValidationOracle.java`
dumps the time-stamp inventory including `dssId` and `timestampedReferences`,
three content/X1/X2 fixtures were added to its `FILES`, and
`jades_upstream_cross_validation_test.go` asserts every column. Reverting either
fix fails that suite (22 and 9 assertions respectively).

## Related finding, NOT fixed here

`CAdESSigProperties.Attributes()` and `XAdESSigProperties.Attributes()` re-allocate
their attribute objects on every call exactly the way the JAdES adapter did, and
`cades_timestamp_source.go` / `xades_timestamp_source.go` both feed
`GetAttributeOrder` into their time-stamp identifier builders. The same defect (2)
is therefore latent in those formats. It was left alone because they are frozen
for this phase and a fix there needs its own broad differential run to confirm —
XAdES in particular wraps a mutable DOM, so a cache would need invalidation that
the JAdES signed properties (an immutable protected header) do not.
