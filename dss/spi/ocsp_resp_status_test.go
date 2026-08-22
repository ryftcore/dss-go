package spi

import "testing"

// TestOCSPRespStatusCodes checks each status carries the RFC 6960 code the Java enum takes
// from org.bouncycastle.cert.ocsp.OCSPResp.
func TestOCSPRespStatusCodes(t *testing.T) {
	expected := []struct {
		status OCSPRespStatus
		code   int
	}{
		{OCSPRespStatusSuccessful, 0},
		{OCSPRespStatusMalformedRequest, 1},
		{OCSPRespStatusInternalError, 2},
		{OCSPRespStatusTryLater, 3},
		{OCSPRespStatusUnknownStatus, 4},
		{OCSPRespStatusSigRequired, 5},
		{OCSPRespStatusUnauthorized, 6},
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
		0:  OCSPRespStatusSuccessful,
		1:  OCSPRespStatusMalformedRequest,
		2:  OCSPRespStatusInternalError,
		3:  OCSPRespStatusTryLater,
		4:  OCSPRespStatusUnknownStatus,
		5:  OCSPRespStatusSigRequired,
		6:  OCSPRespStatusUnauthorized,
		7:  OCSPRespStatusUnknownStatus,
		-1: OCSPRespStatusUnknownStatus,
		99: OCSPRespStatusUnknownStatus,
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
	if OCSPRespStatusSuccessful.StatusCode() != OCSPResponseStatusSuccessful {
		t.Error("SUCCESSFUL is not wired to OCSPResponseStatusSuccessful")
	}
	if OCSPRespStatusUnauthorized.StatusCode() != OCSPResponseStatusUnauthorized {
		t.Error("UNAUTHORIZED is not wired to OCSPResponseStatusUnauthorized")
	}
}
