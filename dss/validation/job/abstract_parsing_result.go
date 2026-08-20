// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/parsing/AbstractParsingResult.java (DSS 6.5.RC1).
package job

// AbstractParsingResult is the abstract parsing result. Concrete parsing results (outside
// this module) embed it to satisfy ParsingResult.
type AbstractParsingResult struct {
	// structureValidationMessages is a list of error messages occurred during a structure
	// validation.
	structureValidationMessages []string
}

// NewAbstractParsingResult creates an AbstractParsingResult with nil values. Port of the
// default constructor.
func NewAbstractParsingResult() *AbstractParsingResult {
	return &AbstractParsingResult{}
}

// StructureValidationMessages returns the structure validation error messages. Port of
// getStructureValidationMessages().
func (a *AbstractParsingResult) StructureValidationMessages() []string {
	return a.structureValidationMessages
}

// SetStructureValidationMessages sets the structure validation error messages. Port of
// setStructureValidationMessages(List).
func (a *AbstractParsingResult) SetStructureValidationMessages(structureValidationMessages []string) {
	a.structureValidationMessages = structureValidationMessages
}
