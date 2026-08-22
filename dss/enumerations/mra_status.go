// Ported from dss-enumerations/.../MRAStatus.java (DSS 6.5.RC1).
//
// It specifies the current status of the MRA for the corresponding trust
// service type identified in the TrustServiceLegalIdentifier field.
package enumerations

import "fmt"

// MRAStatus represents the current status of the MRA. Implements UriBasedEnum.
type MRAStatus string

const (
	// MRAStatusEnacted is used to denote a valid status.
	MRAStatusEnacted MRAStatus = "ENACTED"
	// MRAStatusRepealed is used to denote an invalid status.
	MRAStatusRepealed MRAStatus = "REPEALED"
)

// mraStatusURIs holds the URI for each constant.
var mraStatusURIs = map[MRAStatus]string{
	MRAStatusEnacted:  "http://ec.europa.eu/tools/lotl/mra/enacted",
	MRAStatusRepealed: "http://ec.europa.eu/tools/lotl/mra/repealed",
}

// MRAStatusValues returns all constants in declaration order.
func MRAStatusValues() []MRAStatus {
	return []MRAStatus{
		MRAStatusEnacted,
		MRAStatusRepealed,
	}
}

// MRAStatusValueOf returns the MRAStatus matching the given Java enum name.
func MRAStatusValueOf(name string) (MRAStatus, error) {
	for _, v := range MRAStatusValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant MRAStatus.%s", name)
}

// URI returns the URI of the MRA status. Implements UriBasedEnum.
func (m MRAStatus) URI() string {
	return mraStatusURIs[m]
}

// IsEnacted returns whether the MRA Status corresponds to the enacted Trust
// Service equivalence schema.
func (m MRAStatus) IsEnacted() bool {
	return m == MRAStatusEnacted
}

// compile-time interface assertion.
var _ UriBasedEnum = MRAStatus("")
