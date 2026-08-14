// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/ReferenceOutputType.java (DSS 6.5.RC1).
package xades

// ReferenceOutputType defines possible output types of a transform/reference.
type ReferenceOutputType string

const (
	// ReferenceOutputType_OCTET_STREAM is the octets output type.
	ReferenceOutputType_OCTET_STREAM ReferenceOutputType = "OCTET_STREAM"

	// ReferenceOutputType_NODE_SET is the XML Node output type.
	ReferenceOutputType_NODE_SET ReferenceOutputType = "NODE_SET"
)
