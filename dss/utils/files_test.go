// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1) test
// vectors, cross-checked against known org.apache.commons.io.FileUtils
// semantics.

package utils

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestCleanDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "b.txt"), []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := CleanDirectory(dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty directory, got %v", entries)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("directory itself should still exist: %v", err)
	}
}

func TestCleanDirectoryMissing(t *testing.T) {
	if err := CleanDirectory(filepath.Join(t.TempDir(), "does-not-exist")); err == nil {
		t.Error("expected error for missing directory")
	}
}

func TestListFiles(t *testing.T) {
	dir := t.TempDir()
	must := func(name, content string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	must("a.xml", "1")
	must("b.pdf", "2")
	must("c.txt", "3")

	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "d.xml"), []byte("4"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := ListFiles(dir, []string{"xml"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || filepath.Base(got[0]) != "a.xml" {
		t.Errorf("non-recursive got %v", got)
	}

	got, err = ListFiles(dir, []string{"xml"}, true)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, len(got))
	for i, g := range got {
		names[i] = filepath.Base(g)
	}
	sort.Strings(names)
	want := []string{"a.xml", "d.xml"}
	if len(names) != len(want) || names[0] != want[0] || names[1] != want[1] {
		t.Errorf("recursive got %v, want %v", names, want)
	}

	got, err = ListFiles(dir, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Errorf("expected all 3 top-level files, got %v", got)
	}
}
