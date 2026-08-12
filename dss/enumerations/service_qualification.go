// Ported from dss-enumerations/.../ServiceQualification.java (DSS 6.5.RC1).
package enumerations

// ServiceQualification contains qualification statuses for TrustServices.
type ServiceQualification string

const (
	// ServiceQualification_QC_STATEMENT
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCStatement"):
	// to indicate that all certificates identified by the applicable list
	// of criteria are issued as qualified certificates.
	ServiceQualification_QC_STATEMENT ServiceQualification = "QC_STATEMENT"
	// ServiceQualification_NOT_QUALIFIED
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotQualified"):
	// to indicate that no certificates identified by the applicable list of
	// criteria are to be considered as qualified certificates.
	ServiceQualification_NOT_QUALIFIED ServiceQualification = "NOT_QUALIFIED"
	// ServiceQualification_QC_WITH_SSCD
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithSSCD"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, have
	// their private key residing in an SSCD.
	ServiceQualification_QC_WITH_SSCD ServiceQualification = "QC_WITH_SSCD"
	// ServiceQualification_QC_WITH_QSCD
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, have
	// their private key residing in a QSCD.
	ServiceQualification_QC_WITH_QSCD ServiceQualification = "QC_WITH_QSCD"
	// ServiceQualification_QC_NO_SSCD
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoSSCD"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, do not
	// have their private key residing in an SSCD.
	ServiceQualification_QC_NO_SSCD ServiceQualification = "QC_NO_SSCD"
	// ServiceQualification_QC_NO_QSCD
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoQSCD"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, do not
	// have their private key residing in a QSCD.
	ServiceQualification_QC_NO_QSCD ServiceQualification = "QC_NO_QSCD"
	// ServiceQualification_QC_SSCD_STATUS_AS_IN_CERT
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCSSCDStatusAsInCert"):
	// to indicate that all certificates identified by the applicable list
	// of criteria, when they are claimed or stated as being qualified, do
	// contain proper machine processable information about whether or not
	// their private key residing in an SSCD.
	ServiceQualification_QC_SSCD_STATUS_AS_IN_CERT ServiceQualification = "QC_SSCD_STATUS_AS_IN_CERT"
	// ServiceQualification_QC_QSCD_STATUS_AS_IN_CERT
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDStatusAsInCert"):
	// to indicate that all certificates identified by the applicable list
	// of criteria, when they are claimed or stated as being qualified, do
	// contain proper machine processable information about whether or not
	// their private key residing in a QSCD.
	ServiceQualification_QC_QSCD_STATUS_AS_IN_CERT ServiceQualification = "QC_QSCD_STATUS_AS_IN_CERT"
	// ServiceQualification_QC_QSCD_MANAGED_ON_BEHALF
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDManagedOnBehalf"):
	// to indicate that all certificates identified by the applicable list
	// of criteria, when they are claimed or stated as being qualified, have
	// their private key residing in a QSCD for which the generation and
	// management of that private key is done by the qualified TSP on
	// behalf of the entity whose identity is certified in the certificate.
	ServiceQualification_QC_QSCD_MANAGED_ON_BEHALF ServiceQualification = "QC_QSCD_MANAGED_ON_BEHALF"
	// ServiceQualification_QC_FOR_LEGAL_PERSON
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForLegalPerson"):
	// to indicate that all certificates identified by the applicable list
	// of criteria, when they are claimed or stated as being qualified, are
	// issued to legal persons.
	ServiceQualification_QC_FOR_LEGAL_PERSON ServiceQualification = "QC_FOR_LEGAL_PERSON"
	// ServiceQualification_QC_FOR_ESIG
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESig"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, are
	// issued for electronic signatures.
	ServiceQualification_QC_FOR_ESIG ServiceQualification = "QC_FOR_ESIG"
	// ServiceQualification_QC_FOR_ESEAL
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESeal"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, are
	// issued for electronic seals.
	ServiceQualification_QC_FOR_ESEAL ServiceQualification = "QC_FOR_ESEAL"
	// ServiceQualification_QC_FOR_WSA
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForWSA"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, are
	// issued for web site authentication.
	ServiceQualification_QC_FOR_WSA ServiceQualification = "QC_FOR_WSA"
)

// serviceQualificationURIs holds the URI for each constant.
var serviceQualificationURIs = map[ServiceQualification]string{
	ServiceQualification_QC_STATEMENT:              "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCStatement",
	ServiceQualification_NOT_QUALIFIED:             "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotQualified",
	ServiceQualification_QC_WITH_SSCD:              "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithSSCD",
	ServiceQualification_QC_WITH_QSCD:              "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD",
	ServiceQualification_QC_NO_SSCD:                "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoSSCD",
	ServiceQualification_QC_NO_QSCD:                "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoQSCD",
	ServiceQualification_QC_SSCD_STATUS_AS_IN_CERT: "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCSSCDStatusAsInCert",
	ServiceQualification_QC_QSCD_STATUS_AS_IN_CERT: "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDStatusAsInCert",
	ServiceQualification_QC_QSCD_MANAGED_ON_BEHALF: "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDManagedOnBehalf",
	ServiceQualification_QC_FOR_LEGAL_PERSON:       "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForLegalPerson",
	ServiceQualification_QC_FOR_ESIG:               "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESig",
	ServiceQualification_QC_FOR_ESEAL:              "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESeal",
	ServiceQualification_QC_FOR_WSA:                "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForWSA",
}

