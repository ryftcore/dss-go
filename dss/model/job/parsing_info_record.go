// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/ParsingInfoRecord.java (DSS 6.5.RC1).
package job

// ParsingInfoRecord defines a parsing result record.
type ParsingInfoRecord interface {
	InfoRecord

	// StructureValidationMessages gets a list of error messages when occurred during the
	// structure validation: empty when the structure validation succeeded. Port of
	// getStructureValidationMessages().
	StructureValidationMessages() []string
}
