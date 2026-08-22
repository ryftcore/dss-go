// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAATypeCheck.java (DSS 6.5.RC1).
package checks

import (
	"fmt"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// EAATypeCheck verifies whether the EAA contains an acceptable EAA type.
type EAATypeCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAATypeCheck is the default constructor.
func NewEAATypeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.MultiValuesRule) *EAATypeCheck {
	c := &EAATypeCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		eaa:                          eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAATypeCheck) Process() bool {
	switch c.eaa.EAAType() {
	case enumerations.EAATypeSDJWTVC:
		return c.ProcessValueCheck(c.eaa.EAAVerifiableCredentialsTypeUri())
	case enumerations.EAATypeISOIECMDoc:
		docType := c.eaa.EAADocumentType()
		if docType == "" {
			// Handle IssuerSigned token.
			docType = c.eaa.DocumentType()
		}
		return c.ProcessValueCheck(docType)
	default:
		panic(fmt.Sprintf("The EAA Type '%s' is not supported!", c.eaa.EAAType()))
	}
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAATypeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_ACCEPTABLE_TYPE
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAATypeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_ACCEPTABLE_TYPE_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAATypeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAATypeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
