// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/SignatureValidationAlerter.java (DSS 6.5.RC1).
//
// This type implements the already-landed ValidationAlerter interface (spi/validation/validation_alerter.go),
// whose AssertXxx methods return no error (matching Java's void signature literally). Java's
// underlying alert.alert(status) call is itself void but may throw an unchecked
// AlertException (e.g. via ExceptionOnStatusAlert); the Go alert.Alert(T) port returns an error
// instead of throwing (see alert.Alert's own doc comment: "the returned error carries what Java
// would throw as an AlertException"). Since ValidationAlerter's fixed signature here has no
// error channel to propagate that through, a non-nil error from Alert(status) is repanicked,
// reproducing Java's unchecked-exception propagation out of a void method.
//
// CROSS-CHUNK DEPENDENCY: CertificateVerifier (Java spi.validation.CertificateVerifier) is
// already forward-declared opaquely by five sibling files in this package (see e.g.
// advanced_signature.go, signature_validation_context.go). This file additionally requires the
// seven getAlertOnXxx() accessors below, none of which any already-landed file documents;
// they're spelled out here, matched 1:1 to the Java getters (get-prefix dropped per
// PORTING.md), for CertificateVerifier's owning chunk (VAL-D) to supply:
//
//	AlertOnMissingRevocationData() alert.StatusAlert
//	AlertOnUncoveredPOE() alert.StatusAlert
//	AlertOnInvalidTimestamp() alert.StatusAlert
//	AlertOnRevokedCertificate() alert.StatusAlert
//	AlertOnNoRevocationAfterBestSignatureTime() alert.StatusAlert
//	AlertOnExpiredCertificate() alert.StatusAlert
//	AlertOnNotYetValidCertificate() alert.StatusAlert
//
// getCertificateVerifier() on SignatureValidationContext is already landed as the unexported
// method of that exact name (signature_validation_context.go), callable here since both files
// share this package.
package validation

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// SignatureValidationAlerter uses a SignatureValidationContext to perform validation and
// executes alerts based on the validation result. The configuration of the alerts and their
// behavior is defined within CertificateVerifier. If an alert is not defined, the execution of
// the corresponding check is skipped.
type SignatureValidationAlerter struct {
	// validationContext performs the execution.
	validationContext *SignatureValidationContext

	// signingOperation is the nature of the current operation the verification is done for.
	signingOperation enumerations.SigningOperation
}

// NewSignatureValidationAlerter is the default constructor to instantiate the alerter. Port of
// the SignatureValidationAlerter(SignatureValidationContext) constructor.
func NewSignatureValidationAlerter(validationContext *SignatureValidationContext) *SignatureValidationAlerter {
	return &SignatureValidationAlerter{validationContext: validationContext}
}

// SetSigningOperation (optional) sets the current operation kind to provide a user-friendly
// error message. Port of setSigningOperation(SigningOperation).
func (a *SignatureValidationAlerter) SetSigningOperation(signingOperation enumerations.SigningOperation) {
	a.signingOperation = signingOperation
}

// AssertAllRequiredRevocationDataPresent verifies whether all processed certificates have a
// revocation data.
func (a *SignatureValidationAlerter) AssertAllRequiredRevocationDataPresent() {
	alertOnMissingRevocationData := a.validationContext.getCertificateVerifier().AlertOnMissingRevocationData()
	if alertOnMissingRevocationData == nil {
		return
	}

	status, _ := a.validationContext.allRequiredRevocationDataPresent()
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnMissingRevocationData, status)
	}
}

// AssertAllPOECoveredByRevocationData verifies whether all POE (timestamp tokens) are covered by
// a revocation data.
func (a *SignatureValidationAlerter) AssertAllPOECoveredByRevocationData() {
	alertOnUncoveredPOE := a.validationContext.getCertificateVerifier().AlertOnUncoveredPOE()
	if alertOnUncoveredPOE == nil {
		return
	}

	status, _ := a.validationContext.allPOECoveredByRevocationData()
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnUncoveredPOE, status)
	}
}

// AssertAllTimestampsValid verifies whether all processed timestamps are valid and intact.
func (a *SignatureValidationAlerter) AssertAllTimestampsValid() {
	alertOnInvalidTimestamp := a.validationContext.getCertificateVerifier().AlertOnInvalidTimestamp()
	if alertOnInvalidTimestamp == nil {
		return
	}

	status := a.validationContext.allTimestampsValid()
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnInvalidTimestamp, status)
	}
}

// AssertCertificateNotRevoked verifies whether the certificate is not revoked.
func (a *SignatureValidationAlerter) AssertCertificateNotRevoked(certificateToken *model.CertificateToken) {
	alertOnRevokedCertificate := a.validationContext.getCertificateVerifier().AlertOnRevokedCertificate()
	if alertOnRevokedCertificate == nil {
		return
	}

	status := a.validationContext.certificateNotRevoked(certificateToken)
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnRevokedCertificate, status)
	}
}

