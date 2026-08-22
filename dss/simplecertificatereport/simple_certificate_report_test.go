package simplecertificatereport

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/simplecertificatereport/jaxb"
)

// loadOracle reads a Java-produced oracle dump from the sibling jaxb
// package's testdata/oracle. Most of these fixtures ship in-package; a few
// larger ones live in the external corpus/ instead, so a local miss falls
// through to corpustest.
func loadOracle(t *testing.T, name string) *SimpleCertificateReport {
	t.Helper()
	local := filepath.Join("jaxb", "testdata", "oracle", name)
	data, err := os.ReadFile(local)
	if errors.Is(err, os.ErrNotExist) {
		data, err = os.ReadFile(corpustest.RootPath(t, "simplecertificatereport/jaxb/testdata/oracle/"+name))
	}
	if err != nil {
		t.Fatal(err)
	}
	xmlReport, err := jaxb.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	return NewSimpleCertificateReport(xmlReport)
}

// TestSimpleCertificateReportWrapper checks the wrapper's basic method
// surface against values read directly off a Java-produced oracle dump
// (jaxb/testdata/oracle/good-user.cer.xml): a single leaf certificate with
// no trust chain, no revocation data, and several omitted (not
// self-closed) URL wrapper elements.
func TestSimpleCertificateReportWrapper(t *testing.T) {
	r := loadOracle(t, "good-user.cer.xml")

	const certID = "C-B9B8051A58645938F660EC1261B7534E2DC7422882D873B2876204BBA1078352"

	if ids := r.GetCertificateIds(); len(ids) != 1 || ids[0] != certID {
		t.Fatalf("GetCertificateIds() = %v", ids)
	}
	if got := r.GetCertificateCommonName(certID); got != "good-user" {
		t.Errorf("GetCertificateCommonName() = %q", got)
	}
	if got := r.GetCertificateOrganizationName(certID); got != "Nowina Solutions" {
		t.Errorf("GetCertificateOrganizationName() = %q", got)
	}
	if got := r.GetCertificateCountry(certID); got != "LU" {
		t.Errorf("GetCertificateCountry() = %q", got)
	}
	if got, want := r.GetCertificateIndication(certID), enumerations.IndicationIndeterminate; got != want {
		t.Errorf("GetCertificateIndication() = %q, want %q", got, want)
	}
	if got, want := r.GetCertificateSubIndication(certID), enumerations.SubIndicationNoCertificateChainFound; got != want {
		t.Errorf("GetCertificateSubIndication() = %q, want %q", got, want)
	}
	if got := r.GetCertificateOcspUrls(certID); len(got) != 1 || got[0] != "http://dss.nowina.lu/pki-factory/ocsp/good-ca" {
		t.Errorf("GetCertificateOcspUrls() = %v", got)
	}
	// crlUrls/cpsUrls/pdsUrls/extendedKeyUsages/trustAnchors are entirely
	// omitted in this dump (not self-closed): the wrapper must report them
	// as empty, not panic on a nil pointer.
	if got := r.GetCertificateCrlUrls(certID); len(got) != 0 {
		t.Errorf("GetCertificateCrlUrls() = %v, want empty", got)
	}
	if got := r.GetTrustAnchorVATNumbers(); len(got) != 0 {
		t.Errorf("GetTrustAnchorVATNumbers() = %v, want empty", got)
	}
	errs := r.GetX509ValidationErrors(certID)
	if len(errs) != 1 || errs[0].Key != "BBB_XCV_CCCBB_ANS" {
		t.Errorf("GetX509ValidationErrors() = %v", errs)
	}
	if got := r.GetCertificateIds(); len(got) == 0 {
		t.Fatal("no certificates")
	}
	// A certificateId that resolves nowhere must not panic.
	if got := r.GetCertificateCommonName("does-not-exist"); got != "" {
		t.Errorf("GetCertificateCommonName(missing) = %q, want empty", got)
	}
	if got := r.GetCertificateNotBefore("does-not-exist"); got != nil {
		t.Errorf("GetCertificateNotBefore(missing) = %v, want nil", got)
	}
}

// TestSimpleCertificateReportWrapper_ApprovalStatus checks the
// CertificateApprovalStatus accessors against
// jaxb/testdata/oracle/simple-cert-report-pid.xml, which carries a "PID
// Provider" approval status at both issuance and validation time plus a
// resolved trust chain.
func TestSimpleCertificateReportWrapper_ApprovalStatus(t *testing.T) {
	r := loadOracle(t, "simple-cert-report-pid.xml")

	statuses := r.GetCertificateApprovalStatusAtCertificateIssuance()
	if len(statuses) != 1 {
		t.Fatalf("GetCertificateApprovalStatusAtCertificateIssuance() = %v", statuses)
	}
	if got := statuses[0].Label(); got == "" {
		t.Errorf("approval status label is empty")
	}

	certID := r.GetCertificateIds()[0]
	found := statuses[0]
	errs := r.GetCertificateApprovalStatusErrorsAtIssuanceTime(certID, found)
	warns := r.GetCertificateApprovalStatusWarningsAtIssuanceTime(certID, found)
	info := r.GetCertificateApprovalStatusInfoAtIssuanceTime(certID, found)
	if len(errs) != 0 || len(warns) != 0 || len(info) != 0 {
		t.Errorf("approval status details at issuance = errs=%v warns=%v info=%v, want all empty (Details is self-closed)", errs, warns, info)
	}

	// A CertificateApprovalStatus that does not match anything in the
	// document resolves to no Details at all.
	unrelated := enumerations.NewCertificateApprovalStatus("unrelated", nil, nil, nil)
	if got := r.GetCertificateApprovalStatusErrorsAtIssuanceTime(certID, unrelated); got != nil {
		t.Errorf("GetCertificateApprovalStatusErrorsAtIssuanceTime(unrelated) = %v, want nil", got)
	}

	if got := r.GetCertificateNotBefore(certID); got == nil {
		t.Error("GetCertificateNotBefore() = nil, want a date")
	}
}
