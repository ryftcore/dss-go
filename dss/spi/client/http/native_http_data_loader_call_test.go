package http

import (
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestCallClosesItsConnection guards against leaking one keep-alive connection (and its reader and
// writer goroutines) per Call: each call owns a private http.Transport, so nothing would ever
// reuse or close the idle connection the response leaves behind.
func TestCallClosesItsConnection(t *testing.T) {
	var closed atomic.Int32
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("revocation data"))
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateClosed {
			closed.Add(1)
		}
	}
	server.Start()
	defer server.Close()

	envelope, err := NewNativeHTTPDataLoaderCall(server.URL).Call()
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if got := string(envelope.ResponseBody()); got != "revocation data" {
		t.Fatalf("response body = %q", got)
	}

	// The server observes the client's close asynchronously.
	deadline := time.Now().Add(5 * time.Second)
	for closed.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if closed.Load() == 0 {
		t.Fatal("the connection of a finished Call is still open")
	}
}

// TestCallMaxInputSize checks the response size limit, which is off by default (as upstream).
func TestCallMaxInputSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(make([]byte, 100))
	}))
	defer server.Close()

	unlimited := NewNativeHTTPDataLoaderCall(server.URL)
	if _, err := unlimited.Call(); err != nil {
		t.Fatalf("unlimited Call: %v", err)
	}

	limited := NewNativeHTTPDataLoaderCall(server.URL)
	limited.SetMaxInputSize(10)
	if _, err := limited.Call(); err == nil {
		t.Fatal("a response larger than the maximum input size was accepted")
	}
}
