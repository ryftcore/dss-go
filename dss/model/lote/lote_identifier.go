// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/identifier/LoTEIdentifier.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.lote.identifier is flattened into this lote package.
package lote

import "github.com/ryftcore/dss-go/dss/model"

// loteIdentifierPrefix is the "LoTE-" prefix Identifier passes to AbstractLoTEIdentifier.
const loteIdentifierPrefix = "LoTE-"

// Identifier is the identifier for a List of Trusted Entities.
type Identifier struct {
	AbstractLoTEIdentifier
}

// NewLoTEIdentifier is the default constructor.
func NewLoTEIdentifier(listInfo *Info) *Identifier {
	return &Identifier{
		AbstractLoTEIdentifier: NewAbstractLoTEIdentifier("LoTEIdentifier", loteIdentifierPrefix, listInfo),
	}
}

// compile-time interface assertion.
var _ model.Identifier = (*Identifier)(nil)
