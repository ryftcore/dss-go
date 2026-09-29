package validation

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// Regression tests for the two places where the error of
// SignatureValidationContext.getOrderedCertificateChains used to be dropped
// (T22D1A-SEC-001, T22D1B-SEC-001): Java's CertificateReorderer#getOrderedCertificateChains throws a
// DSSException when the processed certificates cannot be ordered, which aborts the
// required-revocation-data / POE checks; Go reported such a status as empty, i.e. satisfied.

// expectDSSErrorPanic runs fn and requires it to panic with a *model.DSSError carrying message.
func expectDSSErrorPanic(t *testing.T, what, message string, fn func()) {
	t.Helper()
	defer func() {
		t.Helper()
		recovered := recover()
		if recovered == nil {
			t.Fatalf("%s: did not panic, i.e. the unorderable certificate set was reported as a satisfied check", what)
		}
		err, ok := recovered.(error)
		var dssErr *model.DSSError
		if !ok || !errors.As(err, &dssErr) {
			t.Fatalf("%s: panicked with %T (%v), want a *model.DSSError", what, recovered, recovered)
		}
		if dssErr.Message != message {
			t.Errorf("%s: DSSError message = %q, want %q", what, dssErr.Message, message)
		}
	}()
	fn()
}

func TestUnorderableCertificateSetAbortsRevocationChecks(t *testing.T) {
	// No processed certificate at all: CertificateReorderer#getSigningCertificates throws
	// DSSException("No signing certificate found").
	const message = "No signing certificate found"

	newContext := func() *SignatureValidationContext {
		context := NewSignatureValidationContext()
		context.Initialize(NewCommonCertificateVerifierSimple(true))
		return context
	}

	expectDSSErrorPanic(t, "CheckAllRequiredRevocationDataPresent", message, func() {
		newContext().CheckAllRequiredRevocationDataPresent()
	})
	expectDSSErrorPanic(t, "CheckAllPOECoveredByRevocationData", message, func() {
		newContext().CheckAllPOECoveredByRevocationData()
	})
	// The alerter's default configuration alerts on missing revocation data (exception) and on
	// uncovered POE (log), so both assertions run the check.
	expectDSSErrorPanic(t, "AssertAllRequiredRevocationDataPresent", message, func() {
		NewSignatureValidationAlerter(newContext()).AssertAllRequiredRevocationDataPresent()
	})
	expectDSSErrorPanic(t, "AssertAllPOECoveredByRevocationData", message, func() {
		NewSignatureValidationAlerter(newContext()).AssertAllPOECoveredByRevocationData()
	})
	// checkAtLeastOneRevocationDataPresentAfterBestSignatureTime shares the helper.
	expectDSSErrorPanic(t, "getOrderedCertificateChains", message, func() {
		newContext().getOrderedCertificateChains()
	})
}

// freshnessStub is the part of a revocation token checkRevocationForCertificateChainAgainstBestSignatureTime
// reads; the embedded nil *spi.CRLToken supplies the rest of AnyRevocationToken.
type freshnessStub struct {
	*spi.CRLToken
	dssID                string
	relatedCertificateID string
	thisUpdate           time.Time
	nextUpdate           time.Time
}

func (s *freshnessStub) DSSIDAsString() string        { return s.dssID }
func (s *freshnessStub) RelatedCertificateID() string { return s.relatedCertificateID }
func (s *freshnessStub) ThisUpdate() time.Time        { return s.thisUpdate }
func (s *freshnessStub) NextUpdate() time.Time        { return s.nextUpdate }

// testCertificateIssuer is a certificate together with its private key.
type testCertificateIssuer struct {
	certificate *x509.Certificate
	key         *ecdsa.PrivateKey
}

// issueTestCertificate creates a certificate for commonName, self-signed when parent is nil and
// issued by parent otherwise.
func issueTestCertificate(t *testing.T, commonName string, serial int64, ca bool, parent *testCertificateIssuer) *testCertificateIssuer {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(serial),
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  ca,
		BasicConstraintsValid: ca,
	}
	if ca {
		template.KeyUsage = x509.KeyUsageCertSign
	}
	signerCertificate, signerKey := template, key
	if parent != nil {
		signerCertificate, signerKey = parent.certificate, parent.key
	}
	der, err := x509.CreateCertificate(rand.Reader, template, signerCertificate, &key.PublicKey, signerKey)
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return &testCertificateIssuer{certificate: certificate, key: key}
}

