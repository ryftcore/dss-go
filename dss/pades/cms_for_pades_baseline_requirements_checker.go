// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/CMSForPAdESBaselineRequirementsChecker.java
// (DSS 6.5.RC1).
//
// slf4j logging is dropped per PORTING.md; every LOG.warn call site is called out in the
// surrounding comment instead.
//
// DEVIATION: Java's isValidForPAdESBaselineBProfile() calls the protected
// cmsBaselineBRequirements() directly; cades_baseline_requirements_checker.go keeps that method
// unexported (cades package-private), with no exported equivalent, so it is unreachable from
// this package. HasBaselineBProfile() is used instead - it calls cmsBaselineBRequirements()
// internally and additionally enforces requirement (k) (signature-policy-store must not be
// present without a signature-policy-identifier defining sigPolicyHash), a strictly narrower
// check. A signature failing only requirement (k) is therefore rejected here where upstream
// Java would accept it; every other outcome is identical.
//
// getBaselineSignatureForm() is not overridden for the same reason (it is unexported in the
// base) and, per cades_baseline_requirements_checker.go's header, its return value was read only
// by slf4j log statements upstream (dropped per PORTING.md) - so the override would have had no
// observable effect even if it were reachable.
//
// INTEGRATOR NOTE (constructor parameter type conflict between two already-landed sibling
// files): Java's constructor is CMSForPAdESBaselineRequirementsChecker(CAdESSignature) - the
// *cades* base class, matched here as *cades.CAdESSignature (this is what
// pades_with_external_cms_service.go's landed call site at NewCMSForPAdESBaselineRequirements
// Checker(cadesSignature) passes, cadesSignature being exactly a *cades.CAdESSignature built by
// padesWithExternalCMSServiceToCAdESSignature - so that call site compiles unchanged). However
// pades_baseline_requirements_checker.go's header (line ~69) documents an assumed signature of
// NewCMSForPAdESBaselineRequirementsChecker(signature *PAdESSignature) and its landed call site
// (HasBaselineBProfile) passes b.Signature(), a *PAdESSignature - which does NOT satisfy a
// *cades.CAdESSignature parameter (embedding is not Java-style subtyping in Go). That call site
// needs a one-line fix once this file lands: NewCMSForPAdESBaselineRequirementsChecker(
// b.Signature().CAdESSignature) - passing the embedded *cades.CAdESSignature field, which is
// exactly the CAdESSignature Java's polymorphic call passes for a PAdESSignature receiver.
package pades

import (
	"github.com/utain/esig/dss/cades"
	"github.com/utain/esig/dss/enumerations"
)

// CMSForPAdESBaselineRequirementsChecker is used to verify conformance of a CMSSignedData to be
// incorporated to a PDF as a PAdES signature. Port of the class
// CMSForPAdESBaselineRequirementsChecker, extending cades.CAdESBaselineRequirementsChecker.
type CMSForPAdESBaselineRequirementsChecker struct {
	*cades.CAdESBaselineRequirementsChecker
}

// NewCMSForPAdESBaselineRequirementsChecker is the default constructor, used to verify CMS of
// CAdESSignature on conformance to PAdES Baseline-B format.
// Port of the constructor CMSForPAdESBaselineRequirementsChecker(CAdESSignature).
//
// Re-registers the override target as checker itself (not the embedded
// *cades.CAdESBaselineRequirementsChecker cades.NewCAdESBaselineRequirementsChecker already
// self-registered) - exactly the same package-boundary re-registration
// pades_baseline_requirements_checker.go's own NewPAdESBaselineRequirementsChecker performs and
// documents ("STRUCTURE DEVIATION"), needed here so GetBaselineSignatureForm() below (not
// cades.CAdESBaselineRequirementsChecker's CAdES-returning one) is what
// cmsBaselineBRequirements() resolves via BaselineSignatureForm() when it runs for a PDF's
// embedded CMS - see spi/validation.BaselineRequirementsCheckerOverrides.GetBaselineSignatureForm.
func NewCMSForPAdESBaselineRequirementsChecker(signature *cades.CAdESSignature) *CMSForPAdESBaselineRequirementsChecker {
	checker := &CMSForPAdESBaselineRequirementsChecker{
		CAdESBaselineRequirementsChecker: cades.NewCAdESBaselineRequirementsChecker(signature, nil),
	}
	checker.InitBaselineRequirementsChecker(checker)
	return checker
}

// GetBaselineSignatureForm returns the signature form corresponding to the signature: PAdES,
// even though the wrapped signature value is a plain *cades.CAdESSignature (see this file's
// header "INTEGRATOR NOTE" on why b.Signature().CAdESSignature, not a *PAdESSignature, is what
// gets passed in here). Port of the protected getBaselineSignatureForm() override upstream's
// CMSForPAdESBaselineRequirementsChecker.java declares independently of
// PAdESBaselineRequirementsChecker's own identical override, for exactly this reason: this
// checker validates a PDF's embedded CMS on its own, without going through a PAdESSignature.
func (c *CMSForPAdESBaselineRequirementsChecker) GetBaselineSignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureForm_PAdES
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
	return c.HasBaselineBProfile()
}
