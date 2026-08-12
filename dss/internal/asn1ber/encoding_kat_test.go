// Known answers for the BER-to-DER/DL/BER re-encodings, taken from
// org.bouncycastle.asn1.ASN1Primitive#getEncoded(String) run on the same inputs with
// BouncyCastle 1.84. They pin the three rules the CMS layer depends on and that no amount of
// reading X.690 settles on its own: which length form each encoding keeps, when a constructed
// string collapses, and which content octets DER normalises.
package asn1ber

import (
	"encoding/hex"
	"testing"
)

// encodingKAT is one input and the three encodings BouncyCastle answers for it.
type encodingKAT struct {
	// what the case is there to pin down.
	name string
	// input is the hex of the element as it arrives.
	input string
	// der, dl and ber are BouncyCastle's getEncoded(DER), getEncoded(DL) and getEncoded(BER).
	der, dl, ber string
}

// encodingKATs covers the length-form rules and the DER content normalisations.
//
// The length-form rules BouncyCastle applies, which the first six cases pin down: a
// definite-length constructed element writes its whole subtree with definite lengths, because
// ASN1OutputStream switches to a DL sub-stream as soon as it writes one (DLSequence#encode and
// its siblings); an indefinite-length element keeps its indefinite length and writes its
// components in BER; and a constructed OCTET STRING always comes back indefinite, holding the
// input's own segments when it is the definite-length root and one joined segment otherwise.
var encodingKATs = []encodingKAT{
	{
		name:  "a definite-length constructed OCTET STRING keeps its segments at the root",
		input: "240a04030102030403040506",
		der:   "0406010203040506", dl: "0406010203040506",
		ber: "2480040301020304030405060000",
	},
	{
		name:  "an indefinite-length constructed OCTET STRING joins its segments",
		input: "2480040301020304030405060000",
		der:   "0406010203040506", dl: "0406010203040506",
		ber: "248004060102030405060000",
	},
	{
		name:  "a definite-length SEQUENCE writes its subtree in DL, so the string collapses",
		input: "300c240a04030102030403040506",
		der:   "30080406010203040506", dl: "30080406010203040506",
		ber: "30080406010203040506",
	},
	{
		name:  "an indefinite-length SEQUENCE writes its components in BER",
		input: "3080240a040301020304030405060000",
		der:   "30080406010203040506", dl: "30080406010203040506",
		ber: "30802480040601020304050600000000",
	},
	{
		name:  "an indefinite-length SEQUENCE keeps a nested indefinite length",
		input: "308030800403010203000004030405060000",
		der:   "300c300504030102030403040506", dl: "300c300504030102030403040506",
		ber: "308030800403010203000004030405060000",
	},
	{
		name:  "a definite-length SEQUENCE resolves a nested indefinite length",
		input: "300a30800201010201020000",
		der:   "30083006020101020102", dl: "30083006020101020102",
		ber: "30083006020101020102",
	},
	{
		name:  "DER clears the unused bits of a BIT STRING",
		input: "030204ff",
		der:   "030204f0", dl: "030204ff", ber: "030204ff",
	},
	{
		name:  "DER clears seven unused bits",
		input: "030207ff",
		der:   "03020780", dl: "030207ff", ber: "030207ff",
	},
	{
		name:  "a BIT STRING without unused bits is left alone",
		input: "030300ffff",
		der:   "030300ffff", dl: "030300ffff", ber: "030300ffff",
	},
	{
		name:  "the unused bits of a nested BIT STRING are cleared too",
		input: "3004030204ff",
		der:   "3004030204f0", dl: "3004030204ff", ber: "3004030204ff",
	},
	{
		name:  "a constructed BIT STRING collapses, its unused bits cleared in DER",
		input: "2380030204ff0000",
		der:   "030204f0", dl: "030204ff", ber: "030204ff",
	},
	{
		name:  "DER makes BOOLEAN TRUE all-ones",
		input: "010101",
		der:   "0101ff", dl: "010101", ber: "010101",
	},
	{
		name:  "an already canonical BOOLEAN is left alone",
		input: "0101ff",
		der:   "0101ff", dl: "0101ff", ber: "0101ff",
	},
	{
		name:  "a nested BOOLEAN is normalised in DER only",
		input: "3003010101",
		der:   "30030101ff", dl: "3003010101", ber: "3003010101",
	},
}

// TestEncodingsMatchBouncyCastle checks each re-encoding against BouncyCastle's answer.
func TestEncodingsMatchBouncyCastle(t *testing.T) {
	for _, testCase := range encodingKATs {
		t.Run(testCase.name, func(t *testing.T) {
			input, err := hex.DecodeString(testCase.input)
			if err != nil {
				t.Fatal(err)
			}
			element, rest, err := Parse(input)
			if err != nil {
				t.Fatalf("cannot parse %s: %v", testCase.input, err)
			}
			if len(rest) != 0 {
				t.Fatalf("%d bytes left after %s", len(rest), testCase.input)
			}
			for _, encoding := range []struct {
				name     string
				produced []byte
				expected string
			}{
				{"DER", element.DEREncoded(), testCase.der},
				{"DL", element.DLEncoded(), testCase.dl},
				{"BER", element.BEREncoded(), testCase.ber},
			} {
				if got := hex.EncodeToString(encoding.produced); got != encoding.expected {
					t.Errorf("%s: BouncyCastle answers %s, this package answers %s",
						encoding.name, encoding.expected, got)
				}
			}
		})
	}
}
