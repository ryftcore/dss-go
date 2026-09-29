package validation

import (
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
)

// signerInfoWithSigningTime builds a SignerInfo whose only signed attribute is a signing-time
// attribute holding value (a complete DER element).
func signerInfoWithSigningTime(value []byte) *cmscore.SignerInfo {
	return &cmscore.SignerInfo{
		SignedAttributes: cmscore.Attributes{cmscore.NewAttribute(cmscore.OIDSigningTime, value)},
	}
}

// TestSigningTimeIsValidatedLikeBouncyCastle pins the signing-time step of
// timestampTokenVerifySignerInfo to org.bouncycastle.cms.SignerInformation#verify (T22D2-SEC-001):
// getSigningTime() rejects an attribute value that is not a UTCTime/GeneralizedTime with a
// CMSException, and Time#getDate() raises an unchecked IllegalStateException for an unreadable date
// string. Neither is silently skipped, as it would be by DSSASN1Utils.getDate returning null.
func TestSigningTimeIsValidatedLikeBouncyCastle(t *testing.T) {
	// verify runs the signing-time step only: the errors below are all raised before the
	// content, the verifier or the signature are looked at.
	verify := func(signerInfo *cmscore.SignerInfo, candidate *model.CertificateToken) (bool, error) {
		return timestampTokenVerifySignerInfo(signerInfo, nil, candidate, nil, "")
	}

	t.Run("not a Time", func(t *testing.T) {
		verified, err := verify(signerInfoWithSigningTime(asn1ber.WriteTLV(asn1ber.TagOctetString, []byte("240102030405Z"))), nil)
		if verified || err == nil {
			t.Fatalf("verified = %v, err = %v; want an error", verified, err)
		}
		if want := "CMSException : signing-time attribute value not a valid 'Time' structure"; err.Error() != want {
			t.Errorf("error = %q, want %q", err.Error(), want)
		}
	})

	t.Run("unreadable UTCTime", func(t *testing.T) {
		defer func() {
			recovered := recover()
			message, _ := recovered.(string)
			if !strings.HasPrefix(message, "invalid date string: ") {
				t.Errorf("recovered %v, want the invalid date string panic", recovered)
			}
		}()
		verify(signerInfoWithSigningTime(asn1ber.WriteTLV(asn1ber.TagUTCTime, []byte("not-a-date"))), nil)
		t.Error("an unreadable signing time was silently skipped")
	})

	t.Run("certificate not valid at signing time", func(t *testing.T) {
		candidate := issueTestCertificate(t, "signing time candidate", 20, false, nil).token(t)
		// Well before the certificate's validity.
		_, err := verify(signerInfoWithSigningTime(asn1ber.WriteTLV(asn1ber.TagUTCTime, []byte("200102030405Z"))), candidate)
		if err == nil || err.Error() != "CMSVerifierCertificateNotValidException : verifier not valid at signingTime" {
			t.Fatalf("error = %v, want CMSVerifierCertificateNotValidException", err)
		}
	})
}
