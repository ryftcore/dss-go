// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/SignatureRequirementsChecker.java (DSS 6.5.RC1).
package document

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// SignatureRequirementsChecker is used to verify if the signature can be created according to
// the provided requirements in a CertificateVerifier instance, generic over the TP implementation
// of timestamp parameters carried by the checked AbstractSignatureParameters.
type SignatureRequirementsChecker[TP model.SerializableTimestampParameters] struct {
	// certificateVerifier is used for certificates validation.
	certificateVerifier validation.CertificateVerifier

	// signatureParameters is used for signature creation/extension.
	signatureParameters *AbstractSignatureParameters[TP]
}

// NewSignatureRequirementsChecker is the default constructor.
func NewSignatureRequirementsChecker[TP model.SerializableTimestampParameters](certificateVerifier validation.CertificateVerifier, signatureParameters *AbstractSignatureParameters[TP]) *SignatureRequirementsChecker[TP] {
	if certificateVerifier == nil {
		panic("CertificateVerifier cannot be null!")
	}
	if signatureParameters == nil {
		panic("Signature parameters cannot be null!")
	}
	return &SignatureRequirementsChecker[TP]{certificateVerifier: certificateVerifier, signatureParameters: signatureParameters}
}

// AssertSigningCertificateIsValid verifies whether the provided certificate token is acceptable
// for a signature creation against the provided signatureParameters. Port of the CertificateToken
// overload of #assertSigningCertificateIsValid.
func (c *SignatureRequirementsChecker[TP]) AssertSigningCertificateIsValid(certificateToken *model.CertificateToken) {
	validationAlerter := c.initValidationAlerterForCertificate(certificateToken)
	c.assertCertificatesAreYetValid(validationAlerter, certificateToken)
	c.assertCertificatesAreNotExpired(validationAlerter, certificateToken)
	c.assertCertificatesAreNotRevoked(validationAlerter, certificateToken)
}

// AssertSigningCertificatesAreValid verifies a signing certificate for a collection of the given
// signatures. Port of the Collection<AdvancedSignature> overload of
// #assertSigningCertificateIsValid.
func (c *SignatureRequirementsChecker[TP]) AssertSigningCertificatesAreValid(signatures []validation.AdvancedSignature) {
	var signaturesToValidate []validation.AdvancedSignature
	for _, s := range signatures {
		if !c.isSignatureGeneratedWithoutCertificate(s) {
			signaturesToValidate = append(signaturesToValidate, s)
		}
	}
	if len(signaturesToValidate) == 0 {
		return
	}

	validationAlerter := c.initValidationAlerterForSignatures(signatures)
	c.assertCertificatesAreYetValidForSignatures(validationAlerter)
	c.assertCertificatesAreNotExpiredForSignatures(validationAlerter)
	c.assertCertificatesAreNotRevokedForSignatures(validationAlerter)
}

// isSignatureGeneratedWithoutCertificate ports the private #isSignatureGeneratedWithoutCertificate.
func (c *SignatureRequirementsChecker[TP]) isSignatureGeneratedWithoutCertificate(signature validation.AdvancedSignature) bool {
	if c.isSigningCertificateIdentified(signature) {
		return false // signing-certificate is identified
	} else if c.signatureParameters.GenerateTBSWithoutCertificate() {
		return true
	}
	panic(exception.NewIllegalInputException("Signing-certificate token was not found! Unable to verify its validity. " +
		"Provide signing-certificate or use method #setGenerateTBSWithoutCertificate(true) for signature creation without signing-certificate."))
}

// isSigningCertificateIdentified ports the private #isSigningCertificateIdentified.
func (c *SignatureRequirementsChecker[TP]) isSigningCertificateIdentified(signature validation.AdvancedSignature) bool {
	return signature.CertificateSource().NumberOfCertificates() != 0 && signature.SigningCertificateToken() != nil
}

// assertCertificatesAreYetValidForSignatures ports the no-certificateToken private
// #assertCertificatesAreYetValid(ValidationAlerter) overload.
func (c *SignatureRequirementsChecker[TP]) assertCertificatesAreYetValidForSignatures(validationAlerter validation.Alerter) {
	c.assertCertificatesAreYetValid(validationAlerter, nil)
}

