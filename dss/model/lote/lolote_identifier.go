// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/identifier/LoLoTEIdentifier.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.lote.identifier is flattened into this lote package.
package lote

import "github.com/ryftcore/dss-go/dss/model"

// loloteIdentifierPrefix is the "LoLoTE-" prefix LoLoTEIdentifier passes to
// AbstractIdentifier.
const loloteIdentifierPrefix = "LoLoTE-"

// LoLoTEIdentifier is the identifier for a List of Lists of Trusted Entities.
type LoLoTEIdentifier struct {
	AbstractIdentifier
}

// NewLoLoTEIdentifier is the default constructor. Port of the LoLoTEIdentifier(LoTEInfo)
// constructor; Java declares the parameter as LoTEInfo even though the typical caller
// (LoLoTEInfo#buildIdentifier) passes "this", relying on the LoLoTEInfo-extends-LoTEInfo
// upcast. Go has no implicit upcast through embedding, so callers pass the embedded Info
// field explicitly (e.g. &loloteInfo.Info).
func NewLoLoTEIdentifier(listInfo *Info) *LoLoTEIdentifier {
	return &LoLoTEIdentifier{
		AbstractIdentifier: NewAbstractLoTEIdentifier("LoLoTEIdentifier", loloteIdentifierPrefix, listInfo),
	}
}

// compile-time interface assertion.
var _ model.Identifier = (*LoLoTEIdentifier)(nil)
