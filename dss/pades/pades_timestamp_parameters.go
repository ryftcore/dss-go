// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/PAdESTimestampParameters.java
// (DSS 6.5.RC1).
//
// Java's `extends CAdESTimestampParameters` becomes embedding; PAdESSignatureParameters shadows
// (rather than reuses) its own Content/Signature/ArchiveTimestampParameters storage with this
// type - see pades_signature_parameters.go's file header for why (Go has no field-level
// covariance across embedding).
package pades

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PAdESTimestampParametersDefaultContentSize is the default length of a reserved space for the
// timestamp inside a /Contents attribute (9472, from PDFBox). Port of the default timestampSize
// field value.
const PAdESTimestampParametersDefaultContentSize = 9472

// PAdESTimestampParameters holds parameters for a PAdES timestamp creation.
type PAdESTimestampParameters struct {
	cades.CAdESTimestampParameters

	// pdfSignatureCache is the internal signature processing variable.
	pdfSignatureCache *PdfSignatureCache

	// timestampDate is the date of the timestamp.
	timestampDate *time.Time

	// timestampSize is the length of a reserved space for the timestamp inside a /Contents
	// attribute.
	//
	// Default value is 9472 (from PDFBox)
	timestampSize int

	// timestampFilter allows overriding the used Filter for a Timestamp.
	//
	// Default value is Adobe.PPKLite
	timestampFilter string

	// timestampSubFilter allows overriding the used subFilter for a Timestamp.
	//
	// Default value is ETSI.RFC3161
	timestampSubFilter string

	// appName is the signing app name.
	appName string

	// timestampImageParameters is used to create a visible timestamp in PAdES form.
	timestampImageParameters *SignatureImageParameters

	// passwordProtection is the password used to encrypt a PDF.
	passwordProtection []byte
}

// NewPAdESTimestampParameters is the empty constructor.
func NewPAdESTimestampParameters() *PAdESTimestampParameters {
	now := time.Now()
	return &PAdESTimestampParameters{
		CAdESTimestampParameters: *cades.NewCAdESTimestampParameters(),
		timestampDate:            &now,
		timestampSize:            PAdESTimestampParametersDefaultContentSize,
		timestampFilter:          PAdESConstantsTimestampDefaultFilter,
		timestampSubFilter:       PAdESConstantsTimestampDefaultSubFilter,
	}
}

// NewPAdESTimestampParametersWithDigestAlgorithm is the default constructor. Port of
// PAdESTimestampParameters(DigestAlgorithm).
func NewPAdESTimestampParametersWithDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *PAdESTimestampParameters {
	p := NewPAdESTimestampParameters()
	p.SetDigestAlgorithm(digestAlgorithm)
	return p
}

// newPAdESTimestampParametersFromCAdES is used internally to recreate parameters from CAdES
// Timestamp Parameters. Package-private port of the package-private constructor
// PAdESTimestampParameters(CAdESTimestampParameters).
func newPAdESTimestampParametersFromCAdES(cadesTimestampParameters *cades.CAdESTimestampParameters) *PAdESTimestampParameters {
	return NewPAdESTimestampParametersWithDigestAlgorithm(cadesTimestampParameters.DigestAlgorithm())
}

// Filter ports the overridden #getFilter.
func (p *PAdESTimestampParameters) Filter() string {
	return p.timestampFilter
}

// SetFilter sets the filter. Port of #setFilter.
func (p *PAdESTimestampParameters) SetFilter(timestampFilter string) {
	p.timestampFilter = timestampFilter
}

// SubFilter ports the overridden #getSubFilter.
func (p *PAdESTimestampParameters) SubFilter() string {
	return p.timestampSubFilter
}

// SetSubFilter sets the sub filter. Port of #setSubFilter.
func (p *PAdESTimestampParameters) SetSubFilter(timestampSubFilter string) {
	p.timestampSubFilter = timestampSubFilter
}

// AppName ports the overridden #getAppName.
func (p *PAdESTimestampParameters) AppName() string {
	return p.appName
}

// SetAppName sets signing application name. Port of #setAppName.
func (p *PAdESTimestampParameters) SetAppName(appName string) {
	p.appName = appName
}

