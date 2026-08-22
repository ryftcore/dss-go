// Go counterpart of BroadOracle.java: walks a directory of PDFs and dumps the same JSON.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/pades"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

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
		switch c {
		case '"':
			sb.WriteString("\\\"")
		case '\\':
			sb.WriteString("\\\\")
		case '\n':
			sb.WriteString("\\n")
		default:
			sb.WriteByte(c)
		}
	}
	sb.WriteByte('"')
	return sb.String()
}

func main() {
	root := os.Args[1]
	outPath := os.Args[2]
	var files []string
	filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(p), ".pdf") {
			rel, _ := filepath.Rel(root, p)
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	sort.Strings(files)

	var json strings.Builder
	json.WriteString("{\n  \"files\": [\n")
	first := true
	for _, rel := range files {
		one := dumpSafe(root, rel)
		if !first {
			json.WriteString(",\n")
		}
		first = false
		json.WriteString(one)
	}
	json.WriteString("\n  ]\n}\n")
	os.WriteFile(outPath, []byte(json.String()), 0o644)
	fmt.Printf("written %s (%d files)\n", outPath, len(files))
}

func dumpSafe(root, rel string) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = "    {\n      \"path\": " + str(rel) + ",\n      \"error\": " + str(fmt.Sprintf("%v", r)) + "\n    }"
		}
	}()
	var sb strings.Builder
	dumpFile(&sb, root, rel)
	return sb.String()
}

func dumpFile(json *strings.Builder, root, rel string) {
	document, err := model.NewFileDocument(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		panic(err)
	}
	a := pades.NewPDFDocumentAnalyzer(document)
	cv := validation.NewCommonCertificateVerifier()
	a.SetCertificateVerifier(cv)
	a.Validate()
	signatures := a.Signatures()
	detached := a.DetachedTimestamps()

	json.WriteString("    {\n")
	json.WriteString("      \"path\": " + str(rel) + ",\n")
	fmt.Fprintf(json, "      \"signatureCount\": %d,\n", len(signatures))
	json.WriteString("      \"signatures\": [\n")
	for i, s := range signatures {
		dumpSignature(json, s.(*pades.PAdESSignature), cv)
		if i == len(signatures)-1 {
			json.WriteString("\n")
		} else {
			json.WriteString(",\n")
		}
	}
	json.WriteString("      ],\n")
	fmt.Fprintf(json, "      \"detachedTimestampCount\": %d,\n", len(detached))
	json.WriteString("      \"detachedTimestamps\": [\n")
	for i, t := range detached {
		dumpTimestamp(json, t)
		if i == len(detached)-1 {
			json.WriteString("\n")
		} else {
			json.WriteString(",\n")
		}
	}
	json.WriteString("      ]\n")
	json.WriteString("    }")
}

func safeStr(f func() string) (s string) {
	defer func() { recover() }()
	return f()
}

