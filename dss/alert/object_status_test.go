package alert

import "testing"

func TestObjectStatus_EmptyByDefault(t *testing.T) {
	s := NewObjectStatus()
	if !s.IsEmpty() {
		t.Fatalf("expected empty ObjectStatus")
	}
	if got, want := s.String(), "Status : Valid"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestObjectStatus_ErrorStringFormat(t *testing.T) {
	s := NewObjectStatus()
	s.SetMessage("Signature is not valid")
	s.AddRelatedObjectIdentifierAndErrorMessage("sig-1", "no valid revocation data")
	s.AddRelatedObjectIdentifierAndErrorMessage("sig-2", "certificate expired")

	if s.IsEmpty() {
		t.Fatalf("expected non-empty ObjectStatus")
	}

	// Related object map is unordered in Java (HashMap); ObjectStatus sorts keys for
	// determinism, so the expected order here is alphabetical by object id.
	want := "Signature is not valid [sig-1: no valid revocation data; sig-2: certificate expired]"
	if got := s.ErrorString(); got != want {
		t.Fatalf("ErrorString() = %q, want %q", got, want)
	}
	if got := s.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestObjectStatus_ErrorStringNoRelatedObjects(t *testing.T) {
	s := NewObjectStatus()
	s.SetMessage("Signature is not valid")

	want := "Signature is not valid "
	if got := s.ErrorString(); got != want {
		t.Fatalf("ErrorString() = %q, want %q", got, want)
	}
}

func TestObjectStatus_MessageForObjectWithId(t *testing.T) {
	s := NewObjectStatus()
	s.AddRelatedObjectIdentifierAndErrorMessage("id-1", "error-1")

	if got, want := s.MessageForObjectWithId("id-1"), "error-1"; got != want {
		t.Fatalf("MessageForObjectWithId() = %q, want %q", got, want)
	}
	if got := s.MessageForObjectWithId("missing"); got != "" {
		t.Fatalf("MessageForObjectWithId(missing) = %q, want empty", got)
	}
}

func TestObjectStatus_RelatedObjectIds(t *testing.T) {
	s := NewObjectStatus()
	s.AddRelatedObjectIdentifierAndErrorMessage("id-1", "err-1")
	s.AddRelatedObjectIdentifierAndErrorMessage("id-2", "err-2")

	ids := s.RelatedObjectIds()
	if len(ids) != 2 {
		t.Fatalf("RelatedObjectIds() length = %d, want 2", len(ids))
	}
	found := map[string]bool{}
	for _, id := range ids {
		found[id] = true
	}
	if !found["id-1"] || !found["id-2"] {
		t.Fatalf("RelatedObjectIds() = %v, missing expected ids", ids)
	}
}

func TestObjectStatus_ImplementsStatus(t *testing.T) {
	var _ Status = NewObjectStatus()
}

func TestObjectStatus_NotEmptyWithRelatedObjectOnly(t *testing.T) {
	s := NewObjectStatus()
	s.AddRelatedObjectIdentifierAndErrorMessage("id-1", "err-1")
	if s.IsEmpty() {
		t.Fatalf("expected non-empty ObjectStatus when a related object is present, even with no message")
	}
}
