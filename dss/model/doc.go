// Package model ports dss-model (eu.europa.esig.dss.model), the core value
// objects and interfaces every other package in this module builds on: the
// DSSDocument abstraction for in-memory/on-disk/streamed content, digest and
// signature value wrappers, certificate/revocation token types, and the
// signature/timestamp parameter hierarchies signing services consume.
//
// The main entry types are DSSDocument and its implementations
// (InMemoryDocument, FileDocument, DigestDocument), CertificateToken and the
// Token/TokenBase interfaces, Digest and DSSMessageDigest, and the
// AbstractSerializableSignatureParameters/SerializableTimestampParameters
// family that per-format packages (cades, xades, ...) extend.
//
// Subpackages group model types by concern: model/policy (cryptographic
// suite and certificate applicability rules), model/scope (signature
// scopes), model/signature (signature-level value objects such as
// SignaturePolicy), model/timedependent (time-varying value containers),
// model/tsl and model/lote (trusted-list and List of Trusted Entities
// identifiers), model/job (validation job info records), model/http
// (HTTP response envelopes), and model/tls (TLS certificate chains).
package model
