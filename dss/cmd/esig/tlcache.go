package main

import (
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/utain/esig/dss"
)

// tlCacheCertsDir is the subdirectory of a -cache directory (both "tl
// refresh -cache" and "validate -tl-cache") that holds the trusted
// certificates a refresh produced, one PEM file per certificate.
//
// A cache directory intentionally holds trust anchors only, not full
// trusted-list qualification data (which trusted list a certificate came
// from, its declared service type and status, its territory): that data
// lives in the in-memory
// [github.com/utain/esig/dss/spi/tsl.TrustedListsCertificateSource] a
// [github.com/utain/esig/dss/tsl.TLValidationJob] produces, and this CLI
// does not serialize it. A document validated with -tl-cache therefore gets
// an accurate Indication/SubIndication - the chain is genuinely anchored -
// but its SignatureQualification reads "NA": nothing here can tell the
// validation process which trust service qualifies which certificate. A
// process that needs qualified verdicts must wire the tsl package's job
// output into dss.ValidateOptions.CertificateVerifier directly; see the
// examples.
const tlCacheCertsDir = "certs"

// writeTLCache saves certs as one PEM file per certificate under
// dir/tlCacheCertsDir, replacing any cache already there.
func writeTLCache(dir string, certs []*dss.CertificateToken) error {
	certsDir := filepath.Join(dir, tlCacheCertsDir)
	if err := os.RemoveAll(certsDir); err != nil {
		return fmt.Errorf("clearing %s: %w", certsDir, err)
	}
	if err := os.MkdirAll(certsDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", certsDir, err)
	}
	for i, cert := range certs {
		block := &pem.Block{Type: "CERTIFICATE", Bytes: cert.Encoded()}
		path := filepath.Join(certsDir, fmt.Sprintf("%04d.pem", i))
		if err := os.WriteFile(path, pem.EncodeToMemory(block), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// readTLCache loads every certificate a previous "tl refresh -cache dir"
// wrote.
func readTLCache(dir string) ([]*dss.CertificateToken, error) {
	certsDir := filepath.Join(dir, tlCacheCertsDir)
	entries, err := os.ReadDir(certsDir)
	if err != nil {
		return nil, fmt.Errorf("reading trusted-list cache %s: %w (run \"esig tl refresh -cache %s\" first)", certsDir, err, dir)
	}
	certs := make([]*dss.CertificateToken, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(certsDir, entry.Name())
		cert, err := dss.LoadCertificate(path)
		if err != nil {
			return nil, fmt.Errorf("loading cached certificate %s: %w", path, err)
		}
		certs = append(certs, cert)
	}
	return certs, nil
}
