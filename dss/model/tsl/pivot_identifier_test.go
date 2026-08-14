package tsl

import "testing"

func TestPivotIdentifierPrefix(t *testing.T) {
	// Java: new AbstractTLIdentifier("P-", pivotInfo) - the "P-" prefix is what makes a
	// PivotIdentifier's AsXmlID() and String() distinguishable from other TL identifiers.
	if pivotIdentifierPrefix != "P-" {
		t.Fatalf("unexpected pivotIdentifierPrefix: %s", pivotIdentifierPrefix)
	}
}
