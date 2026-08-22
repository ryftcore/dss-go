// Ported from dss-model/.../claim/ClaimDeviceKey.java (DSS 6.5.RC1).
package claim

import (
	"crypto"

	"github.com/ryftcore/dss-go/dss/model"
)

// DeviceKey represents a device key used for creating a key-binding
// signature.
type DeviceKey interface {
	Claim

	// PublicKey gets the public key. Ports ClaimDeviceKey#getPublicKey.
	PublicKey() crypto.PublicKey

	// Certificates gets a list of provided certificates. Ports
	// ClaimDeviceKey#getCertificates.
	Certificates() []*model.CertificateToken

	// CertificateDigests gets a list of certificate digests. Ports
	// ClaimDeviceKey#getCertificateDigests.
	CertificateDigests() []model.Digest

	// CertificateKeyIdentifiers gets a list of certificate key
	// identifiers (KID). Ports
	// ClaimDeviceKey#getCertificateKeyIdentifiers.
	CertificateKeyIdentifiers() []string

	// CertificateUrls gets a list of certificate access URLs. Ports
	// ClaimDeviceKey#getCertificateUrls.
	CertificateUrls() []string

	// AuthorizedNamespaces gets a list of namespaces the key is
	// authorized to sign. Ports ClaimDeviceKey#getAuthorizedNamespaces.
	AuthorizedNamespaces() []string

	// AuthorizedDataElements gets a map of namespaces and applicable data
	// element lists the key is authorized to sign. Ports
	// ClaimDeviceKey#getAuthorizedDataElements.
	AuthorizedDataElements() map[string][]string
}
