// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/SignatureBuilder.java (DSS 6.5.RC1).
package xades

import "github.com/ryftcore/dss-go/dss/model"

// SignatureBuilder builds a XAdES signature of the defined format.
//
// Java's signDocument(byte[]) declares no checked exception, but its sole landed implementer
// (XAdESSignatureBuilder.SignDocument, xades_signature_builder.go) returns a Go error alongside
// the document - the unchecked-exception-to-error convention PORTING.md applies whenever a
// method can fail on malformed input - so the interface carries the same (T, error) shape.
type SignatureBuilder interface {
	// SignDocument signs a document. Ports signDocument(byte[]).
	SignDocument(signatureValue []byte) (model.DSSDocument, error)
}
