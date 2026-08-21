// Package corpustest locates the repository's external oracle/fixture corpus
// (repo-root corpus/, outside the dss module) for tests whose testdata is too
// heavy to ship inside the `go get`-able module zip.
//
// The module directory keeps a small "smoke" subset of each package's
// testdata in place so `go test ./...` on a plain `go get` checkout still
// exercises real code paths. The remaining, heavier fixtures live in
// corpus/<package-relative-path>/testdata/ at the repository root, mirroring
// the package layout under dss/. A test that needs one of those fixtures
// calls Path with the same relative name it would have passed to
// filepath.Join("testdata", ...) before the file moved; Path finds the
// repository's corpus/ directory (present only in a full git checkout,
// never in a downloaded module) and resolves the fixture inside it. A test
// that instead reused another package's testdata by relative path (e.g.
// "../../other/pkg/testdata/x") calls RootPath with the same path measured
// from the dss module root instead.
//
// When corpus/ cannot be found at all (a `go get` of just the module, or any
// checkout without the sibling corpus/ directory) Path and RootPath call
// t.Skip so the calling test is skipped rather than failed. When corpus/ is
// present but the specific fixture is missing, they fail the test instead:
// that combination means the corpus is out of sync with the code and should
// be treated as a bug, not silently skipped.
package corpustest

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

// corpusDirName is the repo-root directory holding the external corpus.
const corpusDirName = "corpus"

// goModFile marks the dss module root while walking up from a package
// directory.
const goModFile = "go.mod"

var (
	rootOnce sync.Once
	rootDir  string // repo root containing corpus/ ("" if not found)
)

// findRepoRoot locates the directory that contains both the dss module (a
// go.mod) and a sibling corpus/ directory, walking up from dir. Every caller
// in a test run lives under the same module checkout, so the answer is a
// process-wide constant; it is computed once and cached.
func findRepoRoot(dir string) string {
	rootOnce.Do(func() {
		d := dir
		for {
			if _, err := os.Stat(filepath.Join(d, goModFile)); err == nil {
				// d is the dss module root; its parent is the repo root
				// that should hold corpus/.
				parent := filepath.Dir(d)
				if info, err := os.Stat(filepath.Join(parent, corpusDirName)); err == nil && info.IsDir() {
					rootDir = parent
				}
				return
			}
			next := filepath.Dir(d)
			if next == d {
				// Reached filesystem root without finding go.mod.
				return
			}
			d = next
		}
	})
	return rootDir
}

// callerPkgDir returns the directory of the source file that called Path,
// i.e. the test package requesting a fixture.
func callerPkgDir() string {
	// Caller(1): the frame that invoked Path (this function is only ever
	// called from Path, one level down).
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}

// Path resolves rel — a path relative to the calling test package's
// testdata/ directory, e.g. "asn1/cert.der" — to a file inside the
// repository's external corpus/ tree, mirroring the calling package's
// location under the dss module.
//
// If the repository's corpus/ directory cannot be found (the common case for
// a consumer's `go get`), Path skips the calling test via t.Skip. If corpus/
// is present but the resolved fixture does not exist, Path fails the test:
// that means the corpus is missing a fixture the code expects.
func Path(t testing.TB, rel string) string {
	t.Helper()
	return resolve(t, callerPkgDir(), rel)
}

// RootPath resolves rel — a path relative to the dss module root, e.g.
// "diagnostic/jaxb/testdata/oracle/cert_01.xml" — directly against the
// repository's external corpus/ tree, without mirroring the calling
// package's own location.
//
// Use RootPath instead of Path when a test's fixture reference already
// crosses package boundaries (a "../../other/pkg/testdata/..." reference
// the module carried before its heavy fixtures moved into corpus/): rel
// should be the same module-root-relative path with the leading "../..."
// climb and any "testdata" segment written out in full, e.g. what
// filepath.Join("..", "..", "otherpkg", "testdata", "x") produced before.
//
// The skip/fail behavior matches Path.
func RootPath(t testing.TB, rel string) string {
	t.Helper()
	return resolveRoot(t, callerPkgDir(), rel)
}

// resolve implements Path against an explicit package directory, so the
// directory-walking logic can be exercised directly in this package's own
// tests without depending on runtime.Caller's compiled-in source paths.
func resolve(t testing.TB, pkgDir, rel string) string {
	t.Helper()

	root := requireCorpusRoot(t, pkgDir, rel)

	moduleRoot, mErr := findModuleRoot(pkgDir)
	if mErr != nil {
		t.Fatalf("corpustest: %v", mErr)
	}
	pkgRel, rErr := filepath.Rel(moduleRoot, pkgDir)
	if rErr != nil {
		t.Fatalf("corpustest: computing package-relative path: %v", rErr)
	}

	full := filepath.Join(root, corpusDirName, pkgRel, "testdata", filepath.FromSlash(rel))
	return statOrFatal(t, full)
}

// resolveRoot implements RootPath against an explicit package directory; see
// resolve for why the directory is a parameter rather than always computed
// via runtime.Caller.
func resolveRoot(t testing.TB, pkgDir, rel string) string {
	t.Helper()

	root := requireCorpusRoot(t, pkgDir, rel)
	full := filepath.Join(root, corpusDirName, filepath.FromSlash(rel))
	return statOrFatal(t, full)
}

// requireCorpusRoot resolves the repo root or skips the test when corpus/
// cannot be found.
func requireCorpusRoot(t testing.TB, pkgDir, rel string) string {
	t.Helper()
	root := findRepoRoot(pkgDir)
	if root == "" {
		t.Skipf("corpustest: corpus/ not found outside the module (heavy fixture %q lives in the repository, not the module zip); skipping", rel)
	}
	return root
}

// statOrFatal returns full if it exists, else fails the test: a resolved
// corpus root missing the specific fixture means the corpus is out of sync
// with the code, not that the fixture is legitimately absent.
func statOrFatal(t testing.TB, full string) string {
	t.Helper()
	if _, statErr := os.Stat(full); statErr != nil {
		t.Fatalf("corpustest: fixture %s not found in corpus (%v); corpus present but out of sync", full, statErr)
	}
	return full
}

// findModuleRoot walks up from dir to the nearest ancestor containing
// go.mod (the dss module root).
func findModuleRoot(dir string) (string, error) {
	d := dir
	for {
		if _, err := os.Stat(filepath.Join(d, goModFile)); err == nil {
			return d, nil
		}
		next := filepath.Dir(d)
		if next == d {
			return "", os.ErrNotExist
		}
		d = next
	}
}
