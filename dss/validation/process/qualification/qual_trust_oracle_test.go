package qualification

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/utain/esig/dss/diagnostic"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/corpustest"
)

// The trust-service KAT for the qualification package: every row of
// testdata/oracle/qual_trust.jsonl is what upstream produces for one designed
// TrustServiceWrapper (or one matrix/vocabulary cell) - see
// testdata/gen/QualTrustOracle.java. It covers the ten TrustServiceChecker
// consistency predicates (and the twelve TrustServiceCondition classes behind
// them), the thirteen TrustServicesFilterFactory filters, QualificationByTL /
// TypeByTL / QSCDByTL, CertQualificationMatrix, SigQualificationMatrix,
// FinalCertificateQualificationCalculator, EIDASUtils, TrustServiceStatus and
// ServiceTypeIdentifier.

const qualTrustOracleLog = "testdata/oracle/qual_trust.jsonl"

type qualTrustInput struct {
	Status      *string  `json:"status"`
	Type        *string  `json:"type"`
	StartDate   *int64   `json:"startDate"`
	EndDate     *int64   `json:"endDate"`
	Qualifiers  []string `json:"qualifiers"`
	Asis        []string `json:"asis"`
	CountryCode *string  `json:"countryCode"`
	TlUrl       *string  `json:"tlUrl"`
	EnactedMRA  *bool    `json:"enactedMRA"`
	MraStart    *int64   `json:"mraStart"`
	MraEnd      *int64   `json:"mraEnd"`
}

type qualTrustChecker struct {
	LegalPerson          bool `json:"legalPerson"`
	QCStatement          bool `json:"qcStatement"`
	QSCD                 bool `json:"qscd"`
	QSCDStatusAsInCert   bool `json:"qscdStatusAsInCert"`
	PostEidasQSCD        bool `json:"postEidasQscd"`
	QualifiersListKnown  bool `json:"qualifiersListKnown"`
	Usage                bool `json:"usage"`
	PreEidasStatus       bool `json:"preEidasStatus"`
	PreEidasQualifierAsi bool `json:"preEidasQualifierAsi"`
	QualifierAsi         bool `json:"qualifierAsi"`
}

type qualTrustFilters struct {
	Granted            bool `json:"granted"`
	CaQc               bool `json:"caQc"`
	QTST               bool `json:"qtst"`
	QEAA               bool `json:"qeaa"`
	ConsistentStatus   bool `json:"consistentStatus"`
	ConsistentQC       bool `json:"consistentQC"`
	ConsistentQSCD     bool `json:"consistentQSCD"`
	ConsistentCertType bool `json:"consistentCertType"`
	MraEnacted         bool `json:"mraEnacted"`
	CountryLU          bool `json:"countryLU"`
	CountriesLUDE      bool `json:"countriesLUDE"`
	Urls               bool `json:"urls"`
}

type qualTrustRow struct {
	Kind string `json:"kind"`

	// kind == "service"
	I                            int               `json:"i"`
	In                           *qualTrustInput   `json:"in"`
	Checker                      *qualTrustChecker `json:"checker"`
	Filters                      *qualTrustFilters `json:"filters"`
	ByDate                       []bool            `json:"byDate"`
	ByMraEquivalenceStartingDate []bool            `json:"byMraEquivalenceStartingDate"`
	QualificationByTL            map[string]string `json:"qualificationByTL"`
	TypeByTL                     map[string]string `json:"typeByTL"`
	QSCDByTL                     map[string]string `json:"qscdByTL"`

	// kind == "certQualificationMatrix"
	QC   string `json:"qc"`
	QSCD string `json:"qscd"`

	// kind == "sigQualificationMatrix"
	Indication        string `json:"indication"`
	CertQualification string `json:"certQualification"`

	// kind == "finalCertQualification"
	AtIssuance string `json:"atIssuance"`
	AtSigning  string `json:"atSigning"`

	// shared output slot for the matrix kinds; also the CertificateType input
	// of certQualificationMatrix.
	Type *string `json:"type"`
	Out  string  `json:"out"`

	// kind == "eidas"
	Millis           *int64 `json:"millis"`
	IsPostEIDAS      bool   `json:"isPostEIDAS"`
	IsPreEIDAS       bool   `json:"isPreEIDAS"`
	IsPostGracePerio bool   `json:"isPostGracePeriod"`

	// kind == "status" / "serviceType" / "statusProbe" / "serviceTypeProbe"
	Name             string  `json:"name"`
	URI              *string `json:"uri"`
	ShortName        string  `json:"shortName"`
	IsPostEidas      bool    `json:"isPostEidas"`
	IsValid          bool    `json:"isValid"`
	AcceptableAfter  bool    `json:"acceptableAfter"`
	AcceptableBefore bool    `json:"acceptableBefore"`
	IsCaQc           bool    `json:"isCaQc"`
	IsQTST           bool    `json:"isQTST"`
	IsQEAA           bool    `json:"isQEAA"`
}

