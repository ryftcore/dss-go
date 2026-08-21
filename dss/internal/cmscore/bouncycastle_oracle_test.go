// The differential test against BouncyCastle. Every field this package exposes is dumped as
// "key=value" lines and compared with testdata/bc-oracle.txt and
// testdata/adversarial/bc-oracle.txt, which testdata/gen/CmsOracle.java produced by reading
// the very same files through org.bouncycastle.cms.CMSSignedData, org.bouncycastle.asn1.cms.*
// and org.bouncycastle.tsp.*. The expectations are therefore BouncyCastle's answers rather
// than this package's own, and a divergence in parsing, in a preserved encoding, in the
// re-encodings CMSObjectUtils writes for an archive time-stamp, or in the RFC 5652 clause 5.1
// version computation fails the test.
package cmscore

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/corpustest"
)

// oracleGoldenFile is the file each corpus keeps its BouncyCastle answers in.
const oracleGoldenFile = "bc-oracle.txt"

// TestBouncyCastleOracle compares the Go reading of both fixture corpora with BouncyCastle's.
//
// Set CMSCORE_REGENERATE=1 to rewrite the golden files from the Go side; do that only after
// checking the result against a fresh run of testdata/gen/CmsOracle.java, since the point of
// the files is that an independent implementation produced them.
func TestBouncyCastleOracle(t *testing.T) {
	for _, rel := range []string{"", "adversarial"} {
		dir := filepath.Join("testdata", rel)
		t.Run(dir, func(t *testing.T) {
			produced := dumpDirectory(t, dir)
			golden := corpustest.Path(t, filepath.Join(rel, oracleGoldenFile))
			if os.Getenv("CMSCORE_REGENERATE") != "" {
				if err := os.WriteFile(golden, []byte(produced), 0o644); err != nil {
					t.Fatal(err)
				}
				t.Skip("the golden file was regenerated")
			}
			expected, err := os.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}
			compareOracleDumps(t, string(expected), produced)
		})
	}
}

// compareOracleDumps reports every line on which the two dumps disagree.
func compareOracleDumps(t *testing.T, expected, produced string) {
	t.Helper()
	expectedLines := oracleLines(expected)
	producedLines := oracleLines(produced)
	keys := make([]string, 0, len(expectedLines))
	seen := map[string]bool{}
	for key := range expectedLines {
		keys = append(keys, key)
		seen[key] = true
	}
	for key := range producedLines {
		if !seen[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		expectedValue, hasExpected := expectedLines[key]
		producedValue, hasProduced := producedLines[key]
		switch {
		case !hasExpected:
			t.Errorf("%s: this package answers %s, BouncyCastle answers nothing", key, producedValue)
		case !hasProduced:
			t.Errorf("%s: BouncyCastle answers %s, this package answers nothing", key, expectedValue)
		case expectedValue != producedValue:
			t.Errorf("%s:\n  BouncyCastle %s\n  cmscore      %s", key, expectedValue, producedValue)
		}
	}
}

// oracleLines splits a dump into its key/value pairs.
func oracleLines(dump string) map[string]string {
	lines := map[string]string{}
	for _, line := range strings.Split(dump, "\n") {
		if key, value, found := strings.Cut(line, "="); found {
			lines[key] = value
		}
	}
	return lines
}

// dumpDirectory reads every fixture of a directory and returns the dump.
func dumpDirectory(t *testing.T, dir string) string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasSuffix(entry.Name(), ".txt") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	var out strings.Builder
	writer := &oracleWriter{out: &out}
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		switch {
		case strings.HasSuffix(name, ".p7s"):
			writer.put(name+".parseOK", writer.dumpCMS(name, data))
		case strings.HasSuffix(name, ".tst"):
			token, err := ParseTimeStampToken(data)
			if err != nil {
				writer.put(name+".parseOK", false)
				continue
			}
			writer.dumpTimeStampToken(name, token)
			writer.dumpCMS(name+".cms", data)
			writer.put(name+".parseOK", true)
		case strings.HasSuffix(name, ".tsr"):
			writer.put(name+".parseOK", writer.dumpTimeStampResp(name, data))
		}
	}
	return out.String()
}

