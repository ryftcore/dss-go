// Ported from the generated JAXB adapters Adapter1..Adapter39 of
// eu.europa.esig.dss.diagnostic.jaxb, which delegate to the
// eu.europa.esig.dss.jaxb.parsers.*Parser classes. Each adapter becomes a named
// Go type over the corresponding enumerations constant so that encoding/xml can
// marshal it through encoding.TextMarshaler with the exact lexical form the
// parser prints, and reject any other lexical form on the way in.

package jaxb

import (
	"fmt"
	"strings"

	"github.com/utain/esig/dss/enumerations"
)

// ASiCContainerTypeValue is the ASiCContainerType adapter: The parser prints toString(), which replaces '_' with '-'.
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

// ArchiveTimestampHashIndexVersionValue is the ArchiveTimestampHashIndexVersion adapter: The parser prints getLabel() and parses via the label.
type ArchiveTimestampHashIndexVersionValue enumerations.ArchiveTimestampHashIndexVersion

// ArchiveTimestampHashIndexVersion returns the underlying enumeration constant.
func (v ArchiveTimestampHashIndexVersionValue) ArchiveTimestampHashIndexVersion() enumerations.ArchiveTimestampHashIndexVersion {
	return enumerations.ArchiveTimestampHashIndexVersion(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v ArchiveTimestampHashIndexVersionValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.ArchiveTimestampHashIndexVersion(v).Label()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *ArchiveTimestampHashIndexVersionValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.ArchiveTimestampHashIndexVersionValues() {
		if enumerations.ArchiveTimestampHashIndexVersion(candidate).Label() == s {
			*v = ArchiveTimestampHashIndexVersionValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant ArchiveTimestampHashIndexVersion.%s", s)
}

// ArchiveTimestampTypeValue is the ArchiveTimestampType adapter: Adapter print()/parse() map to the Java enum name().
type ArchiveTimestampTypeValue enumerations.ArchiveTimestampType

// ArchiveTimestampType returns the underlying enumeration constant.
func (v ArchiveTimestampTypeValue) ArchiveTimestampType() enumerations.ArchiveTimestampType {
	return enumerations.ArchiveTimestampType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v ArchiveTimestampTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *ArchiveTimestampTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.ArchiveTimestampTypeValues() {
		if string(candidate) == s {
			*v = ArchiveTimestampTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant ArchiveTimestampType.%s", s)
}

// COSESignatureTypeValue is the COSESignatureType adapter: The parser prints getLabel() and parses via the label.
type COSESignatureTypeValue enumerations.COSESignatureType

// COSESignatureType returns the underlying enumeration constant.
func (v COSESignatureTypeValue) COSESignatureType() enumerations.COSESignatureType {
	return enumerations.COSESignatureType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v COSESignatureTypeValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.COSESignatureType(v).Label()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *COSESignatureTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.COSESignatureTypeValues() {
		if enumerations.COSESignatureType(candidate).Label() == s {
			*v = COSESignatureTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant COSESignatureType.%s", s)
}

// CertificateOriginValue is the CertificateOrigin adapter: Adapter print()/parse() map to the Java enum name().
type CertificateOriginValue enumerations.CertificateOrigin

// CertificateOrigin returns the underlying enumeration constant.
func (v CertificateOriginValue) CertificateOrigin() enumerations.CertificateOrigin {
	return enumerations.CertificateOrigin(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v CertificateOriginValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *CertificateOriginValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.CertificateOriginValues() {
		if string(candidate) == s {
			*v = CertificateOriginValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant CertificateOrigin.%s", s)
}

// CertificateRefOriginValue is the CertificateRefOrigin adapter: Adapter print()/parse() map to the Java enum name().
type CertificateRefOriginValue enumerations.CertificateRefOrigin

// CertificateRefOrigin returns the underlying enumeration constant.
func (v CertificateRefOriginValue) CertificateRefOrigin() enumerations.CertificateRefOrigin {
	return enumerations.CertificateRefOrigin(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v CertificateRefOriginValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *CertificateRefOriginValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.CertificateRefOriginValues() {
		if string(candidate) == s {
			*v = CertificateRefOriginValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant CertificateRefOrigin.%s", s)
}

// CertificateSourceTypeValue is the CertificateSourceType adapter: Adapter print()/parse() map to the Java enum name().
type CertificateSourceTypeValue enumerations.CertificateSourceType

// CertificateSourceType returns the underlying enumeration constant.
func (v CertificateSourceTypeValue) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v CertificateSourceTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *CertificateSourceTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.CertificateSourceTypeValues() {
		if string(candidate) == s {
			*v = CertificateSourceTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant CertificateSourceType.%s", s)
}

// CertificateStatusValue is the CertificateStatus adapter: Adapter print()/parse() map to the Java enum name().
type CertificateStatusValue enumerations.CertificateStatus

// CertificateStatus returns the underlying enumeration constant.
func (v CertificateStatusValue) CertificateStatus() enumerations.CertificateStatus {
	return enumerations.CertificateStatus(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v CertificateStatusValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *CertificateStatusValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.CertificateStatusValues() {
		if string(candidate) == s {
			*v = CertificateStatusValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant CertificateStatus.%s", s)
}

// CertificationPermissionValue is the CertificationPermission adapter: Adapter print()/parse() map to the Java enum name().
type CertificationPermissionValue enumerations.CertificationPermission

// CertificationPermission returns the underlying enumeration constant.
func (v CertificationPermissionValue) CertificationPermission() enumerations.CertificationPermission {
	return enumerations.CertificationPermission(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v CertificationPermissionValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *CertificationPermissionValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.CertificationPermissionValues() {
		if string(candidate) == s {
			*v = CertificationPermissionValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant CertificationPermission.%s", s)
}

// DigestAlgorithmValue is the DigestAlgorithm adapter: DigestAlgorithmParser.print() writes name(); parse() replaces '-' with '_' first.
type DigestAlgorithmValue enumerations.DigestAlgorithm

// DigestAlgorithm returns the underlying enumeration constant.
func (v DigestAlgorithmValue) DigestAlgorithm() enumerations.DigestAlgorithm {
	return enumerations.DigestAlgorithm(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v DigestAlgorithmValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *DigestAlgorithmValue) UnmarshalText(text []byte) error {
	s := strings.ReplaceAll(string(text), "-", "_")
	for _, candidate := range enumerations.DigestAlgorithmValues() {
		if string(candidate) == s {
			*v = DigestAlgorithmValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant DigestAlgorithm.%s", s)
}

// DigestMatcherTypeValue is the DigestMatcherType adapter: Adapter print()/parse() map to the Java enum name().
type DigestMatcherTypeValue enumerations.DigestMatcherType

// DigestMatcherType returns the underlying enumeration constant.
func (v DigestMatcherTypeValue) DigestMatcherType() enumerations.DigestMatcherType {
	return enumerations.DigestMatcherType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v DigestMatcherTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *DigestMatcherTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.DigestMatcherTypeValues() {
		if string(candidate) == s {
			*v = DigestMatcherTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant DigestMatcherType.%s", s)
}

// EAAPresentationTypeValue is the EAAPresentationType adapter: Adapter print()/parse() map to the Java enum name().
type EAAPresentationTypeValue enumerations.EAAPresentationType

// EAAPresentationType returns the underlying enumeration constant.
func (v EAAPresentationTypeValue) EAAPresentationType() enumerations.EAAPresentationType {
	return enumerations.EAAPresentationType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v EAAPresentationTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EAAPresentationTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EAAPresentationTypeValues() {
		if string(candidate) == s {
			*v = EAAPresentationTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EAAPresentationType.%s", s)
}

// EAARevocationOriginValue is the EAARevocationOrigin adapter: Adapter print()/parse() map to the Java enum name().
type EAARevocationOriginValue enumerations.EAARevocationOrigin

// EAARevocationOrigin returns the underlying enumeration constant.
func (v EAARevocationOriginValue) EAARevocationOrigin() enumerations.EAARevocationOrigin {
	return enumerations.EAARevocationOrigin(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v EAARevocationOriginValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EAARevocationOriginValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EAARevocationOriginValues() {
		if string(candidate) == s {
			*v = EAARevocationOriginValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EAARevocationOrigin.%s", s)
}

// EAAStatusValue is the EAAStatus adapter: Adapter print()/parse() map to the Java enum name().
type EAAStatusValue enumerations.EAAStatus

// EAAStatus returns the underlying enumeration constant.
func (v EAAStatusValue) EAAStatus() enumerations.EAAStatus { return enumerations.EAAStatus(v) }

// MarshalText writes the lexical form produced by the Java parser.
func (v EAAStatusValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EAAStatusValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EAAStatusValues() {
		if string(candidate) == s {
			*v = EAAStatusValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EAAStatus.%s", s)
}

// EAATypeValue is the EAAType adapter: Adapter print()/parse() map to the Java enum name().
type EAATypeValue enumerations.EAAType

// EAAType returns the underlying enumeration constant.
func (v EAATypeValue) EAAType() enumerations.EAAType { return enumerations.EAAType(v) }

// MarshalText writes the lexical form produced by the Java parser.
func (v EAATypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EAATypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EAATypeValues() {
		if string(candidate) == s {
			*v = EAATypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EAAType.%s", s)
}

// EncryptionAlgorithmValue is the EncryptionAlgorithm adapter: EncryptionAlgorithmParser prints getName() and parses via forName().
type EncryptionAlgorithmValue enumerations.EncryptionAlgorithm

// EncryptionAlgorithm returns the underlying enumeration constant.
func (v EncryptionAlgorithmValue) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return enumerations.EncryptionAlgorithm(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v EncryptionAlgorithmValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.EncryptionAlgorithm(v).Name()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EncryptionAlgorithmValue) UnmarshalText(text []byte) error {
	enum, err := enumerations.EncryptionAlgorithmForName(string(text))
	if err != nil {
		return err
	}
	*v = EncryptionAlgorithmValue(enum)
	return nil
}

// EndorsementTypeValue is the EndorsementType adapter: The parser prints getValue() and parses via the value.
type EndorsementTypeValue enumerations.EndorsementType

// EndorsementType returns the underlying enumeration constant.
func (v EndorsementTypeValue) EndorsementType() enumerations.EndorsementType {
	return enumerations.EndorsementType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v EndorsementTypeValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.EndorsementType(v).Value()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EndorsementTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EndorsementTypeValues() {
		if enumerations.EndorsementType(candidate).Value() == s {
			*v = EndorsementTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EndorsementType.%s", s)
}

// EvidenceRecordIncorporationTypeValue is the EvidenceRecordIncorporationType adapter: Adapter print()/parse() map to the Java enum name().
type EvidenceRecordIncorporationTypeValue enumerations.EvidenceRecordIncorporationType

// EvidenceRecordIncorporationType returns the underlying enumeration constant.
func (v EvidenceRecordIncorporationTypeValue) EvidenceRecordIncorporationType() enumerations.EvidenceRecordIncorporationType {
	return enumerations.EvidenceRecordIncorporationType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v EvidenceRecordIncorporationTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EvidenceRecordIncorporationTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EvidenceRecordIncorporationTypeValues() {
		if string(candidate) == s {
			*v = EvidenceRecordIncorporationTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EvidenceRecordIncorporationType.%s", s)
}

// EvidenceRecordOriginValue is the EvidenceRecordOrigin adapter: Adapter print()/parse() map to the Java enum name().
type EvidenceRecordOriginValue enumerations.EvidenceRecordOrigin

// EvidenceRecordOrigin returns the underlying enumeration constant.
func (v EvidenceRecordOriginValue) EvidenceRecordOrigin() enumerations.EvidenceRecordOrigin {
	return enumerations.EvidenceRecordOrigin(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v EvidenceRecordOriginValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EvidenceRecordOriginValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EvidenceRecordOriginValues() {
		if string(candidate) == s {
			*v = EvidenceRecordOriginValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EvidenceRecordOrigin.%s", s)
}

// EvidenceRecordTimestampTypeValue is the EvidenceRecordTimestampType adapter: Adapter print()/parse() map to the Java enum name().
type EvidenceRecordTimestampTypeValue enumerations.EvidenceRecordTimestampType

// EvidenceRecordTimestampType returns the underlying enumeration constant.
func (v EvidenceRecordTimestampTypeValue) EvidenceRecordTimestampType() enumerations.EvidenceRecordTimestampType {
	return enumerations.EvidenceRecordTimestampType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v EvidenceRecordTimestampTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EvidenceRecordTimestampTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EvidenceRecordTimestampTypeValues() {
		if string(candidate) == s {
			*v = EvidenceRecordTimestampTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EvidenceRecordTimestampType.%s", s)
}

// EvidenceRecordTypeEnumValue is the EvidenceRecordTypeEnum adapter: The parser prints getLabel() and parses via the label.
type EvidenceRecordTypeEnumValue enumerations.EvidenceRecordTypeEnum

// EvidenceRecordTypeEnum returns the underlying enumeration constant.
func (v EvidenceRecordTypeEnumValue) EvidenceRecordTypeEnum() enumerations.EvidenceRecordTypeEnum {
	return enumerations.EvidenceRecordTypeEnum(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v EvidenceRecordTypeEnumValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.EvidenceRecordTypeEnum(v).Label()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *EvidenceRecordTypeEnumValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.EvidenceRecordTypeEnumValues() {
		if enumerations.EvidenceRecordTypeEnum(candidate).Label() == s {
			*v = EvidenceRecordTypeEnumValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant EvidenceRecordTypeEnum.%s", s)
}

// GeneralNameTypeValue is the GeneralNameType adapter: The parser prints getLabel() and parses via the label.
type GeneralNameTypeValue enumerations.GeneralNameType

// GeneralNameType returns the underlying enumeration constant.
func (v GeneralNameTypeValue) GeneralNameType() enumerations.GeneralNameType {
	return enumerations.GeneralNameType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v GeneralNameTypeValue) MarshalText() ([]byte, error) {
	return []byte(enumerations.GeneralNameType(v).Label()), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *GeneralNameTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.GeneralNameTypeValues() {
		if enumerations.GeneralNameType(candidate).Label() == s {
			*v = GeneralNameTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant GeneralNameType.%s", s)
}

// JWSSerializationTypeValue is the JWSSerializationType adapter: Adapter print()/parse() map to the Java enum name().
type JWSSerializationTypeValue enumerations.JWSSerializationType

// JWSSerializationType returns the underlying enumeration constant.
func (v JWSSerializationTypeValue) JWSSerializationType() enumerations.JWSSerializationType {
	return enumerations.JWSSerializationType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v JWSSerializationTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *JWSSerializationTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.JWSSerializationTypeValues() {
		if string(candidate) == s {
			*v = JWSSerializationTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant JWSSerializationType.%s", s)
}

// KeyUsageBitValue is the KeyUsageBit adapter: The parser prints getValue() and parses via the value.
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

// PdfLockActionValue is the PdfLockAction adapter: Adapter print()/parse() map to the Java enum name().
type PdfLockActionValue enumerations.PdfLockAction

// PdfLockAction returns the underlying enumeration constant.
func (v PdfLockActionValue) PdfLockAction() enumerations.PdfLockAction {
	return enumerations.PdfLockAction(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v PdfLockActionValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *PdfLockActionValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.PdfLockActionValues() {
		if string(candidate) == s {
			*v = PdfLockActionValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant PdfLockAction.%s", s)
}

// PdfObjectModificationTypeValue is the PdfObjectModificationType adapter: Adapter print()/parse() map to the Java enum name().
type PdfObjectModificationTypeValue enumerations.PdfObjectModificationType

// PdfObjectModificationType returns the underlying enumeration constant.
func (v PdfObjectModificationTypeValue) PdfObjectModificationType() enumerations.PdfObjectModificationType {
	return enumerations.PdfObjectModificationType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v PdfObjectModificationTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *PdfObjectModificationTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.PdfObjectModificationTypeValues() {
		if string(candidate) == s {
			*v = PdfObjectModificationTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant PdfObjectModificationType.%s", s)
}

// RevocationOriginValue is the RevocationOrigin adapter: Adapter print()/parse() map to the Java enum name().
type RevocationOriginValue enumerations.RevocationOrigin

// RevocationOrigin returns the underlying enumeration constant.
func (v RevocationOriginValue) RevocationOrigin() enumerations.RevocationOrigin {
	return enumerations.RevocationOrigin(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v RevocationOriginValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *RevocationOriginValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.RevocationOriginValues() {
		if string(candidate) == s {
			*v = RevocationOriginValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant RevocationOrigin.%s", s)
}

// RevocationReasonValue is the RevocationReason adapter: RevocationReasonParser prints/parses the short name.
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

// RevocationRefOriginValue is the RevocationRefOrigin adapter: Adapter print()/parse() map to the Java enum name().
type RevocationRefOriginValue enumerations.RevocationRefOrigin

// RevocationRefOrigin returns the underlying enumeration constant.
func (v RevocationRefOriginValue) RevocationRefOrigin() enumerations.RevocationRefOrigin {
	return enumerations.RevocationRefOrigin(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v RevocationRefOriginValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *RevocationRefOriginValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.RevocationRefOriginValues() {
		if string(candidate) == s {
			*v = RevocationRefOriginValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant RevocationRefOrigin.%s", s)
}

// RevocationTypeValue is the RevocationType adapter: Adapter print()/parse() map to the Java enum name().
type RevocationTypeValue enumerations.RevocationType

// RevocationType returns the underlying enumeration constant.
func (v RevocationTypeValue) RevocationType() enumerations.RevocationType {
	return enumerations.RevocationType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v RevocationTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *RevocationTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.RevocationTypeValues() {
		if string(candidate) == s {
			*v = RevocationTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant RevocationType.%s", s)
}

// SignatureLevelValue is the SignatureLevel adapter: The parser prints toString(), which replaces '_' with '-'.
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

// SignatureScopeTypeValue is the SignatureScopeType adapter: Adapter print()/parse() map to the Java enum name().
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

// TimestampTypeValue is the TimestampType adapter: Adapter print()/parse() map to the Java enum name().
type TimestampTypeValue enumerations.TimestampType

// TimestampType returns the underlying enumeration constant.
func (v TimestampTypeValue) TimestampType() enumerations.TimestampType {
	return enumerations.TimestampType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v TimestampTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *TimestampTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.TimestampTypeValues() {
		if string(candidate) == s {
			*v = TimestampTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant TimestampType.%s", s)
}

// TimestampedObjectTypeValue is the TimestampedObjectType adapter: Adapter print()/parse() map to the Java enum name().
type TimestampedObjectTypeValue enumerations.TimestampedObjectType

// TimestampedObjectType returns the underlying enumeration constant.
func (v TimestampedObjectTypeValue) TimestampedObjectType() enumerations.TimestampedObjectType {
	return enumerations.TimestampedObjectType(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v TimestampedObjectTypeValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *TimestampedObjectTypeValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.TimestampedObjectTypeValues() {
		if string(candidate) == s {
			*v = TimestampedObjectTypeValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant TimestampedObjectType.%s", s)
}
