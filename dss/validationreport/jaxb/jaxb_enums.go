// Ported from:
//   - specs-validation-report/src/main/java/eu/europa/esig/validationreport/enums/ObjectType.java
//   - .../enums/ConstraintStatus.java
//   - .../enums/TypeOfProof.java
//   - .../enums/SignatureValidationProcessID.java
//   - .../parsers/UriBasedEnumParser.java
//
// (DSS 6.5.RC1).
//
// # Why these hand-written classes live in jaxb, not dss/validationreport
//
// Package layout places enums/ and parsers/ in the parent
// dss/validationreport package (flattened with the root type and the
// facade). But every ConstraintStatusType/POEType/SignatureValidationProcessType/
// ValidationObjectType field these enums back is a generated model field
// here in jaxb, and Facade (in dss/validationreport) already
// imports jaxb for ValidationReportType - so dss/validationreport -> jaxb is
// a fixed edge, and jaxb cannot also import dss/validationreport for the
// enum types without a Go import cycle Java's classpath never had to
// contend with. This is the same shape of problem PORTING.md's internal/
// extraction solves with a type alias, applied in the other direction: the
// concrete types and the parser logic live here where the generated model
// can reach them, and dss/validationreport/object_type.go (etc.) re-exports
// each one as a type alias plus forwarding functions, so callers still get
// working, identical types.
//
// # URI-adapted dss/enumerations wrappers
//
// The bindings.xml customization file routes MainIndication, SubIndication
// and RevocationReason - all already-ported dss/enumerations types - through
// UriBasedEnumParser too (Adapters 4, 7, 8), serializing them by VR URI
// rather than by Java enum name. dss/enumerations.Indication and friends
// cannot gain a MarshalText method from this package (Go permits methods
// only in the type's own package), so URIIndication/URISubIndication/
// URIRevocationReason below are defined types over them purely to carry the
// xs:anyURI text codec; ValueEndorsementType does the same for
// SAOneSignerRoleType's EndorsementType, adapted by value() (EndorsementTypeParser),
// not by URI.
package jaxb

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// ---------------------------------------------------------------- ObjectType

// ObjectType defines object types. Port of enums.ObjectType. Implements
// enumerations.UriBasedEnum.
type ObjectType string

const (
	// ObjectTypeCertificate is a certificate.
	ObjectTypeCertificate ObjectType = "CERTIFICATE"
	// ObjectTypeCRL is a CRL.
	ObjectTypeCRL ObjectType = "CRL"
	// ObjectTypeOCSPResponse is an OCSP response.
	ObjectTypeOCSPResponse ObjectType = "OCSP_RESPONSE"
	// ObjectTypeTimestamp is a TimeStamp.
	ObjectTypeTimestamp ObjectType = "TIMESTAMP"
	// ObjectTypeEvidenceRecord is an evidence record.
	ObjectTypeEvidenceRecord ObjectType = "EVIDENCE_RECORD"
	// ObjectTypePublicKey is a public key.
	ObjectTypePublicKey ObjectType = "PUBLIC_KEY"
	// ObjectTypeSignedData is signed data.
	ObjectTypeSignedData ObjectType = "SIGNED_DATA"
	// ObjectTypeOther is other.
	ObjectTypeOther ObjectType = "OTHER"
)

var objectTypeURI = map[ObjectType]string{
	ObjectTypeCertificate:    "urn:etsi:019102:validationObject:certificate",
	ObjectTypeCRL:            "urn:etsi:019102:validationObject:CRL",
	ObjectTypeOCSPResponse:   "urn:etsi:019102:validationObject:OCSPResponse",
	ObjectTypeTimestamp:      "urn:etsi:019102:validationObject:timestamp",
	ObjectTypeEvidenceRecord: "urn:etsi:019102:validationObject:evidencerecord",
	ObjectTypePublicKey:      "urn:etsi:019102:validationObject:publicKey",
	ObjectTypeSignedData:     "urn:etsi:019102:validationObject:signedData",
	ObjectTypeOther:          "urn:etsi:019102:validationObject:other",
}

// ObjectTypeValues returns all ObjectType constants in declaration order.
func ObjectTypeValues() []ObjectType {
	return []ObjectType{
		ObjectTypeCertificate, ObjectTypeCRL, ObjectTypeOCSPResponse, ObjectTypeTimestamp,
		ObjectTypeEvidenceRecord, ObjectTypePublicKey, ObjectTypeSignedData, ObjectTypeOther,
	}
}

