// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/identifier/LOTLIdentifier.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.tsl.identifier is flattened into this tsl package per
// the Phase 1b cycle-driven flattening table.
package tsl

import "github.com/ryftcore/dss-go/dss/model"

// lotlIdentifierPrefix is the "LOTL-" prefix LOTLIdentifier passes to AbstractTLIdentifier.
const lotlIdentifierPrefix = "LOTL-"

// LOTLIdentifier is the identifier for a LOTL.
type LOTLIdentifier struct {
	AbstractTLIdentifier
}

// NewLOTLIdentifier is the default constructor.
func NewLOTLIdentifier(lotlInfo *LOTLInfo) *LOTLIdentifier {
	return &LOTLIdentifier{
		AbstractTLIdentifier: NewAbstractTLIdentifier("LOTLIdentifier", lotlIdentifierPrefix, lotlInfo),
	}
}

// compile-time interface assertion.
var _ model.Identifier = (*LOTLIdentifier)(nil)
