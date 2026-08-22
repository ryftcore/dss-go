// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/ProfileParameters.java (DSS 6.5.RC1).
package document

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
)

// ProfileParameters manages the internal variables used in the process of creating of a
// signature and which allows accelerating the signature generation.
type ProfileParameters struct {
	// deterministicId is the id created in a deterministic way based on the filled parameters
	// to use in the signature file.
	deterministicId string

	// detachedContents are the cached detached contents (used for DETACHED signature creation
	// or/and ASiC containers signing).
	detachedContents []model.DSSDocument
}

// NewProfileParameters is the default constructor instantiating the object with null values.
func NewProfileParameters() *ProfileParameters {
	return &ProfileParameters{}
}

// DeterministicId ports #getDeterministicId.
func (p *ProfileParameters) DeterministicId() string {
	return p.deterministicId
}

// SetDeterministicId ports #setDeterministicId.
func (p *ProfileParameters) SetDeterministicId(deterministicId string) {
	p.deterministicId = deterministicId
}

// DetachedContents ports #getDetachedContents.
func (p *ProfileParameters) DetachedContents() []model.DSSDocument {
	return p.detachedContents
}

// SetDetachedContents ports #setDetachedContents.
func (p *ProfileParameters) SetDetachedContents(detachedContents []model.DSSDocument) {
	p.detachedContents = detachedContents
}

// String ports #toString.
func (p *ProfileParameters) String() string {
	return fmt.Sprintf("ProfileParameters{deterministicId='%s', detachedContents=%v}", p.deterministicId, p.detachedContents)
}

// Equals ports #equals.
func (p *ProfileParameters) Equals(other *ProfileParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if p.deterministicId != other.deterministicId {
		return false
	}
	return profileParametersDetachedContentsEqual(p.detachedContents, other.detachedContents)
}

// profileParametersDetachedContentsEqual ports Objects.equals(detachedContents, ...) for the
// []model.DSSDocument slice, comparing by identity/length since model.DSSDocument implementations
// do not expose a value Equals method.
func profileParametersDetachedContentsEqual(a, b []model.DSSDocument) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