// ObjectTypeValueOf returns the ObjectType matching the given Java enum name.
func ObjectTypeValueOf(name string) (ObjectType, error) {
	for _, v := range ObjectTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant ObjectType.%s", name)
}

// URI returns the VR URI of the object type. Implements enumerations.UriBasedEnum.
func (o ObjectType) URI() string { return objectTypeURI[o] }

// MarshalText writes the object type's VR URI, as the ObjectType JAXB
// adapter's marshal() does.
func (o ObjectType) MarshalText() ([]byte, error) { return []byte(o.URI()), nil }

// UnmarshalText resolves the VR URI, mirroring UriBasedEnumParser.parseObjectType.
func (o *ObjectType) UnmarshalText(text []byte) error {
	*o = ParseObjectType(string(text))
	return nil
}

// ----------------------------------------------------------- ConstraintStatus

// ConstraintStatus defines the ConstraintStatus type. Port of
// enums.ConstraintStatus. Implements enumerations.UriBasedEnum.
type ConstraintStatus string

const (
	// ConstraintStatusApplied: the constraint has been applied.
	ConstraintStatusApplied ConstraintStatus = "APPLIED"
	// ConstraintStatusDisabled: the constraint has been disabled.
	ConstraintStatusDisabled ConstraintStatus = "DISABLED"
	// ConstraintStatusOverridden: the constraint has been overridden.
	ConstraintStatusOverridden ConstraintStatus = "OVERRIDDEN"
)

var constraintStatusURI = map[ConstraintStatus]string{
	ConstraintStatusApplied:    "urn:etsi:019102:constraintStatus:applied",
	ConstraintStatusDisabled:   "urn:etsi:019102:constraintStatus:disabled",
	ConstraintStatusOverridden: "urn:etsi:019102:constraintStatus:overridden",
}

// ConstraintStatusValues returns all ConstraintStatus constants in declaration order.
func ConstraintStatusValues() []ConstraintStatus {
	return []ConstraintStatus{ConstraintStatusApplied, ConstraintStatusDisabled, ConstraintStatusOverridden}
}

// ConstraintStatusValueOf returns the ConstraintStatus matching the given Java enum name.
func ConstraintStatusValueOf(name string) (ConstraintStatus, error) {
	for _, v := range ConstraintStatusValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant ConstraintStatus.%s", name)
}

// URI returns the VR URI of the constraint status. Implements enumerations.UriBasedEnum.
func (c ConstraintStatus) URI() string { return constraintStatusURI[c] }

// MarshalText writes the constraint status's VR URI.
func (c ConstraintStatus) MarshalText() ([]byte, error) { return []byte(c.URI()), nil }

// UnmarshalText resolves the VR URI, mirroring UriBasedEnumParser.parseConstraintStatus.
func (c *ConstraintStatus) UnmarshalText(text []byte) error {
	*c = ParseConstraintStatus(string(text))
	return nil
}

// -------------------------------------------------------------- TypeOfProof

// TypeOfProof defines a TypeOfProof. Port of enums.TypeOfProof. Implements
// enumerations.UriBasedEnum.
type TypeOfProof string

const (
	// TypeOfProofValidation: the POE has been derived during validation.
	TypeOfProofValidation TypeOfProof = "VALIDATION"
	// TypeOfProofProvided: the POE has been provided to the SVA as an input.
	TypeOfProofProvided TypeOfProof = "PROVIDED"
	// TypeOfProofPolicy: the POE has been derived by the policy.
	TypeOfProofPolicy TypeOfProof = "POLICY"
)

var typeOfProofURI = map[TypeOfProof]string{
	TypeOfProofValidation: "urn:etsi:019102:poetype:validation",
	TypeOfProofProvided:   "urn:etsi:019102:poetype:provided",
	TypeOfProofPolicy:     "urn:etsi:019102:poetype:policy",
}

// TypeOfProofValues returns all TypeOfProof constants in declaration order.
func TypeOfProofValues() []TypeOfProof {
	return []TypeOfProof{TypeOfProofValidation, TypeOfProofProvided, TypeOfProofPolicy}
}

