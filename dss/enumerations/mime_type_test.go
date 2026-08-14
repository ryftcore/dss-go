// Ported from dss-enumerations/.../MimeType.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestMimeTypeGetFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{"simple", "file.txt", "txt"},
		{"multiple dots", "archive.tar.gz", "gz"},
		{"no extension", "file", ""},
		{"empty", "", ""},
		{"blank", "   ", ""},
		{"leading dot only", ".gitignore", ""}, // lastIndexOf('.') == 0, not > 0, so no extension
		{"trailing dot", "file.", ""},
		{"path with dot dir, no ext", "xxx.y/toto", "y/toto"}, // documented Java quirk: not handled specially
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MimeTypeGetFileExtension(tt.fileName)
			if got != tt.want {
				t.Errorf("MimeTypeGetFileExtension(%q) = %q, want %q", tt.fileName, got, tt.want)
			}
		})
	}
}

type fakeMimeType struct {
	mimeTypeString string
	extension      string
}

func (f *fakeMimeType) MimeTypeString() string { return f.mimeTypeString }
func (f *fakeMimeType) Extension() string      { return f.extension }

type fakeMimeTypeLoader struct {
	byString    map[string]MimeType
	byExtension map[string]MimeType
}

func (f *fakeMimeTypeLoader) FromMimeTypeString(s string) MimeType { return f.byString[s] }
func (f *fakeMimeTypeLoader) FromFileExtension(e string) MimeType  { return f.byExtension[e] }

func withMimeTypeLoaders(t *testing.T, loaders ...MimeTypeLoader) {
	t.Helper()
	saved := mimeTypeLoaderRegistry
	mimeTypeLoaderRegistry = nil
	for _, l := range loaders {
		RegisterMimeTypeLoader(l)
	}
	t.Cleanup(func() { mimeTypeLoaderRegistry = saved })
}

func TestMimeTypeFromMimeTypeString_LoaderMatch(t *testing.T) {
	xml := &fakeMimeType{mimeTypeString: "text/xml", extension: "xml"}
	loader := &fakeMimeTypeLoader{byString: map[string]MimeType{"text/xml": xml}}
	withMimeTypeLoaders(t, loader)

	got := MimeTypeFromMimeTypeString("text/xml")
	if got != MimeType(xml) {
		t.Errorf("MimeTypeFromMimeTypeString = %v, want %v", got, xml)
	}
}

func TestMimeTypeFromFileExtension_LoaderOrder(t *testing.T) {
	first := &fakeMimeType{mimeTypeString: "application/pdf", extension: "pdf"}
	loaderA := &fakeMimeTypeLoader{byExtension: map[string]MimeType{}}
	loaderB := &fakeMimeTypeLoader{byExtension: map[string]MimeType{"pdf": first}}
	withMimeTypeLoaders(t, loaderA, loaderB)

	got := MimeTypeFromFileExtension("pdf")
	if got != MimeType(first) {
		t.Errorf("MimeTypeFromFileExtension = %v, want loaderB's result %v", got, first)
	}
}

func TestMimeTypeFromFileName_LowercasesExtension(t *testing.T) {
	pdf := &fakeMimeType{mimeTypeString: "application/pdf", extension: "pdf"}
	loader := &fakeMimeTypeLoader{byExtension: map[string]MimeType{"pdf": pdf}}
	withMimeTypeLoaders(t, loader)

	got := MimeTypeFromFileName("report.PDF")
	if got != MimeType(pdf) {
		t.Errorf("MimeTypeFromFileName(report.PDF) = %v, want %v", got, pdf)
	}
}

func TestMimeTypeFromFilePath_UsesBaseName(t *testing.T) {
	pdf := &fakeMimeType{mimeTypeString: "application/pdf", extension: "pdf"}
	loader := &fakeMimeTypeLoader{byExtension: map[string]MimeType{"pdf": pdf}}
	withMimeTypeLoaders(t, loader)

	got := MimeTypeFromFilePath("/some/dir/report.pdf")
	if got != MimeType(pdf) {
		t.Errorf("MimeTypeFromFilePath = %v, want %v", got, pdf)
	}
}

// TestMimeTypeGetFileExtensionOK_NullVsEmpty pins upstream's distinction between the two
// "no extension" outcomes that both look like "" in Go: getFileExtension returns null only
// for a blank file name, and an empty string for a name with no '.' (or a leading '.').
func TestMimeTypeGetFileExtensionOK_NullVsEmpty(t *testing.T) {
	for _, tt := range []struct {
		fileName string
		want     string
		wantOK   bool
	}{
		{"report.pdf", "pdf", true},
		{"README", "", true},  // Java: "" (not null)
		{".hidden", "", true}, // Java: "" — lastIndexOf('.') == 0 is not > 0
		{"", "", false},       // Java: null
		{"   ", "", false},    // Java: null (blank)
	} {
		got, ok := MimeTypeGetFileExtensionOK(tt.fileName)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("MimeTypeGetFileExtensionOK(%q) = %q, %v; want %q, %v", tt.fileName, got, ok, tt.want, tt.wantOK)
		}
	}
}

// TestMimeTypeFromFileName_ConsultsLoadersForEmptyExtension checks that a name with no
// extension still reaches the registered loaders with an empty extension, as upstream does,
// rather than short-circuiting to BINARY. Only a blank name short-circuits.
func TestMimeTypeFromFileName_ConsultsLoadersForEmptyExtension(t *testing.T) {
	saved := mimeTypeLoaderRegistry
	mimeTypeLoaderRegistry = nil
	t.Cleanup(func() { mimeTypeLoaderRegistry = saved })

	var seen []string
	RegisterMimeTypeLoader(&recordingMimeTypeLoader{seen: &seen})

	MimeTypeFromFileName("README")
	if len(seen) != 1 || seen[0] != "" {
		t.Errorf("after MimeTypeFromFileName(README), loader saw %q, want one call with \"\"", seen)
	}

	seen = nil
	MimeTypeFromFileName("")
	if len(seen) != 0 {
		t.Errorf("after MimeTypeFromFileName(\"\"), loader saw %q, want no calls", seen)
	}
}

type recordingMimeTypeLoader struct{ seen *[]string }

func (r *recordingMimeTypeLoader) FromMimeTypeString(string) MimeType { return nil }
func (r *recordingMimeTypeLoader) FromFileExtension(ext string) MimeType {
	*r.seen = append(*r.seen, ext)
	return nil
}
