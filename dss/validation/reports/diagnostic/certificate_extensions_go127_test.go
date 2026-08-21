//go:build go1.27

package diagnostic

// Go 1.27 tightened crypto/x509's SAN parsing: synthetic_00.der (SAN
// iPAddress containing an IPv4-mapped IPv6 address) is now rejected at
// parse time, where earlier toolchains — and BouncyCastle, and therefore
// Java DSS — accept it. Register the fixture in both skip maps only on
// 1.27+ so older toolchains keep their full KAT/smoke coverage. The spi
// package gates its own copy of this set the same way (see
// spi/certificate_extensions_kat_unparseable_go127_test.go).
func init() {
	certificateExtensionsSmokeUnparseable["synthetic_00.der"] = true
	katKnownDeviation["synthetic_00"] = "crypto/x509 (Go 1.27+) rejects the DER: SAN iPAddress contains IPv4-mapped IPv6 address"
}
