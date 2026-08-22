// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/GetDataToSignASiCWithCAdESHelper.java (DSS 6.5.RC1).
//
// Package layout: the Java packages eu.europa.esig.dss.asic.cades,
// eu.europa.esig.dss.asic.cades.signature{,.asice,.asics,.manifest},
// eu.europa.esig.dss.asic.cades.timestamp, .extract, .merge, .evidencerecord and .validation all
// land in this single Go package (dss/asic/cades). Use
// sites outside it import it as `asiccades "github.com/ryftcore/dss-go/dss/asic/cades"` to avoid
// clashing with the top-level dss/cades package; inside it, the top-level CAdES package is
// imported as `dsscades`.
package cades

import "github.com/ryftcore/dss-go/dss/model"

// GetDataToSignASiCWithCAdESHelper defines a helper to create a ToBeSigned data for an ASiC
// with CAdES.
type GetDataToSignASiCWithCAdESHelper interface {
	// ToBeSigned returns a signed file document.
	//
	// NOTE: In CMS/CAdES, only one file can be signed.
	//
	// Port of getToBeSigned().
	ToBeSigned() model.DSSDocument

	// DetachedContents returns a list of detached documents.
	//
	// NOTE: In case of ASiC-S signature, we need the detached content.
	//
	// Port of getDetachedContents().
	DetachedContents() []model.DSSDocument
}
