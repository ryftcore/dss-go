// Package xades ports dss-xades (eu.europa.esig.dss.xades), XAdES
// (XML Advanced Electronic Signatures, ETSI EN 319 132) signing, extension,
// and validation: building B/T/LT/LTA/legacy-level enveloped, enveloping,
// detached, or internally-detached XML signatures, extending an existing
// signature to a higher level, and parsing an XML-signed document for the
// validation engine.
//
// The main entry types are Service (the signature/extension service
// implementing document.SignatureService), SignatureParameters
// and TimestampParameters (signing configuration), Signature (a
// parsed signature, implementing the validation engine's AdvancedSignature),
// and XMLDocumentValidator/XMLDocumentAnalyzer (the validator entry point
// for XML-signed documents). The DSSReference/DSSTransform builders and the
// SignatureBuilder family assemble the underlying XML Signature structure;
// most callers only need Service and the parameter types.
package xades
