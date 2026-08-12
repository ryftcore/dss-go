// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CommonX509URLCertificateSource.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCIES:
//   - X509URLCertificateSource (spi.x509, flattened into this package, ported in a sibling
//     chunk of phase 2a): assumed to be an interface embedding CertificateSource with a single
//     extra method CertificatesByURL(uri string) []*model.CertificateToken (port of
//     getCertificatesByUrl(String)).
//   - dsshttp.DataLoader (eu.europa.esig.dss.spi.client.http, its own Go package per
//     PORTING.md's package layout, not yet ported in this phase): assumed to expose
//     Get(url string) []byte (port of get(String), returning nil on failure), matching the
//     single method this file actually calls.
package spi

import (
	dsshttp "github.com/utain/esig/dss/http"
	"github.com/utain/esig/dss/model"
)

// CommonX509URLCertificateSource is the common implementation of X509URLCertificateSource
// retrieving X.509 certificates by the given URI. Used for validation of JAdES and CB-AdES
// signatures.
//
// This class provides the following workflows:
//   - Provide a mapping between URL and certificate pairs using AddCertificateForURL and/or
//     AddCertificatesForURL; or
//   - Instantiate the class either with a DataLoader to access the certificates in the
//     runtime, or by setting the DataLoader with SetDataLoader; or
//   - A combination of both.
type CommonX509URLCertificateSource struct {
	CommonCertificateSource

	// dataLoader is used to access 'x5u' certificates.
	dataLoader dsshttp.DataLoader

	// mapByUri maps uris to related certificate tokens.
	mapByUri map[string][]*model.CertificateToken
}

// NewCommonX509URLCertificateSource instantiates a pre-configured certificate source.
// Port of the default constructor.
func NewCommonX509URLCertificateSource() *CommonX509URLCertificateSource {
	return &CommonX509URLCertificateSource{
		CommonCertificateSource: NewCommonCertificateSource(),
		mapByUri:                make(map[string][]*model.CertificateToken),
	}
}

// NewCommonX509URLCertificateSourceWithDataLoader creates an instance of the class with a
// dataLoader to access the certificates from the corresponding 'x5u' location in the runtime.
// Port of the CommonX509URLCertificateSource(DataLoader) constructor.
func NewCommonX509URLCertificateSourceWithDataLoader(dataLoader dsshttp.DataLoader) *CommonX509URLCertificateSource {
	source := NewCommonX509URLCertificateSource()
	source.dataLoader = dataLoader
	return source
}

// SetDataLoader sets the DataLoader to access the certificates from the corresponding 'x5u'
// location in the runtime. Port of setDataLoader(DataLoader).
func (s *CommonX509URLCertificateSource) SetDataLoader(dataLoader dsshttp.DataLoader) {
	s.dataLoader = dataLoader
}

// AddCertificate is not supported in CommonX509URLCertificateSource. Port of the
// addCertificate(CertificateToken) override, which upstream throws
// UnsupportedOperationException from.
//
// Java's UnsupportedOperationException is a RuntimeException; the Go counterpart panics with
// the Java message, since callers of the CertificateSource interface's AddCertificate method
// cannot recover from this any more meaningfully than Java's caller could.
func (s *CommonX509URLCertificateSource) AddCertificate(certificateToAdd *model.CertificateToken) *model.CertificateToken {
	panic("#addCertificate(certificateToAdd) method is not supported in CommonX509URLCertificateSource! " +
		"Please use #addCertificate(uri, certificateToAdd) or #addCertificates(uri, certificatesToAdd) methods.")
}

// AddCertificateForURL adds a certificate for a given 'x5u' URL (JWS/JAdES).
// Port of addCertificate(String, CertificateToken), renamed since Go has no overloading with
// the interface's AddCertificate(CertificateToken).
func (s *CommonX509URLCertificateSource) AddCertificateForURL(uri string, certificate *model.CertificateToken) *model.CertificateToken {
	addedCertificate := s.CommonCertificateSource.AddCertificate(certificate)
	// URI already known: the certificate is added to the existing collection (upstream logs
	// this in debug mode; dropped per PORTING.md's slf4j rule).
	s.mapByUri[uri] = append(s.mapByUri[uri], certificate)
	return addedCertificate
}

