// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/SignatureIntegrityValidator.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi (see PORTING_PLAN.md), so the
// type keeps its Java name unqualified.
//
// Java's abstract protected verify(PublicKey) is implemented by concrete validators defined
// outside this chunk (e.g. CAdES/XAdES signature validators, ported in later phases). This
// follows the same override-registration pattern as model.TokenBase.InitToken: a concrete
// validator embeds SignatureIntegrityValidator and registers itself via
// InitSignatureIntegrityValidator before Validate is called.
//
// ASSUMPTION (flagged for integrator reconciliation, see chunk X509-B which owns
// CandidatesForSigningCertificate and CertificateValidity): they are assumed to expose
//
//	func (c *CandidatesForSigningCertificate) IsEmpty() bool
//	func (c *CandidatesForSigningCertificate) TheBestCandidate() *CertificateValidity
//	func (c *CandidatesForSigningCertificate) CertificateValidityList() []*CertificateValidity
//	func (v *CertificateValidity) PublicKey() *model.PublicKey
//	func (v *CertificateValidity) IsValid() bool
//	func (v *CertificateValidity) CertificateToken() *model.CertificateToken
package spi

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
)

// SignatureIntegrityValidatorOverrides declares the operation Java's abstract class leaves
// abstract, and that the base implementation calls back into.
type SignatureIntegrityValidatorOverrides interface {
	// Verify verifies if the signature has been created with the given public key.
	// Port of the protected abstract verify(PublicKey), which throws DSSException.
	Verify(publicKey *model.PublicKey) (bool, error)
}

// SignatureIntegrityValidator is the contract of a validator that checks signature integrity
// among a provided list of signing certificate candidates. It is the polymorphic half of the
// Java abstract class; the state and concrete method bodies live in
// SignatureIntegrityValidatorBase (mirrors the model.Token / model.TokenBase split).
type SignatureIntegrityValidator interface {
	// Validate verifies validity of a signature across a provided signing certificate
	// candidates list. Port of validate(CandidatesForSigningCertificate).
	Validate(candidates *CandidatesForSigningCertificate) *CertificateValidity
	// ErrorMessages returns error messages after processing of the Validate method, if
	// present. Port of getErrorMessages().
	ErrorMessages() []string
}

// SignatureIntegrityValidatorBase checks signature integrity among a provided list of signing
// certificate candidates.
type SignatureIntegrityValidatorBase struct {
	// overrides points back at the concrete validator; see InitSignatureIntegrityValidator.
	overrides SignatureIntegrityValidatorOverrides

	// errorMessages holds the errors occurred during the signature integrity validation; nil
	// until Validate has run once.
	errorMessages []string
}

// NewSignatureIntegrityValidatorBase instantiates the base state of a signature integrity
// validator.
func NewSignatureIntegrityValidatorBase() SignatureIntegrityValidatorBase {
	return SignatureIntegrityValidatorBase{}
}

// InitSignatureIntegrityValidator registers the concrete validator with its base so that the
// base can dispatch to Verify. It must be called exactly once, by the concrete validator's
// constructor, before Validate.
func (v *SignatureIntegrityValidatorBase) InitSignatureIntegrityValidator(overrides SignatureIntegrityValidatorOverrides) {
	v.overrides = overrides
}

// signatureIntegrityValidatorOverrides returns the registered overrides, panicking when the
// concrete validator forgot to call InitSignatureIntegrityValidator.
func (v *SignatureIntegrityValidatorBase) signatureIntegrityValidatorOverrides() SignatureIntegrityValidatorOverrides {
	if v.overrides == nil {
		panic("SignatureIntegrityValidator was not initialised: the concrete validator must call InitSignatureIntegrityValidator in its constructor")
	}
	return v.overrides
}

// Validate verifies validity of a signature across a provided signing certificate candidates
// list. Port of validate(CandidatesForSigningCertificate).
//
// NOTE: in case of a failed validation, use ErrorMessages() after processing this method for
// more details.
func (v *SignatureIntegrityValidatorBase) Validate(candidates *CandidatesForSigningCertificate) *CertificateValidity {
	v.errorMessages = []string{}

	if candidates.IsEmpty() {
		v.errorMessages = append(v.errorMessages, "There is no signing certificate within the signature or certificate pool.")
	}

	// 1) Process the best found candidate
	bestCandidate := candidates.TheBestCandidate()
	if bestCandidate != nil {
		intact, err := v.isSignatureIntact(bestCandidate)
		if err != nil {
			v.errorMessages = append(v.errorMessages, "Best candidate validation failed : "+err.Error())
		} else if intact {
			return bestCandidate // best candidate either valid or better is not available
		} else {
			v.errorMessages = append(v.errorMessages, "Signature verification failed against the best candidate.")
		}
	}

	// 2) Validate among other candidates
	var bestCertificateValidity *CertificateValidity

	certificateNumber := 0
	for _, certificateValidity := range candidates.CertificateValidityList() {
		if certificateValidity == bestCandidate {
			continue // do not process validation twice
		}
		errorMessagePrefix := fmt.Sprintf("Certificate #%d: ", certificateNumber+1)
		intact, err := v.isSignatureIntact(certificateValidity)
		if err != nil {
			v.errorMessages = append(v.errorMessages, errorMessagePrefix+err.Error())
		} else if intact {
			bestCertificateValidity = certificateValidity
			if certificateValidity.IsValid() {
				break
			}
			// else: certificate candidate does not match a signing certificate reference; the
			// Java WARN log is dropped (slf4j is not load-bearing here).
		} else {
			// upon returning false, santuarioSignature (class XMLSignature) will log
			// "Signature verification failed." with WARN level.
			v.errorMessages = append(v.errorMessages, errorMessagePrefix+"Signature verification failed")
		}
		certificateNumber++
	}

	return bestCertificateValidity
}

// isSignatureIntact ports the private isSignatureIntact(CertificateValidity).
func (v *SignatureIntegrityValidatorBase) isSignatureIntact(certificateValidity *CertificateValidity) (bool, error) {
	publicKey := certificateValidity.PublicKey()
	return v.signatureIntegrityValidatorOverrides().Verify(publicKey)
}

// ErrorMessages returns error messages after processing of the Validate(candidates) method, if
// present. Port of getErrorMessages().
//
// Panics if Validate has not run yet (Java's IllegalStateException).
func (v *SignatureIntegrityValidatorBase) ErrorMessages() []string {
	if v.errorMessages == nil {
		panic("The Validate(candidates) method shall be proceeded before!")
	}
	return v.errorMessages
}
