// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/identifier/LoTEIdentifier.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.lote.identifier is flattened into this lote package
// per the Phase 1b cycle-driven flattening table.
package lote

import "github.com/ryftcore/dss-go/dss/model"

// loteIdentifierPrefix is the "LoTE-" prefix LoTEIdentifier passes to AbstractLoTEIdentifier.
const loteIdentifierPrefix = "LoTE-"

// LoTEIdentifier is the identifier for a List of Trusted Entities.
type LoTEIdentifier struct {
	AbstractLoTEIdentifier
}

// NewLoTEIdentifier is the default constructor.
func NewLoTEIdentifier(listInfo *LoTEInfo) *LoTEIdentifier {
	return &LoTEIdentifier{
		AbstractLoTEIdentifier: NewAbstractLoTEIdentifier("LoTEIdentifier", loteIdentifierPrefix, listInfo),
	}
}

// compile-time interface assertion.
var _ model.Identifier = (*LoTEIdentifier)(nil)
