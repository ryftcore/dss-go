// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/policy/ValidationPolicyWithCryptographicSuite.java (DSS 6.5.RC1).
//
// This class wraps a provided modelpolicy.ValidationPolicy and adds per-scope
// modelpolicy.CryptographicSuite overrides. Java implements ValidationPolicy by hand-writing
// ~190 one-line delegating methods to the wrapped instance, overriding only four
// (SignatureCryptographicConstraint, CertificateCryptographicConstraint,
// EvidenceRecordCryptographicConstraint, EAACryptographicConstraint) plus toString(). This port
// embeds the wrapped modelpolicy.ValidationPolicy interface directly: Go's method promotion
// reproduces every one of those ~190 pure delegations automatically (identical behavior, zero
// transcription risk across an interface this large), and the four methods that need real logic,
// plus String() (toString), are defined explicitly below to shadow the promoted ones. This
// mirrors the DSSDocument-embedding decorator precedent already used in this codebase (e.g.
// asic.DSSZipEntryDocument, jades.HTTPHeadersPayloadBuilder) and PORTING.md's general license to
// collapse Java boilerplate where Go has a native, behavior-preserving equivalent.
//
// Java's private Map<ContextAndSubContext, CryptographicSuite> cryptographicSuitesMap
// (HashMap, so its own iteration order is unspecified and irrelevant - only exact-key lookup is
// used) becomes a Go map keyed directly by the comparable contextAndSubContext struct.
package policy

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
)

// ValidationPolicyWithCryptographicSuite wraps a provided modelpolicy.ValidationPolicy and
// modelpolicy.CryptographicSuite(s), whether applicable. The class manages the returned data.
type ValidationPolicyWithCryptographicSuite struct {
	// validationPolicy is the wrapped validation policy. Embedded so every ValidationPolicy
	// method Java delegates verbatim is promoted for free; see the file header.
	modelpolicy.ValidationPolicy

	// cryptographicSuitesMap is a map of applicability scopes to their cryptographic suite.
	cryptographicSuitesMap map[contextAndSubContext]modelpolicy.CryptographicSuite
}

// compile-time interface assertion.
var _ modelpolicy.ValidationPolicy = (*ValidationPolicyWithCryptographicSuite)(nil)

// NewValidationPolicyWithCryptographicSuite is the default constructor to create a validation
// policy with cryptographic suites wrapper.
//
// Panics when validationPolicy is nil (Java's Objects.requireNonNull("ValidationPolicy cannot
// be null!")).
func NewValidationPolicyWithCryptographicSuite(validationPolicy modelpolicy.ValidationPolicy) *ValidationPolicyWithCryptographicSuite {
	if validationPolicy == nil {
		panic("ValidationPolicy cannot be null!")
	}
	return &ValidationPolicyWithCryptographicSuite{
		ValidationPolicy:       validationPolicy,
		cryptographicSuitesMap: make(map[contextAndSubContext]modelpolicy.CryptographicSuite),
	}
}

// SetCryptographicSuite sets a global cryptographic suite, to be applied when no
// context-specific cryptographic rules are defined. The provided cryptographic suite will
// overwrite the current settings for the global cryptographic suite.
func (p *ValidationPolicyWithCryptographicSuite) SetCryptographicSuite(cryptographicSuite modelpolicy.CryptographicSuite) {
	p.cryptographicSuitesMap[contextAndSubContext{}] = cryptographicSuite
}

// SetCryptographicSuiteForContext sets a cryptographic suite for a specific Context. The
// provided cryptographic suite will overwrite the current settings only for the defined
// applicability scope.
func (p *ValidationPolicyWithCryptographicSuite) SetCryptographicSuiteForContext(cryptographicSuite modelpolicy.CryptographicSuite, context enumerations.Context) {
	p.cryptographicSuitesMap[contextAndSubContext{context: context}] = cryptographicSuite
}

