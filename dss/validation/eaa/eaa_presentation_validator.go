// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/eaa/EAAPresentationValidator.java (DSS 6.5.RC1).
package eaa

import (
	"github.com/utain/esig/dss/spi/eaa/status"
	spivalidation "github.com/utain/esig/dss/spi/validation"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// EAAPresentationValidator is used to validate an Electronic Attestation of Attributes
// presentation.
type EAAPresentationValidator interface {
	dssvalidation.DocumentValidator

	// EAAPresentation gets the EAAPresentation created from the provided document on
	// validation. Port of getEAAPresentation().
	EAAPresentation() spivalidation.EAAPresentation

	// SetEAARevocationSource sets the EAA revocation source providing access to the
	// information about the EAA validity status. Port of setEAARevocationSource(EAARevocationSource).
	SetEAARevocationSource(eaaRevocationSource status.EAARevocationSource)

	// SetEAAValidationParameters sets supplementary data requiring for validation of EAA
	// presentation (format specific). Port of setEAAValidationParameters(EAAValidationParameters).
	SetEAAValidationParameters(eaaValidationParameters spivalidation.EAAValidationParameters)
}
