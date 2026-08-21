// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/ListTypeKnownCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ListTypeKnownCheck verifies whether the list type is known.
type ListTypeKnownCheck struct {
	*process.ChainItemBase[*jaxb.XmlCertificateApprovalStatusProcess]

	// listTypeUri is the List Type URI.
	listTypeUri string
}

// NewListTypeKnownCheck is the default constructor. Port of
// ListTypeKnownCheck(I18nProvider, XmlCertificateApprovalStatusProcess, String, LevelRule).
func NewListTypeKnownCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlCertificateApprovalStatusProcess],
	listTypeUri string, constraint policy.LevelRule) *ListTypeKnownCheck {
	c := &ListTypeKnownCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		listTypeUri:   listTypeUri,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ListTypeKnownCheck) Process() bool {
	listType := enumerations.ListTypeFromURI(c.listTypeUri)
	return listType != nil && listType.Label() != "" // Label is present -> defined
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *ListTypeKnownCheck) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTag_CERTIFICATE_USAGE_LIST_TYPE, c.listTypeUri)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ListTypeKnownCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_LIST_TYPE_KNOWN
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *ListTypeKnownCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_LIST_TYPE_KNOWN_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ListTypeKnownCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *ListTypeKnownCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
