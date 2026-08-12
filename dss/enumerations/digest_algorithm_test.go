package enumerations

import "testing"

type digestAlgorithmCase struct {
	v             DigestAlgorithm
	name          string
	javaName      string
	oid           string
	xmlID         string
	jadesID       string
	httpHeaderID  string
	sdJwtID       string
	srIntegrityID string
	coseID        *int64
	msoID         string
	saltLength    int
}

func digestAlgorithmCases() []digestAlgorithmCase {
	return []digestAlgorithmCase{
		{DigestAlgorithm_SHA1, "SHA1", "SHA-1", "1.3.14.3.2.26", "http://www.w3.org/2000/09/xmldsig#sha1", "", "SHA", "", "", int64p(-14), "", 20},
		{DigestAlgorithm_SHA224, "SHA224", "SHA-224", "2.16.840.1.101.3.4.2.4", "http://www.w3.org/2001/04/xmldsig-more#sha224", "S224", "", "", "", nil, "", 28},
		{DigestAlgorithm_SHA256, "SHA256", "SHA-256", "2.16.840.1.101.3.4.2.1", "http://www.w3.org/2001/04/xmlenc#sha256", "S256", "SHA-256", "sha-256", "sha256", int64p(-16), "SHA-256", 32},
		{DigestAlgorithm_SHA384, "SHA384", "SHA-384", "2.16.840.1.101.3.4.2.2", "http://www.w3.org/2001/04/xmldsig-more#sha384", "S384", "", "sha-384", "sha384", int64p(-43), "SHA-384", 48},
		{DigestAlgorithm_SHA512, "SHA512", "SHA-512", "2.16.840.1.101.3.4.2.3", "http://www.w3.org/2001/04/xmlenc#sha512", "S512", "SHA-512", "sha-512", "sha512", int64p(-44), "SHA-512", 64},
		{DigestAlgorithm_SHA3_224, "SHA3-224", "SHA3-224", "2.16.840.1.101.3.4.2.7", "http://www.w3.org/2007/05/xmldsig-more#sha3-224", "", "", "sha3-224", "", nil, "", 28},
		{DigestAlgorithm_SHA3_256, "SHA3-256", "SHA3-256", "2.16.840.1.101.3.4.2.8", "http://www.w3.org/2007/05/xmldsig-more#sha3-256", "S3-256", "", "sha3-256", "", nil, "", 32},
		{DigestAlgorithm_SHA3_384, "SHA3-384", "SHA3-384", "2.16.840.1.101.3.4.2.9", "http://www.w3.org/2007/05/xmldsig-more#sha3-384", "S3-384", "", "sha3-384", "", nil, "", 48},
		{DigestAlgorithm_SHA3_512, "SHA3-512", "SHA3-512", "2.16.840.1.101.3.4.2.10", "http://www.w3.org/2007/05/xmldsig-more#sha3-512", "S3-512", "", "sha3-512", "", nil, "", 64},
		{DigestAlgorithm_SHAKE128, "SHAKE-128", "SHAKE-128", "2.16.840.1.101.3.4.2.11", "", "", "", "", "", int64p(-18), "", 0},
		{DigestAlgorithm_SHAKE256, "SHAKE-256", "SHAKE-256", "2.16.840.1.101.3.4.2.12", "", "", "", "", "", nil, "", 0},
		{DigestAlgorithm_SHAKE256_512, "SHAKE256-512", "SHAKE256-512", "2.16.840.1.101.3.4.2.18", "", "", "", "", "", int64p(-45), "", 0},
		{DigestAlgorithm_RIPEMD160, "RIPEMD160", "RIPEMD160", "1.3.36.3.2.1", "http://www.w3.org/2001/04/xmlenc#ripemd160", "", "", "", "", nil, "", 0},
		{DigestAlgorithm_MD2, "MD2", "MD2", "1.2.840.113549.2.2", "http://www.w3.org/2001/04/xmldsig-more#md2", "", "", "", "", nil, "", 0},
		{DigestAlgorithm_MD5, "MD5", "MD5", "1.2.840.113549.2.5", "http://www.w3.org/2001/04/xmldsig-more#md5", "", "MD5", "", "", nil, "", 0},
		{DigestAlgorithm_WHIRLPOOL, "WHIRLPOOL", "WHIRLPOOL", "1.0.10118.3.0.55", "http://www.w3.org/2007/05/xmldsig-more#whirlpool", "", "", "", "", nil, "", 0},
	}
}