// oracleWriter accumulates the dump lines.
type oracleWriter struct{ out *strings.Builder }

// put writes one key/value line.
func (w *oracleWriter) put(key string, value any) {
	fmt.Fprintf(w.out, "%s=%v\n", key, value)
}

// putBytes writes an octet string, pinning anything longer than 64 octets by its SHA-256 so
// that the golden files stay readable. CmsOracle.java abbreviates the same way.
func (w *oracleWriter) putBytes(key string, value []byte) {
	switch {
	case value == nil:
		w.put(key, "<nil>")
	case len(value) <= 64:
		w.put(key, hex(value))
	default:
		digest := sha256.Sum256(value)
		w.put(key, "sha256:"+hex(digest[:]))
	}
}

// algorithmDump renders an AlgorithmIdentifier, distinguishing absent parameters from a NULL.
func algorithmDump(algorithm *asn1ber.AlgorithmIdentifier) string {
	if algorithm == nil {
		return "<nil>"
	}
	if algorithm.Parameters == nil {
		return algorithm.Algorithm.String() + "/<absent>"
	}
	return algorithm.Algorithm.String() + "/" + hex(algorithm.Parameters)
}

// dumpCMS writes every field the CMS accessors expose, plus the encodings CMSObjectUtils
// writes for an archive time-stamp message imprint, and reports whether the document parsed.
// A document that did not parse contributes no lines at all, as the Java oracle's rollback of
// a failed dump does.
func (w *oracleWriter) dumpCMS(name string, data []byte) bool {
	cms, err := ParseCMS(data)
	if err != nil {
		return false
	}
	signedData := cms.SignedData()
	w.put(name+".contentType", cms.ContentInfo().ContentType)
	w.put(name+".definiteLength", cms.IsDefiniteLength())
	w.put(name+".version", cms.Version())
	for index, algorithm := range cms.DigestAlgorithmIDs() {
		w.put(fmt.Sprintf("%s.digestAlgorithms[%d]", name, index), algorithmDump(algorithm))
	}
	w.put(name+".eContentType", cms.SignedContentType())
	w.put(name+".detached", cms.IsDetachedSignature())
	w.putBytes(name+".signedContent", cms.SignedContent())

	w.put(name+".certificatesPresent", signedData.Certificates != nil)
	certificateIndex, attributeCertificateIndex := 0, 0
	if signedData.Certificates != nil {
		for index, choice := range signedData.Certificates.Choices {
			w.put(fmt.Sprintf("%s.certChoiceTag[%d]", name, index), choice.TagNo)
			switch choice.TagNo {
			case CertificateChoiceCertificate:
				w.putBytes(fmt.Sprintf("%s.cert[%d]", name, certificateIndex), choice.DER())
				certificateIndex++
			case CertificateChoiceV2AttrCert:
				w.putBytes(fmt.Sprintf("%s.attrCert[%d]", name, attributeCertificateIndex), choice.DER())
				attributeCertificateIndex++
			}
		}
	}
	w.put(name+".certCount", len(cms.Certificates()))
	w.put(name+".attrCertCount", len(cms.AttributeCertificates()))

	w.put(name+".crlsPresent", signedData.CRLs != nil)
	crlIndex, responseIndex, basicIndex := 0, 0, 0
	if signedData.CRLs != nil {
		for index, choice := range signedData.CRLs.Choices {
			if choice.Other == nil {
				w.putBytes(fmt.Sprintf("%s.crl[%d]", name, crlIndex), choice.CRL)
				crlIndex++
				continue
			}
			w.put(fmt.Sprintf("%s.otherRevFormat[%d]", name, index), choice.Other.Format)
			switch {
			case choice.Other.Format.Equal(OIDRIOCSPResponse):
				w.putBytes(fmt.Sprintf("%s.ocspResponse[%d]", name, responseIndex), choice.Other.Info)
				responseIndex++
			case choice.Other.Format.Equal(OIDPKIXOCSPBasic):
				w.putBytes(fmt.Sprintf("%s.ocspBasic[%d]", name, basicIndex), choice.Other.Info)
				basicIndex++
			}
		}
	}
	w.put(name+".crlCount", len(cms.CRLs()))
	w.put(name+".ocspResponseCount", len(cms.OCSPResponses()))
	w.put(name+".ocspBasicCount", len(cms.OCSPBasicResponses()))

	w.put(name+".signerCount", len(cms.SignerInfos()))
	for index, signerInfo := range cms.SignerInfos() {
		w.dumpSignerInfo(fmt.Sprintf("%s.signer[%d]", name, index), signerInfo)
	}

	w.putBytes(name+".derEncoded", cms.DEREncoded())
	w.putBytes(name+".dlEncoded", cms.ContentInfo().Element().DLEncoded())
	w.putBytes(name+".berEncoded", cms.ContentInfo().Element().BEREncoded())
	w.dumpArchiveTimeStampEncodings(name, signedData)
	w.put(name+".recomputedVersion", ComputeSignedDataVersion(signedData.EncapContentInfo.EContentType,
		signedData.Certificates, signedData.CRLs, signedData.SignerInfos))
	return true
}

