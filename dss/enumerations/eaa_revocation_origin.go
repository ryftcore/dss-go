// Ported from dss-enumerations/.../EAARevocationOrigin.java (DSS 6.5.RC1).
package enumerations

// EAARevocationOrigin represents an origin type of the EAA revocation data.
type EAARevocationOrigin string

const (
	// EAARevocationOriginExternal indicates the status data was provided
	// by the user or extracted from online source.
	EAARevocationOriginExternal EAARevocationOrigin = "EXTERNAL"
	// EAARevocationOriginCached indicates the status data was obtained
	// from a local DB or cache.
	EAARevocationOriginCached EAARevocationOrigin = "CACHED"
)

// EAARevocationOriginValues returns all constants in declaration order.
func EAARevocationOriginValues() []EAARevocationOrigin {
	return []EAARevocationOrigin{
		EAARevocationOriginExternal,
		EAARevocationOriginCached,
	}
}
