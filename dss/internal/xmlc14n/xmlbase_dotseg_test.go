package xmlc14n

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

// refRemoveDotSegments is the original strings.Builder form of removeDotSegments and
// dotDotOutput, kept as the reference the byte-slice form must reproduce exactly.
func refRemoveDotSegments(path string) (string, error) {
	input := path
	for strings.Contains(input, "//") {
		input = strings.ReplaceAll(input, "//", "/")
	}
	if input == "" {
		return "", ErrXMLBaseUnjoinable
	}

	var output strings.Builder
	if input[0] == '/' {
		output.WriteByte('/')
		input = input[1:]
	}

	dotDot := func(input string) string {
		s := output.String()
		switch {
		case len(s) == 0:
			output.WriteByte('/')
		case strings.HasSuffix(s, "../"):
			output.WriteString("..")
		case strings.HasSuffix(s, ".."):
			output.WriteString("/..")
		default:
			index := strings.LastIndexByte(s, '/')
			if index == -1 {
				output.Reset()
				if len(input) > 0 && input[0] == '/' {
					input = input[1:]
				}
			} else {
				output.Reset()
				output.WriteString(s[:index])
			}
		}
		return input
	}

	for len(input) != 0 {
		switch {
		case strings.HasPrefix(input, "./"):
			input = input[2:]
		case strings.HasPrefix(input, "../"):
			input = input[3:]
			if output.String() != "/" {
				output.WriteString("../")
			}
		case strings.HasPrefix(input, "/./"):
			input = input[2:]
		case input == "/.":
			input = "/"
		case strings.HasPrefix(input, "/../"):
			input = input[3:]
			input = dotDot(input)
		case input == "/..":
			input = "/"
			input = dotDot(input)
		case input == ".":
			input = ""
		case input == "..":
			if output.String() != "/" {
				output.WriteString("..")
			}
			input = ""
		default:
			var end int
			begin := strings.IndexByte(input, '/')
			if begin == 0 {
				end = strings.IndexByte(input[1:], '/')
				if end >= 0 {
					end++
				}
			} else {
				end = begin
				begin = 0
			}
			var segment string
			if end == -1 {
				segment = input[begin:]
				input = ""
			} else {
				segment = input[begin:end]
				input = input[end:]
			}
			output.WriteString(segment)
		}
	}

	out := output.String()
	if strings.HasSuffix(out, "..") {
		out += "/"
	}
	return out, nil
}

func TestRemoveDotSegmentsMatchesBuilderForm(t *testing.T) {
	rng := rand.New(rand.NewSource(4))
	pieces := []string{"a", "b", "..", ".", "/", "//", "/..", "/.", "../", "./", "/../", "/./", "ab"}
	for i := 0; i < 50000; i++ {
		var sb strings.Builder
		for j, n := 0, rng.Intn(9); j < n; j++ {
			sb.WriteString(pieces[rng.Intn(len(pieces))])
		}
		path := sb.String()
		got, gotErr := removeDotSegments(path)
		want, wantErr := refRemoveDotSegments(path)
		if got != want || gotErr != wantErr {
			t.Fatalf("removeDotSegments(%q) = %q, %v; reference %q, %v", path, got, gotErr, want, wantErr)
		}
	}
}

// TestRemoveDotSegmentsIsLinear is the X08-PERF-003 regression: n/2 segments followed by n/2
// "../" made every "/.." copy the whole output so far (5 s for a 500 KB xml:base).
func TestRemoveDotSegmentsIsLinear(t *testing.T) {
	const k = 100000
	path := "/" + strings.Repeat("a/", k) + strings.Repeat("../", k)
	start := time.Now()
	got, err := removeDotSegments(path)
	if err != nil || got != "/" {
		t.Fatalf("removeDotSegments = %q, %v; want \"/\"", got, err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("removeDotSegments took %v for a %d KB path; want linear", elapsed, len(path)>>10)
	}
}
