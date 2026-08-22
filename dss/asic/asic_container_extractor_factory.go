// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/extract/ASiCContainerExtractorFactory.java
// (DSS 6.5.RC1).
//
// The Java extract sub-package flattens into this Go package.
package asic

import "github.com/ryftcore/dss-go/dss/model"

// ContainerExtractorFactory is used to find and load a corresponding implementation of
// ContainerExtractor for the given DSSDocument ASiC archive.
type ContainerExtractorFactory interface {
	// IsSupported returns whether the format of the given ASiC document is supported by the
	// current ASiCContainerExtractor. Port of isSupported(DSSDocument).
	IsSupported(asicContainer model.DSSDocument) bool

	// Create creates a new ASiCContainerExtractor for the given ZIP-archive container. Port of
	// create(DSSDocument).
	Create(asicContainer model.DSSDocument) ContainerExtractor
}

// asicContainerExtractorFactoryRegistry holds the ContainerExtractorFactory implementations
// registered via RegisterContainerExtractorFactory, consulted in registration order - the Go
// equivalent of Java's ServiceLoader.load(ASiCContainerExtractorFactory.class) iteration (Go has no
// runtime service-provider discovery). The precedent is
// spi/validation/analyzer/document_analyzer_factory.go. The format-specific implementations
// (dss/asic/cades, dss/asic/xades) register themselves here.
var asicContainerExtractorFactoryRegistry []ContainerExtractorFactory

// RegisterContainerExtractorFactory registers an ContainerExtractorFactory to be consulted
// by DefaultASiCContainerExtractorFromDocument. It has no Java counterpart: it replaces the
// META-INF/services provider file each dss-asic-* module ships.
func RegisterContainerExtractorFactory(factory ContainerExtractorFactory) {
	asicContainerExtractorFactoryRegistry = append(asicContainerExtractorFactoryRegistry, factory)
}
