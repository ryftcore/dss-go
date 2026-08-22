// Package checks ports eu.europa.esig.dss.validation.process.eaa.checks and
// eu.europa.esig.dss.validation.process.eaa.status (DSS 6.5.RC1) into a
// single flattened Go package, following the checks-subpackage flattening
// convention used throughout this port (e.g. xcv flattening nine Java
// packages into one Go package; see bbb/xcv/x509_certificate_validation.go's
// header). One exception: KeyBindingSignatureValidationResultCheck
// (eu.europa.esig.dss.validation.process.eaa.checks) is NOT here - it wires
// qualification.SignatureValidationResultCheck, so it lives in the parent
// eaa package instead, which already imports qualification; see this
// package's header for why that split matters.
//
// Split out from the eaa root package: this package holds every
// eaa.checks/eaa.status leaf check, none of which import
// package qualification or vpfswatsp. That keeps it safe for
// bbb/fc/eaa_format_checking.go, bbb/fc/eaa_revocation_format_checking.go,
// bbb/sav/eaa_acceptance_validation.go and
// bbb/sav/eaa_revocation_token_acceptance_validation.go (all `//go:build
// eaa`) to import: those files sit under bbb, which package eaa's own
// dependency on qualification -> vpfswatsp reaches back down to
// (ValidationProcessForSignaturesWithArchivalData -> bbb/sav). Importing the
// combined eaa package from bbb/sav or bbb/fc would close that cycle; this
// package, with no such dependency, does not.
//
// This package is intentionally left untagged (no `//go:build eaa`): it has
// no dependency on anything gated by that tag, so it builds under both the
// default and the `-tags eaa` configurations, the same way xcv/aov/fc/sav
// themselves are untagged despite existing partly to serve eaa-tagged callers.
package checks
