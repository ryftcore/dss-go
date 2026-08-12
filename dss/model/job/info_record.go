// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/InfoRecord.java (DSS 6.5.RC1).
package job

import "time"

// InfoRecord describes a state of a record.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type InfoRecord interface {
	// IsRefreshNeeded gets if the refresh is needed for an entry. Port of isRefreshNeeded().
	IsRefreshNeeded() bool
	// IsDesynchronized gets if the record is desynchronized. Port of isDesynchronized().
	IsDesynchronized() bool
	// IsSynchronized gets if the record is synchronized. Port of isSynchronized().
	IsSynchronized() bool
	// IsError gets if the error is present for the record. Port of isError().
	IsError() bool
	// IsToBeDeleted gets if the record shall be deleted. Port of isToBeDeleted().
	IsToBeDeleted() bool
	// StatusName gets the record's status name. Port of getStatusName().
	StatusName() string
	// LastStateTransitionTime gets the last time when the state of record has been changed.
	// The zero time.Time stands for Java's null. Port of getLastStateTransitionTime().
	LastStateTransitionTime() time.Time
	// LastSuccessSynchronizationTime gets the last time when the record has been synchronized.
	// The zero time.Time stands for Java's null. Port of getLastSuccessSynchronizationTime().
	LastSuccessSynchronizationTime() time.Time
	// ExceptionMessage gets the exception message for an error state. Port of
	// getExceptionMessage().
	ExceptionMessage() string
	// ExceptionStackTrace gets the exception stack trace for an error state. Port of
	// getExceptionStackTrace().
	ExceptionStackTrace() string
	// ExceptionFirstOccurrenceTime gets the first time when the error occurred. The zero
	// time.Time stands for Java's null. Port of getExceptionFirstOccurrenceTime().
	ExceptionFirstOccurrenceTime() time.Time
	// ExceptionLastOccurrenceTime gets the last time when the error occurred. The zero
	// time.Time stands for Java's null. Port of getExceptionLastOccurrenceTime().
	ExceptionLastOccurrenceTime() time.Time
	// IsResultExist gets if a result exists under the record. Port of isResultExist().
	IsResultExist() bool
}
