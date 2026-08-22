// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignaturePolicyValidationResult.java (DSS 6.5.RC1).
package signature

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/model"
)

// PolicyValidationResult contains results of a Policy validation.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type PolicyValidationResult struct {
	// identified indicates if the signature policy has been identified.
	identified bool

	// asn1Processable indicates if the signature policy is of an ASN.1 format.
	asn1Processable bool

	// digestAlgorithmsEqual indicates if digest algorithms match.
	digestAlgorithmsEqual bool

	// digestValid indicates is the signature policy validation result is valid.
	digestValid bool

	// digest is the digest that has been computed on the provided policy document.
	digest model.Digest

	// errors is a list of stored errors occurred during signature policy processing, keyed by
	// error key, preserving insertion order like Java's LinkedHashMap.
	errors []signaturePolicyValidationError
}

// signaturePolicyValidationError is a single entry of the error map, keeping Java's
// LinkedHashMap insertion order without pulling in an ordered-map dependency.
type signaturePolicyValidationError struct {
	key     string
	message string
}

// NewSignaturePolicyValidationResult is the default constructor instantiating the object
// with null (zero) values.
func NewSignaturePolicyValidationResult() *PolicyValidationResult {
	return &PolicyValidationResult{}
}

// IsIdentified returns if the signature policy has been obtained successfully. Port of
// isIdentified().
func (s *PolicyValidationResult) IsIdentified() bool {
	return s.identified
}

// SetIdentified sets if the signature policy has been obtained successfully. Port of
// setIdentified(boolean).
func (s *PolicyValidationResult) SetIdentified(identified bool) {
	s.identified = identified
}

// IsDigestValid returns if the signature policy has been validated successfully. Port of
// isDigestValid().
func (s *PolicyValidationResult) IsDigestValid() bool {
	return s.digestValid
}

// SetDigestValid sets if the signature policy is valid. Port of setDigestValid(boolean).
func (s *PolicyValidationResult) SetDigestValid(digestValid bool) {
	s.digestValid = digestValid
}

// IsAsn1Processable returns if the signature policy is ASN.1 processable. Port of
// isAsn1Processable().
func (s *PolicyValidationResult) IsAsn1Processable() bool {
	return s.asn1Processable
}

// SetAsn1Processable sets if the signature policy is ASN.1 processable. Port of
// setAsn1Processable(boolean).
func (s *PolicyValidationResult) SetAsn1Processable(asn1Processable bool) {
	s.asn1Processable = asn1Processable
}

// IsDigestAlgorithmsEqual returns if the DigestAlgorithm defined in the policy and used for
// the validation do match. Port of isDigestAlgorithmsEqual().
func (s *PolicyValidationResult) IsDigestAlgorithmsEqual() bool {
	return s.digestAlgorithmsEqual
}

// SetDigestAlgorithmsEqual sets if the digest algorithms match. Port of
// setDigestAlgorithmsEqual(boolean).
func (s *PolicyValidationResult) SetDigestAlgorithmsEqual(digestAlgorithmsEqual bool) {
	s.digestAlgorithmsEqual = digestAlgorithmsEqual
}

// Digest returns the Digest that has been computed on the obtained signature policy
// document.
//
// NOTE: returns the zero Digest if a validator was not able to compute the digest.
//
// Port of getDigest().
func (s *PolicyValidationResult) Digest() model.Digest {
	return s.digest
}

// SetDigest sets the Digest computed on the extracted signature policy document. Port of
// setDigest(Digest).
func (s *PolicyValidationResult) SetDigest(digest model.Digest) {
	s.digest = digest
}

// AddError allows adding a new error message occurred during the validation. A repeated
// errorKey overwrites the message the way Java's Map#put would. Port of addError(String,
// String).
func (s *PolicyValidationResult) AddError(errorKey, errorMessage string) {
	for i := range s.errors {
		if s.errors[i].key == errorKey {
			s.errors[i].message = errorMessage
			return
		}
	}
	s.errors = append(s.errors, signaturePolicyValidationError{key: errorKey, message: errorMessage})
}

// ProcessingErrors returns a user-friendly String with the obtained error messages occurred
// during the validation process. Port of getProcessingErrors().
func (s *PolicyValidationResult) ProcessingErrors() string {
	if len(s.errors) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("The errors found on signature policy validation are:")
	for _, e := range s.errors {
		b.WriteString(" at ")
		b.WriteString(e.key)
		b.WriteString(": ")
		b.WriteString(e.message)
		b.WriteString(",")
	}
	result := b.String()
	// Java's StringBuilder#setLength(length - 1) trims the trailing comma.
	return result[:len(result)-1]
}
