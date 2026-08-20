// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/parsing/ParsingResult.java (DSS 6.5.RC1).
package job

// ParsingResult provides an interface to extract information about a parsing task result.
type ParsingResult interface {
	CachedResult

	// StructureValidationMessages returns a list of error messages when occurred during the
	// structure validation: an empty list if the structure validation succeeded. Port of
	// getStructureValidationMessages().
	StructureValidationMessages() []string
}
