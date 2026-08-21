// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/CertificateContentEquivalence.java (DSS 6.5.RC1).
package tsl

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// CertificateContentEquivalence contains information about an MRA equivalence mapping.
//
// java.io.Serializable has no Go counterpart and is dropped.
type CertificateContentEquivalence struct {
	// context defines the context of the certificate content equivalence (i.e. QcCompliance,
	// QcType, etc.).
	context enumerations.MRAEquivalenceContext
	// condition defines rules to trigger the equivalence translation.
	condition Condition
	// contentReplacement contains OIDs of the equivalence rule.
	contentReplacement *QCStatementOids
}

// NewCertificateContentEquivalence instantiates an object with zero values. Port of the
// default constructor.
func NewCertificateContentEquivalence() *CertificateContentEquivalence {
	return &CertificateContentEquivalence{}
}

// Context gets the certificate content equivalence context.
func (c *CertificateContentEquivalence) Context() enumerations.MRAEquivalenceContext {
	return c.context
}

// SetContext sets the context of the certificate content equivalence (i.e. QcCompliance,
// QcType, etc.).
func (c *CertificateContentEquivalence) SetContext(context enumerations.MRAEquivalenceContext) {
	c.context = context
}

// Condition gets the equivalence condition.
func (c *CertificateContentEquivalence) Condition() Condition {
	return c.condition
}

// SetCondition sets the equivalence condition.
func (c *CertificateContentEquivalence) SetCondition(condition Condition) {
	c.condition = condition
}

// ContentReplacement gets the defined OIDs.
func (c *CertificateContentEquivalence) ContentReplacement() *QCStatementOids {
	return c.contentReplacement
}

// SetContentReplacement sets the equivalence OIDs.
func (c *CertificateContentEquivalence) SetContentReplacement(contentReplacement *QCStatementOids) {
	c.contentReplacement = contentReplacement
}

// String returns the Java toString() form.
func (c *CertificateContentEquivalence) String() string {
	return fmt.Sprintf("CertificateContentEquivalence [context=%v, condition=%v, contentReplacement=%v]",
		c.context, c.condition, c.contentReplacement)
}

// Equals ports CertificateContentEquivalence#equals(Object).
//
// Java's Objects.equals(condition, that.condition) reduces to a Condition interface value
// comparison; Go compares the interface values directly (equal when both are nil or hold
// identical dynamic types/values), which matches for the comparable Condition implementations
// this package defines.
func (c *CertificateContentEquivalence) Equals(other *CertificateContentEquivalence) bool {
	if other == nil {
		return false
	}
	if c == other {
		return true
	}
	if c.context != other.context {
		return false
	}
	if c.condition != other.condition {
		return false
	}
	if c.contentReplacement == other.contentReplacement {
		return true
	}
	if c.contentReplacement == nil || other.contentReplacement == nil {
		return false
	}
	return c.contentReplacement.Equals(other.contentReplacement)
}
