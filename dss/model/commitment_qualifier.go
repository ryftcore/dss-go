// Ported from dss-model/.../CommitmentQualifier.java (DSS 6.5.RC1).
package model

// CommitmentQualifier defines a CommitmentTypeQualifier to be incorporated
// within a signature.
type CommitmentQualifier struct {
	// Oid defines the unique commitment qualifier identifier (CAdES/PAdES
	// only). Use: CONDITIONAL (required for CAdES/PAdES).
	oid string

	// content defines the content of the qualifier (required). The content
	// of a qualifier may be of any type, but developers may need to ensure
	// it corresponds to the used signature format (i.e. XML for XAdES,
	// ASN.1 for CAdES, etc.). Use: REQUIRED.
	content DSSDocument
}

// NewCommitmentQualifier instantiates the object with null values. Ports
// the default constructor.
func NewCommitmentQualifier() *CommitmentQualifier {
	return &CommitmentQualifier{}
}

// Oid gets the unique object identifier of the Commitment Qualifier.
func (c *CommitmentQualifier) Oid() string { return c.oid }

// SetOid sets the unique object identifier of the Commitment Qualifier
// (CAdES/PAdES only!). Use: CONDITIONAL (required for CAdES/PAdES).
func (c *CommitmentQualifier) SetOid(oid string) { c.oid = oid }

// Content gets the content of the Commitment Qualifier.
func (c *CommitmentQualifier) Content() DSSDocument { return c.content }

// SetContent sets the content of the Commitment Qualifier. Use: REQUIRED.
func (c *CommitmentQualifier) SetContent(content DSSDocument) { c.content = content }
