// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/ExtendedKeyUsageCheck.java (DSS 6.5.RC1).
package xcv

import (
	"strings"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// ExtendedKeyUsageCheck checks if the extended key usage is acceptable.
type ExtendedKeyUsageCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// context is the execution context (e.g. signature, timestamp, etc.).
	context enumerations.Context

	// subContext is the execution subContext (e.g. signing-certificate, CA certificate).
	subContext enumerations.SubContext
}

// NewExtendedKeyUsageCheck is the default constructor. Port of
// ExtendedKeyUsageCheck(I18nProvider, XmlSubXCV, CertificateWrapper, Context, SubContext, MultiValuesRule).
func NewExtendedKeyUsageCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, context enumerations.Context, subContext enumerations.SubContext,
	constraint policy.MultiValuesRule) *ExtendedKeyUsageCheck {
	c := &ExtendedKeyUsageCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
		context:                      context,
		subContext:                   subContext,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ExtendedKeyUsageCheck) Process() bool {
	return c.ProcessValuesCheck(c.extendedKeyUsageDescriptions())
}

// extendedKeyUsageDescriptions ports the private getExtendedKeyUsageDescriptions().
func (c *ExtendedKeyUsageCheck) extendedKeyUsageDescriptions() []string {
	var result []string
	for _, eku := range c.certificate.ExtendedKeyUsages() {
		if eku.Description != nil {
			result = append(result, *eku.Description)
		} else {
			result = append(result, "")
		}
	}
	return result
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *ExtendedKeyUsageCheck) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTag_EXTENDED_KEY_USAGE,
		javaArrayToString(c.extendedKeyUsageDescriptions()))
	return &message
}

// javaArrayToString ports java.util.Arrays#toString(Object[]).
func javaArrayToString(values []string) string {
	return "[" + strings.Join(values, ", ") + "]"
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ExtendedKeyUsageCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCGEKU
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *ExtendedKeyUsageCheck) BuildErrorMessage() *jaxb.XmlMessage {
	if enumerations.Context_CERTIFICATE == c.context {
		return c.BuildXmlMessage(i18n.MessageTag_BBB_XCV_ISCGEKU_ANS_CERT)
	}
	position, err := process.GetSubContextPosition(c.context, c.subContext)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(i18n.MessageTag_BBB_XCV_ISCGEKU_ANS, position)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ExtendedKeyUsageCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *ExtendedKeyUsageCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}
