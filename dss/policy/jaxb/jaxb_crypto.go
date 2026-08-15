// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule and the LevelConstraint attribute-order finding
// this file also follows.
package jaxb

// Algo is the Go form of the generated JAXB class Algo (complexType Algo): a
// simpleContent extension of xs:string carrying the algorithm name as its
// text value, with optional Size and Date attributes.
type Algo struct {
	Value string  `xml:",chardata"`
	Size  *int    `xml:"Size,attr,omitempty"`
	Date  *string `xml:"Date,attr,omitempty"`
}

// ListAlgo is the Go form of the generated JAXB class ListAlgo (complexType
// ListAlgo).
type ListAlgo struct {
	Algos []*Algo `xml:"Algo,omitempty"`
	LevelConstraint
}

// AlgoExpirationDate is the Go form of the generated JAXB class
// AlgoExpirationDate (complexType AlgoExpirationDate).
type AlgoExpirationDate struct {
	Format           *string    `xml:"Format,attr,omitempty"`
	UpdateDate       *string    `xml:"UpdateDate,attr,omitempty"`
	LevelAfterUpdate LevelValue `xml:"LevelAfterUpdate,attr,omitempty"`
	ListAlgo
}

// CryptographicConstraint is the Go form of the generated JAXB class
// CryptographicConstraint (complexType CryptographicConstraint): global
// constraints about cryptographic usage (encryption, digest, key length,
// algorithm deprecation,...).
type CryptographicConstraint struct {
	AcceptableEncryptionAlgo *ListAlgo           `xml:"AcceptableEncryptionAlgo,omitempty"`
	MiniPublicKeySize        *ListAlgo           `xml:"MiniPublicKeySize,omitempty"`
	AcceptableDigestAlgo     *ListAlgo           `xml:"AcceptableDigestAlgo,omitempty"`
	AlgoExpirationDate       *AlgoExpirationDate `xml:"AlgoExpirationDate,omitempty"`
	LevelConstraint
}