// ServiceQualificationValues returns all constants in declaration order.
func ServiceQualificationValues() []ServiceQualification {
	return []ServiceQualification{
		ServiceQualification_QC_STATEMENT,
		ServiceQualification_NOT_QUALIFIED,
		ServiceQualification_QC_WITH_SSCD,
		ServiceQualification_QC_WITH_QSCD,
		ServiceQualification_QC_NO_SSCD,
		ServiceQualification_QC_NO_QSCD,
		ServiceQualification_QC_SSCD_STATUS_AS_IN_CERT,
		ServiceQualification_QC_QSCD_STATUS_AS_IN_CERT,
		ServiceQualification_QC_QSCD_MANAGED_ON_BEHALF,
		ServiceQualification_QC_FOR_LEGAL_PERSON,
		ServiceQualification_QC_FOR_ESIG,
		ServiceQualification_QC_FOR_ESEAL,
		ServiceQualification_QC_FOR_WSA,
	}
}

// URI gets URI of the AdditionalServiceInformation.
func (s ServiceQualification) URI() string {
	return serviceQualificationURIs[s]
}

// ServiceQualificationGetByUri returns ServiceQualification for the given
// uri, if it exists, or "" (zero value) otherwise, mirroring Java's null
// return.
func ServiceQualificationGetByUri(uri string) ServiceQualification {
	if uri == "" {
		return ""
	}
	for _, v := range ServiceQualificationValues() {
		if serviceQualificationURIs[v] == uri {
			return v
		}
	}
	return ""
}

// serviceQualificationListContains reports whether qualifiers contains the
// URI of any of the given expected ServiceQualification values.
func serviceQualificationListContains(qualifiers []string, expecteds ...ServiceQualification) bool {
	if len(qualifiers) == 0 {
		return false
	}
	for _, expected := range expecteds {
		for _, q := range qualifiers {
			if q == expected.URI() {
				return true
			}
		}
	}
	return false
}

// ServiceQualificationIsQcStatement gets whether the list of TrustService
// qualifiers contains 'QCStatement' identifier.
func ServiceQualificationIsQcStatement(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_STATEMENT)
}

// ServiceQualificationIsNotQualified gets whether the list of TrustService
// qualifiers contains 'NotQualified' identifier.
func ServiceQualificationIsNotQualified(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_NOT_QUALIFIED)
}

// ServiceQualificationIsQcNoQSCD gets whether the list of TrustService
// qualifiers contains 'QCNoQSCD' identifier.
func ServiceQualificationIsQcNoQSCD(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_NO_QSCD)
}

// ServiceQualificationIsQcNoSSCD gets whether the list of TrustService
// qualifiers contains 'QCNoSSCD' identifier.
func ServiceQualificationIsQcNoSSCD(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_NO_SSCD)
}

// ServiceQualificationIsQcForLegalPerson gets whether the list of
// TrustService qualifiers contains 'QCForLegalPerson' identifier.
func ServiceQualificationIsQcForLegalPerson(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_FOR_LEGAL_PERSON)
}

// ServiceQualificationIsQcQSCDStatusAsInCert gets whether the list of
// TrustService qualifiers contains 'QCQSCDStatusAsInCert' identifier.
func ServiceQualificationIsQcQSCDStatusAsInCert(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_QSCD_STATUS_AS_IN_CERT)
}

// ServiceQualificationIsQcSSCDStatusAsInCert gets whether the list of
// TrustService qualifiers contains 'QCSSCDStatusAsInCert' identifier.
func ServiceQualificationIsQcSSCDStatusAsInCert(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_SSCD_STATUS_AS_IN_CERT)
}

// ServiceQualificationIsQcQSCDManagedOnBehalf gets whether the list of
// TrustService qualifiers contains 'QCQSCDManagedOnBehalf' identifier.
func ServiceQualificationIsQcQSCDManagedOnBehalf(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_QSCD_MANAGED_ON_BEHALF)
}

// ServiceQualificationIsQcWithQSCD gets whether the list of TrustService
// qualifiers contains 'QCWithQSCD' identifier.
func ServiceQualificationIsQcWithQSCD(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_WITH_QSCD)
}

// ServiceQualificationIsQcWithSSCD gets whether the list of TrustService
// qualifiers contains 'QCWithSSCD' identifier.
func ServiceQualificationIsQcWithSSCD(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_WITH_SSCD)
}

// ServiceQualificationIsQcForEsig gets whether the list of TrustService
// qualifiers contains 'QCForESig' identifier.
func ServiceQualificationIsQcForEsig(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_FOR_ESIG)
}

// ServiceQualificationIsQcForEseal gets whether the list of TrustService
// qualifiers contains 'QCForESeal' identifier.
func ServiceQualificationIsQcForEseal(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_FOR_ESEAL)
}

// ServiceQualificationIsQcForWSA gets whether the list of TrustService
// qualifiers contains 'QCForWSA' identifier.
func ServiceQualificationIsQcForWSA(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualification_QC_FOR_WSA)
}

// ServiceQualificationGetUsageQualifiers filters qualifiers containing the
// usage type identifiers (i.e. 'QCForESig', 'QCForESeal' or 'QCForWSA').
func ServiceQualificationGetUsageQualifiers(qualifiers []string) []string {
	filtered := []string{}
	if len(qualifiers) == 0 {
		return filtered
	}
	if serviceQualificationSliceContains(qualifiers, ServiceQualification_QC_FOR_ESIG.URI()) {
		filtered = append(filtered, ServiceQualification_QC_FOR_ESIG.URI())
	}
	if serviceQualificationSliceContains(qualifiers, ServiceQualification_QC_FOR_ESEAL.URI()) {
		filtered = append(filtered, ServiceQualification_QC_FOR_ESEAL.URI())
	}
	if serviceQualificationSliceContains(qualifiers, ServiceQualification_QC_FOR_WSA.URI()) {
		filtered = append(filtered, ServiceQualification_QC_FOR_WSA.URI())
	}
	return filtered
}

func serviceQualificationSliceContains(s []string, v string) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}
