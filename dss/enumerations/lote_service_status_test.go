// Ported from dss-enumerations/.../LoTEServiceStatus.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestLoTEServiceStatusFromURI_NoLoaderMatch(t *testing.T) {
	miss := &fakeLoTELoader{}
	withLoTELoaders(t, miss)

	if got := LoTEServiceStatusFromURI("urn:unknown"); got != nil {
		t.Errorf("LoTEServiceStatusFromURI = %v, want nil", got)
	}
}

func TestLoTEServiceStatusFromURI_LoaderMatch(t *testing.T) {
	want := &loteEmptyServiceStatus{uri: "urn:known"}
	loader := &fakeLoTELoader{status: want}
	withLoTELoaders(t, loader)

	got := LoTEServiceStatusFromURI("urn:known")
	if got != LoTEServiceStatus(want) {
		t.Errorf("LoTEServiceStatusFromURI = %v, want %v", got, want)
	}
}
