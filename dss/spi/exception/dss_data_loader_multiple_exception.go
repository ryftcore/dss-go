// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/exception/DSSDataLoaderMultipleException.java (DSS 6.5.RC1).
package exception

import (
	"sort"
	"strings"
)

// DSSDataLoaderMultipleException aggregates the errors encountered while attempting a data
// loader request across multiple URLs. Ports DSSDataLoaderMultipleException (extends
// DSSExternalResourceException, overriding getMessage() to flatten the per-URL causes).
//
// DEVIATION: Java's getStackTrace() override concatenates every wrapped exception's stack
// trace; Go errors carry no stack trace to concatenate (see PORTING.md's error-handling
// conventions), so that override has no port. Callers needing individual failures use
// URLExceptionMap directly, and Unwrap (Go 1.20+ multi-error form) exposes them to
// errors.Is/errors.As.
type DSSDataLoaderMultipleException struct {
	*DSSExternalResourceException

	// URLExceptionMap maps each failed URL to the error encountered fetching it.
	URLExceptionMap map[string]error
}

// NewDSSDataLoaderMultipleException creates a DSSDataLoaderMultipleException from a map
// between failed URLs and the errors that caused them to fail. Port of
// DSSDataLoaderMultipleException(Map<String, Throwable>).
func NewDSSDataLoaderMultipleException(urlExceptionMap map[string]error) *DSSDataLoaderMultipleException {
	return &DSSDataLoaderMultipleException{
		DSSExternalResourceException: newDSSExternalResourceExceptionEmpty(),
		URLExceptionMap:              urlExceptionMap,
	}
}

// Error implements the error interface, overriding the embedded DSSExternalResourceException's
// message with the flattened per-URL failure report. Port of getMessage().
//
// DEVIATION: Java iterates urlExceptionMap in its own (LinkedHashMap or hash) order; Go map
// iteration order is unspecified, so URLs are sorted for deterministic output.
func (e *DSSDataLoaderMultipleException) Error() string {
	if e == nil {
		return ""
	}
	urls := make([]string, 0, len(e.URLExceptionMap))
	for url := range e.URLExceptionMap {
		urls = append(urls, url)
	}
	sort.Strings(urls)

	var b strings.Builder
	for _, url := range urls {
		cause := e.URLExceptionMap[url]
		errorMessage := ""
		if cause != nil {
			errorMessage = cause.Error()
		}
		if extRes, ok := cause.(*DSSExternalResourceException); ok {
			errorMessage = extRes.causeMessage()
		}
		b.WriteString("Failed to get data from URL '")
		b.WriteString(url)
		b.WriteString("'. Reason : [")
		b.WriteString(errorMessage)
		b.WriteString("]. ")
	}
	return b.String()
}

// Unwrap exposes every wrapped per-URL error to errors.Is/errors.As (Go 1.20+ multi-error
// form), standing in for getStackTrace()'s aggregation of the underlying exceptions.
func (e *DSSDataLoaderMultipleException) Unwrap() []error {
	if e == nil {
		return nil
	}
	errs := make([]error, 0, len(e.URLExceptionMap))
	for _, err := range e.URLExceptionMap {
		errs = append(errs, err)
	}
	return errs
}

// causeMessage returns the flattened message, mirroring
// DSSDataLoaderMultipleException#getCauseMessage() (package-private override returning
// getMessage()).
func (e *DSSDataLoaderMultipleException) causeMessage() string {
	return e.Error()
}
