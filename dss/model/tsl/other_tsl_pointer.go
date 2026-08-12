// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/OtherTSLPointer.java (DSS 6.5.RC1).
package tsl

import "github.com/utain/esig/dss/model"

// OtherTSLPointer contains certificates for the url location.
//
// Implements the assumed job.OtherDocumentPointer interface (getLocation/getSdiCertificates,
// ported as Location()/SdiCertificates()) and java.io.Serializable, which has no Go
// counterpart and is dropped.
type OtherTSLPointer struct {
	// sdiCertificates is the list of ServiceDigitalIdentity X509 certificates.
	sdiCertificates []*model.CertificateToken
	// tslLocation is the URL location.
	tslLocation string
	// schemeTerritory is an ISO code of the country or an alliance.
	schemeTerritory string
	// tslType is the type of the Trusted List.
	tslType string
	// mimeType is the MimeType of the Trusted List document.
	mimeType string
	// schemeOperatorNames maps defined scheme operator names between the used languages.
	schemeOperatorNames map[string][]string
	// schemeTypeCommunityRules maps defined type community rules between the used languages.
	schemeTypeCommunityRules map[string][]string
	// mra is the Mutual Recognition Agreement block.
	mra *MRA
}

// NewOtherTSLPointer instantiates an empty object. Port of the default constructor.
func NewOtherTSLPointer() *OtherTSLPointer {
	return &OtherTSLPointer{}
}

// NewOtherTSLPointerFromBuilder instantiates an object from an OtherTSLPointerBuilder. Port of
// the OtherTSLPointer(OtherTSLPointerBuilder) constructor.
func NewOtherTSLPointerFromBuilder(builder *OtherTSLPointerBuilder) *OtherTSLPointer {
	return &OtherTSLPointer{
		sdiCertificates:          builder.SdiCertificates(),
		tslLocation:              builder.TslLocation(),
		schemeTerritory:          builder.SchemeTerritory(),
		tslType:                  builder.TslType(),
		mimeType:                 builder.MimeType(),
		schemeOperatorNames:      builder.SchemeOperatorNames(),
		schemeTypeCommunityRules: builder.SchemeTypeCommunityRules(),
		mra:                      builder.Mra(),
	}
}

// Location ports the OtherDocumentPointer#getLocation() override, which returns the TSL
// location.
func (o *OtherTSLPointer) Location() string {
	return o.TSLLocation()
}

// SdiCertificates gets a list of ServiceDigitalIdentity X509 certificates.
func (o *OtherTSLPointer) SdiCertificates() []*model.CertificateToken {
	return o.sdiCertificates
}

// TSLLocation gets the TSL location url.
func (o *OtherTSLPointer) TSLLocation() string {
	return o.tslLocation
}

// SchemeTerritory gets the scheme territory ISO country code.
func (o *OtherTSLPointer) SchemeTerritory() string {
	return o.schemeTerritory
}

// TslType gets the TSL Type.
func (o *OtherTSLPointer) TslType() string {
	return o.tslType
}

// MimeType gets the MimeType of the referenced document.
func (o *OtherTSLPointer) MimeType() string {
	return o.mimeType
}

// SchemeOperatorNames gets a map of scheme operator names.
func (o *OtherTSLPointer) SchemeOperatorNames() map[string][]string {
	return o.schemeOperatorNames
}

// SchemeTypeCommunityRules gets a map of scheme type community rules.
func (o *OtherTSLPointer) SchemeTypeCommunityRules() map[string][]string {
	return o.schemeTypeCommunityRules
}

// Mra gets a Mutual Recognition Agreement block.
func (o *OtherTSLPointer) Mra() *MRA {
	return o.mra
}

// OtherTSLPointerBuilder builds OtherTSLPointer.
type OtherTSLPointerBuilder struct {
	sdiCertificates          []*model.CertificateToken
	tslLocation              string
	schemeTerritory          string
	tslType                  string
	mimeType                 string
	schemeOperatorNames      map[string][]string
	schemeTypeCommunityRules map[string][]string
	mra                      *MRA
}

