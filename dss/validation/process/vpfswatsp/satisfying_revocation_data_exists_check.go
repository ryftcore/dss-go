// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/vts/checks/SatisfyingRevocationDataExistsCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
//
// The Java class extends bbb.xcv.sub.checks.CertificateRevocationSelectorResultCheck
// and reads that base's protected `crsResult` field in buildAdditionalInfo. The
// Go base keeps the field unexported, so the same XmlCRS the constructor
// already receives is held here too - one pointer, two references, no copy.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/xcv"
)

// SatisfyingRevocationDataExistsCheck checks if an acceptable revocation data
// exists.
type SatisfyingRevocationDataExistsCheck[T any] struct {
	*xcv.CertificateRevocationSelectorResultCheck[T]

	// crsResult is the CRS result; see the file header.
	crsResult *jaxb.XmlCRS

	// certificateWrapper is the concerned certificate token.
	certificateWrapper *diagnostic.CertificateWrapper

	// controlTime is the control time used to find out the revocation data.
	controlTime time.Time
}

// NewSatisfyingRevocationDataExistsCheck is the default constructor. Port of
// SatisfyingRevocationDataExistsCheck(Provider, T, XmlCRS, CertificateWrapper, Date, LevelRule).
func NewSatisfyingRevocationDataExistsCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	crsResult *jaxb.XmlCRS, certificateWrapper *diagnostic.CertificateWrapper, controlTime time.Time,
	constraint policy.LevelRule) *SatisfyingRevocationDataExistsCheck[T] {
	c := &SatisfyingRevocationDataExistsCheck[T]{
		CertificateRevocationSelectorResultCheck: xcv.NewCertificateRevocationSelectorResultCheck(
			i18nProvider, result, crsResult, constraint),
		crsResult:          crsResult,
		certificateWrapper: certificateWrapper,
		controlTime:        controlTime,
	}
	// Re-register with the outer type so the overridden methods below dispatch
	// correctly.
	c.InitChainItem(c)
	return c
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo().
func (c *SatisfyingRevocationDataExistsCheck[T]) BuildAdditionalInfo() *string {
	latestAcceptableRevocationId := c.crsResult.LatestAcceptableRevocationId
	var message string
	if latestAcceptableRevocationId != nil {
		message = c.I18nProvider.GetMessage(i18n.MessageTagCertificateRevocationFound, *latestAcceptableRevocationId,
			c.certificateWrapper.Id(), process.GetFormattedDate(&c.controlTime))
	} else {
		message = c.I18nProvider.GetMessage(i18n.MessageTagCertificateRevocationNotFound, c.certificateWrapper.Id(),
			process.GetFormattedDate(&c.controlTime))
	}
	return &message
}

// MessageTag returns the check's message tag. Port of the overridden
// getMessageTag().
func (c *SatisfyingRevocationDataExistsCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBVTSIRDPFC
}

// ErrorMessageTag returns the check's error message tag. Port of the overridden
// getErrorMessageTag().
func (c *SatisfyingRevocationDataExistsCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBVTSIRDPFCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// the overridden getFailedIndicationForConclusion().
func (c *SatisfyingRevocationDataExistsCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the overridden getFailedSubIndicationForConclusion().
func (c *SatisfyingRevocationDataExistsCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoPOE
}