// qualTrustProbeDates mirrors QualTrustOracle.PROBE_DATES, in its order.
var qualTrustProbeDates = []*time.Time{
	nil,
	msPtr(1370044800000),
	msPtr(1464739200000),
	msPtr(1514764800000),
	msPtr(1577836800000),
}

func msPtr(millis int64) *time.Time {
	t := time.UnixMilli(millis).UTC()
	return &t
}

func str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func readQualTrustOracle(t *testing.T) []*qualTrustRow {
	t.Helper()
	f, err := os.Open(corpustest.Path(t, "oracle/qual_trust.jsonl"))
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer f.Close()
	var rows []*qualTrustRow
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		row := &qualTrustRow{}
		if err := json.Unmarshal(line, row); err != nil {
			t.Fatalf("decode oracle row: %v", err)
		}
		rows = append(rows, row)
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan oracle: %v", err)
	}
	return rows
}

// buildQualTrustService rebuilds the TrustServiceWrapper the Java driver built
// for this row.
func buildQualTrustService(in *qualTrustInput) *diagnostic.TrustServiceWrapper {
	w := diagnostic.NewTrustServiceWrapper()
	w.Status = str(in.Status)
	w.Type = str(in.Type)
	w.StartDate = millisPtr(in.StartDate)
	w.EndDate = millisPtr(in.EndDate)
	if in.Qualifiers != nil {
		qualifiers := make([]*diagnosticjaxb.XmlQualifier, 0, len(in.Qualifiers))
		for _, uri := range in.Qualifiers {
			qualifiers = append(qualifiers, &diagnosticjaxb.XmlQualifier{Value: uri})
		}
		w.CapturedQualifiers = qualifiers
	}
	if in.Asis != nil {
		w.AdditionalServiceInfos = append([]string{}, in.Asis...)
	}
	w.CountryCode = str(in.CountryCode)
	tl := &diagnosticjaxb.XmlTrustedList{}
	if in.TlUrl != nil {
		url := *in.TlUrl
		tl.Url = &url
	}
	w.TrustedList = tl
	w.EnactedMRA = in.EnactedMRA
	w.MraTrustServiceEquivalenceStatusStartingTime = millisPtr(in.MraStart)
	w.MraTrustServiceEquivalenceStatusEndingTime = millisPtr(in.MraEnd)
	return w
}

func millisPtr(millis *int64) *time.Time {
	if millis == nil {
		return nil
	}
	t := time.UnixMilli(*millis).UTC()
	return &t
}

func accepts(filter TrustServiceFilter, single []*diagnostic.TrustServiceWrapper) bool {
	return len(filter.Filter(single)) > 0
}

