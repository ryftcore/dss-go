// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/identifier/AbstractTLIdentifier.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.tsl.identifier is flattened into this tsl package.
package tsl

import "github.com/ryftcore/dss-go/dss/model"

// tlURLProvider is satisfied by TLInfo and LOTLInfo (which embeds TLInfo), giving
// AbstractTLIdentifier access to tlInfo.getUrl() without depending on a concrete type. This
// lets PivotInfo (see pivot_identifier.go) satisfy it structurally too, since PivotInfo embeds
// LOTLInfo and inherits the promoted Url() method.
type tlURLProvider interface {
	Url() string
}

// AbstractTLIdentifier is the abstract base for DSS internal TL/LOTL/pivot identifiers,
// designed for embedding.
type AbstractTLIdentifier struct {
	model.MultipleDigestIdentifier
}

// NewAbstractTLIdentifier builds the identifier over the SHA-256 digest of the TL/LOTL URL
// bytes with the given prefix (e.g. "TL-", "LOTL-", "P-"). className is the Java simple class
// name of the concrete identifier subclass (e.g. "TrustedListIdentifier"), which
// IdentifierBase needs for its Equals/String ports.
//
// DEVIATION: Java's Identifier#toString()/#equals() rely on getClass().getSimpleName(), which
// Go cannot recover generically from an embedded struct, so this constructor takes className
// explicitly. All three concrete identifiers (TrustedListIdentifier, LOTLIdentifier,
// PivotIdentifier) pass their own Java simple class name.
func NewAbstractTLIdentifier(className, prefix string, tlInfo tlURLProvider) AbstractTLIdentifier {
	return AbstractTLIdentifier{
		MultipleDigestIdentifier: model.NewMultipleDigestIdentifier(className, prefix, []byte(tlInfo.Url())),
	}
}