// EncryptionAlgorithm ports the overridden #getEncryptionAlgorithm (not implemented; returns the
// zero value, mirroring Java's `return null`).
func (p *PAdESTimestampParameters) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return ""
}

// ImageParameters ports the overridden #getImageParameters, lazily instantiating it.
func (p *PAdESTimestampParameters) ImageParameters() *SignatureImageParameters {
	if p.timestampImageParameters == nil {
		p.timestampImageParameters = NewSignatureImageParameters()
	}
	return p.timestampImageParameters
}

// SetImageParameters sets the SignatureImageParameters for a visual timestamp creation. Port of
// #setImageParameters.
func (p *PAdESTimestampParameters) SetImageParameters(timestampImageParameters *SignatureImageParameters) {
	p.timestampImageParameters = timestampImageParameters
}

// ContentSize ports the overridden #getContentSize.
func (p *PAdESTimestampParameters) ContentSize() int {
	return p.timestampSize
}

// SetContentSize allows reserving more than the default size for a timestamp.
//
// Default : 9472 bytes
//
// Port of #setContentSize.
func (p *PAdESTimestampParameters) SetContentSize(timestampSize int) {
	p.timestampSize = timestampSize
}

// SigningDate ports the overridden #getSigningDate.
func (p *PAdESTimestampParameters) SigningDate() *time.Time {
	return p.timestampDate
}

// PasswordProtection ports the overridden #getPasswordProtection.
func (p *PAdESTimestampParameters) PasswordProtection() []byte {
	return p.passwordProtection
}

// SetPasswordProtection sets password to the document. Port of #setPasswordProtection.
func (p *PAdESTimestampParameters) SetPasswordProtection(passwordProtection []byte) {
	p.passwordProtection = passwordProtection
}

// DeterministicId ports the overridden #getDeterministicId. Panics on the (practically
// unreachable) MD5 failure path, matching Java's unchecked DSSException.
func (p *PAdESTimestampParameters) DeterministicId() string {
	var signingTime time.Time
	if p.timestampDate != nil {
		signingTime = *p.timestampDate
	}
	deterministicId, err := spi.DSSUtilsDeterministicID(signingTime, nil)
	if err != nil {
		panic(err)
	}
	return deterministicId
}

// PdfSignatureCache ports the overridden #getPdfSignatureCache, lazily instantiating it.
func (p *PAdESTimestampParameters) PdfSignatureCache() *PdfSignatureCache {
	if p.pdfSignatureCache == nil {
		p.pdfSignatureCache = NewPdfSignatureCache()
	}
	return p.pdfSignatureCache
}

// Reinit ports the overridden #reinit.
func (p *PAdESTimestampParameters) Reinit() {
	p.pdfSignatureCache = nil
}

// String ports #toString.
func (p *PAdESTimestampParameters) String() string {
	return fmt.Sprintf("PAdESTimestampParameters [pdfSignatureCache=%v, timestampDate=%v, timestampSize=%v, "+
		"timestampFilter='%s', timestampSubFilter='%s', appName='%s', timestampImageParameters=%v, "+
		"passwordProtection=%v] %s",
		p.pdfSignatureCache, p.timestampDate, p.timestampSize, p.timestampFilter, p.timestampSubFilter,
		p.appName, p.timestampImageParameters, p.passwordProtection, p.CAdESTimestampParameters.String())
}

// Equals ports #equals.
func (p *PAdESTimestampParameters) Equals(other *PAdESTimestampParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.TimestampParameters.Equals(&other.TimestampParameters) {
		return false
	}
	return p.timestampSize == other.timestampSize &&
		p.pdfSignatureCache == other.pdfSignatureCache &&
		pAdESTimestampParametersTimeEquals(p.timestampDate, other.timestampDate) &&
		p.timestampFilter == other.timestampFilter &&
		p.timestampSubFilter == other.timestampSubFilter &&
		p.appName == other.appName &&
		p.timestampImageParameters == other.timestampImageParameters &&
		string(p.passwordProtection) == string(other.passwordProtection)
}

// pAdESTimestampParametersTimeEquals ports Objects.equals(timestampDate, ...) for the *time.Time
// field.
func pAdESTimestampParametersTimeEquals(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equal(*b)
}

var _ PAdESCommonParameters = (*PAdESTimestampParameters)(nil)
