package enumerations

import "testing"

func TestEllipticCurveLabel(t *testing.T) {
	tests := []struct {
		v    EllipticCurve
		want string
	}{
		{EllipticCurve_P_256, "P-256"},
		{EllipticCurve_P_384, "P-384"},
		{EllipticCurve_P_521, "P-521"},
		{EllipticCurve_X25519, "X25519"},
		{EllipticCurve_X448, "X448"},
		{EllipticCurve_ED25519, "Ed25519"},
		{EllipticCurve_ED448, "Ed448"},
		{EllipticCurve_SECP_256K1, "secp256k1"},
		{EllipticCurve_BRAINPOOL_P256_R1, "brainpoolP256r1"},
		{EllipticCurve_BRAINPOOL_P320_R1, "brainpoolP320r1"},
		{EllipticCurve_BRAINPOOL_P384_R1, "brainpoolP384r1"},
		{EllipticCurve_BRAINPOOL_P512_R1, "brainpoolP512r1"},
	}
	for _, tt := range tests {
		if got := tt.v.Label(); got != tt.want {
			t.Errorf("%v.Label() = %q, want %q", tt.v, got, tt.want)
		}
		if got := EllipticCurveForLabel(tt.want); got != tt.v {
			t.Errorf("EllipticCurveForLabel(%q) = %q, want %q", tt.want, got, tt.v)
		}
	}
	if got := EllipticCurveForLabel(""); got != "" {
		t.Errorf("EllipticCurveForLabel(\"\") = %q, want empty", got)
	}
	if got := EllipticCurveForLabel("bogus"); got != "" {
		t.Errorf("EllipticCurveForLabel(bogus) = %q, want empty", got)
	}
}

func TestEllipticCurveSize(t *testing.T) {
	tests := []struct {
		v    EllipticCurve
		want int
	}{
		{EllipticCurve_P_256, 32},
		{EllipticCurve_P_384, 48},
		{EllipticCurve_P_521, 66},
		{EllipticCurve_X25519, 32},
		{EllipticCurve_X448, 56},
		{EllipticCurve_ED25519, 32},
		{EllipticCurve_ED448, 57},
		{EllipticCurve_SECP_256K1, 32},
		{EllipticCurve_BRAINPOOL_P256_R1, 32},
		{EllipticCurve_BRAINPOOL_P320_R1, 40},
		{EllipticCurve_BRAINPOOL_P384_R1, 48},
		{EllipticCurve_BRAINPOOL_P512_R1, 64},
	}
	for _, tt := range tests {
		if got := tt.v.Size(); got != tt.want {
			t.Errorf("%v.Size() = %d, want %d", tt.v, got, tt.want)
		}
	}
}

func TestEllipticCurveOID(t *testing.T) {
	tests := []struct {
		v   EllipticCurve
		oid string
	}{
		{EllipticCurve_P_256, "1.2.840.10045.3.1.7"},
		{EllipticCurve_P_384, "1.3.132.0.34"},
		{EllipticCurve_P_521, "1.3.132.0.35"},
		{EllipticCurve_X25519, "1.3.101.110"},
		{EllipticCurve_X448, "1.3.101.111"},
		{EllipticCurve_ED25519, "1.3.101.112"},
		{EllipticCurve_ED448, "1.3.101.113"},
		{EllipticCurve_SECP_256K1, "1.2.840.10045.3.1.7"},
		{EllipticCurve_BRAINPOOL_P256_R1, "1.3.36.3.3.2.8.1.1.7"},
		{EllipticCurve_BRAINPOOL_P320_R1, "1.3.36.3.3.2.8.1.1.9"},
		{EllipticCurve_BRAINPOOL_P384_R1, "1.3.36.3.3.2.8.1.1.11"},
		{EllipticCurve_BRAINPOOL_P512_R1, "1.3.36.3.3.2.8.1.1.13"},
	}
	for _, tt := range tests {
		if got := tt.v.OID(); got != tt.oid {
			t.Errorf("%v.OID() = %q, want %q", tt.v, got, tt.oid)
		}
	}
	// P_256 and SECP_256K1 share an OID upstream; the last constant in
	// declaration order for that OID (SECP_256K1) wins the reverse lookup.
	if got := EllipticCurveForOID("1.2.840.10045.3.1.7"); got != EllipticCurve_SECP_256K1 {
		t.Errorf("EllipticCurveForOID(shared) = %q, want %q", got, EllipticCurve_SECP_256K1)
	}
	if got := EllipticCurveForOID("1.3.132.0.34"); got != EllipticCurve_P_384 {
		t.Errorf("EllipticCurveForOID(P384) = %q, want %q", got, EllipticCurve_P_384)
	}
	if got := EllipticCurveForOID("bogus"); got != "" {
		t.Errorf("EllipticCurveForOID(bogus) = %q, want empty", got)
	}
}

