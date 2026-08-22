// Ported from dss-enumerations/.../TokenExtractionStrategy.java (DSS 6.5.RC1).
//
// Defines a representation of tokens in the Data (as binaries or
// digests).
package enumerations

// TokenExtractionStrategy defines which token types get extracted into the
// Data.
type TokenExtractionStrategy string

const (
	// TokenExtractionStrategyExtractAll extracts certificates,
	// timestamps and revocation data.
	TokenExtractionStrategyExtractAll TokenExtractionStrategy = "EXTRACT_ALL"
	// TokenExtractionStrategyExtractCertificatesOnly extracts
	// certificates.
	TokenExtractionStrategyExtractCertificatesOnly TokenExtractionStrategy = "EXTRACT_CERTIFICATES_ONLY"
	// TokenExtractionStrategyExtractTimestampsOnly extracts timestamps.
	TokenExtractionStrategyExtractTimestampsOnly TokenExtractionStrategy = "EXTRACT_TIMESTAMPS_ONLY"
	// TokenExtractionStrategyExtractRevocationDataOnly extracts
	// revocation data.
	TokenExtractionStrategyExtractRevocationDataOnly TokenExtractionStrategy = "EXTRACT_REVOCATION_DATA_ONLY"
	// TokenExtractionStrategyExtractEvidenceRecordsOnly extracts
	// evidence records.
	TokenExtractionStrategyExtractEvidenceRecordsOnly TokenExtractionStrategy = "EXTRACT_EVIDENCE_RECORDS_ONLY"
	// TokenExtractionStrategyExtractCertificatesAndTimestamps extracts
	// certificates and timestamps.
	TokenExtractionStrategyExtractCertificatesAndTimestamps TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_TIMESTAMPS"
	// TokenExtractionStrategyExtractCertificatesAndEvidenceRecords
	// extracts certificates and timestamps.
	TokenExtractionStrategyExtractCertificatesAndEvidenceRecords TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategyExtractCertificatesAndTimestampsAndEvidenceRecords
	// extracts certificates, timestamps and evidence records.
	TokenExtractionStrategyExtractCertificatesAndTimestampsAndEvidenceRecords TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategyExtractCertificatesAndRevocationData
	// extracts certificates and revocation data.
	TokenExtractionStrategyExtractCertificatesAndRevocationData TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_REVOCATION_DATA"
	// TokenExtractionStrategyExtractCertificatesAndTimestampsAndRevocationData
	// extracts certificates, timestamps and evidence records.
	TokenExtractionStrategyExtractCertificatesAndTimestampsAndRevocationData TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_REVOCATION_DATA"
	// TokenExtractionStrategyExtractCertificatesAndRevocationDataAndEvidenceRecords
	// extracts certificates, revocation data and evidence records.
	TokenExtractionStrategyExtractCertificatesAndRevocationDataAndEvidenceRecords TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategyExtractTimestampsAndRevocationData
	// extracts timestamps and revocation data.
	TokenExtractionStrategyExtractTimestampsAndRevocationData TokenExtractionStrategy = "EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA"
	// TokenExtractionStrategyExtractTimestampsAndEvidenceRecords
	// extracts timestamps and evidence records.
	TokenExtractionStrategyExtractTimestampsAndEvidenceRecords TokenExtractionStrategy = "EXTRACT_TIMESTAMPS_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategyExtractRevocationDataAndEvidenceRecords
	// extracts revocation data and evidence records.
	TokenExtractionStrategyExtractRevocationDataAndEvidenceRecords TokenExtractionStrategy = "EXTRACT_REVOCATION_DATA_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategyExtractTimestampsAndRevocationDataAndEvidenceRecords
	// extracts timestamps, revocation data and evidence records.
	TokenExtractionStrategyExtractTimestampsAndRevocationDataAndEvidenceRecords TokenExtractionStrategy = "EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategyNone extracts nothing.
	TokenExtractionStrategyNone TokenExtractionStrategy = "NONE"
)

type tokenExtractionStrategyFields struct {
	certificate    bool
	timestamp      bool
	revocationData bool
	evidenceRecord bool
}

