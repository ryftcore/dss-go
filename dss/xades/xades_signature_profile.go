// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/XAdESSignatureProfile.java (DSS 6.5.RC1).
package xades

import "github.com/ryftcore/dss-go/dss/model"

// SignatureProfile is a XAdES signature creation profile.
//
// Java's two signDocument overloads declare no checked exception; both landed implementers
// (LevelBaselineB.SignDocument/SignDocuments, xades_level_baseline_b.go) return a Go error
// alongside the document, so the interface follows suit per PORTING.md's unchecked-exception
// convention.
type SignatureProfile interface {
	// SignDocument creates a signature of the defined profile for signing a document. Ports
	// signDocument(DSSDocument, SignatureParameters, byte[]).
	SignDocument(toSignDocument model.DSSDocument, parameters *SignatureParameters,
		signatureValue []byte) (model.DSSDocument, error)

	// SignDocuments creates a signature of the defined profile for signing a list of documents.
	// Ports signDocument(List<DSSDocument>, XAdESSignatureParameters, byte[]).
	SignDocuments(toSignDocuments []model.DSSDocument, parameters *SignatureParameters,
		signatureValue []byte) (model.DSSDocument, error)
}
