// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/UniqueServiceFilter.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: CertificateQualificationCalculator is ported by the
// qualification/certificate porter into this same package (dss/validation/process/qualification),
// assumed to have the shape
//
//	func NewCertificateQualificationCalculator(endEntityCert *diagnostic.CertificateWrapper,
//	    trustService *diagnostic.TrustServiceWrapper) *CertificateQualificationCalculator
//	// with a Qualification() enumerations.CertificateQualification method
//
// mirroring Java's CertificateQualificationCalculator(CertificateWrapper, TrustServiceWrapper)
// and getQualification().
//
// slf4j logging (the "More than one selected trust services" / "Unable to select..." /
// "All trust services conclude..." records) is dropped per PORTING.md.
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
)

// UniqueServiceFilter is used to select a TrustService that is unambiguous and does not
// have conflicts with other TrustServices. In case of a conflict for the given
// endEntityCert, none of the TrustServices is returned.
type UniqueServiceFilter struct {
	// endEntityCert is the certificate to check TrustServices for.
	endEntityCert *diagnostic.CertificateWrapper
}

// NewUniqueServiceFilter is the default constructor. Port of
// UniqueServiceFilter(CertificateWrapper).
func NewUniqueServiceFilter(endEntityCert *diagnostic.CertificateWrapper) *UniqueServiceFilter {
	return &UniqueServiceFilter{endEntityCert: endEntityCert}
}

// Filter filters a list of TrustServiceWrappers. Port of filter(List).
func (f *UniqueServiceFilter) Filter(trustServices []*diagnostic.TrustServiceWrapper) []*diagnostic.TrustServiceWrapper {
	var selectedTrustService *diagnostic.TrustServiceWrapper

	if utils.CollectionSize(trustServices) == 1 {
		selectedTrustService = trustServices[0]
	} else if utils.IsCollectionNotEmpty(trustServices) {

		qualificationResults := make(map[enumerations.CertificateQualification][]string)
		for _, trustService := range trustServices {
			calculator := NewCertificateQualificationCalculator(f.endEntityCert, trustService)
			certQualification := calculator.Qualification()
			if _, ok := qualificationResults[certQualification]; !ok { // putIfAbsent, as trustService.ServiceNames may be nil
				qualificationResults[certQualification] = trustService.ServiceNames
			}
		}

		if len(qualificationResults) == 1 {
			selectedTrustService = trustServices[0]
		}
		// else: several possible conclusions, selectedTrustService stays nil
	}

	if selectedTrustService != nil {
		return []*diagnostic.TrustServiceWrapper{selectedTrustService}
	}
	return []*diagnostic.TrustServiceWrapper{}
}
