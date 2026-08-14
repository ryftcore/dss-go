package revocation

import "testing"

// crlToken and ocspToken stand in for the concrete tokens dss-spi will add.
type crlToken struct{ id string }

type ocspToken struct{ id string }

func TestRevocationMarkerInterfaces(t *testing.T) {
	// CRL and OCSP both widen Revocation, so a CRL or OCSP value is usable wherever
	// revocation data is expected.
	var crl CRL = crlToken{id: "crl"}
	var ocsp OCSP = ocspToken{id: "ocsp"}

	var revocation Revocation = crl
	if _, ok := revocation.(crlToken); !ok {
		t.Error("a CRL must be usable as Revocation")
	}
	revocation = ocsp
	if _, ok := revocation.(ocspToken); !ok {
		t.Error("an OCSP must be usable as Revocation")
	}
}
