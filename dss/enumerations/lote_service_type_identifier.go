// Ported from dss-enumerations/.../LoTEServiceTypeIdentifier.java (DSS 6.5.RC1).
package enumerations

// LoTEServiceTypeIdentifier represents a LoTE service type identifier.
type LoTEServiceTypeIdentifier interface {
	UriBasedEnum

	// Label gets user-friendly identifier.
	Label() string
}

// LoTEServiceTypeIdentifierFromURI returns a LoTEServiceTypeIdentifier for
// the given URI, or nil if none of the registered LoTELoaders resolve it.
func LoTEServiceTypeIdentifierFromURI(uri string) LoTEServiceTypeIdentifier {
	// Java's fromUri opens with `if (uri == null) return null;`. A Go string
	// cannot be null, and this codebase's diagnostic wrappers already carry an
	// absent service Type/Status as the EMPTY STRING (see
	// validation/process/qualification/trusted_entity_service_by_status_filter.go,
	// which tests `service.Status == ""` for Java's `status == null`), so the
	// empty string is this port's null and the guard has to be reproduced -
	// otherwise LoTEEmptyLoader, which answers EVERY uri with a non-nil
	// wrapper, turns Java's null into a bogus empty-URI constant. That leaked
	// into the certificate-approval-status report as a stray
	// <ServiceStatus></ServiceStatus> and, through
	// CertificateApprovalStatusFromDefinition, as the label "Certificate for
	// Unknown usage" where upstream reports "PID Provider". Found by a full-corpus report
	// byte-parity run on eaa-validation/diag_data_pid.xml.
	if uri == "" {
		return nil
	}
	for _, loader := range loTELoaders() {
		if identifier := loader.ServiceTypeIdentifierFromURI(uri); identifier != nil {
			return identifier
		}
	}
	return nil
}
