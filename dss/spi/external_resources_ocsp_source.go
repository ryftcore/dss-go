// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/ExternalResourcesOCSPSource.java (DSS 6.5.RC1).
//
// DEVIATION: the String-based constructor loads classpath resources upstream
// (getClass().getResourceAsStream). Go has no classpath, so the paths are read from the
// filesystem.
package spi

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/utils"
)

// ExternalResourcesOCSPSource provides a collection of OCSP tokens supplied by the user.
type ExternalResourcesOCSPSource struct {
	OfflineOCSPSourceBase
}

// NewExternalResourcesOCSPSource builds an empty source, ready to be filled by the loaders
// below. The Java class has no such constructor; the Go port needs one because each of its
// three constructors has to register the source with its base first.
func NewExternalResourcesOCSPSource() *ExternalResourcesOCSPSource {
	source := &ExternalResourcesOCSPSource{OfflineOCSPSourceBase: NewOfflineOCSPSourceBase()}
	source.InitOfflineRevocationSource(source)
	return source
}

// NewExternalResourcesOCSPSourceFromPaths loads the OCSP responses found at the given paths.
// Port of the ExternalResourcesOCSPSource(String...) constructor.
func NewExternalResourcesOCSPSourceFromPaths(paths ...string) (*ExternalResourcesOCSPSource, error) {
	source := NewExternalResourcesOCSPSource()
	for _, pathItem := range paths {
		file, err := os.Open(pathItem)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to load OCSP token : %s", err), err)
		}
		if err := source.load(file); err != nil {
			return nil, err
		}
	}
	return source, nil
}

// NewExternalResourcesOCSPSourceFromReaders loads the OCSP responses read from the given
// readers. Port of the ExternalResourcesOCSPSource(InputStream...) constructor; every reader
// that is also an io.Closer is closed, as the try-with-resources upstream does.
func NewExternalResourcesOCSPSourceFromReaders(readers ...io.Reader) (*ExternalResourcesOCSPSource, error) {
	source := NewExternalResourcesOCSPSource()
	for _, reader := range readers {
		if err := source.load(reader); err != nil {
			return nil, err
		}
	}
	return source, nil
}

// NewExternalResourcesOCSPSourceFromDocuments loads the OCSP responses carried by the given
// documents. Port of the ExternalResourcesOCSPSource(DSSDocument...) constructor.
func NewExternalResourcesOCSPSourceFromDocuments(documents ...model.DSSDocument) (*ExternalResourcesOCSPSource, error) {
	source := NewExternalResourcesOCSPSource()
	for _, document := range documents {
		stream, err := document.OpenStream()
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to load OCSP token : %s", err), err)
		}
		if err := source.load(stream); err != nil {
			return nil, err
		}
	}
	return source, nil
}

// load adds the basic OCSP response read from the reader to the general list.
// Port of the private load(InputStream), whose try-with-resources closes the stream whatever
// happens and whose catch-all becomes the returned error.
func (s *ExternalResourcesOCSPSource) load(reader io.Reader) error {
	if closer, ok := reader.(io.Closer); ok {
		defer utils.CloseQuietly(closer)
	}
	binaries, err := utils.ToByteArray(reader)
	if err != nil {
		return model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to load OCSP token : %s", err), err)
	}
	ocspResp, err := NewOCSPRespFromBinaries(binaries)
	if err != nil {
		return model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to load OCSP token : %s", err), err)
	}
	basicOCSPResp, err := ocspResp.ResponseObject()
	if err != nil {
		return model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to load OCSP token : %s", err), err)
	}
	if basicOCSPResp == nil {
		// Java's (BasicOCSPResp) cast raises a ClassCastException for any other response
		// type, which the catch-all turns into the same DSSException.
		cause := errors.New("the response does not encapsulate a BasicOCSPResponse")
		return model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to load OCSP token : %s", cause), cause)
	}
	ocspResponseBinary, err := OCSPResponseBinaryBuild(basicOCSPResp)
	if err != nil {
		return model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to load OCSP token : %s", err), err)
	}
	s.AddBinary(ocspResponseBinary, enumerations.RevocationOrigin_EXTERNAL)
	return nil
}

// RevocationTokens marks every token the base source produced as coming from an external
// resource. Port of the getRevocationTokens(CertificateToken, CertificateToken) override.
func (s *ExternalResourcesOCSPSource) RevocationTokens(certificate *model.CertificateToken,
	issuer *model.CertificateToken) ([]RevocationToken[revocation.OCSP], error) {
	revocationTokens, err := s.OfflineOCSPSourceBase.RevocationTokens(certificate, issuer)
	if err != nil {
		return nil, err
	}
	for _, revocationToken := range revocationTokens {
		revocationToken.SetExternalOrigin(enumerations.RevocationOrigin_EXTERNAL)
	}
	return revocationTokens, nil
}
