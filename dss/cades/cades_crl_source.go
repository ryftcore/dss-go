// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESCRLSource.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
)

// CAdESCRLSource is the CRL source for a CAdES signature. Port of the class CAdESCRLSource,
// extending spi.CMSCRLSource.
type CRLSource struct {
	*spi.CMSCRLSource
}

// NewCAdESCRLSource creates a CAdES CRL source from a CMS and the related unsignedAttributes
// of the signer. Port of the constructor CAdESCRLSource(CMS, AttributeTable).
func NewCAdESCRLSource(cmsObj *cms.CMS, unsignedAttributes cmscore.Attributes) (*CRLSource, error) {
	base, err := spi.NewCMSCRLSource(cmsObj.CRLs(), unsignedAttributes)
	if err != nil {
		return nil, err
	}
	return &CRLSource{CMSCRLSource: base}, nil
}
