// Package fixtures locates the tiny signing fixtures the examples under
// dss/examples share with the root package's own tests
// (dss/testdata/README.md): a self-signed RSA test key store, an EC test key
// acting as a local RFC 3161 time-stamp authority, and a minimal sample PDF.
// Reusing them here - instead of copying fixtures into examples/testdata -
// keeps the module zip small, which is the same goal the repository-root
// corpus/ split serves for heavier per-package testdata.
//
// Path resolves fixture names to an absolute path computed from this file's
// own location via runtime.Caller, not from the process's current working
// directory - so every example works the same whether it is run with
// `go run ./examples/01-sign-pdf-pades` from the module root or from inside
// its own directory.
package fixtures

import (
	"path/filepath"
	"runtime"
)

// moduleRoot is the dss/ module root, computed once from this source file's
// own path.
var moduleRoot = func() string {
	_, file, _, _ := runtime.Caller(0)
	// This file lives at dss/examples/internal/fixtures/fixtures.go; climb
	// three directories to reach dss/.
	return filepath.Join(filepath.Dir(file), "..", "..", "..")
}()

// Path returns the absolute path of the fixture named rel inside
// dss/testdata, for example Path("signer_rsa.p12").
func Path(rel string) string {
	return filepath.Join(moduleRoot, "testdata", rel)
}

// Module returns the absolute path of rel measured from the dss module
// root, for the rare example that reaches into a package's own resources
// instead of the shared testdata - for example
// Module("policy/resources/constraint.xml") to load the ETSI policy the
// library ships as its default.
func Module(rel string) string {
	return filepath.Join(moduleRoot, filepath.FromSlash(rel))
}
