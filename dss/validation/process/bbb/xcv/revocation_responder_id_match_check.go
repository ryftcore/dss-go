// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationResponderIdMatchCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationResponderIdMatchCheck verifies whether the ResponderId property of
// an OCSP response matches the found certificate used to sign the OCSP response.
type RevocationResponderIdMatchCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper
}

// NewRevocationResponderIdMatchCheck is the default constructor. Port of
// RevocationResponderIdMatchCheck(I18nProvider, XmlRAC, RevocationWrapper, LevelRule).
func NewRevocationResponderIdMatchCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlRAC],
	revocationData *diagnostic.RevocationWrapper, constraint policy.LevelRule) *RevocationResponderIdMatchCheck {
	c := &RevocationResponderIdMatchCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationResponderIdMatchCheck) Process() bool {
	// ResponderId is returned as a SIGNING_CERTIFICATE for OCSP token
	relatedSigningCertificates := c.revocationData.FoundCertificates().
		RelatedCertificatesByRefOrigin(enumerations.CertificateRefOriginSigningCertificate)
	signingCertificate := c.revocationData.SigningCertificate()
	if utils.CollectionSize(relatedSigningCertificates) == 1 && signingCertificate != nil {
		return relatedSigningCertificates[0].Id() == signingCertificate.Id()
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationResponderIdMatchCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVRevocRespIDMatch
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationResponderIdMatchCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVRevocRespIDMatchANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationResponderIdMatchCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *RevocationResponderIdMatchCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
