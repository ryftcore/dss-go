// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/dto/condition/PolicyIdCondition.java (DSS 6.5.RC1).
package tsl

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PolicyIdCondition checks if a certificate has a specific policy OID. Objects based on this
// type are instantiated from a trusted list or by SignedDocumentValidator for QCP and QCPPlus.
//
// java.io.Serializable has no Go counterpart and is dropped.
type PolicyIdCondition struct {
	// policyOid is the policy OID to be checked if present in the certificate's policies.
	policyOid string
}

// NewPolicyIdCondition is the default constructor for PolicyIdCondition. Port of
// PolicyIdCondition(String).
//
// Panics with the Java message when policyId is empty (Objects.requireNonNull(policyId,
// "Policy Id must be defined"); the empty string is the Go stand-in for a null String, and no
// caller in the ported tree passes a deliberately empty OID).
func NewPolicyIdCondition(policyId string) *PolicyIdCondition {
	if policyId == "" {
		panic("Policy Id must be defined")
	}
	return &PolicyIdCondition{policyOid: policyId}
}

// PolicyOid returns the policy OID: never empty. Port of the final getPolicyOid().
func (c *PolicyIdCondition) PolicyOid() string {
	return c.policyOid
}

// Check returns true if the condition is evaluated to true for the given certificate. Port of
// check(CertificateToken).
//
// Certificate policies identifier: 2.5.29.32 (IETF RFC 3280); all of the certificate's policies
// are read.
//
// Panics with the Java message when certificateToken is nil
// (Objects.requireNonNull(certificateToken, "Certificate cannot be null")).
func (c *PolicyIdCondition) Check(certificateToken *model.CertificateToken) bool {
	if certificateToken == nil {
		panic("Certificate cannot be null")
	}
	certificatePolicies := spi.CertificateExtensionsUtilsCertificatePolicies(certificateToken)
	if certificatePolicies != nil {
		for _, certificatePolicy := range certificatePolicies.PolicyList() {
			if c.policyOid == certificatePolicy.Oid() {
				return true
			}
		}
	}
	return false
}

// ToString returns a human readable condition using the given indentation. Port of
// toString(String indent).
//
// Go has no null string; the Java `if (indent == null) indent = ""` branch is unreachable here.
func (c *PolicyIdCondition) ToString(indent string) string {
	var builder strings.Builder
	builder.WriteString(indent)
	builder.WriteString("PolicyIdCondition: ")
	builder.WriteString(c.policyOid)
	builder.WriteByte('\n')
	return builder.String()
}

// String ports toString(), i.e. toString("").
func (c *PolicyIdCondition) String() string {
	return c.ToString("")
}

// Equals ports equals(Object): same concrete type and same policy OID.
func (c *PolicyIdCondition) Equals(object any) bool {
	that, ok := object.(*PolicyIdCondition)
	if !ok || that == nil {
		return false
	}
	if c == that {
		return true
	}
	return c.policyOid == that.policyOid
}

var _ tslmodel.Condition = (*PolicyIdCondition)(nil)
