// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/dto/builder/AbstractParsingCacheDTOBuilder.java (DSS 6.5.RC1).
package job

// AbstractParsingCacheDTOBuilderOverrides declares the operations Java's abstract
// AbstractParsingCacheDTOBuilder class leaves abstract/overridable, standing in for the
// virtual dispatch the base needs to reach the concrete builder (outside this module, e.g.
// a TL-specific parsing cache DTO builder). A concrete builder registers itself with
// AbstractParsingCacheDTOBuilder.InitAbstractParsingCacheDTOBuilder.
type AbstractParsingCacheDTOBuilderOverrides interface {
	// Init instantiates the concrete AbstractParsingCacheDTO from the base AbstractCacheDTO.
	// Port of the abstract protected init(AbstractCacheDTO). Returns the ParsingCacheDTO
	// interface, not the concrete *AbstractParsingCacheDTO, so a concrete builder can hand
	// back a narrower wrapper type (Java's covariant return) - see ParsingCacheDTO's doc.
	Init(abstractCacheDTO *AbstractCacheDTO) ParsingCacheDTO

	// BuildParsingCacheDTO fills parsingCacheDTO. Port of the (overridable, non-abstract)
	// protected build(AbstractParsingCacheDTO); a concrete builder wanting the base
	// behaviour first should call
	// AbstractParsingCacheDTOBuilder.BuildParsingCacheDTODefault explicitly, mirroring
	// Java's super.build(dto).
	BuildParsingCacheDTO(parsingCacheDTO ParsingCacheDTO)
}

// AbstractParsingCacheDTOBuilder builds an AbstractParsingCacheDTO.
type AbstractParsingCacheDTOBuilder struct {
	AbstractCacheDTOBuilder[ParsingResult]

	// overrides points back at the concrete builder; see InitAbstractParsingCacheDTOBuilder.
	overrides AbstractParsingCacheDTOBuilderOverrides
}

// NewAbstractParsingCacheDTOBuilder creates an AbstractParsingCacheDTOBuilder for the given
// parsing cache entry. Port of the public constructor.
func NewAbstractParsingCacheDTOBuilder(cachedEntry *CachedEntry[ParsingResult]) AbstractParsingCacheDTOBuilder {
	return AbstractParsingCacheDTOBuilder{AbstractCacheDTOBuilder: NewAbstractCacheDTOBuilder(cachedEntry)}
}

// InitAbstractParsingCacheDTOBuilder registers the concrete builder with its base so that
// the base can dispatch to Init/BuildParsingCacheDTO. It must be called exactly once, by the
// concrete builder's constructor, before any other method.
func (b *AbstractParsingCacheDTOBuilder) InitAbstractParsingCacheDTOBuilder(overrides AbstractParsingCacheDTOBuilderOverrides) {
	b.overrides = overrides
}

func (b *AbstractParsingCacheDTOBuilder) abstractParsingCacheDTOBuilderOverrides() AbstractParsingCacheDTOBuilderOverrides {
	if b.overrides == nil {
		panic("AbstractParsingCacheDTOBuilder was not initialised: the concrete builder must call InitAbstractParsingCacheDTOBuilder in its constructor")
	}
	return b.overrides
}

// Build builds the AbstractParsingCacheDTO. Port of build().
func (b *AbstractParsingCacheDTOBuilder) Build() ParsingCacheDTO {
	parsingCacheDTO := b.abstractParsingCacheDTOBuilderOverrides().Init(b.AbstractCacheDTOBuilder.Build())
	if b.IsResultExist() {
		b.abstractParsingCacheDTOBuilderOverrides().BuildParsingCacheDTO(parsingCacheDTO)
	}
	return parsingCacheDTO
}

// BuildParsingCacheDTODefault is the base (non-abstract) implementation of the
// BuildParsingCacheDTO hook. Port of the protected build(AbstractParsingCacheDTO) default
// body; an overriding concrete builder calls this explicitly to get the base behaviour,
// mirroring Java's super.build(dto).
func (b *AbstractParsingCacheDTOBuilder) BuildParsingCacheDTODefault(parsingCacheDTO ParsingCacheDTO) {
	parsingCacheDTO.SetStructureValidationMessages(b.Result().StructureValidationMessages())
}
