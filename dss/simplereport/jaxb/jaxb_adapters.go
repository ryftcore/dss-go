// Ported from the generated JAXB adapters Adapter1..Adapter8 of
// eu.europa.esig.dss.simplereport.jaxb, which delegate to the
// eu.europa.esig.dss.jaxb.parsers.*Parser classes. Each adapter becomes a
// named Go type over the corresponding enumerations constant so that
// encoding/xml can marshal it through encoding.TextMarshaler with the exact
// lexical form the parser prints, and reject any other lexical form on the
// way in. Pattern and naming convention match dss/diagnostic/jaxb's
// jaxb_adapters.go.
package jaxb

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// ASiCContainerTypeValue is the ASiCContainerType adapter (Adapter1): the
// parser prints toString(), which replaces '_' with '-'.
type ASiCContainerTypeValue enumerations.ASiCContainerType

// ASiCContainerType returns the underlying enumeration constant.
func (v ASiCContainerTypeValue) ASiCContainerType() enumerations.ASiCContainerType {
	return enumerations.ASiCContainerType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v ASiCContainerTypeValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.ASiCContainerType(v).String()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *ASiCContainerTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.ASiCContainerTypeValues() {
		if enumerations.ASiCContainerType(candidate).String() == s {
			*v = ASiCContainerTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant ASiCContainerType.%s", s)
}

// SignatureScopeTypeValue is the SignatureScopeType adapter (Adapter2):
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

// IndicationValue is the Indication adapter (Adapter3): IndicationParser
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

// SubIndicationValue is the SubIndication adapter (Adapter4):
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

// SignatureQualificationValue is the SignatureQualification adapter
// (Adapter5): SignatureQualificationParser prints/parses getReadable().
type SignatureQualificationValue enumerations.SignatureQualification

// SignatureQualification returns the underlying enumeration constant.
func (v SignatureQualificationValue) SignatureQualification() enumerations.SignatureQualification {
	return enumerations.SignatureQualification(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v SignatureQualificationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.SignatureQualification(v).Readable()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *SignatureQualificationValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.SignatureQualificationValues() {
		if enumerations.SignatureQualification(candidate).Readable() == s {
			*v = SignatureQualificationValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant SignatureQualification.%s", s)
}

// SignatureLevelValue is the SignatureFormat adapter (Adapter6): the
// underlying Java type is enumerations.SignatureLevel, and
// SignatureFormatParser prints/parses SignatureLevel.toString(), which
// replaces '_' with '-' (SignatureLevel.String() in this port).
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

// TimestampQualificationValue is the TimestampQualification adapter
// (Adapter7): TimestampQualificationParser prints/parses getReadable().
type TimestampQualificationValue enumerations.TimestampQualification

// TimestampQualification returns the underlying enumeration constant.
func (v TimestampQualificationValue) TimestampQualification() enumerations.TimestampQualification {
	return enumerations.TimestampQualification(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v TimestampQualificationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.TimestampQualification(v).Readable()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *TimestampQualificationValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.TimestampQualificationValues() {
		if enumerations.TimestampQualification(candidate).Readable() == s {
			*v = TimestampQualificationValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant TimestampQualification.%s", s)
}

// EAAQualificationValue is the EAAQualification adapter (Adapter8):
// EAAQualificationParser prints/parses getReadable().
type EAAQualificationValue enumerations.EAAQualification

// EAAQualification returns the underlying enumeration constant.
func (v EAAQualificationValue) EAAQualification() enumerations.EAAQualification {
	return enumerations.EAAQualification(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v EAAQualificationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.EAAQualification(v).Readable()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EAAQualificationValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EAAQualificationValues() {
		if enumerations.EAAQualification(candidate).Readable() == s {
			*v = EAAQualificationValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EAAQualification.%s", s)
}
