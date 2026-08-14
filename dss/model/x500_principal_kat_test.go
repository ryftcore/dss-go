package model

import (
	"encoding/base64"
	"testing"
)

// The vectors below were produced by running javax.security.auth.x500.X500Principal
// on OpenJDK 21 (the JDK dss-model targets) and base64-encoding the UTF-8 bytes of
// getName(RFC2253) and getName(CANONICAL). They pin the Go reimplementation of the
// JDK name-form algorithms to the behaviour dss-model's X500PrincipalHelper relies on.
var x500PrincipalKnownAnswers = []struct {
	name      string
	derB64    string
	rfc2253   string
	canonical string
}{
	{
		name:      "simple",
		derB64:    "MC8xCzAJBgNVBAYTAkJFMQ0wCwYDVQQKEwRBY21lMREwDwYDVQQDEwhKb2huIERvZQ==",
		rfc2253:   "Q049Sm9obiBEb2UsTz1BY21lLEM9QkU=",
		canonical: "Y249am9obiBkb2Usbz1hY21lLGM9YmU=",
	},
	{
		// Multi-valued RDN, including an OID with no RFC 2253 keyword whose
		// UTF8String value is hex-encoded even in the non-canonical form.
		name:      "multi valued RDN",
		derB64:    "MC0xCzAJBgNVBAYTAkZSMR4wCAYDVQQDEwFBMAgGA1UEChMBQjAIBgNVBGMMAVg=",
		rfc2253:   "Q049QStPPUIrMi41LjQuOTk9IzBjMDE1OCxDPUZS",
		canonical: "Y249YStvPWIrMi41LjQuOTk9IzBjMDE1OCxjPWZy",
	},
	{
		// Canonical form collapses internal whitespace runs and NFKD-decomposes
		// the umlauts, so "Ünïcode" becomes "u<COMBINING DIAERESIS>ni<...>code".
		name:      "internal whitespace and non ASCII",
		derB64:    "MEMxCzAJBgNVBAYTAkRFMRcwFQYDVQQKDA7DnG7Dr2NvZGUgR21iSDEbMBkGA1UEAxMSTWl4ZWQgICBDQVNFICBOYW1l",
		rfc2253:   "Q049TWl4ZWQgICBDQVNFICBOYW1lLE89w5xuw69jb2RlIEdtYkgsQz1ERQ==",
		canonical: "Y249bWl4ZWQgY2FzZSBuYW1lLG89dcyIbmnMiGNvZGUgZ21iaCxjPWRl",
	},
	{
		// The canonical escapee set is narrower than the RFC 2253 one: '=' and a
		// non-leading '#' are escaped by getName(RFC2253) but not by CANONICAL.
		name:      "escaped characters",
		derB64:    "MCYxCjAIBgNVBAoTAVgxGDAWBgNVBAMMD2EsYitjPWQ8ZT5mI2c7aA==",
		rfc2253:   "Q049YVwsYlwrY1w9ZFw8ZVw+ZlwjZ1w7aCxPPVg=",
		canonical: "Y249YVwsYlwrYz1kXDxlXD5mI2dcO2gsbz14",
	},
	{
		// emailAddress has no RFC 2253 compliant keyword: it is rendered as a
		// dotted-decimal OID with the hex of the whole BER value.
		name:      "email address OID",
		derB64:    "MCQxCjAIBgNVBAMTAXgxFjAUBgkqhkiG9w0BCQEWB2FAYi5jb20=",
		rfc2253:   "MS4yLjg0MC4xMTM1NDkuMS45LjE9IzE2MDc2MTQwNjIyZTYzNmY2ZCxDTj14",
		canonical: "MS4yLjg0MC4xMTM1NDkuMS45LjE9IzE2MDc2MTQwNjIyZTYzNmY2ZCxjbj14",
	},
	{
		// dnQualifier, surname and title likewise have no compliant keyword.
		name:      "non compliant keywords",
		derB64:    "MDUxCjAIBgNVBAMTAXkxDTALBgNVBAwTBEJvc3MxDDAKBgNVBAQTA1N1cjEKMAgGA1UELhMBQQ==",
		rfc2253:   "Mi41LjQuNDY9IzEzMDE0MSwyLjUuNC40PSMxMzAzNTM3NTcyLDIuNS40LjEyPSMxMzA0NDI2ZjczNzMsQ049eQ==",
		canonical: "Mi41LjQuNDY9IzEzMDE0MSwyLjUuNC40PSMxMzAzNTM3NTcyLDIuNS40LjEyPSMxMzA0NDI2ZjczNzMsY249eQ==",
	},
	{
		name:      "UTF8String values",
		derB64:    "MCUxDzANBgNVBAoMBsOEcmdlcjESMBAGA1UEAwwJw5xuw69jb2Rl",
		rfc2253:   "Q049w5xuw69jb2RlLE89w4RyZ2Vy",
		canonical: "Y249dcyIbmnMiGNvZGUsbz1hzIhyZ2Vy",
	},
	{
		// domainComponent values are IA5Strings, which the canonical form does
		// not treat as a string type: they are hex-encoded, unlike the
		// PrintableString userid.
		name:      "domain component and user id",
		derB64:    "MFMxDTALBgNVBAMTBEpvaG4xFDASBgoJkiaJk/IsZAEBEwRqZG9lMRcwFQYKCZImiZPyLGQBGRYHZXhhbXBsZTETMBEGCgmSJomT8ixkARkWA2NvbQ==",
		rfc2253:   "REM9Y29tLERDPWV4YW1wbGUsVUlEPWpkb2UsQ049Sm9obg==",
		canonical: "ZGM9IzE2MDM2MzZmNmQsZGM9IzE2MDc2NTc4NjE2ZDcwNmM2NSx1aWQ9amRvZSxjbj1qb2hu",
	},
	{
		// 2.5.4.8 resolves to "ST" (not "S") and 2.5.4.9 to "STREET".
		name:      "state and street keywords",
		derB64:    "MEExCjAIBgNVBAMTAXoxDzANBgNVBAcTBk11bmljaDEQMA4GA1UECRMHTWFpbiBSZDEQMA4GA1UECBMHQmF2YXJpYQ==",
		rfc2253:   "U1Q9QmF2YXJpYSxTVFJFRVQ9TWFpbiBSZCxMPU11bmljaCxDTj16",
		canonical: "c3Q9YmF2YXJpYSxzdHJlZXQ9bWFpbiByZCxsPW11bmljaCxjbj16",
	},
	{
		// A value whose first character is '#' is escaped in both forms.
		name:      "leading hash",
		derB64:    "MB8xCjAIBgNVBAoTAXExETAPBgNVBAMMCCNub3RhSGV4",
		rfc2253:   "Q049XCNub3RhSGV4LE89cQ==",
		canonical: "Y249XCNub3RhaGV4LG89cQ==",
	},
	{
		name:      "empty value",
		derB64:    "MBcxCjAIBgNVBAoTAWUxCTAHBgNVBAMTAA==",
		rfc2253:   "Q049LE89ZQ==",
		canonical: "Y249LG89ZQ==",
	},
}

