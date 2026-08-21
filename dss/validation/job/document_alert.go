// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/alerts/DocumentAlert.java (DSS 6.5.RC1).
package job

import (
	"github.com/utain/esig/dss/alert"
	modeljob "github.com/utain/esig/dss/model/job"
)

// DocumentAlert processes events on document processing. D is the current
// modeljob.DocumentInfo, P the parent modeljob.DocumentInfo, mirroring Java's
// "DocumentAlert<D extends DocumentInfo<P>, P extends DocumentInfo<P>>".
type DocumentAlert[D modeljob.DocumentInfo[P], P modeljob.DocumentInfo[P]] struct {
	*alert.AbstractAlert[D]
}

// NewDocumentAlert creates a DocumentAlert from the given detector and handler. Port of the
// default constructor.
func NewDocumentAlert[D modeljob.DocumentInfo[P], P modeljob.DocumentInfo[P]](
	detection alert.AlertDetector[D], handler alert.AlertHandler[D]) *DocumentAlert[D, P] {
	return &DocumentAlert[D, P]{AbstractAlert: alert.NewAbstractAlert(detection, handler)}
}
