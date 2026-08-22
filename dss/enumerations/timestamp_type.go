// Ported from dss-enumerations/.../TimestampType.java (DSS 6.5.RC1).
package enumerations

// TimestampType is the type of timestamp.
type TimestampType string

const (
	// TimestampTypeContentTimestamp: CAdES: id-aa-ets-contentTimestamp,
	// JAdES: adoTst.
	TimestampTypeContentTimestamp TimestampType = "CONTENT_TIMESTAMP"
	// TimestampTypeAllDataObjectsTimestamp: XAdES:
	// AllDataObjectsTimestamp.
	TimestampTypeAllDataObjectsTimestamp TimestampType = "ALL_DATA_OBJECTS_TIMESTAMP"
	// TimestampTypeIndividualDataObjectsTimestamp: XAdES:
	// IndividualDataObjectsTimeStamp.
	TimestampTypeIndividualDataObjectsTimestamp TimestampType = "INDIVIDUAL_DATA_OBJECTS_TIMESTAMP"
	// TimestampTypeSignatureTimestamp: CAdES/PAdES:
	// id-aa-signatureTimeStampToken, XAdES: SignatureTimeStamp, JAdES:
	// sigTst.
	TimestampTypeSignatureTimestamp TimestampType = "SIGNATURE_TIMESTAMP"
	// TimestampTypeVRITimestamp: PAdES: /VRI/TS.
	TimestampTypeVRITimestamp TimestampType = "VRI_TIMESTAMP"
	// TimestampTypeValidationDataRefsOnlyTimestamp: CAdES:
	// id-aa-ets-certCRLTimestamp, XAdES: RefsOnlyTimeStamp, JAdES: rfsTst.
	TimestampTypeValidationDataRefsOnlyTimestamp TimestampType = "VALIDATION_DATA_REFSONLY_TIMESTAMP"
	// TimestampTypeValidationDataTimestamp: CAdES: id-aa-ets-escTimeStamp,
	// XAdES: SigAndRefsTimeStamp, JAdES: sigRTst.
	TimestampTypeValidationDataTimestamp TimestampType = "VALIDATION_DATA_TIMESTAMP"
	// TimestampTypeContainerTimestamp is the ASiC detached timestamp.
	TimestampTypeContainerTimestamp TimestampType = "CONTAINER_TIMESTAMP"
	// TimestampTypeDocumentTimestamp is the PAdES-LTV "document
	// timestamp".
	TimestampTypeDocumentTimestamp TimestampType = "DOCUMENT_TIMESTAMP"
	// TimestampTypeArchiveTimestamp: CAdES: id-aa-ets-archiveTimestamp,
	// XAdES: ArchiveTimeStamp, JAdES: arcTst.
	TimestampTypeArchiveTimestamp TimestampType = "ARCHIVE_TIMESTAMP"
	// TimestampTypeEvidenceRecordTimestamp is an evidence record
	// time-stamp.
	TimestampTypeEvidenceRecordTimestamp TimestampType = "EVIDENCE_RECORD_TIMESTAMP"
)

type timestampTypeFields struct {
	order           int
	coversSignature bool
}

// timestampTypeData holds the (order, coversSignature) tuple for each
// constant. order specifies a presence order of the timestamp in a
// signature: 0 - content timestamps, 1 - signature timestamp, 2 -
// validation data timestamps, 3 - archive timestamps.
var timestampTypeData = map[TimestampType]timestampTypeFields{
	TimestampTypeContentTimestamp:                {0, false},
	TimestampTypeAllDataObjectsTimestamp:         {0, false},
	TimestampTypeIndividualDataObjectsTimestamp:  {0, false},
	TimestampTypeSignatureTimestamp:              {1, true},
	TimestampTypeVRITimestamp:                    {1, true},
	TimestampTypeValidationDataRefsOnlyTimestamp: {2, false},
	TimestampTypeValidationDataTimestamp:         {2, true},
	TimestampTypeContainerTimestamp:              {3, true},
	TimestampTypeDocumentTimestamp:               {3, true},
	TimestampTypeArchiveTimestamp:                {3, true},
	TimestampTypeEvidenceRecordTimestamp:         {3, true},
}

// TimestampTypeValues returns all constants in declaration order.
func TimestampTypeValues() []TimestampType {
	return []TimestampType{
		TimestampTypeContentTimestamp,
		TimestampTypeAllDataObjectsTimestamp,
		TimestampTypeIndividualDataObjectsTimestamp,
		TimestampTypeSignatureTimestamp,
		TimestampTypeVRITimestamp,
		TimestampTypeValidationDataRefsOnlyTimestamp,
		TimestampTypeValidationDataTimestamp,
		TimestampTypeContainerTimestamp,
		TimestampTypeDocumentTimestamp,
		TimestampTypeArchiveTimestamp,
		TimestampTypeEvidenceRecordTimestamp,
	}
}

// IsContentTimestamp checks if the timestamp type is a content timestamp.
func (t TimestampType) IsContentTimestamp() bool {
	return timestampTypeData[t].order == 0
}

// IsSignatureTimestamp checks if the timestamp type is a signature
// timestamp.
func (t TimestampType) IsSignatureTimestamp() bool {
	return timestampTypeData[t].order == 1
}

// IsValidationDataTimestamp checks if the timestamp type is a validation
// data timestamp.
func (t TimestampType) IsValidationDataTimestamp() bool {
	return timestampTypeData[t].order == 2
}

// IsContainerTimestamp checks if the timestamp type is a container
// timestamp (used for ASiC).
func (t TimestampType) IsContainerTimestamp() bool {
	return t == TimestampTypeContainerTimestamp
}

// IsDocumentTimestamp checks if the timestamp type is a document timestamp
// (used for PAdES).
func (t TimestampType) IsDocumentTimestamp() bool {
	return t == TimestampTypeDocumentTimestamp
}

// IsArchivalTimestamp checks if the timestamp type is an archive
// timestamp.
func (t TimestampType) IsArchivalTimestamp() bool {
	return t == TimestampTypeArchiveTimestamp
}

// IsEvidenceRecordTimestamp checks if the timestamp type is an evidence
// record timestamp.
func (t TimestampType) IsEvidenceRecordTimestamp() bool {
	return t == TimestampTypeEvidenceRecordTimestamp
}

// CoversSignature checks if a timestamp of this type covers a signature.
func (t TimestampType) CoversSignature() bool {
	return timestampTypeData[t].coversSignature
}

// Compare compares this TimestampType with the provided timestampType.
// Must be in the order: Content - Signature - Data - Archival.
func (t TimestampType) Compare(timestampType TimestampType) int {
	a := timestampTypeData[t].order
	b := timestampTypeData[timestampType].order
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// TimestampTypeValueOf returns the constant matching the given Java enum
// name.
func TimestampTypeValueOf(name string) (TimestampType, error) {
	for _, v := range TimestampTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &timestampTypeInvalidValueError{name}
}

type timestampTypeInvalidValueError struct {
	name string
}

func (e *timestampTypeInvalidValueError) Error() string {
	return "no enum constant TimestampType." + e.name
}
