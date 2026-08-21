// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ASiCContainerEvidenceRecordParameters.java
// (DSS 6.5.RC1).
//
// hashCode() is dropped (see asic_parameters.go's header).
package asic

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
)

// ASiCContainerEvidenceRecordParameters defines the configuration for creation of an ASiC
// container containing an evidence record document. Java's `extends ASiCParameters` becomes
// embedding.
type ASiCContainerEvidenceRecordParameters struct {
	ASiCParameters

	// asicEvidenceRecordManifest is the ASiC evidence record manifest file to be added within
	// the container.
	asicEvidenceRecordManifest model.DSSDocument
}

// NewASiCContainerEvidenceRecordParameters is the port of the default constructor.
func NewASiCContainerEvidenceRecordParameters() *ASiCContainerEvidenceRecordParameters {
	return &ASiCContainerEvidenceRecordParameters{}
}

// AsicEvidenceRecordManifest gets the ASiCEvidenceRecordManifest file to be added within the
// container. Port of getAsicEvidenceRecordManifest().
func (p *ASiCContainerEvidenceRecordParameters) AsicEvidenceRecordManifest() model.DSSDocument {
	return p.asicEvidenceRecordManifest
}

// SetAsicEvidenceRecordManifest optionally sets a custom ASiCEvidenceRecordManifest to be added
// within the container. When defined, the current manifest file will be used for the evidence
// record incorporation. When not provided, application will create a new
// ASiCEvidenceRecordManifest based on the objects covered by the evidence record. The filename of
// the manifest file will be taken from the document name. The filename of the evidence record
// document will be taken from the manifest signature reference.
//
// Port of setAsicEvidenceRecordManifest(DSSDocument).
func (p *ASiCContainerEvidenceRecordParameters) SetAsicEvidenceRecordManifest(asicEvidenceRecordManifest model.DSSDocument) {
	p.asicEvidenceRecordManifest = asicEvidenceRecordManifest
}

// String ports toString().
func (p *ASiCContainerEvidenceRecordParameters) String() string {
	manifest := "null"
	if p.asicEvidenceRecordManifest != nil {
		manifest = fmt.Sprintf("%v", p.asicEvidenceRecordManifest)
	}
	return "ASiCContainerEvidenceRecordParameters [asicEvidenceRecordManifest=" + manifest + "] " +
		p.ASiCParameters.String()
}

// Equals ports equals(Object). Java's Objects.equals(asicEvidenceRecordManifest, ...) dispatches
// to DSSDocument#equals; model.DSSDocument declares no Equals in its interface (each
// implementation carries its own concrete Equals - see the note in
// pades/pdf_byte_range_document.go), so reference identity is compared here, which is what
// Java's equals() reduces to for the distinct-instance case this is used in.
func (p *ASiCContainerEvidenceRecordParameters) Equals(other *ASiCContainerEvidenceRecordParameters) bool {
	if p == other {
		return true
	}
	if p == nil || other == nil {
		return false
	}
	if !p.ASiCParameters.Equals(&other.ASiCParameters) {
		return false
	}
	return p.asicEvidenceRecordManifest == other.asicEvidenceRecordManifest
}