// dumpSignerInfo writes one SignerInfo.
func (w *oracleWriter) dumpSignerInfo(key string, signerInfo *SignerInfo) {
	w.put(key+".version", signerInfo.Version)
	w.put(key+".sidIsSKI", signerInfo.SID.IsSubjectKeyIdentifier())
	if signerInfo.SID.IsSubjectKeyIdentifier() {
		w.putBytes(key+".ski", signerInfo.SID.SubjectKeyIdentifier)
	} else {
		w.putBytes(key+".issuer", signerInfo.SID.IssuerAndSerialNumber.Issuer)
		w.put(key+".serial", signerInfo.SID.IssuerAndSerialNumber.SerialNumber)
	}
	w.put(key+".digestAlgorithm", algorithmDump(signerInfo.DigestAlgorithm))
	w.put(key+".signatureAlgorithm", algorithmDump(signerInfo.SignatureAlgorithm))
	w.putBytes(key+".signature", signerInfo.Signature)
	w.put(key+".hasSignedAttrs", signerInfo.HasSignedAttributes())
	if signerInfo.HasSignedAttributes() {
		w.putBytes(key+".signedAttrsDER", signerInfo.SignedAttributesDER())
		w.put(key+".signedAttrCount", len(signerInfo.SignedAttributes))
		for index, attribute := range signerInfo.SignedAttributes {
			w.put(fmt.Sprintf("%s.signedAttr[%d].type", key, index), attribute.Type)
			for valueIndex, value := range attribute.Values {
				w.putBytes(fmt.Sprintf("%s.signedAttr[%d].value[%d]", key, index, valueIndex),
					value.DEREncoded())
			}
		}
	}
	w.put(key+".hasUnsignedAttrs", signerInfo.HasUnsignedAttributes())
	if signerInfo.HasUnsignedAttributes() {
		w.putBytes(key+".unsignedAttrsDER", signerInfo.UnsignedAttributes.DERSetEncoded())
	}
}

