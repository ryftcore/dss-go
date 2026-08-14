// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/DSSSecurityProvider.java (DSS 6.5.RC1).
//
// DEVIATION: this is a documented stub. Upstream bootstraps the JCA by installing
// BouncyCastle (and any configured alternatives) into java.security.Security, so that
// MessageDigest/Signature/KeyFactory lookups by algorithm name resolve. Go has no
// provider registry: crypto/x509, crypto/* and golang.org/x/crypto are linked directly by
// the code that needs them, so there is nothing to register and nothing to configure.
//
// The Java class' mutable configuration surface (getSecurityProvider,
// setSecurityProvider, get/setAlternativeSecurityProviders, getSecurityProviderName,
// getAlternativeSecurityProviderNames) is therefore NOT ported: every one of those
// methods exists only to hand a java.security.Provider to a JCA factory call. Only the
// bootstrap entry point survives, so that the ported call sites - DSSASN1Utils' static
// initializer, and any application code that copies upstream's documented startup step -
// keep compiling and keep reading as they do upstream.
package spi

// DSSSecurityProviderInitSystemProviders is the no-op Go counterpart of
// DSSSecurityProvider#initSystemProviders(). Upstream this adds the primary and
// alternative security providers to the JVM-wide java.security.Security list and must be
// called before any provider-backed operation; the Go port needs no such step, so calling
// it is harmless and calling it is never required.
func DSSSecurityProviderInitSystemProviders() {
	// empty: Go resolves crypto implementations at link time, not through a provider list.
}
