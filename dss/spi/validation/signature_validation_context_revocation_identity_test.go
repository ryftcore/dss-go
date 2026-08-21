// Regression test for the membership test behind Java's Set<RevocationToken<?>>
// processedRevocations. RevocationToken#equals compares the DSS Id *and* the related
// certificate, so one CRL covering several certificates of a chain is deliberately kept as one
// entry per certificate. Comparing DSS Ids alone collapsed them, and every certificate but the
// first silently lost its revocation data - which downgraded the detected signature level on
// real CAdES-LTA documents (see cades_upstream_cross_validation_test.go's note).
package validation

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/spi"
)

// revocationIdentityStub implements just enough of AnyRevocationToken for
// signatureValidationContextContainsRevocation, which reads exactly two of its methods. The
// embedded nil *spi.CRLToken supplies the rest of the interface; nothing calls into it.
type revocationIdentityStub struct {
	*spi.CRLToken
	dssID                string
	relatedCertificateID string
}

func (s *revocationIdentityStub) DSSIDAsString() string        { return s.dssID }
func (s *revocationIdentityStub) RelatedCertificateID() string { return s.relatedCertificateID }

func newRevocationIdentityStub(dssID, relatedCertificateID string) AnyRevocationToken {
	return &revocationIdentityStub{dssID: dssID, relatedCertificateID: relatedCertificateID}
}

// TestProcessedRevocationsKeyOnRevocationAndCertificate is the fixed defect: the same CRL
// related to two different certificates must produce two entries, not one.
func TestProcessedRevocationsKeyOnRevocationAndCertificate(t *testing.T) {
	sameCRLForFirstCertificate := newRevocationIdentityStub("R-CRL", "C-first")
	sameCRLForSecondCertificate := newRevocationIdentityStub("R-CRL", "C-second")

	revocations := signatureValidationContextAppendRevocation(nil, sameCRLForFirstCertificate)
	revocations = signatureValidationContextAppendRevocation(revocations, sameCRLForSecondCertificate)

	if len(revocations) != 2 {
		t.Fatalf("one CRL related to two certificates must be kept twice, got %d entr(y|ies)", len(revocations))
	}
	if revocations[0].RelatedCertificateID() != "C-first" || revocations[1].RelatedCertificateID() != "C-second" {
		t.Errorf("related certificates = %q, %q; want C-first, C-second",
			revocations[0].RelatedCertificateID(), revocations[1].RelatedCertificateID())
	}
}

// TestProcessedRevocationsStillDeduplicateExactRepeats checks the other half of the contract:
// the very same (revocation, certificate) pair is still added only once, so the fix did not turn
// the set into a plain append-everything list.
func TestProcessedRevocationsStillDeduplicateExactRepeats(t *testing.T) {
	revocations := signatureValidationContextAppendRevocation(nil, newRevocationIdentityStub("R-CRL", "C-first"))
	revocations = signatureValidationContextAppendRevocation(revocations, newRevocationIdentityStub("R-CRL", "C-first"))

	if len(revocations) != 1 {
		t.Fatalf("the same revocation for the same certificate must be kept once, got %d", len(revocations))
	}
}

// TestProcessedRevocationsDistinguishRevocations checks two different revocations for one
// certificate are both kept (an OCSP response and a CRL commonly both cover the signer).
func TestProcessedRevocationsDistinguishRevocations(t *testing.T) {
	revocations := signatureValidationContextAppendRevocation(nil, newRevocationIdentityStub("R-CRL", "C-first"))
	revocations = signatureValidationContextAppendRevocation(revocations, newRevocationIdentityStub("R-OCSP", "C-first"))

	if len(revocations) != 2 {
		t.Fatalf("two revocations for one certificate must both be kept, got %d", len(revocations))
	}
}
