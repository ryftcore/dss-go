// Known-answer tests for the CAdES archive time-stamp: the ats-hash-index attribute of all
// three versions and the archive-time-stamp-v3 message imprint have to be byte-identical to
// upstream DSS's, because a TSA signs those bytes and a verifier re-computes them.
//
// Every expected value comes from running upstream DSS 6.5.RC1 itself over the documents in
// testdata/upstream/ (copied from dss-cades/src/test/resources/validation/), and from running
// BouncyCastle 1.84's AttributeTable over synthetic attribute sets; see
// testdata/gen/AtsHashIndexOracle.java for how to regenerate both oracles. Nothing here is
// hand-derived.
package cades

import (
	"bufio"
	"encoding/asn1"
	"encoding/hex"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/utain/esig/dss/cms"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/internal/corpustest"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
)

// TestCadesLevelBaselineLTATimestampExtractorOracle replays testdata/ats-hash-index-oracle.txt.
func TestCadesLevelBaselineLTATimestampExtractorOracle(t *testing.T) {
	file, err := os.Open(corpustest.Path(t, "ats-hash-index-oracle.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()

	var (
		document   *cms.CMS
		signerInfo *cmscore.SignerInfo
		extractor  *CadesLevelBaselineLTATimestampExtractor
		original   model.DSSDocument
		name       string
		checked    int
	)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		switch fields[0] {
		case "FILE":
			name = fields[1]
			raw, err := os.ReadFile(cadesFixturePath(t, filepath.Join("upstream", name)))
			if err != nil {
				t.Fatal(err)
			}
			document, err = cms.CMSUtilsParseToCMSBinaries(raw)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			signerInfo, extractor, original = nil, nil, nil

		case "SIG":
			// The SIGNERINFO line that follows selects the signer.

		case "SIGNERINFO":
			signerInfo = nil
			for _, candidate := range document.SignerInfos() {
				digest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithm_SHA256, candidate.DER())
				if err != nil {
					t.Fatal(err)
				}
				if hex.EncodeToString(digest) == fields[1] {
					signerInfo = candidate
				}
			}
			if signerInfo == nil {
				t.Fatalf("%s: no SignerInfo digesting to %s", name, fields[1])
			}
			signature, err := newCadesLTAOracleSignature(document, signerInfo)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			extractor = newCadesLevelBaselineLTATimestampExtractor(signature)

		case "ECONTENTTYPE":
			if got := hex.EncodeToString(extractor.encodedContentType()); got != fields[1] {
				t.Errorf("%s: eContentType = %s, want %s", name, got, fields[1])
			}
			checked++

		case "ORIGINALDOCUMENT":
			original = nil
			if fields[1] == "none" {
				continue
			}
			raw, err := hex.DecodeString(fields[1])
			if err != nil {
				t.Fatal(err)
			}
			original = model.NewInMemoryDocument(raw)

		case "ATS":
			attribute := cadesLTAOracleAtsHashIndex(t, extractor, signerInfo, fields[1], fields[2])
			if got := hex.EncodeToString(attribute.DER()); got != fields[3] {
				t.Errorf("%s ats-hash-index %s/%s:\n got  %s\n want %s", name, fields[1], fields[2], got, fields[3])
			}
			checked++

		case "IMPRINT":
			if original == nil {
				t.Fatalf("%s: an IMPRINT needs the original document", name)
			}
			attribute := cadesLTAOracleAtsHashIndex(t, extractor, signerInfo, fields[1], fields[2])
			imprint, err := extractor.ArchiveTimestampV3MessageImprint(
				signerInfo, attribute, original, enumerations.DigestAlgorithm(fields[2]))
			if err != nil {
				t.Fatalf("%s: ArchiveTimestampV3MessageImprint: %v", name, err)
			}
			if got := hex.EncodeToString(imprint.Value()); got != fields[3] {
				t.Errorf("%s message imprint %s/%s:\n got  %s\n want %s", name, fields[1], fields[2], got, fields[3])
			}
			checked++

		case "VERIFIED":
			token := cadesLTAOracleArchiveTimestamp(t, signerInfo, fields[1])
			attribute, err := extractor.VerifiedAtsHashIndex(signerInfo, token)
			if err != nil {
				t.Fatalf("%s: VerifiedAtsHashIndex: %v", name, err)
			}
			if got := hex.EncodeToString(attribute.DER()); got != fields[2] {
				t.Errorf("%s verified ats-hash-index %s:\n got  %s\n want %s", name, fields[1], got, fields[2])
			}
			checked++

		case "VERIFIEDTHROWS":
			// Upstream raises the recorded exception; this port returns an error instead.
			token := cadesLTAOracleArchiveTimestamp(t, signerInfo, fields[1])
			if _, err := extractor.VerifiedAtsHashIndex(signerInfo, token); err == nil {
				t.Errorf("%s: VerifiedAtsHashIndex succeeded, want the error standing for %s", name, fields[2])
			}
			checked++

		case "VERIFIEDSTATUS":
			token := cadesLTAOracleArchiveTimestamp(t, signerInfo, fields[1])
			// The status is set on the token even when the rebuild fails, so the error is ignored.
			_, _ = extractor.VerifiedAtsHashIndex(signerInfo, token)
			status := token.AtsHashIndexStatus()
			if status == nil {
				t.Fatalf("%s: no ats-hash-index status was set", name)
			}
			version := string(status.Version())
			if version == "" {
				version = "none"
			}
			if version != fields[2] {
				t.Errorf("%s status version = %s, want %s", name, version, fields[2])
			}
			wantMessages := strings.TrimSpace(strings.Join(fields[3:], " "))
			gotMessages := strings.Join(status.ErrorMessages(), " | ")
			if gotMessages != wantMessages {
				t.Errorf("%s status messages:\n got  %q\n want %q", name, gotMessages, wantMessages)
			}
			checked++

		default:
			t.Fatalf("unknown oracle record %q", fields[0])
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if checked == 0 {
		t.Fatal("the oracle yielded no assertion")
	}
	t.Logf("%d oracle assertions checked", checked)
}

// TestCadesLTAAttributeTableOrder replays testdata/attribute-table-order-oracle.txt, pinning the
// AttributeTable#toASN1EncodableVector() ordering the unsignedAttrsHashIndex is built with.
func TestCadesLTAAttributeTableOrder(t *testing.T) {
	file, err := os.Open(cadesFixturePath(t, "attribute-table-order-oracle.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()

	var input cmscore.Attributes
	cases := 0
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		switch fields[0] {
		case "IN":
			input = nil
			for index, attrType := range fields[1:] {
				// The value mirrors the oracle's: an INTEGER holding the input position, so
				// that repeated attribute types stay distinguishable.
				value := asn1ber.EncodeInteger(cadesLTABigInt(index))
				input = append(input, cmscore.NewAttribute(cadesLTAOID(t, attrType), value))
			}
		case "OUT":
			ordered := cadesLTAAttributeTableOrder(input)
			if len(ordered) != len(fields)-1 {
				t.Fatalf("ordering returned %d attributes, want %d", len(ordered), len(fields)-1)
			}
			for index, want := range fields[1:] {
				attrType, position, _ := strings.Cut(want, "#")
				got := ordered[index]
				if got.Type.String() != attrType {
					t.Errorf("position %d: attrType = %s, want %s", index, got.Type, attrType)
				}
				value := got.ValueEncodings()[0]
				element, _, err := asn1ber.Parse(value)
				if err != nil {
					t.Fatal(err)
				}
				if element.Integer().String() != position {
					t.Errorf("position %d: value = %s, want %s", index, element.Integer(), position)
				}
			}
			cases++
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if cases == 0 {
		t.Fatal("the oracle yielded no case")
	}
}

// -----------------------------------------------------------------------------
// The signature the KAT drives the extractor with: it answers the seven accessors
// CadesLevelBaselineLTATimestampExtractor reads, straight off the parsed CMS, so that the
// byte-exact core is exercised without the document analyzer.
// -----------------------------------------------------------------------------

type cadesLTAOracleSignature struct {
	certificateSource *spi.CMSCertificateSource
	completeCerts     *spi.ListCertificateSource
	crlSource         *spi.CMSCRLSource
	ocspSource        *spi.CMSOCSPSource
	completeCRLs      *spi.ListRevocationSource[revocation.CRL]
	completeOCSPs     *spi.ListRevocationSource[revocation.OCSP]
	cms               *cms.CMS
}

func newCadesLTAOracleSignature(document *cms.CMS, signerInfo *cmscore.SignerInfo) (*cadesLTAOracleSignature, error) {
	certificateSource, err := spi.NewCMSCertificateSource(document.SignerInfos(), document.Certificates(), signerInfo)
	if err != nil {
		return nil, err
	}
	crlSource, err := spi.NewCMSCRLSource(document.CRLs(), signerInfo.UnsignedAttributes)
	if err != nil {
		return nil, err
	}
	ocspSource, err := spi.NewCMSOCSPSource(document.OcspResponseStore(), document.OcspBasicStore(),
		signerInfo.UnsignedAttributes)
	if err != nil {
		return nil, err
	}
	return &cadesLTAOracleSignature{
		certificateSource: certificateSource,
		completeCerts:     spi.NewListCertificateSourceFromOne(&certificateSource.SignatureCertificateSource),
		crlSource:         crlSource,
		ocspSource:        ocspSource,
		completeCRLs:      spi.NewListRevocationSourceFrom[revocation.CRL](crlSource),
		completeOCSPs:     spi.NewListRevocationSourceFrom[revocation.OCSP](ocspSource),
		cms:               document,
	}, nil
}

func (s *cadesLTAOracleSignature) CertificateSource() *spi.SignatureCertificateSource {
	return &s.certificateSource.SignatureCertificateSource
}
func (s *cadesLTAOracleSignature) CompleteCertificateSource() *spi.ListCertificateSource {
	return s.completeCerts
}
func (s *cadesLTAOracleSignature) CRLSource() spi.OfflineRevocationSource[revocation.CRL] {
	return s.crlSource
}
func (s *cadesLTAOracleSignature) OCSPSource() spi.OfflineRevocationSource[revocation.OCSP] {
	return s.ocspSource
}
func (s *cadesLTAOracleSignature) CompleteCRLSource() *spi.ListRevocationSource[revocation.CRL] {
	return s.completeCRLs
}
func (s *cadesLTAOracleSignature) CompleteOCSPSource() *spi.ListRevocationSource[revocation.OCSP] {
	return s.completeOCSPs
}
func (s *cadesLTAOracleSignature) CMS() *cms.CMS { return s.cms }

// cadesLTAOracleArchiveTimestamp builds the TimestampToken of the archive time-stamp of the
// signer whose encoding digests to id. The oracle identifies a time-stamp that way rather than
// by its position, which depends on how AdvancedSignature#getArchiveTimestamps() orders them.
func cadesLTAOracleArchiveTimestamp(t *testing.T, signerInfo *cmscore.SignerInfo, id string) *validation.TimestampToken {
	t.Helper()
	for _, attribute := range signerInfo.UnsignedAttributes {
		if !spi.OID_id_aa_ets_archiveTimestampV2.Equal(attribute.Type) &&
			!spi.OID_id_aa_ets_archiveTimestampV3.Equal(attribute.Type) {
			continue
		}
		for _, value := range attribute.Values {
			digest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithm_SHA256, value.Encoded())
			if err != nil {
				t.Fatal(err)
			}
			if hex.EncodeToString(digest) != id {
				continue
			}
			token, err := validation.NewTimestampToken(value.Encoded(),
				enumerations.TimestampType_ARCHIVE_TIMESTAMP)
			if err != nil {
				t.Fatal(err)
			}
			return token
		}
	}
	t.Fatalf("no archive time-stamp digesting to %s", id)
	return nil
}

func cadesLTAOracleAtsHashIndex(t *testing.T, extractor *CadesLevelBaselineLTATimestampExtractor,
	signerInfo *cmscore.SignerInfo, version, digestAlgorithm string) *cmscore.Attribute {
	t.Helper()
	attribute, err := extractor.AtsHashIndex(signerInfo,
		enumerations.DigestAlgorithm(digestAlgorithm), cadesLTAOID(t, version))
	if err != nil {
		t.Fatalf("AtsHashIndex(%s, %s): %v", version, digestAlgorithm, err)
	}
	return attribute
}

func cadesLTAAtoi(t *testing.T, value string) int {
	t.Helper()
	number, err := strconv.Atoi(value)
	if err != nil {
		t.Fatal(err)
	}
	return number
}

// cadesLTAOID parses a dotted OID of the oracle.
func cadesLTAOID(t *testing.T, value string) asn1.ObjectIdentifier {
	t.Helper()
	oid, err := asn1ber.OIDFromString(value)
	if err != nil {
		t.Fatal(err)
	}
	return oid
}

// cadesLTABigInt is the INTEGER value the order oracle gives each synthetic attribute.
func cadesLTABigInt(value int) *big.Int { return big.NewInt(int64(value)) }
