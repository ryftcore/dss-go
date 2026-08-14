package spi

import (
	"bufio"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// crlRefTestKAT reads a "<key>|<value>" fixture file produced by
// testdata/crlocsp/generator/EsfKatGen.java with BouncyCastle.
func crlRefTestKAT(t *testing.T, path string) map[string]string {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<22)
	for scanner.Scan() {
		line := scanner.Text()
		key, value, found := strings.Cut(line, "|")
		if !found {
			continue
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return values
}

// crlRefTestHex decodes a hex fixture value.
func crlRefTestHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	if err != nil {
		t.Fatalf("decode %q: %v", value, err)
	}
	return decoded
}

// crlRefTestMillis decodes a Java Date#getTime() fixture value.
func crlRefTestMillis(t *testing.T, value string) time.Time {
	t.Helper()
	millis, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		t.Fatalf("parse %q: %v", value, err)
	}
	return time.UnixMilli(millis).UTC()
}

// TestParseCrlValidatedIDFull checks a CrlValidatedID carrying a complete CrlIdentifier
// against the BouncyCastle-produced encoding.
func TestParseCrlValidatedIDFull(t *testing.T) {
	kat := crlRefTestKAT(t, "testdata/crlocsp/kat_esf.txt")

	validatedID, err := ParseCrlValidatedID(crlRefTestHex(t, kat["crl.full.der"]))
	if err != nil {
		t.Fatalf("ParseCrlValidatedID: %v", err)
	}

	if got, want := validatedID.CrlHash.HashAlgorithm.Algorithm.String(), kat["crl.full.hashalg"]; got != want {
		t.Errorf("hash algorithm = %s, want %s", got, want)
	}
	if got, want := hex.EncodeToString(validatedID.CrlHash.HashValue), kat["crl.full.hashvalue"]; got != want {
		t.Errorf("hash value = %s, want %s", got, want)
	}
	if validatedID.CrlIdentifier == nil {
		t.Fatal("CrlIdentifier is nil")
	}
	// The crlissuer Name must be handed back byte-identical: the X500Principal built from
	// it is compared and re-serialised elsewhere.
	if got, want := hex.EncodeToString(validatedID.CrlIdentifier.CrlIssuer), kat["crl.full.issuer"]; got != want {
		t.Errorf("crlissuer = %s, want %s", got, want)
	}
	if got, want := validatedID.CrlIdentifier.CrlIssuedTime, crlRefTestMillis(t, kat["crl.full.issuedTime"]); !got.Equal(want) {
		t.Errorf("crlIssuedTime = %s, want %s", got, want)
	}
	if validatedID.CrlIdentifier.CrlNumber == nil {
		t.Fatal("crlNumber is nil")
	}
	if got, want := validatedID.CrlIdentifier.CrlNumber.String(), kat["crl.full.number"]; got != want {
		t.Errorf("crlNumber = %s, want %s", got, want)
	}
}

// TestParseCrlValidatedIDWithoutNumber checks the optional crlNumber is recognised as absent.
func TestParseCrlValidatedIDWithoutNumber(t *testing.T) {
	kat := crlRefTestKAT(t, "testdata/crlocsp/kat_esf.txt")

	validatedID, err := ParseCrlValidatedID(crlRefTestHex(t, kat["crl.nonumber.der"]))
	if err != nil {
		t.Fatalf("ParseCrlValidatedID: %v", err)
	}
	if validatedID.CrlIdentifier == nil {
		t.Fatal("CrlIdentifier is nil")
	}
	if validatedID.CrlIdentifier.CrlNumber != nil {
		t.Errorf("crlNumber = %v, want nil", validatedID.CrlIdentifier.CrlNumber)
	}
	if got, want := hex.EncodeToString(validatedID.CrlIdentifier.CrlIssuer), kat["crl.full.issuer"]; got != want {
		t.Errorf("crlissuer = %s, want %s", got, want)
	}
}

// TestParseCrlValidatedIDHashOnly checks the sha1Hash alternative of OtherHash, which
// implies id-sha1 without parameters, and the absence of the optional CrlIdentifier.
func TestParseCrlValidatedIDHashOnly(t *testing.T) {
	kat := crlRefTestKAT(t, "testdata/crlocsp/kat_esf.txt")

	validatedID, err := ParseCrlValidatedID(crlRefTestHex(t, kat["crl.hashonly.der"]))
	if err != nil {
		t.Fatalf("ParseCrlValidatedID: %v", err)
	}
	if validatedID.CrlIdentifier != nil {
		t.Errorf("CrlIdentifier = %v, want nil", validatedID.CrlIdentifier)
	}
	if got, want := validatedID.CrlHash.HashAlgorithm.Algorithm.String(), kat["crl.hashonly.hashalg"]; got != want {
		t.Errorf("hash algorithm = %s, want %s", got, want)
	}
	if got, want := hex.EncodeToString(validatedID.CrlHash.HashValue), kat["crl.hashonly.hashvalue"]; got != want {
		t.Errorf("hash value = %s, want %s", got, want)
	}
}

// TestParseCrlValidatedIDRejectsTrailingData checks the parser refuses extra data, so that a
// reference is never built from a truncated read.
func TestParseCrlValidatedIDRejectsTrailingData(t *testing.T) {
	kat := crlRefTestKAT(t, "testdata/crlocsp/kat_esf.txt")

	der := append(crlRefTestHex(t, kat["crl.hashonly.der"]), 0x00)
	if _, err := ParseCrlValidatedID(der); err == nil {
		t.Fatal("ParseCrlValidatedID accepted trailing data")
	}
}
