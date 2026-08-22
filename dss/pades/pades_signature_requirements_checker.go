// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESSignatureRequirementsChecker.java (DSS 6.5.RC1).
//
// Java extends document.SignatureRequirementsChecker and overrides the protected
// checkTLevelIsHighest. Go cannot override across embedding, and the base calls its own
// checkTLevelIsHighest from an unexported assertTLevelIsHighest, so the one public entry point
// that reaches it - AssertExtendToTLevelPossible - is re-declared here together with the two
// private steps it performs. Everything else (assertSignaturesValid, assertSigningCertificateIsValid,
// assertExtendToLTLevelPossible, assertCertificateChainValidForLTLevel, ...) is the embedded
// base's, unchanged.
//
// The type parameter is *cades.CAdESTimestampParameters, not *PAdESTimestampParameters: Java's
// SignatureParameters extends SignatureParameters extends
// AbstractSignatureParameters<TimestampParameters>, so that is the timestamp-parameters type
// the checked AbstractSignatureParameters carries.
//
// slf4j is dropped (PORTING.md).
package pades

import (
	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// SignatureRequirementsChecker verifies signature creation or augmentation requirements for
// PAdES signatures.
type SignatureRequirementsChecker struct {
	*document.SignatureRequirementsChecker[*cades.TimestampParameters]

	// certificateVerifier is the base's field of the same name, kept here because the base does
	// not expose it and the re-declared assertions need the augmentation alert.
	certificateVerifier validation.CertificateVerifier
}

// NewPAdESSignatureRequirementsChecker is the default constructor.
// Port of PAdESSignatureRequirementsChecker(CertificateVerifier, PAdESSignatureParameters).
func NewPAdESSignatureRequirementsChecker(certificateVerifier validation.CertificateVerifier,
	signatureParameters *SignatureParameters) *SignatureRequirementsChecker {
	return &SignatureRequirementsChecker{
		SignatureRequirementsChecker: document.NewSignatureRequirementsChecker[*cades.TimestampParameters](
			certificateVerifier, &signatureParameters.AbstractSignatureParameters),
		certificateVerifier: certificateVerifier,
	}
}

// AssertExtendToTLevelPossible verifies whether extension of the signatures to T-level is
// possible. Port of the inherited #assertExtendToTLevelPossible; re-declared so that
// CheckTLevelIsHighest below is the one that runs.
func (c *SignatureRequirementsChecker) AssertExtendToTLevelPossible(signatures []validation.AdvancedSignature) {
	c.assertTLevelIsHighest(signatures)
	c.assertHasNoEmbeddedEvidenceRecords(signatures)
}

// assertTLevelIsHighest checks whether across signatures the T-level is highest and T-level
// augmentation can be performed. Port of the inherited protected #assertTLevelIsHighest.
func (c *SignatureRequirementsChecker) assertTLevelIsHighest(signatures []validation.AdvancedSignature) {
	if c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		c.CheckTLevelIsHighest(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Error on signature augmentation to T-level.")
		padesSignatureRequirementsCheckerAlert(c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel(), status)
	}
}

// assertHasNoEmbeddedEvidenceRecords checks whether none of the signatures is preserved by an
// embedded evidence record. Port of the inherited protected #assertHasNoEmbeddedEvidenceRecords.
func (c *SignatureRequirementsChecker) assertHasNoEmbeddedEvidenceRecords(signatures []validation.AdvancedSignature) {
	if c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		if c.HasEmbeddedEvidenceRecords(signature) {
			status.AddRelatedTokenAndErrorMessage(signature,
				"The signature is preserved by an embedded evidence record.")
		}
	}
	if !status.IsEmpty() {
		status.SetMessage("Error on signature augmentation")
		padesSignatureRequirementsCheckerAlert(c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel(), status)
	}
}

// CheckTLevelIsHighest verifies whether the signature has a maximum B- or T-level. Unlike the
// CAdES/XAdES rule, a PAdES signature carrying a DSS dictionary but no associated timestamp may
// still be extended, so that a best-signature-time can be provided and fresh revocation data
// incorporated. Port of the protected #checkTLevelIsHighest override.
func (c *SignatureRequirementsChecker) CheckTLevelIsHighest(signature validation.AdvancedSignature,
	status *validation.SignatureStatus) {
	if signature.HasLTAProfile() {
		status.AddRelatedTokenAndErrorMessage(signature, "The signature is already extended with a higher level.")

	} else if signature.HasLTProfile() && !signature.AreAllSelfSignedCertificates() {
		if signature.HasTProfile() {
			status.AddRelatedTokenAndErrorMessage(signature, "The signature is already extended with a higher level.")
		}
		// NOTE: Otherwise allow extension, as it may be required to provide a best-signature-time
		// to ensure the best practice of fresh revocation data incorporation.
		// Upstream logs "Signature contains a DSS dictionary, but no associated timestamp.
		// Extension may lead to LTA-level."
	}
}

// padesSignatureRequirementsCheckerAlert is document.mustAlert, which is unexported there: it
// raises the alert and turns a refusing handler's error into Java's propagating exception.
func padesSignatureRequirementsCheckerAlert(statusAlert alert.StatusAlert, status alert.Status) {
	if err := statusAlert.Alert(status); err != nil {
		panic(err)
	}
}
