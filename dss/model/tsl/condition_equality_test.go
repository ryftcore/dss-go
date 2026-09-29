package tsl_test

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	dsstsl "github.com/ryftcore/dss-go/dss/tsl"
)

// TestConditionForQualifiersEqualsIsStructural pins Objects.equals(condition, that.condition)
// semantics: two independently built conditions with the same content are equal (the Condition
// implementations are pointer types, so == would compare identity), and different content is not.
func TestConditionForQualifiersEqualsIsStructural(t *testing.T) {
	qualifiers := []string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithSSCD"}

	a := tslmodel.NewConditionForQualifiers(dsstsl.NewPolicyIdCondition("0.4.0.194112.1.2"), qualifiers)
	b := tslmodel.NewConditionForQualifiers(dsstsl.NewPolicyIdCondition("0.4.0.194112.1.2"), qualifiers)
	if !a.Equals(b) {
		t.Fatal("ConditionForQualifiers with structurally equal conditions must be equal")
	}

	c := tslmodel.NewConditionForQualifiers(dsstsl.NewPolicyIdCondition("0.4.0.194112.1.3"), qualifiers)
	if a.Equals(c) {
		t.Fatal("ConditionForQualifiers with different conditions must not be equal")
	}

	composite := func() tslmodel.Condition {
		cc := dsstsl.NewCompositeCondition()
		cc.AddChild(dsstsl.NewPolicyIdCondition("1.2.3"))
		cc.AddChild(dsstsl.NewPolicyIdCondition("1.2.4"))
		return cc
	}
	d := tslmodel.NewConditionForQualifiers(composite(), qualifiers)
	e := tslmodel.NewConditionForQualifiers(composite(), qualifiers)
	if !d.Equals(e) {
		t.Fatal("ConditionForQualifiers with structurally equal composite conditions must be equal")
	}

	withNil := tslmodel.NewConditionForQualifiers(nil, qualifiers)
	if a.Equals(withNil) || withNil.Equals(a) {
		t.Fatal("a nil condition must not equal a non-nil one")
	}
}

// TestCertificateContentEquivalenceEqualsIsStructural is the CertificateContentEquivalence twin
// of TestConditionForQualifiersEqualsIsStructural.
func TestCertificateContentEquivalenceEqualsIsStructural(t *testing.T) {
	build := func(oid string) *tslmodel.CertificateContentEquivalence {
		c := tslmodel.NewCertificateContentEquivalence()
		c.SetContext(enumerations.MRAEquivalenceContextQCCompliance)
		c.SetCondition(dsstsl.NewPolicyIdCondition(oid))
		return c
	}

	if !build("0.4.0.194112.1.2").Equals(build("0.4.0.194112.1.2")) {
		t.Fatal("CertificateContentEquivalence with structurally equal conditions must be equal")
	}
	if build("0.4.0.194112.1.2").Equals(build("0.4.0.194112.1.3")) {
		t.Fatal("CertificateContentEquivalence with different conditions must not be equal")
	}
}
