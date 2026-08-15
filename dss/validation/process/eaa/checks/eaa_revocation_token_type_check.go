// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/status/EAARevocationTokenTypeCheck.java (DSS 6.5.RC1).
package checks

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

// eaaRevocationTokenTypeMimeTypeApplicationPrefix is the RFC 7519 type prefix.
const eaaRevocationTokenTypeMimeTypeApplicationPrefix = "application/"

// EAARevocationTokenTypeCheck verifies whether the type declared for the EAA
// revocation token is within an acceptable list of values.
type EAARevocationTokenTypeCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlFC]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper
}

// NewEAARevocationTokenTypeCheck is the default constructor.
func NewEAARevocationTokenTypeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlFC],
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper, constraint policy.MultiValuesRule) *EAARevocationTokenTypeCheck {
	c := &EAARevocationTokenTypeCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		eaaStatusToken:               eaaStatusToken,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAARevocationTokenTypeCheck) Process() bool {
	return c.ProcessValueCheck(c.getType())
}

// getType ports the private getType().
//
// TODO : Separate JWT and CWT logic ?
func (c *EAARevocationTokenTypeCheck) getType() string {
	return c.getRFC7519SignatureType(c.eaaStatusToken.Type())
}

// getRFC7519SignatureType ports the private getRFC7519SignatureType(String).
func (c *EAARevocationTokenTypeCheck) getRFC7519SignatureType(mimeType string) string {
	if mimeType == "" {
		return ""
	}
	shortMimeTypeString, err := spi.DSSUtilsStripFirstLeadingOccurrence(mimeType, eaaRevocationTokenTypeMimeTypeApplicationPrefix)
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
func (c *EAARevocationTokenTypeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_TYPE
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationTokenTypeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_TYPE_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationTokenTypeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationTokenTypeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
