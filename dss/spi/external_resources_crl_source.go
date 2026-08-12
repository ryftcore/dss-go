// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/crl/ExternalResourcesCRLSource.java (DSS 6.5.RC1).
//
// DEVIATION: the String-based constructor loads classpath resources upstream
// (getClass().getResourceAsStream). Go has no classpath, so the paths are read from the
// filesystem.
package spi

import (
	"io"
	"os"

	"github.com/utain/esig/dss/crlparser"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/utils"
)

// ExternalResourcesCRLSource provides a CRL source based on a list of external CRLs.
type ExternalResourcesCRLSource struct {
	OfflineCRLSourceBase
}

// NewExternalResourcesCRLSource builds an empty source, ready to be filled by the loaders
// below. The Java class has no such constructor; the Go port needs one because each of its
// three constructors has to register the source with its base first.
func NewExternalResourcesCRLSource() *ExternalResourcesCRLSource {
	source := &ExternalResourcesCRLSource{OfflineCRLSourceBase: NewOfflineCRLSourceBase()}
	source.InitOfflineRevocationSource(source)
	return source
}

// NewExternalResourcesCRLSourceFromPaths builds a CRL source from an array of paths.
// Port of the ExternalResourcesCRLSource(String...) constructor.
func NewExternalResourcesCRLSourceFromPaths(paths ...string) (*ExternalResourcesCRLSource, error) {
	source := NewExternalResourcesCRLSource()
	for _, pathItem := range paths {
		file, err := os.Open(pathItem)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause("Unable to parse the stream (CRL is expected)", err)
		}
		if err := source.addCRLToken(file); err != nil {
			return nil, err
		}
	}
	return source, nil
}

// NewExternalResourcesCRLSourceFromReaders builds a CRL source from an array of readers.
// Port of the ExternalResourcesCRLSource(InputStream...) constructor; every reader that is
// also an io.Closer is closed, as the try-with-resources upstream does.
func NewExternalResourcesCRLSourceFromReaders(readers ...io.Reader) (*ExternalResourcesCRLSource, error) {
	source := NewExternalResourcesCRLSource()
	for _, reader := range readers {
		if err := source.addCRLToken(reader); err != nil {
			return nil, err
		}
	}
	return source, nil
}

// NewExternalResourcesCRLSourceFromDocuments builds a CRL source from an array of documents.
// Port of the ExternalResourcesCRLSource(DSSDocument...) constructor.
func NewExternalResourcesCRLSourceFromDocuments(documents ...model.DSSDocument) (*ExternalResourcesCRLSource, error) {
	source := NewExternalResourcesCRLSource()
	for _, document := range documents {
		stream, err := document.OpenStream()
		if err != nil {
			return nil, model.NewDSSErrorMessageCause("Unable to parse the stream (CRL is expected)", err)
		}
		if err := source.addCRLToken(stream); err != nil {
			return nil, err
		}
	}
	return source, nil
}

// addCRLToken ports the private addCRLToken(InputStream), whose try-with-resources closes
// the stream whatever happens.
func (s *ExternalResourcesCRLSource) addCRLToken(reader io.Reader) error {
	if closer, ok := reader.(io.Closer); ok {
		defer utils.CloseQuietly(closer)
	}
	binaries, err := utils.ToByteArray(reader)
	if err != nil {
		return model.NewDSSErrorMessageCause("Unable to parse the stream (CRL is expected)", err)
	}
	crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(binaries)
	if err != nil {
		return model.NewDSSErrorMessageCause("Unable to parse the stream (CRL is expected)", err)
	}
	s.AddBinary(crlBinary, enumerations.RevocationOrigin_EXTERNAL)
	return nil
}

// RevocationTokens marks every token the base source produced as coming from an external
// resource. Port of the getRevocationTokens(CertificateToken, CertificateToken) override.
func (s *ExternalResourcesCRLSource) RevocationTokens(certificate *model.CertificateToken,
	issuer *model.CertificateToken) ([]RevocationToken[revocation.CRL], error) {
	revocationTokens, err := s.OfflineCRLSourceBase.RevocationTokens(certificate, issuer)
	if err != nil {
		return nil, err
	}
	for _, revocationToken := range revocationTokens {
		revocationToken.SetExternalOrigin(enumerations.RevocationOrigin_EXTERNAL)
	}
	return revocationTokens, nil
}
