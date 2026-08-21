// Package spi ports dss-spi (eu.europa.esig.dss.spi), the crypto and
// PKI service layer shared by every signature format and by the validation
// engine: certificate and revocation (CRL/OCSP) sources, X.509 chain
// building, digest calculation, and the CMS-derived certificate/revocation
// extraction used by CAdES-family formats.
//
// The main entry types are CertificateSource (and its many implementations,
// e.g. CommonCertificateSource, KeyStoreCertificateSource,
// TrustedCertificateSource), RevocationSource[R]/CRLSource/OCSPSource and
// their online/offline/repository-backed implementations, CertificateToken
// chain-building helpers (CandidatesForSigningCertificate,
// BaselineBCertificateSelector), and RevocationToken[R]/CRLToken/OCSPToken.
//
// Subpackages group SPI types by concern: spi/x509/aia and
// spi/x509/evidencerecord (AIA certificate retrieval and evidence-record
// digest support), spi/validation (the shared signature/evidence-record
// validation building blocks per-format validators embed, plus its
// analyzer/executor/scope/timestamp/identifier/tls sub-areas),
// spi/policy (signature policy resolution), spi/tsl and spi/lote
// (trusted-list and List of Trusted Entities certificate sources),
// spi/client/http and spi/client/jdbc (HTTP and JDBC-backed data loaders),
// spi/exception, spi/alerts, spi/eaa (electronic attestation of attributes
// support), and spi/random (deterministic-testing secure random helpers).
package spi
