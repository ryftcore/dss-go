// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/asice/DataToSignASiCEWithXAdESHelper.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// DataToSignASiCEWithXAdESHelper generates a DataToSign for ASiC-E with XAdES.
type DataToSignASiCEWithXAdESHelper struct {
	asic.AbstractGetDataToSignHelper
}

var _ GetDataToSignASiCWithXAdESHelper = (*DataToSignASiCEWithXAdESHelper)(nil)

// NewDataToSignASiCEWithXAdESHelper is the default constructor. Ports
// DataToSignASiCEWithXAdESHelper(ASiCContent).
func NewDataToSignASiCEWithXAdESHelper(asicContent *asic.ASiCContent) *DataToSignASiCEWithXAdESHelper {
	return &DataToSignASiCEWithXAdESHelper{
		AbstractGetDataToSignHelper: asic.NewAbstractGetDataToSignHelper(asicContent),
	}
}

// ToBeSigned ports the @Override getToBeSigned().
func (h *DataToSignASiCEWithXAdESHelper) ToBeSigned() []model.DSSDocument {
	return h.AsicContent().SignedDocuments()
}

// IsOpenDocument ports the @Override isOpenDocument().
func (h *DataToSignASiCEWithXAdESHelper) IsOpenDocument() bool {
	return false
}
