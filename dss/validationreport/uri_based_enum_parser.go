// Ported from specs-validation-report/src/main/java/eu/europa/esig/validationreport/parsers/UriBasedEnumParser.java
// (DSS 6.5.RC1). The registry and parsing logic live in jaxb/jaxb_enums.go
// (they back the generated model's own MarshalText/UnmarshalText, so they
// must live where jaxb can reach them) and are re-exported here as thin
// forwarding functions - see doc.go's "Enum and parser re-exports".
package validationreport

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/validationreport/jaxb"
)

// ParseMainIndication parses the string value and returns the matching
// Indication, or the zero value if none matches.
func ParseMainIndication(v string) enumerations.Indication { return jaxb.ParseMainIndication(v) }

// ParseSubIndication parses the string value and returns the matching
// SubIndication, or the zero value if none matches.
func ParseSubIndication(v string) enumerations.SubIndication { return jaxb.ParseSubIndication(v) }

// ParseObjectType parses the string value and returns the matching
// ObjectType, or the zero value if none matches.
func ParseObjectType(v string) ObjectType { return jaxb.ParseObjectType(v) }

// ParseRevocationReason parses the string value and returns the matching
// RevocationReason, or the zero value if none matches.
func ParseRevocationReason(v string) enumerations.RevocationReason {
	return jaxb.ParseRevocationReason(v)
}

// ParseSignatureValidationProcessID parses the string value and returns the
// matching SignatureValidationProcessID, or the zero value if none matches.
func ParseSignatureValidationProcessID(v string) SignatureValidationProcessID {
	return jaxb.ParseSignatureValidationProcessID(v)
}

// ParseTypeOfProof parses the string value and returns the matching
// TypeOfProof, or the zero value if none matches.
func ParseTypeOfProof(v string) TypeOfProof { return jaxb.ParseTypeOfProof(v) }

// ParseConstraintStatus parses the string value and returns the matching
// ConstraintStatus, or the zero value if none matches.
func ParseConstraintStatus(v string) ConstraintStatus { return jaxb.ParseConstraintStatus(v) }

// Print returns v's URI, or "" if v is nil.
func Print(v enumerations.UriBasedEnum) string { return jaxb.Print(v) }
