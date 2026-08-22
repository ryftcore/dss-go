// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/asics/DataToSignASiCSWithCAdESFromArchive.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// DataToSignASiCSWithCAdESFromArchive generates a DataToSign with ASiC-S with CAdES from an
// existing archive.
type DataToSignASiCSWithCAdESFromArchive struct {
	AbstractGetDataToSignASiCSWithCAdES
}

var _ GetDataToSignASiCWithCAdESHelper = (*DataToSignASiCSWithCAdESFromArchive)(nil)

// NewDataToSignASiCSWithCAdESFromArchive is the default constructor. Ports
// DataToSignASiCSWithCAdESFromArchive(Content).
func NewDataToSignASiCSWithCAdESFromArchive(asicContent *asic.Content) *DataToSignASiCSWithCAdESFromArchive {
	return &DataToSignASiCSWithCAdESFromArchive{
		AbstractGetDataToSignASiCSWithCAdES: NewAbstractGetDataToSignASiCSWithCAdES(asicContent),
	}
}

// ToBeSigned ports the @Override getToBeSigned().
//
// NOTE: in ASiC-S signatures are added within the same signature file, and handling of detached
// document signing is delegated to the CAdES service.
//
// Panics with a *model.DSSError carrying Java's DSSException message when the embedded
// signature cannot be selected.
func (h *DataToSignASiCSWithCAdESFromArchive) ToBeSigned() model.DSSDocument {
	embeddedSignatures := h.AsicContent().SignatureDocuments()
	nbEmbeddedSignatures := utils.CollectionSize(embeddedSignatures)
	if nbEmbeddedSignatures != 1 {
		panic(model.NewDSSError(fmt.Sprintf("Unable to select the embedded signature (nb found:%d)", nbEmbeddedSignatures)))
	}
	return embeddedSignatures[0]
}

// DetachedContents ports the @Override getDetachedContents().
//
// Panics with a *model.DSSError carrying Java's DSSException message when the document to be
// signed cannot be selected.
func (h *DataToSignASiCSWithCAdESFromArchive) DetachedContents() []model.DSSDocument {
	embeddedSignedFiles := h.AsicContent().SignedDocuments()
	nbSignedFiles := utils.CollectionSize(embeddedSignedFiles)
	if nbSignedFiles != 1 {
		panic(model.NewDSSError(fmt.Sprintf("Unable to select the document to be signed (nb found:%d)", nbSignedFiles)))
	}
	return embeddedSignedFiles
}
