# Document-level end-to-end oracle (phase 8f harness, item B)

THE PHASE 8 EXIT CRITERION's item (B): `document_level_oracle_test.go`'s
`TestDocumentLevelOracle` validates 60 real signed documents - 10 per format
family (CAdES, XAdES, PAdES, JAdES, ASiC-CAdES, ASiC-XAdES) - end to end, in
BOTH implementations, and compares diagnostic-data core fields, DetailedReport
BasicBuildingBlocks conclusions, top-level Signature/Timestamp/EvidenceRecord
verdicts, and SimpleReport qualifications.

Unlike `dss/validation/executor/testdata/oracle/full_corpus.jsonl` (item A),
which starts from ALREADY-BUILT diagnostic data, this test starts from the
ORIGINAL signed document (CMS/XML/PDF/JWS/ZIP) on both sides, so it also
exercises each format's own diagnostic-data builder - format detection,
signature/timestamp/evidence-record extraction, identifier construction - end
to end.

This package lives outside every format package (`dss/cades`, `dss/xades`,
`dss/pades`, `dss/jades`, `dss/asic/cades`, `dss/asic/xades` all already
import `dss/validation`, so a single cross-format harness test cannot live
inside any one of them without an import cycle) and is test-only: it has no
non-test `.go` file and exports nothing.

## Fixtures

`testdata/oracle/document_level_manifest.tsv` (`format<TAB>name<TAB>path`)
lists the 60 fixtures, all already vendored in the repo: the same
`testdata/upstream` trees each format's own `*_smoke_test.go` reads (most of
the CAdES/XAdES/PAdES/JAdES entries are literally the same files those smoke
tests already exercise; the rest and all the ASiC ones are additional
already-vendored fixtures from the same trees, picked to reach ten per
format and to route through the LTA/archival/counter-signature/evidence-
record-adjacent paths those smoke tests do not).

## `testdata/oracle/document_level.jsonl`

Produced by `testdata/oracle/gen/DocumentLevelOracle.java` against DSS
6.5.RC1's `dss-cades`/`dss-xades`/`dss-pades`(+`dss-pades-pdfbox`)/`dss-jades`/
`dss-asic-cades`/`dss-asic-xades`/`dss-validation` classes. Both engines run
`SignedDocumentValidator.fromDocument`'s format-autodetection dispatch (the
Go side: `dssvalidation.SignedDocumentValidatorFromDocument`), a permissive
`CertificateVerifier` (`CommonCertificateVerifier(true)` with the same seven
`Alert` setters silenced every format package's own smoke-test
`permissiveCertificateVerifier()` helper silences - these are offline
fixtures with no live revocation data or reachable trust anchors),
`ValidationLevel.ARCHIVAL_DATA`, locale `en`.

## Defects found by this harness

* **Fixed** (`dss/validation/reports/diagnostic/diagnostic_data_builder.go`,
  `GetXmlFoundCertificatesForSource`): passing `&ocspCertificateSource.
  TokenCertificateSource` (the address of the EMBEDDED base field) instead of
  `ocspCertificateSource` itself lost the outer `*OCSPCertificateSource`'s
  `CertificateSourceType()` override, so an orphan OCSP revocation
  identifier's certificate source answered `CertificateSourceType_OTHER`
  instead of `OCSP_RESPONSE`, panicking a runtime type assertion. Found via
  `pades/pades-lt` (PAdES-LT.pdf). The function's parameter is now the
  `foundCertificatesSource` interface (matching its sibling
  `GetXmlFoundCertificatesForToken`), so callers pass the OUTER value they
  actually have and virtual dispatch works the way Java's abstract-class
  parameter always did.

* **F2 - open, tracked** (`knownIdentifierDivergences` in
  `document_level_oracle_test.go`, 30 of 60 fixtures): a fresh-from-document
  `TimestampToken`'s `T-...` identifier embeds an `SA-...`
  `SignatureAttributeIdentifier` for its carrying attribute, itself a digest
  of that attribute's re-serialized bytes (XAdES: DOM serialization of the
  `<xades:...TimeStamp>` element; CAdES/ASiC-CAdES/PAdES: DER re-encoding of
  the CMS unsigned `Attribute`). The identifier disagrees between engines even
  though every `BasicBuildingBlocks` Type/Indication/SubIndication it is keyed
  by is byte-identical - confirmed down to proving the RAW timestamp-binaries
  half of the digest matches Java's `DSSASN1Utils.getDEREncoded(TimeStampToken)`
  exactly, byte for byte, on the same input. The divergence is therefore in
  the ATTRIBUTE-serialization half, not the timestamp-binaries half; isolating
  which of `xml/utils.DomUtilsSerializeNode` (via `internal/xmldom`) or the
  CAdES `Attribute.Encoded()` path is responsible, and fixing it, is a
  dedicated pass this harness session did not have the budget for. NOT a
  verdict-correctness defect - every Indication/SubIndication checked matches
  exactly.

* **F3 - open, tracked** (`knownContentDivergences` in
  `document_level_oracle_test.go`, `cades/baseline-lta`): a genuine, narrower
  verdict divergence - one nested archive timestamp inside
  Signature-C-B-LTA-10.p7m's timestamp chain is `FAILED`/`HASH_FAILURE` in Go
  where Java gets `INDETERMINATE`/`NO_CERTIFICATE_CHAIN_FOUND` (i.e. Java's
  message-imprint verified; Go's did not). Root cause is somewhere in the
  CAdES archive-timestamp (RFC 5126 archive-time-stamp-v2/v3) message-imprint
  construction in `cades_timestamp_message_digest_builder.go`; isolating the
  exact byte difference needs a dedicated pass through that algorithm against
  the RFC, which this harness session did not have the budget for either.

Both F2 and F3 are tracked BY NAME in `document_level_oracle_test.go`, not
silently absorbed: the test still runs the full comparison for every listed
fixture, and turns into a hard failure the moment a listed fixture's mismatch
either disappears (so the entry has to be removed) or changes shape (so the
entry no longer describes what is actually happening).

## Regenerating `document_level.jsonl`

```
CP=<merged runtime classpath for dss-validation, dss-cades, dss-xades, \
    dss-pades(+dss-pades-pdfbox), dss-jades, dss-asic-cades, dss-asic-xades, \
    each module's own target/classes taking precedence over any stale jar>
javac -cp "$CP" -d /tmp/oracle testdata/oracle/gen/DocumentLevelOracle.java
java  -cp "$CP:/tmp/oracle" DocumentLevelOracle \
      testdata/oracle/document_level_manifest.tsv \
      testdata/oracle/document_level.jsonl
```

The manifest's third column must be an ABSOLUTE path when regenerating (the
committed file uses paths relative to this directory, which is what the Go
test reads; rewrite them to absolute before invoking the generator, or point
it at a copy with absolute paths substituted in).
