package asic

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// abstractFilenameFactoryProbe is the Go stand-in for the anonymous
// `new AbstractFilenameFactory() { }` the oracle drives.
type abstractFilenameFactoryProbe struct {
	AbstractFilenameFactory
}

func newAbstractASiCFilenameFactoryProbe() *abstractFilenameFactoryProbe {
	probe := &abstractFilenameFactoryProbe{}
	probe.InitAbstractASiCFilenameFactory(probe)
	return probe
}

// TestAbstractASiCFilenameFactoryNextAvailableDocumentNameMatchesDSS pins the suffix numbering
// against upstream's own getNextAvailableDocumentName, driven through reflection by
// testdata/gen/ZipCoreDssOracle.java. The cases cover the plain increment, the recursion that skips
// an already-taken candidate, and the HashSet de-duplication of the existing-name collection.
func TestAbstractASiCFilenameFactoryNextAvailableDocumentNameMatchesDSS(t *testing.T) {
	oracle := loadDSSOracle(t)
	if len(oracle.FilenameSuffixes) == 0 {
		t.Fatal("the oracle carries no filename-suffix cases")
	}
	factory := newAbstractASiCFilenameFactoryProbe()
	for _, tc := range oracle.FilenameSuffixes {
		tc := tc
		t.Run(tc.Template+"/"+string(rune('0'+len(tc.Existing))), func(t *testing.T) {
			if got := factory.NextAvailableDocumentName(tc.Template, tc.Existing); got != tc.Result {
				t.Fatalf("nextAvailableDocumentName(%q, %v) = %q, DSS = %q", tc.Template, tc.Existing, got, tc.Result)
			}
		})
	}
}

// TestAbstractASiCFilenameFactoryNextAvailableDocumentNamePadsToThreeDigits pins the zero padding
// and its overflow past 999, where the number stops being padded.
func TestAbstractASiCFilenameFactoryNextAvailableDocumentNamePadsToThreeDigits(t *testing.T) {
	factory := newAbstractASiCFilenameFactoryProbe()
	existing := make([]string, 0, 1000)
	for i := 1; i <= 1000; i++ {
		got := factory.NextAvailableDocumentName("META-INF/signature001.p7s", existing)
		var want string
		switch {
		case i < 10:
			want = "META-INF/signature00" + string(rune('0'+i)) + ".p7s"
		case i == 1000:
			want = "META-INF/signature1000.p7s"
		default:
			want = ""
		}
		if want != "" && got != want {
			t.Fatalf("with %d existing names, got %q, want %q", len(existing), got, want)
		}
		existing = append(existing, got)
	}
}

// TestAbstractASiCFilenameFactoryWithMetaInfFolder pins the idempotent "META-INF/" prefixing.
func TestAbstractASiCFilenameFactoryWithMetaInfFolder(t *testing.T) {
	factory := newAbstractASiCFilenameFactoryProbe()
	if got := factory.WithMetaInfFolder("signature001.p7s"); got != "META-INF/signature001.p7s" {
		t.Errorf("withMetaInfFolder = %q", got)
	}
	if got := factory.WithMetaInfFolder("META-INF/signature001.p7s"); got != "META-INF/signature001.p7s" {
		t.Errorf("withMetaInfFolder is not idempotent: %q", got)
	}
}

// TestAbstractASiCFilenameFactoryAssertASiCContentIsValid pins the container-type guard, including
// its two requireNonNull panics.
func TestAbstractASiCFilenameFactoryAssertASiCContentIsValid(t *testing.T) {
	factory := newAbstractASiCFilenameFactoryProbe()

	asicContent := NewContent()
	asicContent.SetContainerType(enumerations.ASiCContainerTypeASiCS)
	if err := factory.AssertASiCContentIsValid(asicContent); err != nil {
		t.Errorf("ASiC-S must be accepted: %v", err)
	}
	asicContent.SetContainerType(enumerations.ASiCContainerTypeASiCE)
	if err := factory.AssertASiCContentIsValid(asicContent); err != nil {
		t.Errorf("ASiC-E must be accepted: %v", err)
	}

	func() {
		defer func() {
			if recovered := recover(); recovered != "Type of ASiC Container shall be defined!" {
				t.Fatalf("panic = %v, want the Java message", recovered)
			}
		}()
		_ = factory.AssertASiCContentIsValid(NewContent())
	}()

	func() {
		defer func() {
			if recovered := recover(); recovered != "ASiCContent shall be provided!" {
				t.Fatalf("panic = %v, want the Java message", recovered)
			}
		}()
		_ = factory.AssertASiCContentIsValid(nil)
	}()
}

