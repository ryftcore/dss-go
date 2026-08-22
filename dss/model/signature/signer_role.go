// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignerRole.java (DSS 6.5.RC1).
package signature

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// SignerRole represents the signer roles extracted from the signature.
type SignerRole struct {
	// role is the role string.
	role string

	// category is the category of the role (certified, claimed, signed).
	category enumerations.EndorsementType

	// notBefore is the certificate's 'notBefore' date for a 'certified' role category. The
	// zero time.Time stands for Java's null.
	notBefore time.Time

	// notAfter is the certificate's 'notAfter' date for a 'certified' role category. The
	// zero time.Time stands for Java's null.
	notAfter time.Time
}

// NewSignerRole is the default constructor.
func NewSignerRole(role string, category enumerations.EndorsementType) *SignerRole {
	return &SignerRole{role: role, category: category}
}

// Role gets the role. Port of getRole().
func (s *SignerRole) Role() string {
	return s.role
}

// Category gets the role category. Port of getCategory().
func (s *SignerRole) Category() enumerations.EndorsementType {
	return s.category
}

// NotBefore gets the certificate's 'notBefore' for a 'certified' role category. Port of
// getNotBefore().
func (s *SignerRole) NotBefore() time.Time {
	return s.notBefore
}

// SetNotBefore sets the certificate's 'notBefore' for a 'certified' role category. Port of
// setNotBefore(Date).
func (s *SignerRole) SetNotBefore(notBefore time.Time) {
	s.notBefore = notBefore
}

// NotAfter gets the certificate's 'notAfter' for a 'certified' role category. Port of
// getNotAfter().
func (s *SignerRole) NotAfter() time.Time {
	return s.notAfter
}

// SetNotAfter sets the certificate's 'notAfter' for a 'certified' role category. Port of
// setNotAfter(Date).
func (s *SignerRole) SetNotAfter(notAfter time.Time) {
	s.notAfter = notAfter
}

// String returns the Java toString() form. Port of toString().
func (s *SignerRole) String() string {
	return fmt.Sprintf("SignerRole [category=%s, role details=%s, notBefore=%s, notAfter=%s]",
		s.category, s.role, signerRoleDateString(s.notBefore), signerRoleDateString(s.notAfter))
}

// signerRoleDateString renders a date the way Java's string concatenation would print a
// null Date ("null") or a set one.
func signerRoleDateString(t time.Time) string {
	if t.IsZero() {
		return "null"
	}
	return t.String()
}

// Equals reports whether both SignerRoles carry the same category, role, notBefore and
// notAfter. Port of equals(Object).
func (s *SignerRole) Equals(other *SignerRole) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	if s.category != other.category {
		return false
	}
	if s.role != other.role {
		return false
	}
	if !s.notBefore.Equal(other.notBefore) {
		return false
	}
	if !s.notAfter.Equal(other.notAfter) {
		return false
	}
	return true
}
