package diagnostic

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/model"
	modelsignature "github.com/utain/esig/dss/model/signature"
)

// The XmlPolicyBuilder nullable-member KAT.
//
// The expected fragment is lifted verbatim out of
// dss/diagnostic/jaxb/testdata/oracle/Signature-C-CZ_SIX-1.p7m.xml, a JAXB dump of the
// upstream validator: a signature policy with an identifier, a URI, a UserNotice carrying
// only an explicit text, and a zero-hash digest. Java's setters take nullable Strings, so
// the absent description and the absent UserNotice organization leave no element at all.
// The Go accessors return plain strings, so the builder has to map "" back to nil - see
// XmlPolicyBuilder.Build.
const xmlPolicyOracleFragment = `<Policy>
    <Id>2.23.134.1.4.1.8.200</Id>
    <Url>http://www.postsignum.cz</Url>
    <UserNotice>
        <ExplicitText>Tento kvalifikovany systemovy certifikat byl vydan podle zakona 227/2000Sb. a navaznych predpisu/This qualified system certificate was issued according to Law No 227/2000Coll. and related regulations</ExplicitText>
    </UserNotice>
    <Identified>true</Identified>
    <Asn1Processable>false</Asn1Processable>
    <DigestAlgoAndValue zeroHash="true" match="true"/>
</Policy>`

func TestXmlPolicyBuilderOmitsAbsentMembers(t *testing.T) {
	userNotice := model.NewUserNotice()
	userNotice.SetExplicitText("Tento kvalifikovany systemovy certifikat byl vydan podle zakona 227/2000Sb. " +
		"a navaznych predpisu/This qualified system certificate was issued according to Law " +
		"No 227/2000Coll. and related regulations")

	validationResult := modelsignature.NewSignaturePolicyValidationResult()
	validationResult.SetIdentified(true)
	validationResult.SetAsn1Processable(false)
	validationResult.SetDigestValid(true)

	signaturePolicy := modelsignature.NewSignaturePolicyWithIdentifier("2.23.134.1.4.1.8.200")
	signaturePolicy.SetURI("http://www.postsignum.cz")
	signaturePolicy.SetUserNotice(userNotice)
	signaturePolicy.SetZeroHash(true)
	signaturePolicy.SetValidationResult(validationResult)

	xmlPolicy := NewXmlPolicyBuilder(signaturePolicy).Build()

	encoded, err := marshalPolicyFragment(xmlPolicy)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if encoded != xmlPolicyOracleFragment {
		t.Errorf("XmlPolicy differs from the Java dump\nexpected:\n%s\n\nactual:\n%s",
			xmlPolicyOracleFragment, encoded)
	}

	// Guard the individual mappings too, so a failure points at the member.
	if xmlPolicy.Description != nil {
		t.Errorf("Description = %q, want nil (Java leaves a null description unset)", *xmlPolicy.Description)
	}
	if xmlPolicy.DocumentationReferences != nil {
		t.Error("DocumentationReferences should stay unset when the policy carries no list")
	}
	if xmlPolicy.UserNotice == nil {
		t.Fatal("UserNotice should be built")
	}
	if xmlPolicy.UserNotice.Organization != nil {
		t.Errorf("UserNotice.Organization = %q, want nil", *xmlPolicy.UserNotice.Organization)
	}
	if xmlPolicy.ProcessingError != nil {
		t.Errorf("ProcessingError = %q, want nil", *xmlPolicy.ProcessingError)
	}
}

// TestXmlPolicyBuilderKeepsEmptyDocumentationReferences pins the other half of the mapping:
// Java's setDocumentationReferences is unconditional, so a non-null but empty list marshals
// as <DocumentationReferences/> - the spelling the corpus carries in
// jades-with-sigPSt-invalid.json.xml.
func TestXmlPolicyBuilderKeepsEmptyDocumentationReferences(t *testing.T) {
	validationResult := modelsignature.NewSignaturePolicyValidationResult()
	signaturePolicy := modelsignature.NewSignaturePolicyWithIdentifier("1.2.3.4.5.6")
	signaturePolicy.SetDocumentationReferences([]string{})
	signaturePolicy.SetValidationResult(validationResult)

	xmlPolicy := NewXmlPolicyBuilder(signaturePolicy).Build()
	if xmlPolicy.DocumentationReferences == nil {
		t.Fatal("an empty (but present) documentation-reference list must still be built")
	}
	if len(xmlPolicy.DocumentationReferences.Items) != 0 {
		t.Errorf("DocumentationReferences carries %d items, want 0",
			len(xmlPolicy.DocumentationReferences.Items))
	}
}

// marshalPolicyFragment renders an XmlPolicy on its own, four-space indented, the way the
// element appears inside a diagnostic-data document.
func marshalPolicyFragment(xmlPolicy *jaxb.XmlPolicy) (string, error) {
	var b strings.Builder
	enc := xml.NewEncoder(&b)
	enc.Indent("", "    ")
	if err := enc.EncodeElement(xmlPolicy, xml.StartElement{Name: xml.Name{Local: "Policy"}}); err != nil {
		return "", err
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	return strings.ReplaceAll(b.String(), "></DigestAlgoAndValue>", "/>"), nil
}
