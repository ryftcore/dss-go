package alert

import (
	"errors"
	"testing"
)

type recordingHandler struct {
	calls *[]string
	name  string
	err   error
}

func (h *recordingHandler) Process(object Status) error {
	*h.calls = append(*h.calls, h.name)
	return h.err
}

func TestCompositeAlertHandler_RunsAllInOrder(t *testing.T) {
	var calls []string
	h := NewCompositeAlertHandler[Status]([]AlertHandler[Status]{
		&recordingHandler{calls: &calls, name: "first"},
		&recordingHandler{calls: &calls, name: "second"},
		&recordingHandler{calls: &calls, name: "third"},
	})

	status := NewMessageStatus()
	status.SetMessage("event")
	if err := h.Process(status); err != nil {
		t.Fatalf("Process() error = %v, want nil", err)
	}
	if got, want := calls, []string{"first", "second", "third"}; !equalStrings(got, want) {
		t.Fatalf("calls = %v, want %v", got, want)
	}
}

func TestCompositeAlertHandler_StopsAtFirstError(t *testing.T) {
	var calls []string
	boom := errors.New("boom")
	h := NewCompositeAlertHandler[Status]([]AlertHandler[Status]{
		&recordingHandler{calls: &calls, name: "first"},
		&recordingHandler{calls: &calls, name: "second", err: boom},
		&recordingHandler{calls: &calls, name: "third"},
	})

	status := NewMessageStatus()
	status.SetMessage("event")
	err := h.Process(status)
	if !errors.Is(err, boom) {
		t.Fatalf("Process() error = %v, want %v", err, boom)
	}
	if got, want := calls, []string{"first", "second"}; !equalStrings(got, want) {
		t.Fatalf("calls = %v, want %v (third handler must not run)", got, want)
	}
}

func TestCompositeAlertHandler_ComposesWithThrowAndSilent(t *testing.T) {
	h := NewCompositeAlertHandler[Status]([]AlertHandler[Status]{
		NewSilentHandler[Status](),
		NewThrowAlertExceptionHandler[Status](),
	})

	status := NewMessageStatus()
	status.SetMessage("boom message")

	err := h.Process(status)
	var alertErr *AlertError
	if !errors.As(err, &alertErr) {
		t.Fatalf("Process() error = %v, want *AlertError", err)
	}
	if got, want := alertErr.Error(), status.String(); got != want {
		t.Fatalf("AlertError.Error() = %q, want %q", got, want)
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
