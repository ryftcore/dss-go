package tsl

import (
	"testing"

	"github.com/utain/esig/dss/model/timedependent"
)

func TestNewTrustPropertiesPanicsOnNilTLInfo(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nil tlInfo")
		}
		if r != "tlInfo cannot be null!" {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()
	NewTrustProperties(nil, NewTrustServiceProvider(), &timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]{})
}

func TestNewTrustPropertiesPanicsOnNilTrustServiceProvider(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nil trustServiceProvider")
		}
		if r != "trustServiceProvider cannot be null!" {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()
	NewTrustProperties(&TLInfo{}, nil, &timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]{})
}

func TestNewTrustPropertiesPanicsOnNilTrustService(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nil trustService")
		}
		if r != "trustService cannot be null!" {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()
	NewTrustProperties(&TLInfo{}, NewTrustServiceProvider(), nil)
}

func TestNewTrustPropertiesRoundTrip(t *testing.T) {
	tlInfo := &TLInfo{}
	tsp := NewTrustServiceProvider()
	trustService := &timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]{}

	tp := NewTrustProperties(tlInfo, tsp, trustService)
	if tp.LOTLInfo() != nil {
		t.Fatalf("expected nil LOTLInfo for the independent-TL constructor, got %v", tp.LOTLInfo())
	}
	if tp.TLInfo() != tlInfo {
		t.Fatal("expected TLInfo() to return the constructor argument")
	}
	if tp.TrustServiceProvider() != tsp {
		t.Fatal("expected TrustServiceProvider() to return the constructor argument")
	}
	if tp.TrustService() != trustService {
		t.Fatal("expected TrustService() to return the constructor argument")
	}

	lotlInfo := &LOTLInfo{}
	tpWithLotl := NewTrustPropertiesWithLOTL(lotlInfo, tlInfo, tsp, trustService)
	if tpWithLotl.LOTLInfo() != lotlInfo {
		t.Fatal("expected LOTLInfo() to return the constructor argument")
	}
}
