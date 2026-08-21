// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/converter/DigitalIdentityListTypeConverter.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// DigitalIdentityListTypeConverter extracts CertificateTokens from a DigitalIdentityListType.
type DigitalIdentityListTypeConverter struct{}

// NewDigitalIdentityListTypeConverter is the default constructor. Port of
// DigitalIdentityListTypeConverter().
func NewDigitalIdentityListTypeConverter() *DigitalIdentityListTypeConverter {
	return &DigitalIdentityListTypeConverter{}
}

// Apply ports apply(DigitalIdentityListType). Errors loading a certificate are dropped (the
// slf4j debug/warn diagnostics upstream logs), matching PORTING.md.
func (c *DigitalIdentityListTypeConverter) Apply(digitalIdentityList *jaxb.DigitalIdentityListType) []*model.CertificateToken {
	var certificates []*model.CertificateToken
	if digitalIdentityList != nil && utils.IsCollectionNotEmpty(digitalIdentityList.DigitalId) {
		for _, digitalIdentity := range digitalIdentityList.DigitalId {
			if digitalIdentity.X509Certificate != nil && utils.IsArrayNotEmpty([]byte(*digitalIdentity.X509Certificate)) {
				if certificate, err := spi.DSSUtilsLoadCertificateFromBinary([]byte(*digitalIdentity.X509Certificate)); err == nil {
					certificates = append(certificates, certificate)
				}
			}
		}
	}
	return certificates
}
