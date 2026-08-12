// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/identifier/AbstractLoTEIdentifier.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.lote.identifier is flattened into this lote package
// per the Phase 1b cycle-driven flattening table.
package lote

import "github.com/utain/esig/dss/model"

// AbstractLoTEIdentifier is the abstract base for DSS internal LoTE/LoLoTE identifiers,
// designed for embedding.
type AbstractLoTEIdentifier struct {
	model.MultipleDigestIdentifier
}

// NewAbstractLoTEIdentifier builds the identifier over the SHA-256 digest of the target
// List's URL bytes with the given prefix (e.g. "LoTE-", "LoLoTE-"). className is the Java
// simple class name of the concrete identifier subclass (e.g. "LoTEIdentifier",
// "LoLoTEIdentifier"), which IdentifierBase needs for its Equals/String ports.
//
// Port of the protected AbstractLoTEIdentifier(String, LoTEInfo) constructor; listInfo is typed
// *LoTEInfo exactly as in Java (LoLoTEInfo embeds LoTEInfo, so callers pass its embedded field,
// e.g. &loloteInfo.LoTEInfo, mirroring Java's implicit upcast of "this").
func NewAbstractLoTEIdentifier(className, prefix string, listInfo *LoTEInfo) AbstractLoTEIdentifier {
	return AbstractLoTEIdentifier{
		MultipleDigestIdentifier: model.NewMultipleDigestIdentifier(className, prefix, []byte(listInfo.Url())),
	}
}
