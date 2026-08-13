// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESCRLSource.java (DSS 6.5.RC1).
package cades

import (
	"github.com/utain/esig/dss/cms"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/spi"
)

// CAdESCRLSource is the CRL source for a CAdES signature. Port of the class CAdESCRLSource,
// extending spi.CMSCRLSource.
type CAdESCRLSource struct {
	*spi.CMSCRLSource
}

// NewCAdESCRLSource creates a CAdES CRL source from a CMS and the related unsignedAttributes
// of the signer. Port of the constructor CAdESCRLSource(CMS, AttributeTable).
func NewCAdESCRLSource(cmsObj *cms.CMS, unsignedAttributes cmscore.Attributes) (*CAdESCRLSource, error) {
	base, err := spi.NewCMSCRLSource(cmsObj.CRLs(), unsignedAttributes)
	if err != nil {
		return nil, err
	}
	return &CAdESCRLSource{CMSCRLSource: base}, nil
}
