// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/PAdESSignatureParameters.java
// (DSS 6.5.RC1).
//
// # Context / timestamp-parameters storage (DEVIATION, flagged for the integrator)
//
// Java's `extends CAdESSignatureParameters` gives PAdESSignatureParameters the SAME protected
// `context`/`contentTimestampParameters`/`signatureTimestampParameters`/
// `archiveTimestampParameters` fields SignatureParameters (via AbstractSignatureParameters)
// declares - statically typed ProfileParameters/TimestampParameters there, but this class's
// getContext()/get*TimestampParameters() overrides narrow them to the runtime ProfileParameters
// / PAdESTimestampParameters objects those overrides themselves construct. Because Java has no
// separate storage per subclass, code anywhere in the hierarchy that reads those fields - even
// through a SignatureParameters-typed reference to the very same object - transparently
// observes the PAdES-specific values.
//
// Go has neither field-level covariance nor virtual dispatch across embedding, and the already
// landed pades/external_cms_service.go and pades/pades_with_external_cms_service.go pass
// `&parameters.SignatureParameters` (the embedded CAdES struct's own address - a DIFFERENT
// object from `parameters`) into cades-package extension code that reads the CAdES-level
// GetContext()/DetachedContents() through THAT address. Shadowing `context` here (the way
// xades/xades_signature_parameters.go shadows XAdES's own, unrelated context field) would sever
// that already-relied-upon sharing: `parameters.GetContext().SetDetachedContents(...)` must keep
// landing in the SAME context object `&parameters.SignatureParameters` sees.
//
// This type therefore does NOT shadow `context`/GetContext()/DetachedContents()/GetDeterministicId
// - those keep resolving to the promoted SignatureParameters/AbstractSignatureParameters
// behaviour, which is what the landed call sites above assume - and instead:
//   - keeps a wholly separate, PAdES-only `padesContext *ProfileParameters` field used
//     exclusively for the pdfToBeSignedCache (nothing else reads or writes it), with PdfSignatureCache
//     lazily creating it, and Reinit additionally clearing it (Java's reinit() clears the single
//     shared context field, which serves both purposes at once; here the two purposes are split
//     across two fields, so both must be cleared)
//   - keeps separate `contentTimestampParameters`/`signatureTimestampParameters`/
//     `archiveTimestampParameters` fields of type *TimestampParameters (shadowing, by
//     Go's shallower-selector-wins rule, both the promoted field of the same name AND the
//     Get-prefixed CAdES accessor is left un-shadowed under its own Get-prefixed name), because
//     nothing in the landed PAdES call sites (pades_service.go, pades_level_baseline_t.go,
//     pades_level_baseline_lta.go) ever reads them through the CAdES-level accessor - every call
//     site uses the Go-idiomatic no-Get name (parameters.ContentTimestampParameters(), etc.)
//     this type defines directly.
//
// This means a caller that reaches into `&parameters.SignatureParameters` and asks it for
// GetSignatureTimestampParameters() (as cades_level_baseline_t.go's CMS-embedding path does, for
// the CAdES-Extended CMS this format wraps) sees CAdES's own lazily-created default
// TimestampParameters, not whatever PAdES-specific SignatureTimestampParameters the caller
// configured through this type's shadowed accessor. This divergence from Java's single-field
// sharing is an unavoidable consequence of Go's lack of field covariance, not a bug to silently
// paper over; it is called out here for the integrator's attention.
//
// java.io.Serializable and serialVersionUID are dropped; hashCode() has no Go counterpart.
// java.util.TimeZone -> *time.Location (native_pdf_signature_service.go already assumes this:
// `signatureParameters.SigningDate().In(signatureParameters.SigningTimeZone())`).
package pades

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// SignatureParametersDefaultSignatureSize is the default preserved space for a signature
// context (9472 - default value in pdfbox). Port of the default signatureSize field value.
const SignatureParametersDefaultSignatureSize = 9472

