// Ported from dss-enumerations/.../ServiceQualification.java (DSS 6.5.RC1).
package enumerations

// ServiceQualification contains qualification statuses for TrustServices.
type ServiceQualification string

const (
	// ServiceQualificationQCStatement
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCStatement"):
	// to indicate that all certificates identified by the applicable list
	// of criteria are issued as qualified certificates.
	ServiceQualificationQCStatement ServiceQualification = "QC_STATEMENT"
	// ServiceQualificationNotQualified
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotQualified"):
	// to indicate that no certificates identified by the applicable list of
	// criteria are to be considered as qualified certificates.
	ServiceQualificationNotQualified ServiceQualification = "NOT_QUALIFIED"
	// ServiceQualificationQCWithSSCD
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithSSCD"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, have
	// their private key residing in an SSCD.
	ServiceQualificationQCWithSSCD ServiceQualification = "QC_WITH_SSCD"
	// ServiceQualificationQCWithQSCD
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, have
	// their private key residing in a QSCD.
	ServiceQualificationQCWithQSCD ServiceQualification = "QC_WITH_QSCD"
	// ServiceQualificationQCNoSSCD
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoSSCD"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, do not
	// have their private key residing in an SSCD.
	ServiceQualificationQCNoSSCD ServiceQualification = "QC_NO_SSCD"
	// ServiceQualificationQCNoQSCD
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoQSCD"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, do not
	// have their private key residing in a QSCD.
	ServiceQualificationQCNoQSCD ServiceQualification = "QC_NO_QSCD"
	// ServiceQualificationQCSSCDStatusAsInCert
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCSSCDStatusAsInCert"):
	// to indicate that all certificates identified by the applicable list
	// of criteria, when they are claimed or stated as being qualified, do
	// contain proper machine processable information about whether or not
	// their private key residing in an SSCD.
	ServiceQualificationQCSSCDStatusAsInCert ServiceQualification = "QC_SSCD_STATUS_AS_IN_CERT"
	// ServiceQualificationQCQSCDStatusAsInCert
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDStatusAsInCert"):
	// to indicate that all certificates identified by the applicable list
	// of criteria, when they are claimed or stated as being qualified, do
	// contain proper machine processable information about whether or not
	// their private key residing in a QSCD.
	ServiceQualificationQCQSCDStatusAsInCert ServiceQualification = "QC_QSCD_STATUS_AS_IN_CERT"
	// ServiceQualificationQCQSCDManagedOnBehalf
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDManagedOnBehalf"):
	// to indicate that all certificates identified by the applicable list
	// of criteria, when they are claimed or stated as being qualified, have
	// their private key residing in a QSCD for which the generation and
	// management of that private key is done by the qualified TSP on
	// behalf of the entity whose identity is certified in the certificate.
	ServiceQualificationQCQSCDManagedOnBehalf ServiceQualification = "QC_QSCD_MANAGED_ON_BEHALF"
	// ServiceQualificationQCForLegalPerson
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForLegalPerson"):
	// to indicate that all certificates identified by the applicable list
	// of criteria, when they are claimed or stated as being qualified, are
	// issued to legal persons.
	ServiceQualificationQCForLegalPerson ServiceQualification = "QC_FOR_LEGAL_PERSON"
	// ServiceQualificationQCForESig
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESig"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, are
	// issued for electronic signatures.
	ServiceQualificationQCForESig ServiceQualification = "QC_FOR_ESIG"
	// ServiceQualificationQCForESeal
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESeal"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, are
	// issued for electronic seals.
	ServiceQualificationQCForESeal ServiceQualification = "QC_FOR_ESEAL"
	// ServiceQualificationQCForWSA
	// ("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForWSA"): to
	// indicate that all certificates identified by the applicable list of
	// criteria, when they are claimed or stated as being qualified, are
	// issued for web site authentication.
	ServiceQualificationQCForWSA ServiceQualification = "QC_FOR_WSA"
)

// serviceQualificationURIs holds the URI for each constant.
var serviceQualificationURIs = map[ServiceQualification]string{
	ServiceQualificationQCStatement:           "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCStatement",
	ServiceQualificationNotQualified:          "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotQualified",
	ServiceQualificationQCWithSSCD:            "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithSSCD",
	ServiceQualificationQCWithQSCD:            "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD",
	ServiceQualificationQCNoSSCD:              "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoSSCD",
	ServiceQualificationQCNoQSCD:              "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoQSCD",
	ServiceQualificationQCSSCDStatusAsInCert:  "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCSSCDStatusAsInCert",
	ServiceQualificationQCQSCDStatusAsInCert:  "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDStatusAsInCert",
	ServiceQualificationQCQSCDManagedOnBehalf: "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDManagedOnBehalf",
	ServiceQualificationQCForLegalPerson:      "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForLegalPerson",
	ServiceQualificationQCForESig:             "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESig",
	ServiceQualificationQCForESeal:            "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESeal",
	ServiceQualificationQCForWSA:              "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForWSA",
}

