// Ported from dss-validation/.../validation/process/bbb/fc/checks/SignedAndTimestampedFilesCoveredCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignedAndTimestampedFilesCoveredCheck checks whether all files signed by the covered
// signatures or timestamped by covered timestamps are covered by the current timestamp as well.
type SignedAndTimestampedFilesCoveredCheck struct {
	AbstractSignedAndTimestampedFilesCoveredCheck[*drjaxb.XmlFC]
}

// NewSignedAndTimestampedFilesCoveredCheck is the default constructor.
func NewSignedAndTimestampedFilesCoveredCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	diagnosticData *diagnostic.DiagnosticData, timestampWrapper *diagnostic.TimestampWrapper,
	constraint policy.LevelRule) *SignedAndTimestampedFilesCoveredCheck {
	c := &SignedAndTimestampedFilesCoveredCheck{}
	c.InitAbstractSignedAndTimestampedFilesCoveredCheck(i18nProvider, result, diagnosticData,
		timestampWrapper.Filename(), constraint)
	c.InitChainItem(c)
	return c
}