func TestEllipticCurveCOSEValue(t *testing.T) {
	tests := []struct {
		v    EllipticCurve
		want int64
	}{
		{EllipticCurve_P_256, 1},
		{EllipticCurve_P_384, 2},
		{EllipticCurve_P_521, 3},
		{EllipticCurve_X25519, 4},
		{EllipticCurve_X448, 5},
		{EllipticCurve_ED25519, 6},
		{EllipticCurve_ED448, 7},
		{EllipticCurve_SECP_256K1, 8},
		{EllipticCurve_BRAINPOOL_P256_R1, 256},
		{EllipticCurve_BRAINPOOL_P320_R1, 257},
		{EllipticCurve_BRAINPOOL_P384_R1, 258},
		{EllipticCurve_BRAINPOOL_P512_R1, 259},
	}
	for _, tt := range tests {
		got, ok := tt.v.COSEValue()
		if !ok {
			t.Errorf("%v.COSEValue() ok = false, want true", tt.v)
		}
		if got != tt.want {
			t.Errorf("%v.COSEValue() = %d, want %d", tt.v, got, tt.want)
		}
		if resolved := EllipticCurveForCOSEValue(tt.want); resolved != tt.v {
			t.Errorf("EllipticCurveForCOSEValue(%d) = %q, want %q", tt.want, resolved, tt.v)
		}
	}
	if got := EllipticCurveForCOSEValue(9999); got != "" {
		t.Errorf("EllipticCurveForCOSEValue(unknown) = %q, want empty", got)
	}
}

func TestEllipticCurveParameterRoundtrip(t *testing.T) {
	withParams := []EllipticCurve{
		EllipticCurve_P_256, EllipticCurve_P_384, EllipticCurve_P_521,
		EllipticCurve_SECP_256K1,
		EllipticCurve_BRAINPOOL_P256_R1, EllipticCurve_BRAINPOOL_P320_R1,
		EllipticCurve_BRAINPOOL_P384_R1, EllipticCurve_BRAINPOOL_P512_R1,
	}
	for _, v := range withParams {
		p, ok := v.Parameter()
		if !ok {
			t.Errorf("%v.Parameter() ok = false, want true", v)
			continue
		}
		if resolved := EllipticCurveForParameter(p); resolved != v {
			t.Errorf("EllipticCurveForParameter(%v) = %q, want %q", v, resolved, v)
		}
	}

	withoutParams := []EllipticCurve{
		EllipticCurve_X25519, EllipticCurve_X448, EllipticCurve_ED25519, EllipticCurve_ED448,
	}
	for _, v := range withoutParams {
		if _, ok := v.Parameter(); ok {
			t.Errorf("%v.Parameter() ok = true, want false", v)
		}
	}
}

func TestEllipticCurveP256PrimeMatchesKnownConstant(t *testing.T) {
	p, ok := EllipticCurve_P_256.Parameter()
	if !ok {
		t.Fatal("P_256 parameter not registered")
	}
	want := "115792089210356248762697446949407573530086143415290314195533631308867097853951"
	if p.P.String() != want {
		t.Errorf("P_256 prime = %s, want %s", p.P.String(), want)
	}
}
