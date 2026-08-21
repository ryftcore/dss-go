// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/evidencerecord/SerializableASiCContainerEvidenceRecordParameters.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/enumerations"

// SerializableASiCContainerEvidenceRecordParameters defines parameters for an ASiC container
// generation with an evidence record document. Ports the Java interface (Serializable has no
// Go analogue and is dropped, per PORTING.md).
type SerializableASiCContainerEvidenceRecordParameters interface {
	// ContainerType gets the target container type. Ports getContainerType().
	ContainerType() enumerations.ASiCContainerType
}
