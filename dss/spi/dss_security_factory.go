// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/security/DSSSecurityFactory.java (DSS 6.5.RC1).
//
// DSSSecurityFactory<I,O> is upstream's abstract base for the DSSCertificateTokenSecurityFactory,
// DSSP7CCertificatesSecurityFactory, DSSPublicKeySecurityFactory, DSSContentVerifierProviderSecurityFactory
// and DSSSignerInformationVerifierSecurityFactory families: it drives buildWithProvider(input,
// securityProvider) first against DSSSecurityProvider's primary java.security.Provider and, on
// failure, against each configured alternative, returning the first success and throwing a
// DSSException only once every provider has failed.
//
// DEVIATION: as documented on DSSSecurityProviderInitSystemProviders (dss_security_provider.go),
// Go has no java.security.Provider registry - crypto/x509, crypto/* and golang.org/x/crypto are
// linked directly by the code that needs them, so there is exactly one "security provider" and
// no alternative-provider fallback list to iterate. buildWithPrimarySecurityProvider and
// buildWithAlternativeSecurityProviders therefore collapse into a single BuildWithProvider call;
// Java's swallow-and-log-then-retry behaviour on a failed provider becomes a returned error, per
// PORTING.md's "data-dependent throw -> (T, error)" rule.
//
// Go's type-parameterised methods cannot themselves be generic, so this is a struct of function
// fields rather than an interface with an overridable buildWithProvider - Go's usual translation
// of Java's template-method inheritance.
package spi

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
)

// DSSSecurityFactory is the generic base used to build an O from an I. FactoryClassName names
// the produced type for the error message (getFactoryClassName()); ToString renders the input
// for diagnostics (toString(I), used by upstream's debug logging - kept for parity though this
// port has no log statements to feed it, per PORTING.md's "slf4j -> dropped unless load-bearing");
// BuildWithProvider performs the actual construction (buildWithProvider(input, securityProvider),
// with the securityProvider parameter dropped per the file DEVIATION above).
type DSSSecurityFactory[I any, O any] struct {
	// FactoryClassName is the produced type's name, as used in the "Unable to load %s..." error.
	FactoryClassName string

	// ToString renders the input for diagnostics. May be nil if unused.
	ToString func(input I) string

	// BuildWithProvider performs the actual construction.
	BuildWithProvider func(input I) (O, error)
}

// Build builds an O from input. Port of DSSSecurityFactory#build(I): upstream tries the primary
// provider then the alternatives before throwing; this port has exactly one provider (see the
// file DEVIATION), so a single BuildWithProvider call plays both roles.
func (f *DSSSecurityFactory[I, O]) Build(input I) (O, error) {
	output, err := f.BuildWithProvider(input)
	if err == nil {
		return output, nil
	}
	var zero O
	return zero, model.NewDSSErrorMessageCause(fmt.Sprintf(
		"Unable to load %s for the given input. All security providers have failed. More detail in debug mode.",
		f.FactoryClassName), err)
}
