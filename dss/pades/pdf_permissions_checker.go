// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfPermissionsChecker.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). slf4j is dropped, per PORTING.md; the sole non-debug log
// (an INFO about the deprecated usage-rights signature) is kept as a comment where it fired.
package pades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/pades/alerts"
)

// PdfPermissionsChecker is used to verify permissions of a PDF document and to check whether
// modifications are allowed. Port of the PdfPermissionsChecker class.
type PdfPermissionsChecker struct {
	// alertOnForbiddenSignatureCreation indicates a behavior for creation of a new signature in
	// a document that does not permit a new signature creation.
	// Default: ProtectedDocumentExceptionOnStatusAlert.
	alertOnForbiddenSignatureCreation alert.StatusAlert
}

// NewPdfPermissionsChecker instantiates the checker with the default configuration.
// Port of the default constructor.
func NewPdfPermissionsChecker() *PdfPermissionsChecker {
	return &PdfPermissionsChecker{
		alertOnForbiddenSignatureCreation: alerts.NewProtectedDocumentExceptionOnStatusAlert(),
	}
}

// SetAlertOnForbiddenSignatureCreation sets a behavior to follow when creating a new signature
// in a document that forbids the creation of new signatures. Default: throw the exception.
// Port of #setAlertOnForbiddenSignatureCreation.
func (c *PdfPermissionsChecker) SetAlertOnForbiddenSignatureCreation(alertOnForbiddenSignatureCreation alert.StatusAlert) {
	c.alertOnForbiddenSignatureCreation = alertOnForbiddenSignatureCreation
}

// CheckDocumentPermissions checks if the document has the necessary permissions for the
// signature operation. Port of #checkDocumentPermissions.
func (c *PdfPermissionsChecker) CheckDocumentPermissions(documentReader PdfDocumentReader, fieldParameters *SignatureFieldParameters) {
	if !documentReader.IsEncrypted() || documentReader.IsOpenWithOwnerAccess() {
		// permissions are applied only for encrypted documents with user-access
		return
	}
	if c.isSignatureFieldFillIn(fieldParameters) {
		if !documentReader.CanFillSignatureForm() {
			c.alertOnForbiddenSignatureCreationMessage("PDF Permissions dictionary does not allow fill in interactive form fields, " +
				"including existing signature fields when document is open with user-access!")
		}
	} else if !documentReader.CanCreateSignatureField() {
		c.alertOnForbiddenSignatureCreationMessage("PDF Permissions dictionary does not allow modification or creation interactive form fields, " +
			"including signature fields when document is open with user-access!")
	}
}

func (c *PdfPermissionsChecker) isSignatureFieldFillIn(fieldParameters *SignatureFieldParameters) bool {
	return fieldParameters.FieldId() != ""
}

// CheckSignatureRestrictionDictionaries verifies whether a new signature is permitted.
// Port of #checkSignatureRestrictionDictionaries.
func (c *PdfPermissionsChecker) CheckSignatureRestrictionDictionaries(documentReader PdfDocumentReader, fieldParameters *SignatureFieldParameters) {
	certificationPermission := documentReader.CertificationPermission()
	if c.isDocumentChangeForbidden(certificationPermission) {
		c.alertOnForbiddenSignatureCreationMessage("DocMDP dictionary does not permit a new signature creation!")
	}
	if documentReader.IsUsageRightsSignaturePresent() {
		// Deprecated. See ISO 32000-2: when a usage rights signature is present, it is up to
		// the PDF processor or to the signature handler to process it or not.
		// Upstream logs "A usage rights signature is present. The feature is deprecated and
		// the entry is not handled." at INFO.
	}

	signatureFieldID := fieldParameters.FieldId()

	sigDictionaries, err := documentReader.ExtractSigDictionaries()
	if err != nil {
		// Upstream logs "An error occurred while reading signature dictionary entries : {}".
		return
	}
	for _, entry := range sigDictionaries {
		fieldMDP := entry.SignatureDictionary.FieldMDP()
		if fieldMDP != nil && c.isSignatureFieldCreationForbidden(fieldMDP, signatureFieldID) {
			c.alertOnForbiddenSignatureCreationMessage("FieldMDP dictionary does not permit a new signature creation!")
		}
	}

	for _, entry := range sigDictionaries {
		for _, signatureField := range entry.Fields {
			lockDict := signatureField.LockDictionary()
			if lockDict != nil && lockDict.CertificationPermission() != "" &&
				c.isSignatureFieldCreationForbidden(lockDict, signatureFieldID) {
				c.alertOnForbiddenSignatureCreationMessage("Lock dictionary does not permit a new signature creation!")
			}
		}
	}
}

// isDocumentChangeForbidden verifies and returns whether changes within a document are
// forbidden according to the defined certificationPermission. Port of #isDocumentChangeForbidden.
func (c *PdfPermissionsChecker) isDocumentChangeForbidden(certificationPermission enumerations.CertificationPermission) bool {
	return certificationPermission == enumerations.CertificationPermissionNoChangePermitted
}

// alertOnForbiddenSignatureCreationMessage executes alertOnForbiddenSignatureCreation with the
// given message. Port of the protected #alertOnForbiddenSignatureCreation(String).
func (c *PdfPermissionsChecker) alertOnForbiddenSignatureCreationMessage(message string) {
	status := alert.NewMessageStatus()
	status.SetMessage(fmt.Sprintf("The creation of new signatures is not permitted in the current document. Reason : %s", message))
	if err := c.alertOnForbiddenSignatureCreation.Alert(status); err != nil {
		panic(err)
	}
}

// isSignatureFieldCreationForbidden checks and returns whether a signature field creation is
// forbidden according to the given configuration of signatureFieldID.
// Port of #isSignatureFieldCreationForbidden.
func (c *PdfPermissionsChecker) isSignatureFieldCreationForbidden(sigFieldPermissions *SigFieldPermissions, signatureFieldID string) bool {
	switch sigFieldPermissions.Action() {
	case enumerations.PdfLockActionAll:
		return true
	case enumerations.PdfLockActionInclude:
		if signatureFieldID == "" {
			return false
		}
		if containsString(sigFieldPermissions.Fields(), signatureFieldID) {
			return true
		}
	case enumerations.PdfLockActionExclude:
		if signatureFieldID == "" {
			return true
		}
		if !containsString(sigFieldPermissions.Fields(), signatureFieldID) {
			return true
		}
	default:
		panic(fmt.Sprintf("The action value '%s' is not supported!", sigFieldPermissions.Action()))
	}
	certificationPermission := sigFieldPermissions.CertificationPermission()
	return c.isDocumentChangeForbidden(certificationPermission)
}
