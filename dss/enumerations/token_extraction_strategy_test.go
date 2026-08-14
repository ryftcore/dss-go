package enumerations

import "testing"

func TestTokenExtractionStrategy(t *testing.T) {
	cases := []struct {
		v                                                      TokenExtractionStrategy
		certificate, timestamp, revocationData, evidenceRecord bool
	}{
		{TokenExtractionStrategy_EXTRACT_ALL, true, true, true, true},
		{TokenExtractionStrategy_EXTRACT_CERTIFICATES_ONLY, true, false, false, false},
		{TokenExtractionStrategy_EXTRACT_TIMESTAMPS_ONLY, false, true, false, false},
		{TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_ONLY, false, false, true, false},
		{TokenExtractionStrategy_EXTRACT_EVIDENCE_RECORDS_ONLY, false, false, false, true},
		{TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS, true, true, false, false},
		{TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_EVIDENCE_RECORDS, true, false, false, true},
		{TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_EVIDENCE_RECORDS, true, true, false, true},
		{TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA, true, false, true, false},
		{TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_TIMESTAMPS_AND_REVOCATION_DATA, true, true, true, false},
		{TokenExtractionStrategy_EXTRACT_CERTIFICATES_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS, true, false, true, true},
		{TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA, false, true, true, false},
		{TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_EVIDENCE_RECORDS, false, true, false, true},
		{TokenExtractionStrategy_EXTRACT_REVOCATION_DATA_AND_EVIDENCE_RECORDS, false, false, true, true},
		{TokenExtractionStrategy_EXTRACT_TIMESTAMPS_AND_REVOCATION_DATA_AND_EVIDENCE_RECORDS, false, true, true, true},
		{TokenExtractionStrategy_NONE, false, false, false, false},
	}
	if len(TokenExtractionStrategyValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(TokenExtractionStrategyValues()))
	}
	for _, c := range cases {
		if got := c.v.IsCertificate(); got != c.certificate {
			t.Errorf("%v.IsCertificate() = %v, want %v", c.v, got, c.certificate)
		}
		if got := c.v.IsTimestamp(); got != c.timestamp {
			t.Errorf("%v.IsTimestamp() = %v, want %v", c.v, got, c.timestamp)
		}
		if got := c.v.IsRevocationData(); got != c.revocationData {
			t.Errorf("%v.IsRevocationData() = %v, want %v", c.v, got, c.revocationData)
		}
		if got := c.v.IsEvidenceRecord(); got != c.evidenceRecord {
			t.Errorf("%v.IsEvidenceRecord() = %v, want %v", c.v, got, c.evidenceRecord)
		}
		if got := TokenExtractionStrategyFromParameters(c.certificate, c.timestamp, c.revocationData, c.evidenceRecord); got != c.v {
			t.Errorf("TokenExtractionStrategyFromParameters(%v,%v,%v,%v) = %v, want %v",
				c.certificate, c.timestamp, c.revocationData, c.evidenceRecord, got, c.v)
		}
	}
}
