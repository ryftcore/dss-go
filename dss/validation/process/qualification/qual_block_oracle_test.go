package qualification

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
)

// The block KAT for the qualification package. Every row of
// testdata/oracle/qual_block.jsonl is the whole
// XmlValidationCertificateQualification upstream's CertQualificationAtTimeBlock
// produces for one (certificate, trust-service list, validation time) triple -
// see testdata/gen/QualBlockOracle.java.
//
// This is the block that carries the normative TrustServiceFilter ORDER of
// TS 119 615, so an out-of-order filter shows up as a differently ordered (or
// differently sized) Constraint list, and a wrong check wiring shows up as a
// different message key. The final CertificateQualification and the block's
// getFilteredServices() outcome are compared too.

const qualBlockOracleLog = "testdata/oracle/qual_block.jsonl"

type qualBlockConstraint struct {
	Name           *string `json:"name"`
	Status         *string `json:"status"`
	Error          *string `json:"error"`
	Warning        *string `json:"warning"`
	Info           *string `json:"info"`
	AdditionalInfo *string `json:"additionalInfo"`
	Id             *string `json:"id"`
}

type qualBlockConclusion struct {
	Indication    *string  `json:"indication"`
	SubIndication *string  `json:"subIndication"`
	Errors        []string `json:"errors"`
	Warnings      []string `json:"warnings"`
	Infos         []string `json:"infos"`
}

type qualBlockResult struct {
	Title                    *string                `json:"title"`
	CertificateQualification *string                `json:"certificateQualification"`
	ValidationTime           *string                `json:"validationTime"`
	Constraints              []*qualBlockConstraint `json:"constraints"`
	Conclusion               *qualBlockConclusion   `json:"conclusion"`
}

type qualBlockRow struct {
	Kind             string           `json:"kind"`
	I                int              `json:"i"`
	Cert             string           `json:"cert"`
	Services         []string         `json:"services"`
	Date             *int64           `json:"date"`
	ValidationTime   string           `json:"validationTime"`
	Result           *qualBlockResult `json:"result"`
	FilteredServices []*string        `json:"filteredServices"`
}

const (
	qualBlockPreEIDAS  int64 = 1370044800000 // 2013-06-01
	qualBlockPostEIDAS int64 = 1514764800000 // 2018-01-01
	qualBlockEarly     int64 = 1420070400000 // 2015-01-01
)

