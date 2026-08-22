// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/EAAPresentation.java (DSS 6.5.RC1).
//
// Java spi.eaa.EAAPresentation lands in this same Go package.
//
// NOTE: Java's javadoc on getElectronicAttestationsOfAttributes() says "a list of
// {@link AdvancedSignature}s", but the declared return type is List<EAA> - the javadoc is
// stale/wrong upstream. The Go signature follows the actual declared type (EAA), not the
// javadoc.
package validation

import "github.com/ryftcore/dss-go/dss/enumerations"

// EAAPresentation represents an EAA Presentation document.
type EAAPresentation interface {
	// EAAPresentationType gets the type of the Electronic Attestation of Attributes
	// presentation. Port of getEAAPresentationType().
	EAAPresentationType() enumerations.EAAPresentationType

	// ElectronicAttestationsOfAttributes gets a list of EAAs used to issue the Electronic
	// Attestation of Attributes. Port of getElectronicAttestationsOfAttributes().
	ElectronicAttestationsOfAttributes() []EAA
}
