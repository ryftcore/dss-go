package model

import (
	"bytes"
	"testing"
)

func TestTimestampBinaryRoundTrip(t *testing.T) {
	data := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	tb := NewTimestampBinary(data)
	if !bytes.Equal(tb.Bytes(), data) {
		t.Fatalf("Bytes() = %x, want %x", tb.Bytes(), data)
	}
}
