// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/job/TLSourceBuilder.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (see xml_download_result.go's header): uses job.CacheKey.
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// TLSourceBuilder builds a list of TLSources.
type TLSourceBuilder struct {
	// lotlList is the LOTL sources to build TLSources from.
	lotlList []*LOTLSource

	// parsingResults is the map of LOTL parsing results.
	parsingResults map[job.CacheKey]*TLParsingCacheDTO
}

// NewTLSourceBuilder is the default constructor.
func NewTLSourceBuilder(lotlList []*LOTLSource, parsingResults map[job.CacheKey]*TLParsingCacheDTO) *TLSourceBuilder {
	return &TLSourceBuilder{lotlList: lotlList, parsingResults: parsingResults}
}

// Build builds a list of TLSources. Port of build().
func (b *TLSourceBuilder) Build() []*TLSource {
	var result []*TLSource
	for _, lotlSource := range b.lotlList {
		cachedResult := b.parsingResults[lotlSource.CacheKey()]
		if cachedResult != nil && cachedResult.IsResultExist() {
			for _, otherTSLPointerDTO := range cachedResult.TlOtherPointers() {
				result = append(result, b.getTLSource(otherTSLPointerDTO, lotlSource))
			}
		}
	}
	return result
}

func (b *TLSourceBuilder) getTLSource(otherTSLPointerDTO *tslmodel.OtherTSLPointer, lotlSource *LOTLSource) *TLSource {
	tlSource := NewTLSource()
	tlSource.SetUrl(otherTSLPointerDTO.TSLLocation())
	tlSource.SetCertificateSource(b.getCertificateSource(otherTSLPointerDTO.SdiCertificates()))
	tlSource.SetTrustServiceProviderPredicate(lotlSource.TrustServiceProviderPredicate())
	tlSource.SetTrustServicePredicate(lotlSource.TrustServicePredicate())
	tlSource.SetTLVersions(lotlSource.TLVersions())
	return tlSource
}

func (b *TLSourceBuilder) getCertificateSource(certificates []*model.CertificateToken) spi.CertificateSource {
	certificateSource := spi.NewCommonCertificateSource()
	for _, certificate := range certificates {
		certificateSource.AddCertificate(certificate)
	}
	return &certificateSource
}
