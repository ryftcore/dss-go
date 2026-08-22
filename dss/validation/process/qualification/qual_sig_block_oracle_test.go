package qualification

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
)

// The top-level KAT for SignatureQualificationBlock. Every row of
// testdata/oracle/qual_sig_block.jsonl is the whole
// XmlValidationSignatureQualification upstream produces for one (signing
// certificate + attached trust-service providers, XmlTLAnalysis panel,
// EN 319 102-1 conclusion, best-signature-time) tuple - see
// testdata/gen/QualSigBlockOracle.java.
//
// This covers the LOTL/TL acceptance loop, the by-URL trust-service filter,
// both nested CertQualificationAtTimeBlocks, the Article-32 checks,
// FinalCertificateQualificationCalculator and SigQualificationMatrix, plus
// setIndication()/collectAdditionalMessages(). Several rows land on the
// "no acceptable trust service survives the filters" shape, on which upstream
// reads getFilteredServices() as an EMPTY LIST rather than null.

const qualSigBlockOracleLog = "testdata/oracle/qual_sig_block.jsonl"

type qualSigNestedCertQual struct {
	ValidationTime           *string                `json:"validationTime"`
	CertificateQualification *string                `json:"certificateQualification"`
	Constraints              []*qualBlockConstraint `json:"constraints"`
	Conclusion               *qualBlockConclusion   `json:"conclusion"`
}

type qualSigBlockResult struct {
	Title                     *string                  `json:"title"`
	SignatureQualification    *string                  `json:"signatureQualification"`
	Constraints               []*qualBlockConstraint   `json:"constraints"`
	Conclusion                *qualBlockConclusion     `json:"conclusion"`
	CertificateQualifications []*qualSigNestedCertQual `json:"certificateQualifications"`
}

type qualSigBlockRow struct {
	Kind              string              `json:"kind"`
	I                 int                 `json:"i"`
	Cert              string              `json:"cert"`
	TSP               string              `json:"tsp"`
	TLAnalyses        string              `json:"tlAnalyses"`
	Etsi              string              `json:"etsi"`
	BestSignatureTime int64               `json:"bestSignatureTime"`
	Result            *qualSigBlockResult `json:"result"`
}

const (
	qualSigTLUrl      = "https://tl.example/LU"
	qualSigLOTLUrl    = "https://lotl.example/EU"
	qualSigOtherTLUrl = "https://tl.example/DE"
)

// qualSigTSPs mirrors QualSigBlockOracle.tsps(String).
func qualSigTSPs(label string) []*diagnosticjaxb.XmlTrustServiceProvider {
	if label == "none" {
		return nil
	}
	tsp := &diagnosticjaxb.XmlTrustServiceProvider{}
	tsp.TSPNames = &diagnosticjaxb.TSPNamesWrapper{
		Items: []*diagnosticjaxb.XmlLangAndValue{langAndValue("en", "TSP "+label)}}

	tlURL := qualSigTLUrl
	if label == "granted-other-tl" {
		tlURL = qualSigOtherTLUrl
	}
	tl := &diagnosticjaxb.XmlTrustedList{}
	tl.Url = &tlURL
	cc := "LU"
	tl.CountryCode = &cc
	tsp.TL = tl
	if label == "granted-lotl" {
		lotlURL := qualSigLOTLUrl
		lotl := &diagnosticjaxb.XmlTrustedList{}
		lotl.Url = &lotlURL
		isLOTL := true
		lotl.LOTL = &isLOTL
		tsp.LOTL = lotl
	}

	const (
		granted   = "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted"
		withdrawn = "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/withdrawn"
		caQc      = "http://uri.etsi.org/TrstSvc/Svctype/CA/QC"
		caPkc     = "http://uri.etsi.org/TrstSvc/Svctype/CA/PKC"
		qcStmt    = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCStatement"
		notQual   = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotQualified"
		withQSCD  = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/QCWithQSCD"
	)
	early := qualBlockEarly
	var services []*diagnosticjaxb.XmlTrustService
	switch label {
	case "two-conflicting":
		services = append(services,
			qualSigTrustService("svc-a", granted, caQc, qualBlockPostEIDAS, nil, []string{qcStmt}),
			qualSigTrustService("svc-b", granted, caQc, qualBlockPostEIDAS, nil, []string{notQual}))
	case "withdrawn-tl":
		services = append(services, qualSigTrustService("svc", withdrawn, caQc, qualBlockPostEIDAS, nil, nil))
	case "capkc-tl":
		services = append(services, qualSigTrustService("svc", granted, caPkc, qualBlockPostEIDAS, nil, nil))
	case "expired-tl":
		services = append(services, qualSigTrustService("svc", granted, caQc, qualBlockPreEIDAS, &early, nil))
	case "granted-tl-qscd":
		services = append(services, qualSigTrustService("svc", granted, caQc, qualBlockPostEIDAS, nil, []string{withQSCD}))
	default:
		services = append(services, qualSigTrustService("svc", granted, caQc, qualBlockPostEIDAS, nil, nil))
	}
	tsp.TrustServices = &diagnosticjaxb.TrustServicesWrapper{Items: services}
	return []*diagnosticjaxb.XmlTrustServiceProvider{tsp}
}