// tokenExtractionStrategyData holds the (certificate, timestamp,
// revocationData, evidenceRecord) tuple for each constant.
var tokenExtractionStrategyData = map[TokenExtractionStrategy]tokenExtractionStrategyFields{
	TokenExtractionStrategyExtractAll:                                             {true, true, true, true},
	TokenExtractionStrategyExtractCertificatesOnly:                                {true, false, false, false},
	TokenExtractionStrategyExtractTimestampsOnly:                                  {false, true, false, false},
	TokenExtractionStrategyExtractRevocationDataOnly:                              {false, false, true, false},
	TokenExtractionStrategyExtractEvidenceRecordsOnly:                             {false, false, false, true},
	TokenExtractionStrategyExtractCertificatesAndTimestamps:                       {true, true, false, false},
	TokenExtractionStrategyExtractCertificatesAndEvidenceRecords:                  {true, false, false, true},
	TokenExtractionStrategyExtractCertificatesAndTimestampsAndEvidenceRecords:     {true, true, false, true},
	TokenExtractionStrategyExtractCertificatesAndRevocationData:                   {true, false, true, false},
	TokenExtractionStrategyExtractCertificatesAndTimestampsAndRevocationData:      {true, true, true, false},
	TokenExtractionStrategyExtractCertificatesAndRevocationDataAndEvidenceRecords: {true, false, true, true},
	TokenExtractionStrategyExtractTimestampsAndRevocationData:                     {false, true, true, false},
	TokenExtractionStrategyExtractTimestampsAndEvidenceRecords:                    {false, true, false, true},
	TokenExtractionStrategyExtractRevocationDataAndEvidenceRecords:                {false, false, true, true},
	TokenExtractionStrategyExtractTimestampsAndRevocationDataAndEvidenceRecords:   {false, true, true, true},
	TokenExtractionStrategyNone:                                                   {false, false, false, false},
}

// TokenExtractionStrategyValues returns all constants in declaration order.
func TokenExtractionStrategyValues() []TokenExtractionStrategy {
	return []TokenExtractionStrategy{
		TokenExtractionStrategyExtractAll,
		TokenExtractionStrategyExtractCertificatesOnly,
		TokenExtractionStrategyExtractTimestampsOnly,
		TokenExtractionStrategyExtractRevocationDataOnly,
		TokenExtractionStrategyExtractEvidenceRecordsOnly,
		TokenExtractionStrategyExtractCertificatesAndTimestamps,
		TokenExtractionStrategyExtractCertificatesAndEvidenceRecords,
		TokenExtractionStrategyExtractCertificatesAndTimestampsAndEvidenceRecords,
		TokenExtractionStrategyExtractCertificatesAndRevocationData,
		TokenExtractionStrategyExtractCertificatesAndTimestampsAndRevocationData,
		TokenExtractionStrategyExtractCertificatesAndRevocationDataAndEvidenceRecords,
		TokenExtractionStrategyExtractTimestampsAndRevocationData,
		TokenExtractionStrategyExtractTimestampsAndEvidenceRecords,
		TokenExtractionStrategyExtractRevocationDataAndEvidenceRecords,
		TokenExtractionStrategyExtractTimestampsAndRevocationDataAndEvidenceRecords,
		TokenExtractionStrategyNone,
	}
}

// IsCertificate returns true if the certificate extraction is enabled.
func (t TokenExtractionStrategy) IsCertificate() bool {
	return tokenExtractionStrategyData[t].certificate
}

// IsTimestamp returns true if the timestamp extraction is enabled.
func (t TokenExtractionStrategy) IsTimestamp() bool {
	return tokenExtractionStrategyData[t].timestamp
}

// IsRevocationData returns true if the revocation data extraction is enabled.
func (t TokenExtractionStrategy) IsRevocationData() bool {
	return tokenExtractionStrategyData[t].revocationData
}

// IsEvidenceRecord returns true if the evidence record extraction is enabled.
func (t TokenExtractionStrategy) IsEvidenceRecord() bool {
	return tokenExtractionStrategyData[t].evidenceRecord
}

// TokenExtractionStrategyFromParameters returns the enumeration value
// depending on parameters. Returns TokenExtractionStrategyNone if no
// constant matches, mirroring the Java fallback.
func TokenExtractionStrategyFromParameters(certificate, timestamp, revocationData, evidenceRecord bool) TokenExtractionStrategy {
	for _, v := range TokenExtractionStrategyValues() {
		f := tokenExtractionStrategyData[v]
		if certificate == f.certificate && timestamp == f.timestamp &&
			revocationData == f.revocationData && evidenceRecord == f.evidenceRecord {
			return v
		}
	}
	return TokenExtractionStrategyNone
}
