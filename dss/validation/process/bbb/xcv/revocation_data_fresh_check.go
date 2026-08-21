// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rfc/checks/RevocationDataFreshCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationDataFreshCheck checks if the revocation data is fresh.
type RevocationDataFreshCheck struct {
	*AbstractRevocationFreshCheck

	// durationRule defines max freshness.
	durationRule policy.DurationRule
}

// NewRevocationDataFreshCheck is the default constructor. Port of
// RevocationDataFreshCheck(I18nProvider, XmlRFC, RevocationWrapper, Date, DurationRule).
func NewRevocationDataFreshCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRFC],
	revocationData *diagnostic.RevocationWrapper, validationDate time.Time, constraint policy.DurationRule) *RevocationDataFreshCheck {
	c := &RevocationDataFreshCheck{
		AbstractRevocationFreshCheck: NewAbstractRevocationFreshCheck(i18nProvider, result, revocationData, validationDate, constraint),
		durationRule:                 constraint,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationDataFreshCheck) Process() bool {
	if c.RevocationData != nil {
		return c.IsThisUpdateTimeAfterValidationTime(c.getMaxFreshness())
	}
	return false
}

// getMaxFreshness ports the protected getMaxFreshness().
func (c *RevocationDataFreshCheck) getMaxFreshness() int64 {
	return c.durationRule.Duration()
}
