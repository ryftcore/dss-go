// Ported from dss-enumerations/.../EncryptionAlgorithm.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestEncryptionAlgorithm(t *testing.T) {
	cases := []struct {
		v       EncryptionAlgorithm
		name    string
		oid     string
		padding string
	}{
		{EncryptionAlgorithm_RSA, "RSA", "1.2.840.113549.1.1.1", "RSA/ECB/PKCS1Padding"},
		{EncryptionAlgorithm_RSASSA_PSS, "RSASSA-PSS", "1.2.840.113549.1.1.10", "RSA/ECB/OAEPPadding"},
		{EncryptionAlgorithm_DSA, "DSA", "1.2.840.10040.4.1", "DSA"},
		{EncryptionAlgorithm_ECDSA, "ECDSA", "1.2.840.10045.2.1", "ECDSA"},
		{EncryptionAlgorithm_PLAIN_ECDSA, "PLAIN-ECDSA", "0.4.0.127.0.7.1.1.4.1", "PLAIN-ECDSA"},
		{EncryptionAlgorithm_X25519, "X25519", "1.3.101.110", "X25519"},
		{EncryptionAlgorithm_X448, "X448", "1.3.101.111", "X448"},
		{EncryptionAlgorithm_EDDSA, "EdDSA", "", "EdDSA"},
		{EncryptionAlgorithm_HMAC, "HMAC", "", ""},
	}
	if len(EncryptionAlgorithmValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(EncryptionAlgorithmValues()))
	}
	for _, c := range cases {
		if got := c.v.Name(); got != c.name {
			t.Errorf("%v.Name() = %q, want %q", c.v, got, c.name)
		}
		if got := c.v.OID(); got != c.oid {
			t.Errorf("%v.OID() = %q, want %q", c.v, got, c.oid)
		}
		if got := c.v.Padding(); got != c.padding {
			t.Errorf("%v.Padding() = %q, want %q", c.v, got, c.padding)
		}
		gotByName, err := EncryptionAlgorithmForName(c.name)
		if err != nil || gotByName != c.v {
			t.Errorf("EncryptionAlgorithmForName(%q) = %v, %v; want %v, nil", c.name, gotByName, err, c.v)
		}
		gotByJavaName, err := EncryptionAlgorithmForName(string(c.v))
		if err != nil || gotByJavaName != c.v {
			t.Errorf("EncryptionAlgorithmForName(%q) = %v, %v; want %v, nil", c.v, gotByJavaName, err, c.v)
		}
	}

	// forOID only has non-empty OIDs for a subset of algorithms.
	oidCases := []struct {
		oid string
		v   EncryptionAlgorithm
	}{
		{"1.2.840.113549.1.1.1", EncryptionAlgorithm_RSA},
		{"1.2.840.113549.1.1.10", EncryptionAlgorithm_RSASSA_PSS},
		{"1.2.840.10040.4.1", EncryptionAlgorithm_DSA},
		{"1.2.840.10045.2.1", EncryptionAlgorithm_ECDSA},
		{"0.4.0.127.0.7.1.1.4.1", EncryptionAlgorithm_PLAIN_ECDSA},
		{"1.3.101.110", EncryptionAlgorithm_X25519},
		{"1.3.101.111", EncryptionAlgorithm_X448},
	}
	for _, c := range oidCases {
		got, err := EncryptionAlgorithmForOID(c.oid)
		if err != nil || got != c.v {
			t.Errorf("EncryptionAlgorithmForOID(%q) = %v, %v; want %v, nil", c.oid, got, err, c.v)
		}
	}
	if _, err := EncryptionAlgorithmForOID("9.9.9"); err == nil {
		t.Error("expected error for unknown OID")
	}
	// EDDSA and HMAC share the "" OID and are not resolvable by forOID
	// (Java's registry map also keeps only the last "" entry inserted:
	// HMAC, since it is declared after EDDSA).
	if got, err := EncryptionAlgorithmForOID(""); err != nil || got != EncryptionAlgorithm_HMAC {
		t.Errorf("EncryptionAlgorithmForOID(\"\") = %v, %v; want %v, nil", got, err, EncryptionAlgorithm_HMAC)
	}

	if _, err := EncryptionAlgorithmForName("nope"); err == nil {
		t.Error("expected error for unknown name")
	}
	if got := EncryptionAlgorithmForNameDefault("nope", EncryptionAlgorithm_RSA); got != EncryptionAlgorithm_RSA {
		t.Errorf("EncryptionAlgorithmForNameDefault(nope) = %v, want %v", got, EncryptionAlgorithm_RSA)
	}
	if got, err := EncryptionAlgorithmForName("EC"); err != nil || got != EncryptionAlgorithm_ECDSA {
		t.Errorf("EncryptionAlgorithmForName(EC) = %v, %v; want %v, nil", got, err, EncryptionAlgorithm_ECDSA)
	}
	if got, err := EncryptionAlgorithmForName("ECC"); err != nil || got != EncryptionAlgorithm_ECDSA {
		t.Errorf("EncryptionAlgorithmForName(ECC) = %v, %v; want %v, nil", got, err, EncryptionAlgorithm_ECDSA)
	}
	if got, err := EncryptionAlgorithmForName("Ed25519"); err != nil || got != EncryptionAlgorithm_EDDSA {
		t.Errorf("EncryptionAlgorithmForName(Ed25519) = %v, %v; want %v, nil", got, err, EncryptionAlgorithm_EDDSA)
	}
	if got, err := EncryptionAlgorithmForName("Ed448"); err != nil || got != EncryptionAlgorithm_EDDSA {
		t.Errorf("EncryptionAlgorithmForName(Ed448) = %v, %v; want %v, nil", got, err, EncryptionAlgorithm_EDDSA)
	}
}

func TestEncryptionAlgorithm_IsEquivalent(t *testing.T) {
	if !EncryptionAlgorithm_RSA.IsEquivalent(EncryptionAlgorithm_RSASSA_PSS) {
		t.Error("expected RSA equivalent to RSASSA_PSS")
	}
	if !EncryptionAlgorithm_ECDSA.IsEquivalent(EncryptionAlgorithm_PLAIN_ECDSA) {
		t.Error("expected ECDSA equivalent to PLAIN_ECDSA")
	}
	if !EncryptionAlgorithm_X25519.IsEquivalent(EncryptionAlgorithm_EDDSA) {
		t.Error("expected X25519 equivalent to EDDSA")
	}
	if EncryptionAlgorithm_RSA.IsEquivalent(EncryptionAlgorithm_ECDSA) {
		t.Error("expected RSA not equivalent to ECDSA")
	}
	if EncryptionAlgorithm_RSA.IsEquivalent("") {
		t.Error("expected RSA not equivalent to empty/nil")
	}
	if !EncryptionAlgorithm_HMAC.IsEquivalent(EncryptionAlgorithm_HMAC) {
		t.Error("expected HMAC equivalent to itself")
	}
}
