// Ported from dss-document/src/main/java/eu/europa/esig/dss/extension/SignedDocumentExtenderFactory.java (DSS 6.5.RC1).
package document

import "github.com/ryftcore/dss-go/dss/model"

// SignedDocumentExtenderFactory is used to analyze the format of the given DSSDocument and
// create a corresponding implementation of SignedDocumentExtender.
type SignedDocumentExtenderFactory interface {
	// IsSupported tests if the current implementation of SignedDocumentExtender supports the
	// given document. Port of #isSupported.
	IsSupported(document model.DSSDocument) bool

	// Create instantiates a SignedDocumentExtender with the given document. Port of #create.
	Create(document model.DSSDocument) SignedDocumentExtender
}
