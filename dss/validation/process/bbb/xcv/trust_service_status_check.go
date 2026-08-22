// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/checks/TrustServiceStatusCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// TrustServiceStatusCheck checks if the certificate's usage time is in the
// validity range of a TrustService with the accepted status.
type TrustServiceStatusCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// usageTime is the timestamp / revocation production time; nil is Java's
	// null (see the package-level date-mapping note in x509_certificate_validation.go).
	usageTime *time.Time

	// context is the validation context.
	context enumerations.Context

	// serviceStatusStr is the service status string.
	serviceStatusStr string
}

// NewTrustServiceStatusCheck is the default constructor. Port of
// TrustServiceStatusCheck(I18nProvider, XmlXCV, CertificateWrapper, Date, Context, MultiValuesRule).
func NewTrustServiceStatusCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlXCV],
	certificate *diagnostic.CertificateWrapper, usageTime *time.Time, context enumerations.Context,
	constraint policy.MultiValuesRule) *TrustServiceStatusCheck {
	c := &TrustServiceStatusCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
		usageTime:                    usageTime,
		context:                      context,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustServiceStatusCheck) Process() bool {
	// do not include Trusted list
	if c.certificate.IsCertificateChainFromTrustedStore() {
		return true
	}

	trustServices := c.certificate.TrustServices()
	if utils.IsCollectionNotEmpty(trustServices) {
		for _, trustService := range trustServices {
			c.serviceStatusStr = utils.Trim(trustService.Status)
			statusStartDate := trustService.StartDate
			if c.ProcessValueCheck(c.serviceStatusStr) && statusStartDate != nil {
				statusEndDate := trustService.EndDate
				// The issuing time of the certificate should be into the validity period of the associated service
				if !c.usageTime.Before(*statusStartDate) && (statusEndDate == nil || c.usageTime.Before(*statusEndDate)) {
					return true
				}
			}
		}
	}

	return false
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *TrustServiceStatusCheck) BuildAdditionalInfo() *string {
	if utils.IsStringNotEmpty(c.serviceStatusStr) {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_TRUSTED_SERVICE_STATUS, c.serviceStatusStr)
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustServiceStatusCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_XCV_TSL_ESP
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustServiceStatusCheck) ErrorMessageTag() i18n.MessageTag {
	switch c.context {
	case enumerations.ContextSignature, enumerations.ContextCounterSignature, enumerations.ContextKeyBindingSignature:
		return i18n.MessageTag_XCV_TSL_ESP_SIG_ANS
	case enumerations.ContextTimestamp:
		return i18n.MessageTag_XCV_TSL_ESP_TSP_ANS
	case enumerations.ContextRevocation:
		return i18n.MessageTag_XCV_TSL_ESP_REV_ANS
	default:
		return i18n.MessageTag_XCV_TSL_ESP_ANS
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustServiceStatusCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *TrustServiceStatusCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoCertificateChainFound
}
