// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/dto/builder/TLParsingCacheDTOBuilder.java (DSS 6.5.RC1).
//
// Extends job.AbstractParsingCacheDTOBuilder (Java
// eu.europa.esig.dss.validation.job.dto.builder.AbstractParsingCacheDTOBuilder). The members it
// relies on, all named after their Java counterparts, are:
//
//	job.CachedEntry[job.ParsingResult]                          - the cached entry the constructor takes
//	job.NewAbstractParsingCacheDTOBuilder(*job.CachedEntry[...]) - the super(cachedEntry) constructor
//	job.AbstractParsingCacheDTOBuilder.InitAbstractParsingCacheDTOBuilder(overrides)
//	                                                             - the Overrides registration this
//	                                                               codebase uses in place of Java
//	                                                               virtual dispatch
//	job.AbstractParsingCacheDTOBuilder.Build() job.ParsingCacheDTO      - super.build()
//	job.AbstractParsingCacheDTOBuilder.BuildParsingCacheDTO(job.ParsingCacheDTO)
//	                                                             - super.build(AbstractParsingCacheDTO)
//	job.AbstractParsingCacheDTOBuilder.Result() job.ParsingResult - the protected final getResult()
//	job.ParsingCacheDTO                                          - the polymorphic view of
//	                                                               AbstractParsingCacheDTO that
//	                                                               init()/build() exchange (Go
//	                                                               cannot pass an abstract struct
//	                                                               polymorphically)
//
// If that chunk chose different spellings, the fix here is mechanical.
package tsl

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// tlParsingResultAccessor is the Go form of `result instanceof AbstractTLParsingResult`:
// AbstractTLParsingResult is an embedded struct, not an interface, so the type test is expressed
// as an assertion against the set of accessors it contributes. Both *TLParsingResult and
// *LOTLParsingResult satisfy it through promotion, exactly as both extend AbstractTLParsingResult
// in Java.
type tlParsingResultAccessor interface {
	TSLType() enumerations.TSLType
	SequenceNumber() *int
	Version() *int
	Territory() string
	IssueDate() time.Time
	NextUpdateDate() time.Time
	DistributionPoints() []string
}

// TLParsingCacheDTOBuilder builds a parsing cache DTO for a Trusted List processing.
type TLParsingCacheDTOBuilder struct {
	job.AbstractParsingCacheDTOBuilder
}

// NewTLParsingCacheDTOBuilder is the default constructor. Port of
// TLParsingCacheDTOBuilder(CachedEntry<ParsingResult>).
func NewTLParsingCacheDTOBuilder(cachedEntry *job.CachedEntry[job.ParsingResult]) *TLParsingCacheDTOBuilder {
	builder := &TLParsingCacheDTOBuilder{
		AbstractParsingCacheDTOBuilder: job.NewAbstractParsingCacheDTOBuilder(cachedEntry),
	}
	builder.InitAbstractParsingCacheDTOBuilder(builder)
	return builder
}

// Build builds the DTO. Port of the covariant build() override, whose Java body is a downcast of
// super.build() - which, thanks to the init() override below, always produced a
// TLParsingCacheDTO.
func (b *TLParsingCacheDTOBuilder) Build() *TLParsingCacheDTO {
	return b.AbstractParsingCacheDTOBuilder.Build().(*TLParsingCacheDTO)
}

// Init instantiates a new TLParsingCacheDTO from an AbstractCacheDTO. Port of the protected
// init(AbstractCacheDTO) override.
func (b *TLParsingCacheDTOBuilder) Init(abstractCacheDTO *job.AbstractCacheDTO) job.ParsingCacheDTO {
	return NewTLParsingCacheDTOFromCacheDTO(abstractCacheDTO)
}