// TypeOfProofValueOf returns the TypeOfProof matching the given Java enum name.
func TypeOfProofValueOf(name string) (TypeOfProof, error) {
	for _, v := range TypeOfProofValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant TypeOfProof.%s", name)
}

// URI returns the VR URI of the proof type. Implements enumerations.UriBasedEnum.
func (t TypeOfProof) URI() string { return typeOfProofURI[t] }

// MarshalText writes the proof type's VR URI.
func (t TypeOfProof) MarshalText() ([]byte, error) { return []byte(t.URI()), nil }

// UnmarshalText resolves the VR URI, mirroring UriBasedEnumParser.parseTypeOfProof.
func (t *TypeOfProof) UnmarshalText(text []byte) error {
	*t = ParseTypeOfProof(string(text))
	return nil
}

// ------------------------------------------------------ SignatureValidationProcessID

// SignatureValidationProcessID defines SignatureValidationProcessID. Port of
// enums.SignatureValidationProcessID. Implements enumerations.UriBasedEnum.
type SignatureValidationProcessID string

const (
	// SignatureValidationProcessIDBasic: the SVA performed the Validation
	// Process for Basic Signatures (ETSI TS 119 102-1, clause 5.3).
	SignatureValidationProcessIDBasic SignatureValidationProcessID = "BASIC"
	// SignatureValidationProcessIDLTVM: the SVA performed the Validation
	// Process for Signatures with Time and LongTerm-Validation Material
	// (ETSI TS 119 102-1, clause 5.5).
	SignatureValidationProcessIDLTVM SignatureValidationProcessID = "LTVM"
	// SignatureValidationProcessIDLTA: the SVA performed the Validation
	// process for Signatures providing Long Term Availability and Integrity
	// of Validation Material (ETSI TS 119 102-1, clause 5.6).
	SignatureValidationProcessIDLTA SignatureValidationProcessID = "LTA"
)

var signatureValidationProcessIDURI = map[SignatureValidationProcessID]string{
	SignatureValidationProcessIDBasic: "urn:etsi:019102:validationprocess:Basic",
	SignatureValidationProcessIDLTVM:  "urn:etsi:019102:validationprocess:LTVM",
	SignatureValidationProcessIDLTA:   "urn:etsi:019102:validationprocess:LTA",
}

// SignatureValidationProcessIDValues returns all SignatureValidationProcessID
// constants in declaration order.
func SignatureValidationProcessIDValues() []SignatureValidationProcessID {
	return []SignatureValidationProcessID{
		SignatureValidationProcessIDBasic, SignatureValidationProcessIDLTVM, SignatureValidationProcessIDLTA,
	}
}

// SignatureValidationProcessIDValueOf returns the SignatureValidationProcessID
// matching the given Java enum name.
func SignatureValidationProcessIDValueOf(name string) (SignatureValidationProcessID, error) {
	for _, v := range SignatureValidationProcessIDValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant SignatureValidationProcessID.%s", name)
}

// URI returns the VR URI of the validation process ID. Implements enumerations.UriBasedEnum.
func (s SignatureValidationProcessID) URI() string { return signatureValidationProcessIDURI[s] }

// MarshalText writes the validation process ID's VR URI.
func (s SignatureValidationProcessID) MarshalText() ([]byte, error) { return []byte(s.URI()), nil }

// UnmarshalText resolves the VR URI, mirroring
// UriBasedEnumParser.parseSignatureValidationProcessID.
func (s *SignatureValidationProcessID) UnmarshalText(text []byte) error {
	*s = ParseSignatureValidationProcessID(string(text))
	return nil
}

// ---------------------------------------------------------- UriBasedEnumParser

// uriToEnum is the map of enum values and corresponding URIs. Port of
// UriBasedEnumParser.URI_TO_ENUM_MAP, populated the same way (Indication,
// ObjectType, RevocationReason, SignatureValidationProcessID, SubIndication,
// TypeOfProof, ConstraintStatus, in that order - no collisions across
// families, so the order has no observable effect, but it is preserved for
// faithfulness).
var uriToEnum = map[string]enumerations.UriBasedEnum{}

