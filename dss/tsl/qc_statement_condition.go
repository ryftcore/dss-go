// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/dto/condition/QCStatementCondition.java (DSS 6.5.RC1).
package tsl

import (
	"strings"

	"github.com/utain/esig/dss/model"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// QCStatementCondition contains the information extracted for a certificate equivalence
// condition.
//
// java.io.Serializable has no Go counterpart and is dropped.
type QCStatementCondition struct {
	// oid is the QcStatement OID.
	oid string

	// qcType is the QcType OID. Named qcType (not type) because "type" is not usable as an
	// identifier here without shadowing; the accessor keeps the Java name.
	qcType string

	// legislation is the QcCClegislation code.
	legislation string
}

// NewQCStatementCondition is the default constructor. Port of QCStatementCondition(String,
// String, String); each argument is the empty string when the corresponding element is absent
// (Java's null), which is exactly what the Utils.isStringNotEmpty guards in Check() test for.
func NewQCStatementCondition(oid, qcType, legislation string) *QCStatementCondition {
	return &QCStatementCondition{oid: oid, qcType: qcType, legislation: legislation}
}

// Oid gets the QcStatement OID. Port of getOid().
func (c *QCStatementCondition) Oid() string {
	return c.oid
}

// Type gets the QcType OID. Port of getType().
func (c *QCStatementCondition) Type() string {
	return c.qcType
}

// Legislation gets the QcCClegislation code. Port of getLegislation().
func (c *QCStatementCondition) Legislation() string {
	return c.legislation
}

// Check returns true if the condition is evaluated to true for the given certificate. Port of
// check(CertificateToken): a certificate carrying no QcStatements extension never matches, and
// each configured constraint that is present must be satisfied.
func (c *QCStatementCondition) Check(certificateToken *model.CertificateToken) bool {
	qcStatements := spi.QcStatementUtilsQcStatements(certificateToken)
	if qcStatements != nil {
		if utils.IsStringNotEmpty(c.oid) && !spi.QcStatementUtilsIsQcStatementPresent(qcStatements, c.oid) {
			return false
		}
		if utils.IsStringNotEmpty(c.qcType) && !spi.QcStatementUtilsIsQcTypePresent(qcStatements, c.qcType) {
			return false
		}
		if utils.IsStringNotEmpty(c.legislation) && !spi.QcStatementUtilsIsQcLegislationPresent(qcStatements, c.legislation) {
			return false
		}
		return true
	}
	return false
}

// ToString returns a human readable condition using the given indentation. Port of
// toString(String indent), whose four lines all carry the same indent.
//
// Go has no null string; the Java `if (indent == null) indent = ""` branch is unreachable here.
// Java renders a null oid/type/legislation as the literal "null" (StringBuilder#append(String)),
// which the empty-string stand-in cannot distinguish from an explicitly empty value - the two
// are equivalent everywhere Check() reads them, and this rendering is debug-only.
func (c *QCStatementCondition) ToString(indent string) string {
	var builder strings.Builder
	builder.WriteString(indent)
	builder.WriteString("QCStatementCondition: ")
	builder.WriteByte('\n')
	builder.WriteString(indent)
	builder.WriteString("oid: ")
	builder.WriteString(c.oid)
	builder.WriteByte('\n')
	builder.WriteString(indent)
	builder.WriteString("type: ")
	builder.WriteString(c.qcType)
	builder.WriteByte('\n')
	builder.WriteString(indent)
	builder.WriteString("legislation: ")
	builder.WriteString(c.legislation)
	builder.WriteByte('\n')
	return builder.String()
}

// String ports toString(), i.e. toString("").
func (c *QCStatementCondition) String() string {
	return c.ToString("")
}

// Equals ports equals(Object): same concrete type and equal oid/type/legislation.
func (c *QCStatementCondition) Equals(object any) bool {
	that, ok := object.(*QCStatementCondition)
	if !ok || that == nil {
		return false
	}
	if c == that {
		return true
	}
	return c.oid == that.oid && c.qcType == that.qcType && c.legislation == that.legislation
}

var _ tslmodel.Condition = (*QCStatementCondition)(nil)