// SetCryptographicSuiteForContextAndSubContext sets a cryptographic suite for a specific
// combination of a Context and a SubContext. The provided cryptographic suite will overwrite the
// current settings only for the defined applicability scope.
//
// Returns an error when context is ContextEvidenceRecord and subContext is not empty (Java's
// thrown IllegalArgumentException("Please use a NULL SubContext for the Context.EVIDENCE_RECORD
// or use #setCryptographicSuite(cryptographicSuite, context) method.")).
func (p *ValidationPolicyWithCryptographicSuite) SetCryptographicSuiteForContextAndSubContext(cryptographicSuite modelpolicy.CryptographicSuite,
	context enumerations.Context, subContext enumerations.SubContext) error {
	if enumerations.ContextEvidenceRecord == context && subContext != "" {
		return fmt.Errorf("please use a NULL SubContext for the Context.EVIDENCE_RECORD or " +
			"use #setCryptographicSuite(cryptographicSuite, context) method")
	}
	p.cryptographicSuitesMap[contextAndSubContext{context: context, subContext: subContext}] = cryptographicSuite
	return nil
}

// SignatureCryptographicConstraint creates the CryptographicSuite corresponding to the context
// parameter.
func (p *ValidationPolicyWithCryptographicSuite) SignatureCryptographicConstraint(context enumerations.Context) modelpolicy.CryptographicSuite {
	if cryptographicSuite := p.cryptographicSuite(context, ""); cryptographicSuite != nil {
		return cryptographicSuite
	}
	return p.ValidationPolicy.SignatureCryptographicConstraint(context)
}

// CertificateCryptographicConstraint creates the CryptographicSuite corresponding to the
// context/subContext parameters.
func (p *ValidationPolicyWithCryptographicSuite) CertificateCryptographicConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.CryptographicSuite {
	if cryptographicSuite := p.cryptographicSuite(context, subContext); cryptographicSuite != nil {
		return cryptographicSuite
	}
	return p.ValidationPolicy.CertificateCryptographicConstraint(context, subContext)
}

// EvidenceRecordCryptographicConstraint returns cryptographic constraints for validation of an
// Evidence Record.
func (p *ValidationPolicyWithCryptographicSuite) EvidenceRecordCryptographicConstraint() modelpolicy.CryptographicSuite {
	if cryptographicSuite := p.cryptographicSuite(enumerations.ContextEvidenceRecord, ""); cryptographicSuite != nil {
		return cryptographicSuite
	}
	return p.ValidationPolicy.EvidenceRecordCryptographicConstraint()
}

// EAACryptographicConstraint returns cryptographic constraints for validation of an EAA.
func (p *ValidationPolicyWithCryptographicSuite) EAACryptographicConstraint() modelpolicy.CryptographicSuite {
	if cryptographicSuite := p.cryptographicSuite(enumerations.ContextEAA, ""); cryptographicSuite != nil {
		return cryptographicSuite
	}
	return p.ValidationPolicy.EAACryptographicConstraint()
}

// cryptographicSuite ports the private getCryptographicSuite(Context, SubContext) helper:
// checks, in order, for an exact Context+SubContext match, then a Context-only match (global
// SubContext), then a fully-global entry, returning nil when none is found.
func (p *ValidationPolicyWithCryptographicSuite) cryptographicSuite(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.CryptographicSuite {
	if len(p.cryptographicSuitesMap) == 0 {
		return nil
	}
	if cs, ok := p.cryptographicSuitesMap[contextAndSubContext{context: context, subContext: subContext}]; ok {
		return cs
	}
	if cs, ok := p.cryptographicSuitesMap[contextAndSubContext{context: context}]; ok {
		return cs
	}
	if cs, ok := p.cryptographicSuitesMap[contextAndSubContext{}]; ok {
		return cs
	}
	return nil
}

// String ports the toString() override.
//
// Java's cryptographicSuitesMap.values() iterates the backing HashMap in its own
// implementation-defined order; this port sorts the names instead so the result is deterministic
// (no map-iteration-ordered output), rather than reproducing an order Java itself does not
// guarantee.
func (p *ValidationPolicyWithCryptographicSuite) String() string {
	names := make([]string, 0, len(p.cryptographicSuitesMap))
	for _, cs := range p.cryptographicSuitesMap {
		names = append(names, cs.PolicyName())
	}
	sort.Strings(names)
	return fmt.Sprintf("ValidationPolicyWithCryptographicSuite [validationPolicyName=%s, cryptographicSuitesNames=[%s]]",
		p.ValidationPolicy.PolicyName(), strings.Join(names, ", "))
}
