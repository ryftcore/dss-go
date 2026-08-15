// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/TrustedEntitiesFilterFactory.java (DSS 6.5.RC1).
package qualification

import "time"

// TrustedEntitiesFilterFactoryCreateFilterByListUrl creates a TrustedEntityService filter
// by list Url. Port of createFilterByListUrl(String).
func TrustedEntitiesFilterFactoryCreateFilterByListUrl(url string) TrustedEntityServiceFilter {
	return NewServiceByTrustedEntityServiceUrlFilter(url)
}

// TrustedEntitiesFilterFactoryCreateFilterByListUrls creates a TrustedEntityService filter
// by list urls. Port of createFilterByListUrls(Collection).
func TrustedEntitiesFilterFactoryCreateFilterByListUrls(urls []string) TrustedEntityServiceFilter {
	return NewServiceByTrustedEntityServiceUrlFilterWithUrls(urls)
}

// TrustedEntitiesFilterFactoryCreateFilterByDate creates a TrustedEntityService filter by
// date. Port of createFilterByDate(Date).
func TrustedEntitiesFilterFactoryCreateFilterByDate(date *time.Time) TrustedEntityServiceFilter {
	return NewTrustedEntityServiceByDateFilter(date)
}

// TrustedEntitiesFilterFactoryCreateFilterByServiceTypeIdentifierUri creates a
// TrustedEntityService filter by STI URI. Port of
// createFilterByServiceTypeIdentifierUri(String).
func TrustedEntitiesFilterFactoryCreateFilterByServiceTypeIdentifierUri(stiUri string) TrustedEntityServiceFilter {
	return NewTrustedEntityServiceByStiFilter(stiUri)
}

// TrustedEntitiesFilterFactoryCreateFilterByServiceStatusUri creates a
// TrustedEntityService filter by service status URI. Port of
// createFilterByServiceStatusUri(String).
func TrustedEntitiesFilterFactoryCreateFilterByServiceStatusUri(statusUri string) TrustedEntityServiceFilter {
	return NewTrustedEntityServiceByStatusFilter(statusUri)
}
