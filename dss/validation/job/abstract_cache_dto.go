// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/dto/AbstractCacheDTO.java (DSS 6.5.RC1).
package job

import (
	"time"

	modeljob "github.com/ryftcore/dss-go/dss/model/job"
)

// AbstractCacheDTO is the abstract cache DTO. It implements modeljob.InfoRecord.
type AbstractCacheDTO struct {
	// cacheState is the state of the record.
	cacheState CacheStateEnum

	// lastStateTransitionTime is the last time of the state change. The zero time.Time
	// stands for Java's null.
	lastStateTransitionTime time.Time

	// lastSuccessSynchronizationTime is the last time of a successful synchronization. The
	// zero time.Time stands for Java's null.
	lastSuccessSynchronizationTime time.Time

	// exceptionMessage is the exception message.
	exceptionMessage string

	// exceptionStackTrace is the exception stack trace.
	exceptionStackTrace string

	// exceptionFirstOccurrenceTime is the first time of the exception occurrence. The zero
	// time.Time stands for Java's null.
	exceptionFirstOccurrenceTime time.Time

	// exceptionLastOccurrenceTime is the last time of the exception occurrence. The zero
	// time.Time stands for Java's null.
	exceptionLastOccurrenceTime time.Time

	// resultExist defines if the result exists.
	resultExist bool
}

// NewAbstractCacheDTO creates an empty AbstractCacheDTO. Port of the empty constructor.
func NewAbstractCacheDTO() *AbstractCacheDTO {
	return &AbstractCacheDTO{}
}

// NewAbstractCacheDTOFrom copies cacheDTO. Port of the copy constructor.
func NewAbstractCacheDTOFrom(cacheDTO *AbstractCacheDTO) *AbstractCacheDTO {
	return &AbstractCacheDTO{
		cacheState:                     cacheDTO.cacheState,
		lastStateTransitionTime:        cacheDTO.lastStateTransitionTime,
		lastSuccessSynchronizationTime: cacheDTO.lastSuccessSynchronizationTime,
		exceptionMessage:               cacheDTO.exceptionMessage,
		exceptionStackTrace:            cacheDTO.exceptionStackTrace,
		exceptionFirstOccurrenceTime:   cacheDTO.exceptionFirstOccurrenceTime,
		exceptionLastOccurrenceTime:    cacheDTO.exceptionLastOccurrenceTime,
		resultExist:                    cacheDTO.resultExist,
	}
}

// CacheState gets the state of the cache. Port of getCacheState().
func (d *AbstractCacheDTO) CacheState() CacheStateEnum {
	return d.cacheState
}

// SetCacheState sets the cache state. Port of setCacheState(CacheStateEnum).
func (d *AbstractCacheDTO) SetCacheState(cacheState CacheStateEnum) {
	d.cacheState = cacheState
}

// LastStateTransitionTime returns the last time of the state change. Port of
// getLastStateTransitionTime().
func (d *AbstractCacheDTO) LastStateTransitionTime() time.Time {
	return d.lastStateTransitionTime
}

// SetLastStateTransitionTime sets the last time of the state change. Port of
// setLastStateTransitionTime(Date).
func (d *AbstractCacheDTO) SetLastStateTransitionTime(lastStateTransitionTime time.Time) {
	d.lastStateTransitionTime = lastStateTransitionTime
}

// LastSuccessSynchronizationTime returns the last time of a successful synchronization.
// Port of getLastSuccessSynchronizationTime().
func (d *AbstractCacheDTO) LastSuccessSynchronizationTime() time.Time {
	return d.lastSuccessSynchronizationTime
}

// SetLastSuccessSynchronizationTime sets the last time of a successful synchronization.
// Port of setLastSuccessSynchronizationTime(Date).
func (d *AbstractCacheDTO) SetLastSuccessSynchronizationTime(lastSuccessSynchronizationTime time.Time) {
	d.lastSuccessSynchronizationTime = lastSuccessSynchronizationTime
}

