package fc

import (
	"testing"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
)

func manifestFile(signatureFilename string, entries ...string) *diagjaxb.XmlManifestFile {
	return &diagjaxb.XmlManifestFile{
		SignatureFilename: &signatureFilename,
		Entries:           &diagjaxb.EntriesWrapper{Items: entries},
	}
}

func filesCoveredCheck(manifests ...*diagjaxb.XmlManifestFile) *AbstractSignedAndTimestampedFilesCoveredCheck[*drjaxb.XmlFC] {
	return &AbstractSignedAndTimestampedFilesCoveredCheck[*drjaxb.XmlFC]{
		DiagnosticData: diagnostic.NewData(&diagjaxb.XmlDiagnosticData{
			ContainerInfo: &diagjaxb.XmlContainerInfo{
				ManifestFiles: &diagjaxb.ManifestFilesWrapper{Items: manifests},
			},
		}),
	}
}

// The recursion over nested manifests (one manifest listing the time-stamp of
// another) is exercised on acyclic layouts: nested and covered, nested and not
// covered, and a diamond that reaches the same manifest along two paths.
func TestCheckManifestFilesCovered_Acyclic(t *testing.T) {
	tests := []struct {
		name      string
		manifests []*diagjaxb.XmlManifestFile
		entries   []string
		want      bool
	}{
		{"no nested manifest", nil, []string{"doc.xml"}, true},
		{"nested manifest covered", []*diagjaxb.XmlManifestFile{
			manifestFile("t2.tst", "doc.xml"),
		}, []string{"doc.xml", "t2.tst"}, true},
		{"nested manifest entry not covered", []*diagjaxb.XmlManifestFile{
			manifestFile("t2.tst", "doc.xml", "other.xml"),
		}, []string{"doc.xml", "t2.tst"}, false},
		{"diamond", []*diagjaxb.XmlManifestFile{
			manifestFile("t2.tst", "t4.tst", "doc.xml"),
			manifestFile("t3.tst", "t4.tst", "doc.xml"),
			manifestFile("t4.tst", "doc.xml"),
		}, []string{"t2.tst", "t3.tst", "t4.tst", "doc.xml"}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := filesCoveredCheck(tc.manifests...)
			if got := c.CheckManifestFilesCovered(tc.entries); got != tc.want {
				t.Fatalf("CheckManifestFilesCovered(%q) = %v, want %v", tc.entries, got, tc.want)
			}
		})
	}
}

// Manifests that reference each other's time-stamp (or their own) make the
// upstream recursion unbounded: a Java StackOverflowError, which the JVM
// survives but a Go stack overflow does not (it is a fatal, unrecoverable
// error that would take the whole process down on a crafted container). The
// port must fail with an ordinary, recoverable panic instead.
func TestCheckManifestFilesCovered_CyclicManifestsPanic(t *testing.T) {
	tests := []struct {
		name      string
		manifests []*diagjaxb.XmlManifestFile
		entries   []string
	}{
		{"manifest lists its own time-stamp", []*diagjaxb.XmlManifestFile{
			manifestFile("t1.tst", "t1.tst"),
		}, []string{"t1.tst"}},
		{"two manifests list each other's time-stamp", []*diagjaxb.XmlManifestFile{
			manifestFile("t1.tst", "t2.tst"),
			manifestFile("t2.tst", "t1.tst"),
		}, []string{"t1.tst", "t2.tst"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := filesCoveredCheck(tc.manifests...)
			defer func() {
				if recover() == nil {
					t.Fatal("expected a recoverable panic for cyclic manifests")
				}
			}()
			c.CheckManifestFilesCovered(tc.entries)
		})
	}
}
