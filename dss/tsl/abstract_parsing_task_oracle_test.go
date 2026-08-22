// Oracle test for abstract_parsing_task.go (and, through it, abstract_tl_parsing_result.go /
// tl_parsing_result.go / lotl_parsing_result.go).
//
// It checks the common SchemeInformation extraction field by field against the two REAL trusted
// lists this package's testdata carries - upstream's src/test/resources/eu-lotl.xml (the
// EU list of the lists, TSLType EUlistofthelists, sequence 248) and sk-tl.xml (a country TL,
// TSLType EUgeneric, sequence 59). Every expected value below is read straight out of those
// documents, i.e. it is what the Java AbstractParsingTask#commonParseSchemeInformation stores into
// the AbstractTLParsingResult for the same bytes. The whole-TLInfo dump comparison this
// complements lives in the harness stage.
//
// The task built here is the minimal concrete AbstractParsingTask a test can make - the two real
// subclasses (TLParsingTask, LOTLParsingTask) additionally need converters and predicates from
// elsewhere in this package, so they are exercised by the harness rather than here.
package tsl

import (
	"reflect"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
)

// schemeInformationTask is the minimal concrete AbstractParsingTask: it inherits the base's
// CreateTrustedListFacade (i.e. the plain Facade, like TLParsingTask does).
type schemeInformationTask struct {
	AbstractParsingTaskBase
}

func newSchemeInformationTask(document model.DSSDocument) *schemeInformationTask {
	task := &schemeInformationTask{AbstractParsingTaskBase: NewAbstractParsingTaskBase(document)}
	task.InitAbstractParsingTask(task)
	return task
}

func TestAbstractParsingTask_CommonParseSchemeInformationOracle(t *testing.T) {
	cases := []struct {
		file               string
		tslType            string
		sequenceNumber     int
		version            int
		territory          string
		issueDate          string
		nextUpdateDate     string
		distributionPoints []string
	}{
		{
			file:               "eu-lotl.xml",
			tslType:            "http://uri.etsi.org/TrstSvc/TrustedList/TSLType/EUlistofthelists",
			sequenceNumber:     248,
			version:            5,
			territory:          "EU",
			issueDate:          "2019-09-04T08:00:00Z",
			nextUpdateDate:     "2020-03-04T00:00:00Z",
			distributionPoints: []string{"https://ec.europa.eu/tools/lotl/eu-lotl.xml"},
		},
		{
			file:           "sk-tl.xml",
			tslType:        "http://uri.etsi.org/TrstSvc/TrustedList/TSLType/EUgeneric",
			sequenceNumber: 59,
			version:        5,
			territory:      "SK",
			issueDate:      "2019-08-23T07:00:00Z",
			nextUpdateDate: "2020-02-19T00:00:00Z",
			distributionPoints: []string{
				"http://ep.nbu.gov.sk/kca/tsl/tsl.xml",
				"http://ep.nbu.gov.sk/kca/tsl/tsl.xml.p7s",
				"http://tl.nbu.gov.sk/kca/tsl/tsl.xml",
				"http://tl.nbu.gov.sk/kca/tsl/tsl.xml.p7s",
				"https://tl.nbu.gov.sk/kca/tsl/tsl.xml",
				"https://tl.nbu.gov.sk/kca/tsl/tsl.xml.p7s",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.file, func(t *testing.T) {
			document, err := model.NewFileDocument(corpustest.Path(t, c.file))
			if err != nil {
				t.Fatalf("unable to load %s: %v", c.file, err)
			}
			task := newSchemeInformationTask(document)
			jaxbObject, err := task.JAXBObject()
			if err != nil {
				t.Fatalf("JAXBObject: %v", err)
			}

			result := NewTLParsingResult()
			task.CommonParseSchemeInformation(&result.AbstractTLParsingResult, jaxbObject.SchemeInformation)

			if result.TSLType() == nil || result.TSLType().URI() != c.tslType {
				t.Errorf("TSLType() = %v, want URI %q", result.TSLType(), c.tslType)
			}
			if result.SequenceNumber() == nil || *result.SequenceNumber() != c.sequenceNumber {
				t.Errorf("SequenceNumber() = %v, want %d", result.SequenceNumber(), c.sequenceNumber)
			}
			if result.Version() == nil || *result.Version() != c.version {
				t.Errorf("Version() = %v, want %d", result.Version(), c.version)
			}
			if result.Territory() != c.territory {
				t.Errorf("Territory() = %q, want %q", result.Territory(), c.territory)
			}
			assertParsedInstant(t, "IssueDate", result.IssueDate(), c.issueDate)
			assertParsedInstant(t, "NextUpdateDate", result.NextUpdateDate(), c.nextUpdateDate)
			if !reflect.DeepEqual(result.DistributionPoints(), c.distributionPoints) {
				t.Errorf("DistributionPoints() = %v, want %v", result.DistributionPoints(), c.distributionPoints)
			}
		})
	}
}

func assertParsedInstant(t *testing.T, name string, got time.Time, wantLexical string) {
	t.Helper()
	want, err := time.Parse(time.RFC3339, wantLexical)
	if err != nil {
		t.Fatalf("bad expected value %q: %v", wantLexical, err)
	}
	if !got.Equal(want) {
		t.Errorf("%s() = %s, want %s", name, got.Format(time.RFC3339), wantLexical)
	}
}

// TestAbstractParsingTask_ConvertToDate covers the XMLGregorianCalendar#toGregorianCalendar#getTime
// port: an offset in the lexical form is honoured, an absent one resolves in the local zone (Java's
// GregorianCalendar default TimeZone), and an unusable form answers the zero instant.
func TestAbstractParsingTask_ConvertToDate(t *testing.T) {
	lexical := func(s string) *string { return &s }

	if got := abstractParsingTaskConvertToDate(nil); !got.IsZero() {
		t.Errorf("nil = %v, want the zero time", got)
	}
	if got := abstractParsingTaskConvertToDate(lexical("not a date")); !got.IsZero() {
		t.Errorf("garbage = %v, want the zero time", got)
	}
	utc := abstractParsingTaskConvertToDate(lexical("2020-02-19T00:00:00Z"))
	if want := time.Date(2020, time.February, 19, 0, 0, 0, 0, time.UTC); !utc.Equal(want) {
		t.Errorf("Z form = %v, want %v", utc, want)
	}
	offset := abstractParsingTaskConvertToDate(lexical("2020-02-19T02:00:00+02:00"))
	if !offset.Equal(utc) {
		t.Errorf("+02:00 form = %v, want the same instant as %v", offset, utc)
	}
	fractional := abstractParsingTaskConvertToDate(lexical("2020-02-19T00:00:00.500Z"))
	if want := time.Date(2020, time.February, 19, 0, 0, 0, 500000000, time.UTC); !fractional.Equal(want) {
		t.Errorf("fractional form = %v, want %v", fractional, want)
	}
	local := abstractParsingTaskConvertToDate(lexical("2020-02-19T00:00:00"))
	if want := time.Date(2020, time.February, 19, 0, 0, 0, 0, time.Local); !local.Equal(want) {
		t.Errorf("timezone-less form = %v, want %v", local, want)
	}
}

// TestAbstractParsingTask_NullDocument covers the Objects.requireNonNull in the protected
// constructor.
func TestAbstractParsingTask_NullDocument(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("a nil document should have panicked")
		}
		if got, want := recovered.(string), "The document is null"; got != want {
			t.Errorf("panic = %q, want %q", got, want)
		}
	}()
	NewAbstractParsingTaskBase(nil)
}

