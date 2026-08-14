// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/resources/AbstractResourcesHandler.java (DSS 6.5.RC1).
package document

import "io"

// AbstractResourcesHandler is used to create OutputStreams (io.Writer here) to be used across DSS
// code. Go has no method overriding across an embedded base, so Java's protected abstract
// buildOutputStream() becomes the exported BuildOutputStream function field, following the
// template-method-as-struct-field idiom already established by spi.DSSSecurityFactory
// (spi/dss_security_factory.go): embedding constructors (InMemoryResourcesHandler,
// TempFileResourcesHandler, both in this same package) set it before use.
type AbstractResourcesHandler struct {
	// BuildOutputStream builds the concrete io.Writer this handler wraps. Must be set by the
	// embedding constructor before CreateOutputStream is called.
	BuildOutputStream func() (io.Writer, error)

	// os is the created OutputStream instance.
	os io.Writer
}

// CreateOutputStream ports #createOutputStream. Panics if the OutputStream has already been
// created (Java's unchecked IllegalStateException).
func (h *AbstractResourcesHandler) CreateOutputStream() (io.Writer, error) {
	if h.os != nil {
		panic("Cannot create OutputStream! The OutputStream has been already created!")
	}
	os, err := h.BuildOutputStream()
	if err != nil {
		return nil, err
	}
	h.os = os
	return os, nil
}

// OutputStream returns the internal OutputStream instance. Port of the protected
// #getOutputStream. Panics if #CreateOutputStream has not been called yet (Java's unchecked
// IllegalStateException).
func (h *AbstractResourcesHandler) OutputStream() io.Writer {
	if h.os == nil {
		panic("Method #createOutputStream() shall be called before!")
	}
	return h.os
}

// Close ports #close: closes the underlying OutputStream when it supports closing.
func (h *AbstractResourcesHandler) Close() error {
	if h.os == nil {
		return nil
	}
	if closer, ok := h.os.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
