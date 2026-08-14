// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OCSPRespStatus.java (DSS 6.5.RC1).
//
// Upstream the constants are taken from org.bouncycastle.cert.ocsp.OCSPResp; the Go port
// reads them from the OCSPResponseStatus* constants the RFC 6960 structures of
// DSSRevocationUtils declare.
package spi

// OCSPRespStatus encapsulates the RFC 6960 OCSPResponseStatus values and offers a lookup
// that never fails.
type OCSPRespStatus string

const (
	// OCSPRespStatus_SUCCESSFUL means the response has valid confirmations.
	OCSPRespStatus_SUCCESSFUL OCSPRespStatus = "SUCCESSFUL"
	// OCSPRespStatus_MALFORMED_REQUEST means the confirmation request was illegal.
	OCSPRespStatus_MALFORMED_REQUEST OCSPRespStatus = "MALFORMED_REQUEST"
	// OCSPRespStatus_INTERNAL_ERROR means an internal error occurred in the issuer.
	OCSPRespStatus_INTERNAL_ERROR OCSPRespStatus = "INTERNAL_ERROR"
	// OCSPRespStatus_TRY_LATER means the requester should try again later.
	OCSPRespStatus_TRY_LATER OCSPRespStatus = "TRY_LATER"
	// OCSPRespStatus_UNKNOWN_STATUS covers the value (4), which RFC 6960 does not use.
	OCSPRespStatus_UNKNOWN_STATUS OCSPRespStatus = "UNKNOWN_STATUS"
	// OCSPRespStatus_SIG_REQUIRED means the request must be signed.
	OCSPRespStatus_SIG_REQUIRED OCSPRespStatus = "SIG_REQUIRED"
	// OCSPRespStatus_UNAUTHORIZED means the request is unauthorized.
	OCSPRespStatus_UNAUTHORIZED OCSPRespStatus = "UNAUTHORIZED"
)

// ocspRespStatusCodes maps every status onto its RFC 6960 code.
var ocspRespStatusCodes = map[OCSPRespStatus]int{
	OCSPRespStatus_SUCCESSFUL:        OCSPResponseStatusSuccessful,
	OCSPRespStatus_MALFORMED_REQUEST: OCSPResponseStatusMalformedRequest,
	OCSPRespStatus_INTERNAL_ERROR:    OCSPResponseStatusInternalError,
	OCSPRespStatus_TRY_LATER:         OCSPResponseStatusTryLater,
	OCSPRespStatus_UNKNOWN_STATUS:    4,
	OCSPRespStatus_SIG_REQUIRED:      OCSPResponseStatusSigRequired,
	OCSPRespStatus_UNAUTHORIZED:      OCSPResponseStatusUnauthorized,
}

// OCSPRespStatusValues returns the statuses in their Java declaration order.
// Port of OCSPRespStatus.values().
func OCSPRespStatusValues() []OCSPRespStatus {
	return []OCSPRespStatus{
		OCSPRespStatus_SUCCESSFUL,
		OCSPRespStatus_MALFORMED_REQUEST,
		OCSPRespStatus_INTERNAL_ERROR,
		OCSPRespStatus_TRY_LATER,
		OCSPRespStatus_UNKNOWN_STATUS,
		OCSPRespStatus_SIG_REQUIRED,
		OCSPRespStatus_UNAUTHORIZED,
	}
}

// OCSPRespStatusFromInt returns the status matching the given code, and UNKNOWN_STATUS when
// no status carries it. Port of fromInt(int).
func OCSPRespStatusFromInt(value int) OCSPRespStatus {
	for _, status := range OCSPRespStatusValues() {
		if ocspRespStatusCodes[status] == value {
			return status
		}
	}
	return OCSPRespStatus_UNKNOWN_STATUS
}

// StatusCode returns the RFC 6960 code of the status. Port of getStatusCode().
func (s OCSPRespStatus) StatusCode() int {
	return ocspRespStatusCodes[s]
}
