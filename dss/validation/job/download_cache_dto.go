// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/dto/DownloadCacheDTO.java (DSS 6.5.RC1).
package job

import (
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	modeljob "github.com/ryftcore/dss-go/dss/model/job"
)

// DownloadCacheDTO is the download record DTO. It implements modeljob.DownloadInfoRecord.
type DownloadCacheDTO struct {
	*AbstractCacheDTO

	// document is the downloaded document.
	document model.DSSDocument

	// sha2ErrorMessages holds error messages occurred during sha2 processing.
	sha2ErrorMessages []string
}

// NewDownloadCacheDTO creates an empty DownloadCacheDTO. Port of the empty constructor.
func NewDownloadCacheDTO() *DownloadCacheDTO {
	return &DownloadCacheDTO{AbstractCacheDTO: NewAbstractCacheDTO()}
}

// NewDownloadCacheDTOFrom copies cacheDTO. Port of the copy constructor.
func NewDownloadCacheDTOFrom(cacheDTO *AbstractCacheDTO) *DownloadCacheDTO {
	return &DownloadCacheDTO{AbstractCacheDTO: NewAbstractCacheDTOFrom(cacheDTO)}
}

// LastDownloadAttemptTime returns the maximum of LastSuccessSynchronizationTime,
// ExceptionLastOccurrenceTime and LastStateTransitionTime. Panics if all three are the zero
// time.Time, mirroring Java's IllegalStateException("All dates are null"). Port of
// getLastDownloadAttemptTime().
func (d *DownloadCacheDTO) LastDownloadAttemptTime() time.Time {
	dates := []time.Time{d.LastSuccessSynchronizationTime(), d.ExceptionLastOccurrenceTime(), d.LastStateTransitionTime()}
	return compareDates(dates)
}

// compareDates returns the maximum non-zero time.Time in dates. Port of
// compareDates(List<Date>), where the zero time.Time stands for Java's null.
func compareDates(dates []time.Time) time.Time {
	var max time.Time
	found := false
	for _, dt := range dates {
		if dt.IsZero() {
			continue
		}
		if !found || dt.After(max) {
			max = dt
			found = true
		}
	}
	if !found {
		panic("All dates are null")
	}
	return max
}

// Document returns the downloaded document. Port of getDocument().
func (d *DownloadCacheDTO) Document() model.DSSDocument {
	return d.document
}

// SetDocument sets the downloaded document. Port of setDocument(DSSDocument).
func (d *DownloadCacheDTO) SetDocument(document model.DSSDocument) {
	d.document = document
}

// Sha2ErrorMessages returns error messages occurred during the sha2 processing. Port of
// getSha2ErrorMessages().
func (d *DownloadCacheDTO) Sha2ErrorMessages() []string {
	return d.sha2ErrorMessages
}

// SetSha2ErrorMessages sets error messages occurred during sha2 file processing. Port of
// setSha2ErrorMessages(List).
func (d *DownloadCacheDTO) SetSha2ErrorMessages(sha2ErrorMessages []string) {
	d.sha2ErrorMessages = sha2ErrorMessages
}

var _ modeljob.DownloadInfoRecord = (*DownloadCacheDTO)(nil)
