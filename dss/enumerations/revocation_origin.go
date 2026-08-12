// Ported from dss-enumerations/.../RevocationOrigin.java (DSS 6.5.RC1).
package enumerations

// RevocationOrigin lists possible revocation data origins.
type RevocationOrigin string

const (
	// RevocationOrigin_CMS_SIGNED_DATA: the revocation data was embedded in
	// the CMS SignedData itself (used in CMS based (CAdES or
	// TimestampToken)).
	RevocationOrigin_CMS_SIGNED_DATA RevocationOrigin = "CMS_SIGNED_DATA"
	// RevocationOrigin_REVOCATION_VALUES: the revocation data was embedded
	// in the signature 'revocation-values' attribute (used in CAdES and
	// XAdES).
	RevocationOrigin_REVOCATION_VALUES RevocationOrigin = "REVOCATION_VALUES"
	// RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES: the revocation data was
	// embedded in the signature 'AttributeRevocationValues' attribute
	// (used in XAdES).
	RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES RevocationOrigin = "ATTRIBUTE_REVOCATION_VALUES"
	// RevocationOrigin_TIMESTAMP_VALIDATION_DATA: the revocation data was
	// embedded in the signature 'TimeStampValidationData' attribute.
	RevocationOrigin_TIMESTAMP_VALIDATION_DATA RevocationOrigin = "TIMESTAMP_VALIDATION_DATA"
	// RevocationOrigin_ANY_VALIDATION_DATA: the revocation data was
	// embedded in the signature 'AnyValidationData' attribute.
	RevocationOrigin_ANY_VALIDATION_DATA RevocationOrigin = "ANY_VALIDATION_DATA"
	// RevocationOrigin_DSS_DICTIONARY: the revocation data was embedded to
	// the contents of DSS PDF dictionary (used in PAdES).
	RevocationOrigin_DSS_DICTIONARY RevocationOrigin = "DSS_DICTIONARY"
	// RevocationOrigin_VRI_DICTIONARY: the revocation data was embedded to
	// VRI dictionary (used in PAdES).
	RevocationOrigin_VRI_DICTIONARY RevocationOrigin = "VRI_DICTIONARY"
	// RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL: the revocation data
	// was obtained from the ADBE attribute.
	RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL RevocationOrigin = "ADBE_REVOCATION_INFO_ARCHIVAL"
	// RevocationOrigin_EVIDENCE_RECORD: revocation data extracted from an
	// Evidence Record's structure.
	RevocationOrigin_EVIDENCE_RECORD RevocationOrigin = "EVIDENCE_RECORD"
	// RevocationOrigin_INPUT_DOCUMENT: the revocation data was embedded in
	// signature or timestamp.
	RevocationOrigin_INPUT_DOCUMENT RevocationOrigin = "INPUT_DOCUMENT"
	// RevocationOrigin_EXTERNAL: the revocation data was provided by the
	// user or online OCSP/CRL.
	RevocationOrigin_EXTERNAL RevocationOrigin = "EXTERNAL"
	// RevocationOrigin_CACHED: the revocation data was obtained from a
	// local DB or cache.
	RevocationOrigin_CACHED RevocationOrigin = "CACHED"
)

// RevocationOriginValues returns all constants in declaration order.
func RevocationOriginValues() []RevocationOrigin {
	return []RevocationOrigin{
		RevocationOrigin_CMS_SIGNED_DATA,
		RevocationOrigin_REVOCATION_VALUES,
		RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES,
		RevocationOrigin_TIMESTAMP_VALIDATION_DATA,
		RevocationOrigin_ANY_VALIDATION_DATA,
		RevocationOrigin_DSS_DICTIONARY,
		RevocationOrigin_VRI_DICTIONARY,
		RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL,
		RevocationOrigin_EVIDENCE_RECORD,
		RevocationOrigin_INPUT_DOCUMENT,
		RevocationOrigin_EXTERNAL,
		RevocationOrigin_CACHED,
	}
}

// IsInternalOrigin checks if the revocation has been obtained from the
// input document.
func (r RevocationOrigin) IsInternalOrigin() bool {
	return RevocationOrigin_EXTERNAL != r && RevocationOrigin_CACHED != r
}
