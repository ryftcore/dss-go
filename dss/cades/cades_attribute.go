// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESAttribute.java (DSS 6.5.RC1).
//
// org.bouncycastle.asn1.tsp.EvidenceRecord has no Go port yet (no existing DSS class mirrors
// it - it is a raw RFC 4998 ASN.1 structure, not the spi/validation.EvidenceRecord model type).
// ToEvidenceRecord below therefore returns the generic parsed *asn1ber.Element for the
// attribute value, which callers can walk field-by-field until a typed EvidenceRecord lands.
package cades

import (
	"encoding/asn1"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation/identifier"
)

// CAdESAttribute represents a CAdES attribute, part of AttributeTable. Port of the class
// CAdESAttribute, implementing spi/validation.SignatureAttribute.
type CAdESAttribute struct {
	// attribute is the Attribute value.
	attribute *cmscore.Attribute
	// order of the attribute within signature properties.
	order *int
	// identifier caches Identifier(); mirrors the protected identifier field.
	identifier *CAdESAttributeIdentifier
}

// NewCAdESAttribute is the port of the package-private CAdESAttribute(Attribute, Integer)
// constructor.
func NewCAdESAttribute(attribute *cmscore.Attribute, order *int) *CAdESAttribute {
	return &CAdESAttribute{attribute: attribute, order: order}
}

// ASN1Oid returns the object identifier. Port of getASN1Oid().
func (a *CAdESAttribute) ASN1Oid() asn1.ObjectIdentifier {
	return a.attribute.Type
}

// AttrValues returns the attribute values set. Port of getAttrValues().
func (a *CAdESAttribute) AttrValues() []*asn1ber.Element {
	return a.attribute.Values
}

// Attribute returns the attribute. Port of getAttribute().
func (a *CAdESAttribute) Attribute() *cmscore.Attribute {
	return a.attribute
}

// ASN1Object returns the inner ASN1Encodable object. Port of getASN1Object().
func (a *CAdESAttribute) ASN1Object() *asn1ber.Element {
	return spi.DSSASN1UtilsAsn1Encodable(a.attribute)
}

// IsTimeStampToken checks if the given CAdESAttribute is a timestamp token.
// Port of isTimeStampToken().
func (a *CAdESAttribute) IsTimeStampToken() bool {
	for _, oid := range CAdESUtilsTimestampOids() {
		if oid.Equal(a.ASN1Oid()) {
			return true
		}
	}
	return false
}

// TimestampTokenType returns type of the timestamp token, when applicable.
// Port of getTimestampTokenType().
func (a *CAdESAttribute) TimestampTokenType() enumerations.TimestampType {
	if a.IsTimeStampToken() {
		return CAdESUtilsTimestampTypeByOid(a.ASN1Oid())
	}
	return ""
}

// ToTimeStampToken returns a TimeStampToken if possible, nil otherwise. Port of
// toTimeStampToken().
func (a *CAdESAttribute) ToTimeStampToken() *cmscore.TimeStampToken {
	if a.IsTimeStampToken() {
		token := CAdESUtilsTimeStampToken(a.attribute)
		if token != nil {
			return token
		}
		// Upstream logs "Unable to build a timestamp token from the attribute [{}] : {}".
	} else {
		// Upstream logs "The given attribute [{}] is not a timestamp! Unable to build a
		// TimeStampToken.".
	}
	return nil
}

// IsEvidenceRecord checks if the given CAdES attribute represents an evidence record.
// Port of isEvidenceRecord().
func (a *CAdESAttribute) IsEvidenceRecord() bool {
	for _, oid := range CAdESUtilsEvidenceRecordOids() {
		if oid.Equal(a.ASN1Oid()) {
			return true
		}
	}
	return false
}

// ToEvidenceRecord returns an EvidenceRecord if possible, nil otherwise. Port of
// toEvidenceRecord(). See the file header on the returned type.
func (a *CAdESAttribute) ToEvidenceRecord() *asn1ber.Element {
	if a.IsEvidenceRecord() {
		object := a.ASN1Object()
		if object != nil {
			return object
		}
		// Upstream logs "Unable to build an evidence record from the attribute [{}] : {}".
	} else {
		// Upstream logs "The given attribute [{}] is not an evidence record! Unable to build an
		// EvidenceRecord.".
	}
	return nil
}

// Order gets order of the CAdES Attribute from the original AttributeTable. Port of the
// protected getOrder().
func (a *CAdESAttribute) Order() *int {
	return a.order
}

// Identifier gets the attribute identifier. Port of getIdentifier(), implementing
// spi/validation.SignatureAttribute.
func (a *CAdESAttribute) Identifier() identifier.SignatureAttributeIdentifier {
	if a.identifier == nil {
		a.identifier = CAdESAttributeIdentifierBuild(a.attribute, a.order)
	}
	return a.identifier.SignatureAttributeIdentifier
}

// String is the port of toString().
func (a *CAdESAttribute) String() string {
	oid := a.ASN1Oid()
	if oid != nil {
		return oid.String()
	}
	return ""
}

// Equals is the port of equals(Object), which CAdESAttribute compares by Identifier.
func (a *CAdESAttribute) Equals(other *CAdESAttribute) bool {
	if a == other {
		return true
	}
	if other == nil {
		return false
	}
	selfID := a.Identifier()
	otherID := other.Identifier()
	return selfID.Equals(&otherID)
}
