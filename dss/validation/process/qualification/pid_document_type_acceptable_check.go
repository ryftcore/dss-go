// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/pid/checks/PIDDocumentTypeAcceptableCheck.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"
	"strings"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// PIDDocumentTypeAcceptableCheck verifies whether the PID contains person
// identification data defined within an eIDAS allowed namespace.
type PIDDocumentTypeAcceptableCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationPIDQualificationProcess]

	// eaa is the EAA presentation to be checked.
	eaa *diagnostic.EAAWrapper
}

// NewPIDDocumentTypeAcceptableCheck is the default constructor. Port of
// PIDDocumentTypeAcceptableCheck(I18nProvider, XmlValidationPIDQualificationProcess, EAAWrapper, LevelRule).
func NewPIDDocumentTypeAcceptableCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationPIDQualificationProcess],
	eaa *diagnostic.EAAWrapper, constraint policy.LevelRule) *PIDDocumentTypeAcceptableCheck {
	c := &PIDDocumentTypeAcceptableCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaa,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *PIDDocumentTypeAcceptableCheck) Process() bool {
	documentType := c.getClaimedDocumentType()
	if documentType == "" {
		return false
	}
	switch c.eaa.EAAType() {
	case enumerations.EAAType_SD_JWT_VC:
		return strings.HasPrefix(documentType, "urn:eudi:pid:")
	case enumerations.EAAType_ISO_IEC_MDOC:
		// TODO : not clear what element is to be checked
		// The attestation type for person identification data in ISO/IEC mdoc format
		// shall be "eu.europa.ec.eudi.pid.1".
		return documentType == "eu.europa.ec.eudi.pid.1"
	default:
		panic(fmt.Sprintf("Not supported EAA Type : '%s'", c.eaa.EAAType()))
	}
}

// getClaimedDocumentType ports the private getClaimedDocumentType().
func (c *PIDDocumentTypeAcceptableCheck) getClaimedDocumentType() string {
	switch c.eaa.EAAType() {
	case enumerations.EAAType_SD_JWT_VC:
		return c.eaa.EAAVerifiableCredentialsTypeUri()
	case enumerations.EAAType_ISO_IEC_MDOC:
		// TODO : not clear what element is to be checked
		// The attestation type for person identification data in ISO/IEC mdoc format
		// shall be "eu.europa.ec.eudi.pid.1".
		return c.eaa.EAADocumentType()
	default:
		panic(fmt.Sprintf("Not supported EAA Type : '%s'", c.eaa.EAAType()))
	}
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PIDDocumentTypeAcceptableCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PID_DOCUMENT_TYPE
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *PIDDocumentTypeAcceptableCheck) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTag_PID_DOCUMENT_TYPE_ANS, c.getClaimedDocumentType())
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *PIDDocumentTypeAcceptableCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *PIDDocumentTypeAcceptableCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
