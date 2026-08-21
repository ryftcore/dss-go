// Package crlparser ports dss-crl-parser and dss-crl-parser-x509crl
// (eu.europa.esig.dss.crl), CRL parsing and revocation-status lookup: a
// streaming parser for X.509 Certificate Revocation Lists that avoids
// loading a large CRL fully into memory just to answer "is this certificate
// revoked", plus the CRLValidity result of validating a CRL's signature and
// issuer.
//
// The main entry types are ICRLUtils (the parsing/validation entry point,
// implemented by CRLUtils), CRLBinary (a CRL's raw bytes plus digest), and
// CRLValidity/CRLEntry (the parsed validity result and individual revoked
// entries).
package crlparser
