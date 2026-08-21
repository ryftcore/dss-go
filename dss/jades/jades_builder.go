// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESBuilder.java (DSS 6.5.RC1).
//
// Java's three methods return bare values and let a runtime exception propagate; every
// implementation in this package raises a checked condition somewhere below (an unsupported
// packaging, a payload that is not URL-safe, an I/O failure while reading a document to sign),
// so Build and BuildDataToBeSigned carry an error return here. getMimeType() cannot fail and
// keeps its bare shape - Java's `get` prefix is dropped per PORTING.md.
package jades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// JAdESBuilder builds a JAdES signature.
type JAdESBuilder interface {
	// Build builds a signature, adding the given SignatureValue to it, and returns the
	// DSSDocument containing the JWS binaries. Port of #build(SignatureValue).
	Build(signatureValue *model.SignatureValue) (model.DSSDocument, error)

	// BuildDataToBeSigned builds the data to be signed, incorporating a detached payload when
	// required (see 5.2.8.3 Mechanism ObjectIdByURI). Port of #buildDataToBeSigned().
	BuildDataToBeSigned() (*model.ToBeSigned, error)

	// MimeType returns the MimeType of the signature produced by the builder.
	// Port of #getMimeType().
	MimeType() enumerations.MimeType
}