// AssertAllSignatureCertificatesNotRevoked verifies recursively whether none of the signature's
// certificate chain certificates are revoked.
func (a *SignatureValidationAlerter) AssertAllSignatureCertificatesNotRevoked() {
	alertOnRevokedCertificate := a.validationContext.getCertificateVerifier().AlertOnRevokedCertificate()
	if alertOnRevokedCertificate == nil {
		return
	}

	status := a.validationContext.allSignatureCertificatesNotRevoked()
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnRevokedCertificate, status)
	}
}

// AssertAllSignatureCertificateHaveFreshRevocationData verifies whether for all signature's
// certificate chain certificates there is a fresh revocation data, after the earliest available
// timestamp token production time.
func (a *SignatureValidationAlerter) AssertAllSignatureCertificateHaveFreshRevocationData() {
	alertOnNoRevocationAfterBestSignatureTime := a.validationContext.getCertificateVerifier().AlertOnNoRevocationAfterBestSignatureTime()
	if alertOnNoRevocationAfterBestSignatureTime == nil {
		return
	}

	status := a.validationContext.allSignatureCertificateHaveFreshRevocationData()
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnNoRevocationAfterBestSignatureTime, status)
	}
}

// AssertAllSignaturesNotExpired verifies whether all signatures added to the ValidationContext
// are not yet expired.
func (a *SignatureValidationAlerter) AssertAllSignaturesNotExpired() {
	alertOnExpiredCertificate := a.validationContext.getCertificateVerifier().AlertOnExpiredCertificate()
	if alertOnExpiredCertificate == nil {
		return
	}

	status := a.validationContext.allSignaturesNotExpired()
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnExpiredCertificate, status)
	}
}

// AssertCertificateNotExpired verifies whether the certificate token is not yet expired.
func (a *SignatureValidationAlerter) AssertCertificateNotExpired(certificateToken *model.CertificateToken) {
	alertOnExpiredCertificate := a.validationContext.getCertificateVerifier().AlertOnExpiredCertificate()
	if alertOnExpiredCertificate == nil {
		return
	}

	status := a.validationContext.certificateNotExpired(certificateToken)
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnExpiredCertificate, status)
	}
}

// AssertAllSignaturesAreYetValid verifies whether all signatures added to the ValidationContext
// have been produced with yet valid certificates.
func (a *SignatureValidationAlerter) AssertAllSignaturesAreYetValid() {
	alertOnNotYetValidCertificate := a.validationContext.getCertificateVerifier().AlertOnNotYetValidCertificate()
	if alertOnNotYetValidCertificate == nil {
		return
	}

	status := a.validationContext.allSignaturesAreYetValid()
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnNotYetValidCertificate, status)
	}
}

// AssertCertificateIsYetValid verifies whether the certificate token is yet valid.
func (a *SignatureValidationAlerter) AssertCertificateIsYetValid(certificateToken *model.CertificateToken) {
	alertOnNotYetValidCertificate := a.validationContext.getCertificateVerifier().AlertOnNotYetValidCertificate()
	if alertOnNotYetValidCertificate == nil {
		return
	}

	status := a.validationContext.certificateIsYetValid(certificateToken)
	if !status.IsEmpty() {
		a.populateMessage(status)
		a.alert(alertOnNotYetValidCertificate, status)
	}
}

// signatureValidationAlerterMessage is the subset of alert.MessageStatus's promoted API this
// file needs to augment a status's message before alerting. Both *TokenStatus and
// *SignatureStatus satisfy it via their embedded alert.ObjectStatus/alert.MessageStatus.
type signatureValidationAlerterMessage interface {
	Message() string
	SetMessage(message string)
}

// populateMessage augments the validation message with the information about the currently
// performing operation kind. Port of populateMessage(MessageStatus).
func (a *SignatureValidationAlerter) populateMessage(status signatureValidationAlerterMessage) {
	if status == nil || a.signingOperation == "" {
		return
	}
	originalMessage := status.Message()
	switch a.signingOperation {
	case enumerations.SigningOperation_SIGN, enumerations.SigningOperation_COUNTER_SIGN:
		status.SetMessage(fmt.Sprintf("Error on signature creation : %s", originalMessage))
	case enumerations.SigningOperation_EXTEND:
		status.SetMessage(fmt.Sprintf("Error on signature augmentation : %s", originalMessage))
	case enumerations.SigningOperation_TIMESTAMP:
		status.SetMessage(fmt.Sprintf("Error on timestamp : %s", originalMessage))
	case enumerations.SigningOperation_ADD_EVIDENCE_RECORD:
		status.SetMessage(fmt.Sprintf("Error on evidence record incorporation : %s", originalMessage))
	case enumerations.SigningOperation_ADD_SIG_POLICY_STORE:
		status.SetMessage(fmt.Sprintf("Error on signature policy store incorporation : %s", originalMessage))
	default:
		panic(fmt.Sprintf("The operation '%s' is not supported!", a.signingOperation))
	}
}

// alert executes the alert on the given status, repanicking any error returned by
// alert.Alert(status) - see this file's header comment for why.
func (a *SignatureValidationAlerter) alert(statusAlert alert.StatusAlert, status alert.Status) {
	if err := statusAlert.Alert(status); err != nil {
		panic(err)
	}
}

// compile-time assertion: a SignatureValidationAlerter is a ValidationAlerter.
var _ ValidationAlerter = (*SignatureValidationAlerter)(nil)
