// Ported from dss-model/.../claim/ClaimRevocationList.java (DSS 6.5.RC1).
package claim

// RevocationList represents a generic EAA revocation list.
type RevocationList interface {
	Claim

	// Uri gets the EAA's Status URI value, when present. Ports
	// ClaimRevocationList#getUri.
	Uri() *String

	// Certificate gets a certificate containing the public key that
	// signed or sealed the top-level certificate in the x5chain element
	// in the MSO revocation list structure. Ports
	// ClaimRevocationList#getCertificate.
	Certificate() *ByteString
}
