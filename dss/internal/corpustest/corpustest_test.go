package corpustest

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// fakeTB is a minimal testing.TB that records Skip/Fatal calls instead of
// acting on them, so resolve's branches can be verified deterministically.
type fakeTB struct {
	testing.TB
	skipped bool
	failed  bool
	msg     string
}

func (f *fakeTB) Helper() {}

func (f *fakeTB) Skipf(format string, args ...any) {
	f.skipped = true
	f.msg = fmt.Sprintf(format, args...)
	panic(skipSignal{})
}

func (f *fakeTB) Fatalf(format string, args ...any) {
	f.failed = true
	f.msg = fmt.Sprintf(format, args...)
	panic(fatalSignal{})
}

type skipSignal struct{}
type fatalSignal struct{}

func runResolve(f *fakeTB, pkgDir, rel string) (path string) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case skipSignal, fatalSignal:
				return
			default:
				panic(r)
			}
		}
	}()
	return resolve(f, pkgDir, rel)
}

// TestResolve_NoCorpusSkips builds an isolated fake module tree with no
// sibling corpus/ directory and verifies resolve skips rather than fails.
func TestResolve_NoCorpusSkips(t *testing.T) {
	tmp := t.TempDir()
	moduleRoot := filepath.Join(tmp, "dss")
	pkgDir := filepath.Join(moduleRoot, "somepkg")
	mustMkdirAll(t, pkgDir)
	mustWriteFile(t, filepath.Join(moduleRoot, "go.mod"), "module example.com/x\n")

	resetRootCache()
	f := &fakeTB{}
	got := runResolve(f, pkgDir, "foo.txt")
	if !f.skipped {
		t.Fatalf("expected Skipf to fire when corpus/ is absent, got failed=%v path=%q", f.failed, got)
	}
}

// TestResolve_CorpusPresentResolves builds a fake repo root with a sibling
// corpus/ directory mirroring the package path and verifies resolve returns
// the fixture without skipping or failing.
func TestResolve_CorpusPresentResolves(t *testing.T) {
	tmp := t.TempDir()
	repoRoot := tmp
	moduleRoot := filepath.Join(repoRoot, "dss")
	pkgDir := filepath.Join(moduleRoot, "somepkg")
	mustMkdirAll(t, pkgDir)
	mustWriteFile(t, filepath.Join(moduleRoot, "go.mod"), "module example.com/x\n")

	fixtureDir := filepath.Join(repoRoot, "corpus", "somepkg", "testdata")
	mustMkdirAll(t, fixtureDir)
	mustWriteFile(t, filepath.Join(fixtureDir, "foo.txt"), "hello")

	resetRootCache()
	f := &fakeTB{}
	got := runResolve(f, pkgDir, "foo.txt")
	if f.skipped || f.failed {
		t.Fatalf("expected resolve to succeed cleanly, got skipped=%v failed=%v msg=%q", f.skipped, f.failed, f.msg)
	}
	want := filepath.Join(fixtureDir, "foo.txt")
	if got != want {
		t.Fatalf("resolve returned %q, want %q", got, want)
	}
}

// TestResolve_NestedPackagePathMirrored verifies a multi-segment package
// path (e.g. validation/process/bbb/xcv) mirrors correctly into corpus/.
func TestResolve_NestedPackagePathMirrored(t *testing.T) {
	tmp := t.TempDir()
	repoRoot := tmp
	moduleRoot := filepath.Join(repoRoot, "dss")
	pkgDir := filepath.Join(moduleRoot, "validation", "process", "bbb", "xcv")
	mustMkdirAll(t, pkgDir)
	mustWriteFile(t, filepath.Join(moduleRoot, "go.mod"), "module example.com/x\n")

	fixtureDir := filepath.Join(repoRoot, "corpus", "validation", "process", "bbb", "xcv", "testdata")
	mustMkdirAll(t, fixtureDir)
	mustWriteFile(t, filepath.Join(fixtureDir, "kat.tsv"), "a\tb\n")

	resetRootCache()
	f := &fakeTB{}
	got := runResolve(f, pkgDir, "kat.tsv")
	if f.skipped || f.failed {
		t.Fatalf("expected resolve to succeed cleanly, got skipped=%v failed=%v msg=%q", f.skipped, f.failed, f.msg)
	}
	want := filepath.Join(fixtureDir, "kat.tsv")
	if got != want {
		t.Fatalf("resolve returned %q, want %q", got, want)
	}
}

// TestResolve_CorpusPresentMissingFixtureFails verifies that a present
// corpus/ tree missing the specific fixture is a hard failure, not a skip:
// that combination means the corpus is out of sync with the code.
func TestResolve_CorpusPresentMissingFixtureFails(t *testing.T) {
	tmp := t.TempDir()
	repoRoot := tmp
	moduleRoot := filepath.Join(repoRoot, "dss")
	pkgDir := filepath.Join(moduleRoot, "somepkg")
	mustMkdirAll(t, pkgDir)
	mustWriteFile(t, filepath.Join(moduleRoot, "go.mod"), "module example.com/x\n")
	mustMkdirAll(t, filepath.Join(repoRoot, "corpus"))

	resetRootCache()
	f := &fakeTB{}
	runResolve(f, pkgDir, "missing.txt")
	if !f.failed {
		t.Fatalf("expected Fatalf when corpus/ exists but the fixture doesn't, got skipped=%v", f.skipped)
	}
}

func runResolveRoot(f *fakeTB, pkgDir, rel string) (path string) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case skipSignal, fatalSignal:
				return
			default:
				panic(r)
			}
		}
	}()
	return resolveRoot(f, pkgDir, rel)
}

// TestResolveRoot_CrossPackageResolves verifies RootPath (via resolveRoot)
// resolves module-root-relative paths directly, ignoring the calling
// package's own location — the pattern a "../../other/pkg/testdata/x"
// cross-package reference needs.
func TestResolveRoot_CrossPackageResolves(t *testing.T) {
	tmp := t.TempDir()
	repoRoot := tmp
	moduleRoot := filepath.Join(repoRoot, "dss")
	callingPkgDir := filepath.Join(moduleRoot, "validation", "process", "bbb", "aov")
	mustMkdirAll(t, callingPkgDir)
	mustWriteFile(t, filepath.Join(moduleRoot, "go.mod"), "module example.com/x\n")

	// The fixture lives under a different package's mirror entirely.
	fixtureDir := filepath.Join(repoRoot, "corpus", "diagnostic", "jaxb", "testdata", "oracle")
	mustMkdirAll(t, fixtureDir)
	mustWriteFile(t, filepath.Join(fixtureDir, "cert_01.xml"), "<x/>")

	resetRootCache()
	f := &fakeTB{}
	got := runResolveRoot(f, callingPkgDir, filepath.Join("diagnostic", "jaxb", "testdata", "oracle", "cert_01.xml"))
	if f.skipped || f.failed {
		t.Fatalf("expected resolveRoot to succeed cleanly, got skipped=%v failed=%v msg=%q", f.skipped, f.failed, f.msg)
	}
	want := filepath.Join(fixtureDir, "cert_01.xml")
	if got != want {
		t.Fatalf("resolveRoot returned %q, want %q", got, want)
	}
}

// resetRootCache clears the process-wide sync.Once cache between subtests
// that each construct their own fake repo root.
func resetRootCache() {
	rootOnce = sync.Once{}
	rootDir = ""
}

func mustMkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", dir, err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}
