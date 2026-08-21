// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/dto/condition/CertSubjectDNAttributeCondition.java (DSS 6.5.RC1).
package tsl

import (
	"encoding/asn1"
	"strings"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CertSubjectDNAttributeCondition implements the CertSubjectDNAttribute criterion.
//
// Presence: This field is optional.
//
// Description: It provides a non empty set of OIDs. Each OID maps to a possible attribute in the
// Subject DN of the certificate. The criteria is matched if all OID refers to an attribute
// present in the DN.
//
// Format: A non-empty sequence of OIDs representing Directory attributes, whose meaning respect
// the description above. For the formal definition see CertSubjectDNAttribute element in the
// schema referenced by clause C.2 (point 3).
//
// java.io.Serializable has no Go counterpart and is dropped.
type CertSubjectDNAttributeCondition struct {
	// subjectAttributeOids is the list of DN attribute OIDs to be checked against the
	// certificate's subject DN.
	subjectAttributeOids []string
}

// NewCertSubjectDNAttributeCondition is the default constructor. Port of
// CertSubjectDNAttributeCondition(List<String>).
func NewCertSubjectDNAttributeCondition(oids []string) *CertSubjectDNAttributeCondition {
	return &CertSubjectDNAttributeCondition{subjectAttributeOids: oids}
}

// AttributeOids returns the list of DN attribute OIDs to be checked against the certificate's
// subject DN: possibly empty, never nil. Port of the final getAttributeOids(), whose Java body
// answers Collections.emptyList() for a null field and an unmodifiable view otherwise; Go has no
// unmodifiable view, so the slice is returned directly (as everywhere else in this port).
func (c *CertSubjectDNAttributeCondition) AttributeOids() []string {
	if c.subjectAttributeOids == nil {
		return []string{}
	}
	return c.subjectAttributeOids
}

// Check returns true if the condition is evaluated to true for the given certificate. Port of
// check(CertificateToken).
//
// NOTE: Java builds an ASN1ObjectIdentifier from each configured OID string, which throws
// IllegalArgumentException for a malformed OID. Condition#check declares no checked exception,
// so the Go port cannot return an error either; a malformed OID yields an empty attribute
// lookup and therefore false - the same answer Java reaches for an OID that simply is not in
// the DN, and the only value-preserving option available at this signature.
func (c *CertSubjectDNAttributeCondition) Check(certificateToken *model.CertificateToken) bool {
	if utils.IsCollectionNotEmpty(c.subjectAttributeOids) {
		subject := certificateToken.Subject()
		for _, oid := range c.subjectAttributeOids {
			identifier, err := certSubjectDNAttributeConditionParseOID(oid)
			if err != nil {
				return false
			}
			attribute := spi.DSSASN1UtilsExtractAttributeFromX500Principal(identifier, subject)
			if utils.IsStringEmpty(attribute) {
				return false
			}
		}
	}
	return true
}

// ToString returns a human readable condition using the given indentation. Port of
// toString(String indent).
//
// Go has no null string; the Java `if (indent == null) indent = ""` branch is therefore
// unreachable here and collapses into the empty-string default the caller already passes.
func (c *CertSubjectDNAttributeCondition) ToString(indent string) string {
	var builder strings.Builder
	builder.WriteString(indent)
	builder.WriteString("CertSubjectDNAttributeCondition: ")
	builder.WriteString(certSubjectDNAttributeConditionJavaListString(c.subjectAttributeOids))
	builder.WriteByte('\n')
	return builder.String()
}

// String ports toString(), i.e. toString("").
func (c *CertSubjectDNAttributeCondition) String() string {
	return c.ToString("")
}

// Equals ports equals(Object): same concrete type and equal OID lists.
func (c *CertSubjectDNAttributeCondition) Equals(object any) bool {
	that, ok := object.(*CertSubjectDNAttributeCondition)
	if !ok || that == nil {
		return false
	}
	if c == that {
		return true
	}
	return certSubjectDNAttributeConditionListEquals(c.subjectAttributeOids, that.subjectAttributeOids)
}

// certSubjectDNAttributeConditionParseOID converts a dotted-decimal OID string into an
// asn1.ObjectIdentifier, standing in for `new ASN1ObjectIdentifier(oid)`.
func certSubjectDNAttributeConditionParseOID(oid string) (asn1.ObjectIdentifier, error) {
	return asn1ber.OIDFromString(oid)
}

// certSubjectDNAttributeConditionJavaListString renders a List<String> the way Java's
// StringBuilder#append(Object) does: "null" for a null list, otherwise "[a, b]".
func certSubjectDNAttributeConditionJavaListString(list []string) string {
	if list == nil {
		return "null"
	}
	return "[" + strings.Join(list, ", ") + "]"
}

// certSubjectDNAttributeConditionListEquals ports Objects.equals(List, List): a null list is
// equal only to another null list, never to an empty one.
func certSubjectDNAttributeConditionListEquals(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
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

var _ tslmodel.Condition = (*CertSubjectDNAttributeCondition)(nil)