// assertCertificatesAreYetValid ports the private #assertCertificatesAreYetValid(ValidationAlerter,
// CertificateToken) overload; certificateToken nil means "validate signatures" instead.
func (c *SignatureRequirementsChecker[TP]) assertCertificatesAreYetValid(validationAlerter validation.Alerter, certificateToken *model.CertificateToken) {
	if c.certificateVerifier.AlertOnNotYetValidCertificate() == nil {
		return
	}
	if certificateToken != nil {
		validationAlerter.AssertCertificateIsYetValid(certificateToken)
	} else {
		validationAlerter.AssertAllSignaturesAreYetValid()
	}
}

// assertCertificatesAreNotExpiredForSignatures ports the no-certificateToken private
// #assertCertificatesAreNotExpired(ValidationAlerter) overload.
func (c *SignatureRequirementsChecker[TP]) assertCertificatesAreNotExpiredForSignatures(validationAlerter validation.Alerter) {
	c.assertCertificatesAreNotExpired(validationAlerter, nil)
}

// assertCertificatesAreNotExpired ports the private #assertCertificatesAreNotExpired(ValidationAlerter,
// CertificateToken) overload.
func (c *SignatureRequirementsChecker[TP]) assertCertificatesAreNotExpired(validationAlerter validation.Alerter, certificateToken *model.CertificateToken) {
	if c.certificateVerifier.AlertOnExpiredCertificate() == nil {
		return
	}
	if certificateToken != nil {
		validationAlerter.AssertCertificateNotExpired(certificateToken)
	} else {
		validationAlerter.AssertAllSignaturesNotExpired()
	}
}

// assertCertificatesAreNotRevokedForSignatures ports the no-certificateToken private
// #assertCertificatesAreNotRevoked(ValidationAlerter) overload.
func (c *SignatureRequirementsChecker[TP]) assertCertificatesAreNotRevokedForSignatures(validationAlerter validation.Alerter) {
	c.assertCertificatesAreNotRevoked(validationAlerter, nil)
}

// assertCertificatesAreNotRevoked ports the private #assertCertificatesAreNotRevoked(ValidationAlerter,
// CertificateToken) overload.
func (c *SignatureRequirementsChecker[TP]) assertCertificatesAreNotRevoked(validationAlerter validation.Alerter, certificateToken *model.CertificateToken) {
	if !c.signatureParameters.CheckCertificateRevocation() {
		return
	}
	if c.certificateVerifier.AlertOnMissingRevocationData() == nil && c.certificateVerifier.AlertOnRevokedCertificate() == nil {
		return
	}

	validationAlerter.AssertAllRequiredRevocationDataPresent()
	if certificateToken != nil {
		validationAlerter.AssertCertificateNotRevoked(certificateToken)
	} else {
		validationAlerter.AssertAllSignatureCertificatesNotRevoked()
	}
}

// initValidationAlerterForCertificate initializes the validation alerter for certificate
// validation. Port of the CertificateToken overload of #initValidationAlerter.
func (c *SignatureRequirementsChecker[TP]) initValidationAlerterForCertificate(certificateToken *model.CertificateToken) validation.Alerter {
	var signingDate time.Time
	if sd := c.signatureParameters.BLevel().SigningDate(); sd != nil {
		signingDate = *sd
	}
	validationContext := validation.NewSignatureValidationContextAtTime(signingDate)
	validationContext.Initialize(c.getCertificateVerifier())

	certificateChain := c.signatureParameters.CertificateChain()
	if len(certificateChain) == 0 {
		if c.signatureParameters.CheckCertificateRevocation() {
			panic("Certificate chain shall be provided for a revocation check! " +
				"Please use parameters.SetCertificateChain(...) method to provide a certificate chain.")
		}
		certificateChain = nil
	}
	validationContext.AddCertificateTokenForVerification(certificateToken)
	for _, certificate := range certificateChain {
		validationContext.AddCertificateTokenForVerification(certificate)
	}

	validationContext.Validate()

	signatureValidationAlerter := validation.NewSignatureValidationAlerter(validationContext)
	signatureValidationAlerter.SetSigningOperation(enumerations.SigningOperationSign)
	return signatureValidationAlerter
}

