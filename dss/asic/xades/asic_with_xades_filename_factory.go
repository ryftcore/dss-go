// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/ASiCWithXAdESFilenameFactory.java
// (DSS 6.5.RC1).
package xades

import "github.com/ryftcore/dss-go/dss/asic"

// ASiCWithXAdESFilenameFactory is used to provide filenames for newly created ZIP-entries
// during a signature creation or extension for ASiC with XAdES containers.
//
// NOTE: Names of signature or manifest files shall be defined with leading "META-INF/" string,
// specifying the target folder of the signature file within a container.
//
// As the same factory is used for ASiC-S and ASiC-E container types, it shall implement logic
// for both container types, when applicable. The type of the container can be obtained from
// asicContent.ContainerType().
type ASiCWithXAdESFilenameFactory interface {
	asic.FilenameFactory
	asic.EvidenceRecordFilenameFactory
}
