// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/vci/checks/SignaturePolicyIdentifierCheck.java (DSS 6.5.RC1).
package vci

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// SignaturePolicyIdentifierCheck checks if the signature policy identifier is
// acceptable.
type SignaturePolicyIdentifierCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlVCI]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper

	// multiValues is the constraint.
	multiValues policy.MultiValuesRule
}

// NewSignaturePolicyIdentifierCheck is the default constructor. Port of
// SignaturePolicyIdentifierCheck(I18nProvider, XmlVCI, SignatureWrapper, MultiValuesRule).
func NewSignaturePolicyIdentifierCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlVCI],
	signature *diagnostic.SignatureWrapper, multiValues policy.MultiValuesRule) *SignaturePolicyIdentifierCheck {
	c := &SignaturePolicyIdentifierCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, multiValues),
		signature:                    signature,
		multiValues:                  multiValues,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignaturePolicyIdentifierCheck) Process() bool {
	policyId := c.signature.PolicyId()
	if containsValue(c.multiValues.Values(), string(enumerations.SignaturePolicyTypeNoPolicy)) &&
		utils.IsStringEmpty(policyId) {
		return true
	} else if containsValue(c.multiValues.Values(), string(enumerations.SignaturePolicyTypeAnyPolicy)) &&
		utils.IsStringNotEmpty(policyId) {
		return true
	} else if containsValue(c.multiValues.Values(), string(enumerations.SignaturePolicyTypeImplicitPolicy)) &&
		utils.AreStringsEqual(string(enumerations.SignaturePolicyTypeImplicitPolicy), policyId) {
		return true
	}
	// oids
	return c.ProcessValueCheck(policyId)
}

// containsValue ports java.util.List#contains(Object) over the constraint values.
func containsValue(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignaturePolicyIdentifierCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBVCIISPK
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignaturePolicyIdentifierCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBVCIISPKANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignaturePolicyIdentifierCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SignaturePolicyIdentifierCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationPolicyProcessingError
}
