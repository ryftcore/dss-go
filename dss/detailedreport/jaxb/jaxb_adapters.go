// Ported from the generated JAXB adapters Adapter1..Adapter11 of
// eu.europa.esig.dss.detailedreport.jaxb, which delegate to the
// eu.europa.esig.dss.jaxb.parsers.*Parser classes (dss-jaxb-parsers). Each
// adapter becomes a named Go type over the corresponding enumerations constant
// so that encoding/xml can marshal it through encoding.TextMarshaler with the
// exact lexical form the parser prints, and reject any other lexical form on
// the way in - the same pattern dss/diagnostic/jaxb/jaxb_adapters.go
// establishes; see that file for the precedent.
//
// Adapter12 (ListType), Adapter13 (LoTEServiceTypeIdentifier) and Adapter14
// (LoTEServiceStatus) bind to enumerations interfaces backed by a dynamic
// LoTE-loader registry rather than a fixed constant set (ListTypeFromURI,
// LoTEServiceTypeIdentifierFromURI, LoTEServiceStatusFromURI), so they cannot
// be expressed as a defined type over a comparable underlying type the way the
// other eleven can. All three are used exclusively by XmlCertificateApprovalStatus
// (jaxb_qualification.go), which implements MarshalXML/UnmarshalXML by hand
// instead.
//
// Deviation from the Java parsers: every parser here returns nil/null on an
// unrecognised lexical form (e.g. RevocationReasonParser.parseShortName,
// CertificateQualificationParser.parse). Consistent with every other adapter in
// this port (see dss/diagnostic/jaxb/jaxb_adapters.go), UnmarshalText returns an
// error instead of silently producing the zero value.
package jaxb

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
)

// IndicationValue is the Indication adapter: IndicationParser.print() writes name(); parse() is Indication.valueOf().
type IndicationValue enumerations.Indication

// Indication returns the underlying enumeration constant.
func (v IndicationValue) Indication() enumerations.Indication { return enumerations.Indication(v) }

// MarshalText writes the lexical form produced by the Java parser.
func (v IndicationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.Indication(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *IndicationValue) UnmarshalText(text []byte) error {
	r, err := enumerations.IndicationValueOf(string(text))
	if err != nil {
		return fmt.Errorf("no enum constant Indication.%s", text)
	}
	*v = IndicationValue(r)
	return nil
}

// SubIndicationValue is the SubIndication adapter: SubIndicationParser.print() writes name(); parse() is SubIndication.forName().
type SubIndicationValue enumerations.SubIndication

// SubIndication returns the underlying enumeration constant.
func (v SubIndicationValue) SubIndication() enumerations.SubIndication {
	return enumerations.SubIndication(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v SubIndicationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.SubIndication(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *SubIndicationValue) UnmarshalText(text []byte) error {
	r, err := enumerations.SubIndicationForName(string(text))
	if err != nil {
		return fmt.Errorf("no enum constant SubIndication.%s", text)
	}
	*v = SubIndicationValue(r)
	return nil
}

// CertificateSourceTypeValue is the CertificateSourceType adapter: CertificateSourceTypeParser.print() writes name(); parse() is CertificateSourceType.valueOf().
type CertificateSourceTypeValue enumerations.CertificateSourceType

// CertificateSourceType returns the underlying enumeration constant.
func (v CertificateSourceTypeValue) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v CertificateSourceTypeValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.CertificateSourceType(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *CertificateSourceTypeValue) UnmarshalText(text []byte) error {
	r, err := enumerations.CertificateSourceTypeValueOf(string(text))
	if err != nil {
		return fmt.Errorf("no enum constant CertificateSourceType.%s", text)
	}
	*v = CertificateSourceTypeValue(r)
	return nil
}

// RevocationReasonValue is the RevocationReason adapter: RevocationReasonParser prints/parses the CRL short name (RevocationReasonParser.printShortName/parseShortName), not name().
type RevocationReasonValue enumerations.RevocationReason

// RevocationReason returns the underlying enumeration constant.
func (v RevocationReasonValue) RevocationReason() enumerations.RevocationReason {
	return enumerations.RevocationReason(v)
}

// MarshalText writes the short-name lexical form produced by the Java parser.
func (v RevocationReasonValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.RevocationReason(v).ShortName()), nil
}

// UnmarshalText resolves the short-name lexical form back to the enumeration constant.
func (v *RevocationReasonValue) UnmarshalText(text []byte) error {
	r := enumerations.RevocationReasonFromValue(string(text))
	if r == "" {
		return fmt.Errorf("no enum constant RevocationReason with short name %q", text)
	}
	*v = RevocationReasonValue(r)
	return nil
}

// CertificateQualificationValue is the CertificateQualification adapter: CertificateQualificationParser prints/parses the readable label (getReadable/fromReadable), not name().
type CertificateQualificationValue enumerations.CertificateQualification

// CertificateQualification returns the underlying enumeration constant.
func (v CertificateQualificationValue) CertificateQualification() enumerations.CertificateQualification {
	return enumerations.CertificateQualification(v)
}

// MarshalText writes the readable lexical form produced by the Java parser.
func (v CertificateQualificationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.CertificateQualification(v).Readable()), nil
}

// UnmarshalText resolves the readable lexical form back to the enumeration constant.
func (v *CertificateQualificationValue) UnmarshalText(text []byte) error {
	r := enumerations.CertificateQualificationFromReadable(string(text))
	if r == "" {
		return fmt.Errorf("no enum constant CertificateQualification with readable %q", text)
	}
	*v = CertificateQualificationValue(r)
	return nil
}

// SignatureQualificationValue is the SignatureQualification adapter: SignatureQualificationParser prints/parses the readable label (getReadable/fromReadable), not name().
type SignatureQualificationValue enumerations.SignatureQualification

// SignatureQualification returns the underlying enumeration constant.
func (v SignatureQualificationValue) SignatureQualification() enumerations.SignatureQualification {
	return enumerations.SignatureQualification(v)
}

// MarshalText writes the readable lexical form produced by the Java parser.
func (v SignatureQualificationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.SignatureQualification(v).Readable()), nil
}

