// Cross-validation harness, direction GO -> UPSTREAM: the second
// end-to-end compatibility proof in the other direction from
// pades_upstream_cross_validation_test.go. It runs testdata/crossgen (a standalone `go run`
// program - see its own doc comment) to sign three corpus PDFs of different xref styles with
// this package's own Service, producing invisible PAdES-B and PAdES-T signatures with a
// real PKCS#12 test key (and, for -T, a real KeyEntityTSPSource-backed timestamp), then hands the
// output to testdata/crossgen/CrossGenValidator.java, which loads each file with upstream DSS
// 6.5.RC1's own SignedDocumentValidator/SignedDocumentDiagnosticDataBuilder and asserts the
// signature is intact, the signing certificate is identified, and the level is recognized.
// Upstream DSS accepting what this port produced is the actual proof; nothing here re-derives the
// answer with this package's own code. Mirrors cades_downstream_cross_validation_test.go and
// xades_downstream_cross_validation_test.go.
//
// Skipped under -short (it shells out to `go run`, `mvn`, `javac` and `java`) and skipped, with
// the exact detection reason logged, when Java/Maven, a built upstream DSS checkout with
// dss-validation, or a built PDF backend (dss-pades-pdfbox) are not available - matching how a
// machine without java or with a half-built dss-upstream checkout should behave (plain `go test`
// never requires either), while still actually running and asserting a hard pass here, where all
// three are present.
package pades

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// crossGenFixtures are the "<file>:<expectedLevel>" arguments CrossGenValidator.java takes, one
// per file testdata/crossgen's generator writes: each of the three corpus PDFs of different xref
// styles (see testdata/crossgen/main.go's own doc comment), signed at both PAdES-BASELINE-B and
// PAdES-BASELINE-T.
var crossGenFixtures = []string{
	"EmptyPage-b.pdf:PAdES_BASELINE_B",
	"EmptyPage-t.pdf:PAdES_BASELINE_T",
	"testdoc-b.pdf:PAdES_BASELINE_B",
	"testdoc-t.pdf:PAdES_BASELINE_T",
	"pdf-xref-streams-b.pdf:PAdES_BASELINE_B",
	"pdf-xref-streams-t.pdf:PAdES_BASELINE_T",
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
// like a built upstream DSS checkout with dss-validation and dss-pades-pdfbox (the PDF backend
// CrossGenValidator.java needs to actually read a PAdES signature) compiled. Returns a non-empty
// skip reason when anything is missing.
func detectJavaAndUpstreamDSS(upstreamHome string) string {
	for _, tool := range []string{"java", "javac", "mvn"} {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Sprintf("%s not found on PATH: %v", tool, err)
		}
	}
	for _, module := range crossGenNeededModules {
		classesDir := filepath.Join(upstreamHome, module, "target", "classes")
		if info, err := os.Stat(classesDir); err != nil || !info.IsDir() {
			// dss-pades-pdfbox's own pom.xml declares a test-scope dependency on dss-pdfa, which
			// pulls in org.verapdf:validation-model-jakarta - unavailable offline and not needed
			// here (see crossGenNeededModules's doc comment) - so `mvn install` on it fails in
			// this environment even though a plain `compile` (skipping dependency resolution for
			// the test-compile phase entirely) succeeds.
			buildHint := "`mvn -o -pl dss-validation,dss-pades,dss-cades,dss-cms,dss-cms-object,dss-document -am install -DskipTests`, plus `mvn -o -pl dss-pades-pdfbox compile`,"
			return fmt.Sprintf(
				"upstream DSS checkout not found or %s not built at %s (set DSS_UPSTREAM_HOME, or run %s in it): %v",
				module, classesDir, buildHint, err)
		}
	}
	return ""
}

