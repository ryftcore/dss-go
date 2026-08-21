// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/OtherListPointer.java (DSS 6.5.RC1).
package lote

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/tsl"
)

// OtherListPointer contains information about a reference to another List, including URL and
// signing certificates.
//
// java.io.Serializable has no Go counterpart and is dropped.
type OtherListPointer struct {
	// sdiCertificates is the list of ServiceDigitalIdentity X509 certificates.
	sdiCertificates []*model.CertificateToken
	// locationUrl is the URL location.
	locationUrl string
	// schemeTerritory is an ISO code of the country or an alliance.
	schemeTerritory string
	// typ is the type of the Trusted List.
	typ string
	// mimeType is the MimeType of the Trusted List document.
	mimeType string
	// schemeOperatorNames maps defined scheme operator names between the used languages.
	schemeOperatorNames map[string][]string
	// schemeTypeCommunityRules maps defined type community rules between the used languages.
	schemeTypeCommunityRules map[string][]string
}

// NewOtherListPointer instantiates an empty object. Port of the default constructor.
func NewOtherListPointer() *OtherListPointer {
	return &OtherListPointer{}
}

// NewOtherListPointerFromBuilder instantiates an object from an OtherListPointerBuilder. Port
// of the OtherListPointer(OtherListPointerBuilder) constructor.
//
// DEVIATION (kept verbatim from upstream): the builder also carries an "mra" field (Mutual
// Recognition Agreement block) with a getter/setter, but the Java constructor never copies it
// into the built OtherListPointer - OtherListPointer has no mra field or getter at all. This
// looks like an upstream oversight; it is preserved here rather than "fixed".
func NewOtherListPointerFromBuilder(builder *OtherListPointerBuilder) *OtherListPointer {
	return &OtherListPointer{
		sdiCertificates:          builder.SdiCertificates(),
		locationUrl:              builder.LocationUrl(),
		schemeTerritory:          builder.SchemeTerritory(),
		typ:                      builder.Type(),
		mimeType:                 builder.MimeType(),
		schemeOperatorNames:      builder.SchemeOperatorNames(),
		schemeTypeCommunityRules: builder.SchemeTypeCommunityRules(),
	}
}

// SdiCertificates gets a list of ServiceDigitalIdentity X509 certificates.
func (o *OtherListPointer) SdiCertificates() []*model.CertificateToken {
	return o.sdiCertificates
}

// LocationUrl gets List location url.
func (o *OtherListPointer) LocationUrl() string {
	return o.locationUrl
}

// SchemeTerritory gets the scheme territory ISO country code.
func (o *OtherListPointer) SchemeTerritory() string {
	return o.schemeTerritory
}

// Type gets the List Type.
func (o *OtherListPointer) Type() string {
	return o.typ
}

// MimeType gets the MimeType of the referenced document.
func (o *OtherListPointer) MimeType() string {
	return o.mimeType
}

// SchemeOperatorNames gets a map of scheme operator names.
func (o *OtherListPointer) SchemeOperatorNames() map[string][]string {
	return o.schemeOperatorNames
}

// SchemeTypeCommunityRules gets a map of scheme type community rules.
func (o *OtherListPointer) SchemeTypeCommunityRules() map[string][]string {
	return o.schemeTypeCommunityRules
}

// OtherListPointerBuilder builds OtherListPointer.
type OtherListPointerBuilder struct {
	sdiCertificates          []*model.CertificateToken
	locationUrl              string
	schemeTerritory          string
	tslType                  string
	mimeType                 string
	schemeOperatorNames      map[string][]string
	schemeTypeCommunityRules map[string][]string
	mra                      *tsl.MRA
}

// NewOtherListPointerBuilder is the default constructor.
func NewOtherListPointerBuilder() *OtherListPointerBuilder {
	return &OtherListPointerBuilder{}
}

// SdiCertificates gets the ServiceDigitalIdentity X509 certificates.
func (b *OtherListPointerBuilder) SdiCertificates() []*model.CertificateToken {
	return b.sdiCertificates
}

// SetSdiCertificates sets the ServiceDigitalIdentity X509 certificates.
func (b *OtherListPointerBuilder) SetSdiCertificates(sdiCertificates []*model.CertificateToken) *OtherListPointerBuilder {
	b.sdiCertificates = sdiCertificates
	return b
}

// LocationUrl gets the List location URL.
func (b *OtherListPointerBuilder) LocationUrl() string {
	return b.locationUrl
}

// SetLocationUrl sets the List location URL.
func (b *OtherListPointerBuilder) SetLocationUrl(locationUrl string) *OtherListPointerBuilder {
	b.locationUrl = locationUrl
	return b
}

// SchemeTerritory gets the scheme territory code.
func (b *OtherListPointerBuilder) SchemeTerritory() string {
	return b.schemeTerritory
}

// SetSchemeTerritory sets the scheme territory code.
func (b *OtherListPointerBuilder) SetSchemeTerritory(schemeTerritory string) *OtherListPointerBuilder {
	b.schemeTerritory = schemeTerritory
	return b
}

// Type gets the TSL Type. Port of the builder's getType(), which reads the tslType field.
func (b *OtherListPointerBuilder) Type() string {
	return b.tslType
}

// SetTslType sets the TSL Type.
func (b *OtherListPointerBuilder) SetTslType(tslType string) *OtherListPointerBuilder {
	b.tslType = tslType
	return b
}

// MimeType gets the MimeType of the Trusted List document.
func (b *OtherListPointerBuilder) MimeType() string {
	return b.mimeType
}

// SetMimeType sets the MimeType of the Trusted List document.
func (b *OtherListPointerBuilder) SetMimeType(mimeType string) *OtherListPointerBuilder {
	b.mimeType = mimeType
	return b
}

// SchemeOperatorNames gets a map of scheme operator names.
func (b *OtherListPointerBuilder) SchemeOperatorNames() map[string][]string {
	return b.schemeOperatorNames
}

// SetSchemeOperatorNames sets a map of scheme operator names.
func (b *OtherListPointerBuilder) SetSchemeOperatorNames(schemeOperatorNames map[string][]string) *OtherListPointerBuilder {
	b.schemeOperatorNames = schemeOperatorNames
	return b
}

// SchemeTypeCommunityRules gets a map of scheme type community rules.
func (b *OtherListPointerBuilder) SchemeTypeCommunityRules() map[string][]string {
	return b.schemeTypeCommunityRules
}

// SetSchemeTypeCommunityRules sets a map of scheme type community rules.
func (b *OtherListPointerBuilder) SetSchemeTypeCommunityRules(schemeTypeCommunityRules map[string][]string) *OtherListPointerBuilder {
	b.schemeTypeCommunityRules = schemeTypeCommunityRules
	return b
}

// Mra gets the MRA (Mutual Recognition Agreement) scheme.
func (b *OtherListPointerBuilder) Mra() *tsl.MRA {
	return b.mra
}

// SetMra sets the MRA (Mutual Recognition Agreement) scheme.
func (b *OtherListPointerBuilder) SetMra(mra *tsl.MRA) *OtherListPointerBuilder {
	b.mra = mra
	return b
}

// Build builds the OtherListPointer.
func (b *OtherListPointerBuilder) Build() *OtherListPointer {
	return NewOtherListPointerFromBuilder(b)
}
