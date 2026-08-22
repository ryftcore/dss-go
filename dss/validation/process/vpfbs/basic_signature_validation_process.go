// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/BasicSignatureValidationProcess.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// BasicSignatureValidationProcess is the signature validation process at
// validation time as per EN 319 102-1 ch. "5.3 Validation process for Basic
// Signatures".
type BasicSignatureValidationProcess struct {
	*AbstractBasicValidationProcess[*jaxb.XmlValidationProcessBasicSignature]

	// xmlTimestamps is the list of timestamps within the signature.
	xmlTimestamps []*jaxb.XmlTimestamp
}

// NewBasicSignatureValidationProcess is the default constructor. Port of
// BasicSignatureValidationProcess(Provider, Data, SignatureWrapper, List, Map).
func NewBasicSignatureValidationProcess(i18nProvider *i18n.Provider, diagnosticData *diagnostic.Data,
	signatureWrapper *diagnostic.SignatureWrapper, xmlTimestamps []*jaxb.XmlTimestamp,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) *BasicSignatureValidationProcess {
	xmlResult := &jaxb.XmlValidationProcessBasicSignature{}
	result := process.NewResult(xmlResult, &xmlResult.XmlConstraintsConclusionWithProofOfExistenceContent.XmlConstraintsConclusionContent,
		&xmlResult.XmlConstraintsConclusionAttrs)
	c := &BasicSignatureValidationProcess{
		AbstractBasicValidationProcess: NewAbstractBasicValidationProcess(i18nProvider, result, diagnosticData, signatureWrapper, bbbs),
		xmlTimestamps:                  xmlTimestamps,
	}
	c.InitAbstractBasicValidationProcess(c)

	result.Value.ProofOfExistence = c.getCurrentTime(diagnosticData)

	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *BasicSignatureValidationProcess) Title() i18n.MessageTag {
	return i18n.MessageTagVPBS
}

// getCurrentTime ports the private getCurrentTime(DiagnosticData).
func (c *BasicSignatureValidationProcess) getCurrentTime(diagnosticData *diagnostic.Data) *jaxb.XmlProofOfExistence {
	proofOfExistence := &jaxb.XmlProofOfExistence{}
	if validationDate := diagnosticData.ValidationDate(); validationDate != nil {
		proofOfExistence.Time = jaxb.XSDateTime(*validationDate)
	}
	return proofOfExistence
}

// ContentTimestamps returns a list of content timestamps. Port of the
// overridden getContentTimestamps().
func (c *BasicSignatureValidationProcess) ContentTimestamps() []*diagnostic.TimestampWrapper {
	signature := c.DiagnosticData.SignatureById(c.Token.Id())
	if signature != nil {
		return signature.ContentTimestamps()
	}
	return nil
}

// TimestampValidation gets the corresponding validation result for a
// timestamp with the given Id. Port of the overridden
// getTimestampValidation(String).
func (c *BasicSignatureValidationProcess) TimestampValidation(timestampId string) *jaxb.XmlValidationProcessBasicTimestamp {
	for _, xmlTimestamp := range c.xmlTimestamps {
		if xmlTimestamp.Id != nil && utils.AreStringsEqual(timestampId, *xmlTimestamp.Id) {
			return xmlTimestamp.ValidationProcessBasicTimestamp
		}
	}
	return nil
}
