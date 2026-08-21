// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qscd/QSCDStrategy.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/enumerations"

// QSCDStrategy is used to extract QSCD status.
type QSCDStrategy interface {
	// QSCDStatus gets QSCD status. Port of getQSCDStatus().
	QSCDStatus() enumerations.QSCDStatus
}
