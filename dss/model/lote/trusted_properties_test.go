package lote

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model/timedependent"
)

func TestNewTrustedPropertiesPanicsOnNilLoTEInfo(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nil listInfo")
		}
		if r != "tlInfo cannot be null!" {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()
	NewTrustedProperties(nil, NewTrustedEntity(), &timedependent.Values[ServiceStatusAndInformationExtensions]{})
}

func TestNewTrustedPropertiesPanicsOnNilTrustedEntity(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nil trustedEntity")
		}
		if r != "trustedEntity cannot be null!" {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()
	NewTrustedProperties(&Info{}, nil, &timedependent.Values[ServiceStatusAndInformationExtensions]{})
}

func TestNewTrustedPropertiesPanicsOnNilTrustedServices(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nil trustedServices")
		}
		if r != "trustedServices cannot be null!" {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()
	NewTrustedProperties(&Info{}, NewTrustedEntity(), nil)
}

func TestNewTrustedPropertiesRoundTrip(t *testing.T) {
	listInfo := &Info{}
	entity := NewTrustedEntity()
	services := &timedependent.Values[ServiceStatusAndInformationExtensions]{}

	tp := NewTrustedProperties(listInfo, entity, services)
	if tp.LoLoTEInfo() != nil {
		t.Fatalf("expected nil LoLoTEInfo() for the independent-list constructor, got %v", tp.LoLoTEInfo())
	}
	if tp.LoTEInfo() != listInfo {
		t.Fatal("expected LoTEInfo() to return the constructor argument")
	}
	if tp.TrustedEntity() != entity {
		t.Fatal("expected TrustedEntity() to return the constructor argument")
	}
	if tp.TrustedServices() != services {
		t.Fatal("expected TrustedServices() to return the constructor argument")
	}

	lolote := &LoLoTEInfo{}
	tpWithLolote := NewTrustedPropertiesWithLoLoTE(lolote, listInfo, entity, services)
	if tpWithLolote.LoLoTEInfo() != lolote {
		t.Fatal("expected LoLoTEInfo() to return the constructor argument")
	}
}