// ServiceQualificationValues returns all constants in declaration order.
func ServiceQualificationValues() []ServiceQualification {
	return []ServiceQualification{
		ServiceQualificationQCStatement,
		ServiceQualificationNotQualified,
		ServiceQualificationQCWithSSCD,
		ServiceQualificationQCWithQSCD,
		ServiceQualificationQCNoSSCD,
		ServiceQualificationQCNoQSCD,
		ServiceQualificationQCSSCDStatusAsInCert,
		ServiceQualificationQCQSCDStatusAsInCert,
		ServiceQualificationQCQSCDManagedOnBehalf,
		ServiceQualificationQCForLegalPerson,
		ServiceQualificationQCForESig,
		ServiceQualificationQCForESeal,
		ServiceQualificationQCForWSA,
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
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCStatement)
}

// ServiceQualificationIsNotQualified gets whether the list of TrustService
// qualifiers contains 'NotQualified' identifier.
func ServiceQualificationIsNotQualified(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationNotQualified)
}

// ServiceQualificationIsQcNoQSCD gets whether the list of TrustService
// qualifiers contains 'QCNoQSCD' identifier.
func ServiceQualificationIsQcNoQSCD(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCNoQSCD)
}

// ServiceQualificationIsQcNoSSCD gets whether the list of TrustService
// qualifiers contains 'QCNoSSCD' identifier.
func ServiceQualificationIsQcNoSSCD(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCNoSSCD)
}

// ServiceQualificationIsQcForLegalPerson gets whether the list of
// TrustService qualifiers contains 'QCForLegalPerson' identifier.
func ServiceQualificationIsQcForLegalPerson(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCForLegalPerson)
}

// ServiceQualificationIsQcQSCDStatusAsInCert gets whether the list of
// TrustService qualifiers contains 'QCQSCDStatusAsInCert' identifier.
func ServiceQualificationIsQcQSCDStatusAsInCert(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCQSCDStatusAsInCert)
}

// ServiceQualificationIsQcSSCDStatusAsInCert gets whether the list of
// TrustService qualifiers contains 'QCSSCDStatusAsInCert' identifier.
func ServiceQualificationIsQcSSCDStatusAsInCert(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCSSCDStatusAsInCert)
}

// ServiceQualificationIsQcQSCDManagedOnBehalf gets whether the list of
// TrustService qualifiers contains 'QCQSCDManagedOnBehalf' identifier.
func ServiceQualificationIsQcQSCDManagedOnBehalf(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCQSCDManagedOnBehalf)
}

// ServiceQualificationIsQcWithQSCD gets whether the list of TrustService
// qualifiers contains 'QCWithQSCD' identifier.
func ServiceQualificationIsQcWithQSCD(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCWithQSCD)
}

// ServiceQualificationIsQcWithSSCD gets whether the list of TrustService
// qualifiers contains 'QCWithSSCD' identifier.
func ServiceQualificationIsQcWithSSCD(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCWithSSCD)
}

// ServiceQualificationIsQcForEsig gets whether the list of TrustService
// qualifiers contains 'QCForESig' identifier.
func ServiceQualificationIsQcForEsig(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCForESig)
}

// ServiceQualificationIsQcForEseal gets whether the list of TrustService
// qualifiers contains 'QCForESeal' identifier.
func ServiceQualificationIsQcForEseal(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCForESeal)
}

// ServiceQualificationIsQcForWSA gets whether the list of TrustService
// qualifiers contains 'QCForWSA' identifier.
func ServiceQualificationIsQcForWSA(qualifiers []string) bool {
	return serviceQualificationListContains(qualifiers, ServiceQualificationQCForWSA)
}

// ServiceQualificationGetUsageQualifiers filters qualifiers containing the
// usage type identifiers (i.e. 'QCForESig', 'QCForESeal' or 'QCForWSA').
func ServiceQualificationGetUsageQualifiers(qualifiers []string) []string {
	filtered := []string{}
	if len(qualifiers) == 0 {
		return filtered
	}
	if serviceQualificationSliceContains(qualifiers, ServiceQualificationQCForESig.URI()) {
		filtered = append(filtered, ServiceQualificationQCForESig.URI())
	}
	if serviceQualificationSliceContains(qualifiers, ServiceQualificationQCForESeal.URI()) {
		filtered = append(filtered, ServiceQualificationQCForESeal.URI())
	}
	if serviceQualificationSliceContains(qualifiers, ServiceQualificationQCForWSA.URI()) {
		filtered = append(filtered, ServiceQualificationQCForWSA.URI())
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