// qualBlockCertificate mirrors QualBlockOracle.certificate(String).
func qualBlockCertificate(label string) *diagnostic.CertificateWrapper {
	post := len(label) < 4 || label[len(label)-4:] != "-pre"
	notBefore := qualBlockPreEIDAS
	if post {
		notBefore = qualBlockPostEIDAS
	}
	compliance := len(label) >= 3 && label[:3] == "qc-"
	sscd := contains(label, "qscd")
	var types []string
	switch {
	case contains(label, "multi-type"):
		types = []string{"0.4.0.1862.1.6.1", "0.4.0.1862.1.6.2"}
	case contains(label, "eseal"):
		types = []string{"0.4.0.1862.1.6.2"}
	case contains(label, "wsa"):
		types = []string{"0.4.0.1862.1.6.3"}
	case contains(label, "esign"):
		types = []string{"0.4.0.1862.1.6.1"}
	default:
		types = []string{}
	}
	cert := buildQualCert(&qualCertInput{
		QcCompliance:        &compliance,
		QcSSCD:              &sscd,
		QcTypes:             types,
		NotBefore:           notBefore,
		QcStatementsPresent: true,
	})
	return cert
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// qualBlockServices mirrors QualBlockOracle.SERVICES, by name, in its order.
var qualBlockServices = map[string]func() *diagnostic.TrustServiceWrapper{}

func init() {
	const (
		granted          = "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted"
		withdrawn        = "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/withdrawn"
		underSupervision = "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/undersupervision"
		caQc             = "http://uri.etsi.org/TrstSvc/Svctype/CA/QC"
		caPkc            = "http://uri.etsi.org/TrstSvc/Svctype/CA/PKC"
		asiEsig          = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSignatures"
		asiEseal         = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSeals"
		qQcStatement     = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCStatement"
		qNotQualified    = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotQualified"
		qWithQSCD        = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD"
		qNoQSCD          = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoQSCD"
		qForESig         = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESig"
		qForESeal        = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESeal"
	)
	type svc struct {
		name       string
		status     string
		svcType    string
		startDate  int64
		endDate    *int64
		qualifiers []string
		asis       []string
		mra        bool
	}
	early := qualBlockEarly
	for _, s := range []svc{
		{"granted-caqc", granted, caQc, qualBlockPostEIDAS, nil, nil, []string{asiEsig}, false},
		{"granted-caqc-qcstatement", granted, caQc, qualBlockPostEIDAS, nil, []string{qQcStatement}, []string{asiEsig}, false},
		{"granted-caqc-notqualified", granted, caQc, qualBlockPostEIDAS, nil, []string{qNotQualified}, []string{asiEsig}, false},
		{"granted-caqc-withqscd", granted, caQc, qualBlockPostEIDAS, nil, []string{qWithQSCD}, []string{asiEsig}, false},
		{"granted-caqc-eseal", granted, caQc, qualBlockPostEIDAS, nil, []string{qForESeal}, []string{asiEseal}, false},
		{"granted-caqc-inconsistent-usage", granted, caQc, qualBlockPostEIDAS, nil,
			[]string{qForESig, qForESeal}, []string{asiEsig}, false},
		{"granted-caqc-inconsistent-qscd", granted, caQc, qualBlockPostEIDAS, nil,
			[]string{qWithQSCD, qNoQSCD}, []string{asiEsig}, false},
		{"withdrawn-caqc", withdrawn, caQc, qualBlockPostEIDAS, nil, nil, []string{asiEsig}, false},
		{"granted-capkc", granted, caPkc, qualBlockPostEIDAS, nil, nil, []string{asiEsig}, false},
		{"undersupervision-caqc-pre", underSupervision, caQc, qualBlockPreEIDAS, nil, nil, nil, false},
		{"granted-caqc-expired", granted, caQc, qualBlockPreEIDAS, &early, nil, []string{asiEsig}, false},
		{"granted-caqc-mra", granted, caQc, qualBlockPostEIDAS, nil, nil, []string{asiEsig}, true},
	} {
		s := s
		qualBlockServices[s.name] = func() *diagnostic.TrustServiceWrapper {
			w := diagnostic.NewTrustServiceWrapper()
			w.ServiceNames = []string{s.name}
			w.Status = s.status
			w.Type = s.svcType
			start := time.UnixMilli(s.startDate).UTC()
			w.StartDate = &start
			if s.endDate != nil {
				end := time.UnixMilli(*s.endDate).UTC()
				w.EndDate = &end
			}
			qualifiers := make([]*diagnosticjaxb.XmlQualifier, 0, len(s.qualifiers))
			for _, uri := range s.qualifiers {
				qualifiers = append(qualifiers, &diagnosticjaxb.XmlQualifier{Value: uri})
			}
			w.CapturedQualifiers = qualifiers
			w.AdditionalServiceInfos = append([]string{}, s.asis...)
			w.CountryCode = "LU"
			url := "https://tl.example/LU"
			tl := &diagnosticjaxb.XmlTrustedList{}
			tl.Url = &url
			if s.mra {
				mra := true
				tl.Mra = &mra
				w.EnactedMRA = &mra
				mraStart := time.UnixMilli(qualBlockPostEIDAS).UTC()
				w.MraTrustServiceEquivalenceStatusStartingTime = &mraStart
			}
			w.TrustedList = tl
			return w
		}
	}
}

// javaStatusToGo maps the Java XmlStatus enum CONSTANT NAME (what the oracle
// dumps) to the lexical value the Go port carries.
var javaStatusToGo = map[string]jaxb.XmlStatus{
	"OK":          jaxb.XmlStatusOK,
	"NOT_OK":      jaxb.XmlStatusNotOK,
	"IGNORED":     jaxb.XmlStatusIgnored,
	"INFORMATION": jaxb.XmlStatusInformation,
	"WARNING":     jaxb.XmlStatusWarning,
}

func readQualBlockOracle(t *testing.T) []*qualBlockRow {
	t.Helper()
	f, err := os.Open(corpustest.Path(t, "oracle/qual_block.jsonl"))
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer f.Close()
	var rows []*qualBlockRow
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for sc.Scan() {
		if len(sc.Bytes()) == 0 {
			continue
		}
		row := &qualBlockRow{}
		if err := json.Unmarshal(sc.Bytes(), row); err != nil {
			t.Fatalf("decode oracle row: %v", err)
		}
		rows = append(rows, row)
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan oracle: %v", err)
	}
	return rows
}

func TestQualBlockOracle(t *testing.T) {
	rows := readQualBlockOracle(t)
	if len(rows) == 0 {
		t.Fatal("empty qualification block oracle corpus")
	}
	provider := i18n.NewI18nProvider()
	for _, row := range rows {
		row := row
		if row.Kind != "certQualAtTime" {
			t.Fatalf("unknown oracle row kind %q", row.Kind)
		}
		if row.Result == nil {
			t.Fatalf("row %d: upstream threw; the Go port has no equivalent recorded", row.I)
		}
		t.Run(fmt.Sprintf("%d/%s/%s/%s", row.I, row.Cert, row.ValidationTime, joinLabels(row.Services)), func(t *testing.T) {
			cert := qualBlockCertificate(row.Cert)
			services := make([]*diagnostic.TrustServiceWrapper, 0, len(row.Services))
			for _, name := range row.Services {
				build, ok := qualBlockServices[name]
				if !ok {
					t.Fatalf("oracle names unknown trust service %q", name)
				}
				services = append(services, build())
			}

			validationTime := enumerations.ValidationTime(row.ValidationTime)
			var block *CertQualificationAtTimeBlock
			if validationTime == enumerations.ValidationTimeCertificateIssuanceTime {
				block = NewCertQualificationAtTimeBlockAtIssuanceTime(provider, validationTime, cert, services)
			} else {
				block = NewCertQualificationAtTimeBlock(provider, validationTime, millisPtr(row.Date), cert, services)
			}
			result := block.Execute()

			want := row.Result
			if want.Title != nil && result.Title != *want.Title {
				t.Errorf("Title = %q, want %q", result.Title, *want.Title)
			}
			gotQualification := ""
			if result.CertificateQualification != nil {
				gotQualification = string(result.CertificateQualification.CertificateQualification())
			}
			wantQualification := ""
			if want.CertificateQualification != nil {
				wantQualification = *want.CertificateQualification
			}
			if gotQualification != wantQualification {
				t.Errorf("CertificateQualification = %q, want %q", gotQualification, wantQualification)
			}
			gotValidationTime := ""
			if result.ValidationTime != nil {
				gotValidationTime = string(result.ValidationTime.ValidationTime())
			}
			wantValidationTime := ""
			if want.ValidationTime != nil {
				wantValidationTime = *want.ValidationTime
			}
			if gotValidationTime != wantValidationTime {
				t.Errorf("ValidationTime = %q, want %q", gotValidationTime, wantValidationTime)
			}

			assertQualBlockConstraints(t, result.Constraint, want.Constraints)
			assertQualBlockConclusion(t, result.Conclusion, want.Conclusion)

			// getFilteredServices(): upstream never returns null after execute().
			filtered := block.FilteredServices()
			if row.FilteredServices == nil {
				t.Fatalf("oracle recorded no filteredServices for a row that did not throw")
			}
			if len(filtered) != len(row.FilteredServices) {
				t.Fatalf("FilteredServices() = %d services, want %d", len(filtered), len(row.FilteredServices))
			}
			for i, w := range row.FilteredServices {
				var got *string
				if names := filtered[i].ServiceNames; len(names) > 0 {
					got = &names[0]
				}
				switch {
				case got == nil && w == nil:
				case got == nil || w == nil:
					t.Errorf("FilteredServices()[%d] = %v, want %v", i, got, w)
				case *got != *w:
					t.Errorf("FilteredServices()[%d] = %q, want %q", i, *got, *w)
				}
			}
		})
	}
}

func joinLabels(labels []string) string {
	if len(labels) == 0 {
		return "none"
	}
	out := labels[0]
	for _, l := range labels[1:] {
		out += "+" + l
	}
	return out
}

func assertQualBlockConstraints(t *testing.T, got []*jaxb.XmlConstraint, want []*qualBlockConstraint) {
	t.Helper()
	if len(got) != len(want) {
		gotNames := make([]string, 0, len(got))
		for _, c := range got {
			gotNames = append(gotNames, keyOf(c.Name))
		}
		wantNames := make([]string, 0, len(want))
		for _, c := range want {
			wantNames = append(wantNames, strOrEmpty(c.Name))
		}
		t.Fatalf("constraint count = %d %v, want %d %v", len(got), gotNames, len(want), wantNames)
	}
	for i := range want {
		g, w := got[i], want[i]
		if keyOf(g.Name) != strOrEmpty(w.Name) {
			t.Errorf("constraint[%d].Name = %q, want %q", i, keyOf(g.Name), strOrEmpty(w.Name))
		}
		wantStatus := jaxb.XmlStatus("")
		if w.Status != nil {
			s, ok := javaStatusToGo[*w.Status]
			if !ok {
				t.Fatalf("oracle carries unknown status %q", *w.Status)
			}
			wantStatus = s
		}
		if g.Status != wantStatus {
			t.Errorf("constraint[%d].Status = %q, want %q", i, g.Status, wantStatus)
		}
		if keyOf(g.Error) != strOrEmpty(w.Error) {
			t.Errorf("constraint[%d].Error = %q, want %q", i, keyOf(g.Error), strOrEmpty(w.Error))
		}
		if keyOf(g.Warning) != strOrEmpty(w.Warning) {
			t.Errorf("constraint[%d].Warning = %q, want %q", i, keyOf(g.Warning), strOrEmpty(w.Warning))
		}
		if keyOf(g.Info) != strOrEmpty(w.Info) {
			t.Errorf("constraint[%d].Info = %q, want %q", i, keyOf(g.Info), strOrEmpty(w.Info))
		}
		if keyOf(g.Name) == "QUAL_HAS_CONF" {
			// IsNoQualificationConflictDetectedCheck's
			// additional info renders Java's Set<CertificateQualification>, a
			// HashSet of a plain enum: its iteration order is JVM
			// identity-hash-bucket order (the corpus carries both
			// "[CERT_FOR_ESIG, QCERT_FOR_ESIG]" and
			// "[QCERT_FOR_ESIG_QSCD, CERT_FOR_ESIG]", i.e. neither insertion
			// nor reverse-insertion order), which no Go port can reproduce and
			// which upstream itself does not guarantee across JVMs. The Go port
			// renders the same values in first-seen order. The SET of values -
			// the only thing the check's own process() looks at, and the only
			// thing a reader of the report can act on - is compared exactly.
			assertResultsSetEqual(t, i, strOrEmpty(g.AdditionalInfo), strOrEmpty(w.AdditionalInfo))
		} else if strOrEmpty(g.AdditionalInfo) != strOrEmpty(w.AdditionalInfo) {
			t.Errorf("constraint[%d].AdditionalInfo = %q, want %q",
				i, strOrEmpty(g.AdditionalInfo), strOrEmpty(w.AdditionalInfo))
		}
		if strOrEmpty(g.Id) != strOrEmpty(w.Id) {
			t.Errorf("constraint[%d].Id = %q, want %q", i, strOrEmpty(g.Id), strOrEmpty(w.Id))
		}
	}
}

func assertQualBlockConclusion(t *testing.T, got *jaxb.XmlConclusion, want *qualBlockConclusion) {
	t.Helper()
	if want == nil {
		if got != nil {
			t.Errorf("Conclusion = %+v, want nil", got)
		}
		return
	}
	if got == nil {
		t.Fatalf("Conclusion = nil, want %+v", want)
	}
	if string(got.Indication.Indication()) != strOrEmpty(want.Indication) {
		t.Errorf("Conclusion.Indication = %q, want %q", got.Indication.Indication(), strOrEmpty(want.Indication))
	}
	gotSub := ""
	if got.SubIndication != nil {
		gotSub = string(got.SubIndication.SubIndication())
	}
	if gotSub != strOrEmpty(want.SubIndication) {
		t.Errorf("Conclusion.SubIndication = %q, want %q", gotSub, strOrEmpty(want.SubIndication))
	}
	assertMessageKeys(t, "Errors", got.Errors, want.Errors)
	assertMessageKeys(t, "Warnings", got.Warnings, want.Warnings)
	assertMessageKeys(t, "Infos", got.Infos, want.Infos)
}

func assertMessageKeys(t *testing.T, what string, got []*jaxb.XmlMessage, want []string) {
	t.Helper()
	if len(got) != len(want) {
		gotKeys := make([]string, 0, len(got))
		for _, m := range got {
			gotKeys = append(gotKeys, keyOf(m))
		}
		t.Fatalf("Conclusion.%s = %v, want %v", what, gotKeys, want)
	}
	for i := range want {
		if keyOf(got[i]) != want[i] {
			t.Errorf("Conclusion.%s[%d] = %q, want %q", what, i, keyOf(got[i]), want[i])
		}
	}
}

// assertResultsSetEqual compares two "Results : [A, B]" additional-info strings
// as multisets, and everything outside the brackets exactly.
func assertResultsSetEqual(t *testing.T, i int, got, want string) {
	t.Helper()
	gotPrefix, gotValues, gotOK := splitResults(got)
	wantPrefix, wantValues, wantOK := splitResults(want)
	if gotOK != wantOK || gotPrefix != wantPrefix {
		t.Errorf("constraint[%d].AdditionalInfo = %q, want %q", i, got, want)
		return
	}
	if !gotOK {
		if got != want {
			t.Errorf("constraint[%d].AdditionalInfo = %q, want %q", i, got, want)
		}
		return
	}
	gotCount := map[string]int{}
	for _, v := range gotValues {
		gotCount[v]++
	}
	for _, v := range wantValues {
		gotCount[v]--
	}
	for _, n := range gotCount {
		if n != 0 {
			t.Errorf("constraint[%d].AdditionalInfo values = %v, want %v (as a set)", i, gotValues, wantValues)
			return
		}
	}
}

// splitResults splits "<prefix>[a, b, c]" into its prefix and its comma-separated values.
func splitResults(s string) (string, []string, bool) {
	open := strings.IndexByte(s, '[')
	if open < 0 || len(s) == 0 || s[len(s)-1] != ']' {
		return s, nil, false
	}
	inner := s[open+1 : len(s)-1]
	var values []string
	if inner != "" {
		for _, v := range strings.Split(inner, ",") {
			values = append(values, strings.TrimSpace(v))
		}
	}
	return s[:open], values, true
}

func keyOf(m *jaxb.XmlMessage) string {
	if m == nil || m.Key == nil {
		return ""
	}
	return *m.Key
}

func strOrEmpty(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
