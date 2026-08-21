// Phase 9 cross-validation harness, contract item (D): end-to-end trust, Phase 9 -> Phase 8.
//
// Ports eu.europa.esig.dss.tsl.validation.SKCertificateTest#skTLTest (dss-tsl-validation 6.5.RC1)
// verbatim: build a TrustedListsCertificateSource offline from the real Slovak trusted list
// (dss/tsl/testdata/sk-tl-sn-95.xml, already vendored for that package's own sha2/parsing
// coverage - the SAME upstream fixture, unchanged), signed by the real "KCA NBU SR 3" issuer
// certificate; then run Phase 8's CertificateValidator on a real Slovak qualified certificate
// (also transcribed verbatim from the Java test) against that trust source, and compare the
// resulting qualification conclusion end to end.
//
// This is deliberately the CERTIFICATE-qualification boundary (CertificateValidator /
// SimpleCertificateReport), not a full SignedDocumentValidator run over a signed document: it is
// the real, already-vendored upstream scenario that exercises exactly what contract (D) asks for
// - "a [] validated against Go-built vs Java-built trusted sources -> identical qualification
// conclusions" - and every qualification conclusion a signed document's own certificate chain
// would reach is computed by this same qualification engine over this same kind of trust input;
// CertificateValidator IS the mechanism SignedDocumentValidator's own qualification step calls
// into per signing certificate. No document fixture in this repository's corpus chains to KCA NBU
// SR 3 (or to any other already-vendored TL's real service certificates), and authoring a
// synthetic signed document plus a matching synthetic trusted list would exercise a chain of the
// harness's own construction rather than a real one, so this fixture pairing - real TL, real
// certificate, both already in the tree - is the more faithful proof available without expanding
// the vendored fixture set.
//
// testdata/oracle/tsl/certificate_qualification.json is a pure Java dump, produced by
// testdata/oracle/tsl/gen/CertificateQualificationOracle.java against DSS 6.5.RC1's TLValidationJob
// + CertificateValidator run the same way, over the same bytes. Compared: the number of trusted
// certificates synchronized from the TL, the certificate's own Indication/SubIndication, its
// CertificateQualification at issuance, and the number of TrustServiceProviders diagnostic data
// associates with it.
package harness

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"testing"

	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"

	"github.com/utain/esig/dss/internal/corpustest"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	spitsl "github.com/utain/esig/dss/spi/tsl"
	spivalidation "github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/tsl"
	"github.com/utain/esig/dss/validation"
)

// cqTLIssuer is the certificate that signs sk-tl-sn-95.xml ("KCA NBU SR 3" -> "TL and Signature
// Policy List 6"), transcribed verbatim from SKCertificateTest.TL_ISSUER.
const cqTLIssuer = "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"

// cqCertificate is the real Slovak qualified certificate whose issuance qualification is
// cross-validated, transcribed verbatim from SKCertificateTest.CERTIFICATE ("SNCA3").
const cqCertificate = "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"

// cqFileLoader is the offline http.DSSFileLoader mapping the logical "sk-tl.xml" url to the
// vendored fixture's real bytes - mirroring SKCertificateTest's MemoryDataLoader over
// DSSUtils.toByteArray(TL_DOC). path is resolved once, synchronously, by the calling test
// (newCqFileLoader) before OfflineRefresh hands GetDocument to worker goroutines: corpustest's
// t.Skip/t.Fatal are only safe to call from the goroutine running the test, not from those
// workers.
type cqFileLoader struct {
	path string
}

// newCqFileLoader resolves the loader's fixture path up front, in the caller's own goroutine.
func newCqFileLoader(t *testing.T) cqFileLoader {
	t.Helper()
	return cqFileLoader{path: resolveManifestFixture(t, "../tsl/testdata/sk-tl-sn-95.xml")}
}

func (l cqFileLoader) GetDocument(url string) (model.DSSDocument, error) {
	if url != "sk-tl.xml" {
		return nil, model.NewDSSError("no fixture mapped for url " + url)
	}
	return model.NewFileDocument(l.path)
}

type cqDump struct {
	NumberOfTrustedCertificates   int      `json:"numberOfTrustedCertificates"`
	Indication                    string   `json:"indication"`
	SubIndication                 string   `json:"subIndication"`
	QualificationAtIssuance       string   `json:"qualificationAtIssuance"`
	NumberOfTrustServiceProviders int      `json:"numberOfTrustServiceProviders"`
	TrustServiceProviders         []cqTSP  `json:"trustServiceProviders"`
	TrustedLists                  []cqList `json:"trustedLists"`
}

