// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESCertificateSource.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
)

// CAdESCertificateSource is a CertificateSource that retrieves items from a CAdES Signature.
// Port of the class CAdESCertificateSource, extending spi.CMSCertificateSource.
type CAdESCertificateSource struct {
	*spi.CMSCertificateSource
}

// NewCAdESCertificateSource creates a CAdES certificate source from a CMS with an additional
// signer id parameter. All certificates are extracted during instantiation.
// Port of the constructor CAdESCertificateSource(CMS, SignerInformation).
func NewCAdESCertificateSource(cmsObj *cms.CMS, signerInformation *cmscore.SignerInfo) (*CAdESCertificateSource, error) {
	base, err := spi.NewCMSCertificateSource(cmsObj.SignerInfos(), cmsObj.Certificates(), signerInformation)
	if err != nil {
		return nil, err
	}
	return &CAdESCertificateSource{CMSCertificateSource: base}, nil
}
