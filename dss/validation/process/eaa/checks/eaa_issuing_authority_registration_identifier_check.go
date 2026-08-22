// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAIssuingAuthorityRegistrationIdentifierCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// EAAIssuingAuthorityRegistrationIdentifierCheck verifies whether the EAA
// issuing authority registration identifier claim contains one of the
// expected values.
type EAAIssuingAuthorityRegistrationIdentifierCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAIssuingAuthorityRegistrationIdentifierCheck is the default constructor.
func NewEAAIssuingAuthorityRegistrationIdentifierCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.MultiValuesRule) *EAAIssuingAuthorityRegistrationIdentifierCheck {
	c := &EAAIssuingAuthorityRegistrationIdentifierCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		eaa:                          eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAIssuingAuthorityRegistrationIdentifierCheck) Process() bool {
	return c.ProcessValueCheck(c.eaa.IssuingRegistrationIdentifier())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAIssuingAuthorityRegistrationIdentifierCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAISSRegID
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAIssuingAuthorityRegistrationIdentifierCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAISSRegIDANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAIssuingAuthorityRegistrationIdentifierCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAIssuingAuthorityRegistrationIdentifierCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