func init() {
	registerURIEnum(indicationValuesAsUriBasedEnum())
	registerURIEnum(objectTypeValuesAsUriBasedEnum())
	registerURIEnum(revocationReasonValuesAsUriBasedEnum())
	registerURIEnum(signatureValidationProcessIDValuesAsUriBasedEnum())
	registerURIEnum(subIndicationValuesAsUriBasedEnum())
	registerURIEnum(typeOfProofValuesAsUriBasedEnum())
	registerURIEnum(constraintStatusValuesAsUriBasedEnum())
}

func registerURIEnum(values []enumerations.UriBasedEnum) {
	for _, v := range values {
		uriToEnum[v.URI()] = v
	}
}

func indicationValuesAsUriBasedEnum() []enumerations.UriBasedEnum {
	vs := enumerations.IndicationValues()
	out := make([]enumerations.UriBasedEnum, len(vs))
	for i, v := range vs {
		out[i] = v
	}
	return out
}

func subIndicationValuesAsUriBasedEnum() []enumerations.UriBasedEnum {
	vs := enumerations.SubIndicationValues()
	out := make([]enumerations.UriBasedEnum, len(vs))
	for i, v := range vs {
		out[i] = v
	}
	return out
}

func revocationReasonValuesAsUriBasedEnum() []enumerations.UriBasedEnum {
	vs := enumerations.RevocationReasonValues()
	out := make([]enumerations.UriBasedEnum, len(vs))
	for i, v := range vs {
		out[i] = v
	}
	return out
}

func objectTypeValuesAsUriBasedEnum() []enumerations.UriBasedEnum {
	vs := ObjectTypeValues()
	out := make([]enumerations.UriBasedEnum, len(vs))
	for i, v := range vs {
		out[i] = v
	}
	return out
}

func signatureValidationProcessIDValuesAsUriBasedEnum() []enumerations.UriBasedEnum {
	vs := SignatureValidationProcessIDValues()
	out := make([]enumerations.UriBasedEnum, len(vs))
	for i, v := range vs {
		out[i] = v
	}
	return out
}

func typeOfProofValuesAsUriBasedEnum() []enumerations.UriBasedEnum {
	vs := TypeOfProofValues()
	out := make([]enumerations.UriBasedEnum, len(vs))
	for i, v := range vs {
		out[i] = v
	}
	return out
}

func constraintStatusValuesAsUriBasedEnum() []enumerations.UriBasedEnum {
	vs := ConstraintStatusValues()
	out := make([]enumerations.UriBasedEnum, len(vs))
	for i, v := range vs {
		out[i] = v
	}
	return out
}

// parse looks the URI up in the registry, mirroring UriBasedEnumParser.parse:
// a miss returns nil, never an error.
func parse(v string) enumerations.UriBasedEnum {
	return uriToEnum[v]
}

// ParseMainIndication parses the URI and returns the matching
// enumerations.Indication, or the zero value if none matches. Port of
// UriBasedEnumParser.parseMainIndication.
func ParseMainIndication(v string) enumerations.Indication {
	if e, ok := parse(v).(enumerations.Indication); ok {
		return e
	}
	return ""
}

// ParseSubIndication parses the URI and returns the matching
// enumerations.SubIndication, or the zero value if none matches. Port of
// UriBasedEnumParser.parseSubIndication.
func ParseSubIndication(v string) enumerations.SubIndication {
	if e, ok := parse(v).(enumerations.SubIndication); ok {
		return e
	}
	return ""
}

// ParseObjectType parses the URI and returns the matching ObjectType, or the
// zero value if none matches. Port of UriBasedEnumParser.parseObjectType.
func ParseObjectType(v string) ObjectType {
	if e, ok := parse(v).(ObjectType); ok {
		return e
	}
	return ""
}

// ParseRevocationReason parses the URI and returns the matching
// enumerations.RevocationReason, or the zero value if none matches. Port of
// UriBasedEnumParser.parseRevocationReason.
func ParseRevocationReason(v string) enumerations.RevocationReason {
	if e, ok := parse(v).(enumerations.RevocationReason); ok {
		return e
	}
	return ""
}

// ParseSignatureValidationProcessID parses the URI and returns the matching
// SignatureValidationProcessID, or the zero value if none matches. Port of
// UriBasedEnumParser.parseSignatureValidationProcessID.
func ParseSignatureValidationProcessID(v string) SignatureValidationProcessID {
	if e, ok := parse(v).(SignatureValidationProcessID); ok {
		return e
	}
	return ""
}

