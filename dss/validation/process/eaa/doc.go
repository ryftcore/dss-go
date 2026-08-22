// Package eaa ports eu.europa.esig.dss.validation.process.eaa (DSS 6.5.RC1):
// EAAValidationBlock, EAAValidationProcess, and the one eaa.checks class that
// belongs here rather than in the sibling eaa/checks package -
// KeyBindingSignatureValidationResultCheck, since it wires
// qualification.SignatureValidationResultCheck and this package already
// imports qualification for EAAValidationBlock/EAAValidationProcess.
//
// eu.europa.esig.dss.validation.process.eaa.checks and
// eu.europa.esig.dss.validation.process.eaa.status were originally flattened
// into this same package, following the checks-subpackage flattening
// convention used throughout this port. That has been undone:
// EAAValidationBlock pulls in package qualification, and qualification ->
// vpfswatsp -> bbb/sav, while bbb/sav's `-tags eaa` files
// (eaa_acceptance_validation.go, eaa_revocation_token_acceptance_validation.go)
// and bbb/fc's equivalents need the eaa.checks/eaa.status check constructors -
// so a single flattened eaa package closes the cycle
// bbb/sav -> eaa -> qualification -> vpfswatsp -> bbb/sav
// under `-tags eaa`. The checks now live in the sibling package
// github.com/ryftcore/dss-go/dss/validation/process/eaa/checks (Java's own eaa vs
// eaa.checks/eaa.status package boundary), which has no qualification
// dependency and so does not participate in the cycle - the same "prefer
// Java's own package boundary" escape hatch used for vpftspwatsp/checks and
// vpfswatsp/evidencerecord.
//
// This package is intentionally left untagged (no `//go:build eaa`): package
// qualification and vpfswatsp are untagged, and eaa itself has no dependency
// on anything gated by that tag, so it builds under both the default and the
// `-tags eaa` configurations - same as xcv/aov/fc/sav.
package eaa
