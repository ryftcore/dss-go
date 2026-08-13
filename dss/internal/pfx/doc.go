// Package pfx is a minimal, dependency-free reader for PKCS#12 (PFX) key stores, replacing
// golang.org/x/crypto/pkcs12 for the two capabilities that package does not offer (see
// dss/token/key_store_signature_token_connection.go's file header for why they are needed):
// decrypting a SafeBag without discarding its private key's type, and parsing every PKCS#8 key
// type DSS's own fixtures use - RSA, EC and Ed25519 through crypto/x509, and DSA (which neither
// x509.ParsePKCS8PrivateKey nor crypto/dsa itself parses) by hand.
//
// It has no Java class to mirror - java.security.KeyStore's PKCS12 provider is a JCA SPI
// implementation, not a class this port translates - so, as PORTING.md prescribes for such
// machinery, it lives under internal/ and states its provenance here rather than claiming a
// "Ported from" source file.
//
// # Scope
//
// RFC 7292 defines a general-purpose container; this package reads exactly the subset openssl
// and the JDK's PKCS12 provider actually produce, which is what every DSS PKCS#12 fixture is:
//
//   - a version-3 PFX PDU whose authSafe is the plain (unencrypted) "data" content type, i.e.
//     password integrity rather than public-key (signed) integrity;
//   - a MacData integrity check computed with the SHA-1/SHA-224/SHA-256/SHA-384/SHA-512 flavour
//     of the RFC 7292 Appendix B.2 HMAC construction (openssl's classic "-macalg SHA1" default
//     and the SHA-256 default of modern openssl/the JDK);
//   - an AuthenticatedSafe holding "data" (unencrypted SafeContents) and/or "encryptedData"
//     (PKCS#7 EncryptedData) ContentInfos;
//   - SafeBags of type keyBag (cleartext PKCS#8 PrivateKeyInfo), pkcs8ShroudedKeyBag
//     (EncryptedPrivateKeyInfo) and certBag (a bare X.509 certificate);
//   - privacy encryption under either the RFC 7292 Appendix B legacy PBE schemes
//     (pbeWithSHAAnd3-KeyTripleDES-CBC, pbeWithSHAAnd2-KeyTripleDES-CBC,
//     pbeWithSHAAnd128/40BitRC2-CBC, pbeWithSHAAnd128/40BitRC4 - openssl's "-legacy" mode and
//     every JDK before the PBES2 default) or PKCS#5 v2.0 PBES2 with a PBKDF2 key derivation
//     function and an AES-CBC (128/192/256) or DES-EDE3-CBC encryption scheme (openssl's and
//     the JDK's modern default).
//
// A signed-integrity PFX, an unrecognised bag or algorithm, or anything else RFC 7292 allows
// but these two producers do not emit is rejected with an error rather than silently ignored.
//
// # Dependencies
//
// The standard library and internal/asn1ber only, per PORTING.md's dependency policy - in
// particular RC2 (crypto/rc2 was never added to the standard library) is implemented from
// RFC 2268 in this package rather than pulled in as a new module dependency.
package pfx
