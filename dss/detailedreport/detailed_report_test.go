// Behavior tests for DetailedReport/DetailedReportMessageCollector against
// Java-dumped answers: every assertion below reads a value straight out of
// one of jaxb's testdata/oracle fixtures (real JAXB reference implementation
// output), so the expected values are what upstream actually produced, not
// values re-derived from this port's own code.
package detailedreport

import (
	"os"
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/corpustest"
)

func loadOracle(t *testing.T, name string) *DetailedReport {
	t.Helper()
	data, err := os.ReadFile(corpustest.RootPath(t, "detailedreport/jaxb/testdata/oracle/"+name))
	if err != nil {
		t.Fatal(err)
	}
	dr, err := jaxb.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal %s: %v", name, err)
	}
	return NewDetailedReport(dr)
}

// TestDetailedReport_dr1 checks navigation against dr1.xml, a real signature
// validation report (see jaxb/testdata/oracle/dr1.xml).
func TestDetailedReport_dr1(t *testing.T) {
	r := loadOracle(t, "dr1.xml")

	const sigId = "S-3933842D6213C62CF3BE126A56FCF3A9EFFC24D93FEA04B66D6DA5BDEC550703"
	const tstId = "T-E70F7C745B498A2208A0096422489744EEAEC0F30FE781A6DC6F3E96467EDEE2"

	if got := r.FirstSignatureId(); got != sigId {
		t.Errorf("FirstSignatureId() = %q, want %q", got, sigId)
	}
	if got := r.SignatureIds(); len(got) != 1 || got[0] != sigId {
		t.Errorf("SignatureIds() = %v, want [%s]", got, sigId)
	}
	if got := r.XmlSignatureById(sigId); got == nil {
		t.Fatal("XmlSignatureById(sigId) = nil")
	}
	if got := r.XmlSignatureById("does-not-exist"); got != nil {
		t.Errorf("XmlSignatureById(unknown) = %v, want nil", got)
	}

	if got := r.FirstTimestampId(); got != tstId {
		t.Errorf("FirstTimestampId() = %q, want %q", got, tstId)
	}
	if got := r.XmlTimestampById(tstId); got == nil {
		t.Fatal("XmlTimestampById(tstId) = nil")
	}

	// The Signature's own top-level Conclusion/Indication is TOTAL_PASSED
	// (the last <Indication> before </Signature> in dr1.xml).
	if got := r.FinalIndication(sigId); got != enumerations.Indication_TOTAL_PASSED {
		t.Errorf("FinalIndication(sigId) = %v, want TOTAL_PASSED", got)
	}
	if got := r.BasicValidationIndication(sigId); got != enumerations.Indication_PASSED {
		t.Errorf("BasicValidationIndication(sigId) = %v, want PASSED", got)
	}

	// <ValidationSignatureQualification SignatureQualification="QESig" ...>
	if got := r.SignatureQualification(sigId); got != enumerations.SignatureQualification_QESIG {
		t.Errorf("SignatureQualification(sigId) = %v, want QESig", got)
	}

	// <ValidationTimestampQualification TimestampQualification="QTSA" ...>
	if got := r.TimestampQualification(tstId); got != enumerations.TimestampQualification_QTSA {
		t.Errorf("TimestampQualification(tstId) = %v, want QTSA", got)
	}

	if got := r.BasicBuildingBlocksNumber(); got != 5 {
		t.Errorf("BasicBuildingBlocksNumber() = %d, want 5", got)
	}
	bbb := r.BasicBuildingBlockById(sigId)
	if bbb == nil {
		t.Fatal("BasicBuildingBlockById(sigId) = nil")
	}
	if got := bbb.Type.Context(); got != enumerations.Context_SIGNATURE {
		t.Errorf("bbb.Type = %v, want SIGNATURE", got)
	}
	if got := r.BasicBuildingBlocksIndication(sigId); got != enumerations.Indication_PASSED {
		t.Errorf("BasicBuildingBlocksIndication(sigId) = %v, want PASSED", got)
	}

	if got := r.IsCertificateValidation(); got {
		t.Error("IsCertificateValidation() = true for a signature-validation report")
	}

	// A token id that doesn't exist anywhere yields the zero value, not a panic.
	if got := r.FinalIndication("no-such-token"); got != "" {
		t.Errorf("FinalIndication(unknown) = %v, want \"\"", got)
	}
	if got := r.BasicBuildingBlocksSubIndication("no-such-token"); got != "" {
		t.Errorf("BasicBuildingBlocksSubIndication(unknown) = %v, want \"\"", got)
	}
}

// TestDetailedReport_dr1_HighestConclusion exercises HighestConclusion and the
// message collector's dedup, both consumed by 8f's executor per the porting
// brief.
func TestDetailedReport_dr1_HighestConclusion(t *testing.T) {
	r := loadOracle(t, "dr1.xml")
	const sigId = "S-3933842D6213C62CF3BE126A56FCF3A9EFFC24D93FEA04B66D6DA5BDEC550703"

	highest := r.HighestConclusion(sigId)
	if highest == nil || highest.Conclusion == nil {
		t.Fatal("HighestConclusion(sigId) has no Conclusion")
	}
	// dr1.xml's signature has no ValidationProcessLongTermData/ArchivalData
	// block (short-term validation only), so the highest reached level is the
	// basic validation process, whose Conclusion.Indication is PASSED.
	if got := highest.Conclusion.Indication.Indication(); got != enumerations.Indication_PASSED {
		t.Errorf("HighestConclusion(sigId).Conclusion.Indication = %v, want PASSED", got)
	}

	// Every constraint in dr1.xml's signature validation is OK/PASSED, so no
	// AdES validation errors should be collected for it.
	if errs := r.AdESValidationErrors(sigId); len(errs) != 0 {
		t.Errorf("AdESValidationErrors(sigId) = %v, want none", errs)
	}
}

