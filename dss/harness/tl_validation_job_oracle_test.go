// Cross-validation harness: offline TLValidationJob parity.
//
// This test runs the SAME TLValidationJob pipeline the Java side runs, entirely offline (a
// caller-supplied file loader stands in for both Java's FileCacheDataLoader/MockDataLoader and
// Go's http.DSSFileLoader - no network access on either side), over a fixed, deterministic set
// of already-vendored fixtures shared with contract item (A)
// (testdata/oracle/tsl/eu-lotl.xml, eu-lotl-broken-sig.xml, eu-lotl-not-parseable.xml, de-tl.xml):
// one LOTLSource (signed by the same certificate tl_validator_task_oracle_test.go's "correctCert"
// case cross-validates as eu-lotl.xml's real signer, reaching Indication TOTAL_PASSED on both
// sides) and three independently-configured TLSources covering three distinct outcomes chosen
// for determinism, not realism: a genuine trust service provider corpus with an EMPTY trust
// anchor (de-tl.xml - parses fine, is not trusted, so its own signature validation is
// INDETERMINATE, but its TSP/TrustService content still reaches TrustedListsCertificateSource,
// since synchronization is gated on PARSING succeeding, not on signature validity - see
// trusted_list_certificate_source_synchronizer.go's isTLParsingDesyncOrError), a broken-signature
// fixture reused from item (A) (eu-lotl-broken-sig.xml, also an empty trust anchor -
// deterministically TOTAL_FAILED/HASH_FAILURE), and a not-parseable fixture also reused from item
// (A) (eu-lotl-not-parseable.xml - deterministically a PARSING cache ERROR, never reaching
// signature validation at all).
//
// testdata/oracle/tsl/tl_validation_job.json is a pure Java dump, produced by
// testdata/oracle/tsl/gen/TLValidationJobOracle.java against DSS 6.5.RC1's TLValidationJob run
// the same way, over the same bytes. Compared: NumberOfProcessedLOTLs/TLs; per LOTL/TL, the
// download/parsing/validation cache StatusName()s and (for validation) Indication/SubIndication;
// which of the four alerts.detections strategies (TLSignatureErrorDetection,
// TLParsingErrorDetection, LOTLLocationChangeDetection, OJUrlChangeDetection) actually FIRED
// during the job's own alerting pass - TLAlert/LOTLAlert wired with a detection strategy and a
// handler that records into this test's own slice, exactly the way Java's real Alert/AlertHandler
// pair does, not re-evaluated after the fact (TLSignatureErrorDetection's precondition,
// downloadCacheInfo.isDesynchronized(), only holds DURING the job's own alerting pass, before
// SynchronizeCertificateSources settles the cache back to SYNCHRONIZED); and the resulting
// TrustedListsCertificateSource content - every trusted certificate's SHA-256, and, per
// certificate, the sorted set of TL URLs, TSP territories and "type|status" service pairs its
// TrustProperties carry.
package harness

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"sort"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	dsspolicy "github.com/ryftcore/dss-go/dss/policy"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
	spitsl "github.com/ryftcore/dss-go/dss/spi/tsl"
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/tsl"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
)

// tvjNoChildTLs is an OtherTSLPointerPredicate matching nothing, used as the LOTLSource's TL
// predicate so the job's own LOTL->child-TL auto-discovery (which would otherwise try to
// offline-fetch every real country TL URL eu-lotl.xml's own pointer list names) never fires;
// this harness's three TLs are independently configured instead, entirely under this test's
// control - see this file's header.
type tvjNoChildTLs struct{}

func (tvjNoChildTLs) Test(*jaxb.OtherTSLPointerType) bool { return false }

// tvjSigningCertificate is the exact certificate tl_validator_task_oracle_test.go's "correctCert"
// case cross-validates as eu-lotl.xml's real signer (Indication TOTAL_PASSED); duplicated here,
// not imported, because it is an unexported test-only constant of package tsl and this harness
// deliberately exercises the public TLValidationJob API from outside that package, the same way
// a composing application would.
const tvjSigningCertificate = "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"

func init() {
	// Same registration tl_validator_task_oracle_test.go's init performs - TLValidatorTask's
	// trustedListValidationPolicy needs a registered ValidationPolicyFactory (see that file's
	// header for why), and TLValidationJob's own validation runnable reaches the same code path.
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
}

