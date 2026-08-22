// Ported from dss-validation/.../validation/process/bbb/fc/checks/EllipticCurveKeySizeCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EllipticCurveKeySizeCheck verifies whether the elliptic curve key size used to create the
// signature corresponds to the value defined within the 'alg' header of the JWA signature, per RFC 7518.
type EllipticCurveKeySizeCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	signature *diagnostic.SignatureWrapper
}

// NewEllipticCurveKeySizeCheck is the default constructor.
func NewEllipticCurveKeySizeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *EllipticCurveKeySizeCheck {
	c := &EllipticCurveKeySizeCheck{signature: signature}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

func (c *EllipticCurveKeySizeCheck) isEncryptionAlgorithmKnown() bool {
	return c.signature.EncryptionAlgorithm() != ""
}

func (c *EllipticCurveKeySizeCheck) isDigestAlgorithmKnown() bool {
	return c.signature.DigestAlgorithm() != ""
}

func (c *EllipticCurveKeySizeCheck) isKeySizeKnown() bool {
	return c.signature.KeyLengthUsedToSignThisToken() != ""
}

func (c *EllipticCurveKeySizeCheck) isDigestAlgorithmAuthorized() bool {
	// Only three DigestAlgorithms are authorized to be used with ECDSA/PLAIN-ECDSA in RFC 7518
	switch c.signature.DigestAlgorithm() {
	case enumerations.DigestAlgorithmSHA256, enumerations.DigestAlgorithmSHA384, enumerations.DigestAlgorithmSHA512:
		return true
	default:
		return false
	}
}

func correspondingKeySize(digestAlgorithm enumerations.DigestAlgorithm) string {
	switch digestAlgorithm {
	case enumerations.DigestAlgorithmSHA256:
		return "256"
	case enumerations.DigestAlgorithmSHA384:
		return "384"
	case enumerations.DigestAlgorithmSHA512:
		return "521"
	default:
		return ""
	}
}

func (c *EllipticCurveKeySizeCheck) keySizeCorrespondsDigestAlgorithm() bool {
	expected := correspondingKeySize(c.signature.DigestAlgorithm())
	return expected != "" && expected == c.signature.KeyLengthUsedToSignThisToken()
}

// Process performs the check.
func (c *EllipticCurveKeySizeCheck) Process() bool {
	if !c.isEncryptionAlgorithmKnown() || !c.isDigestAlgorithmKnown() || !c.isKeySizeKnown() {
		return false
	}
	return !c.signature.EncryptionAlgorithm().IsEquivalent(enumerations.EncryptionAlgorithmECDSA) ||
		(c.isDigestAlgorithmAuthorized() && c.keySizeCorrespondsDigestAlgorithm())
}

// MessageTag returns the constraint message i18n key.
func (c *EllipticCurveKeySizeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCIECKSCDA
}

// ErrorMessageTag returns the error message i18n key.
func (c *EllipticCurveKeySizeCheck) ErrorMessageTag() i18n.MessageTag {
	switch {
	case !c.isEncryptionAlgorithmKnown():
		return i18n.MessageTagBBBFCIECKSCDAANS1
	case !c.isDigestAlgorithmKnown():
		return i18n.MessageTagBBBFCIECKSCDAANS2
	case !c.isKeySizeKnown():
		return i18n.MessageTagBBBFCIECKSCDAANS3
	case !c.isDigestAlgorithmAuthorized():
		return i18n.MessageTagBBBFCIECKSCDAANS4
	case !c.keySizeCorrespondsDigestAlgorithm():
		return i18n.MessageTagBBBFCIECKSCDAANS5
	default:
		return ""
	}
}

// BuildAdditionalInfo builds the additional info message. Port of the overridden
// protected String buildAdditionalInfo().
func (c *EllipticCurveKeySizeCheck) BuildAdditionalInfo() *string {
	if c.isEncryptionAlgorithmKnown() && c.isDigestAlgorithmKnown() && c.isKeySizeKnown() {
		message := c.I18nProvider.GetMessage(i18n.MessageTagSignatureAlgorithmWithKeySize,
			c.signature.SignatureAlgorithm().Name(), c.signature.KeyLengthUsedToSignThisToken())
		return &message
	}
	return nil
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *EllipticCurveKeySizeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *EllipticCurveKeySizeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
