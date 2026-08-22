// Ported from dss-i18n/.../i18n/I18nProviderTest.java (DSS 6.5.RC1).
//
// Upstream's Locale-fallback tests (French/German producing the same
// default-locale message text) collapse into a single assertion here:
// see the defaultBundle doc comment in i18n_provider.go for why every
// locale resolves to the same embedded bundle in this port.
package i18n

import "testing"

func TestI18nProviderGetMessage(t *testing.T) {
	provider := NewI18nProviderForLocale("en")
	got := provider.GetMessage(MessageTagBBBXCVCCCBB)
	want := "Can the certificate chain be built till a trust anchor?"
	if got != want {
		t.Fatalf("GetMessage(BBB_XCV_CCCBB) = %q, want %q", got, want)
	}

	// locale-invariance: see defaultBundle doc comment.
	for _, locale := range []string{"", "en", "fr", "fr_FR", "de"} {
		p := NewI18nProviderForLocale(locale)
		if got := p.GetMessage(MessageTagBBBXCVCCCBB); got != want {
			t.Errorf("locale %q: GetMessage(BBB_XCV_CCCBB) = %q, want %q", locale, got, want)
		}
	}
}

func TestI18nProviderGetMessageNullTag(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for empty MessageTag")
		}
		if r != "messageTag cannot be null!" {
			t.Fatalf("panic value = %v, want %q", r, "messageTag cannot be null!")
		}
	}()
	NewI18nProvider().GetMessage(MessageTag(""))
}

func TestI18nProviderParametrizedTest(t *testing.T) {
	provider := NewI18nProvider()
	message := provider.GetMessage(MessageTagTrustedServiceStatus, "granted")
	want := "Status : granted"
	if message != want {
		t.Fatalf("GetMessage(TRUSTED_SERVICE_STATUS, granted) = %q, want %q", message, want)
	}
}

func TestI18nProviderNestedMessageTagTest(t *testing.T) {
	provider := NewI18nProvider()
	message := provider.GetMessage(MessageTagCertQualificationAtTime, MessageTagVTValidationTime)
	want := "Certificate Qualification at validation time"
	if message != want {
		t.Fatalf("GetMessage(CERT_QUALIFICATION_AT_TIME, VT_VALIDATION_TIME) = %q, want %q", message, want)
	}
}

// TestI18nProviderApostropheTest reproduces
// I18nProviderTest#apostropheTest: MessageFormat's quoting rules treat a
// doubled single quote (”) in the pattern as an escaped literal quote.
func TestI18nProviderApostropheTest(t *testing.T) {
	provider := NewI18nProvider()

	if got, want := provider.GetMessage(MessageTagBBBCVISIT), "Is time-stamp's signature intact?"; got != want {
		t.Errorf("GetMessage(BBB_CV_ISIT) = %q, want %q", got, want)
	}
	if got, want := provider.GetMessage(MessageTagBBBICSISASCP), "Is the signed attribute: 'signing-certificate' present?"; got != want {
		t.Errorf("GetMessage(BBB_ICS_ISASCP) = %q, want %q", got, want)
	}
	if got, want := provider.GetMessage(MessageTagBBBSAVISQPMDOSPP), "Is the signed qualifying property: 'message-digest' or 'SignedProperties' present?"; got != want {
		t.Errorf("GetMessage(BBB_SAV_ISQPMDOSPP) = %q, want %q", got, want)
	}
}

// TestI18nProviderBareApostropheQuirk pins down the upstream
// java.text.MessageFormat quoting quirk described in message_format.go's
// file header: a lone (undoubled) apostrophe in a pattern silently
// disappears from the formatted output instead of being preserved,
// because it opens an unterminated quoted section. This is deliberately
// NOT "fixed" - see PORTING.md's marshal-parity contract and the porter
// brief.
func TestI18nProviderBareApostropheQuirk(t *testing.T) {
	provider := NewI18nProvider()

	got := provider.GetMessage(MessageTagBBBXCVRevocSelfIssuedOCSPANS)
	want := "The checked certificate shall not appear in the OCSP Responders certificate path!"
	if got != want {
		t.Fatalf("GetMessage(BBB_XCV_REVOC_SELF_ISSUED_OCSP_ANS) = %q, want %q (bare-apostrophe quirk)", got, want)
	}
}

func TestI18nProviderUnknownTagFallsBackToId(t *testing.T) {
	// Every MessageTag has a bundle entry (see TestAllMessageTagsResolve
	// in message_tag_test.go); simulate the "missing key" branch directly
	// against the bundle rather than mutating the shared defaultBundle.
	provider := &I18nProvider{bundle: propertyBundle{}}
	tag := MessageTagBBBXCVCCCBB
	if got := provider.GetMessage(tag); got != tag.Id() {
		t.Fatalf("GetMessage with empty bundle = %q, want %q", got, tag.Id())
	}
}
