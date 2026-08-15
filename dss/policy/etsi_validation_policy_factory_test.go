// Ported from dss-policy-jaxb's EtsiValidationPolicyFactoryTest (DSS
// 6.5.RC1) where such a test exists, plus coverage of this port's embedded
// default-policy loading, which upstream's classpath-resource mechanism has
// no direct Go equivalent for (see etsi_validation_policy_factory.go).
package policy

import (
	"bytes"
	"os"
	"testing"

	"github.com/utain/esig/dss/model"
)

func TestEtsiValidationPolicyFactoryLoadDefaultValidationPolicy(t *testing.T) {
	factory := NewEtsiValidationPolicyFactory()
	policy := factory.LoadDefaultValidationPolicy()
	if policy == nil {
		t.Fatal("LoadDefaultValidationPolicy returned nil")
	}
	if got, want := policy.PolicyName(), "QES AES/QC AES TL based"; got != want {
		t.Errorf("PolicyName() = %q, want %q", got, want)
	}

	// The embedded default policy must be byte-identical to upstream's
	// src/main/resources/policy/constraint.xml (also mirrored at
	// jaxb/testdata/policy/constraint.xml).
	onDisk, err := os.ReadFile("jaxb/testdata/policy/constraint.xml")
	if err != nil {
		t.Fatalf("read jaxb/testdata/policy/constraint.xml: %v", err)
	}
	if !bytes.Equal(defaultValidationPolicy, onDisk) {
		t.Error("embedded resources/constraint.xml differs from jaxb/testdata/policy/constraint.xml")
	}
}

func TestEtsiValidationPolicyFactoryIsSupported(t *testing.T) {
	factory := NewEtsiValidationPolicyFactory()

	valid := model.NewInMemoryDocument(defaultValidationPolicy)
	if !factory.IsSupported(valid) {
		t.Error("IsSupported(constraint.xml) = false, want true")
	}

	invalid := model.NewInMemoryDocument([]byte("not xml at all"))
	if factory.IsSupported(invalid) {
		t.Error("IsSupported(garbage) = true, want false")
	}
}

func TestEtsiValidationPolicyFactoryLoadValidationPolicy(t *testing.T) {
	factory := NewEtsiValidationPolicyFactory()
	doc := model.NewInMemoryDocument(defaultValidationPolicy)
	policy := factory.LoadValidationPolicy(doc)
	if got, want := policy.PolicyName(), "QES AES/QC AES TL based"; got != want {
		t.Errorf("PolicyName() = %q, want %q", got, want)
	}
}

func TestEtsiValidationPolicyFactoryLoadValidationPolicyFromReaderPanicsOnGarbage(t *testing.T) {
	factory := NewEtsiValidationPolicyFactory()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic loading a non-policy document, got none")
		}
	}()
	factory.LoadValidationPolicyFromReader(bytes.NewReader([]byte("<not-a-policy/>")))
}
