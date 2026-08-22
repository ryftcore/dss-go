// Cross-validation harness, direction GO -> UPSTREAM: the first end-to-end
// compatibility proof in the other direction from cades_upstream_cross_validation_test.go. It
// runs testdata/crossgen (a standalone `go run` program - see its own doc comment) to sign CAdES-B
// and CAdES-T documents with this package's own CAdESService and a real PKCS#12 test key, then
// hands the output to testdata/crossgen/CrossGenValidator.java, which loads each file with
// upstream DSS 6.5.RC1's own SignedDocumentValidator/SignedDocumentDiagnosticDataBuilder and
// asserts the signature is intact, the signing certificate is identified, and the level is
// recognized. Upstream DSS accepting what this port produced is the actual proof; nothing here
// re-derives the answer with this package's own code.
//
// Skipped under -short (it shells out to `go run`, `mvn`, `javac` and `java`) and skipped,
// with the exact detection reason logged, when Java/Maven or a built upstream DSS checkout with
// dss-validation are not available - matching how a machine without java or with a
// half-built dss-upstream checkout should behave (plain `go test` never requires either), while
// still actually running and asserting a hard pass here, where both are present.
package cades

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// crossGenFixtures are the "<file>:<expectedLevel>[:<detachedContentFile>]" arguments
// CrossGenValidator.java takes, one per file testdata/crossgen's generator writes.
var crossGenFixtures = []string{
	"cades-b-enveloping.p7m:CAdES_BASELINE_B",
	"cades-t-enveloping.p7m:CAdES_BASELINE_T",
	"cades-b-detached.p7s:CAdES_BASELINE_B:cades-b-detached-content.bin",
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

	// A red result here - upstream DSS rejecting a Go-produced signature - is a genuine
	// cross-validation finding, not a test bug: it is reported as a hard failure with the exact
	// upstream error captured above, never silently downgraded to a skip or a weakened assertion.
	if err != nil || !strings.Contains(outputText, "ALL OK") {
		t.Fatalf("upstream DSS did not accept every Go-produced signature (see output above): %v", err)
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

// detectJavaAndUpstreamDSS reports whether java/javac/mvn are on PATH and upstreamHome looks
// like a built upstream DSS checkout with dss-validation compiled (the module
// CrossGenValidator.java needs). Returns a non-empty skip reason when anything is missing.
func detectJavaAndUpstreamDSS(upstreamHome string) string {
	for _, tool := range []string{"java", "javac", "mvn"} {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Sprintf("%s not found on PATH: %v", tool, err)
		}
	}
	dssValidationClasses := filepath.Join(upstreamHome, "dss-validation", "target", "classes")
	if info, err := os.Stat(dssValidationClasses); err != nil || !info.IsDir() {
		return fmt.Sprintf(
			"upstream DSS checkout not found or dss-validation not built at %s (set DSS_UPSTREAM_HOME, or run "+
				"`mvn -o -pl dss-validation,dss-cades,dss-cms-object -am install -DskipTests` in it): %v",
			dssValidationClasses, err)
	}
	return ""
}

// crossGenNeededModules are the upstream modules CrossGenValidator.java needs their
// target/classes directory for directly, rather than through dss-validation's own resolved
// dependency jars: dss-cades/dss-cms/dss-cms-object/dss-document are reactor modules this
// checkout builds locally (source of the CMS/CAdES/generic-signature implementation itself,
// including CAdESUtils's InMemoryResourcesHandlerBuilder dependency on dss-document) and
// dss-validation is the module CrossGenValidator.java calls into.
//
// Deliberately NOT a glob over every built module's target/classes: earlier revisions of this
// test globbed all of them, which pulled dss-evidence-record-asn1's classes onto the classpath
// too. Its ASN1EvidenceRecordAnalyzerFactory registers itself as a java.util.ServiceLoader
// provider that SignedDocumentValidator.fromDocument's isSupported() probe instantiates for
// every document regardless of whether an evidence record is even in play, and that provider
// needs commons-io - a dependency of dss-utils-apache-commons, not of dss-validation itself, so
// it never appears in dss-validation's own runtime dependency:build-classpath output. The
// resulting NoClassDefFoundError had nothing to do with what this test actually exercises.
var crossGenNeededModules = []string{"dss-validation", "dss-cades", "dss-cms", "dss-cms-object", "dss-document"}

// buildUpstreamDSSClasspath assembles the classpath CrossGenValidator.java needs: the modules in
// crossGenNeededModules, plus dss-validation's runtime dependency jars (BouncyCastle, the
// JAXB/XML stack, and every other module - dss-spi, dss-model, dss-document, ... - already
// resolved as local-repository jars since this checkout installs them) resolved offline via
// `mvn -o dependency:build-classpath`.
func buildUpstreamDSSClasspath(upstreamHome string) (string, error) {
	var moduleClasses []string
	for _, module := range crossGenNeededModules {
		classesDir := filepath.Join(upstreamHome, module, "target", "classes")
		if info, err := os.Stat(classesDir); err != nil || !info.IsDir() {
			return "", fmt.Errorf("required module %s is not built at %s: %v", module, classesDir, err)
		}
		moduleClasses = append(moduleClasses, classesDir)
	}

	classpathFile := filepath.Join(os.TempDir(), fmt.Sprintf("dss-crossgen-classpath-%d.txt", os.Getpid()))
	defer func() { _ = os.Remove(classpathFile) }()
	mvn := exec.Command("mvn", "-q", "-o", "-pl", "dss-validation",
		"dependency:build-classpath", "-Dmdep.outputFile="+classpathFile, "-Dmdep.includeScope=runtime")
	mvn.Dir = upstreamHome
	if output, err := mvn.CombinedOutput(); err != nil {
		return "", fmt.Errorf("mvn dependency:build-classpath: %w\n%s", err, output)
	}
	dependencyClasspath, err := os.ReadFile(classpathFile)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", classpathFile, err)
	}

	parts := append(moduleClasses, strings.TrimSpace(string(dependencyClasspath)))
	return strings.Join(parts, string(os.PathListSeparator)), nil
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
