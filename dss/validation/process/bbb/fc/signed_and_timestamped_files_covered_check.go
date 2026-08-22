// Ported from dss-validation/.../validation/process/bbb/fc/checks/SignedAndTimestampedFilesCoveredCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignedAndTimestampedFilesCoveredCheck checks whether all files signed by the covered
// signatures or timestamped by covered timestamps are covered by the current timestamp as well.
type SignedAndTimestampedFilesCoveredCheck struct {
	AbstractSignedAndTimestampedFilesCoveredCheck[*drjaxb.XmlFC]
}

// NewSignedAndTimestampedFilesCoveredCheck is the default constructor.
func NewSignedAndTimestampedFilesCoveredCheck(i18nProvider *i18n.Provider, result *process.Result[*drjaxb.XmlFC],
	diagnosticData *diagnostic.Data, timestampWrapper *diagnostic.TimestampWrapper,
	constraint policy.LevelRule) *SignedAndTimestampedFilesCoveredCheck {
	c := &SignedAndTimestampedFilesCoveredCheck{}
	c.InitAbstractSignedAndTimestampedFilesCoveredCheck(i18nProvider, result, diagnosticData,
		timestampWrapper.Filename(), constraint)
	c.InitChainItem(c)
	return c
}
