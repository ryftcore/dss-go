// Command 08-custom-policy validates the same signature against two
// policies: the ETSI policy the library ships as its default, and a copy of
// it with one constraint relaxed - showing that a policy is a plain XML
// document you can load, edit and pass to Validate yourself, not something
// only the library's authors can change.
//
// Run it from anywhere:
//
//	go run ./examples/08-custom-policy
package main

import (
	"fmt"
	"log"

	"github.com/ryftcore/dss-go/dss"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/examples/internal/fixtures"
	"github.com/ryftcore/dss-go/dss/policy"
	"github.com/ryftcore/dss-go/dss/policy/jaxb"
)

func main() {
	signer, err := dss.OpenPKCS12(fixtures.Path("signer_rsa.p12"), "testpassword")
	if err != nil {
		log.Fatalf("opening key store: %v", err)
	}
	defer signer.Close()

	doc := dss.NewDocument("invoice.xml", []byte("<invoice><total>42</total></invoice>"))
	signed, err := dss.Sign(doc, signer, dss.SignOptions{Format: dss.FormatXAdES, Level: dss.LevelB})
	if err != nil {
		log.Fatalf("signing: %v", err)
	}

	// Deliberately without a trust anchor: the point of this example is the
	// policy, not the trust configuration. See example 02 for what an
	// anchored validation of the same kind of signature looks like.
	withDefault, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatalf("validating with the default policy: %v", err)
	}
	v := withDefault.Verdicts()[0]
	fmt.Println("default ETSI policy:", v.Indication, v.SubIndication)

	// ValidateOptions.Policy accepts any constraint document in the DSS
	// policy schema. Here the starting point is the exact policy the
	// library falls back to when Policy is left unset - loaded with
	// policy.ValidationPolicyFacade, the same type the port's own default
	// uses - so the only difference from "default policy" above is the one
	// constraint this example changes.
	facade := policy.NewValidationPolicyFacade()
	constraints, err := facade.UnmarshalFile(fixtures.Module("policy/resources/constraint.xml"))
	if err != nil {
		log.Fatalf("loading the default policy: %v", err)
	}

	// ProspectiveCertificateChain is what turns an unanchored certificate
	// chain into a hard FAIL. Relaxing it to WARN does not make the
	// signature trustworthy - it still reaches no trust anchor - it changes
	// what the validation process DOES about that fact. A real deployment
	// reaches for this only with a documented reason; it exists here to
	// make the effect of a policy constraint visible without needing a
	// second key store.
	constraints.SignatureConstraints.BasicSignatureConstraints.ProspectiveCertificateChain.Level =
		jaxb.LevelValue(enumerations.LevelWarn)

	relaxed, err := facade.Marshal(constraints)
	if err != nil {
		log.Fatalf("marshalling the relaxed policy: %v", err)
	}

	withCustom, err := dss.Validate(signed, dss.ValidateOptions{
		Policy: dss.NewDocument("relaxed-policy.xml", relaxed),
	})
	if err != nil {
		log.Fatalf("validating with the relaxed policy: %v", err)
	}
	cv := withCustom.Verdicts()[0]
	fmt.Println("relaxed policy:        ", cv.Indication, cv.SubIndication)
	if len(cv.Warnings) > 0 {
		fmt.Println("  warning:", cv.Warnings[0].Value)
	}
}
