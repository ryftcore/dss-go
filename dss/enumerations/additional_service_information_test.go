// Ported from dss-enumerations/.../AdditionalServiceInformation.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestAdditionalServiceInformation(t *testing.T) {
	cases := []struct {
		v   AdditionalServiceInformation
		uri string
	}{
		{AdditionalServiceInformation_FOR_ESIGNATURES, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSignatures"},
		{AdditionalServiceInformation_FOR_ESEALS, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSeals"},
		{AdditionalServiceInformation_FOR_WEB_AUTHENTICATION, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForWebSiteAuthentication"},
	}
	if len(AdditionalServiceInformationValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(AdditionalServiceInformationValues()))
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
		if got := AdditionalServiceInformationGetByUri(c.uri); got != c.v {
			t.Errorf("AdditionalServiceInformationGetByUri(%q) = %v, want %v", c.uri, got, c.v)
		}
	}
	if got := AdditionalServiceInformationGetByUri("unknown"); got != "" {
		t.Errorf("AdditionalServiceInformationGetByUri(unknown) = %v, want \"\"", got)
	}
	if got := AdditionalServiceInformationGetByUri(""); got != "" {
		t.Errorf("AdditionalServiceInformationGetByUri(\"\") = %v, want \"\"", got)
	}
}

func TestAdditionalServiceInformationPredicates(t *testing.T) {
	if !AdditionalServiceInformationIsForeSignatures(AdditionalServiceInformation_FOR_ESIGNATURES.URI()) {
		t.Error("expected IsForeSignatures true")
	}
	if AdditionalServiceInformationIsForeSignatures("other") {
		t.Error("expected IsForeSignatures false")
	}
	if !AdditionalServiceInformationIsForeSeals(AdditionalServiceInformation_FOR_ESEALS.URI()) {
		t.Error("expected IsForeSeals true")
	}
	if !AdditionalServiceInformationIsForWebAuth(AdditionalServiceInformation_FOR_WEB_AUTHENTICATION.URI()) {
		t.Error("expected IsForWebAuth true")
	}

	list := []string{AdditionalServiceInformation_FOR_ESIGNATURES.URI()}
	if !AdditionalServiceInformationIsForeSignaturesList(list) {
		t.Error("expected list to contain eSignatures")
	}
	if !AdditionalServiceInformationIsForeSignaturesOnly(list) {
		t.Error("expected list to only contain eSignatures")
	}
	multi := []string{AdditionalServiceInformation_FOR_ESIGNATURES.URI(), AdditionalServiceInformation_FOR_ESEALS.URI()}
	if AdditionalServiceInformationIsForeSignaturesOnly(multi) {
		t.Error("expected multi list to not be eSignatures-only")
	}
	if AdditionalServiceInformationIsForeSealsOnly(nil) {
		t.Error("expected nil list to not be eSeals-only")
	}
	if !AdditionalServiceInformationIsForWebAuthOnly([]string{AdditionalServiceInformation_FOR_WEB_AUTHENTICATION.URI()}) {
		t.Error("expected list to be web-auth-only")
	}
}
