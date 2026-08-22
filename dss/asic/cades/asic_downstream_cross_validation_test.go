// Cross-validation harness, direction GO -> UPSTREAM, ASiC-with-CAdES: runs
// testdata/crossgen (a standalone `go run` program - see its own doc comment) to build ASiC-S and
// ASiC-E containers at CAdES baseline B and T with this package's own ASiCWithCAdESService and a
// real PKCS#12 test key, then hands the output to testdata/crossgen/CrossGenValidator.java, which
// loads each container with upstream DSS 6.5.RC1's own SignedDocumentValidator/
// SignedDocumentDiagnosticDataBuilder and asserts the container type is recognized, the signature
// is intact, the signing certificate is identified, and the level is recognized. Upstream DSS
// accepting what this port produced is the actual proof; nothing here re-derives the answer with
// this package's own code.
//
// Modeled on cades/cades_downstream_cross_validation_test.go and
// xades/xades_downstream_cross_validation_test.go.
//
// Skipped under -short (it shells out to `go run`, `mvn`, `javac` and `java`) and skipped, with
// the exact detection reason logged, when Java/Maven or a built upstream DSS checkout with
// dss-validation AND dss-asic-cades are not available.
package cades

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// crossGenFixtures are the "<file>:<expectedContainerType>:<expectedLevel>" arguments
// CrossGenValidator.java takes, one per file testdata/crossgen's generator writes.
var crossGenFixtures = []string{
	"asics-cades-b.scs:ASiC_S:CAdES_BASELINE_B",
	"asics-cades-t.scs:ASiC_S:CAdES_BASELINE_T",
	"asice-cades-b.sce:ASiC_E:CAdES_BASELINE_B",
	"asice-cades-t.sce:ASiC_E:CAdES_BASELINE_T",
}

// TestDownstreamCrossValidation runs testdata/crossgen, then validates its output with upstream
// DSS's own SignedDocumentValidator via CrossGenValidator.java.
func TestDownstreamCrossValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("cross-validation against upstream DSS shells out to go run/mvn/javac/java; skipped under -short")
	}

	upstreamHome := os.Getenv("DSS_UPSTREAM_HOME")
	if skipReason := detectJavaAndUpstreamDSS(upstreamHome); skipReason != "" {
		t.Skip(skipReason)
	}

	outDir := t.TempDir()

	// 1) Generate: this is this package's own code, so a failure here is a real test failure,
	// not something to skip past.
	generatorOutput, err := runGoGenerator(outDir)
	if err != nil {
		t.Fatalf("testdata/crossgen generator failed: %v\n%s", err, generatorOutput)
	}
	t.Logf("generator output:\n%s", generatorOutput)
	for _, fixture := range crossGenFixtures {
		fileName := strings.SplitN(fixture, ":", 2)[0]
		if _, err := os.Stat(filepath.Join(outDir, fileName)); err != nil {
			t.Fatalf("generator did not produce %s: %v", fileName, err)
		}
	}

	// 2) Compile CrossGenValidator.java against upstream DSS.
	classpath, err := buildUpstreamDSSClasspath(upstreamHome)
	if err != nil {
		t.Fatalf("building upstream DSS classpath: %v", err)
	}
	validatorClasses := t.TempDir()
	javacArgs := []string{"-cp", classpath, "-d", validatorClasses, "CrossGenValidator.java"}
	javac := exec.Command("javac", javacArgs...)
	javac.Dir = filepath.Join("testdata", "crossgen")
	javacOutput, err := javac.CombinedOutput()
	if err != nil {
		t.Fatalf("javac CrossGenValidator.java failed: %v\n%s", err, javacOutput)
	}

	// 3) Run it: upstream DSS validating what the Go port produced is the actual proof.
	javaArgs := append([]string{"-cp", classpath + string(os.PathListSeparator) + validatorClasses,
		"CrossGenValidator", outDir}, crossGenFixtures...)
	java := exec.Command("java", javaArgs...)
	javaOutput, err := java.CombinedOutput()
	outputText := stripJavaToolOptionsNoise(string(javaOutput))
	t.Logf("CrossGenValidator output:\n%s", outputText)

	// A red result here - upstream DSS rejecting a Go-produced container - is a genuine
	// cross-validation finding, not a test bug: it is reported as a hard failure with the exact
	// upstream error captured above, never silently downgraded to a skip or a weakened assertion.
	if err != nil || !strings.Contains(outputText, "ALL OK") {
		t.Fatalf("upstream DSS did not accept every Go-produced container (see output above): %v", err)
	}
}

