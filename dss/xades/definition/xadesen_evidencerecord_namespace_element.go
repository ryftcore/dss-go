// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xadesen/XAdESEvidencerecordNamespaceElement.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdESEvidencerecordNamespaceElement defines elements specified within
// "http://uri.etsi.org/19132/v1.1.1#" XAdES Evidence Record container namespace (ETSI TS 119 132-3).
type XAdESEvidencerecordNamespaceElement string

// XAdESEvidencerecordNamespaceElement constants, one per element name in the
// XAdES Evidence Record container namespace.
const (
	XAdESEvidencerecordNamespaceElementASN1EvidenceRecord     XAdESEvidencerecordNamespaceElement = "ASN1_EVIDENCE_RECORD"
	XAdESEvidencerecordNamespaceElementEvidenceRecord         XAdESEvidencerecordNamespaceElement = "EVIDENCE_RECORD"
	XAdESEvidencerecordNamespaceElementSealingEvidenceRecords XAdESEvidencerecordNamespaceElement = "SEALING_EVIDENCE_RECORDS"
)

// xadesEvidencerecordNamespaceElementTagNames maps each constant to its wire tag name (getTagName()).
var xadesEvidencerecordNamespaceElementTagNames = map[XAdESEvidencerecordNamespaceElement]string{
	XAdESEvidencerecordNamespaceElementASN1EvidenceRecord:     "ASN1EvidenceRecord",
	XAdESEvidencerecordNamespaceElementEvidenceRecord:         "EvidenceRecord",
	XAdESEvidencerecordNamespaceElementSealingEvidenceRecords: "SealingEvidenceRecords",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e XAdESEvidencerecordNamespaceElement) TagName() string {
	return xadesEvidencerecordNamespaceElementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e XAdESEvidencerecordNamespaceElement) Namespace() *common.DSSNamespace {
	return XAdESNamespaceXAdESEvidencerecordNamespace
}

// URI implements common.DSSElement. Ports getURI().
func (e XAdESEvidencerecordNamespaceElement) URI() string {
	return XAdESNamespaceXAdESEvidencerecordNamespace.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e XAdESEvidencerecordNamespaceElement) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ common.DSSElement = XAdESEvidencerecordNamespaceElement("")