func TestDigestAlgorithmFields(t *testing.T) {
	cases := digestAlgorithmCases()
	if len(DigestAlgorithmValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(DigestAlgorithmValues()))
	}
	for _, c := range cases {
		if got := c.v.Name(); got != c.name {
			t.Errorf("%v.Name() = %q, want %q", c.v, got, c.name)
		}
		if got := c.v.JavaName(); got != c.javaName {
			t.Errorf("%v.JavaName() = %q, want %q", c.v, got, c.javaName)
		}
		if got := c.v.OID(); got != c.oid {
			t.Errorf("%v.OID() = %q, want %q", c.v, got, c.oid)
		}
		if got := c.v.URI(); got != c.xmlID {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.xmlID)
		}
		if got := c.v.JAdESID(); got != c.jadesID {
			t.Errorf("%v.JAdESID() = %q, want %q", c.v, got, c.jadesID)
		}
		if got := c.v.HttpHeaderAlgo(); got != c.httpHeaderID {
			t.Errorf("%v.HttpHeaderAlgo() = %q, want %q", c.v, got, c.httpHeaderID)
		}
		if got := c.v.SDJWTID(); got != c.sdJwtID {
			t.Errorf("%v.SDJWTID() = %q, want %q", c.v, got, c.sdJwtID)
		}
		if got := c.v.SubresourceIntegrityID(); got != c.srIntegrityID {
			t.Errorf("%v.SubresourceIntegrityID() = %q, want %q", c.v, got, c.srIntegrityID)
		}
		gotCose := c.v.CoseID()
		if (gotCose == nil) != (c.coseID == nil) || (gotCose != nil && *gotCose != *c.coseID) {
			t.Errorf("%v.CoseID() = %v, want %v", c.v, gotCose, c.coseID)
		}
		if got := c.v.MSOID(); got != c.msoID {
			t.Errorf("%v.MSOID() = %q, want %q", c.v, got, c.msoID)
		}
		if got := c.v.SaltLength(); got != c.saltLength {
			t.Errorf("%v.SaltLength() = %d, want %d", c.v, got, c.saltLength)
		}
	}
}

func TestDigestAlgorithmForName(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		got, err := DigestAlgorithmForName(c.name)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForName(%q) = %v, %v; want %v, nil", c.name, got, err, c.v)
		}
		if !DigestAlgorithmIsSupportedAlgorithm(c.name) {
			t.Errorf("DigestAlgorithmIsSupportedAlgorithm(%q) = false, want true", c.name)
		}
	}
	if _, err := DigestAlgorithmForName("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
	if DigestAlgorithmIsSupportedAlgorithm("NOPE") {
		t.Error("expected NOPE to be unsupported")
	}
	if got := DigestAlgorithmForNameDefault("NOPE", DigestAlgorithm_SHA256); got != DigestAlgorithm_SHA256 {
		t.Errorf("DigestAlgorithmForNameDefault(NOPE, SHA256) = %v, want SHA256", got)
	}
	if got := DigestAlgorithmForNameDefault("SHA512", DigestAlgorithm_SHA256); got != DigestAlgorithm_SHA512 {
		t.Errorf("DigestAlgorithmForNameDefault(SHA512, SHA256) = %v, want SHA512", got)
	}
}