// SignatureParameters holds parameters to create/extend a PAdES signature.
type SignatureParameters struct {
	cades.SignatureParameters

	// reason is the signature creation reason.
	reason string

	// contactInfo is the contact info.
	contactInfo string

	// location is the signer's location.
	location string

	// signerName allows explicitly specifying the SignerName (name for the Signature). The
	// person or authority signing the document.
	signerName string

	// signatureSize defines the preserved space for a signature context.
	//
	// Default : 9472 (default value in pdfbox)
	signatureSize int

	// signatureFilter allows overriding the used Filter for a Signature.
	//
	// Default value is Adobe.PPKLite
	signatureFilter string

	// signatureSubFilter allows overriding the used subFilter for a Signature.
	//
	// Default value is ETSI.CAdES.detached
	signatureSubFilter string

	// appName is the signing app name.
	appName string

	// signatureImageParameters is used to create a visible signature in PAdES form.
	signatureImageParameters *SignatureImageParameters

	// permission allows creating a "certification signature". That allows removing
	// permission(s) in case of future change(s).
	permission enumerations.CertificationPermission

	// passwordProtection is the password used to encrypt a PDF.
	passwordProtection []byte

	// signingTimeZone is the time-zone used for signature creation.
	//
	// Default: time.Local
	signingTimeZone *time.Location

	// includeVRIDictionary defines whether the VRI dictionary should be included to a PAdES
	// signature on extension within its LT-level revision (DSS-revision).
	//
	// Default: FALSE (VRI dictionary is not included)
	includeVRIDictionary bool

	// padesContext is a PAdES-only signature creation context, used exclusively for the PDF
	// signature cache; see the file header for why this does not shadow the promoted CAdES
	// context.
	padesContext *ProfileParameters

	// contentTimestampParameters holds parameters related to the content timestamp
	// (Baseline-B); see the file header for why this shadows, rather than reuses, the promoted
	// CAdES field of the same name.
	contentTimestampParameters *TimestampParameters

	// signatureTimestampParameters holds parameters related to the signature timestamp
	// (Baseline-T). See contentTimestampParameters doc.
	signatureTimestampParameters *TimestampParameters

	// archiveTimestampParameters holds parameters related to the archive timestamp
	// (Baseline-LTA). See contentTimestampParameters doc.
	archiveTimestampParameters *TimestampParameters
}

// NewSignatureParameters is the default constructor instantiating object with default
// parameters.
func NewSignatureParameters() *SignatureParameters {
	return &SignatureParameters{
		SignatureParameters: *cades.NewSignatureParameters(),
		signatureSize:       SignatureParametersDefaultSignatureSize,
		signatureFilter:     PAdESConstantsSignatureDefaultFilter,
		signatureSubFilter:  PAdESConstantsSignatureDefaultSubFilter,
		signingTimeZone:     time.Local,
	}
}

// SetSignatureLevel ports the overridden #setSignatureLevel. Panics with the Java message when
// signatureLevel is empty or is not a PAdES level (Java's IllegalArgumentException).
func (p *SignatureParameters) SetSignatureLevel(signatureLevel enumerations.SignatureLevel) {
	form, err := signatureLevel.SignatureForm()
	if err != nil || enumerations.SignatureFormPAdES != form {
		panic("Only PAdES form is allowed !")
	}
	p.SignatureParameters.SetSignatureLevel(signatureLevel)
}

// Reason gets the reason. Port of #getReason.
func (p *SignatureParameters) Reason() string {
	return p.reason
}

// SetReason sets the reason. Port of #setReason.
func (p *SignatureParameters) SetReason(reason string) {
	p.reason = reason
}

// ContactInfo gets the contactInfo. Port of #getContactInfo.
func (p *SignatureParameters) ContactInfo() string {
	return p.contactInfo
}

// SetContactInfo sets the contactInfo. Port of #setContactInfo.
func (p *SignatureParameters) SetContactInfo(contactInfo string) {
	p.contactInfo = contactInfo
}

// Location gets location. Port of #getLocation.
func (p *SignatureParameters) Location() string {
	return p.location
}

// SetLocation sets location (The CPU host name or physical location of the signing). Port of
// #setLocation.
func (p *SignatureParameters) SetLocation(location string) {
	p.location = location
}

// SignerName returns the Signer Name. Port of #getSignerName.
func (p *SignatureParameters) SignerName() string {
	return p.signerName
}

// SetSignerName sets the name of the signer. Port of #setSignerName.
func (p *SignatureParameters) SetSignerName(signerName string) {
	p.signerName = signerName
}

// Filter ports the overridden #getFilter.
func (p *SignatureParameters) Filter() string {
	return p.signatureFilter
}

// SetFilter sets the filter. Port of #setFilter.
func (p *SignatureParameters) SetFilter(signatureFilter string) {
	p.signatureFilter = signatureFilter
}

// SubFilter ports the overridden #getSubFilter.
func (p *SignatureParameters) SubFilter() string {
	return p.signatureSubFilter
}

