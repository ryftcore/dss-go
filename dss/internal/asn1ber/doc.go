// Package asn1ber is the BER/DER/DL engine the DSS port uses in place of BouncyCastle's
// ASN.1 layer (org.bouncycastle.asn1.*).
//
// Extracted from dss-spi/src/main/java/eu/europa/esig/dss/spi/DSSASN1Utils.java (DSS 6.5.RC1),
// where the Go port of that class first hosted it. Two conventions come from that origin:
//
//   - Java passes around ASN1Encodable/ASN1Primitive objects; this package represents an
//     ASN.1 value by its encoding ([]byte), because that is what the surrounding code
//     digests, compares and serialises. Wherever identity or a digest depends on the bytes,
//     the original bytes are carried through untouched and never re-encoded - see
//     Element.Encoded.
//   - The BER/DER/DL re-encoding BouncyCastle performs in getEncoded(String) is implemented
//     here (Element.DEREncoded, Element.DLEncoded, Element.BEREncoded): DER output must be
//     byte-exact, since it feeds signature verification and archive-timestamp hashing.
//
// The package is deliberately dependency-free: it imports the standard library only, so that
// it can be shared by every DSS package without dragging the DSS object model along.
package asn1ber