// initValidationAlerterForSignatures initializes the validation alerter for signature
// validation. Port of the Collection<AdvancedSignature> overload of #initValidationAlerter.
func (c *SignatureRequirementsChecker[TP]) initValidationAlerterForSignatures(signatures []validation.AdvancedSignature) validation.Alerter {
	var signingDate time.Time
	if sd := c.signatureParameters.BLevel().SigningDate(); sd != nil {
		signingDate = *sd
	}
	validationContext := validation.NewSignatureValidationContextAtTime(signingDate)
	validationContext.Initialize(c.getCertificateVerifier())

	for _, signature := range signatures {
		validationContext.AddSignatureForVerification(signature)
	}

	validationContext.Validate()

	signatureValidationAlerter := validation.NewSignatureValidationAlerter(validationContext)
	signatureValidationAlerter.SetSigningOperation(enumerations.SigningOperationExtend)
	return signatureValidationAlerter
}

// getCertificateVerifier gets CertificateVerifier to be used for validation context
// verification. Port of the protected #getCertificateVerifier.
func (c *SignatureRequirementsChecker[TP]) getCertificateVerifier() validation.CertificateVerifier {
	if c.signatureParameters.CheckCertificateRevocation() {
		return c.certificateVerifier
	}

	// skip revocation check
	offlineCertificateVerifier := validation.NewCertificateVerifierBuilder(c.certificateVerifier).BuildOfflineCopy()

	acceptAllRevocationDataVerifier := c.createAcceptAllRevocationDataVerifier()
	offlineCertificateVerifier.SetRevocationDataVerifier(acceptAllRevocationDataVerifier)
	timestampTokenVerifier := offlineCertificateVerifier.TimestampTokenVerifier()
	if timestampTokenVerifier == nil {
		timestampTokenVerifier = validation.NewDefaultTimestampTokenVerifier()
	}
	timestampTokenVerifier.SetRevocationDataVerifier(acceptAllRevocationDataVerifier)

	return offlineCertificateVerifier
}

// createAcceptAllRevocationDataVerifier creates a RevocationDataVerifier returning always a
// valid revocation status for a certificate.
//
// NOTE: This method is used internally for a silent revocation data processing check. Port of
// the private #createAcceptAllRevocationDataVerifier.
func (c *SignatureRequirementsChecker[TP]) createAcceptAllRevocationDataVerifier() *validation.RevocationDataVerifier {
	revocationDataVerifier := validation.NewDefaultRevocationDataVerifier()
	revocationDataVerifier.SetAcceptRevocationCertificatesWithoutRevocation(true)
	revocationDataVerifier.SetAcceptTimestampCertificatesWithoutRevocation(true)
	return revocationDataVerifier
}

// AssertExtendToTLevelPossible verifies whether extension of signatures to T-level is possible.
// Port of #assertExtendToTLevelPossible.
func (c *SignatureRequirementsChecker[TP]) AssertExtendToTLevelPossible(signatures []validation.AdvancedSignature) {
	c.assertTLevelIsHighest(signatures)
	c.assertHasNoEmbeddedEvidenceRecords(signatures)
}

// assertTLevelIsHighest checks whether across signatures the T-level is highest and T-level
// augmentation can be performed. Port of the protected #assertTLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) assertTLevelIsHighest(signatures []validation.AdvancedSignature) {
	if c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		c.checkTLevelIsHighest(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Error on signature augmentation to T-level.")
		mustAlert(c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel(), status)
	}
}

// checkTLevelIsHighest verifies whether the signature has maximum B- or T-level. Port of the
// protected #checkTLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) checkTLevelIsHighest(signature validation.AdvancedSignature, status *validation.SignatureStatus) {
	if c.HasLTLevelOrHigher(signature) {
		status.AddRelatedTokenAndErrorMessage(signature, "The signature is already extended with a higher level.")
	}
}

// HasLTLevelOrHigher checks if the signature has LTA-level. Port of #hasLTLevelOrHigher.
func (c *SignatureRequirementsChecker[TP]) HasLTLevelOrHigher(signature validation.AdvancedSignature) bool {
	return signature.HasLTAProfile() ||
		((signature.HasLTProfile() || signature.HasCProfile()) && !signature.AreAllSelfSignedCertificates() && signature.HasTProfile())
}

// AssertExtendToLTLevelPossible verifies whether extension of signatures to LT-level is
// possible. Port of #assertExtendToLTLevelPossible.
func (c *SignatureRequirementsChecker[TP]) AssertExtendToLTLevelPossible(signatures []validation.AdvancedSignature) {
	c.assertLTLevelIsHighest(signatures)
	c.assertHasNoEmbeddedEvidenceRecords(signatures)
}

