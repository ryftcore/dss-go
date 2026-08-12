package spi

import "testing"

// TestDSSSecurityProviderInitSystemProviders documents the stub's contract: calling it is
// safe, repeatable and has no observable effect, because Go resolves crypto implementations
// at link time instead of through a java.security.Security provider list.
func TestDSSSecurityProviderInitSystemProviders(t *testing.T) {
	DSSSecurityProviderInitSystemProviders()
	DSSSecurityProviderInitSystemProviders()
}
