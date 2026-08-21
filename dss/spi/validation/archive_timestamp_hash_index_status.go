// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/ArchiveTimestampHashIndexStatus.java (DSS 6.5.RC1).
package validation

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// ArchiveTimestampHashIndexStatus contains information on the validation status of the
// ats-hash-index(-v3) attribute defined within a timestamp of the archive-time-stamp-v3.
type ArchiveTimestampHashIndexStatus struct {
	// version is the version of the ats-hash-index attribute.
	version enumerations.ArchiveTimestampHashIndexVersion
	// errorMessages contains the error messages that occurred during the timestamp's
	// ats-hash-index-v3 attribute validation.
	errorMessages []string
}

// NewArchiveTimestampHashIndexStatus builds an empty status. Port of the default constructor.
func NewArchiveTimestampHashIndexStatus() *ArchiveTimestampHashIndexStatus {
	return &ArchiveTimestampHashIndexStatus{}
}

// Version gets the version of the ats-hash-index attribute used in the archive-time-stamp-v3.
// Port of getVersion().
func (s *ArchiveTimestampHashIndexStatus) Version() enumerations.ArchiveTimestampHashIndexVersion {
	return s.version
}

// SetVersion sets the version of the ats-hash-index attribute used in the
// archive-time-stamp-v3. Port of setVersion(ArchiveTimestampHashIndexVersion).
func (s *ArchiveTimestampHashIndexStatus) SetVersion(version enumerations.ArchiveTimestampHashIndexVersion) {
	s.version = version
}

// ErrorMessages gets the validation errors that occurred on a structural validation of the
// ats-hash-index(-v3) attribute, when applicable. Port of getErrorMessages(), which lazily
// creates the list.
func (s *ArchiveTimestampHashIndexStatus) ErrorMessages() []string {
	if s.errorMessages == nil {
		s.errorMessages = []string{}
	}
	return s.errorMessages
}

// AddErrorMessage adds a new error message in case of an issue on the ats-hash-index(-v3)
// attribute validation. Port of addErrorMessage(String).
func (s *ArchiveTimestampHashIndexStatus) AddErrorMessage(errorMessage string) {
	s.errorMessages = append(s.ErrorMessages(), errorMessage)
}
