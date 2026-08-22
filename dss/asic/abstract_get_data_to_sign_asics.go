// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/signature/asics/AbstractGetDataToSignASiCS.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.common.signature.asics lands in
// this same Go package (dss/asic).
package asic

// AbstractGetDataToSignASiCS is used to get DataToSign for an ASiC-S container.
type AbstractGetDataToSignASiCS struct {
	AbstractGetDataToSignHelper
}

// NewAbstractGetDataToSignASiCS is the default constructor. Ports
// AbstractGetDataToSignASiCS(ASiCContent).
func NewAbstractGetDataToSignASiCS(asicContent *ASiCContent) AbstractGetDataToSignASiCS {
	return AbstractGetDataToSignASiCS{AbstractGetDataToSignHelper: NewAbstractGetDataToSignHelper(asicContent)}
}
