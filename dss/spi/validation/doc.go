// Package validation ports the dss-spi validation subpackage
// (eu.europa.esig.dss.spi.validation), the shared building blocks every
// per-format AdvancedSignature implementation (Signature,
// Signature, ...) and evidence-record validator is built from:
// signature identifier derivation, embedded evidence-record helpers, and
// the AdvancedSignature interface itself.
//
// The main entry types are AdvancedSignature (the format-independent parsed
// signature interface the validation engine drives),
// AbstractSignatureIdentifierBuilder, and
// AbstractEmbeddedEvidenceRecordHelper. Its analyzer, executor, identifier,
// scope, timestamp, and tls subpackages hold the corresponding
// per-concern helpers used by the validation engine (see package
// validation for the top-level entry points).
package validation
