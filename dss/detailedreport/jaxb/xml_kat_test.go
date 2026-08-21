package jaxb

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/utain/esig/dss/internal/corpustest"
)

// oracleDir returns the directory holding the Java-produced detailed-report
// dumps. The checked-in corpus lives in testdata/oracle; set
// DSS_DETAILEDREPORT_ORACLE_DIR to replay the KAT over a larger local corpus.
func oracleDir(t *testing.T) string {
	t.Helper()
	if dir := os.Getenv("DSS_DETAILEDREPORT_ORACLE_DIR"); dir != "" {
		return dir
	}
	return corpustest.Path(t, "oracle")
}

// TestMarshalParity is the marshal-parity KAT: every Java-produced dump must
// unmarshal into the Go model and marshal back byte for byte.
func TestMarshalParity(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(oracleDir(t), "*.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatalf("no oracle dumps under %s", oracleDir(t))
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			want, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			dr, err := Unmarshal(want)
			if err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			got, err := Marshal(dr)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("round-trip differs:\n%s", firstDiff(want, got))
			}
		})
	}
}

// firstDiff renders the first differing line of two documents.
func firstDiff(want, got []byte) string {
	wl := bytes.Split(want, []byte("\n"))
	gl := bytes.Split(got, []byte("\n"))
	for i := 0; i < len(wl) && i < len(gl); i++ {
		if !bytes.Equal(wl[i], gl[i]) {
			return "line " + itoa(i+1) + "\n  java: " + string(wl[i]) + "\n  go:   " + string(gl[i])
		}
	}
	if len(wl) != len(gl) {
		return "line counts differ: java=" + itoa(len(wl)) + " go=" + itoa(len(gl))
	}
	return "documents differ outside line boundaries"
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
