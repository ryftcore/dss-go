// Ported from dss-i18n/.../i18n/MessageTagTest.java (DSS 6.5.RC1).
package i18n

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// TestAllMessageTagsResolve ports allMessagesPresent/allFRMessagesPresent:
// every MessageTag must resolve without panicking via the (single,
// locale-invariant - see i18n_provider.go) embedded bundle. Upstream
// asserts assertNotNull(message); Go strings are never nil, and one tag
// (EMPTY) legitimately maps to the empty string in
// dss-messages.properties, so the faithful equivalent is "resolving
// falls back to the bundle value, not to Id()".
func TestAllMessageTagsResolve(t *testing.T) {
	provider := NewProvider()
	for _, tag := range MessageTagValues() {
		if msg, want := provider.GetMessage(tag), defaultBundle[tag.Id()]; msg != messageFormat(want, nil) {
			t.Errorf("MessageTag %s resolved to %q, want %q", tag.Id(), msg, messageFormat(want, nil))
		}
	}
}

// TestAllMessageTagsPresentInBundle ports allMessageTagsPresent: every key
// in the embedded properties bundle must correspond to a MessageTag
// constant.
func TestAllMessageTagsPresentInBundle(t *testing.T) {
	known := make(map[string]bool, len(MessageTagValues()))
	for _, tag := range MessageTagValues() {
		known[tag.Id()] = true
	}
	if len(defaultBundle) == 0 {
		t.Fatal("defaultBundle is empty")
	}
	for key := range defaultBundle {
		if !known[key] {
			t.Errorf("MessageTag with a key [%s] does not exist!", key)
		}
	}
}

// TestMessageTagBundleBijection strengthens the upstream test into an
// exhaustive round-trip: every MessageTag constant must also have a
// bundle entry (upstream's suite never separately asserts this direction
// by name, but MessageTagTest#allMessagesPresent implies it via a
// non-null message; asserting key presence directly also guards against
// the fallback-to-Id() path masking a genuinely missing entry).
func TestMessageTagBundleBijection(t *testing.T) {
	for _, tag := range MessageTagValues() {
		if _, ok := defaultBundle[tag.Id()]; !ok {
			t.Errorf("MessageTag %s has no dss-messages.properties entry", tag.Id())
		}
	}
	if got, want := len(defaultBundle), len(MessageTagValues()); got != want {
		t.Errorf("defaultBundle has %d entries, MessageTagValues() has %d", got, want)
	}
}

// TestMessageTagGetSemantic ports allIndicationsSemanticsPresent and
// allSubIndicationsSemanticsPresent.
func TestMessageTagGetSemantic(t *testing.T) {
	for _, indication := range enumerations.IndicationValues() {
		if _, ok := MessageTagGetSemantic(string(indication)); !ok {
			t.Errorf("No SEMANTICS_ MessageTag found for Indication [%s]", indication)
		}
	}
	for _, subIndication := range enumerations.SubIndicationValues() {
		if _, ok := MessageTagGetSemantic(string(subIndication)); !ok {
			t.Errorf("No SEMANTICS_ MessageTag found for SubIndication [%s]", subIndication)
		}
	}
}

func TestMessageTagGetSemanticUnknown(t *testing.T) {
	if _, ok := MessageTagGetSemantic("NOT_A_REAL_CODE"); ok {
		t.Fatal("expected ok=false for an unknown ETSI code")
	}
}

func TestMessageTagValueOf(t *testing.T) {
	got, err := MessageTagValueOf("BBB_XCV_CCCBB")
	if err != nil {
		t.Fatalf("MessageTagValueOf: %v", err)
	}
	if got != MessageTagBBBXCVCCCBB {
		t.Fatalf("MessageTagValueOf(BBB_XCV_CCCBB) = %v, want %v", got, MessageTagBBBXCVCCCBB)
	}
	if _, err := MessageTagValueOf("NOT_A_REAL_TAG"); err == nil {
		t.Fatal("expected error for unknown MessageTag name")
	}
}

func TestMessageTagId(t *testing.T) {
	if got, want := MessageTagBBBXCVCCCBB.Id(), "BBB_XCV_CCCBB"; got != want {
		t.Fatalf("Id() = %q, want %q", got, want)
	}
}
