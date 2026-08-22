// Tests over testdata/adversarial: the CMS documents BouncyCastle writes but no ordinary
// producer does (adv-*.p7s, see gen/Adversarial.java) and the BER pathologies BouncyCastle
// itself will not write (craft-*.p7s, see gen/craft.py). TestBouncyCastleOracle already
// compares every field of these with BouncyCastle's reading; what is checked here is what a
// field-by-field dump cannot express - that the preserved bytes really are the input's own,
// that the signature input keeps every octet it arrived with, and that the documents meant to
// be rejected are.
package cmscore

import (
	"bytes"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// adversarialDir is the corpus these tests read.
const adversarialDir = "adversarial"

// readAdversarial reads one file of the adversarial corpus.
func readAdversarial(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("testdata", adversarialDir, name))
	if err != nil {
		t.Fatalf("cannot read the fixture: %v", err)
	}
	return content
}

// adversarialNames lists the corpus, split into the documents that parse and the ones that
// must not. The craft-bad-* prefix marks the latter.
func adversarialNames(t *testing.T) (accepted []string, rejected []string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join("testdata", adversarialDir))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".p7s") {
			continue
		}
		if strings.HasPrefix(name, "craft-bad-") {
			rejected = append(rejected, name)
		} else {
			accepted = append(accepted, name)
		}
	}
	return accepted, rejected
}

// TestAdversarialDocumentsAreRejected checks the documents whose encoding is broken. Each is
// one BouncyCastle refuses too, which testdata/adversarial/bc-oracle.txt records.
func TestAdversarialDocumentsAreRejected(t *testing.T) {
	_, rejected := adversarialNames(t)
	if len(rejected) == 0 {
		t.Fatal("the corpus holds no craft-bad-* document")
	}
	for _, name := range rejected {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseCMS(readAdversarial(t, name)); err == nil {
				t.Error("the malformed document was accepted")
			}
		})
	}
}

// TestAdversarialPreservedBytesAreSubslices checks that every "as received" accessor still
// hands out a window into the input rather than a re-encoding, whatever the input's length
// forms, tag numbers and component order.
func TestAdversarialPreservedBytesAreSubslices(t *testing.T) {
	accepted, _ := adversarialNames(t)
	for _, name := range accepted {
		t.Run(name, func(t *testing.T) {
			input := readAdversarial(t, name)
			cms, err := ParseCMS(input)
			if err != nil {
				t.Fatalf("ParseCMS: %v", err)
			}
			signedData := cms.SignedData()
			preserved := map[string][]byte{
				"CMS":              cms.Encoded(),
				"SignedData":       signedData.Encoded(),
				"digestAlgorithms": signedData.DigestAlgorithmsElement().Encoded(),
				"encapContentInfo": signedData.EncapContentInfo.Encoded(),
				"certificates":     signedData.Certificates.Encoded(),
				"crls":             signedData.CRLs.Encoded(),
				"signerInfos":      signedData.SignerInfosElement().Encoded(),
			}
			for index, signerInfo := range signedData.SignerInfos {
				preserved["SignerInfo["+itoa(index)+"]"] = signerInfo.Encoded()
				preserved["signedAttrs["+itoa(index)+"]"] = signerInfo.SignedAttributesRaw()
			}
			for index, certificate := range cms.Certificates() {
				preserved["certificate["+itoa(index)+"]"] = certificate
			}
			for index, crl := range cms.CRLs() {
				preserved["CRL["+itoa(index)+"]"] = crl
			}
			for index, response := range cms.OCSPResponses() {
				preserved["OCSPResponse["+itoa(index)+"]"] = response
			}
			for name, preservedBytes := range preserved {
				// An absent optional field has no bytes; a present one must be the input's.
				if len(preservedBytes) == 0 {
					continue
				}
				if !isSubslice(input, preservedBytes) {
					t.Errorf("%s: the bytes handed out are a copy, not the input's own", name)
				}
			}
		})
	}
}

