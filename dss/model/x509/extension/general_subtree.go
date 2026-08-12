// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/GeneralSubtree.java (DSS 6.5.RC1).
package extension

import "math/big"

// GeneralSubtree represents a general subtree element (see "4.2.1.10. Name Constraints" of
// RFC 5280).
type GeneralSubtree struct {
	GeneralName

	// minimum MUST be 0.
	minimum *big.Int

	// maximum MUST be absent.
	maximum *big.Int
}

// NewGeneralSubtree instantiates the object with null values. Ports the default
// constructor.
func NewGeneralSubtree() *GeneralSubtree {
	return &GeneralSubtree{}
}

// Minimum gets the minimum constraint value.
func (g *GeneralSubtree) Minimum() *big.Int {
	return g.minimum
}

// SetMinimum sets the minimum constraint value.
func (g *GeneralSubtree) SetMinimum(minimum *big.Int) {
	g.minimum = minimum
}

// Maximum gets the maximum constraint value.
func (g *GeneralSubtree) Maximum() *big.Int {
	return g.maximum
}

// SetMaximum sets the maximum constraint value.
func (g *GeneralSubtree) SetMaximum(maximum *big.Int) {
	g.maximum = maximum
}
