// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingSignatureSerializationTypeCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLSCertificateBindingSignatureSerializationTypeCheck verifies whether the
// serialization type of the obtained signature is allowed.
type TLSCertificateBindingSignatureSerializationTypeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// signature is the TLS Certificate Binding signature.
	signature *diagnostic.SignatureWrapper
}

// NewTLSCertificateBindingSignatureSerializationTypeCheck is the default
// constructor. Port of
// TLSCertificateBindingSignatureSerializationTypeCheck(Provider, XmlValidationQWACProcess, SignatureWrapper, LevelRule).
func NewTLSCertificateBindingSignatureSerializationTypeCheck(i18nProvider *i18n.Provider,
	result *process.Result[*jaxb.XmlValidationQWACProcess], signature *diagnostic.SignatureWrapper,
	constraint policy.LevelRule) *TLSCertificateBindingSignatureSerializationTypeCheck {
	c := &TLSCertificateBindingSignatureSerializationTypeCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLSCertificateBindingSignatureSerializationTypeCheck) Process() bool {
	return enumerations.JWSSerializationTypeCompactSerialization == c.signature.JWSSerializationType()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLSCertificateBindingSignatureSerializationTypeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingSigSer
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLSCertificateBindingSignatureSerializationTypeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingSigSerANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLSCertificateBindingSignatureSerializationTypeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TLSCertificateBindingSignatureSerializationTypeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
