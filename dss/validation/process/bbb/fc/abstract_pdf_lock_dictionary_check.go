// Ported from dss-validation/.../validation/process/bbb/fc/checks/AbstractPdfLockDictionaryCheck.java (DSS 6.5.RC1).
package fc

import (
	"fmt"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractPdfLockDictionaryCheck is the base for PDF lock dictionary validation. Concrete
// checks (SigFieldLockCheck, FieldMDPCheck) embed it and provide MessageTag()/ErrorMessageTag().
type AbstractPdfLockDictionaryCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	PdfRevision       *diagnostic.PDFRevisionWrapper
	PdfLockDictionary *diagjaxb.XmlPDFLockDictionary
}

// InitAbstractPdfLockDictionaryCheck wires the shared state; called by the concrete
// constructor before InitChainItem.
func (c *AbstractPdfLockDictionaryCheck) InitAbstractPdfLockDictionaryCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*drjaxb.XmlFC], pdfRevision *diagnostic.PDFRevisionWrapper, pdfLockDictionary *diagjaxb.XmlPDFLockDictionary,
	constraint policy.LevelRule) {
	c.PdfRevision = pdfRevision
	c.PdfLockDictionary = pdfLockDictionary
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
}

// Process performs the check. Overridden by SigFieldLockCheck, which calls this as its base
// process() (super.process()).
func (c *AbstractPdfLockDictionaryCheck) Process() bool {
	if !c.PdfRevision.ArePdfObjectModificationsDetected() {
		return true
	}
	if c.PdfLockDictionary == nil {
		return true
	}

	modifiedFieldNames := c.PdfRevision.ModifiedFieldNames()
	if len(modifiedFieldNames) == 0 {
		return true
	}

	lockedFields := c.PdfLockDictionary.Field
	if c.PdfLockDictionary.Action != nil {
		switch c.PdfLockDictionary.Action.PdfLockAction() {
		case enumerations.PdfLockAction_ALL:
			return false
		case enumerations.PdfLockAction_EXCLUDE:
			for _, fieldName := range modifiedFieldNames {
				if !containsString(lockedFields, fieldName) {
					return false
				}
			}
			return true
		case enumerations.PdfLockAction_INCLUDE:
			for _, fieldName := range modifiedFieldNames {
				if containsString(lockedFields, fieldName) {
					return false
				}
			}
			return true
		default:
			panic(fmt.Sprintf("The value '%s' is not supported!", c.PdfLockDictionary.Action.PdfLockAction()))
		}
	}
	return true
}

func containsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
