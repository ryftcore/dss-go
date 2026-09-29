// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/SantuarioInitializer.java (DSS 6.5.RC1).
//
// This is a documented no-op stub. Upstream one-time-initializes Apache Santuario (xmlsec), a
// third-party XML-DSig/XML-Security library, registering its global (process-wide) transform,
// signature-algorithm, canonicalizer and key-resolver tables before first use. The Go port
// depends on no such third-party library and keeps no equivalent process-global registry:
// internal/xmlc14n's seven canonicalizers are plain functions requiring no registration step -
// canonicalization allocates all state per call, with no package-level mutable state - and the
// transform/signature-algorithm registries the future
// xmldsig layer needs are constructed per call site, not process-global. There is therefore
// nothing for SantuarioInitializerInit to do; it and SantuarioInitializerIsInitialized are
// kept only so call sites ported unchanged from Java compile and behave inertly.
package utils

import "sync/atomic"

// santuarioInitializerInitialized mirrors SantuarioInitializer.alreadyInitialized. Since
// Init does no real work, this only tracks whether Init has been called at least once, for
// IsInitialized's benefit. It is atomic because the xades service constructors call
// SantuarioInitializerInit on every construction, possibly from several goroutines.
var santuarioInitializerInitialized atomic.Bool

// SantuarioInitializerIsInitialized reports whether SantuarioInitializerInit has been called.
// Ports isInitialized(); the Init.isInitialized() branch (Santuario's own global flag) has no
// Go counterpart and is dropped, since there is no Santuario library to already be
// initialized by another caller.
func SantuarioInitializerIsInitialized() bool {
	return santuarioInitializerInitialized.Load()
}

// SantuarioInitializerInit is a no-op beyond recording that it ran: see the file-level doc
// comment for why there is nothing to initialize. Ports init()/dynamicInit().
func SantuarioInitializerInit() {
	santuarioInitializerInitialized.Store(true)
}
