package simplereport

import (
	"os"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/simplereport/jaxb"
)

// TestSimpleReportWrapper checks the wrapper's method surface against
// values read directly off a Java-produced oracle dump
// (jaxb/testdata/oracle/counterSig.p7m.xml), a two-signature (one a
// counter-signature) CAdES document with per-signature timestamps.
func TestSimpleReportWrapper(t *testing.T) {
	data, err := os.ReadFile(corpustest.RootPath(t, "simplereport/jaxb/testdata/oracle/counterSig.p7m.xml"))
	if err != nil {
		t.Fatal(err)
	}
	xmlReport, err := jaxb.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	r := NewSimpleReport(xmlReport)

	const parentSig = "S-939C8E0C24EC70765AF7D6D0742A28A03E37980A27D383E4888E876AB0504361"
	const counterSig = "S-5FFD54A32EBD2F81E9FB488BBFBFC1A22C2974A65878B210D3B4215DAC2D3A3F"
	const parentTst = "T-989D00AC9FE4D2BA6F524388871D688BCD6529EC5EFEA8822843568DF0404F8D"

	if got := r.GetDocumentFilename(); got != "counterSig.p7m" {
		t.Errorf("GetDocumentFilename() = %q", got)
	}
	if got, want := r.GetSignaturesCount(), 2; got != want {
		t.Errorf("GetSignaturesCount() = %d, want %d", got, want)
	}
	if got, want := r.GetValidSignaturesCount(), 0; got != want {
		t.Errorf("GetValidSignaturesCount() = %d, want %d", got, want)
	}
	if got := r.GetSignatureIdList(); len(got) != 2 || got[0] != parentSig || got[1] != counterSig {
		t.Errorf("GetSignatureIdList() = %v", got)
	}
	if got := r.GetFirstSignatureId(); got != parentSig {
		t.Errorf("GetFirstSignatureId() = %q, want %q", got, parentSig)
	}
	if got, want := r.GetIndication(parentSig), enumerations.IndicationIndeterminate; got != want {
		t.Errorf("GetIndication(parent) = %q, want %q", got, want)
	}
	if got, want := r.GetSubIndication(parentSig), enumerations.SubIndicationNoCertificateChainFound; got != want {
		t.Errorf("GetSubIndication(parent) = %q, want %q", got, want)
	}
	if r.IsValid(parentSig) {
		t.Error("IsValid(parent) = true, want false")
	}
	if got := r.GetSignedBy(parentSig); got != "EST-COUNTER-SIGNATURE1-OK-EE" {
		t.Errorf("GetSignedBy(parent) = %q", got)
	}
	if got, want := r.GetSignatureFormat(parentSig), enumerations.SignatureLevelCAdEST; got != want {
		t.Errorf("GetSignatureFormat(parent) = %q, want %q", got, want)
	}

	// Counter-signature specific fields.
	sig := r.getSignatureByID(counterSig)
	if sig == nil {
		t.Fatal("getSignatureByID(counterSig) = nil")
	}
	if sig.CounterSignature == nil || !*sig.CounterSignature {
		t.Error("counter-signature's CounterSignature attribute is not true")
	}
	if sig.ParentId == nil || *sig.ParentId != parentSig {
		t.Errorf("counter-signature's ParentId = %v, want %q", sig.ParentId, parentSig)
	}

	// Embedded timestamp lookup by id must descend into Signature.Timestamps.
	tst := r.getTimestampByID(parentTst)
	if tst == nil {
		t.Fatalf("getTimestampByID(%q) = nil", parentTst)
	}
	if got := r.GetProducedBy(parentTst); got != "EST-COUNTER-SIGNATURE1-OK-EE" {
		t.Errorf("GetProducedBy(timestamp) = %q", got)
	}
	tsts := r.GetSignatureTimestamps(parentSig)
	if len(tsts) != 1 || tsts[0].Id != parentTst {
		t.Errorf("GetSignatureTimestamps(parent) = %v", tsts)
	}

	// Errors/warnings on the parent signature's qualification details.
	errs := r.GetQualificationErrors(parentSig)
	if len(errs) != 1 || errs[0].Key != "QUAL_CERT_TRUSTED_LIST_REACHED_ANS" {
		t.Errorf("GetQualificationErrors(parent) = %v", errs)
	}
	warns := r.GetQualificationWarnings(parentSig)
	if len(warns) != 1 || warns[0].Key != "QUAL_IS_ADES_IND" {
		t.Errorf("GetQualificationWarnings(parent) = %v", warns)
	}

	// A tokenId that does not exist anywhere resolves to nothing.
	if got := r.GetIndication("does-not-exist"); got != "" {
		t.Errorf("GetIndication(missing) = %q, want empty", got)
	}
	if got := r.GetSignatureScopes("does-not-exist"); got != nil {
		t.Errorf("GetSignatureScopes(missing) = %v, want nil", got)
	}
}

// TestSimpleReportWrapper_EAA checks EAA-specific accessors against
// jaxb/testdata/oracle/sr-qeaa.xml (reserialised from an upstream
// dss-simple-report-jaxb test fixture): an EAA token with a nested
// EAASignature and an issuing certificate carrying a TrustAnchor.
func TestSimpleReportWrapper_EAA(t *testing.T) {
	data, err := os.ReadFile(corpustest.RootPath(t, "simplereport/jaxb/testdata/oracle/sr-qeaa.xml"))
	if err != nil {
		t.Fatal(err)
	}
	xmlReport, err := jaxb.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	r := NewSimpleReport(xmlReport)

	const eaaID = "EAA-76E2213B89C691030DDA6D366ED62DA76D62FC5DA891BA93C6D60897A44D260F"
	ids := r.GetEAAIdList()
	if len(ids) != 1 || ids[0] != eaaID {
		t.Fatalf("GetEAAIdList() = %v", ids)
	}
	if got, want := r.GetEAAQualification(eaaID), enumerations.EAAQualificationQEAA; got != want {
		t.Errorf("GetEAAQualification() = %q, want %q", got, want)
	}
	quals := r.GetEAAQualifications(eaaID)
	if len(quals) != 1 || quals[0] != enumerations.EAAQualificationQEAA {
		t.Errorf("GetEAAQualifications() = %v", quals)
	}
	sigs := r.GetEAASignatures(eaaID)
	if len(sigs) != 1 {
		t.Fatalf("GetEAASignatures() = %v", sigs)
	}
	if got, want := sigs[0].SignatureFormat.SignatureLevel(), enumerations.SignatureLevelJAdESBaselineB; got != want {
		t.Errorf("EAA signature format = %q, want %q", got, want)
	}
	eaa := r.GetEAAById(eaaID)
	if eaa == nil || eaa.Indication.Indication() != enumerations.IndicationPassed {
		t.Errorf("EAA indication = %v", eaa)
	}
}