func x500PrincipalDecodeB64(t *testing.T, s string) []byte {
	t.Helper()
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Fatalf("cannot decode test vector: %v", err)
	}
	return data
}

func TestX500PrincipalNameFormsMatchJDK(t *testing.T) {
	for _, tc := range x500PrincipalKnownAnswers {
		t.Run(tc.name, func(t *testing.T) {
			der := x500PrincipalDecodeB64(t, tc.derB64)
			principal, err := NewX500Principal(der)
			if err != nil {
				t.Fatalf("NewX500Principal: %v", err)
			}
			wantRFC2253 := string(x500PrincipalDecodeB64(t, tc.rfc2253))
			if got := principal.RFC2253Name(); got != wantRFC2253 {
				t.Errorf("RFC2253Name() = %q, want %q", got, wantRFC2253)
			}
			wantCanonical := string(x500PrincipalDecodeB64(t, tc.canonical))
			if got := principal.Canonical(); got != wantCanonical {
				t.Errorf("Canonical() = %q, want %q", got, wantCanonical)
			}
			// The DER must be handed back byte-for-byte, never re-encoded.
			if got := principal.Encoded(); string(got) != string(der) {
				t.Errorf("Encoded() returned re-encoded DER")
			}
		})
	}
}

// TestX500PrincipalHelperFormsMatchJDK checks the accessors dss-model actually uses.
func TestX500PrincipalHelperFormsMatchJDK(t *testing.T) {
	tc := x500PrincipalKnownAnswers[0]
	der := x500PrincipalDecodeB64(t, tc.derB64)
	principal, err := NewX500Principal(der)
	if err != nil {
		t.Fatalf("NewX500Principal: %v", err)
	}
	helper := NewX500PrincipalHelper(principal)
	if got, want := helper.RFC2253(), string(x500PrincipalDecodeB64(t, tc.rfc2253)); got != want {
		t.Errorf("RFC2253() = %q, want %q", got, want)
	}
	if got, want := helper.Canonical(), string(x500PrincipalDecodeB64(t, tc.canonical)); got != want {
		t.Errorf("Canonical() = %q, want %q", got, want)
	}
	if string(helper.Encoded()) != string(der) {
		t.Errorf("Encoded() must return the original DER")
	}
	// getPrettyPrintRFC2253 substitutes the X520Attributes OID descriptions.
	pretty, err := helper.PrettyPrintRFC2253()
	if err != nil {
		t.Fatalf("PrettyPrintRFC2253: %v", err)
	}
	if pretty == "" {
		t.Errorf("PrettyPrintRFC2253() must not be empty")
	}
}
