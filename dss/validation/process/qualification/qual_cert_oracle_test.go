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
)

// The certificate KAT for the qualification package: every row of
// testdata/oracle/qual_cert.jsonl is what upstream produces for one synthetic
// certificate - see testdata/gen/QualCertOracle.java. It covers
// QCTypeIdentifiers, the six "…ByCertificate…EIDAS" strategies, the full
// CertificateQualificationCalculator over a fixed panel of trust services,
// ServiceByCertificateTypeFilter and UniqueServiceFilter (both its
// "one conclusion" and its "several conclusions" branches).

const qualCertOracleLog = "testdata/oracle/qual_cert.jsonl"

type qualCertInput struct {
	QcCompliance        *bool    `json:"qcCompliance"`
	QcSSCD              *bool    `json:"qcSSCD"`
	QcTypes             []string `json:"qcTypes"`
	Legislations        []string `json:"legislations"`
	Policies            []string `json:"policies"`
	NotBefore           int64    `json:"notBefore"`
	QcStatementsPresent bool     `json:"qcStatementsPresent"`
}

type qualCertTypeIdentifiers struct {
	Esign bool `json:"esign"`
	Eseal bool `json:"eseal"`
	Web   bool `json:"web"`
}

type qualCertFromCert struct {
	Qualification string `json:"qualification"`
	Type          string `json:"type"`
	QSCD          string `json:"qscd"`
}

type qualCertPerService struct {
	Qualification     string `json:"qualification"`
	ByCertificateType *bool  `json:"byCertificateType"`
	QualifiedStatus   string `json:"qualifiedStatus"`
	CertType          string `json:"certType"`
	QSCDStatus        string `json:"qscdStatus"`
}

type qualCertRow struct {
	Kind             string                         `json:"kind"`
	I                int                            `json:"i"`
	In               *qualCertInput                 `json:"in"`
	QCTypeIdentifier *qualCertTypeIdentifiers       `json:"qcTypeIdentifiers"`
	FromCert         *qualCertFromCert              `json:"fromCert"`
	PerService       map[string]*qualCertPerService `json:"perService"`
	UniqueAll        []*string                      `json:"uniqueAll"`
	UniquePairs      [][]*string                    `json:"uniquePairs"`
	UniqueSingleton  []*string                      `json:"uniqueSingleton"`
	UniqueEmpty      []*string                      `json:"uniqueEmpty"`
	ServiceLabels    []string                       `json:"serviceLabels"`
}