// cqTSP is one XmlTrustServiceProvider of the certificate's diagnostic data: the TL content
// Phase 9 handed to Phase 8 for this certificate.
type cqTSP struct {
	TSPNames                   []string     `json:"tspNames"`
	TSPTradeNames              []string     `json:"tspTradeNames"`
	TSPRegistrationIdentifiers []string     `json:"tspRegistrationIdentifiers"`
	TrustServices              []cqTrustSvc `json:"trustServices"`
}

type cqTrustSvc struct {
	ServiceType         string   `json:"serviceType"`
	Status              string   `json:"status"`
	HasStartDate        bool     `json:"hasStartDate"`
	HasEndDate          bool     `json:"hasEndDate"`
	ServiceNames        []string `json:"serviceNames"`
	ServiceSupplyPoints []string `json:"serviceSupplyPoints"`
	CapturedQualifiers  []string `json:"capturedQualifiers"`
}

// cqList is one XmlTrustedList of the diagnostic data - the block upstream's own JUnit suite
// asserts over (url/countryCode/sequenceNumber/version/dates/wellSigned/structural validation).
type cqList struct {
	URL                string   `json:"url"`
	CountryCode        string   `json:"countryCode"`
	SequenceNumber     int      `json:"sequenceNumber"`
	Version            int      `json:"version"`
	HasLastLoading     bool     `json:"hasLastLoading"`
	HasIssueDate       bool     `json:"hasIssueDate"`
	HasNextUpdate      bool     `json:"hasNextUpdate"`
	WellSigned         bool     `json:"wellSigned"`
	LOTL               bool     `json:"lotl"`
	StructuralValid    bool     `json:"structuralValid"`
	StructuralMessages []string `json:"structuralMessages"`
}

// cqNullInt is the sentinel the Java generator emits for a null Integer, matching
// tsl_parsing_oracle_test.go's tpNullInt.
const cqNullInt = -2147483648

func cqSortedStrings(values []string) []string {
	out := make([]string, 0, len(values))
	out = append(out, values...)
	sort.Strings(out)
	return out
}

// cqLangValues flattens a list of XmlLangAndValue into "lang|value" strings, sorted, the way
// CertificateQualificationOracle.java#dumpSortedLangValues does.
func cqLangValues(values []*diagnosticjaxb.XmlLangAndValue) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		lang := ""
		if v.Lang != nil {
			lang = *v.Lang
		}
		out = append(out, lang+"|"+v.Value)
	}
	sort.Strings(out)
	return out
}

func cqQualifiers(values []*diagnosticjaxb.XmlQualifier) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		critical := false
		if v.Critical != nil {
			critical = *v.Critical
		}
		out = append(out, v.Value+"|"+strconv.FormatBool(critical))
	}
	sort.Strings(out)
	return out
}

func cqDeref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func cqDerefInt(p *int) int {
	if p == nil {
		return cqNullInt
	}
	return *p
}

func cqLoadCertificate(t *testing.T, base64Certificate string) *model.CertificateToken {
	t.Helper()
	cert, err := spi.DSSUtilsLoadCertificateFromBase64EncodedString(base64Certificate)
	if err != nil {
		t.Fatalf("loading certificate: %v", err)
	}
	return cert
}

