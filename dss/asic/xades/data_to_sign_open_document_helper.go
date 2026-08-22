// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/asice/DataToSignOpenDocumentHelper.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// DataToSignOpenDocumentHelper generates a DataToSign for an OpenDocument signing.
//
// Java's `extends DataToSignASiCEWithXAdESHelper` becomes struct embedding; nothing calls back
// into this type virtually through the base, so no overrides interface is needed here (unlike
// the abstract base classes elsewhere in this phase) - callers always hold the concrete
// *DataToSignOpenDocumentHelper (or the GetDataToSignASiCWithXAdESHelper interface it
// satisfies), never the embedded DataToSignASiCEWithXAdESHelper.
type DataToSignOpenDocumentHelper struct {
	DataToSignASiCEWithXAdESHelper
}

var _ GetDataToSignASiCWithXAdESHelper = (*DataToSignOpenDocumentHelper)(nil)

// NewDataToSignOpenDocumentHelper is the default constructor. Ports
// DataToSignOpenDocumentHelper(Content).
func NewDataToSignOpenDocumentHelper(asicContent *asic.Content) *DataToSignOpenDocumentHelper {
	return &DataToSignOpenDocumentHelper{
		DataToSignASiCEWithXAdESHelper: *NewDataToSignASiCEWithXAdESHelper(asicContent),
	}
}

// ToBeSigned ports the @Override getToBeSigned().
func (h *DataToSignOpenDocumentHelper) ToBeSigned() []model.DSSDocument {
	return OpenDocumentSupportUtilsGetOpenDocumentCoverage(h.AsicContent())
}

// IsOpenDocument ports the @Override isOpenDocument().
func (h *DataToSignOpenDocumentHelper) IsOpenDocument() bool {
	return true
}
