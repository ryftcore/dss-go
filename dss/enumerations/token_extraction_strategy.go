// Ported from dss-enumerations/.../TokenExtractionStrategy.java (DSS 6.5.RC1).
//
// Defines a representation of tokens in the DiagnosticData (as binaries or
// digests).
package enumerations

// TokenExtractionStrategy defines which token types get extracted into the
// DiagnosticData.
type TokenExtractionStrategy string

const (
	// TokenExtractionStrategy_EXTRACT_ALL extracts certificates,
	// timestamps and revocation data.
	TokenExtractionStrategy_EXTRACT_ALL TokenExtractionStrategy = "EXTRACT_ALL"
	// TokenExtractionStrategy_EXTRACT_CERTIFICATES_ONLY extracts
	// certificates.
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_ONLY TokenExtractionStrategy = "EXTRACT_CERTIFICATES_ONLY"
	// TokenExtractionStrategy_EXTRACT_TIMESTAMPS_ONLY extracts timestamps.
	TokenExtractionStrategy_EXTRACT_TIMESTAMPS_ONLY TokenExtractionStrategy = "EXTRACT_TIMESTAMPS_ONLY"
	// TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_ONLY extracts
	// revocation data.
	TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_ONLY TokenExtractionStrategy = "EXTRACT_REVOCATION_DATA_ONLY"
	// TokenExtractionStrategy_EXTRACT_EVIDENCE_RECORDS_ONLY extracts
	// evidence records.
	TokenExtractionStrategy_EXTRACT_EVIDENCE_RECORDS_ONLY TokenExtractionStrategy = "EXTRACT_EVIDENCE_RECORDS_ONLY"
	// TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS extracts
	// certificates and timestamps.
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_TIMESTAMPS"
	// TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_EVIDENCE_RECORDS
	// extracts certificates and timestamps.
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_EVIDENCE_RECORDS TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_EVIDENCE_RECORDS
	// extracts certificates, timestamps and evidence records.
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_EVIDENCE_RECORDS TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA
	// extracts certificates and revocation data.
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_REVOCATION_DATA"
	// TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_REVOCATION_DATA
	// extracts certificates, timestamps and evidence records.
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_REVOCATION_DATA TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_REVOCATION_DATA"
	// TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS
	// extracts certificates, revocation data and evidence records.
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS TokenExtractionStrategy = "EXTRACT_CERTIFICATES_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA
	// extracts timestamps and revocation data.
	TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA TokenExtractionStrategy = "EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA"
	// TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_EVIDENCE_RECORDS
	// extracts timestamps and evidence records.
	TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_EVIDENCE_RECORDS TokenExtractionStrategy = "EXTRACT_TIMESTAMPS_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_AND_EVIDENCE_RECORDS
	// extracts revocation data and evidence records.
	TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_AND_EVIDENCE_RECORDS TokenExtractionStrategy = "EXTRACT_REVOCATION_DATA_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS
	// extracts timestamps, revocation data and evidence records.
	TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS TokenExtractionStrategy = "EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS"
	// TokenExtractionStrategy_NONE extracts nothing.
	TokenExtractionStrategy_NONE TokenExtractionStrategy = "NONE"
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
	TokenExtractionStrategy_EXTRACT_ALL:                                                   {true, true, true, true},
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_ONLY:                                     {true, false, false, false},
	TokenExtractionStrategy_EXTRACT_TIMESTAMPS_ONLY:                                       {false, true, false, false},
	TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_ONLY:                                  {false, false, true, false},
	TokenExtractionStrategy_EXTRACT_EVIDENCE_RECORDS_ONLY:                                 {false, false, false, true},
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS:                           {true, true, false, false},
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_EVIDENCE_RECORDS:                     {true, false, false, true},
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_EVIDENCE_RECORDS:      {true, true, false, true},
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA:                      {true, false, true, false},
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_REVOCATION_DATA:       {true, true, true, false},
	TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS: {true, false, true, true},
	TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA:                        {false, true, true, false},
	TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_EVIDENCE_RECORDS:                       {false, true, false, true},
	TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_AND_EVIDENCE_RECORDS:                  {false, false, true, true},
	TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS:   {false, true, true, true},
	TokenExtractionStrategy_NONE:                                                          {false, false, false, false},
}

// TokenExtractionStrategyValues returns all constants in declaration order.
func TokenExtractionStrategyValues() []TokenExtractionStrategy {
	return []TokenExtractionStrategy{
		TokenExtractionStrategy_EXTRACT_ALL,
		TokenExtractionStrategy_EXTRACT_CERTIFICATES_ONLY,
		TokenExtractionStrategy_EXTRACT_TIMESTAMPS_ONLY,
		TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_ONLY,
		TokenExtractionStrategy_EXTRACT_EVIDENCE_RECORDS_ONLY,
		TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS,
		TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_EVIDENCE_RECORDS,
		TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_EVIDENCE_RECORDS,
		TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA,
		TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_REVOCATION_DATA,
		TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS,
		TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA,
		TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_EVIDENCE_RECORDS,
		TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_AND_EVIDENCE_RECORDS,
		TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS,
		TokenExtractionStrategy_NONE,
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
// depending on parameters. Returns TokenExtractionStrategy_NONE if no
// constant matches, mirroring the Java fallback.
func TokenExtractionStrategyFromParameters(certificate, timestamp, revocationData, evidenceRecord bool) TokenExtractionStrategy {
	for _, v := range TokenExtractionStrategyValues() {
		f := tokenExtractionStrategyData[v]
		if certificate == f.certificate && timestamp == f.timestamp &&
			revocationData == f.revocationData && evidenceRecord == f.evidenceRecord {
			return v
		}
	}
	return TokenExtractionStrategy_NONE
}
