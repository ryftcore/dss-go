// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/EAAValidationParameters.java (DSS 6.5.RC1).
package validation

// EAAValidationParameters contains supplementary data required for validation of an
// EAA Presentation.
//
// Java declares this as an empty marker interface (`extends Serializable`, with
// Serializable itself dropped per PORTING.md); it carries no methods of its own, so
// implementations satisfy it structurally without needing to reference this type at all.
// It is kept as a named type for signature fidelity with call sites elsewhere that accept
// an EAAValidationParameters.
type EAAValidationParameters any
