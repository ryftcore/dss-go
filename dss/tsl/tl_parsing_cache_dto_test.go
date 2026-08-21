// Tests for tl_parsing_cache_dto.go and parsing_utils.go.
//
// The three counters (getTSPNumber/getTSNumber/getCertNumber) are the only behaviour
// TLParsingCacheDTO carries beyond plain field access, and ParsingUtils is likewise pure
// selection logic, so the expectations below are read directly off the Java bodies rather than off
// a fixture.
package tsl

import (
	"testing"

	"github.com/utain/esig/dss/model"
	tslmodel "github.com/utain/esig/dss/model/tsl"
)

// tlParsingCacheDTOTestProvider builds a TrustServiceProvider carrying certsPerService entries per
// service, one service per element. The certificate slots stay nil: the three counters only read
// len(), never the tokens themselves.
func tlParsingCacheDTOTestProvider(t *testing.T, certsPerService ...int) *tslmodel.TrustServiceProvider {
	t.Helper()
	services := make([]*tslmodel.TrustService, 0, len(certsPerService))
	for _, count := range certsPerService {
		certificates := make([]*model.CertificateToken, count)
		services = append(services, tslmodel.NewTrustService(certificates, nil))
	}
	provider := tslmodel.NewTrustServiceProvider()
	provider.SetServices(services)
	return provider
}

func TestTLParsingCacheDTO_Counters(t *testing.T) {
	dto := NewTLParsingCacheDTO()

	// An unset (null) provider list answers 0 everywhere - the Utils.isCollectionNotEmpty guard.
	if got := dto.TSPNumber(); got != 0 {
		t.Errorf("TSPNumber() = %d, want 0", got)
	}
	if got := dto.TSNumber(); got != 0 {
		t.Errorf("TSNumber() = %d, want 0", got)
	}
	if got := dto.CertNumber(); got != 0 {
		t.Errorf("CertNumber() = %d, want 0", got)
	}

	dto.SetTrustServiceProviders([]*tslmodel.TrustServiceProvider{
		tlParsingCacheDTOTestProvider(t, 1, 2),
		tlParsingCacheDTOTestProvider(t, 3),
		tlParsingCacheDTOTestProvider(t),
	})
	if got := dto.TSPNumber(); got != 3 {
		t.Errorf("TSPNumber() = %d, want 3", got)
	}
	if got := dto.TSNumber(); got != 3 {
		t.Errorf("TSNumber() = %d, want 3", got)
	}
	if got := dto.CertNumber(); got != 6 {
		t.Errorf("CertNumber() = %d, want 6", got)
	}
}

func TestParsingUtils_XMLLOTLPointer(t *testing.T) {
	// A nil DTO, and a DTO whose cached result does not exist, both answer nil.
	if got := ParsingUtilsXMLLOTLPointer(nil); got != nil {
		t.Errorf("nil DTO = %v, want nil", got)
	}
	dto := NewTLParsingCacheDTO()
	dto.SetLotlOtherPointers([]*tslmodel.OtherTSLPointer{tslmodel.NewOtherTSLPointer()})
	if dto.IsResultExist() {
		t.Fatal("a fresh cache DTO must report no result")
	}
	if got := ParsingUtilsXMLLOTLPointer(dto); got != nil {
		t.Errorf("DTO without a result = %v, want nil", got)
	}
}

func TestParsingUtils_LOTLAnnouncedCertificateSource(t *testing.T) {
	pointer := tslmodel.NewOtherTSLPointerBuilder().
		SetSdiCertificates([]*model.CertificateToken{}).Build()
	source := ParsingUtilsLOTLAnnouncedCertificateSource(pointer)
	if source == nil {
		t.Fatal("a certificate source is always returned")
	}
	if got := len(source.Certificates()); got != 0 {
		t.Errorf("Certificates() = %d, want 0", got)
	}
}
