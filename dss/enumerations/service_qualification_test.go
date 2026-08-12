package enumerations

import (
	"reflect"
	"testing"
)

type serviceQualificationCase struct {
	v   ServiceQualification
	uri string
}

func serviceQualificationCases() []serviceQualificationCase {
	return []serviceQualificationCase{
		{ServiceQualification_QC_STATEMENT, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCStatement"},
		{ServiceQualification_NOT_QUALIFIED, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotQualified"},
		{ServiceQualification_QC_WITH_SSCD, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithSSCD"},
		{ServiceQualification_QC_WITH_QSCD, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD"},
		{ServiceQualification_QC_NO_SSCD, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoSSCD"},
		{ServiceQualification_QC_NO_QSCD, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoQSCD"},
		{ServiceQualification_QC_SSCD_STATUS_AS_IN_CERT, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCSSCDStatusAsInCert"},
		{ServiceQualification_QC_QSCD_STATUS_AS_IN_CERT, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDStatusAsInCert"},
		{ServiceQualification_QC_QSCD_MANAGED_ON_BEHALF, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDManagedOnBehalf"},
		{ServiceQualification_QC_FOR_LEGAL_PERSON, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForLegalPerson"},
		{ServiceQualification_QC_FOR_ESIG, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESig"},
		{ServiceQualification_QC_FOR_ESEAL, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESeal"},
		{ServiceQualification_QC_FOR_WSA, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForWSA"},
	}
}

func TestServiceQualificationURI(t *testing.T) {
	cases := serviceQualificationCases()
	if len(ServiceQualificationValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(ServiceQualificationValues()))
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
	}
}

func TestServiceQualificationGetByUri(t *testing.T) {
	for _, c := range serviceQualificationCases() {
		if got := ServiceQualificationGetByUri(c.uri); got != c.v {
			t.Errorf("ServiceQualificationGetByUri(%q) = %v, want %v", c.uri, got, c.v)
		}
	}
	if got := ServiceQualificationGetByUri(""); got != "" {
		t.Errorf("ServiceQualificationGetByUri(\"\") = %v, want zero value", got)
	}
	if got := ServiceQualificationGetByUri("nope"); got != "" {
		t.Errorf("ServiceQualificationGetByUri(unknown) = %v, want zero value", got)
	}
}

func TestServiceQualificationPredicates(t *testing.T) {
	predicates := map[ServiceQualification]func([]string) bool{
		ServiceQualification_QC_STATEMENT:              ServiceQualificationIsQcStatement,
		ServiceQualification_NOT_QUALIFIED:             ServiceQualificationIsNotQualified,
		ServiceQualification_QC_NO_QSCD:                ServiceQualificationIsQcNoQSCD,
		ServiceQualification_QC_NO_SSCD:                ServiceQualificationIsQcNoSSCD,
		ServiceQualification_QC_FOR_LEGAL_PERSON:       ServiceQualificationIsQcForLegalPerson,
		ServiceQualification_QC_QSCD_STATUS_AS_IN_CERT: ServiceQualificationIsQcQSCDStatusAsInCert,
		ServiceQualification_QC_SSCD_STATUS_AS_IN_CERT: ServiceQualificationIsQcSSCDStatusAsInCert,
		ServiceQualification_QC_QSCD_MANAGED_ON_BEHALF: ServiceQualificationIsQcQSCDManagedOnBehalf,
		ServiceQualification_QC_WITH_QSCD:              ServiceQualificationIsQcWithQSCD,
		ServiceQualification_QC_WITH_SSCD:              ServiceQualificationIsQcWithSSCD,
		ServiceQualification_QC_FOR_ESIG:               ServiceQualificationIsQcForEsig,
		ServiceQualification_QC_FOR_ESEAL:              ServiceQualificationIsQcForEseal,
		ServiceQualification_QC_FOR_WSA:                ServiceQualificationIsQcForWSA,
	}
	for v, pred := range predicates {
		if !pred([]string{v.URI()}) {
			t.Errorf("predicate for %v = false with matching list, want true", v)
		}
		if pred([]string{"http://example.org/other"}) {
			t.Errorf("predicate for %v = true with non-matching list, want false", v)
		}
		if pred(nil) {
			t.Errorf("predicate for %v = true with nil list, want false", v)
		}
	}
}

func TestServiceQualificationGetUsageQualifiers(t *testing.T) {
	in := []string{
		ServiceQualification_QC_FOR_WSA.URI(),
		ServiceQualification_QC_STATEMENT.URI(),
		ServiceQualification_QC_FOR_ESIG.URI(),
	}
	want := []string{ServiceQualification_QC_FOR_ESIG.URI(), ServiceQualification_QC_FOR_WSA.URI()}
	got := ServiceQualificationGetUsageQualifiers(in)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ServiceQualificationGetUsageQualifiers(%v) = %v, want %v", in, got, want)
	}
	if got := ServiceQualificationGetUsageQualifiers(nil); len(got) != 0 {
		t.Errorf("ServiceQualificationGetUsageQualifiers(nil) = %v, want empty", got)
	}
}
