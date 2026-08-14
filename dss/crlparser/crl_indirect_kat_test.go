package crlparser

import (
	"math/big"
	"testing"
)

// Known answers captured from upstream DSS 6.5.RC1 (eu.europa.esig.dss.crl.CRLUtils
// #getRevocationInfo, i.e. BouncyCastle's X509CRLObject#getRevokedCertificate) running on the
// two fixtures below.
//
//	testdata/ca_indirect.crl.der         issuingDistributionPoint with indirectCRL TRUE
//	testdata/ca_direct_certissuer.crl.der  the same entries, indirectCRL FALSE
//
// Both carry four entries: 101 declares certificateIssuer "CN=Issuer A,C=LU", 102 declares
// none, 103 declares "CN=Issuer B,C=LU" and 104 declares none. RFC 5280 section 5.3.3 makes an
// entry without a certificateIssuer inherit the preceding one, and BouncyCastle applies that
// only when the CRL itself is indirect - both properties are pinned here.
func TestCRLUtilsRevocationInfo_CertificateIssuerCarryForward(t *testing.T) {
	issuer := crlparserLoadCertificateToken(t, "ca_carry.crt")

	testCases := []struct {
		name     string
		fixture  string
		expected map[int64]string
	}{
		{
			name:    "indirect CRL inherits the preceding certificateIssuer",
			fixture: "ca_indirect.crl.der",
			expected: map[int64]string{
				101: "C=LU,CN=Issuer A",
				102: "C=LU,CN=Issuer A",
				103: "C=LU,CN=Issuer B",
				104: "C=LU,CN=Issuer B",
			},
		},
		{
			name:    "direct CRL reports no certificateIssuer at all",
			fixture: "ca_direct_certissuer.crl.der",
			expected: map[int64]string{
				101: "",
				102: "",
				103: "",
				104: "",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			crlBinary, err := CRLUtilsBuildCRLBinary(crlparserReadFile(t, testCase.fixture))
			if err != nil {
				t.Fatalf("building the CRL binary: %v", err)
			}
			crlValidity, err := CRLUtilsBuildCRLValidity(crlBinary, issuer)
			if err != nil {
				t.Fatalf("building the CRL validity: %v", err)
			}

			for serialNumber, expected := range testCase.expected {
				entry := CRLUtilsRevocationInfo(crlValidity, big.NewInt(serialNumber))
				if entry == nil {
					t.Fatalf("serial %d: expected a revocation entry", serialNumber)
				}
				actual := ""
				if certificateIssuer := entry.CertificateIssuer(); certificateIssuer != nil {
					actual = certificateIssuer.RFC2253Name()
				}
				if actual != expected {
					t.Errorf("serial %d: CertificateIssuer() = %q, want %q", serialNumber, actual, expected)
				}
			}

			if entry := CRLUtilsRevocationInfo(crlValidity, big.NewInt(999)); entry != nil {
				t.Errorf("serial 999: expected no revocation entry, got %v", entry.SerialNumber())
			}
		})
	}
}