func dumpSignature(json *strings.Builder, sig *pades.PAdESSignature, cv validation.CertificateVerifier) {
	sig.InitBaselineRequirementsChecker(cv)
	signingCertificate := sig.SigningCertificateToken()
	signingTime := sig.SigningTime()

	signerInfoDigest := ""
	signerIDIssuerSerial := "null"
	signerIDSki := "null"

	verification := sig.SignatureCryptographicVerification()
	var messageDigestValue []byte
	func() {
		defer func() { recover() }()
		messageDigestValue = sig.MessageDigestValue()
	}()

	json.WriteString("        {\n")
	fmt.Fprintf(json, "          \"signingCertificateFound\": %v,\n", signingCertificate != nil)
	if signingCertificate != nil {
		d, _ := spi.DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, signingCertificate.Encoded())
		json.WriteString("          \"signingCertificateSHA256\": " + str(hexs(d)) + ",\n")
	} else {
		json.WriteString("          \"signingCertificateSHA256\": null,\n")
	}
	if signingTime != nil {
		fmt.Fprintf(json, "          \"claimedSigningTimeMillis\": %d,\n", signingTime.UnixMilli())
	} else {
		json.WriteString("          \"claimedSigningTimeMillis\": null,\n")
	}
	json.WriteString("          \"dataFoundUpToLevel\": " + str(string(sig.DataFoundUpToLevel())) + ",\n")
	if signerInfoDigest != "" {
		json.WriteString("          \"signerInformationDigestSHA256\": " + str(signerInfoDigest) + ",\n")
	} else {
		json.WriteString("          \"signerInformationDigestSHA256\": null,\n")
	}
	json.WriteString("          \"signerIdIssuerSerial\": " + signerIDIssuerSerial + ",\n")
	json.WriteString("          \"signerIdSubjectKeyIdentifier\": " + signerIDSki + ",\n")
	fmt.Fprintf(json, "          \"isCounterSignature\": %v,\n", sig.IsCounterSignature())
	fmt.Fprintf(json, "          \"referenceDataFound\": %v,\n", verification.IsReferenceDataFound())
	fmt.Fprintf(json, "          \"referenceDataIntact\": %v,\n", verification.IsReferenceDataIntact())
	fmt.Fprintf(json, "          \"signatureIntact\": %v,\n", verification.IsSignatureIntact())
	if messageDigestValue != nil {
		json.WriteString("          \"messageDigestValueHex\": " + str(hexs(messageDigestValue)) + ",\n")
	} else {
		json.WriteString("          \"messageDigestValueHex\": null,\n")
	}
	fmt.Fprintf(json, "          \"counterSignatureCount\": %d,\n", len(sig.CounterSignatures()))

	pdfRevision := sig.PdfRevision()
	sigDict := pdfRevision.PdfSigDictInfo()
	byteRange := pdfRevision.ByteRange()

	json.WriteString("          \"subFilter\": " + str(sigDict.SubFilter()) + ",\n")
	writeOpt(json, "filter", sigDict.Filter())
	writeOpt(json, "signerName", sigDict.SignerName())
	writeOpt(json, "reason", sigDict.Reason())
	writeOpt(json, "location", sigDict.Location())
	writeOpt(json, "docMDP", string(sigDict.DocMDP()))
	fmt.Fprintf(json, "          \"byteRange\": [%d, %d, %d, %d],\n",
		byteRange.FirstPartStart(), byteRange.FirstPartEnd(), byteRange.SecondPartStart(), byteRange.SecondPartEnd())
	fmt.Fprintf(json, "          \"areAllOriginalBytesCovered\": %v,\n", pdfRevision.AreAllOriginalBytesCovered())

	fields := pdfRevision.Fields()
	json.WriteString("          \"fieldNames\": [")
	for i, f := range fields {
		json.WriteString(str(f.FieldName()))
		if i != len(fields)-1 {
			json.WriteString(", ")
		}
	}
	json.WriteString("],\n")

	fmt.Fprintf(json, "          \"documentTimestampCount\": %d,\n", len(sig.DocumentTimestamps()))
	fmt.Fprintf(json, "          \"vriTimestampCount\": %d,\n", len(sig.VRITimestamps()))

	md := pdfRevision.ModificationDetection()
	fmt.Fprintf(json, "          \"pdfModificationsDetected\": %v,\n", md != nil && md.AreModificationsDetected())
	if md != nil {
		om := md.ObjectModifications()
		fmt.Fprintf(json, "          \"secureChangeCount\": %d,\n", len(om.SecureChanges()))
		fmt.Fprintf(json, "          \"formFillChangeCount\": %d,\n", len(om.FormFillInAndSignatureCreationChanges()))
		fmt.Fprintf(json, "          \"annotChangeCount\": %d,\n", len(om.AnnotCreationChanges()))
		fmt.Fprintf(json, "          \"undefinedChangeCount\": %d,\n", len(om.UndefinedChanges()))
	} else {
		json.WriteString("          \"secureChangeCount\": 0,\n          \"formFillChangeCount\": 0,\n          \"annotChangeCount\": 0,\n          \"undefinedChangeCount\": 0,\n")
	}
	json.WriteString("          \"id\": " + str(sig.ID()) + "\n")
	json.WriteString("        }")
}

func writeOpt(json *strings.Builder, name, v string) {
	if v == "" {
		json.WriteString("          \"" + name + "\": null,\n")
	} else {
		json.WriteString("          \"" + name + "\": " + str(v) + ",\n")
	}
}

func dumpTimestamp(json *strings.Builder, t *validation.TimestampToken) {
	json.WriteString("        {\n")
	json.WriteString("          \"type\": " + str(string(t.TimeStampType())) + ",\n")
	fmt.Fprintf(json, "          \"messageImprintDataFound\": %v,\n", t.IsMessageImprintDataFound())
	fmt.Fprintf(json, "          \"messageImprintDataIntact\": %v,\n", t.IsMessageImprintDataIntact())
	fmt.Fprintf(json, "          \"signatureIntact\": %v,\n", t.IsSignatureIntact())
	gt := t.GenerationTime()
	if gt.IsZero() {
		json.WriteString("          \"genTimeMillis\": null\n")
	} else {
		fmt.Fprintf(json, "          \"genTimeMillis\": %d\n", gt.UnixMilli())
	}
	json.WriteString("        }")
}