// NewOtherTSLPointerBuilder is the default constructor.
func NewOtherTSLPointerBuilder() *OtherTSLPointerBuilder {
	return &OtherTSLPointerBuilder{}
}

// SdiCertificates gets the ServiceDigitalIdentity X509 certificates.
func (b *OtherTSLPointerBuilder) SdiCertificates() []*model.CertificateToken {
	return b.sdiCertificates
}

// SetSdiCertificates sets the ServiceDigitalIdentity X509 certificates.
func (b *OtherTSLPointerBuilder) SetSdiCertificates(sdiCertificates []*model.CertificateToken) *OtherTSLPointerBuilder {
	b.sdiCertificates = sdiCertificates
	return b
}

// TslLocation gets the TSL location URL.
func (b *OtherTSLPointerBuilder) TslLocation() string {
	return b.tslLocation
}

// SetTslLocation sets the TSL location URL.
func (b *OtherTSLPointerBuilder) SetTslLocation(tslLocation string) *OtherTSLPointerBuilder {
	b.tslLocation = tslLocation
	return b
}

// SchemeTerritory gets the scheme territory code.
func (b *OtherTSLPointerBuilder) SchemeTerritory() string {
	return b.schemeTerritory
}

// SetSchemeTerritory sets the scheme territory code.
func (b *OtherTSLPointerBuilder) SetSchemeTerritory(schemeTerritory string) *OtherTSLPointerBuilder {
	b.schemeTerritory = schemeTerritory
	return b
}

// TslType gets the TSL Type.
func (b *OtherTSLPointerBuilder) TslType() string {
	return b.tslType
}

// SetTslType sets the TSL Type.
func (b *OtherTSLPointerBuilder) SetTslType(tslType string) *OtherTSLPointerBuilder {
	b.tslType = tslType
	return b
}

// MimeType gets the MimeType of the Trusted List document.
func (b *OtherTSLPointerBuilder) MimeType() string {
	return b.mimeType
}

// SetMimeType sets the MimeType of the Trusted List document.
func (b *OtherTSLPointerBuilder) SetMimeType(mimeType string) *OtherTSLPointerBuilder {
	b.mimeType = mimeType
	return b
}

// SchemeOperatorNames gets a map of scheme operator names.
func (b *OtherTSLPointerBuilder) SchemeOperatorNames() map[string][]string {
	return b.schemeOperatorNames
}

// SetSchemeOperatorNames sets a map of scheme operator names.
func (b *OtherTSLPointerBuilder) SetSchemeOperatorNames(schemeOperatorNames map[string][]string) *OtherTSLPointerBuilder {
	b.schemeOperatorNames = schemeOperatorNames
	return b
}

// SchemeTypeCommunityRules gets a map of scheme type community rules.
func (b *OtherTSLPointerBuilder) SchemeTypeCommunityRules() map[string][]string {
	return b.schemeTypeCommunityRules
}

// SetSchemeTypeCommunityRules sets a map of scheme type community rules.
func (b *OtherTSLPointerBuilder) SetSchemeTypeCommunityRules(schemeTypeCommunityRules map[string][]string) *OtherTSLPointerBuilder {
	b.schemeTypeCommunityRules = schemeTypeCommunityRules
	return b
}

// Mra gets the MRA (Mutual Recognition Agreement) scheme.
func (b *OtherTSLPointerBuilder) Mra() *MRA {
	return b.mra
}

// SetMra sets the MRA (Mutual Recognition Agreement) scheme.
func (b *OtherTSLPointerBuilder) SetMra(mra *MRA) *OtherTSLPointerBuilder {
	b.mra = mra
	return b
}

// Build builds the OtherTSLPointer.
func (b *OtherTSLPointerBuilder) Build() *OtherTSLPointer {
	return NewOtherTSLPointerFromBuilder(b)
}
