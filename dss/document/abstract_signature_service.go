// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/AbstractSignatureService.java (DSS 6.5.RC1).
package document

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
)

func init() {
	spi.DSSSecurityProviderInitSystemProviders()
}

// AbstractSignatureService contains the main methods for a signature creation/extension, generic
// over the SP implementation of SerializableSignatureParameters and TP implementation of
// SerializableTimestampParameters.
//
// Java's AbstractSignatureService<SP, TP> implements DocumentSignatureService<SP, TP>; this Go
// port does not assert that conformance, since - like the Java class it ports - it is meant to
// be embedded by a concrete, format-specific service (CAdES, XAdES, ... ported in later phases)
// that supplies the remaining DocumentSignatureService methods (GetDataToSign, SignDocument,
// ExtendDocument, GetContentTimestamp) this type does not implement, matching
// AbstractDocumentExtender's equivalent non-assertion in this same package.
type AbstractSignatureService[SP model.SerializableSignatureParameters, TP model.SerializableTimestampParameters] struct {
	// CertificateVerifier is used for a certificate chain validation.
	CertificateVerifier validation.CertificateVerifier

	// TspSource is the TSPSource to use for timestamp requests.
	TspSource validation.TSPSource
}

// NewAbstractSignatureService constructs a signature service; the CertificateVerifier must be
// set and cannot be nil. Port of the protected constructor.
func NewAbstractSignatureService[SP model.SerializableSignatureParameters, TP model.SerializableTimestampParameters](certificateVerifier validation.CertificateVerifier) AbstractSignatureService[SP, TP] {
	if certificateVerifier == nil {
		panic("CertificateVerifier cannot be null !")
	}
	return AbstractSignatureService[SP, TP]{CertificateVerifier: certificateVerifier}
}

// SetTspSource ports #setTspSource.
func (s *AbstractSignatureService[SP, TP]) SetTspSource(tspSource validation.TSPSource) {
	s.TspSource = tspSource
}

// AssertSigningCertificateValid raises a panic if the signing rules forbid the use of the
// certificate. Port of the protected #assertSigningCertificateValid; Java's thrown
// IllegalArgumentException, being unchecked, becomes a panic (this method has no error channel
// to return through, matching Java's void signature).
func (s *AbstractSignatureService[SP, TP]) AssertSigningCertificateValid(parameters *AbstractSignatureParameters[TP]) {
	signingCertificate := parameters.SigningCertificate()
	if signingCertificate == nil {
		if parameters.GenerateTBSWithoutCertificate() {
			return
		}
		panic("Signing Certificate is not defined! Set signing certificate or use method SetGenerateTBSWithoutCertificate(true).")
	}

	signatureRequirementsChecker := NewSignatureRequirementsChecker[TP](s.CertificateVerifier, parameters)
	signatureRequirementsChecker.AssertSigningCertificateIsValid(signingCertificate)
}

// EnsureSignatureValue ensures the provided signatureValue has the expected
// targetSignatureAlgorithm. Port of the protected #ensureSignatureValue.
func (s *AbstractSignatureService[SP, TP]) EnsureSignatureValue(targetSignatureAlgorithm enumerations.SignatureAlgorithm, signatureValue *model.SignatureValue) (*model.SignatureValue, error) {
	return NewSignatureValueChecker().EnsureSignatureValue(signatureValue, targetSignatureAlgorithm)
}

// GetFinalDocumentName generates and returns a final name for the document to create. Port of
// the protected #getFinalDocumentName(DSSDocument, SigningOperation, SignatureLevel,
// SignaturePackaging, MimeType) - the full 5-argument overload every other GetFinalFileName/
// GetFinalDocumentName variant below delegates to.
func (s *AbstractSignatureService[SP, TP]) GetFinalDocumentName(originalFile model.DSSDocument, operation enumerations.SigningOperation, level enumerations.SignatureLevel, packaging enumerations.SignaturePackaging, containerMimeType enumerations.MimeType) (string, error) {
	return validation.NewFileNameBuilder().
		SetOriginalFilename(originalFile.Name()).
		SetSigningOperation(operation).
		SetSignatureLevel(level).
		SetSignaturePackaging(packaging).
		SetMimeType(containerMimeType).
		Build()
}

// GetFinalFileName returns the final name for the document to create. Port of the
// (DSSDocument, SigningOperation) #getFinalFileName overload.
func (s *AbstractSignatureService[SP, TP]) GetFinalFileName(originalFile model.DSSDocument, operation enumerations.SigningOperation) (string, error) {
	return s.GetFinalFileNameWithLevel(originalFile, operation, "")
}

// GetFinalFileNameWithLevel returns the final name for the document to create. Port of the
// (DSSDocument, SigningOperation, SignatureLevel) #getFinalFileName overload.
func (s *AbstractSignatureService[SP, TP]) GetFinalFileNameWithLevel(originalFile model.DSSDocument, operation enumerations.SigningOperation, level enumerations.SignatureLevel) (string, error) {
	return s.GetFinalFileNameWithPackaging(originalFile, operation, level, "")
}

// GetFinalFileNameWithPackaging returns the final name for the document to create. Port of the
// (DSSDocument, SigningOperation, SignatureLevel, SignaturePackaging) #getFinalFileName overload.
func (s *AbstractSignatureService[SP, TP]) GetFinalFileNameWithPackaging(originalFile model.DSSDocument, operation enumerations.SigningOperation, level enumerations.SignatureLevel, packaging enumerations.SignaturePackaging) (string, error) {
	return s.GetFinalDocumentName(originalFile, operation, level, packaging, nil)
}

// GetFinalDocumentNameWithMimeType generates and returns a final name for the document to
// create. Port of the (DSSDocument, SigningOperation, SignatureLevel, MimeType)
// #getFinalDocumentName overload.
func (s *AbstractSignatureService[SP, TP]) GetFinalDocumentNameWithMimeType(originalFile model.DSSDocument, operation enumerations.SigningOperation, level enumerations.SignatureLevel, containerMimeType enumerations.MimeType) (string, error) {
	return s.GetFinalDocumentName(originalFile, operation, level, "", containerMimeType)
}

// Timestamp ports the default #timestamp(DSSDocument, TP) implementation: unsupported unless a
// concrete format overrides it (later phases). Panics with Java's unchecked
// UnsupportedOperationException.
func (s *AbstractSignatureService[SP, TP]) Timestamp(toTimestampDocument model.DSSDocument, parameters TP) model.DSSDocument {
	panic("Unsupported operation for this file format")
}

// IsValidSignatureValue verifies the signature value against a ToBeSigned and a CertificateToken.
// Port of #isValidSignatureValue. Panics if any argument is nil (Java
// Objects.requireNonNull(...)).
func (s *AbstractSignatureService[SP, TP]) IsValidSignatureValue(toBeSigned *model.ToBeSigned, signatureValue *model.SignatureValue, signingCertificate *model.CertificateToken) bool {
	if toBeSigned == nil {
		panic("ToBeSigned cannot be null!")
	}
	if signatureValue == nil {
		panic("SignatureValue cannot be null!")
	}
	if signingCertificate == nil {
		panic("CertificateToken cannot be null!")
	}

	signature, err := DSSSignatureSecurityFactoryInstance.Build(signatureValue.Algorithm())
	if err != nil {
		return false
	}
	signature.InitVerify(signingCertificate.PublicKey())
	signature.Update(toBeSigned.Bytes())
	valid, err := signature.Verify(signatureValue.Value())
	if err != nil {
		return false
	}
	return valid
}
