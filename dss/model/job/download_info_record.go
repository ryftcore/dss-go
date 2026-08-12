// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/DownloadInfoRecord.java (DSS 6.5.RC1).
package job

import (
	"time"

	"github.com/utain/esig/dss/model"
)

// DownloadInfoRecord defines a download result record.
type DownloadInfoRecord interface {
	InfoRecord

	// Document returns the downloaded document. Port of getDocument().
	Document() model.DSSDocument
	// LastDownloadAttemptTime is the last time when a download attempt has been proceeded.
	// The zero time.Time stands for Java's null. Port of getLastDownloadAttemptTime().
	LastDownloadAttemptTime() time.Time
}
