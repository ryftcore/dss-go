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
	// OCSPRespStatusSuccessful means the response has valid confirmations.
	OCSPRespStatusSuccessful OCSPRespStatus = "SUCCESSFUL"
	// OCSPRespStatusMalformedRequest means the confirmation request was illegal.
	OCSPRespStatusMalformedRequest OCSPRespStatus = "MALFORMED_REQUEST"
	// OCSPRespStatusInternalError means an internal error occurred in the issuer.
	OCSPRespStatusInternalError OCSPRespStatus = "INTERNAL_ERROR"
	// OCSPRespStatusTryLater means the requester should try again later.
	OCSPRespStatusTryLater OCSPRespStatus = "TRY_LATER"
	// OCSPRespStatusUnknownStatus covers the value (4), which RFC 6960 does not use.
	OCSPRespStatusUnknownStatus OCSPRespStatus = "UNKNOWN_STATUS"
	// OCSPRespStatusSigRequired means the request must be signed.
	OCSPRespStatusSigRequired OCSPRespStatus = "SIG_REQUIRED"
	// OCSPRespStatusUnauthorized means the request is unauthorized.
	OCSPRespStatusUnauthorized OCSPRespStatus = "UNAUTHORIZED"
)

// ocspRespStatusCodes maps every status onto its RFC 6960 code.
var ocspRespStatusCodes = map[OCSPRespStatus]int{
	OCSPRespStatusSuccessful:       OCSPResponseStatusSuccessful,
	OCSPRespStatusMalformedRequest: OCSPResponseStatusMalformedRequest,
	OCSPRespStatusInternalError:    OCSPResponseStatusInternalError,
	OCSPRespStatusTryLater:         OCSPResponseStatusTryLater,
	OCSPRespStatusUnknownStatus:    4,
	OCSPRespStatusSigRequired:      OCSPResponseStatusSigRequired,
	OCSPRespStatusUnauthorized:     OCSPResponseStatusUnauthorized,
}

// OCSPRespStatusValues returns the statuses in their Java declaration order.
// Port of OCSPRespStatus.values().
func OCSPRespStatusValues() []OCSPRespStatus {
	return []OCSPRespStatus{
		OCSPRespStatusSuccessful,
		OCSPRespStatusMalformedRequest,
		OCSPRespStatusInternalError,
		OCSPRespStatusTryLater,
		OCSPRespStatusUnknownStatus,
		OCSPRespStatusSigRequired,
		OCSPRespStatusUnauthorized,
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
	return OCSPRespStatusUnknownStatus
}

// StatusCode returns the RFC 6960 code of the status. Port of getStatusCode().
func (s OCSPRespStatus) StatusCode() int {
	return ocspRespStatusCodes[s]
}
