package model

import (
	"sort"
	"testing"
)

func TestTokenComparatorOrdersByDSSIDString(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	leaf := certificateTokenFixture(t, leafCertificateBase64)
	comparator := NewTokenComparator()

	// "C-2E77..." sorts before "C-E46C...".
	if comparator.Compare(leaf, root) >= 0 {
		t.Errorf("expected %q < %q", leaf.DSSIDAsString(), root.DSSIDAsString())
	}
	if comparator.Compare(root, leaf) <= 0 {
		t.Error("the comparison must be antisymmetric")
	}
	if comparator.Compare(root, root) != 0 {
		t.Error("a token must compare equal to itself")
	}

	tokens := []Token{root, leaf}
	sort.Slice(tokens, func(i, j int) bool { return comparator.Less(tokens[i], tokens[j]) })
	if tokens[0] != Token(leaf) {
		t.Errorf("sorted order = %q, %q", tokens[0].DSSIDAsString(), tokens[1].DSSIDAsString())
	}
}
