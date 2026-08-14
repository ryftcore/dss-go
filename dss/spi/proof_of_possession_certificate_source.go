// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/ProofOfPossessionCertificateSource.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: this file references CertificateRef, a spi.x509 type flattened into
// this package and ported in a sibling chunk of phase 2a; see certificate_ref_identifier.go
// for its assumed shape. Java's Set<CertificateRef> is ported as []*CertificateRef, matching
// the representation already used elsewhere in this package for CertificateRef collections
// (e.g. SignatureCertificateSource.SigningCertificateRefs() in signature_certificate_source.go).
package spi

// ProofOfPossessionCertificateSource contains a list of keys, to help the recipient
// cryptographically confirm proof of possession of the key by the token's presenter. Proof of
// possession of a key is also sometimes described as the presenter being a holder-of-key.
type ProofOfPossessionCertificateSource interface {
	CertificateSource

	// AllCertificateRefs returns all certificate references representing the holder's public
	// key certificate. Port of getAllCertificateRefs().
	AllCertificateRefs() []*CertificateRef
}
