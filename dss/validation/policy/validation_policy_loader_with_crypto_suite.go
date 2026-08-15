// Ported from the nested static class
// dss-validation/src/main/java/eu/europa/esig/dss/validation/policy/ValidationPolicyLoader.java
// $ValidationPolicyLoaderWithCryptoSuite (DSS 6.5.RC1).
//
// Java's nested class "extends ValidationPolicyLoader" purely to keep the builder's fluent return
// type, and overrides every inherited with*()/create() method to delegate to a wrapped
// `validationPolicyLoader` field rather than to itself. This port embeds *ValidationPolicyLoader
// directly: Go's method promotion reproduces that full delegation set for free (every With*
// method promoted from the embedded pointer already targets the wrapped loader, and returns a
// fresh *ValidationPolicyLoaderWithCryptoSuite the same way Java's overrides do), so only the
// five level-setter methods this class actually adds need a body - see
// validation_policy_with_cryptographic_suite.go's file header for the identical rationale applied
// there.
package policy

import (
	"github.com/utain/esig/dss/enumerations"
	modelpolicy "github.com/utain/esig/dss/model/policy"
)

// ValidationPolicyLoaderWithCryptoSuite provides a user-friendly configuration of the execution
// levels for the last set cryptographic suite(s). For generic With* methods inherited from
// ValidationPolicyLoader (promoted via embedding), the execution is propagated to the original
// wrapped loader.
type ValidationPolicyLoaderWithCryptoSuite struct {
	// ValidationPolicyLoader is the wrapped loader; embedded to promote its With*/Create methods.
	*ValidationPolicyLoader

	// cryptographicSuites is the cryptographic suite(s) most recently added, whose execution
	// levels the And* methods below configure.
	cryptographicSuites []modelpolicy.CryptographicSuite
}

// newValidationPolicyLoaderWithCryptoSuite ports both nested-class constructors (the single-
// CryptographicSuite overload wraps its argument in a one-element slice in Java; callers here
// pass the slice directly).
func newValidationPolicyLoaderWithCryptoSuite(validationPolicyLoader *ValidationPolicyLoader, cryptographicSuites []modelpolicy.CryptographicSuite) *ValidationPolicyLoaderWithCryptoSuite {
	return &ValidationPolicyLoaderWithCryptoSuite{
		ValidationPolicyLoader: validationPolicyLoader,
		cryptographicSuites:    cryptographicSuites,
	}
}

// AndLevel sets the global execution level for the last provided cryptographic suite(s).
func (l *ValidationPolicyLoaderWithCryptoSuite) AndLevel(level enumerations.Level) *ValidationPolicyLoaderWithCryptoSuite {
	for _, s := range l.cryptographicSuites {
		s.SetLevel(level)
	}
	return l
}

// AndAcceptableDigestAlgorithmsLevel sets the execution level for the acceptable digest
// algorithms check of the last provided cryptographic suite(s).
func (l *ValidationPolicyLoaderWithCryptoSuite) AndAcceptableDigestAlgorithmsLevel(level enumerations.Level) *ValidationPolicyLoaderWithCryptoSuite {
	for _, s := range l.cryptographicSuites {
		s.SetAcceptableDigestAlgorithmsLevel(level)
	}
	return l
}

// AndAcceptableSignatureAlgorithmsLevel sets the execution level for the acceptable signature
// algorithms check of the last provided cryptographic suite(s).
func (l *ValidationPolicyLoaderWithCryptoSuite) AndAcceptableSignatureAlgorithmsLevel(level enumerations.Level) *ValidationPolicyLoaderWithCryptoSuite {
	for _, s := range l.cryptographicSuites {
		s.SetAcceptableSignatureAlgorithmsLevel(level)
	}
	return l
}

// AndAcceptableSignatureAlgorithmsMiniKeySizeLevel sets the execution level for acceptable
// minimum key sizes of encryption algorithms check of the last provided cryptographic suite(s).
func (l *ValidationPolicyLoaderWithCryptoSuite) AndAcceptableSignatureAlgorithmsMiniKeySizeLevel(level enumerations.Level) *ValidationPolicyLoaderWithCryptoSuite {
	for _, s := range l.cryptographicSuites {
		s.SetAcceptableSignatureAlgorithmsMiniKeySizeLevel(level)
	}
	return l
}

// AndAlgorithmsExpirationDateLevel sets the execution level for the expiration of the
// cryptographic algorithms check of the last provided cryptographic suite(s).
func (l *ValidationPolicyLoaderWithCryptoSuite) AndAlgorithmsExpirationDateLevel(level enumerations.Level) *ValidationPolicyLoaderWithCryptoSuite {
	for _, s := range l.cryptographicSuites {
		s.SetAlgorithmsExpirationDateLevel(level)
	}
	return l
}

// AndAlgorithmsExpirationTimeAfterPolicyUpdateLevel sets the execution level for the expiration
// after the cryptographic suite update date of the cryptographic algorithms check of the last
// provided cryptographic suite(s).
func (l *ValidationPolicyLoaderWithCryptoSuite) AndAlgorithmsExpirationTimeAfterPolicyUpdateLevel(level enumerations.Level) *ValidationPolicyLoaderWithCryptoSuite {
	for _, s := range l.cryptographicSuites {
		s.SetAlgorithmsExpirationTimeAfterPolicyUpdateLevel(level)
	}
	return l
}
