package cmscore

import (
	"math/big"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// accuracyInteger encodes an INTEGER, or an IMPLICIT [tag] INTEGER when tag >= 0.
func accuracyInteger(t *testing.T, value *big.Int, tag int) []byte {
	t.Helper()
	encoded := asn1ber.EncodeInteger(value)
	if tag < 0 {
		return encoded
	}
	// Replace the universal INTEGER identifier with the context-specific one.
	return append([]byte{0x80 | byte(tag)}, encoded[1:]...)
}

func accuracyElement(t *testing.T, components ...[]byte) *asn1ber.Element {
	t.Helper()
	var body []byte
	for _, component := range components {
		body = append(body, component...)
	}
	element, rest, err := asn1ber.Parse(asn1ber.WriteSequence(body))
	if err != nil || len(rest) != 0 {
		t.Fatalf("crafted Accuracy does not parse: %v (%d trailing bytes)", err, len(rest))
	}
	return element
}

// TestAccuracyRangeChecks: accuracyFromElement used to read every component with
// int(Integer().Int64()), which turns an INTEGER outside the int64 range into 0 without a word
// and accepted a millis or micros outside (1..999). BouncyCastle's Accuracy(ASN1Sequence) refuses
// the latter - IllegalArgumentException, or ArithmeticException from intValueExact when the value
// does not fit an int - and, as TimeStampToken builds TSTInfo while it is constructed, that
// makes the whole token unparseable there. The verdicts below were taken from a real
// BouncyCastle 1.78.1 run.
func TestAccuracyRangeChecks(t *testing.T) {
	huge := new(big.Int).Lsh(big.NewInt(1), 70)
	cases := []struct {
		name       string
		components [][]byte
		wantError  bool
	}{
		{"seconds only", [][]byte{accuracyInteger(t, big.NewInt(1), -1)}, false},
		{"negative seconds", [][]byte{accuracyInteger(t, big.NewInt(-5), -1)}, false},
		{"all three", [][]byte{
			accuracyInteger(t, big.NewInt(1), -1),
			accuracyInteger(t, big.NewInt(500), 0),
			accuracyInteger(t, big.NewInt(100), 1),
		}, false},
		{"millis lower bound", [][]byte{accuracyInteger(t, big.NewInt(1), 0)}, false},
		{"millis upper bound", [][]byte{accuracyInteger(t, big.NewInt(999), 0)}, false},
		{"micros lower bound", [][]byte{accuracyInteger(t, big.NewInt(1), 1)}, false},
		{"micros upper bound", [][]byte{accuracyInteger(t, big.NewInt(999), 1)}, false},
		{"millis zero", [][]byte{accuracyInteger(t, big.NewInt(0), 0)}, true},
		{"millis 1000", [][]byte{accuracyInteger(t, big.NewInt(1000), 0)}, true},
		{"millis negative", [][]byte{accuracyInteger(t, big.NewInt(-1), 0)}, true},
		{"millis beyond int64", [][]byte{accuracyInteger(t, huge, 0)}, true},
		{"micros zero", [][]byte{accuracyInteger(t, big.NewInt(0), 1)}, true},
		{"micros 1000", [][]byte{accuracyInteger(t, big.NewInt(1000), 1)}, true},
		{"micros beyond int64", [][]byte{accuracyInteger(t, huge, 1)}, true},
		{"seconds beyond int64", [][]byte{accuracyInteger(t, huge, -1)}, true},
		{"constructed millis", [][]byte{{0xA0, 0x03, 0x02, 0x01, 0x05}}, true},
		{"unknown tag", [][]byte{accuracyInteger(t, big.NewInt(5), 2)}, true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			accuracy, err := accuracyFromElement(accuracyElement(t, testCase.components...))
			if testCase.wantError {
				if err == nil {
					t.Fatalf("accepted, got %+v", accuracy)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if testCase.name == "all three" {
				if *accuracy.Seconds != 1 || *accuracy.Millis != 500 || *accuracy.Micros != 100 {
					t.Errorf("accuracy = %d/%d/%d", *accuracy.Seconds, *accuracy.Millis, *accuracy.Micros)
				}
			}
		})
	}
}
