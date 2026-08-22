// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/aia/DefaultAIASource.java (DSS 6.5.RC1).
package aia

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	dsshttp "github.com/ryftcore/dss-go/dss/spi/client/http"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// defaultAIASourceAllProtocols lists every dsshttp.Protocol value, standing in for
// Protocol.values(); used as the default acceptedProtocols.
var defaultAIASourceAllProtocols = []dsshttp.Protocol{
	dsshttp.ProtocolFile,
	dsshttp.ProtocolHTTP,
	dsshttp.ProtocolHTTPS,
	dsshttp.ProtocolLDAP,
	dsshttp.ProtocolFTP,
}

// DefaultAIASource is used to download issuer certificates by AIA from remote sources.
type DefaultAIASource struct {
	// dataLoader is used to download data.
	dataLoader dsshttp.DataLoader

	// acceptedProtocols is the collection of protocols to be accepted and used by the
	// source. Default: all protocols are accepted (FILE, HTTP, HTTPS, LDAP, FTP).
	acceptedProtocols []dsshttp.Protocol
}

// NewDefaultAIASource instantiates a dsshttp.NativeHTTPDataLoader as the default data loader.
func NewDefaultAIASource() *DefaultAIASource {
	return NewDefaultAIASourceWithDataLoader(dsshttp.NewNativeHTTPDataLoader())
}

// NewDefaultAIASourceWithDataLoader is the default constructor with a defined DataLoader.
// Panics if dataLoader is nil (Java Objects.requireNonNull).
func NewDefaultAIASourceWithDataLoader(dataLoader dsshttp.DataLoader) *DefaultAIASource {
	if dataLoader == nil {
		panic("dataLoader cannot be null!")
	}
	return &DefaultAIASource{dataLoader: dataLoader, acceptedProtocols: defaultAIASourceAllProtocols}
}

// SetDataLoader sets the data loader to be used to download a certificate token by AIA. Panics
// if dataLoader is nil.
func (s *DefaultAIASource) SetDataLoader(dataLoader dsshttp.DataLoader) {
	if dataLoader == nil {
		panic("dataLoader cannot be null!")
	}
	s.dataLoader = dataLoader
}

// SetAcceptedProtocols defines a set of protocols to be accepted and used by the AIA Source.
// All protocols which are not defined in the collection will be skipped.
// Default: all protocols are accepted (FILE, HTTP, HTTPS, LDAP, FTP).
func (s *DefaultAIASource) SetAcceptedProtocols(acceptedProtocols []dsshttp.Protocol) {
	s.acceptedProtocols = acceptedProtocols
}

// CertificatesByAIA loads issuer certificates for certificateToken from its AIA caIssuers
// URLs, trying each URL in turn until one yields certificates.
//
// Panics if certificateToken is nil, or if no DataLoader is configured (Java
// Objects.requireNonNull).
func (s *DefaultAIASource) CertificatesByAIA(certificateToken *model.CertificateToken) []*model.CertificateToken {
	if certificateToken == nil {
		panic("CertificateToken cannot be null!")
	}
	if s.dataLoader == nil {
		panic("DataLoader is not provided!")
	}

	caIssuersUrls := s.getCAIssuersUrls(certificateToken)

	for _, caIssuersURL := range caIssuersUrls {
		loadedCertificates, ok := defaultAIASourceTryLoad(s, caIssuersURL)
		if ok {
			for _, certificate := range loadedCertificates {
				certificate.SetSourceURL(caIssuersURL)
			}
			return defaultAIASourceDedup(loadedCertificates)
		}
		// "Unable to retrieve AIA certificates with URL..." LOG.warn dropped, not
		// load-bearing per PORTING.md.
	}

	return nil
}

// defaultAIASourceTryLoad executes the caIssuers request and loads the resulting P7C
// certificates, recovering from any panic raised along the way (e.g. by
// executeCAIssuersRequest or a malformed P7C body) and reporting failure via ok=false. Ports
// the per-URL try/catch(Exception) inside getCertificatesByAIA.
func defaultAIASourceTryLoad(s *DefaultAIASource, caIssuersURL string) (loaded []*model.CertificateToken, ok bool) {
	defer func() {
		if recover() != nil {
			loaded, ok = nil, false
		}
	}()
	bytes := s.executeCAIssuersRequest(caIssuersURL)
	certificates, err := spi.DSSUtilsLoadCertificateFromP7cBinary(bytes)
	if err != nil {
		panic(err)
	}
	return certificates, true
}

// getCAIssuersUrls returns a list of caIssuers URLs for the given certificateToken.
func (s *DefaultAIASource) getCAIssuersUrls(certificateToken *model.CertificateToken) []string {
	urls := spi.CertificateExtensionsUtilsCAIssuersAccessUrls(certificateToken)
	if utils.IsCollectionEmpty(urls) {
		// "There is no AIA extension for certificate download." LOG.info dropped, not
		// load-bearing per PORTING.md.
		return nil
	}
	return s.filterURLs(urls)
}

// filterURLs keeps only the URLs whose protocol is accepted.
func (s *DefaultAIASource) filterURLs(urls []string) []string {
	var filtered []string
	for _, url := range urls {
		if s.isURLAccepted(url) {
			filtered = append(filtered, url)
		}
	}
	return filtered
}

// executeCAIssuersRequest executes a GET request to retrieve caIssuers from caIssuersURL.
// Panics with a *exception.DSSExternalResourceException when the DataLoader returns an empty
// response.
func (s *DefaultAIASource) executeCAIssuersRequest(caIssuersURL string) []byte {
	bytes := s.dataLoader.Get(caIssuersURL)
	if utils.IsArrayNotEmpty(bytes) {
		return bytes
	}
	panic(exception.NewDSSExternalResourceException(
		fmt.Sprintf("AIA DataLoader for certificate with url '%s' responded with an empty byte array!", caIssuersURL)))
}

// isURLAccepted reports whether url's protocol is one of acceptedProtocols.
func (s *DefaultAIASource) isURLAccepted(url string) bool {
	if utils.IsCollectionNotEmpty(s.acceptedProtocols) {
		for _, protocol := range s.acceptedProtocols {
			if protocol.IsTheSame(url) {
				return true
			}
		}
	}
	return false
}

// defaultAIASourceDedup de-duplicates certificateTokens by DSSIDAsString() while preserving
// order, standing in for `new LinkedHashSet<>(...)`.
func defaultAIASourceDedup(certificateTokens []*model.CertificateToken) []*model.CertificateToken {
	seen := make(map[string]struct{}, len(certificateTokens))
	result := make([]*model.CertificateToken, 0, len(certificateTokens))
	for _, certificateToken := range certificateTokens {
		id := certificateToken.DSSIDAsString()
		if _, found := seen[id]; found {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, certificateToken)
	}
	return result
}

var _ AIASource = (*DefaultAIASource)(nil)
