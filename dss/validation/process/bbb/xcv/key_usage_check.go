// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/KeyUsageCheck.java (DSS 6.5.RC1).
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

// KeyUsageCheck checks if the certificate's key usage are acceptable.
type KeyUsageCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// context is the execution context (e.g. signature, timestamp, etc.).
	context enumerations.Context

	// subContext is the execution subContext (e.g. signing-certificate, CA certificate).
	subContext enumerations.SubContext
}

// NewKeyUsageCheck is the default constructor. Port of
// KeyUsageCheck(I18nProvider, XmlSubXCV, CertificateWrapper, Context, SubContext, MultiValuesRule).
func NewKeyUsageCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, context enumerations.Context, subContext enumerations.SubContext,
	constraint policy.MultiValuesRule) *KeyUsageCheck {
	c := &KeyUsageCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
		context:                      context,
		subContext:                   subContext,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *KeyUsageCheck) Process() bool {
	keyUsages := c.certificate.KeyUsages()
	var kubStrings []string
	for _, keyUsageBit := range keyUsages {
		kubStrings = append(kubStrings, keyUsageBit.Value())
	}
	return c.ProcessValuesCheck(kubStrings)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *KeyUsageCheck) BuildAdditionalInfo() *string {
	var names []string
	for _, keyUsageBit := range c.certificate.KeyUsages() {
		names = append(names, string(keyUsageBit))
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_KEY_USAGE, "["+strings.Join(names, ", ")+"]")
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *KeyUsageCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCGKU
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *KeyUsageCheck) BuildErrorMessage() *jaxb.XmlMessage {
	if enumerations.Context_CERTIFICATE == c.context {
		return c.BuildXmlMessage(i18n.MessageTag_BBB_XCV_ISCGKU_ANS_CERT)
	}
	position, err := process.GetSubContextPosition(c.context, c.subContext)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(i18n.MessageTag_BBB_XCV_ISCGKU_ANS, position)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *KeyUsageCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
//
// CA keyCertSign is a part of RFC 5280, while check of sign-cert falls under
// AdES validation process. Port of getFailedSubIndicationForConclusion().
func (c *KeyUsageCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if enumerations.SubContext_CA_CERTIFICATE == c.subContext {
		return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
	}
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}
