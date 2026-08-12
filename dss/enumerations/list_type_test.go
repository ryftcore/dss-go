// Ported from dss-enumerations/.../ListType.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestListTypeFromURI_NoLoaderMatch(t *testing.T) {
	miss := &fakeLoTELoader{}
	withLoTELoaders(t, miss)

	if got := ListTypeFromURI("urn:unknown"); got != nil {
		t.Errorf("ListTypeFromURI = %v, want nil", got)
	}
}

func TestListTypeFromURI_LoaderMatch(t *testing.T) {
	want := &loteEmptyListType{uri: "urn:known"}
	loader := &fakeLoTELoader{listType: want}
	withLoTELoaders(t, loader)

	got := ListTypeFromURI("urn:known")
	if got != ListType(want) {
		t.Errorf("ListTypeFromURI = %v, want %v", got, want)
	}
}
