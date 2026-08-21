package cms

import (
	"encoding/asn1"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
)

// TestCmsAlgorithmProtectionDER checks cmsAlgorithmProtectionDER against a real BouncyCastle
// CMSAlgorithmProtection encoding (see testdata/gen/GenSignatureAlgorithmIdentifiers.java,
// "CMSAlgorithmProtection tag shape"), for both a bare digest AlgorithmIdentifier and one
// carrying an explicit NULL parameter.
func TestCmsAlgorithmProtectionDER(t *testing.T) {
	sha256 := asn1ber.NewAlgorithmIdentifier(asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1})
	sha256WithRSA := asn1ber.NewAlgorithmIdentifier(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 11})

	hexEqual(t, cmsAlgorithmProtectionDER(sha256, sha256WithRSA),
		"301a300b0609608648016503040201a10b06092a864886f70d01010b")

	sha256WithNull := asn1ber.NewAlgorithmIdentifierWithParameters(
		asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}, asn1ber.DERNull)
	hexEqual(t, cmsAlgorithmProtectionDER(sha256WithNull, sha256WithRSA),
		"301c300d06096086480165030402010500a10b06092a864886f70d01010b")
}

// TestCmsSignedAttributeTableGenerate checks the three injected attributes and the
// override-by-type rule ("entries in it for contentType, signingTime, and messageDigest will
// override the generated ones").
func TestCmsSignedAttributeTableGenerate(t *testing.T) {
	digestAlgorithm := asn1ber.NewAlgorithmIdentifierWithParameters(
		asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}, asn1ber.DERNull)
	signatureAlgorithm := asn1ber.NewAlgorithmIdentifier(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 11})
	messageDigest := []byte{1, 2, 3, 4}
	contentType := cmscore.OIDData

	t.Run("empty input gets all three", func(t *testing.T) {
		attributes := cmsSignedAttributeTableGenerate(nil, contentType, messageDigest, digestAlgorithm, signatureAlgorithm)
		if len(attributes) != 3 {
			t.Fatalf("got %d attributes, want 3", len(attributes))
		}
		if attributes.Get(cmscore.OIDContentType) == nil {
			t.Error("missing content-type attribute")
		}
		if attributes.Get(cmscore.OIDMessageDigest) == nil {
			t.Error("missing message-digest attribute")
		}
		if attributes.Get(OID_id_aa_cmsAlgorithmProtect) == nil {
			t.Error("missing cms-algorithm-protection attribute")
		}
	})

	t.Run("nil contentType omits content-type (counter-signature)", func(t *testing.T) {
		attributes := cmsSignedAttributeTableGenerate(nil, nil, messageDigest, digestAlgorithm, signatureAlgorithm)
		if len(attributes) != 2 {
			t.Fatalf("got %d attributes, want 2", len(attributes))
		}
		if attributes.Get(cmscore.OIDContentType) != nil {
			t.Error("content-type attribute should be absent for a nil contentType")
		}
	})

	t.Run("existing message-digest is not overridden", func(t *testing.T) {
		existing := cmscore.Attributes{cmscore.NewAttribute(cmscore.OIDMessageDigest, asn1ber.WriteTLV(asn1ber.TagOctetString, []byte{9, 9}))}
		attributes := cmsSignedAttributeTableGenerate(existing, contentType, messageDigest, digestAlgorithm, signatureAlgorithm)
		if len(attributes) != 3 {
			t.Fatalf("got %d attributes, want 3", len(attributes))
		}
		found := attributes.Get(cmscore.OIDMessageDigest)
		if found == nil {
			t.Fatal("missing message-digest attribute")
		}
		if string(found.ValueEncodings()[0]) != string(existing[0].ValueEncodings()[0]) {
			t.Error("the pre-existing message-digest attribute was overridden")
		}
	})
}
