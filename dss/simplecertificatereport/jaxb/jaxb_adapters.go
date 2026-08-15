// Ported from the generated JAXB adapters Adapter1..Adapter11 of
// eu.europa.esig.dss.simplecertificatereport.jaxb, which delegate to the
// eu.europa.esig.dss.jaxb.parsers.*Parser classes. Each adapter becomes a
// named Go type over the corresponding enumerations constant so that
// encoding/xml can marshal it through encoding.TextMarshaler with the exact
// lexical form the parser prints, and reject any other lexical form on the
// way in. Pattern and naming convention match dss/diagnostic/jaxb's and
// dss/simplereport/jaxb's jaxb_adapters.go.
//
// ListType, LoTEServiceTypeIdentifier and LoTEServiceStatus
// (Adapter9/10/11) are Go interfaces, not named string types (see
// dss/enumerations/list_type.go and friends: they are resolved through a
// pluggable LoTELoader registry, not a closed constant set), so - unlike
// every other adapter here - they cannot be redeclared as a defined type
// with methods (Go forbids methods on a named interface type). Each is
// instead a one-field struct wrapping the interface value.
package jaxb

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
)

// KeyUsageBitValue is the KeyUsageBit adapter (Adapter1): KeyUsageBitParser
// prints/parses getValue().
type KeyUsageBitValue enumerations.KeyUsageBit

// KeyUsageBit returns the underlying enumeration constant.
func (v KeyUsageBitValue) KeyUsageBit() enumerations.KeyUsageBit { return enumerations.KeyUsageBit(v) }

