// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/QWAC2ValidationProcessBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
)

// QWAC2ValidationProcessBlock performs validation of the 2-QWAC profile as
// per ETSI TS 119 411-5 part "4.2 2-QWAC Approach".
type QWAC2ValidationProcessBlock struct {
	*AbstractQWACValidationProcessBlock
}

// NewQWAC2ValidationProcessBlock is the common constructor. Port of
// QWAC2ValidationProcessBlock(I18nProvider, Date, CertificateWrapper, XmlConclusion, XmlCertificateQualificationProcess, String).
func NewQWAC2ValidationProcessBlock(i18nProvider *i18n.I18nProvider, validationTime time.Time,
	certificate *diagnostic.CertificateWrapper, buildingBlocksConclusion *jaxb.XmlConclusion,
	certificateQualification *jaxb.XmlCertificateQualificationProcess, websiteUrl string) *QWAC2ValidationProcessBlock {
	c := &QWAC2ValidationProcessBlock{AbstractQWACValidationProcessBlock: &AbstractQWACValidationProcessBlock{}}
	c.InitAbstractQWACValidationProcessBlock(i18nProvider, validationTime, certificate, buildingBlocksConclusion,
		certificateQualification, websiteUrl, c)
	c.InitChainBase(c)
	return c
}

// QWACProfile gets the current QWAC profile. Port of the overridden public
// QWACProfile getQWACProfile().
func (c *QWAC2ValidationProcessBlock) QWACProfile() enumerations.QWACProfile {
	return enumerations.QWACProfile_QWAC_2
}
