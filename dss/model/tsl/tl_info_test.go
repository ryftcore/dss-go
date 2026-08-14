package tsl

import "testing"

func TestTLInfoConstructorsAndAccessors(t *testing.T) {
	otherPtr := NewOtherTSLPointer()
	parent := &LOTLInfo{TLInfo: *NewTLInfo(nil, nil, nil, "https://example.org/lotl.xml")}
	tlInfo := NewTLInfoFull(nil, nil, nil, "https://example.org/tl.xml", parent, otherPtr)

	if tlInfo.Url() != "https://example.org/tl.xml" {
		t.Fatalf("unexpected Url: %s", tlInfo.Url())
	}
	if tlInfo.Parent() != parent {
		t.Fatalf("unexpected Parent: %v", tlInfo.Parent())
	}
	if tlInfo.OtherTSLPointer() != otherPtr {
		t.Fatalf("unexpected OtherTSLPointer: %v", tlInfo.OtherTSLPointer())
	}
	if tlInfo.DownloadCacheInfo() != nil {
		t.Fatalf("expected nil DownloadCacheInfo")
	}

	plain := NewTLInfo(nil, nil, nil, "https://example.org/plain.xml")
	if plain.Parent() != nil {
		t.Fatalf("expected nil Parent for the 4-arg constructor")
	}
	if plain.OtherTSLPointer() != nil {
		t.Fatalf("expected nil OtherTSLPointer for the 4-arg constructor")
	}
}

func TestTLInfoDSSIDIsCachedAndURLBased(t *testing.T) {
	tlInfo := NewTLInfo(nil, nil, nil, "https://example.org/tl.xml")

	id1 := tlInfo.DSSID()
	id2 := tlInfo.DSSID()
	if id1 != id2 {
		t.Fatalf("expected DSSID() to be cached and return the same instance")
	}

	trustedID, ok := id1.(*TrustedListIdentifier)
	if !ok {
		t.Fatalf("expected DSSID() of a TLInfo to be a *TrustedListIdentifier, got %T", id1)
	}
	if trustedID.AsXmlID()[:3] != "TL-" {
		t.Fatalf("expected TL- prefix, got %s", trustedID.AsXmlID())
	}
	if tlInfo.DSSIDAsString() != trustedID.AsXmlID() {
		t.Fatalf("DSSIDAsString() = %s, want %s", tlInfo.DSSIDAsString(), trustedID.AsXmlID())
	}
}