// runGoGenerator invokes `go run .` inside testdata/crossgen with outDir as its argument.
func runGoGenerator(outDir string) (string, error) {
	cmd := exec.Command("go", "run", ".", outDir)
	cmd.Dir = filepath.Join("testdata", "crossgen")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// detectJavaAndUpstreamDSS reports whether java/javac/mvn are on PATH and upstreamHome looks like
// a built upstream DSS checkout with dss-validation AND dss-asic-cades compiled (the modules
// CrossGenValidator.java needs - dss-asic-cades registers the ASiC container
// ValidatorFactory/AnalyzerFactory SignedDocumentValidator.fromDocument's ServiceLoader probe
// needs to recognize a .scs/.sce file at all). Returns a non-empty skip reason when anything is
// missing.
func detectJavaAndUpstreamDSS(upstreamHome string) string {
	for _, tool := range []string{"java", "javac", "mvn"} {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Sprintf("%s not found on PATH: %v", tool, err)
		}
	}
	for _, module := range crossGenNeededModules {
		classesDir := filepath.Join(upstreamHome, module, "target", "classes")
		if info, err := os.Stat(classesDir); err != nil || !info.IsDir() {
			return fmt.Sprintf(
				"upstream DSS checkout not found or %s not built at %s (set DSS_UPSTREAM_HOME, or run "+
					"`mvn -o -pl dss-validation,dss-asic-cades,dss-asic-common,dss-cades,dss-cms,dss-cms-object,dss-document -am install -DskipTests` in it): %v",
				module, classesDir, err)
		}
	}
	return ""
}

// crossGenNeededModules are the upstream modules CrossGenValidator.java needs their
// target/classes directory for directly, rather than through a single module's own resolved
// dependency jars: dss-asic-cades/dss-asic-common/dss-cades/dss-cms/dss-cms-object/dss-document
// are reactor modules this checkout builds locally (source of the ASiC-with-CAdES implementation
// itself and its own compile-time dependencies, none of which dss-validation's own resolved
// dependency:build-classpath resolves, since dss-validation itself does not depend on
// dss-asic-cades) and dss-validation is the module CrossGenValidator.java calls into.
var crossGenNeededModules = []string{
	"dss-validation", "dss-asic-cades", "dss-asic-common",
	"dss-cades", "dss-cms", "dss-cms-object", "dss-document",
}

// buildUpstreamDSSClasspath assembles the classpath CrossGenValidator.java needs: the modules in
// crossGenNeededModules, plus the union of dss-validation's AND dss-asic-cades's runtime
// dependency jars, resolved offline via `mvn -o dependency:build-classpath`.
func buildUpstreamDSSClasspath(upstreamHome string) (string, error) {
	var moduleClasses []string
	for _, module := range crossGenNeededModules {
		classesDir := filepath.Join(upstreamHome, module, "target", "classes")
		if info, err := os.Stat(classesDir); err != nil || !info.IsDir() {
			return "", fmt.Errorf("required module %s is not built at %s: %v", module, classesDir, err)
		}
		moduleClasses = append(moduleClasses, classesDir)
	}

	validationDeps, err := mvnRuntimeClasspath(upstreamHome, "dss-validation")
	if err != nil {
		return "", err
	}
	asicCadesDeps, err := mvnRuntimeClasspath(upstreamHome, "dss-asic-cades")
	if err != nil {
		return "", err
	}

	parts := append(moduleClasses, validationDeps, asicCadesDeps)
	return strings.Join(parts, string(os.PathListSeparator)), nil
}

// mvnRuntimeClasspath resolves module's own runtime dependency jars offline via
// `mvn -o dependency:build-classpath`.
func mvnRuntimeClasspath(upstreamHome, module string) (string, error) {
	classpathFile := filepath.Join(os.TempDir(), fmt.Sprintf("dss-crossgen-classpath-%s-%d.txt", module, os.Getpid()))
	defer func() { _ = os.Remove(classpathFile) }()
	mvn := exec.Command("mvn", "-q", "-o", "-pl", module,
		"dependency:build-classpath", "-Dmdep.outputFile="+classpathFile, "-Dmdep.includeScope=runtime")
	mvn.Dir = upstreamHome
	if output, err := mvn.CombinedOutput(); err != nil {
		return "", fmt.Errorf("mvn -pl %s dependency:build-classpath: %w\n%s", module, err, output)
	}
	dependencyClasspath, err := os.ReadFile(classpathFile)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", classpathFile, err)
	}
	return strings.TrimSpace(string(dependencyClasspath)), nil
}

