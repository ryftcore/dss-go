// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/SignatureTypeCheck.java (DSS 6.5.RC1).
package sav

import (
	"strings"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// signatureTypeApplicationPrefix is the RFC 7515 signature type prefix.
const signatureTypeApplicationPrefix = "application/"

// SignatureTypeCheck verifies whether the 'typ' (Type) protected header
// parameter has one of the expected values.
type SignatureTypeCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewSignatureTypeCheck is the default constructor. Port of
// SignatureTypeCheck(I18nProvider, XmlSAV, SignatureWrapper, MultiValuesRule).
func NewSignatureTypeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.MultiValuesRule) *SignatureTypeCheck {
	c := &SignatureTypeCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		signature:                    signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignatureTypeCheck) Process() bool {
	return c.ProcessValuesCheck(c.getSignatureType())
}

// getSignatureType ports the private getSignatureType().
func (c *SignatureTypeCheck) getSignatureType() []string {
	signatureTypes := make([]string, 0)
	if c.signature.SignatureType() != "" {
		signatureTypes = append(signatureTypes, c.signature.SignatureType())
	}
	signatureForm, err := c.signature.SignatureFormat().SignatureForm()
	if err != nil {
		panic(err)
	}
	if enumerations.SignatureForm_JAdES == signatureForm && c.signature.SignatureType() != "" {
		signatureTypes = append(signatureTypes, c.getRFC7515SignatureType(c.signature.SignatureType()))
	}
	return signatureTypes
}

// getRFC7515SignatureType ports the private getRFC7515SignatureType(String).
//
// RFC 7515 "4.1.9. "typ" (Type) Header Parameter": To keep messages compact in
// common situations, it is RECOMMENDED that producers omit an "application/"
// prefix of a media type value in a "typ" Header Parameter when no other '/'
// appears in the media type value. A recipient using the media type value MUST
// treat it as if "application/" were prepended to any "typ" value not
// containing a '/'.
func (c *SignatureTypeCheck) getRFC7515SignatureType(signatureType string) string {
	if signatureType == "" {
		return ""
	}
	shortMimeTypeString, err := spi.DSSUtilsStripFirstLeadingOccurrence(signatureType, signatureTypeApplicationPrefix)
	if err != nil {
		panic(err)
	}
	if !strings.Contains(shortMimeTypeString, "/") {
		return shortMimeTypeString
	}
	// return original if contains other '/'
	return signatureType
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignatureTypeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISQPSTYPP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignatureTypeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISQPSTYPP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignatureTypeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SignatureTypeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
