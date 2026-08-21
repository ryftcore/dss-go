package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/utain/esig/dss"
	"github.com/utain/esig/dss/model"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	dsshttp "github.com/utain/esig/dss/spi/client/http"
	spitsl "github.com/utain/esig/dss/spi/tsl"
	"github.com/utain/esig/dss/tsl"
)

// defaultLOTLURL is the European Commission's published List Of Trusted
// Lists location: the same URL Java DSS's own cookbook examples and test
// fixtures use as the EU LOTL.
const defaultLOTLURL = "https://ec.europa.eu/tools/lotl/eu-lotl.xml"

// cmdTL implements "esig tl" and its one subcommand, "refresh".
func cmdTL(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || (args[0] != "refresh" && args[0] != "-h" && args[0] != "-help" && args[0] != "--help") {
		fmt.Fprint(stderr, "Usage: esig tl refresh [flags]\n\nRun \"esig tl refresh -h\" for its flags.\n")
		return exitUsage
	}
	if args[0] != "refresh" {
		fmt.Fprint(stdout, "Usage: esig tl refresh [flags]\n\nRun \"esig tl refresh -h\" for its flags.\n")
		return exitOK
	}
	return cmdTLRefresh(args[1:], stdout, stderr)
}

// cmdTLRefresh implements "esig tl refresh".
func cmdTLRefresh(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("tl refresh", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, `Usage: esig tl refresh -cache <dir> [flags]

Downloads and parses the EU List of Trusted Lists and every member state
Trusted List it points to (github.com/utain/esig/dss/tsl,
TLValidationJob.OnlineRefresh), and writes the certificates it collects to
<dir> for "esig validate -tl-cache" to use as trust anchors.

  -lotl string
    	List of Trusted Lists URL (default %s)
  -cache string
    	directory to write the certificate cache to (required)
  -lotl-cert string
    	trust anchor certificate for the LOTL's own signature (DER or PEM); repeatable

Without -lotl-cert the LOTL's own XAdES signature cannot be anchored, so its
Indication will not be TOTAL_PASSED; the trusted lists it points to are still
downloaded, parsed and synchronized regardless (synchronization is gated on
successful parsing, not on the LOTL's signature validating), and their
certificates are still cached. This mirrors how the underlying
TLValidationJob behaves, not a simplification this command adds; see the
tsl package doc for the trust chain a production deployment needs to supply.
`, defaultLOTLURL)
	}
	var lotlURL, cacheDir string
	var lotlCerts stringList
	fs.StringVar(&lotlURL, "lotl", defaultLOTLURL, "")
	fs.StringVar(&cacheDir, "cache", "", "")
	fs.Var(&lotlCerts, "lotl-cert", "")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if cacheDir == "" {
		fmt.Fprintf(stderr, "esig tl refresh: -cache is required\n\n")
		fs.Usage()
		return exitUsage
	}

	trustAnchors, err := loadTrustAnchors(lotlCerts)
	if err != nil {
		fmt.Fprintf(stderr, "esig tl refresh: %v\n", err)
		return exitRuntime
	}

	lotlSource := tsl.NewLOTLSource()
	lotlSource.SetUrl(lotlURL)
	lotlSource.SetPivotSupport(true)
	if len(trustAnchors) > 0 {
		lotlSource.SetCertificateSource(dss.TrustStore(trustAnchors...))
	}

	job := tsl.NewTLValidationJob()
	job.SetListOfTrustedListSources(lotlSource)
	job.SetOnlineDataLoader(newHTTPFileLoader())
	certSource := spitsl.NewTrustedListsCertificateSource()
	job.SetTrustedListCertificateSource(certSource)

	if err := job.OnlineRefresh(); err != nil {
		fmt.Fprintf(stderr, "esig tl refresh: %v\n", err)
		return exitRuntime
	}

	summary := job.Summary()
	fmt.Fprintf(stdout, "processed %d LOTL(s), %d trusted list(s)\n",
		summary.NumberOfProcessedLOTLs(), summary.NumberOfProcessedTLs())
	for _, l := range summary.LOTLInfos() {
		fmt.Fprintf(stdout, "  LOTL %s: %s\n", l.Url(), tlInfoStatus(&l.TLInfo))
	}
	for _, t := range summary.OtherTLInfos() {
		fmt.Fprintf(stdout, "  TL   %s: %s\n", t.Url(), tlInfoStatus(t))
	}

	certs := certSource.Certificates()
	fmt.Fprintf(stdout, "%d trusted certificate(s) synchronized\n", len(certs))
	if err := writeTLCache(cacheDir, certs); err != nil {
		fmt.Fprintf(stderr, "esig tl refresh: %v\n", err)
		return exitRuntime
	}
	fmt.Fprintf(stdout, "cache written to %s\n", cacheDir)
	return exitOK
}

// tlInfoStatus renders one TLInfo/LOTLInfo's download/parsing/validation
// cache status as a short summary line.
func tlInfoStatus(info *tslmodel.TLInfo) string {
	status := "download=" + cacheInfoStatus(info.DownloadCacheInfo())
	status += " parsing=" + cacheInfoStatus(info.ParsingCacheInfo())
	validation := info.ValidationCacheInfo()
	status += " validation=" + cacheInfoStatus(validation)
	if validation != nil && !validation.IsError() {
		status += fmt.Sprintf(" (%s/%s)", validation.Indication(), validation.SubIndication())
	}
	return status
}

// cacheInfoRecord is the subset of model/job.InfoRecord tlInfoStatus needs.
type cacheInfoRecord interface {
	StatusName() string
}

func cacheInfoStatus(r cacheInfoRecord) string {
	if r == nil {
		return "n/a"
	}
	return r.StatusName()
}

// tlFetchTimeoutMillis bounds each LOTL/TL download: long enough for a slow
// government server, short enough that a hung connection cannot wedge the
// command indefinitely.
const tlFetchTimeoutMillis = 30_000

// httpFileLoader adapts [dsshttp.NativeHTTPDataLoader] to
// [dsshttp.DSSFileLoader], the interface [tsl.TLValidationJob.SetOnlineDataLoader]
// takes: a real network fetch, wrapped as a [model.DSSDocument].
type httpFileLoader struct {
	native dsshttp.NativeHTTPDataLoader
}

// newHTTPFileLoader returns an httpFileLoader with bounded connect/read
// timeouts (the zero value has none - see NativeHTTPDataLoader's own doc
// comments).
func newHTTPFileLoader() *httpFileLoader {
	l := &httpFileLoader{}
	l.native.SetConnectTimeout(tlFetchTimeoutMillis)
	l.native.SetReadTimeout(tlFetchTimeoutMillis)
	return l
}

var _ dsshttp.DSSFileLoader = (*httpFileLoader)(nil)

func (l *httpFileLoader) GetDocument(url string) (model.DSSDocument, error) {
	data, err := l.native.GetWithRefresh(url, false)
	if err != nil {
		return nil, err
	}
	return model.NewInMemoryDocument(data), nil
}
