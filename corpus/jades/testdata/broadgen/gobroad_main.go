// Go counterpart of BroadOracle.java: walks the entire upstream dss-jades test-resources corpus
// and dumps the identical JSON shape from this port's own JWSDocumentAnalyzerFactory /
// JAdESSignature. Diff the two files; every difference is a parity defect on one side.
//
// Kept under testdata/ so `go build ./...` ignores it. Run it with the repo module on the path,
// e.g. from a scratch dir with a go.mod that `replace`s github.com/ryftcore/dss-go => <path to your esig checkout>:
//
//	go run . <corpus root> <output json>
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/jades"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// detachedContent mirrors BroadOracle.java's DETACHED_CONTENT map 1:1 (transcribed from the
// upstream test classes that own each fixture). Entries are "relativeFile[|overrideName]", or
// "inline:NAME:CONTENT" for an in-memory document.
var detachedContent = map[string][]string{
	"validation/simple-detached.json":            {"sample.json"},
	"validation/simple-detached-wrong-algo.json": {"sample.json"},
	"validation/jades-detached-by-uri-encoded-pars.json": {
		"ObjectIdByURI-1.html|https://nowina.lu/pub/JAdES/ObjectIdByURI-1.html",
		"ObjectIdByURI-2.html|https://nowina.lu/pub/JAdES/ObjectIdByURI-2.html",
	},
	"validation/jades-detached-by-uri-hash-encoded-pars.json": {
		"ObjectIdByURIHash-1.html|https://signature-plugtests.etsi.org/pub/JAdES/ObjectIdByURIHash-1.html",
		"ObjectIdByURIHash-2.html|https://signature-plugtests.etsi.org/pub/JAdES/ObjectIdByURIHash-2.html",
	},
	"validation/jades-flattened-BpB-detached-objectByURIHash.json": {"inline:TEST-DOC.txt:TL-039 Test Document"},
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: gobroad <corpus root> <output json>")
		os.Exit(2)
	}
	root := os.Args[1]
	var files []string
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(p), ".json") {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	sort.Strings(files)

	var json strings.Builder
	json.WriteString("{\n  \"files\": [\n")
	first := true
	for _, p := range files {
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		one := dumpFileSafe(root, rel)
		if !first {
			json.WriteString(",\n")
		}
		first = false
		json.WriteString(one)
	}
	json.WriteString("\n  ]\n}\n")
	if err := os.WriteFile(os.Args[2], []byte(json.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("written %s (%d files)\n", os.Args[2], len(files))
}

// dumpFileSafe is the Go analogue of BroadOracle.main's try/catch around dumpFile: a panic (this
// port's rendering of an unchecked Java exception) becomes an "error" record, so a file upstream
// refuses and the port accepts (or vice versa) shows up as a diff rather than aborting the run.
func dumpFileSafe(root, rel string) (out string) {
	defer func() {
		if r := recover(); r != nil {
			var chain []string
			if err, ok := r.(error); ok {
				for e := err; e != nil; e = errors.Unwrap(e) {
					chain = append(chain, e.Error())
				}
			} else {
				chain = append(chain, fmt.Sprint(r))
			}
			out = "    {\n      \"path\": " + str(rel) + ",\n      \"error\": " + str(strings.Join(chain, " | ")) + "\n    }"
		}
	}()
	var json strings.Builder
	dumpFile(&json, root, rel)
	return json.String()
}

func dumpFile(json *strings.Builder, root, rel string) {
	document, err := model.NewFileDocument(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		panic(err)
	}
	factory := jades.NewJWSDocumentAnalyzerFactory()
	analyzer := factory.Create(document)
	certificateVerifier := validation.NewCommonCertificateVerifier()
	analyzer.SetCertificateVerifier(certificateVerifier)

	if entries, ok := detachedContent[rel]; ok {
		var detached []model.DSSDocument
		for _, entry := range entries {
			if strings.HasPrefix(entry, "inline:") {
				rest := strings.TrimPrefix(entry, "inline:")
				i := strings.Index(rest, ":")
				d := model.NewInMemoryDocumentWithName([]byte(rest[i+1:]), rest[:i])
				detached = append(detached, d)
				continue
			}
			file, name := entry, ""
			if bar := strings.Index(entry, "|"); bar >= 0 {
				file, name = entry[:bar], entry[bar+1:]
			}
			d, err := model.NewFileDocument(filepath.Join(root, filepath.FromSlash(file)))
			if err != nil {
				panic(err)
			}
			if name != "" {
				d.SetName(name)
			}
			detached = append(detached, d)
		}
		analyzer.SetDetachedContents(detached)
	}

	signatures := analyzer.Signatures()

	json.WriteString("    {\n")
	json.WriteString("      \"path\": " + str(rel) + ",\n")
	fmt.Fprintf(json, "      \"signatureCount\": %d,\n", len(signatures))
	json.WriteString("      \"signatures\": [\n")
	for i, s := range signatures {
		dumpSignature(json, s.(*jades.JAdESSignature), certificateVerifier)
		if i == len(signatures)-1 {
			json.WriteString("\n")
		} else {
			json.WriteString(",\n")
		}
	}
	json.WriteString("      ]\n")
	json.WriteString("    }")
}

func dumpSignature(json *strings.Builder, sig *jades.JAdESSignature, cv validation.CertificateVerifier) {
	sig.InitBaselineRequirementsChecker(cv)

	certificateToken := sig.SigningCertificateToken()
	signingTime := sig.SigningTime()

	json.WriteString("        {\n")
	fmt.Fprintf(json, "          \"signingCertificateFound\": %v,\n", certificateToken != nil)
	if certificateToken != nil {
		d, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, certificateToken.Encoded())
		if err != nil {
			panic(err)
		}
		json.WriteString("          \"signingCertificateSHA256\": " + str(hexs(d)) + ",\n")
	} else {
		json.WriteString("          \"signingCertificateSHA256\": null,\n")
	}
	if signingTime != nil {
		fmt.Fprintf(json, "          \"claimedSigningTimeMillis\": %d,\n", signingTime.UnixMilli())
	} else {
		json.WriteString("          \"claimedSigningTimeMillis\": null,\n")
	}

	level := ""
	safe(func() { level = string(sig.DataFoundUpToLevel()) })
	json.WriteString("          \"dataFoundUpToLevel\": " + str(level) + ",\n")
	fmt.Fprintf(json, "          \"isCounterSignature\": %v,\n", sig.IsCounterSignature())

	var refFound, refIntact, sigIntact bool
	safe(func() {
		v := sig.SignatureCryptographicVerification()
		refFound, refIntact, sigIntact = v.IsReferenceDataFound(), v.IsReferenceDataIntact(), v.IsSignatureIntact()
	})
	fmt.Fprintf(json, "          \"referenceDataFound\": %v,\n", refFound)
	fmt.Fprintf(json, "          \"referenceDataIntact\": %v,\n", refIntact)
	fmt.Fprintf(json, "          \"signatureIntact\": %v,\n", sigIntact)

	dtbsrAlgorithm, dtbsrHex := "", ""
	hasDTBSR := false
	safe(func() {
		d := sig.DataToBeSignedRepresentation()
		if d.Value() != nil {
			hasDTBSR = true
			dtbsrAlgorithm, dtbsrHex = string(d.Algorithm()), hexs(d.Value())
		}
	})
	if hasDTBSR {
		json.WriteString("          \"dtbsrAlgorithm\": " + str(dtbsrAlgorithm) + ",\n")
		json.WriteString("          \"dtbsrHex\": " + str(dtbsrHex) + ",\n")
	} else {
		json.WriteString("          \"dtbsrAlgorithm\": null,\n")
		json.WriteString("          \"dtbsrHex\": null,\n")
	}

	var sigDMechanism *enumerations.SigDMechanism
	safe(func() { sigDMechanism = sig.SigDMechanism() })
	if sigDMechanism != nil {
		json.WriteString("          \"sigDMechanism\": " + str(string(*sigDMechanism)) + ",\n")
	} else {
		json.WriteString("          \"sigDMechanism\": null,\n")
	}

	json.WriteString("          \"serializationType\": " + str(string(sig.Jws().JwsSerializationType())) + ",\n")
	fmt.Fprintf(json, "          \"hasStructureErrors\": %v,\n", len(sig.StructureValidationResult()) > 0)

	signatureAlgorithm := ""
	safe(func() { signatureAlgorithm = string(sig.SignatureAlgorithm()) })
	json.WriteString("          \"signatureAlgorithm\": " + str(signatureAlgorithm) + ",\n")

	// --- certificate extraction ---
	certCount, signingRefCount, certValuesCount, attrRefCount, completeRefCount := -1, -1, -1, -1, -1
	safe(func() { certCount = len(sig.CertificateSource().Certificates()) })
	safe(func() { signingRefCount = len(sig.CertificateSource().SigningCertificateRefs()) })
	safe(func() { certValuesCount = len(sig.CertificateSource().CertificateValues()) })
	safe(func() { attrRefCount = len(sig.CertificateSource().AttributeCertificateRefs()) })
	safe(func() { completeRefCount = len(sig.CertificateSource().CompleteCertificateRefs()) })
	fmt.Fprintf(json, "          \"certificateCount\": %d,\n", certCount)
	fmt.Fprintf(json, "          \"signingCertificateRefCount\": %d,\n", signingRefCount)
	fmt.Fprintf(json, "          \"certificateValuesCount\": %d,\n", certValuesCount)
	fmt.Fprintf(json, "          \"attributeCertificateRefCount\": %d,\n", attrRefCount)
	fmt.Fprintf(json, "          \"completeCertificateRefCount\": %d,\n", completeRefCount)

	// --- revocation extraction ---
	crlBinaries, crlRefs, ocspBinaries, ocspRefs := -1, -1, -1, -1
	safe(func() { crlBinaries = len(sig.CRLSource().AllRevocationBinaries()) })
	safe(func() {
		// AllRevocationReferences is on the concrete OfflineRevocationSourceBase, not on the
		// OfflineRevocationSource[R] interface JAdESSignature's getter is typed as.
		crlRefs = len(sig.CRLSource().(interface {
			AllRevocationReferences() []spi.RevocationRef[revocation.CRL]
		}).AllRevocationReferences())
	})
	safe(func() { ocspBinaries = len(sig.OCSPSource().AllRevocationBinaries()) })
	safe(func() {
		ocspRefs = len(sig.OCSPSource().(interface {
			AllRevocationReferences() []spi.RevocationRef[revocation.OCSP]
		}).AllRevocationReferences())
	})
	fmt.Fprintf(json, "          \"crlBinaryCount\": %d,\n", crlBinaries)
	fmt.Fprintf(json, "          \"crlRefCount\": %d,\n", crlRefs)
	fmt.Fprintf(json, "          \"ocspBinaryCount\": %d,\n", ocspBinaries)
	fmt.Fprintf(json, "          \"ocspRefCount\": %d,\n", ocspRefs)

	// --- signature scopes ---
	json.WriteString("          \"signatureScopes\": [\n")
	scopes := sig.SignatureScopes()
	for i, s := range scopes {
		json.WriteString("            { \"name\": " + str(s.DocumentName()) + ", \"type\": " + str(string(s.Type())) + " }")
		if i == len(scopes)-1 {
			json.WriteString("\n")
		} else {
			json.WriteString(",\n")
		}
	}
	json.WriteString("          ],\n")

	// --- timestamps ---
	dumpTimestamps(json, "contentTimestamps", sig.ContentTimestamps())
	dumpTimestamps(json, "signatureTimestamps", sig.SignatureTimestamps())
	dumpTimestamps(json, "timestampsX1", sig.TimestampsX1())
	dumpTimestamps(json, "timestampsX2", sig.TimestampsX2())
	dumpTimestamps(json, "archiveTimestamps", sig.ArchiveTimestamps())

	counterSignatures := sig.CounterSignatures()
	fmt.Fprintf(json, "          \"counterSignatureCount\": %d,\n", len(counterSignatures))
	json.WriteString("          \"counterSignatures\": [\n")
	for i, cs := range counterSignatures {
		dumpSignature(json, cs.(*jades.JAdESSignature), cv)
		if i == len(counterSignatures)-1 {
			json.WriteString("\n")
		} else {
			json.WriteString(",\n")
		}
	}
	json.WriteString("          ]\n")
	json.WriteString("        }")
}

func dumpTimestamps(json *strings.Builder, field string, timestamps []*validation.TimestampToken) {
	json.WriteString("          " + str(field) + ": [\n")
	for i, t := range timestamps {
		json.WriteString("            { \"type\": " + str(string(t.TimeStampType())))
		dssID := ""
		safe(func() { dssID = t.DSSIDAsString() })
		json.WriteString(", \"dssId\": " + str(dssID))
		gt := t.GenerationTime()
		if gt.IsZero() {
			json.WriteString(", \"generationTimeMillis\": null")
		} else {
			fmt.Fprintf(json, ", \"generationTimeMillis\": %d", gt.UnixMilli())
		}
		var found, intact, siIntact bool
		var imprint string
		hasImprint := false
		safe(func() { found = t.IsMessageImprintDataFound() })
		safe(func() { intact = t.IsMessageImprintDataIntact() })
		safe(func() { siIntact = t.IsSignatureIntact() })
		safe(func() {
			if mi := t.MessageImprint(); mi.Value() != nil {
				hasImprint, imprint = true, hexs(mi.Value())
			}
		})
		fmt.Fprintf(json, ", \"messageImprintDataFound\": %v", found)
		fmt.Fprintf(json, ", \"messageImprintDataIntact\": %v", intact)
		if hasImprint {
			json.WriteString(", \"messageImprintHex\": " + str(imprint))
		} else {
			json.WriteString(", \"messageImprintHex\": null")
		}
		fmt.Fprintf(json, ", \"signatureIntact\": %v", siIntact)
		// The timestamped-reference set is where a missed virtual dispatch on
		// getSignatureTimestampReferences()/getArchiveTimestampReferences() shows up (JAdES's
		// override folds in getKeyInfoReferences()); sorted so the two sides' insertion orders
		// are not compared, only the set.
		var rendered []string
		safe(func() {
			for _, r := range t.TimestampedReferences() {
				rendered = append(rendered, string(r.Category())+":"+r.ObjectId())
			}
		})
		sort.Strings(rendered)
		json.WriteString(", \"timestampedReferences\": " + str(strings.Join(rendered, ",")))
		json.WriteString(" }")
		if i == len(timestamps)-1 {
			json.WriteString("\n")
		} else {
			json.WriteString(",\n")
		}
	}
	json.WriteString("          ],\n")
}

// safe is the Go analogue of BroadOracle.nullSafe: this port panics where upstream throws, and
// both sides record "absent" for the field instead of failing the whole file.
func safe(f func()) {
	defer func() { _ = recover() }()
	f()
}

func hexs(b []byte) string {
	const hd = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hd[v>>4]
		out[i*2+1] = hd[v&0xf]
	}
	return string(out)
}

func str(s string) string {
	var sb strings.Builder
	sb.WriteByte('"')
	for _, c := range []byte(s) {
		switch {
		case c == '"':
			sb.WriteString("\\\"")
		case c == '\\':
			sb.WriteString("\\\\")
		case c == '\n':
			sb.WriteString("\\n")
		case c == '\r':
			sb.WriteString("\\r")
		case c == '\t':
			sb.WriteString("\\t")
		case c < 0x20:
			fmt.Fprintf(&sb, "\\u%04x", c)
		default:
			sb.WriteByte(c)
		}
	}
	sb.WriteByte('"')
	return sb.String()
}
