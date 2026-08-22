// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/checks/PastValidationAcceptableRevocationDataAvailable.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PastValidationAcceptableRevocationDataAvailable checks if an acceptable
// revocation data is present for a Past Signature Validation process.
type PastValidationAcceptableRevocationDataAvailable[T any] struct {
	*process.ChainItemBase[T]

	// revocationData is the revocation data to check.
	revocationData []*diagnostic.CertificateRevocationWrapper

	// revocationAcceptanceResults is the list of Revocation Acceptance Checks.
	revocationAcceptanceResults []*jaxb.XmlRAC
}

// NewPastValidationAcceptableRevocationDataAvailable is the constructor without
// certificate. Port of
// PastValidationAcceptableRevocationDataAvailable(I18nProvider, T, List, List, LevelRule).
func NewPastValidationAcceptableRevocationDataAvailable[T any](i18nProvider *i18n.I18nProvider,
	result *process.Result[T], revocationData []*diagnostic.CertificateRevocationWrapper,
	revocationAcceptanceResults []*jaxb.XmlRAC,
	constraint policy.LevelRule) *PastValidationAcceptableRevocationDataAvailable[T] {
	c := &PastValidationAcceptableRevocationDataAvailable[T]{
		ChainItemBase:               process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData:              revocationData,
		revocationAcceptanceResults: revocationAcceptanceResults,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *PastValidationAcceptableRevocationDataAvailable[T]) Process() bool {
	return utils.IsCollectionNotEmpty(c.revocationData)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PastValidationAcceptableRevocationDataAvailable[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IARDPFC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *PastValidationAcceptableRevocationDataAvailable[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IARDPFC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *PastValidationAcceptableRevocationDataAvailable[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *PastValidationAcceptableRevocationDataAvailable[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.isCryptoFailure() {
		return enumerations.SubIndication_CRYPTO_CONSTRAINTS_FAILURE_NO_POE
	}
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}

// isCryptoFailure ports the private isCryptoFailure().
func (c *PastValidationAcceptableRevocationDataAvailable[T]) isCryptoFailure() bool {
	for _, xmlRAC := range c.revocationAcceptanceResults {
		// XmlConclusion#getSubIndication(): the generated member is a pointer,
		// whose nil is Java's null.
		var subIndication enumerations.SubIndication
		if xmlRAC.Conclusion.SubIndication != nil {
			subIndication = xmlRAC.Conclusion.SubIndication.SubIndication()
		}
		if enumerations.Indication_INDETERMINATE == xmlRAC.Conclusion.Indication.Indication() &&
			enumerations.SubIndication_CRYPTO_CONSTRAINTS_FAILURE_NO_POE == subIndication {
			return true
		}
	}
	return false
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(), whose null result leaves the element absent.
//
// Java hands java.text.MessageFormat the List<String> itself, which renders it
// with AbstractCollection#toString: "[a, b]". The Go argument is that text.
func (c *PastValidationAcceptableRevocationDataAvailable[T]) BuildAdditionalInfo() *string {
	if utils.IsCollectionNotEmpty(c.revocationData) {
		revocationDataIds := make([]string, 0, len(c.revocationData))
		for _, revocation := range c.revocationData {
			revocationDataIds = append(revocationDataIds, revocation.Id())
		}
		message := c.I18nProvider.GetMessage(i18n.MessageTag_ACCEPTABLE_REVOCATION,
			"["+strings.Join(revocationDataIds, ", ")+"]")
		return &message
	}
	return nil
}
