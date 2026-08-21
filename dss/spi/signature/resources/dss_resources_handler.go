// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/signature/resources/DSSResourcesHandler.java (DSS 6.5.RC1).
package resources

import (
	"io"

	"github.com/ryftcore/dss-go/dss/model"
)

// DSSResourcesHandler is used to create objects required for a document
// signing process (e.g. temporary OutputStream, returned DSSDocument, etc.).
// Ports the Java Closeable interface via io.Closer.
type DSSResourcesHandler interface {
	io.Closer

	// CreateOutputStream creates a new io.Writer to be used as an output
	// for a temporary signature document. Ports #createOutputStream.
	CreateOutputStream() (io.Writer, error)

	// WriteToDSSDocument creates a new DSSDocument representing a signed
	// document, based on the created output stream. Ports
	// #writeToDSSDocument.
	WriteToDSSDocument() (model.DSSDocument, error)
}