func langAndValue(lang, value string) *diagnosticjaxb.XmlLangAndValue {
	v := &diagnosticjaxb.XmlLangAndValue{}
	v.Lang = &lang
	v.Value = value
	return v
}

func qualSigTrustService(name, status, svcType string, start int64, end *int64,
	qualifiers []string) *diagnosticjaxb.XmlTrustService {
	svc := &diagnosticjaxb.XmlTrustService{}
	svc.ServiceNames = &diagnosticjaxb.ServiceNamesWrapper{
		Items: []*diagnosticjaxb.XmlLangAndValue{langAndValue("en", name)}}
	svc.Status = &status
	svc.ServiceType = &svcType
	startDT := diagnosticjaxb.XSDateTime(time.UnixMilli(start).UTC())
	svc.StartDate = &startDT
	if end != nil {
		endDT := diagnosticjaxb.XSDateTime(time.UnixMilli(*end).UTC())
		svc.EndDate = &endDT
	}
	items := make([]*diagnosticjaxb.XmlQualifier, 0, len(qualifiers))
	for _, uri := range qualifiers {
		items = append(items, &diagnosticjaxb.XmlQualifier{Value: uri})
	}
	svc.CapturedQualifiers = &diagnosticjaxb.CapturedQualifiersWrapper{Items: items}
	svc.AdditionalServiceInfoUris = &diagnosticjaxb.AdditionalServiceInfoUrisWrapper{
		Items: []string{"http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSignatures"}}
	sdi := &diagnosticjaxb.XmlCertificate{}
	sdiID := diagnosticjaxb.CollapsedString("C-SDI-" + name)
	sdi.Id = &sdiID
	nb := diagnosticjaxb.XSDateTime(time.UnixMilli(qualBlockPreEIDAS).UTC())
	na := diagnosticjaxb.XSDateTime(time.UnixMilli(1577836800000).UTC())
	sdi.NotBefore = &nb
	sdi.NotAfter = &na
	svc.ServiceDigitalIdentifier = sdi
	return svc
}

// qualSigCertificate mirrors QualSigBlockOracle.certificate(String, String).
func qualSigCertificate(label, tspLabel string) *diagnostic.CertificateWrapper {
	post := len(label) < 4 || label[len(label)-4:] != "-pre"
	notBefore := qualBlockPreEIDAS
	if post {
		notBefore = qualBlockPostEIDAS
	}
	var types []string
	if contains(label, "esign") {
		types = []string{"0.4.0.1862.1.6.1"}
	}
	compliance := len(label) >= 3 && label[:3] == "qc-"
	sscd := contains(label, "qscd")
	xml := buildQualCertXml(&qualCertInput{
		QcCompliance:        &compliance,
		QcSSCD:              &sscd,
		QcTypes:             types,
		NotBefore:           notBefore,
		QcStatementsPresent: true,
	})
	xml.TrustServiceProviders = &diagnosticjaxb.TrustServiceProvidersWrapper{Items: qualSigTSPs(tspLabel)}
	return diagnostic.NewCertificateWrapper(xml)
}

// qualSigTLAnalyses mirrors QualSigBlockOracle.tlAnalyses(String).
func qualSigTLAnalyses(label string) []*jaxb.XmlTLAnalysis {
	switch label {
	case "none":
		return nil
	case "tl-ok":
		return []*jaxb.XmlTLAnalysis{qualSigTLAnalysis(qualSigTLUrl, enumerations.Indication_PASSED, false, false)}
	case "tl-failed":
		return []*jaxb.XmlTLAnalysis{qualSigTLAnalysis(qualSigTLUrl, enumerations.Indication_FAILED, true, false)}
	case "tl-warn":
		return []*jaxb.XmlTLAnalysis{qualSigTLAnalysis(qualSigTLUrl, enumerations.Indication_PASSED, false, true)}
	case "lotl-and-tl-ok":
		return []*jaxb.XmlTLAnalysis{
			qualSigTLAnalysis(qualSigLOTLUrl, enumerations.Indication_PASSED, false, false),
			qualSigTLAnalysis(qualSigTLUrl, enumerations.Indication_PASSED, false, false)}
	case "lotl-failed":
		return []*jaxb.XmlTLAnalysis{
			qualSigTLAnalysis(qualSigLOTLUrl, enumerations.Indication_FAILED, true, false),
			qualSigTLAnalysis(qualSigTLUrl, enumerations.Indication_PASSED, false, false)}
	}
	panic("unknown tl-analysis label " + label)
}

func qualSigTLAnalysis(url string, indication enumerations.Indication, withError, withWarning bool) *jaxb.XmlTLAnalysis {
	analysis := &jaxb.XmlTLAnalysis{}
	analysis.URL = url
	cc := "LU"
	analysis.CountryCode = &cc
	analysis.Title = "TL " + url
	conclusion := &jaxb.XmlConclusion{Indication: jaxb.IndicationValue(indication)}
	if withError {
		key := "QUAL_TL_EXP_ANS"
		conclusion.Errors = append(conclusion.Errors, &jaxb.XmlMessage{Key: &key, Value: "the trusted list has expired"})
	}
	if withWarning {
		key := "QUAL_TL_FRESH_ANS"
		conclusion.Warnings = append(conclusion.Warnings, &jaxb.XmlMessage{Key: &key, Value: "the trusted list is not fresh"})
	}
	analysis.Conclusion = conclusion
	return analysis
}

