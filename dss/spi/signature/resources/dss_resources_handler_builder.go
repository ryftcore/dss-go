// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/signature/resources/DSSResourcesHandlerBuilder.java (DSS 6.5.RC1).
package resources

// DSSResourcesHandlerBuilder builds a new instance of DSSResourcesHandler.
type DSSResourcesHandlerBuilder interface {
	// CreateResourcesHandler instantiates the corresponding factory. Ports
	// #createResourcesHandler.
	CreateResourcesHandler() DSSResourcesHandler
}
