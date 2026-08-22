// Ported from dss-model/.../CommonObjectIdentifier.java (DSS 6.5.RC1).
package model

import "github.com/ryftcore/dss-go/dss/enumerations"

// CommonObjectIdentifier provides a basic implementation of
// enumerations.ObjectIdentifier, allowing creation of a customized
// ObjectIdentifierType signed property.
type CommonObjectIdentifier struct {
	// uri defines the URI of the ObjectIdentifierType (used in XAdES,
	// JAdES). Use: CONDITIONAL (should be present in XAdES, JAdES).
	uri string

	// oid defines the OID identifier of the ObjectIdentifierType (used in
	// CAdES, PAdES, JAdES). Use: CONDITIONAL (shall be present in CAdES,
	// PAdES. May be present in XAdES, JAdES).
	oid string

	// qualifier defines the type of URI when an OID is provided as a URI
	// (used in XAdES, see ETSI EN 319 132-1 for more details). Use:
	// CONDITIONAL (may be present in XAdES).
	qualifier enumerations.ObjectIdentifierQualifier

	// description contains a text describing the ObjectIdentifierType.
	// Use: OPTIONAL.
	description string

	// documentationReferences contains arbitrary references pointing to
	// further explanatory document about the ObjectIdentifierType. Use:
	// OPTIONAL.
	documentationReferences []string
}

var _ enumerations.ObjectIdentifier = (*CommonObjectIdentifier)(nil)

// NewCommonObjectIdentifier instantiates the object with null values.
// Ports the default constructor.
func NewCommonObjectIdentifier() *CommonObjectIdentifier {
	return &CommonObjectIdentifier{}
}

// URI implements enumerations.UriBasedEnum / ObjectIdentifier.
func (c *CommonObjectIdentifier) URI() string { return c.uri }

// SetUri sets the URI identifying the ObjectIdentifierType. Use:
// CONDITIONAL (should be present in XAdES, JAdES).
func (c *CommonObjectIdentifier) SetUri(uri string) { c.uri = uri }

// OID implements enumerations.OidBasedEnum / ObjectIdentifier.
func (c *CommonObjectIdentifier) OID() string { return c.oid }

// SetOid sets the OID identifying the ObjectIdentifierType. Use:
// CONDITIONAL (shall be present in CAdES, PAdES. May be present in XAdES,
// JAdES). Note: when using OID in XAdES, a Qualifier shall be defined via
// SetQualifier. See EN 319 132-1 "5.1.2 The ObjectIdentifierType data type"
// for more information.
func (c *CommonObjectIdentifier) SetOid(oid string) { c.oid = oid }

// Qualifier implements ObjectIdentifier.
func (c *CommonObjectIdentifier) Qualifier() enumerations.ObjectIdentifierQualifier {
	return c.qualifier
}

// SetQualifier sets the Qualifier defining the type of OID identifier used
// for ObjectIdentifierType. See EN 319 132-1 "5.1.2 The ObjectIdentifierType
// data type" for more information. Use: CONDITIONAL (shall be present in
// XAdES when using OID identifier, but not URI). Note: used only in XAdES.
func (c *CommonObjectIdentifier) SetQualifier(qualifier enumerations.ObjectIdentifierQualifier) {
	c.qualifier = qualifier
}

// Description implements enumerations.OidDescription / ObjectIdentifier.
func (c *CommonObjectIdentifier) Description() string { return c.description }

// SetDescription sets text describing the ObjectIdentifierType object.
// Use: OPTIONAL.
func (c *CommonObjectIdentifier) SetDescription(description string) { c.description = description }

// DocumentationReferences implements ObjectIdentifier.
func (c *CommonObjectIdentifier) DocumentationReferences() []string {
	return c.documentationReferences
}

// SetDocumentationReferences sets references pointing to documentation
// describing the ObjectIdentifierType. Use: OPTIONAL.
func (c *CommonObjectIdentifier) SetDocumentationReferences(documentationReferences ...string) {
	c.documentationReferences = documentationReferences
}
