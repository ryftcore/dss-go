// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/ExtendedKeyUsageCheck.java (DSS 6.5.RC1).
package xcv

import (
	"strings"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
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
// ExtendedKeyUsageCheck(Provider, XmlSubXCV, CertificateWrapper, Context, SubContext, MultiValuesRule).
func NewExtendedKeyUsageCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
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
// buildAdditionalInfo(): Arrays.toString(Object[]) prints a null XmlOID
// description as the literal text "null" (not empty), which the []string
// extendedKeyUsageDescriptions() feeds to Process()/ValuesCheck
// cannot represent, so this walks the certificate's ExtendedKeyUsages() again
// to render each entry the way Java's Object[] would.
func (c *ExtendedKeyUsageCheck) BuildAdditionalInfo() *string {
	ekus := c.certificate.ExtendedKeyUsages()
	rendered := make([]string, 0, len(ekus))
	for _, eku := range ekus {
		if eku.Description != nil {
			rendered = append(rendered, *eku.Description)
		} else {
			rendered = append(rendered, "null")
		}
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagExtendedKeyUsage,
		"["+strings.Join(rendered, ", ")+"]")
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ExtendedKeyUsageCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVISCGEKU
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *ExtendedKeyUsageCheck) BuildErrorMessage() *jaxb.XmlMessage {
	if enumerations.ContextCertificate == c.context {
		return c.BuildXmlMessage(i18n.MessageTagBBBXCVISCGEKUANSCert)
	}
	position, err := process.GetSubContextPosition(c.context, c.subContext)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(i18n.MessageTagBBBXCVISCGEKUANS, position)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ExtendedKeyUsageCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *ExtendedKeyUsageCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
