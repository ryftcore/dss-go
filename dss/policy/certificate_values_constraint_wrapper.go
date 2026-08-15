// Ported from dss-policy-jaxb/.../policy/CertificateValuesConstraintWrapper.java (DSS 6.5.RC1).
package policy

import (
	modelpolicy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/policy/jaxb"
)

// CertificateValuesConstraintWrapper wraps
// eu.europa.esig.dss.policy.jaxb.CertificateValuesConstraint into a
// eu.europa.esig.dss.model.policy.CertificateApplicabilityRule.
type CertificateValuesConstraintWrapper struct {
	LevelConstraintWrapper
	constraint *jaxb.CertificateValuesConstraint
}

// NewCertificateValuesConstraintWrapper is the default constructor.
func NewCertificateValuesConstraintWrapper(constraint *jaxb.CertificateValuesConstraint) *CertificateValuesConstraintWrapper {
	var base *jaxb.LevelConstraint
	if constraint != nil {
		base = &constraint.LevelConstraint
	}
	return &CertificateValuesConstraintWrapper{
		LevelConstraintWrapper: LevelConstraintWrapper{constraint: base},
		constraint:             constraint,
	}
}

// CertificateExtensions returns a list of certificate extensions satisfying
// the condition. Ports
// CertificateValuesConstraintWrapper#getCertificateExtensions.
func (w *CertificateValuesConstraintWrapper) CertificateExtensions() modelpolicy.MultiValuesRule {
	if w.constraint != nil {
		return NewMultiValuesConstraintWrapper(w.constraint.CertificateExtensions)
	}
	return NewMultiValuesConstraintWrapper(nil)
}

// CertificatePolicies returns a list of certificate policies satisfying the
// condition. Ports CertificateValuesConstraintWrapper#getCertificatePolicies.
func (w *CertificateValuesConstraintWrapper) CertificatePolicies() modelpolicy.MultiValuesRule {
	if w.constraint != nil {
		return NewMultiValuesConstraintWrapper(w.constraint.CertificatePolicies)
	}
	return NewMultiValuesConstraintWrapper(nil)
}
