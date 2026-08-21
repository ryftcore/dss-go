// Must-reject parity: the acceptance boundary of the reader.
//
// TestOracleCorpus already enforces the boundary in both directions on every
// corpus file — pdfbox failed => Go must fail, pdfbox succeeded => Go must
// succeed. What it does not do is say *why* a document is rejected, or check
// the cases that need an input the oracle never feeds it (a wrong password, a
// public-key security handler). This file documents that boundary explicitly,
// so a future change that starts accepting a document upstream's reader layer
// refuses shows up as a named failure rather than a golden diff.
//
// DESIGN.md §2.6 (rejected encryption), §2.7 ("the one thing pdfbox does not
// recover from") and §4.2 (the error taxonomy) are the contract under test.

package pdf

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// hardFailures are the only two corpus documents pdfbox itself throws on:
// their page-tree root is not a dictionary. DESIGN.md §2.7 pins that we match
// that single hard failure and recover from everything else.
var hardFailures = []string{
	"EmptyPage-corrupted.pdf",
	"EmptyPage-corrupted2.pdf",
}

func TestMustReject_BrokenCatalog(t *testing.T) {
	for _, name := range hardFailures {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(corpusDir(t), name))
			if err != nil {
				t.Skipf("not vendored: %v", err)
			}
			_, err = OpenBytes(data, nil)
			if err == nil {
				t.Fatal("opened a document whose page tree root is not a dictionary")
			}
			if !errors.Is(err, ErrBrokenCatalog) {
				t.Errorf("error is %v, want ErrBrokenCatalog", err)
			}
		})
	}
}

// encryptedCorpus records, per vendored encrypted document, which password
// opens it. It is the same answer the oracle's `password` field carries, but
// spelled out here so the negative cases below can be stated against it.
var encryptedCorpus = []struct {
	path string
	pass string // the password that opens it; "" means the empty user password
	cfm  Name
}{
	{"protected/open_protected.pdf", " ", "AESV2"},
	{"protected/restricted_fields.pdf", "", "AESV3"},
	{"protected/edition_protected_signing_allowed_with_field_signed.pdf", "", "AESV2"},
	{"validation/encrypted.pdf", "", "V2"},
}

// TestMustReject_WrongPassword pins that a password which is neither the user
// nor the owner password is refused with ErrInvalidPassword, and — the part
// that actually matters — that the refusal is not silently downgraded into a
// document with garbage strings and streams.
func TestMustReject_WrongPassword(t *testing.T) {
	for _, tc := range encryptedCorpus {
		t.Run(tc.path, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(corpusDir(t), filepath.FromSlash(tc.path)))
			if err != nil {
				t.Skipf("not vendored: %v", err)
			}
			// The password that should work, works.
			var opts *Options
			if tc.pass != "" {
				opts = &Options{Password: []byte(tc.pass)}
			}
			doc, err := OpenBytes(data, opts)
			if err != nil {
				t.Fatalf("the documented password %q did not open the file: %v", tc.pass, err)
			}
			if !doc.IsEncrypted() {
				t.Fatal("document does not report itself as encrypted")
			}
			if enc := doc.Encryption(); enc == nil {
				t.Fatal("Encryption() is nil for an encrypted document")
			} else if enc.CFM != tc.cfm {
				t.Errorf("CFM %q, want %q", enc.CFM, tc.cfm)
			}

			// A password that is neither user nor owner must be refused.
			_, err = OpenBytes(data, &Options{Password: []byte("definitely-not-the-password")})
			if err == nil {
				t.Fatal("a wrong password opened the document")
			}
			if !errors.Is(err, ErrInvalidPassword) {
				t.Errorf("wrong-password error is %v, want ErrInvalidPassword", err)
			}
		})
	}
}

// TestMustReject_UnsupportedSecurityHandler pins §2.6's rejection of every
// handler that is not /Standard. Public-key handlers (/Adobe.PubSec) carry no
// password at all, so accepting one would mean handing the caller a document
// whose every string is ciphertext.
func TestMustReject_UnsupportedSecurityHandler(t *testing.T) {
	base, err := os.ReadFile(filepath.Join(corpusDir(t), "protected", "open_protected.pdf"))
	if err != nil {
		t.Skipf("not vendored: %v", err)
	}
	// Rewrite /Filter /Standard to a public-key handler of the same length, so
	// every offset in the file — and therefore the whole xref — stays valid and
	// the only thing that changed is the handler name.
	mutated := replaceOnce(t, base, []byte("/Filter/Standard"), []byte("/Filter/PubSec01"))
	if mutated == nil {
		mutated = replaceOnce(t, base, []byte("/Filter /Standard"), []byte("/Filter /PubSec01"))
	}
	if mutated == nil {
		t.Skip("could not locate a literal /Filter /Standard in the fixture")
	}
	_, err = OpenBytes(mutated, &Options{Password: []byte(" ")})
	if err == nil {
		t.Fatal("opened a document with a non-Standard security handler")
	}
	if !errors.Is(err, ErrUnsupportedSecurityHandler) {
		t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
	}
}

// replaceOnce returns a copy with the single occurrence of old replaced, or nil
// when old does not appear exactly once. new must be the same length as old so
// that no byte offset in the document moves.
func replaceOnce(t *testing.T, data, old, new []byte) []byte {
	t.Helper()
	if len(old) != len(new) {
		t.Fatalf("replaceOnce needs equal lengths, got %d and %d", len(old), len(new))
	}
	idx := indexOfBytes(data, old)
	if idx < 0 || indexOfBytes(data[idx+1:], old) >= 0 {
		return nil
	}
	out := make([]byte, len(data))
	copy(out, data)
	copy(out[idx:], new)
	return out
}

func indexOfBytes(hay, needle []byte) int {
	for i := 0; i+len(needle) <= len(hay); i++ {
		match := true
		for j := range needle {
			if hay[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// TestMustReject_NotAPDF pins the one input shape that is refused before any
// leniency applies: no %PDF- marker anywhere in the first 1024 bytes (§2.7 H1).
func TestMustReject_NotAPDF(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"empty", nil},
		{"plain text", []byte("this is not a pdf, not even slightly")},
		{"marker beyond the 1024-byte search window",
			append(make([]byte, 2048), []byte("%PDF-1.4\n%%EOF\n")...)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := OpenBytes(tc.data, nil); !errors.Is(err, ErrNotPDF) {
				t.Errorf("error is %v, want ErrNotPDF", err)
			}
		})
	}
}
