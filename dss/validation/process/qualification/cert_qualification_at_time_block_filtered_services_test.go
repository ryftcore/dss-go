package qualification

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
)

// Upstream's CertQualificationAtTimeBlock.getFilteredServices() throws only
// when execute() has not run yet:
//
//	if (filteredServices == null) throw new IllegalStateException(...)
//
// After execute() the field is never null - it starts as
// `new ArrayList<>(acceptableServices)`, every TrustServiceFilter returns
// `new ArrayList<>()` rather than null, and the "Keep only CA/QC and granted"
// branch assigns `Collections.emptyList()`. SignatureQualificationBlock relies
// on that: it calls getFilteredServices() on both at-time blocks
// unconditionally once an acceptable TL is present, and merely asks whether the
// list is empty.
//
// These cases pin the two ways the list legitimately empties out.
func TestCertQualificationAtTimeBlockFilteredServicesEmptyNotNil(t *testing.T) {
	postEIDAS := time.UnixMilli(1514764800000).UTC() // 2018-01-01
	probe := time.UnixMilli(1420070400000).UTC()     // 2015-01-01, before every service window

	cert := buildQualCert(&qualCertInput{
		QcCompliance:        boolPtr(true),
		QcTypes:             []string{"0.4.0.1862.1.6.1"},
		NotBefore:           postEIDAS.UnixMilli(),
		QcStatementsPresent: true,
	})

	for _, tc := range []struct {
		name     string
		date     *time.Time
		services func() []*diagnostic.TrustServiceWrapper
	}{
		{
			// Every service is dropped by the very first filter (by date), so
			// filteredServices is emptied before any chain item runs.
			name: "all services filtered out by date",
			date: &probe,
			services: func() []*diagnostic.TrustServiceWrapper {
				return []*diagnostic.TrustServiceWrapper{buildQualCertService(1)}
			},
		},
		{
			// The service survives the date filter but is neither CA/QC-and-granted,
			// so the "Keep only CA/QC and granted" branch empties the list at the end.
			name: "selected service neither CA/QC nor granted",
			date: &postEIDAS,
			services: func() []*diagnostic.TrustServiceWrapper {
				svc := buildQualCertService(11) // withdrawn-post
				svc.Type = "http://uri.etsi.org/TrstSvc/Svctype/CA/PKC"
				return []*diagnostic.TrustServiceWrapper{svc}
			},
		},
		{
			// No acceptable services at all.
			name:     "no acceptable services",
			date:     &postEIDAS,
			services: func() []*diagnostic.TrustServiceWrapper { return nil },
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			block := NewCertQualificationAtTimeBlock(i18n.NewI18nProvider(),
				enumerations.ValidationTimeBESTSignatureTime, tc.date, cert, tc.services())
			block.Execute()

			// Must not panic: upstream returns an empty list here, and
			// SignatureQualificationBlock calls this unconditionally.
			filtered := block.FilteredServices()
			if len(filtered) != 0 {
				t.Fatalf("FilteredServices() = %d services, want 0", len(filtered))
			}
			if filtered == nil {
				t.Error("FilteredServices() returned a nil slice; upstream returns an empty list")
			}
		})
	}
}

// The pre-execute() contract is the one upstream really does throw on.
func TestCertQualificationAtTimeBlockFilteredServicesBeforeExecute(t *testing.T) {
	postEIDAS := time.UnixMilli(1514764800000).UTC()
	cert := buildQualCert(&qualCertInput{NotBefore: postEIDAS.UnixMilli()})
	block := NewCertQualificationAtTimeBlock(i18n.NewI18nProvider(),
		enumerations.ValidationTimeBESTSignatureTime, &postEIDAS, cert, nil)

	defer func() {
		if recover() == nil {
			t.Error("FilteredServices() before Execute() did not panic; upstream throws IllegalStateException")
		}
	}()
	block.FilteredServices()
}

func boolPtr(b bool) *bool { return &b }

