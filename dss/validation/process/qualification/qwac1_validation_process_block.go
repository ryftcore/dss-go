// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/QWAC1ValidationProcessBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
)

// QWAC1ValidationProcessBlock performs validation of the 1-QWAC profile as
// per ETSI TS 119 411-5 part "4.1 1-QWAC Approach".
type QWAC1ValidationProcessBlock struct {
	*AbstractQWACValidationProcessBlock
}

// NewQWAC1ValidationProcessBlock is the common constructor. Port of
// QWAC1ValidationProcessBlock(Provider, Date, CertificateWrapper, XmlConclusion, XmlCertificateQualificationProcess, String).
func NewQWAC1ValidationProcessBlock(i18nProvider *i18n.Provider, validationTime time.Time,
	certificate *diagnostic.CertificateWrapper, buildingBlocksConclusion *jaxb.XmlConclusion,
	certificateQualification *jaxb.XmlCertificateQualificationProcess, websiteUrl string) *QWAC1ValidationProcessBlock {
	c := &QWAC1ValidationProcessBlock{AbstractQWACValidationProcessBlock: &AbstractQWACValidationProcessBlock{}}
	c.InitAbstractQWACValidationProcessBlock(i18nProvider, validationTime, certificate, buildingBlocksConclusion,
		certificateQualification, websiteUrl, c)
	c.InitChainBase(c)
	return c
}

// QWACProfile gets the current QWAC profile. Port of the overridden public
// QWACProfile getQWACProfile().
func (c *QWAC1ValidationProcessBlock) QWACProfile() enumerations.QWACProfile {
	return enumerations.QWACProfileQWAC1
}
