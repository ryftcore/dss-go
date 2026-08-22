// Command 07-validate-eu-trusted-lists runs the TSL validation job
// (dss/tsl and dss/validation/job) against the real European
// List Of Trusted Lists (LOTL) and reports what it found. This is the piece
// that turns "the signature verifies" into "the signature is eIDAS
// qualified": qualification is determined from trusted-list content, which
// only a trusted-list source can supply.
//
// This example needs the network, and degrades gracefully without it: on
// any download failure it prints why and exits 0 rather than failing, which
// is also exactly what job.OnlineRefresh does internally for every trusted
// list it is configured with - a download failure is recorded as a cache
// state, not a panic or a returned error.
//
// To keep the example small and fast it fetches ONLY the LOTL document
// itself, not the 27+ country trusted lists it points to - a real
// deployment sets LOTLSource.SetTlPredicate to
// tsl.TLPredicateFactoryCreateEUTLPredicate() (the LOTLSource default) to
// follow all of them, which easily means a hundred HTTP requests.
//
// Run it from anywhere:
//
//	go run ./examples/07-validate-eu-trusted-lists
package main

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
	dsspolicy "github.com/ryftcore/dss-go/dss/policy"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
	spitsl "github.com/ryftcore/dss-go/dss/spi/tsl"
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/tsl"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
)

// init registers the ETSI validation policy factory the TL/LOTL signature
// validation task needs to reach a verdict at all. The dss facade package
// performs this registration automatically on import (see its package doc's
// "Registration" section); this example works one layer below the facade,
// so it registers the same factory itself, exactly as a caller of
// dss/tsl directly has to.
func init() {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
}

// euLOTLURL is the well-known address of the official EU List Of Trusted
// Lists, published under Commission Implementing Decision (EU) 2015/1505.
// It is not something this repository can verify at build time - the
// library ships no default LOTL URL of its own, on purpose (see the dss
// package doc's "Network access" section) - so treat it as an example
// value, not a guarantee this URL still resolves when you read this.
const euLOTLURL = "https://ec.europa.eu/tools/lotl/eu-lotl.xml"

// noChildTLs is a tsl.OtherTSLPointerPredicate matching nothing, used so the
// job fetches only the LOTL itself and none of the trusted lists it points
// to - see the package comment for why.
type noChildTLs struct{}

func (noChildTLs) Test(*jaxb.OtherTSLPointerType) bool { return false }

// fileLoader adapts http.NativeHTTPDataLoader (a raw-bytes DataLoader) to
// http.DSSFileLoader (a model.DSSDocument loader), which is what
// TLValidationJob.SetOnlineDataLoader takes.
type fileLoader struct {
	inner *http.NativeHTTPDataLoader
}

func (l *fileLoader) GetDocument(url string) (model.DSSDocument, error) {
	data, err := l.inner.GetWithRefresh(url, false)
	if err != nil {
		return nil, err
	}
	return model.NewInMemoryDocument(data), nil
}

func main() {
	loader := http.NewNativeHTTPDataLoader()
	loader.SetConnectTimeout(5000) // milliseconds; keep the example responsive when offline
	loader.SetReadTimeout(5000)

	lotlSource := tsl.NewLOTLSource()
	lotlSource.SetUrl(euLOTLURL)
	lotlSource.SetTlPredicate(noChildTLs{})
	// A LOTLSource with no CertificateSource never runs signature
	// validation at all - its validation cache state simply stays
	// unset. Giving it an empty CertificateSource makes the job actually
	// validate the LOTL's XML signature (and, correctly, find no anchor
	// for it - a real deployment supplies the EU Official Journal signing
	// certificates here instead of an empty source).
	lotlSigningCerts := spi.NewCommonCertificateSource()
	lotlSource.SetCertificateSource(&lotlSigningCerts)

	trustedSource := spitsl.NewTrustedListsCertificateSource()

	job := tsl.NewTLValidationJob()
	job.SetOnlineDataLoader(&fileLoader{inner: loader})
	job.SetTrustedListCertificateSource(trustedSource)
	job.SetListOfTrustedListSources(lotlSource)

	// OnlineRefresh downloads, parses and validates every configured
	// source, then synchronizes trustedSource with whatever succeeded. It
	// does not return an error for an unreachable URL - that failure is
	// recorded per source, in the cache - so this is not a network try/
	// catch: it is how the job behaves online or offline alike.
	if err := job.OnlineRefresh(); err != nil {
		// Only a misconfiguration (no data loader, for instance) reaches
		// here; see the package comment on OnlineRefresh's own errors.
		fmt.Println("could not run the TSL validation job:", err)
		return
	}

	summary := job.Summary()
	lotlInfos := summary.LOTLInfos()
	if len(lotlInfos) == 0 {
		fmt.Println("no LOTL was processed")
		return
	}
	lotl := lotlInfos[0]

	download := lotl.DownloadCacheInfo()
	fmt.Println("download:", download.StatusName())
	if download.IsError() {
		// No network access, a proxy or firewall intercepting the request,
		// or the LOTL having moved can all land here - job.OnlineRefresh
		// turns every one of them into this cache state rather than an
		// error or a panic, which is what "degrades gracefully offline"
		// means in practice.
		fmt.Println("  ", download.ExceptionMessage())
		fmt.Println("this example degrades gracefully offline - see the package comment")
		return
	}

	parsing := lotl.ParsingCacheInfo()
	fmt.Println("parsing: ", parsing.StatusName())
	if parsing.IsError() {
		fmt.Println("  ", parsing.ExceptionMessage())
		return
	}

	validation := lotl.ValidationCacheInfo()
	fmt.Printf("validation: %s (indication=%s)\n", validation.StatusName(), validation.Indication())

	fmt.Printf("trusted certificates collected: %d\n", trustedSource.NumberOfCertificates())
	fmt.Println("(0 is expected: this example does not follow the LOTL's own country-TL")
	fmt.Println(" pointers, so no trust service certificates are ever reached - see above)")

	// This is where trustedSource plugs into ValidateOptions in a real
	// deployment - once it actually holds trusted-list content, which needs
	// following the country TLs this example skips:
	//
	//	reports, err := dss.Validate(signedDoc, dss.ValidateOptions{
	//	    TrustedCertificateSources: []dss.CertificateSource{trustedSource},
	//	})
	//
	// A signature whose chain reaches a certificate trustedSource carries
	// gets a real SignatureQualification (QESig, QESeal, ...) instead of the
	// "NA" every other example in this directory prints, because those never
	// hand Validate a trusted-list source at all.
}
