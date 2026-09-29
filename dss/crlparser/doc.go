// Package crlparser ports dss-crl-parser and dss-crl-parser-x509crl
// (eu.europa.esig.dss.crl), CRL parsing and revocation-status lookup: it parses
// X.509 Certificate Revocation Lists (crypto/x509.ParseRevocationList, plus a
// repair for RFC 5280 v1 CRLs that carry no version field) and answers "is this
// certificate revoked" for a serial number, plus the CRLValidity result of
// validating a CRL's signature and issuer.
//
// The CRL is decoded in full and its revoked-certificate entries are kept in
// memory for the lifetime of the CRLValidity: unlike upstream's optional
// streaming implementation (dss-crl-parser-stream), which this port does not
// provide, nothing here avoids loading a large CRL into memory.
//
// The main entry types are ICRLUtils (the parsing/validation entry point,
// implemented by CRLUtils), CRLBinary (a CRL's raw bytes plus digest), and
// CRLValidity/CRLEntry (the parsed validity result and individual revoked
// entries).
package crlparser