// BuildParsingCacheDTO fills the parsing cache DTO. Port of the protected
// build(AbstractParsingCacheDTO) override, which first delegates to super.build(...).
func (b *TLParsingCacheDTOBuilder) BuildParsingCacheDTO(parsingCacheDTO job.ParsingCacheDTO) {
	b.AbstractParsingCacheDTOBuilder.BuildParsingCacheDTODefault(parsingCacheDTO)

	tlParsingCacheDTO := parsingCacheDTO.(*TLParsingCacheDTO)
	tlParsingCacheDTO.SetTSLType(b.tslType())
	tlParsingCacheDTO.SetSequenceNumber(b.sequenceNumber())
	tlParsingCacheDTO.SetVersion(b.version())
	tlParsingCacheDTO.SetTerritory(b.territory())
	tlParsingCacheDTO.SetIssueDate(b.issueDate())
	tlParsingCacheDTO.SetNextUpdateDate(b.nextUpdateDate())
	tlParsingCacheDTO.SetDistributionPoints(b.distributionPoints())
	if b.isLOTL() {
		tlParsingCacheDTO.SetLotlOtherPointers(b.lotlOtherPointers())
		tlParsingCacheDTO.SetTlOtherPointers(b.tlOtherPointers())
		tlParsingCacheDTO.SetPivotUrls(b.pivotUrls())
		tlParsingCacheDTO.SetSigningCertificateAnnouncementUrl(b.signingCertificateAnnouncementUrl())
	} else {
		tlParsingCacheDTO.SetTrustServiceProviders(b.trustServiceProviders())
	}
}

// isLOTL ports the private isLOTL().
func (b *TLParsingCacheDTOBuilder) isLOTL() bool {
	_, ok := b.Result().(*LOTLParsingResult)
	return ok
}

// The private extractors below all follow the same upstream shape: type-test the cached result,
// answer the corresponding field when it matches, and answer a null/empty default otherwise. The
// slf4j debug lines upstream emits on the mismatch branch are dropped, per PORTING.md.

func (b *TLParsingCacheDTOBuilder) tslType() enumerations.TSLType {
	if result, ok := b.Result().(tlParsingResultAccessor); ok {
		return result.TSLType()
	}
	return nil
}

func (b *TLParsingCacheDTOBuilder) sequenceNumber() *int {
	if result, ok := b.Result().(tlParsingResultAccessor); ok {
		return result.SequenceNumber()
	}
	return nil
}

func (b *TLParsingCacheDTOBuilder) version() *int {
	if result, ok := b.Result().(tlParsingResultAccessor); ok {
		return result.Version()
	}
	return nil
}

func (b *TLParsingCacheDTOBuilder) territory() string {
	if result, ok := b.Result().(tlParsingResultAccessor); ok {
		return result.Territory()
	}
	return ""
}

func (b *TLParsingCacheDTOBuilder) issueDate() time.Time {
	if result, ok := b.Result().(tlParsingResultAccessor); ok {
		return result.IssueDate()
	}
	return time.Time{}
}

func (b *TLParsingCacheDTOBuilder) nextUpdateDate() time.Time {
	if result, ok := b.Result().(tlParsingResultAccessor); ok {
		return result.NextUpdateDate()
	}
	return time.Time{}
}

func (b *TLParsingCacheDTOBuilder) distributionPoints() []string {
	if result, ok := b.Result().(tlParsingResultAccessor); ok {
		return result.DistributionPoints()
	}
	// NOTE: upstream returns null here (not the empty list its log line announces).
	return nil
}

func (b *TLParsingCacheDTOBuilder) trustServiceProviders() []*tslmodel.TrustServiceProvider {
	if result, ok := b.Result().(*TLParsingResult); ok {
		return result.TrustServiceProviders()
	}
	return []*tslmodel.TrustServiceProvider{}
}

func (b *TLParsingCacheDTOBuilder) lotlOtherPointers() []*tslmodel.OtherTSLPointer {
	if result, ok := b.Result().(*LOTLParsingResult); ok {
		return result.LotlPointers()
	}
	return []*tslmodel.OtherTSLPointer{}
}

func (b *TLParsingCacheDTOBuilder) tlOtherPointers() []*tslmodel.OtherTSLPointer {
	if result, ok := b.Result().(*LOTLParsingResult); ok {
		return result.TlPointers()
	}
	return []*tslmodel.OtherTSLPointer{}
}

func (b *TLParsingCacheDTOBuilder) pivotUrls() []string {
	if result, ok := b.Result().(*LOTLParsingResult); ok {
		return result.PivotURLs()
	}
	return []string{}
}

func (b *TLParsingCacheDTOBuilder) signingCertificateAnnouncementUrl() string {
	if result, ok := b.Result().(*LOTLParsingResult); ok {
		return result.SigningCertificateAnnouncementURL()
	}
	return ""
}
