// Direct unit tests for the vpfbs package's leaf ChainItem checks.
//
// These are hand-constructed OK/NOT-OK cases
// exercising each check's Process()/FailedIndicationForConclusion()/
// FailedSubIndicationForConclusion()/BuildAdditionalInfo() against jaxb
// structs built directly in Go, not a Java-oracle-generated corpus (the
// testdata/gen driver + JSONL fixture pattern other packages in this port
// use, e.g. bbb/xcv's xcv_direct_oracle_test.go). Producing that corpus
// requires running the upstream Java classes via Maven against the vpfbs
// package tree, which this pass did not have time to set up; this file
// covers the same check-level behavior by direct construction instead, so
// the package does not ship with zero tests, but the full Java-oracle
// parity corpus remains a follow-up.
package vpfbs

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

func newTestI18nProvider() *i18n.I18nProvider {
	return i18n.NewI18nProvider()
}

func newTestResult() (*jaxb.XmlValidationProcessBasicSignature, *process.Result[*jaxb.XmlValidationProcessBasicSignature]) {
	xmlResult := &jaxb.XmlValidationProcessBasicSignature{}
	result := process.NewResult(xmlResult, &xmlResult.XmlConstraintsConclusionWithProofOfExistenceContent.XmlConstraintsConclusionContent,
		&xmlResult.XmlConstraintsConclusionAttrs)
	return xmlResult, result
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
	return conclusionWith(enumerations.IndicationPassed, "")
}

func failLevel() *fixedLevelRule { return &fixedLevelRule{level: enumerations.LevelFail} }

type fixedLevelRule struct{ level enumerations.Level }

func (f *fixedLevelRule) Level() enumerations.Level { return f.level }

func newTestToken(id string) diagnostic.TokenProxy {
	cid := diagjaxb.CollapsedString(id)
	sig := &diagjaxb.XmlSignature{XmlAbstractTokenAttrs: diagjaxb.XmlAbstractTokenAttrs{Id: &cid}}
	return diagnostic.NewSignatureWrapper(sig)
}

// newTestTimestamp builds a minimal TimestampWrapper with the given Id and
// (optional) production time.
func newTestTimestamp(id string, productionTime *time.Time) *diagnostic.TimestampWrapper {
	cid := diagjaxb.CollapsedString(id)
	tst := &diagjaxb.XmlTimestamp{XmlAbstractTokenAttrs: diagjaxb.XmlAbstractTokenAttrs{Id: &cid}}
	if productionTime != nil {
		tst.ProductionTime = diagjaxb.NewXSDateTime(*productionTime)
	}
	return diagnostic.NewTimestampWrapper(tst)
}

func TestFormatCheckingResultCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	token := newTestToken("sig1")

	t.Run("passed", func(t *testing.T) {
		_, result := newTestResult()
		xmlFC := &jaxb.XmlFC{}
		xmlFC.Conclusion = passedConclusion()
		c := NewFormatCheckingResultCheck(i18nProvider, result, xmlFC, token, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true for a PASSED XmlFC")
		}
	})

	t.Run("failed - nil FC", func(t *testing.T) {
		_, result := newTestResult()
		c := NewFormatCheckingResultCheck(i18nProvider, result, nil, token, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false for a nil XmlFC")
		}
		if got := c.FailedIndicationForConclusion(); got != enumerations.IndicationFailed {
			t.Errorf("FailedIndicationForConclusion() = %v, want FAILED", got)
		}
		if got := c.FailedSubIndicationForConclusion(); got != enumerations.SubIndicationFormatFailure {
			t.Errorf("FailedSubIndicationForConclusion() = %v, want FORMAT_FAILURE", got)
		}
	})
}

func TestIdentificationOfSigningCertificateResultCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	token := newTestToken("sig1")

	t.Run("indeterminate on missing ISC", func(t *testing.T) {
		_, result := newTestResult()
		c := NewIdentificationOfSigningCertificateResultCheck(i18nProvider, result, nil, token, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false for a nil XmlISC")
		}
		if got := c.FailedIndicationForConclusion(); got != enumerations.IndicationIndeterminate {
			t.Errorf("FailedIndicationForConclusion() = %v, want INDETERMINATE", got)
		}
		if got := c.FailedSubIndicationForConclusion(); got != enumerations.SubIndicationNoSigningCertificateFound {
			t.Errorf("FailedSubIndicationForConclusion() = %v, want NO_SIGNING_CERTIFICATE_FOUND", got)
		}
	})
}

func TestSigningCertificateNotRevokedCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	token := newTestToken("sig1")

	cases := []struct {
		name string
		xcv  *jaxb.XmlXCV
		want bool
	}{
		{"passed XCV is not revoked", func() *jaxb.XmlXCV {
			x := &jaxb.XmlXCV{}
			x.Conclusion = passedConclusion()
			return x
		}(), true},
		{"REVOKED_NO_POE fails the check", func() *jaxb.XmlXCV {
			x := &jaxb.XmlXCV{}
			x.Conclusion = conclusionWith(enumerations.IndicationIndeterminate, enumerations.SubIndicationRevokedNoPOE)
			return x
		}(), false},
		{"other indeterminate sub-indication passes the check", func() *jaxb.XmlXCV {
			x := &jaxb.XmlXCV{}
			x.Conclusion = conclusionWith(enumerations.IndicationIndeterminate, enumerations.SubIndicationTryLater)
			return x
		}(), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, result := newTestResult()
			c := NewSigningCertificateNotRevokedCheck(i18nProvider, result, tc.xcv, token, failLevel())
			if got := c.Process(); got != tc.want {
				t.Errorf("Process() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestContentTimestampsCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()

	t.Run("non-empty passes", func(t *testing.T) {
		_, result := newTestResult()
		ts := []*diagnostic.TimestampWrapper{newTestTimestamp("tst1", nil)}
		c := NewContentTimestampsCheck(i18nProvider, result, ts, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true for a non-empty timestamp list")
		}
	})

	t.Run("empty fails", func(t *testing.T) {
		_, result := newTestResult()
		c := NewContentTimestampsCheck(i18nProvider, result, nil, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false for an empty timestamp list")
		}
	})
}

func TestTimestampGenerationTimeNotAfterRevocationTimeCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	revocationTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	t.Run("timestamp before revocation passes", func(t *testing.T) {
		_, result := newTestResult()
		productionTime := revocationTime.Add(-time.Hour)
		tst := newTestTimestamp("tst1", &productionTime)
		c := NewTimestampGenerationTimeNotAfterRevocationTimeCheck(i18nProvider, result, tst, &revocationTime, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true when the timestamp is before the revocation time")
		}
	})

	t.Run("timestamp after revocation fails", func(t *testing.T) {
		_, result := newTestResult()
		productionTime := revocationTime.Add(time.Hour)
		tst := newTestTimestamp("tst1", &productionTime)
		c := NewTimestampGenerationTimeNotAfterRevocationTimeCheck(i18nProvider, result, tst, &revocationTime, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false when the timestamp is after the revocation time")
		}
		if got := c.FailedIndicationForConclusion(); got != enumerations.IndicationFailed {
			t.Errorf("FailedIndicationForConclusion() = %v, want FAILED", got)
		}
		if got := c.FailedSubIndicationForConclusion(); got != enumerations.SubIndicationRevoked {
			t.Errorf("FailedSubIndicationForConclusion() = %v, want REVOKED", got)
		}
	})
}

func TestBasicValidationProcessCheck(t *testing.T) {
	i18nProvider := newTestI18nProvider()
	token := newTestToken("sig1")

	t.Run("valid conclusion passes", func(t *testing.T) {
		_, result := newTestResult()
		c := NewBasicValidationProcessCheck(i18nProvider, result, passedConclusion(), token, failLevel())
		if !c.Process() {
			t.Fatalf("expected Process() to be true for a PASSED conclusion")
		}
		if got := c.BuildAdditionalInfo(); got != nil {
			t.Errorf("BuildAdditionalInfo() = %v, want nil for a valid conclusion", *got)
		}
	})

	t.Run("invalid conclusion fails and builds additional info", func(t *testing.T) {
		_, result := newTestResult()
		conclusion := conclusionWith(enumerations.IndicationIndeterminate, enumerations.SubIndicationTryLater)
		c := NewBasicValidationProcessCheck(i18nProvider, result, conclusion, token, failLevel())
		if c.Process() {
			t.Fatalf("expected Process() to be false for an INDETERMINATE conclusion")
		}
		if got := c.BuildAdditionalInfo(); got == nil {
			t.Errorf("BuildAdditionalInfo() = nil, want a non-nil message")
		}
	})
}
