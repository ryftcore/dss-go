// Package eaa ports eu.europa.esig.dss.validation.process.eaa.checks and
// eu.europa.esig.dss.validation.process.eaa.status (DSS 6.5.RC1) into a single
// flattened Go package, following the same checks-subpackage flattening
// convention used throughout this port (e.g. xcv flattening nine Java
// packages into one Go package; see x509_certificate_validation.go's header).
//
// Integration note (phase 8d integration pass): this package resolves the
// forward dependency that fc/eaa_format_checking.go, fc/eaa_revocation_format_checking.go,
// sav/eaa_acceptance_validation.go and sav/eaa_revocation_token_acceptance_validation.go
// (all phase 8c, `//go:build eaa`) declared but left unported - each of those
// files' header documents the exact constructor signatures it assumed, and
// this package supplies them 1:1. Only the check classes those four callers
// actually reference are ported here; KeyBindingSignatureValidationResultCheck
// (eaa.checks) is NOT one of them - it wires eu.europa.esig.dss.validation.process.qualification.signature.checks.SignatureValidationResultCheck
// and eu.europa.esig.dss.detailedreport.jaxb.XmlValidationProcessEAA, an
// unrelated forward dependency (EAA qualification, not EAA format/acceptance
// validation) that no phase 8c caller in this tree pulls in; left unported,
// consistent with "port minimally from upstream."
//
// This package is intentionally left untagged (no `//go:build eaa`): it has
// no dependency on anything gated by that tag, so it builds under both the
// default and the `-tags eaa` configurations, the same way xcv/aov/fc/sav
// themselves are untagged despite existing partly to serve eaa-tagged callers.
package eaa
