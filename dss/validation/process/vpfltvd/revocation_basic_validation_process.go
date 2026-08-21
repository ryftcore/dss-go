// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/RevocationBasicValidationProcess.java (DSS 6.5.RC1).
package vpfltvd

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfbs"
)

// RevocationBasicValidationProcess performs basic validation of a revocation
// data.
type RevocationBasicValidationProcess struct {
	*vpfbs.AbstractBasicValidationProcess[*jaxb.XmlRevocationBasicValidation]

	// revocationData is the revocation data to be validated.
	revocationData *diagnostic.RevocationWrapper
}

// NewRevocationBasicValidationProcess is the default constructor. Port of
// RevocationBasicValidationProcess(I18nProvider, DiagnosticData, RevocationWrapper, Map).
func NewRevocationBasicValidationProcess(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	revocationData *diagnostic.RevocationWrapper, bbbs map[string]*jaxb.XmlBasicBuildingBlocks) *RevocationBasicValidationProcess {
	xmlResult := &jaxb.XmlRevocationBasicValidation{}
	result := process.NewResult(xmlResult, &xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)
	c := &RevocationBasicValidationProcess{
		AbstractBasicValidationProcess: vpfbs.NewAbstractBasicValidationProcess(i18nProvider, result, diagnosticData, revocationData, bbbs),
		revocationData:                 revocationData,
	}
	c.InitAbstractBasicValidationProcess(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *RevocationBasicValidationProcess) Title() i18n.MessageTag {
	return i18n.MessageTag_VPFRVC
}

// AddAdditionalInfo adds additional info to the chain. Port of
// addAdditionalInfo().
func (c *RevocationBasicValidationProcess) AddAdditionalInfo() {
	id := c.revocationData.Id()
	c.Result.Value.Id = &id
}
