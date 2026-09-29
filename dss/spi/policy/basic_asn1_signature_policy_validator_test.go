package policy

import (
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/signature"
)

// The tests below pin BasicASN1SignaturePolicyValidator.validate against the
// Java control flow: everything that is not a fully conforming TR 102 272
// SEQUENCE ends in the catch (Exception) block, which marks the digest
// invalid (T23-SEC-001: a SEQUENCE with fewer than three children used to be
// reported as digestValid=true).

var sha256OID = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}

// policyFixture builds the pieces of a TR 102 272 policy: the SHA-256
// AlgorithmIdentifier, a SignPolicyInfo stand-in and the SignPolicyHash
// computed the way spi.DSSASN1UtilsAsn1SignaturePolicyDigest does
// (digest over DER(alg) || DER(policyInfo)).
func policyFixture() (algID, policyInfo, hash []byte) {
	algID = asn1ber.WriteSequence(asn1ber.EncodeOID(sha256OID))
	policyInfo = asn1ber.WriteSequence(asn1ber.EncodeInteger(big.NewInt(1)))
	sum := sha256.Sum256(append(append([]byte{}, algID...), policyInfo...))
	return algID, policyInfo, sum[:]
}

func concat(parts ...[]byte) []byte {
	var out []byte
	for _, part := range parts {
		out = append(out, part...)
	}
	return out
}

func validate(t *testing.T, policyBytes []byte, digest model.Digest) *signature.PolicyValidationResult {
	t.Helper()
	policy := signature.NewPolicy()
	policy.SetPolicyContent(model.NewInMemoryDocument(policyBytes))
	policy.SetDigest(digest)
	validator := NewBasicASN1SignaturePolicyValidator()
	if !validator.CanValidate(policy) {
		t.Fatalf("CanValidate = false for policy %x", policyBytes)
	}
	return validator.Validate(policy)
}

func TestBasicASN1ValidatorConformingPolicy(t *testing.T) {
	algID, policyInfo, hash := policyFixture()
	digest := model.NewDigest(enumerations.DigestAlgorithmSHA256, hash)
	octets := asn1ber.WriteTLV(asn1ber.TagOctetString, hash)

	der := asn1ber.WriteSequence(concat(algID, policyInfo, octets))
	// BER indefinite-length form: BouncyCastle (and asn1ber) accept it, and the
	// digest is computed over the re-encoded DER of the children either way.
	ber := concat([]byte{0x30, 0x80}, algID, policyInfo, octets, []byte{0x00, 0x00})

	for name, policyBytes := range map[string][]byte{"DER": der, "BER indefinite": ber} {
		t.Run(name, func(t *testing.T) {
			result := validate(t, policyBytes, digest)
			if !result.IsIdentified() || !result.IsAsn1Processable() || !result.IsDigestAlgorithmsEqual() || !result.IsDigestValid() {
				t.Fatalf("conforming policy: identified=%v asn1Processable=%v algorithmsEqual=%v digestValid=%v errors=%q",
					result.IsIdentified(), result.IsAsn1Processable(), result.IsDigestAlgorithmsEqual(), result.IsDigestValid(), result.ProcessingErrors())
			}
		})
	}
}

func TestBasicASN1ValidatorMalformedPolicyIsNotDigestValid(t *testing.T) {
	algID, policyInfo, hash := policyFixture()
	digest := model.NewDigest(enumerations.DigestAlgorithmSHA256, hash)
	octets := asn1ber.WriteTLV(asn1ber.TagOctetString, hash)
	wrongOctets := asn1ber.WriteTLV(asn1ber.TagOctetString, make([]byte, len(hash)))
	conforming := asn1ber.WriteSequence(concat(algID, policyInfo, octets))

	tests := []struct {
		name string
		// policy is the policy content.
		policy []byte
		// asn1Processable is what Java reports: true as soon as the content
		// parsed as a SEQUENCE, whatever happens afterwards.
		asn1Processable bool
	}{
		// Java: asn1Sequence.getObjectAt(0) throws ArrayIndexOutOfBoundsException.
		{"empty SEQUENCE", []byte{0x30, 0x00}, true},
		// Java: getComputedDigest -> getObjectAt(1) throws.
		{"one element", asn1ber.WriteSequence(algID), true},
		// Java: the re-calculated digest matches (digestValid=true), then
		// getObjectAt(2) throws and the catch resets digestValid to false.
		{"two elements (T23-SEC-001)", asn1ber.WriteSequence(concat(algID, policyInfo)), true},
		// getObjectAt(2) is not an ASN1OctetString: ClassCastException.
		{"third element not an OCTET STRING", asn1ber.WriteSequence(concat(algID, policyInfo, policyInfo)), true},
		{"first element not an AlgorithmIdentifier", asn1ber.WriteSequence(concat(octets, policyInfo, octets)), true},
		{"digest value from the policy file differs", asn1ber.WriteSequence(concat(algID, policyInfo, wrongOctets)), true},
		// toASN1Primitive fails / ClassCastException to ASN1Sequence: not processable.
		{"extra data after the SEQUENCE", concat(conforming, []byte{0x00}), false},
		{"truncated SEQUENCE", conforming[:len(conforming)-1], false},
		{"SEQUENCE tag with an invalid body", []byte{0x30, 0x05, 0x02}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validate(t, tt.policy, digest)
			if result.IsDigestValid() {
				t.Errorf("digestValid = true, want false (Java's catch marks the digest invalid)")
			}
			if result.IsAsn1Processable() != tt.asn1Processable {
				t.Errorf("asn1Processable = %v, want %v", result.IsAsn1Processable(), tt.asn1Processable)
			}
			if result.ProcessingErrors() == "" {
				t.Errorf("no processing error recorded")
			}
		})
	}
}

func TestBasicASN1ValidatorDigestAlgorithmMismatch(t *testing.T) {
	algID, policyInfo, hash := policyFixture()
	octets := asn1ber.WriteTLV(asn1ber.TagOctetString, hash)
	policyBytes := asn1ber.WriteSequence(concat(algID, policyInfo, octets))

	result := validate(t, policyBytes, model.NewDigest(enumerations.DigestAlgorithmSHA512, make([]byte, 64)))
	if result.IsDigestValid() || result.IsDigestAlgorithmsEqual() || !result.IsAsn1Processable() {
		t.Fatalf("algorithm mismatch: digestValid=%v algorithmsEqual=%v asn1Processable=%v",
			result.IsDigestValid(), result.IsDigestAlgorithmsEqual(), result.IsAsn1Processable())
	}
}