// TestAbstractParsingTask_UnparseableDocument covers getJAXBObject's error wrapping.
func TestAbstractParsingTask_UnparseableDocument(t *testing.T) {
	task := newSchemeInformationTask(model.NewInMemoryDocument([]byte("this is not XML")))
	if _, err := task.JAXBObject(); err == nil {
		t.Error("a non-XML document should be refused")
	} else if got := err.Error(); len(got) < len("Unable to parse binaries. Reason : ") ||
		got[:len("Unable to parse binaries. Reason : ")] != "Unable to parse binaries. Reason : " {
		t.Errorf("error = %q, want the upstream prefix", got)
	}
}

// TestAbstractTLParsingResult_Defaults covers the "null everywhere" initial state of the abstract
// result, and the TSLType fallback TSLTypeFromURI applies to an unknown URI.
func TestAbstractTLParsingResult_Defaults(t *testing.T) {
	result := NewTLParsingResult()
	if result.TSLType() != nil || result.SequenceNumber() != nil || result.Version() != nil {
		t.Error("a fresh result should carry no TSLType/sequence/version")
	}
	if result.Territory() != "" || !result.IssueDate().IsZero() || !result.NextUpdateDate().IsZero() {
		t.Error("a fresh result should carry no territory/dates")
	}
	if result.DistributionPoints() != nil || result.TrustServiceProviders() != nil {
		t.Error("a fresh result should carry no distribution points/providers")
	}
	if got := enumerations.TSLTypeFromURI("urn:unknown"); got.URI() != "urn:unknown" || got.Label() != "" {
		t.Errorf("TSLTypeFromURI fallback = %v", got)
	}

	lotlResult := NewLOTLParsingResult()
	if lotlResult.LotlPointers() != nil || lotlResult.TlPointers() != nil ||
		lotlResult.PivotURLs() != nil || lotlResult.SigningCertificateAnnouncementURL() != "" {
		t.Error("a fresh LOTL result should carry no pointers/pivots/announcement URL")
	}
}

// TestTLSourceDefaults covers TLSource's DEFAULT_SUPPORTED_TL_VERSIONS and LOTLSource's field
// initialisers.
func TestTLSourceDefaults(t *testing.T) {
	source := NewTLSource()
	if got := source.TLVersions(); !reflect.DeepEqual(got, []int{5, 6}) {
		t.Errorf("TLVersions() = %v, want [5 6]", got)
	}
	if source.TrustServiceProviderPredicate() != nil || source.TrustServicePredicate() != nil ||
		source.TrustAnchorValidityPredicate() != nil {
		t.Error("a fresh TLSource should carry no predicates")
	}
	// The default list is copied per instance, so mutating one source cannot affect the next.
	source.TLVersions()[0] = 99
	if got := NewTLSource().TLVersions(); !reflect.DeepEqual(got, []int{5, 6}) {
		t.Errorf("the shared default was mutated: %v", got)
	}

	lotlSource := NewLOTLSource()
	if lotlSource.IsPivotSupport() || lotlSource.IsMraSupport() {
		t.Error("pivot/MRA support default to false")
	}
	if lotlSource.LotlPredicate() == nil || lotlSource.TlPredicate() == nil {
		t.Error("the LOTL/TL predicates are seeded by TLPredicateFactory")
	}
	if lotlSource.SigningCertificatesAnnouncementPredicate() != nil {
		t.Error("the signing-certificates announcement predicate is optional and starts unset")
	}
	if got := lotlSource.TLVersions(); !reflect.DeepEqual(got, []int{5, 6}) {
		t.Errorf("LOTLSource inherits TLVersions() = %v, want [5 6]", got)
	}
}
