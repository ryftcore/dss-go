package cv

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignatureIntactWithIdCheck against the Java oracle
// (testdata/oracle/cv_direct.jsonl): the check CryptographicVerification never
// wires, driven over a real corpus signature so that its overridden
// additional-info (the token id) is pinned.
func TestSignatureIntactWithIdCheckAgainstJavaOracle(t *testing.T) {
	const file = "CAdESDoubleLTA.p7m.xml"
	const id = "S-09D72AB7D98820997F96ADB5FAB4D05FC9E8CA1895A6609C3DE9ADB251CC3EC7"
	signature := corpusSignature(t, file, id)
	assertDirectRow(t, id,
		func(result *process.Result[*jaxb.XmlCV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlCV] {
			return NewSignatureIntactWithIdCheck(i18nProviderForTests, result, signature,
				enumerations.Context_SIGNATURE, rule)
		})
}