// TestCertificateQualificationOracle is the Phase 9 harness contract item (D).
func TestCertificateQualificationOracle(t *testing.T) {
	issuerSource := spi.NewCommonCertificateSource()
	issuerSource.AddCertificate(cqLoadCertificate(t, cqTLIssuer))

	tlSource := tsl.NewTLSource()
	tlSource.SetUrl("sk-tl.xml")
	tlSource.SetTLVersions([]int{5, 6})
	tlSource.SetCertificateSource(&issuerSource)

	tlValidationJob := tsl.NewTLValidationJob()
	tlValidationJob.SetTrustedListSources(tlSource)
	tlValidationJob.SetOfflineDataLoader(newCqFileLoader(t))
	trustedCertificateSource := spitsl.NewTrustedListsCertificateSource()
	tlValidationJob.SetTrustedListCertificateSource(trustedCertificateSource)

	if err := tlValidationJob.OfflineRefresh(); err != nil {
		t.Fatalf("OfflineRefresh: %v", err)
	}

	certificate := cqLoadCertificate(t, cqCertificate)

	certificateValidator := validation.CertificateValidatorFromCertificate(certificate)
	certificateVerifier := spivalidation.NewCommonCertificateVerifier()
	certificateVerifier.SetTrustedCertSources(trustedCertificateSource)
	certificateValidator.SetCertificateVerifier(certificateVerifier)

	certificateReports, err := certificateValidator.Validate()
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	simpleReport := certificateReports.GetSimpleReport()
	certID := certificate.DSSIDAsString()

	got := cqDump{
		NumberOfTrustedCertificates: len(trustedCertificateSource.Certificates()),
		Indication:                  string(simpleReport.GetCertificateIndication(certID)),
		SubIndication:               string(simpleReport.GetCertificateSubIndication(certID)),
		QualificationAtIssuance:     string(simpleReport.GetQualificationAtCertificateIssuance()),
	}

	diagnosticData := certificateReports.GetDiagnosticData()
	got.TrustServiceProviders = []cqTSP{}
	if cw := diagnosticData.CertificateById(certID); cw != nil {
		got.NumberOfTrustServiceProviders = len(cw.TrustServiceProviders())
		for _, tsp := range cw.TrustServiceProviders() {
			dumped := cqTSP{
				TSPNames:                   cqLangValues(tsp.TSPNames.All()),
				TSPTradeNames:              cqLangValues(tsp.TSPTradeNames.All()),
				TSPRegistrationIdentifiers: cqSortedStrings(tsp.TSPRegistrationIdentifiers.All()),
				TrustServices:              []cqTrustSvc{},
			}
			for _, ts := range tsp.TrustServices.All() {
				dumped.TrustServices = append(dumped.TrustServices, cqTrustSvc{
					ServiceType:         cqDeref(ts.ServiceType),
					Status:              cqDeref(ts.Status),
					HasStartDate:        ts.StartDate != nil,
					HasEndDate:          ts.EndDate != nil,
					ServiceNames:        cqLangValues(ts.ServiceNames.All()),
					ServiceSupplyPoints: cqSortedStrings(ts.ServiceSupplyPoints.All()),
					CapturedQualifiers:  cqQualifiers(ts.CapturedQualifiers.All()),
				})
			}
			sort.Slice(dumped.TrustServices, func(i, j int) bool {
				a, _ := json.Marshal(dumped.TrustServices[i])
				b, _ := json.Marshal(dumped.TrustServices[j])
				return string(a) < string(b)
			})
			got.TrustServiceProviders = append(got.TrustServiceProviders, dumped)
		}
		// See CertificateQualificationOracle.java: both the provider list and each provider's
		// trust-service list come out of a Java HashMap's entrySet(), so neither order is a
		// parity target; both sides emit them in a canonical (JSON-sorted) order.
		sort.Slice(got.TrustServiceProviders, func(i, j int) bool {
			a, _ := json.Marshal(got.TrustServiceProviders[i])
			b, _ := json.Marshal(got.TrustServiceProviders[j])
			return string(a) < string(b)
		})
	}
	got.TrustedLists = []cqList{}
	for _, tl := range diagnosticData.TrustedLists() {
		dumped := cqList{
			URL:                cqDeref(tl.Url),
			CountryCode:        cqDeref(tl.CountryCode),
			SequenceNumber:     cqDerefInt(tl.SequenceNumber),
			Version:            cqDerefInt(tl.Version),
			HasLastLoading:     tl.LastLoading != nil,
			HasIssueDate:       tl.IssueDate != nil,
			HasNextUpdate:      tl.NextUpdate != nil,
			WellSigned:         tl.WellSigned,
			LOTL:               tl.LOTL != nil && *tl.LOTL,
			StructuralMessages: []string{},
		}
		if tl.StructuralValidation != nil {
			dumped.StructuralValid = tl.StructuralValidation.Valid
			dumped.StructuralMessages = cqSortedStrings(tl.StructuralValidation.Message)
		}
		got.TrustedLists = append(got.TrustedLists, dumped)
	}

	f, err := os.Open(corpustest.Path(t, "oracle/tsl/certificate_qualification.json"))
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("empty oracle file")
	}
	var want cqDump
	if err := json.Unmarshal(scanner.Bytes(), &want); err != nil {
		t.Fatalf("unmarshal oracle: %v", err)
	}

	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal got: %v", err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal want: %v", err)
	}
	if string(gotJSON) != string(wantJSON) {
		t.Errorf("mismatch:\n GOT: %s\nWANT: %s", gotJSON, wantJSON)
	}
}
