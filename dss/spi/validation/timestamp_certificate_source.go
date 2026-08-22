// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampCertificateSource.java (DSS 6.5.RC1).
//
// spi.CMSCertificateSource (Java spi.x509.CMSCertificateSource) already turns BouncyCastle's
// SignerInformationStore into a cmscore SignerInfo slice and the Store of X509CertificateHolder
// into the certificate encodings of the CMS, which is exactly what a TimeStampToken hands out.
package validation

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
)

// TimestampCertificateSource is a timestamp CMS certificate source.
type TimestampCertificateSource struct {
	*spi.CMSCertificateSource
}

// NewTimestampCertificateSource extracts the certificates, certificate identifiers and
// certificate references of the given time-stamp token.
// Port of the TimestampCertificateSource(TimeStampToken) constructor.
//
// The source is not re-registered with InitSignatureCertificateSource: this class does not
// override extractCandidatesForSigningCertificate, so the registration NewCMSCertificateSource
// already performed reaches the right implementation.
func NewTimestampCertificateSource(timestampToken *cmscore.TimeStampToken) (*TimestampCertificateSource, error) {
	cms := timestampToken.CMS()
	signerInfos := cms.SignerInfos()
	base, err := spi.NewCMSCertificateSource(signerInfos, cms.Certificates(),
		spi.DSSASN1UtilsFirstSignerInformation(signerInfos))
	if err != nil {
		return nil, err
	}
	return &TimestampCertificateSource{CMSCertificateSource: base}, nil
}

// CertificateSourceType returns TIMESTAMP. Port of the getCertificateSourceType() override.
func (s *TimestampCertificateSource) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceType_TIMESTAMP
}
