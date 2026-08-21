// Oracle test for JavaHashSetOrder / JavaHashSetStringOrder
// (abstract_detailed_report_builder.go).
//
// testdata/oracle/hash_order.tsv is a pure Java dump, produced by
// testdata/oracle/gen/HashOrder.java: for each input sequence it records the
// order a java.util.HashSet iterates it in, once with the
// AbstractTokenProxy hashCode (31 + id.hashCode(), the TOKEN rows) and once
// with a plain String key (the STRING rows).

package executor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/utain/esig/dss/internal/corpustest"
)

// tokenID is a minimal diagnostic.TokenProxy stand-in: JavaHashSetOrder only
// reads Id().
type tokenID string

// Id implements IdentifiedToken.
func (t tokenID) Id() string { return string(t) }

func TestJavaHashSetOrderOracle(t *testing.T) {
	data, err := os.ReadFile(corpustest.Path(t, filepath.Join("oracle", "hash_order.tsv")))
	if err != nil {
		t.Fatalf("reading the oracle dump: %v", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			t.Fatalf("malformed oracle row: %q", line)
		}
		kind, input, want := fields[0], strings.Split(fields[1], "|"), fields[2]
		t.Run(kind+"/"+fields[1], func(t *testing.T) {
			var got string
			switch kind {
			case "TOKEN":
				tokens := make([]tokenID, 0, len(input))
				for _, id := range input {
					tokens = append(tokens, tokenID(id))
				}
				ordered := JavaHashSetOrder(tokens)
				parts := make([]string, 0, len(ordered))
				for _, token := range ordered {
					parts = append(parts, string(token))
				}
				got = strings.Join(parts, "|")
			case "STRING":
				got = strings.Join(JavaHashSetStringOrder(input), "|")
			default:
				t.Fatalf("unknown oracle row kind %q", kind)
			}
			if got != want {
				t.Errorf("order\n got: %s\nwant: %s", got, want)
			}
		})
	}
}