// tvjFileLoader is a trivial offline http.DSSFileLoader over a fixed url->local-path map - the Go
// equivalent of Java's FileCacheDataLoader wired to a MockDataLoader; this harness never touches
// the network on either side.
type tvjFileLoader struct {
	files map[string]string
}

var _ http.DSSFileLoader = (*tvjFileLoader)(nil)

func (l *tvjFileLoader) GetDocument(url string) (model.DSSDocument, error) {
	path, ok := l.files[url]
	if !ok {
		return nil, model.NewDSSError("no fixture mapped for url " + url)
	}
	return model.NewFileDocument(path)
}

// tvjCertificateSource loads zero or more base64 certificates into a fresh CommonCertificateSource,
// mirroring tl_validator_task_oracle_test.go's tlValidatorTaskCertificateSource.
func tvjCertificateSource(t *testing.T, base64Certificates ...string) spi.CertificateSource {
	t.Helper()
	source := spi.NewCommonCertificateSource()
	for _, b64 := range base64Certificates {
		cert, err := spi.DSSUtilsLoadCertificateFromBase64EncodedString(b64)
		if err != nil {
			t.Fatalf("loading certificate: %v", err)
		}
		source.AddCertificate(cert)
	}
	return &source
}

// tvjURLRecorder is an alert.AlertHandler[T] for any T carrying a Url() string, recording the url
// instead of logging (matching the shape of the production Log*AlertHandler types it replaces) -
// this harness cares about WHETHER an alert fired, not what a log line says. Generic (rather than
// a single interface-typed Process) because alert.AlertHandler[T]'s Process(T) needs an EXACT
// type match for each concrete D/L instantiation (*tslmodel.TLInfo, *tslmodel.LOTLInfo).
type tvjURLRecorder[T interface{ Url() string }] struct {
	fired *[]string
}

func (h *tvjURLRecorder[T]) Process(object T) error {
	*h.fired = append(*h.fired, object.Url())
	return nil
}

// tvjInfoRecordDump is the comparable shape of one job.InfoRecord (download/parsing/validation).
type tvjInfoRecordDump struct {
	StatusName     string `json:"statusName"`
	Synchronized   bool   `json:"synchronized"`
	Desynchronized bool   `json:"desynchronized"`
	Error          bool   `json:"error"`
	ResultExist    bool   `json:"resultExist"`
}

// tvjValidationDump adds the Indication/SubIndication on top of tvjInfoRecordDump.
type tvjValidationDump struct {
	tvjInfoRecordDump
	Indication    string `json:"indication"`
	SubIndication string `json:"subIndication"`
}

// tvjTLDump is the comparable shape of one TLInfo/LOTLInfo (contract (C): "cache-state
// summaries").
type tvjTLDump struct {
	URL        string            `json:"url"`
	Download   tvjInfoRecordDump `json:"download"`
	Parsing    tvjInfoRecordDump `json:"parsing"`
	Validation tvjValidationDump `json:"validation"`
}

// tvjCertDump is one trusted certificate's comparable content (contract (C): "every cert +
// associated TL properties").
type tvjCertDump struct {
	SHA256            string   `json:"sha256"`
	TLUrls            []string `json:"tlUrls"`
	TSPTerritories    []string `json:"tspTerritories"`
	ServiceTypeStatus []string `json:"serviceTypeStatus"`
}

// tvjDump is the full comparable shape of one offline TLValidationJob run (contract (C) in full).
type tvjDump struct {
	NumberOfProcessedLOTLs int           `json:"numberOfProcessedLOTLs"`
	NumberOfProcessedTLs   int           `json:"numberOfProcessedTLs"`
	LOTLs                  []tvjTLDump   `json:"lotls"`
	TLs                    []tvjTLDump   `json:"tls"`
	FiredAlerts            []string      `json:"firedAlerts"`
	TrustedCertificates    []tvjCertDump `json:"trustedCertificates"`
}

func tvjDumpInfoRecord(r modeljobInfoRecord) tvjInfoRecordDump {
	if r == nil {
		return tvjInfoRecordDump{}
	}
	return tvjInfoRecordDump{
		StatusName:     r.StatusName(),
		Synchronized:   r.IsSynchronized(),
		Desynchronized: r.IsDesynchronized(),
		Error:          r.IsError(),
		ResultExist:    r.IsResultExist(),
	}
}

