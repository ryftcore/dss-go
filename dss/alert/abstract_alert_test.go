package alert

import "testing"

func TestAbstractAlert_PanicsWithoutDetector(t *testing.T) {
	a := NewAbstractAlert[int](nil, NewSilentHandler[int]())
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic when detector is nil")
		}
		if r != "AlertDetector shall be defined!" {
			t.Fatalf("panic = %v, want %q", r, "AlertDetector shall be defined!")
		}
	}()
	_ = a.Alert(1)
}

func TestAbstractAlert_PanicsWithoutHandler(t *testing.T) {
	a := NewAbstractAlert[int](detectorFunc[int](func(int) bool { return true }), nil)
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic when handler is nil")
		}
		if r != "AlertHandler shall be defined!" {
			t.Fatalf("panic = %v, want %q", r, "AlertHandler shall be defined!")
		}
	}()
	_ = a.Alert(1)
}

func TestAbstractAlert_HandlerNotCalledWhenNotDetected(t *testing.T) {
	called := false
	a := NewAbstractAlert[int](
		detectorFunc[int](func(int) bool { return false }),
		alertHandlerFunc[int](func(int) error { called = true; return nil }),
	)
	if err := a.Alert(42); err != nil {
		t.Fatalf("Alert() error = %v, want nil", err)
	}
	if called {
		t.Fatalf("handler should not run when detector returns false")
	}
}

type alertHandlerFunc[T any] func(T) error

func (f alertHandlerFunc[T]) Process(object T) error { return f(object) }