// TestAdversarialSignedAttributesRawReparse checks that the signedAttrs field handed out as
// received holds exactly the attributes the SignerInfo reports, so that a caller which
// re-parses SignedAttributesRaw - as DSS does when it rebuilds a signature - sees the same set.
func TestAdversarialSignedAttributesRawReparse(t *testing.T) {
	accepted, _ := adversarialNames(t)
	for _, name := range accepted {
		t.Run(name, func(t *testing.T) {
			cms, err := ParseCMS(readAdversarial(t, name))
			if err != nil {
				t.Fatalf("ParseCMS: %v", err)
			}
			for index, signerInfo := range cms.SignerInfos() {
				if !signerInfo.HasSignedAttributes() {
					continue
				}
				raw := signerInfo.SignedAttributesRaw()
				element, err := parseOne(raw, "signedAttrs")
				if err != nil {
					t.Fatalf("signer %d: cannot re-parse the raw signedAttrs: %v", index, err)
				}
				attributes, err := attributesFromElement(element, "signedAttrs")
				if err != nil {
					t.Fatalf("signer %d: %v", index, err)
				}
				if len(attributes) != len(signerInfo.SignedAttributes) {
					t.Fatalf("signer %d: the raw field holds %d attributes, the SignerInfo reports %d",
						index, len(attributes), len(signerInfo.SignedAttributes))
				}
				if got, want := attributes.DERSetEncoded(), signerInfo.SignedAttributesDER(); !bytes.Equal(got, want) {
					t.Errorf("signer %d: the re-parsed attributes encode differently: %s",
						index, firstDifference(want, got))
				}
			}
		})
	}
}

// TestSignedAttributesDERKeepsEveryOctet is the regression for a signature input rebuilt from
// the parsed model rather than re-encoded from the received bytes.
//
// RFC 5652 gives Attribute exactly two components, but craft-attribute-third-component.p7s
// carries an attribute with a third one - and the signature was computed over it, since
// SignerInformation#getEncodedSignedAttributes DER-encodes the SET it received rather than
// BouncyCastle's two-field Attribute model. Rebuilding the attribute from attrType and
// attrValues silently dropped that component and produced a different value to verify against.
func TestSignedAttributesDERKeepsEveryOctet(t *testing.T) {
	input := readAdversarial(t, "craft-attribute-third-component.p7s")
	cms, err := ParseCMS(input)
	if err != nil {
		t.Fatalf("ParseCMS: %v", err)
	}
	signerInfo := cms.SignerInfos()[0]
	// The third component: a UTCTime, following the attribute's attrType and attrValues.
	third := []byte{0x17, 0x0D}
	third = append(third, []byte("240102030405Z")...)
	if !bytes.Contains(signerInfo.SignedAttributesRaw(), third) {
		t.Fatal("the fixture no longer carries an attribute with a third component")
	}
	if !bytes.Contains(signerInfo.SignedAttributesDER(), third) {
		t.Error("SignedAttributesDER dropped the third component of the attribute, so the " +
			"value handed to the signature check is not the one that was signed")
	}
	// Every octet of every attribute has to survive, not just that one.
	for _, attribute := range signerInfo.SignedAttributes {
		if got, want := attribute.DER(), attribute.Element().DEREncoded(); !bytes.Equal(got, want) {
			t.Errorf("attribute %s: DER() differs from the DER of the received element: %s",
				attribute.Type, firstDifference(want, got))
		}
	}
}

// TestParseRejectsOversizedSequences checks the component-count guards on the structures this
// package re-encodes field by field. BouncyCastle's ContentInfo and IssuerAndSerialNumber
// constructors reject the same shapes, and accepting them would mean re-encoding a structure
// whose extra component had silently disappeared.
func TestParseRejectsOversizedSequences(t *testing.T) {
	// A stray component appended to a SEQUENCE that has room for none.
	stray := asn1ber.DERNull
	emptySet := derSetOf([]byte{setIdentifier}, nil)

	t.Run("a ContentInfo with a third component", func(t *testing.T) {
		body := asn1ber.EncodeOID(OIDSignedData)
		body = append(body, encodeContextTagged(0, asn1ber.WriteSequence(nil))...)
		body = append(body, stray...)
		if _, err := ParseCMS(asn1ber.WriteSequence(body)); err == nil {
			t.Error("the oversized ContentInfo was accepted")
		}
	})

	t.Run("an EncapsulatedContentInfo with a third component", func(t *testing.T) {
		encap := asn1ber.EncodeOID(OIDData)
		encap = append(encap, encodeContextTagged(0, encodeOctetString(nil))...)
		encap = append(encap, stray...)
		signedData := encodeInt(CMSVersion1)
		signedData = append(signedData, emptySet...)
		signedData = append(signedData, asn1ber.WriteSequence(encap)...)
		signedData = append(signedData, emptySet...)
		document := NewContentInfo(OIDSignedData, asn1ber.WriteSequence(signedData)).DER()
		if _, err := ParseCMS(document); err == nil {
			t.Error("the oversized EncapsulatedContentInfo was accepted")
		}
	})

	t.Run("an IssuerAndSerialNumber with a third component", func(t *testing.T) {
		body := asn1ber.WriteSequence(nil)
		body = append(body, asn1ber.EncodeInteger(big.NewInt(1))...)
		body = append(body, stray...)
		element, err := parseOne(asn1ber.WriteSequence(body), "sid")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := signerIdentifierFromElement(element); err == nil {
			t.Error("the oversized IssuerAndSerialNumber was accepted")
		}
	})
}
