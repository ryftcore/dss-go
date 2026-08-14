package alert

import "testing"

func TestSilentOnStatusAlert_NeverErrors(t *testing.T) {
	a := NewSilentOnStatusAlert()

	status := NewMessageStatus()
	status.SetMessage("something happened")

	if err := a.Alert(status); err != nil {
		t.Fatalf("Alert() error = %v, want nil", err)
	}
	if err := a.Alert(NewMessageStatus()); err != nil {
		t.Fatalf("Alert() error = %v, want nil for empty status", err)
	}
}

func TestSilentOnAlert_UsesProvidedDetector(t *testing.T) {
	detected := false
	detector := detectorFunc[int](func(v int) bool {
		detected = true
		return v > 10
	})

	a := NewSilentOnAlert[int](detector)

	if err := a.Alert(5); err != nil {
		t.Fatalf("Alert() error = %v, want nil", err)
	}
	if !detected {
		t.Fatalf("expected detector to run")
	}
	if err := a.Alert(20); err != nil {
		t.Fatalf("Alert() error = %v, want nil (silent handler never errors)", err)
	}
}

// detectorFunc adapts a function to the AlertDetector interface for tests.
type detectorFunc[T any] func(T) bool

func (f detectorFunc[T]) Detect(object T) bool { return f(object) }
