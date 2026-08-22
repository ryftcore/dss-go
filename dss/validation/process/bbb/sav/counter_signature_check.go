// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/CounterSignatureCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CounterSignatureCheck checks if a counter signature is present for the
// signature.
type CounterSignatureCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.DiagnosticData

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewCounterSignatureCheck is the default constructor. Port of
// CounterSignatureCheck(I18nProvider, XmlSAV, DiagnosticData, SignatureWrapper, LevelRule).
func NewCounterSignatureCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	diagnosticData *diagnostic.DiagnosticData, signature *diagnostic.SignatureWrapper,
	constraint policy.LevelRule) *CounterSignatureCheck {
	c := &CounterSignatureCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		diagnosticData: diagnosticData,
		signature:      signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CounterSignatureCheck) Process() bool {
	foundCountersignature := false
	currentSignatureId := c.signature.Id()

	signatures := c.diagnosticData.Signatures()
	for _, signatureWrapper := range signatures {
		if signatureWrapper.IsCounterSignature() && currentSignatureId == signatureWrapper.Parent().Id() {
			foundCountersignature = true
			break
		}
	}

	return foundCountersignature
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CounterSignatureCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_IUQPCSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *CounterSignatureCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_IUQPCSP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CounterSignatureCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *CounterSignatureCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}