// crossGenNeededModules are the upstream modules CrossGenValidator.java needs their
// target/classes directory for directly, rather than through dss-validation's own resolved
// dependency jars: dss-pades/dss-pades-pdfbox/dss-cades/dss-cms/dss-cms-object/dss-document are
// reactor modules this checkout builds locally (source of the PDF/CMS/CAdES/generic-signature
// implementation itself) and dss-validation is the module CrossGenValidator.java calls into.
// dss-pades-pdfbox specifically needs a plain `mvn -pl dss-pades-pdfbox compile` (not
// dependency:build-classpath, which also resolves its test-scope dss-pdfa dependency - itself
// unbuildable offline, since it needs org.verapdf:validation-model-jakarta from Maven Central):
// its own compile-scope dependencies (dss-pades, org.apache.pdfbox:pdfbox) are covered by this
// list plus pdfbox's own jars, resolved separately in buildUpstreamDSSClasspath below.
//
// Deliberately NOT a glob over every built module's target/classes: see the identical rationale
// in cades_downstream_cross_validation_test.go (dss-evidence-record-asn1's ServiceLoader
// provider pulling in an unrelated missing dependency).
var crossGenNeededModules = []string{"dss-validation", "dss-pades", "dss-pades-pdfbox", "dss-cades", "dss-cms", "dss-cms-object", "dss-document"}

// crossGenPdfboxArtifacts are the org.apache.pdfbox jars CrossGenValidator.java needs on its
// classpath: dss-pades-pdfbox's own compile-time dependency on org.apache.pdfbox:pdfbox never
// appears in dss-validation's runtime dependency:build-classpath output (dss-validation itself
// has no PDF dependency; see crossGenNeededModules's doc comment on why dss-pades-pdfbox's own
// dependency:build-classpath cannot be used either), so they are found directly in the local
// Maven repository instead. commons-logging is pdfbox's own logging dependency (pulled in
// transitively when built through Maven; resolved the same way here since it is not on
// dss-validation's classpath either).
var crossGenPdfboxArtifacts = []string{"pdfbox", "pdfbox-io", "fontbox"}

// buildUpstreamDSSClasspath assembles the classpath CrossGenValidator.java needs: the modules in
// crossGenNeededModules, the pdfbox jars in crossGenPdfboxArtifacts (plus commons-logging, pdfbox's
// own logging dependency), all found directly under the local Maven repository, plus
// dss-validation's runtime dependency jars (BouncyCastle, the JAXB/XML stack, and every other
// module - dss-spi, dss-model, dss-document, ... - already resolved as local-repository jars
// since this checkout installs them) resolved offline via `mvn -o dependency:build-classpath`.
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
	var extraJars []string
	for _, artifact := range append(append([]string{}, crossGenPdfboxArtifacts...), "commons-logging") {
		jar, err := findNewestJar(m2Repo, artifact)
		if err != nil {
			return "", err
		}
		extraJars = append(extraJars, jar)
	}

	classpathFile := filepath.Join(os.TempDir(), fmt.Sprintf("dss-pades-crossgen-classpath-%d.txt", os.Getpid()))
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

	parts := append(moduleClasses, extraJars...)
	parts = append(parts, strings.TrimSpace(string(dependencyClasspath)))
	return strings.Join(parts, string(os.PathListSeparator)), nil
}

// findNewestJar walks m2Repo for the artifact's directory (org/apache/pdfbox/pdfbox,
// commons-logging/commons-logging, ...) and returns the highest-versioned plain jar (skipping
// -sources.jar/-javadoc.jar) it finds, so this test does not need to hardcode a version this
// checkout's pom.xml might bump.
func findNewestJar(m2Repo, artifact string) (string, error) {
	var groupDir string
	switch artifact {
	case "pdfbox", "pdfbox-io", "fontbox":
		groupDir = filepath.Join(m2Repo, "org", "apache", "pdfbox", artifact)
	case "commons-logging":
		groupDir = filepath.Join(m2Repo, "commons-logging", "commons-logging")
	default:
		return "", fmt.Errorf("findNewestJar: unknown artifact %q", artifact)
	}
	versions, err := os.ReadDir(groupDir)
	if err != nil {
		return "", fmt.Errorf("%s not found under the local Maven repository (%s): %w - run the "+
			"corresponding dependency:build-classpath once so it gets downloaded", artifact, groupDir, err)
	}
	var best, bestJar string
	for _, versionEntry := range versions {
		if !versionEntry.IsDir() {
			continue
		}
		version := versionEntry.Name()
		jar := filepath.Join(groupDir, version, artifact+"-"+version+".jar")
		if _, err := os.Stat(jar); err != nil {
			continue
		}
		if version > best {
			best = version
			bestJar = jar
		}
	}
	if bestJar == "" {
		return "", fmt.Errorf("no usable jar found for %s under %s", artifact, groupDir)
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
