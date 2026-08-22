// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/exception/ProtectedDocumentException.java
// (DSS 6.5.RC1).
//
// Java's `extends DSSException` becomes embedding *model.DSSError, per PORTING.md's
// "DSSException -> model.DSSError" rule and the spi/exception/dss_external_resource_exception.go
// precedent.
package exception

import (
	"github.com/ryftcore/dss-go/dss/model"
)

// ProtectedDocumentException is thrown when the document is protected (the requested operation
// is not permitted).
type ProtectedDocumentException struct {
	*model.DSSError
}

// NewProtectedDocumentException creates a ProtectedDocumentException with a message. Port of
// ProtectedDocumentException(String).
func NewProtectedDocumentException(message string) *ProtectedDocumentException {
	return &ProtectedDocumentException{DSSError: model.NewDSSError(message)}
}
