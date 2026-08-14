// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ZipContainerHandlerBuilder.java
// (DSS 6.5.RC1).
package asic

// ZipContainerHandlerBuilder builds a new instance of ZipContainerHandler.
//
// DEVIATION: Java declares `ZipContainerHandlerBuilder<T extends ZipContainerHandler>` so that
// SecureContainerHandlerBuilder.build() can covariantly narrow its return type to
// SecureContainerHandler. Go has neither covariant returns nor a way to spell Java's
// `ZipContainerHandlerBuilder<?>` wildcard that ZipUtils stores, so the type parameter is
// dropped and Build returns the interface type. Callers needing the concrete handler type
// assert on it.
type ZipContainerHandlerBuilder interface {
	// Build builds a new instance of ZipContainerHandler. Port of build().
	Build() ZipContainerHandler
}
