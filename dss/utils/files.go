// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1),
// filesystem-related methods, matching org.apache.commons.io.FileUtils
// semantics.

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CleanDirectory cleans (deletes all contents of, without deleting the
// directory itself) the given directory.
//
// Mirrors commons-io FileUtils.cleanDirectory: an error is returned if
// directory does not exist or is not a directory.
func CleanDirectory(directory string) error {
	info, err := os.Stat(directory)
	if err != nil {
		return fmt.Errorf("%s does not exist", directory)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", directory)
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(directory, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

// ListFiles lists all files under folder with one of the given
// extensions (matched case-sensitively, without the leading '.').
// A nil/empty extensions list matches all files. When recursive is true,
// sub-directories are descended into.
func ListFiles(folder string, extensions []string, recursive bool) ([]string, error) {
	var result []string

	var walk func(dir string, depth int) error
	walk = func(dir string, depth int) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			full := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				if recursive {
					if err := walk(full, depth+1); err != nil {
						return err
					}
				}
				continue
			}
			if len(extensions) == 0 || hasAnyExtension(entry.Name(), extensions) {
				result = append(result, full)
			}
		}
		return nil
	}

	if err := walk(folder, 0); err != nil {
		return nil, err
	}
	return result, nil
}

func hasAnyExtension(filename string, extensions []string) bool {
	ext := GetFileNameExtension(filename)
	for _, e := range extensions {
		if strings.TrimPrefix(e, ".") == ext {
			return true
		}
	}
	return false
}