func TestDigestAlgorithmForJavaName(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		got, err := DigestAlgorithmForJavaName(c.javaName)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForJavaName(%q) = %v, %v; want %v, nil", c.javaName, got, err, c.v)
		}
	}
	if _, err := DigestAlgorithmForJavaName("NOPE"); err == nil {
		t.Error("expected error for unknown java name")
	}
}

func TestDigestAlgorithmForOID(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		got, err := DigestAlgorithmForOID(c.oid)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForOID(%q) = %v, %v; want %v, nil", c.oid, got, err, c.v)
		}
	}
	if _, err := DigestAlgorithmForOID("9.9.9"); err == nil {
		t.Error("expected error for unknown oid")
	}
}

func TestDigestAlgorithmForXML(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		if c.xmlID == "" {
			continue
		}
		got, err := DigestAlgorithmForXML(c.xmlID)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForXML(%q) = %v, %v; want %v, nil", c.xmlID, got, err, c.v)
		}
	}
	if _, err := DigestAlgorithmForXML("nope"); err == nil {
		t.Error("expected error for unknown xml uri")
	}
}

func TestDigestAlgorithmForJAdES(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		if c.jadesID == "" {
			continue
		}
		got, err := DigestAlgorithmForJAdES(c.jadesID)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForJAdES(%q) = %v, %v; want %v, nil", c.jadesID, got, err, c.v)
		}
	}
	if _, err := DigestAlgorithmForJAdES("nope"); err == nil {
		t.Error("expected error for unknown jades id")
	}
}

func TestDigestAlgorithmForHttpHeader(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		if c.httpHeaderID == "" {
			continue
		}
		got, err := DigestAlgorithmForHttpHeader(c.httpHeaderID)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForHttpHeader(%q) = %v, %v; want %v, nil", c.httpHeaderID, got, err, c.v)
		}
	}
	if _, err := DigestAlgorithmForHttpHeader("nope"); err == nil {
		t.Error("expected error for unknown http header id")
	}
}

func TestDigestAlgorithmForSdJwtId(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		if c.sdJwtID == "" {
			continue
		}
		got, err := DigestAlgorithmForSdJwtId(c.sdJwtID)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForSdJwtId(%q) = %v, %v; want %v, nil", c.sdJwtID, got, err, c.v)
		}
	}
	if _, err := DigestAlgorithmForSdJwtId("nope"); err == nil {
		t.Error("expected error for unknown sd-jwt id")
	}
}

func TestDigestAlgorithmForSrIntegrityId(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		if c.srIntegrityID == "" {
			continue
		}
		got, err := DigestAlgorithmForSrIntegrityId(c.srIntegrityID)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForSrIntegrityId(%q) = %v, %v; want %v, nil", c.srIntegrityID, got, err, c.v)
		}
	}
	if _, err := DigestAlgorithmForSrIntegrityId("nope"); err == nil {
		t.Error("expected error for unknown sr-integrity id")
	}
}

func TestDigestAlgorithmForCOSE(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		if c.coseID == nil {
			continue
		}
		got, err := DigestAlgorithmForCOSE(*c.coseID)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForCOSE(%d) = %v, %v; want %v, nil", *c.coseID, got, err, c.v)
		}
	}
	if _, err := DigestAlgorithmForCOSE(-999); err == nil {
		t.Error("expected error for unknown cose id")
	}
}

func TestDigestAlgorithmForMSO(t *testing.T) {
	for _, c := range digestAlgorithmCases() {
		if c.msoID == "" {
			continue
		}
		got, err := DigestAlgorithmForMSO(c.msoID)
		if err != nil || got != c.v {
			t.Errorf("DigestAlgorithmForMSO(%q) = %v, %v; want %v, nil", c.msoID, got, err, c.v)
		}
	}
	if _, err := DigestAlgorithmForMSO("nope"); err == nil {
		t.Error("expected error for unknown mso id")
	}
}
