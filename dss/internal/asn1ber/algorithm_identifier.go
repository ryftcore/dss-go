package asn1ber

import (
	"bytes"
	"encoding/asn1"
	"errors"
	"fmt"
)

// AlgorithmIdentifier is the X.509 structure
//
//	AlgorithmIdentifier ::= SEQUENCE {
//	    algorithm   OBJECT IDENTIFIER,
//	    parameters  ANY DEFINED BY algorithm OPTIONAL }
//
// replacing org.bouncycastle.asn1.x509.AlgorithmIdentifier. Parameters holds the complete
// DER encoding of the parameters element, or nil when they are absent - the distinction
// matters, since an explicit NULL and an absent parameter are different encodings.
type AlgorithmIdentifier struct {
	// Algorithm is the algorithm OID.
	Algorithm asn1.ObjectIdentifier
	// Parameters is the DER encoding of the parameters, nil when absent.
	Parameters []byte
}

// NewAlgorithmIdentifier builds an AlgorithmIdentifier without parameters.
func NewAlgorithmIdentifier(algorithm asn1.ObjectIdentifier) *AlgorithmIdentifier {
	return &AlgorithmIdentifier{Algorithm: algorithm}
}

// NewAlgorithmIdentifierWithParameters builds an AlgorithmIdentifier carrying the given
// DER-encoded parameters.
func NewAlgorithmIdentifierWithParameters(algorithm asn1.ObjectIdentifier, parameters []byte) *AlgorithmIdentifier {
	return &AlgorithmIdentifier{Algorithm: algorithm, Parameters: parameters}
}

// ParseAlgorithmIdentifier decodes an AlgorithmIdentifier from its DER encoding.
// Port of AlgorithmIdentifier.getInstance(Object).
func ParseAlgorithmIdentifier(der []byte) (*AlgorithmIdentifier, error) {
	element, rest, err := Parse(der)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the AlgorithmIdentifier")
	}
	return AlgorithmIdentifierFromElement(element)
}

// AlgorithmIdentifierFromElement decodes an already parsed AlgorithmIdentifier SEQUENCE.
func AlgorithmIdentifierFromElement(element *Element) (*AlgorithmIdentifier, error) {
	if !element.IsUniversal(TagSequence) || !element.IsConstructed() {
		return nil, errors.New("AlgorithmIdentifier is not a SEQUENCE")
	}
	if len(element.Children()) == 0 || len(element.Children()) > 2 {
		return nil, fmt.Errorf("AlgorithmIdentifier: bad sequence size: %d", len(element.Children()))
	}
	oid, err := element.Children()[0].ObjectIdentifier()
	if err != nil {
		return nil, err
	}
	identifier := &AlgorithmIdentifier{Algorithm: oid}
	if len(element.Children()) == 2 {
		identifier.Parameters = element.Children()[1].DEREncoded()
	}
	return identifier, nil
}

// DER returns the DER encoding of the AlgorithmIdentifier.
func (a *AlgorithmIdentifier) DER() []byte {
	body := EncodeOID(a.Algorithm)
	if a.Parameters != nil {
		body = append(body, a.Parameters...)
	}
	return WriteSequence(body)
}

// Equals compares two AlgorithmIdentifiers by their DER encoding.
func (a *AlgorithmIdentifier) Equals(other *AlgorithmIdentifier) bool {
	if a == other {
		return true
	}
	if a == nil || other == nil {
		return false
	}
	return bytes.Equal(a.DER(), other.DER())
}
