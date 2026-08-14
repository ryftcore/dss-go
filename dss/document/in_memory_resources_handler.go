// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/resources/InMemoryResourcesHandler.java (DSS 6.5.RC1).
package document

import (
	"bytes"
	"io"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/signature/resources"
)

// InMemoryResourcesHandler is an in-memory implementation of DSSResourcesHandler. Using this
// type, all the objects handling during document signing will be stored in memory.
//
// NOTE: this type is used as a default implementation in DSS.
type InMemoryResourcesHandler struct {
	AbstractResourcesHandler

	// buffer is the concrete *bytes.Buffer created by buildOutputStream, kept to extract the
	// written bytes in WriteToDSSDocument (Go's io.Writer, unlike Java's OutputStream, exposes
	// no toByteArray()-equivalent, so the concrete type is retained directly rather than
	// type-asserted off AbstractResourcesHandler.os).
	buffer *bytes.Buffer
}

// NewInMemoryResourcesHandler is the default constructor.
func NewInMemoryResourcesHandler() *InMemoryResourcesHandler {
	h := &InMemoryResourcesHandler{}
	h.BuildOutputStream = h.buildOutputStream
	return h
}

// buildOutputStream ports the protected #buildOutputStream.
func (h *InMemoryResourcesHandler) buildOutputStream() (io.Writer, error) {
	h.buffer = new(bytes.Buffer)
	return h.buffer, nil
}

// WriteToDSSDocument ports #writeToDSSDocument.
func (h *InMemoryResourcesHandler) WriteToDSSDocument() (model.DSSDocument, error) {
	// ensures #CreateOutputStream was called before, matching Java's getOutputStream() contract.
	h.OutputStream()
	defer h.Close()

	if h.buffer == nil {
		panic("The OutputStream shall be an implementation of ByteArrayOutputStream class!")
	}
	return model.NewInMemoryDocument(h.buffer.Bytes()), nil
}

// compile-time interface assertion.
var _ resources.DSSResourcesHandler = (*InMemoryResourcesHandler)(nil)