// modeljobInfoRecord is the subset of model/job.InfoRecord this dump needs; declared locally so
// this file does not need to import model/job (its exported name, "job", would collide with the
// tsl package's own frequent local variable name "job" used below for the TLValidationJob under
// test).
type modeljobInfoRecord interface {
	StatusName() string
	IsSynchronized() bool
	IsDesynchronized() bool
	IsError() bool
	IsResultExist() bool
}

func tvjDumpTLInfo(info *tslmodel.TLInfo) tvjTLDump {
	dump := tvjTLDump{
		URL:      info.Url(),
		Download: tvjDumpInfoRecord(info.DownloadCacheInfo()),
		Parsing:  tvjDumpInfoRecord(info.ParsingCacheInfo()),
	}
	validationCacheInfo := info.ValidationCacheInfo()
	dump.Validation = tvjValidationDump{tvjInfoRecordDump: tvjDumpInfoRecord(validationCacheInfo)}
	if validationCacheInfo != nil {
		dump.Validation.Indication = string(validationCacheInfo.Indication())
		dump.Validation.SubIndication = string(validationCacheInfo.SubIndication())
	}
	return dump
}

func tvjDumpCertSource(source *spitsl.TrustedListsCertificateSource) []tvjCertDump {
	out := []tvjCertDump{}
	for _, cert := range source.Certificates() {
		sum := sha256.Sum256(cert.Encoded())
		d := tvjCertDump{SHA256: hex.EncodeToString(sum[:])}
		tlUrls := map[string]struct{}{}
		territories := map[string]struct{}{}
		typeStatus := map[string]struct{}{}
		for _, tp := range source.TrustServices(cert) {
			if tlInfo := tp.TLInfo(); tlInfo != nil {
				tlUrls[tlInfo.Url()] = struct{}{}
			}
			if tsp := tp.TrustServiceProvider(); tsp != nil {
				territories[tsp.Territory()] = struct{}{}
			}
			if ts := tp.TrustService(); ts != nil {
				for entry := range ts.Iterator() {
					typeStatus[entry.Type()+"|"+entry.Status()] = struct{}{}
				}
			}
		}
		d.TLUrls = tvjSortedKeys(tlUrls)
		d.TSPTerritories = tvjSortedKeys(territories)
		d.ServiceTypeStatus = tvjSortedKeys(typeStatus)
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SHA256 < out[j].SHA256 })
	return out
}

func tvjSortedKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func tvjPrefix(prefix string, urls []string) []string {
	out := make([]string, len(urls))
	for i, u := range urls {
		out[i] = prefix + ":" + u
	}
	return out
}