func TestQualTrustOracle(t *testing.T) {
	rows := readQualTrustOracle(t)
	if len(rows) == 0 {
		t.Fatal("empty qualification trust-service oracle corpus")
	}
	seen := map[string]int{}
	for _, row := range rows {
		seen[row.Kind]++
	}
	// Every kind the driver emits must be present; a silently-dropped kind is
	// exactly the failure mode this assertion exists to catch.
	for _, kind := range []string{"service", "certQualificationMatrix", "sigQualificationMatrix",
		"finalCertQualification", "eidas", "status", "statusProbe", "serviceType", "serviceTypeProbe"} {
		if seen[kind] == 0 {
			t.Fatalf("oracle carries no %q rows", kind)
		}
	}

	for _, row := range rows {
		row := row
		switch row.Kind {
		case "service":
			t.Run(fmt.Sprintf("service/%d", row.I), func(t *testing.T) { assertQualTrustService(t, row) })
		case "certQualificationMatrix":
			t.Run(fmt.Sprintf("certQualificationMatrix/%s/%s/%s", row.QC, str(row.Type), row.QSCD), func(t *testing.T) {
				got := GetCertQualification(enumerations.CertificateQualifiedStatus(row.QC),
					enumerations.CertificateType(str(row.Type)), enumerations.QSCDStatus(row.QSCD))
				if string(got) != row.Out {
					t.Fatalf("GetCertQualification = %q, want %q", got, row.Out)
				}
			})
		case "sigQualificationMatrix":
			t.Run(fmt.Sprintf("sigQualificationMatrix/%s/%s", row.Indication, row.CertQualification), func(t *testing.T) {
				got, panicked := sigQualificationOrPanic(enumerations.Indication(row.Indication),
					enumerations.CertificateQualification(row.CertQualification))
				if panicked {
					if len(row.Out) < 7 || row.Out[:7] != "throws:" {
						t.Fatalf("Go panicked, Java returned %q", row.Out)
					}
					return
				}
				if len(row.Out) >= 7 && row.Out[:7] == "throws:" {
					t.Fatalf("Go returned %q, Java threw %q", got, row.Out)
				}
				if string(got) != row.Out {
					t.Fatalf("SigQualificationMatrixGetSignatureQualification = %q, want %q", got, row.Out)
				}
			})
		case "finalCertQualification":
			t.Run(fmt.Sprintf("finalCertQualification/%s/%s", row.AtIssuance, row.AtSigning), func(t *testing.T) {
				got, panicked := finalQualificationOrPanic(enumerations.CertificateQualification(row.AtIssuance),
					enumerations.CertificateQualification(row.AtSigning))
				if panicked {
					if len(row.Out) < 7 || row.Out[:7] != "throws:" {
						t.Fatalf("Go panicked, Java returned %q", row.Out)
					}
					return
				}
				if len(row.Out) >= 7 && row.Out[:7] == "throws:" {
					t.Fatalf("Go returned %q, Java threw %q", got, row.Out)
				}
				if string(got) != row.Out {
					t.Fatalf("FinalQualification = %q, want %q", got, row.Out)
				}
			})
		case "eidas":
			t.Run(fmt.Sprintf("eidas/%v", row.Millis), func(t *testing.T) {
				d := millisPtr(row.Millis)
				if got := IsPostEIDAS(d); got != row.IsPostEIDAS {
					t.Errorf("IsPostEIDAS = %v, want %v", got, row.IsPostEIDAS)
				}
				if got := IsPreEIDAS(d); got != row.IsPreEIDAS {
					t.Errorf("IsPreEIDAS = %v, want %v", got, row.IsPreEIDAS)
				}
				if got := IsPostGracePeriod(d); got != row.IsPostGracePerio {
					t.Errorf("IsPostGracePeriod = %v, want %v", got, row.IsPostGracePerio)
				}
			})
		case "status":
			t.Run("status/"+row.Name, func(t *testing.T) {
				s := TrustServiceStatus(row.Name)
				if got := s.URI(); got != str(row.URI) {
					t.Errorf("URI = %q, want %q", got, str(row.URI))
				}
				if got := s.ShortName(); got != row.ShortName {
					t.Errorf("ShortName = %q, want %q", got, row.ShortName)
				}
				if got := s.IsPostEidas(); got != row.IsPostEidas {
					t.Errorf("IsPostEidas = %v, want %v", got, row.IsPostEidas)
				}
				if got := s.IsValid(); got != row.IsValid {
					t.Errorf("IsValid = %v, want %v", got, row.IsValid)
				}
			})
		case "statusProbe":
			t.Run("statusProbe/"+str(row.URI), func(t *testing.T) {
				if got := TrustServiceStatusIsAcceptableStatusAfterEIDAS(str(row.URI)); got != row.AcceptableAfter {
					t.Errorf("IsAcceptableStatusAfterEIDAS(%q) = %v, want %v", str(row.URI), got, row.AcceptableAfter)
				}
				if got := TrustServiceStatusIsAcceptableStatusBeforeEIDAS(str(row.URI)); got != row.AcceptableBefore {
					t.Errorf("IsAcceptableStatusBeforeEIDAS(%q) = %v, want %v", str(row.URI), got, row.AcceptableBefore)
				}
			})
		case "serviceType":
			t.Run("serviceType/"+row.Name, func(t *testing.T) {
				if got := ServiceTypeIdentifier(row.Name).URI(); got != str(row.URI) {
					t.Errorf("URI = %q, want %q", got, str(row.URI))
				}
			})
		case "serviceTypeProbe":
			t.Run("serviceTypeProbe/"+str(row.URI), func(t *testing.T) {
				if got := ServiceTypeIdentifierIsCaQc(str(row.URI)); got != row.IsCaQc {
					t.Errorf("IsCaQc = %v, want %v", got, row.IsCaQc)
				}
				if got := ServiceTypeIdentifierIsQTST(str(row.URI)); got != row.IsQTST {
					t.Errorf("IsQTST = %v, want %v", got, row.IsQTST)
				}
				if got := ServiceTypeIdentifierIsQEAA(str(row.URI)); got != row.IsQEAA {
					t.Errorf("IsQEAA = %v, want %v", got, row.IsQEAA)
				}
			})
		default:
			t.Fatalf("unknown oracle row kind %q", row.Kind)
		}
	}
}

