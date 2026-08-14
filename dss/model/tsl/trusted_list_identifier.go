// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/identifier/TrustedListIdentifier.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.tsl.identifier is flattened into this tsl package per
// the Phase 1b cycle-driven flattening table.
package tsl

import "github.com/utain/esig/dss/model"

// trustedListIdentifierPrefix is the "TL-" prefix TrustedListIdentifier passes to
// AbstractTLIdentifier.
const trustedListIdentifierPrefix = "TL-"

// TrustedListIdentifier is the identifier for a TL.
type TrustedListIdentifier struct {
	AbstractTLIdentifier
}

// NewTrustedListIdentifier is the default constructor.
func NewTrustedListIdentifier(tlInfo *TLInfo) *TrustedListIdentifier {
	return &TrustedListIdentifier{
		AbstractTLIdentifier: NewAbstractTLIdentifier("TrustedListIdentifier", trustedListIdentifierPrefix, tlInfo),
	}
}

// compile-time interface assertion.
var _ model.Identifier = (*TrustedListIdentifier)(nil)
