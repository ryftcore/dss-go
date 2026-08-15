// Ported from dss-validation/.../validation/process/bbb/fc/checks/SigFieldLockCheck.java (DSS 6.5.RC1).
package fc

import (
	"fmt"

	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SigFieldLockCheck verifies a signature according to given permissions for the signature field
// in /SigFieldLock.
type SigFieldLockCheck struct {
	AbstractPdfLockDictionaryCheck
}

// NewSigFieldLockCheck is the default constructor.
func NewSigFieldLockCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *SigFieldLockCheck {
	c := &SigFieldLockCheck{}
	c.InitAbstractPdfLockDictionaryCheck(i18nProvider, result, pdfRevision, pdfRevision.SigFieldLock(), constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *SigFieldLockCheck) Process() bool {
	if !c.AbstractPdfLockDictionaryCheck.Process() {
		return false
	}
	if !c.PdfRevision.ArePdfObjectModificationsDetected() {
		return true
	}
	if c.PdfLockDictionary == nil {
		return true
	}

	// optional
	if c.PdfLockDictionary.Permissions != nil {
		switch c.PdfLockDictionary.Permissions.CertificationPermission() {
		case enumerations.CertificationPermission_NO_CHANGE_PERMITTED:
			if len(c.PdfRevision.PdfSignatureOrFormFillChanges()) > 0 ||
				len(c.PdfRevision.PdfAnnotationChanges()) > 0 ||
				len(c.PdfRevision.PdfUndefinedChanges()) > 0 {
				return false
			}
		case enumerations.CertificationPermission_MINIMAL_CHANGES_PERMITTED:
			if len(c.PdfRevision.PdfAnnotationChanges()) > 0 || len(c.PdfRevision.PdfUndefinedChanges()) > 0 {
				return false
			}
		case enumerations.CertificationPermission_CHANGES_PERMITTED:
			if len(c.PdfRevision.PdfUndefinedChanges()) > 0 {
				return false
			}
		default:
			panic(fmt.Sprintf("The value '%s' is not supported!", c.PdfLockDictionary.Permissions.CertificationPermission()))
		}
	}
	return true
}

// MessageTag returns the constraint message i18n key.
func (c *SigFieldLockCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_FC_ISVASFLD }

// ErrorMessageTag returns the error message i18n key.
func (c *SigFieldLockCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISVASFLD_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *SigFieldLockCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *SigFieldLockCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
