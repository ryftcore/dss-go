//go:build !eaa

// eaaBuildTagEnabled tells the full-corpus report byte-parity test whether the
// `eaa` build tag - which gates dss/validation/process's EAA validation blocks
// in - is set for this build.
package executor

const eaaBuildTagEnabled = false