// qualCertServices mirrors QualCertOracle.SERVICES, by label, in its order.
// The "none" entry is the nil TrustServiceWrapper.
var qualCertServices = []struct {
	label      string
	status     string
	svcType    string
	startDate  *int64
	qualifiers []string
	asis       []string
}{
	{label: "none"},
	{"granted-post-plain", statusGranted, typeCaQc, msValue(qualCertPostEIDAS), []string{}, []string{}},
	{"granted-post-qcstatement", statusGranted, typeCaQc, msValue(qualCertPostEIDAS),
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCStatement"}, []string{}},
	{"granted-post-notqualified", statusGranted, typeCaQc, msValue(qualCertPostEIDAS),
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotQualified"}, []string{}},
	{"granted-post-withqscd", statusGranted, typeCaQc, msValue(qualCertPostEIDAS),
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD"}, []string{}},
	{"granted-post-noqscd", statusGranted, typeCaQc, msValue(qualCertPostEIDAS),
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCNoQSCD"}, []string{}},
	{"granted-post-qscdasincert", statusGranted, typeCaQc, msValue(qualCertPostEIDAS),
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCQSCDStatusAsInCert"}, []string{}},
	{"granted-post-foreseal", statusGranted, typeCaQc, msValue(qualCertPostEIDAS),
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESeal"},
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSeals"}},
	{"granted-post-forwsa", statusGranted, typeCaQc, msValue(qualCertPostEIDAS),
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForWSA"},
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForWebSiteAuthentication"}},
	{"granted-pre-plain", statusUnderSupervision, typeCaQc, msValue(qualCertPreEIDAS), []string{}, []string{}},
	{"granted-pre-withsscd", statusUnderSupervision, typeCaQc, msValue(qualCertPreEIDAS),
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithSSCD"}, []string{}},
	{"withdrawn-post", statusWithdrawn, typeCaQc, msValue(qualCertPostEIDAS), []string{}, []string{}},
	{"granted-post-esig-asi", statusGranted, typeCaQc, msValue(qualCertPostEIDAS),
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCForESig"},
		[]string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSignatures"}},
}

const (
	qualCertPreEIDAS  int64 = 1370044800000
	qualCertPostEIDAS int64 = 1514764800000

	statusGranted          = "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted"
	statusUnderSupervision = "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/undersupervision"
	statusWithdrawn        = "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/withdrawn"
	typeCaQc               = "http://uri.etsi.org/TrstSvc/Svctype/CA/QC"
)

func msValue(v int64) *int64 { return &v }

func readQualCertOracle(t *testing.T) []*qualCertRow {
	t.Helper()
	f, err := os.Open(qualCertOracleLog)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer f.Close()
	var rows []*qualCertRow
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for sc.Scan() {
		if len(sc.Bytes()) == 0 {
			continue
		}
		row := &qualCertRow{}
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

func buildQualCert(in *qualCertInput) *diagnostic.CertificateWrapper {
	return diagnostic.NewCertificateWrapper(buildQualCertXml(in))
}

// buildQualCertXml builds the raw XmlCertificate, so callers that also need to
// attach trust-service providers (qual_sig_block_oracle_test.go) can do so
// before the wrapper is built.
func buildQualCertXml(in *qualCertInput) *diagnosticjaxb.XmlCertificate {
	notBefore := diagnosticjaxb.XSDateTime(time.UnixMilli(in.NotBefore).UTC())
	notAfter := diagnosticjaxb.XSDateTime(time.UnixMilli(in.NotBefore + 3*365*86400000).UTC())
	cert := &diagnosticjaxb.XmlCertificate{}
	certID := diagnosticjaxb.CollapsedString("C-SYNTH")
	cert.Id = &certID
	cert.NotBefore = &notBefore
	cert.NotAfter = &notAfter

	var extensions []diagnosticjaxb.XmlCertificateExtensionItem
	if in.QcStatementsPresent {
		qcStatements := &diagnosticjaxb.XmlQcStatements{}
		oid := enumerations.CertificateExtensionEnum_QC_STATEMENTS.OID()
		qcStatements.OID = &oid
		if in.QcCompliance != nil {
			qcStatements.QcCompliance = &diagnosticjaxb.XmlQcCompliance{Present: *in.QcCompliance}
		}
		if in.QcSSCD != nil {
			qcStatements.QcSSCD = &diagnosticjaxb.XmlQcSSCD{Present: *in.QcSSCD}
		}
		if in.QcTypes != nil {
			items := make([]*diagnosticjaxb.XmlOID, 0, len(in.QcTypes))
			for _, o := range in.QcTypes {
				item := &diagnosticjaxb.XmlOID{}
				item.Value = o
				items = append(items, item)
			}
			qcStatements.QcTypes = &diagnosticjaxb.QcTypesWrapper{Items: items}
		}
		if in.Legislations != nil {
			qcStatements.QcCClegislation = &diagnosticjaxb.QcCClegislationWrapper{
				Items: append([]string{}, in.Legislations...)}
		}
		extensions = append(extensions, qcStatements)
	}
	if in.Policies != nil {
		policies := &diagnosticjaxb.XmlCertificatePolicies{}
		oid := enumerations.CertificateExtensionEnum_CERTIFICATE_POLICIES.OID()
		policies.OID = &oid
		for _, o := range in.Policies {
			policy := &diagnosticjaxb.XmlCertificatePolicy{}
			policy.Value = o
			policies.CertificatePolicy = append(policies.CertificatePolicy, policy)
		}
		extensions = append(extensions, policies)
	}
	cert.CertificateExtensions = &diagnosticjaxb.CertificateExtensionsWrapper{Items: extensions}
	return cert
}

func buildQualCertService(i int) *diagnostic.TrustServiceWrapper {
	c := qualCertServices[i]
	if c.status == "" && c.svcType == "" && c.startDate == nil {
		return nil
	}
	w := diagnostic.NewTrustServiceWrapper()
	w.Status = c.status
	w.Type = c.svcType
	w.StartDate = millisPtr(c.startDate)
	qualifiers := make([]*diagnosticjaxb.XmlQualifier, 0, len(c.qualifiers))
	for _, uri := range c.qualifiers {
		qualifiers = append(qualifiers, &diagnosticjaxb.XmlQualifier{Value: uri})
	}
	w.CapturedQualifiers = qualifiers
	w.AdditionalServiceInfos = append([]string{}, c.asis...)
	w.CountryCode = "LU"
	url := "https://tl.example/LU"
	tl := &diagnosticjaxb.XmlTrustedList{}
	tl.Url = &url
	w.TrustedList = tl
	w.ServiceNames = []string{c.label}
	return w
}

func TestQualCertOracle(t *testing.T) {
	rows := readQualCertOracle(t)
	if len(rows) == 0 {
		t.Fatal("empty qualification certificate oracle corpus")
	}
	for _, row := range rows {
		row := row
		if row.Kind != "cert" {
			t.Fatalf("unknown oracle row kind %q", row.Kind)
		}
		t.Run(fmt.Sprintf("cert/%d", row.I), func(t *testing.T) {
			cert := buildQualCert(row.In)

			// QCTypeIdentifiers
			if got := IsQCTypeEsign(cert); got != row.QCTypeIdentifier.Esign {
				t.Errorf("IsQCTypeEsign = %v, want %v", got, row.QCTypeIdentifier.Esign)
			}
			if got := IsQCTypeEseal(cert); got != row.QCTypeIdentifier.Eseal {
				t.Errorf("IsQCTypeEseal = %v, want %v", got, row.QCTypeIdentifier.Eseal)
			}
			if got := IsQCTypeWeb(cert); got != row.QCTypeIdentifier.Web {
				t.Errorf("IsQCTypeWeb = %v, want %v", got, row.QCTypeIdentifier.Web)
			}

			// the three "from certificate" strategies (which pick the pre- or
			// post-eIDAS implementation off the certificate's notBefore)
			if got := CreateQualificationFromCert(cert).QualifiedStatus(); string(got) != row.FromCert.Qualification {
				t.Errorf("createQualificationFromCert = %q, want %q", got, row.FromCert.Qualification)
			}
			if got := CreateTypeFromCert(cert).Type(); string(got) != row.FromCert.Type {
				t.Errorf("createTypeFromCert = %q, want %q", got, row.FromCert.Type)
			}
			if got := CreateQSCDFromCert(cert).QSCDStatus(); string(got) != row.FromCert.QSCD {
				t.Errorf("createQSCDFromCert = %q, want %q", got, row.FromCert.QSCD)
			}

			// the whole calculator, per trust service
			if len(row.PerService) != len(qualCertServices) {
				t.Fatalf("perService has %d entries, want %d", len(row.PerService), len(qualCertServices))
			}
			for i := range qualCertServices {
				label := qualCertServices[i].label
				want, ok := row.PerService[label]
				if !ok {
					t.Fatalf("perService carries no %q entry", label)
				}
				svc := buildQualCertService(i)
				got := NewCertificateQualificationCalculator(cert, svc).Qualification()
				if string(got) != want.Qualification {
					t.Errorf("[%s] CertificateQualificationCalculator = %q, want %q", label, got, want.Qualification)
				}
				if want.ByCertificateType != nil {
					filter := TrustServicesFilterFactoryCreateFilterByCertificateType(cert)
					accepted := len(filter.Filter([]*diagnostic.TrustServiceWrapper{svc})) > 0
					if accepted != *want.ByCertificateType {
						t.Errorf("[%s] createFilterByCertificateType accepted = %v, want %v",
							label, accepted, *want.ByCertificateType)
					}
				}
				qualified := CreateQualificationFromCertAndTL(cert, svc).QualifiedStatus()
				if string(qualified) != want.QualifiedStatus {
					t.Errorf("[%s] createQualificationFromCertAndTL = %q, want %q", label, qualified, want.QualifiedStatus)
				}
				if gotType := CreateTypeFromCertAndTL(cert, svc, qualified).Type(); string(gotType) != want.CertType {
					t.Errorf("[%s] createTypeFromCertAndTL = %q, want %q", label, gotType, want.CertType)
				}
				if gotQSCD := CreateQSCDFromCertAndTL(cert, svc, qualified).QSCDStatus(); string(gotQSCD) != want.QSCDStatus {
					t.Errorf("[%s] createQSCDFromCertAndTL = %q, want %q", label, gotQSCD, want.QSCDStatus)
				}
			}

			// UniqueServiceFilter
			var all []*diagnostic.TrustServiceWrapper
			for i := range qualCertServices {
				if svc := buildQualCertService(i); svc != nil {
					all = append(all, svc)
				}
			}
			if len(all) != len(row.ServiceLabels) {
				t.Fatalf("rebuilt %d non-nil services, oracle recorded %d", len(all), len(row.ServiceLabels))
			}
			assertUniqueLabels(t, "uniqueAll", cert, all, row.UniqueAll)
			if len(row.UniquePairs) != len(all)-1 {
				t.Fatalf("uniquePairs has %d entries, want %d", len(row.UniquePairs), len(all)-1)
			}
			for i := 0; i+1 < len(all); i++ {
				assertUniqueLabels(t, fmt.Sprintf("uniquePairs[%d]", i), cert,
					[]*diagnostic.TrustServiceWrapper{all[i], all[i+1]}, row.UniquePairs[i])
			}
			assertUniqueLabels(t, "uniqueSingleton", cert, all[:1], row.UniqueSingleton)
			assertUniqueLabels(t, "uniqueEmpty", cert, nil, row.UniqueEmpty)
		})
	}
}

func assertUniqueLabels(t *testing.T, what string, cert *diagnostic.CertificateWrapper,
	services []*diagnostic.TrustServiceWrapper, want []*string) {
	t.Helper()
	filtered := TrustServicesFilterFactoryCreateUniqueServiceFilter(cert).Filter(services)
	if len(filtered) != len(want) {
		t.Fatalf("%s: UniqueServiceFilter selected %d services, want %d", what, len(filtered), len(want))
	}
	for i, w := range want {
		var got *string
		if names := filtered[i].ServiceNames; len(names) > 0 {
			got = &names[0]
		}
		switch {
		case got == nil && w == nil:
		case got == nil || w == nil:
			t.Errorf("%s[%d]: selected %v, want %v", what, i, got, w)
		case *got != *w:
			t.Errorf("%s[%d]: selected %q, want %q", what, i, *got, *w)
		}
	}
}
