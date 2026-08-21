// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESOCSPSource.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
)

// CAdESOCSPSource is the OCSP source for a CAdES signature. Port of the class CAdESOCSPSource,
// extending spi.CMSOCSPSource.
type CAdESOCSPSource struct {
	*spi.CMSOCSPSource
}

// NewCAdESOCSPSource creates a CAdES OCSP source from a CMS and the related
// unsignedAttributes of the signer. Port of the constructor CAdESOCSPSource(CMS, AttributeTable).
func NewCAdESOCSPSource(cmsObj *cms.CMS, unsignedAttributes cmscore.Attributes) (*CAdESOCSPSource, error) {
	base, err := spi.NewCMSOCSPSource(cmsObj.OcspResponseStore(), cmsObj.OcspBasicStore(), unsignedAttributes)
	if err != nil {
		return nil, err
	}
	return &CAdESOCSPSource{CMSOCSPSource: base}, nil
}