// TestAbstractASiCFilenameFactoryAssertFilenameValid pins the collision guard and its message.
func TestAbstractASiCFilenameFactoryAssertFilenameValid(t *testing.T) {
	factory := newAbstractASiCFilenameFactoryProbe()
	documents := []model.DSSDocument{model.NewInMemoryDocumentWithName([]byte("x"), "META-INF/signature001.p7s")}

	if err := factory.AssertFilenameValid("META-INF/signature002.p7s", documents); err != nil {
		t.Errorf("a free name must be accepted: %v", err)
	}
	err := factory.AssertFilenameValid("META-INF/signature001.p7s", documents)
	if err == nil {
		t.Fatal("a taken name must be rejected")
	}
	illegalInput, ok := err.(*exception.IllegalInputException)
	want := "The filename 'META-INF/signature001.p7s' cannot be used, as a document of the same name is already present within the container!"
	if !ok || illegalInput.Message != want {
		t.Fatalf("error = %v, want IllegalInputException %q", err, want)
	}
}

// TestAbstractASiCFilenameFactoryValidDataPackageFilename pins the two shape rules a data package
// filename must satisfy.
func TestAbstractASiCFilenameFactoryValidDataPackageFilename(t *testing.T) {
	factory := newAbstractASiCFilenameFactoryProbe()
	asicContent := NewContent()

	if got, err := factory.ValidDataPackageFilename("package.zip", asicContent); err != nil || got != "package.zip" {
		t.Errorf("validDataPackageFilename = %q (err %v)", got, err)
	}
	if _, err := factory.ValidDataPackageFilename("folder/package.zip", asicContent); err == nil {
		t.Error("a data package inside a folder must be rejected")
	} else if err.Error() != "A data package file within ASiC container shall be on the root level!" {
		t.Errorf("error = %q", err.Error())
	}
	if _, err := factory.ValidDataPackageFilename("package.7z", asicContent); err == nil {
		t.Error("a non-zip data package must be rejected")
	} else if err.Error() != "A data package filename within ASiC container shall ends with '.zip'!" {
		t.Errorf("error = %q", err.Error())
	}
	// The extension check is case-insensitive.
	if got, err := factory.ValidDataPackageFilename("PACKAGE.ZIP", asicContent); err != nil || got != "PACKAGE.ZIP" {
		t.Errorf("validDataPackageFilename(PACKAGE.ZIP) = %q (err %v)", got, err)
	}
}

// TestAbstractASiCFilenameFactoryValidEvidenceRecordManifestFilename pins the
// "META-INF/ASiCEvidenceRecordManifest*.xml" template, whose "META-INF/" prefix is optional.
func TestAbstractASiCFilenameFactoryValidEvidenceRecordManifestFilename(t *testing.T) {
	factory := newAbstractASiCFilenameFactoryProbe()
	asicContent := NewContent()

	for _, input := range []string{
		"ASiCEvidenceRecordManifest001.xml",
		"META-INF/ASiCEvidenceRecordManifest001.xml",
	} {
		got, err := factory.ValidEvidenceRecordManifestFilename(input, asicContent)
		if err != nil || got != "META-INF/ASiCEvidenceRecordManifest001.xml" {
			t.Errorf("validEvidenceRecordManifestFilename(%q) = %q (err %v)", input, got, err)
		}
	}

	want := "ASiC evidence record manifest file within ASiC container shall match the template 'META-INF/ASiCEvidenceRecordManifest*.xml'!"
	for _, input := range []string{"META-INF/OtherManifest001.xml", "META-INF/ASiCEvidenceRecordManifest001.p7s"} {
		if _, err := factory.ValidEvidenceRecordManifestFilename(input, asicContent); err == nil {
			t.Errorf("%q must be rejected", input)
		} else if err.Error() != want {
			t.Errorf("error = %q, want %q", err.Error(), want)
		}
	}
}

// abstractFilenameFactoryOverridingProbe overrides isAvailableName, the protected method
// getDocumentNameRecursively and assertFilenameValid call on themselves. Static Go dispatch would
// drop the override; the Overrides interface is what keeps it in play.
type abstractFilenameFactoryOverridingProbe struct {
	AbstractFilenameFactory
	calls int
}

func (p *abstractFilenameFactoryOverridingProbe) IsAvailableName(filename string, restrictedNames []string) bool {
	p.calls++
	// Pretend every name is taken until the fourth candidate.
	return p.calls > 3
}

// TestAbstractASiCFilenameFactoryRoutesSelfCallsThroughOverrides guards the virtual-dispatch bug
// class: a subclass override of isAvailableName must be honoured by getNextAvailableDocumentName.
func TestAbstractASiCFilenameFactoryRoutesSelfCallsThroughOverrides(t *testing.T) {
	probe := &abstractFilenameFactoryOverridingProbe{}
	probe.InitAbstractASiCFilenameFactory(probe)

	got := probe.NextAvailableDocumentName("META-INF/signature001.p7s", nil)
	if probe.calls != 4 {
		t.Fatalf("isAvailableName was called %d times, want 4 - the override was not reached", probe.calls)
	}
	if got != "META-INF/signature004.p7s" {
		t.Fatalf("nextAvailableDocumentName = %q, want META-INF/signature004.p7s", got)
	}
}