// assertLTLevelIsHighest checks whether across signatures the LT-level is highest and LT-level
// augmentation can be performed. Port of the protected #assertLTLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) assertLTLevelIsHighest(signatures []validation.AdvancedSignature) {
	if c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		c.checkLTLevelIsHighest(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Error on signature augmentation to LT-level.")
		mustAlert(c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel(), status)
	}
}

// checkLTLevelIsHighest verifies whether the signature has maximum B-, T- or LT-level. Port of
// the protected #checkLTLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) checkLTLevelIsHighest(signature validation.AdvancedSignature, status *validation.SignatureStatus) {
	if c.HasLTALevelOrHigher(signature) {
		status.AddRelatedTokenAndErrorMessage(signature, "The signature is already extended with a higher level.")
	}
}

// HasLTALevelOrHigher checks if the signature has LTA-level. Port of #hasLTALevelOrHigher.
func (c *SignatureRequirementsChecker[TP]) HasLTALevelOrHigher(signature validation.AdvancedSignature) bool {
	return signature.HasLTAProfile()
}

// AssertCertificateChainValidForLTLevel checks whether across signatures the corresponding
// certificate chains require revocation data for LT-level augmentation. Port of
// #assertCertificateChainValidForLTLevel.
func (c *SignatureRequirementsChecker[TP]) AssertCertificateChainValidForLTLevel(signatures []validation.AdvancedSignature) {
	c.assertCertificateChainValid(signatures, "LT")
}

// AssertCertificateChainValidForCLevel checks whether across signatures the corresponding
// certificate chains require revocation data for C-level augmentation. Port of
// #assertCertificateChainValidForCLevel.
func (c *SignatureRequirementsChecker[TP]) AssertCertificateChainValidForCLevel(signatures []validation.AdvancedSignature) {
	c.assertCertificateChainValid(signatures, "C")
}

// AssertCertificateChainValidForXLLevel checks whether across signatures the corresponding
// certificate chains require revocation data for XL-level augmentation. Port of
// #assertCertificateChainValidForXLLevel.
func (c *SignatureRequirementsChecker[TP]) AssertCertificateChainValidForXLLevel(signatures []validation.AdvancedSignature) {
	c.assertCertificateChainValid(signatures, "XL")
}

// assertCertificateChainValid ports the private #assertCertificateChainValid.
func (c *SignatureRequirementsChecker[TP]) assertCertificateChainValid(signatures []validation.AdvancedSignature, targetLevel string) {
	c.assertCertificatePresent(signatures, targetLevel)
	c.assertCertificatesAreNotSelfSigned(signatures, targetLevel)
}

// assertCertificatePresent ports the private #assertCertificatePresent.
func (c *SignatureRequirementsChecker[TP]) assertCertificatePresent(signatures []validation.AdvancedSignature, targetLevel string) {
	if c.certificateVerifier.AugmentationAlertOnSignatureWithoutCertificates() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		if signature.CertificateSource().NumberOfCertificates() == 0 {
			status.AddRelatedTokenAndErrorMessage(signature, "The signature does not contain certificates.")
		}
	}
	if !status.IsEmpty() {
		status.SetMessage(fmt.Sprintf("Error on signature augmentation to %s-level.", targetLevel))
		mustAlert(c.certificateVerifier.AugmentationAlertOnSignatureWithoutCertificates(), status)
	}
}

// assertCertificatesAreNotSelfSigned ports the private #assertCertificatesAreNotSelfSigned.
func (c *SignatureRequirementsChecker[TP]) assertCertificatesAreNotSelfSigned(signatures []validation.AdvancedSignature, targetLevel string) {
	if c.certificateVerifier.AugmentationAlertOnSelfSignedCertificateChains() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		if signature.AreAllSelfSignedCertificates() {
			status.AddRelatedTokenAndErrorMessage(signature, "The signature contains only self-signed certificate chains.")
		}
	}
	if !status.IsEmpty() {
		status.SetMessage(fmt.Sprintf("Error on signature augmentation to %s-level.", targetLevel))
		mustAlert(c.certificateVerifier.AugmentationAlertOnSelfSignedCertificateChains(), status)
	}
}

