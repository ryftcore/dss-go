// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/EAAPresentation.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.eaa.EAAPresentation lands in this same Go package per
// S2B_BRIEF.md's package layout table.
//
// FORWARD DEPENDENCY: EAA (Java spi.eaa.EAA) is flattened into this same package by a sibling
// chunk of phase 2b and is already referenced opaquely elsewhere in this package (see
// eaa_revocation_token.go, advanced_signature.go), so it is used unqualified here too.
//
// NOTE: Java's javadoc on getElectronicAttestationsOfAttributes() says "a list of
// {@link AdvancedSignature}s", but the declared return type is List<EAA> - the javadoc is
// stale/wrong upstream. The Go signature follows the actual declared type (EAA), not the
// javadoc.
package validation

import "github.com/utain/esig/dss/enumerations"

// EAAPresentation represents an EAA Presentation document.
type EAAPresentation interface {
	// EAAPresentationType gets the type of the Electronic Attestation of Attributes
	// presentation. Port of getEAAPresentationType().
	EAAPresentationType() enumerations.EAAPresentationType

	// ElectronicAttestationsOfAttributes gets a list of EAAs used to issue the Electronic
	// Attestation of Attributes. Port of getElectronicAttestationsOfAttributes().
	ElectronicAttestationsOfAttributes() []EAA
}
