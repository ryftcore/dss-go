package tsl

import (
	"testing"

	"github.com/utain/esig/dss/model/timedependent"
)

func TestMRARoundTrip(t *testing.T) {
	m := NewMRA()
	m.SetTechnicalType("tech")
	m.SetVersion("1")
	m.SetPointingContractingPartyLegislation("pointing")
	m.SetPointedContractingPartyLegislation("pointed")

	equivalences := []*timedependent.MutableTimeDependentValues[*ServiceEquivalence]{}
	m.SetServiceEquivalence(equivalences)

	if m.TechnicalType() != "tech" {
		t.Fatalf("unexpected TechnicalType: %s", m.TechnicalType())
	}
	if m.Version() != "1" {
		t.Fatalf("unexpected Version: %s", m.Version())
	}
	if m.PointingContractingPartyLegislation() != "pointing" {
		t.Fatalf("unexpected PointingContractingPartyLegislation: %s", m.PointingContractingPartyLegislation())
	}
	if m.PointedContractingPartyLegislation() != "pointed" {
		t.Fatalf("unexpected PointedContractingPartyLegislation: %s", m.PointedContractingPartyLegislation())
	}
	if len(m.ServiceEquivalence()) != 0 {
		t.Fatalf("expected empty ServiceEquivalence, got %d", len(m.ServiceEquivalence()))
	}
}