// qualSigEtsiResult mirrors QualSigBlockOracle.etsiResult(String, long).
func qualSigEtsiResult(label string, bestSignatureTime int64) *jaxb.XmlConstraintsConclusionWithProofOfExistence {
	result := &jaxb.XmlConstraintsConclusionWithProofOfExistence{}
	conclusion := &jaxb.XmlConclusion{}
	switch label {
	case "passed":
		conclusion.Indication = jaxb.IndicationValue(enumerations.Indication_PASSED)
	case "indeterminate":
		conclusion.Indication = jaxb.IndicationValue(enumerations.Indication_INDETERMINATE)
		sub := jaxb.SubIndicationValue(enumerations.SubIndication_TRY_LATER)
		conclusion.SubIndication = &sub
	default:
		conclusion.Indication = jaxb.IndicationValue(enumerations.Indication_TOTAL_FAILED)
		sub := jaxb.SubIndicationValue(enumerations.SubIndication_HASH_FAILURE)
		conclusion.SubIndication = &sub
	}
	result.Conclusion = conclusion
	result.ProofOfExistence = &jaxb.XmlProofOfExistence{
		Time: jaxb.XSDateTime(time.UnixMilli(bestSignatureTime).UTC())}
	return result
}

func readQualSigBlockOracle(t *testing.T) []*qualSigBlockRow {
	t.Helper()
	f, err := os.Open(corpustest.Path(t, "oracle/qual_sig_block.jsonl"))
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer f.Close()
	var rows []*qualSigBlockRow
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for sc.Scan() {
		if len(sc.Bytes()) == 0 {
			continue
		}
		row := &qualSigBlockRow{}
		if err := json.Unmarshal(sc.Bytes(), row); err != nil {
			t.Fatalf("decode oracle row: %v", err)
		}
		rows = append(rows, row)
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan oracle: %v", err)
	}
	return rows
}

func TestQualSigBlockOracle(t *testing.T) {
	rows := readQualSigBlockOracle(t)
	if len(rows) == 0 {
		t.Fatal("empty signature qualification block oracle corpus")
	}
	provider := i18n.NewI18nProvider()
	for _, row := range rows {
		row := row
		if row.Kind != "sigQual" {
			t.Fatalf("unknown oracle row kind %q", row.Kind)
		}
		if row.Result == nil {
			t.Fatalf("row %d: upstream threw; the Go port has no equivalent recorded", row.I)
		}
		name := fmt.Sprintf("%d/%s/%s/%s/%s", row.I, row.Cert, row.TSP, row.TLAnalyses, row.Etsi)
		t.Run(name, func(t *testing.T) {
			cert := qualSigCertificate(row.Cert, row.TSP)
			block := NewSignatureQualificationBlock(provider,
				qualSigEtsiResult(row.Etsi, row.BestSignatureTime), cert, qualSigTLAnalyses(row.TLAnalyses))
			result := block.Execute()

			want := row.Result
			if want.Title != nil && result.Title != *want.Title {
				t.Errorf("Title = %q, want %q", result.Title, *want.Title)
			}
			gotQualification := string(result.SignatureQualification.SignatureQualification())
			if gotQualification != strOrEmpty(want.SignatureQualification) {
				t.Errorf("SignatureQualification = %q, want %q",
					gotQualification, strOrEmpty(want.SignatureQualification))
			}
			assertQualBlockConstraints(t, result.Constraint, want.Constraints)
			assertQualBlockConclusion(t, result.Conclusion, want.Conclusion)

			if len(result.ValidationCertificateQualification) != len(want.CertificateQualifications) {
				t.Fatalf("nested certificate qualifications = %d, want %d",
					len(result.ValidationCertificateQualification), len(want.CertificateQualifications))
			}
			for i, w := range want.CertificateQualifications {
				g := result.ValidationCertificateQualification[i]
				gotTime := ""
				if g.ValidationTime != nil {
					gotTime = string(g.ValidationTime.ValidationTime())
				}
				if gotTime != strOrEmpty(w.ValidationTime) {
					t.Errorf("nested[%d].ValidationTime = %q, want %q", i, gotTime, strOrEmpty(w.ValidationTime))
				}
				gotCertQual := ""
				if g.CertificateQualification != nil {
					gotCertQual = string(g.CertificateQualification.CertificateQualification())
				}
				if gotCertQual != strOrEmpty(w.CertificateQualification) {
					t.Errorf("nested[%d].CertificateQualification = %q, want %q",
						i, gotCertQual, strOrEmpty(w.CertificateQualification))
				}
				assertQualBlockConstraints(t, g.Constraint, w.Constraints)
				assertQualBlockConclusion(t, g.Conclusion, w.Conclusion)
			}
		})
	}
}