// SetSubFilter sets the sub filter. Port of #setSubFilter.
func (p *SignatureParameters) SetSubFilter(signatureSubFilter string) {
	p.signatureSubFilter = signatureSubFilter
}

// AppName ports the overridden #getAppName.
func (p *SignatureParameters) AppName() string {
	return p.appName
}

// SetAppName sets signing application name. Port of #setAppName.
func (p *SignatureParameters) SetAppName(appName string) {
	p.appName = appName
}

// ImageParameters ports the overridden #getImageParameters, lazily instantiating it.
func (p *SignatureParameters) ImageParameters() *SignatureImageParameters {
	if p.signatureImageParameters == nil {
		p.signatureImageParameters = NewSignatureImageParameters()
	}
	return p.signatureImageParameters
}

// SetImageParameters sets the SignatureImageParameters for a visual signature creation. Port of
// #setImageParameters.
func (p *SignatureParameters) SetImageParameters(signatureImageParameters *SignatureImageParameters) {
	p.signatureImageParameters = signatureImageParameters
}

// ContentSize ports the overridden #getContentSize.
func (p *SignatureParameters) ContentSize() int {
	return p.signatureSize
}

// SetContentSize defines an amount of bytes to be reserved for a CMS signature contents
// encapsulation.
//
// Default : 9472 bytes
//
// Port of #setContentSize.
func (p *SignatureParameters) SetContentSize(signatureSize int) {
	p.signatureSize = signatureSize
}

// Permission gets the permission for the PDF document modification. Port of #getPermission.
func (p *SignatureParameters) Permission() enumerations.CertificationPermission {
	return p.permission
}

// SetPermission sets the permission for the PDF document modification. Port of #setPermission.
func (p *SignatureParameters) SetPermission(permission enumerations.CertificationPermission) {
	p.permission = permission
}

// PasswordProtection ports the overridden #getPasswordProtection.
func (p *SignatureParameters) PasswordProtection() []byte {
	return p.passwordProtection
}

// SetPasswordProtection sets password to the document. Port of #setPasswordProtection.
func (p *SignatureParameters) SetPasswordProtection(passwordProtection []byte) {
	p.passwordProtection = passwordProtection
}

// SigningDate ports the overridden #getSigningDate.
func (p *SignatureParameters) SigningDate() *time.Time {
	return p.BLevel().SigningDate()
}

// DeterministicId implements PAdESCommonParameters#DeterministicId with the Go-idiomatic no-Get
// name; the promoted AbstractSignatureParameters method underneath keeps the Get prefix
// (document/abstract_signature_parameters.go's GetDeterministicId - PORTING.md's get-prefix-drop
// convention was applied to CommonParameters here but not, historically, to
// AbstractSignatureParameters), so this delegates to it under the interface's expected name.
func (p *SignatureParameters) DeterministicId() string {
	return p.GetDeterministicId()
}

// SetSigningTimeZone sets a TimeZone to use for signature creation. Will be used to define a
// signingTime within a PDF entry with key /M.
//
// Default: time.Local
//
// Port of #setSigningTimeZone.
func (p *SignatureParameters) SetSigningTimeZone(signingTimeZone *time.Location) {
	p.signingTimeZone = signingTimeZone
}

// SigningTimeZone returns a time-zone used to define the signing time. Port of
// #getSigningTimeZone.
func (p *SignatureParameters) SigningTimeZone() *time.Location {
	return p.signingTimeZone
}

// IsIncludeVRIDictionary returns whether the VRI dictionary should be included to the PAdES
// Signature on extension within LT-level revision (DSS revision). Port of
// #isIncludeVRIDictionary.
func (p *SignatureParameters) IsIncludeVRIDictionary() bool {
	return p.includeVRIDictionary
}

// SetIncludeVRIDictionary sets whether corresponding VRI dictionary should be included to the
// PAdES signature on extension to LT-level.
//
// Default: FALSE (VRI dictionary is not included to PAdES signature)
//
// Port of #setIncludeVRIDictionary.
func (p *SignatureParameters) SetIncludeVRIDictionary(includeVRIDictionary bool) {
	p.includeVRIDictionary = includeVRIDictionary
}

// getPAdESContext lazily creates the PAdES-only signature creation context; see the file header.
func (p *SignatureParameters) getPAdESContext() *ProfileParameters {
	if p.padesContext == nil {
		p.padesContext = NewProfileParameters()
	}
	return p.padesContext
}

