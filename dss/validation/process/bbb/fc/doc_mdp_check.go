// Ported from dss-validation/.../validation/process/bbb/fc/checks/DocMDPCheck.java (DSS 6.5.RC1).
package fc

import (
	"fmt"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// DocMDPCheck verifies a signature according to given permissions for the document in /DocMDP.
type DocMDPCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	pdfRevision *diagnostic.PDFRevisionWrapper
}

// NewDocMDPCheck is the default constructor.
func NewDocMDPCheck(i18nProvider *i18n.Provider, result *process.Result[*drjaxb.XmlFC],
	pdfRevision *diagnostic.PDFRevisionWrapper, constraint policy.LevelRule) *DocMDPCheck {
	c := &DocMDPCheck{pdfRevision: pdfRevision}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *DocMDPCheck) Process() bool {
	if !c.pdfRevision.ArePdfObjectModificationsDetected() {
		return true
	}
	docMDPPermissions := c.pdfRevision.DocMDPPermissions()
	if docMDPPermissions == "" {
		return true
	}
	switch docMDPPermissions {
	case enumerations.CertificationPermissionNoChangePermitted:
		if len(c.pdfRevision.PdfSignatureOrFormFillChanges()) > 0 ||
			len(c.pdfRevision.PdfAnnotationChanges()) > 0 ||
			len(c.pdfRevision.PdfUndefinedChanges()) > 0 {
			return false
		}
	case enumerations.CertificationPermissionMinimalChangesPermitted:
		if len(c.pdfRevision.PdfAnnotationChanges()) > 0 || len(c.pdfRevision.PdfUndefinedChanges()) > 0 {
			return false
		}
	case enumerations.CertificationPermissionChangesPermitted:
		if len(c.pdfRevision.PdfUndefinedChanges()) > 0 {
			return false
		}
	default:
		panic(fmt.Sprintf("The value '%s' is not supported!", docMDPPermissions))
	}
	return true
}

// MessageTag returns the constraint message i18n key.
func (c *DocMDPCheck) MessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCISVADMDPD }

// ErrorMessageTag returns the error message i18n key.
func (c *DocMDPCheck) ErrorMessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCISVADMDPDANS }

// FailedIndicationForConclusion returns the Indication on failure.
func (c *DocMDPCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *DocMDPCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
