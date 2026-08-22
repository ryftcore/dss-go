// Direct unit tests for the vpfltvd package's leaf ChainItem checks.
//
// TESTING NOTE (LTVB, phase 8e): see vpfbs/vpfbs_direct_test.go's header -
// same rationale: hand-constructed OK/NOT-OK cases in place of a
// Java-oracle-generated corpus, which this pass did not have time to set up.
package vpfltvd

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

func newTestI18nProvider() *i18n.I18nProvider { return i18n.NewI18nProvider() }

func newTestResult() *process.Result[*jaxb.XmlValidationProcessLongTermData] {
	xmlResult := &jaxb.XmlValidationProcessLongTermData{}
	return process.NewResult(xmlResult, &xmlResult.XmlConstraintsConclusionWithProofOfExistenceContent.XmlConstraintsConclusionContent,
		&xmlResult.XmlConstraintsConclusionAttrs)
}

func conclusionWith(indication enumerations.Indication, subIndication enumerations.SubIndication) *jaxb.XmlConclusion {
	c := &jaxb.XmlConclusion{Indication: jaxb.IndicationValue(indication)}
	if subIndication != "" {
		v := jaxb.SubIndicationValue(subIndication)
		c.SubIndication = &v
	}
	return c
}

func passedConclusion() *jaxb.XmlConclusion {
	return conclusionWith(enumerations.Indication_PASSED, "")
}

type fixedLevelRule struct{ level enumerations.Level }

func (f *fixedLevelRule) Level() enumerations.Level { return f.level }

func failLevel() *fixedLevelRule { return &fixedLevelRule{level: enumerations.Level_FAIL} }

func newTestCertificate(id string, notBefore, notAfter *time.Time) *diagnostic.CertificateWrapper {
	cid := diagjaxb.CollapsedString(id)
	cert := &diagjaxb.XmlCertificate{XmlAbstractTokenAttrs: diagjaxb.XmlAbstractTokenAttrs{Id: &cid}}
	if notBefore != nil {
		cert.NotBefore = diagjaxb.NewXSDateTime(*notBefore)
	}
	if notAfter != nil {
		cert.NotAfter = diagjaxb.NewXSDateTime(*notAfter)
	}
	return diagnostic.NewCertificateWrapper(cert)
}

func newTestCertificateRevocation(id string, revocationDate *time.Time) *diagnostic.CertificateRevocationWrapper {
	cid := diagjaxb.CollapsedString(id)
	rev := &diagjaxb.XmlCertificateRevocation{Revocation: &diagjaxb.XmlRevocation{XmlAbstractTokenAttrs: diagjaxb.XmlAbstractTokenAttrs{Id: &cid}}}
	if revocationDate != nil {
		rev.RevocationDate = diagjaxb.NewXSDateTime(*revocationDate)
	}
	return diagnostic.NewCertificateRevocationWrapper(rev)
}

func newTestSignature(claimedSigningTime *time.Time) *diagnostic.SignatureWrapper {
	sig := &diagjaxb.XmlSignature{}
	if claimedSigningTime != nil {
		sig.ClaimedSigningTime = diagjaxb.NewXSDateTime(*claimedSigningTime)
	}
	return diagnostic.NewSignatureWrapper(sig)
}

func newTestTimestamp(id string, productionTime *time.Time, tstType enumerations.TimestampType) *diagnostic.TimestampWrapper {
	cid := diagjaxb.CollapsedString(id)
	tst := &diagjaxb.XmlTimestamp{XmlAbstractTokenAttrs: diagjaxb.XmlAbstractTokenAttrs{Id: &cid}}
	if productionTime != nil {
		tst.ProductionTime = diagjaxb.NewXSDateTime(*productionTime)
	}
	if tstType != "" {
		v := diagjaxb.TimestampTypeValue(tstType)
		tst.Type = &v
	}
	return diagnostic.NewTimestampWrapper(tst)
}

func TestRevocationDataAcceptableCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()

	t.Run("allowed basic revocation validation passes", func(t *testing.T) {
		c := NewRevocationDataAcceptableCheck(i18nProvider, newTestResult(), "rev1", passedConclusion(), failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true for a PASSED revocation basic validation")
		}
	})

	t.Run("disallowed conclusion fails", func(t *testing.T) {
		conclusion := conclusionWith(enumerations.Indication_INDETERMINATE, enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE)
		c := NewRevocationDataAcceptableCheck(i18nProvider, newTestResult(), "rev1", conclusion, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false for a disallowed conclusion")
		}
		if got := c.FailedIndicationForConclusion(); got != enumerations.Indication_INDETERMINATE {
			t.Errorf("FailedIndicationForConclusion() = %v, want INDETERMINATE", got)
		}
	})
}

func TestAcceptableBasicSignatureValidationCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()

	t.Run("passed conclusion", func(t *testing.T) {
		content := &jaxb.XmlConstraintsConclusionContent{Conclusion: passedConclusion()}
		c := NewAcceptableBasicSignatureValidationCheck(i18nProvider, newTestResult(), content, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true for a PASSED basic signature validation")
		}
	})

	t.Run("nil conclusion fails", func(t *testing.T) {
		content := &jaxb.XmlConstraintsConclusionContent{}
		c := NewAcceptableBasicSignatureValidationCheck(i18nProvider, newTestResult(), content, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false when the conclusion is nil")
		}
	})
}

func TestBestSignatureTimeBeforeCertificateExpirationCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	notAfter := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	cert := newTestCertificate("cert1", nil, &notAfter)

	t.Run("before expiration passes", func(t *testing.T) {
		bst := notAfter.Add(-time.Hour)
		c := NewBestSignatureTimeBeforeCertificateExpirationCheck[*jaxb.XmlValidationProcessLongTermData](
			i18nProvider, newTestResult(), &bst, cert, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true when best-signature-time is before expiration")
		}
	})

	t.Run("after expiration fails", func(t *testing.T) {
		bst := notAfter.Add(time.Hour)
		c := NewBestSignatureTimeBeforeCertificateExpirationCheck[*jaxb.XmlValidationProcessLongTermData](
			i18nProvider, newTestResult(), &bst, cert, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false when best-signature-time is after expiration")
		}
		if got := c.FailedSubIndicationForConclusion(); got != enumerations.SubIndication_OUT_OF_BOUNDS_NOT_REVOKED {
			t.Errorf("FailedSubIndicationForConclusion() = %v, want OUT_OF_BOUNDS_NOT_REVOKED", got)
		}
	})
}

func TestBestSignatureTimeNotBeforeCertificateIssuanceCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	notBefore := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	cert := newTestCertificate("cert1", &notBefore, nil)

	t.Run("after issuance passes", func(t *testing.T) {
		bst := notBefore.Add(time.Hour)
		c := NewBestSignatureTimeNotBeforeCertificateIssuanceCheck[*jaxb.XmlValidationProcessLongTermData](
			i18nProvider, newTestResult(), &bst, cert, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true when best-signature-time is after issuance")
		}
	})

	t.Run("before issuance fails as NOT_YET_VALID", func(t *testing.T) {
		bst := notBefore.Add(-time.Hour)
		c := NewBestSignatureTimeNotBeforeCertificateIssuanceCheck[*jaxb.XmlValidationProcessLongTermData](
			i18nProvider, newTestResult(), &bst, cert, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false when best-signature-time is before issuance")
		}
		if got := c.FailedIndicationForConclusion(); got != enumerations.Indication_FAILED {
			t.Errorf("FailedIndicationForConclusion() = %v, want FAILED", got)
		}
		if got := c.FailedSubIndicationForConclusion(); got != enumerations.SubIndication_NOT_YET_VALID {
			t.Errorf("FailedSubIndicationForConclusion() = %v, want NOT_YET_VALID", got)
		}
	})
}

