// Ported from dss-enumerations/.../RevocationOrigin.java (DSS 6.5.RC1).
package enumerations

// RevocationOrigin lists possible revocation data origins.
type RevocationOrigin string

const (
	// RevocationOriginCMSSignedData: the revocation data was embedded in
	// the CMS SignedData itself (used in CMS based (CAdES or
	// TimestampToken)).
	RevocationOriginCMSSignedData RevocationOrigin = "CMS_SIGNED_DATA"
	// RevocationOriginRevocationValues: the revocation data was embedded
	// in the signature 'revocation-values' attribute (used in CAdES and
	// XAdES).
	RevocationOriginRevocationValues RevocationOrigin = "REVOCATION_VALUES"
	// RevocationOriginAttributeRevocationValues: the revocation data was
	// embedded in the signature 'AttributeRevocationValues' attribute
	// (used in XAdES).
	RevocationOriginAttributeRevocationValues RevocationOrigin = "ATTRIBUTE_REVOCATION_VALUES"
	// RevocationOriginTimestampValidationData: the revocation data was
	// embedded in the signature 'TimeStampValidationData' attribute.
	RevocationOriginTimestampValidationData RevocationOrigin = "TIMESTAMP_VALIDATION_DATA"
	// RevocationOriginAnyValidationData: the revocation data was
	// embedded in the signature 'AnyValidationData' attribute.
	RevocationOriginAnyValidationData RevocationOrigin = "ANY_VALIDATION_DATA"
	// RevocationOriginDSSDictionary: the revocation data was embedded to
	// the contents of DSS PDF dictionary (used in PAdES).
	RevocationOriginDSSDictionary RevocationOrigin = "DSS_DICTIONARY"
	// RevocationOriginVRIDictionary: the revocation data was embedded to
	// VRI dictionary (used in PAdES).
	RevocationOriginVRIDictionary RevocationOrigin = "VRI_DICTIONARY"
	// RevocationOriginAdbeRevocationInfoArchival: the revocation data
	// was obtained from the ADBE attribute.
	RevocationOriginAdbeRevocationInfoArchival RevocationOrigin = "ADBE_REVOCATION_INFO_ARCHIVAL"
	// RevocationOriginEvidenceRecord: revocation data extracted from an
	// Evidence Record's structure.
	RevocationOriginEvidenceRecord RevocationOrigin = "EVIDENCE_RECORD"
	// RevocationOriginInputDocument: the revocation data was embedded in
	// signature or timestamp.
	RevocationOriginInputDocument RevocationOrigin = "INPUT_DOCUMENT"
	// RevocationOriginExternal: the revocation data was provided by the
	// user or online OCSP/CRL.
	RevocationOriginExternal RevocationOrigin = "EXTERNAL"
	// RevocationOriginCached: the revocation data was obtained from a
	// local DB or cache.
	RevocationOriginCached RevocationOrigin = "CACHED"
)

// RevocationOriginValues returns all constants in declaration order.
func RevocationOriginValues() []RevocationOrigin {
	return []RevocationOrigin{
		RevocationOriginCMSSignedData,
		RevocationOriginRevocationValues,
		RevocationOriginAttributeRevocationValues,
		RevocationOriginTimestampValidationData,
		RevocationOriginAnyValidationData,
		RevocationOriginDSSDictionary,
		RevocationOriginVRIDictionary,
		RevocationOriginAdbeRevocationInfoArchival,
		RevocationOriginEvidenceRecord,
		RevocationOriginInputDocument,
		RevocationOriginExternal,
		RevocationOriginCached,
	}
}

// IsInternalOrigin checks if the revocation has been obtained from the
// input document.
func (r RevocationOrigin) IsInternalOrigin() bool {
	return RevocationOriginExternal != r && RevocationOriginCached != r
}
