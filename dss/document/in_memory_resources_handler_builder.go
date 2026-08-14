// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/resources/InMemoryResourcesHandlerBuilder.java (DSS 6.5.RC1).
package document

import "github.com/utain/esig/dss/spi/signature/resources"

// InMemoryResourcesHandlerBuilder creates an InMemoryResourcesHandler to create in-memory
// objects.
//
// NOTE: This implementation is used by default.
type InMemoryResourcesHandlerBuilder struct{}

// NewInMemoryResourcesHandlerBuilder is the default constructor.
func NewInMemoryResourcesHandlerBuilder() *InMemoryResourcesHandlerBuilder {
	return &InMemoryResourcesHandlerBuilder{}
}

// CreateResourcesHandler ports #createResourcesHandler.
func (b *InMemoryResourcesHandlerBuilder) CreateResourcesHandler() resources.DSSResourcesHandler {
	return NewInMemoryResourcesHandler()
}

// compile-time interface assertion.
var _ resources.DSSResourcesHandlerBuilder = (*InMemoryResourcesHandlerBuilder)(nil)
