package spi

import "testing"

// TestOCSPRespStatusCodes checks each status carries the RFC 6960 code the Java enum takes
// from org.bouncycastle.cert.ocsp.OCSPResp.
func TestOCSPRespStatusCodes(t *testing.T) {
	expected := []struct {
		status OCSPRespStatus
		code   int
	}{
		{OCSPRespStatus_SUCCESSFUL, 0},
		{OCSPRespStatus_MALFORMED_REQUEST, 1},
		{OCSPRespStatus_INTERNAL_ERROR, 2},
		{OCSPRespStatus_TRY_LATER, 3},
		{OCSPRespStatus_UNKNOWN_STATUS, 4},
		{OCSPRespStatus_SIG_REQUIRED, 5},
		{OCSPRespStatus_UNAUTHORIZED, 6},
	}

	values := OCSPRespStatusValues()
	if len(values) != len(expected) {
		t.Fatalf("OCSPRespStatusValues() has %d entries, want %d", len(values), len(expected))
	}
	for index, entry := range expected {
		if values[index] != entry.status {
			t.Errorf("values[%d] = %s, want %s", index, values[index], entry.status)
		}
		if got := entry.status.StatusCode(); got != entry.code {
			t.Errorf("%s.StatusCode() = %d, want %d", entry.status, got, entry.code)
		}
	}
}

// TestOCSPRespStatusFromInt checks the lookup resolves every known code and falls back to
// UNKNOWN_STATUS instead of failing, as fromInt(int) does.
func TestOCSPRespStatusFromInt(t *testing.T) {
	cases := map[int]OCSPRespStatus{
		0:  OCSPRespStatus_SUCCESSFUL,
		1:  OCSPRespStatus_MALFORMED_REQUEST,
		2:  OCSPRespStatus_INTERNAL_ERROR,
		3:  OCSPRespStatus_TRY_LATER,
		4:  OCSPRespStatus_UNKNOWN_STATUS,
		5:  OCSPRespStatus_SIG_REQUIRED,
		6:  OCSPRespStatus_UNAUTHORIZED,
		7:  OCSPRespStatus_UNKNOWN_STATUS,
		-1: OCSPRespStatus_UNKNOWN_STATUS,
		99: OCSPRespStatus_UNKNOWN_STATUS,
	}
	for value, want := range cases {
		if got := OCSPRespStatusFromInt(value); got != want {
			t.Errorf("OCSPRespStatusFromInt(%d) = %s, want %s", value, got, want)
		}
	}
}

// TestOCSPRespStatusValuesMatchResponseConstants checks the enum stays wired to the RFC 6960
// constants the ported structures declare.
func TestOCSPRespStatusValuesMatchResponseConstants(t *testing.T) {
	if OCSPRespStatus_SUCCESSFUL.StatusCode() != OCSPResponseStatusSuccessful {
		t.Error("SUCCESSFUL is not wired to OCSPResponseStatusSuccessful")
	}
	if OCSPRespStatus_UNAUTHORIZED.StatusCode() != OCSPResponseStatusUnauthorized {
		t.Error("UNAUTHORIZED is not wired to OCSPResponseStatusUnauthorized")
	}
}
