package enumerations

import "testing"

type qcStatementCase struct {
	v           QCStatement
	description string
	oid         string
}

func qcStatementCases() []qcStatementCase {
	return []qcStatementCase{
		{QCStatement_QC_COMPLIANCE, "qc-compliance", "0.4.0.1862.1.1"},
		{QCStatement_QC_LIMIT_VALUE, "qc-limit-value", "0.4.0.1862.1.2"},
		{QCStatement_QC_RETENTION_PERIOD, "qc-retention-period", "0.4.0.1862.1.3"},
		{QCStatement_QC_SSCD, "qc-sscd", "0.4.0.1862.1.4"},
		{QCStatement_QC_PDS, "qc-pds", "0.4.0.1862.1.5"},
		{QCStatement_QC_TYPE, "qc-type", "0.4.0.1862.1.6"},
		{QCStatement_QC_CCLEGISLATION, "qc-cclegislation", "0.4.0.1862.1.7"},
		{QCStatement_QC_IDENT_METHOD, "qc-identMethod", "0.4.0.1862.1.8"},
		{QCStatement_QC_QSCD_LEGISLATION, "qc-qscdLegislation", "0.4.0.1862.1.9"},
		{QCStatement_QC_PSB, "qc-psb", "0.4.0.194126.1.3"},
	}
}

func TestQCStatementFields(t *testing.T) {
	cases := qcStatementCases()
	if len(QCStatementValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(QCStatementValues()))
	}
	for _, c := range cases {
		if got := c.v.Description(); got != c.description {
			t.Errorf("%v.Description() = %q, want %q", c.v, got, c.description)
		}
		if got := c.v.OID(); got != c.oid {
			t.Errorf("%v.OID() = %q, want %q", c.v, got, c.oid)
		}
	}
}

func TestQCStatementForLabel(t *testing.T) {
	for _, c := range qcStatementCases() {
		if got := QCStatementForLabel(c.description); got != c.v {
			t.Errorf("QCStatementForLabel(%q) = %v, want %v", c.description, got, c.v)
		}
	}
	if got := QCStatementForLabel("nope"); got != "" {
		t.Errorf("QCStatementForLabel(unknown) = %v, want zero value", got)
	}
}

func TestQCStatementForOID(t *testing.T) {
	for _, c := range qcStatementCases() {
		if got := QCStatementForOID(c.oid); got != c.v {
			t.Errorf("QCStatementForOID(%q) = %v, want %v", c.oid, got, c.v)
		}
	}
	if got := QCStatementForOID("9.9.9"); got != "" {
		t.Errorf("QCStatementForOID(unknown) = %v, want zero value", got)
	}
}