// MarshalText writes the lexical form produced by the Java parser.
func (v KeyUsageBitValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.KeyUsageBit(v).Value()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *KeyUsageBitValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.KeyUsageBitValues() {
		if enumerations.KeyUsageBit(candidate).Value() == s {
			*v = KeyUsageBitValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant KeyUsageBit.%s", s)
}

// RevocationReasonValue is the RevocationReason adapter (Adapter2):
// RevocationReasonParser prints/parses the short name.
type RevocationReasonValue enumerations.RevocationReason

// RevocationReason returns the underlying enumeration constant.
func (v RevocationReasonValue) RevocationReason() enumerations.RevocationReason {
	return enumerations.RevocationReason(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v RevocationReasonValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.RevocationReason(v).ShortName()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *RevocationReasonValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.RevocationReasonValues() {
		if enumerations.RevocationReason(candidate).ShortName() == s {
			*v = RevocationReasonValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant RevocationReason.%s", s)
}

// CertificateQualificationValue is the CertificateQualification adapter
// (Adapter3): CertificateQualificationParser prints/parses getReadable().
type CertificateQualificationValue enumerations.CertificateQualification

// CertificateQualification returns the underlying enumeration constant.
func (v CertificateQualificationValue) CertificateQualification() enumerations.CertificateQualification {
	return enumerations.CertificateQualification(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v CertificateQualificationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.CertificateQualification(v).Readable()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *CertificateQualificationValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.CertificateQualificationValues() {
		if enumerations.CertificateQualification(candidate).Readable() == s {
			*v = CertificateQualificationValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant CertificateQualification.%s", s)
}

// QWACProfileValue is the QWACProfile adapter (Adapter4): QWACProfileParser
// prints/parses getReadable().
type QWACProfileValue enumerations.QWACProfile

// QWACProfile returns the underlying enumeration constant.
func (v QWACProfileValue) QWACProfile() enumerations.QWACProfile { return enumerations.QWACProfile(v) }

// MarshalText writes the lexical form produced by the Java parser.
func (v QWACProfileValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.QWACProfile(v).Readable()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *QWACProfileValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.QWACProfileValues() {
		if enumerations.QWACProfile(candidate).Readable() == s {
			*v = QWACProfileValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant QWACProfile.%s", s)
}

// IndicationValue is the Indication adapter (Adapter5): IndicationParser
// prints/parses the Java enum name() verbatim.
type IndicationValue enumerations.Indication

// Indication returns the underlying enumeration constant.
func (v IndicationValue) Indication() enumerations.Indication { return enumerations.Indication(v) }

// MarshalText writes the lexical form produced by the Java parser.
func (v IndicationValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *IndicationValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.IndicationValues() {
		if string(candidate) == s {
			*v = IndicationValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant Indication.%s", s)
}

// SubIndicationValue is the SubIndication adapter (Adapter6):
// SubIndicationParser prints/parses the Java enum name() verbatim.
type SubIndicationValue enumerations.SubIndication

// SubIndication returns the underlying enumeration constant.
func (v SubIndicationValue) SubIndication() enumerations.SubIndication {
	return enumerations.SubIndication(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v SubIndicationValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *SubIndicationValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.SubIndicationValues() {
		if string(candidate) == s {
			*v = SubIndicationValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant SubIndication.%s", s)
}

// SignatureLevelValue is the SignatureFormat adapter (Adapter7): the
// underlying Java type is enumerations.SignatureLevel, and
// SignatureFormatParser prints/parses SignatureLevel.toString() (String()
// in this port), which replaces '_' with '-'.
type SignatureLevelValue enumerations.SignatureLevel

// SignatureLevel returns the underlying enumeration constant.
func (v SignatureLevelValue) SignatureLevel() enumerations.SignatureLevel {
	return enumerations.SignatureLevel(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v SignatureLevelValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.SignatureLevel(v).String()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *SignatureLevelValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.SignatureLevelValues() {
		if enumerations.SignatureLevel(candidate).String() == s {
			*v = SignatureLevelValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant SignatureLevel.%s", s)
}

// SignatureScopeTypeValue is the SignatureScopeType adapter (Adapter8):
// SignatureScopeTypeParser prints/parses the Java enum name() verbatim.
type SignatureScopeTypeValue enumerations.SignatureScopeType

// SignatureScopeType returns the underlying enumeration constant.
func (v SignatureScopeTypeValue) SignatureScopeType() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v SignatureScopeTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *SignatureScopeTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.SignatureScopeTypeValues() {
		if string(candidate) == s {
			*v = SignatureScopeTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant SignatureScopeType.%s", s)
}

// ListTypeValue is the ListType adapter (Adapter9): ListTypeParser
// prints/parses the URI. See the file header for why this is a struct
// rather than a defined type over enumerations.ListType.
type ListTypeValue struct {
	Value enumerations.ListType
}

// MarshalText writes the lexical form produced by the Java parser.
func (v ListTypeValue) MarshalText() ([]byte, error) {
	if v.Value == nil {
		return nil, nil
	}
	return []byte(v.Value.URI()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *ListTypeValue) UnmarshalText(text []byte) error {
	lt := enumerations.ListTypeFromURI(string(text))
	if lt == nil {
		return fmt.Errorf("no ListType for URI %q", text)
	}
	v.Value = lt
	return nil
}

// LoTEServiceTypeIdentifierValue is the LoTEServiceTypeIdentifier adapter
// (Adapter10): LoTEServiceTypeIdentifierParser prints/parses the URI.
type LoTEServiceTypeIdentifierValue struct {
	Value enumerations.LoTEServiceTypeIdentifier
}

// MarshalText writes the lexical form produced by the Java parser.
func (v LoTEServiceTypeIdentifierValue) MarshalText() ([]byte, error) {
	if v.Value == nil {
		return nil, nil
	}
	return []byte(v.Value.URI()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *LoTEServiceTypeIdentifierValue) UnmarshalText(text []byte) error {
	id := enumerations.LoTEServiceTypeIdentifierFromURI(string(text))
	if id == nil {
		return fmt.Errorf("no LoTEServiceTypeIdentifier for URI %q", text)
	}
	v.Value = id
	return nil
}

// LoTEServiceStatusValue is the LoTEServiceStatus adapter (Adapter11):
// LoTEServiceStatusParser prints/parses the URI.
type LoTEServiceStatusValue struct {
	Value enumerations.LoTEServiceStatus
}

// MarshalText writes the lexical form produced by the Java parser.
func (v LoTEServiceStatusValue) MarshalText() ([]byte, error) {
	if v.Value == nil {
		return nil, nil
	}
	return []byte(v.Value.URI()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *LoTEServiceStatusValue) UnmarshalText(text []byte) error {
	st := enumerations.LoTEServiceStatusFromURI(string(text))
	if st == nil {
		return fmt.Errorf("no LoTEServiceStatus for URI %q", text)
	}
	v.Value = st
	return nil
}
