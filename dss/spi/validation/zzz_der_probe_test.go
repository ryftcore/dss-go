package validation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestZZZDERProbe(t *testing.T) {
	path := os.Getenv("PROBE_FILE")
	if path == "" {
		t.Skip("PROBE_FILE not set")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	tok, err := NewTimestampToken(data, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatal(err)
	}
	der := tok.Encoded()
	sum := sha256.Sum256(der)
	fmt.Println("GO len=", len(der), "sha256=", hex.EncodeToString(sum[:]))
	outPath := os.Getenv("PROBE_OUT")
	if outPath != "" {
		os.WriteFile(outPath, der, 0644)
	}
}
