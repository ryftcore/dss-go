// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/identifier/PivotIdentifier.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.tsl.identifier is flattened into this tsl package.
package tsl

// pivotIdentifierPrefix is the "P-" prefix PivotIdentifier passes to AbstractTLIdentifier.
const pivotIdentifierPrefix = "P-"

// PivotIdentifier is the identifier for a Pivot.
type PivotIdentifier struct {
	AbstractTLIdentifier
}

// NewPivotIdentifier is the default constructor.
func NewPivotIdentifier(pivotInfo *PivotInfo) *PivotIdentifier {
	return &PivotIdentifier{
		AbstractTLIdentifier: NewAbstractTLIdentifier("PivotIdentifier", pivotIdentifierPrefix, pivotInfo),
	}
}
