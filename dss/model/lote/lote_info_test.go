package lote

import "testing"

func TestLoTEInfoRoundTrip(t *testing.T) {
	pointer := NewOtherListPointer()
	l := NewInfoFull(nil, nil, nil, "https://example.org/lote.xml", nil, pointer)

	if l.Url() != "https://example.org/lote.xml" {
		t.Fatalf("unexpected Url(): %s", l.Url())
	}
	if l.Parent() != nil {
		t.Fatalf("expected nil Parent(), got %v", l.Parent())
	}
	if l.ListPointer() != pointer {
		t.Fatal("expected ListPointer() to return the constructor argument")
	}
}

func TestLoTEInfoDSSIDIsCachedAndUsesLoTEIdentifier(t *testing.T) {
	l := NewInfo(nil, nil, nil, "https://example.org/lote.xml")

	id := l.DSSID()
	if _, ok := id.(*Identifier); !ok {
		t.Fatalf("expected DSSID() to return a *LoTEIdentifier, got %T", id)
	}
	if l.DSSID() != id {
		t.Fatal("expected DSSID() to cache and return the same identifier on subsequent calls")
	}
	if l.DSSIDAsString() != id.AsXmlID() {
		t.Fatalf("DSSIDAsString() = %q, want %q", l.DSSIDAsString(), id.AsXmlID())
	}
}

func TestLoLoTEInfoDSSIDUsesLoLoTEIdentifier(t *testing.T) {
	// LoLoTEInfo overrides buildIdentifier(); DSSID() called on the LoLoTEInfo must reflect the
	// override, not the embedded Info's own Identifier.
	l := NewLoLoTEInfo(nil, nil, nil, "https://example.org/lolote.xml")

	id := l.DSSID()
	if _, ok := id.(*LoLoTEIdentifier); !ok {
		t.Fatalf("expected DSSID() to return a *LoLoTEIdentifier, got %T", id)
	}

	children := []*Info{NewInfoWithParent(nil, nil, nil, "https://example.org/child.xml", l)}
	l.SetChildrenInfos(children)
	if len(l.ChildrenInfos()) != 1 || l.ChildrenInfos()[0].Parent() != l {
		t.Fatalf("unexpected ChildrenInfos(): %v", l.ChildrenInfos())
	}
}