// UnmarshalText resolves the readable lexical form back to the enumeration constant.
func (v *SignatureQualificationValue) UnmarshalText(text []byte) error {
	r := enumerations.SignatureQualificationFromReadable(string(text))
	if r == "" {
		return fmt.Errorf("no enum constant SignatureQualification with readable %q", text)
	}
	*v = SignatureQualificationValue(r)
	return nil
}

// TimestampQualificationValue is the TimestampQualification adapter: TimestampQualificationParser prints/parses the readable label (getReadable), not name().
type TimestampQualificationValue enumerations.TimestampQualification

// TimestampQualification returns the underlying enumeration constant.
func (v TimestampQualificationValue) TimestampQualification() enumerations.TimestampQualification {
	return enumerations.TimestampQualification(v)
}

// MarshalText writes the readable lexical form produced by the Java parser.
func (v TimestampQualificationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.TimestampQualification(v).Readable()), nil
}

// UnmarshalText resolves the readable lexical form back to the enumeration constant.
func (v *TimestampQualificationValue) UnmarshalText(text []byte) error {
	for _, candidate := range enumerations.TimestampQualificationValues() {
		if candidate.Readable() == string(text) {
			*v = TimestampQualificationValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant TimestampQualification with readable %q", text)
}

// EAAQualificationValue is the EAAQualification adapter: EAAQualificationParser prints/parses the readable label (getReadable/fromReadable), not name().
type EAAQualificationValue enumerations.EAAQualification

// EAAQualification returns the underlying enumeration constant.
func (v EAAQualificationValue) EAAQualification() enumerations.EAAQualification {
	return enumerations.EAAQualification(v)
}

// MarshalText writes the readable lexical form produced by the Java parser.
func (v EAAQualificationValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.EAAQualification(v).Readable()), nil
}

// UnmarshalText resolves the readable lexical form back to the enumeration constant.
func (v *EAAQualificationValue) UnmarshalText(text []byte) error {
	r := enumerations.EAAQualificationFromReadable(string(text))
	if r == "" {
		return fmt.Errorf("no enum constant EAAQualification with readable %q", text)
	}
	*v = EAAQualificationValue(r)
	return nil
}

// ContextValue is the Context adapter: ContextParser.print() writes name() (null-safe); parse() is Context.valueOf() (null-safe).
type ContextValue enumerations.Context

// Context returns the underlying enumeration constant.
func (v ContextValue) Context() enumerations.Context { return enumerations.Context(v) }

// MarshalText writes the lexical form produced by the Java parser.
func (v ContextValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.Context(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *ContextValue) UnmarshalText(text []byte) error {
	for _, candidate := range enumerations.ContextValues() {
		if string(candidate) == string(text) {
			*v = ContextValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant Context.%s", text)
}

// ValidationTimeValue is the ValidationTime adapter: ValidationTimeParser.print() writes name() (null-safe); parse() is ValidationTime.valueOf() (null-safe).
type ValidationTimeValue enumerations.ValidationTime

// ValidationTime returns the underlying enumeration constant.
func (v ValidationTimeValue) ValidationTime() enumerations.ValidationTime {
	return enumerations.ValidationTime(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v ValidationTimeValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.ValidationTime(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *ValidationTimeValue) UnmarshalText(text []byte) error {
	for _, candidate := range enumerations.ValidationTimeValues() {
		if string(candidate) == string(text) {
			*v = ValidationTimeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant ValidationTime.%s", text)
}

// QWACProfileValue is the QWACProfile adapter: QWACProfileParser prints/parses the readable label (getReadable/fromReadable, null-safe), not name().
type QWACProfileValue enumerations.QWACProfile

// QWACProfile returns the underlying enumeration constant.
func (v QWACProfileValue) QWACProfile() enumerations.QWACProfile { return enumerations.QWACProfile(v) }

// MarshalText writes the readable lexical form produced by the Java parser.
func (v QWACProfileValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.QWACProfile(v).Readable()), nil
}

// UnmarshalText resolves the readable lexical form back to the enumeration constant.
func (v *QWACProfileValue) UnmarshalText(text []byte) error {
	r := enumerations.QWACProfileFromReadable(string(text))
	if r == "" {
		return fmt.Errorf("no enum constant QWACProfile with readable %q", text)
	}
	*v = QWACProfileValue(r)
	return nil
}
