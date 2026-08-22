// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/asics/DataToSignASiCSWithXAdESHelper.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// DataToSignASiCSWithXAdESHelper is used to get DataToSign for an ASiC-S with XAdES container.
type DataToSignASiCSWithXAdESHelper struct {
	asic.AbstractGetDataToSignASiCS
}

var _ GetDataToSignASiCWithXAdESHelper = (*DataToSignASiCSWithXAdESHelper)(nil)

// NewDataToSignASiCSWithXAdESHelper is the default constructor. Ports
// DataToSignASiCSWithXAdESHelper(Content).
func NewDataToSignASiCSWithXAdESHelper(asicContent *asic.Content) *DataToSignASiCSWithXAdESHelper {
	return &DataToSignASiCSWithXAdESHelper{
		AbstractGetDataToSignASiCS: asic.NewAbstractGetDataToSignASiCS(asicContent),
	}
}

// ToBeSigned ports the @Override getToBeSigned().
func (h *DataToSignASiCSWithXAdESHelper) ToBeSigned() []model.DSSDocument {
	return h.AsicContent().SignedDocuments()
}

// IsOpenDocument ports the @Override isOpenDocument().
func (h *DataToSignASiCSWithXAdESHelper) IsOpenDocument() bool {
	return false
}
