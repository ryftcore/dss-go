package model

import "testing"

func TestCertificateTokenIdentifierUsesTheCPrefix(t *testing.T) {
	token := certificateTokenFixture(t, leafCertificateBase64)
	identifier := NewCertificateTokenIdentifier(token)

	const want = "C-2E77AC385070F5D3C9AE8DEB5A3496C3076BFD870342AC8692B02330881EE9B8"
	if got := identifier.AsXmlID(); got != want {
		t.Errorf("AsXmlID() = %q, want %q", got, want)
	}
	if got := identifier.String(); got != "CertificateTokenIdentifier:SHA256:#"+want[2:] {
		t.Errorf("String() = %q", got)
	}
	// The token's own identifier must agree with the one built here.
	if got := token.DSSID().AsXmlID(); got != want {
		t.Errorf("token DSSID() = %q, want %q", got, want)
	}
	if !identifier.Equals(NewCertificateTokenIdentifier(token)) {
		t.Error("two identifiers of the same certificate must be equal")
	}
	if identifier.Equals(NewCertificateTokenIdentifier(certificateTokenFixture(t, rootCertificateBase64))) {
		t.Error("identifiers of different certificates must not be equal")
	}
}
