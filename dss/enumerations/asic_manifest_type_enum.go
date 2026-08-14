// Ported from dss-enumerations/.../ASiCManifestTypeEnum.java (DSS 6.5.RC1).
package enumerations

// ASiCManifestTypeEnum defines a type of data object associated with the
// ASiCManifest file.
type ASiCManifestTypeEnum string

const (
	// ASiCManifestTypeEnum_SIGNATURE is the ASiCManifest associated with a
	// signature document.
	ASiCManifestTypeEnum_SIGNATURE ASiCManifestTypeEnum = "SIGNATURE"
	// ASiCManifestTypeEnum_TIMESTAMP is the ASiCManifest associated with a
	// time-stamp document.
	ASiCManifestTypeEnum_TIMESTAMP ASiCManifestTypeEnum = "TIMESTAMP"
	// ASiCManifestTypeEnum_EVIDENCE_RECORD is the ASiCEvidenceRecordManifest
	// associated with an evidence record document.
	ASiCManifestTypeEnum_EVIDENCE_RECORD ASiCManifestTypeEnum = "EVIDENCE_RECORD"
	// ASiCManifestTypeEnum_ARCHIVE_MANIFEST is the ASiCArchiveManifest
	// associated with an archival time-stamp document.
	ASiCManifestTypeEnum_ARCHIVE_MANIFEST ASiCManifestTypeEnum = "ARCHIVE_MANIFEST"
)

// ASiCManifestTypeEnumValues returns all constants in declaration order.
func ASiCManifestTypeEnumValues() []ASiCManifestTypeEnum {
	return []ASiCManifestTypeEnum{
		ASiCManifestTypeEnum_SIGNATURE,
		ASiCManifestTypeEnum_TIMESTAMP,
		ASiCManifestTypeEnum_EVIDENCE_RECORD,
		ASiCManifestTypeEnum_ARCHIVE_MANIFEST,
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
