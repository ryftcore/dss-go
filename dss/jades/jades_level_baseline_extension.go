// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESLevelBaselineExtension.java (DSS 6.5.RC1).
//
// Java's interface extends SignatureExtension<JAdESSignatureParameters>, whose single method is
// extendSignatures(DSSDocument, JAdESSignatureParameters). Every JAdES level also overrides the
// protected extendSignatures(List<AdvancedSignature>, JAdESSignatureParameters), and Go has no
// overloading: as in the XAdES port, the document-taking entry point carries the Document suffix
// and the virtual, list-taking one keeps the plain ExtendSignatures name (see
// jades_level_baseline_t.go). This interface therefore declares ExtendSignaturesDocument rather
// than embedding document.SignatureExtension, whose method name the level types spend on the
// virtual overload; the extra error return is what every Java throw in the extension chain
// becomes.
package jades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// JAdESLevelBaselineExtension is a JAdES extension.
type JAdESLevelBaselineExtension interface {
	// ExtendSignaturesDocument extends the level of the signatures contained in a document.
	// Port of SignatureExtension#extendSignatures(DSSDocument, JAdESSignatureParameters).
	ExtendSignaturesDocument(document model.DSSDocument,
		params *JAdESSignatureParameters) (model.DSSDocument, error)

	// SetOperationKind sets the signing operation.
	//
	// NOTE: the internal variable, used in the signature creation/extension process.
	//
	// Port of #setOperationKind.
	SetOperationKind(signingOperation enumerations.SigningOperation)
}
