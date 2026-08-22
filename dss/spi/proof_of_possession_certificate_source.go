// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/ProofOfPossessionCertificateSource.java (DSS 6.5.RC1).
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
