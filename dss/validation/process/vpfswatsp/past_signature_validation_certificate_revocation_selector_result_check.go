// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/checks/PastSignatureValidationCertificateRevocationSelectorResultCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note, and
// satisfying_revocation_data_exists_check.go for why the XmlCRS the base already
// holds is kept here as well.
package vpfswatsp

import (
	"strings"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb/xcv"
)

// PastSignatureValidationCertificateRevocationSelectorResultCheck verifies the
// validation result of a PastSignatureValidationCertificateRevocationSelector.
type PastSignatureValidationCertificateRevocationSelectorResultCheck struct {
	*xcv.CertificateRevocationSelectorResultCheck[*jaxb.XmlPSV]

	// crsResult is the CRS result; see the file header.
	crsResult *jaxb.XmlCRS
}

// NewPastSignatureValidationCertificateRevocationSelectorResultCheck is the
// default constructor. Port of
// PastSignatureValidationCertificateRevocationSelectorResultCheck(I18nProvider, XmlPSV, XmlCRS, LevelRule).
func NewPastSignatureValidationCertificateRevocationSelectorResultCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlPSV], crsResult *jaxb.XmlCRS,
	constraint policy.LevelRule) *PastSignatureValidationCertificateRevocationSelectorResultCheck {
	c := &PastSignatureValidationCertificateRevocationSelectorResultCheck{
		CertificateRevocationSelectorResultCheck: xcv.NewCertificateRevocationSelectorResultCheck(
			i18nProvider, result, crsResult, constraint),
		crsResult: crsResult,
	}
	// Re-register with the outer type so the overridden methods below dispatch
	// correctly.
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of the overridden
// getBlockType().
func (c *PastSignatureValidationCertificateRevocationSelectorResultCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockType_PSV_CRS
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo(), whose null result leaves the element absent.
//
// Java hands java.text.MessageFormat the List<String> itself, which renders it
// with AbstractCollection#toString: "[a, b]". The Go argument is that text.
func (c *PastSignatureValidationCertificateRevocationSelectorResultCheck) BuildAdditionalInfo() *string {
	acceptableRevocationId := acceptableRevocationIds(c.crsResult)
	if utils.IsCollectionNotEmpty(acceptableRevocationId) {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_ACCEPTABLE_REVOCATION,
			"["+strings.Join(acceptableRevocationId, ", ")+"]")
		return &message
	}
	return nil
}

// acceptableRevocationIds reads XmlCRS#getAcceptableRevocationId(): the
// generated member is a *StringList, whose nil is the empty list JAXB's getter
// would have created.
func acceptableRevocationIds(crsResult *jaxb.XmlCRS) []string {
	if crsResult.AcceptableRevocationId == nil {
		return nil
	}
	return *crsResult.AcceptableRevocationId
}
