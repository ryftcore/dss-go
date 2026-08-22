package enumerations

import "testing"

func TestTokenExtractionStrategy(t *testing.T) {
	cases := []struct {
		v                                                      TokenExtractionStrategy
		certificate, timestamp, revocationData, evidenceRecord bool
	}{
		{TokenExtractionStrategyExtractAll, true, true, true, true},
		{TokenExtractionStrategyExtractCertificatesOnly, true, false, false, false},
		{TokenExtractionStrategyExtractTimestampsOnly, false, true, false, false},
		{TokenExtractionStrategyExtractRevocationDataOnly, false, false, true, false},
		{TokenExtractionStrategyExtractEvidenceRecordsOnly, false, false, false, true},
		{TokenExtractionStrategyExtractCertificatesAndTimestamps, true, true, false, false},
		{TokenExtractionStrategyExtractCertificatesAndEvidenceRecords, true, false, false, true},
		{TokenExtractionStrategyExtractCertificatesAndTimestampsAndEvidenceRecords, true, true, false, true},
		{TokenExtractionStrategyExtractCertificatesAndRevocationData, true, false, true, false},
		{TokenExtractionStrategyExtractCertificatesAndTimestampsAndRevocationData, true, true, true, false},
		{TokenExtractionStrategyExtractCertificatesAndRevocationDataAndEvidenceRecords, true, false, true, true},
		{TokenExtractionStrategyExtractTimestampsAndRevocationData, false, true, true, false},
		{TokenExtractionStrategyExtractTimestampsAndEvidenceRecords, false, true, false, true},
		{TokenExtractionStrategyExtractRevocationDataAndEvidenceRecords, false, false, true, true},
		{TokenExtractionStrategyExtractTimestampsAndRevocationDataAndEvidenceRecords, false, true, true, true},
		{TokenExtractionStrategyNone, false, false, false, false},
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
