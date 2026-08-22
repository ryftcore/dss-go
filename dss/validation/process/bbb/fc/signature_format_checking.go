// Ported from dss-validation/.../validation/process/bbb/fc/SignatureFormatChecking.java (DSS 6.5.RC1).
//
// 5.2.2 Format Checking. This building block checks that the signature to
// validate is conformant to the applicable base format (e.g. CMS, CAdES,
// XML-DSig, XAdES, etc.) prior to any subsequent processing.
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignatureFormatChecking performs "5.2.2 Format Checking" building block execution for a signature.
type SignatureFormatChecking struct {
	AbstractSignatureFormatChecking[*diagnostic.SignatureWrapper]
}

// NewSignatureFormatChecking is the default constructor.
func NewSignatureFormatChecking(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	signature *diagnostic.SignatureWrapper, context enumerations.Context, pol policy.ValidationPolicy) *SignatureFormatChecking {
	c := &SignatureFormatChecking{}
	c.InitAbstractSignatureFormatChecking(i18nProvider, diagnosticData, signature, context, pol, c)
	c.InitChainBase(c)
	return c
}

// InitChain builds the constraint chain. Port of the overridden protected void initChain().
func (c *SignatureFormatChecking) InitChain() {
	var item process.ChainItem[*drjaxb.XmlFC] = c.formatCheck()
	c.FirstItem = item

	item = item.SetNextItem(c.signatureDuplicateCheck())
	item = item.SetNextItem(c.referenceDuplicateCheck())
	item = item.SetNextItem(c.fullScopeCheck())

	// PAdES
	if c.Token.PDFRevision() != nil {
		item = item.SetNextItem(c.signerInformationStoreCheck())
		item = c.GetPDFRevisionValidationChain(item)
	}

	// PDF/A
	if c.DiagnosticData.IsPDFAValidationPerformed() {
		item = c.GetPdfaValidationChain(item)
	}

	// JAdES
	if signatureForm, err := c.Token.SignatureFormat().SignatureForm(); err == nil && signatureForm == enumerations.SignatureFormJAdES {
		if c.Token.EncryptionAlgorithm() != "" && c.Token.EncryptionAlgorithm().IsEquivalent(enumerations.EncryptionAlgorithmECDSA) {
			item = item.SetNextItem(c.ellipticCurveKeySizeCheck())
		}
	}

	// ASiC
	if c.DiagnosticData.IsContainerInfoPresent() {
		item = c.GetASiCContainerValidationChain(item)
		item = item.SetNextItem(c.allFilesSignedCheck()) //nolint:staticcheck // mirrors upstream SignatureFormatChecking#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
	}
}

func (c *SignatureFormatChecking) formatCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.SignatureFormatConstraint(c.Context)
	return NewFormatCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *SignatureFormatChecking) signatureDuplicateCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.SignatureDuplicatedConstraint(c.Context)
	return NewSignatureNotAmbiguousCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *SignatureFormatChecking) referenceDuplicateCheck() process.ChainItem[*drjaxb.XmlFC] {
	return NewReferencesNotAmbiguousCheck(c.I18nProvider, c.Result, c.Token, c.FailLevelRule())
}

func (c *SignatureFormatChecking) fullScopeCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.FullScopeConstraint()
	return NewFullScopeCheck(c.I18nProvider, c.Result, c.Token.SignatureScopes(), constraint)
}

func (c *SignatureFormatChecking) signerInformationStoreCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.SignerInformationStoreConstraint(c.Context)
	return NewSignerInformationStoreCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *SignatureFormatChecking) ellipticCurveKeySizeCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EllipticCurveKeySizeConstraint(c.Context)
	return NewEllipticCurveKeySizeCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

// FilenameAdherenceCheck creates the signature filename adherence check. Port of the overridden
// protected ChainItem<XmlFC> filenameAdherenceCheck().
func (c *SignatureFormatChecking) FilenameAdherenceCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.FilenameAdherenceConstraint()
	return NewSignatureFilenameAdherenceCheck(c.I18nProvider, c.Result, c.DiagnosticData, c.Token, constraint)
}

// ManifestFilenameAdherenceCheck creates the signature manifest filename adherence check. Port
// of the overridden protected ChainItem<XmlFC> manifestFilenameAdherenceCheck().
func (c *SignatureFormatChecking) ManifestFilenameAdherenceCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.FilenameAdherenceConstraint()
	return NewSignatureManifestFilenameAdherenceCheck(c.I18nProvider, c.Result, c.DiagnosticData, c.Token, constraint)
}

func (c *SignatureFormatChecking) allFilesSignedCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.AllFilesSignedConstraint()
	return NewAllFilesSignedCheck(c.I18nProvider, c.Result, c.Token, c.DiagnosticData.ContainerInfo(), constraint)
}
