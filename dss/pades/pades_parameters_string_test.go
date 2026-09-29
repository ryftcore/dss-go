package pades

import (
	"fmt"
	"strings"
	"testing"
)

// The parameters' String() must not disclose the document password (see the DIVERGENCE note on
// both methods).
func TestParametersStringDoesNotDiscloseThePassword(t *testing.T) {
	password := []byte("s3cret-Pw")
	// The renderings upstream's Arrays.toString and Go's %v would produce.
	leaks := []string{string(password), fmt.Sprint(password), "[115 51"}

	signatureParameters := NewSignatureParameters()
	timestampParameters := NewTimestampParameters()
	for name, render := range map[string]func() string{
		"signature": signatureParameters.String,
		"timestamp": timestampParameters.String,
	} {
		if got := render(); !strings.Contains(got, "passwordProtection=null") {
			t.Errorf("%s parameters without a password: %q, want passwordProtection=null", name, got)
		}
	}

	signatureParameters.SetPasswordProtection(password)
	timestampParameters.SetPasswordProtection(password)
	for name, got := range map[string]string{
		"signature": signatureParameters.String(),
		"timestamp": timestampParameters.String(),
	} {
		if !strings.Contains(got, "passwordProtection=[redacted]") {
			t.Errorf("%s parameters: %q, want passwordProtection=[redacted]", name, got)
		}
		for _, leak := range leaks {
			if strings.Contains(got, leak) {
				t.Errorf("%s parameters disclose the password as %q: %q", name, leak, got)
			}
		}
	}
}