// ParseTypeOfProof parses the URI and returns the matching TypeOfProof, or
// the zero value if none matches. Port of UriBasedEnumParser.parseTypeOfProof.
func ParseTypeOfProof(v string) TypeOfProof {
	if e, ok := parse(v).(TypeOfProof); ok {
		return e
	}
	return ""
}

// ParseConstraintStatus parses the URI and returns the matching
// ConstraintStatus, or the zero value if none matches. Port of
// UriBasedEnumParser.parseConstraintStatus.
func ParseConstraintStatus(v string) ConstraintStatus {
	if e, ok := parse(v).(ConstraintStatus); ok {
		return e
	}
	return ""
}

// Print returns v's URI, or "" if v is nil. Port of UriBasedEnumParser.print.
func Print(v enumerations.UriBasedEnum) string {
	if v == nil {
		return ""
	}
	return v.URI()
}

// ------------------------------------------------ URI-adapted enumerations wrappers

// URIIndication adapts enumerations.Indication to the xs:anyURI lexical form
// Adapter7/UriBasedEnumParser.parseMainIndication bind it through (see this
// file's header).
type URIIndication enumerations.Indication

// MarshalText writes the indication's VR URI.
func (v URIIndication) MarshalText() ([]byte, error) {
	return []byte(enumerations.Indication(v).URI()), nil
}

// UnmarshalText resolves the VR URI.
func (v *URIIndication) UnmarshalText(text []byte) error {
	*v = URIIndication(ParseMainIndication(string(text)))
	return nil
}

// Indication returns the wrapped enumerations.Indication.
func (v URIIndication) Indication() enumerations.Indication { return enumerations.Indication(v) }

// URISubIndication adapts enumerations.SubIndication to the xs:anyURI
// lexical form Adapter8/UriBasedEnumParser.parseSubIndication bind it
// through (see this file's header).
type URISubIndication enumerations.SubIndication

// MarshalText writes the sub-indication's VR URI.
func (v URISubIndication) MarshalText() ([]byte, error) {
	return []byte(enumerations.SubIndication(v).URI()), nil
}

// UnmarshalText resolves the VR URI.
func (v *URISubIndication) UnmarshalText(text []byte) error {
	*v = URISubIndication(ParseSubIndication(string(text)))
	return nil
}

// SubIndication returns the wrapped enumerations.SubIndication.
func (v URISubIndication) SubIndication() enumerations.SubIndication {
	return enumerations.SubIndication(v)
}

// URIRevocationReason adapts enumerations.RevocationReason to the xs:anyURI
// lexical form Adapter4/UriBasedEnumParser.parseRevocationReason bind it
// through (see this file's header).
type URIRevocationReason enumerations.RevocationReason

// MarshalText writes the revocation reason's VR URI.
func (v URIRevocationReason) MarshalText() ([]byte, error) {
	return []byte(enumerations.RevocationReason(v).URI()), nil
}

// UnmarshalText resolves the VR URI.
func (v *URIRevocationReason) UnmarshalText(text []byte) error {
	*v = URIRevocationReason(ParseRevocationReason(string(text)))
	return nil
}

// RevocationReason returns the wrapped enumerations.RevocationReason.
func (v URIRevocationReason) RevocationReason() enumerations.RevocationReason {
	return enumerations.RevocationReason(v)
}

// ValueEndorsementType adapts enumerations.EndorsementType to the lowercase
// lexical form Adapter1/eu.europa.esig.dss.jaxb.parsers.EndorsementTypeParser
// bind it through - by value(), not by VR URI (EndorsementType does not
// implement UriBasedEnum upstream either).
type ValueEndorsementType enumerations.EndorsementType

// MarshalText writes the endorsement type's lowercase value.
func (v ValueEndorsementType) MarshalText() ([]byte, error) {
	return []byte(enumerations.EndorsementType(v).Value()), nil
}

// UnmarshalText resolves the lowercase value, mirroring EndorsementTypeParser.parse.
func (v *ValueEndorsementType) UnmarshalText(text []byte) error {
	*v = ValueEndorsementType(enumerations.EndorsementTypeFromString(string(text)))
	return nil
}

// EndorsementType returns the wrapped enumerations.EndorsementType.
func (v ValueEndorsementType) EndorsementType() enumerations.EndorsementType {
	return enumerations.EndorsementType(v)
}
