// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/dto/AbstractParsingCacheDTO.java (DSS 6.5.RC1).
package job

import modeljob "github.com/ryftcore/dss-go/dss/model/job"

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

// ParsingCacheDTO is the polymorphic view of AbstractParsingCacheDTO that
// AbstractParsingCacheDTOBuilder exchanges with a concrete builder's Init/BuildParsingCacheDTO
// overrides, standing in for Java's covariant return types: Java's
// AbstractParsingCacheDTOBuilder<R>#init/#build are declared to return/accept
// AbstractParsingCacheDTO but a concrete subclass builder (e.g. dss/tsl's
// TLParsingCacheDTOBuilder) actually constructs/consumes a narrower DTO type (e.g.
// TLParsingCacheDTO) that embeds *AbstractParsingCacheDTO by pointer. Go has no covariant
// interface-method return types, so the Overrides interface and Build() are typed against this
// interface instead of the concrete *AbstractParsingCacheDTO; *AbstractParsingCacheDTO itself,
// and any wrapper embedding it by pointer, satisfy it through promotion.
type ParsingCacheDTO interface {
	modeljob.ParsingInfoRecord

	// SetStructureValidationMessages sets the structure validation error messages. Port of
	// setStructureValidationMessages(List).
	SetStructureValidationMessages(structureValidationMessages []string)
}

var _ ParsingCacheDTO = (*AbstractParsingCacheDTO)(nil)