// ContentTimestampParameters shadows the promoted field of the same name; see the file header.
func (p *SignatureParameters) ContentTimestampParameters() *TimestampParameters {
	if p.contentTimestampParameters == nil {
		p.contentTimestampParameters = NewTimestampParameters()
	}
	return p.contentTimestampParameters
}

// SetContentTimestampParameters ports the overridden #setContentTimestampParameters. See the
// file header: since *cades.TimestampParameters can never dynamically be a
// *PAdESTimestampParameters in Go (no runtime subtyping the way Java's instanceof check relies
// on), this always wraps, unlike Java's conditional cast.
func (p *SignatureParameters) SetContentTimestampParameters(contentTimestampParameters *cades.TimestampParameters) {
	p.contentTimestampParameters = newTimestampParametersFromCAdES(contentTimestampParameters)
}

// SignatureTimestampParameters shadows the promoted field of the same name; see the file header.
func (p *SignatureParameters) SignatureTimestampParameters() *TimestampParameters {
	if p.signatureTimestampParameters == nil {
		p.signatureTimestampParameters = NewTimestampParameters()
	}
	return p.signatureTimestampParameters
}

// SetSignatureTimestampParameters ports the overridden #setSignatureTimestampParameters. See
// SetContentTimestampParameters's doc for why this always wraps.
func (p *SignatureParameters) SetSignatureTimestampParameters(signatureTimestampParameters *cades.TimestampParameters) {
	p.signatureTimestampParameters = newTimestampParametersFromCAdES(signatureTimestampParameters)
}

// ArchiveTimestampParameters shadows the promoted field of the same name; see the file header.
func (p *SignatureParameters) ArchiveTimestampParameters() *TimestampParameters {
	if p.archiveTimestampParameters == nil {
		p.archiveTimestampParameters = NewTimestampParameters()
	}
	return p.archiveTimestampParameters
}

// SetArchiveTimestampParameters ports the overridden #setArchiveTimestampParameters. See
// SetContentTimestampParameters's doc for why this always wraps.
func (p *SignatureParameters) SetArchiveTimestampParameters(archiveTimestampParameters *cades.TimestampParameters) {
	p.archiveTimestampParameters = newTimestampParametersFromCAdES(archiveTimestampParameters)
}

// PdfSignatureCache ports the overridden #getPdfSignatureCache.
func (p *SignatureParameters) PdfSignatureCache() *PdfSignatureCache {
	return p.getPAdESContext().PdfToBeSignedCache()
}

// Reinit ports the promoted Reinit, additionally clearing the PAdES-only context; see the file
// header for why the base's single reinit() clear must be split across two fields here.
func (p *SignatureParameters) Reinit() {
	p.SignatureParameters.Reinit()
	p.padesContext = nil
}

// String ports #toString.
func (p *SignatureParameters) String() string {
	return fmt.Sprintf("PAdESSignatureParameters [reason='%s', contactInfo='%s', location='%s', signerName='%s', "+
		"signatureSize=%d, signatureFilter='%s', signatureSubFilter='%s', appName='%s', signatureImageParameters=%v, "+
		"permission=%v, passwordProtection=%v, signingTimeZone=%v, includeVRIDictionary=%v] %s",
		p.reason, p.contactInfo, p.location, p.signerName, p.signatureSize, p.signatureFilter, p.signatureSubFilter,
		p.appName, p.signatureImageParameters, p.permission, p.passwordProtection, p.signingTimeZone,
		p.includeVRIDictionary, p.SignatureParameters.String())
}

// Equals ports #equals.
func (p *SignatureParameters) Equals(other *SignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.SignatureParameters.Equals(&other.SignatureParameters) {
		return false
	}
	return p.signatureSize == other.signatureSize &&
		p.includeVRIDictionary == other.includeVRIDictionary &&
		p.reason == other.reason &&
		p.contactInfo == other.contactInfo &&
		p.location == other.location &&
		p.signerName == other.signerName &&
		p.signatureFilter == other.signatureFilter &&
		p.signatureSubFilter == other.signatureSubFilter &&
		p.appName == other.appName &&
		signatureImageParametersEqualPointers(p.signatureImageParameters, other.signatureImageParameters) &&
		p.permission == other.permission &&
		string(p.passwordProtection) == string(other.passwordProtection) &&
		p.signingTimeZone == other.signingTimeZone
}

// signatureImageParametersEqualPointers ports Objects.equals(signatureImageParameters, ...).
func signatureImageParametersEqualPointers(a, b *SignatureImageParameters) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

var _ CommonParameters = (*SignatureParameters)(nil)
