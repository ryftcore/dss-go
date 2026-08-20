// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/dto/AbstractParsingCacheDTO.java (DSS 6.5.RC1).
package job

import modeljob "github.com/utain/esig/dss/model/job"

// AbstractParsingCacheDTO is the parsing record DTO. It implements modeljob.ParsingInfoRecord.
type AbstractParsingCacheDTO struct {
	*AbstractCacheDTO

	// structureValidationMessages is a list of error messages occurred during a structure
	// validation.
	structureValidationMessages []string
}

// NewAbstractParsingCacheDTO creates an empty AbstractParsingCacheDTO. Port of the
// protected default constructor.
func NewAbstractParsingCacheDTO() *AbstractParsingCacheDTO {
	return &AbstractParsingCacheDTO{AbstractCacheDTO: NewAbstractCacheDTO()}
}

// NewAbstractParsingCacheDTOFrom copies cacheDTO. Port of the copy constructor.
func NewAbstractParsingCacheDTOFrom(cacheDTO *AbstractCacheDTO) *AbstractParsingCacheDTO {
	return &AbstractParsingCacheDTO{AbstractCacheDTO: NewAbstractCacheDTOFrom(cacheDTO)}
}

// StructureValidationMessages returns the structure validation error messages. Port of
// getStructureValidationMessages().
func (d *AbstractParsingCacheDTO) StructureValidationMessages() []string {
	return d.structureValidationMessages
}

// SetStructureValidationMessages sets the structure validation error messages. Port of
// setStructureValidationMessages(List).
func (d *AbstractParsingCacheDTO) SetStructureValidationMessages(structureValidationMessages []string) {
	d.structureValidationMessages = structureValidationMessages
}

var _ modeljob.ParsingInfoRecord = (*AbstractParsingCacheDTO)(nil)
