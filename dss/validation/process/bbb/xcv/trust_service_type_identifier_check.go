// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/checks/TrustServiceTypeIdentifierCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// TrustServiceTypeIdentifierCheck checks if the certificate's usage time is
// in the validity range of a TrustService with the accepted type.
type TrustServiceTypeIdentifierCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// usageTime is the timestamp / revocation production time; nil is Java's
	// null (see the package-level date-mapping note in x509_certificate_validation.go).
	usageTime *time.Time

	// context is the validation context.
	context enumerations.Context

	// serviceTypeStr is the service type string.
	serviceTypeStr string
}

// NewTrustServiceTypeIdentifierCheck is the default constructor. Port of
// TrustServiceTypeIdentifierCheck(I18nProvider, XmlXCV, CertificateWrapper, Date, Context, MultiValuesRule).
func NewTrustServiceTypeIdentifierCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlXCV],
	certificate *diagnostic.CertificateWrapper, usageTime *time.Time, context enumerations.Context,
	constraint policy.MultiValuesRule) *TrustServiceTypeIdentifierCheck {
	c := &TrustServiceTypeIdentifierCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
		usageTime:                    usageTime,
		context:                      context,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustServiceTypeIdentifierCheck) Process() bool {
	// do not include Trusted list
	if c.certificate.IsCertificateChainFromTrustedStore() {
		return true
	}

	trustServices := c.certificate.TrustServices()
	for _, trustService := range trustServices {
		c.serviceTypeStr = utils.Trim(trustService.Type)
		statusStartDate := trustService.StartDate
		if c.ProcessValueCheck(c.serviceTypeStr) && statusStartDate != nil {
			statusEndDate := trustService.EndDate
			// The issuing time of the certificate should be into the validity period of the associated service
			if !c.usageTime.Before(*statusStartDate) && (statusEndDate == nil || c.usageTime.Before(*statusEndDate)) {
				return true
			}
		}
	}
	return false
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *TrustServiceTypeIdentifierCheck) BuildAdditionalInfo() *string {
	if utils.IsStringNotEmpty(c.serviceTypeStr) {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_TRUSTED_SERVICE_TYPE, c.serviceTypeStr)
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustServiceTypeIdentifierCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_XCV_TSL_ETIP
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustServiceTypeIdentifierCheck) ErrorMessageTag() i18n.MessageTag {
	switch c.context {
	case enumerations.Context_SIGNATURE, enumerations.Context_COUNTER_SIGNATURE, enumerations.Context_KEY_BINDING_SIGNATURE:
		return i18n.MessageTag_XCV_TSL_ETIP_SIG_ANS
	case enumerations.Context_TIMESTAMP:
		return i18n.MessageTag_XCV_TSL_ETIP_TSP_ANS
	case enumerations.Context_REVOCATION:
		return i18n.MessageTag_XCV_TSL_ETIP_REV_ANS
	default:
		return i18n.MessageTag_XCV_TSL_ETIP_ANS
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustServiceTypeIdentifierCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *TrustServiceTypeIdentifierCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_NO_CERTIFICATE_CHAIN_FOUND
}