// AddCertificatesForURL adds a collection of certificates for a given 'x5u' URL (JWS/JAdES).
// Port of addCertificates(String, Collection<CertificateToken>), renamed for the same reason
// as AddCertificateForURL.
func (s *CommonX509URLCertificateSource) AddCertificatesForURL(uri string, certificates []*model.CertificateToken) []*model.CertificateToken {
	for _, certificate := range certificates {
		addedCertificate := s.CommonCertificateSource.AddCertificate(certificate)
		s.mapByUri[uri] = append(s.mapByUri[uri], addedCertificate)
	}
	return s.mapByUri[uri]
}

// CertificatesByURL gets a collection of CertificateTokens retrieved from the given URI.
// Port of getCertificatesByUrl(String).
func (s *CommonX509URLCertificateSource) CertificatesByURL(uri string) []*model.CertificateToken {
	if certificates := s.mapByUri[uri]; len(certificates) > 0 {
		// Certificates are already known for this 'x5u' value; return the existing values.
		return certificates
	}
	if certificates := s.loadCertificatesFromURL(uri); len(certificates) > 0 {
		return certificates
	}
	return nil
}

// loadCertificatesFromURL loads the certificates using a DataLoader from the given url.
// Port of the protected loadCertificates(String).
//
// Upstream logs and swallows every failure (a nil DataLoader response, or an exception from
// loadCertificates(byte[])); the port mirrors that by returning nil in both cases instead of
// propagating an error, since the Java method has no return channel for it either.
func (s *CommonX509URLCertificateSource) loadCertificatesFromURL(url string) []*model.CertificateToken {
	if s.dataLoader == nil {
		// No DataLoader is configured within the CommonX509URLCertificateSource.
		return nil
	}
	content := s.dataLoader.Get(url)
	if content == nil {
		// No content has been extracted from the 'x5u' protected header with this value.
		return nil
	}
	return s.loadCertificatesFromContent(content)
}

// loadCertificatesFromContent loads certificates from the obtained content.
// Port of the protected loadCertificates(byte[]).
//
// Java wraps DSSUtils.loadCertificateFromP7c's checked exception into an unchecked
// DSSException that propagates out of getCertificatesByUrl; the port mirrors the "unable to
// load" outcome by returning nil, since this method (like loadCertificates(String) above) has
// no error return in the Java API it implements.
func (s *CommonX509URLCertificateSource) loadCertificatesFromContent(content []byte) []*model.CertificateToken {
	certificates, err := DSSUtilsLoadCertificateFromP7cBinary(content)
	if err != nil {
		return nil
	}
	return certificates
}

// FindTokensFromCertRef returns the certificate tokens for the provided CertificateRef, keyed
// by DSSIDAsString(), also considering the certificates known by their 'x5u' URL.
// Port of the findTokensFromCertRef(CertificateRef) override.
func (s *CommonX509URLCertificateSource) FindTokensFromCertRef(certificateRef *CertificateRef) map[string]*model.CertificateToken {
	certificates := s.CommonCertificateSource.FindTokensFromCertRef(certificateRef)
	if certificateRef.X509Url() != "" {
		for _, certificate := range s.mapByUri[certificateRef.X509Url()] {
			certificates[certificate.DSSIDAsString()] = certificate
		}
	}
	return certificates
}

// reset removes all certificates from the source. Port of the protected reset() override.
func (s *CommonX509URLCertificateSource) reset() {
	s.CommonCertificateSource.reset()
	s.mapByUri = make(map[string][]*model.CertificateToken)
}

// compile-time assertion: a CommonX509URLCertificateSource is a CertificateSource.
var _ CertificateSource = (*CommonX509URLCertificateSource)(nil)
