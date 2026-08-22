// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/ContentTypeCheck.java (DSS 6.5.RC1).
package sav

import (
	"strings"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// mimeTypeApplicationPrefix is the RFC 7515 content type prefix.
const mimeTypeApplicationPrefix = "application/"

// ContentTypeCheck checks if the content type is acceptable.
type ContentTypeCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewContentTypeCheck is the default constructor. Port of
// ContentTypeCheck(I18nProvider, XmlSAV, SignatureWrapper, MultiValuesRule).
func NewContentTypeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.MultiValuesRule) *ContentTypeCheck {
	c := &ContentTypeCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		signature:                    signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ContentTypeCheck) Process() bool {
	return c.ProcessValuesCheck(c.getContentType())
}

// getContentType ports the private getContentType().
func (c *ContentTypeCheck) getContentType() []string {
	contentTypes := make([]string, 0)
	if c.signature.ContentType() != "" {
		contentTypes = append(contentTypes, c.signature.ContentType())
	}
	// NOTE: 'cty' (content type) signed header is used to define signed content's MimeType (see 102-2).
	// The check is merged for simplicity of the users.
	signatureForm, err := c.signature.SignatureFormat().SignatureForm()
	if err != nil {
		panic(err)
	}
	if enumerations.SignatureFormJAdES == signatureForm && c.signature.MimeType() != "" {
		contentTypes = append(contentTypes, c.getRFC7515ContentType(c.signature.MimeType()))
	}
	return contentTypes
}

// getRFC7515ContentType ports the private getRFC7515ContentType(String).
//
// RFC 7515 requires to handle the ContentType, as it has an omitted
// "application/" prefix. Therefore, we add additional value to ensure its
// correct processing.
func (c *ContentTypeCheck) getRFC7515ContentType(mimeType string) string {
	if mimeType == "" {
		return ""
	}
	shortMimeTypeString, err := spi.DSSUtilsStripFirstLeadingOccurrence(mimeType, mimeTypeApplicationPrefix)
	if err != nil {
		panic(err)
	}
	if !strings.Contains(shortMimeTypeString, "/") {
		return shortMimeTypeString
	}
	// return original if contains other '/'
	return mimeType
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ContentTypeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISQPCTP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ContentTypeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISQPCTP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ContentTypeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ContentTypeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}