// AssertExtendToCLevelPossible verifies whether extension of signatures to C-level is possible.
// Port of #assertExtendToCLevelPossible.
func (c *SignatureRequirementsChecker[TP]) AssertExtendToCLevelPossible(signatures []validation.AdvancedSignature) {
	c.assertCLevelIsHighest(signatures)
	c.assertHasNoEmbeddedEvidenceRecords(signatures)
}

// assertCLevelIsHighest checks whether across signatures the C-level is highest and C-level
// augmentation can be performed. Port of the protected #assertCLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) assertCLevelIsHighest(signatures []validation.AdvancedSignature) {
	if c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		c.checkCLevelIsHighest(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Error on signature augmentation to C-level.")
		mustAlert(c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel(), status)
	}
}

// checkCLevelIsHighest verifies whether the signature has maximum B-, T- or LT-level. Port of
// the protected #checkCLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) checkCLevelIsHighest(signature validation.AdvancedSignature, status *validation.SignatureStatus) {
	if c.HasXLevelOrHigher(signature) {
		status.AddRelatedTokenAndErrorMessage(signature, "The signature is already extended with a higher level.")
	}
}

// HasXLevelOrHigher checks if the signature has LTA-level. Port of #hasXLevelOrHigher.
func (c *SignatureRequirementsChecker[TP]) HasXLevelOrHigher(signature validation.AdvancedSignature) bool {
	return signature.HasXProfile() || signature.HasAProfile() ||
		(signature.HasXLProfile() && !signature.AreAllSelfSignedCertificates() && signature.HasTProfile())
}

// AssertExtendToXLevelPossible verifies whether extension of signatures to X-level is possible.
// Port of #assertExtendToXLevelPossible.
func (c *SignatureRequirementsChecker[TP]) AssertExtendToXLevelPossible(signatures []validation.AdvancedSignature) {
	c.assertXLevelIsHighest(signatures)
	c.assertHasNoEmbeddedEvidenceRecords(signatures)
}

// assertXLevelIsHighest checks whether across signatures the X-level is highest and X-level
// augmentation can be performed. Port of the protected #assertXLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) assertXLevelIsHighest(signatures []validation.AdvancedSignature) {
	if c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		c.checkXLevelIsHighest(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Error on signature augmentation to X-level.")
		mustAlert(c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel(), status)
	}
}

// checkXLevelIsHighest verifies whether the signature has maximum B-, T- or LT-level. Port of
// the protected #checkXLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) checkXLevelIsHighest(signature validation.AdvancedSignature, status *validation.SignatureStatus) {
	if c.HasXLLevelOrHigher(signature) {
		status.AddRelatedTokenAndErrorMessage(signature, "The signature is already extended with a higher level.")
	}
}

// HasXLLevelOrHigher checks if the signature has LTA-level. Port of #hasXLLevelOrHigher.
func (c *SignatureRequirementsChecker[TP]) HasXLLevelOrHigher(signature validation.AdvancedSignature) bool {
	return signature.HasAProfile() ||
		(signature.HasXLProfile() && !signature.AreAllSelfSignedCertificates() && signature.HasTProfile() && signature.HasXProfile())
}

// AssertExtendToXLLevelPossible verifies whether extension of signatures to XL-level is
// possible. Port of #assertExtendToXLLevelPossible.
func (c *SignatureRequirementsChecker[TP]) AssertExtendToXLLevelPossible(signatures []validation.AdvancedSignature) {
	c.assertXLLevelIsHighest(signatures)
	c.assertHasNoEmbeddedEvidenceRecords(signatures)
}

// assertXLLevelIsHighest checks whether across signatures the XL-level is highest and XL-level
// augmentation can be performed. Port of the protected #assertXLLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) assertXLLevelIsHighest(signatures []validation.AdvancedSignature) {
	if c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		c.checkXLLevelIsHighest(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Error on signature augmentation to XL-level.")
		mustAlert(c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel(), status)
	}
}

// checkXLLevelIsHighest verifies whether the signature has maximum X-level. Port of the
// protected #checkXLLevelIsHighest.
func (c *SignatureRequirementsChecker[TP]) checkXLLevelIsHighest(signature validation.AdvancedSignature, status *validation.SignatureStatus) {
	if c.HasALevelOrHigher(signature) {
		status.AddRelatedTokenAndErrorMessage(signature, "The signature is already extended with a higher level.")
	}
}

