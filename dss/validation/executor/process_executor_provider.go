// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/ProcessExecutorProvider.java
// (DSS 6.5.RC1).

package executor

// ProcessExecutorProvider provides the executor for a validation process.
// Port of the ProcessExecutorProvider<PE extends ProcessExecutor<?>>
// interface. Java's wildcard bound PE extends ProcessExecutor<?> has no Go
// equivalent - Go type parameters cannot be constrained by a generic
// interface with an unbound argument - so PE is left unconstrained here
// and the concrete providers instantiate it with their own executor type.
type ProcessExecutorProvider[PE any] interface {

	// SetProcessExecutor provides the possibility to set the specific
	// process executor. Port of setProcessExecutor(PE).
	SetProcessExecutor(processExecutor PE)

	// DefaultProcessExecutor returns a default validator process executor.
	// Port of getDefaultProcessExecutor().
	DefaultProcessExecutor() PE
}
