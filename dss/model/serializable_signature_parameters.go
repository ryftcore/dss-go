// Ported from dss-model/.../SerializableSignatureParameters.java (DSS 6.5.RC1).
package model

import "github.com/utain/esig/dss/enumerations"

// SerializableSignatureParameters defines signature parameters.
//
// CertificateToken (ported from model.x509, flattened into this package)
// is outside this manifest; assumed to already exist in this package.
type SerializableSignatureParameters interface {
	// SigningCertificate gets the signing certificate. Ports
	// #getSigningCertificate.
	SigningCertificate() *CertificateToken

	// GenerateTBSWithoutCertificate indicates if it is possible to
	// generate ToBeSigned data without the signing certificate. The
	// default value is false. Ports
	// #isGenerateTBSWithoutCertificate.
	GenerateTBSWithoutCertificate() bool

	// CheckCertificateRevocation indicates whether a revocation check
	// shall be performed before -LT level incorporation (i.e. on signing
	// or T-level creation) for a signing certificate and a respectful
	// certificate chain. When false, the revocation check is not
	// performed. When true, a real-time revocation is being requested
	// from external sources (shall be defined in CertificateVerifier) and
	// processed according to alerts set within that CertificateVerifier.
	//
	// Default value: false (no revocation check is performed on
	// signature creation or T-level extension). Ports
	// #isCheckCertificateRevocation.
	CheckCertificateRevocation() bool

	// BLevel gets Baseline B parameters (signed properties). Ports
	// #bLevel.
	BLevel() *BLevelParameters

	// DigestAlgorithm gets the digest algorithm. Ports
	// #getDigestAlgorithm.
	DigestAlgorithm() enumerations.DigestAlgorithm

	// EncryptionAlgorithm gets the encryption algorithm. Ports
	// #getEncryptionAlgorithm.
	EncryptionAlgorithm() enumerations.EncryptionAlgorithm

	// SignatureAlgorithm gets the signature algorithm. Ports
	// #getSignatureAlgorithm.
	SignatureAlgorithm() enumerations.SignatureAlgorithm
}
