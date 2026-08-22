// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/ReferenceOutputType.java (DSS 6.5.RC1).
package xades

// ReferenceOutputType defines possible output types of a transform/reference.
type ReferenceOutputType string

const (
	// ReferenceOutputTypeOctetStream is the octets output type.
	ReferenceOutputTypeOctetStream ReferenceOutputType = "OCTET_STREAM"

	// ReferenceOutputTypeNodeSet is the XML Node output type.
	ReferenceOutputTypeNodeSet ReferenceOutputType = "NODE_SET"
)
