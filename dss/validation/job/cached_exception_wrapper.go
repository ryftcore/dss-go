// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/cache/state/CachedExceptionWrapper.java (DSS 6.5.RC1).
package job

import "time"

// CachedExceptionWrapper wraps an exception for a cache record.
//
// DEVIATION: Java's getStackTrace() renders the wrapped exception's full stack trace via
// printStackTrace(); Go errors carry no stack trace (see PORTING.md), so StackTrace returns
// the wrapped error's Error() message instead, matching the precedent set by
// DSSDataLoaderMultipleException.
type CachedExceptionWrapper struct {
	// date is the first occurrence date of the exception.
	date time.Time

	// exception is the wrapped error.
	exception error

	// lastOccurrenceDate is the last occurrence date of the exception.
	lastOccurrenceDate time.Time
}

// NewCachedExceptionWrapper creates a CachedExceptionWrapper around exception. Panics if
// exception is nil, mirroring Java's Objects.requireNonNull(exception).
func NewCachedExceptionWrapper(exception error) *CachedExceptionWrapper {
	if exception == nil {
		panic("exception cannot be nil")
	}
	now := time.Now()
	return &CachedExceptionWrapper{
		date:               now,
		exception:          exception,
		lastOccurrenceDate: now,
	}
}

// Date returns the first occurrence date of the exception. Port of getDate().
func (c *CachedExceptionWrapper) Date() time.Time {
	return c.date
}

// LastOccurrenceDate returns the last occurrence date of the exception. Port of
// getLastOccurrenceDate().
func (c *CachedExceptionWrapper) LastOccurrenceDate() time.Time {
	return c.lastOccurrenceDate
}

// SetLastOccurrenceDate sets the last occurrence date of the exception. Port of
// setLastOccurrenceDate(Date).
func (c *CachedExceptionWrapper) SetLastOccurrenceDate(lastOccurrenceDate time.Time) {
	c.lastOccurrenceDate = lastOccurrenceDate
}

// Exception returns the wrapped exception. Port of getException().
func (c *CachedExceptionWrapper) Exception() error {
	return c.exception
}

// ExceptionMessage returns the exception message. Port of getExceptionMessage().
func (c *CachedExceptionWrapper) ExceptionMessage() string {
	return c.exception.Error()
}

// StackTrace returns the exception stack trace. Port of getStackTrace(). See the
// DEVIATION note on the type for the Go substitute used here.
func (c *CachedExceptionWrapper) StackTrace() string {
	return c.exception.Error()
}
