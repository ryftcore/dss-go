// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rfc/checks/RevocationDataFreshCheckWithNullConstraint.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// RevocationDataFreshCheckWithNullConstraint checks if the revocation data is
// fresh against its ThisUpdate and NextUpdate time interval.
type RevocationDataFreshCheckWithNullConstraint struct {
	*AbstractRevocationFreshCheck
}

// NewRevocationDataFreshCheckWithNullConstraint is the default constructor.
// Port of RevocationDataFreshCheckWithNullConstraint(I18nProvider, XmlRFC, RevocationWrapper, Date, LevelRule).
func NewRevocationDataFreshCheckWithNullConstraint(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRFC],
	revocationData *diagnostic.RevocationWrapper, validationDate time.Time, constraint policy.LevelRule) *RevocationDataFreshCheckWithNullConstraint {
	c := &RevocationDataFreshCheckWithNullConstraint{
		AbstractRevocationFreshCheck: NewAbstractRevocationFreshCheck(i18nProvider, result, revocationData, validationDate, constraint),
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationDataFreshCheckWithNullConstraint) Process() bool {
	if c.RevocationData != nil && c.RevocationData.NextUpdate() != nil {
		return c.IsThisUpdateTimeAfterValidationTime(c.getMaxFreshness())
	}
	return false
}

// getMaxFreshness ports the protected getMaxFreshness().
func (c *RevocationDataFreshCheckWithNullConstraint) getMaxFreshness() int64 {
	return c.diff(c.RevocationData.NextUpdate(), c.RevocationData.ThisUpdate())
}

// diff ports the private diff(Date, Date), in milliseconds.
func (c *RevocationDataFreshCheckWithNullConstraint) diff(nextUpdate, thisUpdate *time.Time) int64 {
	var nextUpdateTime, thisUpdateTime int64
	if nextUpdate != nil {
		nextUpdateTime = nextUpdate.UnixMilli()
	}
	if thisUpdate != nil {
		thisUpdateTime = thisUpdate.UnixMilli()
	}
	return nextUpdateTime - thisUpdateTime
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationDataFreshCheckWithNullConstraint) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_RFC_IRIF_TUNU
}
