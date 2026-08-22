// Package cades ports dss-cades (eu.europa.esig.dss.cades), CAdES
// (CMS Advanced Electronic Signatures, ETSI EN 319 122) signing, extension,
// and validation: building B/T/LT/LTA-level signatures over a CMS SignedData
// structure, extending an existing signature to a higher level (adding a
// signature timestamp, revocation data, or an archive timestamp), and
// parsing a CAdES-signed document for the validation engine.
//
// The main entry types are Service (the signature/extension service
// implementing document.SignatureService), SignatureParameters
// and TimestampParameters (signing configuration), Signature (a
// parsed signature, implementing the validation engine's AdvancedSignature),
// and CMSDocumentValidator/CMSDocumentAnalyzer (the validator entry point
// for CAdES and bare CMS documents).
package cades
