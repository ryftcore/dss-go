// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/exception/SecurityConfigurationException.java (DSS 6.5.RC1).
package common

// SecurityConfigurationException is raised to catch and re-throw an exception caused by a
// security feature/attribute definition failure. Ports the checked Exception(Exception)
// class; callers in this package match it with errors.As.
type SecurityConfigurationException struct {
	Cause error
}

// NewSecurityConfigurationException creates a SecurityConfigurationException wrapping cause.
// Ports SecurityConfigurationException(Exception).
func NewSecurityConfigurationException(cause error) *SecurityConfigurationException {
	return &SecurityConfigurationException{Cause: cause}
}

// Error implements the error interface. Java's Throwable(Throwable cause) constructor sets
// getMessage() to cause.toString(); Error() here returns cause.Error() as the closest Go
// equivalent.
func (e *SecurityConfigurationException) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return "SecurityConfigurationException"
}

// Unwrap allows errors.Is/errors.As to reach the wrapped Cause.
func (e *SecurityConfigurationException) Unwrap() error {
	return e.Cause
}
