// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/extension/DocumentExtender.java (DSS 6.5.RC1).
package extension

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

// DocumentExtender performs a signature augmentation to the specified
// target augmentation level, respectively to the signature's format, which
// is determined in the runtime during the method execution.
//
// This can be used as an alternative to
// DocumentSignatureService#ExtendDocument, which requires the client to
// know the original signature format before the method execution.
type DocumentExtender interface {
	// SetCertificateVerifier sets a CertificateVerifier providing a
	// configuration for validation of certificates within the original
	// signature or signatures. Ports #setCertificateVerifier.
	SetCertificateVerifier(certificateVerifier validation.CertificateVerifier)

	// SetTspSource sets the source to be used for time-stamp requests,
	// when applicable (e.g. BASELINE-T or BASELINE-LTA profile
	// extension). Ports #setTspSource.
	SetTspSource(tspSource validation.TSPSource)

	// IsSupported checks if the document is supported by the current
	// signature extender. Ports #isSupported.
	IsSupported(dssDocument model.DSSDocument) bool

	// ExtendDocument performs augmentation of all signatures within the
	// provided document according to the specified target
	// signatureProfile. Ports #extendDocument(SignatureProfile).
	ExtendDocument(signatureProfile enumerations.SignatureProfile) model.DSSDocument

	// ExtendDocumentDetached performs augmentation of all signatures of a
	// detached format within the provided document according to the
	// specified target signatureProfile. Ports
	// #extendDocument(SignatureProfile, List).
	ExtendDocumentDetached(signatureProfile enumerations.SignatureProfile, detachedContents []model.DSSDocument) model.DSSDocument

	// ExtendDocumentWithParameters performs augmentation of all
	// signatures within the provided document according to the specified
	// augmentationParameters, with format specific extensionParameters.
	//
	// If explicitParameters are provided, and they contain an
	// implementation of parameters matching the current signature
	// format, those signature parameters will be used explicitly on this
	// signature augmentation. This can be useful when a signature
	// augmentation process for a certain signature format should take
	// into account particular signature constraints (i.e. /Contents size
	// for a PAdES document time-stamp, etc.).
	//
	// NOTE 1: If the matching implementation of extension parameters
	// found, and it contains signatureLevel and/or detachedContents
	// definition, the existing values will be used. Otherwise, the
	// content of the extension parameters will be overwritten
	// respectively.
	// NOTE 2: The definition of extension parameters is order dependent.
	// The first parameters matching the signature format implementation
	// will be used.
	//
	// Ports #extendDocument(SignatureProfile,
	// SerializableSignatureParameters...).
	ExtendDocumentWithParameters(signatureProfile enumerations.SignatureProfile, extensionParameters ...model.SerializableSignatureParameters) model.DSSDocument

	// ExtendDocumentDetachedWithParameters performs augmentation of all
	// signatures of a detached format within the provided document
	// according to the specified target signatureProfile, with format
	// specific extensionParameters.
	//
	// See ExtendDocumentWithParameters for the semantics of
	// extensionParameters.
	//
	// Ports #extendDocument(SignatureProfile, List,
	// SerializableSignatureParameters...).
	ExtendDocumentDetachedWithParameters(signatureProfile enumerations.SignatureProfile, detachedContents []model.DSSDocument, extensionParameters ...model.SerializableSignatureParameters) model.DSSDocument
}
