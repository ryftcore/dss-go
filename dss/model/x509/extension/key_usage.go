// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/KeyUsage.java (DSS 6.5.RC1).
package extension

import "github.com/utain/esig/dss/enumerations"

// KeyUsage is RFC 5280 4.2.1.3. Key Usage.
//
// The key usage extension defines the purpose (e.g., encipherment, signature, certificate
// signing) of the key contained in the certificate. The usage restriction might be
// employed when a key that could be used for more than one operation is to be restricted.
// For example, when an RSA key should be used only to verify signatures on objects other
// than public key certificates and CRLs, the digitalSignature and/or nonRepudiation bits
// would be asserted. Likewise, when an RSA key should be used only for key management, the
// keyEncipherment bit would be asserted.
type KeyUsage struct {
	CertificateExtension

	// keyUsageBits lists defined key usage bits.
	keyUsageBits []enumerations.KeyUsageBit
}

// NewKeyUsage builds a KeyUsage extension.
func NewKeyUsage() *KeyUsage {
	return &KeyUsage{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnum_KEY_USAGE.OID()),
	}
}

// KeyUsageBits returns the key usage bits.
func (k *KeyUsage) KeyUsageBits() []enumerations.KeyUsageBit {
	return k.keyUsageBits
}

// SetKeyUsageBits sets the key usage bits.
func (k *KeyUsage) SetKeyUsageBits(keyUsageBits []enumerations.KeyUsageBit) {
	k.keyUsageBits = keyUsageBits
}
