// Package token ports dss-token (eu.europa.esig.dss.token), the signing-key
// access layer: connections to PKCS#11 devices, PKCS#12/JKS keystores, and
// (where the platform supports it) the OS-native MSCAPI and Apple keychain
// stores, all exposed through one common interface so signing services do
// not need to know which key storage backend produced a private key.
//
// The main entry types are SignatureTokenConnection (the common connection
// interface), DSSPrivateKeyEntry (a selected signing key plus its
// certificate chain), and the concrete connections: Pkcs11SignatureToken,
// Pkcs12SignatureToken, JKSSignatureToken, KeyStoreSignatureTokenConnection,
// and the platform-specific MSCAPISignatureToken/AppleSignatureToken.
package token
