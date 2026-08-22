// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/asics/DataToSignASiCSWithCAdESFromFiles.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// DataToSignASiCSWithCAdESFromFiles generates a DataToSign with ASiC-S with CAdES from files to
// be signed.
type DataToSignASiCSWithCAdESFromFiles struct {
	AbstractGetDataToSignASiCSWithCAdES
}

var _ GetDataToSignASiCWithCAdESHelper = (*DataToSignASiCSWithCAdESFromFiles)(nil)

// NewDataToSignASiCSWithCAdESFromFiles is the default constructor. Ports
// DataToSignASiCSWithCAdESFromFiles(Content).
func NewDataToSignASiCSWithCAdESFromFiles(asicContent *asic.Content) *DataToSignASiCSWithCAdESFromFiles {
	return &DataToSignASiCSWithCAdESFromFiles{
		AbstractGetDataToSignASiCSWithCAdES: NewAbstractGetDataToSignASiCSWithCAdES(asicContent),
	}
}

// ToBeSigned ports the @Override getToBeSigned().
func (h *DataToSignASiCSWithCAdESFromFiles) ToBeSigned() model.DSSDocument {
	return h.AsicContent().SignedDocuments()[0]
}

// DetachedContents ports the @Override getDetachedContents(), which returns
// Collections.emptyList().
func (h *DataToSignASiCSWithCAdESFromFiles) DetachedContents() []model.DSSDocument {
	return []model.DSSDocument{}
}