// ExceptionMessage returns the exception message. Port of getExceptionMessage().
func (d *AbstractCacheDTO) ExceptionMessage() string {
	return d.exceptionMessage
}

// SetExceptionMessage sets the exception message. Port of setExceptionMessage(String).
func (d *AbstractCacheDTO) SetExceptionMessage(exceptionMessage string) {
	d.exceptionMessage = exceptionMessage
}

// ExceptionStackTrace returns the exception stack trace. Port of getExceptionStackTrace().
func (d *AbstractCacheDTO) ExceptionStackTrace() string {
	return d.exceptionStackTrace
}

// SetExceptionStackTrace sets the exception stack trace. Port of
// setExceptionStackTrace(String).
func (d *AbstractCacheDTO) SetExceptionStackTrace(exceptionStackTrace string) {
	d.exceptionStackTrace = exceptionStackTrace
}

// ExceptionFirstOccurrenceTime returns the first time of the exception occurrence. Port of
// getExceptionFirstOccurrenceTime().
func (d *AbstractCacheDTO) ExceptionFirstOccurrenceTime() time.Time {
	return d.exceptionFirstOccurrenceTime
}

// SetExceptionFirstOccurrenceTime sets the first time of the exception occurrence. Port of
// setExceptionFirstOccurrenceTime(Date).
func (d *AbstractCacheDTO) SetExceptionFirstOccurrenceTime(exceptionFirstOccurrenceTime time.Time) {
	d.exceptionFirstOccurrenceTime = exceptionFirstOccurrenceTime
}

// ExceptionLastOccurrenceTime returns the last time of a the exception occurrence. Port of
// getExceptionLastOccurrenceTime().
func (d *AbstractCacheDTO) ExceptionLastOccurrenceTime() time.Time {
	return d.exceptionLastOccurrenceTime
}

// SetExceptionLastOccurrenceTime sets the last time of a the exception occurrence. Port of
// setExceptionLastOccurrenceTime(Date).
func (d *AbstractCacheDTO) SetExceptionLastOccurrenceTime(exceptionLastOccurrenceTime time.Time) {
	d.exceptionLastOccurrenceTime = exceptionLastOccurrenceTime
}

// IsResultExist returns if the cache result exists. Port of isResultExist().
func (d *AbstractCacheDTO) IsResultExist() bool {
	return d.resultExist
}

// SetResultExist sets if the cache result exists. Port of setResultExist(boolean).
func (d *AbstractCacheDTO) SetResultExist(resultExist bool) {
	d.resultExist = resultExist
}

// IsRefreshNeeded reports whether the cache state is REFRESH_NEEDED. Port of
// isRefreshNeeded().
func (d *AbstractCacheDTO) IsRefreshNeeded() bool {
	return CacheStateEnumRefreshNeeded == d.cacheState
}

// IsDesynchronized reports whether the cache state is DESYNCHRONIZED. Port of
// isDesynchronized().
func (d *AbstractCacheDTO) IsDesynchronized() bool {
	return CacheStateEnumDesynchronized == d.cacheState
}

// IsSynchronized reports whether the cache state is SYNCHRONIZED. Port of isSynchronized().
func (d *AbstractCacheDTO) IsSynchronized() bool {
	return CacheStateEnumSynchronized == d.cacheState
}

// IsError reports whether the cache state is ERROR. Port of isError().
func (d *AbstractCacheDTO) IsError() bool {
	return CacheStateEnumError == d.cacheState
}

// IsToBeDeleted reports whether the cache state is TO_BE_DELETED. Port of isToBeDeleted().
func (d *AbstractCacheDTO) IsToBeDeleted() bool {
	return CacheStateEnumToBEDeleted == d.cacheState
}

// StatusName returns the cache state's name. Port of getStatusName().
func (d *AbstractCacheDTO) StatusName() string {
	return string(d.cacheState)
}

var _ modeljob.InfoRecord = (*AbstractCacheDTO)(nil)
