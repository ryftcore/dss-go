// Package definition ports the dss-xades definition subpackage
// (eu.europa.esig.dss.xades.definition), the XML element/attribute name and
// XPath vocabulary for every XML Signature and XAdES schema version (1.1.1
// through 1.4.2) DSS understands, plus a Trusted List element vocabulary
// used when producing/verifying trusted-list signatures.
//
// This is a large, mechanically generated-looking vocabulary: it is
// consulted by name from xades' builders and validators rather than used
// directly by most callers. The main entry types are the DSSElement/
// DSSAttribute/DSSPath-style interfaces (XAdES111Element, XAdES111Attribute,
// XAdES111Path, and their versioned siblings up to XAdES 1.4.2) and
// TrustedListElement.
package definition
