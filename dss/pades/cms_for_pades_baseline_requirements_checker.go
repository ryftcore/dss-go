// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/CMSForPAdESBaselineRequirementsChecker.java
// (DSS 6.5.RC1).
//
// slf4j logging is dropped per PORTING.md; every LOG.warn call site is called out in the
// surrounding comment instead.
//
// Java's isValidForPAdESBaselineBProfile() calls the protected cmsBaselineBRequirements() of the
// CAdES base class directly, which is CMS-level only. The Go base class lives in another package,
// so that method is reached through its exported form, cades.BaselineRequirementsChecker.
// CMSBaselineBRequirements. (Calling HasBaselineBProfile() instead would add CAdES requirement
// (k), signature-policy-store without a signature-policy-identifier defining sigPolicyHash,
// which PAdES does not have: it would reject a CMS that upstream accepts.)
//
// getBaselineSignatureForm() is overridden below, for the same reason as in
// pades_baseline_requirements_checker.go: the base's cmsBaselineBRequirements() has to see the
// PAdES form (it forbids a signing-time attribute) across the package boundary.
//
// INTEGRATOR NOTE (constructor parameter type conflict between two already-landed sibling
// files): Java's constructor is CMSForPAdESBaselineRequirementsChecker(CAdESSignature) - the
// *cades* base class, matched here as *cades.Signature (this is what
// pades_with_external_cms_service.go's landed call site at NewCMSForPAdESBaselineRequirements
// Checker(cadesSignature) passes, cadesSignature being exactly a *cades.Signature built by
// padesWithExternalCMSServiceToCAdESSignature - so that call site compiles unchanged). However
// pades_baseline_requirements_checker.go's header (line ~69) documents an assumed signature of
// NewCMSForPAdESBaselineRequirementsChecker(signature *Signature) and its landed call site
// (HasBaselineBProfile) passes b.Signature(), a *Signature - which does NOT satisfy a
// *cades.CAdESSignature parameter (embedding is not Java-style subtyping in Go). That call site
// needs a one-line fix once this file lands: NewCMSForPAdESBaselineRequirementsChecker(
// b.Signature().Signature) - passing the embedded *cades.Signature field, which is
// exactly the CAdESSignature Java's polymorphic call passes for a PAdESSignature receiver.
package pades

import (
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// CMSForPAdESBaselineRequirementsChecker is used to verify conformance of a CMSSignedData to be
// incorporated to a PDF as a PAdES signature. Port of the class
// CMSForPAdESBaselineRequirementsChecker, extending cades.BaselineRequirementsChecker.
type CMSForPAdESBaselineRequirementsChecker struct {
	*cades.BaselineRequirementsChecker
}

// NewCMSForPAdESBaselineRequirementsChecker is the default constructor, used to verify CMS of
// Signature on conformance to PAdES Baseline-B format.
// Port of the constructor CMSForPAdESBaselineRequirementsChecker(CAdESSignature).
//
// Re-registers the override target as checker itself (not the embedded
// *cades.BaselineRequirementsChecker cades.NewBaselineRequirementsChecker already
// self-registered) - exactly the same package-boundary re-registration
// pades_baseline_requirements_checker.go's own NewBaselineRequirementsChecker performs and
// documents ("STRUCTURE DEVIATION"), needed here so GetBaselineSignatureForm() below (not
// cades.BaselineRequirementsChecker's CAdES-returning one) is what
// cmsBaselineBRequirements() resolves via BaselineSignatureForm() when it runs for a PDF's
// embedded CMS - see spi/validation.BaselineRequirementsCheckerOverrides.GetBaselineSignatureForm.
func NewCMSForPAdESBaselineRequirementsChecker(signature *cades.Signature) *CMSForPAdESBaselineRequirementsChecker {
	checker := &CMSForPAdESBaselineRequirementsChecker{
		BaselineRequirementsChecker: cades.NewBaselineRequirementsChecker(signature, nil),
	}
	checker.InitBaselineRequirementsChecker(checker)
	return checker
}

// GetBaselineSignatureForm returns the signature form corresponding to the signature: PAdES,
// even though the wrapped signature value is a plain *cades.Signature (see this file's
// header "INTEGRATOR NOTE" on why b.Signature().Signature, not a *Signature, is what
// gets passed in here). Port of the protected getBaselineSignatureForm() override upstream's
// CMSForPAdESBaselineRequirementsChecker.java declares independently of
// BaselineRequirementsChecker's own identical override, for exactly this reason: this
// checker validates a PDF's embedded CMS on its own, without going through a Signature.
func (c *CMSForPAdESBaselineRequirementsChecker) GetBaselineSignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureFormPAdES
}

// IsValidForPAdESBaselineBProfile verifies validity of a CMS signature for enveloping within a
// PDF signature of PAdES-BASELINE format. Port of isValidForPAdESBaselineBProfile().
func (c *CMSForPAdESBaselineRequirementsChecker) IsValidForPAdESBaselineBProfile() bool {
	signature := c.Signature()

	// PAdES Part 1 : a DER-encoded SignedData object as specified in ETSI EN 319 122-1 [2]
	// shall be included
	if signature.CMS() == nil {
		// Upstream logs "DER-encoded SignedData object shall be included as the PDF signature
		// in the entry with the key Contents of the Signature Dictionary for {}-BASELINE-B
		// signature (General requirement (a))!".
		return false
	}
	// PAdES Part 1 : There shall only be a single signer in any PDF Signature
	if len(signature.CMS().SignerInfos()) != 1 {
		// Upstream logs "SignedData.signerInfos shall contain one and only one signerInfo for
		// {}-BASELINE-B signature (General requirement (a))!".
		return false
	}
	// PAdES Part 1 : "No data shall be encapsulated in the PKCS#7 SignedData field
	if !signature.CMS().IsDetachedSignature() {
		// Upstream logs "No data shall be encapsulated in the PKCS#7 SignedData field for
		// {}-BASELINE-B signature (General requirement (b))!".
		return false
	}
	return c.CMSBaselineBRequirements()
}
