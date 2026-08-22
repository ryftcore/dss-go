// Ported from dss-policy-jaxb/.../policy/EtsiValidationPolicyFactory.java
// (DSS 6.5.RC1).
package policy

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
)

// defaultValidationPolicy is a byte-identical copy of upstream's
// src/main/resources/policy/constraint.xml, embedded because
// EtsiValidationPolicyFactory#loadDefaultValidationPolicy loads it from the
// classpath at runtime (DEFAULT_VALIDATION_POLICY_LOCATION =
// "/policy/constraint.xml"), unlike policy.xsd (see
// validation_policy_xml_definer.go), whose loader is a deferred stub that
// never reads it.
//
//go:embed resources/constraint.xml
var defaultValidationPolicy []byte

// EtsiValidationPolicyFactory loads a DSS XML Validation Policy based on the
// ETSI TS 102 853 standard. Ports EtsiValidationPolicyFactory.
type EtsiValidationPolicyFactory struct{}

var _ modelpolicy.ValidationPolicyFactory = (*EtsiValidationPolicyFactory)(nil)

// NewEtsiValidationPolicyFactory is the default constructor.
func NewEtsiValidationPolicyFactory() *EtsiValidationPolicyFactory {
	return &EtsiValidationPolicyFactory{}
}

// IsSupported ports EtsiValidationPolicyFactory#isSupported.
//
// Java swallows IOException/JAXBException/XMLStreamException/SAXException;
// this port swallows the corresponding Unmarshal error the same way, per the
// method's job of probing whether validationPolicyDocument parses at all.
func (f *EtsiValidationPolicyFactory) IsSupported(validationPolicyDocument model.DSSDocument) bool {
	rc, err := validationPolicyDocument.OpenStream()
	if err != nil {
		return false
	}
	defer rc.Close()

	cp, err := NewValidationPolicyFacade().Unmarshal(rc)
	return err == nil && cp != nil
}

// LoadDefaultValidationPolicy ports
// EtsiValidationPolicyFactory#loadDefaultValidationPolicy.
//
// Java wraps any failure in an unchecked UnsupportedOperationException; this
// port panics with the same message shape, matching the
// model/policy.ValidationPolicyFactory interface, whose loader methods
// return no error (see that file).
func (f *EtsiValidationPolicyFactory) LoadDefaultValidationPolicy() modelpolicy.ValidationPolicy {
	return f.LoadValidationPolicyFromReader(bytes.NewReader(defaultValidationPolicy))
}

// LoadValidationPolicy ports
// EtsiValidationPolicyFactory#loadValidationPolicy(DSSDocument).
func (f *EtsiValidationPolicyFactory) LoadValidationPolicy(validationPolicyDocument model.DSSDocument) modelpolicy.ValidationPolicy {
	rc, err := validationPolicyDocument.OpenStream()
	if err != nil {
		panic(fmt.Sprintf("Unable to load the default policy document. Reason : %s", err.Error()))
	}
	return f.LoadValidationPolicyFromReader(rc)
}

// LoadValidationPolicyFromReader ports
// EtsiValidationPolicyFactory#loadValidationPolicy(InputStream). Named per
// PORTING.md's Java-overload convention (see
// model/policy.ValidationPolicyFactory's doc comment).
//
// Java wraps any failure (including the try-with-resources close) in an
// unchecked UnsupportedOperationException; ported as a panic, per the
// interface's error-less return (see LoadDefaultValidationPolicy's doc
// comment).
func (f *EtsiValidationPolicyFactory) LoadValidationPolicyFromReader(validationPolicyInputStream io.Reader) modelpolicy.ValidationPolicy {
	if closer, ok := validationPolicyInputStream.(io.Closer); ok {
		defer closer.Close()
	}
	policy, err := NewValidationPolicyFacade().GetValidationPolicy(validationPolicyInputStream)
	if err != nil {
		panic(fmt.Sprintf("Unable to load the default policy document. Reason : %s", err.Error()))
	}
	return policy
}
