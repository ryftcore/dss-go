// Ported from dss-model/.../SpDocSpecification.java (DSS 6.5.RC1).
package model

import (
	"fmt"
	"reflect"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// SpDocSpecification represents an "SpDocSpecification" element.
type SpDocSpecification struct {
	// id is the OID, e.g. OID : 2.2.25.1.
	id string

	// description is optional.
	description string

	// documentationReferences are optional document references.
	documentationReferences []string

	// qualifier is specified in EN 319 132.
	qualifier enumerations.ObjectIdentifierQualifier
}

// NewSpDocSpecification instantiates the object with null values. Ports
// the default constructor.
func NewSpDocSpecification() *SpDocSpecification {
	return &SpDocSpecification{}
}

// Id gets the identifier.
func (s *SpDocSpecification) Id() string { return s.id }

// SetId sets the identifier (URI or OID), e.g. 2.2.25.1 for OID.
func (s *SpDocSpecification) SetId(id string) { s.id = id }

// Description gets the description.
func (s *SpDocSpecification) Description() string { return s.description }

// SetDescription sets the description.
func (s *SpDocSpecification) SetDescription(description string) { s.description = description }

// DocumentationReferences gets the documentation references.
func (s *SpDocSpecification) DocumentationReferences() []string { return s.documentationReferences }

// SetDocumentationReferences sets the documentation references.
func (s *SpDocSpecification) SetDocumentationReferences(documentationReferences ...string) {
	s.documentationReferences = documentationReferences
}

// Qualifier gets the qualifier (used in XAdES).
func (s *SpDocSpecification) Qualifier() enumerations.ObjectIdentifierQualifier { return s.qualifier }

// SetQualifier sets the qualifier (used in XAdES).
func (s *SpDocSpecification) SetQualifier(qualifier enumerations.ObjectIdentifierQualifier) {
	s.qualifier = qualifier
}

// Equals ports SpDocSpecification#equals.
func (s *SpDocSpecification) Equals(other *SpDocSpecification) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	return s.id == other.id &&
		s.description == other.description &&
		reflect.DeepEqual(s.documentationReferences, other.documentationReferences) &&
		s.qualifier == other.qualifier
}

// String ports SpDocSpecification#toString.
func (s *SpDocSpecification) String() string {
	return fmt.Sprintf("SpDocSpecification {id='%s', description='%s', documentationReferences=%v, qualifier=%v}",
		s.id, s.description, s.documentationReferences, s.qualifier)
}