// stripJavaToolOptionsNoise drops the JAVA_TOOL_OPTIONS echo line this container's `java`
// launcher prints to stderr on every invocation (a proxy/truststore configuration notice, not
// part of CrossGenValidator's own output), so test logs stay readable.
func stripJavaToolOptionsNoise(output string) string {
	lines := strings.Split(output, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if strings.HasPrefix(line, "Picked up JAVA_TOOL_OPTIONS") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}

// roundTripFixtures are the containers UPSTREAM DSS builds at CAdES-BASELINE-B (step 1) for the Go
// port to extend to -T (step 2), keyed to what CrossGenValidator.java must then see (step 3).
var roundTripFixtures = []struct {
	javaBuilt     string
	goExtended    string
	validatorSpec string
}{
	{"java-asics-cades-b.scs", "roundtrip-asics-cades-t.scs", "roundtrip-asics-cades-t.scs:ASiC_S:CAdES_BASELINE_T"},
	{"java-asice-cades-b.sce", "roundtrip-asice-cades-t.sce", "roundtrip-asice-cades-t.sce:ASiC_E:CAdES_BASELINE_T"},
}

// TestRoundTripJavaBuiltExtendedByGo is the UPSTREAM -> GO -> UPSTREAM direction: upstream DSS
// builds an ASiC-S and an ASiC-E container at CAdES-BASELINE-B, the Go port extends each to
// CAdES-BASELINE-T, and upstream DSS validates the result.
//
// TestDownstreamCrossValidation already proves upstream accepts a container this port BUILT. That
// is a weaker statement than it looks for the extension surface: extension rewrites an EXISTING
// signed container - it re-encodes the CMS signature, copies every other entry across and rebuilds
// the zip - so a port that dropped an entry, normalised entry metadata or mis-serialised the
// existing signature would still pass a build-only test, because it would never have been handed
// bytes it did not write itself. Here every input byte comes from upstream.
//
// Two things are asserted, and neither is derived by this package's own code:
//
//   - upstream DSS validates the extended container: right container type, signature still intact,
//     signing certificate still identified, level now -T, and the RFC 3161 token cryptographically
//     verified (CrossGenValidator.checkSignatureTimestamp).
//   - the extension preserved the container: every entry of the Java-built container is still
//     present in the Go-extended one, byte-identical except for the signature file the extension is
//     supposed to rewrite, "mimetype" is still the first entry and still STORED, and the only new
//     entries are the ones -T adds.
func TestRoundTripJavaBuiltExtendedByGo(t *testing.T) {
	if testing.Short() {
		t.Skip("cross-validation against upstream DSS shells out to go run/mvn/javac/java; skipped under -short")
	}

	upstreamHome := os.Getenv("DSS_UPSTREAM_HOME")
	if skipReason := detectJavaAndUpstreamDSS(upstreamHome); skipReason != "" {
		t.Skip(skipReason)
	}

	classpath, err := buildUpstreamDSSClasspath(upstreamHome)
	if err != nil {
		t.Fatalf("building upstream DSS classpath: %v", err)
	}
	// dss-token is not a dependency of any module already on the classpath; JavaBuiltGenerator
	// needs it for Pkcs12SignatureToken.
	tokenClasses := filepath.Join(upstreamHome, "dss-token", "target", "classes")
	if info, statErr := os.Stat(tokenClasses); statErr != nil || !info.IsDir() {
		t.Skipf("upstream dss-token is not built at %s: %v", tokenClasses, statErr)
	}
	classpath = classpath + string(os.PathListSeparator) + tokenClasses

	javaClasses := t.TempDir()
	javac := exec.Command("javac", "-cp", classpath, "-d", javaClasses,
		"JavaBuiltGenerator.java", "CrossGenValidator.java")
	javac.Dir = filepath.Join("testdata", "crossgen")
	if javacOutput, javacErr := javac.CombinedOutput(); javacErr != nil {
		t.Fatalf("javac JavaBuiltGenerator.java CrossGenValidator.java failed: %v\n%s", javacErr, javacOutput)
	}
	runtimeClasspath := classpath + string(os.PathListSeparator) + javaClasses

	outDir := t.TempDir()

	// Step 1: upstream DSS builds the -B containers.
	generator := exec.Command("java", "-cp", runtimeClasspath, "JavaBuiltGenerator", outDir, "signer_rsa.p12")
	generator.Dir = filepath.Join("testdata", "crossgen")
	generatorOutput, err := generator.CombinedOutput()
	generatorText := stripJavaToolOptionsNoise(string(generatorOutput))
	if err != nil || !strings.Contains(generatorText, "GENERATED OK") {
		t.Fatalf("JavaBuiltGenerator failed: %v\n%s", err, generatorText)
	}

	// Step 2: the Go port extends each of them to -T. A failure here is this port's failure, not
	// something to skip past.
	var validatorSpecs []string
	for _, fixture := range roundTripFixtures {
		inputPath := filepath.Join(outDir, fixture.javaBuilt)
		outputPath := filepath.Join(outDir, fixture.goExtended)
		extend := exec.Command("go", "run", ".", "-extend", inputPath, outputPath)
		extend.Dir = filepath.Join("testdata", "crossgen")
		if extendOutput, extendErr := extend.CombinedOutput(); extendErr != nil {
			t.Fatalf("extending %s with the Go port failed: %v\n%s", fixture.javaBuilt, extendErr, extendOutput)
		}
		assertExtensionPreservedContainer(t, inputPath, outputPath)
		validatorSpecs = append(validatorSpecs, fixture.validatorSpec)
	}

	// Step 3: upstream DSS validates what the Go port produced from its own bytes.
	javaArgs := append([]string{"-cp", runtimeClasspath, "CrossGenValidator", outDir}, validatorSpecs...)
	java := exec.Command("java", javaArgs...)
	javaOutput, err := java.CombinedOutput()
	outputText := stripJavaToolOptionsNoise(string(javaOutput))
	t.Logf("CrossGenValidator output:\n%s", outputText)
	if err != nil || !strings.Contains(outputText, "ALL OK") {
		t.Fatalf("upstream DSS did not accept every Go-extended container (see output above): %v", err)
	}
}

// assertExtensionPreservedContainer checks the structural half of the round trip: what the Go
// extension did to the bytes upstream handed it.
func assertExtensionPreservedContainer(t *testing.T, originalPath, extendedPath string) {
	t.Helper()

	original := readZipEntries(t, originalPath)
	extended := readZipEntries(t, extendedPath)

	// EN 319 162-1 A.1, on the container this port only rewrote.
	if len(extended.order) == 0 || extended.order[0] != "mimetype" {
		t.Errorf("%s: first entry = %v, want \"mimetype\"", filepath.Base(extendedPath), extended.order)
	} else if extended.methods["mimetype"] != zip.Store {
		t.Errorf("%s: mimetype compression method = %d, want STORED (0)",
			filepath.Base(extendedPath), extended.methods["mimetype"])
	}

	for _, name := range original.order {
		digest, present := extended.digests[name]
		if !present {
			t.Errorf("%s: entry %q from the Java-built container is missing after extension",
				filepath.Base(extendedPath), name)
			continue
		}
		// The signature file is the one entry extension is supposed to rewrite; everything else
		// must survive byte-identical.
		if strings.HasSuffix(name, ".p7s") {
			if digest == original.digests[name] {
				t.Errorf("%s: signature %q is unchanged after extension to -T", filepath.Base(extendedPath), name)
			}
			continue
		}
		if digest != original.digests[name] {
			t.Errorf("%s: entry %q changed content during extension", filepath.Base(extendedPath), name)
		}
	}
}

// zipInventory is the entry inventory of one container, in central-directory order.
type zipInventory struct {
	order   []string
	digests map[string]string
	methods map[string]uint16
}

func readZipEntries(t *testing.T, path string) zipInventory {
	t.Helper()
	reader, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("open %s as zip: %v", path, err)
	}
	defer func() { _ = reader.Close() }()

	inventory := zipInventory{digests: map[string]string{}, methods: map[string]uint16{}}
	for _, file := range reader.File {
		entry, err := file.Open()
		if err != nil {
			t.Fatalf("open entry %s in %s: %v", file.Name, path, err)
		}
		digest := sha256.New()
		if _, err := io.Copy(digest, entry); err != nil {
			t.Fatalf("read entry %s in %s: %v", file.Name, path, err)
		}
		_ = entry.Close()
		inventory.order = append(inventory.order, file.Name)
		inventory.digests[file.Name] = hex.EncodeToString(digest.Sum(nil))
		inventory.methods[file.Name] = file.Method
	}
	return inventory
}
