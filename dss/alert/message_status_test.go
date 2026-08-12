package alert

import "testing"

func TestMessageStatus_EmptyByDefault(t *testing.T) {
	s := NewMessageStatus()
	if !s.IsEmpty() {
		t.Fatalf("expected empty MessageStatus")
	}
	if got, want := s.String(), "Status : Valid"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestMessageStatus_ErrorFormatting(t *testing.T) {
	s := NewMessageStatus()
	s.SetMessage("something went wrong")

	if s.IsEmpty() {
		t.Fatalf("expected non-empty MessageStatus")
	}
	if got, want := s.ErrorString(), "something went wrong"; got != want {
		t.Fatalf("ErrorString() = %q, want %q", got, want)
	}
	if got, want := s.String(), "something went wrong"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestMessageStatus_RelatedObjectIdsUnsupported(t *testing.T) {
	s := NewMessageStatus()
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic from RelatedObjectIds()")
		}
	}()
	s.RelatedObjectIds()
}

func TestMessageStatus_ImplementsStatus(t *testing.T) {
	var _ Status = NewMessageStatus()
}
