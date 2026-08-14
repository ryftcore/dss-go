// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. Per PORTING.md's generated-JAXB rule the
// generated classes are grouped into schema-area files rather than kept
// one-file-per-class; every Java class keeps its name and exact field order so
// encoding/xml reproduces the JAXB RI's element sequence and attribute order
// byte for byte.
//
// # Attribute order across the LevelConstraint hierarchy
//
// Verified against a live JAXB RI oracle (unmarshal-then-remarshal of
// upstream's src/main/resources/policy/*.xml - see doc.go): a subclass's own
// @XmlAttribute fields are written before the Level attribute inherited from
// LevelConstraint (e.g. <RevocationFreshness Unit="DAYS" Value="0"
// Level="IGNORE"/>, <AlgoExpirationDate Format="..." UpdateDate="..."
// LevelAfterUpdate="..." Level="FAIL">). This mirrors the "extension
// attributes before base attributes" rule already documented in
// dss/diagnostic/jaxb/jaxb_token.go for the JAXB RI. Every struct below that
// extends LevelConstraint therefore embeds it as its LAST field.
package jaxb

// LevelConstraint is the Go form of the generated JAXB class LevelConstraint
// (complexType LevelConstraint): a boolean check that follows the specified
// level behavior in case of failure.
type LevelConstraint struct {
	Level LevelValue `xml:"Level,attr,omitempty"`
}

// ValueConstraint is the Go form of the generated JAXB class ValueConstraint
// (complexType ValueConstraint): a value check that follows the specified
// level behavior if the checked element is not equal to the specified value.
type ValueConstraint struct {
	Value *string `xml:"value,attr,omitempty"`
	LevelConstraint
}

// IntValueConstraint is the Go form of the generated JAXB class
// IntValueConstraint (complexType IntValueConstraint): an integer value check
// that follows the specified level behavior when the checked element is
// compliant with the defined constraint value.
type IntValueConstraint struct {
	Value *int `xml:"value,attr,omitempty"`
	LevelConstraint
}

// MultiValuesConstraint is the Go form of the generated JAXB class
// MultiValuesConstraint (complexType MultiValuesConstraint): a multi-values
// check that follows the specified level behavior if the checked element is
// not present in the list. '*' can be used and means any value.
type MultiValuesConstraint struct {
	Id []string `xml:"Id,omitempty"`
	LevelConstraint
}

// CertificateValuesConstraint is the Go form of the generated JAXB class
// CertificateValuesConstraint (complexType CertificateValuesConstraint).
type CertificateValuesConstraint struct {
	CertificateExtensions *MultiValuesConstraint `xml:"CertificateExtensions,omitempty"`
	CertificatePolicies   *MultiValuesConstraint `xml:"CertificatePolicies,omitempty"`
	LevelConstraint
}

// TimeConstraint is the Go form of the generated JAXB class TimeConstraint
// (complexType TimeConstraint): a time-based check that follows the specified
// level behavior if the checked element is over the time limit.
type TimeConstraint struct {
	Unit  TimeUnit `xml:"Unit,attr,omitempty"`
	Value *int     `xml:"Value,attr,omitempty"`
	LevelConstraint
}
