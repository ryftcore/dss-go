package tsl

import (
	"testing"
	"time"

	"github.com/utain/esig/dss/enumerations"
)

func TestServiceEquivalenceBuilderRoundTrip(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	typeAsi := map[ServiceTypeASi]ServiceTypeASi{
		{typ: "t1"}: {typ: "t2"},
	}
	statusEquiv := []StatusEquivalenceMapping{
		{PointedStatuses: []string{"granted"}, PointingStatuses: []string{"recognised"}},
	}
	qualifierEquiv := map[string]string{"q1": "q2"}

	se := NewServiceEquivalenceBuilder().
		SetLegalInfoIdentifier("legal-id").
		SetStatus(enumerations.MRAStatus_ENACTED).
		SetStartDate(start).
		SetEndDate(end).
		SetTypeAsiEquivalence(typeAsi).
		SetStatusEquivalence(statusEquiv).
		SetQualifierEquivalence(qualifierEquiv).
		Build()

	if se.LegalInfoIdentifier() != "legal-id" {
		t.Fatalf("unexpected LegalInfoIdentifier: %s", se.LegalInfoIdentifier())
	}
	if se.Status() != enumerations.MRAStatus_ENACTED {
		t.Fatalf("unexpected Status: %s", se.Status())
	}
	if len(se.TypeAsiEquivalence()) != 1 {
		t.Fatalf("unexpected TypeAsiEquivalence: %v", se.TypeAsiEquivalence())
	}
	if len(se.StatusEquivalence()) != 1 || se.StatusEquivalence()[0].PointedStatuses[0] != "granted" {
		t.Fatalf("unexpected StatusEquivalence: %v", se.StatusEquivalence())
	}
	if se.QualifierEquivalence()["q1"] != "q2" {
		t.Fatalf("unexpected QualifierEquivalence: %v", se.QualifierEquivalence())
	}
}

func TestServiceEquivalenceEquals(t *testing.T) {
	build := func(legalID string) *ServiceEquivalence {
		return NewServiceEquivalenceBuilder().SetLegalInfoIdentifier(legalID).Build()
	}
	a := build("same")
	b := build("same")
	if !a.Equals(b) {
		t.Fatal("expected equal ServiceEquivalence to be Equals()")
	}
	c := build("different")
	if a.Equals(c) {
		t.Fatal("expected different legalInfoIdentifier to not be Equals()")
	}
	if a.Equals(nil) {
		t.Fatal("expected Equals(nil) to be false")
	}
}
