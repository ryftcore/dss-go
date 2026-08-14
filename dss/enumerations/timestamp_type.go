// Ported from dss-enumerations/.../TimestampType.java (DSS 6.5.RC1).
package enumerations

// TimestampType is the type of timestamp.
type TimestampType string

const (
	// TimestampType_CONTENT_TIMESTAMP: CAdES: id-aa-ets-contentTimestamp,
	// JAdES: adoTst.
	TimestampType_CONTENT_TIMESTAMP TimestampType = "CONTENT_TIMESTAMP"
	// TimestampType_ALL_DATA_OBJECTS_TIMESTAMP: XAdES:
	// AllDataObjectsTimestamp.
	TimestampType_ALL_DATA_OBJECTS_TIMESTAMP TimestampType = "ALL_DATA_OBJECTS_TIMESTAMP"
	// TimestampType_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP: XAdES:
	// IndividualDataObjectsTimeStamp.
	TimestampType_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP TimestampType = "INDIVIDUAL_DATA_OBJECTS_TIMESTAMP"
	// TimestampType_SIGNATURE_TIMESTAMP: CAdES/PAdES:
	// id-aa-signatureTimeStampToken, XAdES: SignatureTimeStamp, JAdES:
	// sigTst.
	TimestampType_SIGNATURE_TIMESTAMP TimestampType = "SIGNATURE_TIMESTAMP"
	// TimestampType_VRI_TIMESTAMP: PAdES: /VRI/TS.
	TimestampType_VRI_TIMESTAMP TimestampType = "VRI_TIMESTAMP"
	// TimestampType_VALIDATION_DATA_REFSONLY_TIMESTAMP: CAdES:
	// id-aa-ets-certCRLTimestamp, XAdES: RefsOnlyTimeStamp, JAdES: rfsTst.
	TimestampType_VALIDATION_DATA_REFSONLY_TIMESTAMP TimestampType = "VALIDATION_DATA_REFSONLY_TIMESTAMP"
	// TimestampType_VALIDATION_DATA_TIMESTAMP: CAdES: id-aa-ets-escTimeStamp,
	// XAdES: SigAndRefsTimeStamp, JAdES: sigRTst.
	TimestampType_VALIDATION_DATA_TIMESTAMP TimestampType = "VALIDATION_DATA_TIMESTAMP"
	// TimestampType_CONTAINER_TIMESTAMP is the ASiC detached timestamp.
	TimestampType_CONTAINER_TIMESTAMP TimestampType = "CONTAINER_TIMESTAMP"
	// TimestampType_DOCUMENT_TIMESTAMP is the PAdES-LTV "document
	// timestamp".
	TimestampType_DOCUMENT_TIMESTAMP TimestampType = "DOCUMENT_TIMESTAMP"
	// TimestampType_ARCHIVE_TIMESTAMP: CAdES: id-aa-ets-archiveTimestamp,
	// XAdES: ArchiveTimeStamp, JAdES: arcTst.
	TimestampType_ARCHIVE_TIMESTAMP TimestampType = "ARCHIVE_TIMESTAMP"
	// TimestampType_EVIDENCE_RECORD_TIMESTAMP is an evidence record
	// time-stamp.
	TimestampType_EVIDENCE_RECORD_TIMESTAMP TimestampType = "EVIDENCE_RECORD_TIMESTAMP"
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
	TimestampType_CONTENT_TIMESTAMP:                  {0, false},
	TimestampType_ALL_DATA_OBJECTS_TIMESTAMP:         {0, false},
	TimestampType_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP:  {0, false},
	TimestampType_SIGNATURE_TIMESTAMP:                {1, true},
	TimestampType_VRI_TIMESTAMP:                      {1, true},
	TimestampType_VALIDATION_DATA_REFSONLY_TIMESTAMP: {2, false},
	TimestampType_VALIDATION_DATA_TIMESTAMP:          {2, true},
	TimestampType_CONTAINER_TIMESTAMP:                {3, true},
	TimestampType_DOCUMENT_TIMESTAMP:                 {3, true},
	TimestampType_ARCHIVE_TIMESTAMP:                  {3, true},
	TimestampType_EVIDENCE_RECORD_TIMESTAMP:          {3, true},
}

// TimestampTypeValues returns all constants in declaration order.
func TimestampTypeValues() []TimestampType {
	return []TimestampType{
		TimestampType_CONTENT_TIMESTAMP,
		TimestampType_ALL_DATA_OBJECTS_TIMESTAMP,
		TimestampType_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP,
		TimestampType_SIGNATURE_TIMESTAMP,
		TimestampType_VRI_TIMESTAMP,
		TimestampType_VALIDATION_DATA_REFSONLY_TIMESTAMP,
		TimestampType_VALIDATION_DATA_TIMESTAMP,
		TimestampType_CONTAINER_TIMESTAMP,
		TimestampType_DOCUMENT_TIMESTAMP,
		TimestampType_ARCHIVE_TIMESTAMP,
		TimestampType_EVIDENCE_RECORD_TIMESTAMP,
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
	return t == TimestampType_CONTAINER_TIMESTAMP
}

// IsDocumentTimestamp checks if the timestamp type is a document timestamp
// (used for PAdES).
func (t TimestampType) IsDocumentTimestamp() bool {
	return t == TimestampType_DOCUMENT_TIMESTAMP
}

// IsArchivalTimestamp checks if the timestamp type is an archive
// timestamp.
func (t TimestampType) IsArchivalTimestamp() bool {
	return t == TimestampType_ARCHIVE_TIMESTAMP
}

// IsEvidenceRecordTimestamp checks if the timestamp type is an evidence
// record timestamp.
func (t TimestampType) IsEvidenceRecordTimestamp() bool {
	return t == TimestampType_EVIDENCE_RECORD_TIMESTAMP
}

// CoversSignature checks if a timestamp of this type covers a signature.
func (t TimestampType) CoversSignature() bool {
	return timestampTypeData[t].coversSignature
}

// Compare compares this TimestampType with the provided timestampType.
// Must be in the order: Content - Signature - ValidationData - Archival.
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
