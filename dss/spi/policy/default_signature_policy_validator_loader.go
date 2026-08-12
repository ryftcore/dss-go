// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/DefaultSignaturePolicyValidatorLoader.java (DSS 6.5.RC1).
//
// Deviation: Java discovers implementations via
// ServiceLoader.load(SignaturePolicyValidator.class), driven by
// META-INF/services/eu.europa.esig.dss.spi.policy.SignaturePolicyValidator,
// which upstream lists (in order) as ZeroHashSignaturePolicyValidator,
// NonASN1SignaturePolicyValidator, BasicASN1SignaturePolicyValidator,
// EmptySignaturePolicyValidator. Go has no service-provider mechanism, so
// defaultSignaturePolicyValidatorLoaderDefaultValidators reproduces that
// fixed, ordered registry directly.
package policy

import "github.com/utain/esig/dss/model/signature"

// defaultSignaturePolicyValidatorLoaderDefaultValidators are constructed
// fresh per LoadValidator call (matching ServiceLoader semantics, which
// yields new provider instances per iteration by default) in the order
// upstream's META-INF/services file lists them.
func defaultSignaturePolicyValidatorLoaderDefaultValidators() []SignaturePolicyValidator {
	return []SignaturePolicyValidator{
		NewZeroHashSignaturePolicyValidator(),
		NewNonASN1SignaturePolicyValidator(),
		NewBasicASN1SignaturePolicyValidator(),
		NewEmptySignaturePolicyValidator(),
	}
}

// DefaultSignaturePolicyValidatorLoader loads a relevant
// SignaturePolicyValidator based on the policy content.
type DefaultSignaturePolicyValidatorLoader struct {
	// defaultSignaturePolicyValidator is the validator to be used when only
	// a basic validation according to the signature format is required.
	// NOTE: can be nil (the best corresponding validator will be loaded).
	defaultSignaturePolicyValidator SignaturePolicyValidator

	// supportHashAsInTechnicalSpecification defines whether the
	// SignaturePolicy.hashAsInTechnicalSpecification attribute is
	// supported. Default: true.
	supportHashAsInTechnicalSpecification bool
}

// NewDefaultSignaturePolicyValidatorLoader is the default constructor,
// instantiating the object with a nil SignaturePolicyValidator.
func NewDefaultSignaturePolicyValidatorLoader() *DefaultSignaturePolicyValidatorLoader {
	return &DefaultSignaturePolicyValidatorLoader{supportHashAsInTechnicalSpecification: true}
}

// DefaultSignaturePolicyValidatorLoaderDefaultOnly creates a
// DefaultSignaturePolicyValidatorLoader running the signature policy
// validation using defaultSignaturePolicyValidator. The default
// implementation will be used on all signature policy hash calculations.
func DefaultSignaturePolicyValidatorLoaderDefaultOnly(defaultSignaturePolicyValidator SignaturePolicyValidator) *DefaultSignaturePolicyValidatorLoader {
	loader := NewDefaultSignaturePolicyValidatorLoader()
	loader.SetDefaultSignaturePolicyValidator(defaultSignaturePolicyValidator)
	loader.SetSupportHashAsInTechnicalSpecification(false)
	return loader
}

// DefaultSignaturePolicyValidatorLoaderDefaultUnlessSpecified creates a
// DefaultSignaturePolicyValidatorLoader running the signature policy
// validation using defaultSignaturePolicyValidator. The default
// implementation will be used on all signature policy hash calculations,
// unless a "HashAsInTechnicalSpecification" parameter is set within the
// Signature Policy Identifier.
func DefaultSignaturePolicyValidatorLoaderDefaultUnlessSpecified(defaultSignaturePolicyValidator SignaturePolicyValidator) *DefaultSignaturePolicyValidatorLoader {
	loader := NewDefaultSignaturePolicyValidatorLoader()
	loader.SetDefaultSignaturePolicyValidator(defaultSignaturePolicyValidator)
	loader.SetSupportHashAsInTechnicalSpecification(true)
	return loader
}

// DefaultSignaturePolicyValidatorLoaderPolicyBased creates a
// DefaultSignaturePolicyValidatorLoader running the signature policy
// validation loading the SignaturePolicyValidator based on the signature
// policy's specification. The first SignaturePolicyValidator matching the
// signature policy will be selected. If not defined explicitly, one of the
// default signature policies will be used.
func DefaultSignaturePolicyValidatorLoaderPolicyBased() *DefaultSignaturePolicyValidatorLoader {
	return NewDefaultSignaturePolicyValidatorLoader()
}

// SetDefaultSignaturePolicyValidator sets a SignaturePolicyValidator to be
// used for default signature policy processing according to the signature
// format (when SignaturePolicy.hashAsInTechnicalSpecification == false).
func (l *DefaultSignaturePolicyValidatorLoader) SetDefaultSignaturePolicyValidator(defaultSignaturePolicyValidator SignaturePolicyValidator) {
	l.defaultSignaturePolicyValidator = defaultSignaturePolicyValidator
}

// SetSupportHashAsInTechnicalSpecification sets whether the
// SignaturePolicy.hashAsInTechnicalSpecification attribute is supported. If
// set to true, the behavior of the loader will change based on the
// attribute presence. Otherwise, it is ignored.
//
// Default: true (SignaturePolicy.hashAsInTechnicalSpecification attribute
// is supported).
func (l *DefaultSignaturePolicyValidatorLoader) SetSupportHashAsInTechnicalSpecification(supportHashAsInTechnicalSpecification bool) {
	l.supportHashAsInTechnicalSpecification = supportHashAsInTechnicalSpecification
}

// LoadValidator returns the relevant validator for a SignaturePolicy. See
// the file-level deviation note regarding the ServiceLoader replacement.
func (l *DefaultSignaturePolicyValidatorLoader) LoadValidator(signaturePolicy *signature.SignaturePolicy) SignaturePolicyValidator {
	if l.defaultSignaturePolicyValidator != nil &&
		(!l.supportHashAsInTechnicalSpecification || !signaturePolicy.IsHashAsInTechnicalSpecification()) {
		return l.defaultSignaturePolicyValidator
	}

	for _, candidate := range defaultSignaturePolicyValidatorLoaderDefaultValidators() {
		if candidate.CanValidate(signaturePolicy) {
			return candidate
		}
	}
	// if not empty and no other implementation is found
	return NewNonASN1SignaturePolicyValidator()
}

var _ SignaturePolicyValidatorLoader = (*DefaultSignaturePolicyValidatorLoader)(nil)
