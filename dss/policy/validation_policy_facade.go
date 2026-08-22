// Ported from dss-policy-jaxb/.../policy/ValidationPolicyFacade.java (DSS
// 6.5.RC1).
//
// Java's ValidationPolicyFacade extends the generic dss-jaxb-common
// AbstractJaxbFacade<T>, a module outside S8A_BRIEF.md's manifest (dss-i18n,
// dss-policy-jaxb(+crypto-json/xml), dss-diagnostic-jaxb only). Following the
// precedent set by dss/diagnostic/diagnostic_data_facade.go, this port
// collapses AbstractJaxbFacade's marshal/unmarshal template method directly
// into ValidationPolicyFacade - concretely, by delegating to
// policy/jaxb.Marshal/Unmarshal, which already implement the JAXB-canonical
// byte-for-byte output the marshal-parity contract requires (see that
// package's xml.go). XSD-schema validation (the `validate` boolean Java
// threads through every overload) is out of scope, per
// ValidationPolicyXmlDefiner's deferred Schema() - every unmarshal here is
// schema-unchecked, matching AbstractJaxbFacade#unmarshall(InputStream,
// false).
package policy

import (
	"errors"
	"io"
	"os"

	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/policy/jaxb"
)

// ValidationPolicyFacade is used to read an XML validation policy. Ports
// ValidationPolicyFacade.
type ValidationPolicyFacade struct{}

// NewValidationPolicyFacade initializes a new ValidationPolicyFacade. Ports
// ValidationPolicyFacade#newFacade.
func NewValidationPolicyFacade() *ValidationPolicyFacade {
	return &ValidationPolicyFacade{}
}

// Unmarshal parses r into a ConstraintsParameters document. Ports the
// inherited AbstractJaxbFacade#unmarshall(InputStream), collapsed as
// described in the file header.
func (f *ValidationPolicyFacade) Unmarshal(r io.Reader) (*jaxb.ConstraintsParameters, error) {
	if r == nil {
		return nil, errors.New("InputStream is null")
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return jaxb.Unmarshal(data)
}

// UnmarshalFile parses the document at path into a ConstraintsParameters
// document. Ports the inherited AbstractJaxbFacade#unmarshall(File).
func (f *ValidationPolicyFacade) UnmarshalFile(path string) (*jaxb.ConstraintsParameters, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jaxb.Unmarshal(data)
}

// Marshal writes cp as an XML validation-policy document. Ports the
// inherited AbstractJaxbFacade#marshall(T), which in Java wraps cp via
// ObjectFactory#createConstraintsParameters (see
// jaxb.NewConstraintsParametersElement).
func (f *ValidationPolicyFacade) Marshal(cp *jaxb.ConstraintsParameters) ([]byte, error) {
	return jaxb.Marshal(jaxb.NewConstraintsParametersElement(cp))
}

// GetValidationPolicyFromPath gets the validation policy at path. Ports
// ValidationPolicyFacade#getValidationPolicy(String), whose `path` was a
// classpath resource path (ValidationPolicyFacade.class.getResourceAsStream);
// Go has no classpath, so this reads path directly from the filesystem
// instead - the closest available equivalent, since no embed.FS carries an
// arbitrary caller-supplied set of named policies (only the default policy,
// embedded by EtsiValidationPolicyFactory, is bundled with this port).
func (f *ValidationPolicyFacade) GetValidationPolicyFromPath(path string) (modelpolicy.ValidationPolicy, error) {
	return f.GetValidationPolicyFromFile(path)
}

// GetValidationPolicy gets the validation policy from r. Ports
// ValidationPolicyFacade#getValidationPolicy(InputStream).
//
// Java's Objects.requireNonNull(is, "The provided validation policy is
// null") becomes a returned error per PORTING.md (this method is reached
// through EtsiValidationPolicyFactory#loadValidationPolicy(InputStream),
// which already returns an error on failure - see that file).
func (f *ValidationPolicyFacade) GetValidationPolicy(r io.Reader) (modelpolicy.ValidationPolicy, error) {
	if r == nil {
		return nil, errors.New("the provided validation policy is null")
	}
	cp, err := f.Unmarshal(r)
	if err != nil {
		return nil, err
	}
	return NewEtsiValidationPolicy(cp), nil
}

// GetValidationPolicyFromFile gets the validation policy from the document
// at path. Ports ValidationPolicyFacade#getValidationPolicy(File).
func (f *ValidationPolicyFacade) GetValidationPolicyFromFile(path string) (modelpolicy.ValidationPolicy, error) {
	if path == "" {
		return nil, errors.New("the provided validation policy is null")
	}
	cp, err := f.UnmarshalFile(path)
	if err != nil {
		return nil, err
	}
	return NewEtsiValidationPolicy(cp), nil
}