// TestDetailedReport_Fill checks navigation against dr-fill.xml, a
// hand-built-then-JAXB-marshalled fixture exercising fields dr1.xml never
// reaches: SubXCV.CrossCertificate/EquivalentCertificate, VTS.TrustAnchor,
// CertificateApprovalStatus, ValidationTimestampQualificationAtTime and the
// CounterSignature attribute (see jaxb/testdata/oracle/README.md and
// jaxb_schema_test.go's TestOracleCorpusExercisesModel for how it was built).
func TestDetailedReport_Fill(t *testing.T) {
	r := loadOracle(t, "dr-fill.xml")

	sig := r.XmlSignatureById("S-FILL")
	if sig == nil {
		t.Fatal("XmlSignatureById(S-FILL) = nil")
	}
	if sig.CounterSignature == nil || !*sig.CounterSignature {
		t.Errorf("S-FILL CounterSignature = %v, want true", sig.CounterSignature)
	}

	signing := r.SigningCertificate("S-FILL")
	if signing == nil {
		t.Fatal("SigningCertificate(S-FILL) = nil")
	}
	if signing.Id != "C-FILL" {
		t.Errorf("SigningCertificate(S-FILL).Id = %q, want C-FILL", signing.Id)
	}
	if signing.CrossCertificate == nil || len(*signing.CrossCertificate) != 2 ||
		(*signing.CrossCertificate)[0] != "C-CROSS-1" || (*signing.CrossCertificate)[1] != "C-CROSS-2" {
		t.Errorf("SigningCertificate(S-FILL).CrossCertificate = %v, want [C-CROSS-1 C-CROSS-2]", signing.CrossCertificate)
	}
	if signing.EquivalentCertificate == nil || len(*signing.EquivalentCertificate) != 1 || (*signing.EquivalentCertificate)[0] != "C-EQUIV-1" {
		t.Errorf("SigningCertificate(S-FILL).EquivalentCertificate = %v, want [C-EQUIV-1]", signing.EquivalentCertificate)
	}

	bbb := r.BasicBuildingBlockById("S-FILL")
	if bbb == nil || bbb.VTS == nil || bbb.VTS.TrustAnchor == nil || *bbb.VTS.TrustAnchor != "C-ANCHOR" {
		t.Errorf("BasicBuildingBlockById(S-FILL).VTS.TrustAnchor = %v, want C-ANCHOR", bbb)
	}

	if got := r.TimestampQualification("T-FILL"); got != enumerations.TimestampQualification_QTSA {
		t.Errorf("TimestampQualification(T-FILL) = %v, want QTSA", got)
	}
	if got := r.TimestampQualificationAtTstGenerationTime("T-FILL"); got != enumerations.TimestampQualification_QTSA {
		t.Errorf("TimestampQualificationAtTstGenerationTime(T-FILL) = %v, want QTSA", got)
	}
	// No AtTime block was set for TIMESTAMP_POE_TIME, so this should be "".
	if got := r.TimestampQualificationAtBestPoeTime("T-FILL"); got != "" {
		t.Errorf("TimestampQualificationAtBestPoeTime(T-FILL) = %v, want \"\"", got)
	}

	statuses := r.CertificateApprovalStatussAtIssuanceTime("C-FILL-2")
	if len(statuses) != 1 {
		t.Fatalf("CertificateApprovalStatussAtIssuanceTime(C-FILL-2) = %v, want 1 entry", statuses)
	}
	if statuses[0].ListType() == nil || statuses[0].ServiceTypeIdentifier() == nil {
		t.Errorf("CertificateApprovalStatussAtIssuanceTime(C-FILL-2)[0] missing ListType/ServiceTypeIdentifier: %+v", statuses[0])
	}
	// The validation-time slot was never populated for C-FILL-2, so it must
	// come back empty rather than panicking or reusing the issuance-time slot.
	if statuses := r.CertificateApprovalStatussAtValidationTime("C-FILL-2"); len(statuses) != 0 {
		t.Errorf("CertificateApprovalStatussAtValidationTime(C-FILL-2) = %v, want none", statuses)
	}
}

// TestDetailedReportFacade_RoundTrip is the facade-level marshal-parity check
// (the brief requires this "per module"; jaxb.TestMarshalParity is the
// exhaustive corpus version at the model layer this delegates to).
func TestDetailedReportFacade_RoundTrip(t *testing.T) {
	data, err := os.ReadFile(corpustest.RootPath(t, "detailedreport/jaxb/testdata/oracle/dr1.xml"))
	if err != nil {
		t.Fatal(err)
	}
	f := NewDetailedReportFacade()
	dr, err := f.Unmarshal(string(data))
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	got, err := f.Marshal(dr)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got != string(data) {
		t.Errorf("facade round-trip differs from the oracle dump")
	}
}

// TestDetailedReportFacade_NilGuards ports the Java facade's null-argument
// IllegalArgumentException-style guards.
func TestDetailedReportFacade_NilGuards(t *testing.T) {
	f := NewDetailedReportFacade()
	if _, err := f.Marshal(nil); err == nil {
		t.Error("Marshal(nil) succeeded, want an error")
	}
	if _, err := f.Unmarshal(""); err == nil {
		t.Error("Unmarshal(\"\") succeeded, want an error")
	}
}
