// Package cmscore is the RFC 5652 (CMS) and RFC 3161 (time-stamp) engine the DSS port uses in
// place of BouncyCastle's CMS and TSP layers (org.bouncycastle.asn1.cms.*,
// org.bouncycastle.cms.*, org.bouncycastle.asn1.tsp.*, org.bouncycastle.tsp.*).
//
// It has no Java class to mirror, so - as PORTING.md prescribes for BouncyCastle-replacement
// machinery - it lives under internal/ and states its provenance here. Its consumers are
// dss/cms (the port of the dss-cms module and its CMS interface), dss/cades and the
// time-stamp handling of dss/spi.
//
// # Parse and build
//
// Parsing is BER-tolerant, because that is what real CAdES signatures require: an
// indefinite-length ContentInfo, SignedData or EncapsulatedContentInfo, and an eContent split
// into a constructed OCTET STRING, are all accepted, and are what a streaming producer such as
// "openssl cms -sign -stream" or BouncyCastle's CMSSignedDataStreamGenerator emits.
//
// Building is DER only, and only for SignedData: CAdES signing has to produce a signature, and
// a signature is defined over DER. Building a TimeStampToken is a TSA's job, not a validator's,
// so RFC 3161 is parse-only here.
//
// # Preserved bytes
//
// Wherever identity, a digest or re-serialisation depends on the exact input octets, those
// octets are carried through untouched rather than re-encoded, and are reachable through the
// Encoded method of the structure that owns them (and through Element, which hands out the
// asn1ber view for callers that need a different encoding of the same subtree). That covers
// the whole CMS, the SignedData and each of its four re-encodable fields
// (digestAlgorithms, certificates, crls, signerInfos), the encapContentInfo, every SignerInfo,
// its signed and unsigned attributes, every certificate, CRL and OCSP response, and the whole
// TimeStampToken. Downstream DSS identifiers and archive-timestamp message imprints digest
// exactly these bytes.
//
// A signed attribute set is the one place where two encodings are both legitimate and both
// needed: SignerInfo.SignedAttributesRaw returns the [0] IMPLICIT element as it was received,
// while SignerInfo.SignedAttributesDER returns the DER SET OF re-encoding that RFC 5652
// clause 5.4 makes the signature input.
//
// # Dependencies
//
// The standard library and internal/asn1ber only. In particular the DSS object model is not
// imported: this package deals in OIDs, big integers and raw DER, and leaves the mapping onto
// DigestAlgorithm, CertificateToken and friends to its callers.
//
// # Checked against BouncyCastle
//
// Every field below is compared with BouncyCastle's reading of the same bytes by
// TestBouncyCastleOracle, and every DER this package builds with the DER BouncyCastle builds
// from the same inputs by TestBuilderMatchesBouncyCastle; the expectations in testdata are
// BouncyCastle's output, produced by the generators under testdata/gen. Three differences are
// deliberate and survive that comparison:
//
//   - An OBJECT IDENTIFIER whose arc does not fit a Go int is rejected, where BouncyCastle
//     reads it. encoding/asn1 cannot represent one, and an OID that cannot be represented can
//     never match anything, so the document is refused rather than half-understood.
//
//   - A RevocationInfoChoices member that is neither a CertificateList nor the [1] IMPLICIT
//     OtherRevocationInfoFormat is rejected, where BouncyCastle keeps it and lets the caller
//     ignore it. RFC 5652 clause 10.2.1 defines no third alternative.
//
//   - A CRL, an OCSP response and an OtherRevocationInfoFormat value are handed out as they
//     arrived. BouncyCastle's own accessors re-encode them to DER, which for a value that did
//     not arrive in DER changes the octets its signature was computed over; preserving them is
//     the point of this package.
package cmscore
