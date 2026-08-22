// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/asice/DataToSignASiCEWithCAdESHelper.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.cades.signature.asice lands in
// this same Go package (dss/asic/cades).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// DataToSignASiCEWithCAdESHelper generates a DataToSign with ASiC-E with CAdES.
type DataToSignASiCEWithCAdESHelper struct {
	asic.AbstractGetDataToSignHelper

	// toBeSigned is the cached ToBeSigned document.
	toBeSigned model.DSSDocument
}

var _ GetDataToSignASiCWithCAdESHelper = (*DataToSignASiCEWithCAdESHelper)(nil)

// NewDataToSignASiCEWithCAdESHelper is the default constructor. Ports
// DataToSignASiCEWithCAdESHelper(Content, DSSDocument).
func NewDataToSignASiCEWithCAdESHelper(asicContent *asic.Content, toBeSigned model.DSSDocument) *DataToSignASiCEWithCAdESHelper {
	return &DataToSignASiCEWithCAdESHelper{
		AbstractGetDataToSignHelper: asic.NewAbstractGetDataToSignHelper(asicContent),
		toBeSigned:                  toBeSigned,
	}
}

// ToBeSigned ports the @Override getToBeSigned().
func (h *DataToSignASiCEWithCAdESHelper) ToBeSigned() model.DSSDocument {
	return h.toBeSigned
}

// DetachedContents ports the @Override getDetachedContents(), which returns
// Collections.emptyList().
func (h *DataToSignASiCEWithCAdESHelper) DetachedContents() []model.DSSDocument {
	return []model.DSSDocument{}
}
