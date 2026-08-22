// Cross-validation harness, direction GO -> UPSTREAM, ASiC-with-XAdES: runs
// testdata/crossgen (a standalone `go run` program - see its own doc comment) to build ASiC-S and
// ASiC-E containers at XAdES baseline B and T with this package's own ASiCWithXAdESService and a
// real PKCS#12 test key, then hands the output to testdata/crossgen/CrossGenValidator.java, which
// loads each container with upstream DSS 6.5.RC1's own SignedDocumentValidator/
// SignedDocumentDiagnosticDataBuilder and asserts the container type is recognized, the signature
// is intact, the signing certificate is identified, and the level is recognized. Upstream DSS
// accepting what this port produced is the actual proof; nothing here re-derives the answer with
// this package's own code.
//
// Modeled on cades/cades_downstream_cross_validation_test.go,
// xades/xades_downstream_cross_validation_test.go and its sibling
// asic/cades/asic_downstream_cross_validation_test.go.
//
// Skipped under -short (it shells out to `go run`, `mvn`, `javac` and `java`) and skipped, with
// the exact detection reason logged, when Java/Maven or a built upstream DSS checkout with
// dss-validation AND dss-asic-xades are not available.
package xades

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// crossGenFixtures are the "<file>:<expectedContainerType>:<expectedLevel>" arguments
// CrossGenValidator.java takes, one per file testdata/crossgen's generator writes.
var crossGenFixtures = []string{
	"asics-xades-b.scs:ASiC_S:XAdES_BASELINE_B",
	"asics-xades-t.scs:ASiC_S:XAdES_BASELINE_T",
	"asice-xades-b.sce:ASiC_E:XAdES_BASELINE_B",
	"asice-xades-t.sce:ASiC_E:XAdES_BASELINE_T",
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
// a built upstream DSS checkout with dss-validation AND dss-asic-xades compiled (the modules
// CrossGenValidator.java needs - dss-asic-xades registers the ASiC container
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
					"`mvn -o -pl dss-validation,dss-asic-xades,dss-asic-common,dss-xades,dss-xml-utils,dss-xml-common,specs-trusted-list,specs-xades,specs-xmldsig,dss-cades,dss-cms,dss-cms-object,dss-document -am install -DskipTests` in it): %v",
				module, classesDir, err)
		}
	}
	return ""
}

// crossGenNeededModules are the upstream modules CrossGenValidator.java needs their
// target/classes directory for directly, rather than through a single module's own resolved
// dependency jars: dss-asic-xades/dss-asic-common/dss-xades/dss-xml-utils/dss-xml-common/
// specs-trusted-list/specs-xades/specs-xmldsig/dss-cades/dss-cms/dss-cms-object/dss-document are
// reactor modules this checkout builds locally (source of the ASiC-with-XAdES implementation
// itself and its own compile-time dependencies, none of which dss-validation's own resolved
// dependency:build-classpath resolves, since dss-validation itself does not depend on
// dss-asic-xades) and dss-validation is the module CrossGenValidator.java calls into.
//
// Includes dss-cades/dss-cms/dss-cms-object for the same reason
// xades/xades_downstream_cross_validation_test.go needs them: SignedDocumentValidator.
// fromDocument's isSupported() probe instantiates every registered DocumentValidatorFactory,
// CAdES's included, regardless of which format the document being validated actually is.
var crossGenNeededModules = []string{
	"dss-validation", "dss-asic-xades", "dss-asic-common",
	"dss-xades", "dss-xml-utils", "dss-xml-common",
	"specs-trusted-list", "specs-xades", "specs-xmldsig",
	"dss-cades", "dss-cms", "dss-cms-object", "dss-document",
}

// buildUpstreamDSSClasspath assembles the classpath CrossGenValidator.java needs: the modules in
// crossGenNeededModules, plus the union of dss-validation's AND dss-asic-xades's runtime
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
	asicXadesDeps, err := mvnRuntimeClasspath(upstreamHome, "dss-asic-xades")
	if err != nil {
		return "", err
	}
	xadesDeps, err := mvnRuntimeClasspath(upstreamHome, "dss-xades")
	if err != nil {
		return "", err
	}

	parts := append(moduleClasses, validationDeps, asicXadesDeps, xadesDeps)
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
