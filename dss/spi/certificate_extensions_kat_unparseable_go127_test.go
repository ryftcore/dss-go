//go:build go1.27

package spi

// The crypto/x509 parser in Go 1.27+ rejects one more corpus fixture than
// earlier toolchains (BouncyCastle, and therefore Java DSS, accepts all
// three):
//
//	cert_16.der:     "x509: malformed certificate"
//	cert_19.der:     "x509: RSA key missing NULL parameters"
//	synthetic_00.der: "x509: SAN iPAddress contains IPv4-mapped IPv6 address"
//	                  (new strictness in the Go 1.27 SAN parser)
var certificateExtensionsKATUnparseable = map[string]bool{
	"cert_16.der":      true,
	"cert_19.der":      true,
	"synthetic_00.der": true,
}