// HasALevelOrHigher checks if the signature has A-level. Port of #hasALevelOrHigher.
func (c *SignatureRequirementsChecker[TP]) HasALevelOrHigher(signature validation.AdvancedSignature) bool {
	return c.HasLTALevelOrHigher(signature)
}

// AssertExtendToLTALevelPossible verifies whether extension of signatures to LTA-level is
// possible. Port of #assertExtendToLTALevelPossible.
func (c *SignatureRequirementsChecker[TP]) AssertExtendToLTALevelPossible(signatures []validation.AdvancedSignature) {
	c.assertHasNoEmbeddedEvidenceRecords(signatures)
}

// assertHasNoEmbeddedEvidenceRecords checks whether across signatures the T-level is highest and
// T-level augmentation can be performed. Port of the protected #assertHasNoEmbeddedEvidenceRecords.
func (c *SignatureRequirementsChecker[TP]) assertHasNoEmbeddedEvidenceRecords(signatures []validation.AdvancedSignature) {
	if c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel() == nil {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signatures {
		c.checkHasEmbeddedEvidenceRecords(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Error on signature augmentation")
		mustAlert(c.certificateVerifier.AugmentationAlertOnHigherSignatureLevel(), status)
	}
}

// checkHasEmbeddedEvidenceRecords verifies whether the signature has an embedded evidence
// record. Port of the protected #checkHasEmbeddedEvidenceRecords.
func (c *SignatureRequirementsChecker[TP]) checkHasEmbeddedEvidenceRecords(signature validation.AdvancedSignature, status *validation.SignatureStatus) {
	if c.HasEmbeddedEvidenceRecords(signature) {
		status.AddRelatedTokenAndErrorMessage(signature, "The signature is preserved by an embedded evidence record.")
	}
}

// HasEmbeddedEvidenceRecords checks if the signature has embedded evidence records. Port of
// #hasEmbeddedEvidenceRecords.
func (c *SignatureRequirementsChecker[TP]) HasEmbeddedEvidenceRecords(signature validation.AdvancedSignature) bool {
	return len(signature.EmbeddedEvidenceRecords()) > 0
}

// AssertSignaturesValid verifies cryptographical validity of the signatures. Port of
// #assertSignaturesValid.
func (c *SignatureRequirementsChecker[TP]) AssertSignaturesValid(signatures []validation.AdvancedSignature) {
	if c.certificateVerifier.AlertOnInvalidSignature() == nil {
		return
	}

	var signaturesToValidate []validation.AdvancedSignature
	for _, s := range signatures {
		if !c.isSignatureGeneratedWithoutCertificate(s) {
			signaturesToValidate = append(signaturesToValidate, s)
		}
	}
	if len(signaturesToValidate) == 0 {
		return
	}

	status := validation.NewSignatureStatus()
	for _, signature := range signaturesToValidate {
		signatureCryptographicVerification := signature.SignatureCryptographicVerification()
		if !signatureCryptographicVerification.IsSignatureValid() {
			errorMessage := signatureCryptographicVerification.ErrorMessage()
			suffix := "."
			if errorMessage != "" {
				suffix = " / " + errorMessage
			}
			status.AddRelatedTokenAndErrorMessage(signature, "Cryptographic signature verification has failed"+suffix)
		}
	}
	if !status.IsEmpty() {
		status.SetMessage("Error on signature augmentation")
		mustAlert(c.certificateVerifier.AlertOnInvalidSignature(), status)
	}
}

// mustAlert executes statusAlert.Alert(status), repanicking any returned error. Java's
// StatusAlert#alert is void and may throw an unchecked AlertException (e.g. via
// ExceptionOnStatusAlert); the Go alert.Alert(T) port returns an error instead of throwing (see
// alert.Alert's own doc comment). Every call site here has no error channel to propagate that
// through (matching Java's void-returning assertXxx methods), so a non-nil error is repanicked -
// the same technique validation.SignatureValidationAlerter uses (spi/validation/signature_validation_alerter.go).
func mustAlert(statusAlert alert.StatusAlert, status alert.Status) {
	if err := statusAlert.Alert(status); err != nil {
		panic(err)
	}
}
