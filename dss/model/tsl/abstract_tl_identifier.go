// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/identifier/AbstractTLIdentifier.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.tsl.identifier is flattened into this tsl package per
// the Phase 1b cycle-driven flattening table.
package tsl

import "github.com/utain/esig/dss/model"

// tlURLProvider is satisfied by TLInfo and LOTLInfo (which embeds TLInfo), giving
// AbstractTLIdentifier access to tlInfo.getUrl() without depending on a concrete type. This
// lets PivotInfo (see pivot_identifier.go, already ported) satisfy it structurally too, since
// PivotInfo embeds LOTLInfo and inherits the promoted Url() method.
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
// DEVIATION from pivot_identifier.go (already ported): that file calls
// NewAbstractTLIdentifier(prefix, tlInfo) with two arguments, assuming AbstractTLIdentifier
// derives its className another way. Java's Identifier#toString()/#equals() rely on
// getClass().getSimpleName(), which Go cannot recover generically from an embedded struct, so
// this constructor takes className explicitly instead; pivot_identifier.go's call site will
// need updating by the integrator to pass "PivotIdentifier" as the first argument.
func NewAbstractTLIdentifier(className, prefix string, tlInfo tlURLProvider) AbstractTLIdentifier {
	return AbstractTLIdentifier{
		MultipleDigestIdentifier: model.NewMultipleDigestIdentifier(className, prefix, []byte(tlInfo.Url())),
	}
}
