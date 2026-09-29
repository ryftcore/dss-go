// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/SecureContainerHandlerBuilder.java
// (DSS 6.5.RC1).
package asic

import (
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
)

// SecureContainerHandlerBuilder is the default implementation of a builder, building a new
// instance of SecureContainerHandler.
type SecureContainerHandlerBuilder struct {
	// threshold is the minimum file size to be analyzed on zip bombing.
	threshold int64

	// maxCompressionRatio is the maximum compression ratio.
	maxCompressionRatio int64

	// maxAllowedFilesAmount defines the maximal amount of files that can be inside a ZIP
	// container.
	maxAllowedFilesAmount int

	// maxMalformedFiles is the max iteration over the zip entries.
	maxMalformedFiles int

	// extractComments defines whether comments of ZIP entries shall be extracted.
	// Default : false (not extracted)
	extractComments bool

	// resourcesHandlerBuilder is the builder to be used to create a new DSSResourcesHandler for
	// each internal call, defining a way working with internal resources (e.g. in memory or by
	// using temporary files). The resources are used on a document creation.
	//
	// Default : document.InMemoryResourcesHandler, working with data in memory
	resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder
}

var _ ZipContainerHandlerBuilder = (*SecureContainerHandlerBuilder)(nil)

// NewSecureContainerHandlerBuilder is the default constructor; the Java field initializers become
// the returned struct's values.
func NewSecureContainerHandlerBuilder() *SecureContainerHandlerBuilder {
	return &SecureContainerHandlerBuilder{
		threshold:               1000000, // 1 MB
		maxCompressionRatio:     100,
		maxAllowedFilesAmount:   1000,
		maxMalformedFiles:       100,
		extractComments:         false,
		resourcesHandlerBuilder: document.NewInMemoryResourcesHandlerBuilder(),
	}
}

// SetThreshold sets the maximum allowed threshold after exceeding each the security checks are
// enforced.
//
// Default : 1000000 (1 MB)
//
// Port of setThreshold(long).
func (b *SecureContainerHandlerBuilder) SetThreshold(threshold int64) *SecureContainerHandlerBuilder {
	b.threshold = threshold
	return b
}

// SetMaxCompressionRatio sets the maximum allowed compression ratio. If the container compression
// ratio exceeds the value, an error is being returned.
//
// Default : 100
//
// Port of setMaxCompressionRatio(long).
func (b *SecureContainerHandlerBuilder) SetMaxCompressionRatio(maxCompressionRatio int64) *SecureContainerHandlerBuilder {
	b.maxCompressionRatio = maxCompressionRatio
	return b
}

// SetMaxAllowedFilesAmount sets the maximum allowed amount of files inside a container.
//
// Default : 1000
//
// Port of setMaxAllowedFilesAmount(int).
func (b *SecureContainerHandlerBuilder) SetMaxAllowedFilesAmount(maxAllowedFilesAmount int) *SecureContainerHandlerBuilder {
	b.maxAllowedFilesAmount = maxAllowedFilesAmount
	return b
}

// SetMaxMalformedFiles sets the maximum allowed amount of malformed files.
//
// Default : 100
//
// Port of setMaxMalformedFiles(int).
func (b *SecureContainerHandlerBuilder) SetMaxMalformedFiles(maxMalformedFiles int) *SecureContainerHandlerBuilder {
	b.maxMalformedFiles = maxMalformedFiles
	return b
}

// SetExtractComments sets whether comments of ZIP entries shall be extracted.
//
// Enabling of the feature can be useful when editing an existing archive, in order to preserve the
// existing data (i.e. comments). See SecureContainerHandler.SetExtractComments for the reason it
// is off by default and for its FileDocument-only restriction.
//
// Default : false (not extracted)
//
// Port of setExtractComments(boolean).
func (b *SecureContainerHandlerBuilder) SetExtractComments(extractComments bool) *SecureContainerHandlerBuilder {
	b.extractComments = extractComments
	return b
}

// SetResourcesHandlerBuilder sets the DSSResourcesHandlerBuilder to be used for a
// DSSResourcesHandler creation in internal methods. DSSResourcesHandler defines a way to operate
// with OutputStreams and create DSSDocuments.
//
// Default : document.InMemoryResourcesHandler. Works with data in memory, so creating or merging a
// container holds the whole output in memory; pass a document.TempFileResourcesHandlerBuilder to
// keep large containers on disk instead.
//
// Panics with the Java message when resourcesHandlerBuilder is nil (Objects.requireNonNull).
//
// Port of setResourcesHandlerBuilder(DSSResourcesHandlerBuilder).
func (b *SecureContainerHandlerBuilder) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) *SecureContainerHandlerBuilder {
	if resourcesHandlerBuilder == nil {
		panic("DSSResourcesHandlerBuilder cannot be null!")
	}
	b.resourcesHandlerBuilder = resourcesHandlerBuilder
	return b
}

// Build ports build(). Java narrows the return type to SecureContainerHandler through the
// ZipContainerHandlerBuilder<T> type parameter, which this port drops (see
// zip_container_handler_builder.go).
func (b *SecureContainerHandlerBuilder) Build() ZipContainerHandler {
	secureContainerHandler := NewSecureContainerHandler()
	secureContainerHandler.SetThreshold(b.threshold)
	secureContainerHandler.SetMaxCompressionRatio(b.maxCompressionRatio)
	secureContainerHandler.SetMaxAllowedFilesAmount(b.maxAllowedFilesAmount)
	secureContainerHandler.SetMaxMalformedFiles(b.maxMalformedFiles)
	secureContainerHandler.SetExtractComments(b.extractComments)
	secureContainerHandler.SetResourcesHandlerBuilder(b.resourcesHandlerBuilder)
	return secureContainerHandler
}
