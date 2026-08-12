// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignatureCryptographicVerification.java (DSS 6.5.RC1).
package signature

import "strings"

// SignatureCryptographicVerification contains a result of a signature cryptographic
// validation.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type SignatureCryptographicVerification struct {
	// errorMessages builds the error message if applicable, in the order added. Java backs
	// this with a StringBuilder; the Go port keeps the pieces and joins them lazily in
	// ErrorMessage so that "<br/>\n" is only inserted between two non-empty messages.
	errorMessages []string

	// referenceDataFound defines if (all) references data found.
	referenceDataFound bool

	// referenceDataIntact defines if (all) references data intact.
	referenceDataIntact bool

	// signatureIntact defines if the SignatureValue is valid.
	//
	// NOTE: this can be true but IsSignatureValid can be false.
	signatureIntact bool
}

// NewSignatureCryptographicVerification is the default constructor instantiating the object
// with null (zero) values.
func NewSignatureCryptographicVerification() *SignatureCryptographicVerification {
	return &SignatureCryptographicVerification{}
}

// IsReferenceDataFound gets if (all) references data found. Port of
// isReferenceDataFound().
func (s *SignatureCryptographicVerification) IsReferenceDataFound() bool {
	return s.referenceDataFound
}

// SetReferenceDataFound sets if (all) references data found. Port of
// setReferenceDataFound(boolean).
func (s *SignatureCryptographicVerification) SetReferenceDataFound(referenceDataFound bool) {
	s.referenceDataFound = referenceDataFound
}

// IsReferenceDataIntact gets if (all) references data intact. Port of
// isReferenceDataIntact().
func (s *SignatureCryptographicVerification) IsReferenceDataIntact() bool {
	return s.referenceDataIntact
}

// SetReferenceDataIntact sets if (all) references data intact. Port of
// setReferenceDataIntact(boolean).
func (s *SignatureCryptographicVerification) SetReferenceDataIntact(referenceDataIntact bool) {
	s.referenceDataIntact = referenceDataIntact
}

// IsSignatureIntact gets if the SignatureValue is valid. Port of isSignatureIntact().
func (s *SignatureCryptographicVerification) IsSignatureIntact() bool {
	return s.signatureIntact
}

// SetSignatureIntact sets if the SignatureValue is valid. Port of
// setSignatureIntact(boolean).
func (s *SignatureCryptographicVerification) SetSignatureIntact(signatureIntact bool) {
	s.signatureIntact = signatureIntact
}

// IsSignatureValid returns if the signature is valid: referenceDataFound and
// referenceDataIntact and signatureIntact are all true. Port of isSignatureValid().
func (s *SignatureCryptographicVerification) IsSignatureValid() bool {
	return s.referenceDataFound && s.signatureIntact && s.referenceDataIntact
}

// ErrorMessage returns the error messages obtained during signature cryptographic
// verification, joined as Java's StringBuilder would render them: empty string "" if the
// signature is valid. Port of getErrorMessage().
func (s *SignatureCryptographicVerification) ErrorMessage() string {
	return strings.Join(s.errorMessages, "<br/>\n")
}

// SetErrorMessage adds errorMessage to the error list. Port of setErrorMessage(String).
func (s *SignatureCryptographicVerification) SetErrorMessage(errorMessage string) {
	s.errorMessages = append(s.errorMessages, errorMessage)
}

// SetErrorMessages adds all errorMessages to the error list. Port of
// setErrorMessages(List<String>).
func (s *SignatureCryptographicVerification) SetErrorMessages(errorMessages []string) {
	for _, errorMessage := range errorMessages {
		s.SetErrorMessage(errorMessage)
	}
}

// String returns a summary of the verification result. Port of toString().
func (s *SignatureCryptographicVerification) String() string {
	return "referenceDataFound:" + boolString(s.referenceDataFound) +
		", referenceDataIntact:" + boolString(s.referenceDataIntact) +
		", signatureValid;" + boolString(s.signatureIntact) +
		" / " + s.ErrorMessage()
}

// boolString renders a bool the way Java's string concatenation would ("true"/"false").
func boolString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
