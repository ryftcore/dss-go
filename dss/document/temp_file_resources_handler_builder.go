// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/resources/TempFileResourcesHandlerBuilder.java (DSS 6.5.RC1).
package document

import (
	"fmt"
	"os"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/signature/resources"
)

const (
	// tempFileResourcesHandlerBuilderDefaultPrefix is the default prefix of a temporary created
	// file.
	tempFileResourcesHandlerBuilderDefaultPrefix = "dss-"

	// tempFileResourcesHandlerBuilderDefaultSuffix is the default suffix of a temporary created
	// file.
	tempFileResourcesHandlerBuilderDefaultSuffix = ".tmp"
)

// TempFileResourcesHandlerBuilder creates a TempFileResourcesHandler storing temporary objects to
// temporary filesystem documents.
type TempFileResourcesHandlerBuilder struct {
	// fileNamePrefix is the prefix (beginning) of a filename to be used for created documents.
	// Default: "dss-".
	fileNamePrefix string

	// fileNameSuffix is the suffix (ending/extension) of a filename to be used for created
	// documents. Default: ".tmp".
	fileNameSuffix string

	// tempFileDirectory is the directory containing created documents. Default: temporary
	// system-dependent location.
	tempFileDirectory string

	// handlers is the cached list of created handlers by the current builder.
	handlers []*TempFileResourcesHandler
}

// NewTempFileResourcesHandlerBuilder is the default constructor instantiating an empty list of
// processed handlers.
func NewTempFileResourcesHandlerBuilder() *TempFileResourcesHandlerBuilder {
	return &TempFileResourcesHandlerBuilder{
		fileNamePrefix:    tempFileResourcesHandlerBuilderDefaultPrefix,
		fileNameSuffix:    tempFileResourcesHandlerBuilderDefaultSuffix,
		tempFileDirectory: os.TempDir(),
	}
}

// SetFileNamePrefix sets the filename prefix (beginning) for created temporary documents.
// Default: "dss-".
func (b *TempFileResourcesHandlerBuilder) SetFileNamePrefix(fileNamePrefix string) *TempFileResourcesHandlerBuilder {
	b.fileNamePrefix = fileNamePrefix
	return b
}

// SetFileNameSuffix sets the filename suffix (ending/extension) for created temporary documents.
// Default: ".tmp".
func (b *TempFileResourcesHandlerBuilder) SetFileNameSuffix(fileNameSuffix string) *TempFileResourcesHandlerBuilder {
	b.fileNameSuffix = fileNameSuffix
	return b
}

// SetTempFileDirectory sets a file directory to be used for storing created documents. Default:
// temporary system-dependent location.
func (b *TempFileResourcesHandlerBuilder) SetTempFileDirectory(tempFileDirectory string) *TempFileResourcesHandlerBuilder {
	b.tempFileDirectory = tempFileDirectory
	return b
}

// CreateResourcesHandler ports #createResourcesHandler. Java's DSSException from a failed
// directory creation, and NewTempFileResourcesHandler's own DSSException, are both unchecked and
// propagate out of this no-error-return interface method as a panic, per PORTING.md.
func (b *TempFileResourcesHandlerBuilder) CreateResourcesHandler() resources.DSSResourcesHandler {
	if _, err := os.Stat(b.tempFileDirectory); err != nil {
		if mkErr := os.MkdirAll(b.tempFileDirectory, 0o700); mkErr != nil {
			panic(model.NewDSSError(fmt.Sprintf(
				"Unable to create TempFileResourcesHandler for a file with name '%s'!", b.tempFileDirectory)))
		}
	}
	handler, err := NewTempFileResourcesHandler(b.fileNamePrefix, b.fileNameSuffix, b.tempFileDirectory)
	if err != nil {
		panic(err)
	}
	b.handlers = append(b.handlers, handler)
	return handler
}

// Clear removes all handlers created by the current builder, as well as temporary files from the
// filesystem. This method is not executed in a normal DSS operating, and should be called on
// user's side when the temporary files are no longer needed.
//
// NOTE: do not forget to preserve the output documents, such as a FileDocument returned by a
// #signDocument() method.
func (b *TempFileResourcesHandlerBuilder) Clear() {
	for _, handler := range b.handlers {
		handler.ForceDelete()
	}
	b.handlers = nil
}

// compile-time interface assertions.
var (
	_ resources.DSSResourcesHandlerBuilder = (*TempFileResourcesHandlerBuilder)(nil)
	_ resources.DSSResourcesHandler        = (*TempFileResourcesHandler)(nil)
)
