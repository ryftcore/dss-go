// Ported from dss-enumerations/.../ASiCManifestTypeEnum.java (DSS 6.5.RC1).
package enumerations

// ASiCManifestTypeEnum defines a type of data object associated with the
// ASiCManifest file.
type ASiCManifestTypeEnum string

const (
	// ASiCManifestTypeEnumSignature is the ASiCManifest associated with a
	// signature document.
	ASiCManifestTypeEnumSignature ASiCManifestTypeEnum = "SIGNATURE"
	// ASiCManifestTypeEnumTimestamp is the ASiCManifest associated with a
	// time-stamp document.
	ASiCManifestTypeEnumTimestamp ASiCManifestTypeEnum = "TIMESTAMP"
	// ASiCManifestTypeEnumEvidenceRecord is the ASiCEvidenceRecordManifest
	// associated with an evidence record document.
	ASiCManifestTypeEnumEvidenceRecord ASiCManifestTypeEnum = "EVIDENCE_RECORD"
	// ASiCManifestTypeEnumArchiveManifest is the ASiCArchiveManifest
	// associated with an archival time-stamp document.
	ASiCManifestTypeEnumArchiveManifest ASiCManifestTypeEnum = "ARCHIVE_MANIFEST"
)

// ASiCManifestTypeEnumValues returns all constants in declaration order.
func ASiCManifestTypeEnumValues() []ASiCManifestTypeEnum {
	return []ASiCManifestTypeEnum{
		ASiCManifestTypeEnumSignature,
		ASiCManifestTypeEnumTimestamp,
		ASiCManifestTypeEnumEvidenceRecord,
		ASiCManifestTypeEnumArchiveManifest,
	}
}

// ASiCManifestTypeEnumValueOf returns the constant matching the given Java
// enum name.
func ASiCManifestTypeEnumValueOf(name string) (ASiCManifestTypeEnum, error) {
	for _, v := range ASiCManifestTypeEnumValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &aSiCManifestTypeEnumInvalidValueError{name}
}

type aSiCManifestTypeEnumInvalidValueError struct {
	name string
}

func (e *aSiCManifestTypeEnumInvalidValueError) Error() string {
	return "no enum constant ASiCManifestTypeEnum." + e.name
}
