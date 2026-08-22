// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/summary/ValidationJobSummaryBuilder.java (DSS 6.5.RC1).
package job

import modeljob "github.com/ryftcore/dss-go/dss/model/job"

// ValidationJobSummaryBuilder builds a modeljob.ValidationJobSummary. D is the current
// modeljob.DocumentInfo, L the parent modeljob.DocumentListInfo, mirroring Java's
// "ValidationJobSummaryBuilder<D extends DocumentInfo<L>, L extends DocumentListInfo<L, D>>".
type ValidationJobSummaryBuilder[D modeljob.DocumentInfo[L], L modeljob.DocumentListInfo[L, D]] interface {
	// Build builds the ValidationJobSummary. Port of build().
	Build() modeljob.ValidationJobSummary[D, L]
}
