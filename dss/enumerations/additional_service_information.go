// Ported from dss-enumerations/.../AdditionalServiceInformation.java (DSS 6.5.RC1).

package enumerations

// AdditionalServiceInformation represents an AdditionalServiceInformation
// element content present in a Trusted List.
type AdditionalServiceInformation string

const (
	// AdditionalServiceInformationForESignatures further specifies the
	// "Service type identifier" identified service as being provided for
	// electronic signatures.
	AdditionalServiceInformationForESignatures AdditionalServiceInformation = "FOR_ESIGNATURES"
	// AdditionalServiceInformationForESeals further specifies the
	// "Service type identifier" identified service as being provided for
	// electronic seals.
	AdditionalServiceInformationForESeals AdditionalServiceInformation = "FOR_ESEALS"
	// AdditionalServiceInformationForWebAuthentication further specifies
	// the "Service type identifier" identified service as being provided
	// for web site authentication.
	AdditionalServiceInformationForWebAuthentication AdditionalServiceInformation = "FOR_WEB_AUTHENTICATION"
)

// additionalServiceInformationURIs holds the URI for each constant.
var additionalServiceInformationURIs = map[AdditionalServiceInformation]string{
	AdditionalServiceInformationForESignatures:       "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSignatures",
	AdditionalServiceInformationForESeals:            "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSeals",
	AdditionalServiceInformationForWebAuthentication: "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForWebSiteAuthentication",
}

// AdditionalServiceInformationValues returns all constants in declaration order.
func AdditionalServiceInformationValues() []AdditionalServiceInformation {
	return []AdditionalServiceInformation{
		AdditionalServiceInformationForESignatures,
		AdditionalServiceInformationForESeals,
		AdditionalServiceInformationForWebAuthentication,
	}
}

// URI returns the URI of the AdditionalServiceInformation.
func (a AdditionalServiceInformation) URI() string {
	return additionalServiceInformationURIs[a]
}

// AdditionalServiceInformationGetByUri returns the AdditionalServiceInformation
// for the given uri, or "" if none matches.
func AdditionalServiceInformationGetByUri(uri string) AdditionalServiceInformation {
	if uri != "" {
		for _, v := range AdditionalServiceInformationValues() {
			if uri == v.URI() {
				return v
			}
		}
	}
	return ""
}

// AdditionalServiceInformationIsForeSignatures checks if the given additional
// service info is the "for eSignatures" identifier.
func AdditionalServiceInformationIsForeSignatures(additionalServiceInfo string) bool {
	return string(AdditionalServiceInformationForESignatures.URI()) == additionalServiceInfo
}

// AdditionalServiceInformationIsForeSeals checks if the given additional
// service info is the "for eSeals" identifier.
func AdditionalServiceInformationIsForeSeals(additionalServiceInfo string) bool {
	return AdditionalServiceInformationForESeals.URI() == additionalServiceInfo
}

// AdditionalServiceInformationIsForWebAuth checks if the given additional
// service info is the "for web authentication" identifier.
func AdditionalServiceInformationIsForWebAuth(additionalServiceInfo string) bool {
	return AdditionalServiceInformationForWebAuthentication.URI() == additionalServiceInfo
}

// AdditionalServiceInformationIsForeSignaturesList checks if the given list
// of additional service infos contains the "for eSignatures" identifier.
func AdditionalServiceInformationIsForeSignaturesList(additionalServiceInfos []string) bool {
	return additionalServiceInformationSliceContains(additionalServiceInfos, AdditionalServiceInformationForESignatures.URI())
}

// AdditionalServiceInformationIsForeSealsList checks if the given list of
// additional service infos contains the "for eSeals" identifier.
func AdditionalServiceInformationIsForeSealsList(additionalServiceInfos []string) bool {
	return additionalServiceInformationSliceContains(additionalServiceInfos, AdditionalServiceInformationForESeals.URI())
}

// AdditionalServiceInformationIsForWebAuthList checks if the given list of
// additional service infos contains the "for web authentication" identifier.
func AdditionalServiceInformationIsForWebAuthList(additionalServiceInfos []string) bool {
	return additionalServiceInformationSliceContains(additionalServiceInfos, AdditionalServiceInformationForWebAuthentication.URI())
}

// AdditionalServiceInformationIsForeSignaturesOnly checks if the given list
// of additional service infos only contains the "for eSignatures" identifier.
func AdditionalServiceInformationIsForeSignaturesOnly(additionalServiceInfos []string) bool {
	return additionalServiceInfos != nil && len(additionalServiceInfos) == 1 && AdditionalServiceInformationIsForeSignaturesList(additionalServiceInfos)
}

// AdditionalServiceInformationIsForeSealsOnly checks if the given list of
// additional service infos only contains the "for eSeals" identifier.
func AdditionalServiceInformationIsForeSealsOnly(additionalServiceInfos []string) bool {
	return additionalServiceInfos != nil && len(additionalServiceInfos) == 1 && AdditionalServiceInformationIsForeSealsList(additionalServiceInfos)
}

// AdditionalServiceInformationIsForWebAuthOnly checks if the given list of
// additional service infos only contains the "for web authentication"
// identifier.
func AdditionalServiceInformationIsForWebAuthOnly(additionalServiceInfos []string) bool {
	return additionalServiceInfos != nil && len(additionalServiceInfos) == 1 && AdditionalServiceInformationIsForWebAuthList(additionalServiceInfos)
}

// additionalServiceInformationSliceContains reports whether s contains v.
func additionalServiceInformationSliceContains(s []string, v string) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}