// tvjRunJob builds and offline-refreshes the shared job configuration, returning the dump.
func tvjRunJob(t *testing.T) tvjDump {
	t.Helper()

	loader := &tvjFileLoader{files: map[string]string{
		"LOTL_URL":      corpustest.Path(t, "oracle/tsl/eu-lotl.xml"),
		"DE_TL_URL":     corpustest.Path(t, "oracle/tsl/de-tl.xml"),
		"BROKEN_TL_URL": corpustest.Path(t, "oracle/tsl/eu-lotl-broken-sig.xml"),
		"BAD_TL_URL":    corpustest.Path(t, "oracle/tsl/eu-lotl-not-parseable.xml"),
	}}

	lotlSource := tsl.NewLOTLSource()
	lotlSource.SetUrl("LOTL_URL")
	lotlSource.SetCertificateSource(tvjCertificateSource(t, tvjSigningCertificate))
	lotlSource.SetTlPredicate(tvjNoChildTLs{})

	deTLSource := tsl.NewTLSource()
	deTLSource.SetUrl("DE_TL_URL")
	deTLSource.SetCertificateSource(tvjCertificateSource(t))

	brokenTLSource := tsl.NewTLSource()
	brokenTLSource.SetUrl("BROKEN_TL_URL")
	brokenTLSource.SetCertificateSource(tvjCertificateSource(t))

	badTLSource := tsl.NewTLSource()
	badTLSource.SetUrl("BAD_TL_URL")
	badTLSource.SetCertificateSource(tvjCertificateSource(t))

	tlValidationJob := tsl.NewTLValidationJob()
	tlValidationJob.SetListOfTrustedListSources(lotlSource)
	tlValidationJob.SetTrustedListSources(deTLSource, brokenTLSource, badTLSource)
	tlValidationJob.SetOfflineDataLoader(loader)
	certSource := spitsl.NewTrustedListsCertificateSource()
	tlValidationJob.SetTrustedListCertificateSource(certSource)

	// Wire the four detection strategies for real, with a recording handler each - see this
	// file's header on why re-evaluating after the fact would not observe the same firings.
	//
	// TLSignatureErrorDetection/TLParsingErrorDetection are AlertDetector[*TLInfo]: per
	// ValidationJobAlerter.DetectChanges (validation/job/validation_job_alerter.go), a
	// documentAlerts (TLAlert) entry runs over a document LIST's ChildrenInfos() (the LOTL's own
	// auto-discovered child TLs - none here, tvjNoChildTLs above keeps that list empty) and over
	// OtherDocumentInfos() (this test's three independently-configured TLSources) - never over
	// the LOTL's own cache entry, which is why they are wired as TLAlerts, not LOTLAlerts, and
	// why this harness's LOTL (configured with its real, cross-validated signing certificate) has
	// no signature-error alert to observe. LOTLLocationChangeDetection/OJUrlChangeDetection are
	// the LOTL-specific counterparts, AlertDetector[*LOTLInfo], and run over documentListAlerts.
	var locationFired, ojUrlFired []string
	var tlSigFired, tlParseFired []string
	tlValidationJob.SetLOTLAlerts([]*tsl.LOTLAlert{
		tsl.NewLOTLAlert(tsl.NewLOTLLocationChangeDetection(lotlSource), &tvjURLRecorder[*tslmodel.LOTLInfo]{fired: &locationFired}),
		tsl.NewLOTLAlert(tsl.NewOJUrlChangeDetection(lotlSource), &tvjURLRecorder[*tslmodel.LOTLInfo]{fired: &ojUrlFired}),
	})
	tlValidationJob.SetTLAlerts([]*tsl.TLAlert{
		tsl.NewTLAlert(tsl.NewTLSignatureErrorDetection(), &tvjURLRecorder[*tslmodel.TLInfo]{fired: &tlSigFired}),
		tsl.NewTLAlert(tsl.NewTLParsingErrorDetection(), &tvjURLRecorder[*tslmodel.TLInfo]{fired: &tlParseFired}),
	})

	if err := tlValidationJob.OfflineRefresh(); err != nil {
		t.Fatalf("OfflineRefresh: %v", err)
	}

	summary := tlValidationJob.Summary()

	dump := tvjDump{
		NumberOfProcessedLOTLs: summary.NumberOfProcessedLOTLs(),
		NumberOfProcessedTLs:   summary.NumberOfProcessedTLs(),
		LOTLs:                  []tvjTLDump{},
		TLs:                    []tvjTLDump{},
	}
	for _, l := range summary.LOTLInfos() {
		dump.LOTLs = append(dump.LOTLs, tvjDumpTLInfo(&l.TLInfo))
	}
	for _, tlInfo := range summary.OtherTLInfos() {
		dump.TLs = append(dump.TLs, tvjDumpTLInfo(tlInfo))
	}

	var allFired []string
	allFired = append(allFired, tvjPrefix("LOTL_LOCATION", locationFired)...)
	allFired = append(allFired, tvjPrefix("OJ_URL", ojUrlFired)...)
	allFired = append(allFired, tvjPrefix("TL_SIG", tlSigFired)...)
	allFired = append(allFired, tvjPrefix("TL_PARSE", tlParseFired)...)
	sort.Strings(allFired)
	dump.FiredAlerts = allFired

	dump.TrustedCertificates = tvjDumpCertSource(certSource)

	return dump
}

// TestTLValidationJobOracle runs TLValidationJob offline over the vendored fixtures and compares
// the result against the Java oracle dump.
func TestTLValidationJobOracle(t *testing.T) {
	got := tvjRunJob(t)

	f, err := os.Open(corpustest.Path(t, "oracle/tsl/tl_validation_job.json"))
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1<<20), 16<<20)
	if !scanner.Scan() {
		t.Fatal("empty oracle file")
	}
	var want tvjDump
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
