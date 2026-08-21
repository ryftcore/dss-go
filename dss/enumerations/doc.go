// Package enumerations ports dss-enumerations
// (eu.europa.esig.dss.enumerations), the shared vocabulary of typed
// constants that every other package in this module builds on: signature
// forms and levels, validation Indications and SubIndications, certificate
// and revocation classifications, digest and encryption algorithm
// identifiers, ASN.1 object identifiers, and the many other closed value
// sets ETSI EN 319 102-1 and the related AdES/eIDAS standards define.
//
// Each Java enum becomes a Go named string (or int) type with one constant
// per enum value, following the JavaEnumName_VALUE_NAME spelling so the
// mapping back to the upstream type stays mechanical; most also get a
// Values() accessor returning every constant in declaration order, mirroring
// Java's Enum.values(). Interfaces such as CertificateApprovalStatus or
// QCIdentMethod stand in for enums upstream implements as an "enum with
// members" (each constant carrying extra fields/behaviour) where a plain Go
// const cannot express that.
//
// The main entry types are the signature-related enums (SignatureLevel,
// SignatureForm, SignatureProfile), the validation-result enums (Indication,
// SubIndication), and the certificate/algorithm enums (CertificateOrigin,
// DigestAlgorithm, EncryptionAlgorithm, SignatureAlgorithm).
package enumerations
