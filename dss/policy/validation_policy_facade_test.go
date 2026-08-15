// Ported from dss-policy-jaxb's ValidationPolicyFacadeTest (DSS 6.5.RC1)
// where such a test exists; otherwise round-trip coverage for this file's
// own collapse of AbstractJaxbFacade (see validation_policy_facade.go's
// header).
package policy

import (
	"bytes"
	"os"
	"testing"
)

func TestValidationPolicyFacadeGetValidationPolicy(t *testing.T) {
	data, err := os.ReadFile("jaxb/testdata/policy/constraint.xml")
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	facade := NewValidationPolicyFacade()
	policy, err := facade.GetValidationPolicy(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("GetValidationPolicy: %v", err)
	}
	if got, want := policy.PolicyName(), "QES AES/QC AES TL based"; got != want {
		t.Errorf("PolicyName() = %q, want %q", got, want)
	}
}

func TestValidationPolicyFacadeGetValidationPolicyNilReader(t *testing.T) {
	facade := NewValidationPolicyFacade()
	if _, err := facade.GetValidationPolicy(nil); err == nil {
		t.Error("GetValidationPolicy(nil) = nil error, want an error")
	}
}

func TestValidationPolicyFacadeMarshalUnmarshalRoundTrip(t *testing.T) {
	data, err := os.ReadFile("jaxb/testdata/oracle/constraint.remarshal.xml")
	if err != nil {
		t.Fatalf("read oracle: %v", err)
	}

	facade := NewValidationPolicyFacade()
	cp, err := facade.Unmarshal(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	got, err := facade.Marshal(cp)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Error("Marshal(Unmarshal(oracle)) != oracle; ValidationPolicyFacade should preserve jaxb.Marshal/Unmarshal's byte-exact round-trip")
	}
}

func TestValidationPolicyFacadeGetValidationPolicyFromFile(t *testing.T) {
	facade := NewValidationPolicyFacade()
	policy, err := facade.GetValidationPolicyFromFile("jaxb/testdata/policy/eaa-constraint.xml")
	if err != nil {
		t.Fatalf("GetValidationPolicyFromFile: %v", err)
	}
	if got, want := policy.PolicyName(), "Certificate policy TL based"; got != want {
		t.Errorf("PolicyName() = %q, want %q", got, want)
	}
}