// dumpArchiveTimeStampEncodings writes the five subtree encodings CMSObjectUtils streams into
// an ATSv3 message imprint: the digestAlgorithms SET in the length form it arrived in, the
// encapContentInfo in BER when its eContent is a BER OCTET STRING and in DER otherwise, and
// the certificates, crls and signerInfos fields in BER or DER by the same rule.
func (w *oracleWriter) dumpArchiveTimeStampEncodings(name string, signedData *SignedData) {
	w.putBytes(name+".ats.digestAlgorithms", signedData.DigestAlgorithmsElement().BEREncoded())
	encapContentInfo := signedData.EncapContentInfo
	if content := encapContentInfo.ContentElement(); content != nil && content.IsIndefinite() {
		w.putBytes(name+".ats.contentInfo", encapContentInfo.Element().BEREncoded())
	} else {
		w.putBytes(name+".ats.contentInfo", encapContentInfo.DER())
	}
	if element := signedData.CertificatesElement(); element != nil {
		w.putBytes(name+".ats.certificates", archiveTimeStampEncoded(element))
	}
	if element := signedData.CRLsElement(); element != nil {
		w.putBytes(name+".ats.crls", archiveTimeStampEncoded(element))
	}
	w.putBytes(name+".ats.signerInfos", archiveTimeStampEncoded(signedData.SignerInfosElement()))
}

// archiveTimeStampEncoded applies the "BER when the field arrived indefinite, DER otherwise"
// rule CMSObjectUtils implements with its BERSet check.
func archiveTimeStampEncoded(element *asn1ber.Element) []byte {
	if element.IsIndefinite() {
		return element.BEREncoded()
	}
	return element.DEREncoded()
}

// dumpTimeStampResp writes an RFC 3161 response and reports whether it parsed.
func (w *oracleWriter) dumpTimeStampResp(name string, data []byte) bool {
	response, err := ParseTimeStampResp(data)
	if err != nil {
		return false
	}
	w.put(name+".status", response.Status.Status)
	if response.Status.StatusString == nil {
		w.put(name+".statusString", "<nil>")
	} else {
		w.put(name+".statusString", "["+strings.Join(response.Status.StatusString, ", ")+"]")
	}
	w.put(name+".failInfoPresent", response.Status.FailInfo != nil)
	if response.Status.FailInfo != nil {
		w.putBytes(name+".failInfoBits", response.Status.FailInfo)
		w.put(name+".failInfoPadBits", response.Status.FailInfoUnusedBits)
	}
	w.put(name+".tokenPresent", response.Token != nil)
	if response.Token != nil {
		w.dumpTimeStampToken(name, response.Token)
		w.dumpCMS(name+".cms", response.Token.Encoded())
	}
	return true
}

// dumpTimeStampToken writes the TSTInfo of a token.
func (w *oracleWriter) dumpTimeStampToken(name string, token *TimeStampToken) {
	info := token.TSTInfo()
	w.put(name+".tst.version", info.Version)
	w.put(name+".tst.policy", info.Policy)
	w.put(name+".tst.serial", info.SerialNumber)
	w.put(name+".tst.genTimeMillis", info.GenTime.UnixMilli())
	w.put(name+".tst.genTimeString", info.GenTimeString)
	w.put(name+".tst.imprintAlg", algorithmDump(info.MessageImprint.HashAlgorithm))
	w.putBytes(name+".tst.imprintDigest", info.MessageImprint.HashedMessage)
	w.put(name+".tst.accuracyPresent", info.Accuracy != nil)
	if info.Accuracy != nil {
		w.put(name+".tst.accSeconds", optionalInt(info.Accuracy.Seconds))
		w.put(name+".tst.accMillis", optionalInt(info.Accuracy.Millis))
		w.put(name+".tst.accMicros", optionalInt(info.Accuracy.Micros))
	}
	w.put(name+".tst.ordering", info.Ordering)
	if info.Nonce == nil {
		w.put(name+".tst.nonce", "<nil>")
	} else {
		w.put(name+".tst.nonce", info.Nonce)
	}
	w.put(name+".tst.tsaPresent", info.TSA != nil)
	if info.TSA != nil {
		w.put(name+".tst.tsaTag", info.TSA.TagNo)
		w.putBytes(name+".tst.tsaDER", info.TSA.DER())
	}
	w.put(name+".tst.extensionsPresent", info.Extensions != nil)
	w.putBytes(name+".tst.encoded", token.Encoded())
}

// optionalInt renders an optional INTEGER component.
func optionalInt(value *int) any {
	if value == nil {
		return "<nil>"
	}
	return *value
}
