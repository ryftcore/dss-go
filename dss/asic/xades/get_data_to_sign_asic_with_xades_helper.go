// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/GetDataToSignASiCWithXAdESHelper.java
// (DSS 6.5.RC1).
//
// Package layout: the Java packages eu.europa.esig.dss.asic.xades,
// eu.europa.esig.dss.asic.xades.signature{,.asice,.asics}, .definition, .extract, .merge,
// .evidencerecord and .validation all land in this single Go package (dss/asic/xades). Use
// sites outside it import it as
// `asicxades "github.com/ryftcore/dss-go/dss/asic/xades"` to avoid clashing with the top-level
// dss/xades package; inside it, the top-level XAdES package is imported as `dssxades`.
package xades

import "github.com/ryftcore/dss-go/dss/model"

// GetDataToSignASiCWithXAdESHelper defines a helper to create a ToBeSigned data for an ASiC
// with XAdES.
type GetDataToSignASiCWithXAdESHelper interface {
	// ToBeSigned returns a list of documents to be signed (XAdES allows signing multiple
	// files). Port of getToBeSigned().
	ToBeSigned() []model.DSSDocument

	// IsOpenDocument returns whether the concerned container represents an OpenDocument type.
	// Port of isOpenDocument().
	IsOpenDocument() bool
}
