// Cross-validation harness, direction GO -> UPSTREAM (task #12, JAdES extension): the second
// end-to-end compatibility proof in the other direction from
// jades_upstream_cross_validation_test.go. It runs testdata/crossgen (a standalone `go run`
// program - see its own doc comment) to sign JAdES-B and JAdES-T documents (compact AND flattened/
// full JSON serializations, plus a DETACHED 'sigD' case) with this package's own JAdESService and
// a real PKCS#12 test key (and, for -T, a real KeyEntityTSPSource-backed timestamp), then hands
// the output to testdata/crossgen/CrossGenValidator.java, which loads each file with upstream DSS
// 6.5.RC1's own SignedDocumentValidator/SignedDocumentDiagnosticDataBuilder and asserts the
// signature is intact, the signing certificate is identified, the level is recognized, and (for
// -T) the signature-timestamp itself validates. Upstream DSS accepting what this port produced is
// the actual proof; nothing here re-derives the answer with this package's own code. Mirrors
// cades_downstream_cross_validation_test.go, xades_downstream_cross_validation_test.go and
// pades_downstream_cross_validation_test.go.
//
// Skipped under -short (it shells out to `go run`, `mvn`, `javac` and `java`) and skipped, with
// the exact detection reason logged, when Java/Maven or a built upstream DSS checkout with
// dss-validation/dss-jades/specs-jades/dss-document are not available - matching how a machine
// without java or with a half-built dss-upstream checkout should behave (plain `go test` never
// requires either), while still actually running and asserting a hard pass here, where all are
// present.
package jades

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
// CrossGenValidator.java takes, one per file testdata/crossgen's generator writes (see its own
// doc comment): JAdES-B in compact and flattened JSON serialization, JAdES-T in flattened and
// full JSON serialization (compact carries no unprotected-header slot for a -T's 'etsiU', so
// there is no compact -T fixture - see main.go), and a DETACHED JAdES-B using the
// ObjectIdByURIHash 'sigD' mechanism.
var crossGenFixtures = []string{
	"jades-b-compact.json:JAdES_BASELINE_B",
	"jades-b-flattened.json:JAdES_BASELINE_B",
	"jades-t-flattened.json:JAdES_BASELINE_T",
	"jades-t-full.json:JAdES_BASELINE_T",
	"jades-b-detached.json:JAdES_BASELINE_B:jades-b-detached-content.txt",
}

// TestDownstreamCrossValidation runs testdata/crossgen, then validates its output with upstream
// DSS's own SignedDocumentValidator via CrossGenValidator.java.
func TestDownstreamCrossValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("cross-validation against upstream DSS shells out to go run/mvn/javac/java; skipped under -short")
	}

	upstreamHome := os.Getenv("DSS_UPSTREAM_HOME")
	if upstreamHome == "" {
		upstreamHome = "/home/user/dss-upstream"
	}
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
// like a built upstream DSS checkout with every module in crossGenNeededModules compiled.
// Returns a non-empty skip reason when anything is missing.
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
					"`mvn -o -pl dss-validation,dss-jades,specs-jades,specs-jws,dss-document -am install -DskipTests` in it): %v",
				module, classesDir, err)
		}
	}
	return ""
}

// crossGenNeededModules are the upstream modules CrossGenValidator.java needs their
// target/classes directory for directly, rather than through dss-validation's own resolved
// dependency jars: dss-jades/specs-jades/specs-jws/dss-document are reactor modules this checkout
// builds locally (source of the JAdES/JOSE-header-schema/generic-signature implementation itself)
// and dss-validation is the module CrossGenValidator.java calls into.
//
// Deliberately NOT a glob over every built module's target/classes: see the identical rationale
// in cades_downstream_cross_validation_test.go (dss-evidence-record-asn1's ServiceLoader provider
// pulling in an unrelated missing dependency).
var crossGenNeededModules = []string{"dss-validation", "dss-jades", "specs-jades", "specs-jws", "dss-document"}

// buildUpstreamDSSClasspath assembles the classpath CrossGenValidator.java needs: the modules in
// crossGenNeededModules, the jose4j jar (dss-jades's own compile-time dependency, found directly
// in the local Maven repository - see its own doc comment below), plus dss-validation's runtime
// dependency jars (BouncyCastle, and every other module - dss-spi, dss-model, dss-document, ... -
// already resolved as local-repository jars since this checkout installs them) resolved offline
// via `mvn -o dependency:build-classpath`.
func buildUpstreamDSSClasspath(upstreamHome string) (string, error) {
	var moduleClasses []string
	for _, module := range crossGenNeededModules {
		classesDir := filepath.Join(upstreamHome, module, "target", "classes")
		if info, err := os.Stat(classesDir); err != nil || !info.IsDir() {
			return "", fmt.Errorf("required module %s is not built at %s: %v", module, classesDir, err)
		}
		moduleClasses = append(moduleClasses, classesDir)
	}

	m2Repo := os.Getenv("M2_REPO")
	if m2Repo == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolving home directory for the local Maven repository: %w", err)
		}
		m2Repo = filepath.Join(home, ".m2", "repository")
	}
	jose4jJar, err := findNewestJose4jJar(m2Repo)
	if err != nil {
		return "", err
	}

	classpathFile := filepath.Join(os.TempDir(), fmt.Sprintf("dss-jades-crossgen-classpath-%d.txt", os.Getpid()))
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

	parts := append(moduleClasses, jose4jJar)
	parts = append(parts, strings.TrimSpace(string(dependencyClasspath)))
	return strings.Join(parts, string(os.PathListSeparator)), nil
}

// findNewestJose4jJar walks m2Repo for org/bitbucket/b_c/jose4j and returns the highest-versioned
// plain jar (skipping -sources.jar/-javadoc.jar) it finds: org.bitbucket.b_c:jose4j is dss-jades's
// own compile-scope dependency (the JOSE/JWS implementation library the port itself is modeled
// on), and - the same situation as org.apache.pdfbox:pdfbox for PAdES's crossgen test - never
// appears in dss-validation's own runtime dependency:build-classpath output, since dss-validation
// has no JOSE dependency of its own and dss-jades is never itself a dependency of dss-validation
// (dss-jades depends on dss-validation, not the reverse).
func findNewestJose4jJar(m2Repo string) (string, error) {
	groupDir := filepath.Join(m2Repo, "org", "bitbucket", "b_c", "jose4j")
	versions, err := os.ReadDir(groupDir)
	if err != nil {
		return "", fmt.Errorf("jose4j not found under the local Maven repository (%s): %w - run "+
			"the corresponding dependency:build-classpath once so it gets downloaded", groupDir, err)
	}
	var best, bestJar string
	for _, versionEntry := range versions {
		if !versionEntry.IsDir() {
			continue
		}
		version := versionEntry.Name()
		jar := filepath.Join(groupDir, version, "jose4j-"+version+".jar")
		if _, err := os.Stat(jar); err != nil {
			continue
		}
		if version > best {
			best = version
			bestJar = jar
		}
	}
	if bestJar == "" {
		return "", fmt.Errorf("no usable jose4j jar found under %s", groupDir)
	}
	return bestJar, nil
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
