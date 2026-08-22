// Ported from dss-enumerations/.../SemanticsIdentifier.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// SemanticsIdentifier lists ETSI EN 319 412-1 V1.1.1 semantics identifiers.
//
//	id-etsi-qcs-semantics-identifiers OBJECT IDENTIFIER ::= { itu-t(0)
//	identified-organization(4) etsi(0) id-cert-profile(194121) 1 }
//
// Implements OidDescription.
type SemanticsIdentifier string

const (
	// SemanticsIdentifierQcsSemanticsIdNatural is the semantics identifier for
	// natural person identifier.
	//
	//	id-etsi-qcs-semanticsId-Natural OBJECT IDENTIFIER ::= {
	//	id-etsi-qcs-semantics-identifiers 1 }
	SemanticsIdentifierQcsSemanticsIdNatural SemanticsIdentifier = "qcsSemanticsIdNatural"
	// SemanticsIdentifierQcsSemanticsIdLegal is the semantics identifier for legal
	// person identifier.
	//
	//	id-etsi-qcs-SemanticsId-Legal OBJECT IDENTIFIER ::= {
	//	id-etsi-qcs-semantics-identifiers 2 }
	SemanticsIdentifierQcsSemanticsIdLegal SemanticsIdentifier = "qcsSemanticsIdLegal"
	// SemanticsIdentifierQcsSemanticsIdEIDASNatural is the semantics identifier for
	// eIDAS natural person identifier.
	//
	//	id-etsi-qcs-semanticsId-eIDASNatural OBJECT IDENTIFIER ::= {
	//	id-etsi-qcs-semantics-identifiers 3 }
	SemanticsIdentifierQcsSemanticsIdEIDASNatural SemanticsIdentifier = "qcsSemanticsIdEIDASNatural"
	// SemanticsIdentifierQcsSemanticsIdEIDASLegal is the semantics identifier for
	// eIDAS legal person identifier.
	//
	//	id-etsi-qcs-semanticsId-eIDASNatural OBJECT IDENTIFIER ::= {
	//	id-etsi-qcs-semantics-identifiers 4 }
	SemanticsIdentifierQcsSemanticsIdEIDASLegal SemanticsIdentifier = "qcsSemanticsIdEIDASLegal"
)

type semanticsIdentifierFields struct {
	name        string
	oid         string
	description string
}

// semanticsIdentifierData holds the (name, oid, description) triple for each constant.
var semanticsIdentifierData = map[SemanticsIdentifier]semanticsIdentifierFields{
	SemanticsIdentifierQcsSemanticsIdNatural:      {"qcs-semanticsId-Natural", "0.4.0.194121.1.1", "Semantics identifier for natural person"},
	SemanticsIdentifierQcsSemanticsIdLegal:        {"qcs-SemanticsId-Legal", "0.4.0.194121.1.2", "Semantics identifier for legal person"},
	SemanticsIdentifierQcsSemanticsIdEIDASNatural: {"qcs-semanticsId-eIDASNatural", "0.4.0.194121.1.3", "Semantics identifier for eIDAS natural person"},
	SemanticsIdentifierQcsSemanticsIdEIDASLegal:   {"qcs-SemanticsId-eIDASLegal", "0.4.0.194121.1.4", "Semantics identifier for eIDAS legal person"},
}

// SemanticsIdentifierValues returns all SemanticsIdentifier constants in declaration order.
func SemanticsIdentifierValues() []SemanticsIdentifier {
	return []SemanticsIdentifier{
		SemanticsIdentifierQcsSemanticsIdNatural,
		SemanticsIdentifierQcsSemanticsIdLegal,
		SemanticsIdentifierQcsSemanticsIdEIDASNatural,
		SemanticsIdentifierQcsSemanticsIdEIDASLegal,
	}
}

// SemanticsIdentifierValueOf returns the SemanticsIdentifier matching the given Java enum name.
func SemanticsIdentifierValueOf(name string) (SemanticsIdentifier, error) {
	for _, v := range SemanticsIdentifierValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant SemanticsIdentifier.%s", name)
}

// Name returns the ETSI identifier name.
func (s SemanticsIdentifier) Name() string {
	return semanticsIdentifierData[s].name
}

// OID returns the OID of the semantics identifier. Implements OidDescription.
func (s SemanticsIdentifier) OID() string {
	return semanticsIdentifierData[s].oid
}

// Description returns the human-readable description. Implements OidDescription.
func (s SemanticsIdentifier) Description() string {
	return semanticsIdentifierData[s].description
}

// SemanticsIdentifierFromName returns the SemanticsIdentifier based on the provided
// identifier name, or "" if none matches. Upstream returns null here rather than
// throwing, so this is a plain zero-value return and not an error.
func SemanticsIdentifierFromName(name string) SemanticsIdentifier {
	for _, v := range SemanticsIdentifierValues() {
		if semanticsIdentifierData[v].name == name {
			return v
		}
	}
	return ""
}

// SemanticsIdentifierFromOID returns the SemanticsIdentifier based on the provided
// OID, or "" if none matches. Upstream returns null here rather than throwing, so this
// is a plain zero-value return and not an error.
func SemanticsIdentifierFromOID(oid string) SemanticsIdentifier {
	for _, v := range SemanticsIdentifierValues() {
		if semanticsIdentifierData[v].oid == oid {
			return v
		}
	}
	return ""
}
