// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/MRACertificateEquivalenceApplied.java (DSS 6.5.RC1).
package qualification

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// MRACertificateEquivalenceApplied verifies whether the certificate content
// equivalence information has been applied for the certificate.
type MRACertificateEquivalenceApplied[T any] struct {
	*process.ChainItemBase[T]

	// certificateWrapper is the certificate to be verified.
	certificateWrapper *diagnostic.CertificateWrapper
}

// NewMRACertificateEquivalenceApplied is the default constructor. Port of
// MRACertificateEquivalenceApplied(I18nProvider, T, CertificateWrapper, LevelRule).
func NewMRACertificateEquivalenceApplied[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificateWrapper *diagnostic.CertificateWrapper, constraint policy.LevelRule) *MRACertificateEquivalenceApplied[T] {
	c := &MRACertificateEquivalenceApplied[T]{
		ChainItemBase:      process.NewChainItemBase(i18nProvider, result, constraint),
		certificateWrapper: certificateWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *MRACertificateEquivalenceApplied[T]) Process() bool {
	if !c.certificateWrapper.IsEnactedMRA() {
		return false
	}
	for _, certificateContentEquivalence := range c.certificateWrapper.MRACertificateContentEquivalenceList() {
		if !certificateContentEquivalence.Enacted {
			return false
		}
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *MRACertificateEquivalenceApplied[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_HAS_METS_HCCECBA
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *MRACertificateEquivalenceApplied[T]) BuildErrorMessage() *jaxb.XmlMessage {
	if !c.certificateWrapper.IsEnactedMRA() {
		return c.BuildXmlMessage(i18n.MessageTag_QUAL_HAS_METS_HCCECBA_ANS)
	}
	uriList := c.getFailedCertificateEquivalenceContextUris()
	errorTag := i18n.MessageTag_QUAL_HAS_METS_HCCECBA_ANS_2
	var argument string
	if len(uriList) == 1 {
		argument = uriList[0]
	} else {
		errorTag = i18n.MessageTag_QUAL_HAS_METS_HCCECBA_ANS_3
		argument = "[" + strings.Join(uriList, ", ") + "]"
	}
	return c.BuildXmlMessage(errorTag, argument)
}

// getFailedCertificateEquivalenceContextUris ports the private
// getFailedCertificateEquivalenceContextUris().
func (c *MRACertificateEquivalenceApplied[T]) getFailedCertificateEquivalenceContextUris() []string {
	var result []string
	for _, certificateContentEquivalence := range c.certificateWrapper.MRACertificateContentEquivalenceList() {
		if !certificateContentEquivalence.Enacted {
			if certificateContentEquivalence.Uri != nil {
				result = append(result, *certificateContentEquivalence.Uri)
			} else {
				result = append(result, "")
			}
		}
	}
	return result
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *MRACertificateEquivalenceApplied[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *MRACertificateEquivalenceApplied[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
