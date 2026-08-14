// Ported from dss-enumerations/.../LoTEServiceTypeIdentifier.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestLoTEServiceTypeIdentifierFromURI_NoLoaderMatch(t *testing.T) {
	miss := &fakeLoTELoader{}
	withLoTELoaders(t, miss)

	if got := LoTEServiceTypeIdentifierFromURI("urn:unknown"); got != nil {
		t.Errorf("LoTEServiceTypeIdentifierFromURI = %v, want nil", got)
	}
}

func TestLoTEServiceTypeIdentifierFromURI_LoaderMatch(t *testing.T) {
	want := &loteEmptyServiceTypeIdentifier{uri: "urn:known"}
	loader := &fakeLoTELoader{sti: want}
	withLoTELoaders(t, loader)

	got := LoTEServiceTypeIdentifierFromURI("urn:known")
	if got != LoTEServiceTypeIdentifier(want) {
		t.Errorf("LoTEServiceTypeIdentifierFromURI = %v, want %v", got, want)
	}
}