func (i *testCertificateIssuer) token(t *testing.T) *model.CertificateToken {
	t.Helper()
	token, err := model.NewCertificateToken(i.certificate)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

// freshnessTestLeaf issues a CA and a leaf certificate below it and returns the (not self-signed,
// not trusted) leaf.
func freshnessTestLeaf(t *testing.T) *model.CertificateToken {
	t.Helper()
	ca := issueTestCertificate(t, "freshness test CA", 1, true, nil)
	return issueTestCertificate(t, "freshness test leaf", 2, false, ca).token(t)
}

// TestFreshnessStatusRecordsRevocationNextUpdate pins Java's
// `status instanceof RevocationFreshnessStatus` branch of
// checkRevocationForCertificateChainAgainstBestSignatureTime (T22D1A-STD-001): the nextUpdate of a
// revocation that is too old for the best signature time is recorded on the freshness status,
// which then reports it in its error string.
func TestFreshnessStatusRecordsRevocationNextUpdate(t *testing.T) {
	validationTime := time.Date(2024, time.March, 1, 12, 0, 0, 0, time.UTC)
	context := NewSignatureValidationContextAtTime(validationTime)
	context.Initialize(NewCommonCertificateVerifierSimple(true))

	leaf := freshnessTestLeaf(t)
	context.AddCertificateTokenForVerification(leaf)
	// Mark the leaf as processed so GetRevocationData does not try to fetch revocation data.
	context.isYetVerified(leaf)

	// The revocation data was issued before, and expires after, the validation time, yet it is
	// not fresh with respect to a best signature time that lies after its nextUpdate.
	nextUpdate := validationTime.Add(time.Hour)
	context.processedRevocations = append(context.processedRevocations, &freshnessStub{
		dssID:                "R-freshness",
		relatedCertificateID: leaf.DSSIDAsString(),
		thisUpdate:           validationTime.Add(-time.Hour),
		nextUpdate:           nextUpdate,
	})
	bestSignatureTime := validationTime.Add(2 * time.Hour)

	status := NewRevocationFreshnessStatus()
	context.checkRevocationForCertificateChainAgainstBestSignatureTime(
		[]*model.CertificateToken{leaf}, bestSignatureTime, status, enumerations.ContextSignature)

	if status.IsEmpty() {
		t.Fatal("the stale revocation data was accepted")
	}
	if got := status.MinimalNextUpdateTime(); !got.Equal(nextUpdate) {
		t.Errorf("MinimalNextUpdateTime = %v, want %v (Java records it for a RevocationFreshnessStatus)", got, nextUpdate)
	}
	if got := status.TokenRevocationNextUpdateTime(leaf); !got.Equal(nextUpdate) {
		t.Errorf("TokenRevocationNextUpdateTime = %v, want %v", got, nextUpdate)
	}

	// A plain TokenStatus never records a nextUpdate time, exactly as in Java.
	plain := NewTokenStatus()
	context.checkRevocationForCertificateChainAgainstBestSignatureTime(
		[]*model.CertificateToken{leaf}, bestSignatureTime, plain, enumerations.ContextSignature)
	if plain.IsEmpty() {
		t.Fatal("the stale revocation data was accepted for a TokenStatus")
	}
}

// issuerChainStub is the part of a revocation token isSelfIssuedRevocation reads.
type issuerChainStub struct {
	*spi.CRLToken
	issuer       *model.CertificateToken
	certificates []*model.CertificateToken
}

func (s *issuerChainStub) IssuerCertificateToken() *model.CertificateToken { return s.issuer }
func (s *issuerChainStub) Certificates() []*model.CertificateToken         { return s.certificates }

// TestSelfIssuedRevocationCheckAbortsOnUnorderableIssuerChain pins RevocationDataVerifier's
// isSelfIssuedRevocation (T22D1A-SEC-002): Java's `new CertificateReorderer(issuer, certificates)
// .getOrderedCertificates()` throws a DSSException that nothing catches, so a revocation whose
// certificates cannot be ordered aborts the check instead of being treated as "not self-issued".
func TestSelfIssuedRevocationCheckAbortsOnUnorderableIssuerChain(t *testing.T) {
	ca := issueTestCertificate(t, "issuer CA", 10, true, nil)
	issued := issueTestCertificate(t, "issued by the issuer", 11, false, ca)
	unrelated := issueTestCertificate(t, "unrelated", 12, false, nil)
	verified := issueTestCertificate(t, "verified certificate", 13, false, nil).token(t)

	verifier := NewDefaultRevocationDataVerifier()

	// The issuer signed one of the certificates, so it is not a leaf; the two leaves (the issued
	// and the unrelated certificate) make two chains and neither of them starts at the issuer:
	// "Unable to determine a signing certificate : No pertinent input parameters".
	unorderable := &issuerChainStub{
		issuer:       ca.token(t),
		certificates: []*model.CertificateToken{issued.token(t), unrelated.token(t)},
	}
	expectDSSErrorPanic(t, "isSelfIssuedRevocation", "Unable to determine a signing certificate : No pertinent input parameters", func() {
		verifier.isSelfIssuedRevocation(verified, unorderable)
	})

	// An orderable chain that does not contain the verified certificate is not self-issued, and
	// one that does contain it is.
	orderable := &issuerChainStub{issuer: ca.token(t), certificates: []*model.CertificateToken{issued.token(t)}}
	if verifier.isSelfIssuedRevocation(verified, orderable) {
		t.Error("a revocation whose chain does not contain the certificate was reported as self-issued")
	}
	if !verifier.isSelfIssuedRevocation(issued.token(t), orderable) {
		t.Error("a revocation whose chain contains the certificate was not reported as self-issued")
	}
}
