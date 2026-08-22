// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/alerts/ValidationJobAlerter.java (DSS 6.5.RC1).
package job

import (
	"github.com/ryftcore/dss-go/dss/alert"
	modeljob "github.com/ryftcore/dss-go/dss/model/job"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ValidationJobAlerter processes alerts on ValidationJob. D is the current
// modeljob.DocumentInfo, L the parent modeljob.DocumentListInfo, mirroring Java's
// "ValidationJobAlerter<D extends DocumentInfo<L>, L extends DocumentListInfo<L, D>>".
//
// slf4j warn logging is dropped (execute's catch-and-log becomes catch-and-ignore, see
// execute's DEVIATION note; alert firing itself is preserved as it is load-bearing).
type ValidationJobAlerter[D modeljob.DocumentInfo[L], L modeljob.DocumentListInfo[L, D]] struct {
	// documentListAlerts holds alerts to be applied on document list changes.
	documentListAlerts []alert.Alert[L]

	// documentAlerts holds alerts to be applied on document changes.
	documentAlerts []alert.Alert[D]
}

// NewValidationJobAlerter creates a ValidationJobAlerter. Port of the constructor.
func NewValidationJobAlerter[D modeljob.DocumentInfo[L], L modeljob.DocumentListInfo[L, D]](
	documentListAlerts []alert.Alert[L], documentAlerts []alert.Alert[D]) *ValidationJobAlerter[D, L] {
	return &ValidationJobAlerter[D, L]{documentListAlerts: documentListAlerts, documentAlerts: documentAlerts}
}

// DetectChanges runs alerts on the given ValidationJobSummary. Port of
// detectChanges(ValidationJobSummary).
func (a *ValidationJobAlerter[D, L]) DetectChanges(jobSummary modeljob.ValidationJobSummary[D, L]) {
	for _, docListInfo := range jobSummary.DocumentListInfos() {
		// run document list alerts
		if utils.IsCollectionNotEmpty(a.documentListAlerts) {
			for _, docListAlert := range a.documentListAlerts {
				executeDocumentListAlert(docListAlert, docListInfo)
			}
		}
		// run document alerts
		if utils.IsCollectionNotEmpty(a.documentAlerts) {
			for _, docInfo := range docListInfo.ChildrenInfos() {
				for _, docAlert := range a.documentAlerts {
					executeDocumentAlert(docAlert, docInfo)
				}
			}
		}
	}
	// other documents
	if utils.IsCollectionNotEmpty(a.documentAlerts) {
		for _, docInfo := range jobSummary.OtherDocumentInfos() {
			for _, docAlert := range a.documentAlerts {
				executeDocumentAlert(docAlert, docInfo)
			}
		}
	}
}

// executeDocumentListAlert runs alert on info, swallowing any error. Port of the private
// generic execute(Alert<T>, T) as applied to the document-list branch.
//
// DEVIATION: Java logs the swallowed exception (info.getDSSId().asXmlId() and the error
// message); slf4j logging is dropped per the porting convention for this module, so the
// error is discarded silently rather than logged.
func executeDocumentListAlert[L any](alrt alert.Alert[L], info L) {
	_ = alrt.Alert(info)
}

// executeDocumentAlert runs alert on info, swallowing any error. See
// executeDocumentListAlert for the DEVIATION note (applies identically to the document
// branch).
func executeDocumentAlert[D any](alrt alert.Alert[D], info D) {
	_ = alrt.Alert(info)
}
