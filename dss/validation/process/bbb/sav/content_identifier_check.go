// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/ContentIdentifierCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// ContentIdentifierCheck checks if the content identifier is acceptable.
type ContentIdentifierCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewContentIdentifierCheck is the default constructor. Port of
// ContentIdentifierCheck(Provider, XmlSAV, SignatureWrapper, MultiValuesRule).
func NewContentIdentifierCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.MultiValuesRule) *ContentIdentifierCheck {
	c := &ContentIdentifierCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		signature:                    signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ContentIdentifierCheck) Process() bool {
	contentIdentifier := c.signature.ContentIdentifier()
	return c.ProcessValueCheck(contentIdentifier)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ContentIdentifierCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVISQPCIP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ContentIdentifierCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVISQPCIPANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ContentIdentifierCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ContentIdentifierCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}
