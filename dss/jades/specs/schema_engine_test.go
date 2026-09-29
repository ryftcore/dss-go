package specs

import (
	"fmt"
	"sync"
	"testing"
	"testing/fstest"
)

// TestIsIntegerNumber pins which number literals the "integer" type accepts, as the validator
// upstream (jsonsKema 0.31.0, Validator.findActualNumberType) decides it: a Double whose
// Double.toString() has only zeros after the '.' - which excludes everything from 1e7 on, where
// Double.toString switches to scientific notation (J17B-STD-001).
func TestIsIntegerNumber(t *testing.T) {
	cases := map[string]bool{
		// no '.', 'e' or 'E': an Int / BigInteger
		"0":                       true,
		"-0":                      true,
		"1":                       true,
		"-17":                     true,
		"1700000000":              true,
		"12345678901234567890123": true,
		// whole Doubles below 1e7: Double.toString is "N.0"
		"1.0":       true,
		"1.000":     true,
		"0.0":       true,
		"-0.0":      true,
		"12.0":      true,
		"5e0":       true,
		"1E0":       true,
		"100e-2":    true, // 1.0
		"10e-1":     true, // 1.0
		"1.0e-0":    true,
		"1e2":       true, // 100.0
		"1.5e1":     true, // 15.0
		"9999999.0": true,
		"1e-999":    true, // underflows to 0.0
		// fractions
		"1.5":    false,
		"0.5":    false,
		"-1.25":  false,
		"1e-1":   false,
		"1.5e0":  false,
		"0.0001": false,
		// from 1e7 on Double.toString is "1.0E7", "1.7E9": not an integer upstream
		"1e7":          false,
		"10000000.0":   false,
		"1700000000.0": false,
		"1.7e9":        false,
		"1700000000.5": false,
		// out of range for a Double: BigDecimal.toString, integer for a one-digit coefficient
		"1e999":   true,
		"1.5e999": false,
		"10e998":  false,
	}
	for literal, want := range cases {
		if got := isIntegerNumber(literal); got != want {
			t.Errorf("isIntegerNumber(%q) = %v, want %v", literal, got, want)
		}
	}
}

// TestIntegerTypeAcceptsWholeFloats checks the engine end to end: a schema declaring
// "type":"integer" accepts 1.0 (legal JSON, legal draft-07) and still rejects a real fraction.
func TestIntegerTypeAcceptsWholeFloats(t *testing.T) {
	files := fstest.MapFS{
		"s.json": {Data: []byte(`{"properties":{"iat":{"type":"integer"},"n":{"type":"number"}}}`)},
	}
	defs := map[string]string{"s.json": "s.json"}

	valid := []string{
		`{"iat":1700000000}`, `{"iat":1.0}`, `{"iat":12.0}`, `{"iat":1e2}`, `{"iat":1,"n":1}`, `{"iat":1,"n":1.5}`,
		`{"n":1700000000.0}`,
	}
	for _, doc := range valid {
		if errs := validateJSONAgainstSchema(files, defs, "s.json", doc); len(errs) != 0 {
			t.Errorf("%s: expected no errors, got %v", doc, errs)
		}
	}
	invalid := []string{`{"iat":1.5}`, `{"iat":"1"}`, `{"iat":true}`, `{"iat":1e-1}`, `{"iat":1700000000.0}`}
	for _, doc := range invalid {
		assertErrorFound(t, validateJSONAgainstSchema(files, defs, "s.json", doc), "iat")
	}
}

// TestProtectedHeaderIatAsWholeFloat exercises the embedded protected-header schema, whose 'iat'
// is declared as an integer (rfcs/rfc7519.json).
func TestProtectedHeaderIatAsWholeFloat(t *testing.T) {
	u := JAdESProtectedHeaderUtilsInstance()
	const header = `{"alg":"RS256","x5t#S256":"AAAA","iat":%s}`
	for literal, valid := range map[string]bool{
		"1700000000":   true,
		"1.0":          true,
		"1700000000.5": false,
		"1700000000.0": false, // 1.7E9 for upstream's Double.toString
	} {
		errs := u.ValidateAgainstSchema(fmt.Sprintf(header, literal))
		if valid {
			assertNoErrors(t, errs)
		} else {
			assertErrorFound(t, errs, "iat")
		}
	}
}

// TestEmbeddedSchemaCacheIsConcurrencySafe drives the process-wide decoded-schema and compiled-
// pattern caches from several goroutines (run under -race), checking that verdicts do not change.
func TestEmbeddedSchemaCacheIsConcurrencySafe(t *testing.T) {
	u := JAdESProtectedHeaderUtilsInstance()
	good := `{"alg":"RS256","x5t#S256":"AAAA","sigT":"2020-01-01T00:00:00Z"}`
	bad := `{"alg":"RS256","x5t#S256":"AAAA","iat":"not a number"}`

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				if errs := u.ValidateAgainstSchema(good); len(errs) != 0 {
					t.Errorf("valid header rejected: %v", errs)
					return
				}
				if errs := u.ValidateAgainstSchema(bad); len(errs) == 0 {
					t.Error("invalid header accepted")
					return
				}
			}
		}()
	}
	wg.Wait()
}
