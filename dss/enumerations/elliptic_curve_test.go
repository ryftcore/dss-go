package enumerations

import "testing"

func TestEllipticCurveLabel(t *testing.T) {
	tests := []struct {
		v    EllipticCurve
		want string
	}{
		{EllipticCurveP256, "P-256"},
		{EllipticCurveP384, "P-384"},
		{EllipticCurveP521, "P-521"},
		{EllipticCurveX25519, "X25519"},
		{EllipticCurveX448, "X448"},
		{EllipticCurveED25519, "Ed25519"},
		{EllipticCurveED448, "Ed448"},
		{EllipticCurveSECP256K1, "secp256k1"},
		{EllipticCurveBrainpoolP256R1, "brainpoolP256r1"},
		{EllipticCurveBrainpoolP320R1, "brainpoolP320r1"},
		{EllipticCurveBrainpoolP384R1, "brainpoolP384r1"},
		{EllipticCurveBrainpoolP512R1, "brainpoolP512r1"},
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
		{EllipticCurveP256, 32},
		{EllipticCurveP384, 48},
		{EllipticCurveP521, 66},
		{EllipticCurveX25519, 32},
		{EllipticCurveX448, 56},
		{EllipticCurveED25519, 32},
		{EllipticCurveED448, 57},
		{EllipticCurveSECP256K1, 32},
		{EllipticCurveBrainpoolP256R1, 32},
		{EllipticCurveBrainpoolP320R1, 40},
		{EllipticCurveBrainpoolP384R1, 48},
		{EllipticCurveBrainpoolP512R1, 64},
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
		{EllipticCurveP256, "1.2.840.10045.3.1.7"},
		{EllipticCurveP384, "1.3.132.0.34"},
		{EllipticCurveP521, "1.3.132.0.35"},
		{EllipticCurveX25519, "1.3.101.110"},
		{EllipticCurveX448, "1.3.101.111"},
		{EllipticCurveED25519, "1.3.101.112"},
		{EllipticCurveED448, "1.3.101.113"},
		{EllipticCurveSECP256K1, "1.2.840.10045.3.1.7"},
		{EllipticCurveBrainpoolP256R1, "1.3.36.3.3.2.8.1.1.7"},
		{EllipticCurveBrainpoolP320R1, "1.3.36.3.3.2.8.1.1.9"},
		{EllipticCurveBrainpoolP384R1, "1.3.36.3.3.2.8.1.1.11"},
		{EllipticCurveBrainpoolP512R1, "1.3.36.3.3.2.8.1.1.13"},
	}
	for _, tt := range tests {
		if got := tt.v.OID(); got != tt.oid {
			t.Errorf("%v.OID() = %q, want %q", tt.v, got, tt.oid)
		}
	}
	// P_256 and SECP_256K1 share an OID upstream; the last constant in
	// declaration order for that OID (SECP_256K1) wins the reverse lookup.
	if got := EllipticCurveForOID("1.2.840.10045.3.1.7"); got != EllipticCurveSECP256K1 {
		t.Errorf("EllipticCurveForOID(shared) = %q, want %q", got, EllipticCurveSECP256K1)
	}
	if got := EllipticCurveForOID("1.3.132.0.34"); got != EllipticCurveP384 {
		t.Errorf("EllipticCurveForOID(P384) = %q, want %q", got, EllipticCurveP384)
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
		{EllipticCurveP256, 1},
		{EllipticCurveP384, 2},
		{EllipticCurveP521, 3},
		{EllipticCurveX25519, 4},
		{EllipticCurveX448, 5},
		{EllipticCurveED25519, 6},
		{EllipticCurveED448, 7},
		{EllipticCurveSECP256K1, 8},
		{EllipticCurveBrainpoolP256R1, 256},
		{EllipticCurveBrainpoolP320R1, 257},
		{EllipticCurveBrainpoolP384R1, 258},
		{EllipticCurveBrainpoolP512R1, 259},
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
		EllipticCurveP256, EllipticCurveP384, EllipticCurveP521,
		EllipticCurveSECP256K1,
		EllipticCurveBrainpoolP256R1, EllipticCurveBrainpoolP320R1,
		EllipticCurveBrainpoolP384R1, EllipticCurveBrainpoolP512R1,
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
		EllipticCurveX25519, EllipticCurveX448, EllipticCurveED25519, EllipticCurveED448,
	}
	for _, v := range withoutParams {
		if _, ok := v.Parameter(); ok {
			t.Errorf("%v.Parameter() ok = true, want false", v)
		}
	}
}

func TestEllipticCurveP256PrimeMatchesKnownConstant(t *testing.T) {
	p, ok := EllipticCurveP256.Parameter()
	if !ok {
		t.Fatal("P_256 parameter not registered")
	}
	want := "115792089210356248762697446949407573530086143415290314195533631308867097853951"
	if p.P.String() != want {
		t.Errorf("P_256 prime = %s, want %s", p.P.String(), want)
	}
}
