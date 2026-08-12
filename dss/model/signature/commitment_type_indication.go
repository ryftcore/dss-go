// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/CommitmentTypeIndication.java (DSS 6.5.RC1).
package signature

// CommitmentTypeIndication represents the commitment type indication identifiers extracted
// from the signature.
type CommitmentTypeIndication struct {
	// identifier is the URI or OID identifier.
	identifier string

	// description is the description message.
	description string

	// documentReferences is the list of document references.
	documentReferences []string

	// objectReferences is the list of signed data referenced by the current commitment
	// (XAdES).
	objectReferences []string

	// allDataSignedObjects defines if AllSignedDataObjects element is present (XAdES).
	allDataSignedObjects bool
}

// NewCommitmentTypeIndication is the default constructor.
func NewCommitmentTypeIndication(identifier string) *CommitmentTypeIndication {
	return &CommitmentTypeIndication{identifier: identifier}
}

// Identifier gets the identifier. Port of getIdentifier().
func (c *CommitmentTypeIndication) Identifier() string {
	return c.identifier
}

// Description gets the description. Port of getDescription().
func (c *CommitmentTypeIndication) Description() string {
	return c.description
}

// SetDescription sets the description. Port of setDescription(String).
func (c *CommitmentTypeIndication) SetDescription(description string) {
	c.description = description
}

// DocumentReferences gets the document references. Port of getDocumentReferences().
func (c *CommitmentTypeIndication) DocumentReferences() []string {
	return c.documentReferences
}

// SetDocumentReferences sets the document references. Port of
// setDocumentReferences(List<String>).
func (c *CommitmentTypeIndication) SetDocumentReferences(documentReferences []string) {
	c.documentReferences = documentReferences
}

// ObjectReferences gets a list of signed data objects referenced by the current
// CommitmentType. Port of getObjectReferences().
func (c *CommitmentTypeIndication) ObjectReferences() []string {
	return c.objectReferences
}

// SetObjectReferences sets a list of signed data objects referenced by the current
// CommitmentType. Port of setObjectReferences(List<String>).
func (c *CommitmentTypeIndication) SetObjectReferences(objectReferences []string) {
	c.objectReferences = objectReferences
}

// IsAllDataSignedObjects gets if AllDataSignedObjects are referenced by the current
// CommitmentType (XAdES only). Port of isAllDataSignedObjects().
func (c *CommitmentTypeIndication) IsAllDataSignedObjects() bool {
	return c.allDataSignedObjects
}

// SetAllDataSignedObjects sets if AllDataSignedObjects are referenced by the current
// CommitmentType (XAdES only). Port of setAllDataSignedObjects(boolean).
func (c *CommitmentTypeIndication) SetAllDataSignedObjects(allDataSignedObjects bool) {
	c.allDataSignedObjects = allDataSignedObjects
}
