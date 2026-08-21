//go:build !go1.27

package spi

// The crypto/x509 parser in Go toolchains before 1.27 rejects exactly these
// two corpus fixtures (BouncyCastle, and therefore Java DSS, accepts both):
//
//	cert_16.der: "x509: malformed certificate"
//	cert_19.der: "x509: RSA key missing NULL parameters"
var certificateExtensionsKATUnparseable = map[string]bool{
	"cert_16.der": true,
	"cert_19.der": true,
}
