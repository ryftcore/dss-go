package cv

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignatureIntactCheck in the CERTIFICATE context against the Java oracle
// (testdata/oracle/cv_direct.jsonl): the one message-tag branch the corpus KAT
// cannot reach, since CryptographicVerification is never run for a certificate.
// The SIGNATURE / TIMESTAMP / REVOCATION branches, and both outcomes, are
// exercised by the corpus KAT (BBB_CV_ISI, BBB_CV_ISIT, BBB_CV_ISIR).
func TestSignatureIntactCheckAgainstJavaOracle(t *testing.T) {
	signature := syntheticSignature()
	assertDirectRow(t, "signature-intact-certificate-context",
		func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV] {
			return NewSignatureIntactCheck(i18nProviderForTests, result, signature,
				enumerations.ContextCertificate, rule)
		})
}
