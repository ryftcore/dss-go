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
		{ServiceQualificationQCStatement, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCStatement"},
		{ServiceQualificationNotQualified, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotQualified"},
		{ServiceQualificationQCWithSSCD, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithSSCD"},
		{ServiceQualificationQCWithQSCD, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD"},
		{ServiceQualificationQCNoSSCD, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoSSCD"},
		{ServiceQualificationQCNoQSCD, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoQSCD"},
		{ServiceQualificationQCSSCDStatusAsInCert, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCSSCDStatusAsInCert"},
		{ServiceQualificationQCQSCDStatusAsInCert, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDStatusAsInCert"},
		{ServiceQualificationQCQSCDManagedOnBehalf, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDManagedOnBehalf"},
		{ServiceQualificationQCForLegalPerson, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForLegalPerson"},
		{ServiceQualificationQCForESig, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESig"},
		{ServiceQualificationQCForESeal, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESeal"},
		{ServiceQualificationQCForWSA, "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForWSA"},
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
		ServiceQualificationQCStatement:           ServiceQualificationIsQcStatement,
		ServiceQualificationNotQualified:          ServiceQualificationIsNotQualified,
		ServiceQualificationQCNoQSCD:              ServiceQualificationIsQcNoQSCD,
		ServiceQualificationQCNoSSCD:              ServiceQualificationIsQcNoSSCD,
		ServiceQualificationQCForLegalPerson:      ServiceQualificationIsQcForLegalPerson,
		ServiceQualificationQCQSCDStatusAsInCert:  ServiceQualificationIsQcQSCDStatusAsInCert,
		ServiceQualificationQCSSCDStatusAsInCert:  ServiceQualificationIsQcSSCDStatusAsInCert,
		ServiceQualificationQCQSCDManagedOnBehalf: ServiceQualificationIsQcQSCDManagedOnBehalf,
		ServiceQualificationQCWithQSCD:            ServiceQualificationIsQcWithQSCD,
		ServiceQualificationQCWithSSCD:            ServiceQualificationIsQcWithSSCD,
		ServiceQualificationQCForESig:             ServiceQualificationIsQcForEsig,
		ServiceQualificationQCForESeal:            ServiceQualificationIsQcForEseal,
		ServiceQualificationQCForWSA:              ServiceQualificationIsQcForWSA,
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
		ServiceQualificationQCForWSA.URI(),
		ServiceQualificationQCStatement.URI(),
		ServiceQualificationQCForESig.URI(),
	}
	want := []string{ServiceQualificationQCForESig.URI(), ServiceQualificationQCForWSA.URI()}
	got := ServiceQualificationGetUsageQualifiers(in)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ServiceQualificationGetUsageQualifiers(%v) = %v, want %v", in, got, want)
	}
	if got := ServiceQualificationGetUsageQualifiers(nil); len(got) != 0 {
		t.Errorf("ServiceQualificationGetUsageQualifiers(nil) = %v, want empty", got)
	}
}
