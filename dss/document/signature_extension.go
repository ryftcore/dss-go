// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/SignatureExtension.java (DSS 6.5.RC1).
package document

import "github.com/utain/esig/dss/model"

// SignatureExtension extends the level of AdES signature of a document. After level -B, going
// upper in the signature format level consists of adding unsigned properties to the signature.
// It can be done without breaking the signature. Generic over the SP implementation of signature
// parameters corresponding to the supported signature format.
type SignatureExtension[SP any] interface {
	// ExtendSignatures extends the level of the signatures contained in a document. Port of
	// #extendSignatures.
	ExtendSignatures(document model.DSSDocument, params SP) model.DSSDocument
}