// Upstream's TrustServiceFilter implementations all build their result from
// `new ArrayList<>()`, so filter() never returns null even when nothing is
// accepted. CertQualificationAtTimeBlock stores that result straight into
// filteredServices, whose null-ness is the "execute() not called yet" signal -
// so a nil-returning filter would make getFilteredServices() throw on an
// ordinary empty-result input.
func TestTrustServiceFiltersReturnEmptyNotNil(t *testing.T) {
	// one service that every filter below rejects
	svc := buildQualCertService(11) // withdrawn-post
	svc.Type = "http://uri.etsi.org/TrstSvc/Svctype/CA/PKC"
	svc.CountryCode = "ZZ"
	url := "https://tl.example/other"
	svc.TrustedList.Url = &url
	input := []*diagnostic.TrustServiceWrapper{svc}
	early := time.UnixMilli(1000000000).UTC()

	cert := buildQualCert(&qualCertInput{NotBefore: 1514764800000})

	for name, filter := range map[string]TrustServiceFilter{
		"byGranted":                 TrustServicesFilterFactoryCreateFilterByGranted(),
		"byCaQc":                    TrustServicesFilterFactoryCreateFilterByCaQc(),
		"byQTST":                    TrustServicesFilterFactoryCreateFilterByQTST(),
		"byQEAA":                    TrustServicesFilterFactoryCreateFilterByQEAA(),
		"byDate":                    TrustServicesFilterFactoryCreateFilterByDate(&early),
		"byCountry":                 TrustServicesFilterFactoryCreateFilterByCountry("LU"),
		"byCountries":               TrustServicesFilterFactoryCreateFilterByCountries(map[string]struct{}{"LU": {}}),
		"byUrls":                    TrustServicesFilterFactoryCreateFilterByUrls(map[string]struct{}{"https://tl.example/LU": {}}),
		"mraEnacted":                TrustServicesFilterFactoryCreateMRAEnactedFilter(),
		"byMRAEquivalenceStartDate": TrustServicesFilterFactoryCreateFilterByMRAEquivalenceStartingDate(&early),
	} {
		name, filter := name, filter
		t.Run(name, func(t *testing.T) {
			got := filter.Filter(input)
			if len(got) != 0 {
				t.Fatalf("filter accepted %d services, want 0", len(got))
			}
			if got == nil {
				t.Error("filter returned a nil slice; upstream returns an empty ArrayList")
			}
			if empty := filter.Filter(nil); empty == nil {
				t.Error("filter over an empty input returned a nil slice; upstream returns an empty ArrayList")
			}
		})
	}

	// ServiceByCertificateTypeFilter accepts this service, so it is only checked
	// on the empty input; UniqueServiceFilter is checked on a conflicting pair,
	// the shape on which it deliberately selects nothing.
	byType := TrustServicesFilterFactoryCreateFilterByCertificateType(cert)
	if empty := byType.Filter(nil); empty == nil {
		t.Error("byCertificateType over an empty input returned a nil slice; upstream returns an empty ArrayList")
	}
	unique := TrustServicesFilterFactoryCreateUniqueServiceFilter(cert)
	if empty := unique.Filter(nil); empty == nil {
		t.Error("uniqueService over an empty input returned a nil slice; upstream returns an empty ArrayList")
	}
	conflicting := []*diagnostic.TrustServiceWrapper{buildQualCertService(2), buildQualCertService(3)}
	if selected := unique.Filter(conflicting); len(selected) != 0 {
		t.Errorf("uniqueService selected %d services over a conflicting pair, want 0", len(selected))
	} else if selected == nil {
		t.Error("uniqueService returned a nil slice; upstream returns Collections.emptyList()")
	}

	// the composite (non-AbstractTrustServiceFilter) filters take the same contract
	consistentInput := []*diagnostic.TrustServiceWrapper{buildQualCertService(1)}
	consistentInput[0].CapturedQualifiers = append(consistentInput[0].CapturedQualifiers,
		&diagnosticjaxb.XmlQualifier{Value: "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD"},
		&diagnosticjaxb.XmlQualifier{Value: "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoQSCD"})
	for name, filter := range map[string]TrustServiceFilter{
		"consistentByQSCD":            TrustServicesFilterFactoryCreateConsistentServiceByQSCDFilter(),
		"consistentByStatus":          TrustServicesFilterFactoryCreateConsistentServiceByStatusFilter(),
		"consistentByQC":              TrustServicesFilterFactoryCreateConsistentServiceByQCFilter(),
		"consistentByCertificateType": TrustServicesFilterFactoryCreateConsistentServiceByCertificateTypeFilter(),
	} {
		name, filter := name, filter
		t.Run(name, func(t *testing.T) {
			if empty := filter.Filter(nil); empty == nil {
				t.Error("filter over an empty input returned a nil slice; upstream returns an empty ArrayList")
			}
			if rejected := filter.Filter(consistentInput); rejected == nil {
				t.Error("filter returned a nil slice; upstream returns an empty ArrayList")
			}
		})
	}
}

// The trusted-entity half of the filter family takes the same
// "never return null" contract, and CertificateApprovalStatusAtTimeBlock
// stores the result the same way CertQualificationAtTimeBlock does.
func TestTrustedEntityFiltersReturnEmptyNotNil(t *testing.T) {
	early := time.UnixMilli(1000000000).UTC()
	svc := &diagnostic.TrustedEntityServiceWrapper{
		Status: "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/withdrawn",
		Type:   "http://uri.etsi.org/TrstSvc/Svctype/CA/PKC",
	}
	start := time.UnixMilli(1514764800000).UTC()
	svc.StartDate = &start
	tsl := &diagnosticjaxb.XmlTrustSourceList{}
	url := "https://tl.example/other"
	tsl.Url = &url
	svc.TrustedSourceList = tsl
	input := []*diagnostic.TrustedEntityServiceWrapper{svc}

	for name, filter := range map[string]TrustedEntityServiceFilter{
		"byListUrl":  TrustedEntitiesFilterFactoryCreateFilterByListUrl("https://tl.example/LU"),
		"byListUrls": TrustedEntitiesFilterFactoryCreateFilterByListUrls([]string{"https://tl.example/LU"}),
		"byDate":     TrustedEntitiesFilterFactoryCreateFilterByDate(&early),
		"bySti": TrustedEntitiesFilterFactoryCreateFilterByServiceTypeIdentifierUri(
			"http://uri.etsi.org/TrstSvc/Svctype/CA/QC"),
		"byStatus": TrustedEntitiesFilterFactoryCreateFilterByServiceStatusUri(
			"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted"),
	} {
		name, filter := name, filter
		t.Run(name, func(t *testing.T) {
			got := filter.Filter(input)
			if len(got) != 0 {
				t.Fatalf("filter accepted %d services, want 0", len(got))
			}
			if got == nil {
				t.Error("filter returned a nil slice; upstream returns an empty ArrayList")
			}
			if empty := filter.Filter(nil); empty == nil {
				t.Error("filter over an empty input returned a nil slice; upstream returns an empty ArrayList")
			}
		})
	}
}