func sigQualificationOrPanic(indication enumerations.Indication,
	cq enumerations.CertificateQualification) (out enumerations.SignatureQualification, panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	out = SigQualificationMatrixGetSignatureQualification(indication, cq)
	return
}

func finalQualificationOrPanic(atIssuance,
	atSigning enumerations.CertificateQualification) (out enumerations.CertificateQualification, panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	out = NewFinalCertificateQualificationCalculator(atIssuance, atSigning).FinalQualification()
	return
}

func assertQualTrustService(t *testing.T, row *qualTrustRow) {
	t.Helper()
	w := buildQualTrustService(row.In)
	single := []*diagnostic.TrustServiceWrapper{w}

	// the ten TrustServiceChecker predicates
	for _, tc := range []struct {
		name string
		got  bool
		want bool
	}{
		{"isLegalPersonConsistent", TrustServiceCheckerIsLegalPersonConsistent(w), row.Checker.LegalPerson},
		{"isQCStatementConsistent", TrustServiceCheckerIsQCStatementConsistent(w), row.Checker.QCStatement},
		{"isQSCDConsistent", TrustServiceCheckerIsQSCDConsistent(w), row.Checker.QSCD},
		{"isQSCDStatusAsInCertConsistent", TrustServiceCheckerIsQSCDStatusAsInCertConsistent(w), row.Checker.QSCDStatusAsInCert},
		{"isPostEIDASQSCDConsistent", TrustServiceCheckerIsPostEIDASQSCDConsistent(w), row.Checker.PostEidasQSCD},
		{"isQualifiersListKnownConsistent", TrustServiceCheckerIsQualifiersListKnownConsistent(w), row.Checker.QualifiersListKnown},
		{"isUsageConsistent", TrustServiceCheckerIsUsageConsistent(w), row.Checker.Usage},
		{"isPreEIDASStatusConsistent", TrustServiceCheckerIsPreEIDASStatusConsistent(w), row.Checker.PreEidasStatus},
		{"isPreEIDASQualifierAndAdditionalServiceInfoConsistent",
			TrustServiceCheckerIsPreEIDASQualifierAndAdditionalServiceInfoConsistent(w), row.Checker.PreEidasQualifierAsi},
		{"isQualifierAndAdditionalServiceInfoConsistent",
			TrustServiceCheckerIsQualifierAndAdditionalServiceInfoConsistent(w), row.Checker.QualifierAsi},
	} {
		if tc.got != tc.want {
			t.Errorf("TrustServiceChecker.%s = %v, want %v", tc.name, tc.got, tc.want)
		}
	}

	// the date-independent filters
	for _, tc := range []struct {
		name string
		got  bool
		want bool
	}{
		{"createFilterByGranted", accepts(TrustServicesFilterFactoryCreateFilterByGranted(), single), row.Filters.Granted},
		{"createFilterByCaQc", accepts(TrustServicesFilterFactoryCreateFilterByCaQc(), single), row.Filters.CaQc},
		{"createFilterByQTST", accepts(TrustServicesFilterFactoryCreateFilterByQTST(), single), row.Filters.QTST},
		{"createFilterByQEAA", accepts(TrustServicesFilterFactoryCreateFilterByQEAA(), single), row.Filters.QEAA},
		{"createConsistentServiceByStatusFilter",
			accepts(TrustServicesFilterFactoryCreateConsistentServiceByStatusFilter(), single), row.Filters.ConsistentStatus},
		{"createConsistentServiceByQCFilter",
			accepts(TrustServicesFilterFactoryCreateConsistentServiceByQCFilter(), single), row.Filters.ConsistentQC},
		{"createConsistentServiceByQSCDFilter",
			accepts(TrustServicesFilterFactoryCreateConsistentServiceByQSCDFilter(), single), row.Filters.ConsistentQSCD},
		{"createConsistentServiceByCertificateTypeFilter",
			accepts(TrustServicesFilterFactoryCreateConsistentServiceByCertificateTypeFilter(), single), row.Filters.ConsistentCertType},
		{"createMRAEnactedFilter", accepts(TrustServicesFilterFactoryCreateMRAEnactedFilter(), single), row.Filters.MraEnacted},
		{"createFilterByCountry", accepts(TrustServicesFilterFactoryCreateFilterByCountry("LU"), single), row.Filters.CountryLU},
		{"createFilterByCountries", accepts(TrustServicesFilterFactoryCreateFilterByCountries(
			map[string]struct{}{"LU": {}, "DE": {}}), single), row.Filters.CountriesLUDE},
		{"createFilterByUrls", accepts(TrustServicesFilterFactoryCreateFilterByUrls(
			map[string]struct{}{"https://tl.example/LU": {}, "https://tl.example/FR": {}}), single), row.Filters.Urls},
	} {
		if tc.got != tc.want {
			t.Errorf("TrustServicesFilterFactory.%s accepted = %v, want %v", tc.name, tc.got, tc.want)
		}
	}

	// the date-sensitive filters, at each probe date
	if len(row.ByDate) != len(qualTrustProbeDates) {
		t.Fatalf("byDate has %d entries, want %d", len(row.ByDate), len(qualTrustProbeDates))
	}
	for i, want := range row.ByDate {
		if got := accepts(TrustServicesFilterFactoryCreateFilterByDate(qualTrustProbeDates[i]), single); got != want {
			t.Errorf("createFilterByDate(probe %d) accepted = %v, want %v", i, got, want)
		}
	}
	if len(row.ByMraEquivalenceStartingDate) != len(qualTrustProbeDates) {
		t.Fatalf("byMraEquivalenceStartingDate has %d entries, want %d",
			len(row.ByMraEquivalenceStartingDate), len(qualTrustProbeDates))
	}
	for i, want := range row.ByMraEquivalenceStartingDate {
		got := accepts(TrustServicesFilterFactoryCreateFilterByMRAEquivalenceStartingDate(qualTrustProbeDates[i]), single)
		if got != want {
			t.Errorf("createFilterByMRAEquivalenceStartingDate(probe %d) accepted = %v, want %v", i, got, want)
		}
	}

	// QualificationByTL
	if len(row.QualificationByTL) == 0 {
		t.Fatal("row carries no qualificationByTL entries")
	}
	for key, want := range row.QualificationByTL {
		got := CreateQualificationFromTL(w, stubQualification(enumerations.CertificateQualifiedStatus(key))).QualifiedStatus()
		if string(got) != want {
			t.Errorf("QualificationByTL(%s) = %q, want %q", key, got, want)
		}
	}

	// TypeByTL
	if len(row.TypeByTL) == 0 {
		t.Fatal("row carries no typeByTL entries")
	}
	for key, want := range row.TypeByTL {
		qualified, inCert := splitPair(t, key)
		got := CreateTypeFromTL(w, enumerations.CertificateQualifiedStatus(qualified),
			stubType(enumerations.CertificateType(inCert))).Type()
		if string(got) != want {
			t.Errorf("TypeByTL(%s) = %q, want %q", key, got, want)
		}
	}

	// QSCDByTL
	if len(row.QSCDByTL) == 0 {
		t.Fatal("row carries no qscdByTL entries")
	}
	for key, want := range row.QSCDByTL {
		qualified, inCert := splitPair(t, key)
		got := CreateQSCDFromTL(w, enumerations.CertificateQualifiedStatus(qualified),
			stubQSCD(enumerations.QSCDStatus(inCert))).QSCDStatus()
		if string(got) != want {
			t.Errorf("QSCDByTL(%s) = %q, want %q", key, got, want)
		}
	}
}

func splitPair(t *testing.T, key string) (string, string) {
	t.Helper()
	for i := 0; i < len(key); i++ {
		if key[i] == '|' {
			return key[:i], key[i+1:]
		}
	}
	t.Fatalf("malformed oracle key %q", key)
	return "", ""
}

type stubQualification enumerations.CertificateQualifiedStatus

func (s stubQualification) QualifiedStatus() enumerations.CertificateQualifiedStatus {
	return enumerations.CertificateQualifiedStatus(s)
}

type stubType enumerations.CertificateType

func (s stubType) Type() enumerations.CertificateType { return enumerations.CertificateType(s) }

type stubQSCD enumerations.QSCDStatus

func (s stubQSCD) QSCDStatus() enumerations.QSCDStatus { return enumerations.QSCDStatus(s) }
