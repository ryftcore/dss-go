// Package validation ports dss-validation's core signature and certificate
// validation entry points: the interfaces and base types that turn a signed
// document (or a bare certificate) into a validation report by driving the
// EN 319 102-1 process implemented in validation/process and its
// sub-packages.
//
// # Main entry types
//
// SignedDocumentValidator is the interface consumers obtain from a document
// analyzer to validate a signed document; SignedDocumentValidatorBase
// implements the shared plumbing (certificate verifier wiring, token
// extraction, diagnostic-data building) that format-specific validators in
// cades, xades, pades, jades and asic embed. CertificateValidator validates
// a bare X.509 certificate against a certificate verifier. DocumentValidator
// and DocumentValidatorFactory are the format-agnostic dispatch interfaces
// used by validation/job and by consumers that do not know a document's
// format up front. TrustAnchorVerifierFactory and
// RevocationDataVerifierFactory build the verifier components a validator
// consults while checking trust chains and revocation freshness.
//
// Ported from dss-validation (DSS 6.5.RC1); every file names its own
// upstream source in its header.
package validation
