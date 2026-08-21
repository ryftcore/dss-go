// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/asics/AbstractGetDataToSignASiCSWithCAdES.java (DSS 6.5.RC1).
//
// Package flattening: the Java package eu.europa.esig.dss.asic.cades.signature.asics lands in
// this same Go package (dss/asic/cades) per S7_BRIEF.md's package layout table.
package cades

import "github.com/ryftcore/dss-go/dss/asic"

// AbstractGetDataToSignASiCSWithCAdES generates a DataToSign with ASiC-S with CAdES.
type AbstractGetDataToSignASiCSWithCAdES struct {
	asic.AbstractGetDataToSignASiCS
}

// NewAbstractGetDataToSignASiCSWithCAdES is the default constructor. Ports the protected
// AbstractGetDataToSignASiCSWithCAdES(ASiCContent).
func NewAbstractGetDataToSignASiCSWithCAdES(asicContent *asic.ASiCContent) AbstractGetDataToSignASiCSWithCAdES {
	return AbstractGetDataToSignASiCSWithCAdES{
		AbstractGetDataToSignASiCS: asic.NewAbstractGetDataToSignASiCS(asicContent),
	}
}
