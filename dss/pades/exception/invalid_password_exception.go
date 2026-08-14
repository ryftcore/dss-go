// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/exception/InvalidPasswordException.java
// (DSS 6.5.RC1).
//
// Java's `extends DSSException` becomes embedding *model.DSSError, per PORTING.md's
// "DSSException -> model.DSSError" rule and the spi/exception/dss_external_resource_exception.go
// precedent.
package exception

import (
	"github.com/utain/esig/dss/model"
)

// InvalidPasswordException is thrown if an invalid password has been provided.
type InvalidPasswordException struct {
	*model.DSSError
}

// NewInvalidPasswordException creates an InvalidPasswordException with a message. Port of
// InvalidPasswordException(String).
func NewInvalidPasswordException(message string) *InvalidPasswordException {
	return &InvalidPasswordException{DSSError: model.NewDSSError(message)}
}
