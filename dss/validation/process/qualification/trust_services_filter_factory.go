// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/TrustServicesFilterFactory.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
)

// TrustServicesFilterFactoryCreateFilterByGranted creates a TrustService filter by
// 'granted' status. Port of createFilterByGranted().
func TrustServicesFilterFactoryCreateFilterByGranted() TrustServiceFilter {
	return NewGrantedServiceFilter()
}

// TrustServicesFilterFactoryCreateFilterByCaQc creates a TrustService filter by 'CA/QC'
// identifier. Port of createFilterByCaQc().
func TrustServicesFilterFactoryCreateFilterByCaQc() TrustServiceFilter {
	return NewCaQcServiceFilter()
}

// TrustServicesFilterFactoryCreateFilterByQTST creates a TrustService filter by
// 'TSA/QTST' identifier. Port of createFilterByQTST().
func TrustServicesFilterFactoryCreateFilterByQTST() TrustServiceFilter {
	return NewQTSTServiceFilter()
}

// TrustServicesFilterFactoryCreateFilterByQEAA creates a TrustService filter by 'EAA/Q'
// identifier. Port of createFilterByQEAA().
func TrustServicesFilterFactoryCreateFilterByQEAA() TrustServiceFilter {
	return NewQEAAServiceFilter()
}

// TrustServicesFilterFactoryCreateFilterByDate creates a TrustService filter by date.
// Port of createFilterByDate(Date).
func TrustServicesFilterFactoryCreateFilterByDate(date *time.Time) TrustServiceFilter {
	return NewServiceByDateFilter(date)
}

// TrustServicesFilterFactoryCreateFilterByCountry creates a TrustService filter by
// country code. Port of createFilterByCountry(String).
func TrustServicesFilterFactoryCreateFilterByCountry(countryCode string) TrustServiceFilter {
	return NewServiceByCountryFilter(countryCode)
}

// TrustServicesFilterFactoryCreateFilterByCountries creates a TrustService filter by
// country codes. Port of createFilterByCountries(Set).
func TrustServicesFilterFactoryCreateFilterByCountries(countryCodes map[string]struct{}) TrustServiceFilter {
	return NewServiceByCountryFilterWithCodes(countryCodes)
}

// TrustServicesFilterFactoryCreateFilterByUrls creates a TrustService filter by urls.
// Port of createFilterByUrls(Set).
func TrustServicesFilterFactoryCreateFilterByUrls(urls map[string]struct{}) TrustServiceFilter {
	return NewServiceByTLUrlFilterWithUrls(urls)
}

// TrustServicesFilterFactoryCreateUniqueServiceFilter creates a TrustService filter by
// end-entity certificate. Port of createUniqueServiceFilter(CertificateWrapper).
func TrustServicesFilterFactoryCreateUniqueServiceFilter(endEntityCertificate *diagnostic.CertificateWrapper) TrustServiceFilter {
	return NewUniqueServiceFilter(endEntityCertificate)
}

// TrustServicesFilterFactoryCreateFilterByCertificateType creates a TrustService filter by
// the type as in the given certificate. Port of
// createFilterByCertificateType(CertificateWrapper).
func TrustServicesFilterFactoryCreateFilterByCertificateType(certificate *diagnostic.CertificateWrapper) TrustServiceFilter {
	return NewServiceByCertificateTypeFilter(certificate)
}

// TrustServicesFilterFactoryCreateConsistentServiceByStatusFilter creates a TrustService
// filter by status consistency. Port of createConsistentServiceByStatusFilter().
func TrustServicesFilterFactoryCreateConsistentServiceByStatusFilter() TrustServiceFilter {
	return NewConsistentServiceByStatusFilter()
}

// TrustServicesFilterFactoryCreateConsistentServiceByQCFilter creates a TrustService
// filter by QC consistency. Port of createConsistentServiceByQCFilter().
func TrustServicesFilterFactoryCreateConsistentServiceByQCFilter() TrustServiceFilter {
	return NewConsistentServiceByQCFilter()
}

// TrustServicesFilterFactoryCreateConsistentServiceByCertificateTypeFilter creates a
// TrustService filter by QC consistency. Port of
// createConsistentServiceByCertificateTypeFilter().
func TrustServicesFilterFactoryCreateConsistentServiceByCertificateTypeFilter() TrustServiceFilter {
	return NewConsistentServiceByCertificateTypeFilter()
}

// TrustServicesFilterFactoryCreateConsistentServiceByQSCDFilter creates a TrustService
// filter by QSCD consistency. Port of createConsistentServiceByQSCDFilter().
func TrustServicesFilterFactoryCreateConsistentServiceByQSCDFilter() TrustServiceFilter {
	return NewConsistentServiceByQSCDFilter()
}

// TrustServicesFilterFactoryCreateMRAEnactedFilter creates a TrustService filter by MRA
// enacted. Port of createMRAEnactedFilter().
func TrustServicesFilterFactoryCreateMRAEnactedFilter() TrustServiceFilter {
	return NewServiceByMRAEnactedFilter()
}

// TrustServicesFilterFactoryCreateFilterByMRAEquivalenceStartingDate creates a
// TrustService filter by MRA equivalence starting date. Port of
// createFilterByMRAEquivalenceStartingDate(Date).
func TrustServicesFilterFactoryCreateFilterByMRAEquivalenceStartingDate(date *time.Time) TrustServiceFilter {
	return NewServiceByMRAEquivalenceStartingDateFilter(date)
}