func TestRevocationDateAfterBestSignatureTimeCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	bst := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("revocation after best-signature-time fails (signing cert)", func(t *testing.T) {
		revDate := bst.Add(time.Hour)
		rev := newTestCertificateRevocation("rev1", &revDate)
		c := NewRevocationDateAfterBestSignatureTimeCheck(i18nProvider, newTestResult(), rev, &bst, failLevel(), enumerations.SubContext_SIGNING_CERT)
		if !c.Process() {
			t.Fatalf("expected Process() to be true (i.e. revocation IS after bst) triggering the failure branch")
		}
		if got := c.FailedSubIndicationForConclusion(); got != enumerations.SubIndication_REVOKED_NO_POE {
			t.Errorf("FailedSubIndicationForConclusion() = %v, want REVOKED_NO_POE for SIGNING_CERT", got)
		}
	})

	t.Run("revocation before best-signature-time does not trigger", func(t *testing.T) {
		revDate := bst.Add(-time.Hour)
		rev := newTestCertificateRevocation("rev1", &revDate)
		c := NewRevocationDateAfterBestSignatureTimeCheck(i18nProvider, newTestResult(), rev, &bst, failLevel(), enumerations.SubContext_CA_CERTIFICATE)
		if c.Process() {
			t.Fatalf("expected Process() to be false when revocation predates best-signature-time")
		}
		if got := c.FailedSubIndicationForConclusion(); got != enumerations.SubIndication_REVOKED_CA_NO_POE {
			t.Errorf("FailedSubIndicationForConclusion() = %v, want REVOKED_CA_NO_POE for CA_CERTIFICATE", got)
		}
	})
}

func TestBestSignatureTimeBeforeSuspensionTimeCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	bst := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("bst before suspension passes", func(t *testing.T) {
		revDate := bst.Add(time.Hour)
		rev := newTestCertificateRevocation("rev1", &revDate)
		c := NewBestSignatureTimeBeforeSuspensionTimeCheck(i18nProvider, newTestResult(), rev, &bst, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true when best-signature-time is before the suspension date")
		}
	})

	t.Run("bst after suspension fails", func(t *testing.T) {
		revDate := bst.Add(-time.Hour)
		rev := newTestCertificateRevocation("rev1", &revDate)
		c := NewBestSignatureTimeBeforeSuspensionTimeCheck(i18nProvider, newTestResult(), rev, &bst, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false when best-signature-time is after the suspension date")
		}
	})
}

func TestSigningTimeAttributePresentCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()

	t.Run("present passes", func(t *testing.T) {
		now := time.Now()
		sig := newTestSignature(&now)
		c := NewSigningTimeAttributePresentCheck(i18nProvider, newTestResult(), sig, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true when a claimed signing time is present")
		}
	})

	t.Run("absent fails", func(t *testing.T) {
		sig := newTestSignature(nil)
		c := NewSigningTimeAttributePresentCheck(i18nProvider, newTestResult(), sig, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false when no claimed signing time is present")
		}
	})
}

func TestTimestampCoherenceOrderCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	t0 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("single timestamp always coherent", func(t *testing.T) {
		ts := []*diagnostic.TimestampWrapper{newTestTimestamp("tst1", &t0, enumerations.TimestampType_SIGNATURE_TIMESTAMP)}
		c := NewTimestampCoherenceOrderCheck(i18nProvider, newTestResult(), ts, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true for a single timestamp")
		}
	})

	t.Run("different types with coherent time order passes", func(t *testing.T) {
		t1 := t0.Add(time.Hour)
		ts := []*diagnostic.TimestampWrapper{
			newTestTimestamp("cts1", &t0, enumerations.TimestampType_CONTENT_TIMESTAMP),
			newTestTimestamp("sts1", &t1, enumerations.TimestampType_SIGNATURE_TIMESTAMP),
		}
		c := NewTimestampCoherenceOrderCheck(i18nProvider, newTestResult(), ts, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true when a later timestamp type comes after an earlier one")
		}
	})

	t.Run("out-of-order types fails", func(t *testing.T) {
		t1 := t0.Add(time.Hour)
		// A signature timestamp produced before a content timestamp is incoherent.
		ts := []*diagnostic.TimestampWrapper{
			newTestTimestamp("sts1", &t0, enumerations.TimestampType_SIGNATURE_TIMESTAMP),
			newTestTimestamp("cts1", &t1, enumerations.TimestampType_CONTENT_TIMESTAMP),
		}
		c := NewTimestampCoherenceOrderCheck(i18nProvider, newTestResult(), ts, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false when timestamp types are produced out of order")
		}
		if got := c.FailedSubIndicationForConclusion(); got != enumerations.SubIndication_TIMESTAMP_ORDER_FAILURE {
			t.Errorf("FailedSubIndicationForConclusion() = %v, want TIMESTAMP_ORDER_FAILURE", got)
		}
	})
}
