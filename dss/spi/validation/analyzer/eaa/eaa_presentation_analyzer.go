// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/eaa/EAAPresentationAnalyzer.java (DSS 6.5.RC1).
package eaa

import (
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

// EAAPresentationAnalyzer performs validation of a presentation of Electronic Attestation of
// Attributes (EAA).
type EAAPresentationAnalyzer interface {
	analyzer.DocumentAnalyzer

	// EAAPresentation gets the extracted Electronic Attestation of Attributes (EAA)
	// presentation. Port of getEAAPresentation().
	EAAPresentation() validation.EAAPresentation
}
